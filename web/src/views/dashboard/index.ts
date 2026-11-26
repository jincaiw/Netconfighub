import type { BackupTask } from '@/types'

export interface DailyPoint {
  date: string
  total: number
  success: number
  failed: number
}

export interface RankItem {
  deviceId: number
  deviceName: string
  count: number
}

export interface DashboardStats {
  totalDevices: number
  successRate: number
  failedDevices: number
  unreadAlerts: number
}

export interface RecentTaskRow {
  id: number
  device_name: string
  status: BackupTask['status']
  finished_at: string
}

export interface FailedDeviceRow {
  id: number
  name: string
  ip: string
  reason: string
}

function dayKey(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function lastDayKeys(days: number): string[] {
  const keys: string[] = []
  const today = new Date()
  for (let i = days - 1; i >= 0; i -= 1) {
    const day = new Date(today.getFullYear(), today.getMonth(), today.getDate() - i)
    keys.push(dayKey(day))
  }
  return keys
}

/** 以 BackupTask.created_at（必填）按天聚合近 N 天的成功/失败任务数。 */
export function aggregateByDay(tasks: BackupTask[], days = 7): DailyPoint[] {
  const keys = lastDayKeys(days)
  const buckets = new Map<string, DailyPoint>()
  keys.forEach((date) => buckets.set(date, { date, total: 0, success: 0, failed: 0 }))

  tasks.forEach((task) => {
    if (!task.created_at || typeof task.created_at !== 'string') return
    const key = task.created_at.slice(0, 10)
    const bucket = buckets.get(key)
    if (!bucket) return
    bucket.total += 1
    if (task.status === 'success') bucket.success += 1
    else if (task.status === 'failed') bucket.failed += 1
  })

  return keys.map((key) => buckets.get(key)!)
}

/** 按 device_id 统计失败次数，降序取前 N。 */
export function rankByFailed(tasks: BackupTask[], limit = 5): RankItem[] {
  const counter = new Map<number, RankItem>()
  tasks.forEach((task) => {
    if (task.status !== 'failed') return
    const existing = counter.get(task.device_id)
    if (existing) {
      existing.count += 1
    } else {
      counter.set(task.device_id, {
        deviceId: task.device_id,
        deviceName: task.device_name || `设备#${task.device_id}`,
        count: 1,
      })
    }
  })
  return [...counter.values()].sort((a, b) => b.count - a.count).slice(0, limit)
}

/** 按 device_id 统计备份任务总量，降序取前 N。 */
export function rankByCount(tasks: BackupTask[], limit = 5): RankItem[] {
  const counter = new Map<number, RankItem>()
  tasks.forEach((task) => {
    const existing = counter.get(task.device_id)
    if (existing) {
      existing.count += 1
    } else {
      counter.set(task.device_id, {
        deviceId: task.device_id,
        deviceName: task.device_name || `设备#${task.device_id}`,
        count: 1,
      })
    }
  })
  return [...counter.values()].sort((a, b) => b.count - a.count).slice(0, limit)
}

/** 四项指标口径与迁移前 Dashboard 保持一致。 */
export function buildStats(params: {
  totalDevices: number
  tasks: BackupTask[]
  failedDevices: number
  unreadAlerts: number
}): DashboardStats {
  const total = params.tasks.length
  const success = params.tasks.filter((task) => task.status === 'success').length
  return {
    totalDevices: params.totalDevices,
    successRate: total > 0 ? Math.round((success / total) * 1000) / 10 : 0,
    failedDevices: params.failedDevices,
    unreadAlerts: params.unreadAlerts,
  }
}

/** 可视化元素的数据契约：全部由真实任务纯计算得出，不请求、不伪造。 */
export interface StatusSlice {
  name: string
  value: number
}

export interface TrendDelta {
  direction: 'up' | 'down' | 'flat'
  percent: number
}

export interface InnerCardData {
  backupCount: number
  successRate: number
  countDelta: TrendDelta | null
  rateDelta: TrendDelta | null
}

function parseDay(task: BackupTask): number {
  const t = Date.parse(task.created_at)
  return Number.isNaN(t) ? 0 : t
}

function withinDays(task: BackupTask, days: number): boolean {
  const t = parseDay(task)
  if (!t) return false
  return t >= Date.now() - days * 24 * 60 * 60 * 1000
}

/** TopPanel 迷你图序列：近 7 天的成功/失败逐日数量，复用 aggregateByDay。 */
export function buildSparkline(tasks: BackupTask[]): { success: number[]; failed: number[] } {
  const points = aggregateByDay(tasks, 7)
  return {
    success: points.map((p) => p.success),
    failed: points.map((p) => p.failed),
  }
}

/** MiddleChart 环形饼图：按状态归一化分布（成功/失败/进行中）。 */
export function statusDistribution(tasks: BackupTask[]): StatusSlice[] {
  if (tasks.length === 0) return []
  let success = 0
  let failed = 0
  let ongoing = 0
  tasks.forEach((t) => {
    if (t.status === 'success') success += 1
    else if (t.status === 'failed' || t.status === 'timeout') failed += 1
    else ongoing += 1
  })
  const slices: StatusSlice[] = []
  if (success) slices.push({ name: '成功', value: success })
  if (failed) slices.push({ name: '失败', value: failed })
  if (ongoing) slices.push({ name: '进行中', value: ongoing })
  return slices
}

/** RankList 周/月切换：按 created_at 过滤周期后重新排名。 */
export function rankByPeriod(
  tasks: BackupTask[],
  period: 'week' | 'month',
  metric: 'failed' | 'count',
): RankItem[] {
  const days = period === 'week' ? 7 : 30
  const scoped = tasks.filter((t) => withinDays(t, days))
  return metric === 'failed' ? rankByFailed(scoped, 5) : rankByCount(scoped, 5)
}

/** 由时序前后半段对比得出涨跌方向/百分比；样本不足或除零时返回 null（显示中性）。 */
export function trendDelta(values: number[]): TrendDelta | null {
  if (values.length < 4) return null
  const half = Math.floor(values.length / 2)
  const before = values.slice(0, half).reduce((a, b) => a + b, 0)
  const after = values.slice(half).reduce((a, b) => a + b, 0)
  if (before === 0) return null
  const percent = Math.round(Math.abs(((after - before) / before) * 100) * 10) / 10
  return {
    direction: after > before ? 'up' : after < before ? 'down' : 'flat',
    percent,
  }
}

/** OutputOverview inner-card：本月备份次数与成功率 + 可计算的环比。 */
export function buildInnerCards(tasks: BackupTask[]): InnerCardData {
  const month = tasks.filter((t) => withinDays(t, 30))
  const halfMonth = Math.floor(30 / 2)
  const recent = month.filter((t) => withinDays(t, halfMonth))
  const prior = month.filter((t) => !withinDays(t, halfMonth))
  const success = month.filter((t) => t.status === 'success').length
  const rate = month.length > 0 ? Math.round((success / month.length) * 1000) / 10 : 0
  const countDelta = prior.length > 0 ? trendDelta([prior.length, recent.length]) : null
  const priorRate = prior.length > 0 ? (prior.filter((t) => t.status === 'success').length / prior.length) * 100 : 0
  const recentRate = recent.length > 0 ? (recent.filter((t) => t.status === 'success').length / recent.length) * 100 : 0
  const rateDelta = prior.length > 0 && recent.length > 0 ? trendDelta([priorRate, recentRate]) : null
  return { backupCount: month.length, successRate: rate, countDelta, rateDelta }
}

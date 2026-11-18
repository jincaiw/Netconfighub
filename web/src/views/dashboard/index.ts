import type { BackupTask } from '@/types'

export interface DailyPoint {
  date: string
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
  keys.forEach((date) => buckets.set(date, { date, success: 0, failed: 0 }))

  tasks.forEach((task) => {
    if (!task.created_at || typeof task.created_at !== 'string') return
    const key = task.created_at.slice(0, 10)
    const bucket = buckets.get(key)
    if (!bucket) return
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

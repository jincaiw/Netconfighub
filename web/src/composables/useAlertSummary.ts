import { ref } from 'vue'
import { getAlerts, getUnreadAlertCount, markAlertAsRead, markAllAlertsAsRead } from '@/api/alert'
import type { Alert } from '@/types'

// Module-level shared state: AppHeader and Dashboard consume the same refs and
// the same in-flight request, so a first render issues only one count request.
const unreadCount = ref(0)
const alerts = ref<Alert[]>([])
// 只读失败态：失败时置位并保留上次数据，供 UI 展示“重试”入口，不清零角标/列表。
const pollError = ref(false)
const listError = ref(false)
let inflight: Promise<void> | null = null
let timer: ReturnType<typeof setInterval> | null = null
let consumers = 0

export function useAlertSummary() {
  function fetchUnreadCount(): Promise<void> {
    if (inflight) return inflight
    inflight = (async () => {
      try {
        const res = await getUnreadAlertCount()
        if (res.data) unreadCount.value = res.data.count
        pollError.value = false
      } catch {
        // keep the last known count; polling will retry
        pollError.value = true
      }
    })().finally(() => {
      inflight = null
    })
    return inflight
  }

  async function fetchAlerts(): Promise<void> {
    try {
      const res = await getAlerts({ page: 1, page_size: 20 })
      if (res.data) alerts.value = res.data.items
      listError.value = false
    } catch {
      // keep the last known list
      listError.value = true
    }
  }

  /** 手动重试：重新拉取计数与列表，失败态由各自请求结果复位。 */
  async function refresh(): Promise<void> {
    await Promise.all([fetchUnreadCount(), fetchAlerts()])
  }

  async function markRead(alert: Alert): Promise<boolean> {
    if (alert.is_read) return true
    try {
      await markAlertAsRead(alert.id)
      alert.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
      return true
    } catch {
      return false
    }
  }

  async function markAllRead(): Promise<boolean> {
    try {
      await markAllAlertsAsRead()
      alerts.value.forEach((a) => { a.is_read = true })
      unreadCount.value = 0
      return true
    } catch {
      return false
    }
  }

  function start(): void {
    consumers++
    void fetchUnreadCount()
    if (!timer) {
      timer = setInterval(() => { void fetchUnreadCount() }, 30000)
    }
  }

  function stop(): void {
    if (consumers > 0) consumers--
    if (consumers === 0 && timer) {
      clearInterval(timer)
      timer = null
    }
  }

  return {
    unreadCount,
    alerts,
    pollError,
    listError,
    fetchUnreadCount,
    fetchAlerts,
    refresh,
    markRead,
    markAllRead,
    start,
    stop,
  }
}

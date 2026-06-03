import { get, put, del } from './index'
import type { Alert, PaginatedResponse } from '@/types'

export function getAlerts(params?: Record<string, any>) {
  return get<PaginatedResponse<Alert>>('/alerts', params)
}

export function getUnreadAlertCount() {
  return get<{ count: number }>('/alerts/unread-count')
}

export function markAlertAsRead(id: number) {
  return put<null>(`/alerts/${id}/read`)
}

export function markAllAlertsAsRead() {
  return put<null>('/alerts/read-all')
}

export function deleteAlert(id: number) {
  return del<null>(`/alerts/${id}`)
}

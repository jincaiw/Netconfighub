import { get } from './index'
import type { AuditLog, PaginatedResponse } from '@/types'

export function getAuditLogs(params: { page: number; page_size: number }) {
  return get<PaginatedResponse<AuditLog>>('/audit-logs', params)
}

import { get, post } from './index'
import type { BackupTask, PaginatedResponse } from '@/types'

export function triggerDeviceBackup(deviceId: number) {
  return post<BackupTask>(`/devices/${deviceId}/backup`)
}

export function triggerGroupBackup(groupId: number) {
  return post<BackupTask[]>(`/groups/${groupId}/backup`)
}

export function getBackupTasks(params?: Record<string, any>) {
  return get<PaginatedResponse<BackupTask>>('/backups', params)
}

/** 完整区间的备份历史（用于仪表盘按范围聚合，避免仅取最新 100 条样本）。 */
export function getBackupHistory(params?: { days?: number; limit?: number }) {
  return get<BackupTask[]>('/backups/history', params)
}

export function getBackupDetail(id: number) {
  return get<BackupTask>(`/backups/${id}`)
}

export function getBackupConfig(id: number) {
  return get<{ content: string }>(`/backups/${id}/config`)
}

export function getFailedDevices(params?: Record<string, any>) {
  return get<PaginatedResponse<import('@/types').Device>>('/failed-devices', params)
}

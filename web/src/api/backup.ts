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

export function getBackupDetail(id: number) {
  return get<BackupTask>(`/backups/${id}`)
}

export function getBackupConfig(id: number) {
  return get<{ content: string }>(`/backups/${id}/config`)
}

export function getFailedDevices(params?: Record<string, any>) {
  return get<PaginatedResponse<import('@/types').Device>>('/failed-devices', params)
}

import { get, put } from './index'
import type { ConfigVersion, PaginatedResponse, Deviation, DiffResult } from '@/types'

export function getDeviceConfig(deviceId: number) {
  return get<{ content: string }>(`/devices/${deviceId}/config`)
}

export function getDeviceVersions(deviceId: number) {
  return get<ConfigVersion[]>(`/devices/${deviceId}/versions`)
}

export function getDeviceDeviations(deviceId: number, params?: Record<string, any>) {
  return get<PaginatedResponse<Deviation>>(`/devices/${deviceId}/deviations`, params)
}

export function getConfigDiff(deviceId: number, from: string, to: string) {
  return get<DiffResult>(`/devices/${deviceId}/diff`, { from, to })
}

export function getConfigs() {
  return get<Array<{ key: string; value: string }>>('/configs')
}

export function setConfig(key: string, value: string) {
  return put<null>('/configs', { key, value })
}

import { get, post, put, del } from './index'
import request from './index'
import type { Device, DeviceCreateRequest, Group, GroupCreateRequest, ImportResult, PaginatedResponse, ApiResponse } from '@/types'

export function getDevices(params?: Record<string, any>) {
  return get<PaginatedResponse<Device>>('/devices', params)
}

export function getDevice(id: number) {
  return get<Device>(`/devices/${id}`)
}

export function createDevice(data: DeviceCreateRequest) {
  return post<Device>('/devices', data as unknown as Record<string, any>)
}

export function updateDevice(id: number, data: Partial<DeviceCreateRequest>) {
  return put<Device>(`/devices/${id}`, data as unknown as Record<string, any>)
}

export function deleteDevice(id: number) {
  return del<null>(`/devices/${id}`)
}

export function enableDevice(id: number) {
  return put<Device>(`/devices/${id}/enable`)
}

export function disableDevice(id: number) {
  return put<Device>(`/devices/${id}/disable`)
}

export function importDevices(file: File): Promise<ApiResponse<ImportResult>> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/devices/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }).then((res) => res.data)
}

export function exportDevices(params?: Record<string, any>): Promise<Blob> {
  return request.get('/devices/export', { params, responseType: 'blob' }).then((res) => res.data)
}

export function getGroups(params?: Record<string, any>) {
  return get<PaginatedResponse<Group>>('/groups', params)
}

export function getGroup(id: number) {
  return get<Group>(`/groups/${id}`)
}

export function createGroup(data: GroupCreateRequest) {
  return post<Group>('/groups', data as unknown as Record<string, any>)
}

export function updateGroup(id: number, data: Partial<GroupCreateRequest>) {
  return put<Group>(`/groups/${id}`, data as unknown as Record<string, any>)
}

export function deleteGroup(id: number) {
  return del<null>(`/groups/${id}`)
}

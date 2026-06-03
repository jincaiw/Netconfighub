import { get, post, put, del } from './index'
import request from './index'
import type { Baseline, Deviation, PaginatedResponse } from '@/types'

export function getBaselines(params?: Record<string, any>) {
  return get<PaginatedResponse<Baseline>>('/baselines', params)
}

export function getBaseline(id: number) {
  return get<Baseline>(`/baselines/${id}`)
}

export function createBaseline(data: Partial<Baseline>) {
  return post<Baseline>('/baselines', data as unknown as Record<string, any>)
}

export function updateBaseline(id: number, data: { content: string }) {
  return put<Baseline>(`/baselines/${id}`, data as unknown as Record<string, any>)
}

export function deleteBaseline(id: number) {
  return del<null>(`/baselines/${id}`)
}

export function getDeviations(params?: Record<string, any>) {
  return get<PaginatedResponse<Deviation>>('/deviations', params)
}

export function getDeviation(id: number) {
  return get<Deviation>(`/deviations/${id}`)
}

export function exportDeviations(params?: Record<string, any>): Promise<Blob> {
  return request.get('/deviations/export', { params, responseType: 'blob' }).then((res) => res.data)
}

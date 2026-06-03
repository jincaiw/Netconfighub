import { post, put } from './index'
import type { LoginRequest, LoginResponse, RefreshResponse } from '@/types'

export function login(data: LoginRequest) {
  return post<LoginResponse>('/auth/login', data)
}

export function logout() {
  return post<null>('/auth/logout')
}

export function refreshToken() {
  return post<RefreshResponse>('/auth/refresh')
}

export function changePassword(data: { old_password: string; new_password: string }) {
  return put<null>('/auth/change-password', data)
}

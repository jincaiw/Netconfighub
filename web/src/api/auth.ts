import { post, put } from './index'
import type { LoginRequest, LoginResponse } from '@/types'

// The 401 interceptor in ./index owns the refresh call (with singleton de-dup).
export { refreshToken } from './index'

export function login(data: LoginRequest) {
  return post<LoginResponse>('/auth/login', data)
}

export function logout() {
  return post<null>('/auth/logout')
}

export function changePassword(data: { old_password: string; new_password: string }) {
  return put<null>('/auth/change-password', data)
}

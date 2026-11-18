import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { getToken, setToken, setTokenExpiry, removeToken } from '@/utils/token'
import type { ApiResponse } from '@/types'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

interface RetryableConfig extends InternalAxiosRequestConfig {
  _retry?: boolean
}

// 请求拦截器：自动附加 JWT Token
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// 单例 Promise：并发 401 只触发一次刷新；失败后清空 token 并跳转登录页。
let refreshPromise: Promise<string | null> | null = null
let isRedirecting = false

function isAuthEndpoint(url?: string): boolean {
  return !!url && (url.includes('/auth/refresh') || url.includes('/auth/login'))
}

function redirectToLogin(): void {
  removeToken()
  if (!isRedirecting && !window.location.pathname.includes('/login')) {
    isRedirecting = true
    window.location.href = '/login'
    setTimeout(() => { isRedirecting = false }, 1000)
  }
}

export async function refreshToken(): Promise<string | null> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const response = await request.post<ApiResponse<{ token: string; expires_at: string }>>('/auth/refresh')
        const data = response.data.data
        setToken(data.token)
        if (data.expires_at) {
          setTokenExpiry(data.expires_at)
        }
        return data.token
      } catch {
        removeToken()
        return null
      }
    })().finally(() => {
      // 允许下一次会话再次刷新，但不影响本次并发去重。
      setTimeout(() => { refreshPromise = null }, 0)
    })
  }
  return refreshPromise
}

// 响应拦截器：401 时先尝试刷新一次，成功则重放原请求。
request.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiResponse<unknown>>) => {
    const config = error.config as RetryableConfig | undefined
    const status = error.response?.status
    if (status === 401 && config && !config._retry && !isAuthEndpoint(config.url)) {
      config._retry = true
      const token = await refreshToken()
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
        return request(config)
      }
      redirectToLogin()
      return Promise.reject(error)
    }
    if (status === 401 && isAuthEndpoint(config?.url)) {
      redirectToLogin()
    }
    return Promise.reject(error)
  },
)

// 封装 GET 请求
export function get<T>(url: string, params?: Record<string, any>): Promise<ApiResponse<T>> {
  return request.get(url, { params }).then((res) => res.data)
}

// 封装 POST 请求
export function post<T>(url: string, data?: Record<string, any>): Promise<ApiResponse<T>> {
  return request.post(url, data).then((res) => res.data)
}

// 封装 PUT 请求
export function put<T>(url: string, data?: Record<string, any>): Promise<ApiResponse<T>> {
  return request.put(url, data).then((res) => res.data)
}

// 封装 DELETE 请求
export function del<T>(url: string): Promise<ApiResponse<T>> {
  return request.delete(url).then((res) => res.data)
}

export default request

import axios from 'axios'
import { getToken, removeToken } from '@/utils/token'
import type { ApiResponse } from '@/types'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

// 请求拦截器：自动附加 JWT Token
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 响应拦截器：统一错误处理
let isRedirecting = false

request.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      removeToken()
      if (!isRedirecting && !window.location.pathname.includes('/login')) {
        isRedirecting = true
        window.location.href = '/login'
        setTimeout(() => { isRedirecting = false }, 1000)
      }
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

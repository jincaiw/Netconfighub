import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, logout as logoutApi, refreshToken as refreshTokenApi } from '@/api/auth'
import { reactiveToken, setToken, removeToken, setTokenExpiry, isTokenExpired, getUsername, setUsername } from '@/utils/token'
import type { LoginRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  // 从 utils/token 的共享响应式来源派生，api 层刷新写入后自动同步。
  const token = reactiveToken
  const username = ref<string | null>(getUsername())

  const isLoggedIn = computed(() => !!token.value)

  const isAuthenticated = computed(() => {
    if (!token.value) return false
    return !isTokenExpired()
  })

  async function login(credentials: LoginRequest) {
    const res = await loginApi(credentials)
    setToken(res.data.token)
    if (res.data.expires_at) {
      setTokenExpiry(res.data.expires_at)
    }
    username.value = credentials.username
    setUsername(credentials.username)
  }

  // utils/token 持有单一响应式来源，api 层 401 刷新写入后此处自动派生更新。
  async function refresh() {
    const newToken = await refreshTokenApi()
    if (newToken) {
      setToken(newToken)
    } else {
      removeToken()
      username.value = null
    }
    return token.value
  }

  async function logout() {
    try {
      await logoutApi()
    } finally {
      removeToken()
      username.value = null
    }
  }

  return {
    token,
    username,
    isLoggedIn,
    isAuthenticated,
    login,
    refresh,
    logout,
  }
})

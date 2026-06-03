import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
import { getToken, setToken, removeToken, setTokenExpiry, isTokenExpired, getUsername, setUsername } from '@/utils/token'
import type { LoginRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const username = ref<string | null>(getUsername())

  const isLoggedIn = computed(() => !!token.value)

  const isAuthenticated = computed(() => {
    if (!token.value) return false
    return !isTokenExpired()
  })

  async function login(credentials: LoginRequest) {
    const res = await loginApi(credentials)
    token.value = res.data.token
    setToken(res.data.token)
    if (res.data.expires_at) {
      setTokenExpiry(res.data.expires_at)
    }
    username.value = credentials.username
    setUsername(credentials.username)
  }

  async function logout() {
    try {
      await logoutApi()
    } finally {
      token.value = null
      username.value = null
      removeToken()
    }
  }

  return {
    token,
    username,
    isLoggedIn,
    isAuthenticated,
    login,
    logout,
  }
})

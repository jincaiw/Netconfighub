import { computed, ref } from 'vue'

const TOKEN_KEY = 'nch_token'
const TOKEN_EXPIRY_KEY = 'nch_token_expires'
const USERNAME_KEY = 'nch_username'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

// 单一响应式来源：api 层 401 刷新只写 localStorage，Store 从此处派生，
// 避免 api 层依赖 Pinia（循环依赖 + 需要活动实例）。所有写 token 的入口
// 都必须经过 setToken/removeToken，以保证该 ref 与 localStorage 同步。
const tokenRef = ref<string | null>(getToken())

/** 只读响应式 token，供 Pinia Store 派生使用。 */
export const reactiveToken = computed(() => tokenRef.value)

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
  tokenRef.value = token
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(TOKEN_EXPIRY_KEY)
  localStorage.removeItem(USERNAME_KEY)
  tokenRef.value = null
}

export function getTokenExpiry(): string | null {
  return localStorage.getItem(TOKEN_EXPIRY_KEY)
}

export function setTokenExpiry(expiresAt: string): void {
  localStorage.setItem(TOKEN_EXPIRY_KEY, expiresAt)
}

export function isTokenExpired(): boolean {
  const expiresAt = getTokenExpiry()
  if (!expiresAt) return true
  return new Date(expiresAt).getTime() <= Date.now()
}

export function getUsername(): string | null {
  return localStorage.getItem(USERNAME_KEY)
}

export function setUsername(username: string): void {
  localStorage.setItem(USERNAME_KEY, username)
}

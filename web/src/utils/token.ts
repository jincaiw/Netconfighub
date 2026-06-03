const TOKEN_KEY = 'nch_token'
const TOKEN_EXPIRY_KEY = 'nch_token_expires'
const USERNAME_KEY = 'nch_username'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(TOKEN_EXPIRY_KEY)
  localStorage.removeItem(USERNAME_KEY)
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

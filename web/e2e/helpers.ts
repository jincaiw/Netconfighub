import { Page, request } from '@playwright/test'

const API_BASE = 'http://127.0.0.1:18080/api/v1'

let cachedToken: string | null = null
let cachedPassword = ''

export async function login(page: Page, username = 'admin', password = 'admin') {
  await page.goto('/login')
  await page.locator('input').first().fill(username)
  await page.locator('input[type="password"]').fill(password)
  await page.getByRole('button', { name: '登录' }).click()
  await page.waitForURL('**/dashboard**', { timeout: 10000 })
}

export async function getAuthToken(username = 'admin', password = 'admin'): Promise<string> {
  if (cachedToken && cachedPassword === password) {
    return cachedToken
  }
  const apiContext = await request.newContext()
  const response = await apiContext.post(`${API_BASE}/auth/login`, {
    data: { username, password },
  })
  if (!response.ok()) {
    await apiContext.dispose()
    throw new Error(`Login failed: ${response.status()}`)
  }
  const body = await response.json()
  await apiContext.dispose()
  cachedToken = body.data.token
  cachedPassword = password
  return cachedToken
}

export async function createDevice(token: string, device: Record<string, any>) {
  const apiContext = await request.newContext()
  const payload = { password: 'e2etest123', ...device }
  const response = await apiContext.post(`${API_BASE}/devices`, {
    data: payload,
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok()) {
    await apiContext.dispose()
    throw new Error(`Create device failed: ${response.status()}`)
  }
  const body = await response.json()
  await apiContext.dispose()
  return body.data
}

export async function deleteDevice(token: string, deviceId: number) {
  const apiContext = await request.newContext()
  await apiContext.delete(`${API_BASE}/devices/${deviceId}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  await apiContext.dispose()
}

export async function createGroup(token: string, group: Record<string, any>) {
  const apiContext = await request.newContext()
  const response = await apiContext.post(`${API_BASE}/groups`, {
    data: group,
    headers: { Authorization: `Bearer ${token}` },
  })
  const body = await response.json()
  await apiContext.dispose()
  return body.data
}

export async function deleteGroup(token: string, groupId: number) {
  const apiContext = await request.newContext()
  await apiContext.delete(`${API_BASE}/groups/${groupId}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  await apiContext.dispose()
}

export async function createBaseline(token: string, baseline: Record<string, any>) {
  const apiContext = await request.newContext()
  const response = await apiContext.post(`${API_BASE}/baselines`, {
    data: baseline,
    headers: { Authorization: `Bearer ${token}` },
  })
  const body = await response.json()
  await apiContext.dispose()
  return body.data
}

export async function deleteBaseline(token: string, baselineId: number) {
  const apiContext = await request.newContext()
  await apiContext.delete(`${API_BASE}/baselines/${baselineId}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  await apiContext.dispose()
}

export async function createApiToken(token: string, data: Record<string, any>) {
  const apiContext = await request.newContext()
  const response = await apiContext.post(`${API_BASE}/tokens`, {
    data,
    headers: { Authorization: `Bearer ${token}` },
  })
  const body = await response.json()
  await apiContext.dispose()
  return body.data
}

export async function deleteApiToken(token: string, tokenId: number) {
  const apiContext = await request.newContext()
  await apiContext.delete(`${API_BASE}/tokens/${tokenId}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  await apiContext.dispose()
}

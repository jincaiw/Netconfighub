import { test, expect } from '@playwright/test'
import { createHmac } from 'crypto'
import { getAuthToken } from './helpers'

const API = 'http://127.0.0.1:18080/api/v1'
const E2E_JWT_SECRET = 'e2e-test-secret'

function makeExpiredToken(): string {
  const b64 = (value: object) => Buffer.from(JSON.stringify(value)).toString('base64url')
  const header = b64({ alg: 'HS256', typ: 'JWT' })
  const now = Math.floor(Date.now() / 1000)
  const payload = b64({
    admin_id: 1,
    username: 'admin',
    exp: now - 3600,
    iat: now - 7200,
    iss: 'netconfighub',
  })
  const signature = createHmac('sha256', E2E_JWT_SECRET)
    .update(`${header}.${payload}`)
    .digest('base64url')
  return `${header}.${payload}.${signature}`
}

test.describe('安全功能', () => {
  test('安全响应头存在', async ({ request }) => {
    const response = await request.get(`${API}/health`)
    const headers = response.headers()
    expect(headers['x-content-type-options']).toBe('nosniff')
    expect(headers['x-frame-options']).toBe('DENY')
    expect(headers['referrer-policy']).toBe('no-referrer')
    expect(headers['content-security-policy']).toContain("default-src 'self'")
  })

  test('默认 admin/admin 登录失败（未注入口令的新实例）', async ({ request }) => {
    const response = await request.post(`${API}/auth/login`, {
      data: { username: 'admin', password: 'admin' },
    })
    expect(response.status()).toBe(401)
  })

  test('过期 refresh token 返回 401 且不能换发新令牌', async ({ request }) => {
    const response = await request.post(`${API}/auth/refresh`, {
      headers: { Authorization: `Bearer ${makeExpiredToken()}` },
    })
    expect(response.status()).toBe(401)
    const body = await response.json()
    expect(body.code).toBe(40101)
    expect(body.message).toContain('无效或过期的 JWT Token')
    expect(body.token).toBeUndefined()
  })
  test('CORS头部检查', async ({ request }) => {
    const response = await request.get('http://127.0.0.1:18080/api/v1/health')
    expect(response.headers()['access-control-allow-origin']).toBe('*')
  })

  test('未认证访问被拒绝', async ({ request }) => {
    const response = await request.get('http://127.0.0.1:18080/api/v1/devices')
    expect(response.status()).toBe(401)
  })

  test('API Token认证', async ({ request }) => {
    const token = await getAuthToken()
    const response = await request.get('http://127.0.0.1:18080/api/v1/devices', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    expect(response.ok()).toBeTruthy()
  })

  test('设备密码不在API响应中', async ({ request }) => {
    const token = await getAuthToken()
    const createRes = await request.post('http://127.0.0.1:18080/api/v1/devices', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      data: {
        name: 'security-test-device',
        ip_address: '10.0.0.99',
        vendor: 'cisco',
        model: 'c2960',
        conn_protocol: 'ssh',
        port: 22,
        username: 'admin',
        password: 'secret123',
      },
    })
    expect(createRes.ok()).toBeTruthy()
    const body = await createRes.json()
    expect(body.data.password).toBeUndefined()

    await request.delete(`http://127.0.0.1:18080/api/v1/devices/${body.data.id}`, {
      headers: { 'Authorization': `Bearer ${token}` },
    })
  })

  test('设备导出不含密码', async ({ request }) => {
    const token = await getAuthToken()
    const response = await request.get('http://127.0.0.1:18080/api/v1/devices/export', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    expect(response.ok()).toBeTruthy()
    const text = await response.text()
    expect(text).not.toContain('secret123')
  })
})

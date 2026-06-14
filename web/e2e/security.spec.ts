import { test, expect } from '@playwright/test'
import { getAuthToken } from './helpers'

test.describe('安全功能', () => {
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

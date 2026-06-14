import { test, expect } from '@playwright/test'
import { getAuthToken } from './helpers'

test.describe('系统配置API', () => {
  test('获取配置列表', async ({ request }) => {
    const token = await getAuthToken()
    const response = await request.get('http://127.0.0.1:18080/api/v1/configs', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    expect(response.ok()).toBeTruthy()
    const body = await response.json()
    expect(Array.isArray(body.data)).toBeTruthy()
  })

  test('设置和获取配置', async ({ request }) => {
    const token = await getAuthToken()
    const configKey = `e2e_test_config_${Date.now()}`
    const setRes = await request.put('http://127.0.0.1:18080/api/v1/configs', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      data: {
        key: configKey,
        value: 'test_value',
      },
    })
    expect(setRes.ok()).toBeTruthy()

    const getRes = await request.get('http://127.0.0.1:18080/api/v1/configs', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    expect(getRes.ok()).toBeTruthy()
    const body = await getRes.json()
    const found = body.data.find((c: { key: string; value: string }) => c.key === configKey)
    expect(found).toBeDefined()
    expect(found.value).toBe('test_value')
  })

  test('设置脱敏配置', async ({ request }) => {
    const token = await getAuthToken()
    const response = await request.put('http://127.0.0.1:18080/api/v1/configs', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      data: {
        key: 'sanitize.enabled',
        value: 'true',
      },
    })
    expect(response.ok()).toBeTruthy()
  })
})

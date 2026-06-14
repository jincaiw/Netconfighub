import { test, expect } from '@playwright/test'
import { login, getAuthToken } from './helpers'

const createdHookIds: number[] = []

test.afterAll(async () => {
  const token = await getAuthToken()
  for (const id of createdHookIds) {
    try {
      await fetch(`http://127.0.0.1:18080/api/v1/hooks/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` },
      })
    } catch {}
  }
})

test.describe('Hook管理', () => {
  test('Hook列表页面', async ({ page }) => {
    await login(page)
    await page.goto('/settings')
    await expect(page.getByText('API Token')).toBeVisible()
  })

  test('通过API创建Webhook Hook', async ({ request }) => {
    const token = await getAuthToken()
    const response = await request.post('http://127.0.0.1:18080/api/v1/hooks', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      data: {
        name: 'Test Webhook',
        type: 'webhook',
        config: JSON.stringify({ url: 'https://example.com/hook', secret: 'test-secret' }),
        events: 'backup_success,backup_failed',
        enabled: true,
      },
    })
    expect(response.ok()).toBeTruthy()
    const body = await response.json()
    expect(body.data.name).toBe('Test Webhook')
    expect(body.data.type).toBe('webhook')
    if (body.data.id) createdHookIds.push(body.data.id)
  })

  test('通过API查询Hook列表', async ({ request }) => {
    const token = await getAuthToken()
    const response = await request.get('http://127.0.0.1:18080/api/v1/hooks', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    expect(response.ok()).toBeTruthy()
    const body = await response.json()
    expect(Array.isArray(body.data.items)).toBeTruthy()
  })

  test('通过API删除Hook', async ({ request }) => {
    const token = await getAuthToken()
    const createRes = await request.post('http://127.0.0.1:18080/api/v1/hooks', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      data: {
        name: 'Delete Test Hook',
        type: 'exec',
        config: JSON.stringify({ command: 'echo test' }),
        events: 'backup_success',
        enabled: false,
      },
    })
    const createBody = await createRes.json()
    const hookId = createBody.data.id

    const deleteRes = await request.delete(`http://127.0.0.1:18080/api/v1/hooks/${hookId}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    expect(deleteRes.ok()).toBeTruthy()
  })
})

import { test, expect } from '@playwright/test'
import { login, getAuthToken, createApiToken, deleteApiToken } from './helpers'

test.describe('设置页面', () => {
  let token: string
  const createdTokenIds: number[] = []

  test.beforeAll(async () => {
    token = await getAuthToken()
  })

  test.afterAll(async () => {
    for (const id of createdTokenIds) {
      try {
        await deleteApiToken(token, id)
      } catch {
      // cleanup is best-effort
    }
    }
  })

  test('打开设置页', async ({ page }) => {
    await login(page)
    await page.goto('/settings')
    await expect(page.getByText('API Token 管理')).toBeVisible()
    await expect(page.getByText('备份策略配置')).toBeVisible()
  })

  test('创建 API Token', async ({ page }) => {
    await login(page)
    await page.goto('/settings')
    await expect(page.getByText('API Token 管理')).toBeVisible()
    await page.getByRole('button', { name: '创建 Token' }).click()
    await expect(page.locator('.t-dialog:visible')).toBeVisible()
    await page.getByPlaceholder('请输入 Token 名称').fill(`e2e-token-${Date.now()}`)
    await page.locator('.t-date-picker').locator('input').click()
    const dateCell = page.locator('.t-date-picker__cell:not(.t-date-picker__cell--disabled)').last()
    await expect(dateCell).toBeVisible()
    await dateCell.click()
    const dateConfirm = page.locator('.t-date-picker__footer .t-button').filter({ hasText: '确定' })
    if (await dateConfirm.count() > 0) {
      await expect(dateConfirm).toBeVisible()
      await dateConfirm.click()
    }
    await page.locator('.t-dialog:visible').getByRole('button', { name: '创建', exact: true }).click()
    const successDialog = page.locator('.t-dialog:visible').getByText('Token 创建成功')
    await expect(successDialog).toBeVisible({ timeout: 5000 })
  })

  test('查看 Token 列表', async ({ page }) => {
    const apiToken = await createApiToken(token, {
      name: `e2e-list-token-${Date.now()}`,
      expires_at: '2099-12-31T23:59:59+08:00',
    })
    if (apiToken?.id) createdTokenIds.push(apiToken.id)
    await login(page)
    await page.goto('/settings')
    await expect(page.getByText('API Token 管理')).toBeVisible()
    await expect(page.locator('table')).toBeVisible()
  })

  test('删除 Token', async ({ page }) => {
    const apiToken = await createApiToken(token, {
      name: `e2e-delete-token-${Date.now()}`,
      expires_at: '2099-12-31T23:59:59+08:00',
    })
    await login(page)
    await page.goto('/settings')
    await expect(page.locator('table')).toBeVisible()
    const row = page.locator('tr').filter({ hasText: apiToken.name })
    await row.getByRole('button', { name: '删除' }).click()
    await page.getByRole('button', { name: '删除' }).last().click()
    await expect(page.locator('tr').filter({ hasText: apiToken.name })).toHaveCount(0)
  })

  test('系统配置加载失败时报错并锁定保存，重试后恢复', async ({ page }) => {
    let fail = true
    await page.goto('/login')
    await page.route(/\/api\/v1\/configs(?:\?|$)/, (route) =>
      route.request().method() === 'GET' && fail
        ? route.fulfill({ status: 500, json: { code: 500, message: 'boom' } })
        : route.continue(),
    )
    await login(page)
    await page.goto('/settings')

    await expect(page.locator('.t-alert:visible').filter({ hasText: '系统配置加载失败' })).toBeVisible()
    const saveConfig = page.getByRole('button', { name: '保存配置' })
    const saveSanitize = page.getByRole('button', { name: '保存脱敏配置' })
    await expect(saveConfig).toBeDisabled()
    await expect(saveSanitize).toBeDisabled()

    fail = false
    await page.getByRole('button', { name: '重试' }).click()
    await expect(saveConfig).toBeEnabled({ timeout: 10000 })
    await expect(saveSanitize).toBeEnabled({ timeout: 10000 })
  })
})

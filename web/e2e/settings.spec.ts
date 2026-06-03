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
      } catch {}
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
    await expect(page.locator('.n-dialog')).toBeVisible()
    await page.getByPlaceholder('请输入 Token 名称').fill(`e2e-token-${Date.now()}`)
    await page.locator('.n-date-picker').locator('input').click()
    await page.waitForTimeout(500)
    await page.locator('.n-date-panel-dates .n-date-panel-date:not(.n-date-panel-date--disabled)').last().click()
    await page.waitForTimeout(300)
    await page.locator('.n-date-panel-actions .n-button').filter({ hasText: 'Confirm' }).click()
    await page.waitForTimeout(500)
    await page.locator('.n-dialog').getByRole('button', { name: '创建', exact: true }).click()
    await page.waitForTimeout(3000)
    const successDialog = page.locator('.n-dialog').getByText('Token 创建成功')
    const hasSuccess = await successDialog.isVisible().catch(() => false)
    expect(hasSuccess).toBeTruthy()
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
    await page.waitForTimeout(2000)
  })
})

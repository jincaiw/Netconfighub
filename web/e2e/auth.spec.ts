import { test, expect } from '@playwright/test'
import { login } from './helpers'

test.describe('登录/登出流程', () => {
  test('正确用户名密码登录成功，跳转到仪表盘', async ({ page }) => {
    await login(page, 'admin', 'admin')
    await expect(page).toHaveURL(/\/dashboard/)
    await expect(page.getByText('设备总数')).toBeVisible()
  })

  test('错误密码登录失败，显示错误提示', async ({ page }) => {
    await page.goto('/login')
    await page.locator('input').first().fill('admin')
    await page.locator('input[type="password"]').fill('wrongpassword')
    await page.getByRole('button', { name: '登录' }).click()
    await page.waitForTimeout(2000)
    await expect(page).toHaveURL(/\/login/)
  })

  test('未登录访问受保护页面，跳转到登录页', async ({ page }) => {
    await page.goto('/devices')
    await page.waitForURL('**/login**', { timeout: 10000 })
    await expect(page).toHaveURL(/\/login/)
  })

  test('登出后跳转到登录页', async ({ page }) => {
    await login(page, 'admin', 'admin')
    await expect(page).toHaveURL(/\/dashboard/)
    await page.getByRole('button', { name: '退出' }).click()
    await page.waitForURL('**/login**', { timeout: 10000 })
    await expect(page).toHaveURL(/\/login/)
  })
})

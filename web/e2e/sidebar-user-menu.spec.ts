import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 回归：侧栏左下角用户卡点击应弹出用户菜单（个人设置 / 退出登录），修复“点击无反应”。
test.describe('侧栏用户菜单', () => {
  test('点击左下角用户卡弹出菜单并可进入系统设置', async ({ page }) => {
    await login(page)
    await page.locator('.sider-user').click()
    await expect(page.getByText('个人设置', { exact: true })).toBeVisible()
    await expect(page.getByText('退出登录', { exact: true })).toBeVisible()
    await page.getByText('个人设置', { exact: true }).click()
    await expect(page).toHaveURL(/\/settings/)
  })

  test('点击退出登录回到登录页', async ({ page }) => {
    await login(page)
    await page.locator('.sider-user').click()
    await page.getByText('退出登录', { exact: true }).click()
    await expect(page).toHaveURL(/\/login/, { timeout: 10000 })
  })
})

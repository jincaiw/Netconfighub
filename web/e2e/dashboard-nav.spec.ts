import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 回归：TDesign 迁移审查发现的 P2 —— OutputOverview 的“查看设备/查看全部”按钮
// 曾因 Dashboard.vue 未绑定 @navigate 而成为死按钮。修复后点击应正确路由跳转。
test.describe('首页任务概览跳转（P2 回归）', () => {
  test('“查看设备”按钮跳转到设备管理页', async ({ page }) => {
    await login(page)
    await expect(page.getByText('最近备份任务')).toBeVisible()
    await page.getByRole('button', { name: '查看设备' }).click()
    await expect(page).toHaveURL(/\/devices$/)
    await expect(page.locator('.t-card__title').filter({ hasText: '设备管理' })).toBeVisible()
  })

  test('“查看全部”按钮跳转到失败设备页', async ({ page }) => {
    await login(page)
    await expect(page.getByText('失败设备概览')).toBeVisible()
    await page.getByRole('button', { name: '查看全部' }).click()
    await expect(page).toHaveURL(/\/failed$/)
    await expect(page.locator('.t-card__title').filter({ hasText: '失败设备列表' })).toBeVisible()
  })
})

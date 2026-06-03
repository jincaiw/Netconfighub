import { test, expect } from '@playwright/test'
import { login } from './helpers'

test.describe('失败设备流程', () => {
  test('打开失败设备列表页', async ({ page }) => {
    await login(page)
    await page.goto('/failed')
    await expect(page.getByText('失败设备列表')).toBeVisible()
  })

  test('查看失败设备信息', async ({ page }) => {
    await login(page)
    await page.goto('/failed')
    await expect(page.getByText('失败设备列表')).toBeVisible()
    await page.waitForTimeout(2000)
    const empty = page.getByText('暂无失败设备')
    const table = page.locator('.n-data-table')
    const hasEmpty = await empty.isVisible().catch(() => false)
    if (!hasEmpty) {
      await expect(table).toBeVisible({ timeout: 10000 })
    }
  })

  test('点击重新触发备份按钮', async ({ page }) => {
    await login(page)
    await page.goto('/failed')
    await expect(page.getByText('失败设备列表')).toBeVisible()
    const retryBtn = page.locator('table').getByRole('button', { name: '重新备份' }).first()
    const hasRetryBtn = await retryBtn.isVisible().catch(() => false)
    if (hasRetryBtn) {
      await retryBtn.click()
      await expect(page.getByText('备份任务已触发')).toBeVisible({ timeout: 5000 })
    }
  })
})

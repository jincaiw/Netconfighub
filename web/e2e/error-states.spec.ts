import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 列表加载失败应显示内联错误条 + 重新加载，避免被误显示为“暂无数据”。
test.describe('列表加载失败内联错误态', () => {
  test('设备列表失败显示内联错误与重新加载', async ({ page }) => {
    await login(page)
    await page.route(/\/api\/v1\/devices(\?|$)/, (r) => r.fulfill({ status: 500, json: { code: 500, message: 'boom' } }))
    await page.goto('/devices')
    await page.waitForTimeout(900)
    const alert = page.locator('.list-error')
    await expect(alert).toBeVisible()
    await expect(alert).toContainText('设备列表加载失败')
  })

  test('基线列表失败显示内联错误', async ({ page }) => {
    await login(page)
    await page.route(/\/api\/v1\/baselines(\?|$)/, (r) => r.fulfill({ status: 500, json: { code: 500, message: 'boom' } }))
    await page.goto('/baselines')
    await page.waitForTimeout(900)
    await expect(page.locator('.list-error')).toContainText('基线列表加载失败')
  })

  test('不存在的设备显示“未找到”态且不显示操作按钮', async ({ page }) => {
    await login(page)
    await page.goto('/devices/99999')
    await page.waitForTimeout(900)
    await expect(page.getByText('设备不存在或已被删除')).toBeVisible()
    await expect(page.getByRole('button', { name: '触发备份' })).toHaveCount(0)
  })
})


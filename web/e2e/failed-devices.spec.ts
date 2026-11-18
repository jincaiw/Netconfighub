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
    const empty = page.getByText('暂无失败设备')
    const table = page.locator('.t-table')
    await expect(empty.or(table)).toBeVisible({ timeout: 10000 })
  })

  test('点击重新触发备份按钮', async ({ page }) => {
    await login(page)
    await page.goto('/failed')
    await expect(page.getByText('失败设备列表')).toBeVisible()
    const retryBtn = page.locator('table').getByRole('button', { name: '重新备份' }).first()
    const retryCount = await retryBtn.count()
    if (retryCount > 0) {
      await expect(retryBtn.first()).toBeVisible()
      await retryBtn.first().click()
      await expect(page.getByText('备份任务已触发')).toBeVisible({ timeout: 5000 })
    }
  })

  test('失败设备时间列格式化为 19 字符、不折行且保留原始 title', async ({ page }) => {
    const raw = '2026-01-02T03:04:05+08:00'
    await page.goto('/login')
    await page.route(/\/api\/v1\/failed-devices(?:\?|$)/, (route) =>
      route.fulfill({
        json: {
          code: 0,
          message: 'ok',
          data: {
            items: [
              {
                id: 90001,
                name: 'time-format-device',
                ip: '10.2.3.4',
                vendor: 'cisco',
                last_backup_status: 'failed',
                last_failure_reason: 'timeout',
                retry_count: 2,
                last_backup_at: raw,
              },
            ],
            total: 1,
          },
        },
      }),
    )
    await login(page)
    await page.goto('/failed')
    const cell = page.locator('.datetime-cell').first()
    await expect(cell).toBeVisible()
    const text = (await cell.textContent())?.trim() || ''
    expect(text).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)
    await expect(cell).toHaveAttribute('title', raw)
    expect(await cell.evaluate((el) => getComputedStyle(el).whiteSpace)).toBe('nowrap')
  })
})

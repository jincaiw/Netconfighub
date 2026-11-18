import { test, expect } from '@playwright/test'
import { login } from './helpers'

test.describe('告警计数去重', () => {
  test('首屏只发起一次 unread-count 请求', async ({ page }) => {
    let requests = 0
    page.on('request', (req) => {
      if (req.url().includes('/api/v1/alerts/unread-count')) {
        requests++
      }
    })

    await login(page)
    await expect(page.getByText('设备总数')).toBeVisible()
    await expect.poll(() => requests).toBeGreaterThan(0)
    // 修复前 AppHeader 与 Dashboard 会各发一次（共 2 次），共享 composable 后最多 1 次
    expect(requests).toBeLessThanOrEqual(1)
  })
})

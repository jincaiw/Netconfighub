import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 仪表盘中文化与内容补齐：无英文文案、报表/通知标签有真实内容、导出可用。
test.describe('仪表盘中文化与内容', () => {
  test('页头与标签页均为中文，无英文可见文案', async ({ page }) => {
    await login(page)
    await page.goto('/dashboard')
    await page.waitForSelector('[data-testid="dashboard-loading"]', { state: 'hidden' }).catch(() => {})
    await expect(page.getByRole('heading', { name: '运维总览' })).toBeVisible()
    await expect(page.getByRole('button', { name: '导出报表' })).toBeVisible()
    for (const label of ['概览', '分析', '报表', '通知']) {
      await expect(page.locator('.dashboard-tabs .t-tabs__nav-item', { hasText: label }).first()).toBeVisible()
    }
    // 页面不应再出现旧的英文标签/标题
    await expect(page.getByText('Dashboard', { exact: true })).toHaveCount(0)
    await expect(page.getByText('Overview', { exact: true })).toHaveCount(0)
  })

  test('报表标签页展示汇总与导出按钮', async ({ page }) => {
    await login(page)
    await page.goto('/dashboard')
    await page.waitForSelector('[data-testid="dashboard-loading"]', { state: 'hidden' }).catch(() => {})
    await page.locator('.dashboard-tabs .t-tabs__nav-item', { hasText: '报表' }).first().click()
    await expect(page.getByText('报表汇总')).toBeVisible()
    await expect(page.getByRole('button', { name: /导出设备清单/ })).toBeVisible()
    await expect(page.getByRole('button', { name: /导出备份记录/ })).toBeVisible()
  })

  test('通知标签页展示最近通知卡片', async ({ page }) => {
    await login(page)
    await page.goto('/dashboard')
    await page.waitForSelector('[data-testid="dashboard-loading"]', { state: 'hidden' }).catch(() => {})
    await page.locator('.dashboard-tabs .t-tabs__nav-item', { hasText: '通知' }).first().click()
    await expect(page.getByText('最近通知')).toBeVisible()
  })

  test('导出报表按钮触发 CSV 下载', async ({ page }) => {
    await login(page)
    await page.goto('/dashboard')
    await page.waitForSelector('[data-testid="dashboard-loading"]', { state: 'hidden' }).catch(() => {})
    const [download] = await Promise.all([
      page.waitForEvent('download', { timeout: 8000 }),
      page.getByRole('button', { name: '导出报表' }).click(),
    ])
    expect(download.suggestedFilename()).toContain('.csv')
  })
})

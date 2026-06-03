import { test, expect } from '@playwright/test'
import { login, getAuthToken } from './helpers'

test.describe('告警管理', () => {
  test('告警页面加载', async ({ page }) => {
    await login(page)
    await page.goto('/alerts')
    await expect(page.locator('.n-card-header__main').filter({ hasText: '告警管理' })).toBeVisible()
  })

  test('告警筛选功能', async ({ page }) => {
    await login(page)
    await page.goto('/alerts')
    await expect(page.locator('.n-card-header__main').filter({ hasText: '告警管理' })).toBeVisible()
    const typeSelect = page.locator('.n-form-item').filter({ hasText: '按类型筛选' }).locator('.n-base-selection')
    if (await typeSelect.isVisible()) {
      await typeSelect.click()
      await page.waitForTimeout(300)
      const option = page.locator('.n-base-select-option').first()
      if (await option.isVisible()) {
        await option.click()
      }
    }
  })

  test('全部已读按钮', async ({ page }) => {
    await login(page)
    await page.goto('/alerts')
    await expect(page.locator('.n-card-header__main').filter({ hasText: '告警管理' })).toBeVisible()
    const markAllBtn = page.getByRole('button', { name: '全部已读' })
    if (await markAllBtn.isVisible() && await markAllBtn.isEnabled()) {
      await markAllBtn.click()
    }
  })
})

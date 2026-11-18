import { test, expect } from '@playwright/test'
import { login } from './helpers'

test.describe('告警管理', () => {
  test('告警页面加载', async ({ page }) => {
    await login(page)
    await page.goto('/alerts')
    await expect(page.locator('.t-card__title').filter({ hasText: '告警管理' })).toBeVisible()
  })

  test('告警筛选功能', async ({ page }) => {
    await login(page)
    await page.goto('/alerts')
    await expect(page.locator('.t-card__title').filter({ hasText: '告警管理' })).toBeVisible()
    const typeSelect = page.getByPlaceholder('按类型筛选')
    await expect(typeSelect).toBeVisible()
    await typeSelect.click()
    const option = page.locator('.t-select-option').first()
    await expect(option).toBeVisible()
    await option.click()
  })

  test('全部已读按钮', async ({ page }) => {
    await login(page)
    await page.goto('/alerts')
    await expect(page.locator('.t-card__title').filter({ hasText: '告警管理' })).toBeVisible()
    const markAllBtn = page.getByRole('button', { name: '全部已读' })
    await expect(markAllBtn).toBeVisible()
    // 按钮在无告警时为 disabled，具体点击沿用数据相关的业务步骤
    if (await markAllBtn.isEnabled()) {
      await markAllBtn.click()
    }
  })
})

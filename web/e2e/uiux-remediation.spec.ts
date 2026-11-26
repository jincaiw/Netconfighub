import { test, expect } from '@playwright/test'
import { login } from './helpers'

// UI/UX 整改回归：页面标题随路由变化、设备列表无 1080px 强制宽导致的横向溢出。
test.describe('UI/UX 整改项回归', () => {
  test('浏览器标题随路由变化（meta.title + document.title）', async ({ page }) => {
    await login(page)
    await expect(page).toHaveTitle(/运维总览 · NetConfigHub/)
    await page.goto('/devices')
    await expect(page).toHaveTitle(/设备管理 · NetConfigHub/)
    await page.goto('/settings')
    await expect(page).toHaveTitle(/系统设置 · NetConfigHub/)
  })

  test('设备列表在 1024 视口无横向溢出', async ({ page }) => {
    await login(page)
    await page.goto('/devices')
    await page.setViewportSize({ width: 1024, height: 800 })
    await page.waitForTimeout(800)
    const overflow = await page.evaluate(
      () => document.documentElement.scrollWidth - document.documentElement.clientWidth,
    )
    expect(overflow).toBeLessThanOrEqual(1)
  })
})

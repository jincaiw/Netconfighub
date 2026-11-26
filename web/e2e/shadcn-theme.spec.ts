import { test, expect } from '@playwright/test'
import { login } from './helpers'

// shadcn-admin 主题复刻回归：中性配色令牌、Tab 结构、紧凑 KPI 卡、横向顶栏导航。
test.describe('shadcn-admin 主题复刻', () => {
  test('中性主题令牌生效（brand 为近黑/近白，非品牌蓝）', async ({ page }) => {
    await login(page)
    const brand = await page.evaluate(() =>
      getComputedStyle(document.documentElement).getPropertyValue('--td-brand-color').trim(),
    )
    // shadcn light primary = #27272a（zinc-800），dark = #e4e4e7；均非 TDesign 蓝 #0052d9/#4582e6
    expect(brand).not.toBe('#0052d9')
    expect(brand).not.toBe('#4582e6')
    expect(['#27272a', '#e4e4e7']).toContain(brand)
  })

  test('仪表盘页头与中文标签页（概览/分析/报表/通知）', async ({ page }) => {
    await login(page)
    await expect(page.getByRole('heading', { name: '运维总览' })).toBeVisible()
    await expect(page.getByRole('button', { name: '导出报表' })).toBeVisible()
    await expect(page.getByText('概览', { exact: true }).first()).toBeVisible()
    await expect(page.getByText('分析', { exact: true }).first()).toBeVisible()
    await expect(page.getByText('报表', { exact: true }).first()).toBeVisible()
    await expect(page.getByText('通知', { exact: true }).first()).toBeVisible()
  })

  test('顶部显示当前页面标题（导航收敛到侧栏，无重复菜单）', async ({ page }) => {
    await login(page)
    // 顶部不再有重复的全局导航菜单
    await expect(page.locator('.top-nav')).toHaveCount(0)
    const title = page.locator('.page-title')
    await expect(title).toBeVisible()
    await expect(title).toHaveText('运维总览')
    await page.goto('/devices')
    await expect(page.locator('.page-title')).toHaveText('设备管理')
  })

  test('紧凑 KPI 卡形态保留（4 卡、首卡主色、脚注）', async ({ page }) => {
    await login(page)
    await expect(page.locator('.metric-card')).toHaveCount(4)
    await expect(page.locator('.dashboard-item--main-color')).toHaveCount(1)
    await expect(page.locator('.metric-footnote').first()).toBeVisible()
  })
})

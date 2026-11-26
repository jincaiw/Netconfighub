import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 侧边栏显示合理性：分组清晰、仅保留一个折叠控件、主题与导航不回归。
test.describe('侧边栏显示与分组', () => {
  test('菜单分为总览/运维监控/系统三组，且无底部孤立折叠按钮', async ({ page }) => {
    await login(page)

    // 三个分组标题（限定桌面侧栏，排除移动端抽屉的副本）
    const sider = page.locator('.desktop-sider')
    await expect(sider.locator('.t-menu-group__title').filter({ hasText: '总览' })).toBeVisible()
    await expect(sider.locator('.t-menu-group__title').filter({ hasText: '运维监控' })).toBeVisible()
    await expect(sider.locator('.t-menu-group__title').filter({ hasText: '系统' })).toBeVisible()

    // 7 项菜单归类：运维监控 5 项
    await expect(sider.locator('.app-menu .t-menu__item')).toHaveCount(7)

    // 折叠控件只在顶部（aria-label），侧栏 aside 内不再有折叠按钮
    await expect(page.getByRole('button', { name: /收起侧边栏|展开侧边栏/ })).toHaveCount(1)
    expect(await page.locator('.desktop-sider .t-button').count()).toBe(0)
  })

  test('顶部折叠按钮可折叠/展开侧栏', async ({ page }) => {
    await login(page)
    const toggle = page.getByRole('button', { name: /收起侧边栏|展开侧边栏/ })
    await toggle.click()
    await expect(page.locator('.desktop-sider')).toHaveCSS('width', '64px')
    await toggle.click()
    await expect(page.locator('.desktop-sider')).toHaveCSS('width', '220px')
  })

  test('分组后导航仍正确', async ({ page }) => {
    await login(page)
    await page.locator('.desktop-sider .app-menu .t-menu__item').filter({ hasText: '设备管理' }).click()
    await expect(page).toHaveURL(/\/devices$/)
  })

  test('折叠后顶部 logo 仅显示居中徽标，无 NCH 文字/溢出', async ({ page }) => {
    await login(page)
    await page.getByRole('button', { name: '收起侧边栏' }).click()
    await expect(page.locator('.desktop-sider')).toHaveCSS('width', '64px')
    const logo = page.locator('.desktop-sider .sider-logo')
    await expect(logo).toHaveClass(/sider-logo--collapsed/)
    await expect(logo.locator('.logo-text-short')).toHaveCount(0)
    await expect(logo.locator('.team-badge')).toBeVisible()
    const ov = await logo.evaluate((el) => el.scrollWidth - el.clientWidth)
    expect(ov).toBeLessThanOrEqual(1)
  })
})


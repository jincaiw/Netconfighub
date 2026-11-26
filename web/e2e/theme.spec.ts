import { test, expect, type Page } from '@playwright/test'
import { login } from './helpers'

function getThemeMode(page: Page) {
  return page.evaluate(() => document.documentElement.getAttribute('theme-mode'))
}

test.describe('主题模式', () => {
  test.beforeEach(async ({ page }) => {
    // 清空持久化值，保证每个用例从默认状态开始
    await page.goto('/login')
    await page.evaluate(() => window.localStorage.removeItem('nch-theme-mode'))
    await login(page)
  })

  test('手动切换深色/浅色会更新 theme-mode', async ({ page }) => {
    await page.getByTestId('theme-switch').click()
    const darkItem = page.locator('.t-dropdown__item:visible').filter({ hasText: '深色' })
    await darkItem.waitFor({ state: 'visible', timeout: 15000 })
    await darkItem.click()
    await expect.poll(() => getThemeMode(page), { timeout: 15000 }).toBe('dark')

    await page.getByTestId('theme-switch').click()
    const lightItem = page.locator('.t-dropdown__item:visible').filter({ hasText: '浅色' })
    await lightItem.waitFor({ state: 'visible', timeout: 15000 })
    await lightItem.click()
    await expect.poll(() => getThemeMode(page), { timeout: 15000 }).toBe('')
  })

  test('刷新后保持上次选择', async ({ page }) => {
    await page.getByTestId('theme-switch').click()
    await page.locator('.t-dropdown__item:visible').filter({ hasText: '深色' }).click()
    await expect.poll(() => getThemeMode(page)).toBe('dark')

    await page.reload()
    await expect(page).toHaveURL(/\/dashboard/)
    await expect.poll(() => getThemeMode(page)).toBe('dark')
  })

  test('跟随系统模式随系统偏好变化', async ({ page }) => {
    await page.emulateMedia({ colorScheme: 'dark' })
    await page.getByTestId('theme-switch').click()
    await page.locator('.t-dropdown__item:visible').filter({ hasText: '跟随系统' }).click()
    await expect.poll(() => getThemeMode(page)).toBe('dark')

    await page.emulateMedia({ colorScheme: 'light' })
    await expect.poll(() => getThemeMode(page)).toBe('')
  })

  test('深色模式下侧栏跟随主题', async ({ page }) => {
    await page.getByTestId('theme-switch').click()
    await page.locator('.t-dropdown__item:visible').filter({ hasText: '深色' }).click()
    await expect.poll(() => getThemeMode(page)).toBe('dark')
    await expect(page.locator('.app-menu.t-menu--dark:visible')).toBeVisible()
  })

  test('head 主题脚本先于模块脚本写入首屏主题', async ({ page }) => {
    const response = await page.request.get('/')
    const html = await response.text()
    // 主题初始化脚本改为独立文件（满足生产 CSP script-src 'self'），且在入口模块前加载。
    expect(html).toContain('src="/theme-init.js"')
    expect(html.indexOf('src="/theme-init.js"')).toBeLessThan(html.indexOf('src="/src/main.ts"'))
    const js = await (await page.request.get('/theme-init.js')).text()
    expect(js).toContain('nch-theme-mode')
    expect(js).toContain('prefers-color-scheme: dark')

    await page.goto('/login')
    await page.evaluate(() => window.localStorage.setItem('nch-theme-mode', 'dark'))
    await page.reload()
    expect(await getThemeMode(page)).toBe('dark')

    await page.evaluate(() => window.localStorage.setItem('nch-theme-mode', 'light'))
    await page.reload()
    expect(await getThemeMode(page)).toBe('')
  })

  test('localStorage 不可用时按系统深色降级并保持 DOM 与 Store 一致', async ({ page }) => {
    await page.emulateMedia({ colorScheme: 'dark' })
    // 模拟主题键读取抛错：head 脚本与 setting.ts 的 readStoredMode 均应回退 auto。
    await page.addInitScript(() => {
      const original = Storage.prototype.getItem
      Storage.prototype.getItem = function (key: string) {
        if (key === 'nch-theme-mode') throw new Error('localStorage unavailable')
        return original.call(this, key)
      }
    })
    await page.reload()
    await expect(page).toHaveURL(/\/dashboard/)
    await expect.poll(() => getThemeMode(page)).toBe('dark')
    await expect(page.locator('.app-menu.t-menu--dark:visible')).toBeVisible()
  })

  test('主题下拉仅当前项显示选中标记且随切换更新', async ({ page }) => {
    await page.getByTestId('theme-switch').click()
    await expect(page.getByTestId('theme-selected-mark')).toHaveCount(1)
    const autoItem = page.locator('.t-dropdown__item:visible').filter({ hasText: '跟随系统' })
    await expect(autoItem).toContainText('当前主题')
    await page.locator('.t-dropdown__item:visible').filter({ hasText: '深色' }).click()
    await expect.poll(() => getThemeMode(page)).toBe('dark')
    await page.getByTestId('theme-switch').click()
    await expect(page.getByTestId('theme-selected-mark')).toHaveCount(1)
    const selected = page.locator('.t-dropdown__item:visible').filter({ has: page.getByTestId('theme-selected-mark') })
    await expect(selected).toContainText('深色')
  })
})

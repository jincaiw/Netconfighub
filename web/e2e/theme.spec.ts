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
    await page.locator('.t-dropdown__item:visible').filter({ hasText: '深色' }).click()
    await expect.poll(() => getThemeMode(page)).toBe('dark')

    await page.getByTestId('theme-switch').click()
    await page.locator('.t-dropdown__item:visible').filter({ hasText: '浅色' }).click()
    await expect.poll(() => getThemeMode(page)).toBe('')
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

  test('head 内联主题脚本先于模块脚本写入首屏主题', async ({ page }) => {
    const response = await page.request.get('/')
    const html = await response.text()
    expect(html).toContain('nch-theme-mode')
    expect(html).toContain('prefers-color-scheme: dark')
    // Vite 开发服务器会在 head 顶部注入 @vite/client 模块，因此与应用入口模块比较。
    expect(html.indexOf('nch-theme-mode')).toBeLessThan(html.indexOf('src="/src/main.ts"'))

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
})

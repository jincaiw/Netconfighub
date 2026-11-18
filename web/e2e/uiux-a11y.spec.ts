import { test, expect, type Page } from '@playwright/test'
import { login } from './helpers'

const unreadUrl = /\/api\/v1\/alerts\/unread-count/
const alertsListUrl = /\/api\/v1\/alerts(?:\?|$)/
const markReadUrl = /\/api\/v1\/alerts\/\d+\/read/
const devicesUrl = /\/api\/v1\/devices(?:\?|$)/

function ok<T>(data: T) {
  return { code: 0, message: 'ok', data }
}

async function mockAlerts(page: Page, failCount = false) {
  await page.route(unreadUrl, (route) =>
    failCount
      ? route.fulfill({ status: 500, json: { code: 500, message: 'boom' } })
      : route.fulfill({ json: ok({ count: 1 }) }),
  )
  await page.route(alertsListUrl, (route) =>
    route.fulfill({
      json: ok({
        items: [
          {
            id: 99001,
            type: 'backup_failed',
            device_id: 1,
            device_name: 'a11y-device',
            title: '备份失败告警',
            message: '测试用未读告警',
            severity: 'error',
            is_read: false,
            created_at: new Date().toISOString(),
          },
        ],
        total: 1,
      }),
    }),
  )
  await page.route(markReadUrl, (route) => route.fulfill({ json: ok(null) }))
}

test.describe('UI/UX 整改可访问性与反馈回归', () => {
  test('顶部图标按钮均有无障碍名称', async ({ page }) => {
    await login(page)
    await expect(page.getByRole('button', { name: /收起侧边栏|展开侧边栏/ })).toBeVisible()
    await expect(page.getByRole('button', { name: '告警通知' })).toBeVisible()
    await expect(page.getByRole('button', { name: /主题设置/ })).toBeVisible()
  })

  test('告警项键盘 Enter 可达，未读有非颜色标记', async ({ page }) => {
    await page.goto('/login')
    await mockAlerts(page)
    await login(page)
    await page.getByRole('button', { name: '告警通知' }).click()

    const alertItem = page.getByRole('button', { name: /备份失败告警.*未读/ }).first()
    await expect(alertItem).toBeVisible()
    await expect(page.locator('.alert-unread-dot').first()).toBeVisible()
    await expect(page.locator('.sr-only').filter({ hasText: '未读' })).toHaveCount(1)

    await alertItem.focus()
    await page.keyboard.press('Enter')
    await expect(page.getByText('已标记为已读')).toBeVisible({ timeout: 5000 })
  })

  test('指标卡键盘 Enter 可导航', async ({ page }) => {
    await login(page)
    const card = page.getByRole('link', { name: /设备总数/ })
    await expect(card).toBeVisible()
    await card.focus()
    await page.keyboard.press('Enter')
    await expect(page).toHaveURL(/\/devices/)
  })

  test('全局 :focus-visible 规则已加载', async ({ page }) => {
    await login(page)
    const loaded = await page.evaluate(() => {
      for (const sheet of document.styleSheets) {
        let rules: CSSRuleList
        try {
          rules = sheet.cssRules
        } catch {
          continue
        }
        for (const rule of Array.from(rules)) {
          if (rule.cssText.includes(':focus-visible') && rule.cssText.includes('outline')) return true
        }
      }
      return false
    })
    expect(loaded).toBe(true)
  })

  test('Dashboard 首屏先展示 loading，数据到达后展示真实内容', async ({ page }) => {
    let release!: () => void
    const gate = new Promise<void>((resolve) => { release = resolve })
    await page.goto('/login')
    await page.route(devicesUrl, async (route) => {
      await gate
      await route.continue()
    })
    await login(page)
    await expect(page.getByTestId('dashboard-loading')).toBeVisible()
    release()
    await expect(page.getByTestId('dashboard-loading')).toHaveCount(0, { timeout: 10000 })
    await expect(page.getByText('设备总数')).toBeVisible()
  })

  test('Dashboard 统计失败显示错误并可重新加载', async ({ page }) => {
    let fail = true
    await page.goto('/login')
    await page.route(devicesUrl, (route) =>
      fail
        ? route.fulfill({ status: 500, json: { code: 500, message: 'boom' } })
        : route.continue(),
    )
    await login(page)
    await expect(page.getByTestId('dashboard-error')).toBeVisible()
    await expect(page.getByText(/部分统计数据加载失败/)).toBeVisible()
    fail = false
    await page.getByRole('button', { name: '重新加载' }).click()
    await expect(page.getByTestId('dashboard-error')).toHaveCount(0, { timeout: 10000 })
  })

  test('告警计数轮询失败时面板显示错误与重试', async ({ page }) => {
    let fail = true
    await page.goto('/login')
    await page.route(unreadUrl, (route) =>
      fail
        ? route.fulfill({ status: 500, json: { code: 500, message: 'boom' } })
        : route.fulfill({ json: ok({ count: 1 }) }),
    )
    await page.route(alertsListUrl, (route) => route.fulfill({ json: ok({ items: [], total: 0 }) }))
    await login(page)
    await page.getByRole('button', { name: '告警通知' }).click()
    await expect(page.getByText('告警加载失败，请重试')).toBeVisible()
    const retry = page.getByRole('button', { name: '重试' })
    await expect(retry).toBeVisible()
    fail = false
    await retry.click()
    await expect(page.getByText('告警加载失败，请重试')).toHaveCount(0, { timeout: 10000 })
  })

  test('深色主指标卡脚注文字计算样式对比度达到 WCAG AA（brand-7 on white ≥ 4.5:1）', async ({ page }) => {
    await page.emulateMedia({ colorScheme: 'light' })
    await login(page)
    const footnote = page.locator('.metric-card--primary .metric-footnote').first()
    await expect(footnote).toBeVisible()

    const contrast = await footnote.evaluate((el) => {
      const card = el.closest('.metric-card--primary') as HTMLElement | null
      if (!card) throw new Error('未找到主指标卡容器')
      const toLinear = (value: string) => {
        const match = value.match(/rgba?\(\s*([\d.]+)[,\s]+([\d.]+)[,\s]+([\d.]+)/)
        if (!match) throw new Error(`无法解析颜色: ${value}`)
        return [match[1], match[2], match[3]].map((part) => {
          const channel = Number(part) / 255
          return channel <= 0.03928 ? channel / 12.92 : Math.pow((channel + 0.055) / 1.055, 2.4)
        })
      }
      const luminance = (value: string) => {
        const [r, g, b] = toLinear(value)
        return 0.2126 * r + 0.7152 * g + 0.0722 * b
      }
      const foreground = luminance(getComputedStyle(el).color)
      const background = luminance(getComputedStyle(card).backgroundColor)
      return (Math.max(foreground, background) + 0.05) / (Math.min(foreground, background) + 0.05)
    })
    expect(contrast).toBeGreaterThanOrEqual(4.5)
  })

  test('空态文字计算色为 secondary token 且不是 placeholder', async ({ page }) => {
    await page.goto('/login')
    await page.route(/\/api\/v1\/tokens(?:\?|$)/, (route) =>
      route.fulfill({ json: { code: 0, message: 'ok', data: [] } }),
    )
    await login(page)
    await page.goto('/settings')
    const empty = page.locator('.empty-state-text').first()
    await expect(empty).toBeVisible()
    const colors = await page.evaluate(() => {
      const el = document.querySelector('.empty-state-text') as HTMLElement
      const probe = (color: string) => {
        const span = document.createElement('span')
        span.style.color = color
        document.body.appendChild(span)
        const value = getComputedStyle(span).color
        span.remove()
        return value
      }
      return {
        actual: getComputedStyle(el).color,
        secondary: probe('var(--td-text-color-secondary)'),
        placeholder: probe('var(--td-text-color-placeholder)'),
      }
    })
    expect(colors.actual).toBe(colors.secondary)
    expect(colors.actual).not.toBe(colors.placeholder)
  })
})

import { test, expect, type Page } from '@playwright/test'
import { login } from './helpers'

// 回归：TDesign 迁移审查发现的 P2 —— OutputOverview 的“查看设备/查看全部”按钮
// 曾因 Dashboard.vue 未绑定 @navigate 而成为死按钮。修复后点击应正确路由跳转。
test.describe('首页任务概览跳转（P2 回归）', () => {
  test('“查看设备”按钮跳转到设备管理页', async ({ page }) => {
    await login(page)
    await expect(page.getByText('本月备份次数')).toBeVisible()
    await page.getByRole('button', { name: '查看设备' }).click()
    await expect(page).toHaveURL(/\/devices$/)
    await expect(page.locator('.t-card__title').filter({ hasText: '设备管理' })).toBeVisible()
  })

  test('“查看全部”按钮跳转到失败设备页', async ({ page }) => {
    await login(page)
    await expect(page.getByText('查看全部')).toBeVisible()
    await page.getByRole('button', { name: '查看全部' }).click()
    await expect(page).toHaveURL(/\/failed$/)
    await expect(page.locator('.t-card__title').filter({ hasText: '失败设备列表' })).toBeVisible()
  })
})

// TDesign 断点：xl>=1200 四列 / sm>=576 两列 / xs 单列；768px 兜底由 global.css 覆盖。
test.describe('仪表盘响应式栅格', () => {
  async function cardBoxes(page: Page, selector: string) {
    return page.locator(selector).evaluateAll((els) =>
      els.map((el) => {
        const rect = el.getBoundingClientRect()
        return { y: rect.y, width: rect.width }
      }),
    )
  }

  test('1920 四卡同排、TOP5 并排、任务与失败概览 8:4；1024 两列；390 单列', async ({ page }) => {
    await login(page)
    await expect(page.getByText('设备总数')).toBeVisible()

    await page.setViewportSize({ width: 1920, height: 1080 })
    const cards = await cardBoxes(page, '.metric-card')
    expect(cards).toHaveLength(4)
    for (const card of cards) {
      expect(Math.abs(card.y - cards[0].y)).toBeLessThan(2)
      expect(Math.abs(card.width - cards[0].width)).toBeLessThan(2)
    }
    const ranks = await cardBoxes(page, '.rank-card')
    expect(ranks).toHaveLength(2)
    expect(Math.abs(ranks[0].y - ranks[1].y)).toBeLessThan(2)
    const overview = await cardBoxes(page, '.overview-card')
    expect(overview).toHaveLength(2)
    expect(Math.abs(overview[0].y - overview[1].y)).toBeLessThan(2)
    const ratio = overview[0].width / overview[1].width
    // starter 主题为 9:3 分栏（约 3:1），非旧自定义的 8:4。
    expect(ratio).toBeGreaterThan(2.5)
    expect(ratio).toBeLessThan(4)

    await page.setViewportSize({ width: 1024, height: 768 })
    const midCards = await cardBoxes(page, '.metric-card')
    expect(new Set(midCards.map((c) => Math.round(c.y))).size).toBe(2)
    expect(midCards.filter((c) => Math.abs(c.y - midCards[0].y) < 2)).toHaveLength(2)

    await page.setViewportSize({ width: 390, height: 844 })
    const smallCards = await cardBoxes(page, '.metric-card')
    expect(new Set(smallCards.map((c) => Math.round(c.y))).size).toBe(4)
  })
})

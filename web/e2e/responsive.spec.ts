import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 响应式回归：移动端窄屏关键页面不溢出、卡片标题不竖排、分段 Tabs 不被裁切。
test.describe('响应式（移动端 390px）', () => {
  test.use({ viewport: { width: 390, height: 844 } })

  test('关键页面无横向溢出', async ({ page }) => {
    await login(page)
    for (const r of ['/dashboard', '/devices', '/alerts', '/settings']) {
      await page.goto(r)
      await page.waitForTimeout(500)
      const ov = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
      expect(ov, `溢出路由 ${r}`).toBeLessThanOrEqual(1)
    }
  })

  test('设备列表卡片标题在窄屏不竖排（不换行成列）', async ({ page }) => {
    await login(page)
    await page.goto('/devices')
    await page.waitForTimeout(600)
    const h = await page.evaluate(() => {
      const t = document.querySelector('.t-card__title') as HTMLElement | null
      return t ? Math.round(t.getBoundingClientRect().height) : 0
    })
    // 竖排时高度会远大于单行（约 24px）
    expect(h).toBeGreaterThan(0)
    expect(h).toBeLessThan(40)
  })

  test('仪表盘分段 Tabs 窄屏被约束宽度且不撑破页面', async ({ page }) => {
    await login(page)
    await page.goto('/dashboard')
    // 等待统计加载完成、标签栏渲染后再测量，避免加载态下取不到节点导致偶发
    await page.waitForSelector('[data-testid="dashboard-loading"]', { state: 'hidden', timeout: 15000 }).catch(() => {})
    await page.locator('.dashboard-tabs .t-tabs__nav').first().waitFor({ state: 'visible', timeout: 15000 }).catch(() => {})
    const m = await page.evaluate(() => {
      const n = document.querySelector('.dashboard-tabs .t-tabs__nav') as HTMLElement | null
      return { cw: n ? n.clientWidth : 0, vw: window.innerWidth }
    })
    // 导航容器宽度不超过视口（不撑破页面），超出部分靠横向滚动可达
    expect(m.cw).toBeGreaterThan(0)
    expect(m.cw).toBeLessThanOrEqual(m.vw)
  })
})

// 回归：平板/中等宽度（1024px）设备列表顶部筛选工具栏——标题不竖排、工具栏不溢出裁切。
test.describe('响应式（中等宽度 1024px）', () => {
  test.use({ viewport: { width: 1024, height: 900 } })

  test('设备列表卡片标题不竖排且筛选工具栏完整可见、无横向溢出', async ({ page }) => {
    await login(page)
    await page.goto('/devices')
    await page.waitForTimeout(700)
    const m = await page.evaluate(() => {
      const title = document.querySelector('.t-card__title') as HTMLElement | null
      const search = document.querySelector('.t-card__actions .t-input__wrap') as HTMLElement | null
      const createBtn = [...document.querySelectorAll('.t-card__actions .t-button')].find((b) =>
        (b.textContent || '').includes('添加设备'),
      ) as HTMLElement | undefined
      const titleRect = title ? title.getBoundingClientRect() : null
      const btnRect = createBtn ? createBtn.getBoundingClientRect() : null
      return {
        titleH: titleRect ? Math.round(titleRect.height) : 0,
        overflow: document.documentElement.scrollWidth - document.documentElement.clientWidth,
        searchVisible: !!search && search.getBoundingClientRect().width > 80,
        createInViewport: !!btnRect && btnRect.right <= window.innerWidth + 1 && btnRect.width > 0,
      }
    })
    // 标题单行（竖排会远大于约 24px）
    expect(m.titleH, '标题高度').toBeGreaterThan(0)
    expect(m.titleH, '标题不竖排').toBeLessThan(40)
    // 无横向溢出；搜索框可见；主按钮完整在视口内（未被裁切）
    expect(m.overflow, '横向溢出').toBeLessThanOrEqual(1)
    expect(m.searchVisible, '搜索框可见').toBeTruthy()
    expect(m.createInViewport, '添加设备按钮完整可见').toBeTruthy()
  })
})

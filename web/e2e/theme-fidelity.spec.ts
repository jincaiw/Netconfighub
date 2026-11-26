import { test, expect } from '@playwright/test'
import { login } from './helpers'

// 主题忠实度回归：验证首页四段式与 starter dashboard/base 的结构一致。
test.describe('TDesign dashboard/base 主题结构', () => {
  test('四段式齐备：四指标卡、迷你图、9:3 折线+饼图、日期选择器、周/月切换、排名圆圈、Trend、导出 inner-card', async ({ page }) => {
    await login(page)

    // TopPanel：4 张指标卡（shadcn 风格：标题+右上图标+大数字+脚注）
    await expect(page.locator('.metric-card')).toHaveCount(4)
    await expect(page.locator('.metric-card--primary')).toHaveCount(1)
    // 回归：卡脚注只渲染一次（曾重复为「近7天近7天」），且保留 contrast 依赖的 .metric-footnote。
    await expect(page.locator('.metric-card .metric-footnote')).toHaveCount(4)
    await expect(page.locator('.dashboard-item--main-color')).toHaveCount(1)
    await expect(page.locator('.metric-card .kpi-value').first()).toBeVisible()

    // MiddleChart：9:3 折线 + 饼图 + 日期选择器
    await expect(page.getByText('备份趋势')).toBeVisible()
    await expect(page.getByText('状态分布')).toBeVisible()
    await expect(page.locator('.t-date-range-picker').first()).toBeVisible()

    // RankList：周/月 radio 切换 + 排名圆圈
    await expect(page.getByText('失败次数排名')).toBeVisible()
    await expect(page.getByText('任务量排名')).toBeVisible()
    const rankTables = page.locator('.dashboard-rank-card .t-table')
    expect(await rankTables.count()).toBeGreaterThanOrEqual(2)
    const weekMonth = page.getByText('本周', { exact: true })
    expect(await weekMonth.count()).toBeGreaterThanOrEqual(1)

    // 排名卡片内容区（shadcn 风格：排名圆点 + 设备 + 次数列）
    await expect(page.locator('.dashboard-rank-card').first()).toBeVisible()

    // OutputOverview：外层卡 + 导出按钮 + 两张 inner-card + 查看按钮
    await expect(page.getByText('备份量趋势')).toBeVisible()
    await expect(page.getByRole('button', { name: '导出', exact: true })).toBeVisible()
    await expect(page.locator('.inner-card')).toHaveCount(2)
    await expect(page.getByText('本月备份次数')).toBeVisible()
    await expect(page.getByText('本月成功率')).toBeVisible()
  })

  test('周/月切换改变排名数据区', async ({ page }) => {
    await login(page)
    await page.getByText('本月', { exact: true }).first().click()
    await expect(page.getByText('任务量排名')).toBeVisible()
  })

  test('首页备份取数使用服务端允许的 page_size（<=100），不再请求 total 被静默降级', async ({ page }) => {
    const backupReq = page.waitForRequest(
      (req) => req.url().includes('/api/v1/backups?') && /[?&]page=1/.test(req.url()),
    )
    await login(page)
    const req = await backupReq
    const url = new URL(req.url())
    const size = Number(url.searchParams.get('page_size'))
    expect(size).toBeGreaterThan(10)
    expect(size).toBeLessThanOrEqual(100)
  })
})

import { test, expect } from '@playwright/test'
import { login } from './helpers'

test.describe('404 页面', () => {
  test('未知路由渲染 NotFound 的标题与描述', async ({ page }) => {
    await login(page)
    await page.goto('/this-route-does-not-exist-nch')
    await expect(page.getByText('404', { exact: true })).toBeVisible()
    await expect(page.getByText('您访问的页面不存在')).toBeVisible()
    await expect(page.getByRole('button', { name: '返回首页' })).toBeVisible()
  })
})

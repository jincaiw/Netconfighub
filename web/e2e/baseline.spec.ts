import { test, expect } from '@playwright/test'
import { login, getAuthToken, createDevice, deleteDevice, createBaseline, deleteBaseline } from './helpers'

test.describe('基线管理流程', () => {
  let token: string
  let deviceId: number
  const createdBaselineIds: number[] = []

  test.beforeAll(async () => {
    token = await getAuthToken()
    const device = await createDevice(token, {
      name: `e2e-baseline-dev-${Date.now()}`,
      ip_address: '10.10.0.1',
      vendor: 'cisco',
      model: 'C9300',
      conn_protocol: 'ssh',
      port: 22,
      username: 'admin',
    })
    deviceId = device.id
  })

  test.afterAll(async () => {
    for (const id of createdBaselineIds) {
      try {
        await deleteBaseline(token, id)
      } catch {
      // cleanup is best-effort
    }
    }
    try {
      await deleteDevice(token, deviceId)
    } catch {
      // cleanup is best-effort
    }
  })

  test('打开基线管理页', async ({ page }) => {
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('.t-card__title').filter({ hasText: '基线列表' })).toBeVisible()
    await expect(page.locator('table')).toBeVisible()
  })

  test('创建设备级基线', async ({ page }) => {
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('table')).toBeVisible()
    await page.getByRole('button', { name: '创建基线' }).click()
    await expect(page.locator('.t-dialog:visible')).toBeVisible()
    await page.locator('.t-radio').filter({ hasText: '设备' }).click()
    await page.locator('.t-form__item').filter({ hasText: '选择设备' }).locator('.t-select').click()
    await expect(page.locator('.t-select-option').first()).toBeVisible()
    await page.locator('.t-select-option').first().click()
    await page.locator('.t-form__item').filter({ hasText: '基线内容' }).locator('textarea').fill('hostname e2e-test\ninterface Gig0/1\n ip address 10.0.0.1')
    await page.locator('.t-dialog:visible').getByRole('button', { name: '确定' }).click()
    await expect(page.locator('.t-dialog:visible')).toBeHidden()
  })

  test('查看基线列表', async ({ page }) => {
    const baseline = await createBaseline(token, {
      scope: 'device',
      device_id: deviceId,
      content: 'hostname baseline-list-test\ninterface Gig0/1',
    })
    if (baseline?.id) createdBaselineIds.push(baseline.id)
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('.t-card__title').filter({ hasText: '基线列表' })).toBeVisible()
    await expect(page.locator('table')).toBeVisible()
  })

  test('编辑基线内容', async ({ page }) => {
    const baseline = await createBaseline(token, {
      scope: 'device',
      device_id: deviceId,
      content: 'hostname before-edit',
    })
    if (baseline?.id) createdBaselineIds.push(baseline.id)
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('table')).toBeVisible()
    const editBtn = page.locator('table').getByRole('button', { name: '编辑' }).first()
    await editBtn.click()
    await expect(page.locator('.t-dialog:visible')).toBeVisible()
    const textarea = page.locator('.t-form__item').filter({ hasText: '基线内容' }).locator('textarea')
    await textarea.clear()
    await textarea.fill('hostname after-edit\ninterface Gig0/1\n ip address 10.0.0.1')
    await page.locator('.t-dialog:visible').getByRole('button', { name: '确定' }).click()
    await expect(page.locator('.t-dialog:visible')).toBeHidden()
  })

  test('删除基线', async ({ page }) => {
    await createBaseline(token, {
      scope: 'device',
      device_id: deviceId,
      content: 'hostname to-delete',
    })
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('table')).toBeVisible()
    const deleteResponse = page.waitForResponse(
      (resp) => resp.url().includes('/baselines/') && resp.request().method() === 'DELETE',
      { timeout: 10000 },
    )
    const deleteBtn = page.locator('table').getByRole('button', { name: '删除' }).first()
    await deleteBtn.click()
    await page.getByRole('button', { name: '删除' }).last().click()
    await deleteResponse
  })

  test('页面头部仅创建基线为实心主按钮，导出报告为次级', async ({ page }) => {
    await login(page)
    await page.goto('/baselines')
    const background = async (name: string) =>
      page.getByRole('button', { name, exact: true }).first().evaluate((el) => getComputedStyle(el).backgroundColor)
    const primary = await background('创建基线')
    const secondary = await background('导出报告')
    expect(primary).not.toBe(secondary)
    expect(primary).not.toBe('rgba(0, 0, 0, 0)')
  })
})

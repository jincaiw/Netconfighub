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
      } catch {}
    }
    try {
      await deleteDevice(token, deviceId)
    } catch {}
  })

  test('打开基线管理页', async ({ page }) => {
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('.n-card-header__main').filter({ hasText: '基线列表' })).toBeVisible()
    await expect(page.locator('table')).toBeVisible()
  })

  test('创建设备级基线', async ({ page }) => {
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('table')).toBeVisible()
    await page.getByRole('button', { name: '创建基线' }).click()
    await expect(page.locator('.n-dialog')).toBeVisible()
    await page.locator('.n-radio').filter({ hasText: '设备' }).click()
    await page.locator('.n-form-item').filter({ hasText: '选择设备' }).locator('.n-base-selection').click()
    await page.waitForTimeout(500)
    await page.locator('.n-base-select-option').first().click()
    await page.locator('.n-form-item').filter({ hasText: '基线内容' }).locator('textarea').fill('hostname e2e-test\ninterface Gig0/1\n ip address 10.0.0.1')
    await page.locator('.n-dialog').getByRole('button', { name: '确定' }).click()
    await page.waitForTimeout(2000)
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
    await expect(page.locator('.n-card-header__main').filter({ hasText: '基线列表' })).toBeVisible()
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
    await expect(page.locator('.n-dialog')).toBeVisible()
    const textarea = page.locator('.n-form-item').filter({ hasText: '基线内容' }).locator('textarea')
    await textarea.clear()
    await textarea.fill('hostname after-edit\ninterface Gig0/1\n ip address 10.0.0.1')
    await page.locator('.n-dialog').getByRole('button', { name: '确定' }).click()
    await page.waitForTimeout(2000)
  })

  test('删除基线', async ({ page }) => {
    const baseline = await createBaseline(token, {
      scope: 'device',
      device_id: deviceId,
      content: 'hostname to-delete',
    })
    await login(page)
    await page.goto('/baselines')
    await expect(page.locator('table')).toBeVisible()
    const deleteBtn = page.locator('table').getByRole('button', { name: '删除' }).first()
    await deleteBtn.click()
    await page.getByRole('button', { name: '删除' }).last().click()
    await page.waitForTimeout(2000)
  })
})

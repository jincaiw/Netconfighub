import { test, expect } from '@playwright/test'
import { login, getAuthToken, createDevice, deleteDevice } from './helpers'

test.describe('设备管理流程', () => {
  let token: string
  const createdDeviceIds: number[] = []

  test.beforeAll(async () => {
    token = await getAuthToken()
  })

  test.afterAll(async () => {
    for (const id of createdDeviceIds) {
      try {
        await deleteDevice(token, id)
      } catch {}
    }
  })

  test('打开设备列表页，显示设备表格', async ({ page }) => {
    await login(page)
    await page.goto('/devices')
    await expect(page.locator('.n-card-header__main').filter({ hasText: '设备管理' })).toBeVisible()
    await expect(page.locator('table')).toBeVisible()
  })

  test('点击添加设备按钮，弹出表单', async ({ page }) => {
    await login(page)
    await page.goto('/devices')
    await expect(page.locator('table')).toBeVisible()
    await page.getByRole('button', { name: '添加设备' }).click()
    await expect(page.locator('.n-dialog')).toBeVisible()
    await expect(page.locator('.n-dialog').getByText('设备名称', { exact: true })).toBeVisible()
    await expect(page.locator('.n-dialog').getByText('IP 地址', { exact: true })).toBeVisible()
  })

  test('填写设备信息并提交，设备出现在列表中', async ({ page }) => {
    const deviceName = `e2e-device-${Date.now()}`
    await login(page)
    await page.goto('/devices')
    await expect(page.locator('table')).toBeVisible()
    await page.getByRole('button', { name: '添加设备' }).click()
    await expect(page.locator('.n-dialog')).toBeVisible()
    await page.getByPlaceholder('请输入设备名称').fill(deviceName)
    await page.getByPlaceholder('请输入 IP 地址').fill('192.168.1.100')
    await page.locator('.n-form-item').filter({ hasText: '厂商' }).locator('.n-base-selection').click()
    await page.locator('.n-base-select-option').filter({ hasText: '思科' }).click()
    await page.getByPlaceholder('请输入型号').fill('C9300')
    await page.locator('.n-form-item').filter({ hasText: '连接协议' }).locator('.n-base-selection').click()
    await page.locator('.n-base-select-option').filter({ hasText: 'SSH' }).click()
    await page.locator('.n-form-item').filter({ hasText: '端口' }).locator('input').fill('22')
    await page.getByPlaceholder('请输入用户名').fill('admin')
    await page.getByPlaceholder('请输入密码').fill('testpass123')
    await page.locator('.n-dialog').getByRole('button', { name: '确定' }).click()
    await page.waitForTimeout(2000)
    await page.reload()
    await expect(page.locator('table')).toBeVisible()
    await expect(page.getByText(deviceName)).toBeVisible({ timeout: 5000 })
  })

  test('点击设备名称跳转到详情页', async ({ page }) => {
    const device = await createDevice(token, {
      name: `e2e-detail-${Date.now()}`,
      ip_address: '10.0.0.2',
      vendor: 'cisco',
      model: 'C9300',
      conn_protocol: 'ssh',
      port: 22,
      username: 'admin',
    })
    if (device?.id) createdDeviceIds.push(device.id)
    await login(page)
    await page.goto('/devices')
    await expect(page.locator('table')).toBeVisible()
    const detailBtn = page.locator('table').getByRole('button', { name: '详情' }).first()
    await detailBtn.click()
    await expect(page).toHaveURL(/\/devices\/\d+/, { timeout: 5000 })
  })

  test('在详情页查看设备信息和配置 Tab', async ({ page }) => {
    const device = await createDevice(token, {
      name: `e2e-info-${Date.now()}`,
      ip_address: '10.0.0.3',
      vendor: 'cisco',
      model: 'C9300',
      conn_protocol: 'ssh',
      port: 22,
      username: 'admin',
    })
    if (device?.id) createdDeviceIds.push(device.id)
    await login(page)
    await page.goto(`/devices/${device.id}`)
    await expect(page.getByText('设备信息')).toBeVisible()
    await expect(page.locator('.n-tabs-tab').filter({ hasText: '最新配置' })).toBeVisible()
    await expect(page.locator('.n-tabs-tab').filter({ hasText: '历史版本' })).toBeVisible()
    await expect(page.locator('.n-tabs-tab').filter({ hasText: '基线偏差' })).toBeVisible()
  })

  test('点击触发备份按钮', async ({ page }) => {
    const device = await createDevice(token, {
      name: `e2e-backup-${Date.now()}`,
      ip_address: '10.0.0.4',
      vendor: 'cisco',
      model: 'C9300',
      conn_protocol: 'ssh',
      port: 22,
      username: 'admin',
    })
    if (device?.id) createdDeviceIds.push(device.id)
    await login(page)
    await page.goto(`/devices/${device.id}`)
    await expect(page.getByText('设备信息')).toBeVisible()
    await page.getByRole('button', { name: '触发备份' }).click()
    await page.waitForTimeout(2000)
  })

  test('删除设备', async ({ page }) => {
    const device = await createDevice(token, {
      name: `e2e-delete-${Date.now()}`,
      ip_address: '10.0.0.5',
      vendor: 'cisco',
      model: 'C9300',
      conn_protocol: 'ssh',
      port: 22,
      username: 'admin',
    })
    await login(page)
    await page.goto('/devices')
    await expect(page.locator('table')).toBeVisible()
    const row = page.locator('tr').filter({ hasText: device.name })
    await row.getByRole('button', { name: '删除' }).click()
    await page.getByRole('button', { name: '删除' }).last().click()
    await page.waitForTimeout(2000)
  })

  test('编辑设备功能', async ({ page }) => {
    const token = await getAuthToken()
    const createRes = await fetch('http://localhost:8080/api/v1/devices', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        name: 'edit-test-device',
        ip_address: '10.99.99.99',
        vendor: 'cisco',
        model: 'c2960',
        conn_protocol: 'ssh',
        port: 22,
        username: 'admin',
        password: 'test123',
      }),
    })
    const createData = await createRes.json()
    const deviceId = createData.data.id

    await login(page)
    await page.goto('/devices')
    await page.waitForTimeout(1000)

    const editBtn = page.getByRole('button', { name: '编辑' }).first()
    await expect(editBtn).toBeVisible({ timeout: 10000 })
    await editBtn.click()
    await page.waitForTimeout(500)
    await expect(page.getByText('编辑设备')).toBeVisible()

    const nameInput = page.locator('.n-dialog').locator('input').first()
    await nameInput.clear()
    await nameInput.fill('edited-device')

    await page.getByRole('button', { name: '保存' }).click()
    await page.waitForTimeout(1000)

    await fetch(`http://localhost:8080/api/v1/devices/${deviceId}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` },
    })
  })
})

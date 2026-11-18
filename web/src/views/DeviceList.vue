<template>
  <div class="device-list">
    <t-card title="设备管理">
      <template #actions>
        <DeviceFilterBar
          v-model:search-keyword="searchKeyword"
          v-model:selected-group="selectedGroup"
          :group-options="groupOptions"
          :export-loading="exportLoading"
          @search="handleSearch"
          @group-change="handleGroupChange"
          @export="handleExport"
          @import="showImportModal = true"
          @create="showCreateModal = true"
        />
      </template>

      <DeviceTable
        :devices="devices"
        :loading="loading"
        :vendor-label-map="vendorLabelMap"
        :model-label-map="modelLabelMap"
        @toggle="handleToggleEnabled"
        @detail="(id) => router.push(`/devices/${id}`)"
        @edit="handleEdit"
        @backup="handleTriggerBackup"
        @delete="handleDelete"
      />
    </t-card>

    <t-dialog v-model:visible="showCreateModal" header="添加设备" :width="560">
      <t-form ref="createFormRef" :data="newDevice" :rules="createRules" label-align="left" label-width="90px">
        <t-form-item label="设备名称" name="name">
          <t-input v-model="newDevice.name" placeholder="请输入设备名称" />
        </t-form-item>
        <t-form-item label="IP 地址" name="ip">
          <t-input v-model="newDevice.ip" placeholder="请输入 IP 地址" />
        </t-form-item>
        <t-form-item label="厂商" name="vendor">
          <t-select
            v-model="newDevice.vendor"
            :options="vendorOptions"
            placeholder="请选择厂商"
            @change="handleVendorChange"
          />
        </t-form-item>
        <t-form-item label="型号" name="model">
          <t-input v-model="newDevice.model" placeholder="请输入型号" />
        </t-form-item>
        <t-form-item label="连接协议" name="protocol">
          <t-select
            v-model="newDevice.protocol"
            :options="protocolOptions"
            placeholder="请选择协议"
          />
        </t-form-item>
        <t-form-item label="采集方式" name="collection_mode">
          <t-select v-model="newDevice.collection_mode" :options="collectionModeOptions" />
          <div v-if="newDevice.collection_mode === 'sftp_push'" class="collection-hint">
            该设备不会被平台拉取，需在设备详情中创建 SFTP 上传账号后由设备主动上传配置。
          </div>
        </t-form-item>
        <t-form-item label="端口" name="port">
          <t-input-number v-model="newDevice.port" :min="1" :max="65535" style="width: 100%" />
        </t-form-item>
        <t-form-item label="用户名" name="username">
          <t-input v-model="newDevice.username" placeholder="请输入用户名" />
        </t-form-item>
        <t-form-item label="密码">
          <t-input v-model="newDevice.password" type="password" placeholder="请输入密码" />
        </t-form-item>
        <t-form-item label="SSH 私钥">
          <t-textarea v-model="newDevice.ssh_key" :autosize="{ minRows: 3 }" placeholder="可选，粘贴 SSH 私钥" />
        </t-form-item>
        <t-form-item label="特权模式密码">
          <t-input v-model="newDevice.enable_password" type="password" placeholder="可选，特权模式密码" />
        </t-form-item>
        <t-form-item label="备份间隔">
          <t-input v-model="newDevice.backup_interval" placeholder="如 30m, 1h, 24h，留空使用默认" />
        </t-form-item>
        <t-form-item label="分组">
          <t-select
            v-model="newDevice.group_id"
            :options="groupOptions"
            placeholder="请选择分组"
            clearable
          />
        </t-form-item>
      </t-form>
      <template #footer>
        <t-button variant="outline" @click="handleCancelCreate">取消</t-button>
        <t-button theme="primary" :loading="creating" @click="handleCreate">确定</t-button>
      </template>
    </t-dialog>

    <t-dialog v-model:visible="showEditModal" header="编辑设备" :width="560">
      <t-form ref="editFormRef" :data="editingDevice" :rules="createRules" label-align="left" label-width="90px">
        <t-form-item label="设备名称" name="name">
          <t-input v-model="editingDevice.name" placeholder="请输入设备名称" />
        </t-form-item>
        <t-form-item label="IP 地址" name="ip">
          <t-input v-model="editingDevice.ip" placeholder="请输入 IP 地址" />
        </t-form-item>
        <t-form-item label="厂商" name="vendor">
          <t-select
            v-model="editingDevice.vendor"
            :options="vendorOptions"
            placeholder="请选择厂商"
            @change="handleEditVendorChange"
          />
        </t-form-item>
        <t-form-item label="型号" name="model">
          <t-input v-model="editingDevice.model" placeholder="请输入型号" />
        </t-form-item>
        <t-form-item label="连接协议" name="protocol">
          <t-select
            v-model="editingDevice.protocol"
            :options="protocolOptions"
            placeholder="请选择协议"
          />
        </t-form-item>
        <t-form-item label="采集方式" name="collection_mode">
          <t-select v-model="editingDevice.collection_mode" :options="collectionModeOptions" />
          <div v-if="editingDevice.collection_mode === 'sftp_push'" class="collection-hint">
            该设备不会被平台拉取，需由设备通过 SFTP 主动上传配置。
          </div>
        </t-form-item>
        <t-form-item label="端口" name="port">
          <t-input-number v-model="editingDevice.port" :min="1" :max="65535" style="width: 100%" />
        </t-form-item>
        <t-form-item label="用户名" name="username">
          <t-input v-model="editingDevice.username" placeholder="请输入用户名" />
        </t-form-item>
        <t-form-item label="密码">
          <t-input v-model="editingDevice.password" type="password" placeholder="留空则不修改" />
        </t-form-item>
        <t-form-item label="SSH 私钥">
          <t-textarea v-model="editingDevice.ssh_key" :autosize="{ minRows: 3 }" placeholder="留空则不修改" />
        </t-form-item>
        <t-form-item label="特权模式密码">
          <t-input v-model="editingDevice.enable_password" type="password" placeholder="留空则不修改" />
        </t-form-item>
        <t-form-item label="备份间隔">
          <t-input v-model="editingDevice.backup_interval" placeholder="如 30m, 1h, 24h，留空使用默认" />
        </t-form-item>
        <t-form-item label="分组">
          <t-select
            v-model="editingDevice.group_id"
            :options="groupOptions"
            placeholder="请选择分组"
            clearable
          />
        </t-form-item>
      </t-form>
      <template #footer>
        <t-button variant="outline" @click="handleCancelEdit">取消</t-button>
        <t-button theme="primary" :loading="updating" @click="handleUpdate">保存</t-button>
      </template>
    </t-dialog>

    <DeviceImportExport
      v-model:visible="showImportModal"
      v-model:result-visible="showImportResult"
      :import-result="importResult"
      :import-loading="importLoading"
      :has-file="!!importFile"
      @file-change="handleFileChange"
      @import="handleImport"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  DialogPlugin,
  MessagePlugin,
  type FormInstanceFunctions,
  type FormRules,
  type UploadFile,
} from 'tdesign-vue-next'
import DeviceFilterBar from '@/components/device/DeviceFilterBar.vue'
import DeviceTable from '@/components/device/DeviceTable.vue'
import DeviceImportExport from '@/components/device/DeviceImportExport.vue'
import {
  getDevices,
  getGroups,
  createDevice,
  updateDevice,
  deleteDevice,
  enableDevice,
  disableDevice,
  importDevices,
  exportDevices,
} from '@/api/device'
import { triggerDeviceBackup } from '@/api/backup'
import type { Device, Group, DeviceCreateRequest, ImportResult } from '@/types'

const router = useRouter()

const loading = ref(false)
const creating = ref(false)
const exportLoading = ref(false)
const importLoading = ref(false)
const devices = ref<Device[]>([])
const searchKeyword = ref('')
const selectedGroup = ref<number | null>(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showImportModal = ref(false)
const showImportResult = ref(false)
const createFormRef = ref<FormInstanceFunctions | null>(null)
const editFormRef = ref<FormInstanceFunctions | null>(null)
const importFile = ref<File | null>(null)
const importResult = ref<ImportResult | null>(null)
const updating = ref(false)

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
})

const vendorOptions = [
  { label: '思科', value: 'cisco' },
  { label: '华三', value: 'h3c' },
  { label: '华为', value: 'huawei' },
  { label: '锐捷', value: 'ruijie' },
]

const protocolOptions = [
  { label: 'SSH', value: 'ssh' },
  { label: 'Telnet', value: 'telnet' },
]

// 采集方式决定配置如何进入平台：平台拉取，或设备通过 SFTP 主动上传。
const collectionModeOptions = [
  { label: '平台拉取（SSH/Telnet）', value: 'pull' },
  { label: 'SFTP 主动上传', value: 'sftp_push' },
]

const groupOptions = ref<Array<{ label: string; value: number }>>([])

const vendorModelMap: Record<string, DeviceCreateRequest['model']> = {
  cisco: 'ios',
  huawei: 'vrp',
  h3c: 'comware',
  ruijie: 'rg-os',
}

const newDevice = reactive<DeviceCreateRequest>({
  name: '',
  ip: '',
  vendor: 'cisco',
  model: 'ios',
  protocol: 'ssh',
  port: 22,
  username: '',
  password: '',
  ssh_key: '',
  enable_password: '',
  backup_interval: '',
  group_id: undefined,
  collection_mode: 'pull',
})

const editingDevice = reactive<DeviceCreateRequest & { id: number }>({
  id: 0,
  name: '',
  ip: '',
  vendor: 'cisco',
  model: 'ios',
  protocol: 'ssh',
  port: 22,
  username: '',
  password: '',
  ssh_key: '',
  enable_password: '',
  backup_interval: '',
  group_id: undefined,
  collection_mode: 'pull',
})

const createRules: FormRules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  ip: [{ required: true, message: '请输入 IP 地址', trigger: 'blur' }],
  vendor: [{ required: true, message: '请选择厂商', trigger: 'change' }],
  model: [{ required: true, message: '请选择型号', trigger: 'change' }],
  protocol: [{ required: true, message: '请选择连接协议', trigger: 'change' }],
  port: [{ required: true, message: '请输入端口', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
}

const vendorLabelMap: Record<string, string> = {
  cisco: '思科',
  h3c: '华三',
  huawei: '华为',
  ruijie: '锐捷',
}

const modelLabelMap: Record<string, string> = {
  ios: 'IOS',
  vrp: 'VRP',
  comware: 'Comware',
  'rg-os': 'RG-OS',
}

function handleVendorChange(vendor: string) {
  newDevice.model = vendorModelMap[vendor] || 'ios'
}

function handleEditVendorChange(vendor: string) {
  editingDevice.model = vendorModelMap[vendor] || 'ios'
}

async function handleToggleEnabled(device: Device) {
  try {
    if (device.enabled) {
      await disableDevice(device.id)
      MessagePlugin.success('设备已禁用')
    } else {
      await enableDevice(device.id)
      MessagePlugin.success('设备已启用')
    }
    fetchDevices()
  } catch {
    MessagePlugin.error(device.enabled ? '禁用设备失败' : '启用设备失败')
  }
}

async function fetchDevices() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: paginationReactive.page,
      page_size: paginationReactive.pageSize,
    }
    if (selectedGroup.value) {
      params.group_id = selectedGroup.value
    }
    if (searchKeyword.value) {
      params.search = searchKeyword.value
    }
    const res = await getDevices(params)
    devices.value = res.data.items
    paginationReactive.itemCount = res.data.total
  } catch {
    MessagePlugin.error('获取设备列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchGroups() {
  try {
    const res = await getGroups({ page: 1, page_size: 100 })
    groupOptions.value = res.data.items.map((g: Group) => ({ label: g.name, value: g.id }))
  } catch {
    MessagePlugin.error('获取分组列表失败')
  }
}

function handleSearch() {
  paginationReactive.page = 1
  fetchDevices()
}

function handleGroupChange() {
  paginationReactive.page = 1
  fetchDevices()
}

async function handleTriggerBackup(deviceId: number) {
  try {
    await triggerDeviceBackup(deviceId)
    MessagePlugin.success('备份任务已触发')
  } catch {
    MessagePlugin.error('备份任务触发失败')
  }
}

function handleDelete(device: Device) {
  DialogPlugin.confirm({
    header: '确认删除',
    body: `确定要删除设备 "${device.name}" 吗？`,
    theme: 'warning',
    confirmBtn: '删除',
    cancelBtn: '取消',
    onConfirm: async () => {
      try {
        await deleteDevice(device.id)
        MessagePlugin.success('删除成功')
        fetchDevices()
      } catch {
        MessagePlugin.error('删除失败')
      }
    },
  })
}

function handleEdit(device: Device) {
  editingDevice.id = device.id
  editingDevice.name = device.name
  editingDevice.ip = device.ip
  editingDevice.vendor = device.vendor
  editingDevice.model = device.model
  editingDevice.protocol = device.protocol
  editingDevice.port = device.port
  editingDevice.username = device.username
  editingDevice.password = ''
  editingDevice.ssh_key = ''
  editingDevice.enable_password = ''
  editingDevice.backup_interval = device.backup_interval || ''
  editingDevice.group_id = device.group_id
  editingDevice.collection_mode = device.collection_mode || 'pull'
  showEditModal.value = true
}

async function handleUpdate() {
  const valid = await editFormRef.value?.validate()
  if (valid !== true) return

  updating.value = true
  try {
    const data: Partial<DeviceCreateRequest> = {
      name: editingDevice.name,
      ip: editingDevice.ip,
      vendor: editingDevice.vendor,
      model: editingDevice.model,
      protocol: editingDevice.protocol,
      port: editingDevice.port,
      username: editingDevice.username,
      backup_interval: editingDevice.backup_interval || undefined,
      group_id: editingDevice.group_id,
      collection_mode: editingDevice.collection_mode,
    }
    if (editingDevice.password) {
      data.password = editingDevice.password
    }
    if (editingDevice.ssh_key) {
      data.ssh_key = editingDevice.ssh_key
    }
    if (editingDevice.enable_password) {
      data.enable_password = editingDevice.enable_password
    }
    await updateDevice(editingDevice.id, data)
    MessagePlugin.success('设备更新成功')
    showEditModal.value = false
    fetchDevices()
  } catch {
    MessagePlugin.error('设备更新失败')
  } finally {
    updating.value = false
  }
}

async function handleCreate() {
  const valid = await createFormRef.value?.validate()
  if (valid !== true) return

  creating.value = true
  try {
    await createDevice(newDevice)
    MessagePlugin.success('设备添加成功')
    showCreateModal.value = false
    resetNewDevice()
    fetchDevices()
  } catch {
    MessagePlugin.error('设备添加失败')
  } finally {
    creating.value = false
  }
}

function handleFileChange(files: UploadFile[]) {
  const raw = files[0]?.raw as File | undefined
  importFile.value = raw || null
}

async function handleImport() {
  if (!importFile.value) return
  importLoading.value = true
  try {
    const res = await importDevices(importFile.value)
    importResult.value = res.data
    showImportModal.value = false
    showImportResult.value = true
    importFile.value = null
    fetchDevices()
  } catch {
    MessagePlugin.error('导入设备失败')
  } finally {
    importLoading.value = false
  }
}

async function handleExport() {
  exportLoading.value = true
  try {
    const params: Record<string, any> = {}
    if (searchKeyword.value) params.search = searchKeyword.value
    if (selectedGroup.value) params.group_id = selectedGroup.value
    const blob = await exportDevices(params)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'devices.csv'
    a.click()
    window.URL.revokeObjectURL(url)
    MessagePlugin.success('导出成功')
  } catch {
    MessagePlugin.error('导出设备失败')
  } finally {
    exportLoading.value = false
  }
}

function resetNewDevice() {
  newDevice.name = ''
  newDevice.ip = ''
  newDevice.vendor = 'cisco'
  newDevice.model = 'ios'
  newDevice.protocol = 'ssh'
  newDevice.port = 22
  newDevice.username = ''
  newDevice.password = ''
  newDevice.ssh_key = ''
  newDevice.enable_password = ''
  newDevice.backup_interval = ''
  newDevice.group_id = undefined
}

function handleCancelCreate() {
  showCreateModal.value = false
  resetNewDevice()
}

function handleCancelEdit() {
  showEditModal.value = false
}

onMounted(() => {
  fetchDevices()
  fetchGroups()
})
</script>

<style scoped>
.collection-hint {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--td-text-color-secondary, #888);
}

.device-list {
  width: 100%;
}

.device-table-wrap {
  width: 100%;
  overflow-x: auto;
}

.device-table {
  width: 100%;
  min-width: 1080px;
  border-collapse: collapse;
  font-size: 14px;
}

.device-table th,
.device-table td {
  border-bottom: 1px solid var(--td-component-stroke);
  padding: 12px;
  text-align: left;
  vertical-align: middle;
  white-space: nowrap;
}

.device-table th {
  color: var(--td-text-color-primary);
  font-weight: 600;
  background: var(--td-bg-color-container-hover);
}

.device-table tr:hover td {
  background: var(--td-bg-color-container-hover);
}

.table-empty {
  color: var(--td-text-color-placeholder);
  text-align: center;
  padding: 32px 12px;
}

.import-hint {
  margin-top: 12px;
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}

.result-success {
  color: var(--td-success-color-5);
}

.result-error {
  color: var(--td-error-color-5);
}
</style>

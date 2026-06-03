<template>
  <div class="device-list">
    <n-card title="设备管理">
      <template #header-extra>
        <n-space>
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索设备名称/IP"
            clearable
            style="width: 180px"
            @clear="handleSearch"
            @keydown.enter="handleSearch"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" />
            </template>
          </n-input>
          <n-select
            v-model:value="selectedGroup"
            :options="groupOptions"
            placeholder="按分组筛选"
            clearable
            style="width: 160px"
            @update:value="handleGroupChange"
          />
          <n-button @click="handleExport" :loading="exportLoading">导出</n-button>
          <n-button @click="showImportModal = true">导入</n-button>
          <n-button type="primary" @click="showCreateModal = true">
            <template #icon>
              <n-icon :component="AddOutline" />
            </template>
            添加设备
          </n-button>
        </n-space>
      </template>

      <div class="device-table-wrap">
        <table class="device-table" aria-label="设备列表">
          <thead>
            <tr>
              <th>名称</th>
              <th>IP 地址</th>
              <th>厂商</th>
              <th>型号</th>
              <th>分组</th>
              <th>启用状态</th>
              <th>备份状态</th>
              <th>最近备份</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="9" class="table-empty">加载中...</td>
            </tr>
            <tr v-else-if="devices.length === 0">
              <td colspan="9" class="table-empty">No Data</td>
            </tr>
            <tr v-for="device in devices" v-else :key="device.id">
              <td>{{ device.name }}</td>
              <td>{{ device.ip }}</td>
              <td>{{ vendorLabelMap[device.vendor] || device.vendor }}</td>
              <td>{{ modelLabelMap[device.model] || device.model }}</td>
              <td>{{ device.group_name || '-' }}</td>
              <td>
                <n-switch :value="device.enabled" @update:value="() => handleToggleEnabled(device)" />
              </td>
              <td>
                <backup-status :status="device.last_backup_status" />
              </td>
              <td>{{ device.last_backup_at || '-' }}</td>
              <td>
                <n-space size="small">
                  <n-button size="small" @click="router.push(`/devices/${device.id}`)">详情</n-button>
                  <n-button size="small" @click="handleEdit(device)">编辑</n-button>
                  <n-button size="small" type="primary" @click="handleTriggerBackup(device.id)">备份</n-button>
                  <n-button size="small" :type="device.enabled ? 'warning' : 'success'" @click="handleToggleEnabled(device)">
                    {{ device.enabled ? '禁用' : '启用' }}
                  </n-button>
                  <n-button size="small" type="error" @click="handleDelete(device)">删除</n-button>
                </n-space>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </n-card>

    <n-modal v-model:show="showCreateModal" preset="dialog" title="添加设备" style="width: 560px">
      <n-form ref="createFormRef" :model="newDevice" :rules="createRules" label-placement="left" label-width="90">
        <n-form-item path="name" label="设备名称">
          <n-input v-model:value="newDevice.name" placeholder="请输入设备名称" />
        </n-form-item>
        <n-form-item path="ip" label="IP 地址">
          <n-input v-model:value="newDevice.ip" placeholder="请输入 IP 地址" />
        </n-form-item>
        <n-form-item path="vendor" label="厂商">
          <n-select
            v-model:value="newDevice.vendor"
            :options="vendorOptions"
            placeholder="请选择厂商"
            @update:value="handleVendorChange"
          />
        </n-form-item>
        <n-form-item path="model" label="型号">
          <n-input v-model:value="newDevice.model" placeholder="请输入型号" />
        </n-form-item>
        <n-form-item path="protocol" label="连接协议">
          <n-select
            v-model:value="newDevice.protocol"
            :options="protocolOptions"
            placeholder="请选择协议"
          />
        </n-form-item>
        <n-form-item path="port" label="端口">
          <n-input-number v-model:value="newDevice.port" :min="1" :max="65535" style="width: 100%" />
        </n-form-item>
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="newDevice.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="密码">
          <n-input v-model:value="newDevice.password" type="password" show-password-on="click" placeholder="请输入密码" />
        </n-form-item>
        <n-form-item label="SSH Key">
          <n-input v-model:value="newDevice.ssh_key" type="textarea" placeholder="可选，粘贴 SSH 私钥" :rows="3" />
        </n-form-item>
        <n-form-item label="Enable 密码">
          <n-input v-model:value="newDevice.enable_password" type="password" show-password-on="click" placeholder="可选，特权模式密码" />
        </n-form-item>
        <n-form-item label="备份间隔">
          <n-input v-model:value="newDevice.backup_interval" placeholder="如 30m, 1h, 24h，留空使用默认" />
        </n-form-item>
        <n-form-item label="分组">
          <n-select
            v-model:value="newDevice.group_id"
            :options="groupOptions"
            placeholder="请选择分组"
            clearable
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="handleCancelCreate">取消</n-button>
        <n-button type="primary" :loading="creating" @click="handleCreate">确定</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showEditModal" preset="dialog" title="编辑设备" style="width: 560px">
      <n-form ref="editFormRef" :model="editingDevice" :rules="createRules" label-placement="left" label-width="90">
        <n-form-item path="name" label="设备名称">
          <n-input v-model:value="editingDevice.name" placeholder="请输入设备名称" />
        </n-form-item>
        <n-form-item path="ip" label="IP 地址">
          <n-input v-model:value="editingDevice.ip" placeholder="请输入 IP 地址" />
        </n-form-item>
        <n-form-item path="vendor" label="厂商">
          <n-select v-model:value="editingDevice.vendor" :options="vendorOptions" placeholder="请选择厂商" @update:value="handleEditVendorChange" />
        </n-form-item>
        <n-form-item path="model" label="型号">
          <n-input v-model:value="editingDevice.model" placeholder="请输入型号" />
        </n-form-item>
        <n-form-item path="protocol" label="连接协议">
          <n-select v-model:value="editingDevice.protocol" :options="protocolOptions" placeholder="请选择协议" />
        </n-form-item>
        <n-form-item path="port" label="端口">
          <n-input-number v-model:value="editingDevice.port" :min="1" :max="65535" style="width: 100%" />
        </n-form-item>
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="editingDevice.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="密码">
          <n-input v-model:value="editingDevice.password" type="password" show-password-on="click" placeholder="留空则不修改" />
        </n-form-item>
        <n-form-item label="SSH Key">
          <n-input v-model:value="editingDevice.ssh_key" type="textarea" placeholder="留空则不修改" :rows="3" />
        </n-form-item>
        <n-form-item label="Enable 密码">
          <n-input v-model:value="editingDevice.enable_password" type="password" show-password-on="click" placeholder="留空则不修改" />
        </n-form-item>
        <n-form-item label="备份间隔">
          <n-input v-model:value="editingDevice.backup_interval" placeholder="如 30m, 1h, 24h，留空使用默认" />
        </n-form-item>
        <n-form-item label="分组">
          <n-select v-model:value="editingDevice.group_id" :options="groupOptions" placeholder="请选择分组" clearable />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="handleCancelEdit">取消</n-button>
        <n-button type="primary" :loading="updating" @click="handleUpdate">保存</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showImportModal" preset="dialog" title="导入设备" style="width: 480px">
      <n-upload
        :max="1"
        accept=".csv"
        :default-upload="false"
        @change="handleFileChange"
      >
        <n-upload-dragger>
          <n-text style="font-size: 16px">点击或拖拽 CSV 文件到此区域</n-text>
          <n-p depth="3" style="margin: 8px 0 0 0">
            CSV 格式: name, ip, vendor, model, protocol, port, username, password, group_name
          </n-p>
        </n-upload-dragger>
      </n-upload>
      <template #action>
        <n-button @click="showImportModal = false">取消</n-button>
        <n-button type="primary" :loading="importLoading" :disabled="!importFile" @click="handleImport">导入</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showImportResult" preset="dialog" title="导入结果" style="width: 480px">
      <div v-if="importResult">
        <n-descriptions bordered :column="1">
          <n-descriptions-item label="总数">{{ importResult.total }}</n-descriptions-item>
          <n-descriptions-item label="成功">
            <n-text type="success">{{ importResult.success }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="失败">
            <n-text type="error">{{ importResult.failed }}</n-text>
          </n-descriptions-item>
        </n-descriptions>
        <div v-if="importResult.errors && importResult.errors.length > 0" style="margin-top: 12px">
          <n-text depth="3">失败详情:</n-text>
          <div v-for="(err, idx) in importResult.errors" :key="idx" style="margin-top: 4px">
            <n-text type="error">第 {{ err.row }} 行 ({{ err.name }}): {{ err.reason }}</n-text>
          </div>
        </div>
      </div>
      <template #action>
        <n-button type="primary" @click="showImportResult = false">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NSpace,
  NButton,
  NSelect,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NIcon,
  NText,
  NDescriptions,
  NDescriptionsItem,
  NUpload,
  NUploadDragger,
  NP,
  NSwitch,
  useMessage,
  useDialog,
} from 'naive-ui'
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import { SearchOutline, AddOutline } from '@vicons/ionicons5'
import BackupStatus from '@/components/BackupStatus.vue'
import { getDevices, getGroups, createDevice, updateDevice, deleteDevice, enableDevice, disableDevice, importDevices, exportDevices } from '@/api/device'
import { triggerDeviceBackup } from '@/api/backup'
import type { Device, Group, DeviceCreateRequest, ImportResult } from '@/types'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const creating = ref(false)
const exportLoading = ref(false)
const importLoading = ref(false)
const devices = ref<Device[]>([])
const total = ref(0)
const searchKeyword = ref('')
const selectedGroup = ref<number | null>(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showImportModal = ref(false)
const showImportResult = ref(false)
const createFormRef = ref<FormInst | null>(null)
const editFormRef = ref<FormInst | null>(null)
const importFile = ref<File | null>(null)
const importResult = ref<ImportResult | null>(null)
const updating = ref(false)

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
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
})

const createRules: FormRules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  ip: [{ required: true, message: '请输入 IP 地址', trigger: 'blur' }],
  vendor: [{ required: true, message: '请选择厂商', trigger: 'change' }],
  model: [{ required: true, message: '请选择型号', trigger: 'change' }],
  protocol: [{ required: true, message: '请选择连接协议', trigger: 'change' }],
  port: [{ required: true, type: 'number', message: '请输入端口', trigger: 'blur' }],
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
      message.success('设备已禁用')
    } else {
      await enableDevice(device.id)
      message.success('设备已启用')
    }
    fetchDevices()
  } catch {
    message.error(device.enabled ? '禁用设备失败' : '启用设备失败')
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
    total.value = res.data.total
    paginationReactive.itemCount = res.data.total
  } catch {
    message.error('获取设备列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchGroups() {
  try {
    const res = await getGroups({ page: 1, page_size: 100 })
    groupOptions.value = res.data.items.map((g: Group) => ({ label: g.name, value: g.id }))
  } catch {
    message.error('获取分组列表失败')
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
    message.success('备份任务已触发')
  } catch {
    message.error('备份任务触发失败')
  }
}

function handleDelete(device: Device) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除设备 "${device.name}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteDevice(device.id)
        message.success('删除成功')
        fetchDevices()
      } catch {
        message.error('删除失败')
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
  showEditModal.value = true
}

async function handleUpdate() {
  try {
    await editFormRef.value?.validate()
  } catch {
    return
  }
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
    message.success('设备更新成功')
    showEditModal.value = false
    fetchDevices()
  } catch {
    message.error('设备更新失败')
  } finally {
    updating.value = false
  }
}

async function handleCreate() {
  try {
    await createFormRef.value?.validate()
  } catch {
    return
  }

  creating.value = true
  try {
    await createDevice(newDevice)
    message.success('设备添加成功')
    showCreateModal.value = false
    resetNewDevice()
    fetchDevices()
  } catch {
    message.error('设备添加失败')
  } finally {
    creating.value = false
  }
}

function handleFileChange(data: { fileList: UploadFileInfo[] }) {
  if (data.fileList.length > 0 && data.fileList[0].file) {
    importFile.value = data.fileList[0].file
  } else {
    importFile.value = null
  }
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
    message.error('导入设备失败')
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
    message.success('导出成功')
  } catch {
    message.error('导出设备失败')
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
  border-bottom: 1px solid #efeff5;
  padding: 12px;
  text-align: left;
  vertical-align: middle;
  white-space: nowrap;
}

.device-table th {
  color: #1f2937;
  font-weight: 600;
  background: #fafafa;
}

.device-table tr:hover td {
  background: #fafcff;
}

.table-empty {
  color: #8a8f99;
  text-align: center;
  padding: 32px 12px;
}
</style>

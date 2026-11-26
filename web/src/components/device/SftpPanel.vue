<template>
  <t-card title="SFTP 主动上传" class="sftp-panel">
    <template #actions>
      <t-tag v-if="status && status.enabled" theme="success" variant="light" size="small">
        服务已启用
      </t-tag>
      <t-tag v-else theme="default" variant="light" size="small">服务未启用</t-tag>
    </template>

    <t-alert v-if="status && !status.enabled" theme="warning" class="sftp-alert">
      服务端未启用 SFTP 监听，设备无法上传。请在服务器配置中设置 <code>sftp.enabled=true</code> 后重启。
    </t-alert>
    <t-alert v-if="statusError" theme="error" class="sftp-alert">
      SFTP 服务信息加载失败。
      <t-button size="small" variant="text" theme="primary" @click="fetchStatus">重试</t-button>
    </t-alert>

    <t-descriptions bordered :column="2" class="sftp-conn">
      <t-descriptions-item label="服务地址">
        {{ status?.listen_addr || '未监听' }}
      </t-descriptions-item>
      <t-descriptions-item label="主机密钥指纹">
        <code class="fingerprint">{{ status?.fingerprint || '-' }}</code>
      </t-descriptions-item>
      <t-descriptions-item label="上传账号">
        {{ account?.username || '尚未创建' }}
      </t-descriptions-item>
      <t-descriptions-item label="账号状态">
        <t-tag v-if="account" :theme="account.status === 'active' ? 'success' : 'default'" size="small">
          {{ account.status === 'active' ? '正常' : '已禁用' }}
        </t-tag>
        <span v-else>-</span>
      </t-descriptions-item>
    </t-descriptions>

    <t-space class="sftp-actions">
      <t-button v-if="!account" theme="primary" :loading="accountLoading" @click="handleCreate">
        创建上传账号
      </t-button>
      <template v-else>
        <t-button variant="outline" :loading="accountLoading" @click="handleReset">重置口令</t-button>
        <t-button
          v-if="account.status === 'active'"
          variant="outline"
          :loading="accountLoading"
          @click="handleDisable"
        >
          禁用账号
        </t-button>
        <t-button v-else variant="outline" :loading="accountLoading" @click="handleEnable">
          启用账号
        </t-button>
        <t-button theme="danger" variant="outline" :loading="accountLoading" @click="handleRevoke">
          撤销账号
        </t-button>
      </template>
    </t-space>

    <t-divider>最近上传</t-divider>

    <t-alert v-if="uploadsError" theme="error" class="sftp-alert">
      上传记录加载失败。
      <t-button size="small" variant="text" theme="primary" @click="fetchUploads">重试</t-button>
    </t-alert>

    <t-table
      :data="uploads"
      :columns="uploadColumns"
      :loading="uploadsLoading"
      row-key="id"
      size="small"
      table-layout="fixed"
    >
      <template #received_at="{ row }">{{ formatDateTime(row.received_at) }}</template>
      <template #received_bytes="{ row }">{{ formatBytes(row.received_bytes) }}</template>
      <template #stage="{ row }">
        <t-tag :theme="stageTheme(row.stage)" size="small">{{ stageLabel(row.stage) }}</t-tag>
      </template>
      <template #failure_reason="{ row }">
        <span :title="row.failure_reason || ''">{{ row.failure_reason || '-' }}</span>
      </template>
      <template #actions="{ row }">
        <t-button
          v-if="row.stage === 'failed'"
          size="small"
          variant="outline"
          :loading="retryingId === row.id"
          @click="handleRetry(row.id)"
        >
          重试归档
        </t-button>
        <span v-else>-</span>
      </template>
      <template #empty>
        <span class="empty-state-text">暂无上传记录，设备尚未上传过配置</span>
      </template>
    </t-table>

    <!-- 一次性口令：关闭即丢弃，不写入 store/localStorage -->
    <t-dialog
      v-model:visible="showPasswordDialog"
      header="SFTP 上传口令"
      :footer="false"
      :close-on-overlay-click="false"
      @close="discardPassword"
    >
      <t-alert theme="warning" class="sftp-alert">
        口令只显示这一次，关闭后将无法再次查看。请立即配置到设备上。
      </t-alert>
      <t-descriptions bordered :column="1" class="sftp-conn">
        <t-descriptions-item label="用户名">{{ oneTimeUsername }}</t-descriptions-item>
        <t-descriptions-item label="口令">
          <code class="one-time-password" data-testid="sftp-one-time-password">{{ oneTimePassword }}</code>
        </t-descriptions-item>
      </t-descriptions>
      <t-space class="sftp-actions">
        <t-button theme="primary" @click="discardPassword">我已保存，关闭</t-button>
      </t-space>
    </t-dialog>
  </t-card>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import {
  getSftpStatus,
  getSftpAccount,
  createSftpAccount,
  resetSftpAccount,
  enableSftpAccount,
  disableSftpAccount,
  revokeSftpAccount,
  getSftpUploads,
  retrySftpUpload,
} from '@/api/sftp'
import type { SftpAccount, SftpStatus, SftpUpload } from '@/types'
import { formatDateTime } from '@/utils/datetime'

const props = defineProps<{ deviceId: number }>()

const status = ref<SftpStatus | null>(null)
const account = ref<SftpAccount | null>(null)
const uploads = ref<SftpUpload[]>([])
const accountLoading = ref(false)
const uploadsLoading = ref(false)
const retryingId = ref<number | null>(null)
const statusError = ref(false)
const uploadsError = ref(false)

const showPasswordDialog = ref(false)
const oneTimeUsername = ref('')
const oneTimePassword = ref('')

const uploadColumns: PrimaryTableCol<SftpUpload>[] = [
  { title: '接收时间', colKey: 'received_at', width: 160 },
  { title: '大小', colKey: 'received_bytes', width: 90 },
  { title: '归档状态', colKey: 'stage', width: 100 },
  { title: '失败原因', colKey: 'failure_reason', ellipsis: true },
  { title: '操作', colKey: 'actions', width: 110 },
]

const STAGE_LABELS: Record<string, string> = {
  received: '已接收',
  archiving: '归档中',
  archived: '已归档',
  failed: '归档失败',
  rejected: '已拒绝',
}

function stageLabel(stage: string) {
  return STAGE_LABELS[stage] || stage
}

function stageTheme(stage: string) {
  switch (stage) {
    case 'archived':
      return 'success'
    case 'failed':
      return 'danger'
    case 'rejected':
      return 'warning'
    case 'archiving':
      return 'primary'
    default:
      return 'default'
  }
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`
}

async function fetchStatus() {
  try {
    const res = await getSftpStatus()
    status.value = res.data
    statusError.value = false
  } catch {
    status.value = null
    statusError.value = true
  }
}

async function fetchAccount() {
  try {
    const res = await getSftpAccount(props.deviceId)
    account.value = res.data
  } catch {
    // 404 means no account has been created yet.
    account.value = null
  }
}

async function fetchUploads() {
  uploadsLoading.value = true
  try {
    const res = await getSftpUploads(props.deviceId, { page: 1, page_size: 20 })
    uploads.value = res.data.items || []
    uploadsError.value = false
  } catch {
    uploads.value = []
    uploadsError.value = true
  } finally {
    uploadsLoading.value = false
  }
}

function showCredentials(username: string, password: string) {
  oneTimeUsername.value = username
  // Kept only in component state; cleared as soon as the dialog closes.
  oneTimePassword.value = password
  showPasswordDialog.value = true
}

function discardPassword() {
  oneTimePassword.value = ''
  oneTimeUsername.value = ''
  showPasswordDialog.value = false
}

async function handleCreate() {
  accountLoading.value = true
  try {
    const res = await createSftpAccount(props.deviceId)
    account.value = res.data.account
    showCredentials(res.data.account.username, res.data.password)
    MessagePlugin.success('账号已创建')
  } catch {
    MessagePlugin.error('创建账号失败')
  } finally {
    accountLoading.value = false
  }
}

async function handleReset() {
  accountLoading.value = true
  try {
    const res = await resetSftpAccount(props.deviceId)
    account.value = res.data.account
    showCredentials(res.data.account.username, res.data.password)
    MessagePlugin.success('口令已重置，旧口令立即失效')
  } catch {
    MessagePlugin.error('重置口令失败')
  } finally {
    accountLoading.value = false
  }
}

async function handleEnable() {
  accountLoading.value = true
  try {
    await enableSftpAccount(props.deviceId)
    await fetchAccount()
    MessagePlugin.success('账号已启用')
  } catch {
    MessagePlugin.error('启用账号失败')
  } finally {
    accountLoading.value = false
  }
}

async function handleDisable() {
  accountLoading.value = true
  try {
    await disableSftpAccount(props.deviceId)
    await fetchAccount()
    MessagePlugin.success('账号已禁用')
  } catch {
    MessagePlugin.error('禁用账号失败')
  } finally {
    accountLoading.value = false
  }
}

async function handleRevoke() {
  accountLoading.value = true
  try {
    await revokeSftpAccount(props.deviceId)
    account.value = null
    MessagePlugin.success('账号已撤销')
  } catch {
    MessagePlugin.error('撤销账号失败')
  } finally {
    accountLoading.value = false
  }
}

async function handleRetry(id: number) {
  retryingId.value = id
  try {
    await retrySftpUpload(id)
    MessagePlugin.success('已重新排队归档')
    await fetchUploads()
  } catch {
    MessagePlugin.error('重试失败')
  } finally {
    retryingId.value = null
  }
}

onMounted(async () => {
  await Promise.all([fetchStatus(), fetchAccount(), fetchUploads()])
})

watch(
  () => props.deviceId,
  async () => {
    await Promise.all([fetchAccount(), fetchUploads()])
  },
)
</script>

<style scoped>
.sftp-panel {
  margin-top: 16px;
}

.sftp-alert {
  margin-bottom: 12px;
}

.sftp-conn {
  margin-bottom: 12px;
}

.sftp-actions {
  margin-top: 12px;
}

.fingerprint {
  word-break: break-all;
  font-size: 12px;
}

.one-time-password {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.5px;
  word-break: break-all;
}

.empty-state-text {
  color: var(--td-text-color-placeholder, #999);
}
</style>

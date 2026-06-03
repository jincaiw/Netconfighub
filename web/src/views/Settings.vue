<template>
  <div class="settings">
    <n-grid :x-gap="16" :y-gap="16" :cols="1">
      <n-gi>
        <n-card title="修改密码">
          <n-form :model="passwordForm" label-placement="left" label-width="100">
            <n-form-item label="旧密码">
              <n-input v-model:value="passwordForm.old_password" type="password" show-password-on="click" placeholder="请输入当前密码" style="width: 300px" />
            </n-form-item>
            <n-form-item label="新密码">
              <n-input v-model:value="passwordForm.new_password" type="password" show-password-on="click" placeholder="请输入新密码（至少6位）" style="width: 300px" />
            </n-form-item>
            <n-form-item label="确认新密码">
              <n-input v-model:value="passwordForm.confirm_password" type="password" show-password-on="click" placeholder="请再次输入新密码" style="width: 300px" />
            </n-form-item>
            <n-form-item>
              <n-button type="primary" @click="handleChangePassword" :loading="changingPassword">修改密码</n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-gi>

      <n-gi>
        <n-card title="API Token 管理">
          <template #header-extra>
            <n-button type="primary" @click="openCreateTokenModal">创建 Token</n-button>
          </template>
          <n-data-table :columns="tokenColumns" :data="tokens" :loading="tokensLoading" size="small" />
        </n-card>
      </n-gi>

      <n-gi>
        <n-card title="备份策略配置">
          <n-form :model="backupConfig" label-placement="left" label-width="140">
            <n-form-item label="默认备份间隔(小时)">
              <n-input-number v-model:value="backupConfig.interval" :min="1" :max="720" style="width: 200px" />
            </n-form-item>
            <n-form-item>
              <n-button type="primary" @click="handleSaveConfig" :loading="savingConfig">保存配置</n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-gi>

      <n-gi>
        <n-card title="配置脱敏">
          <n-form :model="sanitizeConfig" label-placement="left" label-width="180">
            <n-form-item label="启用配置脱敏">
              <n-switch v-model:value="sanitizeConfig.enabled" />
            </n-form-item>
            <n-form-item label="脱敏密码字段">
              <n-switch v-model:value="sanitizeConfig.mask_passwords" :disabled="!sanitizeConfig.enabled" />
            </n-form-item>
            <n-form-item label="脱敏 SNMP Community">
              <n-switch v-model:value="sanitizeConfig.mask_snmp_community" :disabled="!sanitizeConfig.enabled" />
            </n-form-item>
            <n-form-item label="脱敏 Enable Secret">
              <n-switch v-model:value="sanitizeConfig.mask_enable_secret" :disabled="!sanitizeConfig.enabled" />
            </n-form-item>
            <n-form-item>
              <n-button type="primary" @click="handleSaveSanitizeConfig" :loading="savingSanitize">保存脱敏配置</n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-gi>
    </n-grid>

    <n-modal v-model:show="showCreateToken" preset="dialog" title="创建 API Token" style="width: 440px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="Token 名称">
          <n-input v-model:value="newTokenName" placeholder="请输入 Token 名称" />
        </n-form-item>
        <n-form-item label="过期时间">
          <n-date-picker
            v-model:value="newTokenExpiresAt"
            type="datetime"
            placeholder="请选择过期时间"
            style="width: 100%"
            :is-date-disabled="(ts: number) => ts < Date.now()"
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showCreateToken = false">取消</n-button>
        <n-button type="primary" :loading="creatingToken" @click="handleCreateToken">创建</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showTokenResult" preset="dialog" title="Token 创建成功" :closable="false" style="width: 520px">
      <n-alert type="warning" style="margin-bottom: 16px">
        请立即复制 Token，关闭后将无法再次查看。
      </n-alert>
      <n-input :value="newlyCreatedToken" readonly type="textarea" :rows="3" />
      <template #action>
        <n-button type="primary" @click="handleCopyToken">复制 Token</n-button>
        <n-button @click="showTokenResult = false">关闭</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import {
  NGrid,
  NGi,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NButton,
  NDataTable,
  NModal,
  NAlert,
  NDatePicker,
  NSwitch,
  useMessage,
  useDialog,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getTokens, createToken, deleteToken } from '@/api/token'
import { getConfigs, setConfig } from '@/api/config'
import { changePassword } from '@/api/auth'
import type { ApiToken } from '@/types'

const message = useMessage()
const dialog = useDialog()

const tokens = ref<ApiToken[]>([])
const tokensLoading = ref(false)
const showCreateToken = ref(false)
const showTokenResult = ref(false)
const newTokenName = ref('')
const newTokenExpiresAt = ref<number | null>(null)
const newlyCreatedToken = ref('')
const creatingToken = ref(false)
const savingSanitize = ref(false)
const savingConfig = ref(false)
const changingPassword = ref(false)

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const backupConfig = reactive({
  interval: 24,
})

const sanitizeConfig = reactive({
  enabled: false,
  mask_passwords: true,
  mask_snmp_community: true,
  mask_enable_secret: true,
})

const tokenColumns: DataTableColumns<ApiToken> = [
  { title: '名称', key: 'name', width: 160 },
  { title: '创建时间', key: 'created_at', width: 180 },
  { title: '过期时间', key: 'expires_at', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: (row) =>
      h(
        NButton,
        { size: 'small', type: 'error', onClick: () => handleDeleteToken(row) },
        { default: () => '删除' },
      ),
  },
]

async function fetchTokens() {
  tokensLoading.value = true
  try {
    const res = await getTokens()
    tokens.value = res.data
  } catch {
    message.error('获取 Token 列表失败')
  } finally {
    tokensLoading.value = false
  }
}

function openCreateTokenModal() {
  newTokenName.value = ''
  newTokenExpiresAt.value = null
  showCreateToken.value = true
}

function handleDeleteToken(token: ApiToken) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除 Token "${token.name}" 吗？删除后使用该 Token 的应用将无法访问 API。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteToken(token.id)
        message.success('Token 已删除')
        fetchTokens()
      } catch {
        message.error('Token 删除失败')
      }
    },
  })
}

async function handleCreateToken() {
  if (!newTokenName.value.trim()) {
    message.warning('请输入 Token 名称')
    return
  }
  if (!newTokenExpiresAt.value) {
    message.warning('请选择过期时间')
    return
  }

  creatingToken.value = true
  try {
    const expiresAt = new Date(newTokenExpiresAt.value).toISOString()
    const res = await createToken({ name: newTokenName.value.trim(), expires_at: expiresAt })
    newlyCreatedToken.value = res.data.token || ''
    showCreateToken.value = false
    showTokenResult.value = true
    fetchTokens()
  } catch {
    message.error('Token 创建失败')
  } finally {
    creatingToken.value = false
  }
}

function handleCopyToken() {
  navigator.clipboard.writeText(newlyCreatedToken.value).then(() => {
    message.success('Token 已复制到剪贴板')
  }).catch(() => {
    message.error('复制失败')
  })
}

async function handleSaveConfig() {
  savingConfig.value = true
  try {
    await setConfig('backup_interval', `${backupConfig.interval}h`)
    message.success('备份策略已保存')
  } catch {
    message.error('保存配置失败')
  } finally {
    savingConfig.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.old_password) {
    message.warning('请输入旧密码')
    return
  }
  if (!passwordForm.new_password || passwordForm.new_password.length < 6) {
    message.warning('新密码至少6位')
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    message.warning('两次输入的新密码不一致')
    return
  }
  changingPassword.value = true
  try {
    await changePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
    message.success('密码修改成功，请重新登录')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
    setTimeout(() => {
      window.location.href = '/login'
    }, 1500)
  } catch (err: any) {
    const msg = err?.response?.data?.message || '密码修改失败'
    message.error(msg)
  } finally {
    changingPassword.value = false
  }
}

async function fetchSystemConfigs() {
  try {
    const res = await getConfigs()
    const configs = res.data
    const findConfig = (key: string) => configs.find((c: { key: string; value: string }) => c.key === key)
    const intervalConfig = findConfig('backup_interval')
    if (intervalConfig) backupConfig.interval = parseInt(intervalConfig.value.replace(/h$/, ''), 10) || 24
    const enabled = findConfig('sanitize.enabled')
    if (enabled) sanitizeConfig.enabled = enabled.value === 'true'
    const maskPasswords = findConfig('sanitize.mask_passwords')
    if (maskPasswords) sanitizeConfig.mask_passwords = maskPasswords.value === 'true'
    const maskSnmp = findConfig('sanitize.mask_snmp_community')
    if (maskSnmp) sanitizeConfig.mask_snmp_community = maskSnmp.value === 'true'
    const maskEnable = findConfig('sanitize.mask_enable_secret')
    if (maskEnable) sanitizeConfig.mask_enable_secret = maskEnable.value === 'true'
  } catch {}
}

async function handleSaveSanitizeConfig() {
  savingSanitize.value = true
  try {
    await Promise.all([
      setConfig('sanitize.enabled', String(sanitizeConfig.enabled)),
      setConfig('sanitize.mask_passwords', String(sanitizeConfig.mask_passwords)),
      setConfig('sanitize.mask_snmp_community', String(sanitizeConfig.mask_snmp_community)),
      setConfig('sanitize.mask_enable_secret', String(sanitizeConfig.mask_enable_secret)),
    ])
    message.success('脱敏配置已保存')
  } catch {
    message.error('保存脱敏配置失败')
  } finally {
    savingSanitize.value = false
  }
}

onMounted(() => {
  fetchTokens()
  fetchSystemConfigs()
})
</script>

<style scoped>
.settings {
  width: 100%;
}
</style>

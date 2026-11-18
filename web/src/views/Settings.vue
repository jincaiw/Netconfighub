<template>
  <div class="settings">
    <t-alert
      v-if="configLoadFailed"
      theme="error"
      message="系统配置加载失败，保存已锁定。请重试后再保存，避免用默认值覆盖服务端配置。"
      style="margin-bottom: 16px"
    >
      <template #operation>
        <t-button size="small" variant="outline" :loading="configLoading" @click="fetchSystemConfigs">重试</t-button>
      </template>
    </t-alert>
    <t-row :gutter="[16, 16]">
      <t-col :xs="12">
        <t-card title="修改密码">
          <t-form :data="passwordForm" label-align="left" label-width="100px">
            <t-form-item label="旧密码">
              <t-input
                v-model="passwordForm.old_password"
                type="password"
                placeholder="请输入当前密码"
                style="width: 300px"
              />
            </t-form-item>
            <t-form-item label="新密码">
              <t-input
                v-model="passwordForm.new_password"
                type="password"
                placeholder="请输入新密码（至少6位）"
                style="width: 300px"
              />
            </t-form-item>
            <t-form-item label="确认新密码">
              <t-input
                v-model="passwordForm.confirm_password"
                type="password"
                placeholder="请再次输入新密码"
                style="width: 300px"
              />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" :loading="changingPassword" @click="handleChangePassword">修改密码</t-button>
            </t-form-item>
          </t-form>
        </t-card>
      </t-col>

      <t-col :xs="12">
        <t-card title="接口令牌（Token）管理">
          <template #actions>
            <t-button theme="primary" @click="openCreateTokenModal">创建接口令牌（Token）</t-button>
          </template>
          <t-table :columns="tokenColumns" :data="tokens" :loading="tokensLoading" row-key="id" size="small">
            <template #empty><span class="empty-state-text">暂无数据</span></template>
          </t-table>
        </t-card>
      </t-col>

      <t-col :xs="12">
        <t-card title="备份策略配置">
          <t-form :data="backupConfig" label-align="left" label-width="140px">
            <t-form-item label="默认备份间隔(小时)">
              <t-input-number v-model="backupConfig.interval" :min="1" :max="720" style="width: 200px" />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" :loading="savingConfig" :disabled="!configLoaded" @click="handleSaveConfig">保存配置</t-button>
            </t-form-item>
          </t-form>
        </t-card>
      </t-col>

      <t-col :xs="12">
        <t-card title="配置脱敏">
          <t-form :data="sanitizeConfig" label-align="left" label-width="180px">
            <t-form-item label="启用配置脱敏">
              <t-switch v-model="sanitizeConfig.enabled" data-testid="sanitize-main-switch" />
            </t-form-item>
            <t-form-item label="脱敏密码字段" :class="{ 'sanitize-subitem--disabled': !sanitizeConfig.enabled }">
              <t-switch v-model="sanitizeConfig.mask_passwords" :disabled="!sanitizeConfig.enabled" data-testid="sanitize-sub-passwords" />
            </t-form-item>
            <t-form-item label="脱敏 SNMP Community" :class="{ 'sanitize-subitem--disabled': !sanitizeConfig.enabled }">
              <t-switch v-model="sanitizeConfig.mask_snmp_community" :disabled="!sanitizeConfig.enabled" data-testid="sanitize-sub-snmp" />
            </t-form-item>
            <t-form-item label="脱敏特权模式密码" :class="{ 'sanitize-subitem--disabled': !sanitizeConfig.enabled }">
              <t-switch v-model="sanitizeConfig.mask_enable_secret" :disabled="!sanitizeConfig.enabled" data-testid="sanitize-sub-enable-secret" />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" :loading="savingSanitize" :disabled="!configLoaded" @click="handleSaveSanitizeConfig">保存脱敏配置</t-button>
            </t-form-item>
          </t-form>
        </t-card>
      </t-col>
    </t-row>

    <t-dialog v-model:visible="showCreateToken" header="创建接口令牌（Token）" :width="440">
      <t-form label-align="left" label-width="80px">
        <t-form-item label="令牌名称">
          <t-input v-model="newTokenName" placeholder="请输入令牌名称" />
        </t-form-item>
        <t-form-item label="过期时间">
          <t-date-picker
            v-model="newTokenExpiresAt"
            enable-time-picker
            format="YYYY-MM-DD HH:mm:ss"
            placeholder="请选择过期时间"
            style="width: 100%"
            :disable-date="(date: any) => new Date(date).getTime() < Date.now()"
          />
        </t-form-item>
      </t-form>
      <template #footer>
        <t-button variant="outline" @click="showCreateToken = false">取消</t-button>
        <t-button theme="primary" :loading="creatingToken" @click="handleCreateToken">创建</t-button>
      </template>
    </t-dialog>

    <t-dialog v-model:visible="showTokenResult" header="令牌创建成功" :close-btn="false" :width="520">
      <t-alert theme="warning" message="请立即复制令牌，关闭后将无法再次查看。" style="margin-bottom: 16px" />
      <t-textarea v-model="newlyCreatedToken" :autosize="{ minRows: 3 }" readonly />
      <template #footer>
        <t-button theme="primary" @click="handleCopyToken">复制令牌</t-button>
        <t-button variant="outline" @click="showTokenResult = false">关闭</t-button>
      </template>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Button, DialogPlugin, MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import { getTokens, createToken, deleteToken } from '@/api/token'
import { getConfigs, setConfig } from '@/api/config'
import { changePassword } from '@/api/auth'
import type { ApiToken } from '@/types'

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
const configLoading = ref(false)
const configLoaded = ref(false)
const configLoadFailed = ref(false)

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

const tokenColumns: PrimaryTableCol<ApiToken>[] = [
  { title: '名称', colKey: 'name', width: 160 },
  { title: '创建时间', colKey: 'created_at', width: 180 },
  { title: '过期时间', colKey: 'expires_at', width: 180 },
  {
    title: '操作',
    colKey: 'actions',
    width: 80,
    cell: (h, { row }) =>
      h(
        Button,
        { size: 'small', theme: 'danger', variant: 'outline', onClick: () => handleDeleteToken(row) },
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
    MessagePlugin.error('获取令牌列表失败')
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
  DialogPlugin.confirm({
    header: '确认删除',
    body: `确定要删除令牌 "${token.name}" 吗？删除后使用该令牌的应用将无法访问 API。`,
    theme: 'warning',
    confirmBtn: '删除',
    cancelBtn: '取消',
    onConfirm: async () => {
      try {
        await deleteToken(token.id)
        MessagePlugin.success('令牌已删除')
        fetchTokens()
      } catch {
        MessagePlugin.error('令牌删除失败')
      }
    },
  })
}

async function handleCreateToken() {
  if (!newTokenName.value.trim()) {
    MessagePlugin.warning('请输入令牌名称')
    return
  }
  if (!newTokenExpiresAt.value) {
    MessagePlugin.warning('请选择过期时间')
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
    MessagePlugin.error('令牌创建失败')
  } finally {
    creatingToken.value = false
  }
}

function handleCopyToken() {
  navigator.clipboard.writeText(newlyCreatedToken.value).then(() => {
    MessagePlugin.success('令牌已复制到剪贴板')
  }).catch(() => {
    MessagePlugin.error('复制失败')
  })
}

async function handleSaveConfig() {
  if (!configLoaded.value) {
    MessagePlugin.warning('系统配置尚未加载成功，无法保存')
    return
  }
  savingConfig.value = true
  try {
    await setConfig('backup_interval', `${backupConfig.interval}h`)
    MessagePlugin.success('备份策略已保存')
  } catch {
    MessagePlugin.error('保存配置失败')
  } finally {
    savingConfig.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.old_password) {
    MessagePlugin.warning('请输入旧密码')
    return
  }
  if (!passwordForm.new_password || passwordForm.new_password.length < 6) {
    MessagePlugin.warning('新密码至少6位')
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    MessagePlugin.warning('两次输入的新密码不一致')
    return
  }
  changingPassword.value = true
  try {
    await changePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
    MessagePlugin.success('密码修改成功，请重新登录')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
    setTimeout(() => {
      window.location.href = '/login'
    }, 1500)
  } catch (err: any) {
    const msg = err?.response?.data?.message || '密码修改失败'
    MessagePlugin.error(msg)
  } finally {
    changingPassword.value = false
  }
}

async function fetchSystemConfigs() {
  configLoading.value = true
  configLoaded.value = false
  configLoadFailed.value = false
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
    configLoaded.value = true
  } catch {
    configLoadFailed.value = true
    MessagePlugin.error('系统配置加载失败，已锁定保存，请重试')
  } finally {
    configLoading.value = false
  }
}

async function handleSaveSanitizeConfig() {
  if (!configLoaded.value) {
    MessagePlugin.warning('系统配置尚未加载成功，无法保存')
    return
  }
  savingSanitize.value = true
  try {
    await Promise.all([
      setConfig('sanitize.enabled', String(sanitizeConfig.enabled)),
      setConfig('sanitize.mask_passwords', String(sanitizeConfig.mask_passwords)),
      setConfig('sanitize.mask_snmp_community', String(sanitizeConfig.mask_snmp_community)),
      setConfig('sanitize.mask_enable_secret', String(sanitizeConfig.mask_enable_secret)),
    ])
    MessagePlugin.success('脱敏配置已保存')
  } catch {
    MessagePlugin.error('保存脱敏配置失败')
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

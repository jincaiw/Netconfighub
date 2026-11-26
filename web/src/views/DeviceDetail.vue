<template>
  <div class="device-detail">
    <div class="page-heading">
      <div>
        <t-button variant="text" @click="router.push('/devices')">
          <template #icon><ChevronLeftIcon /></template>
          返回列表
        </t-button>
        <h1>设备详情</h1>
      </div>
      <t-space v-if="!notFound">
        <t-button
          :theme="device?.enabled ? 'warning' : 'success'"
          variant="outline"
          :disabled="!device"
          @click="handleToggleEnabled"
        >
          {{ device?.enabled ? '禁用设备' : '启用设备' }}
        </t-button>
        <t-button
          theme="primary"
          :loading="backupLoading"
          :disabled="isSftpPush"
          :title="isSftpPush ? '该设备使用 SFTP 主动上传，平台不会拉取' : ''"
          @click="handleBackup"
        >
          触发备份
        </t-button>
      </t-space>
    </div>

    <t-card v-if="notFound" class="device-notfound" style="margin-top: 16px">
      <t-empty description="设备不存在或已被删除">
        <template #action>
          <t-button variant="outline" @click="router.push('/devices')">返回设备列表</t-button>
        </template>
      </t-empty>
    </t-card>
    <t-loading v-else :loading="loading">
      <t-card title="设备信息" style="margin-top: 16px">
        <t-descriptions bordered :column="descColumns">
          <t-descriptions-item label="设备名称">{{ device?.name }}</t-descriptions-item>
          <t-descriptions-item label="IP 地址">{{ device?.ip }}</t-descriptions-item>
          <t-descriptions-item label="厂商">{{ vendorLabel }}</t-descriptions-item>
          <t-descriptions-item label="采集方式">
            <t-tag v-if="isSftpPush" theme="primary" variant="light" size="small">SFTP 主动上传</t-tag>
            <t-tag v-else variant="outline" size="small">平台拉取</t-tag>
          </t-descriptions-item>
          <t-descriptions-item label="型号">{{ modelLabel }}</t-descriptions-item>
          <t-descriptions-item label="连接协议">{{ device?.protocol?.toUpperCase() }}</t-descriptions-item>
          <t-descriptions-item label="端口">{{ device?.port }}</t-descriptions-item>
          <t-descriptions-item label="分组">{{ device?.group_name || '-' }}</t-descriptions-item>
          <t-descriptions-item label="启用状态">
            <t-tag :theme="device?.enabled ? 'success' : 'default'" variant="light" size="small">
              {{ device?.enabled ? '已启用' : '已禁用' }}
            </t-tag>
          </t-descriptions-item>
          <t-descriptions-item label="备份状态">
            <BackupStatus v-if="device" :status="device.last_backup_status" />
          </t-descriptions-item>
          <t-descriptions-item label="最近备份"><span class="datetime-cell" :title="device?.last_backup_at || ''">{{ formatDateTime(device?.last_backup_at) }}</span></t-descriptions-item>
          <t-descriptions-item label="失败原因" :span="2">{{ device?.last_failure_reason || '-' }}</t-descriptions-item>
        </t-descriptions>
      </t-card>

      <t-card style="margin-top: 16px">
        <t-tabs v-model="activeTab">
          <t-tab-panel value="latest" label="最新配置">
            <t-loading :loading="configLoading">
              <div class="config-content">
                <pre><code><span class="empty-state-text">{{ latestConfig || '暂无配置' }}</span></code></pre>
              </div>
            </t-loading>
          </t-tab-panel>
          <t-tab-panel value="history" label="历史版本">
            <t-space v-if="configVersions.length >= 2" :wrap="true" class="version-compare" style="margin-bottom: 12px">
              <t-select
                v-model="selectedFromHash"
                :options="versionOptions"
                placeholder="选择源版本"
                class="version-select"
              />
              <t-select
                v-model="selectedToHash"
                :options="versionOptions"
                placeholder="选择目标版本"
                class="version-select"
              />
              <t-button
                theme="primary"
                :disabled="!selectedFromHash || !selectedToHash || selectedFromHash === selectedToHash"
                :loading="diffLoading"
                @click="handleDiff"
              >
                比对
              </t-button>
            </t-space>
            <t-table :columns="versionColumns" :data="configVersions" :loading="versionsLoading" row-key="hash" />
            <t-dialog v-model:visible="showDiffModal" header="配置比对结果" :width="800" :footer="false">
              <div v-if="diffResult">
                <t-space style="margin-bottom: 8px">
                  <t-tag theme="success" variant="light">新增 {{ diffResult.added_lines }} 行</t-tag>
                  <t-tag theme="danger" variant="light">删除 {{ diffResult.removed_lines }} 行</t-tag>
                </t-space>
                <div class="diff-viewer">
                  <pre
v-for="(line, idx) in diffLines" :key="idx"
                    :class="{
                      'diff-line-added': line.startsWith('+'),
                      'diff-line-removed': line.startsWith('-'),
                      'diff-line-context': line.startsWith(' ')
                    }"
                  >{{ line }}</pre>
                </div>
              </div>
            </t-dialog>
          </t-tab-panel>
          <t-tab-panel value="deviations" label="基线偏差">
            <t-loading :loading="deviationsLoading">
              <t-empty v-if="!deviationsLoading && deviations.length === 0" description="暂无偏差" />
              <div v-else>
                <t-card
                  v-for="deviation in deviations"
                  :key="deviation.id"
                  size="small"
                  style="margin-bottom: 8px"
                >
                  <pre class="diff-content">{{ deviation.diff_content }}</pre>
                  <template #footer>
                    <span class="diff-time">检测时间: <span class="datetime-cell" :title="deviation.detected_at || ''">{{ formatDateTime(deviation.detected_at) }}</span></span>
                  </template>
                </t-card>
              </div>
            </t-loading>
          </t-tab-panel>
        </t-tabs>
      </t-card>

      <SftpPanel v-if="isSftpPush && !isInvalidId" :device-id="deviceId" />
    </t-loading>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import { ChevronLeftIcon } from 'tdesign-icons-vue-next'
import { getDevice, enableDevice, disableDevice } from '@/api/device'
import { triggerDeviceBackup } from '@/api/backup'
import { getDeviceConfig, getDeviceVersions, getDeviceDeviations, getConfigDiff } from '@/api/config'
import type { Device, ConfigVersion, Deviation, DiffResult } from '@/types'
import BackupStatus from '@/components/BackupStatus.vue'
import SftpPanel from '@/components/device/SftpPanel.vue'
import { formatDateTime } from '@/utils/datetime'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const configLoading = ref(false)
const versionsLoading = ref(false)
const deviationsLoading = ref(false)
const backupLoading = ref(false)
const diffLoading = ref(false)
const device = ref<Device | null>(null)
const notFound = ref(false)
const latestConfig = ref('')
const configVersions = ref<ConfigVersion[]>([])
const deviations = ref<Deviation[]>([])
const activeTab = ref('latest')

const deviceId = Number(route.params.id)
const isSftpPush = computed(() => device.value?.collection_mode === 'sftp_push')
const isInvalidId = Number.isNaN(deviceId) || deviceId <= 0

// 响应式：窄屏减少信息表列数，避免标签竖排/内容挤压
const viewportWidth = ref(window.innerWidth)
const onResize = () => { viewportWidth.value = window.innerWidth }
const descColumns = computed(() => (viewportWidth.value < 768 ? 1 : viewportWidth.value < 1100 ? 2 : 3))
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))
const selectedFromHash = ref<string | null>(null)
const selectedToHash = ref<string | null>(null)
const diffResult = ref<DiffResult | null>(null)
const showDiffModal = ref(false)

const vendorMap: Record<string, string> = {
  cisco: '思科',
  h3c: '华三',
  huawei: '华为',
  ruijie: '锐捷',
}

const modelMap: Record<string, string> = {
  ios: 'IOS',
  vrp: 'VRP',
  comware: 'Comware',
  'rg-os': 'RG-OS',
}

const vendorLabel = computed(() => {
  if (!device.value) return ''
  return vendorMap[device.value.vendor] || device.value.vendor
})

const modelLabel = computed(() => {
  if (!device.value) return ''
  return modelMap[device.value.model] || device.value.model
})

const versionOptions = computed(() =>
  configVersions.value.map((v) => ({
    label: `${v.hash.substring(0, 8)} - ${v.message} (${v.date})`,
    value: v.hash,
  })),
)

const diffLines = computed(() => {
  if (!diffResult.value?.diff_content) return []
  return diffResult.value.diff_content.split('\n')
})

const versionColumns: PrimaryTableCol<ConfigVersion>[] = [
  { title: 'Hash', colKey: 'hash', ellipsis: true },
  { title: '提交信息', colKey: 'message' },
  { title: '作者', colKey: 'author', width: 120 },
  { title: '日期', colKey: 'date', width: 180 },
]

async function fetchDevice() {
  loading.value = true
  try {
    const res = await getDevice(deviceId)
    device.value = res.data
    notFound.value = false
  } catch (err: any) {
    // 设备不存在（404）：显示友好“未找到”态，避免渲染空白详情与操作按钮。
    if (err?.response?.status === 404) {
      notFound.value = true
      device.value = null
    } else {
      MessagePlugin.error('获取设备信息失败')
    }
  } finally {
    loading.value = false
  }
}

async function fetchLatestConfig() {
  configLoading.value = true
  try {
    const res = await getDeviceConfig(deviceId)
    latestConfig.value = res.data.content
  } catch (err: any) {
    // 设备尚无配置（404）属正常空态，不弹错误提示，由内联「暂无配置」呈现。
    if (err?.response?.status === 404) latestConfig.value = ''
    else MessagePlugin.error('获取设备配置失败')
  } finally {
    configLoading.value = false
  }
}

async function fetchConfigVersions() {
  versionsLoading.value = true
  try {
    const res = await getDeviceVersions(deviceId)
    configVersions.value = res.data
  } catch (err: any) {
    if (err?.response?.status === 404) configVersions.value = []
    else MessagePlugin.error('获取历史版本失败')
  } finally {
    versionsLoading.value = false
  }
}

async function fetchDeviations() {
  deviationsLoading.value = true
  try {
    const res = await getDeviceDeviations(deviceId, { page: 1, page_size: 20 })
    deviations.value = res.data.items
  } catch (err: any) {
    if (err?.response?.status === 404) deviations.value = []
    else MessagePlugin.error('获取基线偏差失败')
  } finally {
    deviationsLoading.value = false
  }
}

async function handleDiff() {
  if (!selectedFromHash.value || !selectedToHash.value) return
  diffLoading.value = true
  try {
    const res = await getConfigDiff(deviceId, selectedFromHash.value, selectedToHash.value)
    diffResult.value = res.data
    showDiffModal.value = true
  } catch {
    MessagePlugin.error('配置比对失败')
  } finally {
    diffLoading.value = false
  }
}

async function handleBackup() {
  backupLoading.value = true
  try {
    await triggerDeviceBackup(deviceId)
    MessagePlugin.success('备份任务已触发')
  } catch {
    MessagePlugin.error('备份任务触发失败')
  } finally {
    backupLoading.value = false
  }
}

async function handleToggleEnabled() {
  if (!device.value) return
  try {
    if (device.value.enabled) {
      await disableDevice(device.value.id)
      MessagePlugin.success('设备已禁用')
    } else {
      await enableDevice(device.value.id)
      MessagePlugin.success('设备已启用')
    }
    fetchDevice()
  } catch {
    MessagePlugin.error(device.value!.enabled ? '禁用设备失败' : '启用设备失败')
  }
}

onMounted(() => {
  if (isInvalidId) {
    MessagePlugin.error('无效的设备 ID')
    router.push('/devices')
    return
  }
  fetchDevice()
  fetchLatestConfig()
})

watch(activeTab, (tab) => {
  if (tab === 'latest') fetchLatestConfig()
  else if (tab === 'history') fetchConfigVersions()
  else if (tab === 'deviations') fetchDeviations()
})
</script>

<style scoped>
.device-detail {
  width: 100%;
}

.page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.page-heading h1 {
  display: inline-block;
  margin-left: 8px;
  font-size: 22px;
  color: var(--td-text-color-primary);
}

.config-content {
  background-color: var(--td-bg-color-container-hover);
  border-radius: var(--td-radius-default);
  padding: 16px;
  overflow-x: auto;
}

.config-content pre {
  margin: 0;
  font-family: 'SF Mono', 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.diff-content {
  margin: 0;
  font-family: 'SF Mono', 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.diff-viewer {
  background-color: var(--td-bg-color-container-hover);
  border-radius: var(--td-radius-default);
  padding: 12px;
  overflow-x: auto;
  max-height: 500px;
  overflow-y: auto;
}

.diff-viewer pre {
  margin: 0;
  font-family: 'SF Mono', 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  padding: 0 4px;
}

.diff-line-added {
  background-color: var(--td-success-color-1);
}

.diff-line-removed {
  background-color: var(--td-error-color-1);
}

.diff-line-context {
  background-color: transparent;
}

.diff-time {
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}

.version-compare {
  width: 100%;
}

.version-select {
  width: min(300px, 100%);
  flex: 1 1 240px;
}

@media (max-width: 768px) {
  .version-select {
    width: 100%;
    flex-basis: 100%;
  }
}
</style>

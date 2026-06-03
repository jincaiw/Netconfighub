<template>
  <div class="device-detail">
    <n-page-header @back="router.push('/devices')" title="设备详情" subtitle="">
      <template #extra>
        <n-space>
          <n-button @click="router.push('/devices')">返回列表</n-button>
          <n-button :type="device?.enabled ? 'warning' : 'success'" @click="handleToggleEnabled" :disabled="!device">
            {{ device?.enabled ? '禁用设备' : '启用设备' }}
          </n-button>
          <n-button type="primary" :loading="backupLoading" @click="handleBackup">触发备份</n-button>
        </n-space>
      </template>
    </n-page-header>

    <n-spin :show="loading">
      <n-card title="设备信息" style="margin-top: 16px">
        <n-descriptions bordered :column="3">
          <n-descriptions-item label="设备名称">{{ device?.name }}</n-descriptions-item>
          <n-descriptions-item label="IP 地址">{{ device?.ip }}</n-descriptions-item>
          <n-descriptions-item label="厂商">{{ vendorLabel }}</n-descriptions-item>
          <n-descriptions-item label="型号">{{ modelLabel }}</n-descriptions-item>
          <n-descriptions-item label="连接协议">{{ device?.protocol?.toUpperCase() }}</n-descriptions-item>
          <n-descriptions-item label="端口">{{ device?.port }}</n-descriptions-item>
          <n-descriptions-item label="分组">{{ device?.group_name || '-' }}</n-descriptions-item>
          <n-descriptions-item label="启用状态">
            <n-tag :type="device?.enabled ? 'success' : 'default'" size="small">
              {{ device?.enabled ? '已启用' : '已禁用' }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="备份状态">
            <BackupStatus v-if="device" :status="device.last_backup_status" />
          </n-descriptions-item>
          <n-descriptions-item label="最近备份">{{ device?.last_backup_at || '-' }}</n-descriptions-item>
          <n-descriptions-item label="失败原因" :span="2">{{ device?.last_failure_reason || '-' }}</n-descriptions-item>
        </n-descriptions>
      </n-card>

      <n-card style="margin-top: 16px">
        <n-tabs v-model:value="activeTab" type="line">
          <n-tab-pane name="latest" tab="最新配置">
            <n-spin :show="configLoading">
              <div class="config-content">
                <pre><code>{{ latestConfig || '暂无配置' }}</code></pre>
              </div>
            </n-spin>
          </n-tab-pane>
          <n-tab-pane name="history" tab="历史版本">
            <n-space style="margin-bottom: 12px" v-if="configVersions.length >= 2">
              <n-select
                v-model:value="selectedFromHash"
                :options="versionOptions"
                placeholder="选择源版本"
                style="width: 300px"
              />
              <n-select
                v-model:value="selectedToHash"
                :options="versionOptions"
                placeholder="选择目标版本"
                style="width: 300px"
              />
              <n-button
                type="primary"
                :disabled="!selectedFromHash || !selectedToHash || selectedFromHash === selectedToHash"
                :loading="diffLoading"
                @click="handleDiff"
              >
                比对
              </n-button>
            </n-space>
            <n-data-table
              :columns="versionColumns"
              :data="configVersions"
              :loading="versionsLoading"
              size="small"
            />
            <n-modal v-model:show="showDiffModal" preset="dialog" title="配置比对结果" style="width: 800px">
              <div v-if="diffResult">
                <n-space style="margin-bottom: 8px">
                  <n-tag type="success">新增 {{ diffResult.added_lines }} 行</n-tag>
                  <n-tag type="error">删除 {{ diffResult.removed_lines }} 行</n-tag>
                </n-space>
                <div class="diff-viewer">
                  <pre v-for="(line, idx) in diffLines" :key="idx"
                    :class="{
                      'diff-line-added': line.startsWith('+'),
                      'diff-line-removed': line.startsWith('-'),
                      'diff-line-context': line.startsWith(' ')
                    }"
                  >{{ line }}</pre>
                </div>
              </div>
              <template #action>
                <n-button @click="showDiffModal = false">关闭</n-button>
              </template>
            </n-modal>
          </n-tab-pane>
          <n-tab-pane name="deviations" tab="基线偏差">
            <n-spin :show="deviationsLoading">
              <n-empty v-if="!deviationsLoading && deviations.length === 0" description="暂无偏差" />
              <div v-else>
                <n-card
                  v-for="deviation in deviations"
                  :key="deviation.id"
                  size="small"
                  style="margin-bottom: 8px"
                >
                  <pre class="diff-content">{{ deviation.diff_content }}</pre>
                  <template #footer>
                    <n-text depth="3">检测时间: {{ deviation.detected_at }}</n-text>
                  </template>
                </n-card>
              </div>
            </n-spin>
          </n-tab-pane>
        </n-tabs>
      </n-card>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NPageHeader,
  NCard,
  NDescriptions,
  NDescriptionsItem,
  NTabs,
  NTabPane,
  NDataTable,
  NEmpty,
  NText,
  NButton,
  NSpace,
  NSpin,
  NSelect,
  NModal,
  NTag,
  useMessage,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getDevice, enableDevice, disableDevice } from '@/api/device'
import { triggerDeviceBackup } from '@/api/backup'
import { getDeviceConfig, getDeviceVersions, getDeviceDeviations, getConfigDiff } from '@/api/config'
import type { Device, ConfigVersion, Deviation, DiffResult } from '@/types'
import BackupStatus from '@/components/BackupStatus.vue'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const configLoading = ref(false)
const versionsLoading = ref(false)
const deviationsLoading = ref(false)
const backupLoading = ref(false)
const diffLoading = ref(false)
const device = ref<Device | null>(null)
const latestConfig = ref('')
const configVersions = ref<ConfigVersion[]>([])
const deviations = ref<Deviation[]>([])
const activeTab = ref('latest')

const deviceId = Number(route.params.id)
const isInvalidId = Number.isNaN(deviceId) || deviceId <= 0
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
  }))
)

const diffLines = computed(() => {
  if (!diffResult.value?.diff_content) return []
  return diffResult.value.diff_content.split('\n')
})

const versionColumns: DataTableColumns<ConfigVersion> = [
  { title: 'Hash', key: 'hash', ellipsis: { tooltip: true } },
  { title: '提交信息', key: 'message' },
  { title: '作者', key: 'author', width: 120 },
  { title: '日期', key: 'date', width: 180 },
]

async function fetchDevice() {
  loading.value = true
  try {
    const res = await getDevice(deviceId)
    device.value = res.data
  } catch {
    message.error('获取设备信息失败')
  } finally {
    loading.value = false
  }
}

async function fetchLatestConfig() {
  configLoading.value = true
  try {
    const res = await getDeviceConfig(deviceId)
    latestConfig.value = res.data.content
  } catch {
    message.error('获取设备配置失败')
  } finally {
    configLoading.value = false
  }
}

async function fetchConfigVersions() {
  versionsLoading.value = true
  try {
    const res = await getDeviceVersions(deviceId)
    configVersions.value = res.data
  } catch {
    message.error('获取历史版本失败')
  } finally {
    versionsLoading.value = false
  }
}

async function fetchDeviations() {
  deviationsLoading.value = true
  try {
    const res = await getDeviceDeviations(deviceId, { page: 1, page_size: 20 })
    deviations.value = res.data.items
  } catch {
    message.error('获取基线偏差失败')
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
    message.error('配置比对失败')
  } finally {
    diffLoading.value = false
  }
}

async function handleBackup() {
  backupLoading.value = true
  try {
    await triggerDeviceBackup(deviceId)
    message.success('备份任务已触发')
  } catch {
    message.error('备份任务触发失败')
  } finally {
    backupLoading.value = false
  }
}

async function handleToggleEnabled() {
  if (!device.value) return
  try {
    if (device.value.enabled) {
      await disableDevice(device.value.id)
      message.success('设备已禁用')
    } else {
      await enableDevice(device.value.id)
      message.success('设备已启用')
    }
    fetchDevice()
  } catch {
    message.error(device.value!.enabled ? '禁用设备失败' : '启用设备失败')
  }
}

onMounted(() => {
  if (isInvalidId) {
    message.error('无效的设备 ID')
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

.config-content {
  background-color: #f5f5f5;
  border-radius: 4px;
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

.config-content code {
  font-family: inherit;
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
  background-color: #fafafa;
  border-radius: 4px;
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
  background-color: #e6ffec;
}

.diff-line-removed {
  background-color: #ffebe9;
}

.diff-line-context {
  background-color: transparent;
}
</style>

<template>
  <div class="alert-list">
    <n-card title="告警管理">
      <template #header-extra>
        <n-space>
          <n-select
            v-model:value="filterType"
            :options="typeOptions"
            placeholder="按类型筛选"
            clearable
            style="width: 160px"
            @update:value="handleFilter"
          />
          <n-select
            v-model:value="filterSeverity"
            :options="severityOptions"
            placeholder="按级别筛选"
            clearable
            style="width: 140px"
            @update:value="handleFilter"
          />
          <n-button @click="handleMarkAllRead" :disabled="alerts.length === 0">全部已读</n-button>
        </n-space>
      </template>
      <n-data-table
        :columns="columns"
        :data="alerts"
        :loading="loading"
        :pagination="paginationReactive"
        :row-key="(row: Alert) => row.id"
        remote
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NDataTable,
  NSpace,
  NButton,
  NSelect,
  NTag,
  useMessage,
  useDialog,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAlerts, markAlertAsRead, markAllAlertsAsRead, deleteAlert } from '@/api/alert'
import type { Alert } from '@/types'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const alerts = ref<Alert[]>([])
const total = ref(0)
const filterType = ref<string | null>(null)
const filterSeverity = ref<string | null>(null)

const typeOptions = [
  { label: '备份失败', value: 'backup_failed' },
  { label: '基线偏差', value: 'baseline_deviation' },
]

const severityOptions = [
  { label: '错误', value: 'error' },
  { label: '警告', value: 'warning' },
  { label: '信息', value: 'info' },
]

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

const severityColorMap: Record<string, 'default' | 'error' | 'warning' | 'info' | 'success' | 'primary'> = {
  error: 'error',
  warning: 'warning',
  info: 'info',
}

const severityLabelMap: Record<string, string> = {
  error: '错误',
  warning: '警告',
  info: '信息',
}

const typeLabelMap: Record<string, string> = {
  backup_failed: '备份失败',
  baseline_deviation: '基线偏差',
}

const columns: DataTableColumns<Alert> = [
  {
    title: '状态',
    key: 'is_read',
    width: 70,
    render: (row) =>
      h(NTag, { size: 'small', type: row.is_read ? 'default' : 'warning' }, { default: () => row.is_read ? '已读' : '未读' }),
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row) => h(NTag, { size: 'small' }, { default: () => typeLabelMap[row.type] || row.type }),
  },
  {
    title: '级别',
    key: 'severity',
    width: 80,
    render: (row) => h(NTag, { size: 'small', type: severityColorMap[row.severity] || 'default' }, { default: () => severityLabelMap[row.severity] || row.severity }),
  },
  { title: '设备', key: 'device_name', width: 120 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '消息', key: 'message', ellipsis: { tooltip: true } },
  { title: '时间', key: 'created_at', width: 160 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleViewDevice(row), disabled: !row.device_id },
            { default: () => '查看设备' },
          ),
          h(
            NButton,
            { size: 'small', type: row.is_read ? 'default' : 'primary', onClick: () => handleMarkRead(row) },
            { default: () => row.is_read ? '已读' : '标记已读' },
          ),
          h(
            NButton,
            { size: 'small', type: 'error', onClick: () => handleDelete(row) },
            { default: () => '删除' },
          ),
        ],
      }),
  },
]

async function fetchAlerts() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: paginationReactive.page,
      page_size: paginationReactive.pageSize,
    }
    if (filterType.value) params.type = filterType.value
    if (filterSeverity.value) params.severity = filterSeverity.value
    const res = await getAlerts(params)
    alerts.value = res.data.items
    total.value = res.data.total
    paginationReactive.itemCount = res.data.total
  } catch {
    message.error('获取告警列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  paginationReactive.page = page
  fetchAlerts()
}

function handlePageSizeChange(pageSize: number) {
  paginationReactive.pageSize = pageSize
  paginationReactive.page = 1
  fetchAlerts()
}

function handleFilter() {
  paginationReactive.page = 1
  fetchAlerts()
}

function handleViewDevice(alert: Alert) {
  if (alert.device_id) {
    router.push(`/devices/${alert.device_id}`)
  }
}

async function handleMarkRead(alert: Alert) {
  if (alert.is_read) return
  try {
    await markAlertAsRead(alert.id)
    alert.is_read = true
    message.success('已标记为已读')
  } catch {
    message.error('标记已读失败')
  }
}

async function handleMarkAllRead() {
  try {
    await markAllAlertsAsRead()
    alerts.value.forEach(a => { a.is_read = true })
    message.success('已全部标记为已读')
  } catch {
    message.error('标记全部已读失败')
  }
}

function handleDelete(alert: Alert) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除该告警吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteAlert(alert.id)
        message.success('告警已删除')
        fetchAlerts()
      } catch {
        message.error('删除告警失败')
      }
    },
  })
}

onMounted(() => {
  fetchAlerts()
})
</script>

<style scoped>
.alert-list {
  width: 100%;
}
</style>

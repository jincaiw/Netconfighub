<template>
  <div class="alert-list">
    <t-card title="告警管理">
      <template #actions>
        <t-space>
          <t-select
            v-model="filterType"
            :options="typeOptions"
            placeholder="按类型筛选"
            clearable
            style="width: 160px"
            @change="handleFilter"
          />
          <t-select
            v-model="filterSeverity"
            :options="severityOptions"
            placeholder="按级别筛选"
            clearable
            style="width: 140px"
            @change="handleFilter"
          />
          <t-button :disabled="alerts.length === 0" @click="handleMarkAllRead">全部已读</t-button>
        </t-space>
      </template>
      <div v-if="loadError" class="list-error" role="alert">
        <span>告警列表加载失败，请重新加载</span>
        <t-button size="small" variant="text" theme="primary" @click="fetchAlerts">重新加载</t-button>
      </div>
      <t-table
        :columns="columns"
        :data="alerts"
        :loading="loading"
        :pagination="tablePagination"
        row-key="id"
        @page-change="handlePageChange"
      >
        <template #empty><span class="empty-state-text">暂无数据</span></template>
      </t-table>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Button, Space, Tag, DialogPlugin, MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import { getAlerts, markAlertAsRead, markAllAlertsAsRead, deleteAlert } from '@/api/alert'
import type { Alert } from '@/types'

const router = useRouter()

const loading = ref(false)
const loadError = ref(false)
const alerts = ref<Alert[]>([])
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
})

const tablePagination = computed(() => ({
  current: paginationReactive.page,
  pageSize: paginationReactive.pageSize,
  total: paginationReactive.itemCount,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50],
}))

const severityThemeMap: Record<string, 'danger' | 'warning' | 'primary'> = {
  error: 'danger',
  warning: 'warning',
  info: 'primary',
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

const columns: PrimaryTableCol<Alert>[] = [
  {
    title: '状态',
    colKey: 'is_read',
    width: 70,
    cell: (h, { row }) =>
      h(Tag, { size: 'small', theme: row.is_read ? 'default' : 'warning', variant: 'light' }, {
        default: () => (row.is_read ? '已读' : '未读'),
      }),
  },
  {
    title: '类型',
    colKey: 'type',
    width: 100,
    cell: (h, { row }) =>
      h(Tag, { size: 'small', variant: 'light' }, { default: () => typeLabelMap[row.type] || row.type }),
  },
  {
    title: '级别',
    colKey: 'severity',
    width: 80,
    cell: (h, { row }) =>
      h(Tag, { size: 'small', theme: severityThemeMap[row.severity] || 'default', variant: 'light' }, {
        default: () => severityLabelMap[row.severity] || row.severity,
      }),
  },
  { title: '设备', colKey: 'device_name', width: 120 },
  { title: '标题', colKey: 'title', ellipsis: true },
  { title: '消息', colKey: 'message', ellipsis: true },
  { title: '时间', colKey: 'created_at', width: 160 },
  {
    title: '操作',
    colKey: 'actions',
    width: 180,
    cell: (h, { row }) =>
      h(
        Space,
        { size: 'small' },
        {
          default: () => [
            h(
              Button,
              { size: 'small', onClick: () => handleViewDevice(row), disabled: !row.device_id },
              { default: () => '查看设备' },
            ),
            h(
              Button,
              {
                size: 'small',
                theme: row.is_read ? 'default' : 'primary',
                onClick: () => handleMarkRead(row),
              },
              { default: () => (row.is_read ? '已读' : '标记已读') },
            ),
            h(
              Button,
              { size: 'small', theme: 'danger', onClick: () => handleDelete(row) },
              { default: () => '删除' },
            ),
          ],
        },
      ),
  },
]

async function fetchAlerts() {
  loading.value = true
  loadError.value = false
  try {
    const params: Record<string, any> = {
      page: paginationReactive.page,
      page_size: paginationReactive.pageSize,
    }
    if (filterType.value) params.type = filterType.value
    if (filterSeverity.value) params.severity = filterSeverity.value
    const res = await getAlerts(params)
    alerts.value = res.data.items
    paginationReactive.itemCount = res.data.total
  } catch {
    loadError.value = true
    MessagePlugin.error('获取告警列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(pageInfo: { current: number; pageSize: number }) {
  if (pageInfo.pageSize && pageInfo.pageSize !== paginationReactive.pageSize) {
    paginationReactive.pageSize = pageInfo.pageSize
    paginationReactive.page = 1
  } else {
    paginationReactive.page = pageInfo.current
  }
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
    MessagePlugin.success('已标记为已读')
  } catch {
    MessagePlugin.error('标记已读失败')
  }
}

async function handleMarkAllRead() {
  try {
    await markAllAlertsAsRead()
    alerts.value.forEach((a) => {
      a.is_read = true
    })
    MessagePlugin.success('已全部标记为已读')
  } catch {
    MessagePlugin.error('标记全部已读失败')
  }
}

function handleDelete(alert: Alert) {
  DialogPlugin.confirm({
    header: '确认删除',
    body: '确定要删除该告警吗？',
    theme: 'warning',
    confirmBtn: '删除',
    cancelBtn: '取消',
    onConfirm: async () => {
      try {
        await deleteAlert(alert.id)
        MessagePlugin.success('告警已删除')
        fetchAlerts()
      } catch {
        MessagePlugin.error('删除告警失败')
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

.list-error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 12px;
  border: 1px solid var(--td-error-color-2);
  border-radius: var(--td-radius-default);
  color: var(--td-error-color-8);
  background: var(--td-error-color-1);
  font-size: 13px;
}
</style>

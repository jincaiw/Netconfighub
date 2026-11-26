<template>
  <div class="failed-devices">
    <t-card title="失败设备列表">
      <template #actions>
        <t-button theme="primary" :loading="retryingAll" :disabled="devices.length === 0" @click="handleRetryAll">
          全部重试
        </t-button>
      </template>

      <div v-if="loadError" class="list-error" role="alert">
        <span>失败设备列表加载失败，请重新加载</span>
        <t-button size="small" variant="text" theme="primary" @click="fetchFailedDevices">重新加载</t-button>
      </div>
      <t-loading :loading="loading">
        <t-empty v-if="!loadError && !loading && devices.length === 0" description="暂无失败设备，所有设备备份正常">
          <template #action>
            <t-button size="small" variant="outline" @click="router.push('/devices')">查看所有设备</t-button>
          </template>
        </t-empty>
        <t-table
          v-else
          :columns="columns"
          :data="devices"
          :pagination="tablePagination"
          @page-change="handlePageChange"
        />
      </t-loading>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Button, Space, DialogPlugin, MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import { getFailedDevices, triggerDeviceBackup } from '@/api/backup'
import BackupStatus from '@/components/BackupStatus.vue'
import { formatDateTime } from '@/utils/datetime'
import type { Device } from '@/types'

const router = useRouter()

const loading = ref(false)
const loadError = ref(false)
const retryingAll = ref(false)
const devices = ref<Device[]>([])

const vendorLabelMap: Record<string, string> = {
  cisco: '思科',
  h3c: '华三',
  huawei: '华为',
  ruijie: '锐捷',
}

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

const columns: PrimaryTableCol<Device>[] = [
  { title: '设备名称', colKey: 'name' },
  { title: 'IP 地址', colKey: 'ip', width: 140 },
  {
    title: '厂商',
    colKey: 'vendor',
    width: 90,
    cell: (_h, { row }) => vendorLabelMap[row.vendor] || row.vendor,
  },
  {
    title: '备份状态',
    colKey: 'last_backup_status',
    width: 100,
    cell: (h, { row }) => h(BackupStatus, { status: row.last_backup_status }),
  },
  { title: '失败原因', colKey: 'last_failure_reason', ellipsis: true },
  { title: '重试次数', colKey: 'retry_count', width: 90 },
  {
    title: '最近备份时间',
    colKey: 'last_backup_at',
    width: 170,
    cell: (h, { row }) =>
      h('span', { class: 'datetime-cell', title: row.last_backup_at ?? '' }, formatDateTime(row.last_backup_at)),
  },
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
              { size: 'small', theme: 'primary', onClick: () => handleRetry(row.id) },
              { default: () => '重新备份' },
            ),
            h(
              Button,
              { size: 'small', onClick: () => router.push(`/devices/${row.id}`) },
              { default: () => '查看详情' },
            ),
          ],
        },
      ),
  },
]

async function fetchFailedDevices() {
  loading.value = true
  try {
    const res = await getFailedDevices({
      page: paginationReactive.page,
      page_size: paginationReactive.pageSize,
    })
    devices.value = res.data.items
    paginationReactive.itemCount = res.data.total
    loadError.value = false
  } catch {
    loadError.value = true
    MessagePlugin.error('获取失败设备列表失败')
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
  fetchFailedDevices()
}

async function handleRetry(deviceId: number) {
  try {
    await triggerDeviceBackup(deviceId)
    MessagePlugin.success('备份任务已触发')
  } catch {
    MessagePlugin.error('备份任务触发失败')
  }
}

async function handleRetryAll() {
  const instance = DialogPlugin.confirm({
    header: '确认',
    body: `确定要重新备份全部 ${devices.value.length} 台失败设备吗？`,
    confirmBtn: '确定',
    cancelBtn: '取消',
    onConfirm: async () => {
      retryingAll.value = true
      try {
        const results = await Promise.allSettled(devices.value.map((dev) => triggerDeviceBackup(dev.id)))
        const succeeded = results.filter((r) => r.status === 'fulfilled').length
        const failed = results.filter((r) => r.status === 'rejected').length
        if (failed === 0) {
          MessagePlugin.success(`全部 ${succeeded} 台设备备份已触发`)
        } else {
          MessagePlugin.warning(`${succeeded} 台成功，${failed} 台失败`)
        }
        instance.hide()
      } finally {
        retryingAll.value = false
      }
    },
  })
}

onMounted(() => {
  fetchFailedDevices()
})
</script>

<style scoped>
.failed-devices {
  width: 100%;
}
</style>

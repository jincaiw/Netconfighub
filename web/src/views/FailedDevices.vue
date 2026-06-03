<template>
  <div class="failed-devices">
    <n-card title="失败设备列表">
      <template #header-extra>
        <n-button type="primary" @click="handleRetryAll" :loading="retryingAll" :disabled="devices.length === 0">全部重试</n-button>
      </template>

      <n-spin :show="loading">
        <n-empty v-if="!loading && devices.length === 0" description="暂无失败设备，所有设备备份正常">
          <template #extra>
            <n-button size="small" @click="router.push('/devices')">查看所有设备</n-button>
          </template>
        </n-empty>
        <n-data-table
          v-else
          :columns="columns"
          :data="devices"
          :pagination="paginationReactive"
          remote
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NEmpty,
  NSpin,
  useMessage,
  useDialog,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getFailedDevices, triggerDeviceBackup } from '@/api/backup'
import BackupStatus from '@/components/BackupStatus.vue'
import type { Device } from '@/types'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
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
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

const columns: DataTableColumns<Device> = [
  { title: '设备名称', key: 'name' },
  { title: 'IP 地址', key: 'ip', width: 140 },
  {
    title: '厂商',
    key: 'vendor',
    width: 90,
    render: (row) => vendorLabelMap[row.vendor] || row.vendor,
  },
  {
    title: '备份状态',
    key: 'last_backup_status',
    width: 100,
    render: (row) => h(BackupStatus, { status: row.last_backup_status }),
  },
  { title: '失败原因', key: 'last_failure_reason', ellipsis: { tooltip: true } },
  { title: '重试次数', key: 'retry_count', width: 90 },
  { title: '最近备份时间', key: 'last_backup_at', width: 170 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NButton,
            { size: 'small', type: 'primary', onClick: () => handleRetry(row.id) },
            { default: () => '重新备份' },
          ),
          h(
            NButton,
            { size: 'small', onClick: () => router.push(`/devices/${row.id}`) },
            { default: () => '查看详情' },
          ),
        ],
      }),
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
  } catch {
    message.error('获取失败设备列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  paginationReactive.page = page
  fetchFailedDevices()
}

function handlePageSizeChange(pageSize: number) {
  paginationReactive.pageSize = pageSize
  paginationReactive.page = 1
  fetchFailedDevices()
}

async function handleRetry(deviceId: number) {
  try {
    await triggerDeviceBackup(deviceId)
    message.success('备份任务已触发')
  } catch {
    message.error('备份任务触发失败')
  }
}

async function handleRetryAll() {
  dialog.warning({
    title: '确认',
    content: `确定要重新备份全部 ${devices.value.length} 台失败设备吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      retryingAll.value = true
      try {
        const results = await Promise.allSettled(
          devices.value.map(dev => triggerDeviceBackup(dev.id))
        )
        const succeeded = results.filter(r => r.status === 'fulfilled').length
        const failed = results.filter(r => r.status === 'rejected').length
        if (failed === 0) {
          message.success(`全部 ${succeeded} 台设备备份已触发`)
        } else {
          message.warning(`${succeeded} 台成功，${failed} 台失败`)
        }
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

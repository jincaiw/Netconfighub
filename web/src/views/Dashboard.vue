<template>
  <div class="dashboard">
    <n-grid :x-gap="16" :y-gap="16" :cols="4" responsive="screen" item-responsive>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic label="设备总数" :value="stats.totalDevices">
            <template #prefix>
              <n-icon :component="ServerOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic label="备份成功率" :value="stats.successRate">
            <template #prefix>
              <n-icon :component="CheckmarkCircleOutline" />
            </template>
            <template #suffix>%</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic label="失败设备数" :value="stats.failedDevices">
            <template #prefix>
              <n-icon :component="AlertCircleOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic label="未读告警" :value="stats.unreadAlerts">
            <template #prefix>
              <n-icon :component="NotificationsOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <n-grid :x-gap="16" :y-gap="16" :cols="2" style="margin-top: 16px" responsive="screen" item-responsive>
      <n-gi span="2 l:1">
        <n-card title="最近备份任务">
          <n-data-table :columns="taskColumns" :data="recentTasks" size="small" />
        </n-card>
      </n-gi>
      <n-gi span="2 l:1">
        <n-card title="失败设备概览">
          <n-data-table :columns="failedColumns" :data="failedDevices" size="small" />
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import {
  NGrid,
  NGi,
  NCard,
  NStatistic,
  NIcon,
  NDataTable,
  NButton,
  useMessage,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  ServerOutline,
  CheckmarkCircleOutline,
  AlertCircleOutline,
  NotificationsOutline,
} from '@vicons/ionicons5'
import BackupStatus from '@/components/BackupStatus.vue'
import { getDevices } from '@/api/device'
import { getBackupTasks, getFailedDevices } from '@/api/backup'
import { getUnreadAlertCount } from '@/api/alert'
import type { Device, BackupTask } from '@/types'

const router = useRouter()
const message = useMessage()

const stats = reactive({
  totalDevices: 0,
  successRate: 0,
  failedDevices: 0,
  unreadAlerts: 0,
})

interface RecentTaskRow {
  id: number
  device_name: string
  status: 'success' | 'failed' | 'pending' | 'running' | 'retrying' | 'timeout' | 'cancelled' | 'skipped'
  finished_at: string
}

interface FailedDeviceRow {
  id: number
  name: string
  ip: string
  reason: string
}

const recentTasks = ref<RecentTaskRow[]>([])
const failedDevices = ref<FailedDeviceRow[]>([])

const taskColumns: DataTableColumns<RecentTaskRow> = [
  { title: '设备', key: 'device_name' },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => h(BackupStatus, { status: row.status }),
  },
  { title: '完成时间', key: 'finished_at', width: 160 },
]

const failedColumns: DataTableColumns<FailedDeviceRow> = [
  { title: '设备', key: 'name' },
  { title: 'IP', key: 'ip', width: 130 },
  { title: '原因', key: 'reason', ellipsis: { tooltip: true } },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: (row) =>
      h(
        NButton,
        { size: 'small', onClick: () => router.push(`/devices/${row.id}`) },
        { default: () => '详情' },
      ),
  },
]

async function fetchStats() {
  try {
    const res = await getDevices({ page: 1, page_size: 1 })
    stats.totalDevices = res.data.total
  } catch {
    message.error('获取设备统计失败')
  }

  try {
    const backupRes = await getBackupTasks({ page: 1, page_size: 1 })
    const total = backupRes.data.total
    if (total > 0) {
      const allRes = await getBackupTasks({ page: 1, page_size: total })
      const successCount = allRes.data.items.filter((t: BackupTask) => t.status === 'success').length
      stats.successRate = Math.round((successCount / total) * 1000) / 10
    } else {
      stats.successRate = 0
    }
  } catch {
    message.error('获取备份统计失败')
  }

  try {
    const res = await getFailedDevices({ page: 1, page_size: 1 })
    stats.failedDevices = res.data.total
  } catch {
    message.error('获取失败设备统计失败')
  }

  try {
    const res = await getUnreadAlertCount()
    stats.unreadAlerts = res.data.count
  } catch {
    message.error('获取告警统计失败')
  }
}

async function fetchRecentTasks() {
  try {
    const res = await getBackupTasks({ page: 1, page_size: 10 })
    recentTasks.value = res.data.items.map((t: BackupTask) => ({
      id: t.id,
      device_name: t.device_name || `设备#${t.device_id}`,
      status: t.status,
      finished_at: t.finished_at || '-',
    }))
  } catch {
    message.error('获取备份任务列表失败')
  }
}

async function fetchFailedDevices() {
  try {
    const res = await getFailedDevices({ page: 1, page_size: 10 })
    failedDevices.value = res.data.items.map((d: Device) => ({
      id: d.id,
      name: d.name,
      ip: d.ip,
      reason: d.last_failure_reason || '未知原因',
    }))
  } catch {
    message.error('获取失败设备列表失败')
  }
}

onMounted(() => {
  fetchStats()
  fetchRecentTasks()
  fetchFailedDevices()
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}
</style>

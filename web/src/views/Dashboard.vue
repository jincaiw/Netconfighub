<template>
  <div class="dashboard">
    <div class="page-heading">
      <div>
        <h1>运维总览</h1>
        <p>实时掌握设备备份运行状态，及时发现并处理异常。</p>
      </div>
      <n-button type="primary" @click="router.push('/devices')">管理设备</n-button>
    </div>

    <n-grid :x-gap="16" :y-gap="16" :cols="4" responsive="screen" item-responsive>
      <n-gi span="4 m:2 l:1">
        <n-card class="metric-card" hoverable @click="router.push('/devices')">
          <n-statistic label="设备总数" :value="stats.totalDevices">
            <template #prefix>
              <span class="metric-icon metric-icon--green">
                <n-icon :component="ServerOutline" />
              </span>
            </template>
          </n-statistic>
          <div class="metric-footnote">已纳管网络设备</div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="metric-card" hoverable>
          <n-statistic label="备份成功率" :value="stats.successRate">
            <template #prefix>
              <span class="metric-icon metric-icon--green">
                <n-icon :component="CheckmarkCircleOutline" />
              </span>
            </template>
            <template #suffix>%</template>
          </n-statistic>
          <div class="metric-footnote">基于全部备份任务</div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="metric-card" hoverable @click="router.push('/failed')">
          <n-statistic label="失败设备数" :value="stats.failedDevices">
            <template #prefix>
              <span class="metric-icon metric-icon--red">
                <n-icon :component="AlertCircleOutline" />
              </span>
            </template>
          </n-statistic>
          <div class="metric-footnote">需要尽快处理</div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="metric-card" hoverable @click="router.push('/alerts')">
          <n-statistic label="未读告警" :value="stats.unreadAlerts">
            <template #prefix>
              <span class="metric-icon metric-icon--amber">
                <n-icon :component="NotificationsOutline" />
              </span>
            </template>
          </n-statistic>
          <div class="metric-footnote">待确认告警信息</div>
        </n-card>
      </n-gi>
    </n-grid>

    <div class="overview-grid">
      <div>
        <n-card class="overview-card" title="最近备份任务">
          <template #header-extra>
            <n-button text type="primary" @click="router.push('/devices')">查看设备</n-button>
          </template>
          <n-data-table
            :columns="taskColumns"
            :data="recentTasks"
            :scroll-x="560"
            size="small"
          />
        </n-card>
      </div>
      <div>
        <n-card class="overview-card" title="失败设备概览">
          <template #header-extra>
            <n-button text type="primary" @click="router.push('/failed')">查看全部</n-button>
          </template>
          <n-data-table
            :columns="failedColumns"
            :data="failedDevices"
            :scroll-x="520"
            size="small"
          />
        </n-card>
      </div>
    </div>
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
  max-width: 1680px;
  margin: 0 auto;
}

.page-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.page-heading h1 {
  color: #0f172a;
  font-size: 28px;
  line-height: 1.25;
  letter-spacing: -0.02em;
}

.page-heading p {
  margin-top: 7px;
  color: #64748b;
  font-size: 14px;
}

.metric-card {
  min-height: 148px;
  cursor: pointer;
  border-color: #dfe7f1;
}

.metric-card :deep(.n-card__content) {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.metric-card :deep(.n-statistic-value__content) {
  color: #0f172a;
  font-size: 30px;
  font-weight: 700;
}

.metric-card :deep(.n-statistic__label) {
  color: #475569;
  font-size: 14px;
  font-weight: 600;
}

.metric-icon {
  width: 42px;
  height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  border-radius: 50%;
  font-size: 22px;
}

.metric-icon--green {
  color: #059669;
  background: #e7f8f0;
}

.metric-icon--red {
  color: #dc2626;
  background: #fff0f0;
}

.metric-icon--amber {
  color: #d97706;
  background: #fff7e6;
}

.metric-footnote {
  margin-top: 16px;
  color: #94a3b8;
  font-size: 12px;
}

.overview-grid {
  display: grid;
  grid-template-columns: minmax(0, 3fr) minmax(360px, 2fr);
  gap: 16px;
  margin-top: 18px;
}

.overview-card {
  min-height: 430px;
  border-color: #dfe7f1;
}

@media (max-width: 1200px) {
  .overview-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .page-heading {
    align-items: center;
  }

  .page-heading h1 {
    font-size: 24px;
  }

  .page-heading p {
    max-width: 240px;
  }

}
</style>

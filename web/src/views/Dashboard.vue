<template>
  <div class="dashboard">
    <div class="page-heading">
      <div>
        <h1>运维总览</h1>
        <p>实时掌握设备备份运行状态，及时发现并处理异常。</p>
      </div>
      <t-button theme="primary" @click="router.push('/devices')">管理设备</t-button>
    </div>

    <div v-if="statsLoading" class="dashboard-loading" data-testid="dashboard-loading" aria-busy="true">
      <t-loading size="large" text="统计数据加载中..." />
    </div>
    <template v-else>
      <div v-if="failedParts.length > 0" class="dashboard-error" role="alert" data-testid="dashboard-error">
        <span>部分统计数据加载失败：{{ failedParts.join('、') }}</span>
        <t-button size="small" variant="outline" @click="fetchStats">重新加载</t-button>
      </div>
      <TopPanel class="dashboard-section" :stats="stats" @navigate="handleNavigate" />
      <MiddleChart class="dashboard-section" :tasks="allTasks" />
      <RankList class="dashboard-section" :failed-rank="failedRank" :count-rank="countRank" />
      <OutputOverview
        class="dashboard-section"
        :recent-tasks="recentTasks"
        :failed-devices="failedDevices"
        @detail="(id) => router.push(`/devices/${id}`)"
        @navigate="handleNavigate"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { getDevices } from '@/api/device'
import { getBackupTasks, getFailedDevices } from '@/api/backup'
import type { BackupTask, Device } from '@/types'
import { useAlertSummary } from '@/composables/useAlertSummary'
import TopPanel from './dashboard/components/TopPanel.vue'
import MiddleChart from './dashboard/components/MiddleChart.vue'
import RankList from './dashboard/components/RankList.vue'
import OutputOverview from './dashboard/components/OutputOverview.vue'
import {
  buildStats,
  rankByCount,
  rankByFailed,
  type FailedDeviceRow,
  type RecentTaskRow,
} from './dashboard/index'

const router = useRouter()
const { unreadCount: unreadAlerts, start: startAlertPolling, stop: stopAlertPolling } = useAlertSummary()

const totalDevices = ref(0)
const allTasks = ref<BackupTask[]>([])
const failedDeviceTotal = ref(0)
const failedDeviceRows = ref<FailedDeviceRow[]>([])

// 首屏加载态与分片失败收集：加载完成前不渲染 0 值，失败时给出可见反馈与重试入口。
const statsLoading = ref(true)
const failedParts = ref<string[]>([])

function markStatError(part: string) {
  if (!failedParts.value.includes(part)) failedParts.value.push(part)
}

function clearStatError(part: string) {
  failedParts.value = failedParts.value.filter((item) => item !== part)
}

// 派生数据全部来自同一次取数，避免每个子段重复请求
const stats = computed(() =>
  buildStats({
    totalDevices: totalDevices.value,
    tasks: allTasks.value,
    failedDevices: failedDeviceTotal.value,
    unreadAlerts: unreadAlerts.value,
  }),
)

const recentTasks = computed<RecentTaskRow[]>(() =>
  allTasks.value.slice(0, 10).map((task) => ({
    id: task.id,
    device_name: task.device_name || `设备#${task.device_id}`,
    status: task.status,
    finished_at: task.finished_at || '-',
  })),
)

const failedRank = computed(() => rankByFailed(allTasks.value))
const countRank = computed(() => rankByCount(allTasks.value))
const failedDevices = computed(() => failedDeviceRows.value)

type NavigateTarget = 'DeviceList' | 'FailedDevices' | 'AlertList'
const routeMap: Record<NavigateTarget, string> = {
  DeviceList: '/devices',
  FailedDevices: '/failed',
  AlertList: '/alerts',
}

function handleNavigate(target: NavigateTarget) {
  router.push(routeMap[target])
}

async function fetchStats() {
  statsLoading.value = true
  failedParts.value = []
  try {
    const res = await getDevices({ page: 1, page_size: 1 })
    totalDevices.value = res.data.total
    clearStatError('设备总数')
  } catch {
    markStatError('设备总数')
    MessagePlugin.error('获取设备统计失败')
  }

  try {
    const first = await getBackupTasks({ page: 1, page_size: 1 })
    const total = first.data.total
    allTasks.value = total > 0 ? (await getBackupTasks({ page: 1, page_size: total })).data.items : []
    clearStatError('备份任务')
  } catch {
    markStatError('备份任务')
    MessagePlugin.error('获取备份统计失败')
  }

  try {
    const res = await getFailedDevices({ page: 1, page_size: 10 })
    failedDeviceTotal.value = res.data.total
    failedDeviceRows.value = res.data.items.map((device: Device) => ({
      id: device.id,
      name: device.name,
      ip: device.ip,
      reason: device.last_failure_reason || '未知原因',
    }))
    clearStatError('失败设备')
  } catch {
    markStatError('失败设备')
    MessagePlugin.error('获取失败设备统计失败')
  }
  // 首屏计数由 startAlertPolling()/start() 统一发起，此处不再重复请求（保持去重契约）。
  statsLoading.value = false
}

onMounted(() => {
  fetchStats()
  startAlertPolling()
})

onUnmounted(() => {
  stopAlertPolling()
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
  color: var(--td-text-color-primary);
  font-size: 28px;
  line-height: 1.25;
  letter-spacing: -0.02em;
}

.page-heading p {
  margin-top: 7px;
  color: var(--td-text-color-secondary);
  font-size: 14px;
}

.dashboard-section {
  display: block;
  margin-bottom: 16px;
}

.dashboard-loading {
  min-height: 240px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dashboard-error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  padding: 10px 14px;
  border: 1px solid var(--td-error-color-2);
  border-radius: var(--td-radius-default);
  color: var(--td-error-color);
  background: var(--td-error-color-1);
  font-size: 13px;
}

@media (max-width: 768px) {
  .page-heading {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .page-heading h1 {
    font-size: 24px;
  }

  .page-heading p {
    max-width: 240px;
  }
}
</style>

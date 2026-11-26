<template>
  <div class="dashboard">
    <!-- shadcn 风格页头：大标题 + Download -->
    <div class="page-heading">
      <div>
        <h1>运维总览</h1>
        <p>实时掌握设备备份运行状态，及时发现并处理异常。</p>
      </div>
      <t-button theme="primary" @click="exportReport">导出报表</t-button>
    </div>

    <div v-if="statsLoading" class="dashboard-loading" data-testid="dashboard-loading" aria-busy="true">
      <t-loading size="large" text="统计数据加载中..." />
    </div>
    <template v-else>
      <div v-if="failedParts.length > 0" class="dashboard-error" role="alert" data-testid="dashboard-error">
        <span>部分统计数据加载失败：{{ failedParts.join('、') }}</span>
        <t-button size="small" variant="outline" @click="fetchStats">重新加载</t-button>
      </div>

      <t-tabs v-model="activeTab" class="dashboard-tabs">
        <t-tab-panel value="overview" label="概览">
          <TopPanel class="dashboard-section" :stats="stats" :tasks="allTasks" @navigate="handleNavigate" />
          <MiddleChart class="dashboard-section" :tasks="allTasks" />
          <RankList class="dashboard-section" :tasks="allTasks" />
          <OutputOverview class="dashboard-section" :tasks="allTasks" @navigate="handleNavigate" />
        </t-tab-panel>
        <t-tab-panel value="analytics" label="分析">
          <MiddleChart class="dashboard-section" :tasks="allTasks" />
          <RankList class="dashboard-section" :tasks="allTasks" />
        </t-tab-panel>
        <t-tab-panel value="reports" label="报表">
          <ReportsPanel class="dashboard-section" :stats="stats" :tasks="allTasks" />
        </t-tab-panel>
        <t-tab-panel value="notifications" label="通知">
          <NotificationsPanel class="dashboard-section" :active="activeTab === 'notifications'" />
        </t-tab-panel>
      </t-tabs>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { getDevices } from '@/api/device'
import { getBackupHistory, getFailedDevices } from '@/api/backup'
import type { BackupTask } from '@/types'
import { useAlertSummary } from '@/composables/useAlertSummary'
import TopPanel from './dashboard/components/TopPanel.vue'
import MiddleChart from './dashboard/components/MiddleChart.vue'
import RankList from './dashboard/components/RankList.vue'
import OutputOverview from './dashboard/components/OutputOverview.vue'
import ReportsPanel from './dashboard/components/ReportsPanel.vue'
import NotificationsPanel from './dashboard/components/NotificationsPanel.vue'
import { downloadCsv } from '@/utils/csv'
import { buildStats } from './dashboard/index'

const router = useRouter()
const { unreadCount: unreadAlerts, start: startAlertPolling, stop: stopAlertPolling } = useAlertSummary()

const activeTab = ref('overview')
const totalDevices = ref(0)
const allTasks = ref<BackupTask[]>([])
const failedDeviceTotal = ref(0)

const statsLoading = ref(true)
const failedParts = ref<string[]>([])

function markStatError(part: string) {
  if (!failedParts.value.includes(part)) failedParts.value.push(part)
}
function clearStatError(part: string) {
  failedParts.value = failedParts.value.filter((item) => item !== part)
}

const stats = computed(() =>
  buildStats({
    totalDevices: totalDevices.value,
    tasks: allTasks.value,
    failedDevices: failedDeviceTotal.value,
    unreadAlerts: unreadAlerts.value,
  }),
)

type NavigateTarget = 'DeviceList' | 'FailedDevices' | 'AlertList'
const routeMap: Record<NavigateTarget, string> = {
  DeviceList: '/devices',
  FailedDevices: '/failed',
  AlertList: '/alerts',
}
function handleNavigate(target: NavigateTarget) {
  router.push(routeMap[target])
}

function exportReport() {
  const s0 = stats.value
  downloadCsv(`运维报表-${new Date().toISOString().slice(0, 10)}.csv`, [
    ['指标', '数值'],
    ['设备总数', s0.totalDevices],
    ['备份成功率(%)', s0.successRate],
    ['失败设备数', s0.failedDevices],
    ['未读告警', s0.unreadAlerts],
    ['生成时间', new Date().toLocaleString()],
  ])
  MessagePlugin.success('已导出运维报表')
}

async function fetchStats() {
  statsLoading.value = true
  failedParts.value = []
  // 三项统计彼此独立：并行拉取，缩短首屏等待（较原串行 await 更快）。
  await Promise.all([
    (async () => {
      try {
        const res = await getDevices({ page: 1, page_size: 1 })
        totalDevices.value = res.data.total
        clearStatError('设备总数')
      } catch {
        markStatError('设备总数')
        MessagePlugin.error('获取设备统计失败')
      }
    })(),
    (async () => {
      try {
        // 取最近 30 天完整历史（保留窗口内完整），使趋势/状态/排名基于真实数据而非 100 条样本。
        const historyRes = await getBackupHistory({ days: 30, limit: 10000 })
        allTasks.value = historyRes.data
        clearStatError('备份任务')
      } catch {
        markStatError('备份任务')
        MessagePlugin.error('获取备份统计失败')
      }
    })(),
    (async () => {
      try {
        const res = await getFailedDevices({ page: 1, page_size: 1 })
        failedDeviceTotal.value = res.data.total
        clearStatError('失败设备')
      } catch {
        markStatError('失败设备')
        MessagePlugin.error('获取失败设备统计失败')
      }
    })(),
  ])
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
  margin-bottom: 16px;
}
.page-heading h1 {
  color: var(--shadcn-foreground);
  font-size: 28px;
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.02em;
}
.page-heading p {
  margin-top: 7px;
  color: var(--shadcn-muted-foreground);
  font-size: 14px;
}
.dashboard-tabs :deep(.t-tabs__nav-item) {
  font-size: 14px;
}
.dashboard-section {
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
  border-radius: var(--shadcn-radius);
  color: var(--td-error-color);
  background: var(--td-error-color-1);
  font-size: 13px;
}
@media (max-width: 768px) {
  .page-heading { flex-direction: column; align-items: stretch; gap: 12px; }
  .page-heading h1 { font-size: 24px; }
  .page-heading p { max-width: 240px; }
}
</style>

<template>
  <t-row :gutter="16" class="rank-list">
    <t-col :xs="12" :xl="6">
      <t-card title="失败次数排名" class="dashboard-rank-card rank-card" :bordered="false">
        <template #actions>
          <t-radio-group v-model="failedPeriod" variant="default-filled">
            <t-radio-button value="week">本周</t-radio-button>
            <t-radio-button value="month">本月</t-radio-button>
          </t-radio-group>
        </template>
        <t-table :data="failedData" :columns="columns" row-key="deviceId">
          <template #index="{ rowIndex }">
            <span :class="rankClass(rowIndex)">{{ rowIndex + 1 }}</span>
          </template>
          <template #trend="{ row }">
            <trend v-if="row.trend" :type="row.trend.direction" :describe="`${row.trend.percent}%`" />
            <span v-else>—</span>
          </template>
          <template #operation="{ row }">
            <t-link theme="primary" @click="goDetail(row.deviceId)">详情</t-link>
          </template>
          <template #empty><t-empty description="暂无失败记录" /></template>
        </t-table>
      </t-card>
    </t-col>
    <t-col :xs="12" :xl="6">
      <t-card title="任务量排名" class="dashboard-rank-card rank-card" :bordered="false">
        <template #actions>
          <t-radio-group v-model="countPeriod" variant="default-filled">
            <t-radio-button value="week">本周</t-radio-button>
            <t-radio-button value="month">本月</t-radio-button>
          </t-radio-group>
        </template>
        <t-table :data="countData" :columns="columns" row-key="deviceId">
          <template #index="{ rowIndex }">
            <span :class="rankClass(rowIndex)">{{ rowIndex + 1 }}</span>
          </template>
          <template #trend="{ row }">
            <trend v-if="row.trend" :type="row.trend.direction" :describe="`${row.trend.percent}%`" />
            <span v-else>—</span>
          </template>
          <template #operation="{ row }">
            <t-link theme="primary" @click="goDetail(row.deviceId)">详情</t-link>
          </template>
          <template #empty><t-empty description="暂无任务记录" /></template>
        </t-table>
      </t-card>
    </t-col>
  </t-row>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { PrimaryTableCol } from 'tdesign-vue-next'
import Trend from '@/components/trend/index.vue'
import type { BackupTask } from '@/types'
import { rankByPeriod, type RankItem } from '../index'

const props = defineProps<{ tasks: BackupTask[] }>()
const router = useRouter()
const failedPeriod = ref<'week' | 'month'>('week')
const countPeriod = ref<'week' | 'month'>('week')

type Row = RankItem & { trend: null }

// 真实数据无逐设备历史环比时显示中性占位，不伪造涨跌。
function withTrend(items: RankItem[]): Row[] {
  return items.map((item) => ({ ...item, trend: null }))
}

const failedData = computed(() => withTrend(rankByPeriod(props.tasks, failedPeriod.value, 'failed')))
const countData = computed(() => withTrend(rankByPeriod(props.tasks, countPeriod.value, 'count')))

const columns: PrimaryTableCol<Row>[] = [
  { align: 'center', colKey: 'index', title: '排名', width: 70, fixed: 'left' },
  { align: 'left', ellipsis: true, colKey: 'deviceName', title: '设备', width: 150 },
  { align: 'center', colKey: 'trend', title: '趋势', width: 110 },
  { align: 'center', colKey: 'count', title: '次数', width: 80 },
  { align: 'center', colKey: 'operation', title: '操作', width: 80, fixed: 'right' },
]

function rankClass(index: number): (string | Record<string, boolean>)[] {
  return ['dashboard-rank__cell', { 'dashboard-rank__cell--top': index < 3 }]
}
function goDetail(id: number): void { router.push(`/devices/${id}`) }
</script>

<style scoped>
.rank-list { width: 100%; }
.dashboard-rank-card { min-height: 260px; padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl); }
.dashboard-rank-card :deep(.t-card__header) { padding: 0; }
.dashboard-rank-card :deep(.t-card__title) { font: var(--td-font-title-large); font-weight: 400; }
.dashboard-rank-card :deep(.t-card__body) { padding: 0; margin-top: var(--td-comp-margin-xxl); }
.dashboard-rank__cell {
  display: inline-flex; width: 24px; height: 24px; border-radius: 50%;
  color: white; font-size: 14px; background-color: var(--td-gray-color-5);
  align-items: center; justify-content: center; font-weight: 700;
}
.dashboard-rank__cell--top { background: var(--td-brand-color); }
</style>

<template>
  <div class="overview-wrap">
    <t-row :gutter="16">
      <t-col :xs="12" :xl="9">
        <t-card title="备份量趋势" subtitle="近7天" class="ov-card overview-card">
          <template #actions>
            <t-date-range-picker theme="primary" mode="date" :default-value="range" @change="onRangeChange" />
          </template>
          <div v-if="hasData" ref="barEl" class="dashboard-chart-container" style="width: 100%; height: 351px" />
          <t-empty v-else description="暂无备份任务数据" />
        </t-card>
      </t-col>
      <t-col :xs="12" :xl="3">
        <t-card title="本月统计" class="ov-card ov-stats overview-card">
          <template #actions>
            <t-button size="small" variant="outline" @click="exportOverview">导出</t-button>
          </template>
          <div class="stat-block inner-card">
            <div class="stat-label">本月备份次数</div>
            <div class="stat-value">{{ inner.backupCount }}</div>
            <div class="stat-foot">
              较上期
              <trend v-if="inner.countDelta" class="trend-tag" :type="inner.countDelta.direction" :describe="`${inner.countDelta.percent}%`" />
              <span v-else>—</span>
            </div>
            <t-button variant="text" theme="primary" size="small" class="stat-link" @click="navigate('DeviceList')">查看设备</t-button>
          </div>
          <t-divider class="stat-divider" />
          <div class="stat-block inner-card">
            <div class="stat-label">本月成功率</div>
            <div class="stat-value">{{ inner.successRate }}%</div>
            <div class="stat-foot">
              较上期
              <trend v-if="inner.rateDelta" class="trend-tag" :type="inner.rateDelta.direction" :describe="`${inner.rateDelta.percent}%`" />
              <span v-else>—</span>
            </div>
            <t-button variant="text" theme="primary" size="small" class="stat-link" @click="navigate('FailedDevices')">查看全部</t-button>
          </div>
        </t-card>
      </t-col>
    </t-row>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { BarChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { MessagePlugin } from 'tdesign-vue-next'
import Trend from '@/components/trend/index.vue'
import { useSettingStore } from '@/stores/setting'
import { downloadCsv } from '@/utils/csv'
import type { BackupTask } from '@/types'
import { aggregateByDay, buildInnerCards } from '../index'

echarts.use([TooltipComponent, LegendComponent, GridComponent, BarChart, CanvasRenderer])

const props = defineProps<{ tasks: BackupTask[] }>()
const emit = defineEmits<{ navigate: [target: 'DeviceList' | 'FailedDevices'] }>()
function navigate(target: 'DeviceList' | 'FailedDevices'): void { emit('navigate', target) }
const settingStore = useSettingStore()
const barEl = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null
let observer: ResizeObserver | null = null
let rangeDays = 7
function iso(d: Date): string { return d.toISOString().slice(0, 10) }
const range = ref<[string, string]>([iso(new Date(Date.now() - 6 * 86400000)), iso(new Date())])
const inner = computed(() => buildInnerCards(props.tasks))
const points = computed(() => aggregateByDay(props.tasks, rangeDays))
const hasData = computed(() => points.value.some((p) => p.total > 0))

function exportOverview(): void {
  downloadCsv(`备份概览-${new Date().toISOString().slice(0, 10)}.csv`, [
    ['指标', '数值'],
    ['本月备份次数', inner.value.backupCount],
    ['本月成功率(%)', inner.value.successRate],
    ['生成时间', new Date().toLocaleString()],
  ])
  MessagePlugin.success('已导出备份概览')
}

function option(): echarts.EChartsCoreOption {
  const c = settingStore.chartColors
  return {
    color: c.seriesColors,
    tooltip: { trigger: 'item', backgroundColor: c.containerColor, borderColor: c.borderColor, textStyle: { color: c.textColor } },
    xAxis: {
      type: 'category', data: points.value.map((p) => p.date.slice(5)),
      axisLabel: { color: c.placeholderColor }, axisLine: { lineStyle: { color: c.seriesColors[0], width: 1 } },
    },
    yAxis: { type: 'value', axisLabel: { color: c.placeholderColor }, splitLine: { lineStyle: { color: c.borderColor } } },
    grid: { top: '5%', left: '25px', right: 0, bottom: '12%' },
    legend: { icon: 'rect', itemWidth: 12, itemHeight: 4, data: ['备份量', '失败'], textStyle: { fontSize: 12, color: c.placeholderColor }, left: 'center', bottom: 0 },
    series: [
      { name: '备份量', data: points.value.map((p) => p.total), type: 'bar' },
      { name: '失败', data: points.value.map((p) => p.failed), type: 'bar' },
    ],
  }
}
function render(): void {
  if (barEl.value && hasData.value) {
    if (!chart) chart = echarts.init(barEl.value)
    chart.setOption(option(), true)
  }
}
function onRangeChange(value: unknown): void {
  const v = value as string[] | null
  if (!v || v.length < 2) return
  range.value = [v[0], v[1]]
  rangeDays = Math.min(Math.max(1, Math.round((Date.parse(v[1]) - Date.parse(v[0])) / 86400000) + 1), 365)
  render()
}
onMounted(() => {
  nextTick(() => {
    render()
    if (barEl.value) observer = new ResizeObserver(() => chart?.resize())
    if (barEl.value) observer!.observe(barEl.value)
  })
})
watch(() => settingStore.displayMode, render)
watch(points, render, { deep: true })
onBeforeUnmount(() => { observer?.disconnect(); chart?.dispose() })
</script>

<style scoped>
.overview-wrap { width: 100%; }
.overview-wrap :deep(.t-row) { align-items: stretch; }
.overview-wrap :deep(.t-col) { display: flex; }
.ov-card { flex: 1; height: auto; }
.ov-card :deep(.t-card__header) { padding: 0 0 12px; }
.ov-card :deep(.t-card__body) { padding-top: 8px; }
.ov-stats :deep(.t-card__body) { display: flex; flex-direction: column; justify-content: space-between; }
.stat-block { display: flex; flex-direction: column; gap: 4px; padding: 4px 0; }
.stat-label { color: var(--td-text-color-placeholder); font-size: 13px; }
.stat-value { font-size: 28px; font-weight: 600; color: var(--td-text-color-primary); line-height: 1.2; }
.stat-foot { display: flex; align-items: center; color: var(--td-text-color-placeholder); font-size: 12px; }
.stat-foot .trend-tag { margin-left: 6px; }
.stat-link { align-self: flex-start; padding-left: 0; margin-top: 2px; }
.stat-divider { margin: 8px 0; }
</style>

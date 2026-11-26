<template>
  <t-card :bordered="false" class="overview-wrap">
    <t-row>
      <t-col :xs="12" :xl="9">
        <t-card :bordered="false" title="备份量趋势" subtitle="近7天" class="dashboard-overview-card overview-card overview-panel">
          <template #actions>
            <t-date-range-picker theme="primary" mode="date" :default-value="range" @change="onRangeChange" />
          </template>
          <div v-if="hasData" ref="barEl" class="dashboard-chart-container" style="width: 100%; height: 351px" />
          <t-empty v-else description="暂无备份任务数据" />
        </t-card>
      </t-col>
      <t-col :xs="12" :xl="3">
        <t-card :bordered="false" class="dashboard-overview-card overview-card export-panel">
          <template #actions>
            <t-button>导出</t-button>
          </template>
          <t-row>
            <t-col :xs="6" :xl="12">
              <t-card :bordered="false" subtitle="本月备份次数" class="inner-card">
                <div class="inner-card__content">
                  <div class="inner-card__content-title">{{ inner.backupCount }}</div>
                  <div class="inner-card__content-footer">
                    较上期
                    <trend v-if="inner.countDelta" class="trend-tag" :type="inner.countDelta.direction" :describe="`${inner.countDelta.percent}%`" />
                    <span v-else>—</span>
                  </div>
                  <t-button variant="text" theme="primary" size="small" @click="navigate('DeviceList')">查看设备</t-button>
                </div>
              </t-card>
            </t-col>
            <t-col :xs="6" :xl="12">
              <t-card :bordered="false" subtitle="本月成功率" class="inner-card">
                <div class="inner-card__content">
                  <div class="inner-card__content-title">{{ inner.successRate }}%</div>
                  <div class="inner-card__content-footer">
                    较上期
                    <trend v-if="inner.rateDelta" class="trend-tag" :type="inner.rateDelta.direction" :describe="`${inner.rateDelta.percent}%`" />
                    <span v-else>—</span>
                  </div>
                  <t-button variant="text" theme="primary" size="small" @click="navigate('FailedDevices')">查看全部</t-button>
                </div>
              </t-card>
            </t-col>
          </t-row>
        </t-card>
      </t-col>
    </t-row>
  </t-card>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { BarChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import Trend from '@/components/trend/index.vue'
import { useSettingStore } from '@/stores/setting'
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
.overview-wrap :deep(.t-card__body) { padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl); }
.dashboard-overview-card :deep(.t-card__header) { padding: 0; }
.dashboard-overview-card :deep(.t-card__title) { font: var(--td-font-title-large); font-weight: 400; }
.dashboard-overview-card :deep(.t-card__body) { margin-top: var(--td-comp-margin-xxl); padding: 0; }
.dashboard-overview-card.overview-panel { border-right: none; }
.dashboard-overview-card.export-panel { border-left: none; margin-left: calc(var(--td-comp-margin-xxxl) + var(--td-comp-margin-xxxl)); }
.inner-card { margin-top: var(--td-comp-margin-s); margin-bottom: var(--td-comp-margin-xxxxl); }
.inner-card :deep(.t-card__header) { padding-bottom: 0; }
.inner-card :deep(.t-card__body) { margin-top: var(--td-comp-margin-s); }
.inner-card__content-title { font-size: var(--td-font-size-headline-medium); line-height: var(--td-line-height-headline-medium); }
.inner-card__content-footer { display: flex; align-items: center; color: var(--td-text-color-placeholder); margin-top: var(--td-comp-margin-xxl); }
.inner-card__content-footer .trend-tag { margin-left: var(--td-comp-margin-s); }
</style>

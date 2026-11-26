<template>
  <t-row :gutter="16" class="middle-chart">
    <t-col :xs="12" :xl="9">
      <t-card :title="'备份趋势'" :subtitle="currentRange" class="dashboard-chart-card" :bordered="false">
        <template #actions>
          <t-date-range-picker
            class="card-date-picker-container"
            theme="primary"
            mode="date"
            :default-value="last7"
            @change="onRangeChange"
          />
        </template>
        <div v-if="hasTrend" ref="lineEl" class="dashboard-chart-container" :style="{ width: '100%', height: `${resizeTime * 326}px` }" />
        <t-empty v-else description="暂无备份任务数据" />
      </t-card>
    </t-col>
    <t-col :xs="12" :xl="3">
      <t-card :title="'状态分布'" :subtitle="currentRange" class="dashboard-chart-card" :bordered="false">
        <div v-if="hasPie" ref="pieEl" class="dashboard-chart-container" :style="{ width: '100%', height: `${resizeTime * 326}px` }" />
        <t-empty v-else description="暂无状态数据" />
      </t-card>
    </t-col>
  </t-row>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { BarChart, PieChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { useSettingStore } from '@/stores/setting'
import type { BackupTask } from '@/types'
import { aggregateByDay, statusDistribution } from '../index'

echarts.use([TooltipComponent, LegendComponent, PieChart, GridComponent, BarChart, CanvasRenderer])

const props = defineProps<{ tasks: BackupTask[] }>()
const settingStore = useSettingStore()
const resizeTime = ref(1)
const lineEl = ref<HTMLElement | null>(null)
const pieEl = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null
let observer: ResizeObserver | null = null
let rangeDays = 7

function iso(d: Date): string { return d.toISOString().slice(0, 10) }
const last7 = ref<[string, string]>([
  iso(new Date(Date.now() - 6 * 86400000)),
  iso(new Date()),
])
const currentRange = computed(() => `${last7.value[0]} 至 ${last7.value[1]}`)

const lineData = computed(() => aggregateByDay(props.tasks, rangeDays))
const hasTrend = computed(() => lineData.value.some((p) => p.total > 0))
const pieSlices = computed(() => statusDistribution(props.tasks))
const hasPie = computed(() => pieSlices.value.length > 0)

function barOption(): echarts.EChartsCoreOption {
  const c = settingStore.chartColors
  const days = lineData.value.map((p) => p.date.slice(5))
  // shadcn Overview：主色填充的圆角柱状图（primary 近黑/近白，随主题）
  const primary = c.primaryColor
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: c.containerColor,
      borderColor: c.borderColor,
      textStyle: { color: c.textColor },
    },
    grid: { left: 8, right: 16, top: 12, bottom: 8, containLabel: true },
    xAxis: {
      type: 'category',
      data: days,
      axisTick: { show: false },
      axisLine: { lineStyle: { color: c.borderColor } },
      axisLabel: { color: c.placeholderColor, fontSize: 12 },
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: c.placeholderColor, fontSize: 12 },
      splitLine: { lineStyle: { color: c.borderColor } },
    },
    series: [{
      name: '备份量',
      type: 'bar',
      data: lineData.value.map((p) => p.total),
      barMaxWidth: 28,
      itemStyle: { color: primary, borderRadius: [4, 4, 0, 0] },
    }],
  }
}

function pieOption(): echarts.EChartsCoreOption {
  const c = settingStore.chartColors
  const total = pieSlices.value.reduce((a, s) => a + s.value, 0)
  // starter 中心展示「最大通道占比」；这里映射为最大状态占比，避免无成功任务时误标「成功占比 100%」。
  const top = [...pieSlices.value].sort((a, b) => b.value - a.value)[0]
  const topPercent = total > 0 && top ? Math.round((top.value / total) * 100) : 0
  const centerName = top ? `${top.name}占比` : '状态占比'
  return {
    color: c.seriesColors,
    tooltip: { show: false },
    legend: { data: pieSlices.value.map((s) => s.name), textStyle: { color: c.placeholderColor }, left: 'center', bottom: 0, itemWidth: 12, itemHeight: 4 },
    series: [{
      name: '状态分布',
      type: 'pie',
      radius: ['48%', '60%'],
      avoidLabelOverlap: true,
      silent: true,
      itemStyle: { borderColor: c.containerColor, borderWidth: 1 },
      label: {
        show: true,
        position: 'center',
        formatter: [`{value|${topPercent}%}`, `{name|${centerName}}`].join('\n'),
        rich: {
          value: { color: c.textColor, fontSize: 28, fontWeight: 'normal', lineHeight: 46 },
          name: { color: c.placeholderColor, fontSize: 12, lineHeight: 14 },
        },
      },
      labelLine: { show: false },
      data: pieSlices.value,
    }],
  }
}

function render(): void {
  if (lineEl.value && hasTrend.value) {
    if (!lineChart) lineChart = echarts.init(lineEl.value)
    lineChart.setOption(barOption(), true)
  }
  if (pieEl.value && hasPie.value) {
    if (!pieChart) pieChart = echarts.init(pieEl.value)
    pieChart.setOption(pieOption(), true)
  }
}

function onRangeChange(value: unknown): void {
  const v = value as string[] | null
  if (!v || v.length < 2) return
  last7.value = [v[0], v[1]]
  const days = Math.max(1, Math.round((Date.parse(v[1]) - Date.parse(v[0])) / 86400000) + 1)
  rangeDays = Math.min(days, 365)
  render()
}

function updateSize(): void {
  const w = document.documentElement.clientWidth
  resizeTime.value = w >= 1400 && w < 1920 ? Number((w / 2080).toFixed(2)) : w < 1080 ? Number((w / 1080).toFixed(2)) : 1
  lineChart?.resize(); pieChart?.resize()
}

onMounted(() => {
  nextTick(() => {
    render()
    updateSize()
    observer = new ResizeObserver(() => { lineChart?.resize(); pieChart?.resize() })
    if (lineEl.value) observer.observe(lineEl.value)
    if (pieEl.value) observer.observe(pieEl.value)
  })
})
watch(() => settingStore.displayMode, render)
watch([lineData, pieSlices], render, { deep: true })
onBeforeUnmount(() => {
  observer?.disconnect()
  lineChart?.dispose(); pieChart?.dispose()
})
</script>

<style scoped>
.middle-chart { width: 100%; }
.dashboard-chart-card { padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl); }
.dashboard-chart-card :deep(.t-card__header) { padding: 0; }
.dashboard-chart-card :deep(.t-card__body) { padding: 0; margin-top: var(--td-comp-margin-xxl); }
.dashboard-chart-card :deep(.t-card__title) { font: var(--td-font-title-large); font-weight: 400; }
.dashboard-chart-container { display: block; }
</style>

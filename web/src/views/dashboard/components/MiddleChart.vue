<template>
  <t-card title="备份趋势" :bordered="false" class="trend-card">
    <div v-if="hasData" ref="chartEl" class="trend-chart" data-testid="trend-chart" />
    <t-empty v-else description="暂无备份任务数据" />
  </t-card>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { useSettingStore } from '@/stores/setting'
import type { BackupTask } from '@/types'
import { aggregateByDay } from '../index'

echarts.use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{
  tasks: BackupTask[]
}>()

const settingStore = useSettingStore()
const chartEl = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null
let observer: ResizeObserver | null = null

const points = computed(() => aggregateByDay(props.tasks, 7))
const hasData = computed(() => points.value.some((p) => p.success > 0 || p.failed > 0))

function buildOption(): echarts.EChartsCoreOption {
  const colors = settingStore.chartColors
  return {
    color: colors.seriesColors,
    textStyle: { color: colors.textColor },
    tooltip: {
      trigger: 'axis',
      backgroundColor: colors.containerColor,
      borderColor: colors.borderColor,
      textStyle: { color: colors.textColor },
    },
    legend: {
      data: ['成功', '失败'],
      textStyle: { color: colors.textColor },
    },
    grid: { left: 40, right: 20, top: 48, bottom: 32 },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: points.value.map((p) => p.date),
      axisLabel: { color: colors.placeholderColor },
      axisLine: { lineStyle: { color: colors.borderColor } },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: colors.placeholderColor },
      axisLine: { lineStyle: { color: colors.borderColor } },
      splitLine: { lineStyle: { color: colors.borderColor, type: 'dashed' } },
    },
    series: [
      {
        name: '成功',
        type: 'line',
        smooth: true,
        symbolSize: 6,
        data: points.value.map((p) => p.success),
      },
      {
        name: '失败',
        type: 'line',
        smooth: true,
        symbolSize: 6,
        data: points.value.map((p) => p.failed),
      },
    ],
  }
}

function initChart() {
  if (chart || !chartEl.value) return
  chart = echarts.init(chartEl.value)
  chart.setOption(buildOption(), true)
  if (typeof ResizeObserver !== 'undefined') {
    observer = new ResizeObserver(() => chart?.resize())
    observer.observe(chartEl.value)
  }
}

onMounted(() => {
  initChart()
})

// 直接监听配色来源：主题切换与 auto 系统偏好变化都会更新 chartColors，触发重绘。
watch(() => settingStore.chartColors, () => {
  chart?.setOption(buildOption(), true)
})

watch([hasData, () => props.tasks], async () => {
  await nextTick()
  if (!chart) {
    initChart()
  } else {
    chart.setOption(buildOption(), true)
  }
})

onUnmounted(() => {
  observer?.disconnect()
  observer = null
  chart?.dispose()
  chart = null
})
</script>

<style scoped>
.trend-card {
  border-radius: var(--td-radius-medium);
}

.trend-chart {
  width: 100%;
  height: 326px;
}

@media (max-width: 768px) {
  .trend-chart {
    height: 260px;
  }
}
</style>

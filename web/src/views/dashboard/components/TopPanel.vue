<template>
  <t-row :gutter="[16, 16]" class="top-panel">
    <t-col v-for="(card, index) in cards" :key="card.label" :xs="6" :xl="3">
      <div
        class="metric-card dashboard-item"
        :class="{ 'metric-card--primary dashboard-item--main-color': index === 0 }"
        role="link"
        tabindex="0"
        :aria-label="`${card.label}：${card.value}${card.suffix || ''}`"
        @click="emit('navigate', card.route)"
        @keydown.enter="emit('navigate', card.route)"
      >
        <div class="kpi-head">
          <span class="kpi-title">{{ card.label }}</span>
          <span class="kpi-icon"><component :is="card.icon" /></span>
        </div>
        <div class="kpi-value">{{ card.value }}<small v-if="card.suffix">{{ card.suffix }}</small></div>
        <div class="metric-footnote kpi-foot">{{ card.compareLabel }}</div>
      </div>
    </t-col>
  </t-row>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import {
  ServerIcon,
  CheckCircleIcon,
  ErrorCircleIcon,
  NotificationIcon,
} from 'tdesign-icons-vue-next'
import type { DashboardStats } from '../index'

const props = defineProps<{
  stats: DashboardStats
  tasks: unknown
}>()
const emit = defineEmits<{ navigate: [target: 'DeviceList' | 'FailedDevices' | 'AlertList'] }>()

const cards = computed<
  Array<{
    label: string
    value: number
    suffix?: string
    compareLabel: string
    icon: Component
    route: 'DeviceList' | 'FailedDevices' | 'AlertList'
  }>
>(() => {
  const s = props.stats
  return [
    { label: '设备总数', value: s.totalDevices, compareLabel: '近7天', icon: ServerIcon, route: 'DeviceList' },
    { label: '备份成功率', value: s.successRate, suffix: '%', compareLabel: '近30天', icon: CheckCircleIcon, route: 'DeviceList' },
    { label: '失败设备数', value: s.failedDevices, compareLabel: '近7天', icon: ErrorCircleIcon, route: 'FailedDevices' },
    { label: '未读告警', value: s.unreadAlerts, compareLabel: '当前', icon: NotificationIcon, route: 'AlertList' },
  ]
})
</script>

<style scoped>
.top-panel { width: 100%; }

/* shadcn Card：rounded-xl 边框 shadow-sm，白底，p-6 */
.metric-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-height: 118px;
  padding: 22px 24px;
  border: 1px solid var(--shadcn-border);
  border-radius: 12px;
  background: var(--shadcn-card);
  box-shadow: var(--shadcn-card-shadow);
  cursor: pointer;
}
.metric-card:hover { border-color: var(--shadcn-ring); }

.kpi-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.kpi-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--shadcn-muted-foreground);
}
.kpi-icon {
  display: inline-flex;
  align-items: center;
  color: var(--shadcn-muted-foreground);
  font-size: 16px;
}
.kpi-value {
  margin-top: 6px;
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--shadcn-foreground);
}
.kpi-value small { font-size: 14px; font-weight: 600; margin-left: 2px; }
.kpi-foot {
  margin-top: 2px;
  font-size: 12px;
  color: var(--shadcn-muted-foreground);
}

@media (max-width: 768px) { .metric-card { min-height: 100px; padding: 18px; } }
</style>

<template>
  <t-row :gutter="[16, 16]">
    <t-col v-for="card in cards" :key="card.label" :xs="12" :sm="6" :xl="3">
      <t-card
        class="metric-card"
        :class="{ 'metric-card--primary': card.primary }"
        role="link"
        tabindex="0"
        :aria-label="`${card.label}：${card.value}${card.suffix || ''}，${card.footnote}`"
        @click="emit('navigate', card.route)"
        @keydown.enter="emit('navigate', card.route)"
      >
        <div class="metric-body">
          <span class="metric-icon" :class="`metric-icon--${card.tone}`">
            <component :is="card.icon" />
          </span>
          <div>
            <div class="metric-label">{{ card.label }}</div>
            <div class="metric-value">
              {{ card.value }}<span v-if="card.suffix" class="metric-suffix">{{ card.suffix }}</span>
            </div>
          </div>
        </div>
        <div class="metric-footnote">{{ card.footnote }}</div>
      </t-card>
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
}>()

type NavigateTarget = 'DeviceList' | 'FailedDevices' | 'AlertList'

const emit = defineEmits<{
  navigate: [target: NavigateTarget]
}>()

const cards = computed<
  Array<{
    label: string
    value: number
    suffix?: string
    footnote: string
    icon: Component
    tone: string
    route: NavigateTarget
    primary?: boolean
  }>
>(() => [
  {
    label: '设备总数',
    value: props.stats.totalDevices,
    footnote: '已纳管网络设备',
    icon: ServerIcon,
    tone: 'blue',
    route: 'DeviceList',
    primary: true,
  },
  {
    label: '备份成功率',
    value: props.stats.successRate,
    suffix: '%',
    footnote: '基于全部备份任务',
    icon: CheckCircleIcon,
    tone: 'green',
    route: 'DeviceList',
  },
  {
    label: '失败设备数',
    value: props.stats.failedDevices,
    footnote: '需要尽快处理',
    icon: ErrorCircleIcon,
    tone: 'red',
    route: 'FailedDevices',
  },
  {
    label: '未读告警',
    value: props.stats.unreadAlerts,
    footnote: '待确认告警信息',
    icon: NotificationIcon,
    tone: 'amber',
    route: 'AlertList',
  },
])
</script>

<style scoped>
.metric-card {
  min-height: 148px;
  cursor: pointer;
  margin-bottom: 0;
  border-color: var(--td-component-stroke);
}

.metric-card :deep(.t-card__body) {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 112px;
}

.metric-body {
  display: flex;
  align-items: center;
  gap: 12px;
}

.metric-icon {
  width: 42px;
  height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border-radius: 50%;
  font-size: 22px;
}

.metric-icon--blue {
  color: var(--td-brand-color);
  background: var(--td-brand-color-light);
}

.metric-icon--green {
  color: var(--td-success-color-5);
  background: var(--td-success-color-1);
}

.metric-icon--red {
  color: var(--td-error-color-5);
  background: var(--td-error-color-1);
}

.metric-icon--amber {
  color: var(--td-warning-color-5);
  background: var(--td-warning-color-1);
}

.metric-label {
  color: var(--td-text-color-secondary);
  font-size: 14px;
  font-weight: 600;
}

.metric-value {
  color: var(--td-text-color-primary);
  font-size: 30px;
  font-weight: 700;
  line-height: 1.2;
}

.metric-suffix {
  font-size: 18px;
  margin-left: 2px;
}

.metric-footnote {
  margin-top: 12px;
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}

/* 使用 brand-7：浅色与 --td-brand-color 等价；深色下比 brand-8 更深，
   使 12px 白色脚注对比度由 3.76:1 提升到 5.29:1（≥ WCAG AA 4.5:1）。 */
.metric-card--primary {
  background: var(--td-brand-color-7);
  border-color: var(--td-brand-color-7);
}

.metric-card--primary .metric-label,
.metric-card--primary .metric-value,
.metric-card--primary .metric-footnote {
  color: var(--td-text-color-anti);
}

.metric-card--primary .metric-icon {
  color: var(--td-brand-color);
  background: var(--td-text-color-anti);
}

@media (max-width: 768px) {
  .metric-card {
    min-height: 120px;
  }
}
</style>

<template>
  <t-row :gutter="16">
    <t-col :xs="12" :xl="8">
      <t-card title="最近备份任务" :bordered="false" class="overview-card">
        <template #actions>
          <t-button variant="text" theme="primary" @click="emit('navigate', 'DeviceList')">查看设备</t-button>
        </template>
        <t-table
          :columns="taskColumns"
          :data="recentTasks"
          row-key="id"
          :pagination="false"
          size="small"
          table-layout="auto"
        >
          <template #empty><t-empty description="暂无备份任务" /></template>
        </t-table>
      </t-card>
    </t-col>
    <t-col :xs="12" :xl="4">
      <t-card title="失败设备概览" :bordered="false" class="overview-card">
        <template #actions>
          <t-button variant="text" theme="primary" @click="emit('navigate', 'FailedDevices')">查看全部</t-button>
        </template>
        <t-table
          :columns="failedColumns"
          :data="failedDevices"
          row-key="id"
          :pagination="false"
          size="small"
          table-layout="auto"
        >
          <template #empty><t-empty description="暂无失败设备" /></template>
        </t-table>
      </t-card>
    </t-col>
  </t-row>
</template>

<script setup lang="ts">
import { Button, type PrimaryTableCol } from 'tdesign-vue-next'
import BackupStatus from '@/components/BackupStatus.vue'
import type { FailedDeviceRow, RecentTaskRow } from '../index'

defineProps<{
  recentTasks: RecentTaskRow[]
  failedDevices: FailedDeviceRow[]
}>()

const emit = defineEmits<{
  navigate: [target: 'DeviceList' | 'FailedDevices']
  detail: [id: number]
}>()

const taskColumns: PrimaryTableCol<RecentTaskRow>[] = [
  { title: '设备', colKey: 'device_name', ellipsis: true },
  {
    title: '状态',
    colKey: 'status',
    width: 100,
    cell: (h, { row }) => h(BackupStatus, { status: row.status }),
  },
  {
    title: '完成时间',
    colKey: 'finished_at',
    width: 170,
    cell: (h, { row }) => h('span', { class: 'datetime-cell', title: row.finished_at }, row.finished_at),
  },
]

const failedColumns: PrimaryTableCol<FailedDeviceRow>[] = [
  { title: '设备', colKey: 'name', ellipsis: true },
  { title: 'IP', colKey: 'ip', width: 130 },
  { title: '原因', colKey: 'reason', ellipsis: true },
  {
    title: '操作',
    colKey: 'actions',
    width: 80,
    cell: (h, { row }) =>
      h(
        Button,
        { size: 'small', variant: 'text', theme: 'primary', onClick: () => emit('detail', row.id) },
        { default: () => '详情' },
      ),
  },
]

// 失败设备表格沿用原有列语义，操作列复用现有详情跳转
</script>

<style scoped>
.overview-card {
  min-height: 280px;
  border-radius: var(--td-radius-medium);
}
</style>

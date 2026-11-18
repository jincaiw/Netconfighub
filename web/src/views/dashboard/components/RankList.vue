<template>
  <t-row :gutter="16">
    <t-col :xs="12" :xl="6">
      <t-card title="失败次数 TOP5" :bordered="false" class="rank-card">
        <t-table
          :columns="columns"
          :data="failedRank"
          row-key="deviceId"
          :pagination="false"
          size="small"
        >
          <template #empty><t-empty description="暂无失败记录" /></template>
        </t-table>
      </t-card>
    </t-col>
    <t-col :xs="12" :xl="6">
      <t-card title="任务量 TOP5" :bordered="false" class="rank-card">
        <t-table
          :columns="columns"
          :data="countRank"
          row-key="deviceId"
          :pagination="false"
          size="small"
        >
          <template #empty><t-empty description="暂无任务记录" /></template>
        </t-table>
      </t-card>
    </t-col>
  </t-row>
</template>

<script setup lang="ts">
import { h } from 'vue'
import type { PrimaryTableCol } from 'tdesign-vue-next'
import type { RankItem } from '../index'

defineProps<{
  failedRank: RankItem[]
  countRank: RankItem[]
}>()

const columns: PrimaryTableCol<RankItem>[] = [
  {
    title: '排名',
    colKey: 'index',
    width: 70,
    cell: (_h, { rowIndex }) =>
      h('span', { class: ['rank-index', { 'rank-index--top': rowIndex < 3 }] }, String(rowIndex + 1)),
  },
  { title: '设备', colKey: 'deviceName', ellipsis: true },
  { title: '次数', colKey: 'count', width: 90 },
]
</script>

<style scoped>
.rank-card {
  min-height: 240px;
  border-radius: var(--td-radius-medium);
}

.rank-index {
  display: inline-flex;
  width: 24px;
  height: 24px;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--td-text-color-anti);
  background-color: var(--td-gray-color-5);
  font-size: 13px;
  font-weight: 700;
}

.rank-index--top {
  background-color: var(--td-brand-color);
}
</style>

<template>
  <div class="device-table-wrap" aria-label="设备列表">
    <t-table
      :data="devices"
      :columns="columns"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      size="small"
      table-layout="fixed"
      @page-change="(p: any) => emit('page-change', p)"
    >
      <template #name="{ row }">{{ row.name }}</template>
      <template #vendor="{ row }">{{ vendorLabelMap[row.vendor] || row.vendor }}</template>
      <template #model="{ row }">{{ modelLabelMap[row.model] || row.model }}</template>
      <template #group_name="{ row }">{{ row.group_name || '-' }}</template>
      <template #collection_mode="{ row }">
        <t-tag v-if="row.collection_mode === 'sftp_push'" theme="primary" variant="light" size="small">
          SFTP 主动上传
        </t-tag>
        <t-tag v-else variant="outline" size="small">平台拉取</t-tag>
      </template>
      <template #enabled="{ row }">
        <t-switch
          :model-value="row.enabled"
          :aria-label="`${row.enabled ? '禁用' : '启用'}设备 ${row.name}`"
          @change="() => emit('toggle', row)"
        />
      </template>
      <template #last_backup_status="{ row }">
        <BackupStatus :status="row.last_backup_status" />
      </template>
      <template #last_backup_at="{ row }">
        <span class="datetime-cell" :title="row.last_backup_at || ''">{{ formatDateTime(row.last_backup_at) }}</span>
      </template>
      <template #actions="{ row }">
        <t-space size="small" class="device-actions">
          <t-button size="small" variant="outline" @click="emit('detail', row.id)">详情</t-button>
          <t-button size="small" variant="outline" @click="emit('edit', row)">编辑</t-button>
          <t-button
            size="small"
            theme="primary"
            :disabled="row.collection_mode === 'sftp_push'"
            :title="row.collection_mode === 'sftp_push' ? '该设备使用 SFTP 主动上传，平台不会拉取' : ''"
            @click="emit('backup', row.id)"
          >
            备份
          </t-button>
          <t-button size="small" theme="danger" variant="outline" @click="emit('delete', row)">删除</t-button>
        </t-space>
      </template>
      <template #empty><span class="empty-state-text">暂无数据</span></template>
    </t-table>
  </div>
</template>

<script setup lang="ts">
import { type PrimaryTableCol } from 'tdesign-vue-next'
import BackupStatus from '@/components/BackupStatus.vue'
import { formatDateTime } from '@/utils/datetime'
import type { Device } from '@/types'

defineProps<{
  devices: Device[]
  loading: boolean
  vendorLabelMap: Record<string, string>
  modelLabelMap: Record<string, string>
  pagination?: Record<string, any>
}>()

const emit = defineEmits<{
  (e: 'toggle', device: Device): void
  (e: 'detail', id: number): void
  (e: 'edit', device: Device): void
  (e: 'backup', id: number): void
  (e: 'delete', device: Device): void
  (e: 'page-change', payload: { current: number; pageSize: number }): void
}>()

// 保留原生 table/tr 渲染（E2E 依赖 page.locator('table') / page.locator('tr')），
// 首列固定 + 横向滚动替代原 min-width:1080px，桌面结构不变，窄屏不撑破页面。
const columns: PrimaryTableCol<Device>[] = [
  { title: '名称', colKey: 'name', width: 120, fixed: 'left' },
  { title: 'IP 地址', colKey: 'ip', width: 120 },
  { title: '厂商', colKey: 'vendor', width: 70 },
  { title: '型号', colKey: 'model', width: 80 },
  { title: '分组', colKey: 'group_name', width: 90 },
  { title: '采集方式', colKey: 'collection_mode', width: 110 },
  { title: '启用状态', colKey: 'enabled', width: 85 },
  { title: '备份状态', colKey: 'last_backup_status', width: 95 },
  { title: '最近备份', colKey: 'last_backup_at', width: 140 },
  { title: '操作', colKey: 'actions', width: 220 },
]
</script>

<style scoped>
.device-table-wrap {
  width: 100%;
  overflow-x: auto;
}

.device-actions {
  flex-wrap: nowrap;
}

:deep(.t-table__th),
:deep(.t-table__td) {
  white-space: nowrap;
}

@media (max-width: 768px) {
  .device-actions {
    flex-wrap: wrap;
  }
}
</style>

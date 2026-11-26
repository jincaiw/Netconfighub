<template>
  <div class="reports-panel">
    <t-card title="报表汇总" :bordered="false">
      <div class="report-summary">
        <div class="rs-item"><span class="rs-label">设备总数</span><b class="rs-value">{{ stats.totalDevices }}</b></div>
        <div class="rs-item"><span class="rs-label">备份成功率</span><b class="rs-value">{{ stats.successRate }}%</b></div>
        <div class="rs-item"><span class="rs-label">失败设备数</span><b class="rs-value">{{ stats.failedDevices }}</b></div>
        <div class="rs-item"><span class="rs-label">近 30 天备份任务</span><b class="rs-value">{{ tasks.length }}</b></div>
      </div>
      <t-space class="report-actions" :size="12">
        <t-button variant="outline" :loading="exportingDevices" @click="exportDevices">导出设备清单（CSV）</t-button>
        <t-button variant="outline" :loading="exportingBackups" @click="exportBackups">导出备份记录（CSV）</t-button>
      </t-space>
      <p v-if="exportError" class="report-error" role="alert">{{ exportError }}</p>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { getDevices } from '@/api/device'
import { getBackupTasks } from '@/api/backup'
import { downloadCsv } from '@/utils/csv'
import type { BackupTask } from '@/types'
import type { DashboardStats } from '../index'

defineProps<{ stats: DashboardStats; tasks: BackupTask[] }>()

const exportingDevices = ref(false)
const exportingBackups = ref(false)
const exportError = ref('')

const vendorLabel: Record<string, string> = { cisco: '思科', h3c: '华三', huawei: '华为', ruijie: '锐捷' }

async function exportDevices() {
  exportingDevices.value = true
  exportError.value = ''
  try {
    const res = await getDevices({ page: 1, page_size: 1000 })
    const items = res.data.items || []
    const rows: (string | number)[][] = [['名称', 'IP 地址', '厂商', '型号', '协议', '端口', '分组', '启用', '备份状态', '最近备份']]
    items.forEach((d) => {
      rows.push([
        d.name, d.ip, vendorLabel[d.vendor] || d.vendor, d.model, d.protocol.toUpperCase(), d.port,
        d.group_name || '', d.enabled ? '启用' : '禁用', d.last_backup_status || '', d.last_backup_at || '',
      ])
    })
    downloadCsv(`设备清单-${new Date().toISOString().slice(0, 10)}.csv`, rows)
    MessagePlugin.success(`已导出 ${items.length} 台设备`)
  } catch {
    exportError.value = '导出设备清单失败，请重试'
    MessagePlugin.error('导出设备清单失败')
  } finally {
    exportingDevices.value = false
  }
}

async function exportBackups() {
  exportingBackups.value = true
  exportError.value = ''
  try {
    const res = await getBackupTasks({ page: 1, page_size: 1000 })
    const items: BackupTask[] = res.data.items || []
    const rows: (string | number)[][] = [['任务ID', '设备', '状态', '开始时间', '结束时间', '重试次数', '失败原因', '创建时间']]
    items.forEach((t) => {
      rows.push([t.id, t.device_name || t.device_id, t.status, t.started_at || '', t.finished_at || '', t.retry_count, t.failure_reason || '', t.created_at])
    })
    downloadCsv(`备份记录-${new Date().toISOString().slice(0, 10)}.csv`, rows)
    MessagePlugin.success(`已导出 ${items.length} 条备份记录`)
  } catch {
    exportError.value = '导出备份记录失败，请重试'
    MessagePlugin.error('导出备份记录失败')
  } finally {
    exportingBackups.value = false
  }
}
</script>

<style scoped>
.report-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}
.rs-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px;
  border: 1px solid var(--shadcn-border);
  border-radius: 10px;
  background: var(--shadcn-card);
}
.rs-label { font-size: 13px; color: var(--shadcn-muted-foreground); }
.rs-value { font-size: 22px; font-weight: 700; color: var(--shadcn-foreground); }
.report-actions { flex-wrap: wrap; }
.report-error { margin-top: 10px; color: var(--td-error-color-8); font-size: 13px; }
@media (max-width: 768px) {
  .report-actions { width: 100%; }
  .report-actions :deep(.t-button) { width: 100%; }
}
</style>

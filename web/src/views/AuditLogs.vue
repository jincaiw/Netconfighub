<template>
  <div class="audit-logs">
    <n-card title="审计日志">
      <n-data-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="paginationReactive"
        :row-key="(row: AuditLog) => row.id"
        remote
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  NCard,
  NDataTable,
  useMessage,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAuditLogs } from '@/api/audit'
import type { AuditLog } from '@/types'

const message = useMessage()
const loading = ref(false)
const logs = ref<AuditLog[]>([])

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

const columns: DataTableColumns<AuditLog> = [
  { title: '用户', key: 'username', width: 120 },
  { title: '操作', key: 'action', width: 140 },
  { title: '目标类型', key: 'target_type', width: 100 },
  { title: '目标 ID', key: 'target_id', width: 90 },
  { title: '详情', key: 'detail', ellipsis: { tooltip: true } },
  { title: '客户端 IP', key: 'client_ip', width: 130 },
  { title: '时间', key: 'created_at', width: 170 },
]

async function fetchLogs() {
  loading.value = true
  try {
    const res = await getAuditLogs({
      page: paginationReactive.page,
      page_size: paginationReactive.pageSize,
    })
    logs.value = res.data.items
    paginationReactive.itemCount = res.data.total
  } catch {
    message.error('获取审计日志失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  paginationReactive.page = page
  fetchLogs()
}

function handlePageSizeChange(pageSize: number) {
  paginationReactive.pageSize = pageSize
  paginationReactive.page = 1
  fetchLogs()
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.audit-logs {
  width: 100%;
}
</style>

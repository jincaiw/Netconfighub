<template>
  <div class="audit-logs">
    <t-card title="审计日志">
      <div v-if="loadError" class="list-error" role="alert">
        <span>审计日志加载失败，请重新加载</span>
        <t-button size="small" variant="text" theme="primary" @click="fetchLogs">重新加载</t-button>
      </div>
      <t-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="loading ? undefined : tablePagination"
        row-key="id"
        @page-change="handlePageChange"
      >
        <template #empty><span class="empty-state-text">暂无数据</span></template>
      </t-table>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import { getAuditLogs } from '@/api/audit'
import { formatDateTime } from '@/utils/datetime'
import type { AuditLog } from '@/types'

const loading = ref(false)
const loadError = ref(false)
const logs = ref<AuditLog[]>([])

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
})

const tablePagination = computed(() => ({
  current: paginationReactive.page,
  pageSize: paginationReactive.pageSize,
  total: paginationReactive.itemCount,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50],
}))

const columns: PrimaryTableCol<AuditLog>[] = [
  { title: '用户', colKey: 'username', width: 120 },
  { title: '操作', colKey: 'action', width: 140 },
  { title: '目标类型', colKey: 'target_type', width: 100 },
  { title: '目标 ID', colKey: 'target_id', width: 90 },
  { title: '详情', colKey: 'detail', ellipsis: true },
  { title: '客户端 IP', colKey: 'client_ip', width: 130 },
  {
    title: '时间',
    colKey: 'created_at',
    width: 170,
    cell: (h, { row }) =>
      h('span', { class: 'datetime-cell', title: row.created_at ?? '' }, formatDateTime(row.created_at)),
  },
]

async function fetchLogs() {
  loading.value = true
  loadError.value = false
  try {
    const res = await getAuditLogs({
      page: paginationReactive.page,
      page_size: paginationReactive.pageSize,
    })
    logs.value = res.data.items
    paginationReactive.itemCount = res.data.total
  } catch {
    loadError.value = true
    MessagePlugin.error('获取审计日志失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(pageInfo: { current: number; pageSize: number }) {
  if (pageInfo.pageSize && pageInfo.pageSize !== paginationReactive.pageSize) {
    paginationReactive.pageSize = pageInfo.pageSize
    paginationReactive.page = 1
  } else {
    paginationReactive.page = pageInfo.current
  }
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

.list-error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 12px;
  border: 1px solid var(--td-error-color-2);
  border-radius: var(--td-radius-default);
  color: var(--td-error-color-8);
  background: var(--td-error-color-1);
  font-size: 13px;
}
</style>

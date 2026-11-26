<template>
  <div class="notifications-panel">
    <t-card title="最近通知" :bordered="false">
      <template #actions>
        <t-button v-if="alerts.length > 0" size="small" variant="text" theme="primary" :loading="markingAll" @click="markAll">
          全部标记已读
        </t-button>
      </template>

      <div v-if="loading" class="notif-state"><t-loading text="通知加载中..." /></div>
      <div v-else-if="error" class="notif-state notif-error" role="alert">
        <span>通知加载失败</span>
        <t-button size="small" variant="text" theme="primary" @click="fetchAlerts">重试</t-button>
      </div>
      <div v-else-if="alerts.length === 0" class="notif-state">
        <span class="empty-state-text">暂无通知</span>
      </div>
      <ul v-else class="notif-list">
        <li v-for="a in alerts" :key="a.id" class="notif-item" :class="{ 'notif-unread': !a.is_read }">
          <span v-if="!a.is_read" class="notif-dot" aria-hidden="true" />
          <div class="notif-body">
            <div class="notif-title">{{ a.title }}</div>
            <div class="notif-msg">{{ a.message }}</div>
            <div class="notif-time">{{ formatDateTime(a.created_at) }}</div>
          </div>
        </li>
      </ul>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { getAlerts, markAllAlertsAsRead } from '@/api/alert'
import { formatDateTime } from '@/utils/datetime'
import type { Alert } from '@/types'

const alerts = ref<Alert[]>([])
const loading = ref(false)
const error = ref(false)
const markingAll = ref(false)
const loaded = ref(false)

const props = defineProps<{ active?: boolean }>()

async function fetchAlerts() {
  loading.value = true
  error.value = false
  try {
    const res = await getAlerts({ page: 1, page_size: 20 })
    alerts.value = res.data.items || []
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

async function markAll() {
  markingAll.value = true
  try {
    await markAllAlertsAsRead()
    alerts.value.forEach((a) => { a.is_read = true })
    MessagePlugin.success('已全部标记为已读')
  } catch {
    MessagePlugin.error('标记已读失败')
  } finally {
    markingAll.value = false
  }
}

watch(
  () => props.active,
  (v) => { if (v && !loaded.value) { loaded.value = true; fetchAlerts() } },
  { immediate: true },
)
</script>

<style scoped>
.notif-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 32px 0;
}
.notif-error { color: var(--td-error-color-8); font-size: 13px; }
.notif-list { list-style: none; margin: 0; padding: 0; }
.notif-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 4px;
  border-bottom: 1px solid var(--shadcn-border);
}
.notif-item:last-child { border-bottom: none; }
.notif-dot {
  width: 8px; height: 8px; margin-top: 6px; border-radius: 50%;
  background: var(--td-error-color); flex-shrink: 0;
}
.notif-body { min-width: 0; }
.notif-title { font-size: 14px; font-weight: 600; color: var(--shadcn-foreground); }
.notif-msg { font-size: 13px; color: var(--shadcn-muted-foreground); margin-top: 2px; }
.notif-time { font-size: 12px; color: var(--shadcn-muted-foreground); margin-top: 4px; }
</style>

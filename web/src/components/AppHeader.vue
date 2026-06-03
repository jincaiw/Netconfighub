<template>
  <div class="app-header">
    <div class="header-left">
      <n-button quaternary @click="appStore.toggleSidebar">
        <template #icon>
          <n-icon><MenuOutline /></n-icon>
        </template>
      </n-button>
      <n-breadcrumb>
        <n-breadcrumb-item v-for="item in breadcrumbItems" :key="item.label">
          {{ item.label }}
        </n-breadcrumb-item>
      </n-breadcrumb>
    </div>
    <div class="header-right">
      <n-space align="center" :size="12">
        <n-popover trigger="click" placement="bottom-end" :width="380">
          <template #trigger>
            <n-badge :value="unreadCount" :max="99" :show="unreadCount > 0">
              <n-button quaternary size="small">
                <template #icon>
                  <n-icon size="20"><NotificationsOutline /></n-icon>
                </template>
              </n-button>
            </n-badge>
          </template>
          <div class="alert-panel">
            <div class="alert-panel-header">
              <span class="alert-panel-title">告警通知</span>
              <n-button v-if="unreadCount > 0" text type="primary" size="tiny" @click="handleMarkAllRead">
                全部已读
              </n-button>
            </div>
            <n-divider style="margin: 4px 0" />
            <div class="alert-list">
              <div v-if="alerts.length === 0" class="alert-empty">暂无告警</div>
              <div
                v-for="alert in alerts"
                :key="alert.id"
                class="alert-item"
                :class="{ 'alert-unread': !alert.is_read }"
                @click="handleMarkRead(alert)"
              >
                <div class="alert-icon">
                  <n-icon size="18" :color="alert.type === 'backup_failed' ? '#d03050' : '#f0a020'">
                    <AlertCircleOutline v-if="alert.type === 'backup_failed'" />
                    <WarningOutline v-else />
                  </n-icon>
                </div>
                <div class="alert-content">
                  <div class="alert-title">{{ alert.title }}</div>
                  <div class="alert-message">{{ alert.message }}</div>
                  <div class="alert-time">{{ formatTime(alert.created_at) }}</div>
                </div>
              </div>
            </div>
          </div>
        </n-popover>
        <span class="admin-name">
          <n-icon size="18"><PersonOutline /></n-icon>
          {{ authStore.username || '管理员' }}
        </span>
        <n-button quaternary type="error" size="small" @click="handleLogout">
          <template #icon>
            <n-icon><LogOutOutline /></n-icon>
          </template>
          退出
        </n-button>
      </n-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NButton, NIcon, NBreadcrumb, NBreadcrumbItem, NSpace, NBadge, NPopover, NDivider, useMessage } from 'naive-ui'
import { MenuOutline, PersonOutline, LogOutOutline, NotificationsOutline, AlertCircleOutline, WarningOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { getAlerts, getUnreadAlertCount, markAlertAsRead, markAllAlertsAsRead } from '@/api/alert'
import type { Alert } from '@/types'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const authStore = useAuthStore()
const appStore = useAppStore()

const alerts = ref<Alert[]>([])
const unreadCount = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const routeLabelMap: Record<string, string> = {
  Dashboard: '仪表盘',
  DeviceList: '设备管理',
  DeviceDetail: '设备详情',
  FailedDevices: '失败设备',
  AlertList: '告警管理',
  BaselineManage: '基线管理',
  Settings: '系统设置',
}

const breadcrumbItems = computed(() => {
  const items: Array<{ label: string }> = []
  const name = route.name as string

  if (name === 'DeviceDetail') {
    items.push({ label: '设备管理' })
    items.push({ label: '设备详情' })
  } else if (routeLabelMap[name]) {
    items.push({ label: routeLabelMap[name] })
  }

  return items
})

async function fetchUnreadCount() {
  try {
    const res = await getUnreadAlertCount()
    if (res.data) {
      unreadCount.value = res.data.count
    }
  } catch {}
}

async function fetchAlerts() {
  try {
    const res = await getAlerts({ page: 1, page_size: 20 })
    if (res.data) {
      alerts.value = res.data.items
    }
  } catch {}
}

async function handleMarkRead(alert: Alert) {
  if (alert.is_read) return
  try {
    await markAlertAsRead(alert.id)
    alert.is_read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  } catch {
    message.error('标记已读失败')
  }
}

async function handleMarkAllRead() {
  try {
    await markAllAlertsAsRead()
    alerts.value.forEach(a => { a.is_read = true })
    unreadCount.value = 0
  } catch {
    message.error('标记全部已读失败')
  }
}

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  return `${Math.floor(diff / 86400000)} 天前`
}

async function handleLogout() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  await authStore.logout()
  router.push({ name: 'Login' })
}

onMounted(() => {
  fetchUnreadCount()
  fetchAlerts()
  timer = setInterval(() => {
    fetchUnreadCount()
  }, 30000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
})
</script>

<style scoped>
.app-header {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-right {
  display: flex;
  align-items: center;
}

.admin-name {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: rgba(0, 0, 0, 0.65);
}

.alert-panel {
  max-height: 480px;
  display: flex;
  flex-direction: column;
}

.alert-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}

.alert-panel-title {
  font-weight: 600;
  font-size: 14px;
}

.alert-list {
  overflow-y: auto;
  max-height: 400px;
}

.alert-empty {
  padding: 24px 0;
  text-align: center;
  color: rgba(0, 0, 0, 0.35);
  font-size: 13px;
}

.alert-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 4px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.alert-item:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.alert-unread {
  background-color: rgba(24, 160, 88, 0.04);
}

.alert-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.alert-content {
  flex: 1;
  min-width: 0;
}

.alert-title {
  font-size: 13px;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.85);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-message {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.55);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-time {
  font-size: 11px;
  color: rgba(0, 0, 0, 0.35);
  margin-top: 2px;
}
</style>

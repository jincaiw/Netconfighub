<template>
  <div class="app-header">
    <div class="header-left">
      <t-button
        variant="text"
        shape="square"
        :aria-label="appStore.sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'"
        :aria-expanded="!appStore.sidebarCollapsed"
        @click="appStore.toggleSidebar"
      >
        <template #icon>
          <MenuFoldIcon v-if="!appStore.sidebarCollapsed" />
          <MenuUnfoldIcon v-else />
        </template>
      </t-button>
      <t-breadcrumb>
        <t-breadcrumb-item v-for="item in breadcrumbItems" :key="item.label">
          {{ item.label }}
        </t-breadcrumb-item>
      </t-breadcrumb>
    </div>

    <div class="header-right">
      <t-space align="center" :size="12">
        <t-dropdown
          trigger="click"
          placement="bottom-right"
          :popup-props="alertPopupProps"
          :hide-after-item-click="false"
        >
          <t-badge :count="unreadCount" :offset="[4, 4]">
            <t-button variant="text" shape="square" aria-label="告警通知" :aria-expanded="alertDropdownOpen">
              <template #icon>
                <NotificationIcon size="20px" />
              </template>
            </t-button>
          </t-badge>
          <template #dropdown>
            <!-- TDesign 1.20 的 dropdown 仅从 menu item 提取面板内容；纯 div 会被过滤，故包一层 item。 -->
            <t-dropdown-item class="alert-dropdown-item">
              <div class="alert-panel">
              <div class="alert-panel-header">
                <span class="alert-panel-title">告警通知</span>
                <t-button
                  v-if="unreadCount > 0"
                  variant="text"
                  theme="primary"
                  size="small"
                  @click="handleMarkAllRead"
                >
                  全部已读
                </t-button>
              </div>
              <t-divider style="margin: 4px 0" />
              <div v-if="pollError || listError" class="alert-error" role="status">
                <span>告警加载失败，请重试</span>
                <t-button size="small" variant="text" theme="primary" @click="handleRetry">重试</t-button>
              </div>
              <div class="alert-list">
                <div v-if="alerts.length === 0" class="alert-empty">暂无告警</div>
                <div
                  v-for="alert in alerts"
                  :key="alert.id"
                  class="alert-item"
                  :class="{ 'alert-unread': !alert.is_read }"
                  role="button"
                  tabindex="0"
                  :aria-label="`${alert.title}，${alert.message}，${alert.is_read ? '已读' : '未读'}`"
                  @click="handleMarkRead(alert)"
                  @keydown.enter="handleMarkRead(alert)"
                  @keydown.space.prevent="handleMarkRead(alert)"
                >
                  <span v-if="!alert.is_read" class="alert-unread-dot" aria-hidden="true" />
                  <span v-if="!alert.is_read" class="sr-only">未读</span>
                  <div class="alert-icon">
                    <ErrorCircleIcon
                      v-if="alert.type === 'backup_failed'"
                      size="18px"
                      style="color: var(--td-error-color-6)"
                    />
                    <ErrorTriangleIcon v-else size="18px" style="color: var(--td-warning-color-5)" />
                  </div>
                  <div class="alert-content">
                    <div class="alert-title">{{ alert.title }}</div>
                    <div class="alert-message">{{ alert.message }}</div>
                    <div class="alert-time">{{ formatTime(alert.created_at) }}</div>
                  </div>
                </div>
              </div>
              </div>
            </t-dropdown-item>
          </template>
        </t-dropdown>

        <t-dropdown trigger="click" placement="bottom-right" :min-column-width="120" :popup-props="themePopupProps">
          <t-button
            variant="text"
            size="small"
            data-testid="theme-switch"
            :aria-label="`主题设置，当前${themeLabel}`"
            :aria-expanded="themeDropdownOpen"
          >
            <template #icon><PaletteIcon /></template>
            {{ themeLabel }}
          </t-button>
          <template #dropdown>
            <t-dropdown-item
              v-for="option in themeOptions"
              :key="option.value"
              :data-testid="`theme-option-${option.value}`"
              @click="settingStore.changeMode(option.value)"
            >
              {{ option.label }}
            </t-dropdown-item>
          </template>
        </t-dropdown>

        <span class="admin-name">
          <UserIcon size="18px" />
          {{ authStore.username || '管理员' }}
        </span>
        <t-button variant="text" theme="danger" size="small" @click="handleLogout">
          <template #icon><LogoutIcon /></template>
          退出
        </t-button>
      </t-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import {
  MenuFoldIcon,
  MenuUnfoldIcon,
  NotificationIcon,
  ErrorCircleIcon,
  ErrorTriangleIcon,
  PaletteIcon,
  UserIcon,
  LogoutIcon,
} from 'tdesign-icons-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useSettingStore, type ThemeMode } from '@/stores/setting'
import { useAlertSummary } from '@/composables/useAlertSummary'
import type { Alert } from '@/types'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()
const settingStore = useSettingStore()

const {
  unreadCount,
  alerts,
  pollError,
  listError,
  fetchAlerts,
  refresh,
  markRead,
  markAllRead,
  start,
  stop,
} = useAlertSummary()

// 下拉开合状态仅用于 aria-expanded，由 TDesign popup 的可见性回调同步。
const alertDropdownOpen = ref(false)
const themeDropdownOpen = ref(false)
const alertPopupProps = { onVisibleChange: (visible: boolean) => { alertDropdownOpen.value = visible } }
const themePopupProps = { onVisibleChange: (visible: boolean) => { themeDropdownOpen.value = visible } }

const themeOptions: Array<{ label: string; value: ThemeMode }> = [
  { label: '跟随系统', value: 'auto' },
  { label: '浅色', value: 'light' },
  { label: '深色', value: 'dark' },
]

const themeLabel = computed(() => {
  const found = themeOptions.find((option) => option.value === settingStore.mode)
  return found ? found.label : '跟随系统'
})

const routeLabelMap: Record<string, string> = {
  Dashboard: '仪表盘',
  DeviceList: '设备管理',
  DeviceDetail: '设备详情',
  FailedDevices: '失败设备',
  AlertList: '告警管理',
  BaselineManage: '基线管理',
  AuditLogs: '审计日志',
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

async function handleMarkRead(alert: Alert) {
  if (await markRead(alert)) {
    MessagePlugin.success('已标记为已读')
  } else {
    MessagePlugin.error('标记已读失败')
  }
}

async function handleMarkAllRead() {
  if (await markAllRead()) {
    MessagePlugin.success('已全部标记为已读')
  } else {
    MessagePlugin.error('标记全部已读失败')
  }
}

function handleRetry() {
  void refresh()
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
  stop()
  await authStore.logout()
  router.push({ name: 'Login' })
}

onMounted(() => {
  start()
  void fetchAlerts()
})

onUnmounted(() => {
  stop()
})
</script>

<style scoped>
.app-header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
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
  color: var(--td-text-color-secondary);
}

/* 面板被包在 dropdown item 内以满足 TDesign 1.20 的面板提取规则，去掉菜单项的默认内边距。 */
:deep(.alert-dropdown-item) {
  padding: 8px 12px;
  cursor: default;
}

.alert-panel {
  max-height: 480px;
  width: 340px;
  max-width: calc(100vw - 32px);
  display: flex;
  flex-direction: column;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

.alert-error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 6px 4px;
  color: var(--td-error-color);
  font-size: 12px;
}

.alert-unread-dot {
  flex-shrink: 0;
  width: 6px;
  height: 6px;
  margin-top: 7px;
  border-radius: 50%;
  background-color: var(--td-brand-color);
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
  color: var(--td-text-color-placeholder);
  font-size: 13px;
}

.alert-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 4px;
  cursor: pointer;
  border-radius: var(--td-radius-default);
  transition: background-color 0.2s;
}

.alert-item:hover {
  background-color: var(--td-bg-color-container-hover);
}

.alert-unread {
  background-color: var(--td-brand-color-light);
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
  color: var(--td-text-color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-message {
  font-size: 12px;
  color: var(--td-text-color-secondary);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-time {
  font-size: 11px;
  color: var(--td-text-color-placeholder);
  margin-top: 2px;
}

@media (max-width: 768px) {
  .app-header {
    padding: 0 10px;
  }

  .admin-name {
    display: none;
  }
}
</style>

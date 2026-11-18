<template>
  <t-menu
    :value="activeKey"
    :collapsed="collapsed"
    :theme="settingStore.displayMode"
    class="app-menu"
    @change="handleMenuSelect"
  >
    <t-menu-item v-for="item in menuItems" :key="item.key" :value="item.key">
      <template #icon>
        <component :is="item.icon" />
      </template>
      {{ item.label }}
    </t-menu-item>
  </t-menu>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSettingStore } from '@/stores/setting'
import {
  DashboardIcon,
  ServerIcon,
  ErrorTriangleIcon,
  NotificationIcon,
  SecuredIcon,
  FileIcon,
  SettingIcon,
} from 'tdesign-icons-vue-next'

defineProps<{
  collapsed: boolean
}>()

const emit = defineEmits<{
  navigate: []
}>()

const router = useRouter()
const route = useRoute()
const settingStore = useSettingStore()

const activeKey = computed(() => {
  const name = route.name as string
  if (name === 'DeviceDetail') return 'DeviceList'
  return name
})

const menuItems: Array<{ label: string; key: string; icon: Component }> = [
  { label: '仪表盘', key: 'Dashboard', icon: DashboardIcon },
  { label: '设备管理', key: 'DeviceList', icon: ServerIcon },
  { label: '失败设备', key: 'FailedDevices', icon: ErrorTriangleIcon },
  { label: '告警管理', key: 'AlertList', icon: NotificationIcon },
  { label: '基线管理', key: 'BaselineManage', icon: SecuredIcon },
  { label: '审计日志', key: 'AuditLogs', icon: FileIcon },
  { label: '系统设置', key: 'Settings', icon: SettingIcon },
]

function handleMenuSelect(key: string | number) {
  router.push({ name: String(key) })
  emit('navigate')
}
</script>

<style scoped>
.app-menu {
  width: 100%;
  padding: 14px 10px 60px;
}

:deep(.t-menu__item) {
  margin-bottom: 4px;
  border-radius: var(--td-radius-default);
}

:deep(.t-menu__item.t-is-active) {
  font-weight: 600;
}
</style>

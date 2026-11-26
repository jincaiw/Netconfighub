<template>
  <t-menu
    :value="activeKey"
    :collapsed="collapsed"
    :theme="settingStore.displayMode"
    class="app-menu"
    @change="handleMenuSelect"
  >
    <t-menu-group title="总览">
      <t-menu-item value="Dashboard">
        <template #icon><DashboardIcon /></template>
        仪表盘
      </t-menu-item>
    </t-menu-group>
    <t-menu-group title="运维监控">
      <t-menu-item v-for="item in opsItems" :key="item.key" :value="item.key">
        <template #icon><component :is="item.icon" /></template>
        {{ item.label }}
      </t-menu-item>
    </t-menu-group>
    <t-menu-group title="系统">
      <t-menu-item value="Settings">
        <template #icon><SettingIcon /></template>
        系统设置
      </t-menu-item>
    </t-menu-group>
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

defineProps<{ collapsed: boolean }>()
const emit = defineEmits<{ navigate: [] }>()

const router = useRouter()
const route = useRoute()
const settingStore = useSettingStore()

const activeKey = computed(() => {
  const name = route.name as string
  if (name === 'DeviceDetail') return 'DeviceList'
  return name
})

const opsItems: Array<{ label: string; key: string; icon: Component }> = [
  { label: '设备管理', key: 'DeviceList', icon: ServerIcon },
  { label: '失败设备', key: 'FailedDevices', icon: ErrorTriangleIcon },
  { label: '告警管理', key: 'AlertList', icon: NotificationIcon },
  { label: '基线管理', key: 'BaselineManage', icon: SecuredIcon },
  { label: '审计日志', key: 'AuditLogs', icon: FileIcon },
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

:deep(.t-menu-group__title) {
  color: var(--shadcn-muted-foreground);
  font-size: 12px;
  font-weight: 500;
  padding: 12px 8px 4px;
}
</style>

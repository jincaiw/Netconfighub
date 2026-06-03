<template>
  <n-menu
    :collapsed="collapsed"
    :collapsed-width="64"
    :collapsed-icon-size="22"
    :options="menuOptions"
    :value="activeKey"
    @update:value="handleMenuSelect"
  />
</template>

<script setup lang="ts">
import { computed, h } from 'vue'
import type { Component } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NMenu, NIcon } from 'naive-ui'
import {
  SpeedometerOutline,
  ServerOutline,
  WarningOutline,
  NotificationsOutline,
  ShieldCheckmarkOutline,
  SettingsOutline,
  DocumentTextOutline,
} from '@vicons/ionicons5'

defineProps<{
  collapsed: boolean
}>()

const router = useRouter()
const route = useRoute()

const activeKey = computed(() => {
  const name = route.name as string
  if (name === 'DeviceDetail') return 'DeviceList'
  return name
})

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = [
  {
    label: '仪表盘',
    key: 'Dashboard',
    icon: renderIcon(SpeedometerOutline),
  },
  {
    label: '设备管理',
    key: 'DeviceList',
    icon: renderIcon(ServerOutline),
  },
  {
    label: '失败设备',
    key: 'FailedDevices',
    icon: renderIcon(WarningOutline),
  },
  {
    label: '告警管理',
    key: 'AlertList',
    icon: renderIcon(NotificationsOutline),
  },
  {
    label: '基线管理',
    key: 'BaselineManage',
    icon: renderIcon(ShieldCheckmarkOutline),
  },
  {
    label: '审计日志',
    key: 'AuditLogs',
    icon: renderIcon(DocumentTextOutline),
  },
  {
    label: '系统设置',
    key: 'Settings',
    icon: renderIcon(SettingsOutline),
  },
]

function handleMenuSelect(key: string) {
  router.push({ name: key })
}
</script>

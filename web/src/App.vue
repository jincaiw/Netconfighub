<template>
  <t-config-provider>
    <router-view />
  </t-config-provider>
</template>

<script setup lang="ts">
import { onUnmounted } from 'vue'
import { useSettingStore } from '@/stores/setting'

const settingStore = useSettingStore()

onUnmounted(() => {
  settingStore.dispose()
})
</script>

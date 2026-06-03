<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :collapsed="appStore.sidebarCollapsed"
      show-trigger
      @collapse="appStore.sidebarCollapsed = true"
      @expand="appStore.sidebarCollapsed = false"
    >
      <div class="sider-logo" @click="router.push('/dashboard')">
        <n-icon size="28" :component="PulseOutline" class="logo-icon" />
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">NetConfigHub</span>
        <span v-else class="logo-text-short">NCH</span>
      </div>
      <AppSidebar :collapsed="appStore.sidebarCollapsed" />
    </n-layout-sider>
    <n-layout>
      <n-layout-header bordered>
        <AppHeader />
      </n-layout-header>
      <n-layout-content content-style="padding: 24px;">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NIcon } from 'naive-ui'
import { PulseOutline } from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import AppSidebar from '@/components/AppSidebar.vue'
import AppHeader from '@/components/AppHeader.vue'

const router = useRouter()
const appStore = useAppStore()
</script>

<style scoped>
.sider-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  color: #18a058;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  cursor: pointer;
  user-select: none;
}

.logo-icon {
  flex-shrink: 0;
}

.logo-text {
  white-space: nowrap;
}

.logo-text-short {
  font-size: 16px;
}
</style>

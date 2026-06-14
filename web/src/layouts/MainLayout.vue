<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      class="desktop-sider"
      bordered
      inverted
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
    <n-drawer v-model:show="appStore.mobileSidebarOpen" placement="left" :width="280">
      <n-drawer-content
        body-content-style="padding: 0; background: #071a33"
        :native-scrollbar="false"
      >
        <div class="sider-logo mobile-logo" @click="goHome">
          <n-icon size="28" :component="PulseOutline" class="logo-icon" />
          <span class="logo-text">NetConfigHub</span>
        </div>
        <AppSidebar :collapsed="false" @click="appStore.mobileSidebarOpen = false" />
      </n-drawer-content>
    </n-drawer>
    <n-layout>
      <n-layout-header bordered>
        <AppHeader />
      </n-layout-header>
      <n-layout-content class="main-content">
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
import {
  NDrawer,
  NDrawerContent,
  NIcon,
  NLayout,
  NLayoutContent,
  NLayoutHeader,
  NLayoutSider,
} from 'naive-ui'
import { PulseOutline } from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import AppSidebar from '@/components/AppSidebar.vue'
import AppHeader from '@/components/AppHeader.vue'

const router = useRouter()
const appStore = useAppStore()

function goHome() {
  appStore.mobileSidebarOpen = false
  router.push('/dashboard')
}
</script>

<style scoped>
.sider-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  color: #f8fafc;
  background: #071a33;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
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

.mobile-logo {
  justify-content: flex-start;
  padding: 0 22px;
}

.main-content {
  min-width: 0;
  padding: 28px 30px 36px;
  background: #f4f7fb;
}

@media (max-width: 768px) {
  .desktop-sider {
    display: none;
  }

  .main-content {
    padding: 18px 14px 28px;
  }
}
</style>

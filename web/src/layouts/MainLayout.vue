<template>
  <t-layout class="app-layout">
    <t-aside class="desktop-sider" :style="{ width: appStore.sidebarCollapsed ? '64px' : '220px' }">
      <div class="sider-logo" @click="goHome">
        <ActivityIcon size="28px" class="logo-icon" />
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">NetConfigHub</span>
        <span v-else class="logo-text-short">NCH</span>
      </div>
      <AppSidebar :collapsed="appStore.sidebarCollapsed" />
      <t-button variant="text" shape="square" class="collapse-btn" @click="appStore.toggleSidebar">
        <template #icon>
          <MenuUnfoldIcon v-if="appStore.sidebarCollapsed" />
          <MenuFoldIcon v-else />
        </template>
      </t-button>
    </t-aside>

    <t-drawer
      v-model:visible="appStore.mobileSidebarOpen"
      placement="left"
      :width="280"
      :header="false"
      :footer="false"
    >
      <div class="sider-logo mobile-logo" @click="goHome">
        <ActivityIcon size="28px" class="logo-icon" />
        <span class="logo-text">NetConfigHub</span>
      </div>
      <AppSidebar :collapsed="false" @navigate="appStore.mobileSidebarOpen = false" />
    </t-drawer>

    <t-layout class="main-layout">
      <t-header class="header-wrap">
        <AppHeader />
      </t-header>
      <t-content class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </t-content>
    </t-layout>
  </t-layout>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ActivityIcon, MenuFoldIcon, MenuUnfoldIcon } from 'tdesign-icons-vue-next'
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
.app-layout {
  height: 100vh;
  height: 100dvh;
}

.desktop-sider {
  position: relative;
  flex-shrink: 0;
  transition: width 0.2s ease;
  border-right: 1px solid var(--td-component-stroke);
  background: var(--td-bg-color-container);
}

.sider-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  color: var(--td-text-color-anti);
  background: var(--td-brand-color);
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

.collapse-btn {
  position: absolute;
  left: 0;
  bottom: 12px;
  width: 100%;
}

.header-wrap {
  height: 64px;
  line-height: normal;
  padding: 0;
  border-bottom: 1px solid var(--td-component-stroke);
  background: var(--td-bg-color-container);
}

.main-layout {
  min-width: 0;
  max-width: 100%;
  overflow: hidden;
}

.main-content {
  min-width: 0;
  padding: 28px 30px 36px;
  background: var(--td-bg-color-page);
  overflow: auto;
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

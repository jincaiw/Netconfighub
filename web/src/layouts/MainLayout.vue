<template>
  <t-layout class="app-layout">
    <t-aside class="desktop-sider" :style="{ width: appStore.sidebarCollapsed ? '64px' : '220px' }">
      <div class="sider-logo" @click="goHome">
        <span class="team-badge"><ActivityIcon size="18px" class="logo-icon" /></span>
        <span v-if="!appStore.sidebarCollapsed" class="team-meta">
          <span class="team-name">NetConfigHub</span>
          <span class="team-plan">设备配置管理平台</span>
        </span>
        <span v-else class="logo-text-short">NCH</span>
        <ChevronDownIcon v-if="!appStore.sidebarCollapsed" class="team-switch" size="14px" />
      </div>
      <AppSidebar :collapsed="appStore.sidebarCollapsed" />
      <t-dropdown v-if="!appStore.sidebarCollapsed" trigger="click" placement="top-left" :min-column-width="160">
        <div class="sider-user" role="button" tabindex="0" aria-label="用户菜单" aria-haspopup="menu">
          <span class="drawer-avatar">AD</span>
          <span class="drawer-user-meta">
            <span class="drawer-user-name">admin</span>
            <span class="drawer-user-mail">admin@netconfighub.local</span>
          </span>
          <ChevronDownIcon class="sider-user-switch" size="14px" />
        </div>
        <template #dropdown>
          <t-dropdown-item @click="goSettings">个人设置</t-dropdown-item>
          <t-dropdown-item @click="handleLogout">退出登录</t-dropdown-item>
        </template>
      </t-dropdown>
    </t-aside>

    <t-drawer
      v-model:visible="appStore.mobileSidebarOpen"
      placement="left"
      :width="280"
      :header="false"
      :footer="false"
    >
      <div class="sider-logo mobile-logo" @click="goHome">
        <span class="team-badge"><ActivityIcon size="18px" class="logo-icon" /></span>
        <span class="team-meta">
          <span class="team-name">NetConfigHub</span>
          <span class="team-plan">设备配置管理平台</span>
        </span>
      </div>
      <AppSidebar :collapsed="false" @navigate="appStore.mobileSidebarOpen = false" />
      <div class="drawer-user" @click="appStore.mobileSidebarOpen = false">
        <span class="drawer-avatar">admin</span>
        <span class="drawer-user-meta">
          <span class="drawer-user-name">admin</span>
          <span class="drawer-user-mail">admin@netconfighub.local</span>
        </span>
      </div>
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
import { ActivityIcon, ChevronDownIcon } from 'tdesign-icons-vue-next'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import AppSidebar from '@/components/AppSidebar.vue'
import AppHeader from '@/components/AppHeader.vue'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

function goHome() {
  appStore.mobileSidebarOpen = false
  router.push('/dashboard')
}

function goSettings() {
  router.push('/settings')
}

async function handleLogout() {
  await authStore.logout()
  router.push({ name: 'Login' })
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
  display: flex;
  flex-direction: column;
  transition: width 0.2s ease;
  border-right: 1px solid var(--td-component-stroke);
  background: var(--td-bg-color-container);
}

.sider-user {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: auto;
  padding: 12px;
  border-top: 1px solid var(--shadcn-border);
  cursor: pointer;
}
.sider-user-switch {
  margin-left: auto;
  color: var(--shadcn-muted-foreground);
  flex-shrink: 0;
}

.sider-logo {
  height: 64px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  color: var(--shadcn-foreground);
  background: var(--shadcn-card);
  border-bottom: 1px solid var(--shadcn-border);
  cursor: pointer;
  user-select: none;
}

.team-badge {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: var(--shadcn-primary);
  color: var(--shadcn-primary-foreground);
}

.team-meta {
  display: flex;
  flex-direction: column;
  min-width: 0;
  line-height: 1.25;
}
.team-name {
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
}
.team-plan {
  font-size: 11px;
  color: var(--shadcn-muted-foreground);
  white-space: nowrap;
}
.team-switch {
  margin-left: auto;
  flex-shrink: 0;
  color: var(--shadcn-muted-foreground);
}

.logo-icon {
  flex-shrink: 0;
}

.logo-text-short {
  font-size: 13px;
  font-weight: 600;
}

.mobile-logo {
  justify-content: flex-start;
}

.drawer-user {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  border-top: 1px solid var(--shadcn-border);
  cursor: pointer;
}
.drawer-avatar {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 11px;
  background: var(--shadcn-muted);
  color: var(--shadcn-foreground);
}
.drawer-user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
  min-width: 0;
}
.drawer-user-name {
  font-size: 13px;
  font-weight: 600;
}
.drawer-user-mail {
  font-size: 11px;
  color: var(--shadcn-muted-foreground);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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

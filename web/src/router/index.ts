import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { getToken, isTokenExpired } from '@/utils/token'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false, title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: '运维总览' } },
      { path: 'devices', name: 'DeviceList', component: () => import('@/views/DeviceList.vue'), meta: { title: '设备管理' } },
      { path: 'devices/:id', name: 'DeviceDetail', component: () => import('@/views/DeviceDetail.vue'), meta: { title: '设备详情' } },
      { path: 'failed', name: 'FailedDevices', component: () => import('@/views/FailedDevices.vue'), meta: { title: '失败设备' } },
      { path: 'alerts', name: 'AlertList', component: () => import('@/views/AlertList.vue'), meta: { title: '告警管理' } },
      { path: 'baselines', name: 'BaselineManage', component: () => import('@/views/BaselineManage.vue'), meta: { title: '基线管理' } },
      { path: 'audit-logs', name: 'AuditLogs', component: () => import('@/views/AuditLogs.vue'), meta: { title: '审计日志' } },
      { path: 'settings', name: 'Settings', component: () => import('@/views/Settings.vue'), meta: { title: '系统设置' } },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { requiresAuth: false, title: '页面不存在' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = getToken()
  const authenticated = !!token && !isTokenExpired()
  if (to.meta.requiresAuth !== false && !authenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.name === 'Login' && authenticated) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

// 页面标题随路由变化（浏览器标签/读屏可辨识）。默认站名 NetConfigHub。
const APP_NAME = 'NetConfigHub'
router.afterEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} · ${APP_NAME}` : APP_NAME
})

export default router

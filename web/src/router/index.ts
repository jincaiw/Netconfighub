import { createRouter, createWebHistory } from 'vue-router'
import { getToken, isTokenExpired } from '@/utils/token'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
        },
        {
          path: 'devices',
          name: 'DeviceList',
          component: () => import('@/views/DeviceList.vue'),
        },
        {
          path: 'devices/:id',
          name: 'DeviceDetail',
          component: () => import('@/views/DeviceDetail.vue'),
        },
        {
          path: 'failed',
          name: 'FailedDevices',
          component: () => import('@/views/FailedDevices.vue'),
        },
        {
          path: 'alerts',
          name: 'AlertList',
          component: () => import('@/views/AlertList.vue'),
        },
        {
          path: 'baselines',
          name: 'BaselineManage',
          component: () => import('@/views/BaselineManage.vue'),
        },
        {
          path: 'audit-logs',
          name: 'AuditLogs',
          component: () => import('@/views/AuditLogs.vue'),
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('@/views/Settings.vue'),
        },
      ],
    },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { requiresAuth: false },
  },
]
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

export default router

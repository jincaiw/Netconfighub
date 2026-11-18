import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { DARK_CHART_COLORS, LIGHT_CHART_COLORS, type ChartColors } from '@/config/chartColors'

export type ThemeMode = 'light' | 'dark' | 'auto'
export type DisplayMode = 'light' | 'dark'

const STORAGE_KEY = 'nch-theme-mode'
const DARK_QUERY = '(prefers-color-scheme: dark)'

/** 主题模式：light/dark/auto，挂载前初始化，localStorage 持久化，auto 跟随系统偏好。 */
export const useSettingStore = defineStore('setting', () => {
  const mode = ref<ThemeMode>('auto')
  const systemDark = ref(false)
  const chartColors = ref<ChartColors>({ ...LIGHT_CHART_COLORS })
  let mediaQuery: MediaQueryList | null = null
  let bound = false

  const displayMode = computed<DisplayMode>(() => {
    if (mode.value === 'auto') return systemDark.value ? 'dark' : 'light'
    return mode.value
  })

  function readStoredMode(): { mode: ThemeMode; available: boolean } {
    try {
      const stored = window.localStorage.getItem(STORAGE_KEY)
      if (stored === 'light' || stored === 'dark' || stored === 'auto') {
        return { mode: stored, available: true }
      }
      return { mode: 'auto', available: true }
    } catch {
      // localStorage 不可用：与 index.html head 脚本一致回退 auto + 系统偏好，
      // 并继续写入 theme-mode，避免 DOM 深色而 Store 浅色的降级分叉。
      return { mode: 'auto', available: true }
    }
  }

  function applyTheme(display: DisplayMode, writeAttribute = true) {
    const isDark = display === 'dark'
    if (writeAttribute) {
      document.documentElement.setAttribute('theme-mode', isDark ? 'dark' : '')
    }
    chartColors.value = isDark ? { ...DARK_CHART_COLORS } : { ...LIGHT_CHART_COLORS }
  }

  function handleSystemChange(event: MediaQueryListEvent) {
    systemDark.value = event.matches
    if (mode.value === 'auto') {
      applyTheme(event.matches ? 'dark' : 'light')
    }
  }

  function bindMediaListener() {
    unbindMediaListener()
    if (mode.value !== 'auto' || !mediaQuery) return
    if (typeof mediaQuery.addEventListener === 'function') {
      mediaQuery.addEventListener('change', handleSystemChange)
    } else {
      mediaQuery.addListener(handleSystemChange)
    }
    bound = true
  }

  function unbindMediaListener() {
    if (!mediaQuery || !bound) return
    if (typeof mediaQuery.removeEventListener === 'function') {
      mediaQuery.removeEventListener('change', handleSystemChange)
    } else {
      mediaQuery.removeListener(handleSystemChange)
    }
    bound = false
  }

  /** 挂载前调用一次：恢复持久化选择并立即应用，避免首屏主题闪烁。 */
  function initTheme() {
    if (typeof window === 'undefined') return
    try {
      mediaQuery = window.matchMedia(DARK_QUERY)
      systemDark.value = mediaQuery.matches
    } catch {
      mediaQuery = null
      systemDark.value = false
    }

    const stored = readStoredMode()
    mode.value = stored.mode
    const initialDisplay: DisplayMode =
      stored.mode === 'auto' ? (systemDark.value ? 'dark' : 'light') : stored.mode
    applyTheme(initialDisplay, stored.available)
    bindMediaListener()
  }

  function changeMode(next: ThemeMode) {
    mode.value = next
    try {
      window.localStorage.setItem(STORAGE_KEY, next)
    } catch {
      // 存储失败不影响本次会话内的主题切换
    }
    if (next === 'auto' && mediaQuery) {
      systemDark.value = mediaQuery.matches
    }
    applyTheme(displayMode.value)
    bindMediaListener()
  }

  function dispose() {
    unbindMediaListener()
  }

  return {
    mode,
    systemDark,
    displayMode,
    chartColors,
    initTheme,
    changeMode,
    dispose,
  }
})

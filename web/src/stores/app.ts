import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const mobileSidebarOpen = ref(false)

  function toggleSidebar() {
    if (window.matchMedia('(max-width: 768px)').matches) {
      mobileSidebarOpen.value = !mobileSidebarOpen.value
      return
    }
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    sidebarCollapsed,
    mobileSidebarOpen,
    toggleSidebar,
  }
})

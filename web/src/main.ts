import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useSettingStore } from './stores/setting'
import { setupTDesign } from './plugins/tdesign'
import './styles/shadcn-theme.css'
import './styles/global.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
// 挂载前恢复主题，避免首屏主题闪烁
useSettingStore(pinia).initTheme()
setupTDesign(app)
app.use(router)
app.mount('#app')

// 无障碍增强：TDesign 的 input-number 增减为图标按钮且缺 aria-label，
// 在挂载后为其补充可访问名称（减少/增加），覆盖动态渲染。
function enhanceIconButtonA11y(root: ParentNode): void {
  root.querySelectorAll('.t-input-number__decrease').forEach((el) => {
    if (!el.getAttribute('aria-label')) el.setAttribute('aria-label', '减少')
  })
  root.querySelectorAll('.t-input-number__increase').forEach((el) => {
    if (!el.getAttribute('aria-label')) el.setAttribute('aria-label', '增加')
  })
}
enhanceIconButtonA11y(document)
new MutationObserver((mutations) => {
  for (const m of mutations) {
    m.addedNodes.forEach((n) => {
      if (n.nodeType === 1) enhanceIconButtonA11y(n as Element)
    })
  }
}).observe(document.body, { childList: true, subtree: true })

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

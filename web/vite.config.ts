import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

const apiProxyTarget = process.env.NCH_API_PROXY_TARGET || 'http://127.0.0.1:8080'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id: string) {
          if (id.includes('node_modules/echarts') || id.includes('node_modules/zrender')) {
            return 'echarts'
          }
          if (id.includes('node_modules/tdesign-vue-next') || id.includes('node_modules/tdesign-icons-vue-next')) {
            return 'tdesign'
          }
          if (id.includes('node_modules/vue') || id.includes('node_modules/@vue') || id.includes('node_modules/pinia')) {
            return 'vue-vendor'
          }
          return undefined
        },
      },
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
    },
  },
})

import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  retries: 1,
  use: {
    baseURL: 'http://127.0.0.1:13100',
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command: 'cd .. && rm -f /tmp/netconfighub-e2e.db && rm -rf /tmp/netconfighub-e2e-configs && NCH_CONFIG_PATH=configs/config.e2e.yaml go run cmd/api/main.go',
      port: 18080,
      reuseExistingServer: false,
      timeout: 15000,
    },
    {
      command: 'NCH_API_PROXY_TARGET=http://127.0.0.1:18080 npm run dev -- --host 127.0.0.1 --port 13100',
      port: 13100,
      reuseExistingServer: false,
      timeout: 15000,
    },
  ],
})

import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  retries: 1,
  use: {
    baseURL: 'http://localhost:3100',
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command: 'cd .. && rm -f /tmp/netconfighub-e2e.db && rm -rf /tmp/netconfighub-e2e-configs && NCH_CONFIG_PATH=configs/config.e2e.yaml go run cmd/api/main.go',
      port: 8080,
      reuseExistingServer: false,
      timeout: 15000,
    },
    {
      command: 'npm run dev -- --host 127.0.0.1 --port 3100',
      port: 3100,
      reuseExistingServer: false,
      timeout: 15000,
    },
  ],
})

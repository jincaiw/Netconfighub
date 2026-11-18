import { defineConfig } from '@playwright/test'

// The Go webServer and the test workers share one bootstrap password so E2E
// stays deterministic after the literal admin/admin fallback was removed.
const e2eAdminPassword = process.env.NCH_ADMIN_PASSWORD || 'nch-e2e-admin-2024!'
process.env.NCH_ADMIN_PASSWORD = e2eAdminPassword

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
      env: { ...process.env, NCH_ADMIN_PASSWORD: e2eAdminPassword } as Record<string, string>,
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

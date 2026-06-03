# Release Notes

## v0.1.0

Initial NetConfigHub release.

### Highlights

- Go API server with embedded Vue web UI.
- Admin login, JWT auth, and API-token auth.
- Device, group, backup, baseline, deviation, alert, hook, config, failed-device,
  and audit-log APIs.
- Device import/export and password-safe API responses.
- SSH/Telnet collection framework with model-specific commands and sanitization.
- Git-backed configuration storage and diff/history endpoints.
- Scheduler and retry worker flow for backup tasks.
- SQLite default configuration and MySQL configuration support.
- Docker, Docker Compose, and systemd deployment artifacts.
- Playwright e2e suite and comprehensive API smoke script.

### Fixes Included Before Release

- Missing device configuration and version history now return `404` instead of a
  successful empty response.
- Empty group backup requests now return `400`.
- Device create requests accept both current `ip`/`protocol` fields and legacy
  `ip_address`/`conn_protocol` aliases.
- Device list UI renders a stable semantic table and supports arbitrary model
  strings.
- E2E tests use isolated runtime state and a fixed frontend port to avoid stale
  local services.
- Frontend build tooling was upgraded to Vite 8 and `@vitejs/plugin-vue` 6,
  resolving the Vite/esbuild moderate-severity audit findings.

### Verification

The release was verified with:

```bash
go test ./...
go vet ./...
cd web && npm audit --json
cd web && npm run build
go build -o netconfighub-api ./cmd/api
cd web && npx playwright test
bash test_api_v2.sh
```

Observed verification results:

- Go unit/integration tests: passed.
- Go vet: passed.
- npm audit: 0 vulnerabilities.
- Frontend production build: passed.
- Backend binary build: passed.
- Playwright e2e tests: 39 passed.
- Comprehensive API smoke test: 83 passed, 0 failed.
- Browser verification: login and device list page rendered correctly.

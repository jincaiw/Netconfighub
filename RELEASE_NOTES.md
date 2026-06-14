# Release Notes

## v0.1.2

Production hardening, UI refinement, and reproducible deployment release.

### Highlights

- Refined the login, dashboard, navigation, desktop, and mobile layouts.
- Reduced the largest frontend JavaScript chunk from about 1.03 MB to 165 KB.
- Added production startup validation for JWT, encryption, administrator
  password, and CORS settings.
- Added environment-only configuration for single-file and container
  deployments.
- Isolated Playwright and API test ports from local development services.
- Added a non-root, health-checked Docker image and hardened systemd unit.
- Added reproducible Linux amd64 release packaging and checksums.
- Added bilingual deployment guides and three current application screenshots.

### Verification

The release is verified with:

```bash
go test ./...
go vet ./...
cd web && npm audit --json
cd web && npm run build
cd web && npx playwright test
bash test_api_v2.sh
make release-bundle VERSION=v0.1.2
```

Observed verification results:

- Go unit/integration tests: passed.
- Go vet: passed.
- npm audit: 0 vulnerabilities.
- Frontend production build: passed.
- Playwright e2e tests: 39 passed.
- Comprehensive API smoke test: 83 passed, 0 failed.
- Linux amd64 static binary and Docker runtime smoke checks: passed.

### Assets

- `netconfighub-v0.1.2-linux-amd64.tar.gz`
  - SHA256:
    `7484a5862ea7b2d389e53dfe85e2914aa67ff8ce6bb3529e8eb2ccda272515db`

## v0.1.1

Production-readiness and Linux binary release.

### Highlights

- Added configurable CORS policy via `server.cors.*`.
- Hardened embedded frontend routing with explicit path normalization.
- Rejected path-like device and group names before they can be used in Git
  storage paths.
- Added Git storage path-component validation as a defensive second layer.
- Added first-run admin bootstrap environment variables:
  `NCH_ADMIN_USERNAME` and `NCH_ADMIN_PASSWORD`.
- Updated Docker build images to match current Go/Node toolchain requirements
  and use `CGO_ENABLED=0` for the Go binary.
- Added bilingual documentation: English README by default and
  `README.zh-CN.md` for Chinese.
- Published a direct Linux amd64 binary asset.

### Verification

The release was verified with:

```bash
go test ./...
go vet ./...
cd web && npm audit --package-lock-only --json
cd web && npm run build
cd web && npx playwright test
bash test_api_v2.sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/netconfighub-v0.1.1-linux-amd64 ./cmd/api
```

Observed verification results:

- Go unit/integration tests: passed.
- Go vet: passed.
- npm audit: 0 vulnerabilities.
- Frontend production build: passed.
- Playwright e2e tests: 39 passed.
- Comprehensive API smoke test: 83 passed, 0 failed.
- Linux amd64 binary build: passed.

### Assets

- `netconfighub-v0.1.1-linux-amd64.tar.gz`
  - SHA256:
    `89a40e6569d5dfc8356be2f5c50933b49d4b97561dc0872ade69396481275589`

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

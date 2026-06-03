# NetConfigHub

NetConfigHub is a network device configuration backup and management platform for
small and mid-sized network operations teams. It provides a Go backend, embedded
Vue web UI, REST APIs, scheduled/manual backup jobs, Git-backed configuration
history, baseline comparison, alerts, hooks, audit logs, and API-token access.

## Features

- Device CRUD with enable/disable, grouping, import, and export.
- SSH/Telnet collectors for Cisco, Huawei, H3C, and Ruijie style devices.
- Manual and scheduled backup jobs with retry handling.
- Git-backed configuration storage and diff/history APIs.
- Baseline management and deviation reporting.
- Alert, webhook, failed-device, audit-log, and system-config APIs.
- Embedded web UI served by the Go binary.
- SQLite by default, with MySQL configuration support.
- Docker Compose and systemd deployment examples.

## Quick Start

Build and run the single binary:

```bash
go build -o netconfighub ./cmd/api
./netconfighub
```

The default configuration is loaded from `configs/config.yaml`.

Open:

```text
http://localhost:8080
```

Default admin account:

```text
username: admin
password: admin
```

Change the default password and `server.jwt_secret` before using the system in a
real environment.

## Docker

```bash
docker compose up --build
```

The Compose setup exposes the application on `localhost:8080` and stores runtime
data under `./data`.

To start the optional MySQL service:

```bash
docker compose --profile mysql up --build
```

## Configuration

Important settings in `configs/config.yaml`:

- `server.host`, `server.port`: HTTP bind address.
- `server.jwt_secret`: JWT signing secret. Replace the default value.
- `database.driver`: `sqlite` or `mysql`.
- `database.sqlite_path`: SQLite database path.
- `database.mysql_dsn`: MySQL DSN when using MySQL.
- `git.repo_path`: local Git repository used for stored device configs.
- `scheduler.*`: backup interval, worker pool, retry, and timeout settings.
- `sanitize.*`: masking rules for collected secrets.
- `notify.*`: webhook notification settings.

For automated e2e tests, use `configs/config.e2e.yaml`, which writes isolated
runtime state to `/tmp`.

## Development

Backend tests:

```bash
go test ./...
go vet ./...
```

Frontend build:

```bash
cd web
npm ci
npm run build
```

Frontend tooling requires Node.js 20.19 or newer.

Playwright e2e tests:

```bash
cd web
npm ci
npx playwright test
```

Comprehensive API smoke test:

```bash
rm -f /tmp/netconfighub-e2e.db
rm -rf /tmp/netconfighub-e2e-configs
NCH_CONFIG_PATH=configs/config.e2e.yaml go run cmd/api/main.go
bash test_api_v2.sh
```

`web/dist` is intentionally committed because `web/embed.go` embeds it into the
Go binary.

## Release

Current release tag:

```text
v0.1.0
```

See `RELEASE_NOTES.md` for release contents and verification notes.

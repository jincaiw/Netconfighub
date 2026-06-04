# NetConfigHub

中文 | [English](README.md)

NetConfigHub 是面向中小型网络运维团队的网络设备配置备份与管理平台。它包含 Go 后端、内嵌 Vue Web UI、REST API、手动/定时备份任务、Git 配置版本管理、基线比对、告警、Hook、审计日志和 API Token 访问能力。

## 功能

- 设备新增、编辑、删除、启用/禁用、分组、导入和导出。
- 支持 Cisco、Huawei、H3C、Ruijie 风格设备的 SSH/Telnet 采集框架。
- 手动和定时备份任务，支持重试。
- Git 存储设备配置，提供历史和 Diff API。
- 基线管理和偏差报告。
- 告警、Webhook、失败设备、审计日志和系统配置 API。
- Go 单二进制内嵌 Web UI。
- 默认 SQLite，支持 MySQL 配置。
- 提供 Docker Compose 和 systemd 部署示例。

## 快速开始

构建并运行单二进制：

```bash
go build -o netconfighub ./cmd/api
./netconfighub
```

默认配置文件为 `configs/config.yaml`。

打开：

```text
http://localhost:8080
```

默认管理员账号：

```text
username: admin
password: admin
```

生产环境使用前，请修改默认密码、`server.jwt_secret` 和 `server.encryption_key`。

首次启动时可以通过环境变量设置初始管理员：

```bash
export NCH_ADMIN_USERNAME=admin
export NCH_ADMIN_PASSWORD='change-this-strong-password'
```

## Linux 单文件部署

在 Linux 主机上构建，或者从 macOS 交叉编译 Linux 二进制：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o netconfighub ./cmd/api
```

把二进制和配置文件放到标准位置：

```bash
sudo mkdir -p /opt/netconfighub /etc/netconfighub /var/lib/netconfighub
sudo install -m 0755 netconfighub /opt/netconfighub/netconfighub
sudo install -m 0644 configs/config.yaml /etc/netconfighub/config.yaml
```

然后修改 `/etc/netconfighub/config.yaml`，至少完成这些生产配置：

- 将 `database.sqlite_path` 改为 `/var/lib/netconfighub/netconfighub.db`
- 将 `git.repo_path` 改为 `/var/lib/netconfighub/configs`
- 替换 `server.jwt_secret` 和 `server.encryption_key`
- 将 `server.cors.allowed_origins` 收紧为真实浏览器来源

启动方式：

```bash
NCH_CONFIG_PATH=/etc/netconfighub/config.yaml /opt/netconfighub/netconfighub
```

如果你要直接接入 `systemd`，仓库里已经提供了 `deploy/netconfighub.service`。

## systemd 服务

`deploy/netconfighub.service` 适合搭配独立服务账号，并默认使用这些路径：

- 二进制：`/opt/netconfighub/netconfighub`
- 工作目录：`/opt/netconfighub`
- 配置文件：`/etc/netconfighub/config.yaml`

安装并启用：

```bash
sudo useradd --system --home /opt/netconfighub --shell /usr/sbin/nologin nch || true
sudo install -d -o nch -g nch /opt/netconfighub /etc/netconfighub /var/lib/netconfighub
sudo install -m 0755 netconfighub /opt/netconfighub/netconfighub
sudo install -m 0644 configs/config.yaml /etc/netconfighub/config.yaml
sudo cp deploy/netconfighub.service /etc/systemd/system/netconfighub.service
sudo systemctl daemon-reload
sudo systemctl enable --now netconfighub
sudo systemctl status netconfighub
```

## Docker

```bash
docker compose up --build
```

Compose 会将应用暴露在 `localhost:8080`，运行数据保存在 `./data`。

启动可选 MySQL 服务：

```bash
docker compose --profile mysql up --build
```

## 配置

`configs/config.yaml` 中的重要配置：

- `server.host`、`server.port`：HTTP 监听地址。
- `server.jwt_secret`：JWT 签名密钥，生产环境必须替换。
- `server.encryption_key`：设备凭据加密密钥，生产环境必须替换。
- `server.cors.allowed_origins`：允许访问 API 的浏览器来源。生产环境建议使用明确的 HTTPS 来源，不要使用 `*`。
- `database.driver`：`sqlite` 或 `mysql`。
- `database.sqlite_path`：SQLite 数据库路径。
- `database.mysql_dsn`：MySQL DSN。
- `git.repo_path`：保存设备配置的本地 Git 仓库路径。
- `scheduler.*`：备份周期、并发 worker、重试和超时设置。
- `sanitize.*`：采集配置中的敏感信息脱敏规则。
- `notify.*`：Webhook 通知设置。

自动化 e2e 测试使用 `configs/config.e2e.yaml`，运行数据写入 `/tmp`。

## 开发和测试

后端测试：

```bash
go test ./...
go vet ./...
```

前端构建：

```bash
cd web
npm ci
npm run build
```

前端工具链要求 Node.js 20.19 或更新版本。

Playwright e2e 测试：

```bash
cd web
npm ci
npx playwright test
```

综合 API 冒烟测试：

```bash
rm -f /tmp/netconfighub-e2e.db
rm -rf /tmp/netconfighub-e2e-configs
NCH_CONFIG_PATH=configs/config.e2e.yaml go run cmd/api/main.go
bash test_api_v2.sh
```

`web/dist` 会被 `web/embed.go` 内嵌进 Go 二进制，因此会随仓库提交。

从 macOS 或 Linux 直接构建 Linux 二进制：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/netconfighub-linux-amd64 ./cmd/api
```

## 发布

当前发布标签：

```text
v0.1.1
```

发布内容和验证记录见 `RELEASE_NOTES.md`。

# NetConfigHub产品需求文档

**版本**：v0.01

**日期**：2026-05-12

**文档状态**：可执行生产版

**适用对象**：产品、研发、测试、运维、安全、交付

---

## 1. 文档目的

定义 NetConfigHub 的核心功能和实现边界，支持 Cisco/H3C/Huawei/Ruijie 网络设备配置备份与管理，确保生产环境可用。

---

## 2. 产品定位

NetConfigHub 是面向网络运维团队的 **多厂商自动化设备管理与配置备份平台**，核心能力包括：

- 自动采集网络设备配置
- 配置版本化管理（Git）
- 配置变化 Diff 与历史回溯
- 支持手动和定时备份
- 提供 Web 页面与 REST API
- 失败告警与操作审计
- 单文件部署，支持 systemd
- 支持 设备厂商，Cisco/H3C/Huawei/Ruijie 等
- 配置基线管理：可设定基线配置，自动比对 deviations

---

## 3. 总体架构

```text
Web UI / API
     |
Go 后端
     |
+----+----------------------+
|                           |
Device Service           Scheduler / Worker Pool
|                           |
Collector (SSH/Telnet)    Git Writer
|                           |
DeviceModel                Backup Records
|                           |
Store (Git/File)           Alert / Webhook
|
SQLite/MySQL
```

五层核心模块：Device → Collector → DeviceModel → Store → Alert

---

## 4. 用户角色与场景

| 角色 | 核心功能 | 核心场景 |
|------|----------|----------|
| 网络管理员 | 设备管理、配置标准化、基线管理 | 导入设备、分组管理、设置基线、触发备份 |
| 运维工程师 | 自动备份、失败管理、基线偏差检查 | 周期/手动备份、查看失败设备、偏差告警 |
| 审计/合规 | 配置历史、差异追踪、基线审核 | 历史版本、Diff、基线比对报告 |
| API 集成方 | 系统调用 | Webhook/REST API 与现有系统集成 |

---

## 5. 功能模块

### 5.1 设备管理

- 功能：设备新增、编辑、删除、启用/禁用、分组管理
- 字段定义：

| 字段            | 类型     | 必填 | 说明                          |
| --------------- | -------- | ---- | ----------------------------- |
| id              | int      | 是   | 唯一标识                      |
| name            | string   | 是   | 设备名称，唯一                |
| ip              | string   | 是   | 管理 IP                       |
| port            | int      | 是   | 端口                          |
| protocol        | enum     | 是   | ssh / telnet                  |
| vendor          | enum     | 是   | cisco / huawei / h3c / ruijie |
| model           | enum     | 是   | ios / vrp / comware / rg-os   |
| username        | string   | 是   | 登录用户名                    |
| password        | string   | 是   | 登录密码，加密保存            |
| enable_password | string   | 否   | 特权密码                      |
| group_name      | string   | 否   | 分组                          |
| enabled         | bool     | 是   | 是否启用                      |
| created_at      | datetime | 是   | 创建时间                      |
| updated_at      | datetime | 是   | 更新时间                      |

- 功能说明：
  1. 设备管理支持 CRUD。
  2. 启用设备参与采集，禁用设备不采集。
  3. 新增/编辑时校验 IP 和名称唯一性。

---

### 5.2 采集调度与 Worker

- Scheduler 定时触发采集任务
- 支持 Worker Pool 并发采集
- 每台设备采集状态流：pending → running → success / failed → retrying / timeout / cancelled
- 失败任务根据重试策略自动重试

---

### 5.3 Collector（SSH / Telnet）

- 功能：与设备建立连接，执行采集命令，获取配置
- 支持 SSH/Telnet 协议
- 支持 enable 模式（如 Cisco）
- 支持分页控制命令（关闭分页）
- 支持超时和异常处理

---

### 5.4 设备模型（DeviceModel）

- 支持品牌及模型：Cisco IOS、Huawei VRP、H3C Comware、Ruijie OS
- 每个模型包括：
  - 关闭分页命令
  - 配置采集命令
  - 退出命令
  - prompt 正则
  - 脱敏规则
- 示例（Cisco IOS）：

```yaml
model: cisco_ios
vendor: cisco
protocols:
  - ssh
  - telnet
disable_paging:
  - terminal length 0
  - terminal width 0
commands:
  - show running-config
exit:
  - exit
prompt_regex: '[>#]\s*$'
secret_rules:
  - pattern: '^enable secret .+$'
    replace: 'enable secret <removed>'
  - pattern: '^username (\S+) password .+$'
    replace: 'username $1 password <removed>'
```

---

### 5.5 配置存储（Store）

- 支持 Git 存储（生产推荐）和文件存储
- 保存配置变更，支持 Diff、历史回溯
- Git 写入必须串行，避免并发冲突
- 支持自动 push 到远程 Git 仓库

---

### 5.6 配置 Diff

- 提供设备配置差异对比
- 支持最新两次版本对比或任意两个 commit 对比
- API / 页面可展示 unified 或左右对比
- Diff 结果可下载

---

### 5.7 Web/API

- Web 页面功能：
  - 登录
  - 仪表盘
  - 设备管理
  - 配置历史与 Diff
  - 任务列表与手动触发
- REST API 功能：
  - /api/devices CRUD
  - /api/devices/{id}/backup
  - /api/devices/{id}/config
  - /api/devices/{id}/history
  - /api/devices/{id}/diff
  - /api/jobs、/api/jobs/{id}
- API 明确请求参数、返回数据、错误码

---

### 5.8 告警系统（Alert / Hook）

- Webhook 告警，支持失败任务、配置变化通知
- 支持自定义 HTTP 地址
- 支持邮件/企业消息通知（可选扩展）

---

### 5.9 审计日志

- 用户操作审计（新增设备、删除设备、触发备份、修改配置等）
- 保存操作用户、时间、IP、目标对象、动作类型

---

## 6. 数据库设计

- devices、backup_jobs、backup_records、users、audit_logs
- 字段定义详细，含必填、默认值、唯一约束、索引
- 支持 SQLite 默认，MySQL 可选
- 数据库迁移使用 Goose
- SQLC 管理查询和 CRUD

---

## 7. 部署方案

- 单二进制部署，系统服务通过 systemd 启动
- 配置文件使用 YAML
- 日志文件记录采集、任务、Git、告警和审计日志
- 支持 Docker 部署（可选）

---

## 8. 安全要求

- 密码、enable 密码加密存储
- API 查询不返回明文密码
- Git 仓库中不保存明文密码
- 配置脱敏选项（remove_secret）
- SSH 支持 Key 登录
- 支持 TLS/HTTPS API

---

## 9. 技术选型与架构合理性

| 层级 | 技术与说明 |
|------|------------|
| 后端 | Go 1.22+，高性能，多线程任务调度 |
| 前端 | Vue3 + TypeScript + Naive UI SPA |
| 数据库 | SQLite 默认，MySQL 8.0+ 可选 |
| 配置仓库 | Git/Git-Crypt，支持版本化、Diff、敏感信息加密 |
| API | RESTful，支持触发备份、获取配置、Diff、基线偏差 |
| 设备连接 | SSH/Telnet/SCP/HTTP 扩展，支持 Cisco/H3C/Huawei/Ruijie |
| 任务调度 | 内置 Worker Pool，支持周期和手动触发 |
| Hook | Exec / Git Push / Webhook，事件驱动 |
| 部署 | 单文件/二进制 + systemd + Docker，生产可用 |

---

## 10. 功能优先级（最小化实现核心功能）

| 优先级 | 功能 |
|--------|------|
| P0 | 单用户登录、设备管理、周期/手动备份、Git 存储、Web UI、失败设备列表、API Token、配置基线管理 |
| P1 | 配置 Diff、脱敏、基础告警通知、批量命令、插件扩展 |

---

## 11. MVP 不做

- 多租户 / RBAC 权限
- 自动下发配置
- 分布式采集
- AI 分析或报表生成
- SNMP 扫描 / CMDB 集成

---

## 12. 验收标准

- 新增设备可成功保存并参与采集
- SSH/Telnet 采集正常工作
- 配置变更自动生成 Git commit
- 无变化不提交 Git
- 手动备份任务立即生效
- 失败任务按重试策略处理并触发告警
- API 请求返回符合设计
- 页面展示设备列表、历史、Diff、任务状态
- 密码安全、脱敏规则符合要求
- 系统可通过 systemd 启动并稳定运行

---

# 附录

### A.1 支持厂商及模型列表

| 品牌 | 型号 / OS | 支持协议 | 采集命令 | 分页关闭命令 | 退出命令 |
|---|---|---|---|---|---|
| Cisco | IOS | SSH / Telnet | show running-config | terminal length 0 / terminal width 0 | exit |
| Huawei | VRP | SSH / Telnet | display current-configuration | screen-length 0 temporary | quit |
| H3C | Comware | SSH / Telnet | display current-configuration | screen-length disable | quit |
| Ruijie | RG-OS | SSH / Telnet | show running-config | terminal length 0 | exit |

### A.2 Git 并发写入设计

- Worker 并发采集
- Git 写入必须串行
- Worker Pool -> Collect Result Queue -> Git Writer (单线程) -> backup_records

### A.3 YAML 模型示例

详见 [5.4 设备模型（DeviceModel）](#54-设备模型devicemodel) 中的 Cisco IOS 示例。

### A.4 任务状态流

| 状态 | 描述 |
|---|---|
| pending | 等待执行 |
| running | 正在执行 |
| success | 执行成功 |
| failed | 执行失败 |
| retrying | 等待重试 |
| timeout | 超时 |
| cancelled | 取消 |
| skipped | 设备禁用或条件不满足 |

### A.5 错误码及失败分类

| 错误码 | 类型 | 描述 |
|---|---|---|
| 40001 | 参数错误 | API 请求参数不合法 |
| 40002 | 设备名称重复 | 新增设备名称已存在 |
| 40003 | IP 格式错误 | IP 地址无效 |
| 40101 | 未登录 | API 需要认证 |
| 50001 | 数据库错误 | SQLite/MySQL 操作失败 |
| 50002 | SSH 连接失败 | 无法建立 SSH 连接 |
| 50003 | Telnet 连接失败 | 无法建立 Telnet 连接 |
| 50004 | 命令执行失败 | 采集命令执行异常 |
| 50005 | Git 提交失败 | 配置写入 Git 失败 |

### A.6 配置采集异常分类

| 类型 | 描述 |
|---|---|
| connect_timeout | 连接超时 |
| auth_failed | 认证失败 |
| prompt_timeout | 等待提示符超时 |
| command_timeout | 命令执行超时 |
| model_not_found | 设备模型不存在 |
| git_commit_failed | Git 提交失败 |
| config_empty | 获取配置为空 |
| unknown_error | 未知错误 |

### A.7 API 响应示例

- 成功示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "core-sw-01",
    "ip": "192.168.1.1"
  }
}
```
- 错误示例
```json
{
  "code": 40002,
  "message": "设备名称重复"
}
```

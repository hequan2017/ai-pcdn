# ai-pcdn

> Modified for ai-pcdn on 2026-07-21: 合并产品、业务、开发与部署文档。

ai-pcdn 是基于 Go、Gin、Vue 3 和 Vite 构建的 PCDN（P2P CDN）节点供应商管理平台，面向 PCDN 节点供应商和带宽聚合商，统一管理节点接入、流量采集、95 值计算、监控告警及个人门户。

设计文档：[PCDN 系统设计](docs/superpowers/specs/2026-07-21-pcdn-system-design.md)

## 项目定位

- **采购侧**：部署服务器并接入个人或渠道带宽，支持包月、95 计费等结算模式。
- **销售侧**：管理节点接入平台及运行质量，为后续销售对账提供数据基础。
- **运营侧**：关注节点带宽、在线率、流量、95 值和告警，降低异常导致的结算损失。

## 系统架构

![ai-pcdn 系统架构](docs/architecture/ai-pcdn-architecture.svg)

系统由运营后台、个人门户和采集 Agent 三类入口接入 Gin API，核心业务服务统一处理节点、流量、95 值、告警、结算、对账、利润分析和 Agent OTA，并使用 SQLite 持久化数据。定时任务负责离线检查、95 值计算、月账单和告警检查，告警通过钉钉或企业微信 Webhook 通知。

[打开交互式架构图](docs/architecture/ai-pcdn-architecture.html) · [查看架构图源文件](docs/architecture/ai-pcdn.architecture.json)

## 已完成功能

### 数据底座与自助上机

| 模块 | 功能 |
| --- | --- |
| 节点管理 | 节点增删改查、批量删除、在线状态、归属、地域、运营商、接入平台、分组标签和数据权限 |
| 采集 Agent | 每分钟采集网卡峰值、本地 JSONL 持久化、失败重试、30 秒心跳、首次激活和安装命令 |
| 路由与鉴权 | `admin` 使用 JWT、Casbin 和 DataScope；`agent` 使用节点 Token；`portal` 使用个人 JWT |
| 流量与 95 值 | 流量点幂等写入，支持日滚动 95 值和月冻结 95 值 |
| 定时任务 | 节点离线判定、日 95 值滚动计算和月 95 值冻结 |
| 运营后台 | 节点管理、流量曲线和 95 值查看 |
| 个人门户 | 注册、登录、自助添加节点、生成凭证和安装命令、查看节点、流量及账单 |

### 监控告警

| 模块 | 功能 |
| --- | --- |
| 告警规则 | 节点离线、带宽低于阈值、95 值高于阈值和 Agent 上报中断 |
| 告警范围 | 支持全部节点、节点分组和单节点 |
| 告警引擎 | 周期检查、触发与恢复、同规则与节点的告警收敛去重 |
| 通知 | 钉钉和企业微信 Webhook，支持 Markdown 与手机号提醒 |
| 运营后台 | 告警规则增删改查、启停和告警记录查看 |

### 采购结算

| 模块 | 功能 |
| --- | --- |
| 账单生成 | 按账期和贡献者分组，支持包月和 95 计费，月初自动生成上月账单 |
| 账单审核 | 草稿、已审核、已付款、已驳回状态流转 |
| 付款流程 | 记录付款方式、流水号、实付金额和操作人 |
| 个人门户 | 贡献者查看自己的账单 |
| 运营后台 | 账单管理、审核、付款和明细查看 |

### 销售对账

| 模块 | 功能 |
| --- | --- |
| 结算单导入 | 录入大厂结算单，按节点 SN 自动关联 |
| 自动核对 | 用自采集月 95 流量与大厂数据比对，标记一致或差异（阈值 10%） |
| 应收汇总 | 按账期和平台汇总应收收入与核对状态 |
| 运营后台 | 结算单管理、重新核对和应收汇总卡片 |

### 利润大盘

| 模块 | 功能 |
| --- | --- |
| 利润汇总 | 大厂收入减采购成本，展示利润和利润率 |
| 多维明细 | 按平台收入、按贡献者成本 |
| 月度趋势 | 近 6 月收入、成本和利润趋势 |

### Agent OTA 与批量部署

| 模块 | 功能 |
| --- | --- |
| 版本发布 | 管理多版本，标记稳定版和强制升级 |
| 自升级 | Agent 每小时检查最新版本，下载、SHA256 校验、原地替换并重启 |
| 运营后台 | 版本发布、编辑和删除 |

## 技术栈

- **后端**：Go、Gin、GORM、Casbin、JWT
- **数据存储**：SQLite（当前默认，本地和一键部署均无需外部数据库）
- **运营后台**：Vue 3、Vite、Element Plus、UnoCSS、ECharts
- **个人门户**：Vue 3、Vite、Element Plus、Vue Router、Axios
- **采集 Agent**：Go 单二进制，可独立部署

## 目录结构

```text
ai-pcdn/
├── server/
│   ├── plugin/pcdn/             # PCDN 后端核心业务
│   │   ├── model/               # 节点、流量、95 值、告警、账单、结算、版本模型
│   │   ├── service/             # 业务服务、告警引擎、95 值、账单、对账、利润和通知
│   │   ├── api/                 # admin、agent、portal 接口
│   │   ├── router/              # 三组业务路由
│   │   ├── middleware/          # Agent Token 鉴权
│   │   └── initialize/          # 建表、菜单、API 权限点、Casbin/菜单授权、路由和定时任务
│   └── cmd/pcdn-agent/          # 独立采集 Agent（含 OTA 自升级）
├── web/                         # ai-pcdn 运营后台
│   └── src/plugin/pcdn/         # 节点、告警、账单、对账、利润和版本管理页面
├── site/                        # 独立个人门户
├── deploy/docker-compose/       # Docker Compose 配置
├── docs/superpowers/specs/      # 设计文档
└── deploy.sh                    # 一键部署脚本
```

## 一键部署

### 环境要求

- Docker 24 或更高版本
- Docker Compose v2（使用 `docker compose` 命令）
- `curl`
- Linux、macOS、WSL 或 Git Bash

部署前请启动 Docker，并确保目标机器的 `8080` 和 `8888` 端口未被占用。

### 执行部署

首次部署建议显式设置管理员密码：

```bash
AI_PCDN_ADMIN_PASSWORD='请替换为至少6位的强密码' bash ./deploy.sh
```

本地临时验证也可以直接执行：

```bash
bash ./deploy.sh
```

未设置密码时，初始密码为 `123456`。首次登录后应立即修改。

脚本会执行：

1. 检查 Docker、Docker Compose 和 `curl`。
2. 创建 `deploy/docker-compose/runtime/config.yaml`。
3. 构建并启动前后端容器。
4. 等待后端健康检查通过。
5. 首次部署时初始化 SQLite 数据库和基础数据。
6. 验证前端页面及 `/api` 代理。

重复执行脚本会复用已有配置和数据库，不会重新初始化或删除数据。

部署成功后访问：

- 运营后台：<http://127.0.0.1:8080>
- 后端服务：<http://127.0.0.1:8888>
- Swagger：<http://127.0.0.1:8888/swagger/index.html>
- 默认管理员账号：`admin`

当前一键部署包含后端和运营后台；`site/` 个人门户需要按“本地开发”章节单独启动或自行构建发布。

### 自定义参数

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `AI_PCDN_ADMIN_PASSWORD` | `123456` | 首次初始化时的管理员密码，至少 6 位 |
| `AI_PCDN_WEB_PORT` | `8080` | 运营后台宿主机端口 |
| `AI_PCDN_SERVER_PORT` | `8888` | 后端宿主机端口 |
| `AI_PCDN_DB_NAME` | `ai-pcdn` | 首次初始化时的 SQLite 数据库名 |

```bash
AI_PCDN_WEB_PORT=18080 \
AI_PCDN_SERVER_PORT=18888 \
AI_PCDN_ADMIN_PASSWORD='ChangeMe_2026' \
bash ./deploy.sh
```

## 容器运维

查看运行状态：

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml ps
```

查看日志：

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml logs -f --tail=200
```

重新构建并升级：

```bash
bash ./deploy.sh
```

停止服务但保留数据：

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml down
```

SQLite 数据存储在 Docker 卷 `ai-pcdn_server-data`，上传文件存储在 `ai-pcdn_uploads`。不要执行带 `-v` 的 `docker compose down`，否则会删除持久化数据。

## 本地开发

当前 `server/config.yaml` 已配置为 SQLite：数据库名为 `gva`、路径为 `server/`，Redis 默认关闭。

### 后端

```bash
cd server
go run .
```

后端默认地址：<http://127.0.0.1:8888>

### 运营后台

```bash
cd web
npm install
npm run dev
```

运营后台默认地址：<http://127.0.0.1:8080>

### 个人门户

```bash
cd site
npm install
npm run dev
```

个人门户默认地址：<http://localhost:5174>，开发服务器会将 `/pcdn` 请求代理到 `http://127.0.0.1:8888`。

### 采集 Agent

```bash
cd server
go build -o pcdn-agent ./cmd/pcdn-agent/
```

节点获取 SN 和 Token 后运行：

```bash
./pcdn-agent -server http://<后端地址> -sn <节点SN> -token <Token>
```

## 上机闭环

1. 在个人门户注册并登录。
2. 添加节点并填写地域、运营商、接入平台等信息。
3. 获取节点 SN、Token 和安装命令。
4. 在节点服务器启动 Agent。
5. Agent 自动激活并持续上报流量。
6. 在运营后台查看在线状态、流量、95 值和告警。

## API 概览

| 路由组 | 鉴权 | 用途 |
| --- | --- | --- |
| `/pcdn/admin/*` | JWT、Casbin、DataScope | 节点、流量和告警管理 |
| `/pcdn/agent/*` | `X-Node-Sn`、`X-Node-Token` | Agent 激活、上报和心跳 |
| `/pcdn/portal/*` | 公开接口或个人 JWT | 注册、登录、添加节点及个人数据查询 |

## 数据模型

| 表 | 说明 |
| --- | --- |
| `gva_pcdn_node` | 节点主表，包含凭证、归属、位置、平台、状态和计费信息 |
| `gva_pcdn_node_iface` | 节点网卡 |
| `gva_pcdn_node_traffic_point` | 流量分钟峰值点，按节点、窗口和网卡保证幂等 |
| `gva_pcdn_node_95` | 日滚动和月冻结 95 值 |
| `gva_pcdn_alarm_rule` | 告警规则 |
| `gva_pcdn_alarm_record` | 告警触发与恢复记录 |
| `gva_pcdn_bill` | 采购侧月账单（按贡献者汇总，含明细与付款信息） |
| `gva_pcdn_settlement` | 大厂结算单（销售侧应收依据，含核对状态） |
| `gva_pcdn_agent_release` | Agent 版本发布（OTA 升级源） |

## 关键设计决策

- **默认存储**：使用本地 SQLite，降低初始化和部署依赖。
- **采集粒度**：每分钟保存该分钟峰值，兼顾数据量与 95 值精度。
- **可靠上报**：本地持久化后上报，失败由独立任务重试。
- **幂等写入**：节点、时间窗口和网卡组成唯一约束，重复上报不产生重复点。
- **通信方式**：Agent 主动推送，适配 NAT 和家庭宽带场景。
- **95 值**：按周期内分钟峰值排序，去除最高 5% 后取最大值；支持日滚动和月冻结。
- **数据隔离**：个人门户按 `owner_user_id` 限制节点和流量访问范围。
- **告警收敛**：同规则和节点在 firing 期间只通知一次，恢复时再通知一次。
- **插件自动授权**：pcdn 菜单与 API 在初始化时自动授权给超级管理员(888)，增量写入 `sys_authority_menus` 与 `casbin_rule` 并刷新 Casbin enforcer；不使用破坏性的 `SetMenuAuthority`/`UpdateCasbin`（会清空 admin 已有权限）。admin 无需在后台手动分配即可访问，且重启幂等不产生重复记录。

## 验证

```bash
cd server
go test ./...
go build ./...

cd ../web
npm run build

cd ..
docker compose -f deploy/docker-compose/docker-compose.yaml config
```

2026-07-21 本地验证结果：

- 后端 `go build ./...` 全量通过（含采集 Agent），SQLite 健康接口返回 `200 / "ok"`。
- PCDN 核心单测 `go test ./plugin/pcdn/service` 通过：95 分位算法、95 值计算、流量幂等、告警收敛、账单生成（包月+p95+幂等）。
- 菜单（8）与 API（35）均已在 `initialize` 注册，并通过 `initialize/auth.go` 自动授权给超级管理员(888)，写入 `sys_authority_menus` 与 `casbin_rule` 并刷新 Casbin。
- 运营后台和个人门户生产构建通过，三个本地服务均可访问。
- 全量 `go test ./...` 仍存在 MCP 会话、AI Markdown 渲染、自动路由和模板测试失败；这些问题不阻塞当前服务启动，但发布前应单独修复。

## 后续路线图

| 阶段 | 内容 | 状态 |
| --- | --- | --- |
| 阶段 1 | 数据底座与自助上机 | 已完成 |
| 阶段 2 | 监控告警 | 已完成 |
| 阶段 3 | 采购结算：个人月账单和付款流程 | 已完成 |
| 阶段 4 | 销售对账：平台结算单导入和核对 | 已完成 |
| 阶段 5 | 利润大盘：收入与成本多维报表 | 已完成 |
| 阶段 6 | Agent OTA 和批量部署 | 已完成 |

## 授权说明

部署和使用 ai-pcdn 时，请遵循项目实际授权约定，并妥善保管商业授权凭证。

[简体中文](README.md) | [English](README.en.md)

# ai-pcdn

A management platform for PCDN (P2P CDN) node suppliers, built with Go, Gin, Vue 3 and Vite. Designed for PCDN node suppliers and bandwidth aggregators, it unifies node onboarding, traffic collection, 95th-percentile calculation, monitoring & alerting, procurement settlement, sales reconciliation, profit analytics, and Agent OTA.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vuedotjs&logoColor=white)
![SQLite](https://img.shields.io/badge/Storage-SQLite-003B57?logo=sqlite&logoColor=white)

## Introduction

The daily reality of running a PCDN business: you purchase residential/edge bandwidth from individuals and channels, onboard those nodes into major platforms for revenue — but the nodes are scattered and ownership is messy, bandwidth usage is whatever the platform dashboards say, mismatches in 95th-percentile billing translate directly into lost money, and there is no lightweight tool that ties procurement, operations, settlement, reconciliation, and profit together.

ai-pcdn packages that entire chain into a ready-to-run platform:

- **Procurement side**: deploy servers and onboard individual or channel bandwidth, with monthly flat-rate and 95th-percentile billing modes;
- **Sales side**: manage node onboarding and runtime quality, import platform settlement sheets and reconcile them automatically, providing the data foundation for sales reconciliation;
- **Operations side**: keep an eye on node bandwidth, online rate, traffic, 95th percentile, and alarms to reduce settlement losses caused by anomalies.

The system consists of three entry points — an operations console (a customized gin-vue-admin plugin), a personal portal (for contributors to onboard nodes and view bills), and a collection agent (a single Go binary). Data is persisted in SQLite by default; neither local runs nor the one-command deployment require any external database.

Design document: [PCDN System Design](docs/superpowers/specs/2026-07-21-pcdn-system-design.md)

## ✨ Features

### Data Foundation & Self-Service Onboarding

| Module | Capabilities |
| --- | --- |
| Node management | Node CRUD, batch deletion, online status, ownership, region, ISP, onboarding platform, group tags, and data-scope permissions |
| Collection agent | Per-minute NIC peak sampling, local JSONL persistence, retry on failure, 30-second heartbeat, first-time activation, and install commands |
| Routing & auth | `admin` uses JWT + Casbin + DataScope; `agent` uses node tokens (`X-Node-Sn` / `X-Node-Token`); `portal` uses personal JWTs |
| Traffic & 95th percentile | Idempotent traffic-point writes (unique per node + window + NIC), daily rolling 95th percentile and monthly frozen 95th percentile |
| Scheduled tasks | Node offline detection (heartbeat over 3 minutes old), daily rolling 95 calculation, monthly 95 freezing |
| Personal portal | Register, log in, self-service node onboarding, generate credentials and install commands, view own nodes, traffic, and bills |

### Monitoring & Alerting

| Module | Capabilities |
| --- | --- |
| Alarm rules | Node offline, bandwidth below threshold, 95th percentile above threshold, agent reporting interruption |
| Alarm scope | All nodes, node groups, or a single node |
| Alarm engine | Periodic checks, firing & recovery, deduplication/convergence per rule + node |
| Notifications | DingTalk and WeCom (Enterprise WeChat) webhooks, with Markdown and @mobile support |

### Procurement Settlement

| Module | Capabilities |
| --- | --- |
| Bill generation | Grouped by billing period and contributor, supports flat-rate and 95 billing; previous month's bills are generated automatically at month start |
| Bill review | Draft → reviewed → paid → rejected status flow |
| Payment workflow | Records payment method, transaction ID, actual amount, and operator |
| Personal portal | Contributors can view their own bills |

### Sales Reconciliation

| Module | Capabilities |
| --- | --- |
| Settlement import | Enter platform (vendor) settlement sheets, auto-associated by node SN |
| Auto reconciliation | Compares self-collected monthly 95 traffic against the platform's numbers, marking matches or discrepancies (10% threshold) |
| Receivables summary | Aggregates receivable revenue and reconciliation status by period and platform |

### Profit Dashboard

| Module | Capabilities |
| --- | --- |
| Profit summary | Platform revenue minus procurement cost, showing profit and margin |
| Multi-dimensional details | Revenue by platform, cost by contributor |
| Monthly trends | Revenue, cost, and profit trends over the last 6 months |

### Agent OTA & Bulk Deployment

| Module | Capabilities |
| --- | --- |
| Release management | Manage multiple versions, mark stable and force-upgrade releases |
| Self-upgrade | The agent checks for the latest version hourly, downloads it, verifies SHA256, replaces itself in place, and restarts (Linux / Windows) |

## 🏗 Architecture

![ai-pcdn architecture](docs/architecture/ai-pcdn-architecture.svg)

Three entry points — the operations console, the personal portal, and the collection agent — connect to the Gin API. Core business services uniformly handle nodes, traffic, 95th percentile, alarms, settlement, reconciliation, profit analytics, and Agent OTA, with SQLite persistence. Scheduled tasks handle offline detection, 95 calculation, monthly billing, and alarm checks; alarms are delivered via DingTalk or WeCom webhooks.

[Open the interactive architecture diagram](docs/architecture/ai-pcdn-architecture.html) · [View the diagram source](docs/architecture/ai-pcdn.architecture.json)

Key design decisions:

- **Default storage**: local SQLite, minimizing setup and deployment dependencies;
- **Sampling granularity**: the per-minute peak is stored every minute, balancing data volume and 95th-percentile accuracy;
- **Reliable reporting**: data is persisted locally before being reported, with a dedicated retry task; duplicate reports never create duplicate points;
- **Communication model**: the agent pushes proactively, fitting NAT and home-broadband environments;
- **95th-percentile algorithm**: minute peaks within a period are sorted ascending, the top 5% is dropped, and the maximum of the remainder is taken; supports daily rolling and monthly frozen values;
- **Data isolation**: the personal portal restricts node and traffic access by `owner_user_id`;
- **Alarm convergence**: each rule + node notifies once while firing and once on recovery;
- **Plugin auto-authorization**: pcdn menus and APIs are authorized automatically to the super admin (888) at initialization, incrementally written to `sys_authority_menus` and `casbin_rule` with the Casbin enforcer refreshed; restarts are idempotent.

## 🛠 Tech Stack

| Layer | Technologies |
| --- | --- |
| Backend | Go 1.24, Gin, GORM, Casbin, JWT, gin-vue-admin plugin mechanism |
| Storage | SQLite (current default; no external database needed for local runs or one-command deployment) |
| Operations console | Vue 3, Vite, Element Plus, UnoCSS, ECharts |
| Personal portal | Vue 3, Vite, Element Plus, Vue Router, Axios |
| Collection agent | Single Go binary, independently deployable (Linux / Windows) |

## 🚀 Quick Start

### One-command deployment (Docker Compose)

Requirements: Docker 24+, Docker Compose v2, `curl`, on Linux / macOS / WSL / Git Bash; make sure ports `8080` and `8888` are free.

For the first deployment, set the admin password explicitly:

```bash
AI_PCDN_ADMIN_PASSWORD='replace-with-a-strong-password-of-6-plus-chars' bash ./deploy.sh
```

For a quick local trial you can also just run:

```bash
bash ./deploy.sh
```

If no password is set, the initial password is `123456` — change it immediately after the first login.

The script: checks Docker / Compose / curl → generates `deploy/docker-compose/runtime/config.yaml` → builds and starts the frontend and backend containers → waits for the health check → initializes the SQLite database and seed data on first run → verifies the frontend page and the `/api` proxy. Re-running it reuses the existing configuration and database; nothing is re-initialized or deleted.

After a successful deployment:

- Operations console: <http://127.0.0.1:8080>
- Backend service: <http://127.0.0.1:8888>
- Swagger: <http://127.0.0.1:8888/swagger/index.html>
- Default admin account: `admin`

> The one-command deployment currently includes the backend and the operations console; the `site/` personal portal must be started separately (see "Local development") or built and published on your own.

### Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `AI_PCDN_ADMIN_PASSWORD` | `123456` | Admin password used at first initialization, at least 6 characters |
| `AI_PCDN_WEB_PORT` | `8080` | Host port for the operations console |
| `AI_PCDN_SERVER_PORT` | `8888` | Host port for the backend |
| `AI_PCDN_DB_NAME` | `ai-pcdn` | SQLite database name used at first initialization |

```bash
AI_PCDN_WEB_PORT=18080 \
AI_PCDN_SERVER_PORT=18888 \
AI_PCDN_ADMIN_PASSWORD='ChangeMe_2026' \
bash ./deploy.sh
```

### Container operations

```bash
# Check running status
docker compose -f deploy/docker-compose/docker-compose.yaml ps

# View logs
docker compose -f deploy/docker-compose/docker-compose.yaml logs -f --tail=200

# Rebuild and upgrade
bash ./deploy.sh

# Stop services but keep data
docker compose -f deploy/docker-compose/docker-compose.yaml down
```

SQLite data lives in the Docker volume `ai-pcdn_server-data`, and uploaded files in `ai-pcdn_uploads`. Never run `docker compose down` with `-v`, or the persisted data will be deleted.

### Local development

`server/config.yaml` defaults to SQLite (database `gva`, path `server/`), with Redis disabled.

```bash
# Backend (default http://127.0.0.1:8888)
cd server && go run .

# Operations console (default http://127.0.0.1:8080)
cd web && npm install && npm run dev

# Personal portal (default http://localhost:5174; /pcdn requests are proxied to 127.0.0.1:8888)
cd site && npm install && npm run dev

# Collection agent
cd server && go build -o pcdn-agent ./cmd/pcdn-agent/
```

After the node obtains its SN and Token, run the agent:

```bash
./pcdn-agent -server http://<backend-address> -sn <node-SN> -token <Token>
```

Common flags: `-ifaces` NICs to sample (comma-separated, empty = all non-lo), `-interval` sampling/reporting interval (default 60 seconds), `-store` local persistence path (default `/var/lib/pcdn-agent/pending.jsonl`).

### Onboarding workflow

1. Register and log in on the personal portal;
2. Add a node with region, ISP, onboarding platform, etc.;
3. Obtain the node SN, Token, and install command;
4. Start the agent on the node server;
5. The agent activates automatically and keeps reporting traffic;
6. Monitor online status, traffic, 95th percentile, and alarms in the operations console.

## 📁 Directory Structure

```text
ai-pcdn/
├── server/
│   ├── plugin/pcdn/             # PCDN backend core
│   │   ├── model/               # Node, traffic, 95, alarm, bill, settlement, release models
│   │   ├── service/             # Business services, alarm engine, 95, bills, reconciliation, profit, notifications
│   │   ├── api/                 # admin, agent, portal APIs
│   │   ├── router/              # Three business route groups
│   │   ├── middleware/          # Agent token authentication
│   │   └── initialize/          # Table creation, menus, API permissions, Casbin/menu authorization, routes, scheduled tasks
│   └── cmd/pcdn-agent/          # Standalone collection agent (with OTA self-upgrade)
├── web/                         # ai-pcdn operations console
│   └── src/plugin/pcdn/         # Node, alarm, bill, reconciliation, profit, and release pages
├── site/                        # Standalone personal portal
├── deploy/docker-compose/       # Docker Compose configuration
├── docs/superpowers/specs/      # Design documents
└── deploy.sh                    # One-command deployment script
```

## 📡 API Overview

| Route group | Auth | Purpose |
| --- | --- | --- |
| `/pcdn/admin/*` | JWT, Casbin, DataScope | Node, traffic, and alarm management |
| `/pcdn/agent/*` | `X-Node-Sn`, `X-Node-Token` | Agent activation, reporting, and heartbeat |
| `/pcdn/portal/*` | Public endpoints or personal JWT | Registration, login, node onboarding, personal data queries |

## 🗃 Data Model

| Table | Description |
| --- | --- |
| `gva_pcdn_node` | Node master table: credentials, ownership, location, platform, status, and billing info |
| `gva_pcdn_node_iface` | Node network interfaces |
| `gva_pcdn_node_traffic_point` | Per-minute traffic peak points, idempotent per node + window + NIC |
| `gva_pcdn_node_95` | Daily rolling and monthly frozen 95th-percentile values |
| `gva_pcdn_alarm_rule` | Alarm rules |
| `gva_pcdn_alarm_record` | Alarm firing and recovery records |
| `gva_pcdn_bill` | Procurement-side monthly bills (aggregated per contributor, with details and payment info) |
| `gva_pcdn_settlement` | Platform settlement sheets (sales-side receivables basis, with reconciliation status) |
| `gva_pcdn_agent_release` | Agent releases (source for OTA upgrades) |

## ✅ Verification

```bash
cd server
go test ./...
go build ./...

cd ../web
npm run build

cd ..
docker compose -f deploy/docker-compose/docker-compose.yaml config
```

Local verification results from 2026-07-21:

- Backend `go build ./...` passes fully (including the collection agent); the SQLite health endpoint returns `200 / "ok"`;
- Core PCDN unit tests `go test ./plugin/pcdn/service` pass: 95th-percentile algorithm, 95 calculation, traffic idempotency, alarm convergence, bill generation (flat-rate + p95 + idempotency);
- Menus (8) and APIs (35) are registered in `initialize` and auto-authorized to the super admin (888);
- The operations console and personal portal build successfully in production mode, and all three local services are reachable;
- The full `go test ./...` still has failures in MCP sessions, AI Markdown rendering, auto-routing, and template tests; these do not block the current service startup but should be fixed separately before release.

## 🗺 Roadmap

| Phase | Content | Status |
| --- | --- | --- |
| Phase 1 | Data foundation & self-service onboarding | Done |
| Phase 2 | Monitoring & alerting | Done |
| Phase 3 | Procurement settlement: personal monthly bills and payments | Done |
| Phase 4 | Sales reconciliation: settlement import and checking | Done |
| Phase 5 | Profit dashboard: multi-dimensional revenue/cost reports | Done |
| Phase 6 | Agent OTA and bulk deployment | Done |

## 📄 License

The LICENSE file in the repository is empty; no open-source license has been declared yet.

## 📌 Authorization Notice

When deploying and using ai-pcdn, please follow the project's actual licensing terms and keep your commercial authorization credentials safe.

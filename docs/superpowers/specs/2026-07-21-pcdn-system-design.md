# PCDN 管理系统设计文档

- 日期：2026-07-21
- 状态：已定稿，进入开发
- 范围：完整系统蓝图 + 分阶段路线图 + **阶段1（先上线）详细设计**

---

## 1. 背景与商业模式

- **角色**：PCDN 节点供应商 / 带宽聚合商
- **采购侧**：在全国部署服务器、接入个人/渠道带宽，按 **包月 / 95计费** 向贡献者付费
- **销售侧**：节点接入抖音、腾讯等大厂 PCDN，与大厂结算收款
- **盈利**：中间差价 / 服务费
- **规模**：> 500 台服务器，重自动化
- **计费数据来源**：自采集（agent 采集网卡流量），销售侧最终金额以大厂结算单为准（自采集用于核对）

## 2. 系统完整蓝图

| 子域 | 职责 |
|---|---|
| 宣传官网 | 平台介绍、营销页、价格、引导注册 |
| 个人自助接入 | 注册、登录、自助上机（生成凭证 + 一键安装） |
| 节点/服务器管理（A） | 设备台账、分组、地域、归属、状态、批量运维 |
| 流量采集与95值（B） | agent 采集 + 上报 + 95值计算 + 在线率 |
| 监控告警（E） | 在线率/带宽/95值/掉线 规则 → 通知收敛 |
| 采购结算（C） | 个人月账单生成 + 审核 + 付款流程 |
| 销售对账（D） | 大厂结算单导入 + 与自采集核对 + 应收 |
| 利润大盘（F） | 收入−成本=利润，多维报表 |
| 运维自动化 | agent OTA + 批量部署/配置下发 |
| 个人门户 H5（推迟） | 移动端自助查看节点/流量/账单 |

## 3. 分阶段路线图

| 阶段 | 内容 | 状态 |
|---|---|---|
| **阶段1 先上线** | 宣传官网 + 个人注册/登录/自助上机 + 节点管理 + 采集上报/95值 + 运营PC流量查看 | **本文档详写，立即开发** |
| 阶段2 | 监控告警 | 路线图，后续 spec |
| 阶段3 | 采购结算：个人月账单 + 审核 + 付款流程 | 路线图，后续 spec |
| 阶段4 | 销售对账：大厂结算单导入 + 核对 + 应收 | 路线图，后续 spec |
| 阶段5 | 利润大盘 | 路线图，后续 spec |
| 阶段6 | 运维自动化：agent OTA + 批量部署 | 路线图，后续 spec |
| 阶段7（推迟） | 个人门户 H5 | 路线图，后续 spec |

> 原则：**砍功能不砍 schema**。阶段1 把所有后续阶段需要的字段/表预留好，后续阶段直接扩展、不返工。

## 4. 整体架构

```
┌─────────────── 前端三端 ───────────────────────────────────────┐
│  site/   宣传官网 (Vue3+Vite)   │  web/  运营后台 (GVA, Vue3+EP) │  h5/ 个人门户(推迟)
│  营销页/注册/登录/上机           │  节点管理/流量/95值/后续结算    │
└────────────┬────────────────────┴────────────┬─────────────────┘
             │                                    │
             ▼                                    ▼
┌──────────────────────────────────────────────────────────────────┐
│  GVA 后端 (Go+Gin)  独立插件 server/plugin/pcdn/                    │
│  ┌─ /pcdn/portal/*  公开注册/登录 + 个人JWT(按owner隔离)            │
│  ├─ /pcdn/agent/*   agent上报/心跳/激活 (node token + 限流)        │
│  └─ /pcdn/admin/*   运营后台 (JWT+Casbin+DataScope)               │
└──────┬───────────────────┬─────────────────────┬──────────────────┘
       ▼                   ▼                     ▼
   MySQL(业务+流量)     Redis(心跳/限流/去重)   对象存储(可选)
```

```
>500 台 PCDN 服务器（NAT/家庭宽带后）
  采集 Agent (Go 单二进制，与火山/腾讯客户端共存)
  采集器(每分钟)─▶本地持久化─▶上报发送器 ─成功▶ 标记done
                     │      └─失败▶ 保留pending,本次结束
                     └──── 独立定时任务：重试检查器 ─▶ 扫pending重上报
  心跳：每30s
```

## 5. 用户角色与数据隔离

| 角色 | 入口 | 能力 | 数据范围 |
|---|---|---|---|
| 访客 | 官网 | 浏览营销页、注册 | — |
| 个人贡献者 | 官网登录 | 自助上机、看自己节点/流量/95值/账单 | `owner_user_id = 自己` |
| 运营人员 | 运营后台 | 管理所有节点、流量、95值、后续结算 | 全部（受 GVA 部门数据权限约束） |
| 管理员 | 运营后台 | 全部 + 权限/配置 | 全部 |

**数据隔离实现**：
- 个人贡献者 = GVA `sys_users` + Casbin「贡献者」角色
- `pcdn_node.owner_user_id` 关联 `sys_users.id`
- Service 层按角色过滤：运营不加过滤（走 GVA DataScope）；贡献者强制 `owner_user_id = current_user.id`
- API 保持前后端分离，阶段7 H5 直接复用 `/pcdn/portal/*`

## 6. 阶段1 详细设计

### 6.1 宣传官网（site/）

- 独立 Vue3 + Vite 应用，目录 `site/`
- 页面：首页（平台介绍/卖点）、特性、价格（包月/95计费说明）、关于、**注册**、**登录**、**控制台入口**（登录后跳个人上机/节点页，或直接复用 web/ 的个人视图）
- SEO 要求高时后续可换 Nuxt/SSG；阶段1 先 SPA
- 部署：独立静态站点（主域名根 www），后端走 `/pcdn/portal/*`

### 6.2 个人注册 / 登录 / 自助上机

**注册**：手机号/邮箱 + 密码（阶段1），预留短信验证码。注册即建 `sys_users` + 贡献者角色。

**自助上机流程**：
1. 个人用户在控制台点「添加节点」→ 后台预创建 `pcdn_node`（`status=pending`），生成 `node_sn` + `token`（token 仅明文返回一次，库存 hash）
2. 页面展示**一键安装命令**：`curl -fsSL <域名>/pcdn/portal/install/<token>.sh | bash`
3. 用户在自己服务器执行 → 脚本下载 agent 二进制、写入 token、注册 systemd 服务、启动
4. agent 首次上报携带 token → 后台校验 token、绑定 `owner_user_id`、回填硬件信息、`status=online`

**一机一凭证**：每次添加节点生成独立 token，便于管控与吊销。

### 6.3 节点管理（A）数据模型

**`pcdn_node`**（主表）：

| 分组 | 字段 | 说明 |
|---|---|---|
| 标识 | `id`,`node_sn`(unique),`token_hash` | node_sn 唯一；token 存 hash |
| 归属 | `owner_user_id`,`owner_name`,`contact` | 关联 sys_users |
| 位置 | `region`,`isp` | 省/市、运营商 |
| 接入大厂 | `platform`,`platform_node_id` | 销售（阶段4用，先留空） |
| 分组 | `group_id`,`tags`(JSON) | 批量管理 |
| 硬件 | `hostname`,`inner_ip`,`report_ip`,`os` | agent 激活时回填 |
| 状态 | `status`,`last_heartbeat_at`,`agent_version` | pending/online/offline/abnormal/disabled |
| 计费 | `billing_mode`(monthly/p95),`monthly_price`,`contract_period` | 阶段3用，先留字段 |
| 数据权限 | `dept_id` | GVA 行级权限 |
| 公共 | `created_at/updated_at/created_by/updated_by/deleted_at` | GVA 标准 |

**`pcdn_node_iface`**：`node_id`,`iface_name`,`mac`,`enabled` —— 采集维度（一节点多网卡）。

**批量运维（阶段1基础）**：Excel 导入（运营代录）、启停、改归属/分组。自助上机的节点由个人创建。

### 6.4 Agent（Go 单二进制）

- **采集**：每 1 分钟采一次网卡上下行瞬时速率，取该分钟峰值 `max`（保 95 值精度）
- **本地持久化**：写入本地 SQLite/文件，标记 `pending`
- **上报**：每分钟上报新点；失败则本次结束，不阻塞采集
- **重试**：**独立定时任务**扫描 `pending` 重上报（与采集/上报解耦）
- **心跳**：每 30s 上报心跳，携带硬件信息（首次激活时回填）
- **鉴权**：携带 `node_sn + token`（首次激活后改为持久 token）
- **OTA**：阶段6做；阶段1 手动/脚本部署
- **幂等**：后台 `UNIQUE(node_id, window_start, iface)` 去重

### 6.5 上报接入 API（/pcdn/agent/*）

- `POST /pcdn/agent/activate`：首次激活，token 校验 → 绑定 owner、回填硬件、转 online
- `POST /pcdn/agent/report`：批量上报流量点（鉴权 + Redis 限流 + 幂等写入）
- `POST /pcdn/agent/heartbeat`：心跳，刷新 `last_heartbeat_at` + Redis 在线标记

### 6.6 95 值计算

**`pcdn_node_95`**：`node_id`,`period_type`(day/month),`period_start`,`period_end`,`rx_95_bps`,`tx_95_bps`,`combined_95_bps`,`sample_count`,`status`(rolling/frozen),`frozen_at`

- 算法：周期内所有分钟峰值点升序，去掉最高 5% 后取最大值
- **每日凌晨定时任务**：算前一日各节点当日 95 值（`status=rolling`）
- **月底定时任务**：冻结当月 95 值（`status=frozen`），作为阶段3 账单依据
- 复用 GVA timed task 能力

### 6.7 运营后台流量查看（web/）

- 节点列表：在线状态、最新带宽、归属、地域、大厂
- 节点详情：实时/历史流量曲线（ECharts，读 `node_traffic_point`）、95 值（日/月）、心跳
- 数据权限：按部门隔离

## 7. 数据 schema（阶段1 + 后续预留）

**阶段1 表**：
- `pcdn_node`、`pcdn_node_iface`、`pcdn_node_traffic_point`（按月分区）、`pcdn_node_95`

**后续阶段预留表（阶段1 不建，路线图到位再建）**：
- `pcdn_owner_bill`（个人月账单，阶段3）、`pcdn_payment`（付款记录，阶段3）
- `pcdn_platform_settlement`（大厂结算单，阶段4）
- `pcdn_alarm_rule`、`pcdn_alarm_record`（告警，阶段2）
- 预留字段已含在 `pcdn_node`（`billing_mode`/`monthly_price`/`platform`/`platform_node_id`），后续不返工

**数据量与运维**：500节点×2网卡×1440分钟/天×30天 ≈ 4320 万行/月 → `node_traffic_point` 按月分区 + 3 个月明细归档（更久只存冻结 95 值）。

## 8. 技术决策清单

| 决策项 | 选择 | 理由 |
|---|---|---|
| 通信 | 统一 push | NAT 友好，兼容公网/家庭宽带 |
| Agent 语言 | Go 单二进制 | 同语言、跨平台、低占用 |
| 存储 | MySQL（分区表） | 复用 GVA 技术栈，存窗口峰值保精度 |
| 采集频率 | 每分钟 1 点（取分钟峰值） | 保 95 值精度，数据量可控 |
| 上报可靠性 | 本地持久化 + 独立重试任务 | 零丢失，采集与重试解耦 |
| 幂等 | UNIQUE(node_id,window_start,iface) | 重复上报只入库一次 |
| 时序库 | 不引入（MySQL 分区） | 复用现有栈；数据访问层抽象便于未来扩展 |
| 前端 | site/(Vue3) + web/(GVA) + h5/(推迟) | 三端分离，复用同一后端 |
| 落地形态 | 独立插件 server/plugin/pcdn/ | 不侵入 GVA 核心，中间件链对齐 |
| 凭证 | 一机一 token，库存 hash | 可管控可吊销，安全 |

## 9. 上线策略

1. 阶段1 开发完成后，先在 **10–30 台**节点验证完整闭环（注册→上机→采集→上报→95值→断网重试）
2. 验证数据准确性与稳定性后，放量到 500+，并补阶段6 OTA/批量部署
3. 阶段1 不含告警/账单/对账/大盘，但 schema 全部预留

## 10. 工程约束（遵循 AGENT.MD）

- 后端分层 `Router -> API -> Service -> Model`，`enter.go` 作组合入口
- Service 不依赖 `gin.Context`，透传 `c.Request.Context()`
- 统一响应 `{code,data,msg}`，分页 `{page,pageSize,total,list}`，列表走 `request.PageInfo` + `info.LimitOffset()`
- Swagger 注释与真实行为一致，`@Success` 落具体类型
- 插件中间件链与主系统 PrivateGroup 对齐：`JWTAuth -> MustChangePwdGuard -> CasbinHandler -> DataScope`（portal/agent 组按需调整：portal 注册/登录公开、agent 走 token）
- 前端样式优先 UnoCSS 原子类；图标空心优先
- 数据权限走 GORM 全局回调，Service 透传 ctx
- 测试复用 `server/internal/testutil`

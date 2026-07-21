# PCDN 管理系统 - 阶段1（数据底座 + 自助上机）

## 背景
用户为新开 PCDN 公司，定位「PCDN 节点供应商 / 带宽聚合商」：全国部署服务器接入个人带宽（包月/95计费向个人付费），节点接入抖音/腾讯做 PCDN（与大厂结算收款），赚中间差价。规模 >500 台、重自动化，计费数据自采集。

## 阶段1 已交付（2026-07-21）
- **后端插件** `server/plugin/pcdn/`：三组路由（admin 完整鉴权 / agent node-token / portal 个人JWT）、4 张表（gva_pcdn_node / _iface / _traffic_point / _95）、95 值计算、定时任务（离线判定 60s / 日95滚动 1h / 月冻结 月初）
- **采集 agent** `server/cmd/pcdn-agent/`：每分钟网卡峰值采集、本地 JSONL 持久化、上报+独立重试 goroutine、30s 心跳、首次激活、一键安装命令
- **运营后台** `web/src/plugin/pcdn/`：节点列表/编辑/流量 ECharts 曲线/95 值
- **官网** `site/`：独立 Vue3+Vite 应用，首页/注册/登录/控制台（自助上机生成凭证+安装命令、我的节点、流量）

## 关键架构决策
- 存储 MySQL（按月分区，运维项），不引入时序库
- 采集每分钟 1 点取该分钟峰值（保 95 值精度）
- 上报：本地持久化 → 上报失败即结束 → 独立重试任务兜底；幂等靠 (node_id,window_start,iface) 唯一索引 + OnConflict DoNothing
- 通信统一 push（NAT 友好），agent 用 node token（sha256 哈希存库，明文仅返回一次）
- 个人门户复用 `sys_users` + 贡献者角色(authorityId=888)，Service 层按 `owner_user_id` 过滤
- 95 值：升序去最高 5% 取最大；日 rolling / 月 frozen

## 设计文档
`docs/superpowers/specs/2026-07-21-pcdn-system-design.md`（完整蓝图 + 分阶段路线图 + 阶段1 详设）

## 验证状态
- `go build ./...` 全量通过；`go build ./cmd/pcdn-agent/` 通过
- 单测通过：TestPercentile95Basic / TestCalcPeriod95Percentile / TestTrafficUpsertIdempotent
- 端到端实际启动验证（mysql/redis/前后端/agent 联调）待用户环境执行

## 后续阶段（各自再走 spec → 实现）
阶段2 监控告警 → 阶段3 采购账单+付款流程 → 阶段4 销售对账（大厂结算单导入）→ 阶段5 利润大盘 → 阶段6 agent OTA+批量部署

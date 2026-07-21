package initialize

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"gorm.io/datatypes"
)

// SeedDemo 插入演示数据，便于前端联调展示。幂等：已有 DEMO- 开头节点则整体跳过。
// 生产环境如不需要，移除 plugin.go 中的 SeedDemo 调用即可。
func SeedDemo(ctx context.Context) {
	if global.GVA_DB == nil {
		return
	}
	var cnt int64
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{}).Where("node_sn LIKE ?", "DEMO-%").Count(&cnt)
	if cnt > 0 {
		return
	}

	now := time.Now()

	// 1. 节点（覆盖不同状态/平台/计费）
	nodes := []model.PcdnNode{
		{NodeSn: "DEMO-NODE-01", TokenHash: "demo_hash_01", OwnerUserID: 1, OwnerName: "admin", Region: "广东/深圳", Isp: "电信", Platform: "douyin", Status: model.NodeStatusOnline, Hostname: "sz-pcdn-01", OS: "linux", BillingMode: model.BillingModeP95, P95UnitPrice: 0.5, AgentVersion: "1.0.0", LastHeartbeatAt: &now, DeptID: 1, CreatedBy: 1},
		{NodeSn: "DEMO-NODE-02", TokenHash: "demo_hash_02", OwnerUserID: 1, OwnerName: "admin", Region: "江苏/南京", Isp: "联通", Platform: "tencent", Status: model.NodeStatusOffline, Hostname: "nj-pcdn-02", OS: "linux", BillingMode: model.BillingModeMonthly, MonthlyPrice: 100, DeptID: 1, CreatedBy: 1},
		{NodeSn: "DEMO-NODE-03", TokenHash: "demo_hash_03", OwnerUserID: 1, OwnerName: "admin", Region: "四川/成都", Isp: "移动", Platform: "douyin", Status: model.NodeStatusOnline, Hostname: "cd-pcdn-03", OS: "linux", BillingMode: model.BillingModeP95, P95UnitPrice: 0.4, AgentVersion: "1.0.0", LastHeartbeatAt: &now, DeptID: 1, CreatedBy: 1},
		{NodeSn: "DEMO-NODE-04", TokenHash: "demo_hash_04", OwnerUserID: 1, OwnerName: "admin", Region: "北京/北京", Isp: "电信", Platform: "tencent", Status: model.NodeStatusAbnormal, Hostname: "bj-pcdn-04", OS: "linux", BillingMode: model.BillingModeP95, P95UnitPrice: 0.5, DeptID: 1, CreatedBy: 1},
	}
	for i := range nodes {
		global.GVA_DB.WithContext(ctx).Create(&nodes[i])
	}

	// 2. 网卡（每节点 eth0）
	for _, n := range nodes {
		global.GVA_DB.WithContext(ctx).Create(&model.PcdnNodeIface{
			NodeID: n.ID, IfaceName: "eth0", Mac: "02:00:00:00:00:" + n.NodeSn[len(n.NodeSn)-2:], Enabled: true,
		})
	}

	// 3. 流量点（在线节点，最近 24 小时每小时 1 个点，确定性余弦曲线）
	for _, n := range nodes {
		if n.Status != model.NodeStatusOnline {
			continue
		}
		base := 50e6 // 50Mbps
		for h := 23; h >= 0; h-- {
			ws := now.Truncate(time.Hour).Add(-time.Duration(h) * time.Hour)
			factor := 0.4 + 0.6*(math.Sin(float64(h)/3.0)+1)/2
			rx := int64(base * factor)
			tx := int64(base * factor * 0.6)
			global.GVA_DB.WithContext(ctx).Create(&model.PcdnNodeTrafficPoint{
				NodeID: n.ID, IfaceName: "eth0", WindowStart: ws, RxMaxBps: rx, TxMaxBps: tx,
			})
		}
	}

	// 4. 95 值（在线节点的日/月）
	for _, n := range nodes {
		if n.Status != model.NodeStatusOnline {
			continue
		}
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		global.GVA_DB.WithContext(ctx).Create(&model.PcdnNode95{
			NodeID: n.ID, PeriodType: model.PeriodTypeDay, PeriodStart: dayStart, PeriodEnd: dayStart.Add(24 * time.Hour),
			Rx95Bps: int64(75e6), Tx95Bps: int64(45e6), Combined95Bps: int64(120e6), SampleCount: 1440, Status: model.Node95StatusRolling,
		})
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		global.GVA_DB.WithContext(ctx).Create(&model.PcdnNode95{
			NodeID: n.ID, PeriodType: model.PeriodTypeMonth, PeriodStart: monthStart, PeriodEnd: monthStart.AddDate(0, 1, 0),
			Rx95Bps: int64(80e6), Tx95Bps: int64(48e6), Combined95Bps: int64(128e6), SampleCount: 43200, Status: model.Node95StatusRolling,
		})
	}

	// 5. 告警规则
	rules := []model.PcdnAlarmRule{
		{Name: "节点离线告警", ScopeType: model.AlarmScopeAll, Metric: model.AlarmMetricOffline, Enabled: true, DeptID: 1, CreatedBy: 1},
		{Name: "带宽低于30Mbps", ScopeType: model.AlarmScopeAll, Metric: model.AlarmMetricBandwidthLow, Threshold: int64(30e6), Enabled: true, DeptID: 1, CreatedBy: 1},
		{Name: "Agent上报中断5分钟", ScopeType: model.AlarmScopeAll, Metric: model.AlarmMetricAgentDown, DurationSec: 300, Enabled: true, DeptID: 1, CreatedBy: 1},
	}
	for i := range rules {
		global.GVA_DB.WithContext(ctx).Create(&rules[i])
	}

	// 6. 告警记录（1 条 firing + 1 条 resolved）
	firedAt := now.Add(-30 * time.Minute)
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnAlarmRecord{
		RuleID: rules[0].ID, RuleName: rules[0].Name, NodeID: nodes[1].ID, NodeSn: nodes[1].NodeSn,
		Metric: model.AlarmMetricOffline, Status: model.AlarmStatusFiring, FiredAt: firedAt, NotifyCount: 1,
	})
	resolvedAt := now.Add(-2 * time.Hour)
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnAlarmRecord{
		RuleID: rules[0].ID, RuleName: rules[0].Name, NodeID: nodes[3].ID, NodeSn: nodes[3].NodeSn,
		Metric: model.AlarmMetricOffline, Status: model.AlarmStatusResolved, FiredAt: now.Add(-5 * time.Hour), ResolvedAt: &resolvedAt, NotifyCount: 2,
	})

	// 7. 账单（上月，已付款）
	lastMonth := now.AddDate(0, -1, 0)
	period := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location()).Format("2006-01")
	details := []model.BillDetail{
		{NodeID: nodes[0].ID, NodeSn: nodes[0].NodeSn, BillingMode: model.BillingModeP95, Value: 120, UnitPrice: 0.5, Amount: 60},
		{NodeID: nodes[1].ID, NodeSn: nodes[1].NodeSn, BillingMode: model.BillingModeMonthly, Value: 100, UnitPrice: 1, Amount: 100},
		{NodeID: nodes[2].ID, NodeSn: nodes[2].NodeSn, BillingMode: model.BillingModeP95, Value: 80, UnitPrice: 0.4, Amount: 32},
		{NodeID: nodes[3].ID, NodeSn: nodes[3].NodeSn, BillingMode: model.BillingModeP95, Value: 44, UnitPrice: 0.5, Amount: 22},
	}
	detailJSON, _ := json.Marshal(details)
	auditedAt := now.AddDate(0, 0, -2)
	paidAt := now.AddDate(0, 0, -1)
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnBill{
		Period: period, OwnerUserID: 1, OwnerName: "admin", OwnerContact: "admin@example.com",
		NodeCount: 4, Details: datatypes.JSON(detailJSON), TotalAmount: 214, Status: model.BillStatusPaid,
		AuditedBy: 1, AuditedAt: &auditedAt, PaidAmount: 214, PayMethod: "银行转账", PayNo: "DEMO-PAY-001", PaidAt: &paidAt, PaidBy: 1,
		DeptID: 1, CreatedBy: 1,
	})

	// 8. 结算单（上月，1 一致 + 1 差异）
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnSettlement{
		Period: period, Platform: "douyin", NodeID: nodes[0].ID, NodeSn: nodes[0].NodeSn,
		Revenue: 200, TrafficBps: int64(120e6), OurTrafficBps: int64(120e6), DiffPercent: 0, Status: model.SettlementStatusMatched, DeptID: 1, CreatedBy: 1,
	})
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnSettlement{
		Period: period, Platform: "tencent", NodeID: nodes[3].ID, NodeSn: nodes[3].NodeSn,
		Revenue: 150, TrafficBps: int64(80e6), OurTrafficBps: int64(70e6), DiffPercent: -0.125, Status: model.SettlementStatusDiff, DeptID: 1, CreatedBy: 1,
	})

	// 9. 版本
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnAgentRelease{Version: "1.0.0", DownloadURL: "https://demo.example.com/pcdn-agent-1.0.0", Stable: true, Remark: "首个稳定版", DeptID: 1, CreatedBy: 1})
	global.GVA_DB.WithContext(ctx).Create(&model.PcdnAgentRelease{Version: "1.1.0", DownloadURL: "https://demo.example.com/pcdn-agent-1.1.0", Stable: false, Force: true, Remark: "新增 OTA 自升级", DeptID: 1, CreatedBy: 1})
}

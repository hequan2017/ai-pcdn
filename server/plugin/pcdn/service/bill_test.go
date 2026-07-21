package service

import (
	"context"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/testutil"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// TestGenerateMonthlyBill 验证账单生成：包月 + p95 混合计算 + 幂等
func TestGenerateMonthlyBill(t *testing.T) {
	testutil.NewMemoryDB(t, &model.PcdnNode{}, &model.PcdnNode95{}, &model.PcdnBill{})
	ctx := context.Background()
	period := "2026-07"
	start, _, _ := parsePeriod(period)

	// 包月节点 100 元
	if err := global.GVA_DB.Create(&model.PcdnNode{
		NodeSn: "N1", OwnerUserID: 100, OwnerName: "贡献者1",
		BillingMode: model.BillingModeMonthly, MonthlyPrice: 100,
	}).Error; err != nil {
		t.Fatal(err)
	}
	// p95 节点：月95 combined=100Mbps，单价 0.5 元/Mbps → 50 元
	p2 := model.PcdnNode{NodeSn: "N2", OwnerUserID: 100, OwnerName: "贡献者1", BillingMode: model.BillingModeP95, P95UnitPrice: 0.5}
	if err := global.GVA_DB.Create(&p2).Error; err != nil {
		t.Fatal(err)
	}
	if err := global.GVA_DB.Create(&model.PcdnNode95{
		NodeID: p2.ID, PeriodType: model.PeriodTypeMonth, PeriodStart: start,
		Combined95Bps: 100 * 1e6, Status: model.Node95StatusFrozen,
	}).Error; err != nil {
		t.Fatal(err)
	}

	svc := BillService{}
	n, err := svc.GenerateMonthlyBill(ctx, period)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	if n != 1 {
		t.Fatalf("应生成 1 条账单, 得到 %d", n)
	}

	var bill model.PcdnBill
	if err := global.GVA_DB.First(&bill).Error; err != nil {
		t.Fatal(err)
	}
	// 包月 100 + p95 50 = 150
	if bill.TotalAmount != 150 {
		t.Fatalf("总额期望 150, 得到 %f", bill.TotalAmount)
	}
	if bill.NodeCount != 2 {
		t.Fatalf("节点数期望 2, 得到 %d", bill.NodeCount)
	}

	// 幂等：再次生成不应新增
	n2, _ := svc.GenerateMonthlyBill(ctx, period)
	if n2 != 0 {
		t.Fatalf("幂等应 0, 得到 %d", n2)
	}
}

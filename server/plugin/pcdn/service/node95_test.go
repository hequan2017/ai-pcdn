package service

import (
	"context"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/testutil"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// TestPercentile95Basic 验证 95 分位算法：100 个点 1..100，去掉最高 5%，剩最大 = 95
func TestPercentile95Basic(t *testing.T) {
	if percentile95(nil) != 0 {
		t.Fatal("空集应返回 0")
	}
	if percentile95([]int64{42}) != 42 {
		t.Fatal("单元素应返回自身")
	}
	vals := make([]int64, 0, 100)
	for i := int64(1); i <= 100; i++ {
		vals = append(vals, i)
	}
	if got := percentile95(vals); got != 95 {
		t.Fatalf("100点1..100的95分位期望95, 得到 %d", got)
	}
}

// TestCalcPeriod95Percentile 端到端验证：插入流量点 → 计算周期95值 → 断言
func TestCalcPeriod95Percentile(t *testing.T) {
	testutil.NewMemoryDB(t, &model.PcdnNodeTrafficPoint{}, &model.PcdnNode95{})
	ctx := context.Background()
	nodeID := uint(1)
	start := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	pts := make([]model.PcdnNodeTrafficPoint, 0, 100)
	for i := 1; i <= 100; i++ {
		pts = append(pts, model.PcdnNodeTrafficPoint{
			NodeID:      nodeID,
			IfaceName:   "eth0",
			WindowStart: start.Add(time.Duration(i-1) * time.Minute),
			RxMaxBps:    int64(i),
		})
	}
	if err := global.GVA_DB.Create(&pts).Error; err != nil {
		t.Fatalf("插入流量点失败: %v", err)
	}

	svc := Node95Service{}
	if err := svc.CalcAndSavePeriod95(ctx, nodeID, model.PeriodTypeDay, start, start.Add(100*time.Minute), model.Node95StatusRolling); err != nil {
		t.Fatalf("计算95值失败: %v", err)
	}

	var rec model.PcdnNode95
	if err := global.GVA_DB.First(&rec).Error; err != nil {
		t.Fatalf("查询95值失败: %v", err)
	}
	if rec.Rx95Bps != 95 {
		t.Fatalf("Rx95Bps 期望 95, 得到 %d", rec.Rx95Bps)
	}
	if rec.SampleCount != 100 {
		t.Fatalf("SampleCount 期望 100, 得到 %d", rec.SampleCount)
	}
}

package service

import (
	"context"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/testutil"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// TestAlarmConvergence 验证告警收敛与恢复：重复触发不新增，恢复后转 resolved
func TestAlarmConvergence(t *testing.T) {
	testutil.NewMemoryDB(t, &model.PcdnNode{}, &model.PcdnAlarmRule{}, &model.PcdnAlarmRecord{})
	ctx := context.Background()

	// 离线节点 + offline 规则（无 webhook，SendNotify 静默返回）
	node := model.PcdnNode{NodeSn: "PCDN-test", Status: model.NodeStatusOffline}
	if err := global.GVA_DB.Create(&node).Error; err != nil {
		t.Fatal(err)
	}
	rule := model.PcdnAlarmRule{Name: "离线告警", ScopeType: model.AlarmScopeAll, Metric: model.AlarmMetricOffline, Enabled: true}
	if err := global.GVA_DB.Create(&rule).Error; err != nil {
		t.Fatal(err)
	}

	svc := AlarmEngineService{}

	// 第一次检查：应产生 1 条 firing
	svc.CheckAll(ctx)
	var cnt int64
	global.GVA_DB.Model(&model.PcdnAlarmRecord{}).Where("rule_id = ? AND status = ?", rule.ID, model.AlarmStatusFiring).Count(&cnt)
	if cnt != 1 {
		t.Fatalf("首次应 1 条 firing, 得到 %d", cnt)
	}

	// 第二次检查（仍离线）：收敛，不新增
	svc.CheckAll(ctx)
	global.GVA_DB.Model(&model.PcdnAlarmRecord{}).Where("rule_id = ?", rule.ID).Count(&cnt)
	if cnt != 1 {
		t.Fatalf("收敛应仍 1 条, 得到 %d", cnt)
	}

	// 节点恢复 online：应转 resolved
	if err := global.GVA_DB.Model(&node).Update("status", model.NodeStatusOnline).Error; err != nil {
		t.Fatal(err)
	}
	svc.CheckAll(ctx)
	var rec model.PcdnAlarmRecord
	global.GVA_DB.Where("rule_id = ?", rule.ID).First(&rec)
	if rec.Status != model.AlarmStatusResolved {
		t.Fatalf("应恢复为 resolved, status=%s", rec.Status)
	}
	if rec.ResolvedAt == nil {
		t.Fatal("恢复时间应非空")
	}
}

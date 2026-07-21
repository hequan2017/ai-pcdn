package service

import (
	"context"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/testutil"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
)

// TestTrafficUpsertIdempotent 验证重复上报幂等：同一 (node,window,iface) 只入库一条
func TestTrafficUpsertIdempotent(t *testing.T) {
	testutil.NewMemoryDB(t, &model.PcdnNodeTrafficPoint{})
	ctx := context.Background()
	svc := TrafficService{}
	nodeID := uint(1)
	ws := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	pts := []request.TrafficPoint{{IfaceName: "eth0", WindowStart: ws, RxMaxBps: 100, TxMaxBps: 50}}

	if err := svc.BatchUpsertTrafficPoints(ctx, nodeID, pts); err != nil {
		t.Fatalf("首次写入失败: %v", err)
	}
	// 模拟 agent 断网重传，重复上报同一窗口
	if err := svc.BatchUpsertTrafficPoints(ctx, nodeID, pts); err != nil {
		t.Fatalf("重复写入失败: %v", err)
	}

	var cnt int64
	global.GVA_DB.Model(&model.PcdnNodeTrafficPoint{}).Count(&cnt)
	if cnt != 1 {
		t.Fatalf("幂等失败: 期望 1 条, 得到 %d", cnt)
	}
}

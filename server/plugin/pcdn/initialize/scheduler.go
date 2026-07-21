package initialize

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

// 离线判定心跳超时阈值
const offlineHeartbeatTimeout = 3 * time.Minute

// StartScheduler 启动 PCDN 后台调度任务：离线判定、日95滚动、月95冻结、告警检查。
// 由 plugin.Register 通过 go 调用；插件在 DB 就绪后注册，故 global.GVA_DB 可用。
// TODO: 后续可迁移到 GVA 标准 timed task 框架统一管理。
func StartScheduler() {
	go runOfflineCheck()
	go runDaily95()
	go runMonthlyFreeze()
	go runAlarmCheck()
}

// runOfflineCheck 每 60s 将心跳超时的 online 节点置为 offline
func runOfflineCheck() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.NodeService
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()
	for range t.C {
		if _, err := svc.MarkOfflineNodes(ctx, offlineHeartbeatTimeout); err != nil {
			logger.Bg().Mod("pcdn").Err(err).Error("离线判定失败")
		}
	}
}

// runDaily95 每小时滚动计算当日各节点 95 值（rolling）
func runDaily95() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.Node95Service
	_ = svc.CalcAllNodesDaily95(ctx, time.Now())
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for range t.C {
		_ = svc.CalcAllNodesDaily95(ctx, time.Now())
	}
}

// runMonthlyFreeze 每日检查，月初冻结上月 95 值（frozen）
func runMonthlyFreeze() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.Node95Service
	t := time.NewTicker(24 * time.Hour)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		if now.Day() == 1 {
			lastMonth := now.AddDate(0, 0, -1)
			if err := svc.FreezeAllNodesMonthly95(ctx, lastMonth); err != nil {
				logger.Bg().Mod("pcdn").Err(err).Error("月95冻结失败")
			}
		}
	}
}

// runAlarmCheck 每 60s 检查所有启用的告警规则，触发/恢复并通知（带收敛）
func runAlarmCheck() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.AlarmEngineService
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()
	for range t.C {
		if err := svc.CheckAll(ctx); err != nil {
			logger.Bg().Mod("pcdn").Err(err).Error("告警检查失败")
		}
	}
}

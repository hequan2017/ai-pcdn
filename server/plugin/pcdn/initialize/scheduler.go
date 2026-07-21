package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

// 离线判定心跳超时阈值
const offlineHeartbeatTimeout = 3 * time.Minute

// StartScheduler 启动 PCDN 后台调度任务：离线判定、日95滚动、月95冻结、告警检查、月账单生成。
func StartScheduler() {
	go runOfflineCheck()
	go runDaily95()
	go runMonthlyFreeze()
	go runAlarmCheck()
	go runMonthlyBill()
}

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

// runMonthlyBill 每月2号自动生成上月采购账单（上月95已冻结）
func runMonthlyBill() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.BillService
	t := time.NewTicker(24 * time.Hour)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		if now.Day() == 2 {
			last := now.AddDate(0, 0, -1)
			period := fmt.Sprintf("%04d-%02d", last.Year(), int(last.Month()))
			if _, err := svc.GenerateMonthlyBill(ctx, period); err != nil {
				logger.Bg().Mod("pcdn").Err(err).Error("月账单生成失败")
			}
		}
	}
}

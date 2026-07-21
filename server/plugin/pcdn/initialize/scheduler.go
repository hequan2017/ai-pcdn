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

// safeRun 执行一次任务并恢复 panic，避免单次异常终结整个调度循环。
func safeRun(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			logger.Bg().Mod("pcdn").Err(fmt.Errorf("%v", r)).Error("调度任务异常已恢复: " + name)
		}
	}()
	fn()
}

// runOfflineCheck 每 60s 将心跳超时的 online 节点置为 offline
func runOfflineCheck() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.NodeService
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()
	for range t.C {
		safeRun("offline", func() {
			if _, err := svc.MarkOfflineNodes(ctx, offlineHeartbeatTimeout); err != nil {
				logger.Bg().Mod("pcdn").Err(err).Error("离线判定失败")
			}
		})
	}
}

// runDaily95 每小时滚动计算当日各节点 95 值（rolling）
func runDaily95() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.Node95Service
	safeRun("daily95-init", func() { _ = svc.CalcAllNodesDaily95(ctx, time.Now()) })
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for range t.C {
		safeRun("daily95", func() { _ = svc.CalcAllNodesDaily95(ctx, time.Now()) })
	}
}

// nextMonthStart 返回下一个自然月 1 号 0 点（用于固定时刻触发月任务，避免 24h ticker 随重启漂移）
func nextMonthStart(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0)
}

// runMonthlyFreeze 每月 1 号 0 点冻结上月 95 值（固定时刻，重启不错过）
func runMonthlyFreeze() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.Node95Service
	for {
		now := time.Now()
		timer := time.NewTimer(nextMonthStart(now).Sub(now))
		<-timer.C
		safeRun("freeze", func() {
			lastMonth := time.Now().AddDate(0, -1, 0) // 上月
			if err := svc.FreezeAllNodesMonthly95(ctx, lastMonth); err != nil {
				logger.Bg().Mod("pcdn").Err(err).Error("月95冻结失败")
			}
		})
	}
}

// runAlarmCheck 每 60s 检查所有启用的告警规则
func runAlarmCheck() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.AlarmEngineService
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()
	for range t.C {
		safeRun("alarm", func() {
			if err := svc.CheckAll(ctx); err != nil {
				logger.Bg().Mod("pcdn").Err(err).Error("告警检查失败")
			}
		})
	}
}

// runMonthlyBill 每月 1 号 0:05 生成上月账单（冻结完成后）
func runMonthlyBill() {
	ctx := context.Background()
	svc := service.ServiceGroupApp.BillService
	for {
		now := time.Now()
		timer := time.NewTimer(nextMonthStart(now).Add(5 * time.Minute).Sub(now))
		<-timer.C
		safeRun("bill", func() {
			last := time.Now().AddDate(0, -1, 0) // 上月（修复：原 AddDate(0,0,-1) 会得到当月）
			period := fmt.Sprintf("%04d-%02d", last.Year(), int(last.Month()))
			if _, err := svc.GenerateMonthlyBill(ctx, period); err != nil {
				logger.Bg().Mod("pcdn").Err(err).Error("月账单生成失败")
			}
		})
	}
}

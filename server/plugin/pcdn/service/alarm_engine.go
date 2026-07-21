package service

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// AlarmEngineService 告警引擎：周期检查规则 → 触发/恢复 → 通知（带收敛）
type AlarmEngineService struct{}

// CheckAll 检查所有启用的规则（由调度任务周期调用）
func (s *AlarmEngineService) CheckAll(ctx context.Context) error {
	var rules []model.PcdnAlarmRule
	if err := global.GVA_DB.WithContext(ctx).Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return err
	}
	for i := range rules {
		s.checkRule(ctx, &rules[i])
	}
	return nil
}

// checkRule 对单条规则：取范围内节点 → 逐节点评估 → 更新告警记录
func (s *AlarmEngineService) checkRule(ctx context.Context, rule *model.PcdnAlarmRule) {
	nodes := s.getScopeNodes(ctx, rule)
	for i := range nodes {
		triggered, value := s.evaluate(ctx, rule, &nodes[i])
		s.updateRecord(ctx, rule, &nodes[i], triggered, value)
	}
}

// getScopeNodes 按 scope 取节点（排除已停用）
func (s *AlarmEngineService) getScopeNodes(ctx context.Context, rule *model.PcdnAlarmRule) []model.PcdnNode {
	var nodes []model.PcdnNode
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{}).Where("status != ?", model.NodeStatusDisabled)
	switch rule.ScopeType {
	case model.AlarmScopeGroup:
		db = db.Where("group_id = ?", rule.ScopeValue)
	case model.AlarmScopeNode:
		db = db.Where("id = ?", rule.ScopeValue)
	}
	db.Find(&nodes)
	return nodes
}

// evaluate 评估单节点是否满足告警条件，返回是否触发及当前值
func (s *AlarmEngineService) evaluate(ctx context.Context, rule *model.PcdnAlarmRule, node *model.PcdnNode) (bool, int64) {
	switch rule.Metric {
	case model.AlarmMetricOffline:
		return node.Status == model.NodeStatusOffline, 0

	case model.AlarmMetricAgentDown:
		var latest model.PcdnNodeTrafficPoint
		err := global.GVA_DB.WithContext(ctx).Where("node_id = ?", node.ID).Order("window_start desc").First(&latest).Error
		if err != nil {
			return true, 0 // 无上报数据视为中断
		}
		dur := time.Duration(rule.DurationSec) * time.Second
		if dur <= 0 {
			dur = 5 * time.Minute
		}
		gap := time.Since(latest.WindowStart)
		return gap > dur, int64(gap.Seconds())

	case model.AlarmMetricBandwidthLow:
		var latest model.PcdnNodeTrafficPoint
		_ = global.GVA_DB.WithContext(ctx).Where("node_id = ?", node.ID).Order("window_start desc").First(&latest).Error
		cur := latest.RxMaxBps + latest.TxMaxBps
		return cur < rule.Threshold, cur

	case model.AlarmMetricP95High:
		var n95 model.PcdnNode95
		_ = global.GVA_DB.WithContext(ctx).Where("node_id = ? AND period_type = ?", node.ID, model.PeriodTypeDay).Order("period_start desc").First(&n95).Error
		return n95.Combined95Bps > rule.Threshold, n95.Combined95Bps
	}
	return false, 0
}

// updateRecord 收敛逻辑：触发时若已有 firing 记录则跳过；恢复时关闭 firing 并通知
func (s *AlarmEngineService) updateRecord(ctx context.Context, rule *model.PcdnAlarmRule, node *model.PcdnNode, triggered bool, value int64) {
	var firing model.PcdnAlarmRecord
	err := global.GVA_DB.WithContext(ctx).
		Where("rule_id = ? AND node_id = ? AND status = ?", rule.ID, node.ID, model.AlarmStatusFiring).
		First(&firing).Error

	if triggered {
		if err == nil {
			return // 已在 firing，收敛不重复通知
		}
		now := time.Now()
		rec := model.PcdnAlarmRecord{
			RuleID:       rule.ID,
			RuleName:     rule.Name,
			NodeID:       node.ID,
			NodeSn:       node.NodeSn,
			Metric:       rule.Metric,
			TriggerValue: value,
			Status:       model.AlarmStatusFiring,
			FiredAt:      now,
			NotifyCount:  1,
		}
		if err := global.GVA_DB.WithContext(ctx).Create(&rec).Error; err == nil {
			_ = SendNotify(rule, &rec, true)
		}
		return
	}

	if err != nil {
		return // 无 firing 记录，无需恢复
	}
	now := time.Now()
	if err := global.GVA_DB.WithContext(ctx).Model(&firing).
		Updates(map[string]interface{}{"status": model.AlarmStatusResolved, "resolved_at": now}).Error; err == nil {
		firing.Status = model.AlarmStatusResolved
		_ = SendNotify(rule, &firing, false)
	}
}

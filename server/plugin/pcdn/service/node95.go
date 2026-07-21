package service

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"gorm.io/gorm"
)

// Node95Service 95值计算服务
type Node95Service struct{}

// CalcAndSavePeriod95 计算指定节点在 [start,end) 周期的 95 值并落库（先删同周期旧记录再插入）。
// 95值语义：周期内所有分钟峰值点升序，去掉最高 5% 后的最大值。
func (s *Node95Service) CalcAndSavePeriod95(ctx context.Context, nodeID uint, periodType string, start, end time.Time, status string) error {
	var points []model.PcdnNodeTrafficPoint
	err := global.GVA_DB.WithContext(ctx).
		Where("node_id = ? AND window_start >= ? AND window_start < ?", nodeID, start, end).
		Find(&points).Error
	if err != nil {
		return err
	}
	rec := model.PcdnNode95{
		NodeID:        nodeID,
		PeriodType:    periodType,
		PeriodStart:   start,
		PeriodEnd:     end,
		Rx95Bps:       percentile95(pick(points, func(p model.PcdnNodeTrafficPoint) int64 { return p.RxMaxBps })),
		Tx95Bps:       percentile95(pick(points, func(p model.PcdnNodeTrafficPoint) int64 { return p.TxMaxBps })),
		Combined95Bps: percentile95(pick(points, func(p model.PcdnNodeTrafficPoint) int64 { return p.RxMaxBps + p.TxMaxBps })),
		SampleCount:   len(points),
		Status:        status,
	}
	if status == model.Node95StatusFrozen {
		now := time.Now()
		rec.FrozenAt = &now
	}
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("node_id = ? AND period_type = ? AND period_start = ?", nodeID, periodType, start).
			Delete(&model.PcdnNode95{}).Error; err != nil {
			return err
		}
		return tx.Create(&rec).Error
	})
}

// CalcAllNodesDaily95 计算所有非禁用节点的某日 95 值（rolling）
func (s *Node95Service) CalcAllNodesDaily95(ctx context.Context, day time.Time) error {
	var nodeIDs []uint
	if err := global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{}).
		Where("status != ?", model.NodeStatusDisabled).Pluck("id", &nodeIDs).Error; err != nil {
		return err
	}
	loc := day.Location()
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)
	for _, id := range nodeIDs {
		if err := s.CalcAndSavePeriod95(ctx, id, model.PeriodTypeDay, start, end, model.Node95StatusRolling); err != nil {
			continue // 单节点失败不阻塞其余
		}
	}
	return nil
}

// FreezeAllNodesMonthly95 冻结所有非禁用节点的某月 95 值（frozen，作为账单依据）
func (s *Node95Service) FreezeAllNodesMonthly95(ctx context.Context, month time.Time) error {
	var nodeIDs []uint
	if err := global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{}).
		Where("status != ?", model.NodeStatusDisabled).Pluck("id", &nodeIDs).Error; err != nil {
		return err
	}
	loc := month.Location()
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	for _, id := range nodeIDs {
		if err := s.CalcAndSavePeriod95(ctx, id, model.PeriodTypeMonth, start, end, model.Node95StatusFrozen); err != nil {
			continue
		}
	}
	return nil
}

// GetNode95List 查询节点 95 值（按 period_type 过滤，最近 100 条）
func (s *Node95Service) GetNode95List(ctx context.Context, nodeID uint, periodType string) (list []model.PcdnNode95, err error) {
	db := global.GVA_DB.WithContext(ctx).Where("node_id = ?", nodeID)
	if periodType != "" {
		db = db.Where("period_type = ?", periodType)
	}
	err = db.Order("period_start desc").Limit(100).Find(&list).Error
	return
}

// pick 从流量点集合提取指定维度的值列表
func pick(points []model.PcdnNodeTrafficPoint, f func(model.PcdnNodeTrafficPoint) int64) []int64 {
	out := make([]int64, len(points))
	for i, p := range points {
		out[i] = f(p)
	}
	return out
}

// percentile95 升序后去掉最高 5% 取最大值（第 95 百分位）
func percentile95(vals []int64) int64 {
	if len(vals) == 0 {
		return 0
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	n := len(vals)
	idx := int(math.Ceil(float64(n)*0.95)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	return vals[idx]
}

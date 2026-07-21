package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"gorm.io/gorm"
)

// 差异阈值：自采集与大厂差异超过 10% 标记为 diff
const settlementDiffThreshold = 0.10

// SettlementService 大厂结算单服务
type SettlementService struct{}

// Import 导入单条结算单，自动按 node_sn 关联节点并核对自采集流量
func (s *SettlementService) Import(ctx context.Context, req request.SettlementImport) error {
	st := model.PcdnSettlement{
		Period:         req.Period,
		Platform:       req.Platform,
		PlatformNodeID: req.PlatformNodeID,
		Revenue:        req.Revenue,
		TrafficBps:     req.TrafficBps,
		Remark:         req.Remark,
		Status:         model.SettlementStatusPending,
	}
	if req.NodeSn != "" {
		var node model.PcdnNode
		if err := global.GVA_DB.WithContext(ctx).Where("node_sn = ?", req.NodeSn).First(&node).Error; err == nil {
			st.NodeID = node.ID
		}
		st.NodeSn = req.NodeSn
	}
	if st.NodeID > 0 {
		s.fillOurTraffic(ctx, &st)
	}
	return global.GVA_DB.WithContext(ctx).Create(&st).Error
}

// fillOurTraffic 填充自采集流量并计算差异与状态
func (s *SettlementService) fillOurTraffic(ctx context.Context, st *model.PcdnSettlement) {
	start, _, err := parsePeriod(st.Period)
	if err != nil {
		return
	}
	var n95 model.PcdnNode95
	_ = global.GVA_DB.WithContext(ctx).
		Where("node_id = ? AND period_type = ? AND period_start = ?", st.NodeID, model.PeriodTypeMonth, start).
		First(&n95).Error
	st.OurTrafficBps = n95.Combined95Bps
	if st.TrafficBps > 0 {
		st.DiffPercent = float64(st.OurTrafficBps-st.TrafficBps) / float64(st.TrafficBps)
		st.Status = model.SettlementStatusMatched
		if st.DiffPercent > settlementDiffThreshold || st.DiffPercent < -settlementDiffThreshold {
			st.Status = model.SettlementStatusDiff
		}
	}
}

// Recheck 重新核对某账期所有结算单（自采集更新后重算差异）
func (s *SettlementService) Recheck(ctx context.Context, period string) (int, error) {
	var list []model.PcdnSettlement
	if err := global.GVA_DB.WithContext(ctx).Where("period = ? AND node_id > 0", period).Find(&list).Error; err != nil {
		return 0, err
	}
	for i := range list {
		s.fillOurTraffic(ctx, &list[i])
		global.GVA_DB.WithContext(ctx).Model(&model.PcdnSettlement{}).Where("id = ?", list[i].ID).
			Updates(map[string]interface{}{"our_traffic_bps": list[i].OurTrafficBps, "diff_percent": list[i].DiffPercent, "status": list[i].Status})
	}
	return len(list), nil
}

// GetList 分页查询结算单
func (s *SettlementService) GetList(ctx context.Context, info request.SettlementSearch) (list []model.PcdnSettlement, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnSettlement{})
	if info.Period != "" {
		db = db.Where("period = ?", info.Period)
	}
	if info.Platform != "" {
		db = db.Where("platform = ?", info.Platform)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.NodeSn != "" {
		db = db.Where("node_sn LIKE ?", "%"+info.NodeSn+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("id desc").Find(&list).Error
	return
}

// RevenueSummary 应收汇总
type RevenueSummary struct {
	Period       string  `json:"period"`
	Platform     string  `json:"platform"`
	TotalRevenue float64 `json:"totalRevenue"`
	MatchedCount int64   `json:"matchedCount"`
	DiffCount    int64   `json:"diffCount"`
	PendingCount int64   `json:"pendingCount"`
}

func (s *SettlementService) RevenueSummary(ctx context.Context, req request.RevenueSummaryReq) (RevenueSummary, error) {
	sum := RevenueSummary{Period: req.Period, Platform: req.Platform}
	base := func() *gorm.DB {
		db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnSettlement{})
		if req.Period != "" {
			db = db.Where("period = ?", req.Period)
		}
		if req.Platform != "" {
			db = db.Where("platform = ?", req.Platform)
		}
		return db
	}
	var rev struct{ Total float64 }
	base().Where("status != ?", model.SettlementStatusPending).Select("COALESCE(SUM(revenue),0) as total").Scan(&rev)
	sum.TotalRevenue = rev.Total
	base().Where("status = ?", model.SettlementStatusMatched).Count(&sum.MatchedCount)
	base().Where("status = ?", model.SettlementStatusDiff).Count(&sum.DiffCount)
	base().Where("status = ?", model.SettlementStatusPending).Count(&sum.PendingCount)
	return sum, nil
}

// Delete 删除结算单
func (s *SettlementService) Delete(ctx context.Context, id uint) error {
	return global.GVA_DB.WithContext(ctx).Delete(&model.PcdnSettlement{}, "id = ?", id).Error
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// ProfitService 利润大盘（聚合查询，无独立表）
type ProfitService struct{}

// ProfitSummary 利润汇总
type ProfitSummary struct {
	Period          string  `json:"period"`
	Revenue         float64 `json:"revenue"`     // 大厂收入
	Cost            float64 `json:"cost"`        // 采购成本
	Profit          float64 `json:"profit"`      // 利润
	ProfitMargin    float64 `json:"profitMargin"`
	BillCount       int64   `json:"billCount"`
	SettlementCount int64   `json:"settlementCount"`
}

// GetProfitSummary 某账期利润汇总：收入(已核对结算单) - 成本(已审核/已付款账单)
func (s *ProfitService) GetProfitSummary(ctx context.Context, period string) (ProfitSummary, error) {
	var sum ProfitSummary
	sum.Period = period
	var rev struct{ Total float64 }
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnSettlement{}).
		Where("period = ? AND status != ?", period, model.SettlementStatusPending).
		Select("COALESCE(SUM(revenue),0) as total").Scan(&rev)
	sum.Revenue = rev.Total
	var cost struct{ Total float64 }
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{}).
		Where("period = ? AND status IN ?", period, []string{model.BillStatusApproved, model.BillStatusPaid}).
		Select("COALESCE(SUM(total_amount),0) as total").Scan(&cost)
	sum.Cost = cost.Total
	sum.Profit = sum.Revenue - sum.Cost
	if sum.Revenue > 0 {
		sum.ProfitMargin = sum.Profit / sum.Revenue
	}
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{}).Where("period = ?", period).Count(&sum.BillCount)
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnSettlement{}).Where("period = ?", period).Count(&sum.SettlementCount)
	return sum, nil
}

// PlatformRevenue 平台收入明细
type PlatformRevenue struct {
	Platform string  `json:"platform"`
	Revenue  float64 `json:"revenue"`
	Count    int64   `json:"count"`
}

func (s *ProfitService) GetRevenueByPlatform(ctx context.Context, period string) ([]PlatformRevenue, error) {
	var list []PlatformRevenue
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnSettlement{}).
		Select("platform, COALESCE(SUM(revenue),0) as revenue, COUNT(*) as count").
		Where("period = ? AND status != ?", period, model.SettlementStatusPending).
		Group("platform").Scan(&list)
	return list, nil
}

// OwnerCost 贡献者成本明细
type OwnerCost struct {
	OwnerName string  `json:"ownerName"`
	Cost      float64 `json:"cost"`
	Count     int64   `json:"count"`
}

func (s *ProfitService) GetCostByOwner(ctx context.Context, period string) ([]OwnerCost, error) {
	var list []OwnerCost
	global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{}).
		Select("owner_name, COALESCE(SUM(total_amount),0) as cost, COUNT(*) as count").
		Where("period = ?", period).
		Group("owner_name").Scan(&list)
	return list, nil
}

// MonthTrend 月度趋势
type MonthTrend struct {
	Period  string  `json:"period"`
	Revenue float64 `json:"revenue"`
	Cost    float64 `json:"cost"`
	Profit  float64 `json:"profit"`
}

// GetProfitTrend 最近 N 月利润趋势
func (s *ProfitService) GetProfitTrend(ctx context.Context, months int) ([]MonthTrend, error) {
	if months <= 0 {
		months = 6
	}
	now := time.Now()
	trend := make([]MonthTrend, 0, months)
	for i := months - 1; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		period := fmt.Sprintf("%04d-%02d", t.Year(), int(t.Month()))
		sum, _ := s.GetProfitSummary(ctx, period)
		trend = append(trend, MonthTrend{Period: period, Revenue: sum.Revenue, Cost: sum.Cost, Profit: sum.Profit})
	}
	return trend, nil
}

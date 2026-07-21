package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/gin-gonic/gin"
)

// ProfitApi 利润大盘接口
type ProfitApi struct{}

// GetProfitSummary 利润汇总
// @Tags PcdnProfit
// @Summary 利润汇总
// @Security ApiKeyAuth
// @Router /pcdn/admin/profit/summary [get]
func (a *ProfitApi) GetProfitSummary(c *gin.Context) {
	var q request.ProfitQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sum, err := profitService.GetProfitSummary(c.Request.Context(), q.Period)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(sum, "获取成功", c)
}

// GetRevenueByPlatform 按平台收入
// @Tags PcdnProfit
// @Summary 按平台收入
// @Security ApiKeyAuth
// @Router /pcdn/admin/profit/revenueByPlatform [get]
func (a *ProfitApi) GetRevenueByPlatform(c *gin.Context) {
	period := c.Query("period")
	list, err := profitService.GetRevenueByPlatform(c.Request.Context(), period)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetCostByOwner 按贡献者成本
// @Tags PcdnProfit
// @Summary 按贡献者成本
// @Security ApiKeyAuth
// @Router /pcdn/admin/profit/costByOwner [get]
func (a *ProfitApi) GetCostByOwner(c *gin.Context) {
	period := c.Query("period")
	list, err := profitService.GetCostByOwner(c.Request.Context(), period)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetProfitTrend 利润趋势
// @Tags PcdnProfit
// @Summary 利润趋势
// @Security ApiKeyAuth
// @Router /pcdn/admin/profit/trend [get]
func (a *ProfitApi) GetProfitTrend(c *gin.Context) {
	var q request.ProfitQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := profitService.GetProfitTrend(c.Request.Context(), q.Months)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

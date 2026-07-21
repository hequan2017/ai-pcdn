package router

import "github.com/gin-gonic/gin"

// ProfitRouter 利润大盘路由（admin 组）
type ProfitRouter struct{}

// InitProfitRouter 挂载利润路由
func (r *ProfitRouter) InitProfitRouter(group *gin.RouterGroup) {
	group.GET("profit/summary", profitApi.GetProfitSummary)
	group.GET("profit/revenueByPlatform", profitApi.GetRevenueByPlatform)
	group.GET("profit/costByOwner", profitApi.GetCostByOwner)
	group.GET("profit/trend", profitApi.GetProfitTrend)
}

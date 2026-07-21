package router

import "github.com/gin-gonic/gin"

// SettlementRouter 结算单路由（admin 组）
type SettlementRouter struct{}

// InitSettlementRouter 挂载结算单路由
func (r *SettlementRouter) InitSettlementRouter(group *gin.RouterGroup) {
	group.POST("settlement/import", settlementApi.Import)
	group.PUT("settlement/recheck", settlementApi.Recheck)
	group.GET("settlement/list", settlementApi.GetList)
	group.GET("settlement/revenue", settlementApi.RevenueSummary)
	group.DELETE("settlement/delete", settlementApi.Delete)
}

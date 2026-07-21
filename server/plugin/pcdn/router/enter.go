package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/api"

// RouterGroup PCDN 路由聚合入口
type RouterGroup struct {
	NodeRouter
	AgentRouter
	PortalRouter
	AlarmRouter
	BillRouter
	SettlementRouter
	ProfitRouter
	ReleaseRouter
}

// RouterGroupApp 全局路由组实例，供 initialize 层引用
var RouterGroupApp = new(RouterGroup)

// 业务 api 引用
var (
	nodeApi   = api.ApiGroupApp.NodeApi
	agentApi  = api.ApiGroupApp.AgentApi
	portalApi = api.ApiGroupApp.PortalApi
	alarmApi  = api.ApiGroupApp.AlarmApi
	billApi   = api.ApiGroupApp.BillApi
	settlementApi = api.ApiGroupApp.SettlementApi
	profitApi     = api.ApiGroupApp.ProfitApi
	releaseApi    = api.ApiGroupApp.ReleaseApi
)

package api

import (
	pcdnService "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/service"
)

// ApiGroup PCDN API 聚合入口
type ApiGroup struct {
	NodeApi
	AgentApi
	PortalApi
	AlarmApi
}

// ApiGroupApp 全局 API 组实例，供 router 层引用
var ApiGroupApp = new(ApiGroup)

// 业务 service 引用
var (
	nodeService       = pcdnService.ServiceGroupApp.NodeService
	trafficService    = pcdnService.ServiceGroupApp.TrafficService
	agentService      = pcdnService.ServiceGroupApp.AgentService
	portalService     = pcdnService.ServiceGroupApp.PortalService
	node95Service     = pcdnService.ServiceGroupApp.Node95Service
	alarmRuleService  = pcdnService.ServiceGroupApp.AlarmRuleService
	alarmRecordService = pcdnService.ServiceGroupApp.AlarmRecordService
)

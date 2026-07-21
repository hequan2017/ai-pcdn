package service

// ServiceGroup PCDN 业务服务聚合入口
type ServiceGroup struct {
	NodeService
	TrafficService
	AgentService
	PortalService
}

// ServiceGroupApp 全局服务组实例，供 api 层引用
var ServiceGroupApp = new(ServiceGroup)

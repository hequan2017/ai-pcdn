package router

import "github.com/gin-gonic/gin"

// AgentRouter agent 上报路由（agent 组，已挂 AgentTokenAuth）
type AgentRouter struct{}

// InitAgentRouter 挂载 agent 路由
func (r *AgentRouter) InitAgentRouter(group *gin.RouterGroup) {
	group.POST("activate", agentApi.Activate)
	group.POST("report", agentApi.Report)
	group.POST("heartbeat", agentApi.Heartbeat)
}

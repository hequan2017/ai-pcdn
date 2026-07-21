package router

import "github.com/gin-gonic/gin"

// NodeRouter 节点管理路由（admin 组，中间件链已由 initialize/router.go 挂载）
type NodeRouter struct{}

// InitNodeRouter 挂载节点管理路由
func (r *NodeRouter) InitNodeRouter(group *gin.RouterGroup) {
	group.GET("node/list", nodeApi.GetNodeList)
	group.GET("node/find", nodeApi.GetNode)
	group.POST("node/create", nodeApi.CreateNode)
	group.PUT("node/update", nodeApi.UpdateNode)
	group.DELETE("node/delete", nodeApi.DeleteNode)
	group.DELETE("node/deleteByIds", nodeApi.DeleteNodeByIds)
	group.GET("node/traffic", nodeApi.GetNodeTraffic)
	group.GET("node/n95", nodeApi.GetNode95)
}

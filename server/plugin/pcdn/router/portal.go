package router

import "github.com/gin-gonic/gin"

// PortalRouter 个人门户路由（public 公开注册/登录；private 个人JWT，查自己的节点/流量）
type PortalRouter struct{}

// InitPortalRouter 挂载 portal 路由
func (r *PortalRouter) InitPortalRouter(public *gin.RouterGroup, private *gin.RouterGroup) {
	public.POST("register", portalApi.Register)
	public.POST("login", portalApi.Login)
	private.GET("myNodes", portalApi.MyNodes)
	private.GET("myNodeTraffic", portalApi.MyNodeTraffic)
	private.GET("myBills", portalApi.MyBills)
	private.POST("addNode", portalApi.AddNode)
}

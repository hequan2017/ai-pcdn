package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	pcdnMiddleware "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/router"
	"github.com/gin-gonic/gin"
)

// Router 挂载 PCDN 各路由组。每组自建 group 并配置对应中间件链：
//   - /pcdn/admin/* ：完整后台鉴权（与主系统 PrivateGroup 对齐：JWTAuth→MustChangePwdGuard→CasbinHandler→DataScope）
//   - /pcdn/agent/* ：采集 agent 上报，node token 鉴权（AgentTokenAuth），不走 JWT/Casbin/DataScope
//   - /pcdn/portal/* ：个人门户，后续批次补充
func Router(engine *gin.Engine) {
	prefix := engine.Group(global.GVA_CONFIG.System.RouterPrefix)

	// admin 组：运营后台，完整鉴权 + 数据权限（Use 链返回 gin.IRoutes，故先取 *gin.RouterGroup 再挂中间件）
	adminGroup := prefix.Group("pcdn/admin")
	adminGroup.Use(middleware.JWTAuth()).
		Use(middleware.MustChangePwdGuard()).
		Use(middleware.CasbinHandler()).
		Use(middleware.DataScope())
	router.RouterGroupApp.NodeRouter.InitNodeRouter(adminGroup)
	router.RouterGroupApp.AlarmRouter.InitAlarmRouter(adminGroup)
	router.RouterGroupApp.BillRouter.InitBillRouter(adminGroup)
	router.RouterGroupApp.SettlementRouter.InitSettlementRouter(adminGroup)
	router.RouterGroupApp.ProfitRouter.InitProfitRouter(adminGroup)

	// agent 组：采集 agent 上报，node token 鉴权
	agentGroup := prefix.Group("pcdn/agent")
	agentGroup.Use(pcdnMiddleware.AgentTokenAuth())
	router.RouterGroupApp.AgentRouter.InitAgentRouter(agentGroup)

	// portal 组：个人门户（public 公开注册/登录；private 个人 JWT，查自己的节点/流量）
	publicPortal := prefix.Group("pcdn/portal")
	privatePortal := prefix.Group("pcdn/portal")
	privatePortal.Use(middleware.JWTAuth())
	router.RouterGroupApp.PortalRouter.InitPortalRouter(publicPortal, privatePortal)
}

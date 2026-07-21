package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// AgentApi agent 上报接口（/pcdn/agent/*，经 AgentTokenAuth 鉴权）
type AgentApi struct{}

// Activate 首次激活：回填硬件信息并转 online
// @Tags PcdnAgent
// @Summary agent 激活
// @accept application/json
// @Produce application/json
// @Param data body request.AgentActivate true "激活信息"
// @Success 200 {object} response.Response{msg=string} "激活成功"
// @Router /pcdn/agent/activate [post]
func (a *AgentApi) Activate(c *gin.Context) {
	nodeID := middleware.GetNodeIDFromCtx(c)
	var req request.AgentActivate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := agentService.Activate(c.Request.Context(), nodeID, c.ClientIP(), req); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("agent激活失败")
		response.FailWithMessage("激活失败", c)
		return
	}
	response.OkWithMessage("激活成功", c)
}

// Report 批量上报流量分钟峰值点（幂等）
// @Tags PcdnAgent
// @Summary agent 上报流量
// @accept application/json
// @Produce application/json
// @Param data body request.TrafficReport true "流量点集合"
// @Success 200 {object} response.Response{msg=string} "上报成功"
// @Router /pcdn/agent/report [post]
func (a *AgentApi) Report(c *gin.Context) {
	nodeID := middleware.GetNodeIDFromCtx(c)
	var req request.TrafficReport
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := trafficService.BatchUpsertTrafficPoints(c.Request.Context(), nodeID, req.Points); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("agent上报失败")
		response.FailWithMessage("上报失败", c)
		return
	}
	response.OkWithMessage("上报成功", c)
}

// Heartbeat 心跳
// @Tags PcdnAgent
// @Summary agent 心跳
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "ok"
// @Router /pcdn/agent/heartbeat [post]
func (a *AgentApi) Heartbeat(c *gin.Context) {
	nodeID := middleware.GetNodeIDFromCtx(c)
	var req request.AgentHeartbeat
	_ = c.ShouldBindJSON(&req) // 心跳允许空 body
	if err := agentService.Heartbeat(c.Request.Context(), nodeID, c.ClientIP(), req); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("agent心跳失败")
		response.FailWithMessage("心跳失败", c)
		return
	}
	response.OkWithMessage("ok", c)
}

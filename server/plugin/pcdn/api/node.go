package api

import (
	"strconv"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// NodeApi 节点管理接口
type NodeApi struct{}

// CreateNode 创建节点
// @Tags PcdnNode
// @Summary 创建节点
// @Security ApiKeyAuth
// @Param data body model.PcdnNode true "节点信息"
// @Success 200 {object} response.Response{data=model.PcdnNode,msg=string} "创建成功"
// @Router /pcdn/admin/node/create [post]
func (a *NodeApi) CreateNode(c *gin.Context) {
	var node model.PcdnNode
	if err := c.ShouldBindJSON(&node); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := nodeService.CreateNode(c.Request.Context(), &node); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("创建节点失败")
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithDetailed(node, "创建成功", c)
}

// DeleteNode 删除节点
// @Tags PcdnNode
// @Summary 删除节点
// @Security ApiKeyAuth
// @Param data query common.request.GetById true "节点ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /pcdn/admin/node/delete [delete]
func (a *NodeApi) DeleteNode(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := nodeService.DeleteNode(c.Request.Context(), req.Uint()); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("删除节点失败")
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteNodeByIds 批量删除节点
// @Tags PcdnNode
// @Summary 批量删除节点
// @Security ApiKeyAuth
// @Param data body common.request.IdsReq true "节点ID集合"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /pcdn/admin/node/deleteByIds [delete]
func (a *NodeApi) DeleteNodeByIds(c *gin.Context) {
	var req commonReq.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, uint(id))
	}
	if err := nodeService.DeleteNodeByIds(c.Request.Context(), ids); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("批量删除节点失败")
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateNode 更新节点
// @Tags PcdnNode
// @Summary 更新节点
// @Security ApiKeyAuth
// @Param data body model.PcdnNode true "节点信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /pcdn/admin/node/update [put]
func (a *NodeApi) UpdateNode(c *gin.Context) {
	var node model.PcdnNode
	if err := c.ShouldBindJSON(&node); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := nodeService.UpdateNode(c.Request.Context(), node); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("更新节点失败")
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetNode 查询节点详情
// @Tags PcdnNode
// @Summary 查询节点详情
// @Security ApiKeyAuth
// @Param data query common.request.GetById true "节点ID"
// @Success 200 {object} response.Response{data=model.PcdnNode,msg=string} "查询成功"
// @Router /pcdn/admin/node/find [get]
func (a *NodeApi) GetNode(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	node, err := nodeService.GetNode(c.Request.Context(), req.Uint())
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("查询节点失败")
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(node, "查询成功", c)
}

// GetNodeList 分页查询节点列表
// @Tags PcdnNode
// @Summary 分页查询节点列表
// @Security ApiKeyAuth
// @Param data query request.NodeSearch true "查询条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /pcdn/admin/node/list [get]
func (a *NodeApi) GetNodeList(c *gin.Context) {
	var pageInfo request.NodeSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := nodeService.GetNodeList(c.Request.Context(), pageInfo)
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("查询节点列表失败")
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetNodeTraffic 查询节点流量曲线
// @Tags PcdnNode
// @Summary 查询节点流量
// @Security ApiKeyAuth
// @Param data query request.TrafficQuery true "查询条件"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pcdn/admin/node/traffic [get]
func (a *NodeApi) GetNodeTraffic(c *gin.Context) {
	var q request.TrafficQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := trafficService.GetTrafficPoints(c.Request.Context(), q.NodeID, q.Iface, q.Start, q.End)
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("查询节点流量失败")
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetNode95 查询节点 95 值（日/月）
// @Tags PcdnNode
// @Summary 查询节点95值
// @Security ApiKeyAuth
// @Param nodeId query int true "节点ID"
// @Param periodType query string false "周期类型 day/month"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pcdn/admin/node/n95 [get]
func (a *NodeApi) GetNode95(c *gin.Context) {
	nodeID, _ := strconv.ParseUint(c.Query("nodeId"), 10, 64)
	periodType := c.Query("periodType")
	list, err := node95Service.GetNode95List(c.Request.Context(), uint(nodeID), periodType)
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("查询节点95值失败")
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

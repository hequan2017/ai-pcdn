package api

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// PortalApi 个人门户接口
type PortalApi struct{}

// Register 个人注册（公开）
// @Tags PcdnPortal
// @Summary 个人注册
// @accept application/json
// @Produce application/json
// @Param data body request.PortalRegister true "注册信息"
// @Success 200 {object} response.Response{msg=string} "注册成功"
// @Router /pcdn/portal/register [post]
func (a *PortalApi) Register(c *gin.Context) {
	var req request.PortalRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := portalService.Register(c.Request.Context(), req)
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("portal注册失败")
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"userId": user.ID, "username": user.Username}, "注册成功", c)
}

// Login 个人登录（公开）
// @Tags PcdnPortal
// @Summary 个人登录
// @accept application/json
// @Produce application/json
// @Param data body request.PortalLogin true "登录信息"
// @Success 200 {object} response.Response{data=object,msg=string} "登录成功"
// @Router /pcdn/portal/login [post]
func (a *PortalApi) Login(c *gin.Context) {
	var req request.PortalLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, token, err := portalService.Login(c.Request.Context(), req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"user": user, "token": token}, "登录成功", c)
}

// AddNode 自助上机：添加节点（需登录，返回凭证+一键安装命令，凭证仅展示一次）
// @Tags PcdnPortal
// @Summary 添加节点（自助上机）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PortalAddNode true "节点信息"
// @Success 200 {object} response.Response{data=object,msg=string} "添加成功"
// @Router /pcdn/portal/addNode [post]
func (a *PortalApi) AddNode(c *gin.Context) {
	ownerID := utils.GetUserID(c)
	var req request.PortalAddNode
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	node, token, err := portalService.AddNode(c.Request.Context(), ownerID, req)
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("添加节点失败")
		response.FailWithMessage("添加失败", c)
		return
	}
	installScript := fmt.Sprintf("curl -fsSL https://%s/pcdn/portal/install/%s.sh | bash", c.Request.Host, token)
	response.OkWithDetailed(gin.H{
		"nodeSn":        node.NodeSn,
		"token":         token,
		"installScript": installScript,
	}, "添加成功，凭证仅展示一次，请立即保存", c)
}

// MyNodes 我的节点列表（需登录）
// @Tags PcdnPortal
// @Summary 我的节点列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param data query request.PortalNodeSearch true "查询条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /pcdn/portal/myNodes [get]
func (a *PortalApi) MyNodes(c *gin.Context) {
	ownerID := utils.GetUserID(c)
	var info request.PortalNodeSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := portalService.MyNodes(c.Request.Context(), ownerID, info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// MyBills 我的账单（需登录）
// @Tags PcdnPortal
// @Summary 我的账单
// @Security ApiKeyAuth
// @Produce application/json
// @Param data query request.BillSearch true "查询条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /pcdn/portal/myBills [get]
func (a *PortalApi) MyBills(c *gin.Context) {
	ownerID := utils.GetUserID(c)
	var info request.BillSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := billService.GetBillsByOwner(c.Request.Context(), ownerID, info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// MyNodeTraffic 我的节点流量（需登录，强制校验归属）
// @Tags PcdnPortal
// @Summary 我的节点流量
// @Security ApiKeyAuth
// @Produce application/json
// @Param data query request.TrafficQuery true "查询条件"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pcdn/portal/myNodeTraffic [get]
func (a *PortalApi) MyNodeTraffic(c *gin.Context) {
	ownerID := utils.GetUserID(c)
	var q request.TrafficQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ok, err := portalService.CheckNodeOwner(c.Request.Context(), q.NodeID, ownerID)
	if err != nil || !ok {
		response.FailWithMessage("无权访问该节点", c)
		return
	}
	list, err := trafficService.GetTrafficPoints(c.Request.Context(), q.NodeID, q.Iface, q.Start, q.End)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

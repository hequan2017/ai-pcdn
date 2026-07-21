package api

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// SettlementApi 大厂结算单接口
type SettlementApi struct{}

// Import 导入结算单
// @Tags PcdnSettlement
// @Summary 导入结算单
// @Security ApiKeyAuth
// @Router /pcdn/admin/settlement/import [post]
func (a *SettlementApi) Import(c *gin.Context) {
	var req request.SettlementImport
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := settlementService.Import(c.Request.Context(), req); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("导入结算单失败")
		response.FailWithMessage("导入失败", c)
		return
	}
	response.OkWithMessage("导入成功", c)
}

// Recheck 重新核对
// @Tags PcdnSettlement
// @Summary 重新核对账期
// @Security ApiKeyAuth
// @Router /pcdn/admin/settlement/recheck [put]
func (a *SettlementApi) Recheck(c *gin.Context) {
	period := c.Query("period")
	n, err := settlementService.Recheck(c.Request.Context(), period)
	if err != nil {
		response.FailWithMessage("核对失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"count": n}, "核对完成", c)
}

// GetList 结算单列表
// @Tags PcdnSettlement
// @Summary 结算单列表
// @Security ApiKeyAuth
// @Router /pcdn/admin/settlement/list [get]
func (a *SettlementApi) GetList(c *gin.Context) {
	var info request.SettlementSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := settlementService.GetList(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// RevenueSummary 应收汇总
// @Tags PcdnSettlement
// @Summary 应收汇总
// @Security ApiKeyAuth
// @Router /pcdn/admin/settlement/revenue [get]
func (a *SettlementApi) RevenueSummary(c *gin.Context) {
	var req request.RevenueSummaryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sum, err := settlementService.RevenueSummary(c.Request.Context(), req)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(sum, "获取成功", c)
}

// Delete 删除结算单
// @Tags PcdnSettlement
// @Summary 删除结算单
// @Security ApiKeyAuth
// @Router /pcdn/admin/settlement/delete [delete]
func (a *SettlementApi) Delete(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := settlementService.Delete(c.Request.Context(), req.Uint()); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

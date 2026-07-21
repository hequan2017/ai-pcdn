package api

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// BillApi 采购账单接口
type BillApi struct{}

// GenerateBill 生成某账期账单
// @Tags PcdnBill
// @Summary 生成账单
// @Security ApiKeyAuth
// @Param data body request.GenerateBillReq true "账期"
// @Success 200 {object} response.Response{data=object,msg=string} "生成成功"
// @Router /pcdn/admin/bill/generate [post]
func (a *BillApi) GenerateBill(c *gin.Context) {
	var req request.GenerateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	n, err := billService.GenerateMonthlyBill(c.Request.Context(), req.Period)
	if err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("生成账单失败")
		response.FailWithMessage("生成失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"created": n}, "生成成功", c)
}

// GetBillList 账单列表
// @Tags PcdnBill
// @Summary 账单列表
// @Security ApiKeyAuth
// @Router /pcdn/admin/bill/list [get]
func (a *BillApi) GetBillList(c *gin.Context) {
	var info request.BillSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := billService.GetBillList(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// GetBill 账单详情
// @Tags PcdnBill
// @Summary 账单详情
// @Security ApiKeyAuth
// @Router /pcdn/admin/bill/find [get]
func (a *BillApi) GetBill(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bill, err := billService.GetBill(c.Request.Context(), req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(bill, "查询成功", c)
}

// ApproveBill 审核通过
// @Tags PcdnBill
// @Summary 审核通过
// @Security ApiKeyAuth
// @Router /pcdn/admin/bill/approve [put]
func (a *BillApi) ApproveBill(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := billService.ApproveBill(c.Request.Context(), req.Uint(), utils.GetUserID(c)); err != nil {
		response.FailWithMessage("审核失败", c)
		return
	}
	response.OkWithMessage("已审核", c)
}

// RejectBill 驳回
// @Tags PcdnBill
// @Summary 驳回
// @Security ApiKeyAuth
// @Router /pcdn/admin/bill/reject [put]
func (a *BillApi) RejectBill(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := billService.RejectBill(c.Request.Context(), req.ID, utils.GetUserID(c), req.Remark); err != nil {
		response.FailWithMessage("驳回失败", c)
		return
	}
	response.OkWithMessage("已驳回", c)
}

// PayBill 付款
// @Tags PcdnBill
// @Summary 付款
// @Security ApiKeyAuth
// @Router /pcdn/admin/bill/pay [put]
func (a *BillApi) PayBill(c *gin.Context) {
	var req request.PayBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := billService.PayBill(c.Request.Context(), req, utils.GetUserID(c)); err != nil {
		response.FailWithMessage("付款失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已付款", c)
}

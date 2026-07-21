package api

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// ReleaseApi agent 版本发布接口（admin）
type ReleaseApi struct{}

// CreateRelease 发布版本
// @Tags PcdnRelease
// @Summary 发布版本
// @Security ApiKeyAuth
// @Router /pcdn/admin/release/create [post]
func (a *ReleaseApi) CreateRelease(c *gin.Context) {
	var r model.PcdnAgentRelease
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.CreateRelease(c.Request.Context(), &r); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("发布版本失败")
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithDetailed(r, "创建成功", c)
}

// DeleteRelease 删除版本
// @Tags PcdnRelease
// @Summary 删除版本
// @Security ApiKeyAuth
// @Router /pcdn/admin/release/delete [delete]
func (a *ReleaseApi) DeleteRelease(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.DeleteRelease(c.Request.Context(), req.Uint()); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateRelease 更新版本
// @Tags PcdnRelease
// @Summary 更新版本
// @Security ApiKeyAuth
// @Router /pcdn/admin/release/update [put]
func (a *ReleaseApi) UpdateRelease(c *gin.Context) {
	var r model.PcdnAgentRelease
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.UpdateRelease(c.Request.Context(), r); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetRelease 查询版本
// @Tags PcdnRelease
// @Summary 查询版本
// @Security ApiKeyAuth
// @Router /pcdn/admin/release/find [get]
func (a *ReleaseApi) GetRelease(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	r, err := releaseService.GetRelease(c.Request.Context(), req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(r, "查询成功", c)
}

// GetReleaseList 版本列表
// @Tags PcdnRelease
// @Summary 版本列表
// @Security ApiKeyAuth
// @Router /pcdn/admin/release/list [get]
func (a *ReleaseApi) GetReleaseList(c *gin.Context) {
	var info request.ReleaseSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := releaseService.GetReleaseList(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

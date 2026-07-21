package api

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
	"github.com/gin-gonic/gin"
)

// AlarmApi 告警规则与记录接口
type AlarmApi struct{}

// CreateAlarmRule 创建告警规则
// @Tags PcdnAlarm
// @Summary 创建告警规则
// @Security ApiKeyAuth
// @Param data body model.PcdnAlarmRule true "规则"
// @Success 200 {object} response.Response{data=model.PcdnAlarmRule,msg=string} "创建成功"
// @Router /pcdn/admin/alarm/rule/create [post]
func (a *AlarmApi) CreateAlarmRule(c *gin.Context) {
	var rule model.PcdnAlarmRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alarmRuleService.CreateAlarmRule(c.Request.Context(), &rule); err != nil {
		logger.WithCtx(c.Request.Context()).Mod("pcdn").Err(err).Error("创建告警规则失败")
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithDetailed(rule, "创建成功", c)
}

// DeleteAlarmRule 删除告警规则
// @Tags PcdnAlarm
// @Summary 删除告警规则
// @Security ApiKeyAuth
// @Router /pcdn/admin/alarm/rule/delete [delete]
func (a *AlarmApi) DeleteAlarmRule(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alarmRuleService.DeleteAlarmRule(c.Request.Context(), req.Uint()); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAlarmRuleByIds 批量删除告警规则
// @Tags PcdnAlarm
// @Summary 批量删除告警规则
// @Security ApiKeyAuth
// @Router /pcdn/admin/alarm/rule/deleteByIds [delete]
func (a *AlarmApi) DeleteAlarmRuleByIds(c *gin.Context) {
	var req commonReq.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := make([]uint, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id <= 0 {
			response.FailWithMessage("非法ID", c)
			return
		}
		ids = append(ids, uint(id))
	}
	if err := alarmRuleService.DeleteAlarmRuleByIds(c.Request.Context(), ids); err != nil {
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAlarmRule 更新告警规则
// @Tags PcdnAlarm
// @Summary 更新告警规则
// @Security ApiKeyAuth
// @Router /pcdn/admin/alarm/rule/update [put]
func (a *AlarmApi) UpdateAlarmRule(c *gin.Context) {
	var rule model.PcdnAlarmRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alarmRuleService.UpdateAlarmRule(c.Request.Context(), rule); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetAlarmRule 查询告警规则
// @Tags PcdnAlarm
// @Summary 查询告警规则
// @Security ApiKeyAuth
// @Router /pcdn/admin/alarm/rule/find [get]
func (a *AlarmApi) GetAlarmRule(c *gin.Context) {
	var req commonReq.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	rule, err := alarmRuleService.GetAlarmRule(c.Request.Context(), req.Uint())
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(rule, "查询成功", c)
}

// GetAlarmRuleList 告警规则列表
// @Tags PcdnAlarm
// @Summary 告警规则列表
// @Security ApiKeyAuth
// @Router /pcdn/admin/alarm/rule/list [get]
func (a *AlarmApi) GetAlarmRuleList(c *gin.Context) {
	var info request.AlarmRuleSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := alarmRuleService.GetAlarmRuleList(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// GetAlarmRecordList 告警记录列表
// @Tags PcdnAlarm
// @Summary 告警记录列表
// @Security ApiKeyAuth
// @Router /pcdn/admin/alarm/record/list [get]
func (a *AlarmApi) GetAlarmRecordList(c *gin.Context) {
	var info request.AlarmRecordSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := alarmRecordService.GetAlarmRecordList(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

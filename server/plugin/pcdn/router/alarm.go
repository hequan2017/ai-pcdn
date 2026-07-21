package router

import "github.com/gin-gonic/gin"

// AlarmRouter 告警路由（admin 组）
type AlarmRouter struct{}

// InitAlarmRouter 挂载告警规则与记录路由
func (r *AlarmRouter) InitAlarmRouter(group *gin.RouterGroup) {
	group.GET("alarm/rule/list", alarmApi.GetAlarmRuleList)
	group.GET("alarm/rule/find", alarmApi.GetAlarmRule)
	group.POST("alarm/rule/create", alarmApi.CreateAlarmRule)
	group.PUT("alarm/rule/update", alarmApi.UpdateAlarmRule)
	group.DELETE("alarm/rule/delete", alarmApi.DeleteAlarmRule)
	group.DELETE("alarm/rule/deleteByIds", alarmApi.DeleteAlarmRuleByIds)
	group.GET("alarm/record/list", alarmApi.GetAlarmRecordList)
}

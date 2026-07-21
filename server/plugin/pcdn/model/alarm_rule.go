package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// 告警指标
const (
	AlarmMetricOffline      = "offline"       // 节点离线（心跳超时）
	AlarmMetricBandwidthLow = "bandwidth_low" // 实时带宽低于阈值
	AlarmMetricP95High      = "p95_high"      // 当日 95 值高于阈值
	AlarmMetricAgentDown    = "agent_down"    // agent 上报中断
)

// 告警范围
const (
	AlarmScopeAll   = "all"   // 全部节点
	AlarmScopeGroup = "group" // 按分组
	AlarmScopeNode  = "node"  // 单节点
)

// 告警状态
const (
	AlarmStatusFiring   = "firing"   // 触发中
	AlarmStatusResolved = "resolved" // 已恢复
)

// PcdnAlarmRule 告警规则
type PcdnAlarmRule struct {
	global.GVA_MODEL
	Name         string         `json:"name" gorm:"column:name;type:varchar(64);comment:规则名"`
	ScopeType    string         `json:"scopeType" gorm:"column:scope_type;type:varchar(16);comment:范围 all/group/node"`
	ScopeValue   string         `json:"scopeValue" gorm:"column:scope_value;type:varchar(64);comment:范围值 group_id/node_id"`
	Metric       string         `json:"metric" gorm:"column:metric;type:varchar(32);comment:指标"`
	Threshold    int64          `json:"threshold" gorm:"column:threshold;comment:阈值bps(bandwidth_low/p95_high)"`
	DurationSec  int            `json:"durationSec" gorm:"column:duration_sec;comment:持续秒数(agent_down)"`
	NotifyConfig datatypes.JSON `json:"notifyConfig" gorm:"column:notify_config;comment:通知配置;type:json" swaggertype:"object"`
	Enabled      bool           `json:"enabled" gorm:"column:enabled;comment:是否启用"`
	DeptID       uint           `json:"deptId" gorm:"column:dept_id;comment:归属部门(数据权限);index"`
	CreatedBy    uint           `json:"createdBy" gorm:"column:created_by;comment:创建人(数据权限)"`
}

// TableName 告警规则表名
func (PcdnAlarmRule) TableName() string {
	return "gva_pcdn_alarm_rule"
}

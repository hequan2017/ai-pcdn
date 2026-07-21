package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// PcdnAlarmRecord 告警记录（触发/恢复流水）
type PcdnAlarmRecord struct {
	global.GVA_MODEL
	RuleID       uint       `json:"ruleId" gorm:"column:rule_id;index;comment:规则ID"`
	RuleName     string     `json:"ruleName" gorm:"column:rule_name;type:varchar(64);comment:规则名"`
	NodeID       uint       `json:"nodeId" gorm:"column:node_id;index;comment:节点ID"`
	NodeSn       string     `json:"nodeSn" gorm:"column:node_sn;type:varchar(64);comment:节点SN"`
	Metric       string     `json:"metric" gorm:"column:metric;type:varchar(32);comment:指标"`
	TriggerValue int64      `json:"triggerValue" gorm:"column:trigger_value;comment:触发值"`
	Status       string     `json:"status" gorm:"column:status;type:varchar(16);index;comment:状态 firing/resolved"`
	FiredAt      time.Time  `json:"firedAt" gorm:"column:fired_at;comment:触发时间"`
	ResolvedAt   *time.Time `json:"resolvedAt" gorm:"column:resolved_at;comment:恢复时间"`
	NotifyCount  int        `json:"notifyCount" gorm:"column:notify_count;comment:通知次数"`
}

// TableName 告警记录表名
func (PcdnAlarmRecord) TableName() string {
	return "gva_pcdn_alarm_record"
}

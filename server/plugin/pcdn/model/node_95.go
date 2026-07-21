package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 95值周期类型
const (
	PeriodTypeDay   = "day"   // 日
	PeriodTypeMonth = "month" // 月
)

// 95值状态
const (
	Node95StatusRolling = "rolling" // 滚动（周期内每日重算）
	Node95StatusFrozen  = "frozen"  // 冻结（周期结束定型，作为账单依据）
)

// PcdnNode95 节点 95 值结果（日/月）
type PcdnNode95 struct {
	global.GVA_MODEL
	NodeID        uint       `json:"nodeId" gorm:"column:node_id;comment:节点ID;index:idx_period,priority:1"`
	PeriodType    string     `json:"periodType" gorm:"column:period_type;type:varchar(8);comment:周期类型 day/month;index:idx_period,priority:2"`
	PeriodStart   time.Time  `json:"periodStart" gorm:"column:period_start;comment:周期开始;index:idx_period,priority:3"`
	PeriodEnd     time.Time  `json:"periodEnd" gorm:"column:period_end;comment:周期结束"`
	Rx95Bps       int64      `json:"rx95Bps" gorm:"column:rx_95_bps;comment:下行95值bps"`
	Tx95Bps       int64      `json:"tx95Bps" gorm:"column:tx_95_bps;comment:上行95值bps"`
	Combined95Bps int64      `json:"combined95Bps" gorm:"column:combined_95_bps;comment:上下行合计95值bps"`
	SampleCount   int        `json:"sampleCount" gorm:"column:sample_count;comment:采样点数"`
	Status        string     `json:"status" gorm:"column:status;type:varchar(16);comment:状态 rolling/frozen"`
	FrozenAt      *time.Time `json:"frozenAt" gorm:"column:frozen_at;comment:冻结时间"`
}

// TableName 95值表名
func (PcdnNode95) TableName() string {
	return "gva_pcdn_node_95"
}

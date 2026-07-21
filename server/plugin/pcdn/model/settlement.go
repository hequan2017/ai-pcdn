package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 结算单核对状态
const (
	SettlementStatusPending = "pending" // 待核对
	SettlementStatusMatched = "matched" // 已核对（差异在阈值内）
	SettlementStatusDiff    = "diff"    // 有差异
)

// PcdnSettlement 大厂结算单（销售侧应收依据）
type PcdnSettlement struct {
	global.GVA_MODEL
	Period         string  `json:"period" gorm:"column:period;type:varchar(7);comment:账期YYYY-MM;index"`
	Platform       string  `json:"platform" gorm:"column:platform;type:varchar(32);comment:大厂平台;index"`
	NodeID         uint    `json:"nodeId" gorm:"column:node_id;comment:节点ID"`
	NodeSn         string  `json:"nodeSn" gorm:"column:node_sn;type:varchar(64);comment:节点SN"`
	PlatformNodeID string  `json:"platformNodeId" gorm:"column:platform_node_id;type:varchar(128);comment:大厂侧节点ID"`
	Revenue        float64 `json:"revenue" gorm:"column:revenue;comment:大厂结算收入(元)"`
	TrafficBps     int64   `json:"trafficBps" gorm:"column:traffic_bps;comment:大厂侧流量bps"`
	OurTrafficBps  int64   `json:"ourTrafficBps" gorm:"column:our_traffic_bps;comment:自采集流量bps(核对)"`
	DiffPercent    float64 `json:"diffPercent" gorm:"column:diff_percent;comment:差异百分比"`
	Status         string  `json:"status" gorm:"column:status;type:varchar(16);index;comment:pending/matched/diff"`
	Remark         string  `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
	DeptID         uint    `json:"deptId" gorm:"column:dept_id;comment:归属部门(数据权限);index"`
	CreatedBy      uint    `json:"createdBy" gorm:"column:created_by;comment:创建人(数据权限)"`
}

// TableName 结算单表名
func (PcdnSettlement) TableName() string {
	return "gva_pcdn_settlement"
}

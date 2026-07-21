package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// PcdnNodeTrafficPoint 节点流量分钟峰值点（95值计算的原始输入）
// 注：按 (node_id, window_start, iface_name) 复合唯一索引去重，保证重复上报幂等。
// 数据量级：500节点×2网卡×1440分钟/天 ≈ 144万行/天，建议生产环境按 window_start 做月度分区（AutoMigrate 不支持分区，需运维建表时手动分区）。
type PcdnNodeTrafficPoint struct {
	global.GVA_MODEL
	NodeID      uint      `json:"nodeId" gorm:"column:node_id;comment:节点ID;uniqueIndex:uniq_node_win,priority:1"`
	IfaceName   string    `json:"ifaceName" gorm:"column:iface_name;type:varchar(64);comment:网卡名;uniqueIndex:uniq_node_win,priority:2"`
	WindowStart time.Time `json:"windowStart" gorm:"column:window_start;comment:窗口开始(分钟级);uniqueIndex:uniq_node_win,priority:3"`
	RxMaxBps    int64     `json:"rxMaxBps" gorm:"column:rx_max_bps;comment:下行峰值bps"`
	TxMaxBps    int64     `json:"txMaxBps" gorm:"column:tx_max_bps;comment:上行峰值bps"`
}

// TableName 流量点表名
func (PcdnNodeTrafficPoint) TableName() string {
	return "gva_pcdn_node_traffic_point"
}

package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// PcdnNodeIface 节点网卡（采集维度，一节点多网卡）
type PcdnNodeIface struct {
	global.GVA_MODEL
	NodeID    uint   `json:"nodeId" gorm:"column:node_id;comment:节点ID;index"`
	IfaceName string `json:"ifaceName" gorm:"column:iface_name;type:varchar(64);comment:网卡名"`
	Mac       string `json:"mac" gorm:"column:mac;type:varchar(64);comment:MAC地址"`
	Enabled   bool   `json:"enabled" gorm:"column:enabled;comment:是否启用"`
}

// TableName 网卡表名
func (PcdnNodeIface) TableName() string {
	return "gva_pcdn_node_iface"
}

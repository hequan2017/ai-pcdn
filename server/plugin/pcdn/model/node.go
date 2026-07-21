package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// 节点状态
const (
	NodeStatusPending  = "pending"  // 待上机（已生成凭证，agent 未激活）
	NodeStatusOnline   = "online"   // 在线
	NodeStatusOffline  = "offline"  // 离线（心跳超时）
	NodeStatusAbnormal = "abnormal" // 异常（采集/上报异常）
	NodeStatusDisabled = "disabled" // 已停用
)

// 计费模式
const (
	BillingModeMonthly = "monthly" // 包月
	BillingModeP95     = "p95"     // 95 计费
)

// PcdnNode PCDN 节点主表（一切业务的主数据）
type PcdnNode struct {
	global.GVA_MODEL
	NodeSn         string         `json:"nodeSn" gorm:"column:node_sn;type:varchar(64);uniqueIndex;comment:节点序列号"`
	TokenHash      string         `json:"-" gorm:"column:token_hash;type:varchar(128);comment:节点上报token哈希"`
	OwnerUserID    uint           `json:"ownerUserId" gorm:"column:owner_user_id;comment:归属用户ID(贡献者);index"`
	OwnerName      string         `json:"ownerName" gorm:"column:owner_name;type:varchar(64);comment:归属用户名"`
	Contact        string         `json:"contact" gorm:"column:contact;type:varchar(128);comment:联系方式"`
	Region         string         `json:"region" gorm:"column:region;type:varchar(64);comment:地域省/市"`
	Isp            string         `json:"isp" gorm:"column:isp;type:varchar(32);comment:运营商"`
	Platform       string         `json:"platform" gorm:"column:platform;type:varchar(32);comment:接入大厂"`
	PlatformNodeID string         `json:"platformNodeId" gorm:"column:platform_node_id;type:varchar(128);comment:大厂侧节点ID"`
	GroupID        uint           `json:"groupId" gorm:"column:group_id;comment:分组ID;index"`
	Tags           datatypes.JSON `json:"tags" gorm:"column:tags;comment:标签;type:json" swaggertype:"array,string"`
	Hostname       string         `json:"hostname" gorm:"column:hostname;type:varchar(128);comment:主机名"`
	InnerIP        string         `json:"innerIp" gorm:"column:inner_ip;type:varchar(64);comment:内网IP"`
	ReportIP       string         `json:"reportIp" gorm:"column:report_ip;type:varchar(64);comment:上报公网IP"`
	OS             string         `json:"os" gorm:"column:os;type:varchar(64);comment:操作系统"`
	Status         string         `json:"status" gorm:"column:status;type:varchar(32);comment:状态;index"`
	AgentVersion   string         `json:"agentVersion" gorm:"column:agent_version;type:varchar(32);comment:agent版本"`
	LastHeartbeatAt *time.Time    `json:"lastHeartbeatAt" gorm:"column:last_heartbeat_at;comment:最后心跳时间"`
	BillingMode    string         `json:"billingMode" gorm:"column:billing_mode;type:varchar(16);comment:计费模式 monthly/p95"`
	MonthlyPrice   float64        `json:"monthlyPrice" gorm:"column:monthly_price;comment:包月价"`
	P95UnitPrice   float64        `json:"p95UnitPrice" gorm:"column:p95_unit_price;comment:p95单价(元/Mbps)"`
	ContractPeriod string         `json:"contractPeriod" gorm:"column:contract_period;type:varchar(32);comment:合同周期"`
	DeptID         uint           `json:"deptId" gorm:"column:dept_id;comment:归属部门ID(数据权限);index"`
	CreatedBy      uint           `json:"createdBy" gorm:"column:created_by;comment:创建人(数据权限)"`
	Remark         string         `json:"remark" gorm:"column:remark;type:varchar(255);comment:备注"`
}

// TableName 节点表名
func (PcdnNode) TableName() string {
	return "gva_pcdn_node"
}

package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// 账单状态
const (
	BillStatusDraft    = "draft"    // 待审核
	BillStatusApproved = "approved" // 已审核（待付款）
	BillStatusPaid     = "paid"     // 已付款
	BillStatusRejected = "rejected" // 已驳回
)

// BillDetail 账单明细（单节点），存于 PcdnBill.Details JSON 数组
type BillDetail struct {
	NodeID      uint    `json:"nodeId"`
	NodeSn      string  `json:"nodeSn"`
	BillingMode string  `json:"billingMode"` // monthly/p95
	Value       float64 `json:"value"`       // p95=Mbps，monthly=包月价
	UnitPrice   float64 `json:"unitPrice"`   // p95=元/Mbps，monthly=1
	Amount      float64 `json:"amount"`      // 该节点应付金额
}

// PcdnBill 采购侧月账单（按贡献者 owner 汇总）
type PcdnBill struct {
	global.GVA_MODEL
	Period       string         `json:"period" gorm:"column:period;type:varchar(7);comment:账期YYYY-MM;index"`
	OwnerUserID  uint           `json:"ownerUserId" gorm:"column:owner_user_id;index;comment:归属用户ID"`
	OwnerName    string         `json:"ownerName" gorm:"column:owner_name;type:varchar(64);comment:归属用户名"`
	OwnerContact string         `json:"ownerContact" gorm:"column:owner_contact;type:varchar(128);comment:联系方式"`
	NodeCount    int            `json:"nodeCount" gorm:"column:node_count;comment:节点数"`
	Details      datatypes.JSON `json:"details" gorm:"column:details;comment:明细;type:json" swaggertype:"array,object"`
	TotalAmount  float64        `json:"totalAmount" gorm:"column:total_amount;comment:应付总额"`
	Status       string         `json:"status" gorm:"column:status;type:varchar(16);index;comment:draft/approved/paid/rejected"`
	AuditedBy    uint           `json:"auditedBy" gorm:"column:audited_by;comment:审核人"`
	AuditedAt    *time.Time     `json:"auditedAt" gorm:"column:audited_at;comment:审核时间"`
	PaidAmount   float64        `json:"paidAmount" gorm:"column:paid_amount;comment:实付金额"`
	PayMethod    string         `json:"payMethod" gorm:"column:pay_method;type:varchar(32);comment:付款方式"`
	PayNo        string         `json:"payNo" gorm:"column:pay_no;type:varchar(128);comment:付款流水号"`
	PaidAt       *time.Time     `json:"paidAt" gorm:"column:paid_at;comment:付款时间"`
	PaidBy       uint           `json:"paidBy" gorm:"column:paid_by;comment:付款操作人"`
	Remark       string         `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
	DeptID       uint           `json:"deptId" gorm:"column:dept_id;comment:归属部门(数据权限);index"`
	CreatedBy    uint           `json:"createdBy" gorm:"column:created_by;comment:创建人(数据权限)"`
}

// TableName 账单表名
func (PcdnBill) TableName() string {
	return "gva_pcdn_bill"
}

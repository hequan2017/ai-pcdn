package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// BillSearch 账单查询
type BillSearch struct {
	Period    string `json:"period" form:"period"`
	Status    string `json:"status" form:"status"`
	OwnerName string `json:"ownerName" form:"ownerName"`
	commonReq.PageInfo
}

// PayBillReq 付款请求
type PayBillReq struct {
	ID         uint    `json:"id"`
	PayMethod  string  `json:"payMethod"`
	PayNo      string  `json:"payNo"`
	PaidAmount float64 `json:"paidAmount"`
	Remark     string  `json:"remark"`
}

// GenerateBillReq 生成账单请求
type GenerateBillReq struct {
	Period string `json:"period"`
}

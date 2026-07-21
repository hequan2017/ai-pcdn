package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// SettlementSearch 结算单查询
type SettlementSearch struct {
	Period   string `json:"period" form:"period"`
	Platform string `json:"platform" form:"platform"`
	Status   string `json:"status" form:"status"`
	NodeSn   string `json:"nodeSn" form:"nodeSn"`
	commonReq.PageInfo
}

// SettlementImport 单条结算单导入
type SettlementImport struct {
	Period         string  `json:"period"`
	Platform       string  `json:"platform"`
	NodeSn         string  `json:"nodeSn"`
	PlatformNodeID string  `json:"platformNodeId"`
	Revenue        float64 `json:"revenue"`
	TrafficBps     int64   `json:"trafficBps"`
	Remark         string  `json:"remark"`
}

// RevenueSummaryReq 应收汇总查询
type RevenueSummaryReq struct {
	Period   string `json:"period" form:"period"`
	Platform string `json:"platform" form:"platform"`
}

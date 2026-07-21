package request

import (
	"time"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// NodeSearch 节点分页查询
type NodeSearch struct {
	NodeSn         string     `json:"nodeSn" form:"nodeSn"`
	Status         string     `json:"status" form:"status"`
	OwnerName      string     `json:"ownerName" form:"ownerName"`
	Platform       string     `json:"platform" form:"platform"`
	Region         string     `json:"region" form:"region"`
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	commonReq.PageInfo
}

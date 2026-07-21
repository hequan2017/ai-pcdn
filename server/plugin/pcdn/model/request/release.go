package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ReleaseSearch 版本查询
type ReleaseSearch struct {
	Version string `json:"version" form:"version"`
	Stable  *bool  `json:"stable" form:"stable"`
	commonReq.PageInfo
}

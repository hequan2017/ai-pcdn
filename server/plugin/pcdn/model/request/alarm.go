package request

import (
	"time"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// AlarmRuleSearch 告警规则查询
type AlarmRuleSearch struct {
	Name    string `json:"name" form:"name"`
	Metric  string `json:"metric" form:"metric"`
	Enabled *bool  `json:"enabled" form:"enabled"`
	commonReq.PageInfo
}

// AlarmRecordSearch 告警记录查询
type AlarmRecordSearch struct {
	NodeID uint       `json:"nodeId" form:"nodeId"`
	Status string     `json:"status" form:"status"`
	Start  *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	End    *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	commonReq.PageInfo
}

package request

import "time"

// TrafficQuery 流量点查询（admin 与 portal 共用，portal 侧另行校验节点归属）
type TrafficQuery struct {
	NodeID uint      `json:"nodeId" form:"nodeId"`
	Iface  string    `json:"iface" form:"iface"`
	Start  time.Time `json:"start" form:"start" time_format:"2006-01-02 15:04:05"`
	End    time.Time `json:"end" form:"end" time_format:"2006-01-02 15:04:05"`
}

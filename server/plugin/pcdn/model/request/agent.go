package request

import "time"

// TrafficReport agent 流量上报请求（批量分钟峰值点）
type TrafficReport struct {
	Points []TrafficPoint `json:"points"`
}

// TrafficPoint 单个流量分钟峰值点
type TrafficPoint struct {
	IfaceName   string    `json:"ifaceName"`   // 网卡名
	WindowStart time.Time `json:"windowStart"` // 窗口开始（分钟级，对齐到分钟）
	RxMaxBps    int64     `json:"rxMaxBps"`    // 下行峰值 bps
	TxMaxBps    int64     `json:"txMaxBps"`    // 上行峰值 bps
}

// AgentActivate agent 首次激活请求（回填硬件信息）
type AgentActivate struct {
	Hostname string       `json:"hostname"`
	OS       string       `json:"os"`
	InnerIP  string       `json:"innerIp"`
	Ifaces   []AgentIface `json:"ifaces"`
}

// AgentIface 网卡信息
type AgentIface struct {
	IfaceName string `json:"ifaceName"`
	Mac       string `json:"mac"`
	Enabled   bool   `json:"enabled"`
}

// AgentHeartbeat 心跳请求（允许空 body）
type AgentHeartbeat struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

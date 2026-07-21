package config

// Config PCDN 插件私有配置（从 config.yaml 的 pcdn: 段读取）
type Config struct {
	// AgentReportLimitPerMinute 单节点每分钟上报次数上限（限流）
	AgentReportLimitPerMinute int `mapstructure:"agent-report-limit-per-minute" yaml:"agent-report-limit-per-minute"`
	// HeartbeatTimeoutSec 心跳超时秒数，超过判离线
	HeartbeatTimeoutSec int `mapstructure:"heartbeat-timeout-sec" yaml:"heartbeat-timeout-sec"`
	// TrafficDetailRetentionDays 流量明细点保留天数，超出由运维归档
	TrafficDetailRetentionDays int `mapstructure:"traffic-detail-retention-days" yaml:"traffic-detail-retention-days"`
}

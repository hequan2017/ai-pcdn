package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// NotifyConfig 规则通知配置（存于 PcdnAlarmRule.NotifyConfig）
type NotifyConfig struct {
	WebhookURL string   `json:"webhookUrl"`           // 钉钉/企微 webhook
	AtMobiles  []string `json:"atMobiles"`             // @ 手机号
}

var notifyClient = &http.Client{Timeout: 5 * time.Second}

// SendNotify 按规则通知配置发送告警/恢复（当前支持钉钉/企微 webhook；邮件/短信预留扩展）
func SendNotify(rule *model.PcdnAlarmRule, rec *model.PcdnAlarmRecord, isFire bool) error {
	if len(rule.NotifyConfig) == 0 {
		return nil
	}
	var cfg NotifyConfig
	if err := json.Unmarshal(rule.NotifyConfig, &cfg); err != nil || cfg.WebhookURL == "" {
		return nil
	}
	state := "恢复 ✅"
	if isFire {
		state = "触发 🔥"
	}
	text := fmt.Sprintf("### PCDN 告警 %s\n- 规则: %s\n- 节点: %s\n- 指标: %s\n- 触发值: %s\n- 时间: %s",
		state, rule.Name, rec.NodeSn, rec.Metric, formatValue(rec.Metric, rec.TriggerValue), time.Now().Format("2006-01-02 15:04:05"))
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{"title": "PCDN告警", "text": text},
	}
	if len(cfg.AtMobiles) > 0 {
		payload["at"] = map[string]interface{}{"atMobiles": cfg.AtMobiles, "isAtAll": false}
	}
	body, _ := json.Marshal(payload)
	resp, err := notifyClient.Post(cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook status %d", resp.StatusCode)
	}
	return nil
}

// formatValue 按指标格式化触发值
func formatValue(metric string, v int64) string {
	switch metric {
	case model.AlarmMetricBandwidthLow, model.AlarmMetricP95High:
		return fmt.Sprintf("%.2f Mbps", float64(v)/1e6)
	case model.AlarmMetricAgentDown:
		return fmt.Sprintf("%ds 无上报", v)
	default:
		return fmt.Sprintf("%d", v)
	}
}

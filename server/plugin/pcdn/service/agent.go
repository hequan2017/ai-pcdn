package service

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"gorm.io/gorm"
)

// AgentService agent 上报相关服务
type AgentService struct{}

// Activate 激活节点：回填硬件信息、转 online、刷新心跳、重建网卡清单。事务保证一致。
func (s *AgentService) Activate(ctx context.Context, nodeID uint, reportIP string, req request.AgentActivate) error {
	now := time.Now()
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"hostname":          req.Hostname,
			"os":                req.OS,
			"inner_ip":          req.InnerIP,
			"report_ip":         reportIP,
			"status":            model.NodeStatusOnline,
			"last_heartbeat_at": now,
		}
		if err := tx.Model(&model.PcdnNode{}).Where("id = ?", nodeID).Updates(updates).Error; err != nil {
			return err
		}
		if len(req.Ifaces) > 0 {
			if err := tx.Where("node_id = ?", nodeID).Delete(&model.PcdnNodeIface{}).Error; err != nil {
				return err
			}
			rows := make([]model.PcdnNodeIface, 0, len(req.Ifaces))
			for _, f := range req.Ifaces {
				rows = append(rows, model.PcdnNodeIface{
					NodeID:    nodeID,
					IfaceName: f.IfaceName,
					Mac:       f.Mac,
					Enabled:   f.Enabled,
				})
			}
			return tx.Create(&rows).Error
		}
		return nil
	})
}

// Heartbeat 心跳：刷新最后心跳时间与上报IP，非禁用节点恢复 online。
func (s *AgentService) Heartbeat(ctx context.Context, nodeID uint, reportIP string, req request.AgentHeartbeat) error {
	now := time.Now()
	updates := map[string]interface{}{
		"last_heartbeat_at": now,
		"report_ip":         reportIP,
		"status":            model.NodeStatusOnline,
	}
	if req.Hostname != "" {
		updates["hostname"] = req.Hostname
	}
	if req.OS != "" {
		updates["os"] = req.OS
	}
	return global.GVA_DB.WithContext(ctx).
		Model(&model.PcdnNode{}).
		Where("id = ? AND status != ?", nodeID, model.NodeStatusDisabled).
		Updates(updates).Error
}

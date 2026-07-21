package service

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
)

// MarkOfflineNodes 把心跳超时的 online 节点置为 offline，返回受影响行数。
// 由调度任务周期调用（默认阈值见 config / scheduler）。
func (s *NodeService) MarkOfflineNodes(ctx context.Context, timeout time.Duration) (int64, error) {
	threshold := time.Now().Add(-timeout)
	res := global.GVA_DB.WithContext(ctx).
		Model(&model.PcdnNode{}).
		Where("status = ? AND last_heartbeat_at IS NOT NULL AND last_heartbeat_at < ?", model.NodeStatusOnline, threshold).
		Update("status", model.NodeStatusOffline)
	return res.RowsAffected, res.Error
}

package service

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"gorm.io/gorm/clause"
)

// TrafficService 流量点读写服务
type TrafficService struct{}

// BatchUpsertTrafficPoints 批量写入流量分钟峰值点。
// 依赖 (node_id, window_start, iface_name) 复合唯一索引实现幂等：重复上报同一窗口直接跳过。
func (s *TrafficService) BatchUpsertTrafficPoints(ctx context.Context, nodeID uint, points []request.TrafficPoint) error {
	if len(points) == 0 {
		return nil
	}
	rows := make([]model.PcdnNodeTrafficPoint, 0, len(points))
	for _, p := range points {
		rows = append(rows, model.PcdnNodeTrafficPoint{
			NodeID:      nodeID,
			IfaceName:   p.IfaceName,
			WindowStart: p.WindowStart,
			RxMaxBps:    p.RxMaxBps,
			TxMaxBps:    p.TxMaxBps,
		})
	}
	return global.GVA_DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&rows).Error
}

// GetTrafficPoints 查询节点在时间范围内的流量点（用于详情曲线）
func (s *TrafficService) GetTrafficPoints(ctx context.Context, nodeID uint, iface string, start, end time.Time) (list []model.PcdnNodeTrafficPoint, err error) {
	db := global.GVA_DB.WithContext(ctx).
		Where("node_id = ? AND window_start BETWEEN ? AND ?", nodeID, start, end)
	if iface != "" {
		db = db.Where("iface_name = ?", iface)
	}
	err = db.Order("window_start asc").Find(&list).Error
	return
}

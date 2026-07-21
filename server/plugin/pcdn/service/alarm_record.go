package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
)

// AlarmRecordService 告警记录查询
type AlarmRecordService struct{}

func (s *AlarmRecordService) GetAlarmRecordList(ctx context.Context, info request.AlarmRecordSearch) (list []model.PcdnAlarmRecord, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnAlarmRecord{})
	if info.NodeID != 0 {
		db = db.Where("node_id = ?", info.NodeID)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.Start != nil && info.End != nil {
		db = db.Where("fired_at BETWEEN ? AND ?", info.Start, info.End)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("id desc").Find(&list).Error
	return
}

package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
)

// AlarmRuleService 告警规则 CRUD
type AlarmRuleService struct{}

func (s *AlarmRuleService) CreateAlarmRule(ctx context.Context, rule *model.PcdnAlarmRule) error {
	return global.GVA_DB.WithContext(ctx).Create(rule).Error
}

func (s *AlarmRuleService) DeleteAlarmRule(ctx context.Context, id uint) error {
	return global.GVA_DB.WithContext(ctx).Delete(&model.PcdnAlarmRule{}, "id = ?", id).Error
}

func (s *AlarmRuleService) DeleteAlarmRuleByIds(ctx context.Context, ids []uint) error {
	return global.GVA_DB.WithContext(ctx).Delete(&[]model.PcdnAlarmRule{}, "id in ?", ids).Error
}

func (s *AlarmRuleService) UpdateAlarmRule(ctx context.Context, rule model.PcdnAlarmRule) error {
	// 用 Select 显式列，避免 GORM Updates 跳过零值（Enabled=false / Threshold=0 / DurationSec=0）
	return global.GVA_DB.WithContext(ctx).
		Model(&model.PcdnAlarmRule{}).
		Where("id = ?", rule.ID).
		Select("name", "scope_type", "scope_value", "metric", "threshold", "duration_sec", "notify_config", "enabled").
		Updates(&rule).Error
}

func (s *AlarmRuleService) GetAlarmRule(ctx context.Context, id uint) (rule model.PcdnAlarmRule, err error) {
	err = global.GVA_DB.WithContext(ctx).Where("id = ?", id).First(&rule).Error
	return
}

func (s *AlarmRuleService) GetAlarmRuleList(ctx context.Context, info request.AlarmRuleSearch) (list []model.PcdnAlarmRule, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnAlarmRule{})
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Metric != "" {
		db = db.Where("metric = ?", info.Metric)
	}
	if info.Enabled != nil {
		db = db.Where("enabled = ?", *info.Enabled)
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

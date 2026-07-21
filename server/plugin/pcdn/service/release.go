package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
)

// ReleaseService agent 版本发布服务（OTA）
type ReleaseService struct{}

func (s *ReleaseService) CreateRelease(ctx context.Context, r *model.PcdnAgentRelease) error {
	return global.GVA_DB.WithContext(ctx).Create(r).Error
}

func (s *ReleaseService) DeleteRelease(ctx context.Context, id uint) error {
	return global.GVA_DB.WithContext(ctx).Delete(&model.PcdnAgentRelease{}, "id = ?", id).Error
}

func (s *ReleaseService) UpdateRelease(ctx context.Context, r model.PcdnAgentRelease) error {
	// 用 Select 显式列，避免零值 Stable=false / Force=false 被跳过
	return global.GVA_DB.WithContext(ctx).Model(&model.PcdnAgentRelease{}).Where("id = ?", r.ID).
		Select("version", "download_url", "checksum", "stable", "force", "remark").Updates(&r).Error
}

func (s *ReleaseService) GetRelease(ctx context.Context, id uint) (r model.PcdnAgentRelease, err error) {
	err = global.GVA_DB.WithContext(ctx).Where("id = ?", id).First(&r).Error
	return
}

func (s *ReleaseService) GetReleaseList(ctx context.Context, info request.ReleaseSearch) (list []model.PcdnAgentRelease, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnAgentRelease{})
	if info.Version != "" {
		db = db.Where("version LIKE ?", "%"+info.Version+"%")
	}
	if info.Stable != nil {
		db = db.Where("stable = ?", *info.Stable)
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

// GetLatest 最新稳定版（agent 自升级查询）
func (s *ReleaseService) GetLatest(ctx context.Context) (r model.PcdnAgentRelease, err error) {
	err = global.GVA_DB.WithContext(ctx).Where("stable = ?", true).Order("id desc").First(&r).Error
	return
}

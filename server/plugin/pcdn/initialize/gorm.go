package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Gorm 注册 PCDN 业务表（AutoMigrate）
func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.PcdnNode),
		new(model.PcdnNodeIface),
		new(model.PcdnNodeTrafficPoint),
		new(model.PcdnNode95),
		new(model.PcdnAlarmRule),
		new(model.PcdnAlarmRecord),
		new(model.PcdnBill),
	)
	if err != nil {
		err = errors.Wrap(err, "PCDN注册表失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}

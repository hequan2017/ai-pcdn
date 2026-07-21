package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

// Dictionary 注册 PCDN 字典（阶段1 无字典，预留扩展点）
func Dictionary(ctx context.Context) {
	utils.RegisterDictionaries([]system.SysDictionary{}...)
}

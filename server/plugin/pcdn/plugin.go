package pcdn

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/initialize"
	interfaces "github.com/flipped-aurora/gin-vue-admin/server/utils/plugin/v2"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
}

// Register 插件注册入口：依次初始化 API 权限点、菜单、字典、数据表、路由
func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	// initialize.Viper() // 阶段1 暂不读取 config.yaml 的 pcdn 私有配置，相关参数用默认值
	initialize.Api(ctx)
	initialize.Menu(ctx)
	initialize.Dictionary(ctx)
	initialize.Gorm(ctx)
	initialize.Router(group)
}

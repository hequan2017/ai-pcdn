package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

// Menu 注册 PCDN 前端菜单（父菜单 + 子菜单；子菜单 Component 前缀必须是 plugin/pcdn/view/...）
func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{ParentId: 0, Path: "pcdn", Name: "pcdn", Component: "view/routerHolder.vue", Sort: 8, Meta: model.Meta{Title: "PCDN管理", Icon: "cloud-server"}},
		{Path: "pcdnNode", Name: "pcdnNode", Component: "plugin/pcdn/view/node/index.vue", Sort: 1, Meta: model.Meta{Title: "节点管理", Icon: "cloudy"}},
	}
	utils.RegisterMenus(entities...)
}

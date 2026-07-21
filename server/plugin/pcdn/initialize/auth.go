package initialize

import (
	"context"
	"strconv"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

// AdminAuthorityID GVA 默认超级管理员角色 ID（admin 账号主角色，见 source/system/user.go）
const AdminAuthorityID uint = 888

// pcdnMenuNames 与 menu.go 中 RegisterMenus 注册的 Name 一致（按 name 反查 ID）
var pcdnMenuNames = []string{
	"pcdn", "pcdnNode", "pcdnAlarmRule", "pcdnAlarmRecord",
	"pcdnBill", "pcdnSettlement", "pcdnProfit", "pcdnRelease",
}

// pcdnApis 与 api.go 中 RegisterApis 注册的 path+method 一致
var pcdnApis = []struct {
	Path   string
	Method string
}{
	{"/pcdn/admin/node/list", "GET"},
	{"/pcdn/admin/node/find", "GET"},
	{"/pcdn/admin/node/create", "POST"},
	{"/pcdn/admin/node/update", "PUT"},
	{"/pcdn/admin/node/delete", "DELETE"},
	{"/pcdn/admin/node/deleteByIds", "DELETE"},
	{"/pcdn/admin/node/traffic", "GET"},
	{"/pcdn/admin/node/n95", "GET"},
	{"/pcdn/admin/alarm/rule/list", "GET"},
	{"/pcdn/admin/alarm/rule/find", "GET"},
	{"/pcdn/admin/alarm/rule/create", "POST"},
	{"/pcdn/admin/alarm/rule/update", "PUT"},
	{"/pcdn/admin/alarm/rule/delete", "DELETE"},
	{"/pcdn/admin/alarm/rule/deleteByIds", "DELETE"},
	{"/pcdn/admin/alarm/record/list", "GET"},
	{"/pcdn/admin/bill/generate", "POST"},
	{"/pcdn/admin/bill/list", "GET"},
	{"/pcdn/admin/bill/find", "GET"},
	{"/pcdn/admin/bill/approve", "PUT"},
	{"/pcdn/admin/bill/reject", "PUT"},
	{"/pcdn/admin/bill/pay", "PUT"},
	{"/pcdn/admin/settlement/import", "POST"},
	{"/pcdn/admin/settlement/recheck", "PUT"},
	{"/pcdn/admin/settlement/list", "GET"},
	{"/pcdn/admin/settlement/revenue", "GET"},
	{"/pcdn/admin/settlement/delete", "DELETE"},
	{"/pcdn/admin/profit/summary", "GET"},
	{"/pcdn/admin/profit/revenueByPlatform", "GET"},
	{"/pcdn/admin/profit/costByOwner", "GET"},
	{"/pcdn/admin/profit/trend", "GET"},
	{"/pcdn/admin/release/list", "GET"},
	{"/pcdn/admin/release/find", "GET"},
	{"/pcdn/admin/release/create", "POST"},
	{"/pcdn/admin/release/update", "PUT"},
	{"/pcdn/admin/release/delete", "DELETE"},
}

// Auth 把本插件注册的菜单与 API 授权给超管角色 (888)。
//   - 菜单：直接增量 INSERT sys_authority_menus（不调 SetMenuAuthority，因其 Replace 会清空 admin 已有菜单）
//   - API：用 CasbinServiceApp.AddPolicies 增量写 casbin_rule，再 FreshCasbin 让 enforcer 重载（不调 UpdateCasbin，因其 ClearCasbin 会清空 admin 已有策略）
//
// 必须在 Menu(ctx) 之后调用（依赖菜单已落库）。
func Auth(ctx context.Context) {
	if global.GVA_DB == nil {
		zap.L().Warn("pcdn.initialize.Auth: DB 未就绪，跳过授权")
		return
	}

	authIDStr := strconv.Itoa(int(AdminAuthorityID))

	// 1. 菜单 → 角色
	var menus []model.SysBaseMenu
	if err := global.GVA_DB.WithContext(ctx).Where("name IN ?", pcdnMenuNames).Find(&menus).Error; err != nil {
		zap.L().Error("pcdn.initialize.Auth: 查询菜单失败", zap.Error(err))
		return
	}
	if len(menus) == 0 {
		zap.L().Warn("pcdn.initialize.Auth: 未查到 pcdn 菜单，确认 Menu() 是否先执行")
		return
	}
	var existing model.SysAuthorityMenu
	var toCreate []model.SysAuthorityMenu
	for _, m := range menus {
		err := global.GVA_DB.WithContext(ctx).
			Where("sys_authority_authority_id = ? AND sys_base_menu_id = ?", authIDStr, strconv.Itoa(int(m.ID))).
			Take(&existing).Error
		if err == nil {
			continue
		}
		toCreate = append(toCreate, model.SysAuthorityMenu{
			MenuId:      strconv.Itoa(int(m.ID)),
			AuthorityId: authIDStr,
		})
	}
	if len(toCreate) > 0 {
		if err := global.GVA_DB.WithContext(ctx).Create(&toCreate).Error; err != nil {
			zap.L().Error("pcdn.initialize.Auth: 写入 sys_authority_menus 失败", zap.Error(err))
		}
	}

	// 2. API → Casbin 规则
	rules := make([][]string, 0, len(pcdnApis))
	for _, a := range pcdnApis {
		var cnt int64
		global.GVA_DB.WithContext(ctx).Model(&adapter.CasbinRule{}).
			Where("ptype = 'p' AND v0 = ? AND v1 = ? AND v2 = ?", authIDStr, a.Path, a.Method).
			Count(&cnt)
		if cnt > 0 {
			continue
		}
		rules = append(rules, []string{authIDStr, a.Path, a.Method})
	}
	if len(rules) > 0 {
		if err := system.CasbinServiceApp.AddPolicies(global.GVA_DB, rules); err != nil {
			zap.L().Error("pcdn.initialize.Auth: 写入 casbin_rule 失败", zap.Error(err))
		}
		if err := system.CasbinServiceApp.FreshCasbin(); err != nil {
			zap.L().Error("pcdn.initialize.Auth: FreshCasbin 失败", zap.Error(err))
		}
	}
}

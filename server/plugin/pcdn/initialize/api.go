package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

const apiGroupPcdnNode = "PCDN节点"

// Api 注册 PCDN admin 路由的 Casbin 权限点（portal/agent 路由不走 Casbin，不在此注册）
func Api(ctx context.Context) {
	entities := []system.SysApi{
		{Path: "/pcdn/admin/node/list", Description: "查询节点列表", ApiGroup: apiGroupPcdnNode, Method: "GET"},
		{Path: "/pcdn/admin/node/find", Description: "查询节点详情", ApiGroup: apiGroupPcdnNode, Method: "GET"},
		{Path: "/pcdn/admin/node/create", Description: "创建节点", ApiGroup: apiGroupPcdnNode, Method: "POST"},
		{Path: "/pcdn/admin/node/update", Description: "更新节点", ApiGroup: apiGroupPcdnNode, Method: "PUT"},
		{Path: "/pcdn/admin/node/delete", Description: "删除节点", ApiGroup: apiGroupPcdnNode, Method: "DELETE"},
		{Path: "/pcdn/admin/node/deleteByIds", Description: "批量删除节点", ApiGroup: apiGroupPcdnNode, Method: "DELETE"},
		{Path: "/pcdn/admin/node/traffic", Description: "查询节点流量", ApiGroup: apiGroupPcdnNode, Method: "GET"},
		{Path: "/pcdn/admin/node/n95", Description: "查询节点95值", ApiGroup: apiGroupPcdnNode, Method: "GET"},
		{Path: "/pcdn/admin/alarm/rule/list", Description: "告警规则列表", ApiGroup: "PCDN告警", Method: "GET"},
		{Path: "/pcdn/admin/alarm/rule/find", Description: "查询告警规则", ApiGroup: "PCDN告警", Method: "GET"},
		{Path: "/pcdn/admin/alarm/rule/create", Description: "创建告警规则", ApiGroup: "PCDN告警", Method: "POST"},
		{Path: "/pcdn/admin/alarm/rule/update", Description: "更新告警规则", ApiGroup: "PCDN告警", Method: "PUT"},
		{Path: "/pcdn/admin/alarm/rule/delete", Description: "删除告警规则", ApiGroup: "PCDN告警", Method: "DELETE"},
		{Path: "/pcdn/admin/alarm/rule/deleteByIds", Description: "批量删除告警规则", ApiGroup: "PCDN告警", Method: "DELETE"},
		{Path: "/pcdn/admin/alarm/record/list", Description: "告警记录列表", ApiGroup: "PCDN告警", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}

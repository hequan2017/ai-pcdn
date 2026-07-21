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
		{Path: "/pcdn/admin/bill/generate", Description: "生成账单", ApiGroup: "PCDN账单", Method: "POST"},
		{Path: "/pcdn/admin/bill/list", Description: "账单列表", ApiGroup: "PCDN账单", Method: "GET"},
		{Path: "/pcdn/admin/bill/find", Description: "账单详情", ApiGroup: "PCDN账单", Method: "GET"},
		{Path: "/pcdn/admin/bill/approve", Description: "审核账单", ApiGroup: "PCDN账单", Method: "PUT"},
		{Path: "/pcdn/admin/bill/reject", Description: "驳回账单", ApiGroup: "PCDN账单", Method: "PUT"},
		{Path: "/pcdn/admin/bill/pay", Description: "账单付款", ApiGroup: "PCDN账单", Method: "PUT"},
		{Path: "/pcdn/admin/settlement/import", Description: "导入结算单", ApiGroup: "PCDN对账", Method: "POST"},
		{Path: "/pcdn/admin/settlement/recheck", Description: "重新核对", ApiGroup: "PCDN对账", Method: "PUT"},
		{Path: "/pcdn/admin/settlement/list", Description: "结算单列表", ApiGroup: "PCDN对账", Method: "GET"},
		{Path: "/pcdn/admin/settlement/revenue", Description: "应收汇总", ApiGroup: "PCDN对账", Method: "GET"},
		{Path: "/pcdn/admin/settlement/delete", Description: "删除结算单", ApiGroup: "PCDN对账", Method: "DELETE"},
		{Path: "/pcdn/admin/profit/summary", Description: "利润汇总", ApiGroup: "PCDN利润", Method: "GET"},
		{Path: "/pcdn/admin/profit/revenueByPlatform", Description: "按平台收入", ApiGroup: "PCDN利润", Method: "GET"},
		{Path: "/pcdn/admin/profit/costByOwner", Description: "按贡献者成本", ApiGroup: "PCDN利润", Method: "GET"},
		{Path: "/pcdn/admin/profit/trend", Description: "利润趋势", ApiGroup: "PCDN利润", Method: "GET"},
		{Path: "/pcdn/admin/release/list", Description: "版本列表", ApiGroup: "PCDN版本", Method: "GET"},
		{Path: "/pcdn/admin/release/find", Description: "查询版本", ApiGroup: "PCDN版本", Method: "GET"},
		{Path: "/pcdn/admin/release/create", Description: "发布版本", ApiGroup: "PCDN版本", Method: "POST"},
		{Path: "/pcdn/admin/release/update", Description: "更新版本", ApiGroup: "PCDN版本", Method: "PUT"},
		{Path: "/pcdn/admin/release/delete", Description: "删除版本", ApiGroup: "PCDN版本", Method: "DELETE"},
	}
	utils.RegisterApis(entities...)
}

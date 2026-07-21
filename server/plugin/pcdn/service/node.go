package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
)

// NodeService 节点管理服务
type NodeService struct{}

// CreateNode 创建节点（运营代录或自助上机预生成凭证）
func (s *NodeService) CreateNode(ctx context.Context, node *model.PcdnNode) error {
	return global.GVA_DB.WithContext(ctx).Create(node).Error
}

// DeleteNode 删除节点
func (s *NodeService) DeleteNode(ctx context.Context, id uint) error {
	return global.GVA_DB.WithContext(ctx).Delete(&model.PcdnNode{}, "id = ?", id).Error
}

// DeleteNodeByIds 批量删除节点
func (s *NodeService) DeleteNodeByIds(ctx context.Context, ids []uint) error {
	return global.GVA_DB.WithContext(ctx).Delete(&[]model.PcdnNode{}, "id in ?", ids).Error
}

// UpdateNode 更新节点（用 Select 显式列，避免零值 MonthlyPrice=0/P95UnitPrice=0 被跳过；保护 node_sn/token_hash/dept_id/created_by）
func (s *NodeService) UpdateNode(ctx context.Context, node model.PcdnNode) error {
	return global.GVA_DB.WithContext(ctx).
		Model(&model.PcdnNode{}).
		Where("id = ?", node.ID).
		Select("owner_name", "contact", "region", "isp", "platform", "platform_node_id", "group_id", "tags", "status", "billing_mode", "monthly_price", "p95_unit_price", "contract_period", "remark").
		Updates(&node).Error
}

// GetNode 根据 ID 查询节点
func (s *NodeService) GetNode(ctx context.Context, id uint) (node model.PcdnNode, err error) {
	err = global.GVA_DB.WithContext(ctx).Where("id = ?", id).First(&node).Error
	return
}

// GetNodeByNodeSn 根据 node_sn 查询节点（agent 鉴权/激活用）
func (s *NodeService) GetNodeByNodeSn(ctx context.Context, nodeSn string) (node model.PcdnNode, err error) {
	err = global.GVA_DB.WithContext(ctx).Where("node_sn = ?", nodeSn).First(&node).Error
	return
}

// GetNodeList 分页查询节点
func (s *NodeService) GetNodeList(ctx context.Context, info request.NodeSearch) (list []model.PcdnNode, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{})
	if info.NodeSn != "" {
		db = db.Where("node_sn = ?", info.NodeSn)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.OwnerName != "" {
		db = db.Where("owner_name LIKE ?", "%"+info.OwnerName+"%")
	}
	if info.Platform != "" {
		db = db.Where("platform = ?", info.Platform)
	}
	if info.Region != "" {
		db = db.Where("region = ?", info.Region)
	}
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
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

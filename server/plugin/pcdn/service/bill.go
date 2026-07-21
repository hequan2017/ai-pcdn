package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
	"gorm.io/gorm"
)

// BillService 采购侧账单服务
type BillService struct{}

// GenerateMonthlyBill 生成某账期账单：按 owner 分组，monthly 取包月价、p95 取月95 frozen×单价。幂等（同 period+owner 已存在则跳过）。
func (s *BillService) GenerateMonthlyBill(ctx context.Context, period string) (int, error) {
	start, _, err := parsePeriod(period)
	if err != nil {
		return 0, err
	}
	var nodes []model.PcdnNode
	if err := global.GVA_DB.WithContext(ctx).Where("status != ? AND owner_user_id > 0", model.NodeStatusDisabled).Find(&nodes).Error; err != nil {
		return 0, err
	}
	// 按 owner 分组
	groups := map[uint][]model.PcdnNode{}
	for i := range nodes {
		groups[nodes[i].OwnerUserID] = append(groups[nodes[i].OwnerUserID], nodes[i])
	}

	created := 0
	for ownerID, groupNodes := range groups {
		// 幂等：已存在则跳过
		var exist model.PcdnBill
		if err := global.GVA_DB.WithContext(ctx).Where("period = ? AND owner_user_id = ?", period, ownerID).First(&exist).Error; err == nil {
			continue
		}
		// owner 联系方式
		ownerName, contact := ownerInfo(ctx, ownerID, groupNodes)

		details := make([]model.BillDetail, 0, len(groupNodes))
		var total float64
		for i := range groupNodes {
			n := groupNodes[i]
			amount, value, unitPrice := s.calcNodeAmount(ctx, n, start)
			details = append(details, model.BillDetail{
				NodeID: n.ID, NodeSn: n.NodeSn, BillingMode: n.BillingMode,
				Value: value, UnitPrice: unitPrice, Amount: amount,
			})
			total += amount
		}
		detailJSON, _ := json.Marshal(details)
		bill := model.PcdnBill{
			Period:       period,
			OwnerUserID:  ownerID,
			OwnerName:    ownerName,
			OwnerContact: contact,
			NodeCount:    len(groupNodes),
			Details:      detailJSON,
			TotalAmount:  total,
			Status:       model.BillStatusDraft,
		}
		if err := global.GVA_DB.WithContext(ctx).Create(&bill).Error; err == nil {
			created++
		}
	}
	return created, nil
}

// calcNodeAmount 计算单节点账期金额：monthly=包月价；p95=月95(Mbps)×单价
func (s *BillService) calcNodeAmount(ctx context.Context, node model.PcdnNode, periodStart time.Time) (amount, value, unitPrice float64) {
	if node.BillingMode == model.BillingModeMonthly {
		return node.MonthlyPrice, node.MonthlyPrice, 1
	}
	// p95：取该账期月95 frozen
	var n95 model.PcdnNode95
	_ = global.GVA_DB.WithContext(ctx).
		Where("node_id = ? AND period_type = ? AND period_start = ?", node.ID, model.PeriodTypeMonth, periodStart).
		First(&n95).Error
	mbps := float64(n95.Combined95Bps) / 1e6
	return mbps * node.P95UnitPrice, mbps, node.P95UnitPrice
}

// ApproveBill 审核通过（draft → approved）
func (s *BillService) ApproveBill(ctx context.Context, id, auditedBy uint) error {
	now := time.Now()
	return global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{}).
		Where("id = ? AND status = ?", id, model.BillStatusDraft).
		Updates(map[string]interface{}{"status": model.BillStatusApproved, "audited_by": auditedBy, "audited_at": now}).Error
}

// RejectBill 驳回（draft → rejected）
func (s *BillService) RejectBill(ctx context.Context, id, auditedBy uint, remark string) error {
	now := time.Now()
	return global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{}).
		Where("id = ? AND status = ?", id, model.BillStatusDraft).
		Updates(map[string]interface{}{"status": model.BillStatusRejected, "audited_by": auditedBy, "audited_at": now, "remark": remark}).Error
}

// PayBill 付款（approved → paid），事务内校验状态
func (s *BillService) PayBill(ctx context.Context, req request.PayBillReq, paidBy uint) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var bill model.PcdnBill
		if err := tx.Where("id = ? AND status = ?", req.ID, model.BillStatusApproved).First(&bill).Error; err != nil {
			return fmt.Errorf("账单不存在或非已审核状态")
		}
		now := time.Now()
		return tx.Model(&bill).Updates(map[string]interface{}{
			"status":      model.BillStatusPaid,
			"paid_amount": req.PaidAmount,
			"pay_method":  req.PayMethod,
			"pay_no":      req.PayNo,
			"paid_at":     now,
			"paid_by":     paidBy,
			"remark":      req.Remark,
		}).Error
	})
}

// GetBill 查询账单详情
func (s *BillService) GetBill(ctx context.Context, id uint) (bill model.PcdnBill, err error) {
	err = global.GVA_DB.WithContext(ctx).Where("id = ?", id).First(&bill).Error
	return
}

// GetBillList 分页查询账单
func (s *BillService) GetBillList(ctx context.Context, info request.BillSearch) (list []model.PcdnBill, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{})
	if info.Period != "" {
		db = db.Where("period = ?", info.Period)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.OwnerName != "" {
		db = db.Where("owner_name LIKE ?", "%"+info.OwnerName+"%")
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

// GetBillsByOwner 个人门户：查自己的账单
func (s *BillService) GetBillsByOwner(ctx context.Context, ownerUserID uint, info request.BillSearch) (list []model.PcdnBill, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnBill{}).Where("owner_user_id = ?", ownerUserID)
	if info.Period != "" {
		db = db.Where("period = ?", info.Period)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
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

// parsePeriod "2026-07" → (月初, 下月初)
func parsePeriod(period string) (time.Time, time.Time, error) {
	parts := strings.Split(period, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("账期格式应为 YYYY-MM")
	}
	y, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || m < 1 || m > 12 {
		return time.Time{}, time.Time{}, fmt.Errorf("账期非法")
	}
	start := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	return start, end, nil
}

// ownerInfo 取 owner 名称与联系方式（优先节点上的 owner_name/contact，回退查 sys_user）
func ownerInfo(ctx context.Context, ownerUserID uint, nodes []model.PcdnNode) (string, string) {
	name := ""
	if len(nodes) > 0 {
		name = nodes[0].OwnerName
	}
	var u sysModel.SysUser
	global.GVA_DB.WithContext(ctx).Select("nick_name, username, phone, email").First(&u, ownerUserID)
	if name == "" {
		name = u.NickName
		if name == "" {
			name = u.Username
		}
	}
	contact := ""
	if len(nodes) > 0 {
		contact = nodes[0].Contact
	}
	if contact == "" {
		contact = u.Phone
		if contact == "" {
			contact = u.Email
		}
	}
	return name, contact
}

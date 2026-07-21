package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	sysReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	sysService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model/request"
)

// ContributorAuthorityID 个人贡献者角色ID（复用 GVA 默认普通用户角色；生产环境应初始化专门角色）
const ContributorAuthorityID uint = 888

// PortalService 个人门户服务（注册/登录/自助上机/我的节点）
type PortalService struct{}

// Register 个人贡献者注册
func (s *PortalService) Register(ctx context.Context, req request.PortalRegister) (*sysModel.SysUser, error) {
	user := sysModel.SysUser{
		Username:    req.Username,
		Password:    req.Password,
		NickName:    req.NickName,
		Phone:       req.Phone,
		Email:       req.Email,
		AuthorityId: ContributorAuthorityID,
		Enable:      1,
	}
	u, err := sysService.UserServiceApp.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Login 登录校验并签发 JWT
func (s *PortalService) Login(ctx context.Context, req request.PortalLogin) (user *sysModel.SysUser, token string, err error) {
	input := &sysModel.SysUser{Username: req.Username, Password: req.Password}
	user, err = sysService.UserServiceApp.Login(ctx, input)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}
	if user.Enable == 2 {
		return nil, "", fmt.Errorf("账号已被冻结")
	}
	claims := utils.NewJWT().CreateClaims(sysReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		Username:    user.Username,
		NickName:    user.NickName,
		AuthorityId: user.AuthorityId,
	})
	token, err = utils.NewJWT().CreateToken(claims)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// AddNode 自助上机：预创建 pending 节点并生成凭证。
// 明文 token 仅在此处生成并返回一次，数据库只存 sha256 哈希（哈希不可逆，故必须在生成时保留明文）。
func (s *PortalService) AddNode(ctx context.Context, ownerUserID uint, req request.PortalAddNode) (*model.PcdnNode, string, error) {
	var owner sysModel.SysUser
	global.GVA_DB.WithContext(ctx).Select("nick_name, username").First(&owner, ownerUserID)
	ownerName := owner.NickName
	if ownerName == "" {
		ownerName = owner.Username
	}
	plain := genNodeToken()
	node := model.PcdnNode{
		NodeSn:      genNodeSn(),
		TokenHash:   sha256Hex(plain),
		OwnerUserID: ownerUserID,
		OwnerName:   ownerName,
		Region:      req.Region,
		Isp:         req.Isp,
		Platform:    req.Platform,
		Status:      model.NodeStatusPending,
	}
	if err := global.GVA_DB.WithContext(ctx).Create(&node).Error; err != nil {
		return nil, "", err
	}
	return &node, plain, nil
}

// MyNodes 查询个人名下节点
func (s *PortalService) MyNodes(ctx context.Context, ownerUserID uint, info request.PortalNodeSearch) (list []model.PcdnNode, total int64, err error) {
	limit, offset := info.LimitOffset()
	db := global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{}).Where("owner_user_id = ?", ownerUserID)
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.NodeSn != "" {
		db = db.Where("node_sn LIKE ?", "%"+info.NodeSn+"%")
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

// CheckNodeOwner 校验节点归属（portal 查流量前强制校验）
func (s *PortalService) CheckNodeOwner(ctx context.Context, nodeID, ownerUserID uint) (bool, error) {
	var cnt int64
	err := global.GVA_DB.WithContext(ctx).Model(&model.PcdnNode{}).
		Where("id = ? AND owner_user_id = ?", nodeID, ownerUserID).Count(&cnt).Error
	return cnt > 0, err
}

func genNodeSn() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "PCDN-" + hex.EncodeToString(b)
}

func genNodeToken() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

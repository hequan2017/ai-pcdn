package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// PortalRegister 个人贡献者注册
type PortalRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// PortalLogin 个人登录
type PortalLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PortalAddNode 自助上机：添加节点（生成凭证）
type PortalAddNode struct {
	Region   string `json:"region"`
	Isp      string `json:"isp"`
	Platform string `json:"platform"`
}

// PortalNodeSearch 个人门户节点查询
type PortalNodeSearch struct {
	Status string `json:"status" form:"status"`
	NodeSn string `json:"nodeSn" form:"nodeSn"`
	commonReq.PageInfo
}

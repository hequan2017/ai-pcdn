package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	pcdnModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/model"
	pcdnService "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/service"
	"github.com/gin-gonic/gin"
)

// agent 鉴权上下文键
const (
	CtxKeyPcdnNodeID = "pcdn_node_id"
	CtxKeyPcdnNodeSn = "pcdn_node_sn"
)

// AgentTokenAuth 校验采集 agent 上报的节点凭证（X-Node-Sn + X-Node-Token），不走 JWT/Casbin/DataScope。
// token 明文仅生成时返回给用户一次，数据库存 sha256 哈希；此处按 node_sn 查节点后用常数时间比对哈希。
// pending 节点只允许 /activate，拒绝 /report 与 /heartbeat（防止绕过激活流程）。
func AgentTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodeSn := c.GetHeader("X-Node-Sn")
		token := c.GetHeader("X-Node-Token")
		if nodeSn == "" || token == "" {
			response.NoAuth("missing node credentials", c)
			c.Abort()
			return
		}
		node, err := pcdnService.ServiceGroupApp.NodeService.GetNodeByNodeSn(c.Request.Context(), nodeSn)
		if err != nil {
			response.NoAuth("invalid node", c)
			c.Abort()
			return
		}
		if node.TokenHash == "" || node.Status == pcdnModel.NodeStatusDisabled {
			response.NoAuth("node not activated or disabled", c)
			c.Abort()
			return
		}
		sum := sha256.Sum256([]byte(token))
		// 常数时间比较，避免时序攻击逐字节恢复哈希
		if subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(node.TokenHash)) != 1 {
			response.NoAuth("invalid node token", c)
			c.Abort()
			return
		}
		// pending 节点必须先激活，禁止直接上报/心跳
		if node.Status == pcdnModel.NodeStatusPending && !strings.HasSuffix(c.Request.URL.Path, "/activate") {
			response.FailWithMessage("节点未激活，请先调用 activate 接口", c)
			c.Abort()
			return
		}
		c.Set(CtxKeyPcdnNodeID, node.ID)
		c.Set(CtxKeyPcdnNodeSn, node.NodeSn)
		c.Next()
	}
}

// GetNodeIDFromCtx 从上下文取已鉴权的节点ID
func GetNodeIDFromCtx(c *gin.Context) uint {
	if v, ok := c.Get(CtxKeyPcdnNodeID); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

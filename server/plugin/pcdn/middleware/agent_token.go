package middleware

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	pcdnService "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/service"
	"github.com/gin-gonic/gin"
)

// agent 鉴权上下文键
const (
	CtxKeyPcdnNodeID = "pcdn_node_id"
	CtxKeyPcdnNodeSn = "pcdn_node_sn"
)

// AgentTokenAuth 校验采集 agent 上报的节点凭证（X-Node-Sn + X-Node-Token），不走 JWT/Casbin/DataScope。
// token 明文仅生成时返回给用户一次，数据库存 sha256 哈希；此处按 node_sn 查节点后比对哈希。
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
		if node.TokenHash == "" || node.Status == "disabled" {
			response.NoAuth("node not activated or disabled", c)
			c.Abort()
			return
		}
		sum := sha256.Sum256([]byte(token))
		if hex.EncodeToString(sum[:]) != node.TokenHash {
			response.NoAuth("invalid node token", c)
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

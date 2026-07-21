package router

import "github.com/gin-gonic/gin"

// BillRouter 账单路由（admin 组）
type BillRouter struct{}

// InitBillRouter 挂载账单路由
func (r *BillRouter) InitBillRouter(group *gin.RouterGroup) {
	group.POST("bill/generate", billApi.GenerateBill)
	group.GET("bill/list", billApi.GetBillList)
	group.GET("bill/find", billApi.GetBill)
	group.PUT("bill/approve", billApi.ApproveBill)
	group.PUT("bill/reject", billApi.RejectBill)
	group.PUT("bill/pay", billApi.PayBill)
}

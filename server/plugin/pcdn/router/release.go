package router

import "github.com/gin-gonic/gin"

// ReleaseRouter agent 版本发布路由（admin 组）
type ReleaseRouter struct{}

// InitReleaseRouter 挂载版本路由
func (r *ReleaseRouter) InitReleaseRouter(group *gin.RouterGroup) {
	group.GET("release/list", releaseApi.GetReleaseList)
	group.GET("release/find", releaseApi.GetRelease)
	group.POST("release/create", releaseApi.CreateRelease)
	group.PUT("release/update", releaseApi.UpdateRelease)
	group.DELETE("release/delete", releaseApi.DeleteRelease)
}

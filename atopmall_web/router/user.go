package router

import (
	"atopmall_web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	UserRouter := r.Group("user")
	UserRouter.GET("list", api.GerUserList)
}

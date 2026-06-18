package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/user_web/api"
)

func InitUserRouter(r *gin.RouterGroup) {
	UserRouter := r.Group("user")
	UserRouter.GET("list", api.GerUserList)
}

package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/user_web/api"
	"atopmall_web/user_web/middlewares"
)

// 用户相关路由
func InitUserRouter(r *gin.RouterGroup) {
	//UserRouter := r.Group("user").Use(middlewares.JWTAuth(), middlewares.IsAdmin()) //这里的顺序不能改变，先验证JWT这个登录状态，再验证是否是管理员
	UserRouter := r.Group("user")
	UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdmin(), api.GetUserList)
	UserRouter.POST("pwd_login", api.PasswordLogin)
}

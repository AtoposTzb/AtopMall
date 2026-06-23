package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/user_web/api"
)

// 发送验证码的相关路由（注册、登录、重置密码，图片验证码）
func BaseRouteInit(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("send-code", api.SendCode)
	}
}

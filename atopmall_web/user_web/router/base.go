package router

import (
	"atopmall_web/user_web/api"

	"github.com/gin-gonic/gin"
)

func BaseRouteInit(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
	}
}

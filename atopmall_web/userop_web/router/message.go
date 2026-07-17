package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/userop_web/api/message"
	"atopmall_web/userop_web/middlewares"
)

func MessageRouterInit(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.GetMessageList) // 获取留言列表
		MessageRouter.POST("", message.NewMessage)    //新建留言
	}
}

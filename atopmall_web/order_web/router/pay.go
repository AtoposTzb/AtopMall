package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/order_web/api/pay"
)

func PayRouterInit(Router *gin.RouterGroup) {
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify) // 支付宝通知
		PayRouter.GET("alipay/return", pay.Notify)  // 支付宝返回
	}
}

package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/order_web/api/order"
	"atopmall_web/order_web/middlewares"
)

func OrderRouterInit(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.OrderList)       // 订单列表
		OrderRouter.POST("", order.OrderCreate)    // 创建订单
		OrderRouter.GET("/:id", order.OrderDetail) // 获取订单详情

	}
}

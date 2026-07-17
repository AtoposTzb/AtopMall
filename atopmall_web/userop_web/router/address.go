package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/userop_web/api/address"
	"atopmall_web/userop_web/middlewares"
)

func AddressRouterInit(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("", middlewares.JWTAuth(), address.GetAddressList)       // 获取地址列表
		AddressRouter.DELETE("/:id", middlewares.JWTAuth(), address.DeleteAddress) // 删除地址
		AddressRouter.POST("", middlewares.JWTAuth(), address.NewAddress)          //新建地址
		AddressRouter.PUT("/:id", middlewares.JWTAuth(), address.UpdateAddress)    //修改地址信息
	}
}

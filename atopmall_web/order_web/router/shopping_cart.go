package router

import (
	"github.com/gin-gonic/gin"

	shoppingcart "atopmall_web/order_web/api/shopping_cart"
	"atopmall_web/order_web/middlewares"
)

func ShoppingCartRouterInit(Router *gin.RouterGroup) {
	ShoppingCartRouter := Router.Group("shoppingcart").Use(middlewares.JWTAuth())
	{
		ShoppingCartRouter.GET("", shoppingcart.ShoppingCartList)          //获取购物车列表
		ShoppingCartRouter.DELETE("/:id", shoppingcart.DeleteShoppingCart) //删除购物车条目
		ShoppingCartRouter.POST("", shoppingcart.NewShoppingCart)          //添加购物车条目
		ShoppingCartRouter.PATCH("/:id", shoppingcart.UpdateShoppingCart)  //修改购物车条目 put和patch区别：put会覆盖所有字段，patch不会覆盖所有字段，只更新修改的字段
	}
}

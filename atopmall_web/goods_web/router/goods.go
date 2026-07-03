package router

import (
	"atopmall_web/goods_web/api/goods"
	"atopmall_web/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

// 商品相关路由
func GoodsRouterInit(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.GetGoodsList)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.NewGoods)
		GoodsRouter.GET("/:id", goods.GetGoodsDetail)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.DeleteGoods)
		GoodsRouter.GET("/:id/stocks", goods.Stocks) //获取商品的库存

		//更新商品中的是否热门等信息
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateGoodsStatus)
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateGoods)
	}
}

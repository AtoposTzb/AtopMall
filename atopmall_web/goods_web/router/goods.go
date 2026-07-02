package router

import (
	"atopmall_web/goods_web/api/goods"

	"github.com/gin-gonic/gin"
)

// 商品相关路由
func InitGoodsRouter(r *gin.RouterGroup) {
	//UserRouter := r.Group("user").Use(middlewares.JWTAuth(), middlewares.IsAdmin()) //这里的顺序不能改变，先验证JWT这个登录状态，再验证是否是管理员
	GoodsRouter := r.Group("goods")
	GoodsRouter.GET("", goods.GetGoodsList)
	// GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.NewGoods)
	GoodsRouter.POST("", goods.NewGoods)

	GoodsRouter.GET("/:id", goods.GetGoodsDetail)
	// GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.DeleteGoods)
	GoodsRouter.DELETE("/:id", goods.DeleteGoods)

	GoodsRouter.GET("/:id/stocks", goods.Stocks) //获取商品的库存

	//更新商品中的是否热门等信息
	// GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateGoodsStatus)
	// GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateGoods)
	GoodsRouter.PATCH("/:id", goods.UpdateGoodsStatus)
	GoodsRouter.PUT("/:id", goods.UpdateGoods)
}

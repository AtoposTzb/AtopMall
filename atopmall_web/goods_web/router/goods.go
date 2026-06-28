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
}

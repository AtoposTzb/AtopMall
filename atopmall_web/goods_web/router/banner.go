package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/goods_web/api/banners"
	"atopmall_web/goods_web/middlewares"
)

func BannerRouterInit(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banners.GetBannerList)                                                     // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), banners.DeleteBanner) // 删除轮播图
		BannerRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdmin(), banners.NewBanner)          //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), banners.UpdateBanner)    //修改轮播图信息

	}
}

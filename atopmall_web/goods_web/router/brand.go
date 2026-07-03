package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/goods_web/api/brands"
)

// 1. 商品的api接口开发完成
// 2. 图片的坑
func BrandRouterInit(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brands")
	{
		BrandRouter.GET("", brands.BrandList)          // 品牌列表页
		BrandRouter.DELETE("/:id", brands.DeleteBrand) // 删除品牌
		BrandRouter.POST("", brands.NewBrand)          //新建品牌
		BrandRouter.PUT("/:id", brands.UpdateBrand)    //修改品牌信息
	}
}

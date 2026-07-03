package router

import (
	"github.com/gin-gonic/gin"

	categorybrand "atopmall_web/goods_web/api/category_brand"
)

// 1. 商品的api接口开发完成
// 2. 图片的坑
func CategoryBrandRouterInit(Router *gin.RouterGroup) {
	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", categorybrand.CategoryBrandList)          // 类别品牌列表页
		CategoryBrandRouter.DELETE("/:id", categorybrand.DeleteCategoryBrand) // 删除类别品牌
		CategoryBrandRouter.POST("", categorybrand.NewCategoryBrand)          //新建类别品牌
		CategoryBrandRouter.PUT("/:id", categorybrand.UpdateCategoryBrand)    //修改类别品牌
		CategoryBrandRouter.GET("/:id", categorybrand.GetCategoryBrandList)   //获取分类的品牌
	}
}

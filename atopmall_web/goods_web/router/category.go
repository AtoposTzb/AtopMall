package router

import (
	"atopmall_web/goods_web/api/category"

	"github.com/gin-gonic/gin"
)

func CategoryRouterInit(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys")
	{
		CategoryRouter.GET("", category.GetCategoryList)       // 商品类别列表页
		CategoryRouter.DELETE("/:id", category.DeleteCategory) // 删除分类
		CategoryRouter.GET("/:id", category.GetCategoryDetail) // 获取分类详情
		CategoryRouter.POST("", category.NewCategory)          //新建分类
		CategoryRouter.PUT("/:id", category.UpdateCategory)    //修改分类信息
	}
}

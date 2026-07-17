package router

import (
	"github.com/gin-gonic/gin"

	userfav "atopmall_web/userop_web/api/user_fav"
	"atopmall_web/userop_web/middlewares"
)

func UserFavRouterInit(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfavs")
	{
		UserFavRouter.DELETE("/:id", middlewares.JWTAuth(), userfav.DeleteUserFav) // 删除收藏记录
		UserFavRouter.GET("/:id", middlewares.JWTAuth(), userfav.GetUserFavDetail) // 获取收藏记录
		UserFavRouter.POST("", middlewares.JWTAuth(), userfav.NewUserFav)          //新建收藏记录
		UserFavRouter.GET("", middlewares.JWTAuth(), userfav.GetUserFavList)       //获取当前用户的收藏列表
	}
}

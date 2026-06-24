package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"atopmall_web/user_web/middlewares"
	"atopmall_web/user_web/router"
)

func RoutersInit() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	//consul健康检查
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
		})
	})
	ApiRouter := r.Group("/u/v1")
	router.InitUserRouter(ApiRouter)
	router.BaseRouteInit(ApiRouter)
	return r
}

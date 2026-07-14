package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"atopmall_web/order_web/middlewares"
	"atopmall_web/order_web/router"
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

	ApiRouter := r.Group("/o/v1")
	{
		router.OrderRouterInit(ApiRouter)
		router.ShoppingCartRouterInit(ApiRouter)
	}

	return r
}

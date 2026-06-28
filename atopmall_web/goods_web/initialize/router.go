package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"atopmall_web/goods_web/middlewares"
	"atopmall_web/goods_web/router"
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
	ApiRouter := r.Group("/g/v1")
	router.InitGoodsRouter(ApiRouter)
	return r
}

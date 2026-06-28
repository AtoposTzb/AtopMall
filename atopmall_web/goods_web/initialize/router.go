package initialize

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/goods_web/middlewares"
	"atopmall_web/goods_web/router"
)

func RoutersInit() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	ApiRouter := r.Group("/g/v1")
	router.InitGoodsRouter(ApiRouter)
	return r
}

package initialize

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/user_web/middlewares"
	"atopmall_web/user_web/router"
)

func RoutersInit() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	ApiRouter := r.Group("/u/v1")
	router.InitUserRouter(ApiRouter)
	router.BaseRouteInit(ApiRouter)
	return r
}

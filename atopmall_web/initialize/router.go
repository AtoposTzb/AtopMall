package initialize

import (
	"atopmall_web/router"

	"github.com/gin-gonic/gin"
)

func RoutersInit() *gin.Engine {
	r := gin.Default()
	ApiRouter := r.Group("u/v1")
	router.InitUserRouter(ApiRouter)
	return r
}

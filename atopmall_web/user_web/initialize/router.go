package initialize

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/user_web/router"
)

func RoutersInit() *gin.Engine {
	r := gin.Default()
	ApiRouter := r.Group("/u/v1")
	router.InitUserRouter(ApiRouter)
	return r
}

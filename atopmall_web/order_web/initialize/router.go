package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"atopmall_web/order_web/global"
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

	// ========== 新增：OpenTelemetry Gin 追踪中间件 ==========
	// 必须放在业务路由之前，才能拦截所有请求
	// 参数是服务名，和 JaegerTracerInit 中保持一致即可
	// 拦截器放在健康检查路由之后, 这样HTTP 健康检查不会产生大量无意义的 span

	r.Use(otelgin.Middleware(global.ServerConfig.JaegerInfo.Name))

	ApiRouter := r.Group("/o/v1")
	{
		router.OrderRouterInit(ApiRouter)
		router.ShoppingCartRouterInit(ApiRouter)
	}

	return r
}

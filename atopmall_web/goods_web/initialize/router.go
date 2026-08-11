package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"atopmall_web/goods_web/global"
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

	// ========== 新增：OpenTelemetry Gin 追踪中间件 ==========
	// 必须放在业务路由之前，才能拦截所有请求
	// 参数是服务名，和 JaegerTracerInit 中保持一致即可
	// 拦截器放在健康检查路由之后, 这样HTTP 健康检查不会产生大量无意义的 span

	r.Use(otelgin.Middleware(global.ServerConfig.JaegerInfo.Name))
	//kong配置为/g前缀，匹配请求:/g/v1/...
	ApiRouter := r.Group("/v1")
	{
		router.GoodsRouterInit(ApiRouter)
		router.CategoryRouterInit(ApiRouter)
		router.BrandRouterInit(ApiRouter)
		router.BannerRouterInit(ApiRouter)
		router.CategoryBrandRouterInit(ApiRouter)
	}

	return r
}

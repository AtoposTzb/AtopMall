package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"atopmall_web/oss_web/global"
	"atopmall_web/oss_web/middlewares"
	"atopmall_web/oss_web/router"
)

func RoutersInit() *gin.Engine {
	r := gin.Default()
	//配置跨域
	r.Use(middlewares.Cors())
	//consul健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// ========== 新增：OpenTelemetry 追踪中间件 ==========
	// 放在健康检查路由之后，避免健康检查产生大量无意义 Span
	// 参数 "" 会和 Nacos 配置中的 jaeger.name 保持一致

	r.Use(otelgin.Middleware(global.ServerConfig.JaegerInfo.Name))

	r.LoadHTMLGlob("templates/*")
	// 配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	r.StaticFS("/static", http.Dir("./static"))
	// GET：请求方式；/hello：请求的路径
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("/", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	ApiGroup := r.Group("/oss/v1")
	{
		router.OssRouterInit(ApiGroup)
	}

	return r
}

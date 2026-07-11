package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"atopmall_web/oss-web/middlewares"
	"atopmall_web/oss-web/router"
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

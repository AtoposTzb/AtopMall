package router

import (
	"github.com/gin-gonic/gin"

	"atopmall_web/oss_web/handler"
)

func OssRouterInit(Router *gin.RouterGroup) {
	OssRouter := Router.Group("oss")
	{
		// 前端获取MinIO直传签名接口
		OssRouter.GET("token", handler.Token)
		// 原阿里云回调接口，MinIO场景兼容占位
		OssRouter.POST("/callback", handler.HandlerRequest)
		// 清理孤儿文件（上传后超过指定时间未被业务服务认领的文件）
		OssRouter.DELETE("cleanup", handler.CleanupOrphanFiles)
	}
}

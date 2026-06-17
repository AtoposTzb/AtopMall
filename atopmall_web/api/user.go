package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GerUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页")
}

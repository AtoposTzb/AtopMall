package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"atopmall_web/userop_web/models"
)

func IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currrentUser := claims.(*models.CustomClaims)
		if currrentUser.AuthorityID != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "您没有权限访问",
			})
			ctx.Abort() // 中断后续处理
			return
		}
		ctx.Next()
	}
}

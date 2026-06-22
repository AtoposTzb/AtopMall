package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store) // 创建验证码 参数：驱动器、存储器 返回值：验证码对象
	id, b64s, _, err := cp.Generate()             // 生成验证码 返回值：验证码id、base64编码的验证码图片、验证码答案、错误信息
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成验证码失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":        id,
		"picBase64": b64s,
	})

}

package api

import (
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	// 限流
	e, b := sentinel.Entry("captcha", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求频率过快,请稍后重试",
		})
		return
	}
	defer func() {
		if e != nil {
			e.Exit()
		}
	}()
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

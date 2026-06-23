package api

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"

	"atopmall_web/user_web/forms"
	"atopmall_web/user_web/global"
)

// SendEmail 发送邮箱验证码
func SendEmail(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "atopmall注册验证码 <" + global.ServerConfig.EmailInfo.Username + ">"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送"
	e.HTML = []byte("您的注册验证码<b>" + code + "</b>" + "请在" + strconv.Itoa(global.ServerConfig.EmailInfo.Expires/60) + "分钟内输入")
	return e.SendWithTLS(global.ServerConfig.EmailInfo.Host+":"+strconv.Itoa(global.ServerConfig.EmailInfo.Port), smtp.PlainAuth("", global.ServerConfig.EmailInfo.Username, global.ServerConfig.EmailInfo.Password, global.ServerConfig.EmailInfo.Host),
		&tls.Config{InsecureSkipVerify: true, ServerName: global.ServerConfig.EmailInfo.Host})
}

// 随机验证码
func CreateCode() string {
	//生成6位随机数
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

// SendCode 发送邮箱验证码 ，验证码有效期5分钟
func SendCode(ctx *gin.Context) {
	//表单验证
	emailCodeForm := forms.EmailCodeForm{}
	if err := ctx.ShouldBind(&emailCodeForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	email := emailCodeForm.Email //获取邮箱
	if email == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "邮箱不能为空",
		})
		return
	}
	code := CreateCode() //随机生成验证码

	//连接redis 写入验证码
	err := global.RDB.Set(ctx, email, code, time.Second*time.Duration(global.ServerConfig.EmailInfo.Expires)).Err() //验证码有效期expires/60分钟，time.Duration()将int转换为time.Duration类型
	if err != nil {
		log.Println(err) // 写入redis失败，记录日志，用户看不到
		return
	}
	//发送验证码
	err = SendEmail(email, code)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "验证码发送失败:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"code": "验证码发送成功:",
		},
	})
}

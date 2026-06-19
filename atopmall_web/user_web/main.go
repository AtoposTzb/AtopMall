package main

import (
	"strconv"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"atopmall_web/user_web/global"
	"atopmall_web/user_web/initialize"
	myValidator "atopmall_web/user_web/validator"
)

func main() {
	//1.初始化logger
	initialize.LoggerInit()
	//2.初始化配置
	initialize.ConfigInit()
	//3.初始化路由
	r := initialize.RoutersInit()
	//4.初始化翻译器
	initialize.TransInit("zh")

	//5.注册验证器,自定义验证器，配置到form标签中
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myValidator.ValidateMobile)
		//自定义翻译器，参考官网注册翻译器
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)

	// 启动服务器
	if err := r.Run(":" + strconv.Itoa(global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务器失败：", zap.Error(err))
	}

}

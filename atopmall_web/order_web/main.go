package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
	"go.uber.org/zap"

	"atopmall_web/order_web/global"
	"atopmall_web/order_web/initialize"
	"atopmall_web/order_web/utils"
	"atopmall_web/order_web/utils/register/consul"
	myValidator "atopmall_web/order_web/validator"
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
	//5.初始化订单服务的grpc客户端连接
	initialize.SrcClientInitBL() //带负载均衡策略的连接
	//6.动态获取端口号,方便apifox调试使用固定8084
	if debug := initialize.GetEnvInfo(global.Env); !debug {
		orderPort, err := utils.GetAddrPort()
		if err != nil {
			panic(err)
		}
		global.ServerConfig.Port = orderPort
	}

	//5.1注册验证器,自定义验证器，配置到form标签中
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
	//7.初始化consul注册中心
	registryClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	//8.注册服务
	serviceId, _ := uuid.NewV4()
	if err := registryClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId.String()); err != nil {
		zap.S().Panic("注册服务失败：", zap.Error(err))
	}

	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)
	// 启动服务器
	//处理优雅的退出信号
	go func() {
		if err := r.Run(":" + strconv.Itoa(global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动服务器失败：", zap.Error(err))
		}
	}()
	//10.接收退出信号
	quit := make(chan os.Signal, 1)                      //定义一个信号通道，用于接收退出信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //监听SIGINT和SIGTERM信号
	<-quit                                               //等待信号通道接收信号
	if err := registryClient.Deregister(serviceId.String()); err != nil {
		zap.S().Info("注销失败", zap.Error(err))
	} else {
		zap.S().Info("注销成功")
	}

}

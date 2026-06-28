package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
	"go.uber.org/zap"

	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/initialize"
	"atopmall_web/goods_web/utils"
	"atopmall_web/goods_web/utils/register/consul"
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
	//5.初始化商品服务的grpc客户端连接
	initialize.GoodsSrcClientInitBL() //带负载均衡策略的连接
	//6.动态获取端口号,方便apifox调试使用固定8082
	if debug := initialize.GetEnvInfo(global.Env); !debug {
		goodsPort, err := utils.GetAddrPort()
		if err != nil {
			panic(err)
		}
		global.ServerConfig.Port = goodsPort
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
	quit := make(chan os.Signal)                         //定义一个信号通道，用于接收退出信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //监听SIGINT和SIGTERM信号
	<-quit                                               //等待信号通道接收信号
	if err := registryClient.Deregister(serviceId.String()); err != nil {
		zap.S().Info("注销失败", zap.Error(err))
	} else {
		zap.S().Info("注销成功")
	}

}

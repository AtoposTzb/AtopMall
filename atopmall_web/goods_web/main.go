package main

import (
	"strconv"

	"go.uber.org/zap"

	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/initialize"
	"atopmall_web/goods_web/utils"
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
	//6.动态获取端口号,本地调试还是使用配置文件的端口号8081,方便apifox调试
	if debug := initialize.GetEnvInfo(global.Env); !debug {
		goodsPort, err := utils.GetAddrPort()
		if err != nil {
			panic(err)
		}
		global.ServerConfig.Port = goodsPort
	}

	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)

	// 启动服务器
	if err := r.Run(":" + strconv.Itoa(global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务器失败：", zap.Error(err))
	}

}

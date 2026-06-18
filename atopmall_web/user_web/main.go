package main

import (
	"strconv"

	"go.uber.org/zap"

	"atopmall_web/user_web/global"
	"atopmall_web/user_web/initialize"
)

func main() {
	//1.初始化logger
	initialize.LoggerInit()
	//2.初始化配置
	initialize.ConfigInit()
	//3.初始化路由
	r := initialize.RoutersInit()

	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)

	// 启动服务器
	if err := r.Run(":" + strconv.Itoa(global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务器失败：", zap.Error(err))
	}

}

package main

import (
	"atopmall_web/initialize"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	port := 8081
	//1.初始化logger
	initialize.LoggerInit()
	//2.初始化路由
	r := initialize.RoutersInit()

	zap.S().Debugf("启动服务器,端口:%d", port)

	// 启动服务器
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		zap.S().Panic("启动服务器失败：", zap.Error(err))
	}

}

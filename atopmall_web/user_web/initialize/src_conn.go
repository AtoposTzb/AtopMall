package initialize

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"atopmall_web/user_web/global"
	"atopmall_web/user_web/proto"
)

func UserSrcClientInit() {
	//从服务注册中心consul获取用户的信息，主要是用户服务的ip和端口号
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = consulInfo.Host + ":" + strconv.Itoa(consulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" || userSrvPort == 0 {
		zap.S().Errorw("[UserSrcClientInit] 连接【用户服务】失败")
		return
	}
	//调用user_web的user.proto接口 GetUserList 也就是远程
	//跨越问题-- 后端解决 也可以前端解决 这里采用后端解决，跨域问题如何产生？详看有道云笔记
	ip := userSrvHost
	port := userSrvPort
	userCoun, err := grpc.NewClient(ip+":"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[UserSrcClientInit]连接【用户服务失败】",
			"msg", err.Error(),
		)
		return
	}
	// 注意：不能在这里 defer userCoun.Close()，否则函数返回后连接就关闭了
	// global.UserSrvClient 需要保持连接可用
	userSrcClient := proto.NewUserClient(userCoun)
	global.UserSrvClient = userSrcClient
	zap.S().Infow("[UserSrcClientInit] 连接【用户服务】成功", "host", ip, "port", port)
}

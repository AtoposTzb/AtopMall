package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/proto"
)

// GPRC负载均衡连接商品服务
func GoodsSrcClientInitBL() {
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[UserSrcClientInitBL] 连接【用户服务】失败")
		return
	}
	goodsSrcClient := proto.NewGoodsClient(conn)
	global.GoodsSrvClient = goodsSrcClient
	zap.S().Infow("[GoodsSrcClientInitBL] 连接【商品服务】成功", "host", consulInfo.Host, "port", consulInfo.Port)
}

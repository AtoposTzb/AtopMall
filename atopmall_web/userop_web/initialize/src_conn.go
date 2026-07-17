package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"atopmall_web/userop_web/global"
	"atopmall_web/userop_web/proto"
)

// GPRC负载均衡连接订单服务(含订单、购物车、库存、商品服务)
// 复用同一个连接创建所有子服务的客户端，统一赋值给聚合结构体OrderRpcClient
func SrcClientInitBL() {
	consulInfo := global.ServerConfig.ConsulInfo
	// 连接商品服务
	goodsConn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[SrcClientInitBL] 连接【商品服务】失败")
		return
	}
	// 一次性初始化所有子服务客户端，统一赋值给全局变量
	global.GoodsSrvCli = &global.GoodsRpcClient{
		Goods:         proto.NewGoodsClient(goodsConn),
		Brand:         proto.NewBrandClient(goodsConn),
		Category:      proto.NewCategoryClient(goodsConn),
		Banner:        proto.NewBannerClient(goodsConn),
		CategoryBrand: proto.NewCategoryBrandClient(goodsConn),
	}
	zap.S().Infow("[SrcClientInitBL] 连接【商品服务】成功", "host", consulInfo.Host, "port", consulInfo.Port)

	// 连接用户操作服务
	userOpConn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[SrcClientInitBL] 连接【用户操作服务】失败")
		return
	}

	global.AddressSrvCli = proto.NewAddressClient(userOpConn)
	zap.S().Infow("[SrcClientInitBL] 连接【地址服务】成功", "host", consulInfo.Host, "port", consulInfo.Port)

	global.UserFavSrvCli = proto.NewUserFavClient(userOpConn)
	zap.S().Infow("[SrcClientInitBL] 连接【收藏服务】成功", "host", consulInfo.Host, "port", consulInfo.Port)

	global.MessageSrvCli = proto.NewMessageClient(userOpConn)
	zap.S().Infow("[SrcClientInitBL] 连接【消息服务】成功", "host", consulInfo.Host, "port", consulInfo.Port)

}

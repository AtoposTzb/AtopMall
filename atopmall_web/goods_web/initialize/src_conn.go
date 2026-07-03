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

// GPRC负载均衡连接商品服务()
// 复用同一个连接创建所有子服务的客户端，统一赋值给聚合结构体
func GoodsSrcClientInitBL() {
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[GoodsSrcClientInitBL] 连接【商品服务】失败")
		return
	}
	// 一次性初始化所有子服务客户端，统一赋值给全局变量
	global.GoodsSrvCli = &global.GoodsRpcClient{
		Goods:         proto.NewGoodsClient(conn),
		Brand:         proto.NewBrandClient(conn),
		Category:      proto.NewCategoryClient(conn),
		Banner:        proto.NewBannerClient(conn),
		CategoryBrand: proto.NewCategoryBrandClient(conn),
	}
	zap.S().Infow("[GoodsSrcClientInitBL] 连接【商品服务】成功", "host", consulInfo.Host, "port", consulInfo.Port)
}

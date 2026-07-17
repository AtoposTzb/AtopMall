package global

import (
	ut "github.com/go-playground/universal-translator"

	"atopmall_web/userop_web/config"
	"atopmall_web/userop_web/proto"
)

// 聚合商品服务所有子服务客户端
type GoodsRpcClient struct {
	Goods         proto.GoodsClient
	Brand         proto.BrandClient
	Category      proto.CategoryClient
	Banner        proto.BannerClient
	CategoryBrand proto.CategoryBrandClient
}

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	Trans         ut.Translator
	Env           string              = "ATOPMALL_DEBUG"
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvCli   *GoodsRpcClient     = &GoodsRpcClient{}
	MessageSrvCli proto.MessageClient
	AddressSrvCli proto.AddressClient
	UserFavSrvCli proto.UserFavClient
)

/*
	初始化全局变量说明(按从上到下的顺序)
	1. ServerConfig: 服务器配置 用于存储服务器的配置信息，包括数据库连接信息、语言翻译信息等
	2. Trans: 语言翻译 用于翻译错误信息、提示信息等，支持多语言环境
	3. Env: 本地调试用的环境变量{自行在系统环境变量中设置}
	4. NacosConfig: nacos配置 用于存储nacos的配置信息 连接nacos配置中心
	5. GoodsSrvCli: 商品服务客户端 用于调用所有商品相关服务的接口
	6. MessageSrvCli: 消息服务客户端 用于调用所有消息相关服务的接口
	7. AddressSrvCli: 地址服务客户端 用于调用所有地址相关服务的接口
	8. UserFavSrvCli: 用户收藏服务客户端 用于调用所有用户收藏相关服务的接口
*/

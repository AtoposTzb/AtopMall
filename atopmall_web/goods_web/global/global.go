package global

import (
	ut "github.com/go-playground/universal-translator"

	"atopmall_web/goods_web/config"
	"atopmall_web/goods_web/proto"
)

// 聚合所有子服务客户端
type GoodsRpcClient struct {
	Goods         proto.GoodsClient
	Brand         proto.BrandClient
	Category      proto.CategoryClient
	Banner        proto.BannerClient
	CategoryBrand proto.CategoryBrandClient
}

var (
	ServerConfig    *config.ServerConfig = &config.ServerConfig{}
	Trans           ut.Translator
	Env             string              = "ATOPMALL_DEBUG"
	NacosConfig     *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvCli     *GoodsRpcClient     = &GoodsRpcClient{}
	InventorySrvCli proto.InventoryClient
)

/*
	初始化全局变量说明(按从上到下的顺序)
	1. ServerConfig: 服务器配置 用于存储服务器的配置信息，包括数据库连接信息、语言翻译信息等
	2. Trans: 语言翻译 用于翻译错误信息、提示信息等，支持多语言环境
	3. Env: 本地调试用的环境变量{自行在系统环境变量中设置}
	4. NacosConfig: nacos配置 用于存储nacos的配置信息 连接nacos配置中心
	5. GoodsRpcCli: 商品服务客户端 用于调用所有商品相关服务的接口
*/

package global

import (
	ut "github.com/go-playground/universal-translator"

	"atopmall_web/goods_web/config"
	"atopmall_web/goods_web/proto"
)

var (
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	Trans          ut.Translator
	GoodsSrvClient proto.GoodsClient
	Env            string              = "ATOPMALL_DEBUG"
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
)

/*
	初始化全局变量说明(按从上到下的顺序)
	1. ServerConfig: 服务器配置 用于存储服务器的配置信息，包括数据库连接信息、语言翻译信息等
	2. Trans: 语言翻译 用于翻译错误信息、提示信息等，支持多语言环境
	3. GoodsSrvClient: 商品服务客户端 用于调用商品服务的接口
	4. Env: 本地调试用的环境变量{自行在系统环境变量中设置}
	5. NacosConfig: nacos配置 用于存储nacos的配置信息 连接nacos配置中心
*/

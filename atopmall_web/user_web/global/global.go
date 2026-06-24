package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"

	"atopmall_web/user_web/config"
	"atopmall_web/user_web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	Trans         ut.Translator
	RDB           *redis.Client
	UserSrvClient proto.UserClient
)

/*
	初始化全局变量说明(按从上到下的顺序)
	1. ServerConfig: 服务器配置 用于存储服务器的配置信息，包括数据库连接信息、语言翻译信息等
	2. Trans: 语言翻译 用于翻译错误信息、提示信息等，支持多语言环境
	3. RDB: redis客户端 用于连接redis数据库
	4. UserSrvClient: 用户服务客户端 用于调用用户服务的接口
*/

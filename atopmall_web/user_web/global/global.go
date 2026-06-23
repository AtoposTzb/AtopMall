package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"

	"atopmall_web/user_web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
	RDB          *redis.Client
)

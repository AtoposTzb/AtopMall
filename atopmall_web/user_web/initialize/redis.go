package initialize

import (
	"atopmall_web/user_web/global"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func RedisInit() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     global.ServerConfig.RedisInfo.Host + ":" + strconv.Itoa(global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       global.ServerConfig.RedisInfo.DB, // use default DB
	})
}

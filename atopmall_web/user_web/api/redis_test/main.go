package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 测试redis获取key值,测试的时候自己在redis中设置一个key,值为123456
func main() {
	tesrrdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// vlaue, err := global.RDB.Get(context.Background(), "test@163.com").Result()
	value, err := tesrrdb.Get(context.Background(), "test@163.com").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(vlaue) //123456
	if err == redis.Nil {
		fmt.Println("key不存在")
	}
	fmt.Println(value) //123456
}

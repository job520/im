package driver

import (
	"fmt"
	"github.com/go-redis/redis"
	"im/tcp/config"
)

func NewRedisClient() *redis.Client {
	host := config.Config.Redis.Host
	port := config.Config.Redis.Port
	password := config.Config.Redis.Password
	database := config.Config.Redis.Database
	addr := fmt.Sprintf("%s:%d", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})
	return rdb
}

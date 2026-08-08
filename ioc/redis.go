package ioc

import (
	"github.com/jw803/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Get().AppRedis,
	})
	return redisClient
}

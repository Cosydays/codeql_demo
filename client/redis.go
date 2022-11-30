package client

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

var (
	RedisClient *redis.Client
)

func InitRedisClient(ctx context.Context) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "xxx",
		DB:       0,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func SetValue2Redis(key string, value string) {
	RedisClient.Set(key, value, time.Second*10)
}

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
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func SetValue2Redis(ctx context.Context, key string, value string) {
	RedisClient.Set(ctx, key, value, time.Second*10)
}

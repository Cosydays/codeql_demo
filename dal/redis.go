package dal

import (
	"context"
	"fmt"

	"github.com/Cosydays/codeql_demo/client"
	"time"
)

func SetValue2Redis(ctx context.Context, key string, value string) {
	client.RedisClient.Set(ctx, key, value, time.Second*10)
}

func GetRedisValue(ctx context.Context, key string) string {
	fmt.Println(key)
	return client.RedisClient.Get(ctx, key).Val()
}

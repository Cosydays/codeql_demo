package dal

import (
	"fmt"

	"github.com/Cosydays/codeql_demo/client"
	"time"
)

func SetValue2Redis(key string, value string) {
	client.RedisClient.Set(key, value, time.Second*10)
}

func GetRedisValue(key string) string {
	fmt.Println(key)
	return client.RedisClient.Get(key).Val()
}

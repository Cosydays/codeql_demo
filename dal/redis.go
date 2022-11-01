package dal

import (
	"context"

	"github.com/Cosydays/codeql_demo/client"
	"time"
)

func SetValue2Redis(ctx context.Context, key string, value string) {
	client.RedisClient.Set(ctx, key, value, time.Second*10)
}

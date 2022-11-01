package redis

import (
	"fmt"
)

func Set(key string, value string) {
	fmt.Println("redis set", key, value)
}

func Get(key string) string {
	return "redis value"
}

package shared

import (
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitGlobalRdb() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: "",
		DB:       1,
	})
}

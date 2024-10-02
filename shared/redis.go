package shared

import (
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitGlobalRdb() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: "",
		DB:       1,
	})
}

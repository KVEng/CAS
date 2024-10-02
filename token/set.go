package token

import (
	"context"
	"github.com/KVEng/CAS/shared"
	"time"
)

func ActiveToken(token string, username string) error {
	return shared.Redis.Set(context.Background(), token, username, time.Hour*24*7).Err()
}

func RemoveToken(token string) error {
	return shared.Redis.Del(context.Background(), token).Err()
}

func IsTokenValid(token string) bool {
	return shared.Redis.Exists(context.Background(), token).Val() == 1
}

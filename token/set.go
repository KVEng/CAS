package token

import (
	"context"
	"github.com/KVEng/CAS/shared"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func ActiveToken(token string, username string) error {
	return shared.Redis.Set(context.Background(), token, username, time.Hour*24*7).Err()
}

func GetTokenUsername(token string) string {
	return shared.Redis.Get(context.Background(), token).Val()
}

func RemoveToken(token string) error {
	return shared.Redis.Del(context.Background(), token).Err()
}

func IsTokenValid(token string) bool {
	return shared.Redis.Exists(context.Background(), token).Val() == 1
}

func HashPasswd(passwd string) string {
	bs, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(bs)
}

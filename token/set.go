package token

import (
	"context"
	"github.com/KVEng/CAS/shared"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func InvalidByUsername(username string) {
	var cursor uint64
	ctx := context.Background()
	r := shared.Redis

	for {
		var keys []string
		var err error
		keys, cursor, err = r.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return
		}

		for _, key := range keys {
			val := r.Get(ctx, key).Val()
			if val != username {
				continue
			}

			if err = r.Del(ctx, key).Err(); err != nil {
				log.Printf("Error deleting key %s: %v", key, err)
			}
		}

		if cursor == 0 {
			break
		}
	}

}

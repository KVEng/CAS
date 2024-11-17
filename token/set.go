package token

import (
	"github.com/KVEng/CAS/shared"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func ActiveToken(token string, username string) error {
	return shared.PkvSession.SetWithTTL(token, username, 60*60*24*7) // 60sec * 60min * 24hour * 7day
}

func GetTokenUsername(token string) string {
	v, _ := shared.PkvSession.Get(token)
	return v
}

func RemoveToken(token string) error {
	return shared.PkvSession.Del(token)
}

func IsTokenValid(token string) bool {
	_, err := shared.PkvSession.Get(token)
	return err == nil
}

func HashPasswd(passwd string) string {
	bs, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(bs)
}

func InvalidByUsername(username string) {
	pkv := shared.PkvSession
	for {
		keys, err := pkv.Keys()
		if err != nil {
			log.Printf("Error getting keys: %v", err)
			break
		}

		for _, key := range keys {
			if GetTokenUsername(key) != username {
				continue
			}

			if err = pkv.Del(key); err != nil {
				log.Printf("Error deleting key %s: %v", key, err)
			}
		}
	}

}

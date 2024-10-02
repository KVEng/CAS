package auth

import (
	"github.com/KVEng/CAS/shared"
	"golang.org/x/crypto/bcrypt"
)

func Verify(user, password, requiredGroup string) bool {
	u, ok := shared.UserDb[user]
	if !ok {
		return false
	}
	if requiredGroup == "" {
		return verify(password, u.Password)
	}

	for _, g := range u.Group {
		if g == requiredGroup {
			return verify(password, u.Password)
		}
	}

	return false
}

func verify(passwd string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)) == nil
}

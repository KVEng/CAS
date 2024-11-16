package auth

import (
	"github.com/KVEng/CAS/shared"
	"golang.org/x/crypto/bcrypt"
)

func Verify(user, password, requiredGroup string) bool {
	passwd, ok := shared.GetUserPassword(user)
	if !ok {
		return false
	}
	if requiredGroup == "" {
		return verify(password, passwd)
	}

	group, _ := shared.GetUserGroups(user)
	for _, g := range group {
		if g == requiredGroup {
			return verify(password, passwd)
		}
	}

	return false
}

func verify(passwd string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)) == nil
}

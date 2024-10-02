package token

import (
	"crypto/rand"
	"math/big"
)

func TokenGenerator() string {
	return SecureToken(32) + ":" + SecureToken(32)
}

const SECURE_CIPHER_ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func SecureToken(length int) string {
	return SecureString(SECURE_CIPHER_ALPHABET, length)
}

func SecureString(alphabet string, length int) string {
	b := make([]byte, length)
	l := len(alphabet)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(l)))
		if err != nil {
			i--
			continue
		}
		b[i] = alphabet[num.Int64()]
	}

	return string(b)
}

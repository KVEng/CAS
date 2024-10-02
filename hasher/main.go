package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println(HashPasswd("passwd"))
}

func HashPasswd(passwd string) string {
	bs, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(bs)
}

package shared

import (
	"fmt"
	"github.com/KVRes/PiccadillySDK/client"
	"github.com/KVRes/PiccadillySDK/types"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"strings"
)

var pkvUser *client.Client
var pkvGroup *client.Client

func InitPKV(addr string) {
	pkv, err := client.NewClient(addr)
	panicx.NotNilErr(err)
	pkvUser = pkv.Copy()
	pkvGroup = pkv.Copy()

	panicx.NotNilErr(pkvUser.Connect("/CAS/User", types.CreateIfNotExist, types.NoLinear))
	panicx.NotNilErr(pkvUser.Connect("/CAS/Group", types.CreateIfNotExist, types.NoLinear))
}

func GetUserPassword(user string) (string, bool) {
	user = strings.ToLower(user)
	v, err := pkvUser.Get(user)
	fmt.Println("PKV", v, err)
	if err != nil {
		return "", false
	}
	return v, true
}

func GetUserGroups(user string) ([]string, bool) {
	user = strings.ToLower(user)
	v, err := pkvGroup.Get(user)
	if err != nil {
		return nil, false
	}
	return strings.Split(v, ","), true
}

func ChangeUserPassword(user string, password string) error {
	user = strings.ToLower(user)
	return pkvUser.Set(user, password)
}

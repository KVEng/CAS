package shared

import (
	"github.com/KVRes/PiccadillySDK/client"
	"github.com/KVRes/PiccadillySDK/types"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"log"
	"strings"
)

var pkv *client.Client
var pkvUser *client.Client
var pkvGroup *client.Client
var PkvSession *client.Client

func initPKV(addr string) error {
	log.Println("initPKV", addr)
	var err error
	pkv, err = client.NewClient(addr)
	if err != nil {
		return err
	}
	pkvUser = pkv.Copy()
	pkvGroup = pkv.Copy()
	PkvSession = pkv.Copy()

	err = pkvUser.Connect("/CAS/User", types.CreateIfNotExist, types.NoLinear)
	if err != nil {
		return err
	}
	err = pkvGroup.Connect("/CAS/Group", types.CreateIfNotExist, types.NoLinear)
	if err != nil {
		return err
	}

	err = PkvSession.Connect("/CAS/Session", types.CreateIfNotExist, types.NoLinear)
	if err != nil {
		return err
	}

	return nil

}
func InitPKV(addr string) {
	panicx.NotNilErr(initPKV(addr))
}

func GetUserPassword(user string) (string, bool) {
	user = strings.ToLower(user)
	v, err := pkvUser.Get(user)
	if err != nil {
		log.Println(err)
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

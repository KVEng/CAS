package shared

import (
	"fmt"
	"github.com/KVRes/PiccadillySDK/client"
	"github.com/KVRes/PiccadillySDK/types"
	"google.golang.org/grpc/connectivity"
	"log"
	"strings"
	"time"
)

var pkv *client.Client
var pkvUser *client.Client
var pkvGroup *client.Client

func initPKV(addr string) error {
	var err error
	pkv, err = client.NewClient(addr)
	if err != nil {
		return err
	}
	pkvUser = pkv.Copy()
	pkvGroup = pkv.Copy()

	err = pkvUser.Connect("/CAS/User", types.CreateIfNotExist, types.NoLinear)
	if err != nil {
		return err
	}
	err = pkvGroup.Connect("/CAS/Group", types.CreateIfNotExist, types.NoLinear)
	if err != nil {
		return err
	}
	return nil

}
func InitPKV(addr string) {
	err := initPKV(addr)
	if err != nil {
		panic(err)
	}

	go func() {
		conn := pkv.GetConn()
		for {
			state := conn.GetState()
			fmt.Println("Piccadilly connection state:", state)
			if state == connectivity.TransientFailure {
				_pkv, _err := client.NewClient(addr)
				if _err != nil {
					log.Println(_err)
					continue
				}
				conn = pkv.GetConn()
				pkv = _pkv
				pkvGroup = pkv.Copy()
				pkvUser = pkv.Copy()
			}
			time.Sleep(5 * time.Second)
		}
	}()
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

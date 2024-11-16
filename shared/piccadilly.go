package shared

import (
	"fmt"
	"github.com/KVRes/PiccadillySDK/client"
	"github.com/KVRes/PiccadillySDK/types"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"google.golang.org/grpc/connectivity"
	"strings"
	"time"
)

var pkvUser *client.Client
var pkvGroup *client.Client

func InitPKV(addr string) {
	pkv, err := client.NewClient(addr)
	panicx.NotNilErr(err)

	go func() {
		conn := pkv.GetConn()
		for {
			state := conn.GetState()
			if state != connectivity.Ready {
				fmt.Println("Piccadilly is not ready: ", state)
				conn.Connect()
			}
			time.Sleep(10 * time.Second)
		}
	}()
	pkvUser = pkv.Copy()
	pkvGroup = pkv.Copy()

	panicx.NotNilErr(pkvUser.Connect("/CAS/User", types.CreateIfNotExist, types.NoLinear))
	panicx.NotNilErr(pkvGroup.Connect("/CAS/Group", types.CreateIfNotExist, types.NoLinear))
}

func GetUserPassword(user string) (string, bool) {
	user = strings.ToLower(user)
	v, err := pkvUser.Get(user)
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

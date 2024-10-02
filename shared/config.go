package shared

import (
	"encoding/json"
	"github.com/KVEng/CAS/model"
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/KevinZonda/GoX/pkg/panicx"
)

var Config model.Config

var UserDb map[string]model.User

func InitGlobalCfg() {
	bs, err := iox.ReadAllByte("config.json")
	panicx.NotNilErr(err)
	err = json.Unmarshal(bs, &Config)
	panicx.NotNilErr(err)
	for _, u := range Config.User {
		UserDb[u.Username] = u
	}
}

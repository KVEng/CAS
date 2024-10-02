package shared

import (
	"encoding/json"
	"fmt"
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

	UserDb = make(map[string]model.User)
	for _, u := range Config.User {
		UserDb[u.Username] = u
	}
}

func copyDb() map[string]model.User {
	db := make(map[string]model.User)
	for k, v := range UserDb {
		db[k] = v
	}
	return db
}

func ModifyUserDb(f func(db map[string]model.User)) error {
	db := copyDb()
	f(db)
	var users []model.User
	for _, v := range UserDb {
		users = append(users, v)
	}
	bs, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = iox.WriteAllBytes("config.json", bs)

	if err != nil {
		fmt.Println(err)
		return err
	}

	Config.User = users
	UserDb = db
	return nil
}

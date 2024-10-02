package model

type Config struct {
	Debug      bool   `json:"debug"`
	ListenAddr string `json:"listen_addr"`
	RedisAddr  string `json:"redis_addr"`
	User       []User `json:"user"`
}

type User struct {
	Username string   `json:"username"`
	Group    []string `json:"group"`
	Password string   `json:"password"`
}

func (u User) IsInGroup(group string) bool {
	if group == "" {
		group = "admin"
	}
	for _, g := range u.Group {
		if g == group {
			return true
		}
	}
	return false
}

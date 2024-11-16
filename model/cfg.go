package model

type Config struct {
	Debug      bool   `json:"debug"`
	ListenAddr string `json:"listen_addr"`
	RedisAddr  string `json:"redis_addr"`
}

func IsInGroup(uGroups []string, required string) bool {
	if required == "" {
		required = "admin"
	}
	for _, g := range uGroups {
		if g == required || g == "admin" {
			return true
		}
	}
	return false
}

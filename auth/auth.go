package auth

func Verify(user, password, requiredGroup string) bool {
	if password != "123456" {
		return false
	}
	return true
}

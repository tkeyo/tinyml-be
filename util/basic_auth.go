package util

func IsAuthorized(authKey string) bool {
	if authKey != "123" {
		return false
	}
	return authKey == "123"
}

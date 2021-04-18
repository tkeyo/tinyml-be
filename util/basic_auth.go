package util

func IsAuthorized(requestAuthKey string, APIAuthKey string) bool {
	if requestAuthKey != APIAuthKey {
		return false
	}
	return requestAuthKey == APIAuthKey
}

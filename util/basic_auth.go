package util

import (
	"time"
)

func IsAuthorized(requestAuthKey string, APIAuthKey string) bool {
	if requestAuthKey != APIAuthKey {
		return false
	}
	return requestAuthKey == APIAuthKey
}

func GetCurrentTime() int64 {
	time := time.Now().UTC().Unix() * 1000
	return time
}

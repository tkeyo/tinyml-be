package util

import (
	"time"
)

// Simple auth implementation
func IsAuthorized(requestAuthKey string, APIAuthKey string) bool {
	if requestAuthKey != APIAuthKey {
		return false
	}
	return requestAuthKey == APIAuthKey
}

// Gets current UTC time
func GetCurrentTime() int64 {
	time := time.Now().UTC().Unix() * 1000
	return time
}

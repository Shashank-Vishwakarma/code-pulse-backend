package utils

import "time"

func GetCurrentDateTime() int64 {
	return time.Now().Unix()
}

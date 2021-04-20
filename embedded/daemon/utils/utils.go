package utils

import "log"

func ErrLog(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

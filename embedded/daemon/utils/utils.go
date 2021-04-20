package utils

import "log"

func ErrLog(err error) {
	if err != nil {
		log.Println(err)
	}
}

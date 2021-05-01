package utils

import (
	"log"
	"os"
)

func ErrLog(prefix string, err error) bool {
	if err != nil {
		f, err2 := os.OpenFile("dev.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			log.Println(err2)
		}
		defer f.Close()
		logger := log.New(f, prefix, log.LstdFlags)
		logger.Println(err)
		return true
	}
	return false
}

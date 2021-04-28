package utils

import (
	"log"
	"os"

	b "../bootstrap"
)

func ErrLog(prefix string, err error) bool {
	if err != nil {
		f, err := os.OpenFile(b.LOGS,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		logger := log.New(f, prefix, log.LstdFlags)
		logger.Println(err)
		return true
	}
	return false
}

package actions

import (
	"encoding/json"
	"log"
	"os"

	d "../domain"
)

func WriteCache(data d.Report) {
	b, _ := json.Marshal(data)
	f, err := os.Create("cache.json")
	if err != nil {
		log.Println(err)
	}
	f.Write(b)
}

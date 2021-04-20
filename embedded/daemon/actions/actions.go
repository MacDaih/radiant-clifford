package actions

import (
	"encoding/json"
	"os"

	d "../domain"
	u "../utils"
)

func WriteCache(data d.Report) {
	b, err := json.Marshal(data)
	f, err := os.Create("cache.json")
	u.ErrLog(err)
	f.Write(b)
}

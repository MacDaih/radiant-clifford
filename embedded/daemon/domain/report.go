package domain

type Report struct {
	RptAt int64   `json:"time"`
	Temp  float64 `json:"t"`
	Hum   float64 `json:"h"`
	Light int32   `json:"l"`
}

type Reports []Report

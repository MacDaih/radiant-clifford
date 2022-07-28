package domain

type Report struct {
	RptAt int64   `bson:"report_time" json:"time"`
	Temp  float64 `bson:"temp" json:"temp"`
	Hum   float64 `bson:"hum" json:"hum"`
	Light float64 `bson:"ligth" json:"light"`
	Press float64 `bson:"press" json:"press"`
}

type ReportSample struct {
	Metrics Overview
	Reports []Report
}
type Overview struct {
	TempAverage float64 `json:"temp_av"`
	HumAverage  float64 `json:"hum_av"`
	MaxTemp     float64 `json:"max_temp"`
	MinTemp     float64 `json:"min_temp"`
	MaxHum      float64 `json:"max_hum"`
	MinHum      float64 `json:"min_hum"`
}

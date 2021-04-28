package domain

type Report struct {
	RptAt int64   `bson:"report_time" json:"time"`
	Temp  float64 `bson:"temp" json:"t"`
	Hum   float64 `bson:"hum" json:"h"`
	Light int32   `bson:"ligth" json:"l"`
}

type Reports []Report

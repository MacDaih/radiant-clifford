package domain

import (
	"time"
)

type Report struct {
	RptAt time.Time
	Temp  float64
	Hum   float64
	Light int32
}

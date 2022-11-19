package domain

import "math"

const (
	TW4  = 86200
	TWEL = 43200
	SIX  = 21600
	HOUR = 3600
	HALF = 1800
	FFTN = 900

	STR_TW4  = "last_tw4"
	STR_TWEL = "last_twelve"
	STR_SIX  = "last_six_hrs"
	STR_HOUR = "last_hr"
	STR_HALF = "last_half"
	STR_FFTN = "last_fftn"
)

func ToStamp(input string) int64 {
	switch input {
	case STR_FFTN:
		return FFTN
	case STR_HALF:
		return HALF
	case STR_HOUR:
		return HOUR
	case STR_SIX:
		return SIX
	case STR_TWEL:
		return TWEL
	default:
		return TW4
	}
}

func average(input []float64) float64 {
	value := 0.0

	for _, i := range input {
		value += i
	}

	return math.Round((value/float64(len(input)))*100) / 100
}

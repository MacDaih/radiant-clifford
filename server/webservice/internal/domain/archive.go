package domain

import (
	"fmt"
	"time"
)

type Archive struct {
	Label  string
	Period TimeRange
	Report
}

func FormatArchive(archiveRange TimeRange, reports []Report) Archive {

	var temps, humids, pressures []float64

	for _, r := range reports {
		temps = append(temps, r.Temp)
		humids = append(humids, r.Hum)
		pressures = append(pressures, r.Press)
	}

	from := time.Unix(archiveRange.From, 0)

	label := fmt.Sprintf("%s_%d", from.Month().String(), from.Year())

	return Archive{
		Label:  label,
		Period: archiveRange,
		Report: Report{
			ReportedAt: archiveRange.From,
			Temp:       average(temps),
			Hum:        average(humids),
			Press:      average(pressures),
		},
	}
}

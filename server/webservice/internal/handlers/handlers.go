package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net"
	"net/http"
	"time"
	d "webservice/internal/domain"
	"webservice/internal/repository"
)

const (
	TWE = 43200
)

// For dev purpose only
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type serviceHandler struct {
	repository repository.Repository
}
type Handler interface {
	ReadSock(conn net.Conn, e chan error)
	ReportsHandler(w http.ResponseWriter, r *http.Request)
}

func NewServiceHandler(repository repository.Repository) Handler {
	return &serviceHandler{
		repository: repository,
	}
}

func (s *serviceHandler) ReadSock(conn net.Conn, e chan error) {

	var r d.Report
	dc := json.NewDecoder(conn)

	err := dc.Decode(&r)
	if err != nil {
		log.Println("decoding error : ", err)
		e <- err
	}

	err = s.repository.InsertReport(context.Background(), r)
	if err != nil {
		log.Println("failed to report error : ", err)
		e <- err
	}
}

func (s *serviceHandler) ReportsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	t := time.Now().Unix()
	last := t - TWE
	reports, err := s.repository.GetReports(r.Context(), last)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	sample := formatSample(reports)
	json.NewEncoder(w).Encode(sample)
}

func formatSample(reports []d.Report) d.ReportSample {
	if len(reports) > 0 {
		var o d.Overview
		var maxTemp float64 = reports[0].Temp
		var minTemp float64 = reports[0].Temp
		var maxHum float64 = reports[0].Hum
		var minHum float64 = reports[0].Hum
		for _, j := range reports {
			if maxTemp < j.Temp {
				maxTemp = j.Temp
			}
			if minTemp > j.Temp {
				minTemp = j.Temp
			}
			if maxHum < j.Hum {
				maxHum = j.Temp
			}
			if minHum > j.Hum {
				minHum = j.Hum
			}
		}
		avHum, avTemp := average(reports)
		o = d.Overview{
			TempAverage: avTemp,
			HumAverage:  avHum,
			MaxTemp:     maxTemp,
			MinTemp:     minTemp,
			MaxHum:      maxHum,
			MinHum:      minHum,
		}
		return d.ReportSample{
			Metrics: o,
			Reports: reports,
		}
	}
	return d.ReportSample{}
}

func average(r []d.Report) (float64, float64) {
	hum := 0.0
	temp := 0.0
	for _, j := range r {
		hum += j.Hum
		temp += j.Temp
	}
	hum = math.Round((hum/float64(len(r)))*100) / 100
	temp = math.Round((temp/float64(len(r)))*100) / 100
	return hum, temp
}

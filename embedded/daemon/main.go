package main

import (
	"bufio"
	"encoding/json"

	"time"

	a "./actions"
	b "./bootstrap"
	d "./domain"
	u "./utils"
	"github.com/tarm/serial"
)

func main() {
	path := b.SetDevice()
	c := &serial.Config{Name: path, Baud: 9600}
	s, err := serial.OpenPort(c)
	if u.ErrLog("Reading Err : ", err) {
		return
	}
	scanner := bufio.NewScanner(s)
	var reports []d.Report
	for scanner.Scan() {
		r, err := readSerial(scanner.Bytes())
		if u.ErrLog("Report Err : ", err) {
			continue
		}
		reports = append(reports, r)
		if len(reports) == 4 {
			a.InsertReports(reports)
			reports = reports[4:]
		}
		time.Sleep(time.Millisecond * 5000)
	}
	err = scanner.Err()
	u.ErrLog("Scan Err : ", err)
}

func readSerial(s []byte) (d.Report, error) {
	var r d.Report = d.Report{}
	now := time.Now().Unix()
	r.RptAt = now
	err := json.Unmarshal(s, &r)
	if err != nil {
		return d.Report{}, err
	}
	return r, nil
}

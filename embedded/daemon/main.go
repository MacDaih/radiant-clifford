package main

import (
	"bufio"
	"encoding/json"

	"fmt"
	"time"

	a "daemon/actions"
	b "daemon/bootstrap"
	d "daemon/domain"
	u "daemon/utils"

	"github.com/tarm/serial"
)

const (
	HOTTEST = 60.00
	COLDEST = -88.00
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
		time.Sleep(time.Millisecond * 300000)
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
	if r.Temp > HOTTEST {
		hotErr := fmt.Errorf("Recorded Temp. is exceeding a normal treshold (%d °C) : %d", r.Temp)
		return d.Report{}, hotErr
	} else if r.Temp < COLDEST {
		hotErr := fmt.Errorf("Recorded Temp. is below a negative treshold (%d °C) : %d", r.Temp)
		return d.Report{}, hotErr
	}
	return r, nil
}

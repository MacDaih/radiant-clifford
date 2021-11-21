package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	a "daemon/actions"
	d "daemon/domain"
	u "daemon/utils"
)

const (
	SOCK    = "/tmp/thermo.sock"
	HOTTEST = 60.00
	COLDEST = -88.00
	KEY     = "thermo"
)

func readSock(conn io.Reader) {
	buf := make([]byte, 256)
	var reports []d.Report
	for {
		n, err := conn.Read(buf[:])
		if ok := u.ErrLog("reading error : ", err); !ok {
			continue
		}
		report, err := readSerial(buf[0:n])
		if ok := u.ErrLog("serialization error : ", err); !ok {
			continue
		}
		reports = append(reports, report)

		if len(reports) == 4 {
			a.InsertReports(reports)
			reports = reports[4:]
		}
	}
}

func main() {
	conn, err := net.Dial("unix", SOCK)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	go readSock(conn)
	for {
		_, err := conn.Write([]byte(KEY))
		if u.ErrLog("writing err : ", err) {
			return
		}
		time.Sleep(time.Minute * 1)
	}
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
		hotErr := fmt.Errorf("recorded temp. is exceeding a normal treshold (%f °C) : %f", HOTTEST, r.Temp)
		return d.Report{}, hotErr
	} else if r.Temp < COLDEST {
		hotErr := fmt.Errorf("recorded temp. is below a negative treshold (%f °C) : %f", COLDEST, r.Temp)
		return d.Report{}, hotErr
	}
	return r, nil
}

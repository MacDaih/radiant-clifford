package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	d "webservice/domain"
	u "webservice/utils"
)

const (
	HOTTEST = 60.00
	COLDEST = -88.00
)

func ReadSock(conn io.Reader, e chan error) {
	buf := make([]byte, 256)
	var reports []d.Report
	for {
		n, err := conn.Read(buf[:])
		if ok := u.ErrLog("reading error : ", err); !ok {
			continue
		}
		report, err := readSerial(buf[0:n])
		log.Println(report)
		if ok := u.ErrLog("serialization error : ", err); !ok {
			continue
		}
		reports = append(reports, report)

		if len(reports) == 4 {
			err = d.InsertReports(reports)
			if ok := u.ErrLog("insert error : ", err); !ok {
				e <- err
			}
			reports = reports[4:]
		}
	}
}

func readSerial(s []byte) (d.Report, error) {

	r := d.Report{
		RptAt: time.Now().Unix(),
	}

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

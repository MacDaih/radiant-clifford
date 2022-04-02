package collector

import (
	"encoding/json"
	"log"
	"net"

	d "webservice/domain"
)

const (
	HOTTEST = 60.00
	COLDEST = -88.00
)

func ReadSock(conn net.Conn, e chan error) {

	var r d.Report
	dc := json.NewDecoder(conn)

	err := dc.Decode(&r)
	if err != nil {
		log.Println("decoding error : ", err)
		e <- err
	}

	err = d.InsertReport(r)
	if err != nil {
		log.Println("failed to report error : ", err)
		e <- err
	}
}

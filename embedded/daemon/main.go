package main

import (
	"encoding/json"
	"fmt"

	b "./bootstrap"
	d "./domain"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: b.PATH, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("can't open /dev/ttyACM0y: %s \n", err)
	}
	for {
		buf := make([]byte, 128)
		n, err := s.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		// log.Printf("%s", buf[:n])
		var r d.Report
		res := json.Unmarshal(buf[:n], &r)
		fmt.Println(res)
	}
}

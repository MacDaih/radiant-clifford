package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"time"

	a "./actions"
	b "./bootstrap"
	d "./domain"
	"github.com/tarm/serial"
)

func main() {
	path := b.SetDevice()
	c := &serial.Config{Name: path, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("%s \n", err)
	}
	scanner := bufio.NewScanner(s)

	for scanner.Scan() {
		res, err := readSerial(scanner.Bytes())
		if err != nil {
			continue
		} else {
			fmt.Println(res)
		}
		a.WriteCache(res)
		time.Sleep(time.Millisecond * 5000)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
	return r, nil
}

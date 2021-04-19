package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"time"

	b "./bootstrap"
	d "./domain"
	"github.com/tarm/serial"
)

func main() {
	path := b.SetDevice()
	fmt.Println(path)
	c := &serial.Config{Name: path, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("%s \n", err)
	}
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		readSerial(scanner.Bytes())
		time.Sleep(time.Millisecond * 5000)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func readSerial(s []byte) {
	var r d.Report
	now := time.Now().Unix()
	r.RptAt = now
	err := json.Unmarshal(s, &r)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(r)
	}
}

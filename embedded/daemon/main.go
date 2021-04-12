package main

import (
	"os"
	"log"
	"fmt"
)

func main() {
	tty, err := os.Open("/dev/ttyACM0")
	if err != nil {
		log.Fatalf("can't open /dev/tty: %s", err)
	}

	b := make([]byte,5)

	r, err := tty.Read(b)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Result => %s \n", string(b[:r]))
}
package main

import (
	"os"
	"log"
	"fmt"
	"time"
)

func main() {
	for {
		tty, err := os.Open("/dev/ttyACM0")
		if err != nil {
			log.Fatalf("can't open /dev/tty: %s", err)
		}

		b := make([]byte,512)

		r, err := tty.Read(b)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Result => %s \n", string(b[:r]))
		time.Sleep(time.Millisecond * 5000)
	}
}
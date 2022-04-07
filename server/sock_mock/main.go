package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

func main() {
	lst, err := net.Listen("tcp", ":9090")

	if err != nil {
		log.Fatal(err)
	}
	defer lst.Close()
	conn, err := lst.Accept()
	if err != nil {
		log.Fatal(err)
	}
	for {
		if _, err := conn.Write(randomize(0.00, 100.00)); err != nil {
			log.Println(err)
			return
		}
	}
}

func randomize(max float32, min float32) []byte {
	temp := min + rand.Float32()*(max-min)

	hum := min + rand.Float32()*(max-min)

	res := fmt.Sprintf(`{"temp": %.2f, "hum": %.2f}`, temp, hum)

	return []byte(res)
}

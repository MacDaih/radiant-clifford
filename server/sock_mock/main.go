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

	for {
		conn, err := lst.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		b := make([]byte, 1024)
		_, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
			break
		}
		conn.Write(randomize(0.00, 100.00))
	}
}

func randomize(max float32, min float32) []byte {
	temp := min + rand.Float32()*(max-min)

	hum := min + rand.Float32()*(max-min)

	res := fmt.Sprintf(`{"temp": %.2f, "hum": %.2f}`, temp, hum)
	return []byte(res)
}

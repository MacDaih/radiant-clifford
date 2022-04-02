package main

import (
	"log"
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
		if _, err := conn.Write([]byte(`{"temp":0.00,"hum":0.00}`)); err != nil {
			log.Println(err)
			return
		}
	}
}

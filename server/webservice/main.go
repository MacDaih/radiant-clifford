package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
	c "webservice/collector"
	h "webservice/handlers"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	socket := os.Getenv("SENSOR_PORT")
	key := os.Getenv("KEY")

	httpError := make(chan error)
	collError := make(chan error)

	go expose(port, httpError)
	go collect(socket, key, collError)

	select {
	case err := <-httpError:
		log.Fatalf("Http Server error : %s", err)
	case err := <-collError:
		log.Fatalf("Data collector error : %s", err)
	}
}

func expose(p string, e chan error) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/reports", h.ReportsHandler).Methods("GET")
	router.Handle("/", router)

	e <- http.ListenAndServe(p, router)
}

func collect(socket string, key string, e chan error) {
	log.Println("collecting from tpc")
	conn, err := net.Dial("tcp", socket)

	if err != nil {
		e <- err
	}

	go c.ReadSock(conn, e)

	for {
		if _, err := conn.Write([]byte(key)); err != nil {
			log.Println(err)
			e <- err
		}
		time.Sleep(time.Minute * 1)
	}
}
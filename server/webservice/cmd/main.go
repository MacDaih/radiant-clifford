package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"webservice/internal/handlers"

	"webservice/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting webservice")
	port := os.Getenv("PORT")
	socket := os.Getenv("SENSOR_PORT")
	key := os.Getenv("KEY")

	dbName := os.Getenv("DB_NAME")

	repo := repository.NewReportRepository(dbName)

	hdlr := handlers.NewServiceHandler(repo)

	httpError := make(chan error)
	collError := make(chan error)

	go expose(port, hdlr, httpError)
	go collect(socket, key, hdlr, collError)
	log.Println("Collect and expose")
	select {
	case err := <-httpError:
		log.Fatalf("http Server error : %s", err)
	case err := <-collError:
		log.Fatalf("data collector error : %s", err)
	}
}

func expose(p string, h handlers.Handler, e chan error) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/reports", h.ReportsHandler).Methods("GET")
	router.Handle("/", router)

	e <- http.ListenAndServe(p, router)
}

func collect(socket string, key string, h handlers.Handler, e chan error) {
	log.Println("Collector running")
	addr, err := net.ResolveTCPAddr("tcp4", socket)
	if err != nil {
		log.Printf("resolving address error : %s\n", err.Error())
		e <- err
	}
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		log.Printf("dial to address error : %s\n", err.Error())
		e <- err
	}

	defer conn.Close()

	for {
		_, err := conn.Write([]byte(key))
		if err != nil {
			log.Println(err)
			continue
		}
		h.ReadSock(conn)
		<-time.After(time.Second * 10)
	}
}

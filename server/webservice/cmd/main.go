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

	select {
	case err := <-httpError:
		log.Fatalf("Http Server error : %s", err)
	case err := <-collError:
		log.Fatalf("Data collector error : %s", err)
	}
}

func expose(p string, h handlers.Handler, e chan error) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/reports", h.ReportsHandler).Methods("GET")
	router.Handle("/", router)

	e <- http.ListenAndServe(p, router)
}

func collect(socket string, key string, h handlers.Handler, e chan error) {

	conn, err := net.Dial("tcp", socket)

	if err != nil {
		e <- err
	}

	defer conn.Close()

	for {
		_, err := conn.Write([]byte(key))
		time.Sleep(time.Second * 3)
		if err != nil {
			log.Println("tcp err : ", err)
			continue
		}
		go h.ReadSock(conn)
	}
}

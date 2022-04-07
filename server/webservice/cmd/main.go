package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"webservice/internal/database"
	"webservice/internal/handlers"

	"webservice/internal/repository"

	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	socket := os.Getenv("SENSOR_PORT")
	key := os.Getenv("KEY")

	dbName := os.Getenv("DB_NAME")

	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbClient, err := database.ConnectDB(host, dbPort)

	if err != nil {
		log.Fatal(err)
	}

	db := dbClient.Database(dbName)

	repo := repository.NewReportRepository(db)

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

	go h.ReadSock(conn, e)

	for {
		if _, err := conn.Write([]byte(key)); err != nil {
			log.Println(err)
			e <- err
		}
		time.Sleep(time.Minute * 1)
	}
}

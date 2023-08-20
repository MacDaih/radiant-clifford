package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webservice/config"
	"webservice/internal/collector"
	"webservice/internal/handler"
	"webservice/internal/repository"

	"webservice/cmd/server"
	"webservice/cmd/worker"

	"github.com/gorilla/mux"
)

func main() {

	config.Boot()

	log.Println("Starting webservice")

	repo := repository.NewReportRepository(config.GetDBEnv())
	cltr := collector.NewCollector(repo)
	hdlr := handler.NewServiceHandler(repo)

	httpError := make(chan error)
	workerErr := make(chan error)
	sysInt := make(chan os.Signal)

	router := mux.NewRouter()

	router.HandleFunc("/reports/{range}", hdlr.GetReportsFrom).Methods("GET")
	router.HandleFunc("/by_date/{date}", hdlr.GetReportsByDate).Methods("GET")

	webservice := http.Server{
		Addr:    config.GetPort(),
		Handler: router,
	}

	go server.RunWebservice(&webservice, httpError)

	go worker.Process(
		config.GetSensorPort(),
		config.GetSensorKey(),
		cltr,
		workerErr,
	)

	signal.Notify(sysInt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-httpError:
		if err != nil {
			log.Fatalf("http Server error : %s", err.Error())
		}
	case err := <-workerErr:
		if err != nil {
			log.Fatalf("data collector error : %s", err.Error())
		}
	case <-sysInt:
		log.Println("interrupt : webservice is shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		defer cancel()

		if err := webservice.Shutdown(ctx); err != nil {
			log.Printf("error when shutting down server : %s", err)
			log.Println("closing webservice ...")
			webservice.Close()
		}

		return
	}
}

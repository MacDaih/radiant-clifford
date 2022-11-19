package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"webservice/internal/collector"
	"webservice/internal/handler"
	"webservice/internal/repository"

	"webservice/cmd/server"
	"webservice/cmd/worker"
)

func main() {
	log.Println("Starting webservice")

	port := os.Getenv("PORT")
	socket := os.Getenv("SENSOR_PORT")
	key := os.Getenv("KEY")

	dbName := os.Getenv("DB_NAME")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	repo := repository.NewReportRepository(dbName, dbhost, dbport)
	cltr := collector.NewCollector(repo)
	hdlr := handler.NewServiceHandler(repo)

	httpError := make(chan error)
	collError := make(chan error)
	sysInt := make(chan os.Signal, 2)

	go server.RunWebservice(port, hdlr, httpError)

	go worker.Process(socket, key, cltr, collError)

	signal.Notify(sysInt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-httpError:
		log.Fatalf("http Server error : %s", err)
	case err := <-collError:
		log.Fatalf("data collector error : %s", err)
	case <-sysInt:
		log.Println("interrupt : webservice is shutting down")
		return
	}
}

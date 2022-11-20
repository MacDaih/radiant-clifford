package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"webservice/config"
	"webservice/internal/collector"
	"webservice/internal/handler"
	"webservice/internal/repository"

	"webservice/cmd/server"
	"webservice/cmd/worker"
)

func main() {

	config.Boot()

	log.Println("Starting webservice")

	repo := repository.NewReportRepository(config.GetDBEnv())
	cltr := collector.NewCollector(repo)
	hdlr := handler.NewServiceHandler(repo)

	httpError := make(chan error)
	collError := make(chan error)
	sysInt := make(chan os.Signal, 2)

	go server.RunWebservice(config.GetPort(), hdlr, httpError)

	go worker.Process(config.GetSocket(), config.GetSensorKey(), cltr, collError)

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

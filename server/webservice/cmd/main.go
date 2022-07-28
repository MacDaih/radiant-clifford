package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"webservice/internal/collector"
	"webservice/internal/handler"
	"webservice/internal/repository"

	httpserver "webservice/pkg/http_server"
	tcpclient "webservice/pkg/tcp_client"
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

	hdlr := handler.NewServiceHandler(repo)
	cltr := collector.NewCollector(repo)

	httpError := make(chan error)
	collError := make(chan error)
	sysInt := make(chan os.Signal, 2)

	routes := []httpserver.Route{
		{
			Path:   "/reports",
			Fn:     hdlr.ReportsHandler,
			Method: "GET",
		},
	}

	go httpserver.HttpServe(port, routes, httpError)

	go tcpclient.RunTCPCLient(socket, key, cltr.ReadSock, collError)

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

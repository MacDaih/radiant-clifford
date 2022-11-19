package worker

import (
	"time"
	"webservice/internal/collector"
	tcpclient "webservice/pkg/tcp_client"
)

func Process(socket, key string, collector collector.Collector, err chan error) {
	go func() {
		for {
			if time.Now().Day() >= 1 {
				collector.CleanUpWithArchive()
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	err <- tcpclient.RunTCPCLient(socket, key, collector.ReadSock)
}

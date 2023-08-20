package worker

import (
	"fmt"
	"log"
	"time"
	"webservice/internal/collector"
	tcpclient "webservice/pkg/tcp_client"
)

func Process(socket, key string, collector collector.Collector, we chan error) {
	go func() {
		for {
			if time.Now().Day() >= 1 {
				if err := collector.CleanUpWithArchive(); err != nil {
					log.Printf("failed to archive records : %s", err.Error())
				}
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	we <- retry(10, func() error {
		return tcpclient.RunTCPCLient(socket, key, collector.ReadSock)
	})
}

func retry(attempts int, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		} else {
			time.Sleep(30 * time.Second)
		}
	}

	return fmt.Errorf("failed to retry after %d attempts : %w", attempts, err)
}

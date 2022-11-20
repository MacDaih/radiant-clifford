package worker

import (
	"fmt"
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

	err <- retry(10, func() error {
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
			time.Sleep(10 * time.Minute)
		}
	}

	return fmt.Errorf("failed to retry after %d attempts : %w", attempts, err)
}

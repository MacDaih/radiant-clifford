package worker

import (
	"webservice/internal/collector"
	tcpclient "webservice/pkg/tcp_client"
)

func Process(socket, key string, collector collector.Collector, err chan error) {
	err <- tcpclient.RunTCPCLient(socket, key, collector.ReadSock)
}

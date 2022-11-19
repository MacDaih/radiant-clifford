package collector

import (
	"context"
	"encoding/json"
	"net"
	"time"
	"webservice/internal/domain"
	"webservice/internal/repository"
)

type serviceCollector struct {
	repository repository.Repository
}

type Collector interface {
	ReadSock(conn net.Conn) error
}

func NewCollector(repository repository.Repository) Collector {
	return &serviceCollector{
		repository: repository,
	}
}

func (s *serviceCollector) ReadSock(conn net.Conn) error {
	var r domain.Report
	buf := make([]byte, 128)
	rb, err := conn.Read(buf)

	if err != nil {
		return err
	}

	index := 0
	for i, v := range buf[:rb] {
		if string(v) == "}" {
			index = i + 1
		}
	}

	err = json.Unmarshal(buf[:index], &r)
	if err != nil {
		return err
	}

	r.ReportedAt = time.Now().Unix()
	return s.repository.InsertReport(context.Background(), r)
}

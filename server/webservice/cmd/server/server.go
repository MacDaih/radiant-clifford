package server

import (
	"webservice/internal/handler"
	httpserver "webservice/pkg/http_server"
)

func RunWebservice(port string, service handler.Handler, err chan error) {

	routes := []httpserver.Route{
		{
			Path:   "/reports/{range}",
			Fn:     service.GetReportsFrom,
			Method: "GET",
		},
		{
			Path:   "/by_date/{date}",
			Fn:     service.GetReportsByDate,
			Method: "GET",
		},
	}

	err <- httpserver.HttpServe(port, routes)
}

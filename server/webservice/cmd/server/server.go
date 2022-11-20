package server

import (
	"log"
	"net/http"
	"webservice/internal/handler"

	"github.com/gorilla/mux"
)

func RunWebservice(port string, service handler.Handler, err chan error) {

	log.Println("Running HTTP Server")

	router := mux.NewRouter()

	router.HandleFunc("/reports/{range}", service.GetReportsFrom).Methods("GET")
	router.HandleFunc("/by_date/{date}", service.GetReportsByDate).Methods("GET")

	err <- http.ListenAndServe(port, router)
}

package httpserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Path   string
	Fn     func(w http.ResponseWriter, r *http.Request)
	Method string
}

func HttpServe(port string, routes []Route) error {
	log.Println("Running HTTP Server")
	router := mux.NewRouter().StrictSlash(true)

	for _, r := range routes {
		router.HandleFunc(r.Path, r.Fn).Methods(r.Method)
	}

	router.Handle("/", router)

	return http.ListenAndServe(port, router)
}

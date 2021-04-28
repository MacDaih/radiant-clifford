package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")
	router.Handle("/", router)
	http.ListenAndServe(":8080", router)
}

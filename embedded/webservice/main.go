package main

import (
	"fmt"
	"log"
	"net/http"
	h "webservice/handlers"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Webservice")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/reports", h.ReportsHandler).Methods("GET")
	router.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

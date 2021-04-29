package main

import (
	"encoding/json"
	"net/http"
	d "webservice/domain"
	u "webservice/utils"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		reports, err := d.GetReports()
		if u.ErrLog("Get Reports Err : ", err) {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(reports)
		}
	}).Methods("GET")
	router.Handle("/", router)
	http.ListenAndServe(":8080", router)
}

package kaput

import (
	"net/http"

	"github.com/gorilla/mux"
)

func initRouting() {
	r := mux.NewRouter()
	r.HandleFunc("/", helpHandler("/")).Methods("GET")

	r.HandleFunc("/crash", helpHandler("/crash/")).Methods("GET")
	r.HandleFunc("/crash/{it:(?:connection|server)}", handleCrash).Methods("GET")

	r.HandleFunc("/response", helpHandler("/response/")).Methods("GET")
	r.HandleFunc("/response/{code:[12345][0-9]{2}}", handleResponse).Methods("GET")

	r.HandleFunc("/recursive", helpHandler("/recursive")).Methods("GET")
	r.HandleFunc("/recursive/{count:\\d+}", handleRecursive).Methods("GET")
	r.HandleFunc("/recursive/{count:\\d+}", handleRecursive).Methods("GET").Queries("index", "{\\d+}")

	http.Handle("/", r)
}

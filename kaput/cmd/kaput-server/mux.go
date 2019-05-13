package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var mainServeMux = initServeMux()

func initServeMux() *http.ServeMux {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/help", handleHelp).Methods("GET")

	r.HandleFunc("/crash", handleCrash).Methods("GET")
	r.HandleFunc("/fail", handleFail).Methods("GET")

	http.Handle("/", r)
	return http.DefaultServeMux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	url := "/help"
	http.Redirect(w, r, url, http.StatusFound)
}

func handleNotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "not implemented")
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	handleNotImplemented(w, r)
}

func handleFail(w http.ResponseWriter, r *http.Request) {
	// runs into panic recovery of server
	// will not shutdown process
	log.Panic("Going Down")
}

func handleCrash(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Going Down")
}

package main

import (
	"fmt"
	"net/http"
)

var mainServeMux = initServeMux()

func initServeMux() *http.ServeMux {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/help", handleHelp)

	http.HandleFunc("/crash", handleCrash)

	return http.DefaultServeMux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	url := "/help"
	http.Redirect(w, r, url, http.StatusFound)
}

func handleNotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "not implemented")
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handleNotImplemented(w, r)
}

func handleCrash(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handleNotImplemented(w, r)
}

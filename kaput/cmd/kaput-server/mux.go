package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var mainServeMux = initServeMux()

func initServeMux() *http.ServeMux {
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
	return http.DefaultServeMux
}

func helpHandler(path string) func(http.ResponseWriter, *http.Request) {
	const missing = "no help yet"
	h := map[string]string{
		"/":          missing,
		"/crash":     missing,
		"/respone":   missing,
		"/recursive": missing,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, h[path])
	}
}

func handleCrash(w http.ResponseWriter, r *http.Request) {
	switch mux.Vars(r)["it"] {
	case "connection":
		log.Panic("Crash Connection")
	case "server":
		log.Fatal("Crash Server")
	}
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)["code"]
	code, err := strconv.Atoi(value)
	if err != nil { // should not happen
		log.Fatal(err)
	}
	w.WriteHeader(code)
}

// TODO: add timeout
func requestRecursive(url string) chan error {
	c := make(chan error)
	go func() {
		r, err := http.Get(url)
		if err != nil {
			c <- err
			return
		}
		if r.StatusCode != http.StatusOK {
			c <- errors.New("Unexpected Result: " + r.Status)
			return
		}
		c <- nil
	}()
	return c
}

// NOTE: this could fail on wrong env
func createRecursiveURL(count, index int) string {
	return fmt.Sprintf("%s://%s:%s/recursive/%d?index=%d",
		privateRoute.protocol, privateRoute.domain,
		privateRoute.port, count, index)
}

func handleRecursive(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)["count"]
	count, err := strconv.Atoi(value)
	if err != nil { // should not happen
		log.Fatal(err)
	}

	value = r.URL.Query().Get("index")
	index, err := strconv.Atoi(value)
	if err != nil {
		// defaults to 1
		index = 1
	}

	statusCode := http.StatusOK
	switch {
	case count >= 2*index+1: // inner node with two children
		left := requestRecursive(createRecursiveURL(count, 2*index))
		right := requestRecursive(createRecursiveURL(count, 2*index+1))
		for i := 0; i < 2; i++ {
			select {
			case err := <-left:
				if err != nil {
					log.Println(err)
					statusCode = http.StatusInternalServerError
				}
			case err := <-right:
				if err != nil {
					log.Println(err)
					statusCode = http.StatusInternalServerError
				}
			}
		}

	case count == 2*index: // inner node with one child
		left := requestRecursive(createRecursiveURL(count, 2*index))
		err := <-left
		if err != nil {
			log.Println(err)
			statusCode = http.StatusInternalServerError
		}
	}
	w.WriteHeader(statusCode)
}

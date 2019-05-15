package kaput

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// GetEnv return environment variable name if existing
// else return fallback.
func GetEnv(name string, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = fallback
	}
	return value
}

// Crash either tcp connection
// or server process
func handleCrash(w http.ResponseWriter, r *http.Request) {
	switch mux.Vars(r)["it"] {
	case "connection":
		log.Panic("Crash Connection")
	case "server":
		log.Fatal("Crash Server")
	}
}

// Return specified http response status code
func handleResponse(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)["code"]
	code, err := strconv.Atoi(value)
	if err != nil { // should not happen
		log.Panic(err)
	}
	w.WriteHeader(code)
}

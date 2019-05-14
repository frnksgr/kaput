package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/frnksgr/kaput/kaput/pkg/kaput"
)

// middleware doing simple request logging
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got request: %s %s %s \n", r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	address := fmt.Sprintf("0.0.0.0:%s", kaput.GetEnv("PORT", "8080"))

	fmt.Printf("Starting server on %s\n", address)
	log.Fatal(http.ListenAndServe(address, requestLogger(http.DefaultServeMux)))
}

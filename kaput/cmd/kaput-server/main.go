package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = 8080

func getEnv(name string, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = fallback
	}
	return value
}

// middleware doing simple request logging
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got request: %s %s %s \n", r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	address := fmt.Sprintf("0.0.0.0:%s", getEnv("PORT", "8080"))

	fmt.Printf("Starting server on %s\n", address)
	log.Fatal(http.ListenAndServe(address, requestLogger(mainServeMux)))
}

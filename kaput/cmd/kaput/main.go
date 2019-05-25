package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/frnksgr/kaput/kaput/pkg/config"
	"github.com/frnksgr/kaput/kaput/pkg/kaput"
)

var requestCount uint64

// request logging middleware on purpose of debugging
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				http.Error(w, fmt.Sprint(err),
					http.StatusInternalServerError)
				return
			}
			s := strings.ReplaceAll(
				strings.ReplaceAll(string(dump), "\r\n", "\n"),
				"\n", "\n  ")
			requestCount++
			fmt.Fprintf(os.Stderr, "kaput: %d\n  %s\n", requestCount, s)
			next.ServeHTTP(w, r)
		})
}

func init() {
	fmt.Println("Version: ", kaput.Version)
}

func main() {
	address := fmt.Sprintf("%s:%s",
		config.Data.Listening.Host,
		config.Data.Listening.Port)

	fmt.Printf("Starting server on %s\n", address)
	if _, ok := os.LookupEnv("DEBUG"); ok {
		log.Fatal(
			http.ListenAndServe(
				address,
				requestLogger(http.DefaultServeMux)))
	} else {
		log.Fatal(
			http.ListenAndServe(
				address, http.DefaultServeMux))
	}
}

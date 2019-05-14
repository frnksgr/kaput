package kaput

// Support recursive calling of application.
// Each request creates two new asynchrone request until
// specified request count is reached.
// Returns statuscode 500 if any sub request fails.

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const defaultProtocol = "http"
const defaultDomain = "localhost"
const defaultPort = "8080"

type privateRouteT = struct {
	protocol string
	domain   string
	port     string
}

var privateRoute privateRouteT

func initRecursive() {
	privateRoute = privateRouteT{
		GetEnv("PRIVATE_PROTOCOL", defaultProtocol),
		GetEnv("PRIVATE_DOMAIN", defaultDomain),
		GetEnv("PRIVATE_PORT", defaultPort),
	}
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
// TODO: to be fixed in initRecursive
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

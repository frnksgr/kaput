package kaput

// Support recursive calling of application.
// Each request creates two new asynchrone request until
// specified request count is reached.
// Returns statuscode 500 if any sub request fails.

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	defaultProtocol = "http"
	defaultDomain   = "localhost"
	defaultPort     = "8080"
)

type (
	privateRouteT = struct {
		protocol string
		domain   string
		port     string
	}
)

var (
	privateRoute privateRouteT
)

func initRecursive() {
	privateRoute = privateRouteT{
		GetEnv("PRIVATE_PROTOCOL", defaultProtocol),
		GetEnv("PRIVATE_DOMAIN", defaultDomain),
		GetEnv("PRIVATE_PORT", defaultPort),
	}
}

func requestRecursive(url string, payload string) chan error {
	c := make(chan error)
	go func() {
		var r *http.Response
		var err error

		if len(payload) == 0 {
			r, err = http.Get(url)
		} else {
			// NOTE: consider using application/x-sh
			r, err = http.Post(url, "text/plain", strings.NewReader(payload))
		}
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

func getPayload(r *http.Request) (string, error) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	return "", nil
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

	payload, err := getPayload(r)
	if err != nil {
		log.Println("Error reading payload: ", err)
	}

	statusCode := http.StatusOK
	switch {
	case count >= 2*index+1: // inner node with two children
		left := requestRecursive(createRecursiveURL(count, 2*index), payload)
		right := requestRecursive(createRecursiveURL(count, 2*index+1), payload)
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
		left := requestRecursive(createRecursiveURL(count, 2*index), payload)
		err := <-left
		if err != nil {
			log.Println(err)
			statusCode = http.StatusInternalServerError
		}
	}

	w.WriteHeader(statusCode)

	// execute payload synchronously
	if len(payload) != 0 {
		fmt.Printf("Index: %d Running payload...\n", index)
		cmd := exec.Command("/bin/sh", "-c", payload)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Index: %d Excuting payload failed\n", index)
			log.Println(err)
		}
	}
}

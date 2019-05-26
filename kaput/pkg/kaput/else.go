package kaput

import (
	"log"
	"net/http"
	"strconv"

	"github.com/frnksgr/kaput/kaput/pkg/help"
	"github.com/gorilla/mux"
)

const (
	helpCrash = `
/crash/[command] kill current connectio or whole process
Where command is:
    connection      kill current connection (with proper TCP shutdown on uncomplete HTTP response)
    server          kill current process (with SIGTERM)
`
	helpResponse = `
/response/[code] return any response code
Where command is:
    code            will be set as HTTP response status code
`
)

func init() {
	help.Add("/crash", helpCrash)
	help.Add("/response", helpResponse)
}

// Crash either tcp connection
// or server process
func handleCrash(w http.ResponseWriter, r *http.Request) {
	switch mux.Vars(r)["it"] {
	case "connection":
		log.Panic("Crashing Connection...")
	case "server":
		log.Fatal("Crashing Server...")
	}
}

// Return specified http response status code
func handleResponse(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)["code"]
	code, err := strconv.Atoi(value)
	if err != nil { // should not happen
		panic(err)
	}
	w.WriteHeader(code)
}

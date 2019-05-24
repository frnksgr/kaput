package kaput

import (
	"fmt"
	"net/http"

	"github.com/frnksgr/kaput/kaput/pkg/help"
	"github.com/frnksgr/kaput/kaput/pkg/load"
	"github.com/gorilla/mux"
)

const (
	helpRoot = `
/[command] breaking things
Where command is:
    crash           crash something
    response        return arbitrary HTTP response codes
    load            create recursive requests each executing a specified workload
	
Call command URL to get specific help on command.
`
)

func init() {
	fmt.Println("Initializing router ...")

	r := mux.NewRouter()
	r.HandleFunc("/", help.Handler("/")).Methods("GET")

	r.HandleFunc("/crash", help.Handler("/crash")).Methods("GET")
	r.HandleFunc("/crash/{it:(?:connection|server)}", handleCrash).Methods("GET")

	r.HandleFunc("/response", help.Handler("/response")).Methods("GET")
	r.HandleFunc("/response/{code:[12345][0-9]{2}}", handleResponse).Methods("GET")

	r.HandleFunc("/load", help.Handler("/load")).Methods("GET")
	r.HandleFunc("/load/{count:\\d{0,4}}", load.PostHandler).Methods("POST").
		HeadersRegexp("Content-Type", "application/(x-www-form-urlencoded|json)")

	http.Handle("/", r)

	help.Add("/", helpRoot)
}

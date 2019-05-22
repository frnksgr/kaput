package help

import (
	"fmt"
	"net/http"
)

type (
	command struct {
		Name string
		Text string
	}
	help struct {
		Header   string
		Commands []command
	}
)

var (
	helpMap = make(map[string]string)
)

func init() {
	fmt.Println("Initializing help system ...")
}

// Add add help strinbg for specified path
func Add(path string, help string) {
	helpMap[path] = help
}

// Handler serving help URLs
func Handler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		help := helpMap[path]
		if len(help) == 0 {
			help = fmt.Sprintln(path, " no help available")
		}
		fmt.Fprintln(w, help)
	}
}

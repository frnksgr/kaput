package kaput

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

func addHelp(path string, item string) {
	helpMap[path] = item
}

// TODO: provide help; use templating
func helpHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		help := helpMap[path]
		if len(help) == 0 {
			help = fmt.Sprintln(path, " no help available")
		}
		fmt.Fprintln(w, help)
	}
}

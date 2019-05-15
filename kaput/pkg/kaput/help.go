package kaput

import (
	"fmt"
	"net/http"
	"text/template"
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
	helpMap         = make(map[string]*help)
	helpTemplate, _ = template.New("help-template").Parse(helpTemplateRaw)
)

const (
	helpTemplateRaw = `{{.Header}}
    {{range .Commands -}}
    {{.Name}}		{{.Text}}
    {{end}}
`
)

func addHelp(path string, item *help) {
	helpMap[path] = item
}

// TODO: provide help; use templating
func helpHandler(path string) func(http.ResponseWriter, *http.Request) {
	const missing = "no help yet"

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, path, " ")
		data := helpMap[path]
		if data == nil {
			fmt.Fprintln(w, missing)
		} else {
			helpTemplate.Execute(w, helpMap[path])
		}
	}
}

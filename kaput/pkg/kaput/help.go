package kaput

import (
	"fmt"
	"net/http"
)

// TODO: provide help; use templating
func helpHandler(path string) func(http.ResponseWriter, *http.Request) {
	const missing = "no help yet"
	h := map[string]string{
		"/":          missing,
		"/crash":     missing,
		"/respone":   missing,
		"/recursive": missing,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, h[path])
	}
}

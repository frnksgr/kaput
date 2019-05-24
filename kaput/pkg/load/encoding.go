package load

import (
	"encoding/json"
	"io"
	"net/http"
)

// Task describes a specific task.
// It is used by JSON encodings to call
// TaskExec generators
// E.g. '{"cpu", 75.0}' as json can be used to call
// cpuTask(70.5) to create an executable TaskExec
type Task []interface{}

func generate(task Task) TaskExec {
	// we rely on panic in case we fail.
	switch task[0].(string) { // select generator
	case "none":
		return noneTask()
	case "cpu":
		return cpuTask(task[1].(float32))
	case "ram":
		return ramTask(task[1].(uint64))
	default:
		panic("unknown task type")
	}
}

// data is application/x-www-form-urlencoded
// do some value checks
func decodeLoad(r *http.Request) (*Load, error) {
	var l Load
	r.ParseForm()
	return &l, nil
}

// data is application/json
// assume propper values
func decodeNode(r io.Reader) (*Node, error) {
	var n Node
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&n); err != nil {
		return nil, err
	}
	return &n, nil
}

func encodeNode(n *Node, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(n)
}

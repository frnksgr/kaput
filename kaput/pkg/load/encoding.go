package load

import (
	"encoding/json"
	"io"
)

// Task describes a specific task.
// It is used by JSON encodings to call
// TaskExec generators
// E.g. '{"cpu", 75.0}' as json can be used to call
// cpuTask(70.5) to create an executable TaskExec
type Task []interface{}

func generate(task Task) taskExec {
	// we rely on panic in case we fail.
	switch task[0].(string) { // select generator
	case "none":
		return noneTask()
	case "cpu":
		return cpuTask(float32(task[1].(float64)))
	case "ram":
		return ramTask(uint64(task[1].(float64)))
	default:
		panic("unknown task type")
	}
}

// decode initial load on node 1
// do some value checks
func decodeLoad(r io.Reader) (*load, error) {
	var l load
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

// data is application/json
// assume propper values
func decodeNode(r io.Reader) (*node, error) {
	var n node
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&n); err != nil {
		return nil, err
	}
	return &n, nil
}

func encodeNode(n *node, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(n)
}

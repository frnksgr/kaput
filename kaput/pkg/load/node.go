package load

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Node is associated with a set of tasks beeing executed
// for a specific time
type node struct {
	Count   int        `json:"count,omitempty"`
	Index   int        `json:"index,omitempty"`
	Level   int        `json:"level,omitempty"`
	Timeout int64      `json:"timeout"` // timeout in ms (0 is no timeout)
	Tasks   []taskSpec `json:"tasks"`
}

// each node provides a logger to identify logs with this specific node
func (n *node) logger() *log.Logger {
	prefix := fmt.Sprintf("ID: %d ", n.Index)
	return log.New(os.Stdout, prefix, log.Lmicroseconds)
}

// create a fully balanced binary tree
func (n *node) child(which int) *node {
	return &node{
		Count:   n.Count,
		Index:   n.Index + 1<<(uint(n.Level+which)),
		Level:   n.Level + 1,
		Timeout: n.Timeout,
		Tasks:   n.Tasks,
	}
}

// Run specified workload on this Node
func (n *node) run(logger *log.Logger) {
	for _, task := range n.Tasks {
		exec := generate(task)
		go exec(time.After(
			time.Duration(n.Timeout)*time.Millisecond),
			logger)
	}
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

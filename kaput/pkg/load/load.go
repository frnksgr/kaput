package load

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Load represent a set of tasks that can be executed
// asynchrounously for a specific time
type load struct {
	Timeout int64  `json:"timeout"` // timeout in ms (0 is no timeout)
	Tasks   []Task `json:"tasks"`   // tasks to be run
}

// Node metadata of a specific Request
type node struct {
	Count int   `json:"count"`
	Index int   `json:"index"`
	Level int   `json:"level"`
	Load  *load `json:"load"`
}

func (n *node) logger() *log.Logger {
	// logger used for this specific node
	prefix := fmt.Sprintf("ID: %d ", n.Index)
	return log.New(os.Stdout, prefix, log.Lmicroseconds)
}

const (
	childLeft = iota
	childRight
)

// create a fully balanced tree
func (n *node) child(which int) *node {
	return &node{
		Count: n.Count,
		Index: n.Index + 1<<(uint(n.Level+which)),
		Level: n.Level + 1,
		Load:  n.Load,
	}
}

// Run specified workload on this Node
func (n *node) run(logger *log.Logger) {
	l := n.Load
	noTimeout := make(chan time.Time)
	timeout := time.Duration(l.Timeout) * time.Millisecond

	for _, task := range l.Tasks {
		exec := generate(task)
		if l.Timeout >= 0 {
			go exec(time.After(timeout), logger)
		} else {
			go exec(noTimeout, logger)
		}
	}
}

package load

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Load represent a set of tasks that can be executed
// asynchrounously for a specific time
type Load struct {
	Timeout uint   `json:"timeout"` // timeout in ms (0 is no timeout)
	Delay   uint   `json:"delay"`   // delay in ms between triggering tasks
	Tasks   []Task `json:"tasks"`   // tasks to be run
}

// Node metadata of a specific Request
type Node struct {
	Index int   `json:"index"`
	Load  *Load `json:"load"`
}

func (n *Node) logger() *log.Logger {
	// logger used for this specific node
	prefix := fmt.Sprintf("ID: %d ", n.Index)
	return log.New(os.Stdout, prefix, log.Lmicroseconds)
}

const (
	childLeft = iota
	childRight
)

// create a fully balanced tree
func (n *Node) child(which int) *Node {
	// get depth of node in binary tree
	// starting with node 1 on level 0
	// node 2 and 3 on level 1
	// node 4,5, ...7 on level 2
	// node 2^x, ... 2^(x+1)-1 on level x
	level := 31 // chould never be reached index should always be below 9999
	switch ui := uint32(n.Index); ui {
	case 1:
		level = 0
	default:
		for i := uint(1); i < 32; i++ {
			if ui>>i == 0 {
				level = int(i - 1)
				break
			}
		}
	}
	return &Node{
		Index: n.Index + 2 ^ (level + which),
		Load:  n.Load,
	}
}

// Run specified workload on this Node
func (n *Node) Run(logger *log.Logger) {
	l := n.Load
	noTimeout := make(chan time.Time)
	delay := time.Duration(l.Delay) * time.Millisecond
	timeout := time.Duration(l.Timeout) * time.Millisecond

	for _, task := range l.Tasks {
		exec := generate(task)

		if l.Delay > 0 {
			time.Sleep(delay)
		}
		if l.Timeout > 0 {
			go exec(time.After(timeout), logger)
		} else {
			go exec(noTimeout, logger)
		}
	}
}

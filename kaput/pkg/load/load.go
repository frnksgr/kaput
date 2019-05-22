package load

import (
	"fmt"
	"log"
	"time"
)

// TaskType a specific type of task
type TaskType string

// Task describes a specific task
// Used to instantiate a specific work load
type Task struct {
	Type      TaskType
	Paramters []interface{}
}

// Load represent a set of tasks that can be executed
// asynchrounously for a specific time
type Load struct {
	Timeout uint   // timeout in ms (0 is no timeout)
	Delay   uint   // delay in ms between triggering tasks
	Tasks   []Task // tasks to be run
}

// Node metadata of a specific Request
type Node struct {
	Index int
	Load  *Load
}

type taskExec func(timeout <-chan time.Time)

// Printf prefixed fmt.Printf
func (n *Node) Printf(format string, a ...interface{}) error {
	fmt.Printf("ID: %d ", n.Index) // Prefix
	_, err := fmt.Printf(format, a...)
	return err
}

// Println prefixed fmt.Println
func (n *Node) Println(a ...interface{}) error {
	fmt.Printf("ID: %d ", n.Index) // Prefix
	_, err := fmt.Println(a...)
	return err
}

// Run specified workload on this Node
func (n *Node) Run() error {
	l := n.Load
	noTimeout := make(chan time.Time)
	var exec taskExec
	for _, task := range l.Tasks {
		switch task.Type {
		case "cpu":
			percent, ok := task.Paramters[0].(float64)
			if !ok {
				return fmt.Errorf("Error on type assertion for task %s", task.Type)
			}
			exec = cpuTask(percent)
		case "ram":
			bytes, ok := task.Paramters[0].(uint64)
			if !ok {
				return fmt.Errorf("Error on type assertion for task %s", task.Type)
			}
			exec = ramTask(bytes)
		default:
			return fmt.Errorf("Unknown task type: %s", task.Type)
		}

		if l.Delay > 0 {
			time.Sleep(time.Duration(l.Delay) * time.Millisecond)
		}

		n.Printf("executing %s task for %d ms\n", task.Type, l.Timeout)

		if l.Timeout > 0 {
			go exec(time.After(time.Duration(l.Timeout) * time.Millisecond))
		} else {
			go exec(noTimeout)
		}
	}
	return nil
}

func dummyTask(timeout <-chan time.Time) {
	<-timeout
}

func cpuTask(percent float64) taskExec {
	return func(timeout <-chan time.Time) {
		const count = 10 * 1000 * 1000 // should run for a view ms
		for {
			select {
			case <-timeout:
				return
			default:
				x, y := 0, 1
				t := time.Now()
				for i := 1; i <= count; i++ {
					x, y = y, x
				}
				d := time.Since(t)
				time.Sleep(time.Duration((1.0 - percent) * float64(d)))
			}
		}
	}
}

func ramTask(bytes uint64) taskExec {
	return func(timeout <-chan time.Time) {
		buf := make([]byte, bytes)
		for {
			select {
			case <-timeout:
				return
			default:
				// touch each page every 5 ms to stay in memory
				for i := uint64(0); i < bytes; i += 4096 {
					buf[i] = 1
				}
				time.Sleep(5 * time.Millisecond)
			}
		}
	}
}

// DiskIOTask Disk IO bound task
func diskIOTask(io uint64) taskExec {
	log.Panic("Not implemented")
	return dummyTask
}

// NetIOTask Network IO bound task
func netIOTask(bytes uint64) taskExec {
	log.Panic("Not implemented")
	return dummyTask
}

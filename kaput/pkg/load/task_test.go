package load

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func exec(exec taskExec, duration time.Duration) {
	timeout := make(chan time.Time)
	logger := log.New(os.Stdout, "", 0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		exec(timeout, logger)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(duration)
		timeout <- time.Now()
	}()
	wg.Wait()
}

func NoneTask() {
	exec(noneTask(), 0)
}

func ExampleNoneTask() {
	NoneTask()
	// Output:
	// Starting none task
	// Ending none task
}

func CPUTask(percent float32) {
	exec(cpuTask(percent), 0)
}

func ExampleCPUTask() {
	CPUTask(0.5) // 50%
	// Output:
	// Starting cpu task
	// Ending cpu task
}

func RAMTask(bytes uint64) {
	exec(ramTask(bytes), 0)
}

func ExampleRAMTask() {
	RAMTask(10 * 1024 * 1024) // 10 MB
	// Output:
	// Starting ram task
	// Ending ram task
}

func BenchmarkCPUloop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// this should run for about one millisecond
		cpuloop(1000)
	}
}

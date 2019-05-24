package load

import (
	"log"
	"time"
)

// TaskExec is interuptable task that executes until interupted
type TaskExec func(timeout <-chan time.Time, logger *log.Logger)

func noneTask() TaskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting none task")
		<-timeout
		logger.Println("Ending none task")
	}
}

func cpuTask(percent float32) TaskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting cpu task")
		count := 10 * 1000 * 1000 // should run for a view ms
		for done := false; !done; {
			select {
			case <-timeout:
				done = true
			default:
				x, y := 0, 1
				t := time.Now()
				// loop canot be 'preemted' by go scheduler
				// time slice of OS scheduler is about 100 ms
				for i := 1; i < count; i++ {
					x, y = y, x
				}
				d := time.Since(t)
				time.Sleep(time.Duration((1.0 - percent) * float32(d)))
			}
		}
		logger.Println("Ending cpu task")
	}
}

func ramTask(bytes uint64) TaskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting ram task")
		buf := make([]byte, bytes)
		for done := false; !done; {
			select {
			case <-timeout:
				done = true
			default:
				// touch each page every 5 ms to stay in memory
				for i := uint64(0); i < bytes; i += 4096 {
					buf[i] = 1
				}
				time.Sleep(5 * time.Millisecond)
			}
		}
		logger.Println("Ending ram task")
	}
}

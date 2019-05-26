package load

import (
	"log"
	"syscall"
	"time"
)

// TaskExec is interuptable task that executes until interupted
type taskExec func(timeout <-chan time.Time, logger *log.Logger)

func noneTask() taskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting none task")
		<-timeout
		logger.Println("Ending none task")
	}
}

func cpuTask(percent float32) taskExec {
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

// allocate memory using syscall mmap anonymous
// to workaround go memory management
func alloc(bytes int) ([]byte, error) {
	inc := syscall.Getpagesize() / 2
	mem, err := syscall.Mmap(-1, 0, bytes,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANONYMOUS|syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}
	// cause Pagefaults to increase RSS
	for i := 0; i < bytes; i += inc {
		mem[i] = 1
	}
	return mem, nil
}

func free(mem []byte) error {
	return syscall.Munmap(mem)
}

func ramTask(bytes uint64) taskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting ram task")
		mem, err := alloc(int(bytes))
		if err != nil {
			logger.Panic(err)
		}
		<-timeout
		free(mem)
		logger.Println("Ending ram task")
	}
}

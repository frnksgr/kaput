package load

import (
	"log"
	"runtime"
	"syscall"
	"time"
)

// taskExec create some load until done channel is triggered
// output data via logger
type taskExec func(timeout <-chan time.Time, logger *log.Logger)

// taskSpec describes a specific task that can be encoded in json
// E.g. '{"cpu", 75.0}' as json can be used to call
// cpuTask(70.5) to create an executable TaskExec
type taskSpec []interface{}

func noneTask() taskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting none task")
		<-timeout
		logger.Println("Ending none task")
	}
}

// cpuloopcount need to be calibrated
// to cpuloop run about us microseconds
var cpuloopcount = 4 * 1000 * 1000

func cpuloop(us int) {
	count := (us * cpuloopcount) / 1000
	x, y := 0, 1
	for i := 1; i < count; i++ {
		x, y = y, x
	}
}

// calibrate iterations so cpuloop(1000) runs for about 1 millisecond
func calibrate() {
	runtime.LockOSThread()
	for i := 0; i < 3; i++ {
		start := time.Now()
		cpuloop(1000)
		since := time.Since(start)
		ratio := float64(since) / float64(time.Millisecond)
		cpuloopcount = int(float64(cpuloopcount) / ratio)
	}
	runtime.UnlockOSThread()
}

func init() {
	calibrate()
}

func cpuTask(percent float32) taskExec {
	return func(timeout <-chan time.Time, logger *log.Logger) {
		logger.Println("Starting cpu task")
		for done := false; !done; {
			select {
			case <-timeout:
				done = true
			default: // this should run about 10 ms
				t := time.Now()
				// busy
				cpuloop(int(percent * 10000.0))
				// part
				time.Sleep(
					time.Duration(
						(1.0 - percent) * float32(time.Since(t))))
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

func generate(spec taskSpec) taskExec {
	// we rely on panic in case we fail.
	switch spec[0].(string) { // select generator
	case "none":
		return noneTask()
	case "cpu":
		return cpuTask(float32(spec[1].(float64)))
	case "ram":
		return ramTask(uint64(spec[1].(float64)))
	default:
		panic("unknown task type")
	}
}

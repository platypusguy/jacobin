package prof

import (
	"fmt"
	"jacobin/src/trace"
	"os"
	"runtime/pprof"
	"sync"
)

var (
	mu        sync.Mutex
	profFile  *os.File
	profiling bool
)

// Start starts CPU profiling and registers the file so Exit can stop it.
func StartProfiling(path string) {
	mu.Lock()
	if profiling {
		return
	}
	trace.Trace(fmt.Sprintf("prof.StartProfiling: Starting CPU profile in %s", path))
	ff, err := os.Create(path)
	if err != nil {
		mu.Unlock()
		trace.Error(fmt.Sprintf("prof.StartProfiling: os.Create(%s) failed, err: %s", path, err.Error()))
		ExitToOS(1)
	}
	err = pprof.StartCPUProfile(ff)
	if err != nil {
		_ = ff.Close()
		mu.Unlock()
		trace.Error(fmt.Sprintf("prof.StartProfiling: Could not create CPU profile in %s, err: %s", path, err.Error()))
		ExitToOS(1)
	}
	profFile = ff
	profiling = true
	mu.Unlock()
}

// Stop stops profiling if it was started.
func StopProfiling() {
	mu.Lock()
	defer mu.Unlock()
	if !profiling {
		return
	}
	pprof.StopCPUProfile()
	_ = profFile.Sync()
	_ = profFile.Close()
	profFile = nil
	profiling = false
}

// Exit ensures profiling is stopped and then exits with code.
func ExitToOS(code int) {
	// best-effort stop so cpu.out is flushed
	StopProfiling()
	os.Exit(code)
}

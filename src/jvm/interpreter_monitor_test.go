package jvm

import (
	"sync"
	"testing"
	"time"

	"jacobin/src/frames"
	"jacobin/src/object"
)

// Test doMonitorEnter and doMonitorExit with 8 independent threads (frames)
// competing for the same object monitor. Thread 1 acquires first and holds
// briefly; threads 2..8 must block until release, then each should acquire
// and release successfully exactly once.
func TestDoMonitorEnterExit_EightFrames_ContentionAndHandoff(t *testing.T) {
	// Shared object to synchronize on
	obj := object.MakeEmptyObject()

	// Create 8 independent frames with unique thread IDs
	const total = 8
	framesArr := make([]*frames.Frame, total)
	for i := 0; i < total; i++ {
		fr := frames.CreateFrame(8)
		fr.Thread = i + 1 // threads 1..8
		fr.ClName = "LTest;"
		fr.MethName = "test"
		fr.MethType = "()V"
		framesArr[i] = fr
	}

	// Channels to coordinate and record acquisitions (for contenders only)
	acquiredCh := make(chan int, total-1)
	var wg sync.WaitGroup
	wg.Add(total - 1)              // contenders only (threads 2..8)
	startCh := make(chan struct{}) // gate contenders to start after the holder acquires

	// Thread 1 (holder) acquires and holds briefly
	holder := framesArr[0]
	push(holder, obj)
	if got := doMonitorEnter(holder, 0); got != 1 {
		t.Fatalf("holder thread: doMonitorEnter returned %d, want 1", got)
	}

	// Start contenders: threads 2..8 (after holder has acquired)
	for i := 1; i < total; i++ {
		idx := i
		go func() {
			defer wg.Done()
			// Wait until we're signaled to start contending
			<-startCh
			// Push object reference and attempt to enter
			push(framesArr[idx], obj)
			got := doMonitorEnter(framesArr[idx], 0)
			if got != 1 {
				t.Fatalf("thread %d: doMonitorEnter returned %d, want 1", framesArr[idx].Thread, got)
			}
			t.Logf("thread %d acquired object-lock, will sleep 0.5s", framesArr[idx].Thread)
			time.Sleep(500 * time.Millisecond)
			// Record that this thread has acquired the monitor
			acquiredCh <- framesArr[idx].Thread

			// Now exit: push object again (enter popped it) and call exit
			push(framesArr[idx], obj)
			got = doMonitorExit(framesArr[idx], 0)
			if got != 1 {
				t.Fatalf("thread %d: doMonitorExit returned %d, want 1", framesArr[idx].Thread, got)
			}
			t.Logf("thread %d released object-lock", framesArr[idx].Thread)
		}()
	}

	// Release the gate so contenders begin attempting to acquire
	close(startCh)

	// Ensure no contender acquires while the holder still owns the monitor
	select {
	case tid := <-acquiredCh:
		t.Fatalf("contender %d acquired monitor while holder still owns it", tid)
	case <-time.After(20 * time.Millisecond):
		// expected: none acquired yet
	}

	// Release the monitor from the holder to allow contenders to proceed
	// Short sleep to widen the window ensuring contenders are already waiting
	time.Sleep(50 * time.Millisecond)
	push(holder, obj)
	if got := doMonitorExit(holder, 0); got != 1 {
		t.Fatalf("holder thread: doMonitorExit returned %d, want 1", got)
	}

	// Expect all 7 contenders to acquire eventually
	deadline := time.After(10 * time.Second)
	acquired := 0
	for acquired < total-1 {
		select {
		case <-acquiredCh:
			acquired++
		case <-deadline:
			t.Fatalf("timeout: only %d/%d contenders acquired after release", acquired, total-1)
		}
	}

	// Wait for all contenders to finish their exits
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
		// ok
	case <-time.After(10 * time.Second):
		t.Fatalf("timeout waiting for contenders to finish exits")
	}
}

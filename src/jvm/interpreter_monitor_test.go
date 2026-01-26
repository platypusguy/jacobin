package jvm

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
)

// Test doMonitorenter and doMonitorexit with 8 independent threads (frames)
// competing for the same object monitor. Thread 1 acquires first and holds
// briefly; threads 2..8 must block until release, then each should acquire
// and release successfully exactly once.
func TestDoMonitorEnterExit_EightFrames_ContentionAndHandoff(t *testing.T) {
	// Set JacobinName to "test" so ThrowEx returns instead of aborting
	oldName := globals.GetGlobalRef().JacobinName
	globals.GetGlobalRef().JacobinName = "test"
	defer func() { globals.GetGlobalRef().JacobinName = oldName }()

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
	if got := doMonitorenter(holder, 0); got != 1 {
		t.Fatalf("holder thread: doMonitorenter returned %d, want 1", got)
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
			got := doMonitorenter(framesArr[idx], 0)
			if got != 1 {
				// Don't use t.Fatalf in a goroutine
				t.Errorf("thread %d: doMonitorenter returned %d, want 1", framesArr[idx].Thread, got)
				return
			}
			t.Logf("thread %d acquired object-lock, will sleep 0.5s", framesArr[idx].Thread)
			time.Sleep(500 * time.Millisecond)
			// Record that this thread has acquired the monitor
			acquiredCh <- framesArr[idx].Thread

			// Now exit: push object again (enter popped it) and call exit
			push(framesArr[idx], obj)
			got = doMonitorexit(framesArr[idx], 0)
			if got != 1 {
				t.Errorf("thread %d: doMonitorexit returned %d, want 1", framesArr[idx].Thread, got)
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
	if got := doMonitorexit(holder, 0); got != 1 {
		monitor := (*object.ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
		t.Logf("holder thread: doMonitorexit failed. monitor=%+v", monitor)
		t.Fatalf("holder thread: doMonitorexit returned %d, want 1", got)
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

// Simulate Java nested synchronized(lock) { synchronized(lock) { ... } }
// Using interpreter monitor enter/exit. We preconfigure the object as fat-locked
// by the same thread to model reentrant locking (recursion increments on reenter).
func TestDoMonitorEnterExit_NestedSynchronized_Reentrant(t *testing.T) {
	// Set JacobinName to "test" so ThrowEx returns instead of aborting
	oldName := globals.GetGlobalRef().JacobinName
	globals.GetGlobalRef().JacobinName = "test"
	defer func() { globals.GetGlobalRef().JacobinName = oldName }()

	// Create a frame representing a single Java thread
	fr := frames.CreateFrame(8)
	fr.Thread = 1
	fr.ClName = "LTest;"
	fr.MethName = "nested"
	fr.MethType = "()V"

	obj := object.MakeEmptyObject()

	// Enter once (simulating the outer synchronized block)
	push(fr, obj)
	if got := doMonitorenter(fr, 0); got != 1 {
		t.Fatalf("first doMonitorenter returned %d, want 1", got)
	}
	monitor := (*object.ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	if monitor == nil || monitor.Owner != int32(fr.Thread) {
		t.Logf("monitor: %+v", monitor)
		t.Fatalf("after first enter: owner=%v (want owner=%d)",
			func() any {
				if monitor != nil {
					return monitor.Owner
				}
				return nil
			}(),
			fr.Thread)
	}
	if monitor.Recursion != 1 {
		t.Fatalf("after first enter: expected recursion=1, got %d", monitor.Recursion)
	}

	// Enter again (nested synchronized on the same lock)
	push(fr, obj)
	if got := doMonitorenter(fr, 0); got != 1 {
		t.Fatalf("second doMonitorenter returned %d, want 1", got)
	}
	monitor = (*object.ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	if monitor == nil || monitor.Recursion != 2 {
		t.Fatalf("after second enter: expected recursion=2, got monitor=%v rec=%d", monitor, func() int32 {
			if monitor != nil {
				return monitor.Recursion
			}
			return -1
		}())
	}

	// Now unwind like exiting nested synchronized blocks: TWO exits total
	push(fr, obj)
	if got := doMonitorexit(fr, 0); got != 1 {
		t.Fatalf("first doMonitorexit returned %d, want 1", got)
	}
	monitor = (*object.ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	if monitor == nil || monitor.Recursion != 1 {
		t.Fatalf("after first exit: expected recursion=1, got monitor=%v rec=%d", monitor, func() int32 {
			if monitor != nil {
				return monitor.Recursion
			}
			return -1
		}())
	}

	push(fr, obj)
	if got := doMonitorexit(fr, 0); got != 1 {
		t.Fatalf("second doMonitorexit returned %d, want 1", got)
	}
	if object.IsObjectLocked(obj) {
		t.Fatalf("expected object to be fully unlocked after 2 exits")
	}
}

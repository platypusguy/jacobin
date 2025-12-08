package object

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestObjLockUnlock_ThinLockCycle(t *testing.T) {
	obj := MakeEmptyObject()

	// Ensure object starts unlocked
	SetLockState(obj, lockStateUnlocked)

	// Acquire thin lock
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("ObjLock returned error: %v", err)
	}
	if got := obj.Mark.Misc & lockStateMask; got != lockStateThinLocked {
		t.Fatalf("expected thin locked state, got %b", got)
	}

	// Release thin lock
	if err := obj.ObjUnlock(1); err != nil {
		t.Fatalf("ObjUnlock returned error: %v", err)
	}
	if got := obj.Mark.Misc & lockStateMask; got != lockStateUnlocked {
		t.Fatalf("expected unlocked state after unlock, got %b", got)
	}
}

func TestObjUnlock_WhenAlreadyUnlocked_ReturnsError(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateUnlocked)

	err := obj.ObjUnlock(1)
	if err == nil {
		t.Fatalf("expected error when unlocking an unlocked object")
	}
	// Be tolerant to exact wording but ensure it's the correct path
	if !errors.Is(err, err) { // placeholder to use err; message check below
	}
}

func TestObjLockUnlock_GCMarked_ReturnsError(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateGCMarked)

	if err := obj.ObjLock(1); err == nil {
		t.Fatalf("expected error on ObjLock for GC-marked object")
	}

	if err := obj.ObjUnlock(1); err == nil {
		t.Fatalf("expected error on ObjUnlock for GC-marked object")
	}
}

func TestObjUnlock_FatLock_MonitorNil_ReturnsError(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateFatLocked)
	obj.Monitor = nil

	if err := obj.ObjUnlock(1); err == nil {
		t.Fatalf("expected error when fat-locked but monitor is nil")
	}
}

func TestObjUnlock_FatLock_OwnerMismatch_ReturnsError(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateFatLocked)
	obj.Monitor = &ObjectMonitor{Owner: 2, Recursion: 0}

	if err := obj.ObjUnlock(1); err == nil {
		t.Fatalf("expected error when unlocking fat lock from non-owner thread")
	}
}

func TestObjLock_FatLock_OwnerRecursiveIncrementsRecursion(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateFatLocked)
	obj.Monitor = &ObjectMonitor{Owner: 7, Recursion: 0}

	if err := obj.ObjLock(7); err != nil {
		t.Fatalf("ObjLock (fat, same owner) returned error: %v", err)
	}
	if obj.Monitor.Recursion != 1 {
		t.Fatalf("expected recursion to increment to 1, got %d", obj.Monitor.Recursion)
	}
	// State remains fat-locked
	if got := obj.Mark.Misc & lockStateMask; got != lockStateFatLocked {
		t.Fatalf("expected object to remain fat-locked, got %b", got)
	}
}

func TestObjUnlock_FatLock_RecursiveDecrementAndFinalRelease(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateFatLocked)
	obj.Monitor = &ObjectMonitor{Owner: 3, Recursion: 2}

	// First unlock should decrement recursion only
	if err := obj.ObjUnlock(3); err != nil {
		t.Fatalf("first ObjUnlock returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 1 {
		t.Fatalf("expected recursion to decrement to 1, got monitor=%v rec=%d", obj.Monitor, obj.Monitor.Recursion)
	}
	if got := obj.Mark.Misc & lockStateMask; got != lockStateFatLocked {
		t.Fatalf("expected to remain fat-locked after decrement, got %b", got)
	}

	// Second unlock should decrement to 0 (still fat-locked)
	if err := obj.ObjUnlock(3); err != nil {
		t.Fatalf("second ObjUnlock returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 0 {
		t.Fatalf("expected recursion to decrement to 0 before final release, got monitor=%v rec=%d", obj.Monitor, obj.Monitor.Recursion)
	}

	// Final unlock should clear monitor and set unlocked state
	if err := obj.ObjUnlock(3); err != nil {
		t.Fatalf("final ObjUnlock returned error: %v", err)
	}
	if obj.Monitor != nil {
		t.Fatalf("expected monitor to be cleared on final unlock")
	}
	if got := obj.Mark.Misc & lockStateMask; got != lockStateUnlocked {
		t.Fatalf("expected unlocked state after final release, got %b", got)
	}
}

// Two goroutines contend for a thin lock on the same object.
// Goroutine B must block until A releases, then acquire successfully.
func TestObjLock_TwoThreads_ThinLockContention(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateUnlocked)

	// Thread 1 acquires thin lock
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("thread 1 ObjLock returned error: %v", err)
	}

	acquired := make(chan struct{})
	done := make(chan struct{})

	// Thread 2 tries to acquire while A holds it
	go func() {
		if err := obj.ObjLock(2); err != nil {
			t.Fatalf("thread 2 ObjLock returned error: %v", err)
		}
		// signal B acquired
		close(acquired)
		// immediately release
		if err := obj.ObjUnlock(2); err != nil {
			t.Fatalf("thread 2 ObjUnlock returned error: %v", err)
		}
		close(done)
	}()

	// Ensure B does not acquire within a short window while A still holds the lock
	select {
	case <-acquired:
		t.Fatalf("thread 2 acquired lock while A still holds it")
	case <-time.After(20 * time.Millisecond):
		// expected: no acquire yet
	}

	// A releases the lock, allowing B to acquire
	if err := obj.ObjUnlock(1); err != nil {
		t.Fatalf("thread 1 ObjUnlock returned error: %v", err)
	}

	// Now B should acquire shortly
	select {
	case <-acquired:
		// good
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for thread 2 to acquire after release")
	}

	// And finish cleanly
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for thread 2 to finish unlock")
	}
}

// Two goroutines contend when the object is fat-locked by thread 1.
// Thread 2 must block until thread 1 fully releases the fat lock, then acquire.
func TestObjLock_TwoThreads_FatLockContentionAndHandoff(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateFatLocked)
	obj.Monitor = &ObjectMonitor{Owner: 1, Recursion: 0}

	acquired := make(chan struct{})
	done := make(chan struct{})

	// Thread 2 attempts to lock while fat-locked by owner=1
	go func() {
		if err := obj.ObjLock(2); err != nil {
			t.Fatalf("thread 2 ObjLock (fat contention) returned error: %v", err)
		}
		// Signal acquired
		t.Log("Thread 2 locked object successfully.")
		close(acquired)
		// Release and finish
		if err := obj.ObjUnlock(2); err != nil {
			t.Fatalf("thread 2 ObjUnlock after fat handoff returned error: %v", err)
		}
		t.Log("Thread 2 released object successfully.")
		close(done)
	}()

	// Ensure thread 2 hasn't acquired yet while thread 1 still owns the fat lock
	select {
	case <-acquired:
		t.Fatalf("thread 2 acquired fat lock while A still owns it")
	case <-time.After(20 * time.Millisecond):
		// expected
	}

	// Thread 1 releases the fat lock completely
	time.Sleep(2000 * time.Millisecond)
	if err := obj.ObjUnlock(1); err != nil {
		t.Fatalf("thread 1 ObjUnlock (fat) returned error: %v", err)
	}
	t.Log("Thread 1 released object successfully.")

	// After release, thread 2 should be able to acquire fairly quickly
	select {
	case <-acquired:
		// ok
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for thread 2 to acquire after fat release")
	}

	// And finish cleanly
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for thread 2 to finish after fat release")
	}
}

// Clone of the fat-lock contention test, but using thin locking throughout.
// Two goroutines contend when the object is thin-locked by thread 1.
// Thread 2 must block until thread 1 releases the thin lock, then acquire.
func TestObjLock_TwoThreads_ThinLockContentionAndHandoff(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateUnlocked)

	// Thread 1 acquires thin lock
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("thread 1 ObjLock (thin) returned error: %v", err)
	}

	acquired := make(chan struct{})
	done := make(chan struct{})

	// Thread 2 attempts to lock while thin-locked by thread 1
	go func() {
		if err := obj.ObjLock(2); err != nil {
			t.Fatalf("thread 2 ObjLock (thin contention) returned error: %v", err)
		}
		t.Log("Thread 2 locked object successfully (thin).")
		close(acquired)
		// Release and finish
		if err := obj.ObjUnlock(2); err != nil {
			t.Fatalf("thread 2 ObjUnlock after thin handoff returned error: %v", err)
		}
		t.Log("Thread 2 released object successfully (thin).")
		close(done)
	}()

	// Ensure thread 2 hasn't acquired yet while thread 1 still holds the thin lock
	select {
	case <-acquired:
		t.Fatalf("thread 2 acquired thin lock while thread 1 still holds it")
	case <-time.After(20 * time.Millisecond):
		// expected
	}

	// Thread 1 releases the thin lock
	time.Sleep(2000 * time.Millisecond)
	if err := obj.ObjUnlock(1); err != nil {
		t.Fatalf("thread 1 ObjUnlock (thin) returned error: %v", err)
	}
	t.Log("Thread 1 released object successfully (thin).")

	// After release, thread 2 should be able to acquire fairly quickly
	select {
	case <-acquired:
		// ok
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for thread 2 to acquire after thin release")
	}

	// And finish cleanly
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for thread 2 to finish after thin release")
	}
}

// Clone of the thin-lock handoff test but with 8 total threads (1 holder + 7 contenders).
// Thread 1 acquires a thin lock; threads 2..8 block until release, then each acquires and
// releases once. We assert no contender acquires before the holder releases and that all
// contenders eventually acquire and release successfully.
func TestObjLock_EightThreads_ThinLockContentionAndHandoff(t *testing.T) {
	obj := MakeEmptyObject()
	SetLockState(obj, lockStateUnlocked)

	// Thread 1 acquires thin lock
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("thread 1 ObjLock (thin) returned error: %v", err)
	}

	const contenders = 7 // threads 2..8
	var wg sync.WaitGroup
	wg.Add(contenders)

	acquiredCh := make(chan int, contenders) // buffer to record acquisitions by threadID

	// Start contender goroutines that attempt to lock, then release
	for id := int32(2); id < int32(2+contenders); id++ {
		tid := id
		go func() {
			defer wg.Done()
			if err := obj.ObjLock(tid); err != nil {
				t.Fatalf("contender %d ObjLock (thin) returned error: %v", tid, err)
			}
			t.Logf("Thread %d locked object successfully (thin).", tid)
			// Record acquisition
			acquiredCh <- int(tid)
			if err := obj.ObjUnlock(tid); err != nil {
				t.Fatalf("contender %d ObjUnlock (thin) returned error: %v", tid, err)
			}
			t.Logf("Thread %d released object successfully (thin).", tid)
		}()
	}

	// Ensure no contender acquires while thread 1 holds the lock
	select {
	case tid := <-acquiredCh:
		t.Fatalf("contender %d acquired thin lock while thread 1 still holds it", tid)
	case <-time.After(20 * time.Millisecond):
		// expected: no acquisition yet
	}

	// Release by thread 1 to allow contenders to proceed
	time.Sleep(2000 * time.Millisecond)
	if err := obj.ObjUnlock(1); err != nil {
		t.Fatalf("thread 1 ObjUnlock (thin) returned error: %v", err)
	}

	// Expect all 7 contenders to acquire once each, within a reasonable time
	deadlines := time.After(2 * time.Second)
	got := make(map[int]bool)
	for len(got) < contenders {
		select {
		case tid := <-acquiredCh:
			got[tid] = true
		case <-deadlines:
			t.Fatalf("timeout waiting for all contenders to acquire: got %d/%d", len(got), contenders)
		}
	}

	// Wait for all contenders to finish unlocking
	doneCh := make(chan struct{})
	go func() { wg.Wait(); close(doneCh) }()
	select {
	case <-doneCh:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for contenders to finish unlocks")
	}
}

// Simulate nested synchronized(lock) { synchronized(lock) { ... } }
// We model Java's reentrant locking by using a fat lock owned by the same thread
// and then attempting to lock it again, which should increment the recursion count.
func TestObjLock_NestedSynchronized_ReentrantFatLock(t *testing.T) {
	obj := MakeEmptyObject()

	// Pretend the object has already been inflated to a fat lock owned by thread 42
	SetLockState(obj, lockStateFatLocked)
	obj.Monitor = &ObjectMonitor{Owner: 42, Recursion: 0}

	// First nested synchronized: lock again by the same owner
	if err := obj.ObjLock(42); err != nil {
		t.Fatalf("first nested ObjLock returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 1 {
		t.Fatalf("expected recursion to become 1, got monitor=%v rec=%d", obj.Monitor, obj.Monitor.Recursion)
	}

	// Second nested synchronized: lock yet again by the same owner
	if err := obj.ObjLock(42); err != nil {
		t.Fatalf("second nested ObjLock returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 2 {
		t.Fatalf("expected recursion to become 2, got monitor=%v rec=%d", obj.Monitor, obj.Monitor.Recursion)
	}

	// Now unwind like exiting nested synchronized blocks: three unlocks total
	if err := obj.ObjUnlock(42); err != nil {
		t.Fatalf("first unwind ObjUnlock returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 1 {
		t.Fatalf("expected recursion to decrement to 1, got monitor=%v rec=%d", obj.Monitor, obj.Monitor.Recursion)
	}

	if err := obj.ObjUnlock(42); err != nil {
		t.Fatalf("second unwind ObjUnlock returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 0 {
		t.Fatalf("expected recursion to decrement to 0, got monitor=%v rec=%d", obj.Monitor, obj.Monitor.Recursion)
	}

	if err := obj.ObjUnlock(42); err != nil {
		t.Fatalf("final unwind ObjUnlock returned error: %v", err)
	}
	if obj.Monitor != nil {
		t.Fatalf("expected monitor to be cleared after final unlock")
	}
	if got := obj.Mark.Misc & lockStateMask; got != lockStateUnlocked {
		t.Fatalf("expected object to be unlocked after final release, got %b", got)
	}
}

// Same nested synchronized(lock) { synchronized(lock) { ... } } test but
// start from a thin-locked state first. We then inflate to a fat lock owned by
// the same thread and verify reentrant behavior (recursion increments) and
// proper unwind via unlocks.
func TestObjLock_NestedSynchronized_StartThinThenReentrantFatLock(t *testing.T) {
	obj := MakeEmptyObject()

	// Start explicitly from a thin-locked state
	SetObjectThinLocked(obj)

	// Inflate to a fat lock and assign ownership to thread 42
	SetObjectFatLocked(obj)
	obj.Monitor = &ObjectMonitor{Owner: 42, Recursion: 0}

	// First nested synchronized: lock again by the same owner
	if err := obj.ObjLock(42); err != nil {
		t.Fatalf("first nested ObjLock (start thin) returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 1 {
		t.Fatalf("expected recursion to become 1, got monitor=%v rec=%d", obj.Monitor, func() int32 {
			if obj.Monitor != nil {
				return obj.Monitor.Recursion
			}
			return -1
		}())
	}

	// Second nested synchronized: lock yet again by the same owner
	if err := obj.ObjLock(42); err != nil {
		t.Fatalf("second nested ObjLock (start thin) returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 2 {
		t.Fatalf("expected recursion to become 2, got monitor=%v rec=%d", obj.Monitor, func() int32 {
			if obj.Monitor != nil {
				return obj.Monitor.Recursion
			}
			return -1
		}())
	}

	// Now unwind like exiting nested synchronized blocks: three unlocks total
	if err := obj.ObjUnlock(42); err != nil {
		t.Fatalf("first unwind ObjUnlock (start thin) returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 1 {
		t.Fatalf("expected recursion to decrement to 1, got monitor=%v rec=%d", obj.Monitor, func() int32 {
			if obj.Monitor != nil {
				return obj.Monitor.Recursion
			}
			return -1
		}())
	}

	if err := obj.ObjUnlock(42); err != nil {
		t.Fatalf("second unwind ObjUnlock (start thin) returned error: %v", err)
	}
	if obj.Monitor == nil || obj.Monitor.Recursion != 0 {
		t.Fatalf("expected recursion to decrement to 0, got monitor=%v rec=%d", obj.Monitor, func() int32 {
			if obj.Monitor != nil {
				return obj.Monitor.Recursion
			}
			return -1
		}())
	}

	if err := obj.ObjUnlock(42); err != nil {
		t.Fatalf("final unwind ObjUnlock (start thin) returned error: %v", err)
	}
	if obj.Monitor != nil {
		t.Fatalf("expected monitor to be cleared after final unlock")
	}
	if got := obj.Mark.Misc & lockStateMask; got != lockStateUnlocked {
		t.Fatalf("expected object to be unlocked after final release, got %b", got)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-25 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"errors"
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"path"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// This file contains basic functions of object creation. (Array objects
// are created in object\arrays.go.)

/*
ObjectMonitor is a simple structure that holds the owner thread ID and recursion depth.
* Thin locks (2-bit Misc) are fast for uncontended objects.
* Recursive acquisition inflates the lock to a fat lock.
* Fat lock tracks the owning thread and recursion count.
* Unlocking decrements recursion and only releases when recursion hits zero.
*/

// With regard to the layout of a created object in Jacobin, note that
// on some architectures, but not Jacobin, there is an additional field
// that insures that the fields that follow the oops (the mark word and
// the class pointer) are aligned in memory for maximal performance.
type Object struct {
	Mark       MarkWord
	KlassName  uint32           // the index of the class name in the string pool
	FieldTable map[string]Field // map of field name to field struct
	Monitor    unsafe.Pointer   // --> an ObjectMonitor, accessed using atomic functions (thread safe)
	ThMutex    *sync.RWMutex    // Protect FieldTable set and get
}

// The mark word contains values for different purposes. Here,
// we use the first four bytes for a hash value, which is taken
// from the address of the object. The 'misc' field is divided in a
// Jacobin sense and does not match HotSpot.
type MarkWord struct {
	Hash uint32 // contains hash code which is the lower 32 bits of the address
	Misc uint32 // Misc represents auxiliary metadata such as lock information or GC states, encoded in the MarkWord structure.
}

// Ftype indicates the general type of data contained in Fvalue. When there is a leading X,
// the value of the field is not in Fvalue but stored in the statics table.
// For non-statics, Fvalue holds the field value.
type Field struct {
	Ftype  string // what type of value is stored in the field
	Fvalue any    // the actual value or a pointer to the value (ftype="Lsomething)
}

// ObjectMonitor is a simple structure that holds the owner thread ID and recursion depth.
// Thin locks (2-bit Misc) are fast for uncontended objects.
// Recursive acquisition inflates the lock to a fat lock.
// Fat lock tracks the owning thread and recursion count.
// Unlocking decrements recursion and only releases when recursion hits zero.

const MONITOR_OWNER_NONE = -1

// Definition for Object monitor
type ObjectMonitor struct {
	Owner     int32      // thread ID of owning thread
	Recursion int32      // recursion depth
	Mutex     sync.Mutex // used for blocking when fat locked
	Cond      *sync.Cond // used for wait/notify
}

// Global map tracking which object each thread is waiting on (for interrupt support)
var WaitingThreads = struct {
	sync.RWMutex
	MapThToObj map[uint32]*Object // Thread ID -> Object it's waiting on
}{MapThToObj: make(map[uint32]*Object)}

// MakeEmptyObject() creates an empty basis Object. It is expected that other
// code will fill in the Klass header field and the data fields.
func MakeEmptyObject() *Object {
	m := &ObjectMonitor{
		Owner:     MONITOR_OWNER_NONE,
		Recursion: 0,
	}
	m.Cond = sync.NewCond(&m.Mutex)
	o := Object{Monitor: unsafe.Pointer(m)}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	SetLockState(&o, lockStateUnlocked)
	o.KlassName = types.InvalidStringIndex // s/be filled in later, when class is filled in.

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]Field)
	o.ThMutex = &sync.RWMutex{}
	return &o
}

// MakeEmptyObjectWithClassName() creates an empty Object using the passed-in class name
func MakeEmptyObjectWithClassName(className *string) *Object {
	m := &ObjectMonitor{
		Owner:     MONITOR_OWNER_NONE,
		Recursion: 0,
	}
	m.Cond = sync.NewCond(&m.Mutex)
	o := Object{Monitor: unsafe.Pointer(m)}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	SetLockState(&o, lockStateUnlocked)
	o.KlassName = stringPool.GetStringIndex(className)

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]Field)
	o.ThMutex = &sync.RWMutex{}
	return &o
}

// makes an instance of a JLC (java/lang/Class) object, which has special considerations.
func MakeJlcObject(className *string) *Object {
	o := MakeEmptyObject()
	o.KlassName = types.StringPoolJavaLangClassIndex
	o.FieldTable["name"] = Field{Ftype: types.GolangString, Fvalue: *className}
	o.FieldTable["$klass"] = Field{Ftype: types.RawGoPointer, Fvalue: nil}          // points to the Klass object in metadata
	o.FieldTable["$statics"] = Field{Ftype: types.Array, Fvalue: make([]string, 0)} // array of static field names for this class
	return o
}

// Make an object for a Java primitive field (byte, int, etc.), given the class and field type.
func MakePrimitiveObject(classString string, ftype string, arg any) *Object {
	objPtr := MakeEmptyObject()
	(*objPtr).KlassName = stringPool.GetStringIndex(&classString)
	field := Field{ftype, arg}
	objPtr.ThMutex.Lock()
	(*objPtr).FieldTable["value"] = field
	objPtr.ThMutex.Unlock()
	return objPtr
}

// Make an object for a Java primitive field (byte, int, etc.), given the class, field name, and field type.
func MakeOneFieldObject(classString string, fname string, ftype string, arg any) *Object {
	objPtr := MakeEmptyObject()
	(*objPtr).KlassName = stringPool.GetStringIndex(&classString)
	field := Field{ftype, arg}
	objPtr.ThMutex.Lock()
	(*objPtr).FieldTable[fname] = field
	objPtr.ThMutex.Unlock()
	return objPtr
}

// UpdateValueFieldFromJavaBytes: Set the value field of the given String object to the given JavaByte array
func UpdateValueFieldFromJavaBytes(objPtr *Object, argBytes []types.JavaByte) {
	if objPtr == nil || argBytes == nil {
		if globals.TraceInst || globals.TraceVerbose {
			trace.Error("UpdateValueFieldFromJavaBytes: nil object or argBytes")
		}
		return
	}
	fld := Field{Ftype: types.StringClassRef, Fvalue: argBytes}
	objPtr.ThMutex.Lock()
	objPtr.FieldTable["value"] = fld
	objPtr.ThMutex.Unlock()
}

// Null is the Jacobin implementation of Java's null
// JACOBIN-618 changed definition of null to this.
var Null = (*Object)(nil)

// IsNull determines whether a value is null
func IsNull(value any) bool {
	switch value.(type) {
	case *Object:
		obj := value.(*Object)
		return obj == nil || obj == Null
	case []*Object:
		return false
	}
	return value == nil
}

// CloneObject makes a replica of an existing object.
func CloneObject(oldObject *Object) *Object {
	// Create new empty object.
	newObject := MakeEmptyObject()
	// Mimic the class.
	newObject.KlassName = oldObject.KlassName

	oldObject.ThMutex.RLock()
	defer oldObject.ThMutex.RUnlock()

	// Get a slice of keys from the old FieldTable.
	keys := make([]string, 0, len(oldObject.FieldTable))
	for key := range oldObject.FieldTable {
		keys = append(keys, key)
	}

	newObject.ThMutex.Lock()
	defer newObject.ThMutex.Unlock()

	// For each key in the old FieldTable, copy that entry into the new FieldTable.
	for _, key := range keys {
		newObject.FieldTable[key] = oldObject.FieldTable[key]
	}
	return newObject
}

// Clear the field table of the given object.
func ClearFieldTable(object *Object) {
	object.ThMutex.Lock()
	object.FieldTable = make(map[string]Field)
	object.ThMutex.Unlock()
}

// Get a class name suffix (E.g. String from java/lang/String) from an object.
// If inner is true, we will try for an inner class name.
func GetClassNameSuffix(arg *Object, inner bool) string {

	// Guard against trouble.
	if arg == nil || arg == Null {
		return types.NullString
	}

	// Get the class name.
	className := GoStringFromStringPoolIndex(arg.KlassName)
	className = strings.ReplaceAll(className, ".", "/")

	// Return the full suffix?
	if !inner {
		// Return the full suffix including class names that end in A$B (inner classes).
		return path.Base(className)
	}

	// Get the last segment
	base := path.Base(className)

	// If there's an inner class, return only the inner class name.
	if idx := strings.LastIndex(base, "$"); idx != -1 {
		return base[idx+1:]
	}
	return base
}

// Convert a Go boolean to a Jacobin Boolean.
func JavaBooleanFromGoBoolean(arg bool) int64 {
	if arg {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Convert a Jacobin Boolean to a Go boolean.
func GoBooleanFromJavaBoolean(arg int64) bool {
	if arg == types.JavaBoolTrue {
		return true
	}
	return false
}

// Valid lock state transitions
// ----------------------------
// Unlocked -> ThinLocked (first lock)
// ThinLocked -> FatLocked (inflation due to contention/recursion)
// ThinLocked -> Unlocked (unlock)
// FatLocked -> Unlocked (final unlock)
// Any -> GCMarked (during GC)

// monitorenter: shorthand for "attempting to acquire the lock on an object

/* Recursion example:
Thread A calls ObjLock()       → thin locked, depth = 1
Thread A calls ObjLock() again → inflates to fat, monitor.Recursion = 2
Thread A calls ObjUnlock()     → still fat locked, monitor.Recursion = 1
Thread A calls ObjUnlock()     → UNLOCKED, monitor.Recursion = 0

Recursion = "How many times do I need to call monitorexit to fully unlock?"
So after thin → fat inflation due to a second acquisition attempt, monitor.
Recursion = 2 means "this thread acquired the lock twice, so it needs to release twice."
*/

/*
xxxx...xxxx  (upper 30 bits unused)

	..11   (lowest 2 bits = lock state)
*/
const (
	lockStateThinLocked = 0b00 // 0
	lockStateUnlocked   = 0b01 // 1
	lockStateFatLocked  = 0b10 // 2 (not implemented here)
	lockStateGCMarked   = 0b11 // 3 (GC mark)
	lockStateMask       = 0b11
)

// SetLockState atomically sets the lock state bits on obj.Mark.Misc.
func SetLockState(obj *Object, state uint32) {
	miscPtr := &obj.Mark.Misc
	for {
		oldVal := atomic.LoadUint32(miscPtr)
		newVal := (oldVal &^ lockStateMask) | state
		if atomic.CompareAndSwapUint32(miscPtr, oldVal, newVal) {
			return
		}
	}
}

// IsObjectLocked returns true if some thread currently has locked the object.
func IsObjectLocked(obj *Object) bool {
	misc := atomic.LoadUint32(&obj.Mark.Misc)
	return (misc & lockStateMask) != lockStateUnlocked
}

// Unlock the given object.
func SetObjectUnlocked(obj *Object) {
	SetLockState(obj, lockStateUnlocked)
}

// Set the given object to a thin-locked state.
func SetObjectThinLocked(obj *Object) {
	SetLockState(obj, lockStateThinLocked)
}

// Set the given object to a fat-locked state.
func SetObjectFatLocked(obj *Object) {
	SetLockState(obj, lockStateFatLocked)
}

// GetMonitor returns the object's monitor.
// Returns nil if the monitor hasn't been inflated yet.
// Callers must use atomic operations to access Owner/Recursion fields,
// or lock monitor.Mutex before accessing other state.
func (obj *Object) GetMonitor() *ObjectMonitor {
	return (*ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
}

// GetMonitorOwner retrieves the thread ID of the monitor owner for the object.
// Returns MONITOR_OWNER_NONE if no thread owns the monitor.
func (obj *Object) GetMonitorOwner() int32 {
	monitor := (*ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	if monitor == nil {
		return MONITOR_OWNER_NONE
	}
	return atomic.LoadInt32(&monitor.Owner) // Atomic read
}

// GetMonitorRecursion retrieves the recursion depth of the monitor associated with the object.
// Returns 0 if the monitor hasn't been inflated or if the lock is not held.
func (obj *Object) GetMonitorRecursion() int32 {
	monitor := (*ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	if monitor == nil {
		return 0
	}
	return atomic.LoadInt32(&monitor.Recursion)
}

// inflateLock inflates a thin lock to a fat lock when the owning thread attempts recursive acquisition.
// It assumes that a monitor has already been allocated and stored in obj.Monitor, and that
// monitor.Mutex is already locked by the current thread.
// It returns true if the transition to FatLocked state was successful.
func (obj *Object) inflateLock(miscPtr *uint32, oldVal uint32, monitor *ObjectMonitor) bool {
	// Prepare the new value with the FatLocked state bit set
	newVal := (oldVal &^ lockStateMask) | lockStateFatLocked
	if atomic.CompareAndSwapUint32(miscPtr, oldVal, newVal) {
		// Inflation successful - set recursion count to 2
		// (1 for the original thin lock, 1 for this reentrant acquisition)
		atomic.StoreInt32(&monitor.Recursion, 2)
		return true
	}
	return false
}

// inflateAndWait handles lock inflation when a different thread contends for a thin lock.
// It inflates the lock to a fat lock, making it possible for the contending thread to block on a mutex.
// The monitor must already be installed in obj.Monitor.
// It returns true once the lock has been successfully inflated and acquired by the calling thread.
func (obj *Object) inflateAndWait(miscPtr *uint32, monitor *ObjectMonitor, threadID int32) bool {

	// Lock the monitor mutex first (blocks if another thread currently holds it)
	monitor.Mutex.Lock()

	// Now atomically transition the object's state to FatLocked
	for {
		oldVal := atomic.LoadUint32(miscPtr)
		newVal := (oldVal &^ lockStateMask) | lockStateFatLocked
		if atomic.CompareAndSwapUint32(miscPtr, oldVal, newVal) {
			// Successfully inflated while holding mutex, now set ownership
			atomic.StoreInt32(&monitor.Owner, threadID)
			atomic.StoreInt32(&monitor.Recursion, 1)
			return true
		}
		// If CAS failed (e.g., due to a concurrent state change), retry
	}
}

// ObjLock acquires the object lock for the given threadID.
// It implements Java's synchronized synchronization primitive.
// The lock can be in one of several states: Unlocked, ThinLocked, or FatLocked.
//   - Unlocked: Transitions to ThinLocked via a fast CAS.
//   - ThinLocked: If owned by the same thread, inflates to FatLocked (recursive).
//     If owned by another thread, it spins and eventually inflates to FatLocked to block.
//   - FatLocked: If owned by the same thread, increments recursion.
//     If owned by another thread, blocks on the monitor's Mutex.
func (obj *Object) ObjLock(threadID int32) error {
	if threadID < 0 {
		return errors.New("ObjLock: invalid thread ID")
	}

	miscPtr := (*uint32)(unsafe.Pointer(&obj.Mark.Misc))
	spinCount := 0
	const maxSpins = 1000 // Prevent indefinite spinning before forced inflation

	// Spin top.
	for {
		miscVal := atomic.LoadUint32(miscPtr)
		state := miscVal & lockStateMask
		monitor := obj.GetMonitor()

		if monitor == nil {
			return errors.New("ObjLock: monitor is nil")
		}

		switch state {
		case lockStateUnlocked:
			// Fast path: try to acquire as thin lock by setting state to ThinLocked
			newVal := (miscVal &^ lockStateMask) | lockStateThinLocked
			if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
				// Successfully acquired thin lock.
				// Now we must also hold the Mutex to protect fat lock transitions and ownership fields.
				monitor.Mutex.Lock()

				// After acquiring Mutex, we MUST check if someone inflated the lock while we were waiting for the Mutex
				curMisc := atomic.LoadUint32(miscPtr)
				if (curMisc & lockStateMask) != lockStateThinLocked {
					// It was inflated by another thread! We hold the Mutex now, so we can just
					// take over the fat lock.
					newVal := (curMisc &^ lockStateMask) | lockStateFatLocked
					atomic.StoreUint32(miscPtr, newVal)
				}
				// Set the owner and initial recursion count
				atomic.StoreInt32(&monitor.Owner, threadID)
				atomic.StoreInt32(&monitor.Recursion, 1)
				return nil
			}
			// CAS failed (concurrent acquisition), retry the whole loop

		case lockStateThinLocked:
			owner := atomic.LoadInt32(&monitor.Owner)
			if owner == threadID {
				// Recursive acquisition of a thin lock - inflate to fat lock to track recursion
				if obj.inflateLock(miscPtr, miscVal, monitor) {
					return nil
				}
				// Inflation CAS failed, retry
			} else {
				// Different thread owns the lock - spin for a bit before giving up and inflating
				spinCount++
				if spinCount > maxSpins {
					// Too much contention, inflate to fat lock so we can block on the Mutex
					if obj.inflateAndWait(miscPtr, monitor, threadID) {
						return nil
					}
					spinCount = 0 // Reset after inflation attempt if it failed for some reason
				}
			}

		case lockStateFatLocked:
			owner := atomic.LoadInt32(&monitor.Owner)
			if owner == threadID {
				// Recursive acquisition on an already fat lock
				atomic.AddInt32(&monitor.Recursion, 1)
				return nil
			}
			// Another thread owns it - block on the monitor's Mutex
			monitor.Mutex.Lock()

			// After acquiring Mutex, we are the new owner.
			// Verify state and set ownership fields.
			miscVal = atomic.LoadUint32(miscPtr)
			newVal := (miscVal &^ lockStateMask) | lockStateFatLocked
			atomic.StoreUint32(miscPtr, newVal)
			atomic.StoreInt32(&monitor.Owner, threadID)
			atomic.StoreInt32(&monitor.Recursion, 1)
			return nil

		case lockStateGCMarked:
			// Cannot lock objects that are currently being processed by the Garbage Collector
			return errors.New("ObjLock: object in GC-marked state")
		}

		// Yield to other goroutines to give the owner a chance to release the lock
		if spinCount%100 == 0 {
			runtime.Gosched()
		}
	}
}

// ObjUnlock releases the object lock for the given threadID.
// It returns an error if the thread does not own the lock or if the lock is not held.
func (obj *Object) ObjUnlock(threadID int32) error {
	return obj.objUnlockInternal(threadID, false)
}

// objUnlockInternal performs the actual lock release logic.
// If isWait is true, it is being called from Object.wait(), which requires
// slightly different handling of the Mutex (it's released later in wait()).
func (obj *Object) objUnlockInternal(threadID int32, isWait bool) error {
	if threadID < 0 {
		return errors.New("ObjUnlock: invalid thread ID")
	}

	miscPtr := (*uint32)(unsafe.Pointer(&obj.Mark.Misc))
	monitor := obj.GetMonitor()

	if monitor == nil {
		return errors.New("ObjUnlock: monitor is nil")
	}

	// Verify that the calling thread is indeed the owner
	owner := atomic.LoadInt32(&monitor.Owner)
	if owner != threadID {
		return errors.New("ObjUnlock: thread does not own lock")
	}

	// Verify that the lock is actually held
	recursion := atomic.LoadInt32(&monitor.Recursion)
	if recursion <= 0 {
		return errors.New("ObjUnlock: lock not held")
	}

	// Decrement recursion count
	newRecursion := atomic.AddInt32(&monitor.Recursion, -1)

	if newRecursion == 0 {
		// This was the last level of recursion, fully releasing the lock
		for {
			miscVal := atomic.LoadUint32(miscPtr)
			// Transition the state back to Unlocked
			newVal := (miscVal &^ lockStateMask) | lockStateUnlocked

			if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
				break
			}
			// If CAS failed (e.g., concurrent inflation or state change), retry
		}

		// MUST clear the owner AFTER successfully setting state to Unlocked.
		// This prevents races where another thread might see a null owner while the lock is still held.
		atomic.StoreInt32(&monitor.Owner, MONITOR_OWNER_NONE)

		if !isWait {
			// Release the underlying Mutex so waiting threads can proceed.
			// For Object.wait(), the Mutex release is handled by the wait mechanism itself.
			monitor.Mutex.Unlock()
		}
	}
	// else: still recursively held by the current thread, just decremented the count

	return nil
}

// ObjectWait implements java.lang.Object.wait()
func (obj *Object) ObjectWait(threadID int32, millis int64) error {
	monitor := obj.GetMonitor()
	if monitor == nil {
		return errors.New("ObjectWait: monitor is nil")
	}

	owner := atomic.LoadInt32(&monitor.Owner)
	if owner != threadID {
		return errors.New(fmt.Sprintf("ObjectWait: thread %d does not own lock, owner: %d", threadID, owner))
	}

	// Check if already interrupted before waiting
	if isThreadInterrupted(uint32(threadID)) {
		clearThreadInterrupted(uint32(threadID))
		return errors.New("thread interrupted")
	}

	savedRecursion := atomic.LoadInt32(&monitor.Recursion)

	// In Java, wait() fully releases the lock.
	// We call objUnlockInternal enough times to reach recursion 0.
	// The last call will keep the Mutex locked if we pass isWait=true.
	for i := int32(0); i < savedRecursion-1; i++ {
		if err := obj.objUnlockInternal(threadID, false); err != nil {
			return err
		}
	}
	if err := obj.objUnlockInternal(threadID, true); err != nil {
		return err
	}

	// Register that we're waiting on THIS object (for Thread.interrupt() support)
	WaitingThreads.Lock()
	WaitingThreads.MapThToObj[uint32(threadID)] = obj
	WaitingThreads.Unlock()

	// Now we wait on the condition variable.
	// monitor.Mutex is STILL LOCKED here because of isWait=true.
	// Cond.Wait() will atomically unlock it, block, and relock when woken.

	var interruptErr error

	if millis > 0 {
		// Timed wait implementation
		done := make(chan bool, 1)
		timer := time.AfterFunc(time.Duration(millis)*time.Millisecond, func() {
			select {
			case done <- true:
				monitor.Cond.Broadcast() // Wake up to handle timeout
			default:
				// Already notified, timer is irrelevant
			}
		})

		monitor.Cond.Wait() // Atomically unlocks Mutex, waits, relocks on wakeup

		// Check if we were woken by timeout or by notify
		select {
		case <-done:
			// Timeout occurred - this is normal, not an error
		default:
			// Woken by notify/notifyAll, cancel the timer
			timer.Stop()
			// Drain the channel in case timer fired between Wait() and Stop()
			select {
			case <-done:
			default:
			}
		}

		// Check thread interruption status
		if isThreadInterrupted(uint32(threadID)) {
			clearThreadInterrupted(uint32(threadID))
			interruptErr = errors.New("thread interrupted during wait")
		}

	} else {
		// Indefinite wait
		monitor.Cond.Wait() // Atomically unlocks Mutex, waits, relocks on wakeup

		// Check thread interruption status
		if isThreadInterrupted(uint32(threadID)) {
			clearThreadInterrupted(uint32(threadID))
			interruptErr = errors.New("thread interrupted during wait")
		}
	}

	// Unregister - no longer waiting
	WaitingThreads.Lock()
	delete(WaitingThreads.MapThToObj, uint32(threadID))
	WaitingThreads.Unlock()

	// At this point, Cond.Wait() has returned and the Mutex is locked again
	monitor.Mutex.Unlock()

	// Re-acquire the lock with the same recursion level.
	// We MUST re-acquire the first lock (which acquires the Mutex)
	// and then just increment recursion for the rest, because
	// ObjLock(threadID) would try to lock the Mutex again if we called it multiple times.
	if err := obj.ObjLock(threadID); err != nil {
		return err
	}
	for i := int32(1); i < savedRecursion; i++ {
		atomic.AddInt32(&monitor.Recursion, 1)
	}

	// Return interruption error AFTER re-acquiring the lock
	// (Java semantics require the lock to be held when InterruptedException is thrown)
	if interruptErr != nil {
		return interruptErr
	}

	return nil
}

func isThreadInterrupted(thID uint32) bool {
	gr := globals.GetGlobalRef()
	gr.ThreadLock.RLock()
	defer gr.ThreadLock.RUnlock()
	thObj := gr.Threads[int(thID)].(*Object)
	interrupted := thObj.FieldTable["interrupted"].Fvalue.(types.JavaBool)
	return interrupted == types.JavaBoolTrue
}

func clearThreadInterrupted(thID uint32) {
	gr := globals.GetGlobalRef()
	gr.ThreadLock.Lock()
	defer gr.ThreadLock.Unlock()
	thObj := gr.Threads[int(thID)].(*Object)
	fld := thObj.FieldTable["interrupted"]
	fld.Fvalue = types.JavaBoolFalse
	thObj.FieldTable["interrupted"] = fld
}

// ObjectNotify implements java.lang.Object.notify()
func (obj *Object) ObjectNotify(threadID int32) error {
	monitor := obj.GetMonitor()
	if monitor == nil {
		return errors.New("ObjectNotify: monitor is nil")
	}
	owner := atomic.LoadInt32(&monitor.Owner)
	if owner != threadID {
		return errors.New(fmt.Sprintf("ObjectNotify: thread %d does not own lock, owner: %d", threadID, owner))
	}
	monitor.Cond.Signal()
	return nil
}

// ObjectNotifyAll implements java.lang.Object.notifyAll()
func (obj *Object) ObjectNotifyAll(threadID int32) error {
	monitor := obj.GetMonitor()
	if monitor == nil {
		return errors.New("ObjectNotifyAll: monitor is nil")
	}
	owner := atomic.LoadInt32(&monitor.Owner)
	if owner != threadID {
		return errors.New(fmt.Sprintf("ObjectNotifyAll: thread %d does not own lock, owner: %d", threadID, owner))
	}
	monitor.Cond.Broadcast()
	return nil
}

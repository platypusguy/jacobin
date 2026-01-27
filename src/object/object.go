/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-25 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"errors"
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
type ObjectMonitor struct {
	Owner     int32      // thread ID of owning thread
	Recursion int32      // recursion depth
	mutex     sync.Mutex // used for blocking when fat locked
	cond      *sync.Cond // used for wait/notify
}

// With regard to the layout of a created object in Jacobin, note that
// on some architectures, but not Jacobin, there is an additional field
// that insures that the fields that follow the oops (the mark word and
// the class pointer) are aligned in memory for maximal performance.
type Object struct {
	Mark       MarkWord
	KlassName  uint32           // the index of the class name in the string pool
	FieldTable map[string]Field // map mapping field name to field
	Monitor    unsafe.Pointer   // *ObjectMonitor, accessed atomically
	ThMutex    *sync.RWMutex    // non-nil ONLY for thread FieldTable processing
}

// These mark word contains values for different purposes. Here,
// we use the first four bytes for a hash value, which is taken
// from the address of the object. The 'misc' field is divided in a
// Jacobin sense and does not match HotSpot.
type MarkWord struct {
	Hash uint32 // contains hash code which is the lower 32 bits of the address
	Misc uint32 // Misc represents auxiliary metadata such as lock information or GC states, encoded in the MarkWord structure.
}

// We need to know the type of the field only to tell whether
// it occupies one or two slots on the stack when getfield and
// putfield bytecodes are executed. The type also flags static
// fields (with a leading X in the field type, which tells us
// to locate the value in the statics table.
type Field struct {
	Ftype  string // what type of value is stored in the field
	Fvalue any    // the actual value or a pointer to the value (ftype="[something)
}

// MakeEmptyObject() creates an empty basis Object. It is expected that other
// code will fill in the Klass header field and the data fields.
func MakeEmptyObject() *Object {
	m := &ObjectMonitor{
		Owner:     -1,
		Recursion: 0,
	}
	m.cond = sync.NewCond(&m.mutex)
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
		Owner:     -1,
		Recursion: 0,
	}
	m.cond = sync.NewCond(&m.mutex)
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

/*
   1. Thin Locks
       ◦ Only 2 bits (lockStateThinLocked)
       ◦ Cannot track owner --> simple spin until free
       ◦ CAS used for acquisition to avoid races
   2. Fat Locks
       ◦ Only allocated if Monitor exists (i.e., recursive lock)
       ◦ Tracks owner (Owner) and recursion (Recursion)
       ◦ Unlock decrements recursion and frees monitor when recursion reaches 0
   3. Spin & Yield
       ◦ runtime.Gosched() ensures CPU time for other goroutines
       ◦ No busy-waiting loop burns CPU
   4. Error Handling
       ◦ Returns error instead of panic for invalid operations
       ◦ Examples: unlocking unlocked object, unlocking a GC-marked object, or unlocking fat lock by non-owner

TODO: Track the owner thread even in thin-locking (Hotspot). Then we can support the Thread.holdsLock() query.

According to chatGPT on 12/3/2025:

What is “owner tracking in thin locking”

* The concept of Thin lock (also known as “lightweight lock”) for Java was described in the paper
Thin Locks: Featherweight Synchronization for Java by Bacon, Konuru, Murthy & Serrano.
They describe a header (“lockword”) that encodes the owner thread identifier and a nested lock count when the object is thin-locked.
* Specifically: when a thread acquires the lock, its thread ID is stored in bits in the object’s header along with a “count” for nested reentrant locking.
That is literally tracking the “owner (thread)” of the lock in the thin-lock.
* If there is contention, the thin-lock can “inflate” to a fat (heavyweight) lock.
* Hence “owner tracking” — storing which thread currently owns a lock in the object header — is a fundamental part of the thin-lock idea.
*/

// Set the lock state bits on obj.Mark.Misc.
func SetLockState(obj *Object, state uint32) {
	obj.Mark.Misc = (obj.Mark.Misc &^ lockStateMask) | state
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

func IsObjectLocked(obj *Object) bool {
	return (obj.Mark.Misc & lockStateMask) != lockStateUnlocked
}

// GetMonitor safely loads the monitor pointer
func (obj *Object) GetMonitor() *ObjectMonitor {
	return (*ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
}

// getMonitor safely loads the monitor pointer
func (obj *Object) getMonitor() *ObjectMonitor {
	return obj.GetMonitor()
}

func (obj *Object) GetMonitorOwner() int32 {
	monitor := (*ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	return monitor.Owner
}

func (obj *Object) GetMonitorRecursion() int32 {
	monitor := (*ObjectMonitor)(atomic.LoadPointer(&obj.Monitor))
	return monitor.Recursion
}

// inflateLock inflates from thin to fat lock for recursive acquisition
func (obj *Object) inflateLock(miscPtr *uint32, oldVal uint32, monitor *ObjectMonitor, threadID int32) bool {
	newVal := (oldVal &^ lockStateMask) | lockStateFatLocked
	if atomic.CompareAndSwapUint32(miscPtr, oldVal, newVal) {
		// Successfully inflated - increment recursion
		atomic.AddInt32(&monitor.Recursion, 1)
		// Mutex is already locked from thin lock
		return true
	}
	return false
}

// inflateAndWait inflates to fat lock and waits to acquire
func (obj *Object) inflateAndWait(miscPtr *uint32, oldVal uint32, monitor *ObjectMonitor, threadID int32) bool {
	newVal := (oldVal &^ lockStateMask) | lockStateFatLocked
	if atomic.CompareAndSwapUint32(miscPtr, oldVal, newVal) {
		// Successfully inflated - now block on mutex
		monitor.mutex.Lock()

		// After acquiring mutex, we MUST ensure the state is FatLocked
		// because the previous owner might have set it to Unlocked upon release.
		for {
			mVal := atomic.LoadUint32(miscPtr)
			nVal := (mVal &^ lockStateMask) | lockStateFatLocked
			if atomic.CompareAndSwapUint32(miscPtr, mVal, nVal) {
				break
			}
		}

		atomic.StoreInt32(&monitor.Owner, threadID)
		atomic.StoreInt32(&monitor.Recursion, 1)
		return true
	}
	return false
}

// ObjLock acquires the object lock for the given thread
func (obj *Object) ObjLock(threadID int32) error {
	if threadID < 0 {
		return errors.New("ObjLock: invalid thread ID")
	}

	miscPtr := (*uint32)(unsafe.Pointer(&obj.Mark.Misc))
	spinCount := 0
	const maxSpins = 1000 // Prevent indefinite spinning

	for {
		miscVal := atomic.LoadUint32(miscPtr)
		state := miscVal & lockStateMask
		monitor := obj.getMonitor()

		if monitor == nil {
			return errors.New("ObjLock: monitor is nil")
		}

		switch state {
		case lockStateUnlocked:
			// Fast path: try to acquire as thin lock
			newVal := (miscVal &^ lockStateMask) | lockStateThinLocked
			if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
				// Successfully acquired thin lock
				monitor.mutex.Lock()
				// After acquiring mutex, we MUST check if someone inflated the lock while we were waiting for the mutex
				curMisc := atomic.LoadUint32(miscPtr)
				if (curMisc & lockStateMask) != lockStateThinLocked {
					// It was inflated! We hold the mutex now, so we can just take over the fat lock
					newVal := (curMisc &^ lockStateMask) | lockStateFatLocked
					atomic.StoreUint32(miscPtr, newVal)
				}
				atomic.StoreInt32(&monitor.Owner, threadID)
				atomic.StoreInt32(&monitor.Recursion, 1)
				return nil
			}
			// CAS failed, retry

		case lockStateThinLocked:
			owner := atomic.LoadInt32(&monitor.Owner)
			if owner == threadID {
				// Recursive acquisition - inflate to fat lock
				if obj.inflateLock(miscPtr, miscVal, monitor, threadID) {
					return nil
				}
				// Inflation failed, retry
			} else {
				// Different thread owns lock - spin or inflate
				spinCount++
				if spinCount > maxSpins {
					// Too much contention, inflate to fat lock for blocking
					if obj.inflateAndWait(miscPtr, miscVal, monitor, threadID) {
						return nil
					}
					spinCount = 0 // Reset after inflation attempt
				}
			}

		case lockStateFatLocked:
			owner := atomic.LoadInt32(&monitor.Owner)
			if owner == threadID {
				// Recursive acquisition on fat lock
				atomic.AddInt32(&monitor.Recursion, 1)
				return nil
			}
			// Another thread owns it - block on mutex
			monitor.mutex.Lock()
			// After acquiring mutex, verify state and set ownership
			miscVal = atomic.LoadUint32(miscPtr)
			newVal := (miscVal &^ lockStateMask) | lockStateFatLocked
			atomic.StoreUint32(miscPtr, newVal)
			atomic.StoreInt32(&monitor.Owner, threadID)
			atomic.StoreInt32(&monitor.Recursion, 1)
			return nil

		case lockStateGCMarked:
			return errors.New("ObjLock: object in GC-marked state")
		}

		// Yield to other goroutines
		if spinCount%100 == 0 {
			runtime.Gosched()
		}
	}
}

func (obj *Object) ObjUnlock(threadID int32) error {
	return obj.objUnlockInternal(threadID, false)
}

// objUnlockInternal releases the object lock.
// If isWait is true, it means we are unlocking for Object.wait(),
// so we don't clear the owner until we actually exit the wait.
// Actually, in Java, wait() releases the lock entirely, so owner becomes -1.
func (obj *Object) objUnlockInternal(threadID int32, isWait bool) error {
	if threadID < 0 {
		return errors.New("ObjUnlock: invalid thread ID")
	}

	miscPtr := (*uint32)(unsafe.Pointer(&obj.Mark.Misc))
	monitor := obj.getMonitor()

	if monitor == nil {
		return errors.New("ObjUnlock: monitor is nil")
	}

	owner := atomic.LoadInt32(&monitor.Owner)
	if owner != threadID {
		return errors.New("ObjUnlock: thread does not own lock")
	}

	recursion := atomic.LoadInt32(&monitor.Recursion)
	if recursion <= 0 {
		return errors.New("ObjUnlock: lock not held")
	}

	// Decrement recursion count
	newRecursion := atomic.AddInt32(&monitor.Recursion, -1)

	if newRecursion == 0 {
		// Fully releasing the lock
		for {
			miscVal := atomic.LoadUint32(miscPtr)
			// Release fat lock or thin lock
			newVal := (miscVal &^ lockStateMask) | lockStateUnlocked

			if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
				break
			}
			// If CAS failed, it might be because someone else inflated it or something.
			// Retry until we successfully set it to Unlocked.
		}

		// MUST store -1 AFTER setting state to Unlocked, to avoid race in doMonitorexit checks
		atomic.StoreInt32(&monitor.Owner, -1)
		if !isWait {
			monitor.mutex.Unlock() // Always unlock mutex on final release
		}
	}
	// else: still recursively held, just decremented count

	return nil
}

// ObjectWait implements java.lang.Object.wait()
func (obj *Object) ObjectWait(threadID int32, millis int64) error {
	monitor := obj.getMonitor()
	if monitor == nil {
		return errors.New("ObjectWait: monitor is nil")
	}

	owner := atomic.LoadInt32(&monitor.Owner)
	if owner != threadID {
		return errors.New("ObjectWait: thread does not own lock")
	}

	savedRecursion := atomic.LoadInt32(&monitor.Recursion)

	// In Java, wait() fully releases the lock.
	// We call objUnlockInternal enough times to reach recursion 0.
	// The last call will keep the mutex locked if we pass isWait=true.
	for i := int32(0); i < savedRecursion-1; i++ {
		if err := obj.objUnlockInternal(threadID, false); err != nil {
			return err
		}
	}
	if err := obj.objUnlockInternal(threadID, true); err != nil {
		return err
	}

	// Now we wait on the condition variable.
	// monitor.mutex is STILL LOCKED here because of isWait=true.

	if millis > 0 {
		timeout := false
		timer := time.AfterFunc(time.Duration(millis)*time.Millisecond, func() {
			monitor.mutex.Lock()
			timeout = true
			monitor.cond.Broadcast() // Wake up to check timeout
			monitor.mutex.Unlock()
		})

		for !timeout {
			monitor.cond.Wait()
			// If we are here, we were either signaled OR the timer fired.
			if timer.Stop() {
				// We were signaled before timeout.
				break
			}
			// If Stop() returned false, timer already fired.
			// timeout will be true next iteration.
		}
	} else {
		monitor.cond.Wait()
	}
	monitor.mutex.Unlock()

	// Re-acquire the lock with the same recursion level.
	// We MUST re-acquire the first lock (which acquires the mutex)
	// and then just increment recursion for the rest, because
	// ObjLock(threadID) would try to lock the mutex again if we called it multiple times.
	if err := obj.ObjLock(threadID); err != nil {
		return err
	}
	for i := int32(1); i < savedRecursion; i++ {
		atomic.AddInt32(&monitor.Recursion, 1)
	}

	return nil
}

// ObjectNotify implements java.lang.Object.notify()
func (obj *Object) ObjectNotify(threadID int32) error {
	monitor := obj.getMonitor()
	if monitor == nil {
		return errors.New("ObjectNotify: monitor is nil")
	}
	if atomic.LoadInt32(&monitor.Owner) != threadID {
		return errors.New("ObjectNotify: thread does not own lock")
	}
	monitor.cond.Signal()
	return nil
}

// ObjectNotifyAll implements java.lang.Object.notifyAll()
func (obj *Object) ObjectNotifyAll(threadID int32) error {
	monitor := obj.getMonitor()
	if monitor == nil {
		return errors.New("ObjectNotifyAll: monitor is nil")
	}
	if atomic.LoadInt32(&monitor.Owner) != threadID {
		return errors.New("ObjectNotifyAll: thread does not own lock")
	}
	monitor.cond.Broadcast()
	return nil
}

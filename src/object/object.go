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
	"sync/atomic"
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
	Owner     int32 // thread ID of owning thread
	Recursion int32 // recursion depth
}

// With regard to the layout of a created object in Jacobin, note that
// on some architectures, but not Jacobin, there is an additional field
// that insures that the fields that follow the oops (the mark word and
// the class pointer) are aligned in memory for maximal performance.
type Object struct {
	Mark       MarkWord
	KlassName  uint32           // the index of the class name in the string pool
	FieldTable map[string]Field // map mapping field name to field
	Monitor    *ObjectMonitor   // needed if fat locking the object
}

// These mark word contains values for different purposes. Here,
// we use the first four bytes for a hash value, which is taken
// from the address of the object. The 'misc' field is divided in a
// Jacobin sense and does not match HotSpot.
type MarkWord struct {
	Hash uint32 // contains hash code which is the lower 32 bits of the address
	Misc uint32 //
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
	o := Object{}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	SetLockState(&o, lockStateUnlocked)
	o.KlassName = types.InvalidStringIndex // s/be filled in later, when class is filled in.

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]Field)
	return &o
}

// MakeEmptyObjectWithClassName() creates an empty Object using the passed-in class name
func MakeEmptyObjectWithClassName(className *string) *Object {
	o := Object{}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	SetLockState(&o, lockStateUnlocked)
	o.KlassName = stringPool.GetStringIndex(className)

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]Field)
	return &o
}

// Make an object for a Java primitive field (byte, int, etc.), given the class and field type.
func MakePrimitiveObject(classString string, ftype string, arg any) *Object {
	objPtr := MakeEmptyObject()
	(*objPtr).KlassName = stringPool.GetStringIndex(&classString)
	field := Field{ftype, arg}
	(*objPtr).FieldTable["value"] = field
	return objPtr
}

// Make an object for a Java primitive field (byte, int, etc.), given the class, field name, and field type.
func MakeOneFieldObject(classString string, fname string, ftype string, arg any) *Object {
	objPtr := MakeEmptyObject()
	(*objPtr).KlassName = stringPool.GetStringIndex(&classString)
	field := Field{ftype, arg}
	(*objPtr).FieldTable[fname] = field
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
	objPtr.FieldTable["value"] = fld
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
	// Get a slice of keys from the old FieldTable.
	keys := make([]string, 0, len(oldObject.FieldTable))
	for key := range oldObject.FieldTable {
		keys = append(keys, key)
	}
	// For each key in the old FieldTable, copy that entry into the new FieldTable.
	for _, key := range keys {
		newObject.FieldTable[key] = oldObject.FieldTable[key]
	}
	return newObject
}

// Clear the field table of the given object.
func ClearFieldTable(object *Object) {
	object.FieldTable = make(map[string]Field)
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

// Lock the object to the specified thread.
func (obj *Object) ObjLock(threadID int32) error {
	miscPtr := (*uint32)(unsafe.Pointer(&obj.Mark.Misc))

	for {
		miscVal := atomic.LoadUint32(miscPtr)
		state := miscVal & lockStateMask

		switch state {

		case lockStateUnlocked:
			// Fast path: object is unlocked --> try to acquire thin lock
			newVal := (miscVal &^ lockStateMask) | lockStateThinLocked
			if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
				// Lock acquired successfully as thin; record owner for potential reentry inflation
				obj.Monitor = &ObjectMonitor{Owner: threadID, Recursion: 0}
				return nil
			}

			// CAS failed --> retry in next for-loop iteration

		case lockStateThinLocked:
			// If the same thread re-enters while thin-locked, inflate to fat and
			// treat as recursive acquisition.
			monitor := obj.Monitor
			if monitor != nil && monitor.Owner == threadID {
				// Attempt to atomically flip thin -> fat in the header first.
				newVal := (miscVal &^ lockStateMask) | lockStateFatLocked
				if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
					// Successful inflation. Increment recursion to account for reentry.
					atomic.AddInt32(&monitor.Recursion, 1)
					return nil
				}
				// If CAS failed, loop and retry as state may have changed.
				break
			}
			// Different thread (or unknown owner) --> spin until lock becomes free

		case lockStateFatLocked:
			monitor := obj.Monitor
			if monitor == nil {
				// Another thread may be in the middle of releasing the fat lock.
				// Yield and retry until state or monitor becomes consistent.
				runtime.Gosched()
				break
			}

			if monitor.Owner == threadID {
				// Recursive acquisition --> increment recursion count
				atomic.AddInt32(&monitor.Recursion, 1)
				return nil
			}

			// Another thread owns the monitor --> spin and retry

		case lockStateGCMarked:
			// GC-marked object --> cannot lock
			return errors.New("ObjLock: object in GC-marked state")
		}

		// Let another thread run.
		runtime.Gosched()

	}
}

// Release the object lock from the specified thread.
func (obj *Object) ObjUnlock(threadID int32) error {
	miscPtr := (*uint32)(unsafe.Pointer(&obj.Mark.Misc))

	for {
		miscVal := atomic.LoadUint32(miscPtr)
		state := miscVal & lockStateMask

		switch state {

		case lockStateThinLocked:
			// Thin lock --> release by setting unlocked bits
			newVal := (miscVal &^ lockStateMask) | lockStateUnlocked
			if atomic.CompareAndSwapUint32(miscPtr, miscVal, newVal) {
				// Clear any thin-owner tracking to avoid stale ownership
				obj.Monitor = nil
				return nil
			}
			// CAS failed --> retry

		case lockStateFatLocked:
			monitor := obj.Monitor
			if monitor == nil {
				return errors.New("ObjUnlock: fat lock exists but monitor is nil")
			}

			if monitor.Owner != threadID {
				return errors.New("ObjUnlock: current thread does not own the monitor")
			}

			rec := atomic.LoadInt32(&monitor.Recursion)
			if rec > 0 {
				// Recursive lock --> decrement recursion count
				atomic.AddInt32(&monitor.Recursion, -1)
				return nil
			}

			// Last unlock --> release fat lock.
			// First flip state to unlocked so contenders won't observe fat+nil.
			atomic.StoreUint32(miscPtr, (miscVal&^lockStateMask)|lockStateUnlocked)
			// Then clear the monitor, but only if it hasn't been replaced by a contender
			// that acquired the lock immediately after we marked it unlocked.
			if obj.Monitor == monitor {
				obj.Monitor = nil
			}
			return nil

		case lockStateUnlocked:
			return errors.New("ObjUnlock: object is already unlocked")

		case lockStateGCMarked:
			return errors.New("ObjUnlock: object in GC-marked state")
		}

		// Yield CPU and retry if CAS failed or lock not available
		runtime.Gosched()
	}
}

// Tree View Object (TVO) provides a comprehensive debug view of an Object
// for the GoLand debugger.
//
// Call from "Evaluate Expression": obj.TVObject()
// Displays:
// - Mark.Misc value
// - Class name from string pool
// - String values for []int8 fields
// - Integer values for integer fields
func (obj *Object) TVO() string {
	if obj == nil {
		return "Object: nil"
	}

	var sb strings.Builder
	sb.WriteString("=== Tree View Object ===\n")

	// Display Mark.Misc value
	sb.WriteString(fmt.Sprintf("Mark.Misc: %d (0x%08X)\n", obj.Mark.Misc, obj.Mark.Misc))

	// Display the class name from the string pool.
	className := GoStringFromStringPoolIndex(obj.KlassName)
	sb.WriteString(fmt.Sprintf("Class: %s\n", className))

	// Display fields
	if len(obj.FieldTable) > 0 {

		// Create a slice of keys.
		keys := make([]string, 0, len(obj.FieldTable))
		for key := range obj.FieldTable {
			keys = append(keys, key)
		}

		// Sort the keys, case-insensitive.
		globals.SortCaseInsensitive(&keys)

		sb.WriteString("Fields:\n")

		// For each field .....
		for _, fieldName := range keys {

			field := obj.FieldTable[fieldName]
			value := field.Fvalue
			// Check for integer types
			switch value.(type) {

			case int64:
				if field.Ftype == types.Bool {
					var str string
					if field.Fvalue == types.JavaBoolTrue {
						str = "true"
					} else {
						str = "false"
					}
					sb.WriteString(fmt.Sprintf("  %s [%s]: %s\n", fieldName, field.Ftype, str))
				} else {
					sb.WriteString(fmt.Sprintf("  %s [%s]: %d\n", fieldName, field.Ftype, value))
				}

			case []types.JavaByte:
				if field.Ftype == types.ByteArray || field.Ftype == types.StringClassName || field.Ftype == types.StringClassRef {
					str := GoStringFromJavaByteArray(value.([]types.JavaByte))
					sb.WriteString(fmt.Sprintf("  %s [%s]: %s\n", fieldName, field.Ftype, str))
				} else {
					sb.WriteString(fmt.Sprintf("  %s [%s]: %v\n", fieldName, field.Ftype, field.Fvalue))
				}

			case *Object:
				clname := GoStringFromStringPoolIndex(value.(*Object).KlassName)
				if clname == types.StringClassName {
					str := GoStringFromStringObject(value.(*Object))
					sb.WriteString(fmt.Sprintf("  %s [object %s]: %s\n", fieldName, field.Ftype, str))
				} else {
					sb.WriteString(fmt.Sprintf("  %s [%s]: class %s\n", fieldName, field.Ftype, clname))
				}

			case int8, int16, int32, uint8, uint16, uint32, uint64:
				sb.WriteString(fmt.Sprintf("  %s [%s]: %d\n", fieldName, field.Ftype, value))

			default:
				// For other types, show type and value
				sb.WriteString(fmt.Sprintf("  %s [%s]: %v\n", fieldName, field.Ftype, field.Fvalue))
			}
		}
	} else {
		sb.WriteString("Fields: (none)\n")
	}

	return sb.String()
}

// STR provides a string view of a Java byte array ([]types.JavaByte)
// for the GoLand debugger.
//
// Call from "Evaluate Expression": object.STR(array)
func STR(array []types.JavaByte) string {
	return GoStringFromJavaByteArray(array)
}

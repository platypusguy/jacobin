package javaLang

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"runtime"
	"sync"
	"time"
)

func cloneNotSupportedException(_ []interface{}) interface{} {
	errMsg := "cloneNotSupportedException: Not supported for threads"
	return ghelpers.GetGErrBlk(excNames.CloneNotSupportedException, errMsg)
}

// Get the thread state and return it to caller.
func GetThreadState(th *object.Object) int {
	thStateObj, ok := th.FieldTable["state"].Fvalue.(*object.Object)
	if !ok {
		return UNDEFINED
	}
	return thStateObj.FieldTable["value"].Fvalue.(int)
}

// Has the given thread been interrupted?
func isInterrupted(th *object.Object) bool {
	interruptedObj, ok := th.FieldTable["interrupted"].Fvalue.(*object.Object)
	if !ok {
		return false
	}
	interruptedVal, ok := interruptedObj.FieldTable["value"].Fvalue.(int)
	return ok && interruptedVal != 0
}

// Populate the thread object with default values.
// Note that the thread number is incremented in the call to threadNumberingNext().
func populateThreadObject(t *object.Object) {

	idField := object.Field{Ftype: types.Int, Fvalue: threadNumberingNext(nil).(int64)}
	t.FieldTable["ID"] = idField

	// the JDK defaults to "Thread-N" where N is the thread number
	// the sole exception is the main thread, which is called "main"
	defaultName := fmt.Sprintf("Thread-%d", idField.Fvalue)
	nameField := object.Field{Ftype: types.ByteArray, Fvalue: object.StringObjectFromGoString(defaultName)}
	t.FieldTable["name"] = nameField

	stateField := object.Field{Ftype: types.Ref, Fvalue: threadStateCreateWithValue([]any{NEW})}
	t.FieldTable["state"] = stateField

	daemonField := object.Field{Ftype: types.Int, Fvalue: types.JavaBoolFalse}
	t.FieldTable["daemon"] = daemonField

	interruptedField := object.Field{Ftype: types.Int, Fvalue: types.JavaBoolFalse}
	t.FieldTable["interrupted"] = interruptedField

	globals.GetGlobalRef().TGLock.RLock()
	tg, ok := globals.GetGlobalRef().ThreadGroups["main"].(*object.Object)
	globals.GetGlobalRef().TGLock.RUnlock()
	if !ok {
		panic("populateThreadObject: globals.GetGlobalRef().ThreadGroups[\"main\"] does not exist")
	}

	// The default thread group is the main thread group
	threadGroup := object.Field{Ftype: types.Ref, Fvalue: tg}
	t.FieldTable["threadgroup"] = threadGroup

	priority := object.Field{Ftype: types.Int, Fvalue: statics.GetStaticValue("java/lang/Thread", "NORM_PRIORITY").(int64)}
	t.FieldTable["priority"] = priority

	frameStack := object.Field{Ftype: types.LinkedList, Fvalue: nil}
	t.FieldTable["framestack"] = frameStack

}

// Add the specified thread to the global registry of threads.
// TODO: Unused. Why?
func RegisterThread(t *object.Object) {
	glob := globals.GetGlobalRef()
	ID := int(t.FieldTable["ID"].Fvalue.(int64))
	glob.ThreadLock.Lock()
	glob.Threads[ID] = t
	glob.ThreadLock.Unlock()
}

// Create a new Runnable object.
// Store it in the target field of the thread object.
// The class name comes from the top-most frame on the frame stack.
// The method name and type are fixed: run()V
func storeThreadRunnable(t *object.Object, fs *list.List) {
	frame := *fs.Front().Value.(*frames.Frame)
	runnable := NewRunnable(
		object.JavaByteArrayFromGoString(frame.ClName),
		object.JavaByteArrayFromGoString("run"),
		object.JavaByteArrayFromGoString("()V"))
	t.FieldTable["target"] = object.Field{Ftype: types.Ref, Fvalue: runnable}
}

// Set the thread state to the supplied value unconditionally.
// Returns:
// * Previous state or -1 if unknown
// * Result
//   - nil (success)
//   - *ghelpers.GErrBlk (oops)
func SetThreadState(th *object.Object, newState int) (interface{}, interface{}) {
	// Returns (previousState int, error interface{})
	if err := th.ObjLock(0); err != nil {
		return -1, ghelpers.GetGErrBlk(excNames.VirtualMachineError,
			"SetThreadState: cannot lock thread object")
	}
	defer th.ObjUnlock(0)

	// Retrieve the 'state' field
	thStateObj, ok := th.FieldTable["state"].Fvalue.(*object.Object)
	if !ok || thStateObj == nil {
		// Create the new state object if missing.
		stateField := object.Field{
			Ftype:  types.Ref,
			Fvalue: threadStateCreateWithValue([]any{newState}),
		}
		th.FieldTable["state"] = stateField
		return -1, nil // no previous state
	}

	// Get previous state
	prevVal, ok := thStateObj.FieldTable["value"].Fvalue.(int)
	if !ok {
		prevVal = -1 // unknown previous state
	}

	// Only update if different
	if prevVal != newState {
		thStateObj.FieldTable["value"] = object.Field{
			Ftype:  types.Int,
			Fvalue: newState,
		}
	}

	return prevVal, nil
}

// Should we need to create a thread (as in tests), here is the instantiable implementation
func ThreadCreateNoarg(_ []interface{}) any {
	t := object.MakeEmptyObjectWithClassName(&types.ClassNameThread)
	populateThreadObject(t)
	return t
}

// Wait for a thread to be TERMINATED up to maxTime milliseconds.
// Returns:
// * nil if the thread terminated within maxTime milliseconds
// *     or the thread terminated after maxTime milliseconds
// * ghelpers.GetGErrBlk(IllegalArgumentException) if:
//   - interrupted while waiting
//   - cannot thin-lock the target thread object
//   - the thread state is not an object
//   - the thread state is missing the value field
//   - max wait time <= 0
//
// Sentinel for continuing the loop
var continueLoop = struct{}{}

func waitForTermination(waitingThread, targetThread *object.Object, maxTime int64) interface{} {
	if maxTime <= 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"joinThread: max wait time <= 0")
	}

	start := time.Now().UnixMilli()
	targetThreadID := int32(targetThread.FieldTable["ID"].Fvalue.(int64))

	for {
		// Lock the target thread to safely access its state
		if err := targetThread.ObjLock(targetThreadID); err != nil {
			return ghelpers.GetGErrBlk(excNames.IllegalMonitorStateException,
				"waitForTermination: unable to lock target thread object: "+err.Error())
		}

		// Locked the target thread object.
		result := func() interface{} {

			defer targetThread.ObjUnlock(targetThreadID)

			// .Get the target thread state value.
			stateObj, ok := targetThread.FieldTable["state"].Fvalue.(*object.Object)
			if !ok {
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
					"waitForTermination: field state is not an object")
			}
			stateVal, ok := stateObj.FieldTable["value"].Fvalue.(int)
			if !ok {
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
					"waitForTermination: state object is missing a value field")
			}

			// Terminated --> normal return
			if stateVal == TERMINATED {
				return nil
			}

			// Did anyone interrupt me?
			if isInterrupted(waitingThread) {
				return ghelpers.GetGErrBlk(excNames.InterruptedException,
					"waitForTermination: waiting thread was interrupted")
			}

			return continueLoop

		}()

		if result != continueLoop {
			return result
		}

		// Continue waiting.
		// Yield to allow the target thread to run.
		runtime.Gosched()

		// Timeout check.
		if time.Now().UnixMilli()-start >= maxTime {
			return nil // In Java, timeout is NORMAL!
		}
	}
}

// =========================================== THREAD ID (NUMBERING) FUNCTIONS =================================

// Note that ThreadNumbering is a private static class in java/lang/Thread (hotpot).
// setInitialThreadNumberingValue guarantees that the thread numbering is initialized only once.
var setInitialThreadNumberingValue = sync.OnceValue(func() any {
	threadNumber = int64(0)
	return nil
})

func threadNumbering(_ []any) any { // initialize thread numbering
	setInitialThreadNumberingValue()
	return nil
}

func threadNumberingNext(_ []any) any {
	threadNumberingMutex.Lock()
	threadNumber += 1
	//trace.Trace(fmt.Sprintf("threadNumberingNext: thread numbering incremented to %d", threadNumber))
	threadNumberingMutex.Unlock()
	return threadNumber
}

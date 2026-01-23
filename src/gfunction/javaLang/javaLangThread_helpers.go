package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
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
func GetThreadState(th *object.Object) int64 {
	th.ThMutex.RLock()
	defer th.ThMutex.RUnlock()
	state, ok := th.FieldTable["state"].Fvalue.(int64)
	if !ok {
		return -1
	}
	return state
}

// Set the thread state to the supplied value unconditionally.
// Returns:
// * Previous state or -1 if unknown
// * Result
//   - nil (success)
//   - *ghelpers.GErrBlk (oops)
func SetThreadState(th *object.Object, newState int64) (int64, interface{}) {
	th.ThMutex.Lock()
	defer th.ThMutex.Unlock()

	// Retrieve the current 'state' field if it exists.
	curVal, ok := th.FieldTable["state"].Fvalue.(int64)
	if !ok {
		// Create state field if missing (should normally not happen).
		fld := object.Field{Ftype: types.Int, Fvalue: newState}
		th.FieldTable["state"] = fld
		// Return value of -1 indicates no previous value.
		return -1, nil
	}

	// Update only if different
	if curVal != newState {
		fld := object.Field{Ftype: types.Int, Fvalue: newState}
		th.FieldTable["state"] = fld
	}

	// We're done, unless state is TERMINATED, in which case we remove it from the
	// global registry of threads.
	if newState == TERMINATED {
		globals.GetGlobalRef().ThreadLock.Lock()
		delete(globals.GetGlobalRef().Threads, int(th.FieldTable["ID"].Fvalue.(int64)))
		globals.GetGlobalRef().ThreadLock.Unlock()
	}

	return curVal, nil
}

// Has the given thread been interrupted?
func isInterrupted(th *object.Object) bool {
	th.ThMutex.RLock()
	interruptedObj, ok := th.FieldTable["interrupted"].Fvalue.(*object.Object)
	th.ThMutex.RUnlock()
	if !ok {
		return false
	}
	interruptedObj.ThMutex.RLock()
	interruptedVal, ok := interruptedObj.FieldTable["value"].Fvalue.(int)
	interruptedObj.ThMutex.RUnlock()
	return ok && interruptedVal != 0
}

// Populate the thread object with default values.
// Note that the thread number is incremented in the call to threadNumberingNext().
func populateThreadObject(t *object.Object) {

	t.ThMutex.Lock()
	defer t.ThMutex.Unlock()

	idField := object.Field{Ftype: types.Int, Fvalue: threadNumberingNext(nil).(int64)}
	t.FieldTable["ID"] = idField

	// the JDK defaults to "Thread-N" where N is the thread number
	// the sole exception is the main thread, which is called "main"
	defaultName := fmt.Sprintf("Thread-%d", idField.Fvalue)
	nameField := object.Field{Ftype: types.ByteArray, Fvalue: object.StringObjectFromGoString(defaultName)}
	t.FieldTable["name"] = nameField

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

	fld := object.Field{Ftype: types.Int, Fvalue: NEW}
	t.FieldTable["state"] = fld

}

// Add the specified thread to the global registry of threads.
func RegisterThread(th *object.Object) {
	th.ThMutex.RLock()
	ID := int(th.FieldTable["ID"].Fvalue.(int64))
	th.ThMutex.RUnlock()
	glob := globals.GetGlobalRef()
	glob.ThreadLock.Lock()
	glob.Threads[ID] = th
	glob.ThreadLock.Unlock()
}

// Should we need to create a thread (as in tests), here is the instantiable implementation
func ThreadCreateNoarg(_ []interface{}) any {
	th := object.MakeEmptyObjectWithClassName(&types.ClassNameThread)
	populateThreadObject(th)
	return th
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
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"waitForTermination: max wait time <= 0",
		)
	}

	start := time.Now().UnixMilli()

	for {

		// Get current target thread state.
		stateVal := GetThreadState(targetThread)

		// TERMINATED -> normal return
		if stateVal == TERMINATED {
			return nil
		}

		// Interrupted?
		if isInterrupted(waitingThread) {
			return ghelpers.GetGErrBlk(
				excNames.InterruptedException,
				"waitForTermination: waiting thread was interrupted",
			)
		}

		// Yield to allow target thread to run
		runtime.Gosched()

		// Timeout -> normal Java behavior
		if time.Now().UnixMilli()-start >= maxTime {
			return nil
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

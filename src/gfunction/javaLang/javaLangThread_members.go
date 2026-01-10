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
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"runtime"
	"time"
)

// Get the number of active threads in Jacobin.
func threadActiveCount(_ []interface{}) any {
	return int64(len(globals.GetGlobalRef().Threads))
}

// threadCurrentThread retrieves the current Thread Object from the frame stack provided in the parameters.
// Returns an error block if the input is invalid or if the required context data is not a valid frame stack.
func threadCurrentThread(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadCurrentThread: Expected context data, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	fStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadCurrentThread: Expected context data to be a frame stack"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	frame := *fStack.Front().Value.(*frames.Frame)
	thID := frame.Thread
	gr := globals.GetGlobalRef()
	gr.ThreadLock.RLock()
	defer gr.ThreadLock.RUnlock()
	th := gr.Threads[thID].(*object.Object)
	return th
}

// threadDumpStack dumps the JVM stack trace to stderr. Returns an error block if invalid parameters are provided.
func threadDumpStack(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadDumpStack: Expected context data, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jvmStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadDumpStack: Expected context data to be a frame"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	globalRef := globals.GetGlobalRef()
	if globalRef.StrictJDK { // if strictly following HotSpot, ...
		_, _ = fmt.Fprintln(os.Stderr, "java.lang.Exception: Stack trace")
	} else { // TODO: add the source line numbers to both variants
		// we print more data than HotSpot does, starting with the thread name
		o := *jvmStack.Front().Value.(*frames.Frame)
		threadID := o.Thread
		globalRef.ThreadLock.RLock()
		defer globalRef.ThreadLock.RUnlock()
		th := globalRef.Threads[threadID].(*object.Object)
		raws := th.FieldTable["name"].Fvalue.(*object.Object)
		threadName := object.GoStringFromStringObject(raws)
		_, _ = fmt.Fprintf(os.Stderr, "Stack trace (thread %s)\n", threadName)
	}

	for e := jvmStack.Front(); e != nil; e = e.Next() {
		fr := *e.Value.(*frames.Frame)
		if globalRef.StrictJDK {
			_, _ = fmt.Fprintf(os.Stderr, "\tat %s.%s\n", fr.ClName, fr.MethName)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "\tat %s.%s(PC: %d)\n",
				fr.ClName, fr.MethName, fr.PC)
		}
	}
	return nil
}

// java/lang/Thread.enumerate([Ljava/lang/Thread;)I
// per Javadoc: Copies into the specified array every live platform thread in this thread group and its subgroups.
// Virtual threads are not enumerated by this method.
func threadEnumerate(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadEnumerate expected a thread object, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	arrObj := params[0].(*object.Object)
	arr := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	count := len(arr)

	globalRef := globals.GetGlobalRef()
	threadCount := len(globalRef.Threads)
	count = min(count, threadCount)
	i := 0
	for _, value := range globalRef.Threads {
		arr[i] = value.(*object.Object)
		i += 1
	}
	return count
}

// threadGetId extracts and returns the ID field from a given non-null Thread object.
// Returns an error block if input validation fails or the ID field is missing/invalid.
func threadGetId(params []interface{}) any {
	// Expect exactly one parameter: the Thread object (this)
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetId: Expected 1 parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the parameter is a non-null Thread object
	t, ok := params[0].(*object.Object)
	if !ok || object.IsNull(t) {
		errMsg := "threadGetId: Expected first parameter to be a non-null Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract the ID field and ensure it is int64
	idField, ok := t.FieldTable["ID"]
	if !ok {
		errMsg := "threadGetId: Thread object missing 'ID' field"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}
	ID, ok := idField.Fvalue.(int64)
	if !ok {
		errMsg := "threadGetId: 'ID' field has unexpected type"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}

	return ID
}

// threadGetName retrieves the "name" field value of the given Thread object.
// Returns an error block if the input is invalid or not a Thread object.
func threadGetName(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetName: Expected no parameters, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGetName: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return t.FieldTable["name"].Fvalue.(*object.Object)
}

// threadGetPriority retrieves the "priority" field value from the provided object.
// Returns an error block if the input parameter count is invalid.
func threadGetPriority(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetPriority: Expected no parameters, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t := params[0].(*object.Object)
	return t.FieldTable["priority"].Fvalue.(int64)
}

// threadGetStackTrace retrieves the stack trace of a thread from the provided context parameters.
// It requires the JVM frame stack as context data.
// If the arguments are invalid or an error occurs during stack trace population, an appropriate exception is returned.
// The function returns the populated stack trace object.
func threadGetStackTrace(params []interface{}) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadGetStackTrace: Expected context data, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jvmFrameStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadGetStackTrace: Expected context data to be a frame stack"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	stackTrace := object.MakeEmptyObject()
	stackTrace.KlassName = object.StringPoolIndexFromGoString("[java/lang/StackTraceElement")
	ret := FillInStackTrace([]interface{}{jvmFrameStack, stackTrace})
	if ret == nil {
		errMsg := "threadGetStackTrace: Call to gfunction.FillInStackTrace() failed to fill in stack trace"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}
	traceObj := stackTrace.FieldTable["stackTrace"].Fvalue.(*object.Object)
	return traceObj
}

// threadGetState retrieves the "state" field of a given thread object.
// Returns an IllegalArgumentException error block if input is invalid.
func threadGetState(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetState: Expected 1 parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGetState: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	state := t.FieldTable["state"].Fvalue.(*object.Object)
	return state
}

// threadGetThreadGroup retrieves the thread group associated with the given thread object.
// Returns the thread group object if found, or an error block for invalid input or missing thread group information.
func threadGetThreadGroup(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetThreadGroup: Expected 1 parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t := params[0].(*object.Object)
	threadGroup, ok := t.FieldTable["threadgroup"].Fvalue.(*object.Object)
	if !ok {
		errMsg := "threadGetThreadGroup: Expected threadgroup to be an object"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}
	return threadGroup
}

// Is the specified thread alive?
func threadIsAlive(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadIsAlive: Expected 1 parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadIsAlive: Expected thread to be an object"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}
	state := GetThreadState(t)
	if state > NEW && state < TERMINATED {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Has the specified thread been interrupted?
func threadIsInterrupted(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadIsInterrupted: Expected 1 parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadIsInterrupted: Expected thread to be an object"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}
	return t.FieldTable["interrupted"].Fvalue.(types.JavaBool)
}

// Has the specified thread terminated?
func threadIsTerminated(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadIsTerminated: Expected 1 parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadIsTerminated: Expected thread to be an object"
		return ghelpers.GetGErrBlk(excNames.InternalException, errMsg)
	}
	state := GetThreadState(t)
	return state == TERMINATED
}

// threadJoin synchronizes the current thread with a target thread, optionally waiting for a specified time duration.
// Accepts parameters including the current thread's frame stack, the target thread, and optional wait time in millis/nanos.
// Returns an error block for invalid inputs or a termination synchronization result.
func threadJoin(params []interface{}) any {
	fStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadJoin: Expected context data to be a frame stack"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	frame := *fStack.Front().Value.(*frames.Frame)
	thID := frame.Thread
	gr := globals.GetGlobalRef()
	gr.ThreadLock.RLock()
	currentThread := gr.Threads[thID].(*object.Object)
	gr.ThreadLock.RUnlock()

	targetThread, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadJoin: Missing/erroneous target thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	millis := int64(ghelpers.MaxIntValue) // Hotspot waits forever
	if len(params) > 2 {
		millis = params[2].(int64)
		nanos := int64(0)
		if len(params) > 3 {
			nanos = params[3].(int64)
			if nanos > 0 {
				millis += 1 // not precise
			}
		}
	}

	return waitForTermination(currentThread, targetThread, millis)
}

// threadRun validates and processes a thread object passed as a single parameter and logs its start information.
// Returns a custom error block if the parameter is missing, incorrect, or invalid.
// Note: This function should never be called since the user should provide their own Thread.run() function.
func threadRun(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadRun: Expected only the thread object parameter, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadRun: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name := t.FieldTable["name"].Fvalue.(*object.Object)
	id := t.FieldTable["ID"].Fvalue.(int64)
	warnMsg := fmt.Sprintf("threadRun nil-function name: %s, ID: %d started", object.GoStringFromStringObject(name), id)
	trace.Warning(warnMsg)
	return nil
}

// threadSetName sets the name of a thread to a specified Java String.
// The function expects exactly two parameters: a thread object and a non-null Java String object for the name.
// Returns an error block if any parameter is invalid or updates the thread's name field otherwise.
func threadSetName(params []interface{}) any {
	// Expect exactly two parameters: the thread object and the Java String name
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadSetName: Expected 2 parameters, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the first parameter is the Thread object
	th, ok := params[0].(*object.Object)
	if !ok || object.IsNull(th) {
		errMsg := "threadSetName: Expected first parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the second parameter is the Java String object (non-null)
	nameObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(nameObj) {
		errMsg := "threadSetName: name must not be null"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Update the thread's name field (stored as a Java byte string)
	th.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: nameObj}

	return nil
}

// threadSetPriority sets the priority of a specified thread.
// Accepts a slice of two parameters: a Thread object and an int64 priority value.
// Raises IllegalArgumentException if parameters are invalid or priority is out of the valid range for threads.
func threadSetPriority(params []interface{}) any {
	// Expect exactly two parameters: the Thread object and the priority (int64)
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadSetPriority: Expected 2 parameters, got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the first parameter is the Thread object
	th, ok := params[0].(*object.Object)
	if !ok || object.IsNull(th) {
		errMsg := "threadSetPriority: Expected first parameter to be a non-null Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract and validate the second parameter (priority as int64)
	priority, ok := params[1].(int64)
	if !ok {
		errMsg := "threadSetPriority: priority must be an int64 (long)"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Fetch bounds from statics (java/lang/Thread.MIN_PRIORITY and MAX_PRIORITY)
	minP := statics.GetStaticValue("java/lang/Thread", "MIN_PRIORITY").(int64)
	maxP := statics.GetStaticValue("java/lang/Thread", "MAX_PRIORITY").(int64)

	if priority < minP || priority > maxP {
		errMsg := fmt.Sprintf("threadSetPriority: priority %d out of range [%d..%d]", priority, minP, maxP)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Update the thread's priority field (Jacobin stores it as an int64 under type types.Int)
	th.FieldTable["priority"] = object.Field{Ftype: types.Int, Fvalue: priority}
	return nil
}

// threadSleep pauses the current thread for the duration specified in milliseconds by the first parameter.
// If the parameter is not an int64, it returns an IOException error block.
func threadSleep(params []interface{}) interface{} {
	sleepTime, ok := params[0].(int64)
	if !ok {
		errMsg := "threadSleep: Parameter must be an int64 (long)"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return nil
}

// threadStart starts a new thread based on the given Thread object provided as a parameter.
// Expects a single parameter of type *object.Object, representing the Thread instance to be started.
// Returns nil on success or an error block when the input is invalid (e.g., missing Runnable object).
func threadStart(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadStart: Expected only the thread object parameter, got %d", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get thread object.
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadStart: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get runnable object.
	runnable, ok := t.FieldTable["target"].Fvalue.(*object.Object)
	if !ok {
		errMsg := "threadStart: Expected Runnable target field to be an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract class name, method name, and method type from the runnable object.
	var clName, methName, methType string
	ftbl := runnable.FieldTable
	fld, ok := ftbl["clName"]
	if !ok {
		errMsg := "threadStart: Missing the clName field in the runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	clName = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
	fld, ok = ftbl["methName"]
	if !ok {
		errMsg := "threadStart: Missing the methName field in the runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	methName = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
	fld, ok = ftbl["methType"]
	if !ok {
		errMsg := "threadStart: Missing the methType field in the runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	methType = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Spawn RunJavaThread to interpret bytecode of run()
	args := []interface{}{t, clName, methName, methType}
	if clName == "main$CounterTask" && methName == "run" { // JACOBIN-824 experimentation
		globals.GetGlobalRef().FuncRunThread(args)
		return nil
	} else {
		go globals.GetGlobalRef().FuncRunThread(args)
		runtime.Gosched()
		return nil
	}
	//
	// go globals.GetGlobalRef().FuncRunThread(args)
	// runtime.Gosched()
	//
	// return nil
}

// threadYield allows the current goroutine to relinquish the processor, enabling other goroutines to run.
// It uses runtime.Gosched() to yield execution and returns nil.
func threadYield([]interface{}) interface{} {
	runtime.Gosched()
	return nil
}

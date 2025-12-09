package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/frames"
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

// "java/lang/Thread.currentThread()Ljava/lang/Thread;"
func threadCurrentThread(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadCurrentThread: Expected context data, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	fStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadCurrentThread: Expected context data to be a frame stack"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	frame := *fStack.Front().Value.(*frames.Frame)
	thID := frame.Thread
	th := globals.GetGlobalRef().Threads[thID].(*object.Object)
	return th
}

// java/lang/Thread.dumpStack()V
func threadDumpStack(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadDumpStack: Expected context data, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jvmStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadDumpStack: Expected context data to be a frame"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	globalRef := globals.GetGlobalRef()
	if globalRef.StrictJDK { // if strictly following HotSpot, ...
		_, _ = fmt.Fprintln(os.Stderr, "java.lang.Exception: Stack trace")
	} else { // TODO: add the source line numbers to both variants
		// we print more data than HotSpot does, starting with the thread name
		o := *jvmStack.Front().Value.(*frames.Frame)
		threadID := o.Thread
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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

// "java/lang/Thread.getId()J"
func threadGetId(params []interface{}) any {
	// Expect exactly one parameter: the Thread object (this)
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetId: Expected 1 parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the parameter is a non-null Thread object
	t, ok := params[0].(*object.Object)
	if !ok || object.IsNull(t) {
		errMsg := "threadGetId: Expected first parameter to be a non-null Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract the ID field and ensure it is int64
	idField, ok := t.FieldTable["ID"]
	if !ok {
		errMsg := "threadGetId: Thread object missing 'ID' field"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	ID, ok := idField.Fvalue.(int64)
	if !ok {
		errMsg := "threadGetId: 'ID' field has unexpected type"
		return getGErrBlk(excNames.InternalException, errMsg)
	}

	return ID
}

// "java/lang/Thread.getName()Ljava/lang/String;"
func threadGetName(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetName: Expected no parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGetName: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return t.FieldTable["name"].Fvalue.(*object.Object)
}

// "java/lang/Thread.getPriority()I"
func threadGetPriority(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetPriority: Expected no parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t := params[0].(*object.Object)
	return t.FieldTable["priority"].Fvalue
}

// java/lang/Thread.getStackTrace()[Ljava/lang/StackTraceElement;
func threadGetStackTrace(params []interface{}) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadGetStackTrace: Expected context data, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jvmFrameStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadGetStackTrace: Expected context data to be a frame stack"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	stackTrace := object.MakeEmptyObject()
	stackTrace.KlassName = object.StringPoolIndexFromGoString("[java/lang/StackTraceElement")
	ret := FillInStackTrace([]interface{}{jvmFrameStack, stackTrace})
	if ret == nil {
		errMsg := "threadGetStackTrace: Call to gfunction.FillInStackTrace() failed to fill in stack trace"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	traceObj := stackTrace.FieldTable["stackTrace"].Fvalue.(*object.Object)
	return traceObj
}

func threadGetState(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetState: Expected 1 parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGetState: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	state := t.FieldTable["state"].Fvalue.(*object.Object)
	return state
}

// java/lang/Thread.getThreadGroup()Ljava/lang/ThreadGroup;
func threadGetThreadGroup(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGetThreadGroup: Expected 1 parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t := params[0].(*object.Object)
	threadGroup, ok := t.FieldTable["threadgroup"].Fvalue.(*object.Object)
	if !ok {
		errMsg := "threadGetThreadGroup: Expected threadgroup to be an object"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	return threadGroup
}

func threadIsAlive(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadIsAlive: Expected 1 parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadIsAlive: Expected thread to be an object"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	state := GetThreadState(t)
	if state > NEW && state < TERMINATED {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func threadIsInterrupted(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadIsInterrupted: Expected 1 parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadIsInterrupted: Expected thread to be an object"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	return t.FieldTable["interrupted"].Fvalue
}

func threadIsTerminated(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadIsTerminated: Expected 1 parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadIsTerminated: Expected thread to be an object"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	state := GetThreadState(t)
	return state == TERMINATED
}

func threadJoin(params []interface{}) any {
	fStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadJoin: Expected context data to be a frame stack"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	frame := *fStack.Front().Value.(*frames.Frame)
	thID := frame.Thread
	currentThread := globals.GetGlobalRef().Threads[thID].(*object.Object)

	targetThread, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadJoin: Missing/erroneous target thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	millis := int64(0)
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

func threadRun(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadRun: Expected only the thread object parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadRun: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name := t.FieldTable["name"].Fvalue.(*object.Object)
	id := t.FieldTable["ID"].Fvalue.(int64)
	warnMsg := fmt.Sprintf("threadRun name:%s, ID: %d started", object.GoStringFromStringObject(name), id)
	trace.Warning(warnMsg)
	return nil
}

// "java/lang/Thread.setName(Ljava/lang/String;)V"
func threadSetName(params []interface{}) any {
	// Expect exactly two parameters: the thread object and the Java String name
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadSetName: Expected 2 parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the first parameter is the Thread object
	th, ok := params[0].(*object.Object)
	if !ok || object.IsNull(th) {
		errMsg := "threadSetName: Expected first parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the second parameter is the Java String object (non-null)
	nameObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(nameObj) {
		errMsg := "threadSetName: name must not be null"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Update the thread's name field (stored as a Java byte string)
	th.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: nameObj}

	return nil
}

// "java/lang/Thread.setPriority(I)V"
func threadSetPriority(params []interface{}) any {
	// Expect exactly two parameters: the Thread object and the priority (int64)
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadSetPriority: Expected 2 parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the first parameter is the Thread object
	th, ok := params[0].(*object.Object)
	if !ok || object.IsNull(th) {
		errMsg := "threadSetPriority: Expected first parameter to be a non-null Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract and validate the second parameter (priority as int64)
	priority, ok := params[1].(int64)
	if !ok {
		errMsg := "threadSetPriority: priority must be an int64 (long)"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Fetch bounds from statics (java/lang/Thread.MIN_PRIORITY and MAX_PRIORITY)
	minP := statics.GetStaticValue("java/lang/Thread", "MIN_PRIORITY").(int64)
	maxP := statics.GetStaticValue("java/lang/Thread", "MAX_PRIORITY").(int64)

	if priority < minP || priority > maxP {
		errMsg := fmt.Sprintf("threadSetPriority: priority %d out of range [%d..%d]", priority, minP, maxP)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Update the thread's priority field (Jacobin stores it as an int64 under type types.Int)
	th.FieldTable["priority"] = object.Field{Ftype: types.Int, Fvalue: priority}
	return nil
}

// "java/lang/Thread.sleep(J)V"
func threadSleep(params []interface{}) interface{} {
	sleepTime, ok := params[0].(int64)
	if !ok {
		errMsg := "threadSleep: Parameter must be an int64 (long)"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return nil
}

func threadStart(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadStart: Expected only the thread object parameter, got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get thread object.
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadStart: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Determine receiver: target Runnable or thread itself
	var runnable *object.Object
	var clName, methName, methType string
	f, ok := t.FieldTable["target"]
	if ok && f.Fvalue != nil {
		runnable, ok = f.Fvalue.(*object.Object)
		if !ok {
			errMsg := "threadStart: Expected Runnable target field to be an object"
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		ftbl := runnable.FieldTable
		clName = object.GoStringFromJavaByteArray(ftbl["clName"].Fvalue.([]types.JavaByte))
		methName = object.GoStringFromJavaByteArray(ftbl["methName"].Fvalue.([]types.JavaByte))
		methType = object.GoStringFromJavaByteArray(ftbl["methType"].Fvalue.([]types.JavaByte))
	} else {
		clName = object.GoStringFromJavaByteArray(t.FieldTable["clName"].Fvalue.([]types.JavaByte))
		methName = "run"
		methType = "()V"
	}

	// Spawn RunJavaThread to interpret bytecode of run()
	args := []interface{}{t, clName, methName, methType}
	go globals.GetGlobalRef().FuncRunThread(args)
	runtime.Gosched()

	return nil
}

func threadYield([]interface{}) interface{} {
	runtime.Gosched()
	return nil
}

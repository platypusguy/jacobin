/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/thread"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"sync"
	"time"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

type ThreadGroup struct {
	Name string
}

type PrivateFields struct {
	Target                   interface{}
	ThreadLocals             map[string]interface{}
	InheritableLocals        map[string]interface{}
	UncaughtExceptionHandler func(thread *PublicFields, err error)
	ContextClassLoader       interface{}
	StackTrace               []string
	ParkBlocker              interface{}
	NativeThreadID           int64
	Alive                    bool
	Interrupted              bool
	Holder                   interface{}  // Added previously missing `holder` field
	Daemon                   bool         // Reflects the `daemon` field
	Priority                 int          // Reflects the `priority` field
	ThreadGroup              *ThreadGroup // Reflects the `group` field
	Name                     string       // Reflects the `name` field
	Started                  bool         // Reflects the `started` field
	Stillborn                bool         // Reflects the `stillborn` field
	Interruptible            bool         // Reflects the `interruptible` field
}

type PublicFields struct {
	ID          int64
	Name        string
	Priority    int
	IsDaemon    bool
	ThreadGroup *ThreadGroup
	State       string // Enum-like representation of Thread.State
}

func Load_Lang_Thread() {

	// constructors (followed by alpha list of public methods)
	MethodSignatures["java/lang/Thread.Thread()Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadCreateNoarg,
		}

	MethodSignatures["java/lang/Thread.Thread(Ljava/lang/String;)Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  threadCreateWithName,
		}

	MethodSignatures["java/lang/Thread.Thread(Ljava/lang/Runnable;Ljava/lang/String;)Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  threadCreateWithRunnableAndName,
		}

	// remaining methods are in alpha order by Java FQN string

	MethodSignatures["java/lang/Thread.activeCount()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadActiveCount,
		}

	MethodSignatures["java/lang/Thread.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Thread.clone()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  cloneNotSupportedException,
		}

	MethodSignatures["java/lang/Thread.currentThread()Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots:   0,
			GFunction:    threadCurrentThread,
			NeedsContext: true,
		}

	MethodSignatures["java/lang/Thread.dumpStack()V"] =
		GMeth{
			ParamSlots:   0,
			GFunction:    threadDumpStack,
			NeedsContext: true,
		}

	MethodSignatures["java/lang/Thread.enumerate([Ljava/lang/Thread;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  threadEnumerate,
		}

	MethodSignatures["java/lang/Thread.getName()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadGetName,
		}

	MethodSignatures["java/lang/Thread.getPriority()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadGetPriority,
		}

	MethodSignatures["java/lang/Thread.getNextThreadIdOffset()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.getStackTrace()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.getThreadGroup()Ljava/lang/ThreadGroup;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.holdsLock(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Thread.interrupt()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.interrupted()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Thread.run()V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  run,
		}

	MethodSignatures["java/lang/Thread.setName(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.setPriority(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.setScopedValueCache([Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.sleep(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  threadSleep,
		}

	MethodSignatures["java/lang/Thread.sleepNanos(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Thread.start()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}
	MethodSignatures["java/lang/Thread.ThreadNumbering()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadNumbering,
		}

	MethodSignatures["java/lang/Thread.ThreadNumberingNext()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadNumberingNext,
		}

	MethodSignatures["java/lang/Thread.yield()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}
}

var classname = "java/lang/Thread"

func threadActiveCount(_ []interface{}) any {
	return int64(len(globals.GetGlobalRef().Threads))
}

func threadCreateNoarg(_ []interface{}) any {

	t := object.MakeEmptyObjectWithClassName(&classname)

	idField := object.Field{Ftype: types.Int,
		Fvalue: threadNumberingNext(nil).(int64)}
	t.FieldTable["ID"] = idField

	// the JDK defaults to "Thread-N" where N is the thread number
	// the sole exception is the main thread, which is called "main"
	defaultName := fmt.Sprintf("Thread-%d", idField.Fvalue)
	nameField := object.Field{Ftype: types.GolangString, Fvalue: defaultName}
	t.FieldTable["name"] = nameField

	stateField := object.Field{Ftype: types.Int, Fvalue: thread.NEW}
	t.FieldTable["state"] = stateField

	daemonFiled := object.Field{
		Ftype: types.Int, Fvalue: types.JavaBoolFalse}
	t.FieldTable["daemon"] = daemonFiled

	threadGroup := object.Field{
		Ftype: types.Ref, Fvalue: nil}
	t.FieldTable["threadgroup"] = threadGroup

	priority := object.Field{
		Ftype: types.Int, Fvalue: int64(thread.NORM_PRIORITY)}
	t.FieldTable["priority"] = priority

	frameStack := object.Field{
		Ftype: types.LinkedList, Fvalue: nil}
	t.FieldTable["framestack"] = frameStack

	// task is the runnable that is executed if the run() method is called
	t.FieldTable["task"] = object.Field{Ftype: types.Ref, Fvalue: nil}

	return t
}

func threadCreateWithName(params []interface{}) any {
	t := threadCreateNoarg(nil).(*object.Object)
	t.FieldTable["name"] = object.Field{
		Ftype: types.GolangString, Fvalue: params[0].(string)}
	return t
}

func threadCreateWithRunnable(params []interface{}) any {
	t := threadCreateNoarg(nil).(*object.Object)
	t.FieldTable["task"] = object.Field{
		Ftype: types.Ref, Fvalue: params[0].(*object.Object)}
	return t
}

func threadCreateWithRunnableAndName(params []interface{}) any {
	t := threadCreateNoarg(nil).(*object.Object)
	t.FieldTable["task"] = object.Field{
		Ftype: types.Ref, Fvalue: params[0].(*object.Object)}
	t.FieldTable["name"] = object.Field{
		Ftype: types.GolangString, Fvalue: params[1].(string)}
	return t
}

// "java/lang/Thread.currentThread()Ljava/lang/Thread;"
func threadCurrentThread(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("CurrentThread: Expected context data, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	fStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "CurrentThread: Expected context data to be a frame"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	frame := fStack.Front().Value.(*frames.Frame)
	thID := frame.Thread
	th := globals.GetGlobalRef().Threads[thID].(*object.Object)
	return th
}

// java/lang/Thread.dumpStack()V
func threadDumpStack(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("DumpStack: Expected context data, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jvmStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "DumpStack: Expected context data to be a frame"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	globalRef := globals.GetGlobalRef()
	if globalRef.StrictJDK { // if strictly following HotSpot, ...
		_, _ = fmt.Fprintln(os.Stderr, "java.lang.Exception: Stack trace")
	} else { // TODO: add the source line numbers to both variants
		// we print more data than HotSpot does, starting with the thread name
		threadID := jvmStack.Front().Value.(*frames.Frame).Thread
		th := globalRef.Threads[threadID].(*object.Object)
		threadName := th.FieldTable["name"].Fvalue.(string)
		_, _ = fmt.Fprintf(os.Stderr, "Stack trace (thread %s)\n", threadName)
	}

	for e := jvmStack.Front(); e != nil; e = e.Next() {
		fr := e.Value.(*frames.Frame)
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
func threadEnumerate(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("getName: Expected no parameters, got %d parameters", len(params))
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

// "java/lang/Thread.getName()Ljava/lang/String;"
func threadGetName(params []interface{}) any {
	if len(params) != 0 {
		errMsg := fmt.Sprintf("getName: Expected no parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t := params[0].(*object.Object)
	return t.FieldTable["name"].Fvalue
}

// "java/lang/Thread.getPriority()I"
func threadGetPriority(params []interface{}) any {
	if len(params) != 0 {
		errMsg := fmt.Sprintf("getName: Expected no parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t := params[0].(*object.Object)
	return t.FieldTable["priority"].Fvalue
}

// "java/lang/Thread.run()V" This is the function for starting a thread. In sequence:
// 1. Fetch the run method
// 2. Create the frame stack
// 3. Create the frame
// 4. Push the frame onto the frame stack
// 5. Register the thread
// 6. Instantiate the class
// 7. Run the thread

func run(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("Run: Expected thread parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t := params[0].(*object.Object)
	runObj := t.FieldTable["task"].Fvalue
	// if the runnable is nil, then just return (per the JDK spec)
	if runObj == nil {
		return nil
	}

	// get the method to run (identified by Runnable's three fields)
	runnable := *runObj.(*object.Object)
	runFields := runnable.FieldTable
	clName := runFields["clName"].Fvalue.(string)
	methName := runFields["methName"].Fvalue.(string)
	methType := runFields["signature"].Fvalue.(string)

	m, err := classloader.FetchMethodAndCP( // resume here, with _ replaced by meth
		clName, methName, methType)

	if err != nil {
		errMsg := fmt.Sprintf("Run: Could not find run method: %v", err)
		return getGErrBlk(excNames.NoSuchMethodError, errMsg)
	}

	tID := t.FieldTable["ID"].Fvalue.(int64)
	meth := m.Meth.(classloader.JmEntry)
	f := frames.CreateFrame(meth.MaxStack + types.StackInflator) // experiment with stack size. See JACOBIN-494
	f.Thread = int(tID)
	f.ClName = clName
	f.MethName = methName
	f.MethType = methType

	f.CP = meth.Cp                        // add its pointer to the class CP
	f.Meth = append(f.Meth, meth.Code...) // copy the bytecodes over

	// allocate the local variables
	for k := 0; k < meth.MaxLocals; k++ {
		f.Locals = append(f.Locals, 0)
	}

	if tID == 1 { // if thread is the main thread, then load the CLI args into the first local
		var objArray []*object.Object
		for _, str := range globals.GetGlobalRef().AppArgs {
			sobj := object.StringObjectFromGoString(str)
			objArray = append(objArray, sobj)
		}
		f.Locals[0] = object.MakePrimitiveObject("[Ljava/lang/String", types.RefArray, objArray)
	}

	t.FieldTable["frame"] = object.Field{Ftype: types.Ref, Fvalue: f}
	fs := frames.CreateFrameStack()
	t.FieldTable["framestack"] = object.Field{Ftype: types.LinkedList, Fvalue: fs}

	if frames.PushFrame(fs, f) != nil {
		errMsg := fmt.Sprintf("Memory error allocating frame on thread: %d", tID)
		exceptions.ThrowEx(excNames.OutOfMemoryError, errMsg, nil)
	}

	// must first instantiate the class, so that any static initializers are run
	_, instantiateError := globals.GetGlobalRef().FuncInstantiateClass(clName, fs)
	if instantiateError != nil {
		errMsg := "Error instantiating: " + clName + ".main()"
		exceptions.ThrowEx(excNames.InstantiationException, errMsg, nil)
	}

	// threads are registered only when they are started
	thread.RegisterThread(t)

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("StartExec: class=%s, meth=%s%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, f.MethType, meth.MaxStack, meth.MaxLocals, len(meth.Code))
		trace.Trace(traceInfo)
	}

	return globals.GetGlobalRef().FuncRunThread(t)
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

func cloneNotSupportedException(_ []interface{}) interface{} {
	errMsg := "cloneNotSupportedException: Not supported for threads"
	return getGErrBlk(excNames.CloneNotSupportedException, errMsg)
}

// ========= ThreadNumbering is a private static class in java/lang/Thread

// this guarantees that the thread numbering is initialized only once
var setInitialThreadNumberingValue = sync.OnceValue(func() any {
	thread.ThreadNumber = int64(0)
	return nil
})

func threadNumbering(_ []any) any { // initialize thread numbering
	setInitialThreadNumberingValue()
	return nil
}

// avoid contention when creating threads
var threadNumberingMutex sync.Mutex

func threadNumberingNext(_ []any) any {
	threadNumberingMutex.Lock()
	thread.ThreadNumber += 1
	threadNumberingMutex.Unlock()
	return int64(thread.ThreadNumber)
}

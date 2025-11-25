/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/thread"
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
 creating the new frame and a function). All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

func Load_Lang_Thread() {

	// -------------------------
	// <clinit>
	// -------------------------
	MethodSignatures["java/lang/Thread.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadClinit,
		}

	// -------------------------
	// Constructors in invocation order
	// -------------------------
	MethodSignatures["java/lang/Thread.<init>()V"] =
		GMeth{ParamSlots: 0, GFunction: threadInitNull}

	MethodSignatures["java/lang/Thread.<init>(Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 1, GFunction: threadInitWithName}

	MethodSignatures["java/lang/Thread.<init>(Ljava/lang/Runnable;Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 2, GFunction: threadInitWithRunnableAndName}

	MethodSignatures["java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 2, GFunction: threadInitWithThreadGroupAndName}

	MethodSignatures["java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;)V"] =
		GMeth{ParamSlots: 2, GFunction: threadInitWithThreadGroupRunnable}

	MethodSignatures["java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 3, GFunction: threadInitWithThreadGroupRunnableAndName}

	// Long constructor
	args := "(Ljava/lang/ThreadGroup;" +
		"Ljava/lang/String;" +
		"I" +
		"Ljava/lang/Runnable;" +
		"J" +
		"Ljava/Security/AccessControlContext;" +
		")V"

	MethodSignatures["java/lang/Thread.<init>"+args] =
		GMeth{ParamSlots: 6, GFunction: threadInitFromPackageConstructor}

	// -------------------------
	// Methods in strict alphabetical order
	// -------------------------

	MethodSignatures["java/lang/Thread.activeCount()I"] =
		GMeth{ParamSlots: 0, GFunction: threadActiveCount}

	MethodSignatures["java/lang/Thread.blockedOn(Ljava/nio/channels/Interruptible;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.checkAccess()V"] =
		GMeth{ParamSlots: 0, GFunction: trapDeprecated}

	MethodSignatures["java/lang/Thread.clearInterrupt()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.clone()Ljava/lang/Object;"] =
		GMeth{ParamSlots: 0, GFunction: cloneNotSupportedException}

	MethodSignatures["java/lang/Thread.countStackFrames()I"] =
		GMeth{ParamSlots: 0, GFunction: trapDeprecated}

	MethodSignatures["java/lang/Thread.currentThread()Ljava/lang/Thread;"] =
		GMeth{ParamSlots: 0, GFunction: threadCurrentThread, NeedsContext: true}

	MethodSignatures["java/lang/Thread.destroy()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.dumpStack()V"] =
		GMeth{ParamSlots: 0, GFunction: threadDumpStack, NeedsContext: true}

	MethodSignatures["java/lang/Thread.enumerate([Ljava/lang/Thread;)I"] =
		GMeth{ParamSlots: 1, GFunction: threadEnumerate}

	MethodSignatures["java/lang/Thread.exit()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.getContextClassLoader()Ljava/lang/ClassLoader;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.getContinuation()Ljdk/internal/vm/Continuation;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.getDefaultUncaughtExceptionHandler()Ljava/lang/Thread$UncaughtExceptionHandler;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.getId()J"] =
		GMeth{ParamSlots: 0, GFunction: threadGetId}

	MethodSignatures["java/lang/Thread.getName()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: threadGetName}

	MethodSignatures["java/lang/Thread.getNextThreadIdOffset()J"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.getPriority()I"] =
		GMeth{ParamSlots: 0, GFunction: threadGetPriority}

	MethodSignatures["java/lang/Thread.getStackTrace()[Ljava/lang/StackTraceElement;"] =
		GMeth{ParamSlots: 0, GFunction: threadGetStackTrace, NeedsContext: true}

	MethodSignatures["java/lang/Thread.getState()Ljava/lang/Thread$State;"] =
		GMeth{ParamSlots: 0, GFunction: threadGetState}

	MethodSignatures["java/lang/Thread.getThreadGroup()Ljava/lang/ThreadGroup;"] =
		GMeth{ParamSlots: 0, GFunction: threadGetThreadGroup}

	MethodSignatures["java/lang/Thread.getUncaughtExceptionHandler()Ljava/lang/Thread$UncaughtExceptionHandler;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.holdsLock(Ljava/lang/Object;)Z"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.interrupt()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.interrupted()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.isAlive()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.isCCLOverridden(Ljava/lang/Class;)Z"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.isDaemon()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.isInterrupted()Z"] =
		GMeth{ParamSlots: 0, GFunction: threadIsInterrupted}

	MethodSignatures["java/lang/Thread.isVirtual()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.join()V"] =
		GMeth{ParamSlots: 0, GFunction: threadJoin}

	MethodSignatures["java/lang/Thread.join(J)V"] =
		GMeth{ParamSlots: 1, GFunction: threadJoin}

	MethodSignatures["java/lang/Thread.join(JI)V"] =
		GMeth{ParamSlots: 2, GFunction: threadJoin}

	MethodSignatures["java/lang/Thread.join(Ljava/time/Duration;)Z"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.onSpinWait()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.registerNatives()V"] =
		GMeth{ParamSlots: 0, GFunction: justReturn}

	MethodSignatures["java/lang/Thread.resume()V"] =
		GMeth{ParamSlots: 0, GFunction: trapDeprecated}

	MethodSignatures["java/lang/Thread.run()V"] =
		GMeth{ParamSlots: 1, GFunction: threadRun}

	MethodSignatures["java/lang/Thread.setContextClassLoader(Ljava/lang/ClassLoader;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.setDaemon(Z)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.setDefaultUncaughtExceptionHandler(Ljava/lang/Thread$UncaughtExceptionHandler;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.setName(Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 1, GFunction: threadSetName}

	MethodSignatures["java/lang/Thread.setPriority(I)V"] =
		GMeth{ParamSlots: 1, GFunction: threadSetPriority}

	MethodSignatures["java/lang/Thread.setScopedValueCache([Ljava/lang/Object;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.setUncaughtExceptionHandler(Ljava/lang/Thread$UncaughtExceptionHandler;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.sleep(J)V"] =
		GMeth{ParamSlots: 1, GFunction: threadSleep}

	MethodSignatures["java/lang/Thread.sleepNanos(J)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/lang/Thread.start()V"] =
		GMeth{ParamSlots: 0, GFunction: threadStart}

	MethodSignatures["java/lang/Thread.stop()V"] =
		GMeth{ParamSlots: 0, GFunction: trapDeprecated}

	MethodSignatures["java/lang/Thread.stop(Ljava/lang/Throwable;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapDeprecated}

	MethodSignatures["java/lang/Thread.suspend()V"] =
		GMeth{ParamSlots: 0, GFunction: trapDeprecated}

	MethodSignatures["java/lang/Thread.ThreadNumbering()J"] =
		GMeth{ParamSlots: 0, GFunction: threadNumbering}

	MethodSignatures["java/lang/Thread.ThreadNumberingNext()J"] =
		GMeth{ParamSlots: 0, GFunction: threadNumberingNext}

	MethodSignatures["java/lang/Thread.threadId()J"] =
		GMeth{ParamSlots: 0, GFunction: threadGetId}

	MethodSignatures["java/lang/Thread.yield()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	// finalize <clinit>
	threadClinit(nil)
}

var classname = "java/lang/Thread"

func threadActiveCount(_ []interface{}) any {
	return int64(len(globals.GetGlobalRef().Threads))
}

func _threadNameGen() *object.Object {
	num := threadNumberingNext(nil).(int64)
	return object.StringObjectFromGoString(fmt.Sprintf("Thread-%d", num))
}

// our clinit method simply specifies static constants
func threadClinit(_ []interface{}) any {
	_ = statics.AddStatic("java/lang/Thread.MIN_PRIORITY",
		statics.Static{Type: types.Int, Value: int64(thread.MIN_PRIORITY)})
	_ = statics.AddStatic("java/lang/Thread.NORM_PRIORITY",
		statics.Static{Type: types.Int, Value: int64(thread.NORM_PRIORITY)})
	_ = statics.AddStatic("java/lang/Thread.MAX_PRIORITY",
		statics.Static{Type: types.Int, Value: int64(thread.MAX_PRIORITY)})
	return nil
}

// Handles package-private constructor with these parameters:
// thread group:    Ljava/lang/ThreadGroup;
// name:            Ljava/lang/String;
// characteristics: I
// task:            Ljava/lang/Runnable;
// stack size:      J (0 = ignore)
// access control   java/Security/AccessControlContext;
// Validates each parameter, then calls threadCreateWithRunnableAndName()
// passing the 4th (Runnable) and 2nd (String name) parameters, in that order.
func threadInitFromPackageConstructor(params []interface{}) any {
	const where = "threadCreateFromPackageConstructor"

	// Expect object + 6 parameters
	if len(params) != 7 {
		errMsg := fmt.Sprintf("%s: Expected thread object + 6 parameters, got %d parameters",
			where, len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var ok bool
	var th, threadGroup *object.Object
	// 0: Threadg object
	if params[0] != nil {
		if th, ok = params[0].(*object.Object); !ok {
			errMsg := fmt.Sprintf("%s: Expected first parameter to be a Thread object (or null)", where)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// 1: Threadgroup (object may be null)
	if params[1] != nil {
		if threadGroup, ok = params[1].(*object.Object); !ok {
			errMsg := fmt.Sprintf("%s: Expected first argument to be a ThreadGroup object (or null)", where)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// 2: Name (String)
	name, ok := params[2].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("%s: Expected second argument to be a String name", where)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// 3: Priority (int). Accept common integer types.
	switch params[3].(type) {
	case int, int32, int64:
		// ok; we don't use it here but we validate presence/type
	default:
		errMsg := fmt.Sprintf("%s: Expected third argument to be an int priority", where)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// 4: Runnable (object, may be null)
	var runnable *object.Object
	if params[4] != nil {
		var ok bool
		runnable, ok = params[4].(*object.Object)
		if !ok {
			errMsg := fmt.Sprintf("%s: Expected fourth argument to be a Runnable object (or null)", where)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// 5: Long (J)
	if _, ok := params[5].(int64); !ok {
		errMsg := fmt.Sprintf("%s: Expected fifth argument to be a long (int64)", where)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// 6: AccessControlContext (object, may be null)
	if params[6] != nil {
		if _, ok := params[5].(*object.Object); !ok {
			errMsg := fmt.Sprintf("%s: Expected sixth argument to be an AccessControlContext object (or null)", where)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// Delegate: threadCreateWithRunnableAndName expects [runnable, name]
	threadInitWithRunnableAndName([]interface{}{th, runnable, name})
	tg := object.Field{ // default thread group is the main thread group
		Ftype: types.Ref, Fvalue: threadGroup}
	th.FieldTable["threadgroup"] = tg
	return nil
}

// Should we need to create a thread (as in tests), here is the instantiable implementation
func ThreadCreateNoarg(_ []interface{}) any {

	t := object.MakeEmptyObjectWithClassName(&classname)

	idField := object.Field{Ftype: types.Int,
		Fvalue: threadNumberingNext(nil).(int64)}
	t.FieldTable["ID"] = idField

	// the JDK defaults to "Thread-N" where N is the thread number
	// the sole exception is the main thread, which is called "main"
	defaultName := fmt.Sprintf("Thread-%d", idField.Fvalue)
	nameField := object.Field{Ftype: types.Ref,
		Fvalue: object.StringObjectFromGoString(defaultName)}
	t.FieldTable["name"] = nameField

	stateField := object.Field{Ftype: types.Ref,
		Fvalue: threadStateCreateWithValue([]any{NEW})}
	t.FieldTable["state"] = stateField

	daemonField := object.Field{
		Ftype: types.Int, Fvalue: types.JavaBoolFalse}
	t.FieldTable["daemon"] = daemonField

	interruptedField := object.Field{
		Ftype: types.Int, Fvalue: types.JavaBoolFalse}
	t.FieldTable["interrupted"] = interruptedField

	InitializeGlobalThreadGroups()
	tg := globals.GetGlobalRef().ThreadGroups["main"].(*object.Object)
	threadGroup := object.Field{ // default thread group is the main thread group
		Ftype: types.Ref, Fvalue: tg}
	t.FieldTable["threadgroup"] = threadGroup

	priority := object.Field{
		Ftype:  types.Int,
		Fvalue: statics.GetStaticValue("java/lang/Thread", "NORM_PRIORITY").(int64)}
	t.FieldTable["priority"] = priority

	frameStack := object.Field{
		Ftype: types.LinkedList, Fvalue: nil}
	t.FieldTable["framestack"] = frameStack

	// task is the runnable that is executed if the run() method is called
	t.FieldTable["task"] = object.Field{Ftype: types.Ref, Fvalue: nil}

	return t
}

// java/lang/Thread.<init>()V
func threadInitNull(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadInitNull: Expected 1 parameter, "+
			"(the thread object), got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitNull(: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: _threadNameGen()}
	return nil
}

// java/lang/Thread.<init>(Ljava/lang/String;)V
func threadInitWithName(params []interface{}) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadInitWithName: Expected 2 parameters, "+
			"(the thread object and name), got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithName: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithName: Expected name parameter to be a String"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["name"] = object.Field{
		Ftype: types.ByteArray, Fvalue: name}
	return nil
}

func ThreadInitWithName(params []interface{}) any { // exported version
	return threadInitWithName(params)
}

// java/lang/Thread.<init>(Ljava/lang/Runnable;)V
func threadInitWithRunnable(params []interface{}) any {
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnable: Expected thread object to be created"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	runnable, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnable: Expected parameter to be a Runnable object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["task"] = object.Field{Ftype: types.Ref, Fvalue: runnable}
	t.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: _threadNameGen()}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/Runnable;Ljava/lang/String;)V
func threadInitWithRunnableAndName(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf("threadInitWithRunnableAndName: "+
			"Expected 2 parameters plus thread object, got %d parameters",
			len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnableAndName: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	runnable, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnableAndName: Expected parameter to be a Runnable object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadCreateWithRunnableAndName: Expected  parameter to be a String"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["task"] = object.Field{
		Ftype: types.Ref, Fvalue: runnable}

	t.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: name}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V
func threadInitWithThreadGroupAndName(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf("threadInitWithThreadGroupAndName: "+
			"Expected 2 parameters plus thread object, got %d parameters",
			len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupAndName: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	threadGroup, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupAndName: Expected parameter to be a ThreadGroup object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupAndName: Expected parameter to be a String"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["threadgroup"] = object.Field{
		Ftype: types.Ref, Fvalue: threadGroup}
	t.FieldTable["name"] = object.Field{
		Ftype: types.Ref, Fvalue: name}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;Ljava/lang/String;)V
func threadInitWithThreadGroupRunnable(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf("threadInitWithThreadGroupRunnable: "+
			"Expected 2 parameters plus thread object, got %d parameters",
			len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnable: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	threadGroup, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnable: Expected parameter to be a ThreadGroup object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runnable, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnable: Expected parameter to be a Runnable object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["task"] = object.Field{
		Ftype: types.Ref, Fvalue: runnable}
	t.FieldTable["threadgroup"] = object.Field{
		Ftype: types.Ref, Fvalue: threadGroup}
	t.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: _threadNameGen()}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;Ljava/lang/String;)V
func threadInitWithThreadGroupRunnableAndName(params []interface{}) any {
	if len(params) != 4 {
		errMsg := fmt.Sprintf("threadInitWithThreadGroupRunnableAndName: "+
			"Expected 3 parameters plus thread object, got %d parameters",
			len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	threadGroup, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a ThreadGroup object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runnable, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a Runnable object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[3].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a String"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t.FieldTable["task"] = object.Field{
		Ftype: types.Ref, Fvalue: runnable}
	t.FieldTable["threadgroup"] = object.Field{
		Ftype: types.Ref, Fvalue: threadGroup}
	t.FieldTable["name"] = object.Field{
		Ftype: types.Ref, Fvalue: name}

	return nil
}

// "java/lang/Thread.currentThread()Ljava/lang/Thread;"
func threadCurrentThread(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadCurrentThread: Expected context data, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	fStack, ok := params[0].(*list.List)
	if !ok {
		errMsg := "threadCurrentThread: Expected context data to be a frame"
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

func threadStart(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadStart: Expected only the thread object parameter, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadRun: Expected parameter to be a Thread object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	runObj := t.FieldTable["task"].Fvalue
	// if the runnable is nil, then just return (per the JDK spec)
	if runObj == nil {
		return nil
	}

	// get the method to run (identified by Runnable's three fields)
	runnable := *runObj.(*object.Object)
	runFields := runnable.FieldTable
	clname := runFields["clName"].Fvalue.([]types.JavaByte)
	clName := object.GoStringFromJavaByteArray(clname)
	methname := runFields["methName"].Fvalue.([]types.JavaByte)
	methName := object.GoStringFromJavaByteArray(methname)
	methtype := runFields["signature"].Fvalue.([]types.JavaByte)
	methType := object.GoStringFromJavaByteArray(methtype)

	// Run jvm/run.go::RunJavaThread(t).
	args := []interface{}{t, clName, methName, methType}
	go globals.GetGlobalRef().FuncRunThread(args)

	return nil
}

func threadJoin(params []interface{}) any {
	th := params[0].(*object.Object)
	millis := int64(0)
	if len(params) > 1 {
		millis = params[1].(int64)
		nanos := int64(0)
		if len(params) > 2 {
			nanos = params[2].(int64)
			if nanos > 0 {
				millis += 1 // not precise
			}
		}
	}
	joinThread(th, millis)
	return nil
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

	SetThreadState(t, TERMINATED)
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
	th.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: nameObj}

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

// =========== Support functions for the functions above ===========
func RegisterThread(t *object.Object) {
	glob := globals.GetGlobalRef()
	ID := int(t.FieldTable["ID"].Fvalue.(int64))
	glob.ThreadLock.Lock()
	glob.Threads[ID] = t
	glob.ThreadLock.Unlock()
}

// Set the thread state to the supplied value unconditionally.
func SetThreadState(th *object.Object, newState int) {
	thStateObj, ok := th.FieldTable["state"].Fvalue.(*object.Object)
	if !ok {
		stateField := object.Field{Ftype: types.Ref, Fvalue: threadStateCreateWithValue([]any{newState})}
		th.FieldTable["state"] = stateField
	}
	thStateObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: newState}
}

// Wait for a thread to be TERMINATED up to maxTime milliseconds.
// Returns:
// * true if the thread terminated within maxTime milliseconds
// * false if the thread terminated after maxTime milliseconds
// * getGErrBlk(IllegalArgumentException) if:
//   - the thread state is not an object
//   - the thread state is missing the value field
//   - max wait time <= 0
func joinThread(th *object.Object, maxTime int64) interface{} {
	if maxTime <= 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "joinThread: max wait time <= 0")
	}
	var t1, t2 int64
	t1 = time.Now().UnixMilli()
	for {
		thStateObj, ok := th.FieldTable["state"].Fvalue.(*object.Object)
		if !ok {
			return getGErrBlk(excNames.IllegalArgumentException, "joinThread: field state is not an object")
		}
		thStateInt, ok := thStateObj.FieldTable["value"].Fvalue.(int)
		if !ok {
			return getGErrBlk(excNames.IllegalArgumentException, "joinThread: state object is missing a value field")
		}
		if thStateInt == TERMINATED {
			return true
		}
		t2 = time.Now().UnixMilli()
		if t2-t1 >= maxTime {
			return false
		}
	}
}

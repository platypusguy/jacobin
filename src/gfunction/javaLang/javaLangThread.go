/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"sync"
)

/*
 Each object or library that has Go methods contains a reference to ghelpers.MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. ghelpers.MethodSignatures is a map whose key is the fully qualified name and
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

const (
	MIN_PRIORITY  = 1
	NORM_PRIORITY = 5
	MAX_PRIORITY  = 10
)

func Load_Lang_Thread() {

	// -------------------------
	// <clinit>
	// -------------------------
	ghelpers.MethodSignatures["java/lang/Thread.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  threadClinit,
		}

	// ================================ Constructors

	ghelpers.MethodSignatures["java/lang/Thread.<init>()V"] =
		ghelpers.GMeth{ParamSlots: 0, NeedsContext: true, GFunction: threadInitNull}

	ghelpers.MethodSignatures["java/lang/Thread.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, NeedsContext: true, GFunction: ThreadInitWithName}

	ghelpers.MethodSignatures["java/lang/Thread.<init>(Ljava/lang/Runnable;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadInitWithRunnable}

	ghelpers.MethodSignatures["java/lang/Thread.<init>(Ljava/lang/Runnable;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: threadInitWithRunnableAndName}

	ghelpers.MethodSignatures["java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, NeedsContext: true, GFunction: threadInitWithThreadGroupAndName}

	ghelpers.MethodSignatures["java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: threadInitWithThreadGroupRunnable}

	ghelpers.MethodSignatures["java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: threadInitWithThreadGroupRunnableAndName}

	args := "(Ljava/lang/ThreadGroup;" +
		"Ljava/lang/String;" +
		"I" +
		"Ljava/lang/Runnable;" +
		"J" +
		"Ljava/Security/AccessControlContext;" +
		")V"
	ghelpers.MethodSignatures["java/lang/Thread.<init>"+args] =
		ghelpers.GMeth{ParamSlots: 6, GFunction: threadInitFromPackageConstructor}

	// ============================= Member functions

	ghelpers.MethodSignatures["java/lang/Thread.activeCount()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadActiveCount}

	ghelpers.MethodSignatures["java/lang/Thread.blockedOn(Ljava/nio/channels/Interruptible;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.checkAccess()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Thread.clearInterrupt()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: cloneNotSupportedException}

	ghelpers.MethodSignatures["java/lang/Thread.countStackFrames()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Thread.currentThread()Ljava/lang/Thread;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadCurrentThread, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Thread.destroy()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.dumpStack()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadDumpStack, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Thread.enumerate([Ljava/lang/Thread;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadEnumerate}

	ghelpers.MethodSignatures["java/lang/Thread.exit()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.getContextClassLoader()Ljava/lang/ClassLoader;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.getContinuation()Ljdk/internal/vm/Continuation;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.getDefaultUncaughtExceptionHandler()Ljava/lang/Thread$UncaughtExceptionHandler;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.getId()J"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetId}

	ghelpers.MethodSignatures["java/lang/Thread.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetName}

	ghelpers.MethodSignatures["java/lang/Thread.getNextThreadIdOffset()J"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.getPriority()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetPriority}

	ghelpers.MethodSignatures["java/lang/Thread.getStackTrace()[Ljava/lang/StackTraceElement;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetStackTrace, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Thread.getState()Ljava/lang/Thread$State;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetState}

	ghelpers.MethodSignatures["java/lang/Thread.getThreadGroup()Ljava/lang/ThreadGroup;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetThreadGroup}

	ghelpers.MethodSignatures["java/lang/Thread.getUncaughtExceptionHandler()Ljava/lang/Thread$UncaughtExceptionHandler;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.holdsLock(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.interrupt()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.interrupted()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.isAlive()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadIsAlive}

	ghelpers.MethodSignatures["java/lang/Thread.isCCLOverridden(Ljava/lang/Class;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.isDaemon()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.isInterrupted()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadIsInterrupted}

	ghelpers.MethodSignatures["java/lang/Thread.isTerminated()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadIsTerminated}

	ghelpers.MethodSignatures["java/lang/Thread.isVirtual()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ReturnTrue}

	ghelpers.MethodSignatures["java/lang/Thread.join()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadJoin, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Thread.join(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadJoin, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Thread.join(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: threadJoin, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Thread.join(Ljava/time/Duration;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.onSpinWait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadYield}

	ghelpers.MethodSignatures["java/lang/Thread.registerNatives()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.JustReturn}

	ghelpers.MethodSignatures["java/lang/Thread.resume()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Thread.run()V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadRun}

	ghelpers.MethodSignatures["java/lang/Thread.setContextClassLoader(Ljava/lang/ClassLoader;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.setDaemon(Z)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.setDefaultUncaughtExceptionHandler(Ljava/lang/Thread$UncaughtExceptionHandler;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.setName(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadSetName}

	ghelpers.MethodSignatures["java/lang/Thread.setPriority(I)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadSetPriority}

	ghelpers.MethodSignatures["java/lang/Thread.setScopedValueCache([Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.setUncaughtExceptionHandler(Ljava/lang/Thread$UncaughtExceptionHandler;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.sleep(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: threadSleep}

	ghelpers.MethodSignatures["java/lang/Thread.sleepNanos(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Thread.start()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadStart}

	ghelpers.MethodSignatures["java/lang/Thread.stop()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Thread.stop(Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Thread.suspend()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Thread.ThreadNumbering()J"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadNumbering}

	ghelpers.MethodSignatures["java/lang/Thread.ThreadNumberingNext()J"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadNumberingNext}

	ghelpers.MethodSignatures["java/lang/Thread.threadId()J"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadGetId}

	ghelpers.MethodSignatures["java/lang/Thread.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ThreadToString}

	ghelpers.MethodSignatures["java/lang/Thread.yield()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: threadYield}

	// finalize <clinit>
	threadClinit(nil)
}

// Thread numbering is a static counter that increments for each thread created.
// Amendment is under the control of a mutex.
var threadNumber int64 = 0
var threadNumberingMutex sync.Mutex

// our clinit method simply specifies static constants
func threadClinit(_ []interface{}) any {
	_ = statics.AddStatic("java/lang/Thread.MIN_PRIORITY",
		statics.Static{Type: types.Int, Value: int64(MIN_PRIORITY)})
	_ = statics.AddStatic("java/lang/Thread.NORM_PRIORITY",
		statics.Static{Type: types.Int, Value: int64(NORM_PRIORITY)})
	_ = statics.AddStatic("java/lang/Thread.MAX_PRIORITY",
		statics.Static{Type: types.Int, Value: int64(MAX_PRIORITY)})
	return nil
}

// Handles package-private constructor with these parameters:
// thread group:    Ljava/lang/ThreadGroup;
// name:            Ljava/lang/String;
// characteristics: I
// target:          Ljava/lang/Runnable;
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var ok bool
	var th, threadGroup *object.Object
	// 0: Threadg object
	if params[0] != nil {
		if th, ok = params[0].(*object.Object); !ok {
			errMsg := fmt.Sprintf("%s: Expected first parameter to be a Thread object (or null)", where)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// 1: Threadgroup (object may be null)
	if params[1] != nil {
		if threadGroup, ok = params[1].(*object.Object); !ok {
			errMsg := fmt.Sprintf("%s: Expected first argument to be a ThreadGroup object (or null)", where)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// 2: Name (String)
	name, ok := params[2].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("%s: Expected second argument to be a String name", where)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// 3: Priority (int). Accept common integer types.
	switch params[3].(type) {
	case int, int32, int64:
		// ok; we don't use it here but we validate presence/type
	default:
		errMsg := fmt.Sprintf("%s: Expected third argument to be an int priority", where)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// 4: Runnable (object, may be null)
	var runnable *object.Object
	if params[4] != nil {
		var ok bool
		runnable, ok = params[4].(*object.Object)
		if !ok {
			errMsg := fmt.Sprintf("%s: Expected fourth argument to be a Runnable object (or null)", where)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// 5: Long (J)
	if _, ok := params[5].(int64); !ok {
		errMsg := fmt.Sprintf("%s: Expected fifth argument to be a long (int64)", where)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// 6: AccessControlContext (object, may be null)
	if params[6] != nil {
		if _, ok := params[5].(*object.Object); !ok {
			errMsg := fmt.Sprintf("%s: Expected sixth argument to be an AccessControlContext object (or null)", where)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	// Delegate: threadCreateWithRunnableAndName expects [runnable, name]
	threadInitWithRunnableAndName([]interface{}{th, runnable, name})
	idField := object.Field{Ftype: types.Int, Fvalue: threadNumberingNext(nil).(int64)}
	th.FieldTable["ID"] = idField

	tg := object.Field{ // default thread group is the main thread group
		Ftype: types.Ref, Fvalue: threadGroup}
	th.FieldTable["threadgroup"] = tg
	return nil
}

// java/lang/Thread.<init>()V
func threadInitNull(params []interface{}) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadInitNull: Expected 2 parameter2, "+
			"(frame stack + the thread object), got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the thread object and populate it.
	t, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitNull(: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/String;)V
func ThreadInitWithName(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf("ThreadInitWithName: Expected 2 parameters, "+
			"(the thread object and name), got %d parameters", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "ThreadInitWithName: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	// Get thread name.
	name, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "ThreadInitWithName: Expected name parameter to be a String"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: name}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/Runnable;)V
func threadInitWithRunnable(params []interface{}) any {
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnable: Expected thread object to be created"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	runnable, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnable: Expected parameter to be a Runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runClassName := object.GoStringFromStringPoolIndex(runnable.KlassName)
	ret := setUpRunnable(runnable, runClassName)
	if ret != nil {
		return ret
	}

	t.FieldTable["target"] = object.Field{Ftype: types.Ref, Fvalue: runnable}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/Runnable;Ljava/lang/String;)V
func threadInitWithRunnableAndName(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf("threadInitWithRunnableAndName: "+
			"Expected 2 parameters plus thread object, got %d parameters",
			len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnableAndName: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	runnable, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithRunnableAndName: Expected parameter to be a Runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runClassName := object.GoStringFromStringPoolIndex(runnable.KlassName)
	ret := setUpRunnable(runnable, runClassName)
	if ret != nil {
		return ret
	}

	name, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadCreateWithRunnableAndName: Expected  parameter to be a String"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["target"] = object.Field{
		Ftype: types.Ref, Fvalue: runnable}

	t.FieldTable["name"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: name}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V
func threadInitWithThreadGroupAndName(params []interface{}) any {
	if len(params) != 4 {
		errMsg := fmt.Sprintf("threadInitWithThreadGroupAndName: "+
			"Expected 2 parameters plus thread object, got %d parameters",
			len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "ThreadInitWithName: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	// Get the class name "java/lang/Thread" or the user's own subclass of Thread.
	// frameStack := params[0].(*list.List)
	// storeThreadClassName(t, frameStack)

	threadGroup, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupAndName: Expected parameter to be a ThreadGroup object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[3].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupAndName: Expected parameter to be a String"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	t.FieldTable["threadgroup"] = object.Field{Ftype: types.Ref, Fvalue: threadGroup}
	t.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: name}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;Ljava/lang/String;)V
func threadInitWithThreadGroupRunnable(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf("threadInitWithThreadGroupRunnable: "+
			"Expected 2 parameters plus thread object, got %d parameters",
			len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnable: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	threadGroup, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnable: Expected parameter to be a ThreadGroup object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runnable, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnable: Expected parameter to be a Runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runClassName := object.GoStringFromStringPoolIndex(runnable.KlassName)
	ret := setUpRunnable(runnable, runClassName)
	if ret != nil {
		return ret
	}

	t.FieldTable["target"] = object.Field{
		Ftype: types.Ref, Fvalue: runnable}
	t.FieldTable["threadgroup"] = object.Field{
		Ftype: types.Ref, Fvalue: threadGroup}

	return nil
}

// java/lang/Thread.<init>(Ljava/lang/ThreadGroup;Ljava/lang/Runnable;Ljava/lang/String;)V
func threadInitWithThreadGroupRunnableAndName(params []interface{}) any {
	if len(params) != 4 {
		errMsg := fmt.Sprintf("threadInitWithThreadGroupRunnableAndName: "+
			"Expected 3 parameters plus thread object, got %d parameters",
			len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a Thread object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	populateThreadObject(t)

	threadGroup, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a ThreadGroup object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runnable, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a Runnable object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	runClassName := object.GoStringFromStringPoolIndex(runnable.KlassName)
	ret := setUpRunnable(runnable, runClassName)
	if ret != nil {
		return ret
	}

	name, ok := params[3].(*object.Object)
	if !ok {
		errMsg := "threadInitWithThreadGroupRunnableAndName: Expected parameter to be a String"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	t.FieldTable["target"] = object.Field{
		Ftype: types.Ref, Fvalue: runnable}
	t.FieldTable["threadgroup"] = object.Field{
		Ftype: types.Ref, Fvalue: threadGroup}
	t.FieldTable["name"] = object.Field{
		Ftype: types.ByteArray, Fvalue: name}

	return nil
}

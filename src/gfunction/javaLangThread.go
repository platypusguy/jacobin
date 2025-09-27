/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/thread"
	"jacobin/src/types"
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
	// vanilla constructor
	MethodSignatures["java/lang/Thread.Thread()Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadCreateNoarg,
		}

	// constructor with name
	MethodSignatures["java/lang/Thread.Thread(Ljava/lang/String;)Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  threadCreateWithName,
		}

	MethodSignatures["java/lang/Thread.currentThread(Ljava/lang/Runnable;)Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadCreateWithRunnable,
		}

	MethodSignatures["java/lang/Thread.currentThread(Ljava/lang/Runnable;Ljava/lang/String;)Ljava/lang/Thread;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadCreateWithRunnableAndName,
		}
	MethodSignatures["java/lang/Thread.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Thread.sleep(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  threadSleep,
		}

	// various methods
	MethodSignatures["java/lang/Thread.clone()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  cloneNotSupportedException,
		}

	MethodSignatures["java/lang/Thread.run()V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  run,
		}

	// ThreadNumbering is a private static class in java/lang/Thread
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
}

var classname = "java/lang/Thread"

func threadCreateNoarg(params []interface{}) any {

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

	// task is the runnable that is executed if the run() method is called
	t.FieldTable["task"] = object.Field{Ftype: types.Ref, Fvalue: nil}

	return &t
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
	runnObj := t.FieldTable["task"].Fvalue.(*object.Object)
	runFields := runnObj.FieldTable
	_, err := classloader.FetchMethodAndCP( // resume here, with _ replaced by meth
		runFields["clName"].Fvalue.(string),
		runFields["methName"].Fvalue.(string),
		runFields["signature"].Fvalue.(string))
	if err != nil {
		errMsg := fmt.Sprintf("Run: Could not find run method: %v", err)
		return getGErrBlk(excNames.NoSuchMethodError, errMsg)
	}

	// threads are registered only when they are started
	thread.RegisterThread(t)
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

func cloneNotSupportedException(params []interface{}) interface{} {
	errMsg := "cloneNotSupportedException: Not supported for threads"
	return getGErrBlk(excNames.CloneNotSupportedException, errMsg)
}

// ========= ThreadNumbering is a private static class in java/lang/Thread
func threadNumbering(params []any) any { // initialize thread numbering
	thread.ThreadNumber = int64(0)
	return thread.ThreadNumber
}

func threadNumberingNext(params []any) any {
	thread.ThreadNumber += 1
	return int64(thread.ThreadNumber)
}

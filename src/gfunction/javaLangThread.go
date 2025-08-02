/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/excNames"
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

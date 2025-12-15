/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/trace"
	"jacobin/src/util"
	"os"
	"runtime/debug"
)

// This file contains the functions for throwing exceptions from within
// Jacobin. That is, situations in which Jacobin itself is throwing the error,
// rather than the application. Typically, this is for errors/exceptions in the
// operation of the JVM and for a few occasional user errors, such as
// divide by zero.

const (
	Caught    = true
	NotCaught = false
)

// ThrowExNil simply calls ThrowEx with a nil pointer for the frame.
func ThrowExNil(which int, msg string) bool {
	return ThrowEx(which, msg, nil)
}

// ThrowEx throws an exception. It is used primarily for exceptions and
// errors thrown by Jacobin, rather than by the application. (The latter
// would generally use the ATHROW bytecode.)
//
// Important: if you change the name of this function, you need to update
// exceptions.ShowGoStackTrace(), which explicitly tests for this function name.
func ThrowEx(which int, msg string, f *frames.Frame) bool {
	if globals.TraceVerbose {
		infoMsg := fmt.Sprintf("[ThrowEx] %s, msg: %s", excNames.JVMexceptionNames[which], msg)
		trace.Trace(infoMsg)
	}

	// If in a unit test, log a severe message and return.
	glob := globals.GetGlobalRef()
	if glob.JacobinName == "test" {
		var errMsg string
		if f != nil {
			errMsg = fmt.Sprintf("%s in %s, %s", excNames.JVMexceptionNames[which], frames.FormatFQN(f), msg)
		} else { // if we got here by a call to ThrowExNil()
			errMsg = fmt.Sprintf("%s: %s", excNames.JVMexceptionNames[which], msg)
		}
		_, _ = fmt.Fprintln(os.Stderr, errMsg)
		return NotCaught
	}

	// Frame pointer provided?
	if f == nil {
		MinimalAbort(which, msg) // this calls exit()
		return NotCaught         // only occurs in tests
	}

	// the name of the exception as shown to the user
	var exceptionNameForUser string
	if glob.StrictJDK {
		exceptionNameForUser = excNames.JVMexceptionNames[which]
	} else {
		exceptionNameForUser = excNames.JVMexceptionNamesJacobin[which]
	}

	// the internal format used in the constant pool
	exceptionCPname := util.ConvertClassFilenameToInternalFormat(excNames.JVMexceptionNames[which])

	// capture the PC where the exception was thrown, if it hasn't been captured yet.
	// (saved b/c later we modify the value of f.PC)
	if f.ExceptionPC == -1 {
		f.ExceptionPC = f.PC
	}

	var fs *list.List
	glob.ThreadLock.RLock()
	th, ok := glob.Threads[f.Thread].(*object.Object)
	glob.ThreadLock.RUnlock()
	if !ok {
		errMsg := fmt.Sprintf("[ThrowEx] glob.Threads index not found or entry corrupted, thread index: %d, FQN: %s",
			f.Thread, frames.FormatFQN(f))
		MinimalAbort(excNames.InternalException, errMsg)
	}
	fs = th.FieldTable["framestack"].Fvalue.(*list.List)

	// find out if the exception is caught and if so point to the catch code
	catchFrame, catchPC := FindCatchFrame(fs, exceptionCPname, f.ExceptionPC)
	if catchFrame != nil {
		// at this point, we know that the exception was caught
		// and that the returned frame is the frame
		// containing the catch logic, referred to here as the catchFrame.
		// now, set up the execution of the catch code by:
		// 0. popping off the frames that are above the catch frame,
		//    if any--so that top frame in the frame stack is the catch frame
		// 1. creating a new objRef for the exception
		// 2. pushing the objRef on the op stack of the frame
		// 3. setting the PC to point to the catch code (which expects the objRef at TOS)
		if globals.TraceVerbose {
			infoMsg := fmt.Sprintf("[ThrowEx] caught %s, FQN: %s, msg: %s", exceptionCPname, frames.FormatFQN(f), msg)
			trace.Trace(infoMsg)
		}

		// remove the frames we examined that did not have the catch logic
		for fs.Len() > 0 {
			fr := fs.Front().Value
			if fr == catchFrame {
				break
			} else {
				fs.Remove(fs.Front())
			}
		}

		// Set up the catch frame.
		objRef, _ := glob.FuncInstantiateClass(exceptionCPname, fs)

		// Fill in the stack trace for the exception object
		throwObj := objRef.(*object.Object)
		params := []any{fs, throwObj}
		glob.FuncFillInStackTrace(params)

		catchFrame.TOS = 0
		catchFrame.OpStack[0] = objRef // push the objRef
		catchFrame.PC = catchPC

		// the exception logic might throw another exception, in which case that will be
		// the new ExceptionPC. However, it won't be updated to that value unless ExceptionPC
		// is reset to -1. So, at this point, the exception has been caught, so we can reset
		// ExeptionPC to -1. See JACOBIN-534
		f.ExceptionPC = -1
		return Caught
	}

	// ---- if exception is not caught ----

	throwObject, err := glob.FuncInstantiateClass(exceptionCPname, fs)
	if err != nil {
		fmt.Printf("InstantiateClass failed, FQN: %s, %s", frames.FormatFQN(f), err.Error())
		// if throwObject != nil {
		_, _ = fmt.Fprintf(os.Stderr, "throwObject: %v\n", throwObject)
		_ = shutdown.Exit(shutdown.JVM_EXCEPTION)
		// }
	}

	throwObj := throwObject.(*object.Object)
	params := []any{fs, throwObj}
	glob.FuncFillInStackTrace(params)

	excInfo := fmt.Sprintf("%s: FQN: %s, %s", exceptionNameForUser, frames.FormatFQN(f), msg)
	_, _ = fmt.Fprintln(os.Stderr, excInfo)

	stackTrace := throwObj.FieldTable["stackTrace"].Fvalue.(*object.Object)
	traceEntries := stackTrace.FieldTable["value"].Fvalue.([]*object.Object)

	// now print out the JVM stack
	for _, traceEntry := range traceEntries {
		// HotSpot uses a slightly different format for method names:
		// package.class.method, we prefer package/class.method, so we format
		// method name according to whether -strictJDK is in force
		var declaringClass string
		if glob.StrictJDK {
			declaringClass = util.ConvertInternalClassNameToUserFormat(
				traceEntry.FieldTable["declaringClass"].Fvalue.(string))
		} else {
			declaringClass = traceEntry.FieldTable["declaringClass"].Fvalue.(string)
		}

		var traceInfo string
		if traceEntry.FieldTable["sourceLine"].Fvalue.(string) == "" {
			// if no source line number is available, just print the method name
			traceInfo = fmt.Sprintf("  at %s.%s(%s)",
				declaringClass,
				traceEntry.FieldTable["methodName"].Fvalue.(string),
				traceEntry.FieldTable["fileName"].Fvalue.(string))
		} else {
			traceInfo = fmt.Sprintf("  at %s.%s(%s:%s)",
				// otherwise, print the method name and source file + line number
				declaringClass,
				traceEntry.FieldTable["methodName"].Fvalue.(string),
				traceEntry.FieldTable["fileName"].Fvalue.(string),
				traceEntry.FieldTable["sourceLine"].Fvalue.(string))
		}
		trace.AsIs(traceInfo)
	}

	if !glob.StrictJDK {
		// the next statement disables showing the line that identifies
		// the cause of a golang panic, because if we got here, there
		// was no panic, rather just an uncaught exception. So we show
		// the golang stack without implying there was a panic.
		glob.PanicCauseShown = true
		ShowGoStackTrace("")
	}

	_ = shutdown.Exit(shutdown.JVM_EXCEPTION) // in test mode, this call returns
	return NotCaught                          // only applies to tests
}

// MinimalAbort is the exception thrown when the frame info is not available,
// such as during start-up, if the main class can't be found, etc.
func MinimalAbort(whichException int, msg string) {
	var stack string
	bytes := debug.Stack()
	if len(bytes) > 0 {
		stack = string(bytes)
	} else {
		stack = ""
	}
	glob := globals.GetGlobalRef()
	glob.ErrorGoStack = stack
	errMsg := fmt.Sprintf("%s: %s", excNames.JVMexceptionNames[whichException], msg)
	_, _ = fmt.Fprintln(os.Stderr, errMsg)
	ShowGoStackTrace(nil)
	_ = shutdown.Exit(shutdown.APP_EXCEPTION)
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/trace"
	"slices"
)

// Execution of Java native functions.
// Called by the INVOKE* op codes in run.go.
// This effort begins with the JACOBIN-582 task.
//
// Parameters:
// * fs : the frame stack of the current thread
// * className : class name
// * functionName : function name (method name in the jvm package)
// * methodType : parameters and return value expressed as (parameter-types)value-type
// * params : a slice of parameters being passed to the method
// * tracing : a boolean such that when true, trace-prints should be performed
//
// Returns :
// * a value if the native function returned a value (success)
// * an error if the called native function returned an error indication
// * an errorBlock if an exception occurred
//
// Note that RunNativeFunction will determine whether a native function is supported (yet).
// If not yet supported, then the following will happen:
//
// errMsg := "RunNativeFunction: Unsupported native method requested: " + className + "." + methodName + methodType
// status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, frame)
// if status != exceptions.Caught {
//      return errors.New(errMsg) // applies only if in test
// }

func RunNativeFunction(fs *list.List, className, nativeFunctionName, methodType string, params *[]interface{}, tracing bool) interface{} {

	frame := fs.Front().Value.(*frames.Frame)

	// Compute the parameter count.
	var paramCount int
	if params == nil {
		paramCount = 0
	} else {
		paramCount = len(*params)
	}

	// Form the full method name.
	fullMethName := fmt.Sprintf("%s.%s%s", className, nativeFunctionName, methodType)
	if tracing {
		infoMsg := fmt.Sprintf("RunNativeFunction: %s, paramSlots: %d", fullMethName, paramCount)
		trace.Trace(infoMsg)
		// TODO jvm.LogTraceStack(templateFunction)
	}

	// Reverse the parameter order. Last appended will be fetched first.
	if paramCount > 0 {
		slices.Reverse(*params)
	}

	// Discern between thread-safe G functions and ordinary ones.
	// No matter what, ret = the result from the G function.
	var ret any

	/*

	   During initialization (not part of this function),
	   * The NfLibXrefTable is built by either a POSIX loader or a Windows loader. Note that both the library path and handle are populated.
	   * The nfToTmplTable remains nil.

	   At run-time, RunNativeFunction will do the following in order to get (1) a native function handle
	   and (2) the corresponding template function address:
	   * Look up the funcName in the nfToTmplTable.
	   * If not found,
	        - Look up funcName in nfToLibTable. Not found ---> error.
	        - Derive the template function to use for this methodName based on the methodType.
	        - Store the template function handle in nfToTmplTable.
	   * Call the template function (by address) with arguments: library handle and the function name.

	*/

	var templateFunction typeTemplateFunction
	var libHandle uintptr
	var ok bool

	// Get the library handle.
	libHandle, ok = nfToLibTable[nativeFunctionName]
	if !ok {
		errMsg := fmt.Sprintf("RunNativeFunction: Function %s is not in the function-to-library-table", nativeFunctionName)
		return NativeErrBlk{ExceptionType: excNames.VirtualMachineError, ErrMsg: errMsg}
	}

	// Get the template function handle.
	templateFunction, ok = nfToTmplTable[nativeFunctionName]
	if !ok {

		// Does not yet have a template function handle.
		// Get the template function handle.
		templateFunction, ok = mapToTemplateHandle(methodType)
		if !ok {
			errMsg := fmt.Sprintf("RunNativeFunction: mapToTemplateHandle(%s) not found", methodType)
			return NativeErrBlk{ExceptionType: excNames.VirtualMachineError, ErrMsg: errMsg}
		}

		// Update nfToTmplTable with the template function handle for the next time.
		nfToTmplTable[nativeFunctionName] = templateFunction

	}

	// Call the template function.
	ret = templateFunction(libHandle, nativeFunctionName, *params, tracing)

	// Check the type of function completion.
	switch ret.(type) {

	case *NativeErrBlk: // Native error block was returned.
		// Convenience capture of the error block.
		errBlk := *ret.(*NativeErrBlk)

		// Get the thread name.
		var threadName string
		if frame.Thread == 1 {
			threadName = "main"
		} else {
			threadName = fmt.Sprintf("%d", frame.Thread)
		}

		// Build the exception message and return it.
		errMsg := fmt.Sprintf("%s in thread: %s", errBlk.ErrMsg, threadName)
		status := exceptions.ThrowEx(errBlk.ExceptionType, errMsg, frame)
		if status != exceptions.Caught {
			return errors.New(errMsg + " " + errBlk.ErrMsg) // applies only if in test
		} else {
			// if the exception was caught, tell calling function to execute the catching logic
			return CaughtNativeFunctionException
		}

	case error: // Go error object returned.
		// Build the error message for the native function failure and return it.
		errMsg := (ret.(error)).Error()
		status := exceptions.ThrowEx(excNames.NativeMethodException, errMsg, frame)
		if status != exceptions.Caught {
			return ret.(error) // applies only if in test
		} else {
			return nil // return nothing if the error was caught
		}
	}

	// At this point, we have a successful outcome.
	// Return the value generated by the native function.
	return ret
}

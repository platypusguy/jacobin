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
	"jacobin/log"
	"slices"
)

// Native function error block.
type NativeErrBlk struct {
	ExceptionType int
	ErrMsg        string
}

var CaughtNativeFunctionException = errors.New("caught native function exception")

// Execution of Java native functions.
// Called by the INVOKE* op codes in run.go.
// This effort begins with the JACOBIN-582 task.
//
// Parameters:
// * fs : the frame stack of the current thread
// * className, methodName : class and method names
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

func RunNativeFunction(fs *list.List, className, methodName, methodType string, params *[]interface{}, tracing bool) any {

	f := fs.Front().Value.(*frames.Frame)

	// Compute the parameter count.
	var paramCount int
	if params == nil {
		paramCount = 0
	} else {
		paramCount = len(*params)
	}

	// Form the full method name.
	fullMethName := fmt.Sprintf("%s.%s%s", className, methodName, methodType)
	if tracing {
		traceInfo := fmt.Sprintf("RunNativefunction: %s, paramSlots: %d", fullMethName, paramCount)
		_ = log.Log(traceInfo, log.TRACE_INST)
		// TODO jvm.LogTraceStack(f)
	}

	// Reverse the parameter order. Last appended will be fetched first.
	if paramCount > 0 {
		slices.Reverse(*params)
	}

	// Discern between thread-safe G functions and ordinary ones.
	// No matter what, ret = the result from the G function.
	var ret any

	// ****************************************************
	// TODO Figure out how to call the native function.
	// ****************************************************
	// Select a template function based on method type.
	// Each template function uses purego and can return
	// one of the following:
	// - a value computed by the native function (success)
	// - an error returned by the native function (failure)
	// - a pointer to a NativeErrBlk (exception occurred)
	// ****************************************************
	// TODO Integrate this with tables under investigation by ALB
	// ****************************************************

	// Check the type of function completion.
	switch ret.(type) {

	case *NativeErrBlk:
		// Convenience capture of the error block.
		errBlk := *ret.(*NativeErrBlk)

		// Get the thread name.
		var threadName string
		if f.Thread == 1 {
			threadName = "main"
		} else {
			threadName = fmt.Sprintf("%d", f.Thread)
		}

		// Build the exception message and return it.
		errMsg := fmt.Sprintf("%s in thread: %s, method: %s",
			errBlk.ErrMsg, threadName, fullMethName)
		status := exceptions.ThrowEx(errBlk.ExceptionType, errMsg, f)
		if status != exceptions.Caught {
			return errors.New(errMsg + " " + errBlk.ErrMsg) // applies only if in test
		} else {
			// if the exception was caught, tell calling function to execute the catching logic
			return CaughtNativeFunctionException
		}

	case error:
		// Build the error message for the native function failure and return it.
		errMsg := (ret.(error)).Error()
		status := exceptions.ThrowEx(excNames.NativeMethodException, errMsg, f)
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

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/trace"
	"slices"
)

var CaughtGfunctionException = errors.New("caught gfunction exception")

// Execution of gfunctions (that is, Java functions ported to golang).
// As part of JACOBIN-519, this code seeks to replace the previous set of
// functions (e.g., runGframe() and runGmethod()) with a simpler streamlined
// function.
//
// Parameters: mt is the gmethod we're about to run, fs is the frame stack
// className, methodName, methodType should be self-explanatory, params is
// a slice of parameters/arguments being passed to the gmethod, and objRef
// is a boolean indicating whether a pointer to the object whose method is
// being called was pushed on to the stack (true) or not (false)
//
// Returns an errorBlock if an exception occured, an error if the gfunction
// returned an error but did not throw an exception, or a value if the
// gfunction returned a value.

func RunGfunction(mt classloader.MTentry, fs *list.List,
	className, methodName, methodType string,
	params *[]interface{}, objRef bool, tracing bool) any {

	f := fs.Front().Value.(*frames.Frame)

	// If the method needs context (i.e., if mt.Meth.NeedsContext == true),
	// then add the pointer to the JVM frame stack to the parameter list here.
	entry := mt.Meth.(ghelpers.GMeth)
	if entry.NeedsContext {
		*params = append(*params, fs)
	}

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
		infoMsg := fmt.Sprintf("RunGfunction: %s, objectRef: %v, paramSlots: %d",
			fullMethName, objRef, paramCount)
		trace.Trace(infoMsg)
		// TODO jvm.LogTraceStack(f)
	}

	// Reverse the parameter order. The last appended will be fetched first.
	if paramCount > 1 {
		slices.Reverse(*params)
	}

	// No matter what, ret = the result from the G function.
	var ret any
	gmeth := mt.Meth.(ghelpers.GMeth)

	// Call the G function, passing it a pointer to the slice of arguments.
	if paramCount == 0 {
		ret = gmeth.GFunction(nil)
	} else {
		ret = gmeth.GFunction(*params)
	}

	// if an error occured
	switch ret.(type) {
	case *ghelpers.GErrBlk:
		// var errorDetails string
		errBlk := *ret.(*ghelpers.GErrBlk)

		var threadName string
		if f.Thread == 1 {
			threadName = "main"
		} else {
			threadName = fmt.Sprintf("%d", f.Thread)
		}
		if f.Thread == 0 {
			errMsg := fmt.Sprintf("in main thread initialization, %s reported by G-function: %s", errBlk.ErrMsg, fullMethName)
			exceptions.MinimalAbort(errBlk.ExceptionType, errMsg)
		}
		errMsg := fmt.Sprintf("in thread: %s, in thread %s, reported by G-function: %s", errBlk.ErrMsg, threadName, fullMethName)
		status := exceptions.ThrowEx(errBlk.ExceptionType, errMsg, f)
		if status != exceptions.Caught {
			return errors.New(errMsg + " " + errBlk.ErrMsg) // applies only if in test
		}

		// if the exception was caught, tell calling function to execute the catching logic
		return CaughtGfunctionException

	case error:
		errMsg := (ret.(error)).Error()
		status := exceptions.ThrowEx(excNames.NativeMethodException, errMsg, f)
		if status != exceptions.Caught {
			return ret.(error) // applies only if in test
		}
		return nil // return nothing if the error was caught
	}

	// if return is not an errBlk or an error, then it's a legitimate
	// return value, so return it.
	return ret
}

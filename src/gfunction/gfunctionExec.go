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
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"slices"
)

var CaughtGfunctionException = errors.New("caught gfunction exception")

// RunGfunction - Execution of gfunctions (that is, Java functions ported to Golang).
//
// Parameters:
// * mt holds the MTentry we're about to run
// * fs is the frame stack
// * params is a pointer to a slice of parameters/arguments being passed to the gmethod
// * objRef is a boolean indicating whether a pointer to the object whose method is
//   being called was included in the parameters
// * tracing is a boolean governing G function executive tracing
//
// Returns normally:
// * int64 (Java boolean, byte, char, short, int, or long)
// * float64 (Java float or double)
// * *object.Object (Java object)
// * nil (Java void)
//
// Returns a Golang error struct, indicating an exception condition has occurred:
// * CaughtGfunctionException (a special error instance) if the exception should be caught
// * a Golang error struct with some other value if the exception cannot be caught

func RunGfunction(mt classloader.MTentry, fs *list.List,
	params *[]any, objRef bool, tracing bool) any {

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
	if tracing {
		infoMsg := fmt.Sprintf("RunGfunction: %s.%s%s, objectRef: %v, paramSlots: %d",
			*(stringPool.GetStringPointer(mt.MethClass)),
			*(stringPool.GetStringPointer(mt.MethName)),
			*(stringPool.GetStringPointer(mt.MethType)),
			objRef, paramCount)
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

	// What was returned from the G function?
	switch ret.(type) {
	case *ghelpers.GErrBlk: // G function error block (bad news)
		errBlk := *ret.(*ghelpers.GErrBlk)

		var threadName string
		if f.Thread == 1 {
			threadName = "main"
		} else {
			threadName = fmt.Sprintf("%d", f.Thread)
		}
		var errMsg string
		if f.Thread == 0 {
			errMsg = fmt.Sprintf("in main thread initialization, %s reported by G-function: %s.%s%s",
				errBlk.ErrMsg,
				*(stringPool.GetStringPointer(mt.MethClass)),
				*(stringPool.GetStringPointer(mt.MethName)),
				*(stringPool.GetStringPointer(mt.MethType)))
		} else {
			errMsg = fmt.Sprintf("in thread: %s, in thread %s, reported by G-function: %s.%s%s",
				errBlk.ErrMsg, threadName,
				*(stringPool.GetStringPointer(mt.MethClass)),
				*(stringPool.GetStringPointer(mt.MethName)),
				*(stringPool.GetStringPointer(mt.MethType)))
		}
		status := exceptions.ThrowEx(errBlk.ExceptionType, errMsg, f)
		if status != exceptions.Caught {
			return errors.New(errMsg + " " + errBlk.ErrMsg) // applies only if in test
		}

		// if the exception was caught, tell calling function to execute the catching logic
		return CaughtGfunctionException

	default:
		// No errors.
		// Return int64, float64, *object.Object, or nil.
		// TODO: Come the revolution, this default will be transformed.
		return ret

	}

}

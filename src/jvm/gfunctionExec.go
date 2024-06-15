/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/log"
	"slices"
)

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

func runGfunction(mt classloader.MTentry, fs *list.List,
	className, methodName, methodType string,
	params *[]interface{}, objRef bool) any {

	f := fs.Front().Value.(*frames.Frame)
	var paramCount int
	if params == nil {
		paramCount = 0
	} else {
		paramCount = len(*params)
	}

	fullMethName := fmt.Sprintf("%s.%s%s", className, methodName, methodType)
	if MainThread.Trace {
		traceInfo := fmt.Sprintf("runGfunction: %s, objectRef: %v, paramSlots: %d",
			fullMethName, objRef, paramCount)
		_ = log.Log(traceInfo, log.FINE)
		logTraceStack(f)
	}

	if paramCount > 0 {
		slices.Reverse(*params)
	}

	var ret any
	// call the function, passing it a pointer to the slice of arguments
	if paramCount == 0 {
		ret = mt.Meth.(gfunction.GMeth).GFunction(nil)
	} else {
		ret = mt.Meth.(gfunction.GMeth).GFunction(*params)
	}

	// if an error occured
	switch ret.(type) {
	case *gfunction.GErrBlk:
		// var errorDetails string
		errBlk := *ret.(*gfunction.GErrBlk)
		// parts := strings.SplitN(errBlk.ErrMsg, ":", 2)
		// if len(parts) == 2 {
		// 	funcName = parts[0]
		// 	errorDetails = parts[1]
		// } else {
		// 	funcName = "{MISSINGCOLON}"
		// 	errorDetails = errBlk.ErrMsg
		// }

		var threadName string
		if f.Thread == 1 {
			threadName = "main"
		} else {
			threadName = fmt.Sprintf("%d", f.Thread)
		}
		errMsg := fmt.Sprintf("%s in thread: %s, method: %s",
			errBlk.ErrMsg, threadName, fullMethName)
		status := exceptions.ThrowEx(errBlk.ExceptionType, errMsg, f)
		if status != exceptions.Caught {
			return errors.New(errMsg + " " + errBlk.ErrMsg) // applies only if in test
		} else {
			return nil // return nothing if the exception was caught
		}

	case error:
		errMsg := (ret.(error)).Error()
		status := exceptions.ThrowEx(excNames.NativeMethodException, errMsg, f)
		if status != exceptions.Caught {
			return ret.(error) // applies only if in test
		} else {
			return nil // return nothing if the error was caught
		}
	}

	// if return is not an errBlk or an error, then it's a legitimate
	// return value, so return it.
	return ret
}

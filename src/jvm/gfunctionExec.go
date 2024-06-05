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
	"strings"
)

// Execution of gfunctions (that is, Java functions ported to golang).
// As part of JACOBIN-519, this code seeks to replace the previous set of
// functions (e.g., runGframe() and runGmethod()) with a simpler streamlined
// function.

func runGfunction(
	mt classloader.MTentry,                   // the method we're about to run
	fs *list.List,                            // the frame stack
	className, methodName, methodType string, // self-explanatory
	params *[]interface{},                    // the parameters, including the object reference
	objRef bool) /* flag indicating the object reference is/is not included */ any {

	f := fs.Front().Value.(*frames.Frame)
	var paramCount int
	if params == nil {
		paramCount = 0
	} else {
		paramCount = len(*params)
	}

	if localDebugging || MainThread.Trace {
		traceInfo := fmt.Sprintf("runGfunction: %s.%s%s, objectRef: %v, paramSlots: %d",
			className, methodName, methodType, objRef, paramCount)
		_ = log.Log(traceInfo, log.WARNING)
		logTraceStack(f)
	}

	if paramCount > 0 {
		slices.Reverse(*params)
	}

	var ret any
	// call the function passing a pointer to the slice of arguments
	if paramCount == 0 {
		ret = mt.Meth.(gfunction.GMeth).GFunction(nil)
	} else {
		ret = mt.Meth.(gfunction.GMeth).GFunction(*params)
	}

	// if an error occured
	switch ret.(type) {
	case *gfunction.GErrBlk:
		var funcName, errorDetails string
		errBlk := *ret.(*gfunction.GErrBlk)
		parts := strings.SplitN(errBlk.ErrMsg, ":", 2)
		if len(parts) == 2 {
			funcName = parts[0]
			errorDetails = parts[1]
		} else {
			funcName = "{MISSINGCOLON}"
			errorDetails = errBlk.ErrMsg
		}

		var threadName string
		if f.Thread == 1 {
			threadName = "main"
		} else {
			threadName = fmt.Sprintf("%d", f.Thread)
		}
		errMsg := fmt.Sprintf("com.sun.jdi.NativeMethodException in thread: %s, %s():\n",
			threadName, funcName)
		errMsg = errMsg + errorDetails
		status := exceptions.ThrowEx(errBlk.ExceptionType, errorDetails, f)
		if status != exceptions.Caught {
			return errors.New(errMsg) // applies only if in test
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

	// if it's not an errBlk or an error, then return whatever it is
	return ret
}

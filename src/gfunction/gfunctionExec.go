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
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/trace"
	"slices"
	"sync"
	"time"
)

var CaughtGfunctionException = errors.New("caught gfunction exception")

var thSafeMap sync.Map
var dummy uint8

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
	// then add pointer to the JVM frame stack to the parameter list here.
	entry := mt.Meth.(GMeth)
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

	// Reverse the parameter order. Last appended will be fetched first.
	if paramCount > 0 {
		slices.Reverse(*params)
	}

	// Discern between thread-safe G functions and ordinary ones.
	// No matter what, ret = the result from the G function.
	var ret any
	gmeth := mt.Meth.(GMeth)
	if gmeth.ThreadSafe {
		var loaded = true
		// Make sure that an object reference is the first parameter.
		if !objRef {
			errMsg := "Thread-safe G function requested but no object reference was supplied"
			exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, f)
		}
		// Get key = object pointer.
		key := (*(params))[0].(*object.Object)
		// Lock the key.
	lockloop:
		_, loaded = thSafeMap.LoadOrStore(key, dummy)
		if loaded {
			time.Sleep(globals.SleepMsecs * time.Millisecond) // sleep awhile
			goto lockloop
		}
		// The key is locked to me.
		// Call the G function, passing it a pointer to the slice of arguments.
		if paramCount == 0 {
			ret = gmeth.GFunction(nil)
		} else {
			ret = gmeth.GFunction(*params)
		}
		// Unlock thw key.
		thSafeMap.Delete(key)
	} else {
		// Call the function, passing it a pointer to the slice of arguments.
		if paramCount == 0 {
			ret = gmeth.GFunction(nil)
		} else {
			ret = gmeth.GFunction(*params)
		}
	}

	// if an error occured
	switch ret.(type) {
	case *GErrBlk:
		// var errorDetails string
		errBlk := *ret.(*GErrBlk)

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
			// if the exception was caught, tell calling function to execute the catching logic
			return CaughtGfunctionException
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

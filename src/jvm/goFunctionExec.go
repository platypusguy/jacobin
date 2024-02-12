/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/log"
	"strings"
)

// Similar to global tracing but just for this source file.
var localDebugging bool = false

// This function is called from run(). It executes a frame whose method is
// a native method implemented in golang. It copies the parameters from the
// operand stack and passes them to the golang function, called GFunction,
// as an array of interface{}, which can be nil if there are no arguments.
// Any return value from the method is returned to run() as an interface{}
// (which is nil in the case of a void function), where it is placed
// by run() on the operand stack of the calling function.
func runGframe(fs *list.List, fr *frames.Frame) (interface{}, int, error) {
	if localDebugging || MainThread.Trace {
		traceInfo := fmt.Sprintf("runGframe %s.%s, f.OpStack:", fr.ClName, fr.MethName)
		_ = log.Log(traceInfo, log.WARNING)
		logTraceStack(fr)
	}

	// get the go method from the MTable
	me := classloader.MTable[fr.ClName+"."+fr.MethName]
	if me.Meth == nil {
		return nil, 0, errors.New("runGframe: go method not found: " +
			fr.ClName + "." + fr.MethName)
	}

	// pull arguments for the function off the frame's operand stack and put them in a slice
	var params = new([]interface{})
	for _, v := range fr.OpStack {
		*params = append(*params, v)
	}

	// call the function passing a pointer to the slice of arguments
	ret := me.Meth.(gfunction.GMeth).GFunction(*params)

	// Was a G error block was returned?
	switch ret.(type) {
	case *gfunction.GErrBlk:
		// Get the G error block
		ge := *ret.(*gfunction.GErrBlk)
		// Pop the G frame off the frame stack.
		fs.Remove(fs.Front())
		// Get a pointer to the previous frame.
		fprev := fs.Front().Value.(*frames.Frame)
		// Throw an exception in the previous frame.
		exceptions.ThrowEx(ge.ExceptionType, ge.ErrMsg, fprev)
		// Create an error object to return to caller.
		var err = errors.New(ge.ErrMsg)
		// Return to caller a nil G function result, 0 slots, and an error object.
		return nil, 0, err
	}

	// how many slots does the return value consume on the op stack?
	// the last char in the method name indicates the data type of the return
	// value. If it's 'J' (a long) or 'D' (a double), it will require two
	// slots on the op stack of the calling function. If the return value
	// is nil, then no slots will be required. Otherwise, it's one slot
	// (such as for ints, shorts, boolean, etc.)
	var slotCount int
	if ret == nil {
		slotCount = 0
	} else if strings.HasSuffix(fr.MethName, "J") || strings.HasSuffix(fr.MethName, "D") {
		slotCount = 2
	} else {
		slotCount = 1
	}

	return ret, slotCount, nil
}

// This function creates a new frame for the go-style function, loads its arguments onto
// its stack, pushes the frame onto the head of the frame stack and then calls run() to
// execute it. This eventually calls runGFrame(), which handles any return value. After
// the function is run, this method pops the frame off the frame stack and returns.
// The parameter, objRef, points to the object whose method is being called. It's used
// principally (exclusively?) by INVOKEVIRTUAL and INVOKESPECIAL (See JVM spec).
func runGmethod(mt classloader.MTentry, fs *list.List, className, methodName,
	methodType string, params *[]interface{}, objRef bool) (*frames.Frame, error) {

	f := fs.Front().Value.(*frames.Frame)

	// if the method needs context (i.e., if mt.Meth.NeedsContext == true),
	// then add pointer to the JVM frame stack to the parameter list here.
	entry := mt.Meth.(gfunction.GMeth)
	if entry.NeedsContext {
		*params = append(*params, fs)
	}

	var paramCount int
	if params == nil {
		paramCount = 0
	} else {
		paramCount = len(*params)
	}

	if localDebugging || MainThread.Trace {
		traceInfo := fmt.Sprintf("runGmethod %s.%s%s, objectRef: %v, paramSlots: %d",
			className, methodName, methodType, objRef, paramCount)
		_ = log.Log(traceInfo, log.WARNING)
		logTraceStack(f)
	}

	// create a frame (gf for 'go frame') for this function
	var gf *frames.Frame

	gf = frames.CreateFrame(paramCount)
	gf.Thread = f.Thread
	gf.MethName = methodName + methodType
	gf.ClName = className
	gf.Meth = nil
	gf.CP = nil
	gf.Locals = nil
	gf.Ftype = 'G' // a golang function

	// Current frame stack is one of 2 forms:
	// (1) { pn | ... | p1 | p0 } where TOS is p0                    Note: calls from INVOKESTATIC
	// (2) { pn | ... | p1 | p0 | object Ref }  TOS is object Ref    Note: calls from INVOKEVIRTUAL and INVOKESPECIAL
	//
	// The object ref in #2 is the first argument passed in cases where objRef == true.
	// This object reference points to the object whose method is being called. For example,
	// if String.toUpperCase(), object reference points to the String instance. Note that this
	// item is always at params[0] of the receiving gfunction. Its presence is not counted in
	// the definitions used by LoadLib, as that parameter value is the number of true arguments
	// passed by the original Java method.
	//
	// There exists one exception to the above. If the method has NeedsContext set to true
	// in its definiton, then a pointer to JVM frame stack for the present thread is pushed.
	// It will always appear as the last parameter. There are not many functions in which
	// this is the case.

	// Push the arguments in reverse order onto the Go op stack.
	// If there was an extra parameter, it's at the Go op stack[0].
	for j := paramCount - 1; j >= 0; j-- {
		push(gf, (*params)[j])
	}

	// Set the Go frame TOS to point to first parameter.
	gf.TOS = len(gf.OpStack) - 1
	if localDebugging || MainThread.Trace {
		_ = log.Log("runGmethod G method OpStack:", log.WARNING)
		logTraceStack(gf)
	}

	// push this new frame onto the frame stack for this thread
	fs.PushFront(gf)                     // push the new frame
	f = fs.Front().Value.(*frames.Frame) // point f to the new head

	// then run the frame, which will call run(), which will eventually call runGFrame()
	err := runFrame(fs)

	// If an error object is returned from runFrame,
	// * The G frame has already been popped off,
	//   ensuring that the previous frame is at the head of the frame stack.
	// * ThrowEx has already executed, setting up the next op code to be ATHROW
	//   which will take care of catching and exception reporting.
	// * Return a pointer to the previous frame and a nil error object.
	if err != nil {
		f = fs.Front().Value.(*frames.Frame) // point f the head again
		return f, nil
	}

	// No errors.
	// Pop off the G frame from the frame stack which
	// makes the previous frame the current frame.
	fs.Remove(fs.Front())                // pop off the G frame
	f = fs.Front().Value.(*frames.Frame) // point f to the head (previous frame)
	return f, nil
}

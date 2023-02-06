/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"container/list"
	"errors"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/log"
	"strings"
)

// This function is called from run(). It executes a frame whose
// method is a golang method. It copies the parameters from the
// operand stack and passes them to the go function, here called Fu,
// as an array of interface{}, which can be nil if there are no arguments.
// Any return value from the method is returned to run() as an interface{}
// (which is nil in the case of a void function), where it is placed
// by run() on the operand stack of the calling function.
func runGframe(fr *frames.Frame) (interface{}, int, error) {
	// get the go method from the MTable
	me := classloader.MTable[fr.MethName]
	if me.Meth == nil {
		return nil, 0, errors.New("go method not found: " + fr.MethName)
	}

	// pull arguments for the function off the frame's operand stack and put them in a slice
	var params = new([]interface{})
	for _, v := range fr.OpStack {
		*params = append(*params, v)
	}

	// call the function passing a pointer to the slice of arguments
	ret := me.Meth.(classloader.GmEntry).Fu(*params)

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
func runGmethod(mt classloader.MTentry, fs *list.List, className, methodName, methodType string) (*frames.Frame, error) {
	f := fs.Front().Value.(*frames.Frame)

	// create a frame (gf for 'go frame') for this function
	paramSlots := mt.Meth.(classloader.GmEntry).ParamSlots
	gf := frames.CreateFrame(paramSlots)
	gf.Thread = f.Thread

	gf.MethName = methodName + methodType
	gf.ClName = className
	gf.Meth = nil
	gf.CP = nil
	gf.Locals = nil
	gf.Ftype = 'G' // a golang function

	// get the args (if any) from the operand stack of the current frame(f)
	// then push them onto the stack of the go function
	var argList []interface{}

	for i := 0; i < paramSlots; i++ {
		arg := pop(f)
		intArg := arg
		argList = append(argList, intArg)
	}
	for j := len(argList) - 1; j >= 0; j-- {
		push(gf, argList[j])
	}
	gf.TOS = len(gf.OpStack) - 1

	// push this new frame onto the frame stack for this thread
	fs.PushFront(gf)                     // push the new frame
	f = fs.Front().Value.(*frames.Frame) // point f to the new head

	// then run the frame, which will call run(), which will eventually call runGFrame()
	err := runFrame(fs)
	if err != nil {
		_ = log.Log("Error: "+err.Error(), log.SEVERE)
		return nil, err
	}

	// now that the go function is done, pop the frame off the stack and
	// point the previous frame as the current frame
	fs.Remove(fs.Front())                // pop the frame off
	f = fs.Front().Value.(*frames.Frame) // point f the head again
	return f, nil
}

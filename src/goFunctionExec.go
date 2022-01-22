/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"container/list"
	"errors"
	"jacobin/classloader"
)

// This function is called from main.run(). It execuates a frame whose
// method is a golang method. It copies the parameters from the
// operand stack and passes them to the go function, here called Fu,
// as an array of interface{}, which can be nil if there are no arguments.
// Any return value from the method is returned to run() as an interface{}
// (which is nil in the case of a void function), where it is placed
// by run() on the operand stack of the calling function.
func runGframe(fr *frame) (interface{}, error) {
	// get the go method from the MTable
	me := classloader.MTable[fr.methName]
	if me.Meth == nil {
		return nil, errors.New("go method not found: " + fr.methName)
	}

	// pull arguments for the function off the frame's operand stack and put them in a slice
	var params = new([]interface{})
	for _, v := range fr.opStack {
		*params = append(*params, v)
	}

	// call the function passing a pointer to the slice of arguments
	ret := me.Meth.(classloader.GmEntry).Fu(*params)
	return ret, nil
}

// This function creates a new frame for the go-style function, loads its arguments onto
// its stack, pushes the frame onto the head of the frame stack and then calls run() to
// execute it. This eventually calls runGFrame(), which handles any return value. After
// the function is run, this method pops the frame off the frame stack and returns.
func runGmethod(mt classloader.MTentry, fs *list.List, className, methodName, methodType string) error {
	f := fs.Front().Value.(*frame)

	// create a frame (gf for 'go frame') for this function
	paramSlots := mt.Meth.(classloader.GmEntry).ParamSlots
	gf := createFrame(paramSlots)
	gf.thread = f.thread
	gf.methName = className + "." + methodName + methodType
	gf.clName = className
	gf.meth = nil
	gf.cp = nil
	gf.locals = nil
	gf.ftype = 'G' // a golang function

	// get the args (if any) from the operand stack of the current frame(f)
	// then push them onto the stack of the go function
	var argList []int64
	for i := 0; i < paramSlots; i++ {
		arg := pop(f)
		argList = append(argList, arg)
	}
	for j := len(argList) - 1; j >= 0; j-- {
		push(gf, argList[j])
	}
	gf.tos = len(gf.opStack) - 1

	// push this new frame onto the frame stack for this thread
	fs.PushFront(gf)              // push the new frame
	f = fs.Front().Value.(*frame) // point f to the new head

	// then run the frame, which will call run(), which will eventually call runGFrame()
	err := runFrame(fs)
	if err != nil {
		return err
	}

	// now that the go function is done, pop the frame off the stack and
	// point the previous frame as the current frame
	fs.Remove(fs.Front())         // pop the frame off
	f = fs.Front().Value.(*frame) // point f the head again
	return nil
}

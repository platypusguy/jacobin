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

// runs a frame whose method is a golang method. It copies the parameters
// from the operand stack and passes them to the go function, here called Fu.
// Any return value from the method is returned to the call from run(), where
// it is placed on the stack of the calling function.
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

func runGmethod(mt classloader.MTentry, fs *list.List, className, methodName, methodType string) error {
	f := fs.Front().Value.(*frame)

	paramSlots := mt.Meth.(classloader.GmEntry).ParamSlots
	gf := createFrame(paramSlots)
	gf.thread = f.thread
	gf.methName = className + "." + methodName + methodType
	gf.clName = className
	gf.meth = nil
	gf.cp = nil
	gf.locals = nil
	gf.ftype = 'G' // a golang function

	var argList []int64
	for i := 0; i < paramSlots; i++ {
		arg := pop(f)
		argList = append(argList, arg)
	}
	for j := len(argList) - 1; j >= 0; j-- {
		push(gf, argList[j])
	}
	gf.tos = len(gf.opStack) - 1

	fs.PushFront(gf)              // push the new frame
	f = fs.Front().Value.(*frame) // point f to the new head

	err := runFrame(fs) // this will eventually find its way to runGFrame()
	if err != nil {
		return err
	}

	fs.Remove(fs.Front())         // pop the frame off
	f = fs.Front().Value.(*frame) // point f the head again
	return nil
}

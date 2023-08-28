/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/thread"
	"strconv"
)

// Initialization blocks are code blocks that for all intents are methods. They're gathered up by the
// Java compiler into a method called <clinit>, which must be run at class instantiation--that is,
// before any constructor. Because that code might well call other methods, it will need to be run
// just like a regular method with stack frames and depending on the interpreter in run.go
// In addition, we have to make sure that the initialization blocks of superclasses have been
// executed as well.
//
// CURR: Implement the above logic here.
func runInitializationBlock(k *classloader.Klass, idx int) error {

	// msg := fmt.Sprintf("<clinit> found in %s, method #%d\n", k.Data.Name, idx)
	// _ = log.Log(msg, log.FINE)
	// fmt.Print(msg)

	className := k.Data.Name
	me, err := classloader.FetchMethodAndCP(className, "<clinit>", "()V")
	if err != nil {
		// in case of error, user will be notified in calling function
		return errors.New("Method not found: " + className + "<clinit>()")
	}

	switch me.MType {
	case 'J': // it's a Java initializer (the most common case)
		return runJavaInitializer(me.Meth, k)
	case 'G': // it's a native (that is, golang) initializer
		return runNativeInitializer(me.Meth, k)
	}
	// m := me.Meth.(classloader.JmEntry)
	// f := frames.CreateFrame(m.MaxStack) // create a new frame
	// f.MethName = "<clinit>"
	// f.ClName = className
	// f.CP = m.Cp                        // add its pointer to the class CP
	// for i := 0; i < len(m.Code); i++ { // copy the bytecodes over
	// 	f.Meth = append(f.Meth, m.Code[i])
	// }
	//
	// // allocate the local variables
	// for j := 0; j < m.MaxLocals; j++ {
	// 	f.Locals = append(f.Locals, 0)
	// }
	//
	// // create the first thread and place its first frame on it
	// glob := globals.GetGlobalRef()
	// clInitThread := thread.CreateThread()
	// clInitThread.Stack = frames.CreateFrameStack()
	// clInitThread.ID = thread.AddThreadToTable(&clInitThread, &glob.Threads)
	//
	// clInitThread.Trace = MainThread.Trace
	// f.Thread = clInitThread.ID
	//
	// if frames.PushFrame(clInitThread.Stack, f) != nil {
	// 	_ = log.Log("Memory exceptions allocating frame on thread: "+strconv.Itoa(clInitThread.ID),
	// 		log.SEVERE)
	// 	return errors.New("outOfMemory Exception")
	// }
	//
	// if clInitThread.Trace {
	// 	traceInfo := fmt.Sprintf("StartExec: f.MethName=%s, m.MaxStack=%d, m.MaxLocals=%d, len(m.Code)=%d",
	// 		f.MethName, m.MaxStack, m.MaxLocals, len(m.Code))
	// 	_ = log.Log(traceInfo, log.TRACE_INST)
	// }
	//
	// err = runThread(&clInitThread)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func runJavaInitializer(m classloader.MData, k *classloader.Klass) error {
	meth := m.(classloader.JmEntry)
	f := frames.CreateFrame(meth.MaxStack) // create a new frame
	f.MethName = "<clinit>"
	f.ClName = k.Data.Name
	f.CP = meth.Cp                        // add its pointer to the class CP
	for i := 0; i < len(meth.Code); i++ { // copy the bytecodes over
		f.Meth = append(f.Meth, meth.Code[i])
	}

	// allocate the local variables
	for j := 0; j < meth.MaxLocals; j++ {
		f.Locals = append(f.Locals, 0)
	}

	// create the first thread and place its first frame on it
	glob := globals.GetGlobalRef()
	clInitThread := thread.CreateThread()
	clInitThread.Stack = frames.CreateFrameStack()
	clInitThread.ID = thread.AddThreadToTable(&clInitThread, &glob.Threads)

	clInitThread.Trace = MainThread.Trace
	f.Thread = clInitThread.ID

	if frames.PushFrame(clInitThread.Stack, f) != nil {
		_ = log.Log("Memory exceptions allocating frame on thread: "+strconv.Itoa(clInitThread.ID),
			log.SEVERE)
		return errors.New("outOfMemory Exception")
	}

	if clInitThread.Trace {
		traceInfo := fmt.Sprintf("StartExec: f.MethName=%s, m.MaxStack=%d, m.MaxLocals=%d, len(m.Code)=%d",
			f.MethName, meth.MaxStack, meth.MaxLocals, len(meth.Code))
		_ = log.Log(traceInfo, log.TRACE_INST)
	}

	err := runThread(&clInitThread)
	if err != nil {
		return err
	}
	return nil
}

func runNativeInitializer(meth classloader.MData, k *classloader.Klass) error {
	return nil
}

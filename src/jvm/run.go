/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
	"runtime"
	"runtime/debug"
)

// Run as an independent thread to completion (TERMINATED).
// Launched with: go globals.GetGlobalRef().FuncRunThread(t, clName, methName, methType)
func RunJavaThread(args []any) {

	if len(args) != 4 {
		errMsg := fmt.Sprintf("RunJavaThread: Expected 4 arguments, observed %d: %v", len(args), args)
		exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, nil)
	}

	// Set up arguments.
	t := args[0].(*object.Object)
	clName := args[1].(string)
	methName := args[2].(string)
	methType := args[3].(string)

	defer func() int {
		// only an untrapped panic gets us here
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack
			exceptions.ShowPanicCause(r)
			exceptions.ShowFrameStack(t) // update for new thread model
			exceptions.ShowGoStackTrace(nil)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	// Set up the method and CPool for the method's class.
	mte, err := classloader.FetchMethodAndCP(clName, methName, methType)
	if err != nil {
		errMsg := fmt.Sprintf("RunJavaThread: Could not find run method (%s.%s%s): %v", clName, methName, methType, err)
		exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, nil)
	}

	//trace.Trace(fmt.Sprintf("DEBUG RunJavaThread: Found run method: %s.%s%s", clName, methName, methType))

	// Get Mtable entry for the method.
	meth := mte.Meth.(classloader.JmEntry)

	// Create the frame stack for this thread.
	fs := frames.CreateFrameStack()

	// Create the initial frame for this thread.
	f := frames.CreateFrame(meth.MaxStack + types.StackInflator) // experiment with stack size. See JACOBIN-494
	t.ThMutex.RLock()
	tID := t.FieldTable["ID"].Fvalue.(int64)
	t.ThMutex.RUnlock()
	f.Thread = int(tID)
	f.ClName = clName
	f.MethName = methName
	f.MethType = methType
	f.AccessFlags = meth.AccessFlags

	f.CP = meth.Cp                        // add its pointer to the class CP
	f.Meth = append(f.Meth, meth.Code...) // copy the bytecodes over

	// Allocate the method's local variables for this frame.
	for k := 0; k < meth.MaxLocals; k++ {
		f.Locals = append(f.Locals, int64(0))
	}

	// If this is the program entry point `main(String[] args)`, initialize local 0
	// with a proper Java String[] built from the CLI application arguments.
	// Without this, bytecodes like ARRAYLENGTH on `args` would see an uninitialized
	// int (zero) and throw an IllegalArgumentException.
	if methName == "main" && methType == "([Ljava/lang/String;)V" && len(f.Locals) > 0 {
		appArgs := globals.GetGlobalRef().AppArgs
		// Build a Java String object array from Go strings
		strObjs := object.StringObjectArrayFromGoStringArray(appArgs)
		// Wrap it into a Java reference array object of type java/lang/String
		strArrayObj := object.Make1DimRefArray("java/lang/String", int64(len(strObjs)))
		// Populate the array elements
		if val, ok := strArrayObj.FieldTable["value"]; ok {
			if arr, ok2 := val.Fvalue.([]*object.Object); ok2 && len(arr) == len(strObjs) {
				copy(arr, strObjs)
				// store back (not strictly necessary since slice is by reference)
				strArrayObj.FieldTable["value"] = object.Field{Ftype: val.Ftype, Fvalue: arr}
			}
		}
		// Assign to local variable 0
		f.Locals[0] = strArrayObj
	} else {
		// Not the main thread.
		// Try for a runnable object in the thread's field table.
		t.ThMutex.RLock()
		fld, ok := t.FieldTable["target"]
		t.ThMutex.RUnlock()
		if ok {
			// Got a target field (runnable object).
			// Initialize local 0 with the runnable object so that it can be accessed within the interpreter
			f.Locals[0] = fld.Fvalue
		} else {
			// Not target field present.
			// Initialize local 0 with the thread object so that it can be accessed within the interpreter
			f.Locals[0] = t
		}
	}

	// JACOBIN-824:
	// cl := classloader.MethAreaFetch(clName) // JACOBIN-824
	// // if cl == nil {
	// // 	errMsg := fmt.Sprintf("RunJavaThread: Could not load class %s", clName)
	// // 	exceptions.ThrowEx(excNames.ClassNotFoundException, errMsg, nil)
	// // 	return
	// // }
	// if cl != nil {
	// 	f.Locals[0] = cl.Data
	// }

	// Add the initial frame and the frame stack to the thread's field table.
	t.ThMutex.Lock()
	t.FieldTable["frame"] = object.Field{Ftype: types.Ref, Fvalue: f}
	t.FieldTable["framestack"] = object.Field{Ftype: types.LinkedList, Fvalue: fs}
	t.ThMutex.Unlock()

	// Push the frame on the stack.
	if frames.PushFrame(fs, f) != nil {
		errMsg := fmt.Sprintf("RunJavaThread: frames.PushFrame failed on thread: %d", tID)
		exceptions.ThrowEx(excNames.OutOfMemoryError, errMsg, nil)
	}

	// Mark the thread RUNNABLE and register it.
	_, ret := javaLang.SetThreadState(t, javaLang.RUNNABLE)
	if ret != nil {
		errMsg := "RunJavaThread: SetThreadState(RUNNABLE) failed"
		exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, nil)
	}

	// Register this non-main thread.
	javaLang.RegisterThread(t)

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("threadRun: class=%s, meth=%s%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, f.MethType, meth.MaxStack, meth.MaxLocals, len(meth.Code))
		trace.Trace(traceInfo)
	}

	// Execute the thread's frame set.
	for fs.Len() > 0 {
		interpret(fs)
		runtime.Gosched()
	}

	// The End.
	javaLang.SetThreadState(t, javaLang.TERMINATED)

	// Notify all threads waiting for this thread to terminate (e.g., in Thread.join()).
	// Per JVM spec, this is equivalent to t.notifyAll(), which requires holding the lock on t.
	if err := t.ObjLock(int32(tID)); err == nil {
		_ = t.ObjectNotifyAll(int32(tID))
		_ = t.ObjUnlock(int32(tID))
	}
}

// multiply two numbers
func multiply[N frames.Number](num1, num2 N) N {
	return num1 * num2
}

func subtract[N frames.Number](num1, num2 N) N {
	return num1 - num2
}

// create a new frame and load up the local variables with the passed
// arguments, set up the stack, and all the remaining items to begin execution
// Note: the includeObjectRef parameter is a boolean. When true, it indicates
// that in addition to the method parameter, an object reference is also on
// the stack and needs to be popped off the caller's opStack and passed in.
// (This would be the case for invokevirtual, among others.) When false, no
// object pointer is needed (for invokestatic, among others).
func createAndInitNewFrame(
	className string, methodName string, methodType string,
	m *classloader.JmEntry,
	includeObjectRef bool,
	currFrame *frames.Frame) (*frames.Frame, error) {

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("createAndInitNewFrame: class=%s, meth=%s%s, includeObjectRef=%v, maxStack=%d, maxLocals=%d",
			className, methodName, methodType, includeObjectRef, m.MaxStack, m.MaxLocals)
		trace.Trace(traceInfo)
	}

	f := currFrame

	stackSize := m.MaxStack + types.StackInflator // Experimental addition, see JACOBIN-494
	if stackSize < 1 {
		stackSize = 2
	}

	fram := frames.CreateFrame(stackSize)
	fram.Thread = currFrame.Thread
	fram.FrameStack = currFrame.FrameStack
	fram.ClName = className
	fram.MethName = methodName
	fram.MethType = methodType
	fram.CP = m.Cp                           // add its pointer to the class CP
	fram.Meth = append(fram.Meth, m.Code...) // copy the method's bytecodes over
	fram.AccessFlags = m.AccessFlags

	// pop the parameters off the present stack and put them in
	// the new frame's locals. This is done in reverse order so
	// that the parameters are pushed in the right order to be
	// popped off by the receiving function
	var argList []interface{}
	paramsToPass :=
		util.ParseIncomingParamsFromMethTypeString(methodType)

	// primitives use a single byte/letter, but arrays can be many bytes:
	// a minimum of two (e.g., [I for array of ints). If the array
	// is multidimensional, the bytes will be [[I with one instance
	// of [ for every dimension. In the case of multidimensional
	// arrays, the arrays are always pushed as arrays of references,
	// and we simply mark off the number of [. For single-dimensional
	// arrays, we pass the kind of pointer that applies and mark off
	// a single instance of [
	for j := len(paramsToPass) - 1; j > -1; j-- {
		param := paramsToPass[j]
		primitive := param[0]

		arrayDimensions := 0
		if primitive == '[' {
			i := 0
			for i = 0; i < len(param); i++ {
				if param[i] == '[' {
					arrayDimensions += 1
				} else {
					break
				}
			}
			// param[i] now holds the primitive of the array
			primitive = param[i]
		}

		if arrayDimensions > 1 { // a multidimensional array
			// if the array is multidimensional, then we are
			// passing in a pointer to an array of references
			// to objects (lower arrays) regardless of the
			// lowest level of primitive in the array
			arg := pop(f).(*object.Object)
			argList = append(argList, arg)
			continue
		}

		if arrayDimensions == 1 { // a single-dimension array
			// a bunch of Java functions return raw arrays (like String.toCharArray()), which
			// are not really viewed by the JVM as objects in the full sense of the term. These
			// almost invariably are single-dimension arrays. So we test for these here and
			// return the corresponding object entity.
			value := pop(f)
			arg := object.MakeArrayFromRawArray(value)
			argList = append(argList, arg)
			continue
		}

		switch primitive { // it's not an array
		case 'D': // double
			arg := pop(f).(float64)
			argList = append(argList, arg)
		case 'F': // float
			arg := pop(f).(float64)
			argList = append(argList, arg)
		case 'B', 'C', 'I', 'S': // byte, char, integer, short
			arg := pop(f)
			switch arg.(type) {
			case int: // the arg should be int64, but is occasionally int. Tracking this down.
				arg = int64(arg.(int))
			}
			argList = append(argList, arg)
		case 'J': // long
			arg := pop(f).(int64)
			argList = append(argList, arg)
		case 'L': // pointer/reference
			arg := pop(f) // can't be *Object b/c the arg could be nil, which would panic
			argList = append(argList, arg)
		default:
			arg := pop(f)
			argList = append(argList, arg)
		}
	}

	// Initialize lenLocals = max (m.MaxLocals, len(argList)) but at least 1
	lenArgList := len(argList)
	lenLocals := m.MaxLocals
	if lenArgList > m.MaxLocals {
		lenLocals = lenArgList
	}
	if lenLocals < 1 {
		lenLocals = 1
	}

	// allocate the local variables
	for k := 0; k < lenLocals; k++ {
		fram.Locals = append(fram.Locals, int64(0))
	}

	// if includeObjectRef is true then objectRef != nil.
	// Insert it in the local[0]
	// This is used in invokevirtual, invokespecial, and invokeinterface.
	destLocal := 0
	if includeObjectRef {
		fram.Locals[0] = pop(f)
		fram.Locals = append(fram.Locals, int64(0)) // add the slot taken up by objectRef
		destLocal = 1                               // The first parameter starts at index 1
		lenLocals++                                 // There is 1 more local needed
		if fram.AccessFlags&classloader.ACC_SYNCHRONIZED > 0 {
			obj := fram.Locals[0].(*object.Object)
			fram.ObjSync = obj
			err := obj.ObjLock(int32(fram.Thread))
			if err != nil {
				fqn := fram.ClName + "." + fram.MethName + fram.MethType
				errMsg := fmt.Sprintf("createAndInitNewFrame: ObjLock error, PC: %d, FQN: %s", fram.PC, fqn)
				return nil, errors.New(errMsg)
			}
			if globals.TraceInst {
				traceInfo := fmt.Sprintf("\tcreateAndInitNewFrame: Locked object %s",
					object.GoStringFromStringPoolIndex(obj.KlassName))
				trace.Trace(traceInfo)
			}
		}
	}

	if globals.TraceVerbose {
		traceInfo := fmt.Sprintf("\tcreateAndInitNewFrame: lenArgList=%d, lenLocals=%d, stackSize=%d",
			lenArgList, lenLocals, stackSize)
		trace.Trace(traceInfo)
	}

	ptpx := 0
	for j := lenArgList - 1; j >= 0; j-- {
		fram.Locals[destLocal] = argList[j]
		switch paramsToPass[ptpx] {
		case "D", "J":
			destLocal += 2
		default:
			destLocal += 1
		}
		ptpx++
	}

	fram.TOS = -1

	return fram, nil
}

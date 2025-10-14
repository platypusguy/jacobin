/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/config"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/statics"
	"jacobin/src/thread"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
	"os"
	"runtime/debug"
	"strconv"
)

var MainThread thread.ExecThread

// StartExec is where execution begins. It initializes various structures, such as
// the MTable, then using the passed-in name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string, mainThread *thread.ExecThread, globalStruct *globals.Globals) {

	MainThread = *mainThread

	me, err := classloader.FetchMethodAndCP(className, "main", "([Ljava/lang/String;)V")
	if err != nil {
		errMsg := "Class not found: " + className + ".main()"
		exceptions.ThrowEx(excNames.ClassNotFoundException, errMsg, nil)
	}

	m := me.Meth.(classloader.JmEntry)
	f := frames.CreateFrame(m.MaxStack + types.StackInflator) // experiment with stack size. See JACOBIN-494
	f.Thread = MainThread.ID
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"
	f.ClName = className
	f.CP = m.Cp                        // add its pointer to the class CP
	f.Meth = append(f.Meth, m.Code...) // copy the bytecodes over

	// allocate the local variables
	for k := 0; k < m.MaxLocals; k++ {
		f.Locals = append(f.Locals, 0)
	}

	// Create an array of string objects for any CLI args in locals[0].
	var objArray []*object.Object
	for _, str := range globalStruct.AppArgs {
		// sobj := object.NewStringFromGoString(str) // deprecated by JACOBIN-480
		sobj := object.StringObjectFromGoString(str)
		objArray = append(objArray, sobj)
	}
	f.Locals[0] = object.MakePrimitiveObject("[Ljava/lang/String", types.RefArray, objArray)

	// create the first thread and place its first frame on it
	MainThread.Stack = frames.CreateFrameStack()
	mainThread.Stack = MainThread.Stack

	// moved here as part of JACOBIN-554. Was previously after the InstantiateClass() call next
	if frames.PushFrame(MainThread.Stack, f) != nil {
		errMsg := "Memory error allocating frame on thread: " + strconv.Itoa(MainThread.ID)
		exceptions.ThrowEx(excNames.OutOfMemoryError, errMsg, nil)
	}

	// must first instantiate the class, so that any static initializers are run
	_, instantiateError := InstantiateClass(className, MainThread.Stack)
	if instantiateError != nil {
		errMsg := "Error instantiating: " + className + ".main()"
		exceptions.ThrowEx(excNames.InstantiationException, errMsg, nil)
	}

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("StartExec: class=%s, meth=%s%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, f.MethType, m.MaxStack, m.MaxLocals, len(m.Code))
		trace.Trace(traceInfo)
	}

	err = runThread(&MainThread)

	if globals.TraceVerbose {
		statics.DumpStatics("StartExec end", statics.SelectUser, "")
		_ = config.DumpConfig(os.Stderr)
	}
}

// Point the thread to the top of the frame stack and tell it to run from there.
func runThread(t *thread.ExecThread) error {

	defer func() int {
		// only an untrapped panic gets us here
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack
			exceptions.ShowPanicCause(r)
			exceptions.ShowFrameStack(t)
			exceptions.ShowGoStackTrace(nil)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	for t.Stack.Len() > 0 {
		interpret(t.Stack)
	}

	if t.Stack.Len() == 0 { // true when the last executed frame was main()
		return nil
	}
	return nil
}

func runJavaThread(thObj *object.Object) error {
	t := thObj.FieldTable

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

	fs := t["framestack"].Fvalue.(*list.List)
	for fs.Len() > 0 {
		interpret(fs)
	}

	if fs.Len() == 0 { // true when the last executed frame was main()
		return nil
	}
	return nil
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

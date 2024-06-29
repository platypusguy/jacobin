/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/shutdown"
	"jacobin/stringPool"
	"jacobin/thread"
	"jacobin/util"
	"os"
	"runtime/debug"
)

// This file contains the functions for throwing exceptions from within
// Jacobin. That is, situations in which Jacobin itself is throwing the error,
// rather than the application. Typically, this is for errors/exceptions in the
// operation of the JVM, and for a few occasional user errors, such as
// divide by zero.

const (
	Caught    = true
	NotCaught = false
)

// ThrowExNil simply calls ThrowEx with a nil pointer for the frame.
func ThrowExNil(which int, msg string) {
	ThrowEx(which, msg, nil)
}

// ThrowEx throws an exception. It is used primarily for exceptions and
// errors thrown by Jacobin, rather than by the application. (The latter
// would generally use the ATHROW bytecode.)
//
// Important: if you change the name of this function, you need to update
// exceptions.ShowGoStackTrace(), which explicitly tests for this function name.
func ThrowEx(which int, msg string, f *frames.Frame) bool {
	traceMsg := fmt.Sprintf("[ThrowEx] %s, msg: %s", excNames.JVMexceptionNames[which], msg)
	_ = log.Log(traceMsg, log.TRACE_INST)

	// If in a unit test, log a severe message and return.
	glob := globals.GetGlobalRef()
	if glob.JacobinName == "test" {
		errMsg := fmt.Sprintf("%s in %s.%s, %s",
			excNames.JVMexceptionNames[which],
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName, msg)
		fmt.Fprintln(os.Stderr, errMsg)
		return NotCaught
	}

	// Frame pointer provided?
	if f == nil {
		minimalAbort(which, msg) // this calls exit()
	}

	// the name of the exception as shown to the user
	exceptionNameForUser := excNames.JVMexceptionNames[which]

	// the internal format used in the constant pool
	exceptionCPname := util.ConvertClassFilenameToInternalFormat(exceptionNameForUser)

	// capture the PC where the exception was thrown, if it hasn't been captured yet.
	// (saved b/c later we modify the value of f.PC)
	if f.ExceptionPC == -1 {
		f.ExceptionPC = f.PC
	}

	th, ok := glob.Threads[f.Thread].(*thread.ExecThread)
	if !ok {
		errMsg := fmt.Sprintf("[ThrowEx] glob.Threads index not found or entry corrupted, thread index: %d", f.Thread)
		minimalAbort(excNames.InternalException, errMsg)
	}
	fs := th.Stack

	// find out if the exception is caught and if so point to the catch code
	// catchFrame, catchPC := FindExceptionFrame(f, exceptionCPname, f.ExceptionPC)
	catchFrame, catchPC := FindCatchFrame(fs, exceptionCPname, f.ExceptionPC)
	if catchFrame != nil {
		// at this point, we know that the exception was caught
		// and that the returned frame is the frame
		// containing the catch logic, referred to here as the catchFrame.
		// now, set up the execution of the catch code by:
		// 0. popping off the frames that are above the catch frame,
		//    if any--so that top frame in the frame stack is the catch frame
		// 1. creating a new objRef for the exception
		// 2. pushing the objRef on the op stack of the frame
		// 3. setting the PC to point to the catch code (which expects the objRef at TOS)
		caughtMsg := fmt.Sprintf("[ThrowEx] caught %s, msg: %s", exceptionCPname, msg)
		log.Log(caughtMsg, log.TRACE_INST)

		th = glob.Threads[f.Thread].(*thread.ExecThread)
		fs = th.Stack
		for fs.Len() > 0 { // remove the frames we examined that did not have the catch logic
			fr := fs.Front().Value
			if fr == catchFrame {
				break
			} else {
				fs.Remove(fs.Front())
			}
		}

		objRef, _ := glob.FuncInstantiateClass(exceptionCPname, fs)
		catchFrame.TOS = 0
		catchFrame.OpStack[0] = objRef // push the objRef
		// catchFrame.PC = catchPC - 1    // -1 because the loop in run.go will increment PC after this code block's return
		catchFrame.PC = catchPC

		// the exception logic might throw another exception, in which case that will be
		// the new ExceptionPC. However, it won't be updated to that value unless ExceptionPC
		// is reset to -1. So, at this point, the exception's been caught, so we can reset
		// ExeptionPC to -1. See JACOBIN-534
		f.ExceptionPC = -1
		return Caught
	}

	// ---- if exception is not caught ----

	throwObject, err := glob.FuncInstantiateClass(exceptionCPname, fs)
	if err != nil {
		println(err.Error())
		if throwObject != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", throwObject)
			_ = shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	}

	throwObj := throwObject.(*object.Object)
	params := []any{fs, throwObj}
	glob.FuncFillInStackTrace(params)

	excInfo := fmt.Sprintf("%s: %s", exceptionNameForUser, msg)
	fmt.Fprintln(os.Stderr, excInfo)

	stackTrace := throwObj.FieldTable["stackTrace"].Fvalue.(*object.Object)
	traceEntries := stackTrace.FieldTable["value"].Fvalue.([]*object.Object)

	// now print out the JVM stack
	for _, traceEntry := range traceEntries {
		// HotSpot uses a slightly different format for method names:
		// package.class.method, we prefer package/class.method, so we format
		// method name according to whether -strictJDK is in force
		var declaringClass string
		if glob.StrictJDK {
			declaringClass = util.ConvertInternalClassNameToUserFormat(
				traceEntry.FieldTable["declaringClass"].Fvalue.(string))
		} else {
			declaringClass = traceEntry.FieldTable["declaringClass"].Fvalue.(string)
		}

		traceInfo := fmt.Sprintf("  at %s.%s(%s:%s)",
			declaringClass,
			traceEntry.FieldTable["methodName"].Fvalue.(string),
			traceEntry.FieldTable["fileName"].Fvalue.(string),
			traceEntry.FieldTable["sourceLine"].Fvalue.(string))
		fmt.Fprintln(os.Stderr, traceInfo)
	}

	if !glob.StrictJDK {
		// the next statement disables showing the line that identifies
		// the cause of a golang panic, because if we got here, there
		// was no panic, rather just an uncaught exception. So we show
		// the golang stack without implying there was a panic.
		glob.PanicCauseShown = true
		ShowGoStackTrace("")
	}

	_ = shutdown.Exit(shutdown.JVM_EXCEPTION) // in test mode, this call returns
	return NotCaught
}

func generateThrowBytecodes(f *frames.Frame, exceptionCPname string, msg string) []byte {
	// the functionality we generate bytecodes for is (using a NPE as an example):
	// 0: new           #7                  // class java/lang/NullPointerException
	// 3: dup
	// 4: ldc           #9                  // String  (the msg passed into this function)
	// 6: invokespecial #11                 // Method java/lang/NullPointerException."<init>":(Ljava/lang/String;)V
	// 9: athrow
	//
	// Note that to do this, we need to twiddle with the constant pool as well

	CP := f.CP.(*classloader.CPool)
	CP.Utf8Refs = append(CP.Utf8Refs, exceptionCPname) // probably not needed due to use of string pool
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.UTF8, Slot: uint16(len(CP.Utf8Refs) - 1)})
	// then add a classref entry for the exception
	nameIndex := stringPool.GetStringIndex(&exceptionCPname)
	CP.ClassRefs = append(CP.ClassRefs, nameIndex) // point to the string pool entry
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.ClassRef, Slot: uint16(len(CP.ClassRefs) - 1)})
	exceptionClassCPindex := uint16(len(CP.CpIndex) - 1)

	// start converting previous work into bytecodes
	var genCode []byte
	genCode = append(genCode, opcodes.NOP) // the first bytecode is skipped by the JVM
	genCode = append(genCode, opcodes.NEW)
	hiByte := uint8((len(CP.CpIndex) - 1) >> 8)
	loByte := uint8(len(CP.CpIndex) - 1)
	genCode = append(genCode, hiByte)
	genCode = append(genCode, loByte)
	genCode = append(genCode, opcodes.DUP)

	// now load the error message, if any
	if msg != "" {
		CP.Utf8Refs = append(CP.Utf8Refs, msg)
		utf8MsgIndex := uint16(len(CP.Utf8Refs) - 1)
		CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
			Type: classloader.UTF8, Slot: utf8MsgIndex})
		CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
			Type: classloader.StringConst, Slot: uint16(len(CP.CpIndex) - 1)})
		stringMsgIndex := uint16(len(CP.CpIndex) - 1)
		if stringMsgIndex < 256 {
			genCode = append(genCode, opcodes.LDC)
			genCode = append(genCode, uint8(stringMsgIndex))
		} else {
			// if the index is > 255, we need to use LDC_W and a two-byte index
			hiByte = uint8(stringMsgIndex >> 8)
			loByte = uint8(stringMsgIndex)
			genCode = append(genCode, opcodes.LDC_W)
			genCode = append(genCode, hiByte)
			genCode = append(genCode, loByte)
		}
	}

	// now, set up the CP entries for INVOKESPECIAL. This includes a MethodRef
	// which points to the previous classRef and to a name and type record, which
	// itself points to a UTF8 entry for the method name and a UTF8 entry for
	// the method's signature. We start with the NameAndType entry.
	CP.Utf8Refs = append(CP.Utf8Refs, "<init>")
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.UTF8, Slot: uint16(len(CP.Utf8Refs) - 1)})
	if msg != "" {
		CP.Utf8Refs = append(CP.Utf8Refs, "(Ljava/lang/String;)V")
	} else {
		CP.Utf8Refs = append(CP.Utf8Refs, "()V")
	}
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.UTF8, Slot: uint16(len(CP.Utf8Refs) - 1)})
	CP.NameAndTypes = append(CP.NameAndTypes, classloader.NameAndTypeEntry{
		NameIndex: uint16(len(CP.CpIndex) - 2),
		DescIndex: uint16(len(CP.CpIndex) - 1)})
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.NameAndType, Slot: uint16(len(CP.NameAndTypes) - 1)})

	// now the MethodRef entry
	CP.MethodRefs = append(CP.MethodRefs, classloader.MethodRefEntry{
		ClassIndex:  exceptionClassCPindex,
		NameAndType: uint16(len(CP.CpIndex) - 1)})
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.MethodRef, Slot: uint16(len(CP.MethodRefs) - 1)})
	methodCPindex := uint16(len(CP.CpIndex) - 1)

	genCode = append(genCode, opcodes.INVOKESPECIAL)
	hiByte = uint8(methodCPindex >> 8)
	loByte = uint8(methodCPindex)
	genCode = append(genCode, hiByte)
	genCode = append(genCode, loByte)
	genCode = append(genCode, opcodes.ATHROW)
	return genCode
}

// minimalAbort is the exception thrown when the frame info is not available,
// such as during start-up, when the main class can't be found, etc.
func minimalAbort(whichException int, msg string) {
	var stack string
	bytes := debug.Stack()
	if len(bytes) > 0 {
		stack = string(bytes)
	} else {
		stack = ""
	}
	glob := globals.GetGlobalRef()
	glob.ErrorGoStack = stack
	errMsg := fmt.Sprintf("%s: %s", excNames.JVMexceptionNames[whichException], msg)
	fmt.Fprintln(os.Stderr, errMsg)
	// errMsg := fmt.Sprintf("[ThrowEx][minimalAbort] %s", msg)
	// ShowPanicCause(errMsg)
	// ShowFrameStack(&thread.ExecThread{})
	ShowGoStackTrace(nil)
	_ = shutdown.Exit(shutdown.APP_EXCEPTION)
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/opcodes"
	"jacobin/shutdown"
	"jacobin/stringPool"
	"jacobin/thread"
	"jacobin/util"
	"runtime/debug"
)

// This file contains support functions for throwing exceptions from within
// Jacobin. That is, situations in which Jacobin itself is throwing the error,
// rather than the application. Typically, this is for errors/exceptions in the
// operation of the JVM, and for a few occasional user errors, such as
// divide by zero.
//
// We are here duplicating how in-application throws/catches are handled. To
// accomplish this, we generate bytecodes which are then placed in the frame of
// the current thread.
func ThrowEx(which int, msg string, f *frames.Frame) {

	// If tracing, announce.
	helloMsg := fmt.Sprintf("[ThrowEx] Arrived, which: %d, msg: %s", which, msg)
	//log.Log(helloMsg, log.TRACE_INST)
	log.Log(helloMsg, log.SEVERE)

	// If in a unit test, log a severe message and return.
	glob := globals.GetGlobalRef()
	if glob.JacobinName == "test" {
		errMsg := fmt.Sprintf("[ThrowEx][test] %s", msg)
		log.Log(errMsg, log.SEVERE)
		return
	}

	// Frame pointer provided?
	if f == nil {
		minimalAbort(msg)
	}

	// the name of the exception as shown to the user
	exceptionNameForUser := JVMexceptionNames[which]

	// // the name of the class that implements this exception
	// exceptionClassName := util.ConvertInternalClassNameToFilename(exceptionNameForUser)

	// the internal format used in the constant pool
	exceptionCPname := util.ConvertClassFilenameToInternalFormat(exceptionNameForUser)

	// capture the PC where the exception was thrown (saved b/c later we modify the value of f.PC)
	f.ExceptionPC = f.PC

	th := glob.Threads[f.Thread].(*thread.ExecThread)
	fs := th.Stack

	// find out if the exception is caught and if so point to the catch code
	// catchFrame, catchPC := FindExceptionFrame(f, exceptionCPname, f.ExceptionPC)
	catchFrame, catchPC := FindCatchFrame(fs, exceptionCPname, f.ExceptionPC)
	if catchFrame != nil {
		// at this point, we know that the exception was caught
		// and that the top of the frame stack holds the frame
		// containing the catch logic, referred to here as the catchFrame.
		// now, set up the execution of the catch code by:
		// 1. creating a new objRef for the exception
		// 2. pushing the objRef on the stack of the frame
		// 3. setting the PC to point to the catch code (which expects the objRef at TOS)
		th = glob.Threads[f.Thread].(*thread.ExecThread)
		fs = th.Stack
		objRef, _ := glob.FuncInstantiateClass(exceptionCPname, fs)
		catchFrame.TOS = 0
		catchFrame.OpStack[0] = objRef // push the objRef
		catchFrame.PC = catchPC - 1    // -1 because the loop in run.go will increment PC after this code block's return
		return
	}

	// if the exception was not caught...
	genCode := generateThrowBytecodes(f, exceptionCPname, msg)

	// append the genCode to the bytecode of the current method in the frame
	// and set the PC to point to it.
	endPoint := len(f.Meth)
	f.Meth = append(f.Meth, genCode...)
	f.PC = endPoint
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

	// start converting previous work into bytecode
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

func minimalAbort(msg string) {
	var stack string
	bytes := debug.Stack()
	if len(bytes) > 0 {
		stack = string(bytes)
	} else {
		stack = ""
	}
	glob := globals.GetGlobalRef()
	glob.ErrorGoStack = stack
	errMsg := fmt.Sprintf("[ThrowEx][minimalAbort] %s", msg)
	ShowPanicCause(errMsg)
	ShowFrameStack(&thread.ExecThread{})
	ShowGoStackTrace(nil)
	_ = shutdown.Exit(shutdown.APP_EXCEPTION)
}

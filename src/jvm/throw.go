/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/opcodes"
	"jacobin/util"
	"os"
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
func throw(which int, msg string, f *frames.Frame) {
	// the name of the exception as shown to the user
	exceptionNameForUser := exceptions.JVMexceptionNames[which]

	// the name of the class that implements this exception
	exceptionClassName := util.ConvertInternalClassNameToFilename(exceptionNameForUser)

	// the internal format used in the constant pool
	exceptionCPname := util.ConvertClassFilenameToInternalFormat(exceptionNameForUser)

	// the functionality we generate bytecodes for is (using a NPE as an example):
	// 0: new           #7                  // class java/lang/NullPointerException
	// 3: dup
	// 4: ldc           #9                  // String  (the msg passed into this function)
	// 6: invokespecial #11                 // Method java/lang/NullPointerException."<init>":(Ljava/lang/String;)V
	// 9: athrow
	//
	// Note that to do this, we need to twiddle with the constant pool as well

	CP := f.CP.(*classloader.CPool)
	// first add an entry to the UTF8 entries containing the exception class name
	CP.Utf8Refs = append(CP.Utf8Refs, exceptionCPname)
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.UTF8, Slot: uint16(len(CP.Utf8Refs) - 1)})
	exceptionClassCPindex := uint16(len(CP.CpIndex) - 1)
	// then add a classref entry for the exception
	CP.ClassRefs = append(CP.ClassRefs, uint16(len(CP.CpIndex)-1)) // point to the UTF8 entry
	CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
		Type: classloader.ClassRef, Slot: uint16(len(CP.ClassRefs) - 1)})

	// start converting previous work into bytecode
	var genCode []byte
	genCode = append(genCode, opcodes.NOP) // the first bytecode is skipped by the JVM
	genCode = append(genCode, opcodes.NEW)
	genCode = append(genCode, uint8(len(CP.CpIndex)-2))
	genCode = append(genCode, opcodes.DUP)

	// now load the error message, if any
	if msg != "" {
		CP.Utf8Refs = append(CP.Utf8Refs, msg)
		ut8MsgIndex := uint16(len(CP.Utf8Refs) - 1)
		CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
			Type: classloader.UTF8, Slot: ut8MsgIndex})
		CP.CpIndex = append(CP.CpIndex, classloader.CpEntry{
			Type: classloader.StringConst, Slot: uint16(len(CP.CpIndex) - 1)})
		stringMsgIndex := uint16(len(CP.CpIndex) - 1)
		if stringMsgIndex < 256 {
			genCode = append(genCode, opcodes.LDC)
			genCode = append(genCode, uint8(stringMsgIndex))
		} else {
			// if the index is > 255, we need to use LDC_W and a two-byte index
			hiByte := uint8(stringMsgIndex >> 8)
			loByte := uint8(stringMsgIndex)
			genCode = append(genCode, opcodes.LDC_W)
			genCode = append(genCode, hiByte)
			genCode = append(genCode, loByte)
		}
	}

	// now, set up the CP entries for INVOKESPECIAL. This includes a MethodRef
	// which points to the previous classRef and to a name and type record, which
	// itself pointos to a UTF8 entry for the method name and a UTF8 entry for
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
	// methodCPindex := uint16(len(CP.CpIndex) - 1)  // CURR: resume with set up of invokespecial

	fmt.Fprintf(os.Stderr, "Throwing exception: %s, internal name: %s\n",
		exceptionClassName, exceptionCPname)
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/frames"
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
	exceptionInternalName := exceptions.JVMexceptionNames[which]
	exceptionClassName := util.ConvertInternalClassNameToFilename(exceptionInternalName)
	exceptionCPname := util.ConvertClassFilenameToInternalFormat(exceptionInternalName)
	// the functionality we generate bytecodes for is (using NPE as an example):
	// 0: new           #7                  // class java/lang/NullPointerException
	// 3: dup
	// 4: invokespecial #9                  // Method java/lang/NullPointerException."<init>":()V
	// 7: athrow
	//
	// Note that to do this, we need to twiddle with the constant pool as well

	fmt.Fprintf(os.Stderr, "Throwing exception: %s, internal name: %s\n",
		exceptionClassName, exceptionCPname)
}

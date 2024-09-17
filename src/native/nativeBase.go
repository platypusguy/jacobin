/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"errors"
	"os"
	"unsafe"
)

var CaughtNativeFunctionException = errors.New("caught native function exception")

// Native function to library handle table
// Input: Native function name.
// Output: Library file handle.
var nfToLibTable = map[string]uintptr{}

// Native function to template function handle table
// Input: Native function name.
// Output: Template function handle.
var nfToTmplTable = map[string]typeTemplateFunction{}

// Native function error block.
type NativeErrBlk struct {
	ExceptionType int
	ErrMsg        string
}

// Type definition for all the template functions
type typeTemplateFunction func(libHandle uintptr, functionName string, params []interface{}, tracing bool) interface{}

// Argument types for template functions.
type (
	NFboolean   uint8
	NFbyte      uint8
	NFchar      uint16
	NFshort     int16
	NFint       int32
	NFuint      uint32
	NFlong      int64
	NFfloat     float32
	NFdouble    float64
	NFbyteArray unsafe.Pointer
	NFobject    unsafe.Pointer
)

// Struct for CreateJvm.
type t_JavaVMInitArgs struct {
	version            NFint
	nOptions           NFint
	JavaVMOption       uintptr
	ignoreUnrecognized NFboolean
}

// JVM initialisation parameters.
var JavaVMInitArgs = t_JavaVMInitArgs{version: 0x00090000, nOptions: 0, JavaVMOption: 0, ignoreUnrecognized: 0}

// O/S stuff.
var OperSys string                           // One of: "darwin", "linux", "unix", "windows"
var WindowsOS = false                        // true only if OperSys = "windows"
var PathDirLibs string                       // Directory of the more common JVM libraries (E.g. libzip.so)
var PathLibjvm string                        // Full path of libjvm.so
var PathLibjava string                       // Full path of libjava.so
var FileExt string                           // File extension of a library file: "so" (Linux and Unix), "dll" (Windows), "dylib" (MacOS)
var SepPathString = string(os.PathSeparator) // ";" (Windows) or ":" (everybody else)
var HandleLibjvm uintptr                     // Handle of the open libjvm
var HandleLibjava uintptr                    // Handle of the open libjava
var HandleJVM uintptr                        // Handle of the created JVM
var HandleENV uintptr                        // Handle of the JNI environment

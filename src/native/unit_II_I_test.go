/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"testing"
)

func storeLibHandle(t *testing.T, argLib, argFunction string) bool {
	var lib string
	if WindowsOS {
		lib = PathDirLibs + SepPathString + argLib + "." + FileExt
	} else {
		lib = PathDirLibs + SepPathString + "lib" + argLib + "." + FileExt
	}
	handle := ConnectLibrary(lib)
	if handle == 0 {
		return false
	}
	t.Logf("storeLibHandle: lib: %s, function: %s\n", lib, argFunction)
	nfToLibTable[argFunction] = handle
	return true
}

func Test_II_I(t *testing.T) {
	tracing := true

	// Initialize jacobin and set up a dummy frame stack.
	globals.InitGlobals("test")
	log.Init()
	if tracing {
		log.SetLogLevel(log.TRACE_INST)
	}

	// Perform native initialisation.
	if !nativeInit() {
		t.Error("nativeInit() failed")
	}
	t.Log("nativeInit ok")

	// SIMULATION: Store some library handles.
	if !storeLibHandle(t, "awt", "apples") {
		t.Error("storeLibHandle() failed")
	}
	if !storeLibHandle(t, "net", "bananas") {
		t.Error("storeLibHandle() failed")
	}
	if !storeLibHandle(t, "zip", "Java_java_util_zip_CRC32_update") {
		t.Error("storeLibHandle() failed")
	}

	// Create a stack frame with one frame.
	frame := frames.CreateFrame(10)
	frame.Thread = 1
	fs := frames.CreateFrameStack()
	fs.PushFront(frame)

	// Call RunNativeFunction.
	params := make([]interface{}, 2)
	params[1] = NFint(0)
	params[0] = NFint('A')
	expected := NFuint(0xd3d99e8b)
	ret := RunNativeFunction(fs, "CRC32", "Java_java_util_zip_CRC32_update", "(II)I", &params, tracing)
	switch ret.(type) {
	case NFint:
		observed := NFuint(ret.(NFint))
		if observed != expected {
			t.Errorf("Oops, expected: 0x%08x, observed: 0x%08x\n", expected, NFuint(observed))
		} else {
			t.Logf("Success, observed = expected = 0x%08x\n", expected)
		}
	case NativeErrBlk:
		t.Errorf(ret.(NativeErrBlk).ErrMsg)
	default:
		t.Errorf("Unexpected observed type: %T\n", ret)
	}
}

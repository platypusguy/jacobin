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

func Test_II_I(t *testing.T) {
	tracing := true
	functionName := "Java_java_util_zip_CRC32_update"

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
	if !storeLibHandle("awt", "apples") {
		t.Error("storeLibHandle() failed")
	}
	if !storeLibHandle("net", "bananas") {
		t.Error("storeLibHandle() failed")
	}

	// tell purego that zip lib/DLL contains functionName
	if !storeLibHandle("zip", functionName) {
		t.Error("storeLibHandle() failed")
	}

	// Create a stack frame with one frame.
	frame := frames.CreateFrame(10)
	frame.Thread = 1
	fs := frames.CreateFrameStack()
	fs.PushFront(frame)

	// Call RunNativeFunction.
	params := make([]interface{}, 2)
	params[1] = int64(0)   // starting point for the CRC computation
	params[0] = int64('A') // 'A' the value we want the CRC computed for
	expected := NFuint(0xd3d99e8b)
	// fs = frame stack,
	// "CRC32" stand-in for the actual class name
	// "Java_java_util_zip_CRC32_update" native function name
	// (II)I native function type
	// ptr param array
	// tracing parameter for Jacobin's tracing purposes
	ret := RunNativeFunction(fs, "java/util/CRC32", functionName, "(II)I", &params, tracing)
	switch ret.(type) {
	case int64:
		observed := NFuint(ret.(int64))
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

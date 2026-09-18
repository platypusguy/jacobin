/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"strings"
	"testing"
)

// --- ThrowEx: "test" mode branch (glob.JacobinName == "test") ---

func TestThrowEx_TestModeWithFrame(t *testing.T) {
	globals.InitGlobals("test")

	f := frames.CreateFrame(1)
	f.ClName = "com/test/SomeClass"
	f.MethName = "someMethod"
	f.MethType = "()V"

	out := captureStderr(func() {
		result := ThrowEx(excNames.ArithmeticException, "boom", f)
		if result != NotCaught {
			t.Errorf("expected NotCaught in test mode, got %v", result)
		}
	})
	if !strings.Contains(out, "java.lang.ArithmeticException") || !strings.Contains(out, "boom") {
		t.Errorf("expected message with exception name and msg, got: %s", out)
	}
}

// --- ThrowEx: nil/invalid-thread frame -> MinimalAbort path ---

func TestThrowEx_FrameWithInvalidThread(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")

	f := frames.CreateFrame(1)
	f.ClName = "com/test/SomeClass"
	f.MethName = "someMethod"
	f.MethType = "()V"
	f.Thread = 0 // < 1, so MinimalAbort should be called

	out := captureStderr(func() {
		result := ThrowEx(excNames.ArithmeticException, "boom", f)
		if result != NotCaught {
			t.Errorf("expected NotCaught, got %v", result)
		}
	})
	if !strings.Contains(out, "java.lang.ArithmeticException") || !strings.Contains(out, "boom") {
		t.Errorf("expected MinimalAbort message, got: %s", out)
	}
}

// --- ThrowEx: thread not registered in glob.Threads ---
// This is a documented pre-existing edge case: after MinimalAbort logs the
// "thread not found" error, execution continues (MinimalAbort does not stop
// the current goroutine in test mode) and subsequently panics on a nil
// *object.Object dereference. We recover from that panic to keep the test
// process alive while still verifying the logged error message.
func TestThrowEx_ThreadNotFoundInGlobals(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	glob := globals.GetGlobalRef()
	glob.Threads = make(map[int]interface{})

	f := frames.CreateFrame(1)
	f.ClName = "com/test/SomeClass"
	f.MethName = "someMethod"
	f.MethType = "()V"
	f.Thread = 12345 // not present in glob.Threads

	out := captureStderr(func() {
		defer func() {
			_ = recover() // the code continues after MinimalAbort and later panics
		}()
		ThrowEx(excNames.ArithmeticException, "boom", f)
	})
	if !strings.Contains(out, "glob.Threads index not found") {
		t.Errorf("expected 'glob.Threads index not found' message, got: %s", out)
	}
}

// --- ThrowEx: thread entry present but of the wrong type ---
func TestThrowEx_ThreadEntryWrongType(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	glob := globals.GetGlobalRef()
	glob.Threads = make(map[int]interface{})
	glob.Threads[999] = "not-an-object"

	f := frames.CreateFrame(1)
	f.ClName = "com/test/SomeClass"
	f.MethName = "someMethod"
	f.MethType = "()V"
	f.Thread = 999

	out := captureStderr(func() {
		defer func() {
			_ = recover()
		}()
		ThrowEx(excNames.ArithmeticException, "boom", f)
	})
	if !strings.Contains(out, "glob.Threads entry corrupted") {
		t.Errorf("expected 'glob.Threads entry corrupted' message, got: %s", out)
	}
}

// --- ThrowEx: full path, exception caught ---

func TestThrowEx_CaughtException(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	glob := globals.GetGlobalRef()

	glob.FuncInstantiateClass = func(_ string, _ *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		return obj, nil
	}
	glob.FuncFillInStackTrace = func(_ []any) any { return nil }

	cp := buildCPwithClassRef("java/lang/ArithmeticException")
	excs := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 42, CatchType: 1},
	}

	// build the catch frame with a matching exception handler
	fqn := "com/test/CaughtOwner.method()V"
	registerJavaMethod(fqn, cp, excs)

	f := frames.CreateFrame(1)
	f.ClName = "com/test/CaughtOwner"
	f.MethName = "method"
	f.MethType = "()V"
	f.CP = cp
	f.PC = 3
	f.ExceptionPC = -1

	threadObj := object.MakeEmptyObject()
	fs := frames.CreateFrameStack()
	_ = frames.PushFrame(fs, f)
	threadObj.FieldTable["framestack"] = object.Field{Fvalue: fs}
	f.Thread = 7

	glob.Threads = make(map[int]interface{})
	glob.Threads[7] = threadObj

	out := captureStderr(func() {
		result := ThrowEx(excNames.ArithmeticException, "boom", f)
		if result != Caught {
			t.Errorf("expected Caught, got %v", result)
		}
	})
	_ = out

	if f.PC != 42 {
		t.Errorf("expected frame PC to be set to handler pc 42, got %d", f.PC)
	}
	if f.TOS != 0 {
		t.Errorf("expected TOS to be reset to 0, got %d", f.TOS)
	}
	if f.OpStack[0] == nil {
		t.Errorf("expected the exception object reference to be pushed on the op stack")
	}
}

// --- ThrowEx: full path, exception NOT caught ---

func TestThrowEx_UncaughtException(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	glob := globals.GetGlobalRef()
	glob.StrictJDK = false

	glob.FuncInstantiateClass = func(_ string, _ *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		stackTraceObj := object.MakeEmptyObject()
		stackTraceObj.FieldTable["value"] = object.Field{Fvalue: []*object.Object{}}
		obj.FieldTable["stackTrace"] = object.Field{Fvalue: stackTraceObj}
		return obj, nil
	}
	glob.FuncFillInStackTrace = func(_ []any) any { return nil }

	// no exception table at all -> never caught
	fqn := "com/test/UncaughtOwner.method()V"
	registerJavaMethod(fqn, nil, nil)

	f := frames.CreateFrame(1)
	f.ClName = "com/test/UncaughtOwner"
	f.MethName = "method"
	f.MethType = "()V"
	f.PC = 3
	f.ExceptionPC = -1
	f.Thread = 8

	threadObj := object.MakeEmptyObject()
	fs := frames.CreateFrameStack()
	_ = frames.PushFrame(fs, f)
	threadObj.FieldTable["framestack"] = object.Field{Fvalue: fs}

	glob.Threads = make(map[int]interface{})
	glob.Threads[8] = threadObj

	out := captureStderr(func() {
		result := ThrowEx(excNames.ArithmeticException, "boom", f)
		if result != NotCaught {
			t.Errorf("expected NotCaught, got %v", result)
		}
	})
	if !strings.Contains(out, "ArithmeticException") {
		t.Errorf("expected uncaught exception message, got: %s", out)
	}
}

// --- ShowJVMstackTrace ---

func TestShowJVMstackTrace(t *testing.T) {
	globals.InitGlobals("test")
	glob := globals.GetGlobalRef()
	glob.StrictJDK = false

	entry := object.MakeEmptyObject()
	entry.FieldTable["declaringClass"] = object.Field{Fvalue: "com/test/SomeClass"}
	entry.FieldTable["methodName"] = object.Field{Fvalue: "someMethod"}
	entry.FieldTable["fileName"] = object.Field{Fvalue: "SomeClass.java"}
	entry.FieldTable["sourceLine"] = object.Field{Fvalue: "10"}
	entry.FieldTable["PC"] = object.Field{Fvalue: "3"}

	out := captureStderr(func() {
		ShowJVMstackTrace([]*object.Object{entry}, glob)
	})
	if !strings.Contains(out, "com/test/SomeClass.someMethod") || !strings.Contains(out, "SomeClass.java:10") {
		t.Errorf("expected formatted stack trace entry, got: %s", out)
	}
}

func TestShowJVMstackTrace_NoSourceLine(t *testing.T) {
	globals.InitGlobals("test")
	glob := globals.GetGlobalRef()
	glob.StrictJDK = true

	entry := object.MakeEmptyObject()
	entry.FieldTable["declaringClass"] = object.Field{Fvalue: "com/test/SomeClass"}
	entry.FieldTable["methodName"] = object.Field{Fvalue: "someMethod"}
	entry.FieldTable["fileName"] = object.Field{Fvalue: "SomeClass.java"}
	entry.FieldTable["sourceLine"] = object.Field{Fvalue: ""}
	entry.FieldTable["PC"] = object.Field{Fvalue: "3"}

	out := captureStderr(func() {
		ShowJVMstackTrace([]*object.Object{entry}, glob)
	})
	if !strings.Contains(out, "com.test.SomeClass.someMethod") {
		t.Errorf("expected user-formatted declaring class in strict JDK mode, got: %s", out)
	}
}

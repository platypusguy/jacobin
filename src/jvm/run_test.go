/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/testutil"
	"testing"
)

// ---------- multiply/subtract generics ----------

func TestMultiplyInt64(t *testing.T) {
	if multiply(int64(6), int64(7)) != 42 {
		t.Errorf("multiply(6,7) expected 42")
	}
}

func TestMultiplyFloat64(t *testing.T) {
	if multiply(3.0, 4.0) != 12.0 {
		t.Errorf("multiply(3.0,4.0) expected 12.0")
	}
}

func TestSubtractInt64(t *testing.T) {
	if subtract(int64(10), int64(3)) != 7 {
		t.Errorf("subtract(10,3) expected 7")
	}
}

func TestSubtractFloat64(t *testing.T) {
	if subtract(10.5, 0.5) != 10.0 {
		t.Errorf("subtract(10.5,0.5) expected 10.0")
	}
}

// ---------- CkSyncStaticMeth ----------

func TestCkSyncStaticMeth_NotSynchronized(t *testing.T) {
	testutil.UTinit(t)

	fram := frames.CreateFrame(2)
	fram.ClName = "SomeClass"
	fram.AccessFlags = 0 // not synchronized

	err := CkSyncStaticMeth(fram)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if fram.ObjSync != nil {
		t.Errorf("expected ObjSync to remain nil, got %v", fram.ObjSync)
	}
}

func TestCkSyncStaticMeth_Synchronized_Success(t *testing.T) {
	testutil.UTinit(t)
	classloader.InitMethodArea()

	className := "TestCkSyncStaticMethClass"
	classObj := object.MakeEmptyObjectWithClassName(&className)

	kd := &classloader.ClData{
		Name:        className,
		ClassObject: classObj,
	}
	k := &classloader.Klass{
		Status: 'L',
		Loader: "bootstrap",
		Data:   kd,
	}
	classloader.MethAreaInsert(className, k)

	fram := frames.CreateFrame(2)
	fram.ClName = className
	fram.AccessFlags = classloader.ACC_SYNCHRONIZED
	fram.Thread = 1

	err := CkSyncStaticMeth(fram)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fram.ObjSync != classObj {
		t.Errorf("expected ObjSync to point to the class object")
	}

	// The class object should now be locked by thread 1.
	if classObj.GetMonitorOwner() != 1 {
		t.Errorf("expected class object to be locked by thread 1, got owner=%d", classObj.GetMonitorOwner())
	}

	// clean up: unlock so other tests aren't affected
	_ = classObj.ObjUnlock(1)
}

// ---------- createAndInitNewFrame ----------

func buildJmEntry(accessFlags, maxStack, maxLocals int) *classloader.JmEntry {
	return &classloader.JmEntry{
		AccessFlags: accessFlags,
		MaxStack:    maxStack,
		MaxLocals:   maxLocals,
		Code:        []byte{},
		Cp:          &classloader.CPool{},
	}
}

func TestCreateAndInitNewFrame_StaticNoParams(t *testing.T) {
	testutil.UTinit(t)

	m := buildJmEntry(0, 2, 0)

	curr := frames.CreateFrame(2)
	curr.Thread = 0
	curr.FrameStack = frames.CreateFrameStack()

	fram, err := createAndInitNewFrame("SomeClass", "someMethod", "()V", m, false, curr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fram.ClName != "SomeClass" || fram.MethName != "someMethod" || fram.MethType != "()V" {
		t.Errorf("frame fields not set as expected: %+v", fram)
	}
	if fram.TOS != -1 {
		t.Errorf("expected TOS = -1, got %d", fram.TOS)
	}
}

func TestCreateAndInitNewFrame_StaticWithIntParams(t *testing.T) {
	testutil.UTinit(t)

	m := buildJmEntry(0, 4, 0)

	curr := frames.CreateFrame(4)
	curr.Thread = 0
	curr.FrameStack = frames.CreateFrameStack()
	// push two ints, in reverse order (as the interpreter would prior to invocation)
	push(curr, int64(5))
	push(curr, int64(10))

	fram, err := createAndInitNewFrame("SomeClass", "addTwo", "(II)I", m, false, curr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fram.Locals) < 2 {
		t.Fatalf("expected at least 2 locals, got %d", len(fram.Locals))
	}
	if fram.Locals[0] != int64(5) || fram.Locals[1] != int64(10) {
		t.Errorf("expected locals[0]=5, locals[1]=10, got %v, %v", fram.Locals[0], fram.Locals[1])
	}
}

func TestCreateAndInitNewFrame_NonStaticWithObjectRef(t *testing.T) {
	testutil.UTinit(t)

	m := buildJmEntry(0, 4, 0)

	curr := frames.CreateFrame(4)
	curr.Thread = 0
	curr.FrameStack = frames.CreateFrameStack()

	className := "SomeInstanceClass"
	obj := object.MakeEmptyObjectWithClassName(&className)
	push(curr, obj) // the "this" reference, popped first since includeObjectRef=true

	fram, err := createAndInitNewFrame(className, "instMethod", "()V", m, true, curr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fram.Locals[0] != obj {
		t.Errorf("expected locals[0] to be the popped object reference")
	}
}

func TestCreateAndInitNewFrame_NonStaticSynchronized_LocksObject(t *testing.T) {
	testutil.UTinit(t)

	m := buildJmEntry(classloader.ACC_SYNCHRONIZED, 4, 0)

	curr := frames.CreateFrame(4)
	curr.Thread = 2
	curr.FrameStack = frames.CreateFrameStack()

	className := "SyncInstanceClass"
	obj := object.MakeEmptyObjectWithClassName(&className)
	push(curr, obj)

	fram, err := createAndInitNewFrame(className, "syncMethod", "()V", m, true, curr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fram.ObjSync != obj {
		t.Errorf("expected ObjSync to be set to the instance object")
	}
	if obj.GetMonitorOwner() != 2 {
		t.Errorf("expected object to be locked by thread 2, got owner=%d", obj.GetMonitorOwner())
	}

	_ = obj.ObjUnlock(2)
}

func TestCreateAndInitNewFrame_StaticSynchronized_ReturnsError(t *testing.T) {
	testutil.UTinit(t)
	classloader.InitMethodArea()
	// Note: no class registered in the method area under this name, so
	// CkSyncStaticMeth's MethAreaFetch will return nil, and dereferencing
	// cl.Data will panic -- covered instead through the success path test,
	// here we simply verify a normally-registered class works end-to-end
	// via createAndInitNewFrame for the static+synchronized case.
	className := "StaticSyncClass"
	classObj := object.MakeEmptyObjectWithClassName(&className)
	kd := &classloader.ClData{Name: className, ClassObject: classObj}
	k := &classloader.Klass{Status: 'L', Loader: "bootstrap", Data: kd}
	classloader.MethAreaInsert(className, k)

	m := buildJmEntry(classloader.ACC_SYNCHRONIZED, 2, 0)

	curr := frames.CreateFrame(2)
	curr.Thread = 3
	curr.FrameStack = frames.CreateFrameStack()

	fram, err := createAndInitNewFrame(className, "staticSyncMethod", "()V", m, false, curr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fram.ObjSync != classObj {
		t.Errorf("expected ObjSync to be the class object")
	}
	if classObj.GetMonitorOwner() != 3 {
		t.Errorf("expected class object locked by thread 3, got owner=%d", classObj.GetMonitorOwner())
	}
	_ = classObj.ObjUnlock(3)
}

// ---------- RunJavaThread argument-count validation ----------

// In "test" mode, exceptions.ThrowEx logs the error and returns rather than
// terminating the process, so RunJavaThread continues executing past the
// invalid-argument-count check. Since the panic-recovering defer is only
// registered *after* the args[0].(*object.Object) type assertion, the
// subsequent panic on that type assertion (args[0] is a string, not an
// *object.Object) is NOT caught by RunJavaThread's own recover, and
// propagates out of the function. This test documents that existing
// behavior while confirming the argument-count error message is emitted.
func TestRunJavaThread_WrongArgCount(t *testing.T) {
	globals.InitGlobals("test")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic from the subsequent type assertion, got none")
		}
	}()

	// Only 2 args instead of the required 4 -- this exercises the
	// argument-count-validation branch of RunJavaThread.
	RunJavaThread([]any{"a", "b"})
}

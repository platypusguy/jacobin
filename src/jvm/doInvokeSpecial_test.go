/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/gfunction/jacobinSrc"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"os"
	"strings"
	"testing"
)

// This file contains unit tests for doInvokeSpecial.go (opcode 0xB7 INVOKESPECIAL).

// helper: build a minimal CP with a single MethodRef entry pointing at
// className.methodName.methodType, plus the surrounding ClassRef/NameAndType/UTF8 entries.
func buildInvokeSpecialCP(className, methodName, methodType string) *classloader.CPool {
	CP := &classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	CP.MethodRefs[0] = classloader.MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	CP.ClassRefs[0] = stringPool.GetStringIndex(&className)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{NameIndex: 4, DescIndex: 5}

	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = methodName
	CP.Utf8Refs[1] = methodType

	classloader.ResolveCPmethRefs(CP)
	return CP
}

// ensure the test class "jacobin/src/test/Object" exists in the MethArea so that
// FetchMethodAndCP() does not try to load it from disk.
func ensureTestClassLoaded() {
	jacobinSrc.CheckTestGfunctionsLoaded()
}

// ensure the test class "jacobin/src/test/Object" is present in the MethArea with
// a non-nil Data/MethodTable, so FetchMethodAndCP can safely search it when the
// requested method is not already present in the global MTable.
func ensureTestClassLoadedWithMethodTable() {
	className := "jacobin/src/test/Object"
	clData := classloader.ClData{
		Name:            className,
		NameIndex:       stringPool.GetStringIndex(&className),
		SuperclassIndex: 0,
		MethodTable:     make(map[string]*classloader.Method),
		CP:              classloader.CPool{},
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "bootstrap",
		Data:   &clData,
	}
	classloader.MethAreaInsert(className, &k)
}

// INVOKESPECIAL: method not found anywhere -> NoSuchMethodError
func TestDoInvokeSpecial_MethodNotFound(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureTestClassLoadedWithMethodTable()

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeSpecialCP("jacobin/src/test/Object", "noSuchMethodAtAll", "()V")
	f.CP = CP

	push(f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKESPECIAL: Class method not found") {
		t.Errorf("Expected NoSuchMethodError, got: %s", errMsg)
	}
}

// INVOKESPECIAL: method is native -> UnsupportedOperationException
func TestDoInvokeSpecial_NativeMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureTestClassLoaded()

	methName := "nativeTestMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: classloader.ACC_NATIVE,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeSpecialCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	push(f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKESPECIAL: Native method requested") {
		t.Errorf("Expected UnsupportedOperationException for native method, got: %s", errMsg)
	}
}

// INVOKESPECIAL: method has empty code (abstract) -> AbstractMethodError
func TestDoInvokeSpecial_AbstractMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureTestClassLoaded()

	methName := "abstractTestMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{}, // empty code => abstract
		},
	})

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeSpecialCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	push(f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKESPECIAL: J class method code is empty") {
		t.Errorf("Expected AbstractMethodError, got: %s", errMsg)
	}
}

// INVOKESPECIAL: successful dispatch of a J method -> new frame is pushed and runs,
// returning control to the caller frame without error.
func TestDoInvokeSpecial_SuccessfulJMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureTestClassLoaded()

	methName := "simpleVoidMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			MaxStack:    2,
			MaxLocals:   1,
			Code:        []byte{opcodes.RETURN}, // just return
		},
	})

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeSpecialCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	push(f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("Expected no error for successful J method invocation, got: %s", errMsg)
	}

	if f.TOS != -1 {
		t.Errorf("Expected op stack empty (TOS -1) after invocation, got: %d", f.TOS)
	}
}

// INVOKESPECIAL: method ref decoded via the interface-ref path (private/default interface method)
func TestDoInvokeSpecial_InterfaceRefPath(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureTestClassLoaded()

	methName := "ifaceDefaultMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := &classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.Interface, Slot: 0}
	CP.InterfaceRefs = make([]classloader.InterfaceRefEntry, 1)
	CP.InterfaceRefs[0] = classloader.InterfaceRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	className := "jacobin/src/test/Object"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&className)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{NameIndex: 4, DescIndex: 5}

	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = methName
	CP.Utf8Refs[1] = methType

	f.CP = CP

	push(f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("Expected no error for interface-ref INVOKESPECIAL, got: %s", errMsg)
	}
}

// INVOKESPECIAL: null object reference for a G function -> NullPointerException
func TestDoInvokeSpecial_GFunctionNullObjectRef(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeSpecialCP("jacobin/src/test/Object", "test", "()Ljava/lang/Object;")
	f.CP = CP

	// push a nil reference, instead of a valid *object.Object
	push(f, nil)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKESPECIAL: Stack reference object is nil") {
		t.Errorf("Expected NullPointerException for nil receiver, got: %s", errMsg)
	}
}

// INVOKESPECIAL: the CP slot has already been cached as a CachedMeth by a previous
// INVOKEVIRTUAL/INVOKESTATIC on the same call site -> the fast path should be used
// and the invocation should still succeed.
func TestDoInvokeSpecial_CachedMethFastPath(t *testing.T) {
	globals.InitGlobals("test")
	globals.CacheMeths = true

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureTestClassLoaded()

	methName := "cachedVoidMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeSpecialCP("jacobin/src/test/Object", methName, methType)

	// Pre-populate the CachedMethods slot and mark CP.CpIndex[1] as CachedMeth,
	// simulating a previous resolution for this call site.
	className := "jacobin/src/test/Object"
	cachedEntry := classloader.MTentry{
		MType:     'J',
		Meth:      classloader.JmEntry{AccessFlags: 0, Code: []byte{opcodes.RETURN}},
		MethClass: stringPool.GetStringIndex(&className),
		MethName:  stringPool.GetStringIndex(&methName),
		MethType:  stringPool.GetStringIndex(&methType),
	}
	CP.CachedMethods = append(CP.CachedMethods, cachedEntry)
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.CachedMeth, Slot: uint16(len(CP.CachedMethods) - 1)}

	f.CP = CP

	push(f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("Expected no error for cached-method fast path, got: %s", errMsg)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
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

// This file contains unit tests for doInvokeVirtual.go (opcode 0xB6 INVOKEVIRTUAL).

// helper: build a minimal CP with a single MethodRef entry pointing at
// className.methodName.methodType, plus the surrounding ClassRef/NameAndType/UTF8 entries.
func buildInvokeVirtualCP(className, methodName, methodType string) *classloader.CPool {
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

// INVOKEVIRTUAL: method not found anywhere -> NoSuchMethodError
// The class must have at least one interface entry so the code enters the
// interface-search block and ultimately throws NoSuchMethodError.
func TestDoInvokeVirtual_MethodNotFound(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}

	// Register a class that has one (dummy) interface so the interface-search
	// block is entered and NoSuchMethodError is thrown when the method is absent.
	ifaceName := "jacobin/src/test/IFace"
	className := "jacobin/src/test/WithIface"
	clData := classloader.ClData{
		Name:            className,
		NameIndex:       stringPool.GetStringIndex(&className),
		SuperclassIndex: 0,
		Interfaces:      []uint16{uint16(stringPool.GetStringIndex(&ifaceName))},
		MethodTable:     make(map[string]*classloader.Method),
		CP:              classloader.CPool{},
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "bootstrap",
		Data:   &clData,
	}
	classloader.MethAreaInsert(className, &k)

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP(className, "noSuchMethodAtAll", "()V")
	f.CP = CP

	recvObj := object.MakeEmptyObject()
	recvObj.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEVIRTUAL") {
		t.Errorf("Expected INVOKEVIRTUAL error for missing method, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: method is native -> UnsupportedOperationException
func TestDoInvokeVirtual_NativeMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	methName := "nativeVirtualMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: classloader.ACC_NATIVE,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	recvObj := object.MakeEmptyObject()
	recvObj.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEVIRTUAL: Native method requested") {
		t.Errorf("Expected UnsupportedOperationException for native method, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: method has empty code (abstract) -> AbstractMethodError
func TestDoInvokeVirtual_AbstractMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	methName := "abstractVirtualMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{}, // empty code => abstract
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	recvObj := object.MakeEmptyObject()
	recvObj.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEVIRTUAL: Empty code segment") {
		t.Errorf("Expected AbstractMethodError, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: successful dispatch of a J method -> new frame is pushed and runs,
// returning control to the caller frame without error.
func TestDoInvokeVirtual_SuccessfulJMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	methName := "simpleVirtualVoidMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			MaxStack:    2,
			MaxLocals:   1,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	recvObj := object.MakeEmptyObject()
	recvObj.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj)

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

// INVOKEVIRTUAL: successful dispatch of a G function
func TestDoInvokeVirtual_GFunction(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	// "test" is a G function registered by CheckTestGfunctionsLoaded
	CP := buildInvokeVirtualCP("jacobin/src/test/Object", "test", "()Ljava/lang/Object;")
	f.CP = CP

	className := "jacobin/src/test/Object"
	recvObj := object.MakeEmptyObject()
	recvObj.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("Expected no error for G function invocation, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: cached method fast path (CachedMeth CP entry) -> succeeds without re-lookup
func TestDoInvokeVirtual_CachedMethFastPath(t *testing.T) {
	globals.InitGlobals("test")
	globals.CacheMeths = true

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	methName := "cachedVirtualMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP("jacobin/src/test/Object", methName, methType)

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

	recvObj2 := object.MakeEmptyObject()
	recvObj2.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj2)

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

// INVOKEVIRTUAL: null receiver (non-*object.Object on stack) -> NullPointerException
func TestDoInvokeVirtual_NullReceiver(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	methName := "nullReceiverMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	// Push nil instead of a valid *object.Object
	push(f, nil)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEVIRTUAL") {
		t.Errorf("Expected NullPointerException for nil receiver, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: method caching is triggered on first call (shouldCacheMeth path)
func TestDoInvokeVirtual_MethodCachingOnFirstCall(t *testing.T) {
	globals.InitGlobals("test")
	globals.CacheMeths = true

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	methName := "firstCallCacheMethod"
	methType := "()V"
	fqn := "jacobin/src/test/Object." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			MaxStack:    2,
			MaxLocals:   1,
			Code:        []byte{opcodes.RETURN},
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	// CP entry is a plain MethodRef (not CachedMeth), so shouldCacheMeth will be set true
	CP := buildInvokeVirtualCP("jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	recvObj := object.MakeEmptyObject()
	recvObj.KlassName = stringPool.GetStringIndex(&className)
	push(f, recvObj)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("Expected no error on first-call caching path, got: %s", errMsg)
	}

	// After the call, the CP entry should have been rewritten to CachedMeth
	CP2 := f.CP.(*classloader.CPool)
	CP2.Mutex.RLock()
	entryType := CP2.CpIndex[1].Type
	CP2.Mutex.RUnlock()
	if entryType != classloader.CachedMeth {
		t.Errorf("Expected CP entry to be rewritten to CachedMeth after first call, got type: %d", entryType)
	}
}

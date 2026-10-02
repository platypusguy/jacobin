/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"container/list"
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
func buildInvokeVirtualCP(t *testing.T, className, methodName, methodType string) *classloader.CPool {
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

	err := classloader.ResolveCPmethRefs(CP)
	if err != nil {
		fqn := className + "." + methodName + methodType
		t.Fatalf("buildInvokeVirtualCP: FQN: %s, error: %v", fqn, err)
	}
	return CP
}

// ensure the test class "jacobin/src/test/Object" exists in the MethArea so that
// FetchMethodAndCP() does not try to load it from disk.
func ensureVirtualTestClassLoaded() {
	jacobinSrc.CheckTestGfunctionsLoaded()
}

// ensure the test class "jacobin/src/test/Object" is present in the MethArea with
// a non-nil Data/MethodTable, so FetchMethodAndCP can safely search it when the
// requested method is not already present in the global MTable.
func ensureVirtualTestClassLoadedWithMethodTable() {
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

// registerStaleDispatchClass inserts a bare class entry (no methods of its own) into the
// MethArea for className, so that classloader.FetchMethodAndCP() finds the class already
// loaded and goes on to look up the method in the global MTable instead of trying (and
// failing) to load className from disk.
func registerStaleDispatchClass(className string) {
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

// INVOKEVIRTUAL: method not found anywhere -> NoSuchMethodError
func TestDoInvokeVirtual_MethodNotFound(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureVirtualTestClassLoadedWithMethodTable()

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP(t, "jacobin/src/test/Object", "noSuchVirtualMethodAtAll", "()V")
	f.CP = CP

	className := "jacobin/src/test/Object"
	push(f, object.MakeEmptyObjectWithClassName(&className))

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEVIRTUAL: Class method not found") {
		t.Errorf("Expected NoSuchMethodError, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: declared method is native -> UnsupportedOperationException
func TestDoInvokeVirtual_NativeMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureVirtualTestClassLoaded()

	methName := "nativeVirtualTestMethod"
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

	CP := buildInvokeVirtualCP(t, "jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	push(f, object.MakeEmptyObjectWithClassName(&className))

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

// INVOKEVIRTUAL: declared method has empty code (abstract) -> AbstractMethodError
func TestDoInvokeVirtual_AbstractMethod(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureVirtualTestClassLoaded()

	methName := "abstractVirtualTestMethod"
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

	CP := buildInvokeVirtualCP(t, "jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	push(f, object.MakeEmptyObjectWithClassName(&className))

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEVIRTUAL: J class method code is empty") {
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
	ensureVirtualTestClassLoaded()

	methName := "simpleVirtualVoidMethod"
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

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP(t, "jacobin/src/test/Object", methName, methType)
	f.CP = CP

	className := "jacobin/src/test/Object"
	push(f, object.MakeEmptyObjectWithClassName(&className))

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

// INVOKEVIRTUAL: receiver reference on the operand stack is nil -> NullPointerException
func TestDoInvokeVirtual_NullObjectRef(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureVirtualTestClassLoaded()

	methName := "simpleVirtualVoidMethod2"
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

	CP := buildInvokeVirtualCP(t, "jacobin/src/test/Object", methName, methType)
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
	if !strings.Contains(errMsg, "INVOKEVIRTUAL: Stack reference object is nil") {
		t.Errorf("Expected NullPointerException for nil receiver, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL: the CP slot has already been cached as a CachedMeth by a previous
// call on the same call site -> the fast path should be used and the invocation
// should still succeed.
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
	ensureVirtualTestClassLoaded()

	methName := "cachedVirtualVoidMethod"
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

	CP := buildInvokeVirtualCP(t, "jacobin/src/test/Object", methName, methType)

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

	push(f, object.MakeEmptyObjectWithClassName(&className))

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

// TestInvokeVirtual_StaleDispatchCache_A populates dispatchCache for a receiver class
// "StaleDispatchClassA" with a *successful* run()V method. The next test,
// TestInvokeVirtual_StaleDispatchCache_B, calls globals.InitGlobals("test") again, which
// resets the string pool (back to the same starting index) without clearing dispatchCache.
// If the two tests intern their class/method/descriptor strings in the same relative
// order -- which they do, since both follow the exact same setup sequence -- the string-pool
// index assigned to "StaleDispatchClassB" collides with the one previously assigned to
// "StaleDispatchClassA", and likewise for "run"/"()V". That makes the dispatchKey for B's
// call identical to A's, so dispatchCache.Load returns A's stale *dispatchVal instead of
// resolving B's (different) method.
func TestInvokeVirtual_StaleDispatchCache_A(t *testing.T) {
	globals.InitGlobals("test")

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureVirtualTestClassLoaded()

	className := "StaleDispatchClassA"
	methName := "run"
	methType := "()V"
	registerStaleDispatchClass(className)
	fqn := className + "." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{opcodes.RETURN}, // a valid, non-abstract method
		},
	})

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00, 0x01)

	CP := buildInvokeVirtualCP(t, className, methName, methType)
	f.CP = CP

	push(f, object.MakeEmptyObjectWithClassName(&className))

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	if errMsg := string(msg); errMsg != "" {
		t.Fatalf("Expected no error for successful setup call, got: %s", errMsg)
	}
}

// TestInvokeVirtual_StaleDispatchCache_B defines an unrelated class "StaleDispatchClassB"
// whose run()V method is abstract (empty code, so it MUST raise AbstractMethodError).
// Without clearing dispatchCache after globals.InitGlobals resets the string pool, the
// stale entry left behind by TestInvokeVirtual_StaleDispatchCache_A is returned instead,
// so the call wrongly "succeeds". Calling jvm.ResetDispatchCaches() clears the stale state,
// after which the same call correctly resolves StaleDispatchClassB.run()V and throws
// AbstractMethodError.
func TestInvokeVirtual_StaleDispatchCache_B(t *testing.T) {
	globals.InitGlobals("test") // resets the string pool to its initial state

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	ensureVirtualTestClassLoaded()

	className := "StaleDispatchClassB"
	methName := "run"
	methType := "()V"
	registerStaleDispatchClass(className)
	fqn := className + "." + methName + methType
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{
		MType: 'J',
		Meth: classloader.JmEntry{
			AccessFlags: 0,
			Code:        []byte{}, // empty code => abstract
		},
	})

	buildFrame := func() *list.List {
		f := newFrame(opcodes.INVOKEVIRTUAL)
		f.Meth = append(f.Meth, 0x00, 0x01)
		CP := buildInvokeVirtualCP(t, className, methName, methType)
		f.CP = CP
		push(f, object.MakeEmptyObjectWithClassName(&className))
		fs := frames.CreateFrameStack()
		fs.PushFront(f)
		return fs
	}

	// --- Demonstrate the stale-cache bug: dispatchCache was never cleared, so this
	// call reuses TestInvokeVirtual_StaleDispatchCache_A's cached (successful) result
	// instead of resolving StaleDispatchClassB.run()V, which is abstract.
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	interpret(buildFrame())

	_ = w.Close()
	staleMsg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	if strings.Contains(string(staleMsg), "AbstractMethodError") {
		t.Fatalf("expected the stale dispatchCache entry to mask the AbstractMethodError, "+
			"but it was correctly raised: %s", string(staleMsg))
	}

	// --- Now clear the memoized dispatch results and retry: the call must now resolve
	// StaleDispatchClassB.run()V correctly and raise AbstractMethodError.
	ResetDispatchCaches()

	normalStderr = os.Stderr
	r, w, _ = os.Pipe()
	os.Stderr = w

	interpret(buildFrame())

	_ = w.Close()
	freshMsg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	if !strings.Contains(string(freshMsg), "INVOKEVIRTUAL: J class method code is empty") {
		t.Errorf("Expected AbstractMethodError for StaleDispatchClassB.run()V after "+
			"ResetDispatchCaches, got: %s", string(freshMsg))
	}
}

// argSlots: verifies the memoized parameter-count helper for a couple of descriptors.
func TestArgSlots(t *testing.T) {
	globals.InitGlobals("test")

	if n := argSlots("()V"); n != 0 {
		t.Errorf("Expected 0 arg slots for ()V, got %d", n)
	}
	if n := argSlots("(II)V"); n != 2 {
		t.Errorf("Expected 2 arg slots for (II)V, got %d", n)
	}
	// second call should hit the cache and return the same result
	if n := argSlots("(II)V"); n != 2 {
		t.Errorf("Expected 2 arg slots for (II)V on cached call, got %d", n)
	}
}

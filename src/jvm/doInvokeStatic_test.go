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
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"os"
	"strings"
	"testing"
)

// helper: build a minimal CP with a single MethodRef or InterfaceMethodRef entry pointing at
// className.methodName.methodType.
func buildInvokeStaticCP(className, methodName, methodType string, isInterface bool) *classloader.CPool {
	CP := &classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	if isInterface {
		CP.CpIndex[1] = classloader.CpEntry{Type: classloader.Interface, Slot: 0}
	} else {
		CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
	}

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

func ensureInvokeStaticTestClass(className string, methods map[string]*classloader.Method) {
	clData := classloader.ClData{
		Name:            className,
		NameIndex:       stringPool.GetStringIndex(&className),
		SuperclassIndex: 0,
		MethodTable:     methods,
		CP:              classloader.CPool{},
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "bootstrap",
		Data:   &clData,
	}
	classloader.MethAreaInsert(className, &k)
}

func TestDoInvokeStatic_MethodNotFound(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	className := "jacobin/src/test/StaticTestObj"
	ensureInvokeStaticTestClass(className, make(map[string]*classloader.Method))

	f := newFrame(opcodes.INVOKESTATIC)
	f.Meth = append(f.Meth, 0x00, 0x01)
	f.CP = buildInvokeStaticCP(className, "missingMethod", "()V", false)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKESTATIC: Class method not found: "+className+".missingMethod()V") {
		t.Fatalf("Expected NoSuchMethodError with full FQN, got: %s", errMsg)
	}
}

func TestDoInvokeStatic_NotStaticMethod_ThrowsIncompatibleClassChangeError(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	className := "jacobin/src/test/StaticTestObj"
	methodName := "instanceMethod"
	methodType := "()V"

	methods := map[string]*classloader.Method{
		methodName + methodType: {
			AccessFlags: classloader.ACC_PUBLIC, // NOT static
			CodeAttr:    classloader.CodeAttrib{Code: []byte{byte(opcodes.RETURN)}},
		},
	}
	ensureInvokeStaticTestClass(className, methods)

	f := newFrame(opcodes.INVOKESTATIC)
	f.Meth = append(f.Meth, 0x00, 0x01)
	f.CP = buildInvokeStaticCP(className, methodName, methodType, false)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Method is not static") {
		t.Fatalf("Expected IncompatibleClassChangeError, got: %s", errMsg)
	}
}

func TestDoInvokeStatic_AbstractMethod_ThrowsAbstractMethodError(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	className := "jacobin/src/test/StaticTestObj"
	methodName := "abstractStaticMethod"
	methodType := "()V"

	methods := map[string]*classloader.Method{
		methodName + methodType: {
			AccessFlags: classloader.ACC_PUBLIC | classloader.ACC_STATIC | classloader.ACC_ABSTRACT,
		},
	}
	ensureInvokeStaticTestClass(className, methods)

	f := newFrame(opcodes.INVOKESTATIC)
	f.Meth = append(f.Meth, 0x00, 0x01)
	f.CP = buildInvokeStaticCP(className, methodName, methodType, false)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Abstract method requested") {
		t.Fatalf("Expected AbstractMethodError, got: %s", errMsg)
	}
}

func TestDoInvokeStatic_NativeMethod_ThrowsUnsatisfiedLinkError(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := classloader.Init()
	if err != nil {
		t.Fatalf("Failure to load classes: %s", err.Error())
	}
	jacobinSrc.CheckTestGfunctionsLoaded()

	className := "jacobin/src/test/StaticTestObj"
	methodName := "nativeStaticMethod"
	methodType := "()V"

	methods := map[string]*classloader.Method{
		methodName + methodType: {
			AccessFlags: classloader.ACC_PUBLIC | classloader.ACC_STATIC | classloader.ACC_NATIVE,
		},
	}
	ensureInvokeStaticTestClass(className, methods)

	f := newFrame(opcodes.INVOKESTATIC)
	f.Meth = append(f.Meth, 0x00, 0x01)
	f.CP = buildInvokeStaticCP(className, methodName, methodType, false)

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Native method requested") {
		t.Fatalf("Expected UnsatisfiedLinkError, got: %s", errMsg)
	}
}

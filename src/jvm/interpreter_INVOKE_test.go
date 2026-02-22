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
	"jacobin/src/gfunction"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
	"strings"
	"sync"
	"testing"
)

// This contains all the unit tests for the INVOKE family of bytecodes. They would normally
// appear in run_II-LD_test.go, but they would make that an enormous file. So, they're extracted here.

// INVOKEINTERFACE: CP entry does not point to an interface -> IncompatibleClassChangeError
func TestInvokeInterface_NotPointingToInterface(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEINTERFACE)
	// CP slot 1, count=1, zeroByte=0
	f.Meth = append(f.Meth, 0x00, 0x01, 0x01, 0x00)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	// WRONG: MethodRef instead of Interface
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
	f.CP = &CP

	// push a dummy object (won't be used because we fail earlier)
	push(&f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEINTERFACE: CP entry type") {
		t.Fatalf("Expected IncompatibleClassChangeError about CP type, got: %s", errMsg)
	}
}

// INVOKEINTERFACE: Non-zero zeroByte (5th byte) -> IncompatibleClassChangeError
func TestInvokeInterface_NonZeroZeroByte(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEINTERFACE)
	// CP slot 1, count=1, zeroByte=1 (invalid)
	f.Meth = append(f.Meth, 0x00, 0x01, 0x01, 0x01)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.Interface, Slot: 0}
	CP.InterfaceRefs = make([]classloader.InterfaceRefEntry, 1)
	// Minimal, will not be reached
	CP.InterfaceRefs[0] = classloader.InterfaceRefEntry{ClassIndex: 0, NameAndType: 0}
	f.CP = &CP

	push(&f, object.MakeEmptyObject())

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEINTERFACE: CP entry type") {
		t.Fatalf("Expected IncompatibleClassChangeError due to zero byte, got: %s", errMsg)
	}
}

// INVOKEINTERFACE: Null object reference -> NullPointerException
func TestInvokeInterface_NullObjectRef(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEINTERFACE)
	// CP slot 1, count=1 (no args; objRef at TOS), zeroByte=0
	f.Meth = append(f.Meth, 0x00, 0x01, 0x01, 0x00)

	// Build minimal CP so the path reaches objRef check and can compose method id in message
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.Interface, Slot: 0}
	CP.InterfaceRefs = make([]classloader.InterfaceRefEntry, 1)
	CP.InterfaceRefs[0] = classloader.InterfaceRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	iname := "pkg/MyIface"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&iname)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{NameIndex: 4, DescIndex: 5}

	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.Utf8Refs = make([]string, 6)
	CP.Utf8Refs[0] = "m"
	CP.Utf8Refs[1] = "()V"

	f.CP = &CP

	// Push a nil interface{} directly, so objRef == nil triggers
	push(&f, nil)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEINTERFACE: object whose method") ||
		!strings.Contains(errMsg, "is invoked is null") {
		t.Fatalf("Expected NullPointerException for null objRef, got: %s", errMsg)
	}
}

/*
// INVOKEINTERFACE: Object class cannot be loaded -> NoClassDefFoundError path (reported to tests)
func TestInvokeInterface_ObjectClassNotFound(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEINTERFACE)
	// CP slot 1, count=1, zeroByte=0
	f.Meth = append(f.Meth, 0x00, 0x01, 0x01, 0x00)

	// Minimal interface ref to get past early checks
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.Interface, Slot: 0}
	CP.InterfaceRefs = make([]classloader.InterfaceRefEntry, 1)
	CP.InterfaceRefs[0] = classloader.InterfaceRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	iname := "pkg/MyIface"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&iname)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{NameIndex: 4, DescIndex: 5}

	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.Utf8Refs = make([]string, 6)
	CP.Utf8Refs[0] = "m"
	CP.Utf8Refs[1] = "()V"
	f.CP = &CP

	// Push an object whose class name does not exist -> LoadClassFromNameOnly should fail
	bogusClass := "no/such/Class"
	obj := object.MakeEmptyObject()
	obj.KlassName = stringPool.GetStringIndex(&bogusClass)
	push(&f, obj)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	// The specific wording originates in runUtils (NoClassDefFoundError). We only assert that an
	// INVOKEINTERFACE-related error surfaced, since exact string may vary.
	if !strings.Contains(string(msg), "INVOKEINTERFACE:") {
		t.Fatalf("Expected an INVOKEINTERFACE error due to missing class, got: %s", string(msg))
	}
}

// INVOKEINTERFACE: Interface name resolves but is not an interface -> IncompatibleClassChangeError
func TestInvokeInterface_TargetNotAnInterface(t *testing.T) {
	globals.InitGlobals("test")

	// Initialize CP so that the interface name points to a normal class (not an interface).
	// We do not need to actually load real classes; the path fails during interface validation.

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEINTERFACE)
	// CP slot 1, count=1, zeroByte=0
	f.Meth = append(f.Meth, 0x00, 0x01, 0x01, 0x00)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.Interface, Slot: 0}

	CP.InterfaceRefs = make([]classloader.InterfaceRefEntry, 1)
	CP.InterfaceRefs[0] = classloader.InterfaceRefEntry{ClassIndex: 2, NameAndType: 3}

	// ClassRef -> points to a non-interface class name (any name is fine; validation will fail later)
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	notIface := "jacobin/src/test/Object" // known class in tests (not an interface)
	CP.ClassRefs[0] = stringPool.GetStringIndex(&notIface)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{NameIndex: 4, DescIndex: 5}

	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.Utf8Refs = make([]string, 6)
	CP.Utf8Refs[0] = "m"
	CP.Utf8Refs[1] = "()V"

	f.CP = &CP

	// Push a minimal object for objRef; class loading for objRef may fail in CI, so set it to a harmless name
	// The path we assert here fails during interface validation before needing the objRef class.
	someClass := "pkg/SomeClass"
	obj := object.MakeEmptyObject()
	obj.KlassName = stringPool.GetStringIndex(&someClass)
	push(&f, obj)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "INVOKEINTERFACE:") || !strings.Contains(errMsg, "not an interface") {
		t.Fatalf("Expected IncompatibleClassChangeError for non-interface, got: %s", errMsg)
	}
}

// INVOKESPECIAL should do nothing and report no errors
func TestNewInvokeSpecialJavaLangObject(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialJavaLangObject")
	}
	classloader.LoadBaseClasses() // must follow classloader.Init()

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	CP.MethodRefs[0] = classloader.MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	CP.ClassRefs[0] = types.StringPoolObjectIndex

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "<init>"

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "()V"

	f.CP = &CP
	classloader.ResolveCPmethRefs(&CP)
	classname := "java/lang/Object"
	push(&f, object.MakeEmptyObjectWithClassName(&classname))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame

	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("INVOKESPECIAL: Got unexpected error: %s", errMsg)
	}

	if f.TOS != 0 {
		t.Errorf("INVOKESPECIAL: Expected TOS after return to be 0, got %d", f.TOS)
	}
}
*/

// INVOKESPECIAL: verify that a call to a gmethod works correctly (passing in nothing, getting a link back)
func TestNewInvokeSpecialGmethodNoParams(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialGmethodNoParams")
	}

	gfunction.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	CP.MethodRefs[0] = classloader.MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	classname := "jacobin/src/test/Object"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&classname)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "test"

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "()Ljava/lang/Object;"

	f.CP = &CP
	classloader.ResolveCPmethRefs(&CP)
	obj := object.MakeEmptyObject()
	push(&f, obj) // INVOKESPECIAL expects a pointer to an object on the op stack

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("INVOKESPECIAL: Got unexpected error: %s", errMsg)
	}

	if f.TOS != 0 { // it's 0 b/c the gfunction returns a value, that is pushed onto the op stack
		t.Errorf("Expecting TOS to be 0, got %d", f.TOS)
	}
}

// INVOKESPECIAL: verify call to a gmethod works correctly and pushes the returned D twice
func TestNewInvokeSpecialGmethodNoParamsReturnsD(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialGmethodReturnsD")
	}

	gfunction.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	CP.MethodRefs[0] = classloader.MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	classname := "jacobin/src/test/Object"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&classname)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "test"

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "()D"

	f.CP = &CP
	classloader.ResolveCPmethRefs(&CP)
	obj := object.MakeEmptyObject()
	push(&f, obj) // INVOKESPECIAL expects a pointer to an object on the op stack

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("INVOKESPECIAL: Got unexpected error: %s", errMsg)
	}

	if f.TOS != 0 {
		t.Errorf("Expecting TOS to be 0, got %d", f.TOS)
	}
}

// INVOKESPECIAL: Test proper operation of a method that reports an error
func TestNewInvokeSpecialGmethodErrorReturn(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialGmethodErrorReturn")
	}

	gfunction.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKESPECIAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	CP.MethodRefs[0] = classloader.MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	classname := "jacobin/src/test/Object"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&classname)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "test"

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "(D)E"

	f.CP = &CP
	classloader.ResolveCPmethRefs(&CP)
	obj := object.MakeEmptyObject()
	push(&f, obj)        // INVOKESPECIAL expects a pointer to an object on the op stack
	push(&f, int64(999)) // push the one param

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("INVOKESPECIAL: Expected an error returned, got none")
	} else {
		if !strings.Contains(errMsg, "intended return of test error") {
			t.Errorf("INVOKESPECIAL: Expected error message re 'intended return of test error', got: %s", errMsg)
		}
	}
}

// INVOKESTATIC: verify that a call to a gmethod works correctly (passing in nothing, getting a link back)
func TestNewInvokeStaticGmethodNoParams(t *testing.T) {
	globals.InitGlobals("test")
	className := "jacobin/src/test/Object"

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeStaticGmethodNoParams")
	}

	gfunction.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKESTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	CP.MethodRefs[0] = classloader.MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	classname := "jacobin/src/test/Object"
	CP.ClassRefs[0] = stringPool.GetStringIndex(&classname)

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "test"

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "()Ljava/lang/Object;"

	f.CP = &CP
	classloader.ResolveCPmethRefs(&CP)

	// INVOKESTATIC needs a parsed/loaded object in the MethArea to function
	clData := classloader.ClData{
		Name:      className,
		NameIndex: CP.ClassRefs[0],
		// Superclass:      "java/lang/Object",
		SuperclassIndex: types.StringPoolObjectIndex,
		Module:          "",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		// Methods:         nil,
		Attributes: nil,
		SourceFile: "",
		CP:         classloader.CPool{},
		Access:     classloader.AccessFlags{},
		ClInit:     types.ClInitRun,
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "boostrap",
		Data:   &clData,
	}

	classloader.MethAreaInsert(className, &k)
	jlc := classloader.Jlc{
		Lock:     sync.RWMutex{},
		Statics:  []string{},
		KlassPtr: nil,
	}
	globals.JLCmap[className] = &jlc

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("INVOKESTATIC: Got unexpected error: %s", errMsg)
	}

	if f.TOS != 0 {
		t.Errorf("INVOKESTATIC: Expecting TOS to be 0, got %d", f.TOS)
	}
}

// INVOKESTATIC: verify that a call to a gmethod works correctly (passing in nothing, getting a link back)
func TestNewInvokeStaticGmethodErrorReturn(t *testing.T) {
	globals.InitGlobals("test")
	className := "jacobin/src/test/Object"

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeStaticGmethodNoParams")
	}

	gfunction.CheckTestGfunctionsLoaded()

	f := newFrame(opcodes.INVOKESTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
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
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "test"

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "(D)E"

	f.CP = &CP
	classloader.ResolveCPmethRefs(&CP)

	push(&f, int64(999)) // push the one param

	// INVOKESTATIC needs a parsed/loaded object in the MethArea to function
	clData := classloader.ClData{
		Name:      "jacobin/src/test/Object",
		NameIndex: CP.ClassRefs[0],
		// Superclass:      "java/lang/Object",
		SuperclassIndex: types.StringPoolObjectIndex,
		Module:          "",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Attributes:      nil,
		SourceFile:      "",
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          types.ClInitRun,
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "boostrap",
		Data:   &clData,
	}

	classloader.MethAreaInsert(className, &k)
	jlc := classloader.Jlc{
		Lock:     sync.RWMutex{},
		Statics:  []string{},
		KlassPtr: nil,
	}
	globals.JLCmap[className] = &jlc

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("INVOKESTATIC: Expected an error returned, got none")
	} else {
		if !strings.Contains(errMsg, "intended return of test error") {
			t.Errorf("INVOKESTATIC: Expected error message re 'intended return of test error', got: %s", errMsg)
		}
	}
}

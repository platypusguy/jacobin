/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/stringPool"
	"jacobin/types"
	"os"
	"strings"
	"testing"
)

// This contains all the unit tests for the INVOKE family of bytecodes. They would normally
// appear in run_II-LD_test.go, but they would make that an enormous file. So, they're extracted here.

// INVOKEINTERFACE: Invalid passed parameter
func TestInvokeInterfaceInvalid(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialJavaLangObject")
	}
	classloader.LoadBaseClasses() // must follow classloader.Init()

	f := newFrame(opcodes.INVOKEINTERFACE)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	f.Meth = append(f.Meth, 0x00) // the param count (which cannot be zero--this causes the error)
	f.Meth = append(f.Meth, 0x00)

	// create a dummy CP with 2 entries so that the CP slot index above does not cause an error.
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err == nil {
		t.Errorf("INVOKEINTERFACE: Should have returned an error for non-existent method, but didn't")
	} else {
		if !strings.Contains(err.Error(), "Invalid values for INVOKEINTERFACE bytecode") {
			t.Errorf("INVOKEINTERFACE: Got unexpected error message: %s", err.Error())
		}
	}
	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr

}

// INVOKESPECIAL should do nothing and report no errors
func TestInvokeSpecialJavaLangObject(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
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
	CP.ClassRefs[0] = types.ObjectPoolStringIndex

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
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err != nil {
		t.Errorf("INVOKESPECIAL: Got unexpected error: %s", err.Error())
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKESPECIAL of non-existent class. (throws exception in real code; returns error in tests)
func TestInvokeSpecialNonExistentMethod(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialNonExistentMethod")
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
	CP.ClassRefs[0] = types.ObjectPoolStringIndex

	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "no-such-method" // the non-existent method

	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "()V"

	f.CP = &CP
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err == nil {
		t.Errorf("INVOKESPECIAL: Should have returned an error for non-existent method, but didn't")
	} else {
		if !strings.Contains(err.Error(),
			"INVOKESPECIAL: Class method not found: java/lang/Object.no-such-method()V") {
			t.Errorf("INVOKESPECIAL: Got unexpected error: %s", err.Error())
		}
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKESPECIAL: verify that a call to a gmethod works correctly (passing in nothing, getting a link back)
func TestInvokeSpecialGmethodNoParams(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
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
	classname := "jacobin/test/Object"
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
	obj := object.MakeEmptyObject()
	push(&f, obj) // INVOKESPECIAL expects a pointer to an object on the op stack

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err != nil {
		t.Errorf("INVOKESPECIAL: Got unexpected error: %s", err.Error())
	}

	if f.TOS != 0 {
		t.Errorf("Expecting TOS to be 0, got %d", f.TOS)
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKESPECIAL: verify call to a gmethod works correctly and pushes the returned D twice
func TestInvokeSpecialGmethodNoParamsReturnsD(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
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
	classname := "jacobin/test/Object"
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
	obj := object.MakeEmptyObject()
	push(&f, obj) // INVOKESPECIAL expects a pointer to an object on the op stack

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err != nil {
		t.Errorf("INVOKESPECIAL: Got unexpected error: %s", err.Error())
	}

	if f.TOS != 1 { // should be 1 b/c a returned D occupies two slots on the op stack
		t.Errorf("Expecting TOS to be 1, got %d", f.TOS)
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKESPECIAL: Test proper operation of a method that reports an error
func TestInvokeSpecialGmethodErrorReturn(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
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
	classname := "jacobin/test/Object"
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
	obj := object.MakeEmptyObject()
	push(&f, obj)        // INVOKESPECIAL expects a pointer to an object on the op stack
	push(&f, int64(999)) // push the one param

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err == nil {
		t.Errorf("INVOKESPECIAL: Expected an error returned, got none")
	} else {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "intended return of test error") {
			t.Errorf("INVOKESPECIAL: Expected error message re 'intended return of test error', got: %s", errMsg)
		}
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKESTATIC: verify that a call to a gmethod works correctly (passing in nothing, getting a link back)
func TestInvokeStaticGmethodNoParams(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
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
	classname := "jacobin/test/Object"
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

	// INVOKESTATIC needs a parsed/loaded object in the MethArea to function
	clData := classloader.ClData{
		Name:            "jacobin/test/Object",
		NameIndex:       CP.ClassRefs[0],
		Superclass:      "java/lang/Object",
		SuperclassIndex: 0,
		Module:          "",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Methods:         nil,
		Attributes:      nil,
		SourceFile:      "",
		Bootstraps:      nil,
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          types.ClInitRun,
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "boostrap",
		Data:   &clData,
	}

	classloader.MethAreaInsert("jacobin/test/Object", &k)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err != nil {
		t.Errorf("INVOKESTATIC: Got unexpected error: %s", err.Error())
	}

	if f.TOS != 0 {
		t.Errorf("INVOKESTATIC: Expecting TOS to be 0, got %d", f.TOS)
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKESTATIC: verify that a call to a gmethod works correctly (passing in nothing, getting a link back)
func TestInvokeStaticGmethodErrorReturn(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
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
	classname := "jacobin/test/Object"
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

	push(&f, int64(999)) // push the one param

	// INVOKESTATIC needs a parsed/loaded object in the MethArea to function
	clData := classloader.ClData{
		Name:            "jacobin/test/Object",
		NameIndex:       CP.ClassRefs[0],
		Superclass:      "java/lang/Object",
		SuperclassIndex: 0,
		Module:          "",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Methods:         nil,
		Attributes:      nil,
		SourceFile:      "",
		Bootstraps:      nil,
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          types.ClInitRun,
	}
	k := classloader.Klass{
		Status: 'X',
		Loader: "boostrap",
		Data:   &clData,
	}

	classloader.MethAreaInsert("jacobin/test/Object", &k)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err = runFrame(fs)

	if err == nil {
		t.Errorf("INVOKESTATIC: Expected an error returned, got none")
	} else {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "intended return of test error") {
			t.Errorf("INVOKESTATIC: Expected error message re 'intended return of test error', got: %s", errMsg)
		}
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

// INVOKEVIRTUAL : invoke method -- here testing for error
func TestInvokevirtualInvalid(t *testing.T) {

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0} // should be a method ref
	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.FieldRefEntry, 1)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 0, NameAndType: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	if err == nil {
		t.Errorf("INVOKEVIRTUAL: Expected error but did not get one.")
	} else {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "Expected a method ref, but got") {
			t.Errorf("INVOKEVIRTUAL: Did not get expected error message, got: %s", errMsg)
		}
	}

	// restore stderr
	_ = w.Close()
	os.Stderr = normalStderr
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"math"
	"os"
	"strings"
	"testing"
)

// This file contains unit tests for the array bytecodes. Array operation primitives
// are tested in object.arrays_test.go

func TestNewNewJdkArrayTypeToJacobinType(t *testing.T) {

	a := object.JdkArrayTypeToJacobinType(object.T_BOOLEAN)
	if a != object.BYTE {
		t.Errorf("Expected Jacobin type of %d, got: %d", object.BYTE, a)
	}

	b := object.JdkArrayTypeToJacobinType(object.T_CHAR)
	if b != object.INT {
		t.Errorf("Expected Jacobin type of %d, got: %d", object.INT, b)
	}

	c := object.JdkArrayTypeToJacobinType(object.T_DOUBLE)
	if c != object.FLOAT {
		t.Errorf("Expected Jacobin type of %d, got: %d", object.FLOAT, c)
	}

	d := object.JdkArrayTypeToJacobinType(999)
	if d != object.ERROR {
		t.Errorf("Expected Jacobin type of %d, got: %d", object.ERROR, d)
	}
}

// AALOAD: Test fetching and pushing the value of an element in a reference array
// The logic here is effectively identical to IALOAD. This code also tests AASTORE.
func TestNewAaload(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.ANEWARRAY)
	push(&f, int64(30)) // make an array of 30 elements
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // use the classRef at CP[1] as the type of reference

	// now create the CP.
	// CP[0] is perforce 0
	// CP[1] is a ClassRef that points to string pool entry for the class name
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, types.StringPoolStringIndex) // use string pool
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.AASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	oPtr := object.MakeEmptyObject()
	push(&f, oPtr) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.AALOAD) // now fetch the value in array[20]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	os.Stderr = normalStderr

	res := pop(&f)
	if res != oPtr {
		t.Errorf("AALOAD: Expected loaded array value = %v, got: %v", oPtr, res)
	}

	if f.TOS != -1 {
		t.Errorf("AALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// AALOAD: Test with a nil
func TestNewAaloadWithNil(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fs := frames.CreateFrameStack()

	f := newFrame(opcodes.AALOAD)
	push(&f, nil)       // push the reference to the array -- here nil
	push(&f, int64(20)) // index to array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("AALOAD: Expecting error for nil refernce, but got none")
	}

	if !strings.Contains(errMsg, "Invalid (null) reference") {
		t.Errorf("AALOAD: Did not get expected error msg, got: %s", errMsg)
	}
}

// AALOAD: using an invalid subscript into the array
func TestAaloadInvalidSubscript(t *testing.T) {
	globals.InitGlobals("test")

	refArr := object.Make1DimRefArray(types.ObjectClassName, 10)

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.AALOAD) // now fetch the value
	push(&f, refArr)              // push the reference to the array
	push(&f, int64(200))          // get contents in array[200] which is invalid
	ret := doAaload(&f, 0)

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	if ret != math.MaxInt32 { // = exceptions.ERROR_OCCURRED. Literal not used due to circularity.
		t.Errorf("AALOAD: Expecting error code, got: %d", ret)
	}

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("DALOAD: Did not get expected err msg for invalid subscript, got: %s",
			errMsg)
	}
}

// ANEWARRAY: create an array of T_REF.
// AASTORE: store a value in the array.
//
// Create an array of 30 String elements and store ptr value in array[20].
func TestNewAastore(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.ANEWARRAY)
	push(&f, int64(30)) // make an array of 30 elements
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // use the classRef at CP[2] as the type of reference

	// now create the CP.
	// CP[0] is perforce 0
	// [1] is a ClassRef that points to a string pool entry
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, types.StringPoolStringIndex) // point to string pool
	f.CP = &CP

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("TestAastore: Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("TestAastore: Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	statics.LoadStaticsString()

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)
	f = newFrame(opcodes.AASTORE)

	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // index to array[20]
	// objRef := object.NewStringFromGoString("test")
	objRef := object.StringObjectFromGoString("test")
	push(&f, objRef) // store the address of a string

	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// now retrieve the updated element
	array := ptr.FieldTable["value"].Fvalue.([]*object.Object)
	udpatedElement := array[20]
	if udpatedElement != objRef { // check that the element is actually updated
		t.Errorf("TestAastore: Expected array[20]=test, observed: %v", udpatedElement)
	}
}

// AASTORE: Test error conditions: invalid array address
func TestNewAastoreInvalid1(t *testing.T) {
	globals.InitStringPool()
	f := newFrame(opcodes.AASTORE)
	obj := object.Make1DimArray(object.REF, 10)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, obj)                   // the value to insert

	globals.InitGlobals("test")
	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null)") {
		t.Errorf("AASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// AASTORE: Test error conditions: wrong type of array (not [I)
func TestNewAastoreInvalid2(t *testing.T) {

	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(opcodes.AASTORE)
	push(&f, o)        // this should point to an array of refs, not ints, will here cause the error
	push(&f, int64(5)) // the index into the array
	push(&f, o)        // the value to insert

	trace.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "field type must start with '[L',") {
		t.Errorf("AASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// AASTORE: Test error conditions: index out of range
func TestNewAastoreInvalid3(t *testing.T) {
	objType := types.ObjectClassName
	o := object.Make1DimRefArray(objType, 10)
	f := newFrame(opcodes.AASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, o)         // the value to insert

	trace.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "but array index is") {
		t.Errorf("AASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// ANEWARRAY: creation of array for references to strings
func TestNewAnewrray(t *testing.T) {
	f := newFrame(opcodes.ANEWARRAY)
	push(&f, int64(13)) // make an array of 13 elements
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // use the classRef at CP[1] as the type of reference

	// now create the CP.
	// [0] is First entry is perforce 0
	// [1] is a ClassRef that points to a string pool entry for "java/lang/String"
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, types.StringPoolStringIndex) // point to string pool entry
	f.CP = &CP

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, test the length of the array, which should be 13
	element := g.ArrayAddressList.Front()
	ptr := element.Value.(*object.Object)
	o := ptr.FieldTable["value"]
	array := o.Fvalue.([]*object.Object)
	if len(array) != 13 {
		t.Errorf("ANEWARRAY: Expecting array length of 13, got %d", len(array))
	}
}

// ANEWARRAY: creation of array for references; test contents of Klass field
func TestNewAnewrrayKlassField(t *testing.T) {
	f := newFrame(opcodes.ANEWARRAY)
	push(&f, int64(13)) // make an array of 13 elements
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x02) // use the classRef at CP[2] as the type of reference

	// now create the CP. First entry is perforce 0
	// [1] entry points to a UTF8 entry with the class name (should be "java/lang/String")
	// [2] is a ClassRef that points to string pool entry for same
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, types.StringPoolStringIndex) // point to string pool entry
	CP.Utf8Refs = append(CP.Utf8Refs, types.StringClassName)
	f.CP = &CP

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, test the length of the array, which should be 13
	element := g.ArrayAddressList.Front()
	ptr := element.Value.(*object.Object)
	klassString := stringPool.GetStringPointer(ptr.KlassName)
	if !strings.HasPrefix(*klassString, types.RefArray) {
		t.Errorf("ANEWARRAY: Expecting class to start with '[L', got %s", *klassString)
	}

	if !strings.HasSuffix(*klassString, types.StringClassName) {
		t.Errorf("ANEWARRAY: Expecting class to end with 'java/lang/String', got %s", *klassString)
	}
}

// ANEWARRAY: creation of array for references; test invalid array size
func TestNewAnewrrayInvalidSize(t *testing.T) {
	f := newFrame(opcodes.ANEWARRAY)
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	push(&f, int64(-1)) // make the array an invalid size

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg == "" {
		t.Errorf("ANEWARRAY: Did not get expected error")
	}

	if !strings.Contains(errMsg, "java.lang.NegativeArraySizeException") {
		t.Errorf("ANEWARRAY: Expecting different error msg, got %s", msg)
	}
}

// ARRAYLENGTH: Test length of byte array
// First, we create the array of 13 elements, then we push the reference
// to it and execute the ARRAYLENGTH bytecode using the address stored
// in the global array address list
func TestNewByteArrayLength(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(13))                    // make the array 13 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(opcodes.ARRAYLENGTH)
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	size := pop(&f).(int64)
	if size != 13 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 13, got %d", size)
	}
}

// ARRAYLENGTH: Test length of int array
func TestNewIntArrayLength(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(22))                   // make the array 22 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(opcodes.ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	size := pop(&f).(int64)
	if size != 22 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 13, got %d", size)
	}
}

// ARRAYLENGTH: Test length of float array
func TestNewFloatArrayLength(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(34))                      // make the array 34 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(opcodes.ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 34, got %d", size)
	}
}

// ARRAYLENGTH: Test length of array of longs
func TestNewLongArrayLength(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(34))                    // make the array 34 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(opcodes.ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 34, got %d", size)
	}
}

// ARRAYLENGTH: Test length of array of references
func TestNewRefArrayLength(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(34))                   // make the array 34 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of references

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(opcodes.ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 34, got %d", size)
	}
}

// ARRAYLENGTH: Test length of raw byte array
func TestNewRawByteArrayLength(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	array := []byte{'a', 'b', 'c'}
	f := newFrame(opcodes.ARRAYLENGTH)
	push(&f, &array) // push the reference to the array
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("ARRAYLENGTH: Got unexpected error message: %s", errMsg)
	}

	length := pop(&f).(int64)
	if length != 3 {
		t.Errorf("ARRAYLENGTH: Expecting length of 3, got: %d", length)
	}
}

// ARRAYLENGTH: Test length of raw int8 array
func TestNewRawInt8ArrayLength(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	array := []uint8{'a', 'b', 'c'}
	f := newFrame(opcodes.ARRAYLENGTH)
	push(&f, &array) // push the reference to the array
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("TestRawInt8ArrayLength: Got unexpected error message: %s", errMsg)
	}

	length := pop(&f).(int64)
	if length != 3 {
		t.Errorf("TestRawInt8ArrayLength: Expecting length of 3, got: %d", length)
	}
}

// ARRAYLENGTH: Test length of nil array -- should return an error
func TestNewNilArrayLength(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.ARRAYLENGTH)
	push(&f, nil) // push the reference to the array
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("ARRAYLENGTH: Expecting an error message, but got none")
	}

	if !strings.Contains(errMsg, "ARRAYLENGTH: Invalid (null) reference to an array") {
		t.Errorf("ARRAYLENGTH: Expecting different error msg, got: %s", errMsg)
	}
}

// BALOAD: Test fetching and pushing the value of an element in a byte/boolean array
// The logic here is effectively identical to IALOAD. This code also tests BASTORE.
func TestNewBaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f)

	f = newFrame(opcodes.BASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, byte(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.BALOAD) // now fetch the value in array[20]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("BALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("BALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// BALOAD: Test exception on nil array address
func TestNewBaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.BALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("BALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// BALOAD: using an invalid subscript into the array
func TestNewBaloadInvalidSubscript(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	trace.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.BALOAD) // now fetch the value
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(200))         // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("BALOAD: Did not get expected err msg for invalid subscript, got: %s",
			errMsg)
	}
}

// BASTORE: store value in array of bytes
// Create an array of 30 elements, store value 100 in array[20], then
// sum all the elements in the array, and test for a sum of 100.
// Note the value we store must be an int64 value--not a byte
func TestNewBastore(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.BASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, byte(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	o := ptr.FieldTable["value"]
	array := o.Fvalue.([]types.JavaByte) // get the array
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(array[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// BASTORE: Tests whether storing an int64 into a byte array does the right thing
func TestNewBastoreInt64(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.BASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	o := ptr.FieldTable["value"]
	array := o.Fvalue.([]types.JavaByte) // get the array
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(array[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// BASTORE: Test error conditions: invalid array address
func TestNewBastoreInvalid1(t *testing.T) {
	f := newFrame(opcodes.BASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, int64(20))             // the value to insert

	trace.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("BASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// BASTORE: Test error conditions: wrong type of array (not [I)
func TestNewBastoreInvalid2(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.BASTORE)
	o := object.Make1DimArray(object.FLOAT, 10)
	push(&f, o)         // this should point to an array of ints, not floats, will here cause the error
	push(&f, int64(30)) // the index into the array
	push(&f, int64(20)) // the value to insert

	trace.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "field type expected=[B") {
		t.Errorf("BASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// BASTORE: Test error conditions: index out of range
func TestNewBastoreInvalid3(t *testing.T) {

	globals.InitGlobals("test")
	o := object.Make1DimArray(object.BYTE, 10)
	f := newFrame(opcodes.BASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, int64(20)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "but array index is") {
		t.Errorf("BASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// BASTORE: Test storing into a direct []types.JavaByte array (not wrapped in object.Object)
func TestNewBastoreJavaByteArray(t *testing.T) {
	globals.InitGlobals("test")

	// Create a direct []types.JavaByte array
	javaByteArray := make([]types.JavaByte, 30)

	f := newFrame(opcodes.BASTORE)
	push(&f, javaByteArray) // push the direct JavaByte array reference
	push(&f, int64(20))     // in array[20]
	push(&f, int64(100))    // the value we're storing

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// Verify the value was stored correctly
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(javaByteArray[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}

	// Verify the value at index 20 is specifically 100
	if javaByteArray[20] != 100 {
		t.Errorf("BASTORE: Expected javaByteArray[20] to be 100, got: %d", javaByteArray[20])
	}
}

// BASTORE: Test error conditions: unexpected reference type (triggers default case)
func TestNewBastoreInvalidType(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.BASTORE)

	// Push an unexpected type (e.g., a string instead of an array)
	push(&f, "not an array") // this will trigger the default case
	push(&f, int64(20))      // the index into the array
	push(&f, int64(100))     // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "unexpected reference type") {
		t.Errorf("BASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// CALOAD: Test fetching and pushing the value of an element in an char array
// Chars in Java are two bytes; we accord each one an int64 element. As a result,
// the logic here is effectively identical to IALOAD. This code also tests CASTORE.
func TestNewCaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_CHAR) // make it an array of chars

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.CASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.CALOAD) // now fetch the value in array[20]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("CALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("CALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// DALOAD: Test fetching and pushing the value of an element in an float array
func TestNewDaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                      // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.DASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.DALOAD) // now fetch the value in array[30]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	res := pop(&f).(float64)
	if res != 100.0 {
		t.Errorf("FALOAD: Expected loaded array value of 100, got: %e", res)
	}

	if f.TOS != -1 {
		t.Errorf("DALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// DALOAD: Test exception on nil array address
func TestNewDaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.DALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid object pointer") {
		t.Errorf("DALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// DALOAD: using an invalid subscript into the array
func TestNewLaDoadInvalidSubscript(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                      // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	trace.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.DALOAD) // now fetch the value
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(200))         // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("DALOAD: Did not get expected err msg for invalid subscript, got: %s",
			errMsg)
	}
}

// DASTORE: store value in array of doubles
// See comments for IASTORE for the logic of this test
func TestNewDastore(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                      // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.DASTORE)
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // in array[20]
	push(&f, 100_000_000_000.25) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}

	oa := ptr.FieldTable["value"]
	array := oa.Fvalue.([]float64)
	var sum float64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100_000_000_000.25 {
		t.Errorf("DASTORE: Expected sum of doubles array to be 100,000,000,000.25, got: %f",
			sum)
	}
}

// DASTORE: Test error conditions: invalid array address
func TestNewDastoreInvalid1(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.DASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, float64(20.0))         // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("DASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// DASTORE: Test error conditions: wrong type of array (not [I)
func TestNewDastoreInvalid2(t *testing.T) {
	globals.InitGlobals("test")
	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(opcodes.DASTORE)
	push(&f, o)             // this should point to an array of floats, not ints, will here cause the error
	push(&f, int64(30))     // the index into the array
	push(&f, float64(20.0)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "field type expected") {
		t.Errorf("DASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// DASTORE: Test error conditions: index out of range
func TestNewDastoreInvalid3(t *testing.T) {

	globals.InitGlobals("test")
	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(opcodes.DASTORE)
	push(&f, o)             // an array of 10 ints, not floats
	push(&f, int64(30))     // the index into the array: it's too big, causing error
	push(&f, float64(20.0)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, " but array index is") {
		t.Errorf("DASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// FALOAD: Test fetching and pushing the value of an element in an float array
func TestNewFaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                     // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_FLOAT) // make it an array of floats

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.FASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.FALOAD) // now fetch the value in array[30]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	res := pop(&f).(float64)
	if res != 100.0 {
		t.Errorf("FALOAD: Expected loaded array value of 100, got: %e", res)
	}

	if f.TOS != -1 {
		t.Errorf("FALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// FALOAD: Test with raw float array
func TestFaLoadWithRawFloatArray(t *testing.T) {
	f := newFrame(opcodes.FALOAD)
	fArray := []float64{1.0, 2.0, 3.0, 4.0, 50.}
	push(&f, fArray)   // push the reference to the array, here a raw byte array
	push(&f, int64(2)) // get contents in array[2]

	globals.InitGlobals("test")

	// execute the bytecode
	ret := doFaload(&f, 0)

	if ret != 1 {
		t.Errorf("FALOAD: Expected error return of 1, got %d", ret)
	}

	fl := pop(&f).(float64)
	if fl != 3.0 {
		t.Errorf("FALOAD: Expected 3.0, got %f", fl)
	}
}

// FALOAD: Test exception on nil array address
func TestNewFaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.FALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid object pointer") {
		t.Errorf("FALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// FALOAD: using an invalid subscript into the array
func TestNewFaloadInvalidSubscript(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                     // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_FLOAT) // make it an array of floats

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	trace.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.FALOAD) // now fetch the value
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(200))         // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("DALOAD: Did not get expected err msg for invalid subscript, got: %s",
			errMsg)
	}
}

func TestFaLoadWhenNotAValidArray(t *testing.T) {
	f := newFrame(opcodes.FALOAD)
	badArray := []byte{1, 2, 3, 4, 5}
	push(&f, badArray)  // push the reference to the array, here a raw byte array
	push(&f, int64(20)) // get contents in array[20]

	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// execute the bytecode
	ret := doFaload(&f, 0)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	if ret != math.MaxInt32 { // = exceptions.ERROR_OCCURRED but can't use here due to circular import
		t.Errorf("FALOAD: Expected error return of , got %d", ret)
	}

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Reference invalid type of array") {
		t.Errorf("FALOAD: Got unexpected error message: %s", errMsg)
	}
}

// FASTORE: store value in array of floats
// Create an array of 30 elements, store value 100.0 in array[20], then
// sum all the elements in the array, and test for a sum of 100.0
func TestNewFastore(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                     // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_FLOAT) // make it an array of floats

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.FASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	oa := ptr.FieldTable["value"]
	array := oa.Fvalue.([]float64)
	var fsum float64
	for i := 0; i < 30; i++ {
		fsum += array[i]
	}
	if fsum != 100.0 {
		t.Errorf("FASTORE: Expected sum of array entries to be 100, got: %e", fsum)
	}
}

// FASTORE: Test error conditions: invalid array address
func TestNewFastoreInvalid1(t *testing.T) {
	f := newFrame(opcodes.FASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, float64(20.0))         // the value to insert

	trace.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("FASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// FASTORE: Test error conditions: wrong type of array (not [I)
func TestNewFastoreInvalid2(t *testing.T) {
	globals.InitGlobals("test")
	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(opcodes.FASTORE)
	push(&f, o)             // this should point to an array of floats, not ints, will here cause the error
	push(&f, int64(30))     // the index into the array
	push(&f, float64(20.0)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "field type expected") {
		t.Errorf("FASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// FASTORE: Test error conditions: index out of range
func TestNewFastoreInvalid3(t *testing.T) {
	globals.InitGlobals("test")
	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(opcodes.FASTORE)
	push(&f, o)             // an array of 10 ints, not floats
	push(&f, int64(30))     // the index into the array: it's too big, causing error
	push(&f, float64(20.0)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, " but array index is") {
		t.Errorf("FASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// IALOAD: Test fetching and pushing the value of an element in an int array
func TestNewIaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.IALOAD) // now fetch the value in array[20]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("IALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("IALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// IALOAD: Test exception on nil array address
func TestNewIaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.IALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid null reference to an array") {
		t.Errorf("IALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// IALOAD: using an invalid subscript into the array
func TestNewIaloadInvalidSubscript(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	trace.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.IALOAD) // now fetch the value
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(200))         // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("IALOAD: Did not get expected err msg for invalid subscript, got: %s",
			errMsg)
	}
}

func TestIaLoadWhenNotAValidArray(t *testing.T) {
	f := newFrame(opcodes.IALOAD)
	badArray := []byte{1, 2, 3, 4, 5}
	push(&f, badArray)  // push the reference to the array, here a raw byte array
	push(&f, int64(20)) // get contents in array[20]

	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// execute the bytecode
	ret := doIaload(&f, 0)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	if ret != math.MaxInt32 { // = exceptions.ERROR_OCCURRED but can't use here due to circular import
		t.Errorf("IALOAD: Expected error return of , got %d", ret)
	}

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Invalid reference to an array") {
		t.Errorf("IALOAD: Got unexpected error message: %s", errMsg)
	}
}

// IASTORE: store value in array of ints
// Create an array of 30 elements, store value 100 in array[20], then
// sum all the elements in the array, and test for a sum of 100.
func TestNewIastore(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	ao := ptr.FieldTable["value"].Fvalue
	array := ao.([]int64)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("IASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// IASTORE: Test error conditions: invalid array address
func TestNewIastoreInvalid1(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.IASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, int64(20))             // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("IASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// IASTORE: Test error conditions: wrong type of array (not [I)
func TestNewIastoreInvalid2(t *testing.T) {

	globals.InitGlobals("test")
	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(opcodes.IASTORE)
	push(&f, o)         // this should point to an array of ints, not floats, will here cause the error
	push(&f, int64(30)) // the index into the array
	push(&f, int64(20)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "field type expected") {
		t.Errorf("IASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// IASTORE: Test error conditions: index out of range
func TestNewIastoreInvalid3(t *testing.T) {

	globals.InitGlobals("test")
	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(opcodes.IASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, int64(20)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, " but array index is") {
		t.Errorf("IASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// LALOAD: Test fetching and pushing the value of an element into a long array
func TestNewLaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.LASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.LALOAD) // now fetch the value in array[20]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	if f.TOS != 0 {
		t.Errorf("LALOAD: Top of stack, expected 0, got: %d", f.TOS)
	}

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("LALOAD: Expected loaded array value of 100, got: %d", res)
	}
}

// LALOAD: Test exception on nil array address
func TestNewLaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.LALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid null reference to an array") {
		t.Errorf("LALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// LALOAD: using an invalid subscript into the array
func TestNewLaloadInvalidSubscript(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	trace.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.LALOAD) // now fetch the value
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(200))         // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("LALOAD: Did not get expected err msg for invalid subscript, got: %s",
			errMsg)
	}
}

// LASTORE: store value in array of longs
// See comments for IASTORE for the logic of this test
func TestNewLastore(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.LASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}

	oa := ptr.FieldTable["value"]
	array := oa.Fvalue.([]int64)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("LASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// LASTORE: Test error conditions: invalid array address
func TestNewLastoreInvalid1(t *testing.T) {
	f := newFrame(opcodes.LASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, int64(20))             // the value to insert

	trace.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("LASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// LASTORE: Test error conditions: wrong type of array (not [I)
func TestNewLastoreInvalid2(t *testing.T) {

	globals.InitGlobals("test")
	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(opcodes.LASTORE)
	push(&f, o)         // this should point to an array of ints, not floats, will here cause the error
	push(&f, int64(30)) // the index into the array
	push(&f, int64(20)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "field type expected") {
		t.Errorf("LASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// LASTORE: Test error conditions: index out of range
func TestNewLastoreInvalid3(t *testing.T) {

	globals.InitGlobals("test")
	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(opcodes.LASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, int64(20)) // the value to insert

	trace.Init()
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, " but array index is") {
		t.Errorf("LASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// MULTIANEWARRAY: test creation of a two-dimensional byte array
func TestNew2DimArray1(t *testing.T) {
	globals.InitGlobals("test")
	arr, err := object.Make2DimArray(3, 4, object.BYTE)
	if err != nil {
		t.Error("Error creating 2-dimensional array")
	}

	o := arr.FieldTable["value"]
	arrLevelArrayPtr := o.Fvalue.([]*object.Object)
	if len(arrLevelArrayPtr) != 3 {
		t.Errorf("MULTIANEWARRAY: Expected length of pointer array of 3, got: %d",
			len(arrLevelArrayPtr))
	}

	oa := arrLevelArrayPtr[0].FieldTable["value"]
	leafLevelArrayPtr := (oa.Fvalue).([]types.JavaByte)
	arrLen := len(leafLevelArrayPtr)
	if arrLen != 4 {
		t.Errorf("MULTIANEWARRAY: Expected length of leaf array of 4got: %d", arrLen)
	}
}

// MULTIANEWARRAY: test creation of a two-dimensional byte array and its Klass field
func TestNew2DimArrayKlassField(t *testing.T) {
	globals.InitGlobals("test")
	arr, err := object.Make2DimArray(3, 4, object.BYTE)
	if err != nil {
		t.Error("Error creating 2-dimensional array")
	}

	if arr.KlassName == types.InvalidStringIndex {
		t.Errorf("Array Klass field was invalid")
		return
	}

	arrKlass := stringPool.GetStringPointer(arr.KlassName)
	if *arrKlass != "[B" {
		t.Errorf("Expecting array with Klass of '[B', got: %s", *arrKlass)
	}
}

// MULTINEWARRAY: Test a straightforward 3x3x4 array of int64's
func TestNew3DimArray1(t *testing.T) {
	g := globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.

	// create the constant pool we'll point to
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	arrayType := "[[[I"
	nameIndex := stringPool.GetStringIndex(&arrayType)
	CP.ClassRefs = append(CP.ClassRefs, nameIndex)
	CP.Utf8Refs = append(CP.Utf8Refs, "[[[I")

	// create the frame
	f := newFrame(opcodes.MULTIANEWARRAY)
	f.Meth = append(f.Meth, 0x00) // this byte and next form index into CP
	f.Meth = append(f.Meth, 0x02)
	f.Meth = append(f.Meth, 0x03) // the number of dimensions
	push(&f, int64(0x03))         // size of the three dimensions: 4x3x2
	push(&f, int64(0x03))
	push(&f, int64(0x04))
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode
	if f.TOS != 0 {
		t.Errorf("MULTIANEWARRAY: Top of stack, expected 0, got: %d", f.TOS)
	}

	arrayPtr := pop(&f)
	if arrayPtr == nil {
		t.Error("MULTIANEWARRAY: Expected a pointer to an array, got nil")
	}

	topLevelArray := *(arrayPtr.(*object.Object))
	if topLevelArray.FieldTable["value"].Ftype != "[L" {
		t.Errorf("MULTIANEWARRAY: Expected 1st dim to be type '[L', got %s",
			topLevelArray.FieldTable["value"].Ftype)
	}

	dim1 := topLevelArray.FieldTable["value"].Fvalue.([]*object.Object)
	if len(dim1) != 3 {
		t.Errorf("MULTINEWARRAY: Expected 1st dim to have 3 elements, got: %d",
			len(dim1))
	}

	dim2type := dim1[0].FieldTable["value"].Ftype
	if dim2type != "[[I" {
		t.Errorf("MULTIANEWARRAY: Expected 2nd dim to be type '[[I', got %s",
			dim2type)
	}

	dim2 := dim1[0].FieldTable["value"].Fvalue.([]*object.Object)
	if len(dim2) != 3 {
		t.Errorf("MULTINEWARRAY: Expected 2nd dim to have 3 elements, got: %d",
			len(dim2))
	}

	dim3type := dim2[0].FieldTable["value"].Ftype
	if dim3type != "[I" {
		t.Errorf("MULTIANEWARRAY: Expected leaf dim to be type '[I', got %s",
			dim3type)
	}

	dim3 := dim2[0].FieldTable["value"].Fvalue.([]int64)
	if len(dim3) != 4 {
		t.Errorf("MULTINEWARRAY: Expected leaf dim to have 4 elements, got: %d",
			len(dim3))
	}

	elementValue := dim3[2] // an element in the leaf array
	if elementValue != 0 {
		t.Errorf("Expected element value to be 0, got %d", elementValue)
	}
}

// MULTINEWARRAY: Test an array 4x0x3 array of int64's. The zero
// size of the second dimension should result in an single-dimension
// array of int64s
func TestNew3DimArray2(t *testing.T) {
	g := globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	// create the constant pool we'll point to
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	arrayType := "[[[I"
	nameIndex := stringPool.GetStringIndex(&arrayType)
	CP.ClassRefs = append(CP.ClassRefs, nameIndex)
	CP.Utf8Refs = append(CP.Utf8Refs, "[[[I")

	// create the frame
	f := newFrame(opcodes.MULTIANEWARRAY)
	f.Meth = append(f.Meth, 0x00) // this byte and next form index into CP
	f.Meth = append(f.Meth, 0x02)
	f.Meth = append(f.Meth, 0x03) // the number of dimensions
	push(&f, int64(0x03))         // size of the three dimensions: 4x0x3
	push(&f, int64(0x00))
	push(&f, int64(0x04))
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	_ = w.Close()
	os.Stderr = normalStderr

	if f.TOS != 0 {
		t.Errorf("MULTIANEWARRAY: Top of stack, expected 0, got: %d", f.TOS)
	}

	arrayPtr := pop(&f)
	if arrayPtr == nil {
		t.Error("MULTIANEWARRAY: Expected a pointer to an array, got nil")
	}

	topLevelArray := *(arrayPtr.(*object.Object))
	if topLevelArray.FieldTable["value"].Ftype != "[I" {
		t.Errorf("MULTIANEWARRAY: Expected 1st dim to be type '[I', got %s",
			topLevelArray.FieldTable["value"].Ftype)
	}

	dim1 := topLevelArray.FieldTable["value"].Fvalue.([]int64)
	if len(dim1) != 4 {
		t.Errorf("MULTINEWARRAY: Expected 1st dim to have 4 elements, got: %d",
			len(dim1))
	}
}

// NEWARRAY: creation of array for primitive values
func TestNewNewrray(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(13))                    // make the array 13 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("NEWARRAY: Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, test the length of the array, which should be 13
	element := g.ArrayAddressList.Front()
	ptr := element.Value.(*object.Object)
	array := ptr.FieldTable["value"].Fvalue.([]int64)
	if len(array) != 13 {
		t.Errorf("NEWARRAY: Expecting array length of 13, got %d", len(array))
	}
}

// NEWARRAY: Create new array of 13 bytes
func TestNewNewrrayForByteArray(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(13))                    // size
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("NEWARRAY: Got unexpected error: %s", errMsg)
	}

	arrayPtr := pop(&f).(*object.Object)
	array := arrayPtr.FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(array) != 13 {
		t.Errorf("NEWARRAY: Got unexpected array size: %d", len(array))
	}
}

// NEWARRAY: Create new array -- test with invalid size
func TestNewNewArrayInvalidSize(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(-13))                   // invalid size (less than 0)
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("NEWARRAY: Expected an error message, but got none")
	}

	if !strings.Contains(errMsg, "Invalid size for array") {
		t.Errorf("NEWARRAY: Got unexpected error message: %s", errMsg)
	}
}

// NEWARRAY: Create new array -- test with invalid type ERROR
func TestNewNewrrayInvalidTypeError(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(13))                   // size
	f.Meth = append(f.Meth, object.ERROR) // invalid type

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("TestNewrrayInvalidTypeError: Expected an error message, but got none")
	}

	if !strings.Contains(errMsg, "Invalid array type specified") {
		t.Errorf("TestNewrrayInvalidTypeError: Got unexpected error message: %s", errMsg)
	}
}

// NEWARRAY: Create new array -- test with invalid type T_REF
func TestNewNewrrayInvalidTypeRef(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(13))                   // size
	f.Meth = append(f.Meth, object.T_REF) // invalid type

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("TestNewrrayInvalidTypeRef: Expected an error message, but got none")
		return
	}

	if !strings.Contains(errMsg, "Invalid array type specified") {
		t.Errorf("TestNewrrayInvalidTypeRef: Got unexpected error message: %s", errMsg)
	}
}

// SALOAD: Test fetching and pushing the value of an element in a short array
func TestNewSaload(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.SASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	f = newFrame(opcodes.SALOAD) // now fetch the value in array[30]
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("SALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("SALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// SASTORE: store value in array of shorts
// See comments for IASTORE for the logic of this test
func TestNewSastore(t *testing.T) {
	f := newFrame(opcodes.NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(opcodes.SASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)    // execute the bytecode

	array := ptr.FieldTable["value"].Fvalue.([]int64)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("SASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

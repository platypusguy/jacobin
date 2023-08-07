/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/types"
	"os"
	"strings"
	"testing"
)

func TestJdkArrayTypeToJacobinType(t *testing.T) {

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
func TestAaload(t *testing.T) {
	f := newFrame(ANEWARRAY)
	push(&f, int64(30)) // make the array 30 elements big

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(AASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	oPtr := object.MakeEmptyObject()
	push(&f, oPtr) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(AALOAD) // now fetch the value in array[20]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f)
	if res != oPtr {
		t.Errorf("AALOAD: Expected loaded array value = %v, got: %v", oPtr, res)
	}

	if f.TOS != -1 {
		t.Errorf("AALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// AALOAD: Test with a nil
func TestAaloadWithNil(t *testing.T) {
	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()

	f := newFrame(AALOAD)
	push(&f, nil)       // push the reference to the array -- here nil
	push(&f, int64(20)) // index to array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f)    // push the new frame
	err := runFrame(fs) // execute the bytecode

	if err == nil {
		t.Errorf("AALOAD: Expecting error for nil refernce, but got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Invalid (null) reference") {
		t.Errorf("AALOAD: Did not get expected error msg, got: %s", errMsg)
	}
}

// AASTORE: store value in array of bytes
// Create an array of 30 elements, store ptr value in array[20],
// then go through all the elements in the array, and test for
// a non-nil value. Should result in a single non-nil value.
func TestAastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_REF) // make it an array of references

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(AASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // index to array[20]
	push(&f, ptr)       // store any viable address
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue.(*[]*object.Object))
	var total int64
	for i := 0; i < 30; i++ {
		if array[i] != nil {
			total += 1
		}
	}
	if total != 1 {
		t.Errorf("AASTORE: Expected 1 value not to be nil, got: %d", total)
	}
}

// AASTORE: Test error conditions: invalid array address
func TestAastoreInvalid1(t *testing.T) {
	f := newFrame(AASTORE)
	push(&f, (*object.Object)(nil))                // this should point to an array, will here cause the error
	push(&f, int64(30))                            // the index into the array
	push(&f, object.Make1DimArray(object.REF, 10)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("AASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// AASTORE: Test error conditions: wrong type of array (not [I)
func TestAastoreInvalid2(t *testing.T) {

	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(AASTORE)
	push(&f, o)        // this should point to an array of refs, not ints, will here cause the error
	push(&f, int64(5)) // the index into the array
	push(&f, o)        // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Attempt to access array of incorrect type") {
		t.Errorf("AASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// AASTORE: Test error conditions: index out of range
func TestAastoreInvalid3(t *testing.T) {
	o := object.Make1DimArray(object.REF, 10)
	f := newFrame(AASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, o)         // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "AASTORE: Invalid array subscript") {
		t.Errorf("AASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// ANEWARRAY: creation of array for references
func TestAnewrray(t *testing.T) {
	f := newFrame(ANEWARRAY)
	push(&f, int64(13)) // make the array 13 elements big

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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
	arrayPtr := ptr.Fields[0].Fvalue.(*[]*object.Object)
	if len(*arrayPtr) != 13 {
		t.Errorf("ANEWARRAY: Expecting array length of 13, got %d", len(*arrayPtr))
	}
}

// ANEWARRAY: creation of array for references; test contents of Klass field
func TestAnewrrayKlassField(t *testing.T) {
	f := newFrame(ANEWARRAY)
	push(&f, int64(13)) // make the array 13 elements big

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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
	klassString := ptr.Klass
	if !strings.HasPrefix(*klassString, types.RefArray) {
		t.Errorf("ANEWARRAY: Expecting class to start with '[L', got %s", *klassString)
	}
}

// ANEWARRAY: creation of array for references; test invalid array size
func TestAnewrrayInvalidSize(t *testing.T) {
	f := newFrame(ANEWARRAY)
	push(&f, int64(-1)) // make the array an invalid size

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)
	if err == nil {
		t.Errorf("ANEWARRAY: Did not get expected error")
	}

	msg := err.Error()
	if !(msg == "ANEWARRAY: Invalid size for array") {
		t.Errorf("ANEWARRAY: Expecting different error msg, got %s", msg)
	}
}

// ARRAYLENGTH: Test length of byte array
// First, we create the array of 13 elements, then we push the reference
// to it and execute the ARRAYLENGTH bytecode using the address stored
// in the global array address list
func TestByteArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))                    // make the array 13 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(ARRAYLENGTH)
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 13 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 13, got %d", size)
	}
}

// ARRAYLENGTH: Test length of int array
func TestIntArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(22))                   // make the array 22 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 22 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 13, got %d", size)
	}
}

// ARRAYLENGTH: Test length of float array
func TestFloatArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(34))                      // make the array 34 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 34, got %d", size)
	}
}

// ARRAYLENGTH: Test length of array of longs
func TestLongArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(34))                    // make the array 34 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 34, got %d", size)
	}
}

// ARRAYLENGTH: Test length of array of references
func TestRefArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(34))                   // make the array 34 elements big
	f.Meth = append(f.Meth, object.T_REF) // make it an array of references

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(ARRAYLENGTH)
	// uptr := uintptr(unsafe.Pointer(ptr))
	push(&f, ptr) // push the reference to the array
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	size := pop(&f).(int64)
	if size != 34 {
		t.Errorf("ARRAYLENGTH: Expecting array length of 34, got %d", size)
	}
}

// ARRAYLENGTH: Test length of raw byte array
func TestRawByteArrayLength(t *testing.T) {
	array := []byte{'a', 'b', 'c'}
	f := newFrame(ARRAYLENGTH)
	push(&f, &array) // push the reference to the array
	fs := frames.CreateFrameStack()
	fs.PushFront(&f)    // push the new frame
	err := runFrame(fs) // execute the bytecode

	if err != nil {
		t.Errorf("ARRAYLENGTH: Got unexpected error message: %s", err.Error())
	}

	length := pop(&f).(int64)
	if length != 3 {
		t.Errorf("ARRAYLENGTH: Expecting length of 3, got: %d", length)
	}
}

// ARRAYLENGTH: Test length of raw int8 array
func TestRawInt8ArrayLength(t *testing.T) {
	array := []int8{'a', 'b', 'c'}
	f := newFrame(ARRAYLENGTH)
	push(&f, &array) // push the reference to the array
	fs := frames.CreateFrameStack()
	fs.PushFront(&f)    // push the new frame
	err := runFrame(fs) // execute the bytecode

	if err != nil {
		t.Errorf("ARRAYLENGTH: Got unexpected error message: %s", err.Error())
	}

	length := pop(&f).(int64)
	if length != 3 {
		t.Errorf("ARRAYLENGTH: Expecting length of 3, got: %d", length)
	}
}

// ARRAYLENGTH: Test length of nil array -- should return an error
func TestNilArrayLength(t *testing.T) {
	f := newFrame(ARRAYLENGTH)
	push(&f, nil) // push the reference to the array
	fs := frames.CreateFrameStack()
	fs.PushFront(&f)    // push the new frame
	err := runFrame(fs) // execute the bytecode

	if err == nil {
		t.Errorf("ARRAYLENGTH: Expecting an error message, but got none")
	}

	errMsg := err.Error()
	if errMsg != "ARRAYLENGTHY: invalid (null) reference to an array" {
		t.Errorf("ARRAYLENGTH: Expecting different error msg, got: %s", errMsg)
	}
}

// BALOAD: Test fetching and pushing the value of an element in a byte/boolean array
// The logic here is effectively identical to IALOAD. This code also tests BASTORE.
func TestBaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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
	// ptr := pop(&f).(unsafe.Pointer)
	ptr := pop(&f)

	f = newFrame(BASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, byte(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(BALOAD) // now fetch the value in array[20]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("BALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("BALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// BALOAD: Test exception on nil array address
func TestBaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(BALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode -- should generate exception

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
func TestBaloadInvalidSubscript(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	log.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(BALOAD) // now fetch the value
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(200)) // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

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
func TestBastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(BASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, byte(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue.(*[]byte)) // changed in JACOBIN-282
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(array[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// BASTORE: Tests whether storing an int64 into a byte array does the right thing
func TestBastoreInt64(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(BASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue.(*[]byte))
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(array[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// BASTORE: Test error conditions: invalid array address
func TestBastoreInvalid1(t *testing.T) {
	f := newFrame(BASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, int64(20))             // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

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
func TestBastoreInvalid2(t *testing.T) {

	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(BASTORE)
	push(&f, o)         // this should point to an array of ints, not floats, will here cause the error
	push(&f, int64(30)) // the index into the array
	push(&f, int64(20)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Attempt to access array of incorrect type") {
		t.Errorf("BASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// BASTORE: Test error conditions: index out of range
func TestBastoreInvalid3(t *testing.T) {

	o := object.Make1DimArray(object.BYTE, 10)
	f := newFrame(BASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, int64(20)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "BASTORE: Invalid array subscript") {
		t.Errorf("BASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// CALOAD: Test fetching and pushing the value of an element in an char array
// Chars in Java are two bytes; we accord each one an int64 element. As a result,
// the logic here is effectively identical to IALOAD. This code also tests CASTORE.
func TestCaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_CHAR) // make it an array of chars

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(CASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(CALOAD) // now fetch the value in array[20]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("CALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("CALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// DALOAD: Test fetching and pushing the value of an element in an float array
func TestDaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                      // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(DASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	push(&f, 100.0)     //     pushed twice because it's 64-bits wide
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(DALOAD) // now fetch the value in array[30]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f).(float64)
	if res != 100.0 {
		t.Errorf("FALOAD: Expected loaded array value of 100, got: %e", res)
	}

	if f.TOS != 0 {
		t.Errorf("DALOAD: Top of stack, expected 0, got: %d", f.TOS)
	}
}

// DALOAD: Test exception on nil array address
func TestDaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(DALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("DALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// DALOAD: using an invalid subscript into the array
func TestLaDoadInvalidSubscript(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                      // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	log.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(DALOAD) // now fetch the value
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(200)) // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

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
func TestDastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                      // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_DOUBLE) // make it an array of doubles

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(DASTORE)
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // in array[20]
	push(&f, 100_000_000_000.25) // the value we're storing
	push(&f, 100_000_000_000.25) //   pushed twice due to being 64 bits
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}

	array := *(ptr.Fields[0].Fvalue).(*[]float64)
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
func TestDastoreInvalid1(t *testing.T) {
	f := newFrame(DASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, float64(20.0))         // the value to insert
	push(&f, float64(20.0))

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

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
func TestDastoreInvalid2(t *testing.T) {
	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(DASTORE)
	push(&f, o)             // this should point to an array of floats, not ints, will here cause the error
	push(&f, int64(30))     // the index into the array
	push(&f, float64(20.0)) // the value to insert
	push(&f, float64(20.0))

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Attempt to access array of incorrect type") {
		t.Errorf("DASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// DASTORE: Test error conditions: index out of range
func TestDastoreInvalid3(t *testing.T) {

	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(DASTORE)
	push(&f, o)             // an array of 10 ints, not floats
	push(&f, int64(30))     // the index into the array: it's too big, causing error
	push(&f, float64(20.0)) // the value to insert
	push(&f, float64(20.0))

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("DASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// FALOAD: Test fetching and pushing the value of an element in an float array
func TestFaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                     // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_FLOAT) // make it an array of floats

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(FASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(FALOAD) // now fetch the value in array[30]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f).(float64)
	if res != 100.0 {
		t.Errorf("FALOAD: Expected loaded array value of 100, got: %e", res)
	}

	if f.TOS != -1 {
		t.Errorf("FALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// FALOAD: Test exception on nil array address
func TestFaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(FALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("FALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// FALOAD: using an invalid subscript into the array
func TestFaloadInvalidSubscript(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                     // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_FLOAT) // make it an array of floats

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	log.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(FALOAD) // now fetch the value
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(200)) // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

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

// FASTORE: store value in array of floats
// Create an array of 30 elements, store value 100.0 in array[20], then
// sum all the elements in the array, and test for a sum of 100.0
func TestFastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                     // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_FLOAT) // make it an array of floats

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(FASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue).(*[]float64)
	var fsum float64
	for i := 0; i < 30; i++ {
		fsum += array[i]
	}
	if fsum != 100.0 {
		t.Errorf("FASTORE: Expected sum of array entries to be 100, got: %e", fsum)
	}
}

// FASTORE: Test error conditions: invalid array address
func TestFastoreInvalid1(t *testing.T) {
	f := newFrame(FASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, float64(20.0))         // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

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
func TestFastoreInvalid2(t *testing.T) {
	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(FASTORE)
	push(&f, o)             // this should point to an array of floats, not ints, will here cause the error
	push(&f, int64(30))     // the index into the array
	push(&f, float64(20.0)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Attempt to access array of incorrect type") {
		t.Errorf("FASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// FASTORE: Test error conditions: index out of range
func TestFastoreInvalid3(t *testing.T) {

	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(FASTORE)
	push(&f, o)             // an array of 10 ints, not floats
	push(&f, int64(30))     // the index into the array: it's too big, causing error
	push(&f, float64(20.0)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("FASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// IALOAD: Test fetching and pushing the value of an element in an int array
func TestIaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(IALOAD) // now fetch the value in array[20]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("IALOAD: Expected loaded array value of 100, got: %d", res)
	}

	if f.TOS != -1 {
		t.Errorf("IALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// IALOAD: Test exception on nil array address
func TestIaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(IALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("IALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// IALOAD: using an invalid subscript into the array
func TestIaloadInvalidSubscript(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	log.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(IALOAD) // now fetch the value
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(200)) // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

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

// IASTORE: store value in array of ints
// Create an array of 30 elements, store value 100 in array[20], then
// sum all the elements in the array, and test for a sum of 100.
func TestIastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue).(*[]int64)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("IASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// IASTORE: Test error conditions: invalid array address
func TestIastoreInvalid1(t *testing.T) {
	f := newFrame(IASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, int64(20))             // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

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
func TestIastoreInvalid2(t *testing.T) {

	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(IASTORE)
	push(&f, o)         // this should point to an array of ints, not floats, will here cause the error
	push(&f, int64(30)) // the index into the array
	push(&f, int64(20)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Attempt to access array of incorrect type") {
		t.Errorf("IASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// IASTORE: Test error conditions: index out of range
func TestIastoreInvalid3(t *testing.T) {

	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(IASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, int64(20)) // the value to insert

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "IA/CA/SATORE: Invalid array subscript") {
		t.Errorf("IASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// LALOAD: Test fetching and pushing the value of an element into a long array
func TestLaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(LASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	push(&f, int64(100)) //    push twice due to being 64-bits wide
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(LALOAD) // now fetch the value in array[20]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	// the loaded item should take two slots on the stack, so TOS s/ = 1
	if f.TOS != 1 {
		t.Errorf("LALOAD: Top of stack, expected 1, got: %d", f.TOS)
	}

	res := pop(&f).(int64)
	if res != 100 {
		t.Errorf("LALOAD: Expected loaded array value of 100, got: %d", res)
	}
}

// LALOAD: Test exception on nil array address
func TestLaloadNilArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(LALOAD)
	push(&f, object.Null) // push the reference to the array, here nil
	push(&f, int64(20))   // get contents in array[20]
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode -- should generate exception

	// restore stderr to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	if !strings.Contains(errMsg, "Invalid (null) reference to an array") {
		t.Errorf("LALOAD: Did not get expected err msg for nil array, got: %s",
			errMsg)
	}
}

// LALOAD: using an invalid subscript into the array
func TestLaloadInvalidSubscript(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	globals.InitGlobals("test")
	log.Init()
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// now, get the reference to the array
	ptr := pop(&f).(*object.Object)

	f = newFrame(LALOAD) // now fetch the value
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(200)) // get contents in array[200] which is invalid
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

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
func TestLastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                    // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(LASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	push(&f, int64(100)) //   pushed twice due to being 64 bits
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}

	array := *(ptr.Fields[0].Fvalue).(*[]int64)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("LASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// LASTORE: Test error conditions: invalid array address
func TestLastoreInvalid1(t *testing.T) {
	f := newFrame(LASTORE)
	push(&f, (*object.Object)(nil)) // this should point to an array, will here cause the error
	push(&f, int64(30))             // the index into the array
	push(&f, int64(20))             // the value to insert
	push(&f, int64(20))

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

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
func TestLastoreInvalid2(t *testing.T) {

	o := object.Make1DimArray(object.FLOAT, 10)
	f := newFrame(LASTORE)
	push(&f, o)         // this should point to an array of ints, not floats, will here cause the error
	push(&f, int64(30)) // the index into the array
	push(&f, int64(20)) // the value to insert
	push(&f, int64(20))

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Attempt to access array of incorrect type") {
		t.Errorf("LASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// LASTORE: Test error conditions: index out of range
func TestLastoreInvalid3(t *testing.T) {

	o := object.Make1DimArray(object.INT, 10)
	f := newFrame(LASTORE)
	push(&f, o)         // an array of 10 ints, not floats
	push(&f, int64(30)) // the index into the array: it's too big, causing error
	push(&f, int64(20)) // the value to insert
	push(&f, int64(20))

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "Invalid array subscript") {
		t.Errorf("LASTORE: Did not get expected error msg, got: %s", errMsg)
	}
}

// MULTIANEWARRAY: test creation of a two-dimensional byte array
func Test2DimArray1(t *testing.T) {
	arr, err := object.Make2DimArray(3, 4, object.BYTE)
	if err != nil {
		t.Error("Error creating 2-dimensional array")
	}

	arrLevelArrayPtr := (arr.Fields[0].Fvalue).(*[]*object.Object)
	if len(*arrLevelArrayPtr) != 3 {
		t.Errorf("MULTIANEWARRAY: Expected length of pointer array of 3, got: %d",
			len(*arrLevelArrayPtr))
	}

	leafLevelArrayPtr := ((*arrLevelArrayPtr)[0].Fields[0].Fvalue).(*[]byte)
	arrLen := len(*leafLevelArrayPtr)
	if arrLen != 4 {
		t.Errorf("MULTIANEWARRAY: Expected length of leaf array of 4got: %d",
			arrLen)
	}
}

// MULTIANEWARRAY: test creation of a two-dimensional byte array and its Klass field
func Test2DimArrayKlassField(t *testing.T) {
	arr, err := object.Make2DimArray(3, 4, object.BYTE)
	if err != nil {
		t.Error("Error creating 2-dimensional array")
	}

	if arr.Klass == nil {
		t.Errorf("Array Klass field was nil")
		return
	}

	arrKlass := arr.Klass
	if *arrKlass != "[B" {
		t.Errorf("Expecting array with Klass of '[B', got: %s", *arrKlass)
	}
}

// MULTINEWARRAY: Test a straightforward 3x3x4 array of int64's
func Test3DimArray1(t *testing.T) {
	g := globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	_ = log.SetLogLevel(log.SEVERE)

	// create the constant pool we'll point to
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, 0)
	CP.Utf8Refs = append(CP.Utf8Refs, "[[[I")

	// create the frame
	f := newFrame(MULTIANEWARRAY)
	f.Meth = append(f.Meth, 0x00) // this byte and next form index into CP
	f.Meth = append(f.Meth, 0x02)
	f.Meth = append(f.Meth, 0x03) // the number of dimensions
	push(&f, int64(0x03))         // size of the three dimensions: 4x3x2
	push(&f, int64(0x03))
	push(&f, int64(0x04))
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode
	if f.TOS != 0 {
		t.Errorf("MULTIANEWARRAY: Top of stack, expected 0, got: %d", f.TOS)
	}

	arrayPtr := pop(&f)
	if arrayPtr == nil {
		t.Error("MULTIANEWARRAY: Expected a pointer to an array, got nil")
	}

	topLevelArray := *(arrayPtr.(*object.Object))
	if topLevelArray.Fields[0].Ftype != "[L" {
		t.Errorf("MULTIANEWARRAY: Expected 1st dim to be type '[L', got %s",
			topLevelArray.Fields[0].Ftype)
	}

	dim1 := *(topLevelArray.Fields[0].Fvalue.(*[]*object.Object))
	if len(dim1) != 3 {
		t.Errorf("MULTINEWARRAY: Expected 1st dim to have 3 elements, got: %d",
			len(dim1))
	}

	dim2type := dim1[0].Fields[0].Ftype
	if dim2type != "[[I" {
		t.Errorf("MULTIANEWARRAY: Expected 2nd dim to be type '[[I', got %s",
			dim2type)
	}

	dim2 := *(dim1[0].Fields[0].Fvalue.(*[]*object.Object))
	if len(dim2) != 3 {
		t.Errorf("MULTINEWARRAY: Expected 2nd dim to have 3 elements, got: %d",
			len(dim2))
	}

	dim3type := dim2[0].Fields[0].Ftype
	if dim3type != "[I" {
		t.Errorf("MULTIANEWARRAY: Expected leaf dim to be type '[I', got %s",
			dim3type)
	}

	dim3 := *(dim2[0].Fields[0].Fvalue.(*[]int64))
	if len(dim3) != 4 {
		t.Errorf("MULTINEWARRAY: Expected leaf dim to have 4 elements, got: %d",
			len(dim3))
	}

	elementValue := dim3[2] // an element in the leaf array
	if elementValue != 0 {
		t.Errorf("Expected element value to be 0, got %d", elementValue)
	}
}

// MULTINEWARRAY: Test an array 4x3x3 array of int64's. The zero
// size of the second dimension should result in an single-dimension
// array of int64s
func Test3DimArray2(t *testing.T) {
	g := globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	_ = log.SetLogLevel(log.SEVERE)

	// create the constant pool we'll point to
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, 0)
	CP.Utf8Refs = append(CP.Utf8Refs, "[[[I")

	// create the frame
	f := newFrame(MULTIANEWARRAY)
	f.Meth = append(f.Meth, 0x00) // this byte and next form index into CP
	f.Meth = append(f.Meth, 0x02)
	f.Meth = append(f.Meth, 0x03) // the number of dimensions
	push(&f, int64(0x03))         // size of the three dimensions: 4x3x2
	push(&f, int64(0x00))
	push(&f, int64(0x04))
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode
	if f.TOS != 0 {
		t.Errorf("MULTIANEWARRAY: Top of stack, expected 0, got: %d", f.TOS)
	}

	arrayPtr := pop(&f)
	if arrayPtr == nil {
		t.Error("MULTIANEWARRAY: Expected a pointer to an array, got nil")
	}

	topLevelArray := *(arrayPtr.(*object.Object))
	if topLevelArray.Fields[0].Ftype != "[I" {
		t.Errorf("MULTIANEWARRAY: Expected 1st dim to be type '[I', got %s",
			topLevelArray.Fields[0].Ftype)
	}

	dim1 := *(topLevelArray.Fields[0].Fvalue.(*[]int64))
	if len(dim1) != 4 {
		t.Errorf("MULTINEWARRAY: Expected 1st dim to have 4 elements, got: %d",
			len(dim1))
	}

}

// NEWARRAY: creation of array for primitive values
func TestNewrray(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))                    // make the array 13 elements big
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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
	arrayPtr := ptr.Fields[0].Fvalue.(*[]int64)
	if len(*arrayPtr) != 13 {
		t.Errorf("NEWARRAY: Expecting array length of 13, got %d", len(*arrayPtr))
	}
}

// NEWARRAY: Create new array of 13 bytes
func TestNewrrayForByteArray(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))                    // size
	f.Meth = append(f.Meth, object.T_BYTE) // make it an array of bytes

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	if err != nil {
		t.Errorf("NEWARRAY: Got unexpected error: %s", err.Error())
	}

	arrayPtr := pop(&f).(*object.Object)
	array := arrayPtr.Fields[0].Fvalue.(*[]byte)
	if len(*array) != 13 {
		t.Errorf("NEWARRAY: Got unexpected array size: %d", len(*array))
	}
}

// NEWARRAY: Create new array -- test with invalid size
func TestNewrrayInvalidSize(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(-13))                   // invalid size (less than 0)
	f.Meth = append(f.Meth, object.T_LONG) // make it an array of longs

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	if err == nil {
		t.Errorf("NEWARRAY: Expected an error message, but got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Invalid size for array") {
		t.Errorf("NEWARRAY: Got unexpected error message: %s", errMsg)
	}
}

// NEWARRAY: Create new array -- test with invalid type
func TestNewrrayInvalidType(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))                   // size
	f.Meth = append(f.Meth, object.ERROR) // invalid type

	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	if err == nil {
		t.Errorf("NEWARRAY: Expected an error message, but got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Invalid array type specified") {
		t.Errorf("NEWARRAY: Got unexpected error message: %s", errMsg)
	}
}

// SALOAD: Test fetching and pushing the value of an element in a short array
func TestSaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(SASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(SALOAD) // now fetch the value in array[30]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

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
func TestSastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))                   // make the array 30 elements big
	f.Meth = append(f.Meth, object.T_INT) // make it an array of ints

	globals.InitGlobals("test")
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
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

	f = newFrame(SASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue).(*[]int64)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("SASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

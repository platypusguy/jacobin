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
	"jacobin/javaTypes"
	"jacobin/log"
	"jacobin/object"
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
	oPtr := object.NewString()
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
	push(&f, ptr)                     // push the reference to the array
	push(&f, int64(20))               // in array[20]
	push(&f, javaTypes.JavaByte(100)) // the value we're storing
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
	push(&f, ptr)                     // push the reference to the array
	push(&f, int64(20))               // in array[20]
	push(&f, javaTypes.JavaByte(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	array := *(ptr.Fields[0].Fvalue.(*[]javaTypes.JavaByte))
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(array[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// Tests whether storing an int64 into a byte array does the right thing
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

	array := *(ptr.Fields[0].Fvalue.(*[]javaTypes.JavaByte))
	var sum int64
	for i := 0; i < 30; i++ {
		sum += int64(array[i])
	}
	if sum != 100 {
		t.Errorf("BASTORE: Expected sum of array entries to be 100, got: %d", sum)
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

// MULTIANEWARRAY: test creation of a two-dimensional array
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

	leafLevelArrayPtr := ((*arrLevelArrayPtr)[0].Fields[0].Fvalue).(*[]javaTypes.JavaByte)
	arrLen := len(*leafLevelArrayPtr)
	if arrLen != 4 {
		t.Errorf("MULTIANEWARRAY: Expected length of leaf array of 4got: %d",
			arrLen)
	}
}

// MULTINEWARRAY: Test a straightforward 3x3x4 array of int64's
func Test3DimArray1(t *testing.T) {
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
	//
	// dim2type := dim1[0].Fields[0].Ftype
	// if dim2type != "[[I" {
	//     t.Errorf("MULTIANEWARRAY: Expected 2nd dim to be type '[[I', got %s",
	//         dim2type)
	// }
	//
	// dim2 := *(dim1[0].Fields[0].Fvalue.(*[]*object.Object))
	// if len(dim2) != 3 {
	//     t.Errorf("MULTINEWARRAY: Expected 2nd dim to have 3 elements, got: %d",
	//         len(dim2))
	// }
	//
	// dim3type := dim2[0].Fields[0].Ftype
	// if dim3type != "[I" {
	//     t.Errorf("MULTIANEWARRAY: Expected leaf dim to be type '[I', got %s",
	//         dim3type)
	// }
	//
	// dim3 := *(dim2[0].Fields[0].Fvalue.(*[]int64))
	// if len(dim3) != 4 {
	//     t.Errorf("MULTINEWARRAY: Expected leaf dim to have 4 elements, got: %d",
	//         len(dim3))
	// }
	//
	// elementValue := dim3[2] // an element in the leaf array
	// if elementValue != 0 {
	//     t.Errorf("Expected element value to be 0, got %d", elementValue)
	// }
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

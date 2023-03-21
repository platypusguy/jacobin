/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/frames"
	"jacobin/globals"
	"testing"
	"unsafe"
)

func TestJdkArrayTypeToJacobinType(t *testing.T) {

	a := jdkArrayTypeToJacobinType(T_BOOLEAN)
	if a != BYTE {
		t.Errorf("Expected Jacobin type of %d, got: %d", BYTE, a)
	}

	b := jdkArrayTypeToJacobinType(T_CHAR)
	if b != INT {
		t.Errorf("Expected Jacobin type of %d, got: %d", INT, b)
	}

	c := jdkArrayTypeToJacobinType(T_DOUBLE)
	if c != FLOAT {
		t.Errorf("Expected Jacobin type of %d, got: %d", FLOAT, c)
	}

	d := jdkArrayTypeToJacobinType(999)
	if d != ERROR {
		t.Errorf("Expected Jacobin type of %d, got: %d", ERROR, d)
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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(AASTORE)
	push(&f, ptr)                // push the reference to the array
	push(&f, int64(20))          // in array[20]
	push(&f, unsafe.Pointer(&f)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	f = newFrame(AALOAD) // now fetch the value in array[20]
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // get contents in array[20]
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	res := pop(&f).(unsafe.Pointer)
	if res != unsafe.Pointer(&f) {
		t.Errorf("AALOAD: Expected loaded array value = address of frame, got: %X", res)
	}

	if f.TOS != -1 {
		t.Errorf("AALOAD: Top of stack, expected -1, got: %d", f.TOS)
	}
}

// ANEWARRAY: creation of array for primitive values
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
	ptr := element.Value.(*JacobinRefArray)
	if len(*ptr.Arr) != 13 {
		t.Errorf("ANEWARRAY: Expecting array length of 13, got %d", len(*ptr.Arr))
	}
}

// ARRAYLENGTH: Test length of byte array
// First, we create the array of 13 elements, then we push the reference
// to it and execute the ARRAYLENGTH bytecode using the address stored
// in the global array address list
func TestByteArrayLength(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))             // make the array 13 elements big
	f.Meth = append(f.Meth, T_BYTE) // make it an array of bytes

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
	push(&f, int64(22))            // make the array 22 elements big
	f.Meth = append(f.Meth, T_INT) // make it an array of ints

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
	push(&f, int64(34))               // make the array 34 elements big
	f.Meth = append(f.Meth, T_DOUBLE) // make it an array of doubles

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
	push(&f, int64(30))             // make the array 30 elements big
	f.Meth = append(f.Meth, T_BYTE) // make it an array of bytes

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(BASTORE)
	push(&f, ptr)           // push the reference to the array
	push(&f, int64(20))     // in array[20]
	push(&f, JavaByte(100)) // the value we're storing
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

// BASTORE: store value in array of ints
// Create an array of 30 elements, store value 100 in array[20], then
// sum all the elements in the array, and test for a sum of 100.
// Note the value we store must be an int64 value--not a byte
func TestBastore(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))             // make the array 30 elements big
	f.Meth = append(f.Meth, T_BYTE) // make it an array of bytes

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(BASTORE)
	push(&f, ptr)           // push the reference to the array
	push(&f, int64(20))     // in array[20]
	push(&f, JavaByte(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	byteRef := (*JacobinByteArray)(ptr)
	array := *(byteRef.Arr)
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
	push(&f, int64(30))             // make the array 30 elements big
	f.Meth = append(f.Meth, T_BYTE) // make it an array of bytes

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(BASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	byteRef := (*JacobinByteArray)(ptr)
	array := *(byteRef.Arr)
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
	push(&f, int64(30))             // make the array 30 elements big
	f.Meth = append(f.Meth, T_CHAR) // make it an array of chars

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
	ptr := pop(&f).(unsafe.Pointer)

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
	push(&f, int64(30))               // make the array 30 elements big
	f.Meth = append(f.Meth, T_DOUBLE) // make it an array of doubles

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
	ptr := pop(&f).(unsafe.Pointer)

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
	push(&f, int64(30))               // make the array 30 elements big
	f.Meth = append(f.Meth, T_DOUBLE) // make it an array of doubles

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
	ptr := pop(&f).(unsafe.Pointer)

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

	floatRef := (*JacobinFloatArray)(ptr)
	array := *(floatRef.Arr)
	var sum float64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100_000_000_000.25 {
		t.Errorf("DASTORE: Expected sum of doubles array to be 100,000,000,000.25, got: %f",
			sum)
	}
}

// FALOAD: Test fetching and pushing the value of an element in an float array
func TestFaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))              // make the array 30 elements big
	f.Meth = append(f.Meth, T_FLOAT) // make it an array of floats

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
	ptr := pop(&f).(unsafe.Pointer)

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
	push(&f, int64(30))              // make the array 30 elements big
	f.Meth = append(f.Meth, T_FLOAT) // make it an array of floats

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(FASTORE)
	push(&f, ptr)       // push the reference to the array
	push(&f, int64(20)) // in array[20]
	push(&f, 100.0)     // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	floatRef := (*JacobinFloatArray)(ptr)
	array := *(floatRef.Arr)
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
	push(&f, int64(30))            // make the array 30 elements big
	f.Meth = append(f.Meth, T_INT) // make it an array of ints

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
	ptr := pop(&f).(unsafe.Pointer)

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
	push(&f, int64(30))            // make the array 30 elements big
	f.Meth = append(f.Meth, T_INT) // make it an array of ints

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(IASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	intRef := (*JacobinIntArray)(ptr)
	array := *(intRef.Arr)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("IASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// LALOAD: Test fetching and pushing the value of an element in an long array
func TestLaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))             // make the array 30 elements big
	f.Meth = append(f.Meth, T_LONG) // make it an array of longs

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
	ptr := pop(&f).(unsafe.Pointer)

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
	push(&f, int64(30))             // make the array 30 elements big
	f.Meth = append(f.Meth, T_LONG) // make it an array of longs

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
	ptr := pop(&f).(unsafe.Pointer)

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

	intRef := (*JacobinIntArray)(ptr)
	array := *(intRef.Arr)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("LASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

// NEWARRAY: creation of array for primitive values
func TestNewrray(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))             // make the array 13 elements big
	f.Meth = append(f.Meth, T_LONG) // make it an array of longs

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
	ptr := element.Value.(*JacobinIntArray)
	if len(*ptr.Arr) != 13 {
		t.Errorf("NEWARRAY: Expecting array length of 13, got %d", len(*ptr.Arr))
	}
}

// SALOAD: Test fetching and pushing the value of an element in a short array
func TestSaload(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(30))            // make the array 30 elements big
	f.Meth = append(f.Meth, T_INT) // make it an array of ints

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(IASTORE)
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
	push(&f, int64(30))            // make the array 30 elements big
	f.Meth = append(f.Meth, T_INT) // make it an array of ints

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
	ptr := pop(&f).(unsafe.Pointer)

	f = newFrame(SASTORE)
	push(&f, ptr)        // push the reference to the array
	push(&f, int64(20))  // in array[20]
	push(&f, int64(100)) // the value we're storing
	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs) // execute the bytecode

	intRef := (*JacobinIntArray)(ptr)
	array := *(intRef.Arr)
	var sum int64
	for i := 0; i < 30; i++ {
		sum += array[i]
	}
	if sum != 100 {
		t.Errorf("SASTORE: Expected sum of array entries to be 100, got: %d", sum)
	}
}

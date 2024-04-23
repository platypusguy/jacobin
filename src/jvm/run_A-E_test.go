/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/opcodes"
	"math"
	"os"
	"strings"
	"testing"
	"unsafe"
)

// These tests test the individual bytecode instructions. They are presented
// here in alphabetical order of the instruction name.
// THIS FILE CONTAINS TESTS FOR ALL BYTECODES UP TO DUP2_X1.
// All other bytecodes are in run_*_test.go files except
// for array bytecodes, which are located in arrayBytecodes_test.go

// set up function to create a frame with a method with the single instruction
// that's being tested
func newFrame(code byte) frames.Frame {
	f := frames.CreateFrame(6)
	f.Ftype = 'J'
	f.Meth = append(f.Meth, code)
	return *f
}

var zero = int64(0)
var zerof = float64(0)

var maxFloatDiff = .00001

func validateFloatingPoint(t *testing.T, op string, expected float64, actual float64) {
	if math.Abs(expected-actual) > maxFloatDiff {
		t.Errorf("%s: expected a result of %f, but got: %f", op, expected, actual)
	}
}

// ---- tests ----

// ACONST_NULL: Load null onto opStack
func TestAconstNull(t *testing.T) {
	f := newFrame(opcodes.ACONST_NULL)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := peek(&f)
	if x != object.Null {
		t.Errorf("ACONST_NULL: Expecting nil on stack, got: %d", x)
	}
	if f.TOS != 0 {
		t.Errorf("ACONST_NULL: Expecting TOS = 0, but tos is: %d", f.TOS)
	}
}

// ALOAD: test load of reference in locals[index] on to stack
func TestAload(t *testing.T) {
	f := newFrame(opcodes.ALOAD)
	f.Meth = append(f.Meth, 0x04) // use local var #4
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x1234562)) // put value in locals[4]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	if x != 0x1234562 {
		t.Errorf("ALOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("ALOAD: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
	if f.PC != 2 {
		t.Errorf("ALOAD: Expected pc to be pointing at byte 2, got: %d", f.PC)
	}
}

// ALOAD_0: test load of reference in locals[0] on to stack
func TestAload0(t *testing.T) {
	f := newFrame(opcodes.ALOAD_0)
	f.Locals = append(f.Locals, int64(0x1234560)) // put value in locals[0]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	if x != 0x1234560 {
		t.Errorf("ALOAD_0: Expecting 0x1234560 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("ALOAD_0: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ALOAD_1: test load of reference in locals[1] on to stack
func TestAload1(t *testing.T) {
	f := newFrame(opcodes.ALOAD_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x1234561)) // put value in locals[1]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	if x != 0x1234561 {
		t.Errorf("ALOAD_1: Expecting 0x1234561 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("ALOAD_1: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ALOAD_2: test load of reference in locals[2] on to stack
func TestAload2(t *testing.T) {
	f := newFrame(opcodes.ALOAD_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x1234562)) // put value in locals[2]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	if x != 0x1234562 {
		t.Errorf("ALOAD_2: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("ALOAD_2: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ALOAD_3: test load of reference in locals[3] on to stack
func TestAload3(t *testing.T) {
	f := newFrame(opcodes.ALOAD_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x1234563)) // put value in locals[3]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	if x != 0x1234563 {
		t.Errorf("ALOAD_3: Expecting 0x1234563 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("ALOAD_3: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ARETURN: Return a long from a function
func TestAreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, unsafe.Pointer(&f0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f0)

	// create a new frame which does an ARETURN of pointer to f1
	f1 := newFrame(opcodes.ARETURN)
	push(&f1, unsafe.Pointer(&f1))
	fs.PushFront(&f1)
	_ = runFrame(fs)

	// now that the ARETURN has completed, pop that frame (the one that did the ARETURN)
	_ = frames.PopFrame(fs)

	// and see whether the pointer at the frame's top of stack points to f1
	f2 := fs.Front().Value.(*frames.Frame)
	newVal := pop(f2).(unsafe.Pointer)
	if newVal != unsafe.Pointer(&f1) {
		t.Error("ARETURN: did not get expected value of reference")
	}
}

// ASTORE: Store reference in local var specified by following byte.
func TestAstore(t *testing.T) {
	f := newFrame(opcodes.ASTORE)
	f.Meth = append(f.Meth, 0x03) // use local var #4

	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22223))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[3] != int64(0x22223) {
		t.Errorf("ASTORE: Expecting 0x22223 in locals[3], got: 0x%x", f.Locals[3])
	}
	if f.TOS != -1 {
		t.Errorf("ASTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ASTORE_0: test store of reference from stack into locals[0]
func TestAstore0(t *testing.T) {
	f := newFrame(opcodes.ASTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22220))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[0] != int64(0x22220) {
		t.Errorf("ASTORE_0: Expecting 0x22220 on stack, got: 0x%x", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ASTORE_0: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ASTORE_1: test store of reference from stack into locals[1]
func TestAstore1(t *testing.T) {
	f := newFrame(opcodes.ASTORE_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22221))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[1] != int64(0x22221) {
		t.Errorf("ASTORE_1: Expecting 0x22221 on stack, got: 0x%x", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ASTORE_1: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ASTORE_2: test store of reference from stack into locals[2]
func TestAstore2(t *testing.T) {
	f := newFrame(opcodes.ASTORE_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22222))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2] != int64(0x22222) {
		t.Errorf("ASTORE_2: Expecting 0x22222 on stack, got: 0x%x", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ASTORE_2: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ASTORE3: store of reference from stack into locals[3]
func TestAstore3(t *testing.T) {
	f := newFrame(opcodes.ASTORE_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22223))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[3] != int64(0x22223) {
		t.Errorf("ASTORE_3: Expecting 0x22223 on stack, got: 0x%x", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ASTORE_3: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// BASTORE is tested in arrayBytecodes_test.go

// BIPUSH
func TestBipush(t *testing.T) {
	f := newFrame(opcodes.BIPUSH)
	f.Meth = append(f.Meth, 0x05)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 5 {
		t.Errorf("BIPUSH: Expected popped value to be 5, got: %d", value)
	}
}

// BIPUSH with negative value
func TestBipushNeg(t *testing.T) {
	f := newFrame(opcodes.BIPUSH)
	val := -5
	f.Meth = append(f.Meth, byte(val))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != -5 {
		t.Errorf("BIPUSH: Expected popped value to be -5, got: %d", value)
	}
}

// CHECKCAST: This bytecode uses similar logic to INSTANCEOF, except how
// it handles exceptional conditions.
func TestCheckcastOfString(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	log.Init()

	classloader.Init()
	// classloader.LoadBaseClasses()
	classloader.MethAreaInsert("java/lang/String",
		&(classloader.Klass{
			Status: 'X', // use a status that's not subsequently tested for.
			Loader: "bootstrap",
			Data:   nil,
		}))
	s := object.StringObjectFromGoString("hello world")

	f := newFrame(opcodes.CHECKCAST)
	f.Meth = append(f.Meth, 0) // point to entry [2] in CP
	f.Meth = append(f.Meth, 1) // " "

	// now create the CP.
	// [0] First entry is perforce 0
	// [1] is a ClassRef that points to the UTF8 string in [1]
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, object.StringPoolStringIndex) // point to string pool entry
	f.CP = &CP

	push(&f, s)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(*object.Object)
	if value != s { // if the stack is unchanged, we got a match
		t.Errorf(" CHECKCAST: Expected stack not found on successful check")
	}
}

// CHECKCAST: Test for nil. This should simply move the PC forward by 3 bytes. Nothing else.
func TestCheckcastOfNil(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	log.Init()

	classloader.Init()
	// classloader.LoadBaseClasses()
	classloader.MethAreaInsert("java/lang/String",
		&(classloader.Klass{
			Status: 'X', // use a status that's not subsequently tested for.
			Loader: "bootstrap",
			Data:   nil,
		}))

	f := newFrame(opcodes.CHECKCAST)
	push(&f, nil) // this should cause the error

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 {
		t.Errorf("CHECKCAST: Expected TOS to be 0, got %d", f.TOS)
	}

	if f.PC != 3 { // skip two bytes are error is discovered, +1 to get to next bytecode
		t.Errorf(" CHECKCAST: Expected PC to be at 3, got %d", f.PC)
	}
}

// CHECKCAST: Test for null -- this should simply move the PC forward by 3 bytes. Nothing else.
func TestCheckcastOfNull(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	log.Init()

	f := newFrame(opcodes.CHECKCAST)
	push(&f, object.Null) // this should cause the error

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 {
		t.Errorf("CHECKCAST: Expected TOS to be 0, got %d", f.TOS)
	}

	if f.PC != 3 { // skip two bytes are error is discovered, +1 to get to next bytecode
		t.Errorf(" CHECKCAST: Expected PC to be at 3, got %d", f.PC)
	}
}

// CHECKCAST: Test for non-object pointer -- this should result in an exception
func TestCheckcastOfInvalidReference(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	log.Init()
	log.SetLogLevel(log.SEVERE)

	// redirect stderr to avoid printing error message to console
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.CHECKCAST)
	push(&f, float64(42.0)) // this should cause the error

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	os.Stderr = normalStderr // restore stderr

	if err == nil {
		t.Errorf("CHECKCAST: Expected an error, but did not get one")
	}

	if f.TOS != 0 {
		t.Errorf("CHECKCAST: Expected TOS to be 0, got %d", f.TOS)
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "CHECKCAST: Invalid class reference") {
		t.Errorf("CHECKCAST: Expected different error message. Got: %s", errMsg)
	}
}

// D2F: test convert double to float
func TestD2f(t *testing.T) {
	f := newFrame(opcodes.D2F)
	push(&f, 2.9)
	push(&f, 2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(float64)
	if math.Abs(val-2.9) > maxFloatDiff {
		t.Errorf("D2F: expected a result of 2.9, but got: %f", val)
	}
	if f.TOS != -1 {
		t.Errorf("D2F: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// D2I: test convert double to int, positive
func TestD2iPositive(t *testing.T) {
	f := newFrame(opcodes.D2I)
	push(&f, 2.9)
	push(&f, 2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != 2 {
		t.Errorf("D2I: expected a result of 2, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("D2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// D2I: test convert double to int, negative
func TestD2iNegative(t *testing.T) {
	f := newFrame(opcodes.D2I)
	push(&f, -2.9)
	push(&f, -2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != -2 {
		t.Errorf("D2I: expected a result of -2, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("D2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// D2L: test convert double to long, positive
func TestD2lPositive(t *testing.T) {
	f := newFrame(opcodes.D2L)
	push(&f, 2.9)
	push(&f, 2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	val := pop(&f).(int64)
	if val != 2 {
		t.Errorf("D2L: expected a result of 2, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("D2L: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// D2L: test convert double to long, negative
func TestD2lNegative(t *testing.T) {
	f := newFrame(opcodes.D2L)
	push(&f, -2.9)
	push(&f, -2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	val := pop(&f).(int64)
	if val != -2 {
		t.Errorf("D2L: expected a result of -2, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("D2L: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DADD: test add two doubles
func TestDadd(t *testing.T) {
	f := newFrame(opcodes.DADD)
	push(&f, 15.3)
	push(&f, 15.3)
	push(&f, 22.1)
	push(&f, 22.1)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	pop(&f)
	val := pop(&f).(float64)
	validateFloatingPoint(t, "DADD", 37.4, val)
	if f.TOS != -1 {
		t.Errorf("DADD: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DADD: test add two doubles, one is NaN
func TestDaddNan(t *testing.T) {
	f := newFrame(opcodes.DADD)
	push(&f, 15.3)
	push(&f, 15.3)
	push(&f, math.NaN())
	push(&f, math.NaN())

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	pop(&f)
	val := pop(&f).(float64)

	if !math.IsNaN(val) {
		t.Errorf("DADD: Expected NaN, got: %f", val)
	}

	if f.TOS != -1 {
		t.Errorf("DADD: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DADD: test add two doubles, one is Inf
func TestDaddInf(t *testing.T) {
	f := newFrame(opcodes.DADD)
	push(&f, 15.3)
	push(&f, 15.3)
	push(&f, math.Inf(1))
	push(&f, math.Inf(1))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	pop(&f)
	val := pop(&f).(float64)

	if !math.IsInf(val, 1) {
		t.Errorf("DADD: Expected Inf, got: %f", val)
	}

	if f.TOS != -1 {
		t.Errorf("DADD: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DCMP0: compare two doubles, 1 on NaN
func TestDcmpg1(t *testing.T) {
	f := newFrame(opcodes.DCMPG)
	push(&f, 3.0)
	push(&f, 3.0)
	push(&f, 2.0)
	push(&f, 2.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 1 {
		t.Errorf("DCMPG: Expected value to be 1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("DCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DCMPG: compare two doubles
func TestDcmpgMinus1(t *testing.T) {
	f := newFrame(opcodes.DCMPG)
	push(&f, 2.0)
	push(&f, 2.0)
	push(&f, 3.0)
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != -1 {
		t.Errorf("DCMPG: Expected value to be -1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("DCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DCMP0: compare two doubles
func TestDcmpg0(t *testing.T) {
	f := newFrame(opcodes.DCMPG)
	push(&f, 3.0)
	push(&f, 3.0)
	push(&f, 3.0)
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 0 {
		t.Errorf("DCMPG: Expected value to be 0, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("DCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

func TestDcmpgNan(t *testing.T) {
	f := newFrame(opcodes.DCMPG)
	push(&f, math.NaN())
	push(&f, math.NaN())
	push(&f, 3.0)
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 1 {
		t.Errorf("DCMPG: Expected value to be 1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("DCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DCMPL
func TestDcmplNan(t *testing.T) {
	f := newFrame(opcodes.DCMPL)
	push(&f, math.NaN())
	push(&f, math.NaN())
	push(&f, 3.0)
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != -1 {
		t.Errorf("DCMPL: Expected value to be -1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("DCMPL: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DCONST_0
func TestDconst0(t *testing.T) {
	f := newFrame(opcodes.DCONST_0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)

	if value != 0.0 {
		t.Errorf("DCONST_0: Expected popped value to be 0.0, got: %f", value)
	}

	if f.TOS != -1 {
		t.Errorf("DCONST_0: Expected empty stack, got: %d", f.TOS)
	}
}

// DCONST_1
func TestDconst1(t *testing.T) {
	f := newFrame(opcodes.DCONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)

	if value != 1.0 {
		t.Errorf("DCONST_1: Expected popped value to be 1.0, got: %f", value)
	}

	if f.TOS != -1 {
		t.Errorf("Expected empty stack, got: %d", f.TOS)
	}
}

// DDIV: double divide of.TOS-1 by tos, push result
func TestDdiv(t *testing.T) {
	f := newFrame(opcodes.DDIV)
	push(&f, 3.0)
	push(&f, 3.0)
	push(&f, 2.0)
	push(&f, 2.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)
	validateFloatingPoint(t, "DDIV", 1.5, value)
	if f.TOS != -1 {
		t.Errorf("DDIV: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// DDIV: with divide zero by zero, should = NaN
func TestDdivDivideZeroByZero(t *testing.T) {
	f := newFrame(opcodes.DDIV)
	push(&f, float64(0))
	push(&f, float64(0))

	push(&f, float64(0))
	push(&f, float64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	ret := pop(&f)

	if !math.IsNaN(ret.(float64)) {
		t.Errorf("DDIV: Did not get an expected NaN")
	}
}

// DDIV: with divide positive number by zero, should = +Inf
func TestDdivDividePosNumberByZero(t *testing.T) {
	f := newFrame(opcodes.DDIV)
	push(&f, float64(10))
	push(&f, float64(10))

	push(&f, float64(0))
	push(&f, float64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	ret := pop(&f)

	if !math.IsInf(ret.(float64), 1) {
		t.Errorf("DDIV: Did not get an expected +Infinity")
	}
}

// DLOAD: test load of double in locals[index] on to stack
func TestDload(t *testing.T) {
	f := newFrame(opcodes.DLOAD)
	f.Meth = append(f.Meth, 0x04) // use local var #4
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, float64(0x1234562)) // put value in locals[4]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(float64)
	pop(&f) // pop twice due to two entries on op stack due to 64-bit width of data type
	if x != float64(0x1234562) {
		t.Errorf("DLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("DLOAD: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
	if f.PC != 2 {
		t.Errorf("DLOAD: Expected pc to be pointing at byte 2, got: %d", f.PC)
	}
}

// DLOAD_0: load of double in locals[0] onto stack
func TestDload0(t *testing.T) {
	f := newFrame(opcodes.DLOAD_0)
	f.Locals = append(f.Locals, 1.2)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)
	validateFloatingPoint(t, "DLOAD_0", 1.2, value)

	if f.TOS != -1 {
		t.Errorf("DLOAD_0: Expected empty stack, got: %d", f.TOS)
	}
}

// DLOAD_1: load of double in locals[1] onto stack
func TestDload1(t *testing.T) {
	f := newFrame(opcodes.DLOAD_1)
	f.Locals = append(f.Locals, 1.3)
	f.Locals = append(f.Locals, 1.2)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)
	validateFloatingPoint(t, "DLOAD_1", 1.2, value)

	if f.TOS != -1 {
		t.Errorf("DLOAD_1: Expected empty stack, got: %d", f.TOS)
	}
}

// DLOAD_2: load of double in locals[2] onto stack
func TestDload2(t *testing.T) {
	f := newFrame(opcodes.DLOAD_2)
	f.Locals = append(f.Locals, 1.3)
	f.Locals = append(f.Locals, 1.3)
	f.Locals = append(f.Locals, 1.2)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)
	validateFloatingPoint(t, "DLOAD_2", 1.2, value)

	if f.TOS != -1 {
		t.Errorf("DLOAD_2: Expected empty stack, got: %d", f.TOS)
	}
}

// DLOAD_3: load of double in locals[3] onto stack
func TestDload3(t *testing.T) {
	f := newFrame(opcodes.DLOAD_3)
	f.Locals = append(f.Locals, 1.3)
	f.Locals = append(f.Locals, 1.3)
	f.Locals = append(f.Locals, 1.3)
	f.Locals = append(f.Locals, 1.2)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	value := pop(&f).(float64)
	validateFloatingPoint(t, "DLOAD_3", 1.2, value)

	if f.TOS != -1 {
		t.Errorf("DLOAD_3: Expected empty stack, got: %d", f.TOS)
	}
}

// DMUL (pop 2 doubles, multiply them, push result)
func TestDmul(t *testing.T) {
	f := newFrame(opcodes.DMUL)
	push(&f, 1.5)
	push(&f, 1.5)
	push(&f, 2.0)
	push(&f, 2.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	validateFloatingPoint(t, "DMUL", 3.0, pop(&f).(float64))

	if f.TOS != -1 {
		t.Errorf("DMUL, Top of stack, expected 0, got: %d", f.TOS)
	}
}

// DNEG Negate a double
func TestDneg(t *testing.T) {
	f := newFrame(opcodes.DNEG)
	push(&f, 1.5)
	push(&f, 1.5)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	validateFloatingPoint(t, "DNEG", -1.5, pop(&f).(float64))

	if f.TOS != -1 {
		t.Errorf("DNEG, Top of stack, expected 0, got: %d", f.TOS)
	}
}

// DNEG Negate a double - infinity
func TestDnegInf(t *testing.T) {
	f := newFrame(opcodes.DNEG)
	push(&f, math.Inf(1))
	push(&f, math.Inf(1))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	pop(&f)
	val := pop(&f).(float64)

	if math.Inf(-1) != val {
		t.Errorf("Expected negative infinity, got %f", val)
	}

	if f.TOS != -1 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
}

// DREM: remainder of float division (the % operator)
func TestDrem(t *testing.T) {
	f := newFrame(opcodes.DREM)
	push(&f, 23.5)
	push(&f, 23.5)
	push(&f, 3.3)
	push(&f, 3.3)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 1 {
		t.Errorf("DREM, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(float64)
	if math.Abs(value-0.40000033) > maxFloatDiff {
		t.Errorf("DREM: Expected popped value to be 0.40000033, got: %f", value)
	}
}

// DRETURN: Return a long from a function
func TestDreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, float64(20))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f0)
	f1 := newFrame(opcodes.DRETURN)
	push(&f1, float64(21))
	push(&f1, float64(21))
	fs.PushFront(&f1)
	_ = runFrame(fs)
	_ = frames.PopFrame(fs)
	f3 := fs.Front().Value.(*frames.Frame)
	newVal := pop(f3).(float64)
	if newVal != 21.0 {
		t.Errorf("After DRETURN, expected a value of 21 in previous frame, got: %f", newVal)
	}
	pop(f3) // popped a second time due to longs taking two slots

	prevVal := pop(f3).(float64)
	if prevVal != 20 {
		t.Errorf("After DRETURN, expected a value of 20 in 2nd place of previous frame, got: %f", prevVal)
	}
}

// DSTORE: Store double from stack into local specified by following byte.
func TestDstore(t *testing.T) {
	f := newFrame(opcodes.DSTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)

	push(&f, float64(0x22223)) // pushed twice due to double using two slots
	push(&f, float64(0x22223))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2] != float64(0x22223) {
		t.Errorf("DSTORE: Expecting 0x22223 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.Locals[3] != float64(0x22223) {
		t.Errorf("DSTORE: Expecting 0x22223 in locals[3], got: 0x%x", f.Locals[3])
	}

	if f.TOS != -1 {
		t.Errorf("DSTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// DSTORE_0: Store double from stack into localVar[0]
func TestDstore0(t *testing.T) {
	f := newFrame(opcodes.DSTORE_0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	push(&f, 1.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[0].(float64) != 1.0 {
		t.Errorf("DSTORE_0: expected locals[0] to be 1.0, got: %f", f.Locals[0].(float64))
	}

	if f.TOS != -1 {
		t.Errorf("DSTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// DSTORE_1: Store double from stack into localVar[1]
func TestDstore1(t *testing.T) {
	f := newFrame(opcodes.DSTORE_1)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	push(&f, 1.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[1].(float64) != 1.0 {
		t.Errorf("DSTORE_1: expected locals[1] to be 1.0, got: %f", f.Locals[1].(float64))
	}

	if f.TOS != -1 {
		t.Errorf("DSTORE_1: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// DSTORE_2: Store double from stack into localVar[2]
func TestDstore2(t *testing.T) {
	f := newFrame(opcodes.DSTORE_2)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	push(&f, 1.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2].(float64) != 1.0 {
		t.Errorf("DSTORE_2: expected locals[2] to be 1.0, got: %f", f.Locals[2].(float64))
	}

	if f.TOS != -1 {
		t.Errorf("DSTORE_2: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// DSTORE_3: Store double from stack into localVar[3]
func TestDstore3(t *testing.T) {
	f := newFrame(opcodes.DSTORE_3)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	push(&f, 1.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[3].(float64) != 1.0 {
		t.Errorf("DSTORE_3: expected locals[3] to be 1.0, got: %f", f.Locals[3].(float64))
	}

	if f.TOS != -1 {
		t.Errorf("DSTORE_3: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// DSUB: double subtraction
func TestDsub(t *testing.T) {
	f := newFrame(opcodes.DSUB)
	push(&f, 1.0)
	push(&f, 1.0)
	push(&f, 0.7)
	push(&f, 0.7)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(float64)
	pop(&f)

	if math.Abs(value-0.3) > maxFloatDiff {
		t.Errorf("DSUB: Expected popped value to be 0.3, got: %f", value)
	}

	if f.TOS != -1 {
		t.Errorf("DSUB, Empty stack expected, got: %d", f.TOS)
	}
}

// DUP: Push a duplicate of the top item on the stack
func TestDup(t *testing.T) {
	f := newFrame(opcodes.DUP)
	push(&f, int64(0x22223))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS < 1 {
		t.Errorf("DUP: stack should have two elements with tos > 0, tos was: %d", f.TOS)
	}

	a := pop(&f).(int64)
	b := pop(&f).(int64)
	if a != 0x22223 || b != 0x22223 {
		t.Errorf(
			"DUP: popped values are incorrect. Expecting 0x22223, got: %X and %X", a, b)
	}
}

// DUP2: Push duplicate of the top two items on the stack
func TestDup2(t *testing.T) {
	f := newFrame(opcodes.DUP2)
	push(&f, int64(0x22))
	push(&f, int64(0x11)) // this is TOS

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 3 {
		t.Errorf("DUP2: stack should have four elements, got tos was: %d", f.TOS)
	}

	a := pop(&f).(int64)
	b := pop(&f).(int64)
	c := pop(&f).(int64)
	d := pop(&f).(int64)
	if a != 0x11 || c != 0x11 {
		t.Errorf(
			"DUP2: popped values are incorrect. Expecting 0x11, got: %X and %X", a, c)
	}
	if b != 0x22 || d != 0x22 {
		t.Errorf(
			"DUP2: popped values are incorrect. Expecting 0x22, got: %X and %X", b, d)
	}
}

// DUP_X1: Duplicate the top stack value and insert it two slots down
func TestDupX1(t *testing.T) {
	f := newFrame(opcodes.DUP_X1)
	push(&f, int64(0x3))
	push(&f, int64(0x2))
	push(&f, int64(0x1))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 3 {
		t.Errorf("DUP_X1: Expecting a top of stack = 3 (so stack size 4), got: %d", f.TOS)
	}

	a := pop(&f).(int64)
	b := pop(&f).(int64)
	c := pop(&f).(int64)
	if a != 1 || c != 1 {
		t.Errorf(
			"DUP_X1: popped values are incorrect. Expecting value of 1, got: %X and %X", a, b)
	}
}

// DUP_X2: Duplicate the top stack value and insert it three slots down
func TestDupX2(t *testing.T) {
	f := newFrame(opcodes.DUP_X2)
	push(&f, int64(0x3))
	push(&f, int64(0x2))
	push(&f, int64(0x1)) // this will be the dup'ed value
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 3 {
		t.Errorf("DUP_X2: Expecting a top of stack = 3 (so stack size 4), got: %d", f.TOS)
	}

	a := pop(&f).(int64)
	pop(&f)
	pop(&f)
	d := pop(&f).(int64)
	if a != 1 || d != 1 {
		t.Errorf(
			"DUP_X2: popped values are incorrect. Expecting value of 1, got: %X and %X", a, d)
	}
}

// DUP2_X1: Duplicate the top 2 stack values and insert them 3 slots down
func TestDup2X1(t *testing.T) {
	f := newFrame(opcodes.DUP2_X1)
	push(&f, int64(0x3))
	push(&f, int64(0x2))
	push(&f, int64(0x1)) // this is nowdir TOS
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 4 {
		t.Errorf("DUP2_X1: Expecting a top of stack = 4 (so stack size 5), got: %d", f.TOS)
	}

	a := pop(&f).(int64)
	pop(&f)
	pop(&f)
	d := pop(&f).(int64)
	if a != 1 || d != 1 {
		t.Errorf(
			"DUP2_X1: popped values are incorrect. Expecting value of 1, got: %X and %X", a, d)
	}
}

// DUP2_X1: Duplicate the top 2 stack values and insert them 4 slots down
func TestDup2X2(t *testing.T) {
	f := newFrame(opcodes.DUP2_X2)
	push(&f, int64(0x4))
	push(&f, int64(0x3))
	push(&f, int64(0x2))
	push(&f, int64(0x1)) // this is now TOS
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 5 {
		t.Errorf("DUP2_X2: Expecting a top of stack = 5 (so stack size 6), got: %d", f.TOS)
	}

	a := pop(&f).(int64)
	pop(&f)
	pop(&f)
	pop(&f)
	e := pop(&f).(int64)
	if a != 1 || e != 1 {
		t.Errorf(
			"DUP2_X2: popped values are incorrect. Expecting value of 1, got: %X and %X", a, e)
	}
}

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
	"jacobin/statics"
	"jacobin/types"
	"math"
	"os"
	"strings"
	"testing"
	"unsafe"
)

// These tests test the individual bytecode instructions. They are presented
// here in alphabetical order of the instruction name.
// THIS FILE CONTAINS TESTS FOR ALL BYTECODES UP TO AND INCLUDING IFNULL.
// All other bytecodes that come after IFNULL are in run_part2_test.go *except
// for array bytecodes*, which are located in arrays_test.go

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
	s := object.NewStringFromGoString("hello world")

	f := newFrame(opcodes.CHECKCAST)
	f.Meth = append(f.Meth, 0) // point to entry [2] in CP
	f.Meth = append(f.Meth, 2) // " "

	// now create the CP. First entry is perforce 0
	// [1] entry points to a UTF8 entry with the class name
	// [2] is a ClassRef that points to the UTF8 string in [1]
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, 1) // point to record 1 in CP: Utf8 for class name
	CP.Utf8Refs = append(CP.Utf8Refs, "java/lang/String")
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

// F2D: test convert float to double
func TestF2d(t *testing.T) {
	f := newFrame(opcodes.F2D)
	push(&f, 2.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(float64)
	if val != 2.0 {
		t.Errorf("F2D: expected a result of 2.0, but got: %f", val)
	}
	if f.TOS != 0 {
		t.Errorf("F2D: Expected stack with 1 item, but got a TOS of: %d", f.TOS)
	}
}

// F2I: test convert float to int
func TestF2iPositive(t *testing.T) {
	f := newFrame(opcodes.F2I)
	push(&f, 2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != 2 {
		t.Errorf("F2I: expected a result of 2, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("F2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

func TestF2iNegative(t *testing.T) {
	f := newFrame(opcodes.F2I)
	push(&f, -2.9)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != -2 {
		t.Errorf("F2I: expected a result of 2, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("F2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// F2L: test convert float to long
func TestF2l(t *testing.T) {
	f := newFrame(opcodes.F2L)
	push(&f, 2.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != 2 {
		t.Errorf("F2L: expected a result of 2.0, but got: %d", val)
	}
	if f.TOS != 0 {
		t.Errorf("F2L: Expected stack with 1 item, but got a TOS of: %d", f.TOS)
	}
}

// FADD: Add two floats
func TestFadd(t *testing.T) {
	f := newFrame(opcodes.FADD)
	push(&f, 2.1)
	push(&f, 3.1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(float64)
	if math.Abs(value-5.2) > maxFloatDiff {
		t.Errorf("FADD: expected a result of 5.2, but got: %f", value)
	}
	if f.TOS != -1 {
		t.Errorf("FADD: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// FCMPG: compare two floats
func TestFcmpg1(t *testing.T) {
	f := newFrame(opcodes.FCMPG)
	push(&f, 3.0)
	push(&f, 2.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 1 {
		t.Errorf("FCMPG: Expected value to be 1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("FCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// FCMPG: compare two floats
func TestFcmpgMinus1(t *testing.T) {
	f := newFrame(opcodes.FCMPG)
	push(&f, 2.0)
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != -1 {
		t.Errorf("FCMPG: Expected value to be -1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("FCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// FCMPG: compare two floats
func TestFcmpg0(t *testing.T) {
	f := newFrame(opcodes.FCMPG)
	push(&f, 3.0)
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 0 {
		t.Errorf("FCMPG: Expected value to be 0, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("FCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// FCMPG
func TestFcmpgNan(t *testing.T) {
	f := newFrame(opcodes.FCMPG)
	push(&f, math.NaN())
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 1 {
		t.Errorf("FCMPG: Expected value to be 1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("FCMPG: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// FCMPL
func TestFcmplNan(t *testing.T) {
	f := newFrame(opcodes.FCMPL)
	push(&f, math.NaN())
	push(&f, 3.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != -1 {
		t.Errorf("FCMPL: Expected value to be -1, got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("FCMPL: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// FCONST_0
func TestFconst0(t *testing.T) {
	f := newFrame(opcodes.FCONST_0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 0.0 {
		t.Errorf("FCONST_0: Expected popped value to be 0.0, got: %f", value)
	}
}

// FCONST_1
func TestFconst1(t *testing.T) {
	f := newFrame(opcodes.FCONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 1.0 {
		t.Errorf("FCONST_1: Expected popped value to be 1.0, got: %f", value)
	}
}

// FCONST_2
func TestFconst2(t *testing.T) {
	f := newFrame(opcodes.FCONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 2.0 {
		t.Errorf("FCONST_2: Expected popped value to be 2.0, got: %f", value)
	}
}

// FDIV: float divide of.TOS-1 by tos, push result
func TestFdiv(t *testing.T) {
	f := newFrame(opcodes.FDIV)
	push(&f, 3.0)
	push(&f, 2.0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(float64)
	if value != 1.5 {
		t.Errorf("FDIV: expected a result of 1.5, but got: %f", value)
	}
}

// FDIV: with divide zero by zero, should = NaN
func TestFdivDivideZeroByZero(t *testing.T) {
	f := newFrame(opcodes.FDIV)
	push(&f, float64(0))
	push(&f, float64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	ret := pop(&f)

	if !math.IsNaN(ret.(float64)) {
		t.Errorf("FDIV: Did not get an expected NaN")
	}
}

// FDIV: with divide positive number by zero, should = +Inf
func TestFdivDividePosNumberByZero(t *testing.T) {
	f := newFrame(opcodes.FDIV)
	push(&f, float64(10))
	push(&f, float64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	ret := pop(&f)

	if !math.IsInf(ret.(float64), 1) {
		t.Errorf("FDIV: Did not get an expected +Infinity")
	}
}

// FLOAD: test load of float in locals[index] on to stack
func TestFload(t *testing.T) {
	f := newFrame(opcodes.FLOAD)
	f.Meth = append(f.Meth, 0x04) // use local var #4
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, float64(0x1234562)) // put value in locals[4]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(float64)
	if x != float64(0x1234562) {
		t.Errorf("FLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("FLOAD: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
	if f.PC != 2 {
		t.Errorf("FLOAD: Expected pc to be pointing at byte 2, got: %d", f.PC)
	}
}

// FLOAD_0: load of float in locals[0] onto stack
func TestFload0(t *testing.T) {
	f := newFrame(opcodes.FLOAD_0)
	f.Locals = append(f.Locals, 1.2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 1.2 {
		t.Errorf("FLOAD_0: Expected popped value to be 1.2, got: %f", value)
	}
}

// FLOAD_1: load of float in locals[1] onto stack
func TestFload1(t *testing.T) {
	f := newFrame(opcodes.FLOAD_1)
	f.Locals = append(f.Locals, 1.1)
	f.Locals = append(f.Locals, 1.2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 1.2 {
		t.Errorf("FLOAD_1: Expected popped value to be 1.2, got: %f", value)
	}
}

// FLOAD_2: load of float in locals[2] onto stack
func TestFload2(t *testing.T) {
	f := newFrame(opcodes.FLOAD_2)
	f.Locals = append(f.Locals, 1.1)
	f.Locals = append(f.Locals, 1.1)
	f.Locals = append(f.Locals, 1.2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 1.2 {
		t.Errorf("FLOAD_2: Expected popped value to be 1.2, got: %f", value)
	}
}

// FLOAD_3: load of fload in locals[3] onto stack
func TestFload3(t *testing.T) {
	f := newFrame(opcodes.FLOAD_3)
	f.Locals = append(f.Locals, 1.1)
	f.Locals = append(f.Locals, 1.1)
	f.Locals = append(f.Locals, 1.1)
	f.Locals = append(f.Locals, 1.2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 1.2 {
		t.Errorf("FLOAD_3: Expected popped value to be 1.2, got: %f", value)
	}
}

// FMUL (pop 2 floats, multiply them, push result)
func TestFmul(t *testing.T) {
	f := newFrame(opcodes.FMUL)
	push(&f, 1.5)
	push(&f, 2.0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("FMUL, Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 3.0 {
		t.Errorf("FMUL: Expected popped value to be 3.0, got: %f", value)
	}
}

// FNEG: negate a float
func TestFneg(t *testing.T) {
	f := newFrame(opcodes.FNEG)
	push(&f, 10.0)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 {
		t.Errorf("FNEG, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(float64)
	if value != -10.0 {
		t.Errorf("FNEG: Expected popped value to be -10.0, got: %f", value)
	}
}

// FREM: remainder of float division (the % operator)
func TestFrem(t *testing.T) {
	f := newFrame(opcodes.FREM)
	push(&f, 23.5)

	push(&f, 3.3)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 {
		t.Errorf("FREM, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(float64)
	if math.Abs(value-0.40000033) > maxFloatDiff {
		t.Errorf("FREM: Expected popped value to be 0.40000033, got: %f", value)
	}
}

// FSTORE: Store float from stack into local specified by following byte.
func TestFstore(t *testing.T) {
	f := newFrame(opcodes.FSTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	f.Locals = append(f.Locals, zerof)
	push(&f, float64(0x22223))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2] != float64(0x22223) {
		t.Errorf("FSTORE: Expecting 0x22223 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("FSTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// FSTORE_0: Store float from stack into localVar[0]
func TestFstore0(t *testing.T) {
	f := newFrame(opcodes.FSTORE_0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[0].(float64) != 1.0 {
		t.Errorf("FSTORE_0: expected lcoals[0] to be 1.0, got: %f", f.Locals[0].(float64))
	}
	if f.TOS != -1 {
		t.Errorf("FSTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// FSTORE_1: Store float from stack into localVar[0]
func TestFstore1(t *testing.T) {
	f := newFrame(opcodes.FSTORE_1)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[1].(float64) != 1.0 {
		t.Errorf("FSTORE_1: expected lcoals[1] to be 1.0, got: %f", f.Locals[1].(float64))
	}
	if f.TOS != -1 {
		t.Errorf("FSTORE_1: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// FSTORE_2: Store float from stack into localVar[2]
func TestFstore2(t *testing.T) {
	f := newFrame(opcodes.FSTORE_2)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[2].(float64) != 1.0 {
		t.Errorf("FSTORE_2: expected lcoals[2] to be 1.0, got: %f", f.Locals[2].(float64))
	}
	if f.TOS != -1 {
		t.Errorf("FSTORE_2: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// FSTORE_3: Store float from stack into localVar[3]
func TestFstore3(t *testing.T) {
	f := newFrame(opcodes.FSTORE_3)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	f.Locals = append(f.Locals, 0.0)
	push(&f, 1.0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[3].(float64) != 1.0 {
		t.Errorf("FSTORE_3: expected lcoals[3] to be 1.0, got: %f", f.Locals[3].(float64))
	}
	if f.TOS != -1 {
		t.Errorf("FSTORE_3: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// FSUB:float subtraction
func TestFsub(t *testing.T) {
	f := newFrame(opcodes.FSUB)
	push(&f, 1.0)
	push(&f, 0.7)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(float64)

	if math.Abs(value-0.3) > maxFloatDiff {
		t.Errorf("FSUB: Expected popped value to be 0.3, got: %f", value)
	}

	if f.TOS != -1 {
		t.Errorf("DSUB, Empty stack expected, got: %d", f.TOS)
	}
}

// GETFIELD: Get a field from an object
func TestGetField(t *testing.T) {
	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: 9, Slot: 0} // point to a fieldRef
	CP.FieldRefs = make([]classloader.FieldRefEntry, 1, 1)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 0, NameAndType: 0}
	f.CP = &CP

	// push the string whose field[0] we'll be getting
	str := object.NewString()
	str.Fields[0].Fvalue = "hello"
	push(&f, str)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// preceding should mean that the field value is on the stack
	ret := pop(&f)
	if ret != "hello" {
		t.Errorf("GETFIELD: did not get expected pointer to a string 'hello'")
	}

	if f.TOS != -1 {
		t.Errorf("GETFIELD: Expected an empty op stack, got TOS: %d", f.TOS)
	}
}

// GETFIELD: Get a long field, make sure that it's value is pushed twice
func TestGetFieldWithLong(t *testing.T) {
	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: 9, Slot: 0} // point to a fieldRef
	CP.FieldRefs = make([]classloader.FieldRefEntry, 1, 1)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 0, NameAndType: 0}
	f.CP = &CP

	// push the string whose field[0] we'll be getting
	obj := object.MakeEmptyObject()
	obj.Fields = make([]object.Field, 1, 1)
	obj.Fields[0].Fvalue = int64(222)
	obj.Fields[0].Ftype = types.Long
	push(&f, obj)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// preceding should mean that the field value is on the stack
	ret := pop(&f).(int64)
	if ret != 222 {
		t.Errorf("GETFIELD: expected popped value of 222, got: %d", ret)
	}

	if f.TOS != 0 {
		t.Errorf("GETFIELD: Expected 1 remaining value op stack, got TOS: %d", f.TOS)
	}
}

// GETFIELD: Get a field from an object (here, with error that it's not a fieldref)
func TestGetFieldInvalidFieldEntry(t *testing.T) {
	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	// pointing to the next CP entry, which s/be a FieldRef but is a UTF8 record
	CP.CpIndex[0] = classloader.CpEntry{Type: 1, Slot: 0}
	f.CP = &CP
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	ret := runFrame(fs)
	if !strings.Contains(ret.Error(), "Expected a field ref, but got") {
		t.Errorf("GETFIELD: Expected a different error, got: %s",
			ret.Error())
	}
}

// GETSTATIC: Get a static field's value (here, with error that it's not a fieldref)
func TestGetStaticInvalidFieldEntry(t *testing.T) {
	f := newFrame(opcodes.GETSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	// pointing to the next CP entry, which s/be a FieldRef but is a UTF8 record
	CP.CpIndex[0] = classloader.CpEntry{Type: 1, Slot: 0}
	f.CP = &CP
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	ret := runFrame(fs)
	if !strings.Contains(ret.Error(), "Expected a field ref, but got") {
		t.Errorf("GETFIELD: Expected a different error, got: %s",
			ret.Error())
	}
}

// GETSTATIC: Get a static field's value (here, a boolean in the String class, set to true)
func TestGetStaticBoolean(t *testing.T) {
	f := newFrame(opcodes.GETSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	statics.PreloadStatics() // load the statics table with the String class

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 2}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}

	CP.FieldRefs = make([]classloader.FieldRefEntry, 2, 2)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 2, NameAndType: 4}

	CP.Utf8Refs = make([]string, 5, 5)
	CP.Utf8Refs[0] = "java/lang/String"
	CP.Utf8Refs[1] = "COMPACT_STRINGS"

	CP.ClassRefs = make([]uint16, 5, 5)
	CP.ClassRefs[0] = 2 // point to CpIndex[2] -- need to validate this is right

	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 5, 5)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 5, // field name as UTF8 entry, here the CPindex index
		DescIndex: 0,
	}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	ret := runFrame(fs)
	if ret != nil {
		t.Errorf("GETSTATIC: Expected a different error, got: %s",
			ret.Error())
	}

	retVal := pop(&f).(int64)
	if retVal != 1 {
		t.Errorf("GETSTATIC: Expected a return of 1 (true) for a boolean, got: %d", retVal)
	}
}

// GOTO: in forward direction (to a later bytecode)
func TestGotoForward(t *testing.T) {
	f := newFrame(opcodes.GOTO)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x03)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.NOP)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN {
		t.Errorf("GOTO forward: Expected pc to point to RETURN, but instead it points to : %s", opcodes.BytecodeNames[f.Meth[f.PC]])
	}
}

// GOTO: go to instruction in backward direction (to an earlier bytecode)
func TestGotoBackward(t *testing.T) {
	f := newFrame(opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.GOTO)
	f.Meth = append(f.Meth, 0xFF) // should be -1
	f.Meth = append(f.Meth, 0xFF)
	f.Meth = append(f.Meth, opcodes.BIPUSH)
	f.PC = 1 // skip over the return instruction to start, catch it on the backward goto
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN {
		t.Errorf("GOTO backeard Expected pc to point to RETURN, but instead it points to : %s", opcodes.BytecodeNames[f.Meth[f.PC]])
	}
}

// I2B: convert int to Java char (16-bit value)
func TestI2B(t *testing.T) {
	f := newFrame(opcodes.I2B)
	push(&f, int64(2100))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != 52 {
		t.Errorf("I2B: expected a result of 52, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("I2B: Expected stack with 1 entry, but got a TOS of: %d", f.TOS)
	}
}

// I2B: convert int to Java char (16-bit value) using a negative value
func TestI2Bneg(t *testing.T) { // TODO: check that this matches Java result
	f := newFrame(opcodes.I2B)
	push(&f, int64(-2100))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != -204 { // looks like 256-52
		t.Errorf("I2B: expected a result of -204, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("I2B: Expected stack with 1 entry, but got a TOS of: %d", f.TOS)
	}
}

// I2C: convert int to Java char (16-bit value)
func TestI2C(t *testing.T) {
	f := newFrame(opcodes.I2C)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != 21 {
		t.Errorf("I2C: expected a result of 21, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("I2C: Expected stack with 1 entry, but got a TOS of: %d", f.TOS)
	}
}

// I2D: Convert int to double
// Note that while ints are stored in one opStack slot,
// doubles use two slots.
func TestI2D(t *testing.T) {
	f := newFrame(opcodes.I2D)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(float64)
	if value != 21.0 {
		t.Errorf("I2D: expected a result of 21.0, but got: %f", value)
	}
	if f.TOS != 0 {
		t.Errorf("I2D: Expected stack with 1 entry, but got a TOS of: %d", f.TOS)
	}
}

// I2F: convert int to short
func TestI2f(t *testing.T) {
	f := newFrame(opcodes.I2F)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(float64)
	if value != 21.0 {
		t.Errorf("I2F: expected a result of 21.0, but got: %f", value)
	}
	if f.TOS != -1 {
		t.Errorf("I2F: Expected stack with 0 entry, but got a TOS of: %d", f.TOS)
	}
}

// I2L: Convert int to long
// Note that since ints in Jacobin are int64--which is the same size as a long--
// so no conversion takes place. However, while ints are stored in one opStack
// slot, longs use two slots. So this is the primary test here.
func TestI2l(t *testing.T) {
	f := newFrame(opcodes.I2L)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != 21 {
		t.Errorf("I2L: expected a result of 21, but got: %d", value)
	}
	if f.TOS != 0 {
		t.Errorf("I2L: Expected stack with 1 entry, but got a TOS of: %d", f.TOS)
	}
}

// I2S: convert int to short
func TestI2s(t *testing.T) {
	f := newFrame(opcodes.I2S)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != 21 {
		t.Errorf("I2S: expected a result of 21, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("I2S: Expected stack with 0 entry, but got a TOS of: %d", f.TOS)
	}
}

// IADD: Add two integers
func TestIadd(t *testing.T) {
	f := newFrame(opcodes.IADD)
	push(&f, int64(21))
	push(&f, int64(22))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != 43 {
		t.Errorf("IADD: expected a result of 43, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IADD: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// IAND: Logical and of two ints, push result
func TestIand(t *testing.T) {
	f := newFrame(opcodes.IAND)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 20 { // 21 & 22 = 20
		t.Errorf("IAND: expected a result of 20, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IAND: Expected an empty stack, but got a TOS of: %d", f.TOS)
	}
}

// IDIV: integer divide of.TOS-1 by tos, push result
func TestIdiv(t *testing.T) {
	f := newFrame(opcodes.IDIV)
	push(&f, int64(220))
	push(&f, int64(22))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f).(int64)
	if value != 10 {
		t.Errorf("IDIV: expected a result of 10, but got: %d", value)
	}
}

// IDIV: make sure that divide by zero generates an Arithmetic Exception and
// displays an error message.
/* TEMPORARILY DISABLED because of reworking how exceptions are thrown.
func TestIdivDivideByZero(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	log.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// Init classloader and load base classes
	err := classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Error initiating environment for IDIV test")
	}
	classloader.LoadBaseClasses() // load base classes

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadNatives(&classloader.MTable) // load native classes

	f := newFrame(opcodes.IDIV)
	f.ClName = "testClass"
	f.MethName = "testMethod"
	var CP = classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.Utf8Refs = []string{"java/lang/ArithmeticException"}
	CP.ClassRefs = []uint16{0}
	f.CP = &CP

	push(&f, int64(220))
	push(&f, int64(0))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame

	// need to create a thread to catch the exception
	thread := thread.CreateThread()
	thread.Stack = fs
	thread.AddThreadToTable(g)
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "java.lang.ArithmeticException") {
		t.Errorf("IDIV: Did not get expected error msg, got: %s", errMsg)
	}
} */

// ICONST_M1:
func TestIconstN1(t *testing.T) {
	f := newFrame(opcodes.ICONST_M1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	var value = pop(&f).(int64)
	if value != -1 {
		t.Errorf("ICONST_M1: Expected popped value to be -1, got: %d", value)
	}
}

// ICONST_0
func TestIconst0(t *testing.T) {
	f := newFrame(opcodes.ICONST_0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("ICONST_0: Expected popped value to be 0, got: %d", value)
	}
}

// ICONST_1
func TestIconst1(t *testing.T) {
	f := newFrame(opcodes.ICONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 1 {
		t.Errorf("ICONST_1: Expected popped value to be 1, got: %d", value)
	}
}

// ICONST_2
func TestIconst2(t *testing.T) {
	f := newFrame(opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 2 {
		t.Errorf("ICONST_2: Expected popped value to be 2, got: %d", value)
	}
}

// ICONST_3
func TestIconst3(t *testing.T) {
	f := newFrame(opcodes.ICONST_3)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 3 {
		t.Errorf("ICONST_3: Expected popped value to be 3, got: %d", value)
	}
}

// ICONST_4
func TestIconst4(t *testing.T) {
	f := newFrame(opcodes.ICONST_4)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 4 {
		t.Errorf("ICONST_4: Expected popped value to be 4, got: %d", value)
	}
}

// ICONST_5:
func TestIconst5(t *testing.T) {
	f := newFrame(opcodes.ICONST_5)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 5 {
		t.Errorf("ICONST_5: Expected popped value to be 5, got: %d", value)
	}
}

// IF_ACMPEQ: jump if two addresses are equal
func TestIfAcmpEq(t *testing.T) {
	f := newFrame(opcodes.IF_ACMPEQ)
	push(&f, int64(0xFF8899))
	push(&f, int64(0xFF8899))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IF_ACMPEQ: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ACMPEQ: jump if two addresses are equal (this tests addresses being unequal)
func TestIfAcmpeqFail(t *testing.T) {
	f := newFrame(opcodes.IF_ACMPEQ)
	push(&f, int64(0xFF8899))
	push(&f, int64(0xFF889A))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST_2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN) // the failed test should drop to this
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("IF_ICMPEQ: expecting fall-through to RETURN instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ACMPNE: jump if two addresses are not equal
func TestIfAcmpNe(t *testing.T) {
	f := newFrame(opcodes.IF_ACMPNE)
	push(&f, int64(0xFF8899))
	push(&f, int64(0xFF889A))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IF_ACMPNE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ACMPNE: jump if two addresses are equal (this tests addresses being equal)
func TestIfAcmpneFail(t *testing.T) {
	f := newFrame(opcodes.IF_ACMPNE)
	push(&f, int64(0xFF8899))
	push(&f, int64(0xFF8899))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST_2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN) // the failed test should drop to this
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("IF_ICMPNE: expecting fall-through to RETURN instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPEQ: jump if val1 == val2 (both ints, both popped off stack)
func TestIfIcmpeq(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPEQ)
	push(&f, int64(9)) // pushed two equal values, so jump should be made.
	push(&f, int64(9))

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPEQ: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPEQ: jump if val1 == val2; here test with unequal value
func TestIfIcmpeqUnequal(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPEQ)
	push(&f, int64(9)) // pushed two unequal values, so no jump should be made.
	push(&f, int64(-9))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	if err != nil {
		t.Errorf("IF_ICMPEQ: Got unexpected error: %s", err.Error())
	}

	if f.PC != 3 { // 2 for the jump due to inequality above, +1 for fetch of next bytecode
		t.Errorf("IF_ICMPEQ: Expected PC to be 2, got %d", f.PC)
	}

	if f.TOS != -1 { // stack should be empty
		t.Errorf("IF_CIMPEQ: Expected an empty stack, got TOS of: %d", f.TOS)
	}
}

// IF_CMPGE: if integer compare val 1 >= val 2. Here test for = (next test for >)
func TestIfIcmpge1(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPGE)
	push(&f, int64(9))
	push(&f, int64(9))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPGE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPGE: if integer compare val 1 >= val 2. Here test for > (previous test for =)
func TestIfIcmpge2(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPGE)
	push(&f, int64(9))
	push(&f, int64(8))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPGE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPGE: if integer compare val 1 >= val 2 //test when condition fails
func TestIfIcmgetFail(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPGE)
	push(&f, int64(8))
	push(&f, int64(9))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN) // the failed test should drop to this
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("ICMPGE: expecting fall-through to RETURN instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPGE: if integer compare val 1 >= val 2. Here test for > (previous test for =)
func TestIfIcmple2(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPLE)
	push(&f, int64(8))
	push(&f, int64(9))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IF_ICMPLE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPGT: jump if val1 > val2 (both ints, both popped off stack)
func TestIfIcmpgt(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPGT)
	push(&f, int64(9)) // val1 > val2, so jump should be made.
	push(&f, int64(8))

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IF_ICMPNE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPLE: if integer compare val 1 ! <= val 2 //test when condition fails
func TestIfIcmpletFail(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPLE)
	push(&f, int64(9))
	push(&f, int64(8))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN) // the failed test should drop to this
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("IF_ICMPLE: expecting fall-through to RETURN instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPLE: if integer compare val 1 <= val 2. Here testing for =
func TestIfIcmple1(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPLE)
	push(&f, int64(9))
	push(&f, int64(9))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPLE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPLT: if integer compare val 1 < val 2
func TestIfIcmplt(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPLT)
	push(&f, int64(8))
	push(&f, int64(9))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.ICONST_1)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPLT: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPLT: if integer compare val 1 < val 2 //test when condition fails
func TestIfIcmpltFail(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPLT)
	push(&f, int64(9))
	push(&f, int64(9))
	// note that the byte passed in newframe() is at f.Meth[0]
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN) // the failed test should drop to this
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC] != opcodes.RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("ICMPLT: expecting fall-through to RETURN instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPNE: jump if val1 != val2 (both ints, both popped off stack)
func TestIfIcmpne(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPNE)
	push(&f, int64(9)) // pushed two unequal values, so jump should be made.
	push(&f, int64(8))

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IF_ICMPNE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IF_ICMPNE: jump if val1 != val2 Here tests when they are equal
func TestIfIcmpneAreEqual(t *testing.T) {
	f := newFrame(opcodes.IF_ICMPNE)
	push(&f, int64(9)) // pushed two equal values, so jump should not be made.
	push(&f, int64(9))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	if err != nil {
		t.Errorf("IF_ICMPNE: Got unexpected error: %s", err.Error())
	}
	if f.PC != 3 { // PC+= 2 when test fails, +1 for the next bytecode
		t.Errorf("IF_ICMPNE: PC to be 3, got: %d", f.PC)
	}
	if f.TOS != -1 {
		t.Errorf("IF_ICMPNE: Expected stack to be empty, TOS was: %d", f.TOS)
	}
}

// IFEQ: jump if int popped off TOS is = 0
func TestIfeq(t *testing.T) {
	f := newFrame(opcodes.IFEQ)
	push(&f, int64(0)) // pushed 0, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFEQ: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFEQ: jump if int popped off TOS is = 0; here != 0
func TestIfeqFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFEQ)
	push(&f, int64(23)) // pushed 23, so jump should not be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFEQ: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFGE: jump if int popped off TOS is >= 0
func TestIfge(t *testing.T) {
	f := newFrame(opcodes.IFGE)
	push(&f, int64(66)) // pushed 66, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFGE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFGE: jump if int popped off TOS is >= 0, here = 0
func TestIfgeEqual0(t *testing.T) {
	f := newFrame(opcodes.IFGE)
	push(&f, int64(0)) // pushed 0, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFGE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFGE: jump if int popped off TOS is >= 0; here < 0
func TestIfgeFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFGE)
	push(&f, int64(-1)) // pushed -1, so jump should not be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFGE: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFGT: jump if int popped off TOS is > 0
func TestIfgt(t *testing.T) {
	f := newFrame(opcodes.IFGT)
	push(&f, int64(66)) // pushed 66, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFGT: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFGT: jump if int popped off TOS is > 0; here = 0
func TestIfgtFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFGT)
	push(&f, int64(0)) // pushed 0, so jump should not be made.

	f.Meth = append(f.Meth, 0)
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFGT: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFLE: jump if int popped off TOS is <= 0
func TestIfle(t *testing.T) {
	f := newFrame(opcodes.IFLE)
	push(&f, int64(-66)) // pushed -66, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFLE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFLE: jump if int popped off TOS is <= 0; here = 0
func TestIfleTest0(t *testing.T) {
	f := newFrame(opcodes.IFLE)
	push(&f, int64(0)) // pushed 0, so jump should be made.

	f.Meth = append(f.Meth, 0)
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFLE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFLE: jump if int popped off TOS is <= 0; here > 0, so no jump
func TestIfleFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFLE)
	push(&f, int64(66)) // pushed 66, so jump should not be made.

	f.Meth = append(f.Meth, 0)
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFLE: Invalid jump when expecting fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFLT: jump if int popped off TOS is < 0
func TestIflt(t *testing.T) {
	f := newFrame(opcodes.IFLT)
	push(&f, int64(-66)) // pushed -66, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFLT: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFLT: jump if int popped off TOS is < 0; here = 0
func TestIfltFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFLT)
	push(&f, int64(0)) // pushed 0, so jump should not be made.

	f.Meth = append(f.Meth, 0)
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFLT: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFNE: jump if int popped off TOS is != 0
func TestIfne(t *testing.T) {
	f := newFrame(opcodes.IFNE)
	push(&f, int64(1)) // pushed 1, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFNE: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFNE: jump if int popped off TOS is != 0; here it is = 0
func TestIfneFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFNE)
	push(&f, int64(0)) // pushed 0, so jump should not be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFNE: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFNONNULL: jump if TOS holds a non-null address
func TestIfn0nnull(t *testing.T) {
	f := newFrame(opcodes.IFNONNULL)
	o := object.NewString()
	push(&f, o) // pushed a valid address, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFNONNULL: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFNONNULL: jump if TOS holds a non-null address; here it is null
func TestIfnonnullFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFNONNULL)
	var oAddr *object.Object
	oAddr = nil
	push(&f, oAddr)
	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFNONNULL: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFNULL: jump if TOS holds null address
func TestIfnull(t *testing.T) {
	f := newFrame(opcodes.IFNULL)
	var oAddr *object.Object
	oAddr = nil     // note either nil or object.Null will give same result
	push(&f, oAddr) // pushed null, so jump should be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.NOP)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] != opcodes.ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFNULL: expecting a jump to ICONST_2 instuction, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

// IFNULL: jump if TOS address is null; here not null
func TestIfnullFallThrough(t *testing.T) {
	f := newFrame(opcodes.IFNULL)
	o := object.MakeEmptyObject()
	push(&f, o) // pushed non-null address, so jump should not be made.

	f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.Meth = append(f.Meth, 4)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.ICONST_2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Meth[f.PC-1] == opcodes.IFNULL { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("IFNULL: Invalid fall-through, got: %s",
			opcodes.BytecodeNames[f.PC])
	}
}

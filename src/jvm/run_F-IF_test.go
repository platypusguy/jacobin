/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"math"
	"strings"
	"testing"
)

// These tests test the individual bytecode instructions. They are presented
// here in alphabetical order of the instruction name.
// THIS FILE CONTAINS TESTS FOR ALL BYTECODES FROM F2D THROUGH IFNULL.
// All other bytecodes are in run_*_test.go files except
// for array bytecodes, which are located in arrays_test.go

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
	CP.CpIndex[1] = classloader.CpEntry{Type: 9, Slot: 0} // point to fieldRef[0]

	CP.FieldRefs = make([]classloader.FieldRefEntry, 1, 1)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 0, NameAndType: 0}

	CP.ClassRefs = make([]uint32, 1, 1)
	CP.ClassRefs[0] = 0 // classRefs are not used to access a field

	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 1, 1)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{
		NameIndex: 0, // UTF8: "value"
		DescIndex: 1} // UTF8: "Ljava/lang/String;"

	CP.Utf8Refs = make([]string, 2, 2)
	CP.Utf8Refs[0] = "value"
	CP.Utf8Refs[1] = "Ljava/lang/String;"
	f.CP = &CP

	// push the string whose field[0] we'll be getting
	str := object.NewStringObject()
	str.FieldTable["value"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: "hello",
	}

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
	globals.InitGlobals("test")

	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: 9, Slot: 0} // point to a fieldRef

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.FieldRefEntry, 1, 1)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 0, NameAndType: 0}
	f.CP = &CP

	// now create the NameAndType records
	CP.NameAndTypes = make([]classloader.NameAndTypeEntry, 1, 1)
	CP.NameAndTypes[0] = classloader.NameAndTypeEntry{NameIndex: 0, DescIndex: 1}

	// and finally the UTF8 records pointed to by the NameAndType entry above
	CP.Utf8Refs = make([]string, 2)
	CP.Utf8Refs[0] = "value"
	CP.Utf8Refs[1] = types.Long
	f.CP = &CP

	// push the string whose field[0] we'll be getting
	obj := object.MakePrimitiveObject("java/lang/Long", types.Long, int64(222))
	// obj.DumpObject("TestGetFieldWithLong", 0)
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
	globals.InitGlobals("test")
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
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}

	CP.FieldRefs = make([]classloader.FieldRefEntry, 2, 2)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 2, NameAndType: 4}

	CP.Utf8Refs = make([]string, 5, 5)
	// class = "java/lang/String" in string pool
	classNameIndex := stringPool.GetStringIndex(&object.StringClassName)
	CP.Utf8Refs[0] = "COMPACT_STRINGS"

	CP.ClassRefs = make([]uint32, 5, 5)
	CP.ClassRefs[0] = classNameIndex
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

// IDIV: Testing the exception is done in TestHexIDIVexception.go

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
func TestIfnonnull(t *testing.T) {
	f := newFrame(opcodes.IFNONNULL)
	o := object.NewStringObject()
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

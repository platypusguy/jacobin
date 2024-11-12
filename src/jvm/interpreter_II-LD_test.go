/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"strings"
	"testing"
)

// These tests test the individual bytecode instructions. They are presented
// here in alphabetical order of the instruction name.
// THIS FILE CONTAINS TESTS FOR ALL BYTECODES FROM IINC to LDIV.
// All other bytecodes are in run_*_test.go files except
// for array bytecodes, which are located in arrayBytecodes_test.go

// IINC: increment local variable
func TestNewIinc(t *testing.T) {
	f := newFrame(opcodes.IINC)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(10)) // initialize local variable[1] to 10
	f.Meth = append(f.Meth, 1)             // increment local variable[1]
	f.Meth = append(f.Meth, 27)            // increment it by 27
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}
	value := f.Locals[1]
	if value != int64(37) {
		t.Errorf("IINC: Expected popped value to be 37, got: %d", value)
	}
}

// IINC: increment local variable by negative value
func TestNewIincNeg(t *testing.T) {
	f := newFrame(opcodes.IINC)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(10)) // initialize local variable[1] to 10
	f.Meth = append(f.Meth, 1)             // increment local variable[1]
	val := -27
	f.Meth = append(f.Meth, byte(val)) // "increment" it by -27
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}
	value := f.Locals[1]
	if value != int64(-17) {
		t.Errorf("IINC: Expected popped value to be -17, got: %d", value)
	}
}

// ILOAD: test load of int in locals[index] on to stack
func TestNewIload(t *testing.T) {
	f := newFrame(opcodes.ILOAD)
	f.Meth = append(f.Meth, 0x04) // use local var #4
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x1234562)) // put value in locals[4]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	x := pop(&f).(int64)
	if x != 0x1234562 {
		t.Errorf("ILOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("ILOAD: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
	if f.PC != 2 {
		t.Errorf("ILOAD: Expected pc to be pointing at byte 2, got: %d", f.PC)
	}
}

// ILOAD_0: load of int in locals[0] onto stack
func TestNewIload0(t *testing.T) {
	f := newFrame(opcodes.ILOAD_0)
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != int64(27) {
		t.Errorf("ILOAD_0: Expected popped value to be 27, got: %d", value)
	}
}

// ILOAD_1: load of int in locals[1] onto stack
func TestNewIload1(t *testing.T) {
	f := newFrame(opcodes.ILOAD_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 27 {
		t.Errorf("ILOAD_1: Expected popped value to be 27, got: %d", value)
	}
}

// ILOAD_2: load of int in locals[2] onto stack
func TestNewIload2(t *testing.T) {
	f := newFrame(opcodes.ILOAD_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(1))
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != int64(27) {
		t.Errorf("ILOAD_2: Expected popped value to be 27, got: %d", value)
	}
}

// ILOAD_3: load of int in locals[3] onto stack
func TestNewIload3(t *testing.T) {
	f := newFrame(opcodes.ILOAD_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(1))
	f.Locals = append(f.Locals, int64(2))
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 27 {
		t.Errorf("ILOAD_3: Expected popped value to be 27, got: %d", value)
	}
}

// Test IMUL (pop 2 values, multiply them, push result)
func TestNewImul(t *testing.T) {
	f := newFrame(opcodes.IMUL)
	push(&f, int64(10))
	push(&f, int64(7))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("IMUL, Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 70 {
		t.Errorf("IMUL: Expected popped value to be 70, got: %d", value)
	}
}

// INEG: negate an int
func TestNewIneg(t *testing.T) {
	f := newFrame(opcodes.INEG)
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("INEG, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != -10 {
		t.Errorf("INEG: Expected popped value to be -10, got: %d", value)
	}
}

// INSTANCEOF: Is the TOS item an instance of a particular class?
func TestNewInstanceofNilAndNull(t *testing.T) {
	f := newFrame(opcodes.INSTANCEOF)
	push(&f, nil)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("INSTANCEOF: Expected nil to return a 0, got %d", value)
	}

	f = newFrame(opcodes.INSTANCEOF)
	push(&f, object.Null)

	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value = pop(&f).(int64)
	if value != 0 {
		t.Errorf("INSTANCEOF: Expected null to return a 0, got %d", value)
	}
}

// INSTANCEOF for a string
func TestNewInstanceofString(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	trace.Init()

	_ = classloader.Init()
	// classloader.LoadBaseClasses()
	classloader.MethAreaInsert(types.StringClassName,
		&(classloader.Klass{
			Status: 'X', // use a status that's not subsequently tested for.
			Loader: "bootstrap",
			Data:   nil,
		}))
	s := object.StringObjectFromGoString("hello world")

	f := newFrame(opcodes.INSTANCEOF)
	f.Meth = append(f.Meth, 0) // point to entry [1] in CP
	f.Meth = append(f.Meth, 1) // " "

	// now create the CP.
	// [0] First entry is perforce 0
	// [1] is a ClassRef that points to string pool entry for java/lang/String
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.ClassRefs = append(CP.ClassRefs, types.StringPoolStringIndex) // point to string pool entry for java/lang/String
	f.CP = &CP

	push(&f, s)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)
	if value != 1 { // a 1 = it's a match between class and object
		t.Errorf("INSTANCEOF: Expected string to return a 1, got %d", value)
	}
}

// IOR: Logical OR of two ints
func TestNewIor(t *testing.T) {
	f := newFrame(opcodes.IOR)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)

	if value != 23 { // 21 | 22 = 23
		t.Errorf("IOR: expected a result of 23, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// IREM: int modulo
func TestNewIrem(t *testing.T) {
	f := newFrame(opcodes.IREM)
	push(&f, int64(74))
	push(&f, int64(6))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 { // product is pushed twice b/c it's a long, which occupies 2 slots
		t.Errorf("IREM, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != 2 {
		t.Errorf("IREM: Expected result to be 2, got: %d", value)
	}
}

// IREM: int modulo -- divide by zero
// Because this test requires a full class set up due to IREM now throwing a full exception,
// the test code has been moved to ThrowIREMexception.go in wholeClassTests.

// IRETURN: push an int on to the op stack of the calling method and exit the present method/frame
func TestNewIreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, int64(20))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f0)
	f1 := newFrame(opcodes.IRETURN)
	push(&f1, int64(21))
	fs.PushFront(&f1)
	interpret(fs)

	f3 := fs.Front().Value.(*frames.Frame)
	newVal := pop(f3).(int64)
	if newVal != 21 {
		t.Errorf("After IRETURN, expected a value of 21 in previous frame, got: %d", newVal)
	}
	prevVal := pop(f3).(int64)
	if prevVal != 20 {
		t.Errorf("After IRETURN, expected a value of 20 in 2nd place of previous frame, got: %d", prevVal)
	}
}

// ISHL: Left shift of long
func TestNewIshl(t *testing.T) {
	f := newFrame(opcodes.ISHL)
	push(&f, int64(22)) // longs require two slots, so pushed twice
	push(&f, int64(3))  // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 176 { // 22 << 3 = 176
		t.Errorf("ISHL: expected a result of 176, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("ISHL: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// ISHR: Right shift of int
func TestNewIshr(t *testing.T) {
	f := newFrame(opcodes.ISHR)
	push(&f, int64(200))
	push(&f, int64(3)) // shift right 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("ISHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("ISHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// ISHR: Right shift of negative int
func TestNewIshrNeg(t *testing.T) {
	f := newFrame(opcodes.ISHR)
	push(&f, int64(-200))
	push(&f, int64(3)) // shift right 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != -25 { // 200 >> 3 = -25
		t.Errorf("ISHR: expected a result of -25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("ISHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// ISTORE: Store integer from stack into local specified by following byte.
func TestNewIstore(t *testing.T) {
	f := newFrame(opcodes.ISTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22223))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[2] != int64(0x22223) {
		t.Errorf("ISTORE: Expecting 0x22223 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("ISTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ISTORE: Store byte value from stack into local specified by following byte.
func TestNewIstoreByte(t *testing.T) {
	f := newFrame(opcodes.ISTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, uint8(0x22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[2] != int64(0x22) {
		t.Errorf("ISTORE: Expecting int64 of 0x222 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("ISTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ISTORE_0: Store integer from stack into localVar[0]
func TestNewIstore0(t *testing.T) {
	f := newFrame(opcodes.ISTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(220))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[0] != int64(220) {
		t.Errorf("ISTORE_0: expected locals[0] to be 220, got: %d", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE_0: Store byte value from stack into localVar[0]
// Note: the logic for this bytecode is the same as ISTORE_1, ISTORE_2, ISTORE_3,
// so this test is not duplicated for those bytecodes
func TestNewIstore0Byte(t *testing.T) {
	f := newFrame(opcodes.ISTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, byte(220))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[0] != int64(220) {
		t.Errorf("ISTORE_0: expected locals[0] to be int64 of value 220, got value of: %d", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE_0: Store uint32 value from stack into localVar[0]
// Note: the logic for this bytecode is the same as ISTORE_1, ISTORE_2, ISTORE_3,
// so this test is not duplicated for those bytecodes
func TestNewIstore0Uint32(t *testing.T) {
	f := newFrame(opcodes.ISTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, uint32(220))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[0] != int64(220) {
		t.Errorf("ISTORE_0: expected locals[0] to be int64 of value 220, got value of: %d", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

func TestNewIstore0BooleanTrue(t *testing.T) {
	f := newFrame(opcodes.ISTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, true)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[0] != types.JavaBoolTrue {
		t.Errorf("ISTORE_0: expected locals[0] to be int64 of value 1, got value of: %d", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

func TestNewIstore0BooleanFalse(t *testing.T) {
	f := newFrame(opcodes.ISTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, false)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[0] != types.JavaBoolFalse {
		t.Errorf("ISTORE_0: expected locals[0] to be int64 of value 0, got value of: %d", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE_1
func TestNewIstore1(t *testing.T) {
	f := newFrame(opcodes.ISTORE_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(221))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[1] != int64(221) {
		t.Errorf("ISTORE_1: expected locals[1] to be 221, got: %d", f.Locals[1])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_1: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE_2
func TestNewIstore2(t *testing.T) {
	f := newFrame(opcodes.ISTORE_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(222))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[2] != int64(222) {
		t.Errorf("ISTORE_2: expected locals[2] to be 222, got: %d", f.Locals[2])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_2: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE_3
func TestNewIstore3(t *testing.T) {
	f := newFrame(opcodes.ISTORE_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(223))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.Locals[3] != int64(223) {
		t.Errorf("ISTORE_3: expected locals[3] to be 223, got: %d", f.Locals[3])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_3: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISUB: integer subtraction
func TestNewIsub(t *testing.T) {
	f := newFrame(opcodes.ISUB)
	push(&f, int64(10))
	push(&f, int64(7))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("ISUB, Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 3 {
		t.Errorf("ISUB: Expected popped value to be 3, got: %d", value)
	}
}

// IUSHR: unsigned right shift of int
func TestNewIushr(t *testing.T) {
	f := newFrame(opcodes.IUSHR)
	push(&f, int64(-200))
	push(&f, int64(3)) // shift right 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("IUSHR: expected a result of 25, but got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("IUSHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// IXOR: Logical XOR of two ints
func TestNewIxor(t *testing.T) {
	f := newFrame(opcodes.IXOR)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)
	if value != 3 { // 21 ^ 22 = 3
		t.Errorf("IXOR: expected a result of 3, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IXOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// JSR: compute the offset of a jump, push it, and do the jump
// In this test we use JSR to jump over the RETURN bytecode to a second JSR bytecode.
// That second bytecode, jumps backwards to the RETURN bytecode. In this way, we test
// forward and backward jumps. We then test that the jump offsets were appropriately
// placed on the op stack.
/*
func TestNewJSR(t *testing.T) {
	f := newFrame(opcodes.JSR) // jump to next JSR
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x04)
	f.Meth = append(f.Meth, opcodes.RETURN)
	f.Meth = append(f.Meth, opcodes.JSR) // jump backward to RETURN
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0xFF)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.OpStack[0] != int64(4) {
		t.Errorf("JSR: expected opstack[0] to be 4, got: %d", f.OpStack[0])
	}

	if f.OpStack[1] != int64(-1) {
		t.Errorf("JSR: expected opstack[1] to be -1, got: %d", f.OpStack[1])
	}
}
*/
// L2D: Convert long to double
func TestNewL2d(t *testing.T) {
	f := newFrame(opcodes.L2D)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	val := pop(&f).(float64)
	if val != 21.0 {
		t.Errorf("L2D: expected a result of 21.0, but got: %f", val)
	}
	if f.TOS != 0 {
		t.Errorf("L2D: Expected stack with 1 item, but got a TOS of: %d", f.TOS)
	}
}

// L2F: Convert long to float
func TestNewL2f(t *testing.T) {
	f := newFrame(opcodes.L2F)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	val := pop(&f).(float64)
	if val != 21.0 {
		t.Errorf("L2D: expected a result of 21.0, but got: %f", val)
	}

	if f.TOS != -1 {
		t.Errorf("L2D: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// L2I: Convert long to int
func TestNewL2i(t *testing.T) {
	f := newFrame(opcodes.L2I)
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	val := pop(&f).(int64)
	if val != 21 {
		t.Errorf("L2I: expected a result of 21, but got: %d", val)
	}

	if f.TOS != -1 {
		t.Errorf("L2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// L2I: Convert long to int (test with negative value)
func TestNewL2ineg(t *testing.T) {
	f := newFrame(opcodes.L2I)
	push(&f, int64(-21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	val := pop(&f).(int64)
	if val != -21 {
		t.Errorf("L2I: expected a result of -21, but got: %d", val)
	}

	if f.TOS != -1 {
		t.Errorf("L2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// LADD: Add two longs
func TestNewLadd(t *testing.T) {
	f := newFrame(opcodes.LADD)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 43 {
		t.Errorf("LADD: expected a result of 43, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LADD: Expected an empty stack, but got a TOS of: %d", f.TOS)
	}
}

// LAND: Logical and of two longs, push result
func TestNewLand(t *testing.T) {
	f := newFrame(opcodes.LAND)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)

	if value != 20 { // 21 & 22 = 20
		t.Errorf("LAND: expected a result of 20, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LAND: Expected an empty stack, but got a TOS of: %d", f.TOS)
	}
}

// LCMP: compare two longs (using two equal values)
func TestNewLcmpEQ(t *testing.T) {
	f := newFrame(opcodes.LCMP)
	push(&f, int64(21))
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("LCMP: Expected comparison to result in 0, got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LCMP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LCMP: compare two longs (with val1 > val2)
func TestNewLcmpGT(t *testing.T) {
	f := newFrame(opcodes.LCMP)
	push(&f, int64(22)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)
	if value != 1 {
		t.Errorf("LCMP: Expected comparison to result in 1, got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LCMP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LCMP: compare two longs (using val1 < val2)
func TestNewLcmpLT(t *testing.T) {
	f := newFrame(opcodes.LCMP)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)
	if value != -1 {
		t.Errorf("LCMP: Expected comparison to result in -1, got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LCMP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LCONST_0: push a long 0 onto opStack
func TestNewLconst0(t *testing.T) {
	f := newFrame(opcodes.LCONST_0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("LCONST_0: Expected popped value to be 0, got: %d", value)
	}
}

// LCONST_1: push a long 1 onto opStack
func TestNewLconst1(t *testing.T) {
	f := newFrame(opcodes.LCONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 1 {
		t.Errorf("LCONST_1: Expected popped value to be 1, got: %d", value)
	}
}

// LDC: get CP entry indexed by following byte
func TestNewLdc(t *testing.T) {
	f := newFrame(opcodes.LDC)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var ints = make([]int32, 1)
	CP.IntConsts = ints
	CP.IntConsts[0] = 25

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	doubleEntry := classloader.CpEntry{
		Type: classloader.IntConst, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 25 {
		t.Errorf("LDC_W: Expected popped value to be 25, got: %d", value)
	}
}

// LDC: get CP string entry indexed by following byte. Returns a string object
// whose value field contains an index into the string pool
func TestNewLdcTest2(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.LDC)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var utf8s = make([]string, 1)
	CP.Utf8Refs = utf8s
	CP.Utf8Refs[0] = "hello"

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	stringEntry := classloader.CpEntry{
		Type: classloader.UTF8, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, stringEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	strObj := pop(&f).(*object.Object)
	str := string(strObj.FieldTable["value"].Fvalue.([]byte))
	index := stringPool.GetStringIndex(&str)
	checkStrPtr := stringPool.GetStringPointer(index)
	if *checkStrPtr != "hello" {
		t.Errorf("LDC_W: Expected popped value to be 'hello', got %s", *checkStrPtr)
	}
}

// LDC cannot load a double. This tests that it generates the right error.
func TestNewLdcInvalidDouble(t *testing.T) {
	globals.InitGlobals("test")

	// hide the error message to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.LDC)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var doubles = make([]float64, 2)
	CP.Doubles = doubles
	CP.Doubles[0] = 1.234

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	stringEntry := classloader.CpEntry{
		Type: classloader.DoubleConst, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, stringEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		if !strings.Contains(errMsg, "LDC: Invalid type") {
			t.Errorf("Did not get expected error from LDC with double value, got: %s", errMsg)
		}
	} else {
		t.Errorf("Did not get expected error from LDC with double value")
	}
}

// Test LDC_W: get int64 CP entry indexed by two bytes
func TestNewLdcw(t *testing.T) {
	f := newFrame(opcodes.LDC_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var ints = make([]int32, 1)
	CP.IntConsts = ints
	CP.IntConsts[0] = 25

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	doubleEntry := classloader.CpEntry{
		Type: classloader.IntConst, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 25 {
		t.Errorf("LDC_W: Expected popped value to be 25, got: %d", value)
	}
}

// LDC_W: get float64 CP entry indexed by two bytes
func TestNewLdcwFloat(t *testing.T) {
	f := newFrame(opcodes.LDC_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var floats = make([]float32, 1)
	CP.Floats = floats
	CP.Floats[0] = 25.0

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	floatEntry := classloader.CpEntry{
		Type: classloader.FloatConst, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, floatEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 25.0 {
		t.Errorf("LDC_W: Expected popped value to be 25.0, got: %f", value)
	}
}

// LDC2_W: get CP entry for double indexed by following 2 bytes
func TestNewLdc2wForDouble(t *testing.T) {
	f := newFrame(opcodes.LDC2_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var doubles = make([]float64, 1)
	CP.Doubles = doubles
	CP.Doubles[0] = 25.0

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	doubleEntry := classloader.CpEntry{
		Type: classloader.DoubleConst, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 25.0 {
		t.Errorf("LDC2_W: Expected popped value to be 25.0, got: %f", value)
	}
}

// LDC2_W: get CP entry for long indexed by following 2 bytes
func TestNewLdc2wForLong(t *testing.T) {
	f := newFrame(opcodes.LDC2_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var longs = make([]int64, 1)
	CP.LongConsts = longs
	CP.LongConsts[0] = 25

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{} // CP 0
	doubleEntry := classloader.CpEntry{
		Type: classloader.LongConst, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != 25. {
		t.Errorf("LDC2_W: Expected popped value to be 25, got: %d", value)
	}
}

// LDC2_W can only be used for doubles and longs. Here we test its error repsonse when used on a string object.
func TestNewLdc2wInvalidForString(t *testing.T) {
	globals.InitGlobals("test")

	// hide the error message to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.LDC2_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	CP := f.CP.(*classloader.CPool)
	// now create a skeletal, two-entry CP
	var utf8s = make([]string, 1)
	CP.Utf8Refs = utf8s
	CP.Utf8Refs[0] = "hello"

	CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	stringEntry := classloader.CpEntry{
		Type: classloader.UTF8, Slot: 0,
	}
	CP.CpIndex = append(CP.CpIndex, dummyEntry)
	CP.CpIndex = append(CP.CpIndex, stringEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		if !strings.Contains(errMsg, "LDC2_W: Invalid type") {
			t.Errorf("Did not get expected error from LDC with double value, got: %s", errMsg)
		}
	} else {
		t.Errorf("Did not get expected error message in TestLdc2wInvalidForString()")
	}
}

// LDIV: (pop 2 longs, divide second term by top of stack, push result)
func TestNewLdiv(t *testing.T) {
	f := newFrame(opcodes.LDIV)
	push(&f, int64(70))
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("LDIV, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != 7 {
		t.Errorf("LDIV: Expected popped value to be 70, got: %d", value)
	}
}

// LDIV: with divide by zero error. This is handled in the wholeClassTests package

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Jacobin authors. All rights reserved.
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
	"os"
	"strings"
	"testing"
	"unsafe"
)

// Bytecodes tested in alphabetical order. Non-bytecode tests at ene of file.
// Note: array bytecodes are in array_test.go. All bytecodes from ACONST_NULL
// to IFNULL are in run_test.go. The remaining bytecodes are in this file.

// IINC: increment local variable
func TestIinc(t *testing.T) {
	f := newFrame(IINC)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(10)) // initialize local variable[1] to 10
	f.Meth = append(f.Meth, 1)             // increment local variable[1]
	f.Meth = append(f.Meth, 27)            // increment it by 27
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}
	value := f.Locals[1]
	if value != int64(37) {
		t.Errorf("IINC: Expected popped value to be 37, got: %d", value)
	}
}

// IINC: increment local variable by negative value
func TestIincNeg(t *testing.T) {
	f := newFrame(IINC)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(10)) // initialize local variable[1] to 10
	f.Meth = append(f.Meth, 1)             // increment local variable[1]
	val := -27
	f.Meth = append(f.Meth, byte(val)) // "increment" it by -27
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}
	value := f.Locals[1]
	if value != int64(-17) {
		t.Errorf("IINC: Expected popped value to be -17, got: %d", value)
	}
}

// ILOAD: test load of int in locals[index] on to stack
func TestIload(t *testing.T) {
	f := newFrame(ILOAD)
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
func TestIload0(t *testing.T) {
	f := newFrame(ILOAD_0)
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != int64(27) {
		t.Errorf("ILOAD_0: Expected popped value to be 27, got: %d", value)
	}
}

// ILOAD_1: load of int in locals[1] onto stack
func TestIload1(t *testing.T) {
	f := newFrame(ILOAD_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 27 {
		t.Errorf("ILOAD_1: Expected popped value to be 27, got: %d", value)
	}
}

// ILOAD_2: load of int in locals[2] onto stack
func TestIload2(t *testing.T) {
	f := newFrame(ILOAD_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(1))
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != int64(27) {
		t.Errorf("ILOAD_2: Expected popped value to be 27, got: %d", value)
	}
}

// ILOAD_3: load of int in locals[3] onto stack
func TestIload3(t *testing.T) {
	f := newFrame(ILOAD_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(1))
	f.Locals = append(f.Locals, int64(2))
	f.Locals = append(f.Locals, int64(27))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 27 {
		t.Errorf("ILOAD_3: Expected popped value to be 27, got: %d", value)
	}
}

// Test IMUL (pop 2 values, multiply them, push result)
func TestImul(t *testing.T) {
	f := newFrame(IMUL)
	push(&f, int64(10))
	push(&f, int64(7))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("IMUL, Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 70 {
		t.Errorf("IMUL: Expected popped value to be 70, got: %d", value)
	}
}

// INEG: negate an int
func TestIneg(t *testing.T) {
	f := newFrame(INEG)
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 {
		t.Errorf("INEG, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != -10 {
		t.Errorf("INEG: Expected popped value to be -10, got: %d", value)
	}
}

// INSTANCEOF: Is the TOS item an instance of a particular class?
func TestInstanceofNilAndNull(t *testing.T) {
	f := newFrame(INSTANCEOF)
	push(&f, nil)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("INSTANCEOF: Expected nil to return a 0, got %d", value)
	}

	f = newFrame(INSTANCEOF)
	push(&f, object.Null)

	fs = frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value = pop(&f).(int64)
	if value != 0 {
		t.Errorf("INSTANCEOF: Expected null to return a 0, got %d", value)
	}
}

// INSTANCEOF for a string
func TestInstanceofString(t *testing.T) {
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
	s := classloader.NewStringFromGoString("hello world")

	f := newFrame(INSTANCEOF)
	f.Meth = append(f.Meth, 0) // point to entry [2] in CP
	f.Meth = append(f.Meth, 2) // " "

	// now create the CP. First entry is perforce 0
	// [1] entry points to a UTF8 entry with the class name
	// [2] is a ClassRef that points to the UTF8 string in [1]
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 1}
	CP.ClassRefs = append(CP.ClassRefs, 0) // point to record 0 in Utf8Refs
	CP.Utf8Refs = append(CP.Utf8Refs, "java/lang/String")
	f.CP = &CP

	push(&f, s)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	if value != 1 { // a 1 = it's a match between class and object
		t.Errorf("INSTANCEOF: Expected string to return a 1, got %d", value)
	}
}

// IOR: Logical OR of two ints
func TestIor(t *testing.T) {
	f := newFrame(IOR)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)

	if value != 23 { // 21 | 22 = 23
		t.Errorf("IOR: expected a result of 23, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// IREM: int modulo
func TestIrem(t *testing.T) {
	f := newFrame(IREM)
	push(&f, int64(74))
	push(&f, int64(6))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 { // product is pushed twice b/c it's a long, which occupies 2 slots
		t.Errorf("IREM, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != 2 {
		t.Errorf("IREM: Expected result to be 2, got: %d", value)
	}
}

// IREM: int modulo -- divide by zero
func TestIremDivideByZero(t *testing.T) {
	f := newFrame(IREM)
	push(&f, int64(6))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)
	errMsg := err.Error()
	if !strings.Contains(errMsg, "divide by zero") {
		t.Errorf("IREM: Expected divide by zero error msg, got: %s", errMsg)
	}
}

// IRETURN: push an int on to the op stack of the calling method and exit the present method/frame
func TestIreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, int64(20))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f0)
	f1 := newFrame(IRETURN)
	push(&f1, int64(21))
	fs.PushFront(&f1)
	_ = runFrame(fs)
	_ = frames.PopFrame(fs)
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
func TestIshl(t *testing.T) {
	f := newFrame(ISHL)
	push(&f, int64(22)) // longs require two slots, so pushed twice
	push(&f, int64(3))  // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 176 { // 22 << 3 = 176
		t.Errorf("ISHL: expected a result of 176, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("ISHL: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// ISHR: Right shift of int
func TestIshr(t *testing.T) {
	f := newFrame(ISHR)
	push(&f, int64(200))
	push(&f, int64(3)) // shift right 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("ISHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("ISHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// ISHR: Right shift of negative int
func TestIshrNeg(t *testing.T) {
	f := newFrame(ISHR)
	push(&f, int64(-200))
	push(&f, int64(3)) // shift right 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != -25 { // 200 >> 3 = -25
		t.Errorf("ISHR: expected a result of -25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("ISHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
	/*
		// The following code runs correctly and prints -25 to the
		// console during test results.
		var printArray = make([]interface{}, 2)
		printArray[0] = 0
		printArray[1] = value
		classloader.PrintlnI(printArray)
	*/

}

// ISTORE: Store integer from stack into local specified by following byte.
func TestIstore(t *testing.T) {
	f := newFrame(ISTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22223))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2] != int64(0x22223) {
		t.Errorf("ISTORE: Expecting 0x22223 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("ISTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// ISTORE_0: Store integer from stack into localVar[0]
func TestIstore0(t *testing.T) {
	f := newFrame(ISTORE_0)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(220))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[0] != int64(220) {
		t.Errorf("ISTORE_0: expected lcoals[0] to be 220, got: %d", f.Locals[0])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE1
func TestIstore1(t *testing.T) {
	f := newFrame(ISTORE_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(221))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[1] != int64(221) {
		t.Errorf("ISTORE_1: expected locals[1] to be 221, got: %d", f.Locals[1])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_1: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISTORE2
func TestIstore2(t *testing.T) {
	f := newFrame(ISTORE_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(222))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[2] != int64(222) {
		t.Errorf("ISTORE_2: expected locals[2] to be 222, got: %d", f.Locals[2])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_2: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

func TestIstore3(t *testing.T) {
	f := newFrame(ISTORE_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(223))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.Locals[3] != int64(223) {
		t.Errorf("ISTORE_3: expected locals[3] to be 223, got: %d", f.Locals[3])
	}
	if f.TOS != -1 {
		t.Errorf("ISTORE_3: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// ISUB: integer subtraction
func TestIsub(t *testing.T) {
	f := newFrame(ISUB)
	push(&f, int64(10))
	push(&f, int64(7))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("ISUB, Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 3 {
		t.Errorf("ISUB: Expected popped value to be 3, got: %d", value)
	}
}

// IUSHR: unsigned right shift of int
func TestIushr(t *testing.T) {
	f := newFrame(IUSHR)
	push(&f, int64(-200))
	push(&f, int64(3)) // shift right 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("IUSHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IUSHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// IXOR: Logical XOR of two ints
func TestIxor(t *testing.T) {
	f := newFrame(IXOR)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	if value != 3 { // 21 ^ 22 = 3
		t.Errorf("IXOR: expected a result of 3, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("IXOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// L2D: Convert long to double
func TestL2d(t *testing.T) {
	f := newFrame(L2D)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(float64)
	if val != 21.0 {
		t.Errorf("L2D: expected a result of 21.0, but got: %f", val)
	}
	if f.TOS != 0 {
		t.Errorf("L2D: Expected stack with 1 item, but got a TOS of: %d", f.TOS)
	}
}

// L2F: Convert long to float
func TestL2f(t *testing.T) {
	f := newFrame(L2F)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(float64)
	if val != 21.0 {
		t.Errorf("L2D: expected a result of 21.0, but got: %f", val)
	}
	if f.TOS != -1 {
		t.Errorf("L2D: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// L2I: Convert long to int
func TestL2i(t *testing.T) {
	f := newFrame(L2I)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != 21 {
		t.Errorf("L2I: expected a result of 21, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("L2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// L2I: Convert long to int (test with negative value)
func TestL2ineg(t *testing.T) {
	f := newFrame(L2I)
	push(&f, int64(-21)) // longs require two slots, so pushed twice
	push(&f, int64(-21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	val := pop(&f).(int64)
	if val != -21 {
		t.Errorf("L2I: expected a result of -21, but got: %d", val)
	}
	if f.TOS != -1 {
		t.Errorf("L2I: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
	}
}

// LADD: Add two longs
func TestLadd(t *testing.T) {
	f := newFrame(LADD)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	push(&f, int64(22))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 43 {
		t.Errorf("LADD: expected a result of 43, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LADD: Expected an empty stack, but got a TOS of: %d", f.TOS)
	}
}

// LAND: Logical and of two longs, push result
func TestLand(t *testing.T) {
	f := newFrame(LAND)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	push(&f, int64(22))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 20 { // 21 & 22 = 20
		t.Errorf("LAND: expected a result of 20, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LAND: Expected an empty stack, but got a TOS of: %d", f.TOS)
	}
}

// LCMP: compare two longs (using two equal values)
func TestLcmpEQ(t *testing.T) {
	f := newFrame(LCMP)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	push(&f, int64(21))
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("LCMP: Expected comparison to result in 0, got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LCMP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LCMP: compare two longs (with val1 > val2)
func TestLcmpGT(t *testing.T) {
	f := newFrame(LCMP)
	push(&f, int64(22)) // longs require two slots, so pushed twice
	push(&f, int64(22))

	push(&f, int64(21))
	push(&f, int64(21))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	if value != 1 {
		t.Errorf("LCMP: Expected comparison to result in 1, got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LCMP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LCMP: compare two longs (using val1 < val2)
func TestLcmpLT(t *testing.T) {
	f := newFrame(LCMP)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	push(&f, int64(22))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	if value != -1 {
		t.Errorf("LCMP: Expected comparison to result in -1, got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LCMP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LCONST_0: push a long 0 onto opStack
func TestLconst0(t *testing.T) {
	f := newFrame(LCONST_0)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 1 {
		t.Errorf("Top of stack, expected 1, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 0 {
		t.Errorf("LCONST_0: Expected popped value to be 0, got: %d", value)
	}
}

// LCONST_1: push a long 1 onto opStack
func TestLconst1(t *testing.T) {
	f := newFrame(LCONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 1 {
		t.Errorf("Top of stack, expected 1, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 1 {
		t.Errorf("LCONST_1: Expected popped value to be 1, got: %d", value)
	}
}

// LDC_W: get CP entry indexed by following byte
func TestLdc(t *testing.T) {
	f := newFrame(LDC)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	// now create a skeletal, two-entry CP
	var ints = make([]int32, 1)
	f.CP.IntConsts = ints
	f.CP.IntConsts[0] = 25

	f.CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	doubleEntry := classloader.CpEntry{
		Type: classloader.IntConst, Slot: 0,
	}
	f.CP.CpIndex = append(f.CP.CpIndex, dummyEntry)
	f.CP.CpIndex = append(f.CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 25 {
		t.Errorf("LDC_W: Expected popped value to be 25, got: %d", value)
	}
}

// Test LDC_W: get int64 CP entry indexed by two bytes
func TestLdcw(t *testing.T) {
	f := newFrame(LDC_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	// now create a skeletal, two-entry CP
	var ints = make([]int32, 1)
	f.CP.IntConsts = ints
	f.CP.IntConsts[0] = 25

	f.CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	doubleEntry := classloader.CpEntry{
		Type: classloader.IntConst, Slot: 0,
	}
	f.CP.CpIndex = append(f.CP.CpIndex, dummyEntry)
	f.CP.CpIndex = append(f.CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 25 {
		t.Errorf("LDC_W: Expected popped value to be 25, got: %d", value)
	}
}

// LDC_W: get float64 CP entry indexed by two bytes
func TestLdcwFloat(t *testing.T) {
	f := newFrame(LDC_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	// now create a skeletal, two-entry CP
	var floats = make([]float32, 1)
	f.CP.Floats = floats
	f.CP.Floats[0] = 25.0

	f.CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	floatEntry := classloader.CpEntry{
		Type: classloader.FloatConst, Slot: 0,
	}
	f.CP.CpIndex = append(f.CP.CpIndex, dummyEntry)
	f.CP.CpIndex = append(f.CP.CpIndex, floatEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 25.0 {
		t.Errorf("LDC_W: Expected popped value to be 25.0, got: %f", value)
	}
}

// LDC2_W: get CP entry for long or double indexed by following 2 bytes
func TestLdc2w(t *testing.T) {
	f := newFrame(LDC2_W)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01)

	cp := classloader.CPool{}
	f.CP = &cp
	// now create a skeletal, two-entry CP
	var doubles = make([]float64, 1)
	f.CP.Doubles = doubles
	f.CP.Doubles[0] = 25.0

	f.CP.CpIndex = []classloader.CpEntry{}
	dummyEntry := classloader.CpEntry{}
	doubleEntry := classloader.CpEntry{
		Type: classloader.DoubleConst, Slot: 0,
	}
	f.CP.CpIndex = append(f.CP.CpIndex, dummyEntry)
	f.CP.CpIndex = append(f.CP.CpIndex, doubleEntry)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 1 {
		t.Errorf("Top of stack, expected 1, got: %d", f.TOS)
	}
	value := pop(&f).(float64)
	if value != 25.0 {
		t.Errorf("LDC2_W: Expected popped value to be 25.0, got: %f", value)
	}
}

// LDIV: (pop 2 longs, divide second term by top of stack, push result)
func TestLdiv(t *testing.T) {
	f := newFrame(LDIV)
	push(&f, int64(70))
	push(&f, int64(70))

	push(&f, int64(10))
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 1 { // product is pushed twice b/c it's a long, which occupies 2 slots
		t.Errorf("LDIV, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	pop(&f)
	if value != 7 {
		t.Errorf("LDIV: Expected popped value to be 70, got: %d", value)
	}
}

// LDIV: with divide by zero error
func TestLdivDivideByZero(t *testing.T) {
	f := newFrame(LDIV)
	push(&f, int64(10))
	push(&f, int64(10))

	push(&f, int64(0))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	res := runFrame(fs)

	if !strings.Contains(res.Error(), "Divide by zero") {
		t.Errorf("LDIV: Expected err msg re divide by zero, got %s", res.Error())
	}
}

// LLOAD: test load of long in locals[index] on to stack
func TestLload(t *testing.T) {
	f := newFrame(LLOAD)
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
	pop(&f) // pop twice due to two entries on op stack due to 64-bit width of data type
	if x != 0x1234562 {
		t.Errorf("LLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("LLOAD: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
	if f.PC != 2 {
		t.Errorf("LLOAD: Expected pc to be pointing at byte 2, got: %d", f.PC)
	}
}

// LLOAD_0: Load long from locals[0]
func TestLload0(t *testing.T) {
	f := newFrame(LLOAD_0)

	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[0]
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[1] // lload uses two local consecutive

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	pop(&f) // due to longs taking 2 slots
	if x != 0x12345678 {
		t.Errorf("LLOAD_0: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.Locals[1] != x {
		t.Errorf("LLOAD_0: Local variable[1] holds invalid value: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_0: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LLOAD_1: Load long from locals[1]
func TestLload1(t *testing.T) {
	f := newFrame(LLOAD_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[1]
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[2] // lload uses two local consecutive

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	pop(&f) // due to longs taking two slots
	if x != 0x12345678 {
		t.Errorf("LLOAD_1: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.Locals[2] != x {
		t.Errorf("LLOAD_1: Local variable[2] holds invalid value: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_1: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LLOAD_2: Load long from locals[2]
func TestLload2(t *testing.T) {
	f := newFrame(LLOAD_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[2]
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[3] // lload uses two local consecutive

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	pop(&f) // due to longs taking two slots
	if x != 0x12345678 {
		t.Errorf("LLOAD_12: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.Locals[3] != x {
		t.Errorf("LLOAD_2: Local variable[3] holds invalid value: 0x%x", f.Locals[3])
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_1: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LLOAD_3: Load long from locals[3]
func TestLload3(t *testing.T) {
	f := newFrame(LLOAD_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[3]
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[4] // lload uses two local consecutive

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f).(int64)
	pop(&f) // due to longs taking two slots
	if x != 0x12345678 {
		t.Errorf("LLOAD_3: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.Locals[4] != x {
		t.Errorf("LLOAD_3: Local variable[4] holds invalid value: 0x%x", f.Locals[4])
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_3: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LMUL: pop 2 longs, multiply them, push result
func TestLmul(t *testing.T) {
	f := newFrame(LMUL)
	push(&f, int64(10))
	push(&f, int64(10))

	push(&f, int64(7))
	push(&f, int64(7))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 1 { // product is pushed twice b/c it's a long, which occupies 2 slots
		t.Errorf("LMUL, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	pop(&f)
	if value != 70 {
		t.Errorf("LMUL: Expected popped value to be 70, got: %d", value)
	}
}

// LNEG: negate a long
func TestLneg(t *testing.T) {
	f := newFrame(LNEG)
	push(&f, int64(10))
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 1 { // product is pushed twice b/c it's a long, which occupies 2 slots
		t.Errorf("LNEG, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	pop(&f)
	if value != -10 {
		t.Errorf("LNEG: Expected popped value to be -10, got: %d", value)
	}
}

// LOR: Logical OR of two longs
func TestLor(t *testing.T) {
	f := newFrame(LOR)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	push(&f, int64(22))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 23 { // 21 | 22 = 23
		t.Errorf("LOR: expected a result of 23, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LREM: remainder of long division (the % operator)
func TestLrem(t *testing.T) {
	f := newFrame(LREM)
	push(&f, int64(74))
	push(&f, int64(74))

	push(&f, int64(6))
	push(&f, int64(6))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 1 { // product is pushed twice b/c it's a long, which occupies 2 slots
		t.Errorf("LREM, Top of stack, expected 1, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	pop(&f)
	if value != 2 {
		t.Errorf("LREM: Expected popped value to be 2, got: %d", value)
	}
}

// LREM: long modulo -- divide by zero
func TestLremDivideByZero(t *testing.T) {
	f := newFrame(LREM)
	push(&f, int64(6))
	push(&f, int64(6))
	push(&f, int64(0))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)
	errMsg := err.Error()
	if !strings.Contains(errMsg, "divide by zero") {
		t.Errorf("LREM: Expected divide by zero error msg, got: %s", errMsg)
	}
}

// LRETURN: Return a long from a function
func TestLreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, int64(20))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f0)
	f1 := newFrame(LRETURN)
	push(&f1, int64(21))
	push(&f1, int64(21))
	fs.PushFront(&f1)
	_ = runFrame(fs)
	_ = frames.PopFrame(fs)
	f3 := fs.Front().Value.(*frames.Frame)
	newVal := pop(f3).(int64)
	if newVal != 21 {
		t.Errorf("After LRETURN, expected a value of 21 in previous frame, got: %d", newVal)
	}
	pop(f3) // popped a second time due to longs taking two slots

	prevVal := pop(f3).(int64)
	if prevVal != 20 {
		t.Errorf("After LRETURN, expected a value of 20 in 2nd place of previous frame, got: %d", prevVal)
	}
}

// LSHL: Left shift of long
func TestLshl(t *testing.T) {
	f := newFrame(LSHL)
	push(&f, int64(22)) // longs require two slots, so pushed twice
	push(&f, int64(22))

	push(&f, int64(3)) // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 176 { // 22 << 3 = 176
		t.Errorf("LSHL: expected a result of 176, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LSHL: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LSHR: Right shift of long
func TestLshr(t *testing.T) {
	f := newFrame(LSHR)
	push(&f, int64(200)) // longs require two slots, so pushed twice
	push(&f, int64(200))

	push(&f, int64(3)) // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("LSHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LSHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LSTORE: Store long from stack into local specified by following byte, and the local var after it.
func TestLstore(t *testing.T) {
	f := newFrame(LSTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22223))
	push(&f, int64(0x22223)) // push twice due to longs using two slots

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2] != int64(0x22223) {
		t.Errorf("LSTORE: Expecting 0x22223 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.Locals[3] != int64(0x22223) {
		t.Errorf("LSTORE: Expecting 0x22223 in locals[3], got: 0x%x", f.Locals[3])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LSTORE_0: Store long from stack in localVar[0] and again in localVar[1]
func TestLstore0(t *testing.T) {
	f := newFrame(LSTORE_0)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[0] != int64(0x12345678) {
		t.Errorf("LSTORE_0: expected locals[0] to be 0x12345678, got: %d", f.Locals[0])
	}

	if f.Locals[1] != int64(0x12345678) {
		t.Errorf("LSTORE_0: expected locals[1] to be 0x12345678, got: %d", f.Locals[1])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSTORE_1: Store long from stack in localVar[1] and again in localVar[2]
func TestLstore1(t *testing.T) {
	f := newFrame(LSTORE_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[1] != int64(0x12345678) {
		t.Errorf("LSTORE_1: expected locals[1] to be 0x12345678, got: %d", f.Locals[1])
	}

	if f.Locals[2] != int64(0x12345678) {
		t.Errorf("LSTORE_1: expected locals[2] to be 0x12345678, got: %d", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_1: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSTORE_2: Store long from stack in localVar[2] and again in localVar[3]
func TestLstore2(t *testing.T) {
	f := newFrame(LSTORE_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[2] != int64(0x12345678) {
		t.Errorf("LSTORE_2: expected locals[2] to be 0x12345678, got: %d", f.Locals[2])
	}

	if f.Locals[3] != int64(0x12345678) {
		t.Errorf("LSTORE_2: expected locals[3] to be 0x12345678, got: %d", f.Locals[3])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_2: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSTORE_3: Store long from stack in localVar[3] and again in localVar[]
func TestLstore3(t *testing.T) {
	f := newFrame(LSTORE_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.Locals[3] != int64(0x12345678) {
		t.Errorf("LSTORE_3: expected locals[3] to be 0x12345678, got: %d", f.Locals[3])
	}

	if f.Locals[4] != int64(0x12345678) {
		t.Errorf("LSTORE_3: expected locals[4] to be 0x12345678, got: %d", f.Locals[4])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_3: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSUB: Subtract two longs
func TestLsub(t *testing.T) {
	f := newFrame(LSUB)
	push(&f, int64(10)) // longs occupy two slots, hence the double pops and pushes
	push(&f, int64(10))

	push(&f, int64(7))
	push(&f, int64(7))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64)
	pop(&f)

	if f.TOS != -1 {
		t.Errorf("LSUB, Top of stack, expected -1, got: %d", f.TOS)
	}

	if value != 3 {
		t.Errorf("LSUB: Expected popped value to be 3, got: %d", value)
	}
}

// LUSHR: Right unsigned shift of long
func TestLushr(t *testing.T) {
	f := newFrame(LUSHR)
	push(&f, int64(200)) // longs require two slots, so pushed twice
	push(&f, int64(200))

	push(&f, int64(3)) // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("LUSHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LUSHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LXOR: Logical XOR of two longs
func TestLxor(t *testing.T) {
	f := newFrame(LXOR)
	push(&f, int64(21)) // longs require two slots, so pushed twice
	push(&f, int64(21))

	push(&f, int64(22))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice
	pop(&f)

	if value != 3 { // 21 ^ 22 = 3
		t.Errorf("LXOR: expected a result of 3, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LXOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// MONITORENTER: The JDK JVM does not implement this, nor do we. So just pop the ref off stack
func TestMonitorEnter(t *testing.T) {
	f := newFrame(MONITORENTER)
	push(&f, &f) // push any value and make sure it gets popped off

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != -1 {
		t.Errorf("MONITORENTER: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// MONITOREXIT: The JDK JVM does not implement this, nor do we. So just pop the ref off stack
func TestMonitorExit(t *testing.T) {
	f := newFrame(MONITOREXIT)
	push(&f, &f) // push any value and make sure it gets popped off

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != -1 {
		t.Errorf("MONITOREXIT: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// POP: pop item off stack and discard it
func TestPop(t *testing.T) {
	f := newFrame(POP)
	push(&f, int64(34)) // push three different values
	push(&f, int64(21))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 1 {
		t.Errorf("POP: Expected stack with 2 items, but got a tos of: %d", f.TOS)
	}

	top := pop(&f).(int64)

	if top != 21 {
		t.Errorf("POP: expected top's value to be 21, but got: %d", top)
	}
}

// POP2: pop two items off stack and discard them
func TestPop2(t *testing.T) {
	f := newFrame(POP2)
	push(&f, int64(34)) // push three different values; 34 at bottom
	push(&f, int64(21))
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.TOS != 0 {
		t.Errorf("POP2: Expected stack with 1 item, but got a tos of: %d", f.TOS)
	}

	top := pop(&f).(int64)

	if top != 34 {
		t.Errorf("POP2: expected top's value to be 34, but got: %d", top)
	}
}

// PUTFIELD: Update a field in an object -- error doesn't point to a field
func TestPutFieldNonFieldCPentry(t *testing.T) {
	f := newFrame(PUTFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: 8, Slot: 0} // point to non-fieldRef
	CP.FieldRefs = make([]classloader.FieldRefEntry, 1, 1)
	CP.FieldRefs[0] = classloader.FieldRefEntry{ClassIndex: 0, NameAndType: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	err := runFrame(fs)

	msg := err.Error()
	if !strings.Contains(msg, "PUTFIELD: Expected a field ref, but") {
		t.Errorf("PUTFIELD: Did not get expected error msg: %s", msg)
	}
}

// RETURN: Does a function return correctly?
func TestReturn(t *testing.T) {
	f := newFrame(RETURN)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	ret := runFrame(fs)
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}

	if ret != nil {
		t.Error("RETURN: Expected popped value to be nil, got: " + ret.Error())
	}
}

// SIPUSH: create int from next two bytes and push the int
func TestSipush(t *testing.T) {
	f := newFrame(SIPUSH)
	f.Meth = append(f.Meth, 0x01)
	f.Meth = append(f.Meth, 0x02)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 258 {
		t.Errorf("SIPUSH: Expected popped value to be 258, got: %d", value)
	}

	if f.PC != 3 {
		t.Errorf("SIPUSH: Expected PC to be 3, got: %d", f.PC)
	}
}

// SIPUSH: create a negative int from next two bytes and push the int
func TestSipushNegative(t *testing.T) {
	f := newFrame(SIPUSH)
	val := -1
	f.Meth = append(f.Meth, byte(val))
	f.Meth = append(f.Meth, 0x02)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value >= 0 {
		t.Errorf("SIPUSH: Expected popped value to be negative, got: %d", value)
	}

	if f.PC != 3 {
		t.Errorf("SIPUSH: Expected PC to be 3, got: %d", f.PC)
	}
}

// SWAP: Swap top two items on stack
func TestSwap(t *testing.T) {
	f := newFrame(SWAP)
	push(&f, int64(34)) // push two different values
	push(&f, int64(21)) // TOS now = 21

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	top := pop(&f).(int64)
	next := pop(&f).(int64)

	if top != 34 {
		t.Errorf("SWAP: expected top's value to be 34, but got: %d", top)
	}

	if next != 21 {
		t.Errorf("SWAP: expected next's value to be 21, but got: %d", next)
	}

	if f.TOS != -1 {
		t.Errorf("SWAP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

func TestInvalidInstruction(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global := globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)
	LoadOptionsTable(Global)

	// to avoid cluttering the test results, redirect stdout
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(252)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	ret := runFrame(fs)
	if ret == nil {
		t.Errorf("Invalid instruction: Expected an error returned, but got nil.")
	}

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Invalid bytecode") {
		t.Errorf("Error message for invalid bytecode not as expected, got: %s", msg)
	}
}

func TestConvertInterfaceToUint64(t *testing.T) {
	var i64 int64 = 200
	var f64 float64 = 345.0
	var ptr = unsafe.Pointer(&f64)

	ret := convertInterfaceToUint64(i64)
	if ret != 200 {
		t.Errorf("Expected TestConvertInterfaceToUint64() to retun 200, got %d\n",
			ret)
	}

	ret = convertInterfaceToUint64(f64)
	if ret != 345 {
		t.Errorf("Expected TestConvertInterfaceToUint64() to retun 345, got %d\n",
			ret)
	}

	ret = convertInterfaceToUint64(ptr)
	if ret == 0 { // a minimal test
		t.Error("Expected TestConvertInterfaceToUint64() to !=0, got 0\n")
	}
}

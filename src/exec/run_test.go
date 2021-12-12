/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exec

import (
	"testing"
)

// These tests test the individual bytecode instructions. They are presented here in
// alphabetical order of the instruction name.

// set up function to create a frame with a method with the single instruction
// that's being tested
func newFrame(code byte) frame {
	f := createFrame(6)
	f.ftype = 'J'
	f.meth = append(f.meth, code)
	return *f
}

// ---- tests ----

func TestBipush(t *testing.T) {
	f := newFrame(BIPUSH)
	f.meth = append(f.meth, 0x05)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 5 {
		t.Errorf("BIPUSH: Expected popped value to be 5, got: %d", value)
	}
}

func TestIadd(t *testing.T) {
	f := newFrame(IADD)
	push(&f, 21)
	push(&f, 22)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f)
	if value != 43 {
		t.Errorf("IADD: expected a result of 43, but got: %d", value)
	}
	if f.tos != -1 {
		t.Errorf("IADD: Expected an empty stack, but got a tos of: %d", f.tos)
	}
}

func TestIconstN1(t *testing.T) {
	f := newFrame(ICONST_N1)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != -1 {
		t.Errorf("ICONST_N1: Expected popped value to be -1, got: %d", value)
	}
}

func TestIconst0(t *testing.T) {
	f := newFrame(ICONST_0)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 0 {
		t.Errorf("ICONST_0: Expected popped value to be 0, got: %d", value)
	}
}

func TestIconst1(t *testing.T) {
	f := newFrame(ICONST_1)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 1 {
		t.Errorf("ICONST_1: Expected popped value to be 1, got: %d", value)
	}
}

func TestIconst2(t *testing.T) {
	f := newFrame(ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 2 {
		t.Errorf("ICONST_2: Expected popped value to be 2, got: %d", value)
	}
}

func TestIconst3(t *testing.T) {
	f := newFrame(ICONST_3)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 3 {
		t.Errorf("ICONST_3: Expected popped value to be 3, got: %d", value)
	}
}

func TestIconst4(t *testing.T) {
	f := newFrame(ICONST_4)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 4 {
		t.Errorf("ICONST_4: Expected popped value to be 4, got: %d", value)
	}
}

func TestIconst5(t *testing.T) {
	f := newFrame(ICONST_5)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 5 {
		t.Errorf("ICONST_5: Expected popped value to be 5, got: %d", value)
	}
}

// ICMPGE: if integer compare val 1 >= val 2. Here test for = (next test for >)
func TestIfIcmpge1(t *testing.T) {
	f := newFrame(IF_ICMPGE)
	push(&f, 9)
	push(&f, 9)
	// note that the byte passed in newframe() is at f.meth[0]
	f.meth = append(f.meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.meth = append(f.meth, 4)
	f.meth = append(f.meth, ICONST_1)
	f.meth = append(f.meth, ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPGE: expecting a jump to ICONST_2 instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

// ICMPGE: if integer compare val 1 >= val 2. Here test for > (previous test for =)
func TestIfIcmpge21(t *testing.T) {
	f := newFrame(IF_ICMPGE)
	push(&f, 9)
	push(&f, 8)
	// note that the byte passed in newframe() is at f.meth[0]
	f.meth = append(f.meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.meth = append(f.meth, 4)
	f.meth = append(f.meth, ICONST_1)
	f.meth = append(f.meth, ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPGE: expecting a jump to ICONST_2 instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

// ICMPGE: if integer compare val 1 >= val 2 //test when condition fails
func TestIfIcmgetFail(t *testing.T) {
	f := newFrame(IF_ICMPGE)
	push(&f, 8)
	push(&f, 9)
	// note that the byte passed in newframe() is at f.meth[0]
	f.meth = append(f.meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.meth = append(f.meth, 4)
	f.meth = append(f.meth, RETURN) // the failed test should drop to this
	f.meth = append(f.meth, ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc] != RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("ICMPGE: expecting fall-through to RETURN instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

// ICMPLT: if integer compare val 1 < val 2
func TestIfIcmplt(t *testing.T) {
	f := newFrame(IF_ICMPLT)
	push(&f, 8)
	push(&f, 9)
	// note that the byte passed in newframe() is at f.meth[0]
	f.meth = append(f.meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.meth = append(f.meth, 4)
	f.meth = append(f.meth, ICONST_1)
	f.meth = append(f.meth, ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
		t.Errorf("ICMPLT: expecting a jump to ICONST_2 instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

// ICMPLT: if integer compare val 1 < val 2 //test when condition fails
func TestIfIcmpltFail(t *testing.T) {
	f := newFrame(IF_ICMPLT)
	push(&f, 9)
	push(&f, 9)
	// note that the byte passed in newframe() is at f.meth[0]
	f.meth = append(f.meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.meth = append(f.meth, 4)
	f.meth = append(f.meth, RETURN) // the failed test should drop to this
	f.meth = append(f.meth, ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc] != RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("ICMPLT: expecting fall-through to RETURN instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

func TestIinc(t *testing.T) {
	f := newFrame(IINC)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 10) // initialize local variable[1] to 10
	f.meth = append(f.meth, 1)      // increment local variable[1]
	f.meth = append(f.meth, 27)     // increment it by 27
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.tos)
	}
	value := f.locals[1]
	if value != 37 {
		t.Errorf("IINC: Expected popped value to be 37, got: %d", value)
	}
}

func TestIload0(t *testing.T) {
	f := newFrame(ILOAD_0)
	f.locals = append(f.locals, 27)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 27 {
		t.Errorf("ILOAD_0: Expected popped value to be 27, got: %d", value)
	}
}

func TestIload1(t *testing.T) {
	f := newFrame(ILOAD_1)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 27)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 27 {
		t.Errorf("ILOAD_1: Expected popped value to be 27, got: %d", value)
	}
}

func TestIload2(t *testing.T) {
	f := newFrame(ILOAD_2)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 1)
	f.locals = append(f.locals, 27)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 27 {
		t.Errorf("ILOAD_2: Expected popped value to be 27, got: %d", value)
	}
}

func TestIload3(t *testing.T) {
	f := newFrame(ILOAD_3)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 1)
	f.locals = append(f.locals, 2)
	f.locals = append(f.locals, 27)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 27 {
		t.Errorf("ILOAD_3: Expected popped value to be 27, got: %d", value)
	}
}

// Test IMUL (pop 2 values, multiply them, push result)
func TestImul(t *testing.T) {
	f := newFrame(IMUL)
	push(&f, 10)
	push(&f, 7)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("IMUL, Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 70 {
		t.Errorf("IMUL: Expected popped value to be 70, got: %d", value)
	}
}

func TestIstore0(t *testing.T) {
	f := newFrame(ISTORE_0)
	f.locals = append(f.locals, 0)
	push(&f, 220)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.locals[0] != 220 {
		t.Errorf("After ISTORE_0, expected lcoals[2] to be 220, got: %d", f.locals[0])
	}
	if f.tos != -1 {
		t.Errorf("ISTORE_0: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestIstore1(t *testing.T) {
	f := newFrame(ISTORE_1)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	push(&f, 221)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.locals[1] != 221 {
		t.Errorf("After ISTORE_1, expected lcoals[1] to be 221, got: %d", f.locals[1])
	}
	if f.tos != -1 {
		t.Errorf("ISTORE_1: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestIstore2(t *testing.T) {
	f := newFrame(ISTORE_2)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	push(&f, 222)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.locals[2] != 222 {
		t.Errorf("After ISTORE_2, expected lcoals[2] to be 222, got: %d", f.locals[2])
	}
	if f.tos != -1 {
		t.Errorf("ISTORE_2: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestIstore3(t *testing.T) {
	f := newFrame(ISTORE_3)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	push(&f, 223)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.locals[3] != 223 {
		t.Errorf("After ISTORE_3, expected lcoals[0] to be 223, got: %d", f.locals[3])
	}
	if f.tos != -1 {
		t.Errorf("ISTORE_3: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestIsub(t *testing.T) {
	f := newFrame(ISUB)
	push(&f, 10)
	push(&f, 7)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("ISUB, Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 3 {
		t.Errorf("ISUB: Expected popped value to be 3, got: %d", value)
	}
}

func TestLdc(t *testing.T) {
	f := newFrame(LDC)
	f.meth = append(f.meth, 0x05)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 5 {
		t.Errorf("LDC: Expected popped value to be 5, got: %d", value)
	}
}

func TestReturn(t *testing.T) {
	f := newFrame(RETURN)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	ret := runFrame(fs)
	if f.tos != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.tos)
	}

	if ret != nil {
		t.Error("RETURN: Expected popped value to be 2, got: " + ret.Error())
	}
}

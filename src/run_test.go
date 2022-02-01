/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"io/ioutil"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
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

// test load of reference in locals[index] on to stack
func TestAload(t *testing.T) {
	f := newFrame(ALOAD)
	f.meth = append(f.meth, 0x04) // use local var #4
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234562) // put value in locals[4]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234562 {
		t.Errorf("ALOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("ALOAD: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
	if f.pc != 2 {
		t.Errorf("ALOAD: Expected pc to be pointing at byte 2, got: %d", f.pc)
	}
}

// test load of reference in locals[0] on to stack
func TestAload0(t *testing.T) {
	f := newFrame(ALOAD_0)
	f.locals = append(f.locals, 0x1234560) // put value in locals[0]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234560 {
		t.Errorf("ALOAD_0: Expecting 0x1234560 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("ALOAD_0: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test load of reference in locals[1] on to stack
func TestAload1(t *testing.T) {
	f := newFrame(ALOAD_1)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234561) // put value in locals[1]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234561 {
		t.Errorf("ALOAD_1: Expecting 0x1234561 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("ALOAD_1: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test load of reference in locals[2] on to stack
func TestAload2(t *testing.T) {
	f := newFrame(ALOAD_2)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234562) // put value in locals[2]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234562 {
		t.Errorf("ALOAD_2: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("ALOAD_2: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test load of reference in locals[3] on to stack
func TestAload3(t *testing.T) {
	f := newFrame(ALOAD_3)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234563) // put value in locals[3]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234563 {
		t.Errorf("ALOAD_3: Expecting 0x1234563 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("ALOAD_3: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test store of reference from stack into locals[0]
func TestAstore0(t *testing.T) {
	f := newFrame(ASTORE_0)
	f.locals = append(f.locals, 0)
	push(&f, 0x22220)

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[0] != 0x22220 {
		t.Errorf("ASTORE_0: Expecting 0x22220 on stack, got: 0x%x", f.locals[0])
	}
	if f.tos != -1 {
		t.Errorf("ASTORE_0: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test store of reference from stack into locals[1]
func TestAstore1(t *testing.T) {
	f := newFrame(ASTORE_1)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	push(&f, 0x22221)

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[1] != 0x22221 {
		t.Errorf("ASTORE_1: Expecting 0x22221 on stack, got: 0x%x", f.locals[0])
	}
	if f.tos != -1 {
		t.Errorf("ASTORE_1: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test store of reference from stack into locals[2]
func TestAstore2(t *testing.T) {
	f := newFrame(ASTORE_2)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	push(&f, 0x22222)

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[2] != 0x22222 {
		t.Errorf("ASTORE_2: Expecting 0x22222 on stack, got: 0x%x", f.locals[0])
	}
	if f.tos != -1 {
		t.Errorf("ASTORE_2: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// test store of reference from stack into locals[3]
func TestAstore3(t *testing.T) {
	f := newFrame(ASTORE_3)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	push(&f, 0x22223)

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[3] != 0x22223 {
		t.Errorf("ASTORE_3: Expecting 0x22223 on stack, got: 0x%x", f.locals[0])
	}
	if f.tos != -1 {
		t.Errorf("ASTORE_3: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

func TestBipush(t *testing.T) {
	f := newFrame(BIPUSH)
	f.meth = append(f.meth, 0x05)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 5 {
		t.Errorf("BIPUSH: Expected popped value to be 5, got: %d", value)
	}
}

// DLOAD: test load of double in locals[index] on to stack
func TestDload(t *testing.T) {
	f := newFrame(DLOAD)
	f.meth = append(f.meth, 0x04) // use local var #4
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234562) // put value in locals[4]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234562 {
		t.Errorf("DLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("DLOAD: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
	if f.pc != 2 {
		t.Errorf("DLOAD: Expected pc to be pointing at byte 2, got: %d", f.pc)
	}
}

// FLOAD: test load of float in locals[index] on to stack
func TestFload(t *testing.T) {
	f := newFrame(FLOAD)
	f.meth = append(f.meth, 0x04) // use local var #4
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234562) // put value in locals[4]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234562 {
		t.Errorf("FLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("FLOAD: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
	if f.pc != 2 {
		t.Errorf("FLOAD: Expected pc to be pointing at byte 2, got: %d", f.pc)
	}
}

// test of GOTO instruction -- in forward direction (to a later bytecode)
func TestGotoForward(t *testing.T) {
	f := newFrame(GOTO)
	f.meth = append(f.meth, 0x00)
	f.meth = append(f.meth, 0x03)
	f.meth = append(f.meth, RETURN)
	f.meth = append(f.meth, NOP)
	f.meth = append(f.meth, NOP)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc] != RETURN {
		t.Errorf("GOTO forward: Expected pc to point to RETURN, but instead it points to : %s", BytecodeNames[f.meth[f.pc]])
	}
}

// test of GOTO instruction -- in backward direction (to an earlier bytecode)
func TestGotoBackward(t *testing.T) {
	f := newFrame(RETURN)
	f.meth = append(f.meth, GOTO)
	f.meth = append(f.meth, 0xFF) // should be -1
	f.meth = append(f.meth, 0xFF)
	f.meth = append(f.meth, BIPUSH)
	f.pc = 1 // skip over the return instruction to start, catch it on the backward goto
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc] != RETURN {
		t.Errorf("GOTO backeard Expected pc to point to RETURN, but instead it points to : %s", BytecodeNames[f.meth[f.pc]])
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
func TestIfIcmpge2(t *testing.T) {
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

// IF_ICMPLE: if integer compare val 1 <= val 2. Here testing for =
func TestIfIcmple1(t *testing.T) {
	f := newFrame(IF_ICMPLE)
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
		t.Errorf("ICMPLE: expecting a jump to ICONST_2 instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

// ICMPGE: if integer compare val 1 >= val 2. Here test for > (previous test for =)
func TestIfIcmple2(t *testing.T) {
	f := newFrame(IF_ICMPLE)
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
		t.Errorf("IF_ICMPLE: expecting a jump to ICONST_2 instuction, got: %s",
			BytecodeNames[f.pc])
	}
}

// IF_ICMPLE: if integer compare val 1 <>>= val 2 //test when condition fails
func TestIfIcmletFail(t *testing.T) {
	f := newFrame(IF_ICMPLE)
	push(&f, 9)
	push(&f, 8)
	// note that the byte passed in newframe() is at f.meth[0]
	f.meth = append(f.meth, 0) // where we are jumping to, byte 4 = ICONST2
	f.meth = append(f.meth, 4)
	f.meth = append(f.meth, RETURN) // the failed test should drop to this
	f.meth = append(f.meth, ICONST_2)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.meth[f.pc] != RETURN { // b/c we return directly, we don't subtract 1 from pc
		t.Errorf("IF_ICMPLE: expecting fall-through to RETURN instuction, got: %s",
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

// ILOAD: test load of int in locals[index] on to stack
func TestIload(t *testing.T) {
	f := newFrame(ILOAD)
	f.meth = append(f.meth, 0x04) // use local var #4
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234562) // put value in locals[4]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234562 {
		t.Errorf("ILOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("ILOAD: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
	if f.pc != 2 {
		t.Errorf("ILOAD: Expected pc to be pointing at byte 2, got: %d", f.pc)
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

// IRETURN: push an int on to the op stack of the calling method and exit the present method/frame
func TestIreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, 20)
	fs := createFrameStack()
	fs.PushFront(&f0)
	f1 := newFrame(IRETURN)
	push(&f1, 21)
	fs.PushFront(&f1)
	_ = runFrame(fs)
	_ = popFrame(fs)
	f3 := fs.Front().Value.(*frame)
	newVal := pop(f3)
	if newVal != 21 {
		t.Errorf("After IRETURN, expected a value of 21 in previous frame, got: %d", newVal)
	}
	prevVal := pop(f3)
	if prevVal != 20 {
		t.Errorf("After IRETURN, expected a value of 20 in 2nd place of previous frame, got: %d", prevVal)
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
		t.Errorf("ISTORE_0: expected lcoals[0] to be 220, got: %d", f.locals[0])
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
		t.Errorf("ISTORE_1: expected locals[1] to be 221, got: %d", f.locals[1])
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
		t.Errorf("ISTORE_2: expected locals[2] to be 222, got: %d", f.locals[2])
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
		t.Errorf("ISTORE_3: expected locals[3] to be 223, got: %d", f.locals[3])
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

func TestLadd(t *testing.T) {
	f := newFrame(LADD)
	push(&f, 21)
	push(&f, 22)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	value := pop(&f)
	if value != 43 {
		t.Errorf("LADD: expected a result of 43, but got: %d", value)
	}
	if f.tos != -1 {
		t.Errorf("LADD: Expected an empty stack, but got a tos of: %d", f.tos)
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

func TestLconst0(t *testing.T) {
	f := newFrame(LCONST_0)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 0 {
		t.Errorf("LCONST_0: Expected popped value to be 0, got: %d", value)
	}
}

func TestLconst1(t *testing.T) {
	f := newFrame(LCONST_1)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 1 {
		t.Errorf("LCONST_1: Expected popped value to be 1, got: %d", value)
	}
}

// LLOAD: test load of lon in locals[index] on to stack
func TestLload(t *testing.T) {
	f := newFrame(LLOAD)
	f.meth = append(f.meth, 0x04) // use local var #4
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x1234562) // put value in locals[4]

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x1234562 {
		t.Errorf("LLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.tos != -1 {
		t.Errorf("LLOAD: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
	if f.pc != 2 {
		t.Errorf("LLOAD: Expected pc to be pointing at byte 2, got: %d", f.pc)
	}
}

func TestLload0(t *testing.T) {
	f := newFrame(LLOAD_0)

	f.locals = append(f.locals, 0x12345678) // put value in locals[0]
	f.locals = append(f.locals, 0x12345678) // put value in locals[1] // lload uses two local consecutive

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x12345678 {
		t.Errorf("LLOAD_0: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.locals[1] != x {
		t.Errorf("LLOAD_0: Local variable[1] holds invalid value: 0x%x", f.locals[2])
	}

	if f.tos != -1 {
		t.Errorf("LLOAD_0: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

func TestLload1(t *testing.T) {
	f := newFrame(LLOAD_1)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x12345678) // put value in locals[1]
	f.locals = append(f.locals, 0x12345678) // put value in locals[2] // lload uses two local consecutive

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x12345678 {
		t.Errorf("LLOAD_1: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.locals[2] != x {
		t.Errorf("LLOAD_1: Local variable[2] holds invalid value: 0x%x", f.locals[2])
	}

	if f.tos != -1 {
		t.Errorf("LLOAD_1: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

func TestLload2(t *testing.T) {
	f := newFrame(LLOAD_2)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x12345678) // put value in locals[2]
	f.locals = append(f.locals, 0x12345678) // put value in locals[3] // lload uses two local consecutive

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x12345678 {
		t.Errorf("LLOAD_12: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.locals[3] != x {
		t.Errorf("LLOAD_2: Local variable[3] holds invalid value: 0x%x", f.locals[3])
	}

	if f.tos != -1 {
		t.Errorf("LLOAD_1: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

func TestLload3(t *testing.T) {
	f := newFrame(LLOAD_3)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0x12345678) // put value in locals[3]
	f.locals = append(f.locals, 0x12345678) // put value in locals[4] // lload uses two local consecutive

	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	x := pop(&f)
	if x != 0x12345678 {
		t.Errorf("LLOAD_3: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.locals[4] != x {
		t.Errorf("LLOAD_3: Local variable[4] holds invalid value: 0x%x", f.locals[4])
	}

	if f.tos != -1 {
		t.Errorf("LLOAD_3: Expecting an empty stack, but tos points to item: %d", f.tos)
	}
}

// Test LMUL (pop 2 longs, multiply them, push result)
func TestLmul(t *testing.T) {
	f := newFrame(LMUL)
	push(&f, 10)
	push(&f, 7)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("LMUL, Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 70 {
		t.Errorf("LMUL: Expected popped value to be 70, got: %d", value)
	}
}

func TestLstore0(t *testing.T) {
	f := newFrame(LSTORE_0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0) // LSTORE instructions fill two local variables (with the same value)
	push(&f, 0x12345678)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[0] != 0x12345678 {
		t.Errorf("LSTORE_0: expected locals[0] to be 0x12345678, got: %d", f.locals[0])
	}

	if f.locals[1] != 0x12345678 {
		t.Errorf("LSTORE_0: expected locals[1] to be 0x12345678, got: %d", f.locals[1])
	}

	if f.tos != -1 {
		t.Errorf("LSTORE_0: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestLstore1(t *testing.T) {
	f := newFrame(LSTORE_1)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0) // LSTORE instructions fill two local variables (with the same value)
	push(&f, 0x12345678)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[1] != 0x12345678 {
		t.Errorf("LSTORE_1: expected locals[1] to be 0x12345678, got: %d", f.locals[1])
	}

	if f.locals[2] != 0x12345678 {
		t.Errorf("LSTORE_1: expected locals[2] to be 0x12345678, got: %d", f.locals[2])
	}

	if f.tos != -1 {
		t.Errorf("LSTORE_1: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestLstore2(t *testing.T) {
	f := newFrame(LSTORE_2)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0) // LSTORE instructions fill two local variables (with the same value)
	push(&f, 0x12345678)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[2] != 0x12345678 {
		t.Errorf("LSTORE_2: expected locals[2] to be 0x12345678, got: %d", f.locals[2])
	}

	if f.locals[3] != 0x12345678 {
		t.Errorf("LSTORE_2: expected locals[3] to be 0x12345678, got: %d", f.locals[3])
	}

	if f.tos != -1 {
		t.Errorf("LSTORE_2: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

func TestLstore3(t *testing.T) {
	f := newFrame(LSTORE_3)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 0) // LSTORE instructions fill two local variables (with the same value)
	push(&f, 0x12345678)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	if f.locals[3] != 0x12345678 {
		t.Errorf("LSTORE_3: expected locals[3] to be 0x12345678, got: %d", f.locals[3])
	}

	if f.locals[4] != 0x12345678 {
		t.Errorf("LSTORE_3: expected locals[4] to be 0x12345678, got: %d", f.locals[4])
	}

	if f.tos != -1 {
		t.Errorf("LSTORE_3: Expected op stack to be empty, got tos: %d", f.tos)
	}
}

// LSUB: Subtract two longs
func TestLsub(t *testing.T) {
	f := newFrame(LSUB)
	push(&f, 10)
	push(&f, 7)
	fs := createFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.tos != 0 {
		t.Errorf("LSUB, Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 3 {
		t.Errorf("LSUB: Expected popped value to be 3, got: %d", value)
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
	fs := createFrameStack()
	fs.PushFront(&f)
	ret := runFrame(fs)
	if ret == nil {
		t.Errorf("Invalid instruction: Expected an error returned, but got nil.")
	}

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Invalid bytecode") {
		t.Errorf("Error message for invalid bytecode not as expected, got: %s", msg)
	}
}

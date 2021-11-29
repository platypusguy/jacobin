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
	return f
}

// ---- tests ----

func TestBipush(t *testing.T) {
	f := newFrame(BIPUSH)
	f.meth = append(f.meth, 0x05)
	_ = runFrame(&f)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 5 {
		t.Errorf("BIPUSH: Expected popped value to be 5, got: %d", value)
	}
}

func TestIconstN1(t *testing.T) {
	f := newFrame(ICONST_N1)
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 5 {
		t.Errorf("ICONST_5: Expected popped value to be 5, got: %d", value)
	}
}

func TestIinc(t *testing.T) {
	f := newFrame(IINC)
	f.locals = append(f.locals, 0)
	f.locals = append(f.locals, 10) // initialize local variable[1] to 10
	f.meth = append(f.meth, 1)      // increment local variable[1]
	f.meth = append(f.meth, 27)     // increment it by 27
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	_ = runFrame(&f)
	if f.tos != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.tos)
	}
	value := pop(&f)
	if value != 27 {
		t.Errorf("ILOAD_3: Expected popped value to be 27, got: %d", value)
	}
}

func TestIsub(t *testing.T) {
	f := newFrame(ISUB)
	push(&f, 10)
	push(&f, 7)
	_ = runFrame(&f)
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
	_ = runFrame(&f)
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
	ret := runFrame(&f)
	if f.tos != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.tos)
	}

	if ret != nil {
		t.Error("RETURN: Expected popped value to be 2, got: " + ret.Error())
	}
}

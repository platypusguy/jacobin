/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
    "io"
    "jacobin/classloader"
    "jacobin/frames"
    "jacobin/globals"
    "jacobin/log"
    "jacobin/thread"
    "math"
    "os"
    "strings"
    "testing"
    "unsafe"
)

// These tests test the individual bytecode instructions. They are presented here in
// alphabetical order of the instruction name.
// Note: Test for bytecodes related to array operations are located in arrays_test.go

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
    f := newFrame(ACONST_NULL)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    x := peek(&f).(int64)
    if x != 0 {
        t.Errorf("ACONST_NULL: Expecting 0 on stack, got: %d", x)
    }
    if f.TOS != 0 {
        t.Errorf("ACONST_NULL: Expecting TOS = 0, but tos is: %d", f.TOS)
    }
}

// ALOAD: test load of reference in locals[index] on to stack
func TestAload(t *testing.T) {
    f := newFrame(ALOAD)
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
    f := newFrame(ALOAD_0)
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
    f := newFrame(ALOAD_1)
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
    f := newFrame(ALOAD_2)
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
    f := newFrame(ALOAD_3)
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
    f1 := newFrame(ARETURN)
    push(&f1, unsafe.Pointer(&f1))
    fs.PushFront(&f1)
    _ = runFrame(fs)
    _ = frames.PopFrame(fs)
    f3 := fs.Front().Value.(*frames.Frame)
    newVal := pop(f3).(unsafe.Pointer)
    if newVal != unsafe.Pointer(&f1) {
        t.Error("ARETURN: did not get expected value of reference")
    }
}

// ASTORE: Store reference in local var specified by following byte.
func TestAstore(t *testing.T) {
    f := newFrame(ASTORE)
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
    f := newFrame(ASTORE_0)
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
    f := newFrame(ASTORE_1)
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
    f := newFrame(ASTORE_2)
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
    f := newFrame(ASTORE_3)
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
    f := newFrame(BIPUSH)
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

// D2F: test convert double to float
func TestD2f(t *testing.T) {
    f := newFrame(D2F)
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
    f := newFrame(D2I)
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
    f := newFrame(D2I)
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
    f := newFrame(D2L)
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
    f := newFrame(D2L)
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
    f := newFrame(DADD)
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
    f := newFrame(DADD)
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
    f := newFrame(DADD)
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
    f := newFrame(DCMPG)
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
        t.Errorf("DDIV: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
    }
}

// DCMP0: compare two doubles
func TestDcmpgMinus1(t *testing.T) {
    f := newFrame(DCMPG)
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
        t.Errorf("DDIV: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
    }
}

// DCMP0: compare two doubles
func TestDcmpg0(t *testing.T) {
    f := newFrame(DCMPG)
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
        t.Errorf("DDIV: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
    }
}

func TestDcmpgNan(t *testing.T) {
    f := newFrame(DCMPG)
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
        t.Errorf("DDIV: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
    }
}

// DCMPL
func TestDcmplNan(t *testing.T) {
    f := newFrame(DCMPL)
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
        t.Errorf("DDIV: Expected stack with 0 items, but got a TOS of: %d", f.TOS)
    }
}

// DCONST_0
func TestDconst0(t *testing.T) {
    f := newFrame(DCONST_0)
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
    f := newFrame(DCONST_1)
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
    f := newFrame(DDIV)
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

// DLOAD: test load of double in locals[index] on to stack
func TestDload(t *testing.T) {
    f := newFrame(DLOAD)
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
    f := newFrame(DLOAD_0)
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
    f := newFrame(DLOAD_1)
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
    f := newFrame(DLOAD_2)
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
    f := newFrame(DLOAD_3)
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

// Test DMUL (pop 2 doubles, multiply them, push result)
func TestDmul(t *testing.T) {
    f := newFrame(DMUL)
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

// Test DNEG Negate a double
func TestDneg(t *testing.T) {
    f := newFrame(DNEG)
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

// Test DNEG Negate a double - infinity
func TestDnegInf(t *testing.T) {
    f := newFrame(DNEG)
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
    f := newFrame(DREM)
    push(&f, 23.5)
    push(&f, 23.5)
    push(&f, 3.3)
    push(&f, 3.3)

    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)

    if f.TOS != 0 {
        t.Errorf("DREM, Top of stack, expected 0, got: %d", f.TOS)
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
    f1 := newFrame(DRETURN)
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
    f := newFrame(DSTORE)
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
    f := newFrame(DSTORE_0)
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
    f := newFrame(DSTORE_1)
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
    f := newFrame(DSTORE_2)
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
    f := newFrame(DSTORE_3)
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
    f := newFrame(DSUB)
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
    f := newFrame(DUP)
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
    f := newFrame(DUP2)
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
    f := newFrame(DUP_X1)
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
    f := newFrame(DUP_X2)
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
    f := newFrame(DUP2_X1)
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
    f := newFrame(DUP2_X2)
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
    f := newFrame(F2D)
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
    f := newFrame(F2I)
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
    f := newFrame(F2I)
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
    f := newFrame(F2L)
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
    f := newFrame(FADD)
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

// FCONST_0
func TestFconst0(t *testing.T) {
    f := newFrame(FCONST_0)
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
    f := newFrame(FCONST_1)
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
    f := newFrame(FCONST_2)
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
    f := newFrame(FDIV)
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

// FLOAD: test load of float in locals[index] on to stack
func TestFload(t *testing.T) {
    f := newFrame(FLOAD)
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
    f := newFrame(FLOAD_0)
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
    f := newFrame(FLOAD_1)
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
    f := newFrame(FLOAD_2)
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
    f := newFrame(FLOAD_3)
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

// Test FMUL (pop 2 floats, multiply them, push result)
func TestFmul(t *testing.T) {
    f := newFrame(FMUL)
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
    f := newFrame(FNEG)
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
    f := newFrame(FREM)
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
    f := newFrame(FSTORE)
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
    f := newFrame(FSTORE_0)
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
    f := newFrame(FSTORE_1)
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
    f := newFrame(FSTORE_2)
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
    f := newFrame(FSTORE_3)
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
    f := newFrame(FSUB)
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

// GOTO: in forward direction (to a later bytecode)
func TestGotoForward(t *testing.T) {
    f := newFrame(GOTO)
    f.Meth = append(f.Meth, 0x00)
    f.Meth = append(f.Meth, 0x03)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, NOP)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN {
        t.Errorf("GOTO forward: Expected pc to point to RETURN, but instead it points to : %s", BytecodeNames[f.Meth[f.PC]])
    }
}

// GOTO: go to instruction in backward direction (to an earlier bytecode)
func TestGotoBackward(t *testing.T) {
    f := newFrame(RETURN)
    f.Meth = append(f.Meth, GOTO)
    f.Meth = append(f.Meth, 0xFF) // should be -1
    f.Meth = append(f.Meth, 0xFF)
    f.Meth = append(f.Meth, BIPUSH)
    f.PC = 1 // skip over the return instruction to start, catch it on the backward goto
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN {
        t.Errorf("GOTO backeard Expected pc to point to RETURN, but instead it points to : %s", BytecodeNames[f.Meth[f.PC]])
    }
}

// I2B: convert int to Java char (16-bit value)
func TestI2B(t *testing.T) {
    f := newFrame(I2B)
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
    f := newFrame(I2B)
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
    f := newFrame(I2C)
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
    f := newFrame(I2D)
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
    f := newFrame(I2F)
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
    f := newFrame(I2L)
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
    f := newFrame(I2S)
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
    f := newFrame(IADD)
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
    f := newFrame(IAND)
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
    f := newFrame(IDIV)
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
func TestIdivDivideByZero(t *testing.T) {
    g := globals.GetGlobalRef()
    globals.InitGlobals("test")
    // g.Threads = list.New()
    g.JacobinName = "test" // prevents a shutdown when the exception hits.
    log.Init()

    // redirect stderr & stdout to capture results from stderr
    normalStderr := os.Stderr
    r, w, _ := os.Pipe()
    os.Stderr = w

    normalStdout := os.Stdout
    _, wout, _ := os.Pipe()
    os.Stdout = wout

    f := newFrame(IDIV)
    f.ClName = "testClass"
    f.MethName = "testMethod"
    push(&f, int64(220))
    push(&f, int64(0))
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame

    // need to create a thread to catch the exception
    hread := thread.CreateThread()
    hread.Stack = fs
    hread.ID = thread.AddThreadToTable(&hread, &g.Threads)
    _ = runFrame(fs)

    // restore stderr and stdout to what they were before
    _ = w.Close()
    out, _ := io.ReadAll(r)
    os.Stderr = normalStderr

    errMsg := string(out[:])

    _ = wout.Close()
    os.Stdout = normalStdout

    if !strings.Contains(errMsg, "Arithmetic Exception") {
        t.Errorf("IDIV: Did not get expected error msg, got: %s", errMsg)
    }
}

// ICONST:
func TestIconstN1(t *testing.T) {
    f := newFrame(ICONST_N1)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.TOS != 0 {
        t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
    }
    var value = pop(&f).(int64)
    if value != -1 {
        t.Errorf("ICONST_N1: Expected popped value to be -1, got: %d", value)
    }
}

// ICONST_0
func TestIconst0(t *testing.T) {
    f := newFrame(ICONST_0)
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
    f := newFrame(ICONST_1)
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
    f := newFrame(ICONST_2)
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
    f := newFrame(ICONST_3)
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
    f := newFrame(ICONST_4)
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
    f := newFrame(ICONST_5)
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
    f := newFrame(IF_ACMPEQ)
    push(&f, int64(0xFF8899))
    push(&f, int64(0xFF8899))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IF_ACMPEQ: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ACMPEQ: jump if two addresses are equal (this tests addresses being unequal)
func TestIfAcmpeqFail(t *testing.T) {
    f := newFrame(IF_ACMPEQ)
    push(&f, int64(0xFF8899))
    push(&f, int64(0xFF889A))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST_2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN) // the failed test should drop to this
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN { // b/c we return directly, we don't subtract 1 from pc
        t.Errorf("IF_ICMPEQ: expecting fall-through to RETURN instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ACMPNE: jump if two addresses are not equal
func TestIfAcmpNe(t *testing.T) {
    f := newFrame(IF_ACMPNE)
    push(&f, int64(0xFF8899))
    push(&f, int64(0xFF889A))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IF_ACMPNE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ACMPNE: jump if two addresses are equal (this tests addresses being equal)
func TestIfAcmpneFail(t *testing.T) {
    f := newFrame(IF_ACMPNE)
    push(&f, int64(0xFF8899))
    push(&f, int64(0xFF8899))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST_2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN) // the failed test should drop to this
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN { // b/c we return directly, we don't subtract 1 from pc
        t.Errorf("IF_ICMPNE: expecting fall-through to RETURN instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPEQ: jump if val1 == val2 (both ints, both popped off stack)
func TestIfIcmpeq(t *testing.T) {
    f := newFrame(IF_ICMPEQ)
    push(&f, int64(9)) // pushed two equal values, so jump should be made.
    push(&f, int64(9))

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("ICMPEQ: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_CMPGE: if integer compare val 1 >= val 2. Here test for = (next test for >)
func TestIfIcmpge1(t *testing.T) {
    f := newFrame(IF_ICMPGE)
    push(&f, int64(9))
    push(&f, int64(9))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("ICMPGE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPGE: if integer compare val 1 >= val 2. Here test for > (previous test for =)
func TestIfIcmpge2(t *testing.T) {
    f := newFrame(IF_ICMPGE)
    push(&f, int64(9))
    push(&f, int64(8))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("ICMPGE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPGE: if integer compare val 1 >= val 2 //test when condition fails
func TestIfIcmgetFail(t *testing.T) {
    f := newFrame(IF_ICMPGE)
    push(&f, int64(8))
    push(&f, int64(9))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN) // the failed test should drop to this
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN { // b/c we return directly, we don't subtract 1 from pc
        t.Errorf("ICMPGE: expecting fall-through to RETURN instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPGE: if integer compare val 1 >= val 2. Here test for > (previous test for =)
func TestIfIcmple2(t *testing.T) {
    f := newFrame(IF_ICMPLE)
    push(&f, int64(8))
    push(&f, int64(9))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IF_ICMPLE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPGT: jump if val1 > val2 (both ints, both popped off stack)
func TestIfIcmpgt(t *testing.T) {
    f := newFrame(IF_ICMPGT)
    push(&f, int64(9)) // val1 > val2, so jump should be made.
    push(&f, int64(8))

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IF_ICMPNE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPLE: if integer compare val 1 ! <= val 2 //test when condition fails
func TestIfIcmpletFail(t *testing.T) {
    f := newFrame(IF_ICMPLE)
    push(&f, int64(9))
    push(&f, int64(8))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN) // the failed test should drop to this
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN { // b/c we return directly, we don't subtract 1 from pc
        t.Errorf("IF_ICMPLE: expecting fall-through to RETURN instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPLE: if integer compare val 1 <= val 2. Here testing for =
func TestIfIcmple1(t *testing.T) {
    f := newFrame(IF_ICMPLE)
    push(&f, int64(9))
    push(&f, int64(9))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("ICMPLE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPLT: if integer compare val 1 < val 2
func TestIfIcmplt(t *testing.T) {
    f := newFrame(IF_ICMPLT)
    push(&f, int64(8))
    push(&f, int64(9))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, ICONST_1)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("ICMPLT: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPLT: if integer compare val 1 < val 2 //test when condition fails
func TestIfIcmpltFail(t *testing.T) {
    f := newFrame(IF_ICMPLT)
    push(&f, int64(9))
    push(&f, int64(9))
    // note that the byte passed in newframe() is at f.Meth[0]
    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN) // the failed test should drop to this
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC] != RETURN { // b/c we return directly, we don't subtract 1 from pc
        t.Errorf("ICMPLT: expecting fall-through to RETURN instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IF_ICMPNE: jump if val1 != val2 (both ints, both popped off stack)
func TestIfIcmpne(t *testing.T) {
    f := newFrame(IF_ICMPNE)
    push(&f, int64(9)) // pushed two unequal values, so jump should be made.
    push(&f, int64(8))

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IF_ICMPNE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFEQ: jump if int popped off TOS is = 0
func TestIfeq(t *testing.T) {
    f := newFrame(IFEQ)
    push(&f, int64(0)) // pushed 0, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFEQ: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFEQ: jump if int popped off TOS is = 0; here != 0
func TestIfeqFallThrough(t *testing.T) {
    f := newFrame(IFEQ)
    push(&f, int64(23)) // pushed 23, so jump should not be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFEQ: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFGE: jump if int popped off TOS is >= 0
func TestIfge(t *testing.T) {
    f := newFrame(IFGE)
    push(&f, int64(66)) // pushed 66, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFGE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFGE: jump if int popped off TOS is >= 0, here = 0
func TestIfgeEqual0(t *testing.T) {
    f := newFrame(IFGE)
    push(&f, int64(0)) // pushed 0, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFGE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFGE: jump if int popped off TOS is >= 0; here < 0
func TestIfgeFallThrough(t *testing.T) {
    f := newFrame(IFGE)
    push(&f, int64(-1)) // pushed -1, so jump should not be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFGE: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFGT: jump if int popped off TOS is > 0
func TestIfgt(t *testing.T) {
    f := newFrame(IFGT)
    push(&f, int64(66)) // pushed 66, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFGT: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFGT: jump if int popped off TOS is > 0; here = 0
func TestIfgtFallThrough(t *testing.T) {
    f := newFrame(IFGT)
    push(&f, int64(0)) // pushed 0, so jump should not be made.

    f.Meth = append(f.Meth, 0)
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFGT: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFLE: jump if int popped off TOS is <= 0
func TestIfle(t *testing.T) {
    f := newFrame(IFLE)
    push(&f, int64(-66)) // pushed -66, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFLE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFLE: jump if int popped off TOS is <= 0; here = 0
func TestIfleTest0(t *testing.T) {
    f := newFrame(IFLE)
    push(&f, int64(0)) // pushed 0, so jump should be made.

    f.Meth = append(f.Meth, 0)
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFLE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFLE: jump if int popped off TOS is <= 0; here > 0, so no jump
func TestIfleFallThrough(t *testing.T) {
    f := newFrame(IFLE)
    push(&f, int64(66)) // pushed 66, so jump should not be made.

    f.Meth = append(f.Meth, 0)
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFLE: Invalid jump when expecting fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFLT: jump if int popped off TOS is < 0
func TestIflt(t *testing.T) {
    f := newFrame(IFLT)
    push(&f, int64(-66)) // pushed -66, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFLT: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFLT: jump if int popped off TOS is < 0; here = 0
func TestIfltFallThrough(t *testing.T) {
    f := newFrame(IFLT)
    push(&f, int64(0)) // pushed 0, so jump should not be made.

    f.Meth = append(f.Meth, 0)
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFLT: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFNE: jump if int popped off TOS is != 0
func TestIfne(t *testing.T) {
    f := newFrame(IFNE)
    push(&f, int64(1)) // pushed 1, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFNE: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFNE: jump if int popped off TOS is != 0; here it is = 0
func TestIfneFallThrough(t *testing.T) {
    f := newFrame(IFNE)
    push(&f, int64(0)) // pushed 0, so jump should not be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFNE: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFNONNULL: jump if TOS holds a non-null address
func TestIfn0nnull(t *testing.T) {
    f := newFrame(IFNONNULL)
    push(&f, int64(1)) // pushed 1, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFNONNULL: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFNONNULL: jump if TOS holds a non-null address; here it is null
func TestIfnonnullFallThrough(t *testing.T) {
    f := newFrame(IFNONNULL)
    push(&f, int64(0)) // pushed 0, so jump should not be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFNONNULL: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFNULL: jump if TOS holds null address
func TestIfnull(t *testing.T) {
    f := newFrame(IFNULL)
    push(&f, int64(0)) // pushed null, so jump should be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, NOP)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] != ICONST_2 { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFNULL: expecting a jump to ICONST_2 instuction, got: %s",
            BytecodeNames[f.PC])
    }
}

// IFNULL: jump if TOS address is null; here not null
func TestIfnullFallThrough(t *testing.T) {
    f := newFrame(IFNULL)
    push(&f, int64(23)) // pushed 23, so jump should not be made.

    f.Meth = append(f.Meth, 0) // where we are jumping to, byte 4 = ICONST2
    f.Meth = append(f.Meth, 4)
    f.Meth = append(f.Meth, RETURN)
    f.Meth = append(f.Meth, ICONST_2)
    fs := frames.CreateFrameStack()
    fs.PushFront(&f) // push the new frame
    _ = runFrame(fs)
    if f.Meth[f.PC-1] == IFNULL { // -1 b/c the run loop adds 1 before exiting
        t.Errorf("IFNULL: Invalid fall-through, got: %s",
            BytecodeNames[f.PC])
    }
}

// IINC:
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

// Test LDC_W: get float64 CP entry indexed by two bytes
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

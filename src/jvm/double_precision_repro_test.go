package jvm

import (
	"fmt"
	"jacobin/src/frames"
	"jacobin/src/opcodes"
	"testing"
)

func TestDoublePrecision(t *testing.T) {
	// DADD precision test
	// 1.0 + 1e-16 should be different from 1.0 in double, but same in float32
	f := newFrame(opcodes.DADD)
	v1 := 1.0
	v2 := 1e-15
	push(&f, v1)
	push(&f, v2)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	res := pop(&f).(float64)
	fmt.Printf("DADD result: %v\n", res)
	if res == 1.0 {
		t.Errorf("DADD: lost precision, 1.0 + 1e-15 resulted in 1.0")
	}

	// DREM precision test
	f = newFrame(opcodes.DREM)
	push(&f, 1.0)
	push(&f, 0.3)
	fs = frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	res = pop(&f).(float64)
	fmt.Printf("DREM result: %v\n", res)
	// 1.0 % 0.3 should be 0.1 (roughly)
	if float32(res) != float32(1.0-0.9) {
		t.Errorf("DREM: unexpected result %v", res)
	}
}

func TestD2fConversion(t *testing.T) {
	f := newFrame(opcodes.D2F)
	// A value that changes when converted to float32
	v := 1.0000000000000002
	push(&f, v)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	res := pop(&f).(float64)
	fmt.Printf("D2F result: %v\n", res)
	if res == v {
		t.Errorf("D2F: failed to round to float32 precision")
	}
	if res != float64(float32(v)) {
		t.Errorf("D2F: expected %v, got %v", float64(float32(v)), res)
	}
}

func TestF2dConversion(t *testing.T) {
	// This is more subtle. F2D should not lose any more precision than it already has.
	// If we have a float32, converting it to double should preserve the float32 value.
	// In Jacobin, everything is float64 on stack.
	// FADD should produce a rounded-to-float32 float64.

	f := newFrame(opcodes.FADD)
	v1 := float64(float32(1.0))
	v2 := float64(float32(1e-16))
	push(&f, v1)
	push(&f, v2)

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	res := pop(&f).(float64)
	if res != 1.0 {
		t.Errorf("FADD: expected 1.0 due to float32 precision limits, got %v", res)
	}
}

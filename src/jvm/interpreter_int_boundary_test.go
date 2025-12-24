/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/types"
	"math"
	"testing"
)

func TestIaddBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v1       int32
		v2       int32
		expected int32
	}{
		{"MaxInt + 1", math.MaxInt32, 1, math.MinInt32},
		{"MinInt - 1", math.MinInt32, -1, math.MaxInt32},
		{"MaxInt + MaxInt", math.MaxInt32, math.MaxInt32, -2},
		{"MinInt + MinInt", math.MinInt32, math.MinInt32, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.IADD)
			push(&f, int64(tt.v1))
			push(&f, int64(tt.v2))
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestIsubBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v1       int32
		v2       int32
		expected int32
	}{
		{"MinInt - 1", math.MinInt32, 1, math.MaxInt32},
		{"MaxInt - (-1)", math.MaxInt32, -1, math.MinInt32},
		{"MinInt - MaxInt", math.MinInt32, math.MaxInt32, 1},
		{"MaxInt - MinInt", math.MaxInt32, math.MinInt32, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.ISUB)
			push(&f, int64(tt.v1))
			push(&f, int64(tt.v2))
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestImulBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v1       int32
		v2       int32
		expected int32
	}{
		{"MaxInt * 2", math.MaxInt32, 2, -2},
		{"MinInt * -1", math.MinInt32, -1, math.MinInt32},
		{"MaxInt * MaxInt", math.MaxInt32, math.MaxInt32, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.IMUL)
			push(&f, int64(tt.v1))
			push(&f, int64(tt.v2))
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestInegBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int32
		expected int32
	}{
		{"MinInt", math.MinInt32, math.MinInt32},
		{"MaxInt", math.MaxInt32, -math.MaxInt32},
		{"Zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.INEG)
			push(&f, int64(tt.v))
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestIshlBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int32
		dist     int64
		expected int32
	}{
		{"1 << 31", 1, 31, math.MinInt32},
		{"1 << 32", 1, 32, 1},
		{"1 << -1", 1, -1, math.MinInt32}, // -1 & 0x1F = 31
		{"-1 << 1", -1, 1, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.ISHL)
			push(&f, int64(tt.v))
			push(&f, tt.dist)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestIshrBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int32
		dist     int64
		expected int32
	}{
		{"MinInt >> 1", math.MinInt32, 1, -1073741824},
		{"MinInt >> 31", math.MinInt32, 31, -1},
		{"-1 >> 31", -1, 31, -1},
		{"1 >> 1", 1, 1, 0},
		{"1 >> 32", 1, 32, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.ISHR)
			push(&f, int64(tt.v))
			push(&f, tt.dist)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestIushrBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int32
		dist     int64
		expected int32
	}{
		{"-1 >>> 1", -1, 1, 2147483647},
		{"MinInt >>> 1", math.MinInt32, 1, 1073741824},
		{"-1 >>> 31", -1, 31, 1},
		{"1 >>> 32", 1, 32, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.IUSHR)
			push(&f, int64(tt.v))
			push(&f, tt.dist)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if int32(res) != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, int32(res))
			}
		})
	}
}

func TestIbitwiseBoundary(t *testing.T) {
	t.Run("IAND Neg", func(t *testing.T) {
		f := newFrame(opcodes.IAND)
		push(&f, int64(-1))
		push(&f, int64(123))
		fs := frames.CreateFrameStack()
		fs.PushFront(&f)
		interpret(fs)
		res := pop(&f).(int64)
		if int32(res) != 123 {
			t.Errorf("IAND: expected 123, got %d", int32(res))
		}
	})
	t.Run("IOR Neg", func(t *testing.T) {
		f := newFrame(opcodes.IOR)
		push(&f, int64(-1))
		push(&f, int64(123))
		fs := frames.CreateFrameStack()
		fs.PushFront(&f)
		interpret(fs)
		res := pop(&f).(int64)
		if int32(res) != -1 {
			t.Errorf("IOR: expected -1, got %d", int32(res))
		}
	})
	t.Run("IXOR Neg", func(t *testing.T) {
		f := newFrame(opcodes.IXOR)
		push(&f, int64(-1))
		push(&f, int64(0))
		fs := frames.CreateFrameStack()
		fs.PushFront(&f)
		interpret(fs)
		res := pop(&f).(int64)
		if int32(res) != -1 {
			t.Errorf("IXOR: expected -1, got %d", int32(res))
		}
	})
}

func TestI2bBoundary(t *testing.T) {
	tests := []struct {
		v        int32
		expected int8
	}{
		{127, 127},
		{128, -128},
		{-1, -1},
		{-128, -128},
		{-129, 127},
		{255, -1},
		{256, 0},
	}

	for _, tt := range tests {
		f := newFrame(opcodes.I2B)
		push(&f, int64(tt.v))
		fs := frames.CreateFrameStack()
		fs.PushFront(&f)
		interpret(fs)
		res := pop(&f).(int64)
		if int8(res) != tt.expected {
			t.Errorf("I2B(%d): expected %d, got %d", tt.v, tt.expected, int8(res))
		}
	}
}

func TestI2cBoundary(t *testing.T) {
	tests := []struct {
		v        int32
		expected uint16
	}{
		{65535, 65535},
		{65536, 0},
		{-1, 65535},
	}

	for _, tt := range tests {
		f := newFrame(opcodes.I2C)
		push(&f, int64(tt.v))
		fs := frames.CreateFrameStack()
		fs.PushFront(&f)
		interpret(fs)
		res := pop(&f).(int64)
		if uint16(res) != tt.expected {
			t.Errorf("I2C(%d): expected %d, got %d", tt.v, tt.expected, uint16(res))
		}
	}
}

func TestI2sBoundary(t *testing.T) {
	tests := []struct {
		v        int32
		expected int16
	}{
		{32767, 32767},
		{32768, -32768},
		{-1, -1},
		{65535, -1},
		{65536, 0},
	}

	for _, tt := range tests {
		f := newFrame(opcodes.I2S)
		push(&f, int64(tt.v))
		fs := frames.CreateFrameStack()
		fs.PushFront(&f)
		interpret(fs)
		res := pop(&f).(int64)
		if int16(res) != tt.expected {
			t.Errorf("I2S(%d): expected %d, got %d", tt.v, tt.expected, int16(res))
		}
	}
}

func TestLaddBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v1       int64
		v2       int64
		expected int64
	}{
		{"MaxLong + 1", math.MaxInt64, 1, math.MinInt64},
		{"MinLong - 1", math.MinInt64, -1, math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.LADD)
			push(&f, tt.v1)
			push(&f, tt.v2)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if res != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, res)
			}
		})
	}
}

func TestLshlBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int64
		dist     int64
		expected int64
	}{
		{"1 << 63", 1, 63, math.MinInt64},
		{"1 << 64", 1, 64, 1},
		{"1 << -1", 1, -1, math.MinInt64}, // -1 & 0x3F = 63
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.LSHL)
			push(&f, tt.v)
			push(&f, tt.dist)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if res != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, res)
			}
		})
	}
}

func TestLshrBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int64
		dist     int64
		expected int64
	}{
		{"MinLong >> 1", math.MinInt64, 1, -4611686018427387904},
		{"MinLong >> 63", math.MinInt64, 63, -1},
		{"-1 >> 63", -1, 63, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.LSHR)
			push(&f, tt.v)
			push(&f, tt.dist)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if res != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, res)
			}
		})
	}
}

func TestLushrBoundary(t *testing.T) {
	tests := []struct {
		name     string
		v        int64
		dist     int64
		expected int64
	}{
		{"-1 >>> 1", -1, 1, 9223372036854775807},
		{"MinLong >>> 1", math.MinInt64, 1, 4611686018427387904},
		{"-1 >>> 63", -1, 63, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFrame(opcodes.LUSHR)
			push(&f, tt.v)
			push(&f, tt.dist)
			fs := frames.CreateFrameStack()
			fs.PushFront(&f)
			interpret(fs)
			res := pop(&f).(int64)
			if res != tt.expected {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, res)
			}
		})
	}
}

func TestIaloadBoundary(t *testing.T) {
	globals.InitGlobals("test")

	// IALOAD
	f := newFrame(opcodes.IALOAD)
	arr := []int64{1, 2, 3}
	obj := object.MakePrimitiveObject(types.IntArray, types.IntArray, arr)
	push(&f, obj)
	push(&f, int64(1))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	res := pop(&f).(int64)
	if int32(res) != 2 {
		t.Errorf("IALOAD: expected 2, got %d", int32(res))
	}

	// SALOAD (short array load, should sign extend)
	f = newFrame(opcodes.SALOAD)
	sarr := []int64{-1, 2, 3}
	sobj := object.MakePrimitiveObject(types.ShortArray, types.ShortArray, sarr)
	push(&f, sobj)
	push(&f, int64(0))
	fs = frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	res = pop(&f).(int64)
	if int32(res) != -1 {
		t.Errorf("SALOAD: expected -1, got %d", int32(res))
	}

	// CALOAD (char array load, should zero extend)
	f = newFrame(opcodes.CALOAD)
	carr := []int64{0xFFFF, 2, 3}
	cobj := object.MakePrimitiveObject(types.CharArray, types.CharArray, carr)
	push(&f, cobj)
	push(&f, int64(0))
	fs = frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	res = pop(&f).(int64)
	if int32(res) != 65535 {
		t.Errorf("CALOAD: expected 65535, got %d", int32(res))
	}
}

func TestIastoreBoundary(t *testing.T) {
	globals.InitGlobals("test")

	// IASTORE
	f := newFrame(opcodes.IASTORE)
	arr := []int64{0, 0, 0}
	obj := object.MakePrimitiveObject(types.IntArray, types.IntArray, arr)
	push(&f, obj)
	push(&f, int64(1))
	push(&f, int64(-123))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	if arr[1] != -123 {
		t.Errorf("IASTORE: expected -123, got %d", arr[1])
	}

	// SASTORE
	f = newFrame(opcodes.SASTORE)
	sarr := []int64{0, 0, 0}
	sobj := object.MakePrimitiveObject(types.ShortArray, types.ShortArray, sarr)
	push(&f, sobj)
	push(&f, int64(2))
	push(&f, int64(0x12345678)) // should be truncated to 0x5678
	fs = frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	if sarr[2] != 0x5678 {
		t.Errorf("SASTORE: expected 0x5678, got 0x%x", sarr[2])
	}

	// CASTORE
	f = newFrame(opcodes.CASTORE)
	carr := []int64{0, 0, 0}
	cobj := object.MakePrimitiveObject(types.CharArray, types.CharArray, carr)
	push(&f, cobj)
	push(&f, int64(0))
	push(&f, int64(-1)) // should be truncated to 0xFFFF
	fs = frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)
	if carr[0] != 0xFFFF {
		t.Errorf("CASTORE: expected 0xFFFF, got 0x%x", carr[0])
	}
}

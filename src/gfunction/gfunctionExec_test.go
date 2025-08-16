/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"errors"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"testing"
)

// helper to build a minimal frame stack with one frame at the front
func makeFrameStack() *list.List {
	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(0)
	f.Thread = 1 // main thread for deterministic messages
	f.ClName = "java/lang/Test"
	f.MethName = "m"
	f.MethType = "()V"
	_ = frames.PushFrame(fs, f)
	return fs
}

func TestRunGfunction_ParamOrderAndContext(t *testing.T) {
	globals.InitGlobals("test")

	fs := makeFrameStack()

	// Capture the received params in the GFunction
	var received []interface{}
	gm := GMeth{
		ParamSlots:   2,
		NeedsContext: true,
		GFunction: func(in []interface{}) interface{} {
			// save a copy of the slice
			received = append([]interface{}{}, in...)
			return nil
		},
	}
	mt := classloader.MTentry{Meth: gm, MType: 'G'}

	// params in forward order as they would be pushed by bytecode: last arg at top
	p := []interface{}{int64(1), int64(2)}
	objRef := false

	_ = RunGfunction(mt, fs, "pkg/Clazz", "method", "(II)V", &p, objRef, false)

	if len(received) != 3 { // two args + context
		t.Fatalf("expected 3 params passed to GFunction (2 args + context), got %d", len(received))
	}

	// Because NeedsContext=true, fs is appended then the slice is reversed, so fs should be first
	if received[0] != fs {
		t.Fatalf("expected context frame stack as first param after reversal; got %T", received[0])
	}
	// The remaining args should be reversed order of original
	if v, ok := received[1].(int64); !ok || v != int64(2) {
		t.Fatalf("expected second param to be 2, got %v (%T)", received[1], received[1])
	}
	if v, ok := received[2].(int64); !ok || v != int64(1) {
		t.Fatalf("expected third param to be 1, got %v (%T)", received[2], received[2])
	}
}

func TestRunGfunction_ReturnsValue_NonThreadSafe(t *testing.T) {
	globals.InitGlobals("test")

	fs := makeFrameStack()

	gm := GMeth{
		ParamSlots: 1,
		GFunction: func(in []interface{}) interface{} {
			if len(in) != 1 || in[0] != "x" {
				t.Fatalf("unexpected params to GFunction: %#v", in)
			}
			return int64(42)
		},
	}
	mt := classloader.MTentry{Meth: gm, MType: 'G'}

	p := []interface{}{"x"}
	ret := RunGfunction(mt, fs, "A/B", "c", "(Ljava/lang/String;)I", &p, false, false)

	if v, ok := ret.(int64); !ok || v != 42 {
		t.Fatalf("expected int64(42) return, got %v (%T)", ret, ret)
	}
}

func TestRunGfunction_ReturnsError_PassesThrough(t *testing.T) {
	globals.InitGlobals("test")

	fs := makeFrameStack()

	gm := GMeth{GFunction: func([]interface{}) interface{} {
		return errors.New("native boom")
	}}
	mt := classloader.MTentry{Meth: gm, MType: 'G'}

	var nilParams []interface{}
	ret := RunGfunction(mt, fs, "P/Q", "r", "()V", &nilParams, false, false)

	if err, ok := ret.(error); !ok {
		t.Fatalf("expected error return, got %T: %v", ret, ret)
	} else if err.Error() != "native boom" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestRunGfunction_GErrBlk_ReturnsErrorInTestMode(t *testing.T) {
	globals.InitGlobals("test")

	fs := makeFrameStack()

	gm := GMeth{GFunction: func([]interface{}) interface{} {
		return &GErrBlk{ExceptionType: excNames.ArrayIndexOutOfBoundsException, ErrMsg: "array oob"}
	}}
	mt := classloader.MTentry{Meth: gm, MType: 'G'}

	params := []interface{}{}
	ret := RunGfunction(mt, fs, "X/Y", "z", "()V", &params, false, false)

	if err, ok := ret.(error); !ok {
		t.Fatalf("expected error return for GErrBlk in test mode, got %T: %v", ret, ret)
	} else {
		// The error message should contain our method FQN and original message
		if !contains(err.Error(), "array oob") || !contains(err.Error(), "X/Y.z()V") {
			t.Fatalf("error message missing expected content: %q", err.Error())
		}
	}
}

func TestRunGfunction_ThreadSafe_WithObjRef(t *testing.T) {
	globals.InitGlobals("test")

	fs := makeFrameStack()

	// The thread-safe function will simply return a constant
	gm := GMeth{
		ThreadSafe: true,
		GFunction: func(in []interface{}) interface{} {
			// First param must be the object reference
			if len(in) == 0 {
				t.Fatalf("expected at least one param for thread-safe call")
			}
			if _, ok := in[0].(*object.Object); !ok {
				t.Fatalf("first parameter to thread-safe GFunction not an *object.Object: %T", in[0])
			}
			return "ok"
		},
	}
	mt := classloader.MTentry{Meth: gm, MType: 'G'}

	obj := object.MakeEmptyObject()
	params := []interface{}{obj}

	ret := RunGfunction(mt, fs, "TS/C", "m", "()V", &params, true, false)

	if s, ok := ret.(string); !ok || s != "ok" {
		t.Fatalf("expected \"ok\" return from thread-safe GFunction, got %v (%T)", ret, ret)
	}
}

// contains is a tiny helper to avoid importing strings just for Contains
func contains(haystack, needle string) bool {
	return len(needle) == 0 || (len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0)
}

// indexOf naive substring search
func indexOf(s, sub string) int {
	// very small helper; ok for tests
	outer := []rune(s)
	inner := []rune(sub)
	if len(inner) == 0 {
		return 0
	}
	for i := 0; i+len(inner) <= len(outer); i++ {
		match := true
		for j := 0; j < len(inner); j++ {
			if outer[i+j] != inner[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

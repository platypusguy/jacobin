/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"container/list"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"testing"
)

func TestStringConcatFactory_MakeConcatWithConstants_IncorrectParamCount(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []any{
		object.StringObjectFromGoString("caller/MyClass"),
		object.StringObjectFromGoString("makeConcatWithConstants"),
	}

	result := stringConcatFactory_makeConcatWithConstants(params)
	if result == nil {
		t.Errorf("Expected error block for incorrect parameter count, got nil")
	}
}

func TestConcat1param_InvalidParams(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Case 1: Incorrect param count
	res1 := ConcatStringElements([]any{list.New()})
	if res1 == nil {
		t.Errorf("Expected error block for incorrect parameter count, got nil")
	}

	// Case 2: params[1] is not a list
	res2 := ConcatStringElements([]any{"arg0", "not-a-list"})
	if _, ok := res2.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for invalid frame stack, got %T", res2)
	}

	// Case 3: empty frame stack
	emptyFs := list.New()
	res3 := ConcatStringElements([]any{"arg0", emptyFs})
	if _, ok := res3.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for empty frame stack, got %T", res3)
	}
}

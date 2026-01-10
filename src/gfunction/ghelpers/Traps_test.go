/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package ghelpers

import (
	"reflect"
	"strings"
	"testing"

	"jacobin/src/excNames"
)

func TestLoad_Traps_RegistersSomeMethods(t *testing.T) {
	saved := MethodSignatures
	defer func() { MethodSignatures = saved }()
	MethodSignatures = make(map[string]GMeth)

	Load_Traps()
	Load_Traps_Java_Io()

	// Representative subset across class, function, deprecated
	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/io/BufferedOutputStream.<clinit>()V", 0, TrapClass},
		{"java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;", 0, TrapFunction},
		{"java/io/FileDescriptor.valid()Z", 0, TrapFunction},
		//{"java/rmi/RMISecurityManager.<clinit>()V", 0, TrapDeprecated},
		//{"java/rmi/RMISecurityManager.<init>()V", 0, TrapDeprecated},
	}

	for _, c := range checks {
		got, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if got.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
		}
		if got.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		p1 := reflect.ValueOf(got.GFunction).Pointer()
		p2 := reflect.ValueOf(c.fn).Pointer()
		if p1 != p2 {
			t.Logf("%s GFunction mismatch, p1: %v, p2: %v", c.key, p1, p2)
		}
	}
}

func TestTrapFunctions_ReturnUnsupported(t *testing.T) {
	// TrapClass
	if blk, ok := TrapClass(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
		t.Fatalf("TrapClass expected UnsupportedOperationException with TRAP: message, got %+v", blk)
	}
	// TrapFunction
	if blk, ok := TrapFunction(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
		t.Fatalf("TrapFunction expected UnsupportedOperationException with TRAP: message, got %+v", blk)
	}
	// TrapDeprecated
	if blk, ok := TrapDeprecated(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
		t.Fatalf("TrapDeprecated expected UnsupportedOperationException with TRAP: message, got %+v", blk)
	}
	// TrapUndocumented
	if blk, ok := TrapUndocumented(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
		t.Fatalf("TrapUndocumented expected UnsupportedOperationException with TRAP: message, got %+v", blk)
	}
	// TrapProtected
	if blk, ok := TrapProtected(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
		t.Fatalf("TrapProtected expected UnsupportedOperationException with TRAP: message, got %+v", blk)
	}
}

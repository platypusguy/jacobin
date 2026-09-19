/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package ghelpers

import (
	"strings"
	"testing"

	"jacobin/src/excNames"
)

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

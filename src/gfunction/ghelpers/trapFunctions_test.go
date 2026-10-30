/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package ghelpers

import (
	"jacobin/src/excNames"
	"strings"
	"testing"
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

// Trap Preview
func TestTrapPreview(t *testing.T) {
	if blk, ok := TrapPreview(nil).(*GErrBlk); !ok ||
		blk.ExceptionType != excNames.UnsupportedClassVersionError || !strings.Contains(blk.ErrMsg, "Preview features") {
		t.Fatalf("TrapPreview expected UnsupportedClassVersionError message containing 'Preview features', got %+v", blk)
	}
}

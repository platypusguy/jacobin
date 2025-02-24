/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"io"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/globals"
	"os"
	"strings"
	"testing"
)

// test the interpreter methods *other than those that execute bytecodes*
// bytecodes are tested in interpreter_xxx_test.go

func TestEmptyCodeSegment(t *testing.T) {
	fr := frames.CreateFrame(1)
	fr.Meth = []byte{} // empty code segment

	globals.InitGlobals("test")
	globalPtr := globals.GetGlobalRef()
	globalPtr.FuncMinimalAbort = exceptions.MinimalAbort
	globalPtr.FuncThrowException = exceptions.ThrowExNil

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fs := frames.CreateFrameStack()
	frames.PushFrame(fs, fr)
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Empty code segment") {
		t.Errorf("TestEmptyCodeSegment: did not get expected error message, got: %s", errMsg)
	}
}

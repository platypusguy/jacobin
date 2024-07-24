/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"io"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
	"testing"
)

func TestThrowExNil(t *testing.T) {
	globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)

	// to inspect log messages, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	ThrowExNil(excNames.UnknownError, "just a test")
	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])
	if !strings.Contains(msg, "java.lang.UnknownError") || !strings.Contains(msg, "just a test") {
		t.Errorf("Got unexpected output: %s", msg)
	}
}

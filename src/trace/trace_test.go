/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package trace

import (
	"io"
	"jacobin/globals"
	"os"
	"strings"
	"testing"
)

func initialize() {
	globals.InitGlobals("test")
	Init()
}

func TestErrorMessage(t *testing.T) {

	initialize()

	// Save existing stderr
	savedStderr := os.Stderr

	// Capture the writing done to stderr in a pipe
	rdr, wrtr, _ := os.Pipe()
	os.Stderr = wrtr
	Error("Only YOU can prevent forest fires!")
	_ = wrtr.Close()

	// Restore stderr to what it was before
	os.Stderr = savedStderr

	// Collect stderr output bytes --> string
	outBytes, _ := io.ReadAll(rdr)
	outString := string(outBytes[:])

	// What we expected?
	if !strings.Contains(outString, "ERROR") { // No
		t.Errorf("Empty trace message failed: expected an error message but saw [%s]\n", outString)
	}

}

func TestWarningMessage(t *testing.T) {

	initialize()

	// Save existing stderr
	savedStderr := os.Stderr

	// Capture the writing done to stderr in a pipe
	rdr, wrtr, _ := os.Pipe()
	os.Stderr = wrtr
	Warning("Woe is me!")
	_ = wrtr.Close()

	// Restore stderr to what it was before
	os.Stderr = savedStderr

	// Collect stderr output bytes --> string
	outBytes, _ := io.ReadAll(rdr)
	outString := string(outBytes[:])

	// What we expected?
	if !strings.Contains(outString, "WARNING") { // No
		t.Errorf("Empty trace message failed: expected a warning message but saw [%s]\n", outString)
	}

}

func TestValidTraceMessage(t *testing.T) {

	initialize()

	const expected = "Mary had a little lamb whose fleece was white as snow"

	// Save existing stderr
	savedStderr := os.Stderr

	// Capture the writing done to stderr in a pipe
	rdr, wrtr, _ := os.Pipe()
	os.Stderr = wrtr
	Trace(expected)
	_ = wrtr.Close()

	// Restore stderr to what it was before
	os.Stderr = savedStderr

	// Collect stderr output bytes --> string
	outBytes, _ := io.ReadAll(rdr)
	outString := string(outBytes[:])

	// What we expected?
	if !strings.Contains(outString, expected) { // No
		t.Errorf("Nonempty trace message failed: expected [%s] as a subset of [%s]\n", expected, outString)
	}

}

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

func TestEmptyTraceMessage(t *testing.T) {

	initialize()

	// Save existing stderr
	savedStderr := os.Stderr

	// Capture the writing done to stderr in a pipe
	rdr, wrtr, _ := os.Pipe()
	os.Stderr = wrtr
	Trace("")
	_ = wrtr.Close()

	// Restore stderr to what it was before
	os.Stderr = savedStderr

	// Collect stderr output bytes --> string
	outBytes, _ := io.ReadAll(rdr)
	outString := string(outBytes[:])

	// What we expected?
	if !strings.Contains(outString, EmptyMsg) { // No
		t.Errorf("Empty trace message failed: expected [%s] as a subset of [%s]\n", EmptyMsg, outString)
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
		t.Errorf("Nonempty trace message failed: expected [%s] as a subset of [%s]\n", EmptyMsg, outString)
	}

}

func TestFailoverToStdout(t *testing.T) {

	initialize()

	const expected = "Mary had a little lamb whose fleece was white as snow"

	// Save existing stderr
	savedStderr := os.Stderr
	savedStdout := os.Stdout

	// Set up stderr from a pipe and then close it
	_, wrtrErr, _ := os.Pipe()
	os.Stderr = wrtrErr
	err := os.Stderr.Close()
	if err != nil {
		os.Stderr = savedStderr
		os.Stdout = savedStdout
		t.Errorf("Failed to close os.Stderr, err: %v\n", err)
	}

	// Capture the writing done to stdout in a pipe
	rdrOut, wrtrOut, _ := os.Pipe()
	os.Stdout = wrtrOut

	Trace(expected)
	_ = wrtrOut.Close()

	// Restore stderr to what it was before
	os.Stderr = savedStderr
	os.Stdout = savedStdout

	// Collect stderr output bytes --> string
	outBytes, _ := io.ReadAll(rdrOut)
	outString := string(outBytes[:])

	// What we expected?
	if !strings.Contains(outString, expected) { // No
		t.Errorf("Stdout trace message failed: expected [%s] as a subset of [%s]\n", EmptyMsg, outString)
	}

}

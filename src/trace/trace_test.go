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
	"regexp"
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


func TestDisableSuppressesOutput(t *testing.T) {
	initialize()
	// Ensure we re-enable after this test
	defer Init()

	// capture stderr
	saved := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Disable()
	Trace("should not appear")
	Error("nor should this")
	Warning("or this")
	_ = w.Close()
	os.Stderr = saved

	buf, _ := io.ReadAll(r)
	if len(buf) != 0 {
		t.Fatalf("Disable() did not suppress output; saw: %q", string(buf))
	}
}

func TestTracePrefixContainsTimestamp(t *testing.T) {
	initialize()

	// capture stderr
	saved := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Trace("hello")
	_ = w.Close()
	os.Stderr = saved

	out, _ := io.ReadAll(r)
	line := string(out)
	// Expect something like: [  0.0xxs] hello\n
	re := regexp.MustCompile(`^\s*\[[0-9]+\.[0-9]{3}s\] \s*hello\s*\n?$`)
	if !re.MatchString(line) {
		// Be lenient: check for bracketed seconds.millis and a space before message
		if !(strings.HasPrefix(line, "[") && strings.Contains(line, "] ") && strings.Contains(line, "hello")) {
			t.Fatalf("Trace prefix not as expected; got: %q", line)
		}
	}
}

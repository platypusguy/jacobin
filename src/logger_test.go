/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLogLevelTooLow(t *testing.T) {
	Global = initGlobals("test")
	err := SetLogLevel(0)
	if err == nil {
		t.Error("setting logging level to 0 did not generate an error")
	}
}

func TestLogLevelTooHigh(t *testing.T) {
	Global = initGlobals("test")
	err := SetLogLevel(99)
	if err == nil {
		t.Error("setting logging level to 99 did not generate an error")
	}
}

func TestEmptyLogMessage(t *testing.T) {
	Global = initGlobals("test")
	SetLogLevel(WARNING)
	err := Log("", SEVERE)
	if err == nil {
		t.Error("trying to log an empty message did not generate an error")
	}
}

func TestValidLogMessageFineLevel(t *testing.T) {
	Global = initGlobals("test")
	SetLogLevel(FINE)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Log("Test message (FINE)", FINE)

	// reset stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Test message (FINE)") ||
		!strings.HasPrefix(msg, "[") { // a FINE message should start with elapsed time between [ ]'s
		t.Error("valid FINE logging message was not logged properly")
	}
}

func TestValidLogMessageWarningLevel(t *testing.T) {
	Global = initGlobals("test")
	SetLogLevel(WARNING)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Log("Test message (WARNING)", WARNING)

	// reset stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Test message (WARNING)") ||
		strings.HasPrefix(msg, "[") { // if the global log level is warning, no elapsed time should be logged
		t.Error("valid WARNING logging message was not logged properly")
	}
}

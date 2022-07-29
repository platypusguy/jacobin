/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package log

import (
	"io/ioutil"
	"jacobin/globals"
	"os"
	"strings"
	"testing"
)

func TestLoggerInit(t *testing.T) {
	globals.InitGlobals("test")
	Init()

	if Level != WARNING {
		t.Error("log init() did not set Level to WARNING")
	}
}

func TestGlobalSetLogLevelTooLow(t *testing.T) {
	globals.InitGlobals("test")
	err := SetLogLevel(0)
	if err == nil {
		t.Error("setting logging level to 0 did not generate an error")
	}
}

func TestGlobalSetLogLevelTooHigh(t *testing.T) {
	globals.InitGlobals("test")
	err := SetLogLevel(99)
	if err == nil {
		t.Error("setting logging level to 99 did not generate an error")
	}
}

// you cannot set logging level to log.SEVERE (which would hide log.WARNINGs), so
// attempting to do so should generate an error
func TestLogSetLogLevelTologSevere(t *testing.T) {
	globals.InitGlobals("test")
	err := SetLogLevel(SEVERE)
	if err == nil {
		t.Error("setting logging level to log.SEVERE did not generate an error")
	}
}

func TestSettingLogLevels(t *testing.T) {
	globals.InitGlobals("test") // this sets the Level to log.WARNING (the default value)
	err := SetLogLevel(CLASS)
	if err != nil || (Level != CLASS) {
		t.Error("setting logging level to CLASS did not work correctly")
	}
	err = SetLogLevel(FINE)
	if err != nil || (Level != FINE) {
		t.Error("setting logging level to FINE did not work correctly")
	}

	err = SetLogLevel(FINEST)
	if err != nil || (Level != FINEST) {
		t.Error("setting logging level to FINEST did not work correctly")
	}
}

func TestEmptyLogMessage(t *testing.T) {
	globals.InitGlobals("test")
	_ = SetLogLevel(WARNING)
	err := Log("", SEVERE)
	if err == nil {
		t.Error("trying to log an empty message did not generate an error")
	}
}

func TestValidLogMessageFineLevel(t *testing.T) {
	globals.InitGlobals("test")
	_ = SetLogLevel(FINE)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	_ = Log("Test message (FINE)", FINE)

	// reset stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Test message (FINE)") ||
		!strings.HasPrefix(msg, "[") { // a FINE message should start with elapsed time between [ ]'s
		t.Error("valid FINE logging message was not logged properly")
	}
}

func TestValidLogMessagelogWarningLevel(t *testing.T) {
	globals.InitGlobals("test")
	_ = SetLogLevel(WARNING)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	_ = Log("Test message (log.WARNING)", WARNING)

	// reset stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Test message (log.WARNING)") ||
		strings.HasPrefix(msg, "[") { // if the global log level is log.WARNING, no elapsed time should be logged
		t.Error("valid log.WARNING logging message was not logged properly")
	}
}

func TestLoggingMessageAtInvalidLoggingLevel(t *testing.T) {
	globals.InitGlobals("test")
	_ = SetLogLevel(WARNING)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := Log("Test message (log.WARNING)", 0)

	// reset stderr to what it was before
	_ = w.Close()
	os.Stdout = normalStderr

	if err == nil {
		t.Error("logging message at invalid logging level did not generate ane error")
	}
}

func TestThatMsgWithFinerLoggingLevelThanAllowedPrintsNothing(t *testing.T) {
	globals.InitGlobals("test")
	_ = SetLogLevel(WARNING)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	_ = Log("Test message (log.WARNING)", FINEST)

	// reset stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if len(msg) > 0 {
		t.Errorf("Test should not have logged anything, but it did: %s", msg)
	}
}

func TestThatTraceLoggingWithoutCLIsettingPrintsNothing(t *testing.T) {
	globals.InitGlobals("test")
	_ = SetLogLevel(TRACE_INST)

	// to test the error message, capture the writing done to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	_ = Log("Test message (log.WARNING)", TRACE_INST)

	// reset stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if len(msg) > 0 {
		t.Errorf("Test should not have logged anything, but it did: %s", msg)
	}
}

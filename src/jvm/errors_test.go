/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"container/list"
	"encoding/binary"
	"errors"
	"io"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/thread"
	"os"
	"runtime/debug"
	"strings"
	"testing"
)

func TestFormatOverflowError(t *testing.T) {
	fs := frames.CreateFrameStack() // create a frame stack
	fr := frames.CreateFrame(3)     // create a new frame
	fr.MethName = "main"
	fr.ClName = "error_test"
	fr.CP = nil
	fr.PC = 4 // set hypothetical program counter to 4
	fr.Meth = make([]byte, 6)
	_ = frames.PushFrame(fs, fr) // push the minimal frame on to the frame stack

	formatStackOverflowError(fr)

	if fr.PC != 0 {
		t.Errorf("Expecting PC to be 0, got: %d", fr.PC)
	}

	if fr.Meth[1] != IMPDEP2 {
		t.Errorf("Expecting bytecode to be IMDEP2 (%X), got: %X", IMPDEP2, fr.Meth[1])
	}

	if fr.Meth[2] != 0x01 {
		t.Errorf("Expecting error code to be 0x01, got: %X", fr.Meth[2])
	}

	bytes := make([]byte, 2)
	bytes[0] = fr.Meth[3]
	bytes[1] = fr.Meth[4]
	location := int16(binary.BigEndian.Uint16(bytes))
	if location != 4 {
		t.Errorf("Expecting saved PC to be 4, got %d", location)
	}
}

// if the JVM frame stack has already been displayed, then
// don't display it again.
func TestShowFrameStackWhenPreviouslyShown(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	th := thread.ExecThread{}
	globals.GetGlobalRef().JvmFrameStackShown = true // should prevent any output
	showFrameStack(&th)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	if string(msg) != "" {
		t.Errorf("Got following output when expecting none: %s", string(msg))
	}
}

// if the JVM stack is empty, then notify the user that
// no additional data is available
func TestShowFrameStackWithEmptyStack(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	th := thread.CreateThread()
	th.Stack = list.New()
	globals.GetGlobalRef().JvmFrameStackShown = false
	showFrameStack(&th)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if errMsg != "no further data available\n" {
		t.Errorf("Got this when expecting 'no further data available': %s", errMsg)
	}
}

func TestShowFrameStackWithOneEntry(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	f := frames.CreateFrame(1) // create a new frame
	f.MethName = "main"
	f.ClName = "testClass"
	f.PC = 42

	th := thread.CreateThread()
	th.Stack = frames.CreateFrameStack()
	_ = frames.PushFrame(th.Stack, f)

	globals.GetGlobalRef().JvmFrameStackShown = false
	showFrameStack(&th)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if errMsg != "Method: testClass.main                           PC: 042\n" {
		t.Errorf("Got this when expecting 'Method: testClass.main                           PC: 042': %s",
			errMsg)
	}
}

// check that when a Go stack if it has not been previously been captured
func TestShowGoStackWhenNotPreviouslyCaptured(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().GoStackShown = false

	showGoStackTrace(nil)
	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	contents := string(msg)
	if !strings.Contains(contents, "goroutine") {
		t.Errorf("Go stack did not contain expected entries: %s", contents)
	}

	if globals.GetGlobalRef().GoStackShown != true {
		t.Errorf("after showing golang stack, globals.GoStackShown was still false")
	}
}

// check that when a Go stack is captured, it is shown when we call showGoStackTrace()
func TestShowGoStackWhenPreviouslyCaptured(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().GoStackShown = false
	capturedGoStack := debug.Stack()
	stackAsString := string(capturedGoStack)
	globals.GetGlobalRef().ErrorGoStack = stackAsString
	entries := strings.Split(stackAsString, "\n")
	firstEntry := entries[0]

	showGoStackTrace(nil)
	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	contents := string(msg)
	if !strings.Contains(contents, firstEntry) {
		t.Errorf("Go stack did not contain expected entry: %s", contents)
	}
}

// check that when a Go stack is not shown a second time when we call showGoStackTrace()
func TestShowGoStackWhenPreviouslyShown(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().GoStackShown = true
	capturedGoStack := debug.Stack()
	stackAsString := string(capturedGoStack)
	globals.GetGlobalRef().ErrorGoStack = stackAsString

	showGoStackTrace(nil)
	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	contents := string(msg)
	if len(contents) != 0 {
		t.Errorf("Expected empty string, got: %s", contents)
	}
}

// showPanicCause() should correctly report an error's content
func TestShowPanicCause(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().PanicCauseShown = false
	cause := errors.New("error causing panic")
	showPanicCause(cause)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if !strings.Contains(errMsg, "error causing panic") {
		t.Errorf("Got unexpected message re panic cause: %s", errMsg)
	}
}

// showPanicCause() should show nothing if it's already been called
func TestShowPanicCauseAfterAlreadyShown(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().PanicCauseShown = true // should prevent showing
	cause := errors.New("error causing panic")
	showPanicCause(cause)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("Expected empty string, got: %s", errMsg)
	}
}

// if showPanicCause() is passed nil, it should state causal data is not available
func TestShowPanicCauseNil(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().PanicCauseShown = false
	showPanicCause(nil)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if !strings.Contains(errMsg, "error: go panic -- cause unknown") {
		t.Errorf("Got unexpected message for nil panic cause: %s", errMsg)
	}
}

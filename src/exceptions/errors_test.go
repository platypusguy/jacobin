/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package exceptions

import (
	"container/list"
	"errors"
	"io"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/thread"
	"jacobin/src/trace"
	"os"
	"runtime/debug"
	"strings"
	"testing"
)

// if the JVM frame stack has already been displayed, then
// don't display it again.
func TestShowFrameStackWhenPreviouslyShown(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	th := thread.ExecThread{}
	globals.GetGlobalRef().JvmFrameStackShown = true // should prevent any output
	ShowFrameStack(&th)

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
	g := globals.InitGlobals("test")
	g.StrictJDK = false

	trace.Init()

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
	ShowFrameStack(&th)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if !strings.Contains(errMsg, "no further data available") {
		t.Errorf("Got this when expecting 'no further data available': %s", errMsg)
	}
}

func TestShowFrameStackWithOneEntry(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	trace.Init()

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
	ShowFrameStack(&th)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Method: testClass.main") {
		t.Errorf("Got this when expecting 'Method: testClass.main                           PC: 042': %s",
			errMsg)
	}
}

// check that when a Go stack is not shown a second time when we call showGoStackTrace()
func TestShowGoStackWhenPreviouslyShown(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test"
	g.StrictJDK = false

	trace.Init()

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

	ShowGoStackTrace(nil)
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

	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().PanicCauseShown = false
	cause := errors.New("error causing panic")
	ShowPanicCause(cause)

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

	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().PanicCauseShown = true // should prevent showing
	cause := errors.New("error causing panic")
	ShowPanicCause(cause)

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

	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.GetGlobalRef().PanicCauseShown = false
	ShowPanicCause(nil)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
	msg, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout

	errMsg := string(msg)
	if !strings.Contains(errMsg, "go panic -- cause unknown") {
		t.Errorf("Got unexpected message for nil panic cause: %s", errMsg)
	}
}

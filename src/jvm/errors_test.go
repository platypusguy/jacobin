/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"container/list"
	"io"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/thread"
	"os"
	"testing"
)

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
	globals.GetGlobalRef().JvmFrameStackShown = false // should prevent any output
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

	thread := thread.ExecThread{}
	globals.GetGlobalRef().JvmFrameStackShown = true // should prevent any output
	showFrameStack(&thread)

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

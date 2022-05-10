/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package shutdown

import (
	"io/ioutil"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
	"testing"
)

func TestShutdownOK(t *testing.T) {
	globals.InitGlobals("test")
	gl := globals.GetGlobalRef()
	gl.JacobinName = "test"

	_ = log.SetLogLevel(log.FINE)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	Exit(UNKNOWN_ERROR)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "shutdown") {
		t.Errorf("Expecting shutdown message, but got: %s", msg)
	}
}

func TestShutdownReturn(t *testing.T) {
	globals.InitGlobals("test")
	gl := globals.GetGlobalRef()
	gl.JacobinName = "test"

	_ = log.SetLogLevel(log.FINE)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	ret := Exit(OK) // should return from Exit with a 0

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if ret != 0 {
		t.Errorf("Expecting exit() return value of 0, but got %d", ret)
	}
}

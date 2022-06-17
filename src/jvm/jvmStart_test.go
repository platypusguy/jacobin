/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io/ioutil"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
	"testing"
)

func TestNoExecutable(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.FINE)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	JVMrun()

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "No executable program specified.") {
		t.Errorf("jvmRun() with no executable specified did not get expected error msg, got: %s", errMsg)
	}
}

func TestNoMainClassInJar(t *testing.T) {
	
}

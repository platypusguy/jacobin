/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"bytes"
	"io"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoExecutable(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errC := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		errC <- buf.String()
	}()

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	JVMrun()

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	errMsg := <-errC

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "No executable program specified.") {
		t.Errorf("jvmRun() with no executable specified did not get expected error msg, got: %s", errMsg)
	}
}

func TestNoMainClassInJar(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Error("Error getting current working directory")
		return
	}

	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.JacobinName = "test" // prevents a shutdown when the exception hits.
	g.StartingJar = filepath.Join(cwd, "..", "..", "testdata", "nomanifest.jar")
	g.StrictJDK = false

	log.Init()
	_ = log.SetLogLevel(log.INFO)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errC := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		errC <- buf.String()
	}()

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	JVMrun()

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	errMsg := <-errC

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(errMsg, "no main manifest attribute") {
		t.Errorf("jvmRun() with a jar that has no manifest should have given no main manifest attribute error, got %s", errMsg)
	}
}

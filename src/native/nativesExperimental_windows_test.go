//go:build windows

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"jacobin/globals"
	"jacobin/log"
	"os"
	"testing"
)

// simple test to exercise the code, rather than validating.
// validation tests will come in time.
func TestExports(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stdout to avoid printing error message to console
	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	normalStderr := os.Stderr
	_, werr, _ := os.Pipe()
	os.Stderr = werr

	_ = log.SetLogLevel(log.FINE)
	err := CreateNativeFunctionTable("")
	if err != nil {
		t.Error(err)
	}

	_ = w.Close()
	// msg, _ := io.ReadAll(r)
	os.Stdout = normalStdout // restore stdout

	_ = werr.Close()
	// msgErr, _ := io.ReadAll(rerr)
	os.Stderr = normalStderr

	// if string(msg) == "" {
	// 	t.Error("expected list of DLL files, but got none")
	// }
}

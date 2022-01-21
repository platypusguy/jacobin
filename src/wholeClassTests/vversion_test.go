/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package wholeClassTests

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

/*
 * Tests for: jacobin --version
 *     which should print info about the JVM to stdout (which go test reroutes to stderr (dang it!))
 */

func initVarsVversion() {
	_JACOBIN = "d:\\GoogleDrive\\Dev\\jacobin\\src\\jacobin.exe"
	_JVM_ARGS = "--version"
	_TESTCLASS = "" // the class to test
	_APP_ARGS = ""
}

func TestRunVversion(t *testing.T) {
	initVarsVersion()
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	// test that executable exists
	if _, err := os.Stat(_JACOBIN); err != nil {
		t.Errorf("Missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	// run the various combinations of args. This is necessary b/c the empty string is viewed as
	// an actual specified option on the command line.
	if len(_JVM_ARGS) > 0 {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS, _APP_ARGS)
		} else {
			if len(_TESTCLASS) > 0 {
				cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS)
			} else {
				cmd = exec.Command(_JACOBIN, _JVM_ARGS)
			}
		}
	} else {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _TESTCLASS)
		}
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// run the command
	if err = cmd.Start(); err != nil {
		t.Errorf("Got error running Jacobin: %s", err.Error())
	}

	slurpErr, _ := io.ReadAll(stderr) // the output is written to stderr, which go test reroutes to stderr
	if len(slurpErr) == 0 {
		t.Errorf("Expected output to stderr, but got none")
	}

	msg := string(slurpErr)
	if !strings.HasPrefix(msg, "Jacobin VM") {
		t.Errorf("Output did not begin with Jacobin name, instead: %s", string(msg))
	}
	//
	if !strings.Contains(msg, "64-bit server VM") {
		t.Errorf("Did not get expected output to stderr. Got: %s", msg)
	}
}

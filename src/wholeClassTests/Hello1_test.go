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

var JACOBIN = "d:\\GoogleDrive\\Dev\\jacobin\\src\\jacobin.exe"
var JVM_ARGS = ""
var TESTCLASS = "d:\\GoogleDrive\\Dev\\jacobin\\testdata\\Hello.class" // the class to test
var APP_ARGS = ""

func TestRunHello1(t *testing.T) {
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	// test that executable exists
	if _, err := os.Stat(JACOBIN); err != nil {
		t.Errorf("Missing Jacobin executable, which was specified as %s", JACOBIN)
	}

	// run the various combinations of args. This is necessary b/c the empty string is viewed as
	// an actual specified option on the command line.
	if len(JVM_ARGS) > 0 {
		if len(APP_ARGS) > 0 {
			cmd = exec.Command(JACOBIN, JVM_ARGS, TESTCLASS, APP_ARGS)
		} else {
			cmd = exec.Command(JACOBIN, JVM_ARGS, TESTCLASS)
		}
	} else {
		if len(APP_ARGS) > 0 {
			cmd = exec.Command(JACOBIN, TESTCLASS, APP_ARGS)
		} else {
			cmd = exec.Command(JACOBIN, TESTCLASS)
		}
	}

	// get the stdout and stderr contents from the file execution
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// run the command
	if err = cmd.Start(); err != nil {
		t.Errorf("Got error running Jacobin: %s", err.Error())
	}

	// Here begin the actual tests on the output to stderr and stdout
	slurp, _ := io.ReadAll(stderr)
	if len(slurp) != 0 {
		t.Errorf("Got unexpected output to stderr: %s", string(slurp))
	}

	slurp, _ = io.ReadAll(stdout)
	if !strings.HasPrefix(string(slurp), "Jacobin VM") {
		t.Errorf("Stdout did not begin with Jacobin copyright, instead: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "Hello from Hello.main!") {
		t.Errorf("Did not get expected output to stdout. Got: %s", string(slurp))
	}
}

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
 * Tests for: jacobin -client -version
 *     which should print info about the client version of the JVM to stderr
 * These tests check the output.
 */

func initVarsClientVersion() {
	_JACOBIN = "d:\\GoogleDrive\\Dev\\jacobin\\src\\jacobin.exe"
}

func TestRunClientVersion(t *testing.T) {
	initVarsClientVersion()
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	// test that executable exists
	if _, err := os.Stat(_JACOBIN); err != nil {
		t.Errorf("Missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	cmd = exec.Command(_JACOBIN, "-client", "-version")

	// get the stdout and stderr contents from the file execution
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// run the command
	if err = cmd.Start(); err != nil {
		t.Errorf("Got error running Jacobin: %s", err.Error())
	}

	// Here begin the actual tests on the output to stderr and stdout

	slurp, _ := io.ReadAll(stderr)
	if !strings.HasPrefix(string(slurp), "Jacobin VM") {
		t.Errorf("Stderr did not begin with Jacobin name, instead: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "64-bit client VM") {
		t.Errorf("Did not get expected output to stderr. Got: %s", string(slurp))
	}
}

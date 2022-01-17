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

// Test for Hello3 class, which has loop and function calls 3 levels deep. Source code:
//
// class Hello3 {
//
// 	public static void main( String[] args) {
// 		int x;
// 		for( int i = 0; i < 10; i++) {
// 			x = addTwo(i, i-1);
// 			System.out.println( x );
// 		}
// 	}
//
// 	static int addTwo(int j, int k) {
// 		int m = multTwo(j, k);
// 		return m+1;
// 	}
//
// 	static int multTwo(int m, int n){
// 		return m*n;
// }
// }
//
// These tests check the output with various options for verbosity and features set on the command line.

func initVarsHello3() {
	_JACOBIN = "d:\\GoogleDrive\\Dev\\jacobin\\src\\jacobin.exe"
	_JVM_ARGS = ""
	_TESTCLASS = "d:\\GoogleDrive\\Dev\\jacobin\\testdata\\Hello3.class" // the class to test
	_APP_ARGS = ""
}

func TestRunHello3(t *testing.T) {
	initVarsHello3()
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
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS)
		}
	} else {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _TESTCLASS)
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
	slurpErr := string(slurp)
	if len(slurp) != 0 {
		t.Errorf("Got unexpected output to stderr: %s", slurpErr)
	}

	slurp, _ = io.ReadAll(stdout)
	slurpOut := string(slurp)
	if !strings.HasPrefix(slurpOut, "Jacobin VM") {
		t.Errorf("Stdout did not begin with Jacobin copyright, instead: %s", slurpOut)
	}

	if !strings.Contains(slurpOut, "57") && !strings.Contains(string(slurp), "73") {
		t.Errorf("Did not get expected output to stdout. Got: %s", slurpOut)
	}
}

func TestRunHello3VerboseClass(t *testing.T) {
	initVarsHello3()
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	// test that executable exists
	if _, err := os.Stat(_JACOBIN); err != nil {
		t.Errorf("Missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	_JVM_ARGS = "-verbose:class"
	// run the various combinations of args. This is necessary b/c the empty string is viewed as
	// an actual specified option on the command line.
	if len(_JVM_ARGS) > 0 {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS)
		}
	} else {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _TESTCLASS)
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
	slurpErr := string(slurp)
	if !strings.Contains(slurpErr, "Class: Hello3, loader: bootstrap") {
		t.Errorf("Got unexpected output to stderr: %s", slurpErr)
	}

	slurp, _ = io.ReadAll(stdout)
	slurpOut := string(slurp)
	if !strings.HasPrefix(string(slurp), "Jacobin VM") {
		t.Errorf("Stdout did not begin with Jacobin copyright, instead: %s", slurpOut)
	}

	if !strings.Contains(slurpOut, "31") && !strings.Contains(slurpOut, "43") {
		t.Errorf("Did not get expected output to stdout. Got: %s", slurpOut)
	}
}

func TestRunHello3VerboseFinest(t *testing.T) {
	initVarsHello3()
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	// test that executable exists
	if _, err := os.Stat(_JACOBIN); err != nil {
		t.Errorf("Missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	_JVM_ARGS = "-verbose:finest"
	// run the various combinations of args. This is necessary b/c the empty string is viewed as
	// an actual specified option on the command line.
	if len(_JVM_ARGS) > 0 {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS)
		}
	} else {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _TESTCLASS)
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
	slurpErr := string(slurp)
	if !strings.Contains(slurpErr, "Class Hello3 has been format-checked.") {
		t.Errorf("Got unexpected output to stderr: %s", slurpErr)
	}

	slurp, _ = io.ReadAll(stdout)
	slurpOut := string(slurp)
	if !strings.HasPrefix(slurpOut, "Jacobin VM") {
		t.Errorf("Stdout did not begin with Jacobin copyright, instead: %s", slurpOut)
	}

	if !strings.Contains(slurpOut, "21") {
		t.Errorf("Did not get expected output to stdout. Got: %s", slurpOut)
	}
}

func TestRunHello3TraceInst(t *testing.T) {
	initVarsHello3()
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	// test that executable exists
	if _, err := os.Stat(_JACOBIN); err != nil {
		t.Errorf("Missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	_JVM_ARGS = "-trace:inst"
	// run the various combinations of args. This is necessary b/c the empty string is viewed as
	// an actual specified option on the command line.
	if len(_JVM_ARGS) > 0 {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _JVM_ARGS, _TESTCLASS)
		}
	} else {
		if len(_APP_ARGS) > 0 {
			cmd = exec.Command(_JACOBIN, _TESTCLASS, _APP_ARGS)
		} else {
			cmd = exec.Command(_JACOBIN, _TESTCLASS)
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
	slurpErr := string(slurp)
	if !strings.Contains(slurpErr, "class: Hello3, meth: main, pc: 29, inst: RETURN, tos: -1") {
		t.Errorf("Got unexpected output to stderr: %s", slurpErr)
	}

	slurp, _ = io.ReadAll(stdout)
	slurpOut := string(slurp)
	if !strings.HasPrefix(slurpOut, "Jacobin VM") {
		t.Errorf("Stdout did not begin with Jacobin copyright, instead: %s", slurpOut)
	}

	if !strings.Contains(slurpOut, "13") {
		t.Errorf("Did not get expected output to stdout. Got: %s", slurpOut)
	}
}

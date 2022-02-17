/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package wholeClassTests

import (
	"fmt"
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
// To run your class, enter its name in _TESTCLASS, any args in their respective variables and then run the tests.
// This test harness expects that environmental variable JACOBIN_EXE gives the full name and path of the executable
// we're running the tests on. The folder which contains the test class should be specified in the environmental
// variable JACOBIN_TESTDATA (without a terminating slash).
func initVarsHello3() error {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		return fmt.Errorf("test not run due to -short")
	}

	_JACOBIN = os.Getenv("JACOBIN_EXE") // returns "" if JACOBIN_EXE has not been specified.
	_JVM_ARGS = ""
	_TESTCLASS = "Hello3.class" // the class to test
	_APP_ARGS = ""

	if _JACOBIN == "" {
		return fmt.Errorf("test failure due to missing Jacobin executable. Please specify it in JACOBIN_EXE")
	} else if _, err := os.Stat(_JACOBIN); err != nil {
		return fmt.Errorf("missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	testClass := os.Getenv("JACOBIN_TESTDATA") + string(os.PathSeparator) + _TESTCLASS
	if _, err := os.Stat(testClass); err != nil {
		return fmt.Errorf("nissing class to test, which was specified as %s", testClass)
	} else {
		_TESTCLASS = testClass
	}
	return nil
}

func TestRunHello3(t *testing.T) {
	initErr := initVarsHello3()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}
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
	initErr := initVarsHello3()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}
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
	initErr := initVarsHello3()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}
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
	initErr := initVarsHello3()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}
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

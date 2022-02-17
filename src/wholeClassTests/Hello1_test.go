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

/*
 * Tests for Hello.class, which is the first class Jacobin executed. Its source code:
 *
 *	class Hello {
 *		public static void main( String[] args) {
 *			for( int i = 0; i < 10; i++)
 *				System.out.println( "Hello from Hello.main!" );
 *		}
 *	}
 *
 * These tests check the output with various options for verbosity and features set on the command line.
 */

// To run your class, enter its name in _TESTCLASS, any args in their respective variables and then run the tests.
// This test harness expects that environmental variable JACOBIN_EXE gives the full name and path of the executable
// we're running the tests on. The folder which contains the test class should be specified in the environmental
// variable JACOBIN_TESTDATA (without a terminating slash).
func initVarsHello() error {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		return fmt.Errorf("test not run due to -short")
	}

	_JACOBIN = os.Getenv("JACOBIN_EXE") // returns "" if JACOBIN_EXE has not been specified.
	_JVM_ARGS = ""
	_TESTCLASS = "Hello.class" // the class to test
	_APP_ARGS = ""

	if _JACOBIN == "" {
		return fmt.Errorf("test failure due to missing Jacobin executable. Please specify it in JACOBIN_EXE")
	} else if _, err := os.Stat(_JACOBIN); err != nil {
		return fmt.Errorf("missing Jacobin executable, which was specified as %s", _JACOBIN)
	}

	if _TESTCLASS != "" {
		testClass := os.Getenv("JACOBIN_TESTDATA") + string(os.PathSeparator) + _TESTCLASS
		if _, err := os.Stat(testClass); err != nil {
			return fmt.Errorf("missing class to test, which was specified as %s", testClass)
		} else {
			_TESTCLASS = testClass
		}
	}
	return nil
}

func TestRunHello(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsHello()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}

	var cmd *exec.Cmd
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

func TestRunHelloVerboseClass(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsHello()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}

	var cmd *exec.Cmd

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
	if !strings.Contains(string(slurp), "Class: Hello, loader: bootstrap") {
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

func TestRunHelloVerboseFinest(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsHello()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}

	var cmd *exec.Cmd

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
	if !strings.Contains(string(slurp), "Class Hello has been format-checked.") {
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

func TestRunHelloTraceInst(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsHello()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}

	_JVM_ARGS = "-trace:inst"

	// run the various combinations of args. This is necessary b/c the empty string is viewed as
	// an actual specified option on the command line.

	var cmd *exec.Cmd
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
	if !strings.Contains(string(slurp), "class: Hello, meth: main, pc: 22, inst: RETURN, tos: -1") {
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

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. All rights reserved.
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
 * Tests for LOOKUPSWITCH processing. Surce code:
 *
 * class lookupswitch {
 * // class to test functioning of lookupswitch bytecode
 *   public static void main(String[] args) {
 *       var len = args.length;
 *       switch (len) {
 *           case 0: System.out.println("zero args"); break;
 *           case -100: System.out.println("100 args"); break;
 *           case 250: System.out.println("250 args"); break;
 *           default: System.out.println("args != -100, 0, or 250"); break;
 *       }
 *   }
 * }
 *
 * This test checks the output with various options on the command line.
 */

// To run your class, enter its name in _TESTCLASS, any args in their respective variables and then run the tests.
// This test harness expects that environmental variable JACOBIN_EXE gives the full name and path of the executable
// we're running the tests on. The folder which contains the test class should be specified in the environmental
// variable JACOBIN_TESTDATA (without a terminating slash).
func initVarsLookupswtch() error {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		return fmt.Errorf("test not run due to -short")
	}

	_JACOBIN = os.Getenv("JACOBIN_EXE") // returns "" if JACOBIN_EXE has not been specified.
	_JVM_ARGS = ""
	_TESTCLASS = "lookupswitch.class" // the class to test
	_APP_ARGS = ""

	if _JACOBIN == "" {
		return fmt.Errorf("missing Jacobin executable. Please specify it in JACOBIN_EXE")
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

func TestLookupSwitchNoArgs(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsLookupswtch()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// run the command
	if err = cmd.Start(); err != nil {
		t.Errorf("Got error running Jacobin: %s", err.Error())
	}

	// Here begin the actual tests on the output to stderr and stdout
	slurp, _ := io.ReadAll(stdout)
	if len(slurp) == 0 {
		t.Errorf("Did not get error output to stdout")
	}

	if !strings.Contains(string(slurp), "zero args") {
		t.Errorf("Did not get expected output to stderr. Got: %s", string(slurp))
	}
}

// same as previous test, but with one argument
func TestLookupSwitchOneArg(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsLookupswtch()
	if initErr != nil {
		t.Fatalf("Test failure due to: %s", initErr.Error())
	}
	var cmd *exec.Cmd

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	_APP_ARGS = "allo!"
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// run the command
	if err = cmd.Start(); err != nil {
		t.Errorf("Got error running Jacobin: %s", err.Error())
	}

	// Here begin the actual tests on the output to stderr and stdout
	slurp, _ := io.ReadAll(stdout)
	if len(slurp) == 0 {
		t.Errorf("Did not get error output to stderr")
	}

	if !strings.Contains(string(slurp), "args != -100, 0, or 250") {
		t.Errorf("Did not get expected output to stderr. Got: %s", string(slurp))
	}
}

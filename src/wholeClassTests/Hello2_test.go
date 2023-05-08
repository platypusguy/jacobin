/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
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
 * Tests for Hello2.class, which is one of the first classes Jacobin executed. Source code:
 *
 *		public static void main( String[] args) {
 *			int x;
 *			for( int i = 0; i < 10; i++) {
 *				x = addTwo(i, i-1);
 *				System.out.println( x );
 *          }
 *      }
 *
 *	    static int addTwo(int j, int k) {
 *		    return j + k;
 *	    }
 *
 * These tests check the output with various options for verbosity and features set on the command line.
 */

// To run your class, enter its name in _TESTCLASS, any args in their respective variables and then run the tests.
// This test harness expects that environmental variable JACOBIN_EXE gives the full name and path of the executable
// we're running the tests on. The folder which contains the test class should be specified in the environmental
// variable JACOBIN_TESTDATA (without a terminating slash).
func initVarsHello2() error {
    if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
        return fmt.Errorf("test not run due to -short")
    }

    _JACOBIN = os.Getenv("JACOBIN_EXE") // returns "" if JACOBIN_EXE has not been specified.
    _JVM_ARGS = ""
    _TESTCLASS = "Hello2.class" // the class to test
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

func TestRunHello2(t *testing.T) {
    if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
        t.Skip()
    }

    initErr := initVarsHello2()
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

    if !strings.Contains(string(slurp), "-1") && !strings.Contains(string(slurp), "17") {
        t.Errorf("Did not get expected output to stdout. Got: %s", string(slurp))
    }
}

func TestRunHello2VerboseClass(t *testing.T) {
    if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
        t.Skip()
    }

    initErr := initVarsHello2()
    if initErr != nil {
        t.Fatalf("Test failure due to: %s", initErr.Error())
    }
    var cmd *exec.Cmd

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
    if !strings.Contains(string(slurp), "Class: Hello2, loader: bootstrap") {
        t.Errorf("Got unexpected output to stderr: %s", string(slurp))
    }

    slurp, _ = io.ReadAll(stdout)
    if !strings.Contains(string(slurp), "-1") && !strings.Contains(string(slurp), "17") {
        t.Errorf("Did not get expected output to stdout. Got: %s", string(slurp))
    }
}

func TestRunHello2VerboseFinest(t *testing.T) {
    if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
        t.Skip()
    }

    initErr := initVarsHello2()
    if initErr != nil {
        t.Fatalf("Test failure due to: %s", initErr.Error())
    }
    var cmd *exec.Cmd

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
    if !strings.Contains(string(slurp), "Class Hello2 has been format-checked.") {
        t.Errorf("Got unexpected output to stderr: %s", string(slurp))
    }

    slurp, _ = io.ReadAll(stdout)
    if !strings.Contains(string(slurp), "13") {
        t.Errorf("Did not get expected output to stdout. Got: %s", string(slurp))
    }
}

func TestRunHello2TraceInst(t *testing.T) {
    if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
        t.Skip()
    }

    initErr := initVarsHello2()
    if initErr != nil {
        t.Fatalf("Test failure due to: %s", initErr.Error())
    }
    var cmd *exec.Cmd

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
    if !strings.Contains(string(slurp), "Hello2     meth: main       PC:   5, ILOAD_2       TOS:  - ") {
        t.Errorf("Got unexpected output to stderr: %s", string(slurp))
    }

    slurp, _ = io.ReadAll(stdout)
    if !strings.Contains(string(slurp), "15") {
        t.Errorf("Did not get expected output to stdout. Got: %s", string(slurp))
    }
}

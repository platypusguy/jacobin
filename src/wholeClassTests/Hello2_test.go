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
 * The bytecode for this is:
 *
 * stack=3, locals=3, args_size=1
 *        0: iconst_0
 *        1: istore_2
 *        2: goto          23
 *        5: iload_2
 *        6: iload_2
 *        7: iconst_1
 *        8: isub
 *        9: invokestatic  #16                 // Method addTwo:(II)I
 *       12: istore_1
 *       13: getstatic     #20                 // Field java/lang/System.out:Ljava/io/PrintStream;
 *       16: iload_1
 *       17: invokevirtual #26                 // Method java/io/PrintStream.println:(I)V
 *       20: iinc          2, 1
 *       23: iload_2
 *       24: bipush        10
 *       26: if_icmplt     5
 *       29: return
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

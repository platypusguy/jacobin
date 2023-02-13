/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
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
 * Tests for testArrays.class, which tests one-dimensional arrays of primitives. Source code:
 *
public static void main( String[] args) {
        int[] intArray = {10, 20, 30, 40};
        int product = intArray[1] * intArray[2];

        System.out.print("intArray: "); System.out.println(product);

        float[] floatArray = { 100.0f, 200.0f, 300.0f};
        float fsum = floatArray[0] + floatArray[2];
        System.out.print("floatArray: "); System.out.println(fsum);

        long[] longArray = { 1000, 2000, 3000};
        long lsum = longArray[0] + longArray[1] + longArray[2] - longArray[0];
        System.out.print( "longArray: "); System.out.println(lsum);

        double[] doubleArray = { 100_000_000.0f, 200.0f, 300.05f};
        double dsum = doubleArray[0] + doubleArray[2];
        System.out.print("doubleArray: "); System.out.println(dsum);

        boolean[] boolArray = {true, false, true, true};
        var trueCount = 0;
        for (boolean tf : boolArray) {
            if( tf == true)
                trueCount += 1;
        }
        System.out.print( "booleanArray: "); System.out.println(trueCount);

        byte[] byteArray = {5,6, 7, 8};
        var byteSum = byteArray[1]+byteArray[3]-byteArray[0];
        System.out.print( "byteArray: "); System.out.println(byteSum);
    }
    These tests check the output with various options for verbosity and features set on the command line.
*/

// To run your class, enter its name in _TESTCLASS, any args in their respective variables and then run the tests.
// This test harness expects that environmental variable JACOBIN_EXE gives the full name and path of the executable
// we're running the tests on. The folder which contains the test class should be specified in the environmental
// variable JACOBIN_TESTDATA (without a terminating slash).
func initVarsPrimitiveArrays() error {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		return fmt.Errorf("test not run due to -short")
	}

	_JACOBIN = os.Getenv("JACOBIN_EXE") // returns "" if JACOBIN_EXE has not been specified.
	_JVM_ARGS = ""
	_TESTCLASS = "testArrays.class" // the class to test
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

func TestRunPrimitiveArrays(t *testing.T) {
	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	initErr := initVarsPrimitiveArrays()
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

	if !strings.Contains(string(slurp), "intArray: 600") {
		t.Errorf("Did not get expected output for intArray. Got: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "floatArray: 400") {
		t.Errorf("Did not get expected output for floatArray. Got: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "longArray: 5000") {
		t.Errorf("Did not get expected output for longArray. Got: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "doubleArray: 1.000003") {
		t.Errorf("Did not get expected output for doubleArray. Got: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "booleanArray: 3") {
		t.Errorf("Did not get expected output for booleanArray. Got: %s", string(slurp))
	}

	if !strings.Contains(string(slurp), "byteArray: 9") {
		t.Errorf("Did not get expected output for byteArray. Got: %s", string(slurp))
	}
}

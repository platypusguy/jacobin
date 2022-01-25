/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"io/ioutil"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
	"testing"
)

// unset all of the JVM environment variables and make sure
// collecting them results in an empty string
func TestGetJVMenvVariablesWhenAbsent(t *testing.T) {
	_ = os.Unsetenv("JAVA_TOOL_OPTIONS")
	_ = os.Unsetenv("_JAVA_OPTIONS")
	_ = os.Unsetenv("JDK_JAVA_OPTIONS")

	javaEnvVars := getEnvArgs()
	if javaEnvVars != "" {
		t.Error("getting non-existent Java environment options failed")
	}
}

// set two of the JVM environment variables and make sure
// they are fetched correctly and a space is inserted between them
func TestGetJVMenvVariablesWhenTwoArePresent(t *testing.T) {
	_ = os.Unsetenv("JAVA_TOOL_OPTIONS")
	_ = os.Setenv("_JAVA_OPTIONS", "Hello,")
	_ = os.Setenv("JDK_JAVA_OPTIONS", "Jacobin!")

	javaEnvVars := getEnvArgs()
	if javaEnvVars != "Hello, Jacobin!" {
		t.Error("getting two set Java environment options failed: " + javaEnvVars)
	}

	// clean up the environment
	_ = os.Unsetenv("_JAVA_OPTIONS")
	_ = os.Unsetenv("JDK_JAVA_OPTIONS")
}

// verify the output to stderr -help option is used
func TestHandleUsageMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	global := globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)
	LoadOptionsTable(global)

	// to avoid cluttering the test results, redirect stdout
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-help"}
	_ = HandleCli(args, &global)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Usage:") ||
		!strings.Contains(msg, "where options include") {
		t.Error("jacobin -help did not generate the usage message to stderr. msg was: " + msg)
	}

	if global.ExitNow != true {
		t.Error("'jacobin -help' should have set Global.exitNow to true to signal end of processing")
	}
}

func TestShowUsageMessageExitsProperlyWith__Help(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	global := globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)
	LoadOptionsTable(global)

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	_, _ = showHelpStdoutAndExit(0, "--help", &global)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if global.ExitNow != true {
		t.Error("'jacobin --help' should set Global.exitNow to true but did not")
	}
}

func TestShowVersionMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	global := globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)

	// to avoid cluttering the test results, redirect stdout
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	LoadOptionsTable(global)
	args := []string{"jacobin", "-showversion", "main.clas"}

	_ = HandleCli(args, &global)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin -version did not generate the correct message to stderr. msg was: " + msg)
	}
}

func TestShow__VersionUsingOptionTable(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	normalStdout := os.Stdout
	r, wout, _ := os.Pipe()
	os.Stdout = wout

	_, _ = versionStdoutThenExit(0, "--version", &global)

	_ = wout.Close()
	os.Stdout = normalStdout
	out, _ := ioutil.ReadAll(r)

	os.Stdout = normalStdout
	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin --version did not generate the correct msg to stdout. msg was: " + msg)
	}

	if global.ExitNow != true {
		t.Error("--version did not set exitNow value to exit. Should be set.")
	}
}

func TestChangeLoggingLevels(t *testing.T) {
	global := globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)
	LoadOptionsTable(global)

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-verbose:info", "main.class"}
	_ = HandleCli(args, &global)

	// reset stdout and stderr to what they were before redirection
	_ = w.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.Level != log.INFO {
		t.Error("Setting log level to INFO via command line failed")
	}

	// --- now test with FINE

	_ = log.SetLogLevel(log.WARNING)

	normalStdout = os.Stdout
	_, wout, _ = os.Pipe()
	os.Stdout = wout

	normalStderr = os.Stderr
	_, w, _ = os.Pipe()
	os.Stderr = w

	LoadOptionsTable(global)
	args = []string{"jacobin", "-verbose:fine", "main.class"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.Level != log.FINE {
		t.Error("Setting log level to FINE via command line failed")
	}

	// --- now try with FINEST

	_ = log.SetLogLevel(log.WARNING)

	normalStdout = os.Stdout
	_, wout, _ = os.Pipe()
	os.Stdout = wout

	normalStderr = os.Stderr
	_, w, _ = os.Pipe()
	os.Stderr = w

	LoadOptionsTable(global)
	args = []string{"jacobin", "-verbose:finest", "main.class"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.Level != log.FINEST {
		t.Error("Setting log level to FINEST via command line failed")
	}

	// --- finally test with CLASS

	_ = log.SetLogLevel(log.WARNING)

	normalStdout = os.Stdout
	_, wout, _ = os.Pipe()
	os.Stdout = wout

	normalStderr = os.Stderr
	_, w, _ = os.Pipe()
	os.Stderr = w

	LoadOptionsTable(global)
	args = []string{"jacobin", "-verbose:class", "main.class"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.Level != log.CLASS {
		t.Error("Setting log level to CLASS via command line failed")
	}
}

func TestInvalidLoggingLevel(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)
	_ = log.SetLogLevel(log.WARNING)

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	_, err := verbosityLevel(0, "severe", &global)

	_ = w.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if err == nil {
		t.Error("Setting log level to SEVERE via command line did not generate expected error")
	}
}

func TestSpecifyClientVM(t *testing.T) {

	global := globals.InitGlobals("test")
	LoadOptionsTable(global)
	if global.VmModel != "server" {
		t.Error("Initialization of Global.vmModel was not set to 'server' Got: " +
			global.VmModel)
	}

	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "-client"}
	_ = HandleCli(args, &global)

	// restore stdout to what it was before
	_ = w.Close()
	os.Stdout = normalStdout

	if global.VmModel != "client" {
		t.Error("global.vmModel should be set to 'client'. Instead got: " +
			global.VmModel)
	}
}

func TestSpecifyValidButUnsupportedOption(t *testing.T) {

	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	// redirect stdout to avoid cluttering test results with copyright notice
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "--dry-run", "main.class"}
	_ = HandleCli(args, &global)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "not currently supported") {
		t.Error("Unsupported but valid option not identified. Instead got" + msg)
	}
}

func TestShowCopyright(t *testing.T) {
	globals.InitGlobals("test")
	_ = log.SetLogLevel(log.WARNING)

	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showCopyright()

	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "All rights reserved.") ||
		!strings.Contains(msg, "2021") {
		t.Error("Copyright does not contain expected terms")
	}
}

func TestFoundClassFileWithNoArgs(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	// redirecting stdout to avoid clutter in the test results
	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "main.class"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	os.Stdout = normalStdout

	if global.StartingClass != "main.class" {
		t.Error("main.class not identified as starting class. Got: " +
			global.StartingClass)
	}

	if len(global.AppArgs) != 0 {
		t.Error("app arg to main.class should be empty, but got: " +
			global.AppArgs[0])
	}
}

// make sure that if a file path to the executable has an embedded :
// (as it might under Windows), that it's not mistaken for an option
// with an embedded argument (JACOBIN-2)
func TestClassFileColonIFilePath(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	// redirecting stdout to avoid clutter in the test results
	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "d:main.class"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	os.Stdout = normalStdout

	if global.StartingClass != "d:main.class" {
		t.Error("d:main.class not identified as starting class. Got: " +
			global.StartingClass)
	}

	if len(global.AppArgs) != 0 {
		t.Error("app arg to main.class should be empty, but got: " +
			global.AppArgs[0])
	}
}

func TestFoundClassFileWithArgs(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	// redirecting stdout to avoid clutter in the test results
	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "main.class", "appArg1"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	os.Stdout = normalStdout

	if global.StartingClass != "main.class" {
		t.Error("main.class not identified as starting class. Got: " +
			global.StartingClass)
	}

	if global.AppArgs[0] != "appArg1" {
		t.Error("app arg to main.class not correct. Got: " +
			global.AppArgs[0])
	}
}

func TestGetJarFilename(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-jar", "pinkle.jar", "appArg1"}

	_ = HandleCli(args, &global)

	_ = w.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if global.StartingJar != "pinkle.jar" {
		t.Error("Name of JAR file not correctly extracted from CLI")
	}

	if global.AppArgs[0] != "appArg1" {
		t.Error("JAR file arg not correctly extracted from CLI. Expected: appArg1, got: " +
			global.AppArgs[0])
	}
}

func TestMissingJARfilename(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)
	global.Args = []string{"jacobin", "-jar"}

	_, err := getJarFilename(1, "-jar", &global)
	if err != os.ErrInvalid {
		t.Error("Missing JAR filename after -jar did not trigger the right error")
	}
}

func TestEmptyOptionForEmbeddedArg(t *testing.T) {
	_, _, err := getOptionRootAndArgs("")
	if err == nil {
		t.Error("Empty option should fail test for embedded args, but did not.")
	}
}

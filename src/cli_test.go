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
	os.Unsetenv("JAVA_TOOL_OPTIONS")
	os.Unsetenv("_JAVA_OPTIONS")
	os.Unsetenv("JDK_JAVA_OPTIONS")

	javaEnvVars := getEnvArgs()
	if javaEnvVars != "" {
		t.Error("getting non-existent Java enviroment options failed")
	}
}

// set two of the JVM environment variables and make sure
// they are fetched correctly and a space is inserted between them
func TestGetJVMenvVariablesWhenTwoArePresent(t *testing.T) {
	os.Unsetenv("JAVA_TOOL_OPTIONS")
	os.Setenv("_JAVA_OPTIONS", "Hello,")
	os.Setenv("JDK_JAVA_OPTIONS", "Jacobin!")

	javaEnvVars := getEnvArgs()
	if javaEnvVars != "Hello, Jacobin!" {
		t.Error("getting two set Java enviroment options failed: " + javaEnvVars)
	}

	// clean up the environment
	os.Unsetenv("_JAVA_OPTIONS")
	os.Unsetenv("JDK_JAVA_OPTIONS")
}

// verify the output to stderr -help option is used
func TestHandleUsageMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global := globals.InitGlobals("test")
	log.SetLogLevel(log.WARNING)
	LoadOptionsTable(Global)

	// to avoid cluttering the test results, redirect stdout
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-help"}
	HandleCli(args, &Global)

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)

	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Usage:") ||
		!strings.Contains(msg, "where options include") {
		t.Error("jacobin -help did not generate the usage message to stderr. msg was: " + msg)
	}

	if Global.ExitNow != true {
		t.Error("'jacobin -help' should have set Global.exitNow to true to signal end of processing")
	}
}

func TestShowUsageMessageExitsProperlyWith__Help(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global := globals.InitGlobals("test")
	log.SetLogLevel(log.WARNING)
	LoadOptionsTable(Global)

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	showHelpStdoutAndExit(0, "--help", &Global)

	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if Global.ExitNow != true {
		t.Error("'jacobin --help' should set Global.exitNow to true but did not")
	}
}

func TestShowVersionMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global := globals.InitGlobals("test")
	log.SetLogLevel(log.WARNING)

	// to avoid cluttering the test results, redirect stdout
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	LoadOptionsTable(Global)
	args := []string{"jacobin", "-showversion", "main.clas"}

	HandleCli(args, &Global)

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)

	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin -version did not generate the correct message to stderr. msg was: " + msg)
	}
}

func TestShow__VersionUsingOptionTable(t *testing.T) {
	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)

	normalStdout := os.Stdout
	r, wout, _ := os.Pipe()
	os.Stdout = wout

	versionStdoutThenExit(0, "--version", &Global)

	wout.Close()
	os.Stdout = normalStdout
	out, _ := ioutil.ReadAll(r)

	os.Stdout = normalStdout
	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin --version did not generate the correct msg to stdout. msg was: " + msg)
	}

	if Global.ExitNow != true {
		t.Error("--version did not set exitNow value to exit. Should be set.")
	}
}

func TestChangeLoggingLevels(t *testing.T) {
	Global := globals.InitGlobals("test")
	log.SetLogLevel(log.WARNING)
	LoadOptionsTable(Global)

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-verbose:info", "main.class"}
	HandleCli(args, &Global)

	// reset stdout and stderr to what they were before redirection
	w.Close()
	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.LogLevel != log.INFO {
		t.Error("Setting log level to INFO via command line failed")
	}

	// --- now test with FINE

	log.SetLogLevel(log.WARNING)

	normalStdout = os.Stdout
	_, wout, _ = os.Pipe()
	os.Stdout = wout

	normalStderr = os.Stderr
	_, w, _ = os.Pipe()
	os.Stderr = w

	LoadOptionsTable(Global)
	args = []string{"jacobin", "-verbose:fine", "main.class"}
	HandleCli(args, &Global)

	w.Close()
	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.LogLevel != log.FINE {
		t.Error("Setting log level to FINE via command line failed")
	}

	// --- now try with FINEST

	log.SetLogLevel(log.WARNING)

	normalStdout = os.Stdout
	_, wout, _ = os.Pipe()
	os.Stdout = wout

	normalStderr = os.Stderr
	_, w, _ = os.Pipe()
	os.Stderr = w

	LoadOptionsTable(Global)
	args = []string{"jacobin", "-verbose:finest", "main.class"}
	HandleCli(args, &Global)

	w.Close()
	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.LogLevel != log.FINEST {
		t.Error("Setting log level to FINEST via command line failed")
	}

	// --- finally test with CLASS

	log.SetLogLevel(log.WARNING)

	normalStdout = os.Stdout
	_, wout, _ = os.Pipe()
	os.Stdout = wout

	normalStderr = os.Stderr
	_, w, _ = os.Pipe()
	os.Stderr = w

	LoadOptionsTable(Global)
	args = []string{"jacobin", "-verbose:class", "main.class"}
	HandleCli(args, &Global)

	w.Close()
	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if log.LogLevel != log.CLASS {
		t.Error("Setting log level to CLASS via command line failed")
	}
}

func TestInvalidLoggingLevel(t *testing.T) {
	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)
	log.SetLogLevel(log.WARNING)

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	_, err := verbosityLevel(0, "severe", &Global)

	w.Close()
	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if err == nil {
		t.Error("Setting log level to SEVERE via command line did not generate expected error")
	}
}

func TestSpecifyClientVM(t *testing.T) {

	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)
	if Global.VmModel != "server" {
		t.Error("Initialization of Global.vmModel was not set to 'server' Got: " +
			Global.VmModel)
	}

	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "-client"}
	HandleCli(args, &Global)

	// restore stdout to what it was before
	w.Close()
	os.Stdout = normalStdout

	if Global.VmModel != "client" {
		t.Error("Global.vmModel should be set to 'client'. Instead got: " +
			Global.VmModel)
	}
}

func TestSpecifyValidButUnsupportedOption(t *testing.T) {

	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)

	// redirect stdout to avoid cluttering test results with copyright notice
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "--dry-run", "main.class"}
	HandleCli(args, &Global)

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "not currently supported") {
		t.Error("Unsupported but valid option not identified. Instead got" + msg)
	}
}

func TestShowCopyright(t *testing.T) {
	globals.InitGlobals("test")
	log.SetLogLevel(log.WARNING)

	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showCopyright()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "All rights reserved.") ||
		!strings.Contains(msg, "2021") {
		t.Error("Copyright does not contain expected terms")
	}
}

func TestFoundClassFileWithNoArgs(t *testing.T) {
	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)

	// redirecting stdout to avoid clutter in the test results
	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "main.class"}
	HandleCli(args, &Global)

	w.Close()
	os.Stdout = normalStdout

	if Global.StartingClass != "main.class" {
		t.Error("main.class not identified as starting class. Got: " +
			Global.StartingClass)
	}

	if len(Global.AppArgs) != 0 {
		t.Error("app arg to main.class should be empty, but got: " +
			Global.AppArgs[0])
	}
}

func TestFoundClassFileWithArgs(t *testing.T) {
	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)

	// redirecting stdout to avoid clutter in the test results
	normalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"jacobin", "main.class", "appArg1"}
	HandleCli(args, &Global)

	w.Close()
	os.Stdout = normalStdout

	if Global.StartingClass != "main.class" {
		t.Error("main.class not identified as starting class. Got: " +
			Global.StartingClass)
	}

	if Global.AppArgs[0] != "appArg1" {
		t.Error("app arg to main.class not correct. Got: " +
			Global.AppArgs[0])
	}
}

func TestGetJarFilename(t *testing.T) {
	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-jar", "pinkle.jar", "appArg1"}

	HandleCli(args, &Global)

	w.Close()
	wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if Global.StartingJar != "pinkle.jar" {
		t.Error("Name of JAR file not correctly extracted from CLI")
	}

	if Global.AppArgs[0] != "appArg1" {
		t.Error("JAR file arg not correctly extracted from CLI. Expected: appArg1, got: " +
			Global.AppArgs[0])
	}
}

func TestMissingJARfilename(t *testing.T) {
	Global := globals.InitGlobals("test")
	LoadOptionsTable(Global)
	Global.Args = []string{"jacobin", "-jar"}

	_, err := getJarFilename(1, "-jar", &Global)
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

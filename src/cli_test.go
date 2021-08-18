/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"io/ioutil"
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

// verify the output to stderr when only usage info is requested (i.e., jacobin -help)
func TestHandleUsageMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global = initGlobals(os.Args[0])
	SetLogLevel(WARNING)
	LoadOptionsTable(Global)

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	args := []string{"jacobin", "-help"}
	HandleCli(args)

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Usage:") ||
		!strings.Contains(msg, "where options include") {
		t.Error("jacobin -help did not generate the usage message to stderr. msg was: " + msg)
	}

	if Global.exitNow != true {
		t.Error("'jacobin -help' should have set Global.exitNow to true to signal end of processing")
	}
}

func TestHandleShowVersionMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global = initGlobals(os.Args[0])
	SetLogLevel(WARNING)

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	LoadOptionsTable(Global)
	args := []string{"jacobin", "-showversion"}

	HandleCli(args)

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin -version did not generate the correct message to stderr. msg was: " + msg)
	}
}

func TestChangeLoggingLevel(t *testing.T) {
	Global = initGlobals(os.Args[0])
	SetLogLevel(WARNING)

	LoadOptionsTable(Global)
	args := []string{"jacobin", "-verbose:info"}
	HandleCli(args)

	if Global.logLevel != INFO {
		t.Error("Setting log level to INFO via command line failed")
	}
}

func TestInvalidLoggingLevel(t *testing.T) {
	Global = initGlobals(os.Args[0])
	LoadOptionsTable(Global)
	_, err := verbosityLevel(0, "severe")
	if err == nil {
		t.Error("Setting log level to SEVERE via command line did not generate expected error")
	}
}

func TestShowCopyright(t *testing.T) {
	Global = initGlobals(os.Args[0])
	SetLogLevel(WARNING)

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
	Global = initGlobals(os.Args[0])
	LoadOptionsTable(Global)
	args := []string{"jacobin", "main.class"}
	HandleCli(args)

	if Global.startingClass != "main.class" {
		t.Error("main.class not identified as starting class. Got: " +
			Global.startingClass)
	}

	if len(Global.appArgs) != 0 {
		t.Error("app arg to main.class should be empty, but got: " +
			Global.appArgs[0])
	}
}

func TestFoundClassFileWithArgs(t *testing.T) {
	Global = initGlobals(os.Args[0])
	LoadOptionsTable(Global)
	args := []string{"jacobin", "main.class", "appArg1"}
	HandleCli(args)

	if Global.startingClass != "main.class" {
		t.Error("main.class not identified as starting class. Got: " +
			Global.startingClass)
	}

	if Global.appArgs[0] != "appArg1" {
		t.Error("app arg to main.class not correct. Got: " +
			Global.appArgs[0])
	}
}

func TestGetJarFilename(t *testing.T) {
	Global = initGlobals(os.Args[0])
	LoadOptionsTable(Global)
	args := []string{"jacobin", "-jar", "pinkle.jar", "appArg1"}

	HandleCli(args)
	if Global.startingJar != "pinkle.jar" {
		t.Error("Name of JAR file not correctly extracted from CLI")
	}

	if Global.appArgs[0] != "appArg1" {
		t.Error("JAR file arg not correctly extracted from CLI. Expected: appArg1, got: " +
			Global.appArgs[0])
	}
}

func TestMissingJARfilename(t *testing.T) {
	Global = initGlobals(os.Args[0])
	LoadOptionsTable(Global)
	Global.args = []string{"jacobin", "-jar"}

	_, err := getJarFilename(1, "-jar")
	if err != os.ErrInvalid {
		t.Error("Missing JAR filename after -jar did not trigger the right error")
	}
}

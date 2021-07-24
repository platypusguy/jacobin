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
func TestHandleUserMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global = initGlobals(os.Args[0])
	SetLogLevel(WARNING)

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	handleUserMessages("jacobin -help")

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Usage:") ||
		!strings.Contains(msg, "where options include") {
		t.Error("jacobin -help did not generate the usage message to stderr. msg was: " + msg)
	}
}

func TestHandleUserMessageSignalsShutdown(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	Global = initGlobals(os.Args[0])
	SetLogLevel(WARNING)

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	stopProcessing := handleUserMessages("jacobin -help")

	// restore stderr to what it was before
	w.Close()
	os.Stdout = normalStderr

	if stopProcessing != true {
		t.Error("'jacobin -help' should have returned true to signal end of processing")
	}
}

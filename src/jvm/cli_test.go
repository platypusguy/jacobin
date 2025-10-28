/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io"
	"jacobin/src/globals"
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
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Usage:") ||
		!strings.Contains(msg, "where options include") {
		t.Error("jacobin -help did not generate the usage message to stderr. msg was: " + msg)
	}

	if global.ExitNow != true {
		t.Error("'jacobin -help' should have set globPtr.exitNow to true to signal end of processing")
	}
}

func TestShowUsageMessageExitsProperlyWith__Help(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	global := globals.InitGlobals("test")
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
		t.Error("'jacobin --help' should set globPtr.exitNow to true but did not")
	}
}

func TestShowVersionMessage(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test
	global := globals.InitGlobals("test")

	// to avoid cluttering the test results, redirect stdout
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// redirect stderr to capture writing to it
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	LoadOptionsTable(global)
	args := []string{"jacobin", "-showversion", " clas"}

	_ = HandleCli(args, &global)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin -showversion did not generate the correct message to stderr. msg was: " + msg)
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
	out, _ := io.ReadAll(r)

	os.Stdout = normalStdout
	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin VM v.") {
		t.Error("jacobin --version did not generate the correct msg to stdout. msg was: " + msg)
	}

	if global.ExitNow != true {
		t.Error("--version did not set exitNow value to exit. Should be set.")
	}
}

func TestInvalidTraceSelection(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)
	var err error

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, werr, _ := os.Pipe()
	os.Stderr = werr

	options := "-trace:inst" + TraceSep + "class" + TraceSep + "mickey"
	args := []string{"jacobin", options}
	err = HandleCli(args, &global)

	_ = werr.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if err == nil {
		t.Errorf("%s failed to generate the expected error", options)
	}
	t.Logf("HandleCli err: %v\n", err)
}

func TestValidTraceSelection(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)
	var err error

	// to avoid cluttering the test results, redirect stdout and stderr
	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	normalStderr := os.Stderr
	_, werr, _ := os.Pipe()
	os.Stderr = werr

	options := "-trace:inst" + TraceSep + "class" + TraceSep + "inst"
	args := []string{"jacobin", options}
	err = HandleCli(args, &global)

	_ = werr.Close()
	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr

	if err != nil {
		t.Errorf("HandleCli err: %v\n", err)
	}
	if globals.TraceInst && globals.TraceClass && (!globals.TraceVerbose) {
		return
	}
	t.Errorf("globals.TraceInst = %t (true), globals.TraceClass = %t (true), globals.TraceVerbose = %t (false)\n",
		globals.TraceInst, globals.TraceClass, globals.TraceVerbose)
}

func TestSpecifyClientVM(t *testing.T) {

	global := globals.InitGlobals("test")
	LoadOptionsTable(global)
	if global.VmModel != "server" {
		t.Error("Initialization of globPtr.vmModel was not set to 'server' Got: " +
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

	args := []string{"jacobin", "--dry-run", " class"}
	_ = HandleCli(args, &global)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "not currently supported") {
		t.Error("Unsupported but valid option not identified. Instead got" + msg)
	}
}

func TestShowCopyrightInVersion(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.StrictJDK = false // Copyright is shown in a run only when not in strictJDK mode

	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// copyright appears only if one of the -version family of
	// commands has been specified.
	g.CommandLine = "-version"
	showCopyright(g)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "Copyright") ||
		!strings.Contains(msg, "2021") {
		t.Error("Copyright does not contain expected terms")
	}
}

func TestShowCopyrightWithStrictJDKswitch(t *testing.T) {
	g := globals.GetGlobalRef()
	globals.InitGlobals("test")
	g.StrictJDK = true // Copyright is shown in a run only when not in strictJDK mode

	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showCopyright(g)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = normalStdout

	msg := string(out[:])

	if msg != "" {
		t.Errorf("Expected no copyright notice, but got: %s", msg)
	}
}

// Command line:  jacobin  a.class
func TestFoundClassFileWithNoArgs(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	args := []string{"jacobin", "a.class"}
	_ = HandleCli(args, &global)

	if global.StartingClass != "a.class" {
		t.Errorf("Expected global.StartingClass = \"a.class\", observed: \"%s\"", global.StartingClass)
	}

	if len(global.AppArgs) != 0 {
		t.Errorf("Expected global.AppArgs = 0, observed: %d", len(global.AppArgs))
	}
}

// Command line:  jacobin  a.class  apple  banana  peach
func TestFoundClassFileWithArgs(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	args := []string{"jacobin", "a.class", "apple", "banana", "peach"}
	_ = HandleCli(args, &global)

	if global.StartingClass != "a.class" {
		t.Errorf("Expected global.StartingClass = \"a.class\", observed: \"%s\"", global.StartingClass)
	}

	if len(global.AppArgs) != 3 {
		t.Errorf("Expected global.AppArgs = 3, observed: %d", len(global.AppArgs))
	}

	if global.AppArgs[0] != "apple" {
		t.Errorf("Expected global.AppArgs[0] = \"apple\", observed: \"%s\"", global.AppArgs[0])
	}

	if global.AppArgs[1] != "banana" {
		t.Errorf("Expected global.AppArgs[1] = \"banana\", observed: \"%s\"", global.AppArgs[1])
	}

	if global.AppArgs[2] != "peach" {
		t.Errorf("Expected global.AppArgs[2] = \"peach\", observed: \"%s\"", global.AppArgs[2])
	}
}

// Command line:  jacobin  Starter  --double-dash  apple  banana  peach
func TestFoundClassFileWithArgsAlt(t *testing.T) {
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	args := []string{"jacobin", "Starter", "--double-dash", "apple", "banana", "peach"} // no .class suffix
	_ = HandleCli(args, &global)

	if global.StartingClass != "Starter.class" {
		t.Errorf("Expected global.StartingClass = \"Starter.class\", observed: \"%s\"", global.StartingClass)
	}

	if len(global.AppArgs) != 4 {
		t.Errorf("Expected global.AppArgs = 4, observed: %d", len(global.AppArgs))
	} else {
		if global.AppArgs[0] != "--double-dash" {
			t.Errorf("Expected global.AppArgs[0] = \"--double-dash\", observed: \"%s\"", global.AppArgs[0])
		}
		if global.AppArgs[1] != "apple" {
			t.Errorf("Expected global.AppArgs[0] = \"apple\", observed: \"%s\"", global.AppArgs[1])
		}
		if global.AppArgs[2] != "banana" {
			t.Errorf("Expected global.AppArgs[1] = \"banana\", observed: \"%s\"", global.AppArgs[2])
		}
		if global.AppArgs[3] != "peach" {
			t.Errorf("Expected global.AppArgs[2] = \"peach\", observed: \"%s\"", global.AppArgs[3])
		}
	}
}

// Command line:  jacobin  -classpath Mercury:Venus:Earth:Mars  Starter  --double-dash  apple  banana  peach
func _execWithClasspath(t *testing.T, optName string) {

	t.Logf("_execWithClasspath: option name = \"%s\"", optName)
	global := globals.InitGlobals("test")
	LoadOptionsTable(global)

	args := []string{"jacobin", optName, "Mercury:Venus:Earth:Mars", "Starter", "--double-dash", "apple", "banana", "peach"} // no .class suffix
	_ = HandleCli(args, &global)

	if global.StartingClass != "Starter.class" {
		t.Errorf("Expected global.StartingClass = \"Starter.class\", observed: \"%s\"", global.StartingClass)
	}

	if global.ClasspathRaw != "Mercury:Venus:Earth:Mars" {
		t.Errorf("Expected global.ClasspathRaw = \"Mercury:Venus:Earth:Mars\", observed: \"%s\"", global.ClasspathRaw)
	}

	if len(global.Classpath) != 4 {
		t.Errorf("Expected global.Classpath = 4, observed: %d", len(global.Classpath))
	} else {
		if global.Classpath[0] != "Mercury/" {
			t.Errorf("Expected global.Classpath[0] = \"Mercury/\", observed: \"%s\"", global.Classpath[0])
		}
		if global.Classpath[1] != "Venus/" {
			t.Errorf("Expected global.Classpath[1] = \"Venus/\", observed: \"%s\"", global.Classpath[1])
		}
		if global.Classpath[2] != "Earth/" {
			t.Errorf("Expected global.Classpath[2] = \"Earth/\", observed: \"%s\"", global.Classpath[2])
		}
		if global.Classpath[3] != "Mars/" {
			t.Errorf("Expected global.Classpath[3] = \"Mars/\", observed: \"%s\"", global.Classpath[3])
		}
	}

	if len(global.AppArgs) != 4 {
		t.Errorf("Expected global.AppArgs = 4, observed: %d", len(global.AppArgs))
	} else {
		if global.AppArgs[0] != "--double-dash" {
			t.Errorf("Expected global.AppArgs[0] = \"--double-dash\", observed: \"%s\"", global.AppArgs[0])
		}
		if global.AppArgs[1] != "apple" {
			t.Errorf("Expected global.AppArgs[0] = \"apple\", observed: \"%s\"", global.AppArgs[1])
		}
		if global.AppArgs[2] != "banana" {
			t.Errorf("Expected global.AppArgs[1] = \"banana\", observed: \"%s\"", global.AppArgs[2])
		}
		if global.AppArgs[3] != "peach" {
			t.Errorf("Expected global.AppArgs[2] = \"peach\", observed: \"%s\"", global.AppArgs[3])
		}
	}
}

func TestWithCp1(t *testing.T) {
	_execWithClasspath(t, "-cp")
	_execWithClasspath(t, "-classpath")
	_execWithClasspath(t, "--class-path")
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

	args := []string{"jacobin", "d:a.class"}
	_ = HandleCli(args, &global)

	_ = w.Close()
	os.Stdout = normalStdout

	if global.StartingClass != "d:a.class" {
		t.Error("d:a.class not identified as starting class. Got: " +
			global.StartingClass)
	}

	if len(global.AppArgs) != 0 {
		t.Error("app arg to  class should be empty, but got: " +
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

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"jacobin/src/globals"
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

func TestEnableAssertions(t *testing.T) {
	global := globals.InitGlobals("test")
	statics.LoadProgramStatics()

	global.Args = []string{"-ea"}

	pos, err := enableAssertions(0, "", &global)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assertStatus := statics.GetStaticValue("main", "$assertionsDisabled")
	if assertStatus.(int64) != (types.JavaBoolFalse) {
		t.Error("Expected assertions to be enabled, but it is not.")
	}

	if pos != 0 {
		t.Errorf("Expected position 0, got %d", pos)
	}
}

func TestExpandClasspathWithJarFile(t *testing.T) {
	globals.InitGlobals("test")

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	// Ensure the directory is deleted after the test
	defer os.RemoveAll(tempDir)

	// Create an empty file named abc.JAR in the temp directory
	tempFileName := tempDir + string(os.PathSeparator) + "abc.JAR"
	file, err := os.Create(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Close the file after creation
	file.Close()

	gl := globals.GetGlobalRef()
	gl.Classpath = make([]string, 0)
	gl.ClasspathRaw = tempFileName + string(os.PathListSeparator) + "a"
	expandClasspth(gl)

	expected := strings.Split(gl.ClasspathRaw, string(os.PathListSeparator))
	for i, path := range expected {
		if !strings.HasSuffix(path, string(os.PathSeparator)) &&
			!strings.HasSuffix(path, ".JAR") && !strings.HasSuffix(path, ".jar") {
			expected[i] = path + string(os.PathSeparator)
		}
	}
	if !equalSlices(gl.Classpath, expected) {
		t.Errorf("Expected classpath %v, got %v", expected, gl.Classpath)
	}
}

func TestExpandClasspathWithWildcard(t *testing.T) {
	globals.InitGlobals("test")

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	// Ensure the directory is deleted after the test
	defer os.RemoveAll(tempDir)

	// Create an empty file named abc.JAR in the temp directory
	tempFileName := tempDir + string(os.PathSeparator) + "abc.JAR"
	file, err := os.Create(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Close the file after creation
	file.Close()

	gl := globals.GetGlobalRef()
	gl.Classpath = make([]string, 0)
	gl.ClasspathRaw = tempDir + string(os.PathSeparator) + "*" + string(os.PathListSeparator) + "a"
	expandClasspth(gl)

	expected := make([]string, 2)
	expected[0] = tempDir + string(os.PathSeparator) + "abc.JAR"
	expected[1] = "a" + string(os.PathSeparator)

	if !equalSlices(gl.Classpath, expected) {
		t.Errorf("Expected classpath %v, got %v", expected, gl.Classpath)
	}
}

func TestGetClasspathValidInput(t *testing.T) {
	global := globals.InitGlobals("test")
	separator := string(os.PathListSeparator)
	pathArg := "a" + separator + "b" + separator + "c"
	global.Args = []string{"-cp", pathArg}

	pos, err := getClasspath(0, pathArg, &global)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Split(pathArg, string(os.PathListSeparator))
	for i, path := range expected {
		if !strings.HasSuffix(path, string(os.PathSeparator)) {
			expected[i] = path + string(os.PathSeparator)
		}
	}
	if !equalSlices(global.Classpath, expected) {
		t.Errorf("Expected classpath %v, got %v", expected, global.Classpath)
	}

	if pos != 1 {
		t.Errorf("Expected position 1, got %d", pos)
	}
}

func TestGetClasspathMissingArgument(t *testing.T) {
	global := globals.InitGlobals("test")
	global.Args = []string{"-cp"}

	pos, err := getClasspath(0, "", &global)
	if err == nil || err.Error() != "missing classpath after -cp or -classpath option" {
		t.Errorf("Expected error for missing classpath, got: %v", err)
	}

	if pos != 0 {
		t.Errorf("Expected position 0, got %d", pos)
	}
}

func TestGetJarFilenameValid(t *testing.T) {
	global := globals.InitGlobals("test")
	global.Args = []string{"-jar", "filename.jar"}

	pos, _ := getJarFilename(0, "", &global)
	if global.StartingJar != "filename.jar" {
		t.Errorf("Did not get expected JAR file name, got: %s", global.StartingJar)
	}

	if pos != 2 {
		t.Errorf("Expected position 2, got %d", pos)
	}
}

func TestGetJarFilenameWithAppArgs(t *testing.T) {
	global := globals.InitGlobals("test")
	global.Args = []string{"-jar", "filename.jar", "arg1", "arg2"}

	pos, _ := getJarFilename(0, "", &global)
	if global.StartingJar != "filename.jar" {
		t.Errorf("Did not get expected JAR file name, got: %s", global.StartingJar)
	}

	if global.AppArgs[0] != "arg1" {
		t.Errorf("Did not get expected app arg, got: %s", global.AppArgs[0])
	}

	if pos != 4 {
		t.Errorf("Expected position 2, got %d", pos)
	}
}
func TestGetJarFilenameMissingArgument(t *testing.T) {
	global := globals.InitGlobals("test")
	global.Args = []string{"-jar"}

	pos, err := getJarFilename(0, "", &global)
	if err != os.ErrInvalid {
		t.Errorf("Expected error for missing jar file name, got: %v", err)
	}

	if pos != 0 {
		t.Errorf("Expected position 0, got %d", pos)
	}
}

// Helper function to compare slices
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

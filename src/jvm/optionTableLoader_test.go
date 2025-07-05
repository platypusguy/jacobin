/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"jacobin/globals"
	"os"
	"strings"
	"testing"
)

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

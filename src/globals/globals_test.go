/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"fmt"
	"os"
	"testing"
)

var foobar string
var foo string

func nameFooBar(t *testing.T) bool {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("nameFooBar: os.UserHomeDir failed: %s", err.Error())
		return false
	}
	foo = userHomeDir + string(os.PathSeparator) + "foo"
	foobar = foo + string(os.PathSeparator) + "bar"
	return true
}

func TestGlobalsInit(t *testing.T) {
	g := InitGlobals("testInit")

	if g.JacobinName != "testInit" {
		t.Errorf("Expecting globals init to set Jacobin name to 'testInit', got: %s", g.JacobinName)
	}

	if g.VmModel != "server" {
		t.Errorf("Expected globals init to set VmModel to 'server', got: %s", g.VmModel)
	}
}

// make sure the JAVA_HOME environment variable is extracted and the embedded slashes
// are reformatted correctly
func TestJavaHomeFormat(t *testing.T) {
	if !nameFooBar(t) {
		return
	}
	defer os.RemoveAll(foo)

	origJavaHome := os.Getenv("JAVA_HOME")
	_ = os.Setenv("JAVA_HOME", foobar)
	InitJavaHome()
	ret := JavaHome()
	expectedPath := foobar
	if ret != expectedPath {
		t.Errorf("Expecting a JAVA_HOME of '%s', got: %s", expectedPath, ret)
	}
	_ = os.Setenv("JAVA_HOME", origJavaHome)
}

// test original java home + version
func TestJavaHomeAndVersion(t *testing.T) {
	InitJavaHome()
	home := JavaHome()
	version := JavaVersion()
	if home == "" {
		t.Errorf("JAVA_HOME is nil")
	}
	if version == "" {
		t.Errorf("JAVA_VERSION is nil")
	}
	fmt.Printf("TestJavaHomeAndVersion: JAVA_HOME=%s, JAVA_VERSION=%s\n", home, version)
}

// verify that a trailing slash in JAVA_HOME is removed
func TestJavaHomeRemovalOfTrailingSlash(t *testing.T) {
	if !nameFooBar(t) {
		return
	}
	defer os.RemoveAll(foobar)

	origJavaHome := os.Getenv("JAVA_HOME")
	_ = os.Setenv("JAVA_HOME", foobar)
	InitJavaHome()
	ret := JavaHome()
	expectedPath := foobar
	if ret != expectedPath {
		t.Errorf("Expecting a JAVA_HOME of '%s', got: %s", expectedPath, ret)
	}
	_ = os.Setenv("JAVA_HOME", origJavaHome)
}

// make sure the JACOBIN_HOME environment variable is extracted and reformatted correctly
// Per JACOBIN-184, the trailing slash is removed.
func TestJacobinHomeFormat(t *testing.T) {
	if !nameFooBar(t) {
		return
	}
	defer os.RemoveAll(foo)

	origJavaHome := os.Getenv("JACOBIN_HOME")
	_ = os.Setenv("JACOBIN_HOME", foobar)
	InitJacobinHome()
	ret := JacobinHome()
	expectedPath := foobar
	if ret != expectedPath {
		t.Errorf("Expecting a JACOBIN_HOME of '%s', got: %s", expectedPath, ret)
	}
	_ = os.Setenv("JACOBIN_HOME", origJavaHome)
}

// verify that a trailing slash in JAVA_HOME is removed
func TestJacobinHomeRemovalOfTrailingSlash(t *testing.T) {
	if !nameFooBar(t) {
		return
	}
	defer os.RemoveAll(foo)

	origJavaHome := os.Getenv("JACOBIN_HOME")
	_ = os.Setenv("JACOBIN_HOME", foobar)
	InitJacobinHome()
	ret := JacobinHome()
	expectedPath := foobar
	if ret != expectedPath {
		t.Errorf("Expecting a JACOBIN_HOME of '%s', got: %s", expectedPath, ret)
	}
	_ = os.Setenv("JACOBIN_HOME", origJavaHome)
}

func TestVariousInitialDefaultValues(t *testing.T) {
	InitGlobals("testInit")
	gl := GetGlobalRef()
	if gl.StrictJDK != false ||
		gl.ExitNow != false ||
		!(gl.MaxJavaVersion >= 11) {
		t.Errorf("Some global variables intialized to unexpected values.")
	}
}

func TestGetSystemProperty(t *testing.T) {
	InitGlobals("testInit")
	buildGlobalProperties()
	ret := GetSystemProperty("java.version")
	if ret < "21" {
		t.Errorf("Expecting a java.version of 21 or more, got: %s", ret)
	}
}

func TestGetSystemPropertyNotFound(t *testing.T) {
	InitGlobals("testInit")
	buildGlobalProperties()
	ret := GetSystemProperty("java.version.notfound")
	if ret != "" {
		t.Errorf("Expecting a java.version.notfound of '', got: %s", ret)
	}
}

func TestSetSystemProperty(t *testing.T) {
	InitGlobals("testInit")
	buildGlobalProperties()
	SetSystemProperty("java.version", "22")
	ret := GetSystemProperty("java.version")
	if ret != "22" {
		t.Errorf("Expecting a java.version of 22, got: %s", ret)
	}
}

func TestRemoveSystemProperty(t *testing.T) {
	InitGlobals("testInit")
	buildGlobalProperties()
	SetSystemProperty("java.version", "22")
	ret := GetSystemProperty("java.version")
	if ret != "22" {
		t.Errorf("Expecting a java.version of 22, got: %s", ret)
	}
	RemoveSystemProperty("java.version")
	ret = GetSystemProperty("java.version")
	if ret != "" {
		t.Errorf("Expecting a java.version of '', got: %s", ret)
	}
}

func TestReplaceSystemProperties(t *testing.T) {
	InitGlobals("testInit")
	buildGlobalProperties()
	SetSystemProperty("java.version", "22")
	ret := GetSystemProperty("java.version")
	if ret != "22" {
		t.Errorf("Expecting a java.version of 22, got: %s", ret)
	}

	newMap := make(map[string]string)
	newMap["java.version"] = "23"
	ReplaceSystemProperties(newMap)
	ret = GetSystemProperty("java.version")
	if ret != "23" {
		t.Errorf("Expecting a java.version of 23, got: %s", ret)
	}
}

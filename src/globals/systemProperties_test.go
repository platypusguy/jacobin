/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package globals

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

// tests are in alphabetical order

func TestGetFileEncoding(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("file.encoding")
	if runtime.GOOS == "windows" {
		if ret != "windows-1252" && ret != "UTF-8" {
			t.Errorf("Expecting a file.encoding of windows-1252 or UTF-8 on Windows, got: %s", ret)
		}
	} else if ret != "UTF-8" {
		t.Errorf("Expecting a file.encoding of UTF-8, got: %s", ret)
	}
}

func TestGetFileNameEncoding(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("sun.jnu.encoding")
	if ret != "UTF-8" {
		t.Errorf("Expecting a filename encoding (sun.jnu.encoding) of UTF-8, got: %s", ret)
	}
}

func TestGetJDKmajorVersionInvalid(t *testing.T) {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	prevJavaHomeEnv := os.Getenv("JAVA_HOME")
	_ = os.Setenv("JAVA_HOME", "nonexistent")
	InitGlobals("test")
	ret := GetSystemProperty("jdk.major.version")
	if ret != "" { // should be empty if JAVA_HOME is invalid
		t.Errorf("Expecting a jdk.major.version of '', got: %s", ret)
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr
	errMsg := string(msg)

	if !strings.Contains(errMsg, "Cannot find the specified path: ") {
		t.Errorf("Expected error message containing 'Cannot find the specified path: ', got: %s", errMsg)
	}
	_ = os.Setenv("JAVA_HOME", prevJavaHomeEnv)
}

func TestGetSystemClasspath(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("java.class.path")
	if ret != "." {
		t.Errorf("Expecting a java.class.path of ., got: %s", ret)
	}
}

func TestGetSystemProperty(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("java.vm.vendor")
	if ret < "Jacobin" {
		t.Errorf("Expecting java.vm.vendor = 'Jacobin', got: %s", ret)
	}
}

func TestGetSystemPropertyNonExistent(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("non.existent.property")
	if ret != "" {
		t.Errorf("Expecting a non.existent.property of '', got: %s", ret)
	}
}

func TestGetSystemPropertyJNUencoding(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("sun.jnu.encoding")
	if ret != "UTF-8" {
		t.Errorf("Expecting a sun.jnu.encoding of UTF-8, got: %s", ret)
	}
}

func TestRemoveSystemProperty(t *testing.T) {
	InitGlobals("test")
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
	InitGlobals("test")
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

func TestSetSystemProperty(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	SetSystemProperty("java.version", "22")
	ret := GetSystemProperty("java.version")
	if ret != "22" {
		t.Errorf("Expecting a java.version of 22, got: %s", ret)
	}
}

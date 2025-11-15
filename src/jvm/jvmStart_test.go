/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"archive/zip"
	"io"
	"jacobin/src/globals"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeTempJar creates a minimal jar file with an optional manifest
// and arbitrary entries. Returns the path and a cleanup function.
func makeTempJar(t *testing.T, manifest map[string]string, files map[string][]byte) (string, func()) {
	t.Helper()

	dir := t.TempDir()
	jarPath := filepath.Join(dir, "test.jar")
	f, err := os.Create(jarPath)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	zw := zip.NewWriter(f)
	if manifest != nil {
		mw, err := zw.Create("META-INF/MANIFEST.MF")
		if err != nil {
			t.Fatalf("manifest create: %v", err)
		}
		for k, v := range manifest {
			if _, err := io.WriteString(mw, k+": "+v+"\n"); err != nil {
				t.Fatalf("manifest write: %v", err)
			}
		}
	}
	for p, content := range files {
		w, err := zw.Create(p)
		if err != nil {
			t.Fatalf("entry create %s: %v", p, err)
		}
		if _, err := w.Write(content); err != nil {
			t.Fatalf("entry write %s: %v", p, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zw close: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("file close: %v", err)
	}
	return jarPath, func() { _ = os.Remove(jarPath) }
}

// withArgs temporarily replaces os.Args during the callback and restores them afterward.
func withArgs(t *testing.T, args []string, fn func()) {
	t.Helper()
	old := os.Args
	os.Args = args
	defer func() { os.Args = old }()
	fn()
}

// captureStdout captures stdout during the function, returning the captured string.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old
	return string(out)
}

func TestJVMrun_ExitNow_WithVersionFlag(t *testing.T) {
	// Initialize globals in test mode so JVMrun won't re-init them.
	globals.InitGlobals("test")

	withArgs(t, []string{"jacobin", "-version"}, func() {
		code := JVMrun()
		// In test mode, shutdown.Exit maps OK to TEST_OK which returns 0
		if code != 0 {
			t.Fatalf("expected TEST_OK (0) for -version early exit, got %d", code)
		}
	})
}

func TestJVMrun_HandleCliError_UnknownOption(t *testing.T) {
	globals.InitGlobals("test")

	withArgs(t, []string{"jacobin", "-notAnOption"}, func() {
		code := JVMrun()
		// In test mode, non-OK exit maps to TEST_ERR which returns 1
		if code != 1 {
			t.Fatalf("expected TEST_ERR (1) for CLI error, got %d", code)
		}
	})
}

func TestJVMrun_NoStartingTarget_ShowsUsage_AndAppException(t *testing.T) {
	globals.InitGlobals("test")

	out := captureStdout(t, func() {
		withArgs(t, []string{"jacobin"}, func() {
			code := JVMrun()
			if code != 1 { // APP_EXCEPTION -> TEST_ERR (1)
				t.Fatalf("expected TEST_ERR (1) when no starting class/jar, got %d", code)
			}
		})
	})

	if !strings.Contains(out, "Usage: jacobin") {
		t.Fatalf("expected usage text on stdout when no start target, got: %q", out)
	}
}

func TestJVMrun_JarWithoutMainClass_AppException(t *testing.T) {
	globals.InitGlobals("test")

	// Create a jar with no Main-Class manifest attribute
	jarPath, cleanup := makeTempJar(t, map[string]string{"Class-Path": "lib/a.jar"}, map[string][]byte{})
	defer cleanup()

	withArgs(t, []string{"jacobin", "-jar", jarPath}, func() {
		code := JVMrun()
		if code != 1 { // APP_EXCEPTION -> TEST_ERR (1)
			t.Fatalf("expected TEST_ERR (1) for jar without Main-Class, got %d", code)
		}
	})
}

// jacobin -jar testdata/hello.jar
// Expectation: the jar has a Main-Class and a simple HelloWorld class; JVMrun should
// initialize successfully and return TEST_OK (0) in test mode.
func TestJVMrun_Jar_Hello_FromTestdata_OK(t *testing.T) {
	globals.InitGlobals("test")

	// Build path to repo testdata/hello.jar relative to this package dir (src/jvm)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	jarPath := filepath.Join(pwd, "..", "..", "testdata", "hello.jar")

	withArgs(t, []string{"jacobin", "-jar", jarPath}, func() {
		code := JVMrun()
		if code != 0 {
			t.Fatalf("expected TEST_OK(0) running hello.jar, got %d", code)
		}
	})
}

// jacobin -jar testdata/nomanifest.jar
// Expectation: manifest lacks Main-Class; JVMrun should map APP_EXCEPTION to TEST_ERR (1).
func TestJVMrun_Jar_NoManifest_FromTestdata_AppException(t *testing.T) {
	globals.InitGlobals("test")

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	jarPath := filepath.Join(pwd, "..", "..", "testdata", "nomanifest.jar")

	withArgs(t, []string{"jacobin", "-jar", jarPath}, func() {
		code := JVMrun()
		if code != 1 {
			t.Fatalf("expected TEST_ERR(1) for nomanifest.jar, got %d", code)
		}
	})
}

// jacobin -cp ../../testdata/nomanifest.jar Hello.class
// Expectation: class is loaded from the jar on the classpath; JVMrun should return TEST_OK (0) in test mode.
func TestJVMrun_Class_FromClasspathJar_OK(t *testing.T) {
	globals.InitGlobals("test")

	// Build path to repo testdata/nomanifest.jar relative to this package dir (src/jvm)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	cpJar := filepath.Join(pwd, "..", "..", "testdata", "nomanifest.jar")

	withArgs(t, []string{"jacobin", "-cp", cpJar, "Hello.class"}, func() {
		code := JVMrun()
		if code != 0 {
			t.Fatalf("expected TEST_OK(0) loading Hello.class from classpath jar, got %d", code)
		}
	})
}

// jacobin -jar ../../testdata/jarring.jar
// Expectations:
// * main.class is loaded from the jar.
// * middle/calculator/Calculator.class is found and executed
// * JVMrun should return TEST_OK (0) in test mode.
func TestJVMrun_jarring(t *testing.T) {
	globals.InitGlobals("test")

	// Build path to repo testdata/nomanifest.jar relative to this package dir (src/jvm)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	execJar := filepath.Join(pwd, "..", "..", "testdata", "jarring.jar")

	withArgs(t, []string{"jacobin", "-jar", execJar}, func() {
		code := JVMrun()
		if code != 0 {
			t.Fatalf("expected TEST_OK(0) loading Hello.class from classpath jar, got %d", code)
		}
	})
}

// jacobin -jar ../../testdata/jar1.jar
//
//	(../../testdata/jar2.jar is in the Class-Path of jar1.jar)
//
// Expectations:
// * main.class is loaded from the jar.
// * middle/calculator/Calculator.class is found and executed
// * JVMrun should return TEST_OK (0) in test mode.
func TestJVMrun_jar1_and_jar2(t *testing.T) {
	globals.InitGlobals("test")

	// Build path to repo testdata/nomanifest.jar relative to this package dir (src/jvm)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	execJar := filepath.Join(pwd, "..", "..", "testdata", "jar1.jar")

	withArgs(t, []string{"jacobin", "-jar", execJar}, func() {
		code := JVMrun()
		if code != 0 {
			t.Fatalf("expected TEST_OK(0) loading Hello.class from classpath jar, got %d", code)
		}
	})
}

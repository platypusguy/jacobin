/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"bytes"
	"jacobin/src/globals"
	"os"
	"path/filepath"
	"testing"
)

// projectRoot returns the absolute path to the jacobin project root, so that
// tests can reference the real testdata/jmod fixtures regardless of the
// working directory used when `go test` is invoked.
func projectRoot(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() failed: %v", err)
	}
	// src/classloader -> project root is two levels up
	return filepath.Dir(filepath.Dir(wd))
}

// setupJavaHomeWithJmod copies testdata/jmod/<srcJmodName> into a temporary
// directory laid out as <tmp>/jmods/<destJmodName>, and points
// global.JavaHome at that temporary directory. It returns the JavaHome path.
func setupJavaHomeWithJmod(t *testing.T, srcJmodName, destJmodName string) string {
	root := projectRoot(t)
	srcPath := filepath.Join(root, "testdata", "jmod", srcJmodName)
	content, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read fixture jmod %s: %v", srcPath, err)
	}

	tmpDir := t.TempDir()
	jmodsDir := filepath.Join(tmpDir, "jmods")
	if err = os.MkdirAll(jmodsDir, 0755); err != nil {
		t.Fatalf("failed to create jmods dir: %v", err)
	}
	destPath := filepath.Join(jmodsDir, destJmodName)
	if err = os.WriteFile(destPath, content, 0644); err != nil {
		t.Fatalf("failed to write fixture jmod to %s: %v", destPath, err)
	}

	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	global.JavaHome = tmpDir
	return tmpDir
}

// --- GetBaseJmodBytes ---

func TestGetBaseJmodBytes_Success(t *testing.T) {
	setupJavaHomeWithJmod(t, "jacobin.jmod", BaseJmodFileName)

	global := globals.GetGlobalRef()
	global.JmodBaseBytes = nil

	GetBaseJmodBytes()

	if len(global.JmodBaseBytes) == 0 {
		t.Fatalf("Expected GetBaseJmodBytes to populate JmodBaseBytes, got empty slice")
	}

	root := projectRoot(t)
	expected, err := os.ReadFile(filepath.Join(root, "testdata", "jmod", "jacobin.jmod"))
	if err != nil {
		t.Fatalf("failed to read reference fixture: %v", err)
	}
	if !bytes.Equal(global.JmodBaseBytes, expected) {
		t.Errorf("JmodBaseBytes does not match the fixture file contents")
	}
}

// --- GetClassBytes ---

func TestGetClassBytes_BaseJmod_Success(t *testing.T) {
	setupJavaHomeWithJmod(t, "jacobin.jmod", BaseJmodFileName)
	global := globals.GetGlobalRef()
	global.JmodBaseBytes = nil
	GetBaseJmodBytes()

	classBytes, err := GetClassBytes(BaseJmodFileName, "org/jacobin/test/Hello")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(classBytes) == 0 {
		t.Fatalf("Expected non-empty class bytes")
	}

	root := projectRoot(t)
	expected, err := os.ReadFile(filepath.Join(root, "testdata", "jmod", "classes", "org", "jacobin", "test", "Hello.class"))
	if err != nil {
		t.Fatalf("failed to read reference fixture: %v", err)
	}
	if !bytes.Equal(classBytes, expected) {
		t.Errorf("GetClassBytes returned bytes that do not match the fixture Hello.class")
	}
}

func TestGetClassBytes_BaseJmod_ClassNotFound(t *testing.T) {
	setupJavaHomeWithJmod(t, "jacobin.jmod", BaseJmodFileName)
	global := globals.GetGlobalRef()
	global.JmodBaseBytes = nil
	GetBaseJmodBytes()

	_, err := GetClassBytes(BaseJmodFileName, "org/jacobin/test/DoesNotExist")
	if err == nil {
		t.Fatalf("Expected an error for a class that does not exist in the jmod, got nil")
	}
}

func TestGetClassBytes_NonBaseJmod_Success(t *testing.T) {
	const otherJmodName = "myModule.jmod"
	setupJavaHomeWithJmod(t, "jacobin.jmod", otherJmodName)

	classBytes, err := GetClassBytes(otherJmodName, "org/jacobin/test/Hello")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(classBytes) == 0 {
		t.Fatalf("Expected non-empty class bytes")
	}

	root := projectRoot(t)
	expected, err := os.ReadFile(filepath.Join(root, "testdata", "jmod", "classes", "org", "jacobin", "test", "Hello.class"))
	if err != nil {
		t.Fatalf("failed to read reference fixture: %v", err)
	}
	if !bytes.Equal(classBytes, expected) {
		t.Errorf("GetClassBytes returned bytes that do not match the fixture Hello.class")
	}
}

func TestGetClassBytes_NonBaseJmod_FileNotFound(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	global.JavaHome = t.TempDir() // no jmods dir/file created here

	_, err := GetClassBytes("nonexistent.jmod", "org/jacobin/test/Hello")
	if err == nil {
		t.Fatalf("Expected an error when the jmod file does not exist, got nil")
	}
}

func TestGetClassBytes_NonBaseJmod_InvalidMagicNumber(t *testing.T) {
	tmpDir := t.TempDir()
	jmodsDir := filepath.Join(tmpDir, "jmods")
	if err := os.MkdirAll(jmodsDir, 0755); err != nil {
		t.Fatalf("failed to create jmods dir: %v", err)
	}
	badJmodName := "bad.jmod"
	// content that does not start with the expected magic number bytes
	if err := os.WriteFile(filepath.Join(jmodsDir, badJmodName), []byte{0, 0, 0, 0, 1, 2, 3, 4}, 0644); err != nil {
		t.Fatalf("failed to write bad jmod file: %v", err)
	}

	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	global.JavaHome = tmpDir

	// NOTE: GetClassBytes has a pre-existing bug: on an invalid magic number
	// it returns "nil, err" where err is still nil (the os.ReadFile call
	// that set it succeeded), so no error is actually surfaced here. This
	// test documents the current behavior rather than the ideal behavior.
	classBytes, err := GetClassBytes(badJmodName, "org/jacobin/test/Hello")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if classBytes != nil {
		t.Errorf("Expected nil class bytes for invalid magic number, got %v", classBytes)
	}
}

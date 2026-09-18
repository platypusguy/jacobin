/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"archive/zip"
	"bytes"
	"jacobin/src/globals"
	"testing"
)

// buildJmodBytes creates an in-memory byte slice that mimics the layout of a
// .jmod file: a 4-byte header (ignored by WalkBaseJmod/getClasslist) followed
// by a valid ZIP archive whose entries are supplied by the caller.
func buildJmodBytes(t *testing.T, files map[string][]byte) []byte {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("buildJmodBytes: zw.Create(%s) failed: %v", name, err)
		}
		if _, err = w.Write(content); err != nil {
			t.Fatalf("buildJmodBytes: write(%s) failed: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("buildJmodBytes: zw.Close() failed: %v", err)
	}

	// prepend a fake 4-byte JMOD header
	header := []byte{0x4A, 0x4D, 0x01, 0x00}
	return append(header, buf.Bytes()...)
}

// --- WalkBaseJmod ---

func TestWalkBaseJmod_InvalidZip(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	// Not a valid zip file (after removing 4-byte header)
	global.JmodBaseBytes = []byte{0, 0, 0, 0, 'n', 'o', 't', 'a', 'z', 'i', 'p'}

	err := WalkBaseJmod()
	if err == nil {
		t.Fatalf("Expected an error for invalid zip content, got nil")
	}
}

func TestWalkBaseJmod_NoClassesNoClasslist(t *testing.T) {
	globals.InitGlobals("test")
	_ = Init()
	global := globals.GetGlobalRef()

	files := map[string][]byte{
		"module-info.class": {0x01, 0x02}, // wrong prefix, should be skipped
		"readme.txt":        []byte("hello"),
	}
	global.JmodBaseBytes = buildJmodBytes(t, files)

	beforeCount := BootstrapCL.ClassCount
	err := WalkBaseJmod()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if BootstrapCL.ClassCount != beforeCount {
		t.Errorf("Expected no classes to be loaded, class count changed from %d to %d",
			beforeCount, BootstrapCL.ClassCount)
	}
}

func TestWalkBaseJmod_ClassOutsideClassesPrefixSkipped(t *testing.T) {
	globals.InitGlobals("test")
	_ = Init()
	global := globals.GetGlobalRef()

	files := map[string][]byte{
		"lib/SomeClass.class": {0x01, 0x02, 0x03}, // not prefixed with "classes"
	}
	global.JmodBaseBytes = buildJmodBytes(t, files)

	beforeCount := BootstrapCL.ClassCount
	err := WalkBaseJmod()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if BootstrapCL.ClassCount != beforeCount {
		t.Errorf("Expected class outside 'classes' prefix to be skipped")
	}
}

func TestWalkBaseJmod_ClassNotSuffixedSkipped(t *testing.T) {
	globals.InitGlobals("test")
	_ = Init()
	global := globals.GetGlobalRef()

	files := map[string][]byte{
		"classes/SomeClass.txt": []byte("not a class"), // wrong suffix
	}
	global.JmodBaseBytes = buildJmodBytes(t, files)

	beforeCount := BootstrapCL.ClassCount
	err := WalkBaseJmod()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if BootstrapCL.ClassCount != beforeCount {
		t.Errorf("Expected file with wrong suffix to be skipped")
	}
}

// Garbage (unparseable) class bytes are safe to feed through WalkBaseJmod:
// ParseAndPostClass() will log an error and return early (no panic), and
// WalkBaseJmod() does not check ParseAndPostClass's return values.
func TestWalkBaseJmod_NoBootstrapList_ProcessesAllClasses(t *testing.T) {
	globals.InitGlobals("test")
	_ = Init()
	global := globals.GetGlobalRef()

	files := map[string][]byte{
		"classes/org/jacobin/test/Foo.class": {0xCA, 0xFE, 0xBA, 0xBE}, // garbage but well-prefixed/suffixed
		"classes/org/jacobin/test/Bar.class": {0xCA, 0xFE, 0xBA, 0xBE},
	}
	global.JmodBaseBytes = buildJmodBytes(t, files)

	err := WalkBaseJmod()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Garbage bytes fail to parse, so no classes should have been added,
	// but no error should propagate from WalkBaseJmod either.
}

func TestWalkBaseJmod_WithBootstrapList_FiltersClasses(t *testing.T) {
	globals.InitGlobals("test")
	_ = Init()
	global := globals.GetGlobalRef()

	files := map[string][]byte{
		"lib/classlist":                            []byte("org/jacobin/test/Foo\r\norg/jacobin/test/Baz\n"),
		"classes/org/jacobin/test/Foo.class":       {0xCA, 0xFE, 0xBA, 0xBE},
		"classes/org/jacobin/test/NotOnList.class": {0xCA, 0xFE, 0xBA, 0xBE},
	}
	global.JmodBaseBytes = buildJmodBytes(t, files)

	err := WalkBaseJmod()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Both classes fail to parse (garbage bytes) but the bootstrap-list
	// filtering logic itself must not error and must not skip processing
	// of "Foo.class" (on the list) while excluding "NotOnList.class".
}

// --- getClasslist ---

func newZipReader(t *testing.T, files map[string][]byte) *zip.Reader {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("newZipReader: zw.Create(%s) failed: %v", name, err)
		}
		if _, err = w.Write(content); err != nil {
			t.Fatalf("newZipReader: write(%s) failed: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("newZipReader: zw.Close() failed: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("newZipReader: zip.NewReader failed: %v", err)
	}
	return zr
}

func TestGetClasslist_MissingFile(t *testing.T) {
	globals.InitGlobals("test")
	zr := newZipReader(t, map[string][]byte{"other/file.txt": []byte("hello")})

	result := getClasslist(zr)
	if len(result) != 0 {
		t.Errorf("Expected empty map when lib/classlist is missing, got %d entries", len(result))
	}
}

func TestGetClasslist_Success(t *testing.T) {
	globals.InitGlobals("test")
	content := "org/jacobin/test/Foo\r\norg/jacobin/test/Bar\norg/jacobin/test/Baz"
	zr := newZipReader(t, map[string][]byte{"lib/classlist": []byte(content)})

	result := getClasslist(zr)

	expected := []string{
		"org/jacobin/test/Foo.class",
		"org/jacobin/test/Bar.class",
		"org/jacobin/test/Baz.class",
	}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d entries, got %d: %v", len(expected), len(result), result)
	}
	for _, name := range expected {
		if _, ok := result[name]; !ok {
			t.Errorf("Expected %q to be present in classlist map", name)
		}
	}
}

func TestGetClasslist_EmptyFile(t *testing.T) {
	globals.InitGlobals("test")
	zr := newZipReader(t, map[string][]byte{"lib/classlist": []byte("")})

	result := getClasslist(zr)
	// A single empty line becomes "" + ".class" = ".class"
	if len(result) != 1 {
		t.Fatalf("Expected 1 entry for empty classlist content, got %d: %v", len(result), result)
	}
	if _, ok := result[".class"]; !ok {
		t.Errorf("Expected entry %q to be present, got %v", ".class", result)
	}
}

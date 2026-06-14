/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"os"
	"path/filepath"
	"testing"
)

func Test_BasicFileAttributes(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_BasicFileAttributes()
	Load_Nio_File_Attribute_FileTime()

	dir := t.TempDir()
	f := filepath.Join(dir, "test.txt")
	data := []byte("hello world")
	if err := os.WriteFile(f, data, 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}

	info, err := os.Stat(f)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	attrs := newBasicFileAttributes(info)

	// Test isRegularFile
	if res := bfaIsRegularFile([]interface{}{attrs}); res.(int64) != 1 {
		t.Errorf("expected isRegularFile to be true")
	}

	// Test isDirectory
	if res := bfaIsDirectory([]interface{}{attrs}); res.(int64) != 0 {
		t.Errorf("expected isDirectory to be false")
	}

	// Test size
	if res := bfaSize([]interface{}{attrs}); res.(int64) != int64(len(data)) {
		t.Errorf("expected size %d, got %d", len(data), res.(int64))
	}

	// Test lastModifiedTime
	ft := bfaLastModifiedTime([]interface{}{attrs}).(*object.Object)
	milli := fileTimeToMillis([]interface{}{ft})
	if milli.(int64) != info.ModTime().UnixMilli() {
		t.Errorf("expected lastModifiedTime %d, got %d", info.ModTime().UnixMilli(), milli.(int64))
	}

	// Test directory
	dirInfo, _ := os.Stat(dir)
	dirAttrs := newBasicFileAttributes(dirInfo)
	if res := bfaIsDirectory([]interface{}{dirAttrs}); res.(int64) != 1 {
		t.Errorf("expected isDirectory to be true for directory")
	}
	if res := bfaIsRegularFile([]interface{}{dirAttrs}); res.(int64) != 0 {
		t.Errorf("expected isRegularFile to be false for directory")
	}
}

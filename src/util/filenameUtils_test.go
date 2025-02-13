/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"os"
	"strings"
	"testing"
)

func TestConvertInternalClassNameToFilename(t *testing.T) {
	fs := os.PathSeparator

	s := ConvertInternalClassNameToFilename("sponge")
	if s != "sponge.class" {
		t.Errorf("In ConvertInternalClassNameToFilename(), expected 'sponge.class', got: %s", s)
	}

	s = ConvertInternalClassNameToFilename("sponge/bob")
	if fs == '/' {
		if s != "sponge/bob.class" {
			t.Error("From ConvertInternalClassNameToFilename() expected 'sponge/bob.class, got: " + s)
		}
	} else if fs == '\\' {
		if s != "sponge\\bob.class" {
			t.Error("From ConvertInternalClassNameToFilename() expected 'sponge\\bob.class, got: " + s)
		}
	}

	s = ConvertInternalClassNameToFilename("sponge/bob\\square.pants")
	if fs == '/' {
		if s != "sponge/bob/square/pants.class" {
			t.Error("From ConvertInternalClassNameToFilename() expected 'sponge/bob/square/pants.class', got: " + s)
		}
	} else if fs == '\\' {
		if s != "sponge\\bob\\square\\pants.class" {
			t.Error("From ConvertInternalClassNameToFilename() expected 'sponge\\bob\\square\\pants.classs', got: " + s)
		}
	}
}

func TestConvertClassFilenameToInternalFormat(t *testing.T) {
	if ConvertClassFilenameToInternalFormat("sponge") != "sponge" {
		t.Error("Unexpected result in call ConvertClassFilenameToInternalFormat()")
	}

	s := ConvertClassFilenameToInternalFormat("sponge.bob.class")
	if s != "sponge/bob" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}

	s = ConvertClassFilenameToInternalFormat("sponge/bob/square.Pants.class")
	if s != "sponge/bob/square/Pants" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}
}

func TestConvertInternalClassNameToUserFormat(t *testing.T) {
	s := ConvertInternalClassNameToUserFormat("java/lang/Object")
	if s != "java.lang.Object" {
		t.Errorf("Expected 'java.lang.Object', got: %s", s)
	}

	s = ConvertInternalClassNameToUserFormat("com.example.MyClass")
	if s != "com.example.MyClass" {
		t.Errorf("Expected 'com.example.MyClass', got: %s", s)
	}
}

func TestConvertFilenameWithPlatformPathSeparator(t *testing.T) {
	if os.PathSeparator == '\\' {
		s := ConvertToPlatformPathSeparators("snoop/dog/the/man")
		if strings.ContainsRune(s, '/') {
			t.Errorf("Expected a path with no / slashes, got: %s", s)
		}
		if strings.Count(s, string("\\")) != 3 {
			t.Error("Expected 3 backslashes in path, got: ", strings.Count(s, "\\"))
		}
	} else {
		var s string
		if os.PathSeparator == '/' {
			s = ConvertToPlatformPathSeparators("snoop\\dog\\the\\man")
			if strings.ContainsRune(s, '\\') {
				t.Errorf("Expected a path with no \\ slashes, got: %s", s)
			}
		}
		if strings.Count(s, string("/")) != 3 {
			t.Error("Expected 3 forward slashes in path, got: ", strings.Count(s, "/"))
		}
	}
}

// test whether a file is part of the JDK, based on its prefix

func TestIsFilePartOfJDK_JdkPrefix(t *testing.T) {
	filename := "jdk/internal/reflect/Reflection.class"
	if !IsFilePartOfJDK(&filename) {
		t.Errorf("Expected true for filename with 'jdk' prefix, got false")
	}
}

func TestIsFilePartOfJDK_SunPrefix(t *testing.T) {
	filename := "sun/misc/Unsafe.class"
	if !IsFilePartOfJDK(&filename) {
		t.Errorf("Expected true for filename with 'sun' prefix, got false")
	}
}

func TestIsFilePartOfJDK_NoPrefix(t *testing.T) {
	filename := "com/example/MyClass.class"
	if IsFilePartOfJDK(&filename) {
		t.Errorf("Expected false for filename with no JDK prefix, got true")
	}
}

// test searching for a file by extension

func TestSearchDirByFileExtension_FindsFiles(t *testing.T) {
	tempDir := t.TempDir()

	tempFile1, err := os.CreateTemp(tempDir, "*.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = tempFile1.Close()

	tempFile2, err := os.CreateTemp(tempDir, "*.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = tempFile2.Close()

	tempFile3, err := os.CreateTemp(tempDir, "*.log")
	if err != nil {
		t.Fatal(err)
	}
	_ = tempFile3.Close()

	result := SearchDirByFileExtension(tempDir, "txt")
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(*result) != 2 {
		t.Errorf("Expected 2 files, got %d", len(*result))
	}
}

func TestSearchDirByFileExtension_NoFilesFound(t *testing.T) {
	tempDir := t.TempDir()

	result := SearchDirByFileExtension(tempDir, "txt")
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(*result) != 0 {
		t.Errorf("Expected 0 files, got %d", len(*result))
	}
}

func TestSearchDirByFileExtension_HandlesErrors(t *testing.T) {
	tempDir := " no such directory -- 778899 "
	result := SearchDirByFileExtension(tempDir, "txt")
	if result != nil {
		t.Fatal("Expected nil result due to absent directory")
	}
}

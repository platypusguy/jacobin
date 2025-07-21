/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

// test listing jar files in a directory
func TestListJarFiles_ValidDirectoryWithJars(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "test_jars")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	// Create dummy files
	jar1Path := filepath.Join(tmpDir, "app.jar")
	jar2Path := filepath.Join(tmpDir, "lib.JAR")
	txtPath := filepath.Join(tmpDir, "notes.txt")
	subDirPath := filepath.Join(tmpDir, "subdir")
	subDirJarPath := filepath.Join(subDirPath, "nested.jar")

	// Create files and a subdirectory
	if err := ioutil.WriteFile(jar1Path, []byte("jar content"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", jar1Path, err)
	}
	if err := ioutil.WriteFile(jar2Path, []byte("jar content"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", jar2Path, err)
	}
	if err := ioutil.WriteFile(txtPath, []byte("text content"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", txtPath, err)
	}
	if err := os.Mkdir(subDirPath, 0755); err != nil {
		t.Fatalf("Failed to create subdir %s: %v", subDirPath, err)
	}
	if err := ioutil.WriteFile(subDirJarPath, []byte("nested jar content"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", subDirJarPath, err)
	}

	expectedJars := []string{jar1Path, jar2Path, subDirJarPath}

	foundJars, err := ListJarFiles(tmpDir)
	if err != nil {
		t.Errorf("ListJarFiles returned an unexpected error: %v", err)
	}

	if len(foundJars) != len(expectedJars) {
		t.Errorf("Expected %d jar files, got %d. Found: %v, Expected: %v",
			len(expectedJars), len(foundJars), foundJars, expectedJars)
	}

	// Check if all expected jars are found, regardless of order
	foundMap := make(map[string]bool)
	for _, jar := range foundJars {
		foundMap[jar] = true
	}

	for _, expectedJar := range expectedJars {
		if _, ok := foundMap[expectedJar]; !ok {
			t.Errorf("Expected jar file %s not found in results: %v", expectedJar, foundJars)
		}
	}
}

func TestListJarFiles_EmptyDirectory(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	foundJars, err := ListJarFiles(tmpDir)
	if err != nil {
		t.Errorf("ListJarFiles returned an unexpected error for empty directory: %v", err)
	}
	if len(foundJars) != 0 {
		t.Errorf("Expected 0 jar files for empty directory, got %d: %v", len(foundJars), foundJars)
	}
}

func TestListJarFiles_NoJarsInDirectory(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test_no_jars")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := ioutil.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(tmpDir, "document.pdf"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	foundJars, err := ListJarFiles(tmpDir)
	if err != nil {
		t.Errorf("ListJarFiles returned an unexpected error for directory with no jars: %v", err)
	}
	if len(foundJars) != 0 {
		t.Errorf("Expected 0 jar files, got %d: %v", len(foundJars), foundJars)
	}
}

func TestListJarFiles_NonExistentDirectory(t *testing.T) {
	nonExistentDir := filepath.Join(os.TempDir(), "non_existent_dir_12345") // Use a unique name
	// Ensure it doesn't exist before the test
	os.RemoveAll(nonExistentDir)

	_, err := ListJarFiles(nonExistentDir)
	if err == nil {
		t.Errorf("Expected an error for non-existent directory, but got none")
	}
	// Check if the error message indicates a "no such file or directory" error
	if !os.IsNotExist(err) && !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Expected 'no such file or directory' error, but got: %v", err)
	}
}

func TestListJarFiles_OnlySubdirectories(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test_only_subdirs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.Mkdir(filepath.Join(tmpDir, "subdir1"), 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tmpDir, "subdir2"), 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	foundJars, err := ListJarFiles(tmpDir)
	if err != nil {
		t.Errorf("ListJarFiles returned an unexpected error: %v", err)
	}
	if len(foundJars) != 0 {
		t.Errorf("Expected 0 jar files, got %d: %v", len(foundJars), foundJars)
	}
}

func TestListJarFiles_MixedContent(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test_mixed_content")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files and directories
	jarRoot := filepath.Join(tmpDir, "root.jar")
	txtRoot := filepath.Join(tmpDir, "root.txt")
	dir1 := filepath.Join(tmpDir, "dir1")
	jarDir1 := filepath.Join(dir1, "dir1_app.jar")
	dir2 := filepath.Join(tmpDir, "dir2")
	txtDir2 := filepath.Join(dir2, "dir2_file.txt")
	dir3 := filepath.Join(tmpDir, "dir3") // Empty directory

	if err := ioutil.WriteFile(jarRoot, []byte("root jar"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", jarRoot, err)
	}
	if err := ioutil.WriteFile(txtRoot, []byte("root txt"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", txtRoot, err)
	}
	if err := os.Mkdir(dir1, 0755); err != nil {
		t.Fatalf("Failed to create %s: %v", dir1, err)
	}
	if err := ioutil.WriteFile(jarDir1, []byte("dir1 jar"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", jarDir1, err)
	}
	if err := os.Mkdir(dir2, 0755); err != nil {
		t.Fatalf("Failed to create %s: %v", dir2, err)
	}
	if err := ioutil.WriteFile(txtDir2, []byte("dir2 txt"), 0644); err != nil {
		t.Fatalf("Failed to create %s: %v", txtDir2, err)
	}
	if err := os.Mkdir(dir3, 0755); err != nil {
		t.Fatalf("Failed to create %s: %v", dir3, err)
	}

	expectedJars := []string{jarRoot, jarDir1} // The order might vary, so we'll check content
	foundJars, err := ListJarFiles(tmpDir)
	if err != nil {
		t.Errorf("ListJarFiles returned an unexpected error: %v", err)
	}

	if len(foundJars) != len(expectedJars) {
		t.Errorf("Expected %d jar files, got %d. Found: %v, Expected: %v",
			len(expectedJars), len(foundJars), foundJars, expectedJars)
	}

	foundMap := make(map[string]bool)
	for _, jar := range foundJars {
		foundMap[jar] = true
	}

	for _, expectedJar := range expectedJars {
		if _, ok := foundMap[expectedJar]; !ok {
			t.Errorf("Expected jar file %s not found in results: %v", expectedJar, foundJars)
		}
	}
}

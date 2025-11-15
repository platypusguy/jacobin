/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-3 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"io"
	"jacobin/src/globals"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// Most of the functionality in classloader package is tested in other files, such as
// * cpParser_test.go (constant pool parser)
// * formatCheck_test.go (the format checking)
// * parser_test.go (the class parsing)
// etc.
// This files tests remaining routines.

func TestInitOfClassloaders(t *testing.T) {
	globals.InitGlobals("test")
	// set the logger to low granularity, so that logging messages are not also captured in this test

	_ = Init()

	// check that the classloader hierarchy is set up correctly
	if BootstrapCL.Parent != "" {
		t.Errorf("Expecting parent of Boostrap classloader to be empty, got: %s",
			BootstrapCL.Parent)
	}

	if ExtensionCL.Parent != "bootstrap" {
		t.Errorf("Expecting parent of Extension classloader to be Boostrap, got: %s",
			ExtensionCL.Parent)
	}

	if AppCL.Parent != "extension" {
		t.Errorf("Expecting parent of Application classloader to be Extension, got: %s",
			AppCL.Parent)
	}

	// check that the classloaders have empty tables ready
	if BootstrapCL.ClassCount != 0 {
		t.Errorf("Expected size of bootstrap CL's table to be 0, got: %d",
			BootstrapCL.ClassCount)
	}

	if ExtensionCL.ClassCount != 0 {
		t.Errorf("Expected size of extension CL's table to be 0, got: %d",
			ExtensionCL.ClassCount)
	}

	if AppCL.ClassCount != 0 {
		t.Errorf("Expected size of application CL's table to be 0, got: %d",
			AppCL.ClassCount)
	}
}

func TestWalkWithError(t *testing.T) {
	e := errors.New("test error")
	err := walk("", nil, e)
	if err != e {
		t.Errorf("Expected an error = to 'test error', got %s",
			err.Error())
	}
}

// when walk() encounters an invalid file, it is simply skipped
// with no error generated as it's not clear that entry in jmod
// will be necessary. If it is, when it's invoked, it will be loaded
// then and any errors in finding the file will be returned then.
func TestJmodWalkWithInvalidDirAndFile(t *testing.T) {
	err := os.Mkdir("subdir", 0755)
	defer os.RemoveAll("subdir")
	_ = os.WriteFile("subdir/file1", []byte(""), 0644)

	dirEntry, err := os.ReadDir("subdir")
	err = walk("gherkin", dirEntry[0], nil)
	if err != nil {
		t.Errorf("Expected no error on invalid file in walk(), but got %s",
			err.Error())
	}
}

func TestLoadClassFromFileInvalidName(t *testing.T) {
	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	nameIndex, _, err := LoadClassFromFile(Classloader{}, "noSuchFile")

	if nameIndex != types.InvalidStringIndex {
		t.Errorf("Expected empty filename due to error, got: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Expected an error message for invalid file name, but got none")
	}

	_ = w.Close()
	_, _ = io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

// remove leading [L and delete trailing;, eliminate all other entries with [prefix
func TestNormalizingClassReference(t *testing.T) {
	s := normalizeClassReference("[Ljava/test/java.String;")
	if s != "java/test/java.String" {
		t.Error("Unexpected normalized class reference: " + s)
	}

	s = normalizeClassReference(types.ByteArray)
	if s != "" {
		t.Error("Unexpected normalized class reference: " + s)
	}

	s = normalizeClassReference(types.ObjectClassName)
	if s != types.ObjectClassName {
		t.Error("Unexpected normalized class reference: " + s)
	}
}

func TestConvertToPostableClassStringRefs(t *testing.T) {
	// Testing the changes made as a result of JACOBIN-103
	globals.InitGlobals("test")
	trace.Init()

	// set up a class with a constant pool containing the one
	// StringConst we want to make sure is converted to a UTF8
	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{StringConst, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})

	klass.stringRefs = append(klass.stringRefs, stringConstantEntry{index: 0})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{content: "Hello string"})

	klass.cpCount = 3

	postableClass := convertToPostableClass(&klass)
	if len(postableClass.CP.Utf8Refs) != 1 {
		t.Errorf("Expecting a UTF8 slice of length 1, got %d",
			len(postableClass.CP.Utf8Refs))
	}

	// cpIndex[1] is a StringConst above, should now be a UTF8
	utf8 := postableClass.CP.CpIndex[1]
	if utf8.Type != UTF8 {
		t.Errorf("Expecting StringConst entry to have become UTF8 entry,"+
			"but instead is of type: %d", utf8.Type)
	}
}

func TestGetInvalidJarName(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	_, err := getArchiveFile(BootstrapCL, "")
	if err == nil {
		t.Errorf("expected err msg for fetching an invalid JAR, but got err=nil")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "ERROR") {
		t.Error("Got unexpected error msg: " + msg)
	}
}

func TestGetClassFromInvalidJar(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	_, _, err := LoadClassFromArchive(BootstrapCL, "pickle", "gherkin")
	if err == nil {
		t.Errorf("expected err msg for loading invalid class from invalid JAR, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "ERROR") {
		t.Error("Got unexpected error msg: " + msg)
	}
}

func TestMainClassFromInvalidJar(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	_, _, err := GetMainClassFromJar(BootstrapCL, "gherkin")
	if err == nil {
		t.Errorf("expected err msg for loading main class from invalid JAR, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "ERROR") {
		t.Error("Got unexpected error msg: " + msg)
	}
}

func TestInvalidMagicNumberViaParseAndPostFunction(t *testing.T) {

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("test")
	trace.Init()

	err := Init()

	testBytes := []byte{
		0xCB, 0xFE, 0xBA, 0xBE,
	}

	_, _, err = ParseAndPostClass(&BootstrapCL, "Hello2", testBytes)
	if err == nil {
		t.Error("Expected an error, but got none.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(string(msg), "invalid magic number") {
		t.Errorf("Expected error message to contain in part 'invalid magic number', got: %s", string(msg))
	}
}

var Hello2Bytes = []byte{
	0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x37, 0x00, 0x2B, 0x07, 0x00, 0x02, 0x01, 0x00, 0x06,
	0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x32, 0x07, 0x00, 0x04, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61,
	0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x01, 0x00, 0x06, 0x3C,
	0x69, 0x6E, 0x69, 0x74, 0x3E, 0x01, 0x00, 0x03, 0x28, 0x29, 0x56, 0x01, 0x00, 0x04, 0x43, 0x6F,
	0x64, 0x65, 0x0A, 0x00, 0x03, 0x00, 0x09, 0x0C, 0x00, 0x05, 0x00, 0x06, 0x01, 0x00, 0x0F, 0x4C,
	0x69, 0x6E, 0x65, 0x4E, 0x75, 0x6D, 0x62, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x01, 0x00,
	0x12, 0x4C, 0x6F, 0x63, 0x61, 0x6C, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6C, 0x65, 0x54, 0x61,
	0x62, 0x6C, 0x65, 0x01, 0x00, 0x04, 0x74, 0x68, 0x69, 0x73, 0x01, 0x00, 0x08, 0x4C, 0x48, 0x65,
	0x6C, 0x6C, 0x6F, 0x32, 0x3B, 0x01, 0x00, 0x04, 0x6D, 0x61, 0x69, 0x6E, 0x01, 0x00, 0x16, 0x28,
	0x5B, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69,
	0x6E, 0x67, 0x3B, 0x29, 0x56, 0x0A, 0x00, 0x01, 0x00, 0x11, 0x0C, 0x00, 0x12, 0x00, 0x13, 0x01,
	0x00, 0x06, 0x61, 0x64, 0x64, 0x54, 0x77, 0x6F, 0x01, 0x00, 0x05, 0x28, 0x49, 0x49, 0x29, 0x49,
	0x09, 0x00, 0x15, 0x00, 0x17, 0x07, 0x00, 0x16, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61, 0x2F,
	0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6D, 0x0C, 0x00, 0x18, 0x00, 0x19,
	0x01, 0x00, 0x03, 0x6F, 0x75, 0x74, 0x01, 0x00, 0x15, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69,
	0x6F, 0x2F, 0x50, 0x72, 0x69, 0x6E, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x3B, 0x0A, 0x00,
	0x1B, 0x00, 0x1D, 0x07, 0x00, 0x1C, 0x01, 0x00, 0x13, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F,
	0x2F, 0x50, 0x72, 0x69, 0x6E, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x0C, 0x00, 0x1E, 0x00,
	0x1F, 0x01, 0x00, 0x07, 0x70, 0x72, 0x69, 0x6E, 0x74, 0x6C, 0x6E, 0x01, 0x00, 0x04, 0x28, 0x49,
	0x29, 0x56, 0x01, 0x00, 0x04, 0x61, 0x72, 0x67, 0x73, 0x01, 0x00, 0x13, 0x5B, 0x4C, 0x6A, 0x61,
	0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x3B, 0x01,
	0x00, 0x01, 0x78, 0x01, 0x00, 0x01, 0x49, 0x01, 0x00, 0x01, 0x69, 0x01, 0x00, 0x0D, 0x53, 0x74,
	0x61, 0x63, 0x6B, 0x4D, 0x61, 0x70, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x07, 0x00, 0x21, 0x01, 0x00,
	0x01, 0x6A, 0x01, 0x00, 0x01, 0x6B, 0x01, 0x00, 0x0A, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x46,
	0x69, 0x6C, 0x65, 0x01, 0x00, 0x0B, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x32, 0x2E, 0x6A, 0x61, 0x76,
	0x61, 0x00, 0x20, 0x00, 0x01, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
	0x05, 0x00, 0x06, 0x00, 0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x2F, 0x00, 0x01, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x05, 0x2A, 0xB7, 0x00, 0x08, 0xB1, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0A, 0x00, 0x00,
	0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x05, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x09, 0x00, 0x0E, 0x00, 0x0F,
	0x00, 0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x81, 0x00, 0x03, 0x00, 0x03, 0x00, 0x00, 0x00, 0x1E,
	0x03, 0x3D, 0xA7, 0x00, 0x15, 0x1C, 0x1C, 0x04, 0x64, 0xB8, 0x00, 0x10, 0x3C, 0xB2, 0x00, 0x14,
	0x1B, 0xB6, 0x00, 0x1A, 0x84, 0x02, 0x01, 0x1C, 0x10, 0x0A, 0xA1, 0xFF, 0xEB, 0xB1, 0x00, 0x00,
	0x00, 0x03, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x16, 0x00, 0x05, 0x00, 0x00, 0x00, 0x06, 0x00, 0x05,
	0x00, 0x07, 0x00, 0x0D, 0x00, 0x08, 0x00, 0x14, 0x00, 0x06, 0x00, 0x1D, 0x00, 0x0A, 0x00, 0x0B,
	0x00, 0x00, 0x00, 0x20, 0x00, 0x03, 0x00, 0x00, 0x00, 0x1E, 0x00, 0x20, 0x00, 0x21, 0x00, 0x00,
	0x00, 0x0D, 0x00, 0x0A, 0x00, 0x22, 0x00, 0x23, 0x00, 0x01, 0x00, 0x02, 0x00, 0x1B, 0x00, 0x24,
	0x00, 0x23, 0x00, 0x02, 0x00, 0x25, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x02, 0xFF, 0x00, 0x05, 0x00,
	0x03, 0x07, 0x00, 0x26, 0x00, 0x01, 0x00, 0x00, 0x11, 0x00, 0x08, 0x00, 0x12, 0x00, 0x13, 0x00,
	0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x38, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x1A,
	0x1B, 0x60, 0xAC, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x0D, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x16, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x00,
	0x27, 0x00, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x28, 0x00, 0x23, 0x00, 0x01, 0x00,
	0x01, 0x00, 0x29, 0x00, 0x00, 0x00, 0x02, 0x00, 0x2A,
}

func TestLoadFullyParsedClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	fullyParsedClass, err := parse(Hello2Bytes)
	if err != nil {
		t.Errorf("Got unexpected error from parse of Hello2.class: %s", err.Error())
	}
	classToPost := convertToPostableClass(&fullyParsedClass)
	if len(classToPost.MethodTable) < 1 {
		t.Errorf("Invalid number of methods in Hello2.class: %d", len(classToPost.MethodTable))
	}
}

// === the following tests were generated by Jetbrains Junie to fill in test coverage gaps ===

// Helper to reset global structures that tests rely on and keep isolation between tests.
func resetClassloaderState() {
	MethArea = &sync.Map{}
	// Reset classloader counts
	AppCL.ClassCount = 0
	BootstrapCL.ClassCount = 0
	ExtensionCL.ClassCount = 0
}

func TestLoadClassFromNameOnly_EmptyClassName(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	err := LoadClassFromNameOnly("")
	if err == nil {
		t.Fatalf("expected error for empty class name, got nil")
	}
	if !strings.Contains(err.Error(), "null class name is invalid") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestLoadClassFromNameOnly_ClassNameWithSemicolon(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	// Capture stderr for error logging
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := LoadClassFromNameOnly("com/example/Test;")
	_ = w.Close()
	_, _ = io.ReadAll(r)
	os.Stderr = normalStderr

	if err == nil {
		t.Fatalf("expected error for class name with semicolon, got nil")
	}
	if !strings.Contains(err.Error(), "invalid class name") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestLoadClassFromFile_AbsolutePath(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	// Create a temporary file with valid class structure
	tempFile, err := os.CreateTemp("", "TestClass*.class")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the Hello2 class bytes (which are valid and already used in existing tests)
	_, err = tempFile.Write(Hello2Bytes)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Test with absolute path - this should work since we're using real class bytes
	_, _, err = LoadClassFromFile(AppCL, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error loading absolute path: %v", err)
	}
}

func TestLoadClassFromFile_ClasspathIteration(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	// Create temporary directories
	tempDir1, err := os.MkdirTemp("", "classpath1")
	if err != nil {
		t.Fatalf("failed to create temp dir 1: %v", err)
	}
	defer os.RemoveAll(tempDir1)

	tempDir2, err := os.MkdirTemp("", "classpath2")
	if err != nil {
		t.Fatalf("failed to create temp dir 2: %v", err)
	}
	defer os.RemoveAll(tempDir2)

	// Create a class file in the second directory using Hello2 bytes
	classFile := filepath.Join(tempDir2, "Hello2.class")
	err = os.WriteFile(classFile, Hello2Bytes, 0644)
	if err != nil {
		t.Fatalf("failed to write class file: %v", err)
	}

	// Set classpath with both directories
	originalClasspath := globals.GetGlobalRef().Classpath
	defer func() { globals.GetGlobalRef().Classpath = originalClasspath }()
	globals.GetGlobalRef().Classpath = []string{tempDir1, tempDir2}

	// Test loading - should find file in second directory
	_, _, err = LoadClassFromFile(AppCL, "Hello2")
	if err != nil {
		t.Fatalf("unexpected error loading from classpath: %v", err)
	}
}

func TestLoadClassFromFile_FileNotFound(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	// Set a non-existent classpath
	originalClasspath := globals.GetGlobalRef().Classpath
	defer func() { globals.GetGlobalRef().Classpath = originalClasspath }()
	globals.GetGlobalRef().Classpath = []string{"/nonexistent/path"}

	// Capture stderr for error logging
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	_, _, err := LoadClassFromFile(AppCL, "NonExistentClass")
	_ = w.Close()
	_, _ = io.ReadAll(r)
	os.Stderr = normalStderr

	if err == nil {
		t.Fatalf("expected error for non-existent class file, got nil")
	}
}

func TestLoadFromLoaderChannel_ClassAlreadyLoaded(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	// Insert a class into MethArea
	testClass := &Klass{Status: 'F', Loader: "test", Data: &ClData{Name: "TestClass"}}
	MethAreaInsert("TestClass", testClass)

	// Create a channel and send the class name
	channel := make(chan string, 1)
	channel <- "TestClass"
	close(channel)

	// Set up WaitGroup
	globals.LoaderWg.Add(1)

	// This should skip the already loaded class and call Done()
	LoadFromLoaderChannel(channel)

	// Verify the class is still there and unchanged
	retrieved := MethAreaFetch("TestClass")
	if retrieved == nil {
		t.Fatalf("expected class to still be in MethArea")
	}
	if retrieved.Loader != "test" {
		t.Fatalf("expected class loader to be unchanged, got: %s", retrieved.Loader)
	}
}

func TestGetCountOfLoadedClasses(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	// Test initial count
	count := AppCL.GetCountOfLoadedClasses()
	if count != 0 {
		t.Fatalf("expected initial count 0, got %d", count)
	}

	// Increment count manually
	AppCL.ClassCount = 5

	count = AppCL.GetCountOfLoadedClasses()
	if count != 5 {
		t.Fatalf("expected count 5, got %d", count)
	}
}

func TestCFE_ErrorFormatting(t *testing.T) {
	globals.InitGlobals("test")

	err := CFE("test error message")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Class Format Error: test error message") {
		t.Fatalf("unexpected error format: %s", errMsg)
	}
	if !strings.Contains(errMsg, "detected by file:") {
		t.Fatalf("error should contain file location: %s", errMsg)
	}
}

func TestCfe_ErrorFormatting(t *testing.T) {
	globals.InitGlobals("test")

	err := cfe("another test error")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Class Format Error: another test error") {
		t.Fatalf("unexpected error format: %s", errMsg)
	}
	if !strings.Contains(errMsg, "detected by file:") {
		t.Fatalf("error should contain file location: %s", errMsg)
	}
}

func TestNormalizeClassReference_Array(t *testing.T) {
	// Test array reference normalization
	result := normalizeClassReference(types.ByteArray)
	if result != "" {
		t.Fatalf("expected empty string for byte array, got: %s", result)
	}
}

func TestNormalizeClassReference_RefArray(t *testing.T) {
	// Test reference array normalization
	result := normalizeClassReference("[Ljava/lang/String;")
	if result != "java/lang/String" {
		t.Fatalf("expected 'java/lang/String', got: %s", result)
	}
}

func TestLoadClassFromNameOnly_SuperclassRecursion(t *testing.T) {
	globals.InitGlobals("test")
	resetClassloaderState()

	Init()
	LoadBaseClasses()

	// Create temporary directories for classpath
	tempDir, err := os.MkdirTemp("", "classfiles")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up classpath
	originalClasspath := globals.GetGlobalRef().Classpath
	defer func() { globals.GetGlobalRef().Classpath = originalClasspath }()
	globals.GetGlobalRef().Classpath = []string{tempDir}

	// Write Hello2 class file (which has java/lang/Object as superclass)
	classFile := filepath.Join(tempDir, "Hello2.class")
	err = os.WriteFile(classFile, Hello2Bytes, 0644)
	if err != nil {
		t.Fatalf("failed to write class file: %v", err)
	}

	// Ensure java/lang/Object is available (it should be from InitMethodArea)
	err = LoadClassFromNameOnly("Hello2")
	if err != nil {
		t.Fatalf("unexpected error in class loading: %v", err)
	}

	// Verify both Hello2 and java/lang/Object are loaded
	hello2Class := MethAreaFetch("Hello2")
	if hello2Class == nil {
		t.Fatalf("Hello2 class should be loaded")
	}

	objectClass := MethAreaFetch(types.ObjectClassName)
	if objectClass == nil {
		t.Fatalf("java/lang/Object should be available")
	}
}

// Ensure getArchiveFile caches the opened archive within the classloader
func TestGetArchiveFile_CachesArchive(t *testing.T) {
    // Build a simple jar with a manifest and a single class entry
    jarPath, cleanup := makeTempJar(t, map[string]string{"Main-Class": "com.example.Main"}, map[string][]byte{
        "com/example/Main.class": {0xCA, 0xFE, 0xBA, 0xBE},
    })
    defer cleanup()

    // Create a fresh classloader instance
    cl := Classloader{Name: "test", Parent: "", Archives: make(map[string]*Archive)}

    // First fetch should open and cache
    a1, err := getArchiveFile(cl, jarPath)
    if err != nil {
        t.Fatalf("first getArchiveFile failed: %v", err)
    }
    if a1 == nil {
        t.Fatalf("expected non-nil archive")
    }

    // Second fetch should retrieve the same pointer from cache
    a2, err := getArchiveFile(cl, jarPath)
    if err != nil {
        t.Fatalf("second getArchiveFile failed: %v", err)
    }
    if a1 != a2 {
        t.Fatalf("expected cached archive pointer to be reused")
    }

    // Cache should have exactly one entry
    if len(cl.Archives) != 1 {
        t.Fatalf("expected exactly 1 cached archive, got %d", len(cl.Archives))
    }
}

// GetMainClassFromJar should not error when the manifest lacks Main-Class; it should return ""
func TestGetMainClassFromJar_NoMainClass_NoError(t *testing.T) {
    // Jar with a manifest but without Main-Class
    jarPath, cleanup := makeTempJar(t, map[string]string{"Class-Path": "lib/a.jar lib/b.jar"}, map[string][]byte{})
    defer cleanup()

    cl := Classloader{Name: "test", Parent: "", Archives: make(map[string]*Archive)}
    mainClass, archive, err := GetMainClassFromJar(cl, jarPath)
    if err != nil {
        t.Fatalf("GetMainClassFromJar returned error for jar without Main-Class: %v", err)
    }
    if archive == nil {
        t.Fatalf("expected non-nil archive")
    }
    if mainClass != "" {
        t.Fatalf("expected empty main class, got %q", mainClass)
    }
}

// LoadClassFromArchive should successfully parse and post a class from a jar
func TestLoadClassFromArchive_Success(t *testing.T) {
    globals.InitGlobals("test")
    trace.Init()
    _ = Init()
    LoadBaseClasses()

    // Create a jar that contains Hello2.class at the root
    jarPath, cleanup := makeTempJar(t, map[string]string{"Main-Class": "Hello2"}, map[string][]byte{
        "Hello2.class": Hello2Bytes,
    })
    defer cleanup()

    // Use a fresh classloader instance so archive cache starts empty
    cl := Classloader{Name: "bootstrap", Archives: make(map[string]*Archive)}

    nameIdx, superIdx, err := LoadClassFromArchive(cl, "Hello2", jarPath)
    if err != nil {
        t.Fatalf("LoadClassFromArchive failed: %v", err)
    }
    if nameIdx == types.InvalidStringIndex || superIdx == types.InvalidStringIndex {
        t.Fatalf("unexpected invalid indices: name=%d super=%d", nameIdx, superIdx)
    }

    // Verify class is present in method area
    if kc := MethAreaFetch("Hello2"); kc == nil {
        t.Fatalf("expected Hello2 to be posted to method area")
    }
}

// Additional normalization edge cases beyond existing tests
func TestNormalizeClassReference_MultiDimensionalAndMalformed(t *testing.T) {
    // Multi-dimensional reference array -> current behavior skips arrays (returns empty)
    got := normalizeClassReference("[[[Ljava/util/List;")
    if got != "" {
        t.Fatalf("expected empty for multi-dimensional ref arrays, got %q", got)
    }

    // Multi-dimensional primitive array -> should be skipped (empty)
    got = normalizeClassReference("[[I")
    if got != "" {
        t.Fatalf("expected empty for primitive arrays, got %q", got)
    }

    // Malformed reference array missing trailing ';' -> returns the remainder as-is
    got = normalizeClassReference("[Lbad/Ref")
    if got != "bad/Ref" {
        t.Fatalf("expected 'bad/Ref' for malformed ref array, got %q", got)
    }
}

// === end of tests generated by Jetbrains Junie ===

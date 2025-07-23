/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	// "fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	g := InitGlobals("test")

	if g.JacobinName != "test" {
		t.Errorf("Expecting globals init to set Jacobin name to 'testInit', got: %s", g.JacobinName)
	}

	if g.VmModel != "server" {
		t.Errorf("Expected globals init to set VmModel to 'server', got: %s", g.VmModel)
	}
}

func TestInitClasspath(t *testing.T) {
	g := InitGlobals("test")
	pwd, _ := os.Getwd()
	if g.ClasspathRaw != pwd {
		t.Errorf("Expected ClassPath to be set to current working directory, got: %s", g.ClasspathRaw)
	}

	if g.Classpath[0] != pwd {
		t.Errorf("Expected ClassPath[0] to be set to current working directory, got: %s", g.Classpath[0])
	}
}

func TestInitClasspathWithEnv(t *testing.T) {
	origClasspath := os.Getenv("CLASSPATH")
	_ = os.Setenv("CLASSPATH", "home")
	defer os.Setenv("CLASSPATH", origClasspath)

	g := InitGlobals("test")
	if g.ClasspathRaw != "home" {
		t.Errorf("Expected ClassPath to be set to 'home', got: %s", g.ClasspathRaw)
	}
}

func TestJavaHomeInvalidDir(t *testing.T) {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	if !nameFooBar(t) {
		return
	}
	defer os.RemoveAll(foo)

	origJavaHome := os.Getenv("JAVA_HOME")
	_ = os.Setenv("JAVA_HOME", foobar)
	InitJavaHome()

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "The system cannot find the path specified.") {
		t.Errorf("Expecting error message containing 'The system cannot find the path specified.', got: %s",
			errMsg)
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
	// fmt.Printf("TestJavaHomeAndVersion: JAVA_HOME=%s, JAVA_VERSION=%s\n", home, version)
}

// verify that there is no trailing slash in JAVA_HOME
func TestJavaHomeRemovalOfTrailingSlash(t *testing.T) {
	InitJavaHome()
	ret := JavaHome()
	if strings.HasSuffix(ret, string(os.PathSeparator)) {
		t.Errorf("JAVA_HOME should not have a trailing slash, got: %s", ret)
	}
}

func TestJavaHomeNotSet(t *testing.T) {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Save the original JAVA_HOME environment variable
	origJavaHome := os.Getenv("JAVA_HOME")
	// Unset JAVA_HOME
	_ = os.Unsetenv("JAVA_HOME")

	// Initialize JavaHome
	InitJavaHome()

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "Environment variable JAVA_HOME missing but is required. Exiting.") {
		t.Errorf("Expected error message 'JAVA_HOME environment variable is not set.', got: %s", errMsg)
	}

	// Restore the original JAVA_HOME environment variable
	_ = os.Setenv("JAVA_HOME", origJavaHome)
}

func TestInitJavaHomeWithInvalidReleasFile(t *testing.T) {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Create a temporary directory
	tempDir := t.TempDir()

	// Define the path for the release file
	releaseFilePath := filepath.Join(tempDir, "release")

	// Write "nonsense" to the release file
	err := os.WriteFile(releaseFilePath, []byte("nonsense\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create release file: %v", err)
	}

	// Set JAVA_HOME to the temporary directory
	prevJavaHomeEnv := os.Getenv("JAVA_HOME")
	_ = os.Setenv("JAVA_HOME", tempDir)

	// Initialize JavaHome
	InitJavaHome()

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "InitJavaHome: Did not find the JAVA_VERSION record") {
		t.Errorf("Expected error message 'InitJavaHome: Did not find the JAVA_VERSION record', got: %s", errMsg)
	}
	defer os.Setenv("JAVA_HOME", prevJavaHomeEnv)
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
	InitGlobals("test")
	gl := GetGlobalRef()
	if gl.StrictJDK != false ||
		gl.ExitNow != false ||
		!(gl.MaxJavaVersion >= 11) {
		t.Errorf("Some global variables intialized to unexpected values.")
	}
}

func TestGetSystemProperty(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("java.version")
	if ret < "21" {
		t.Errorf("Expecting a java.version of 21 or more, got: %s", ret)
	}
}

func TestGetSystemPropertyNotFound(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("java.version.notfound")
	if ret != "" {
		t.Errorf("Expecting a java.version.notfound of '', got: %s", ret)
	}
}

func TestGetSystemClasspath(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("java.class.path")
	if ret != "." {
		t.Errorf("Expecting a java.class.path of ., got: %s", ret)
	}
}

func TestGetFileEncoding(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("file.encoding")
	if ret == "UTF-8" {
		t.Errorf("Expecting a file.encoding of UTF-8, got: %s", ret)
	}
}

func TestGetFileNameEncoding(t *testing.T) {
	InitGlobals("test")
	buildGlobalProperties()
	ret := GetSystemProperty("sun.jnu.encoding")
	if ret == "UTF-8" {
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

	if !strings.Contains(errMsg, "The system cannot find the path specified.") {
		t.Errorf("Expecting error message 'The system cannot find the path specified.', got: %s", errMsg)
	}
	_ = os.Setenv("JAVA_HOME", prevJavaHomeEnv)
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

func TestGetJDKversionSuccess(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_java_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a release file with JAVA_VERSION
	releaseFilePath := tempDir + string(os.PathSeparator) + "release"
	err = os.WriteFile(releaseFilePath, []byte("JAVA_VERSION=\"17.0.1\"\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create release file: %v", err)
	}

	// Set JAVA_HOME to the temporary directory
	os.Setenv("JAVA_HOME", tempDir)
	defer os.Unsetenv("JAVA_HOME")
	gl := InitGlobals("test")
	gl.JavaHome = tempDir // Ensure the global variable is set to the temp directory

	// Call GetJDKversion
	_, version := GetJDKmajorVersion()
	if version != "17.0.1" {
		t.Errorf("Expected JAVA_VERSION '17.0.1', got '%s'", version)
	}
}

func TestGetJDKversionFileNotFound(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_java_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set JAVA_HOME to the temporary directory
	_ = os.Setenv("JAVA_HOME", tempDir)
	defer os.Unsetenv("JAVA_HOME")

	// Call GetJDKversion
	_, version := GetJDKmajorVersion()
	if version != "" {
		t.Errorf("Expected empty JAVA_VERSION, got '%s'", version)
	}
}

func TestGetJDKversionMalformedFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_java_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a malformed release file
	releaseFilePath := tempDir + string(os.PathSeparator) + "release"
	err = os.WriteFile(releaseFilePath, []byte("MALFORMED_LINE\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create release file: %v", err)
	}

	// Set JAVA_HOME to the temporary directory
	_ = os.Setenv("JAVA_HOME", tempDir)
	defer os.Unsetenv("JAVA_HOME")

	// Call GetJDKversion
	_, version := GetJDKmajorVersion()
	if version != "" {
		t.Errorf("Expected empty JAVA_VERSION for malformed file, got '%s'", version)
	}
}

func TestGetJDKversionScannerError(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_java_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a release file with a large line to simulate scanner error
	releaseFilePath := tempDir + string(os.PathSeparator) + "release"
	largeLine := make([]byte, 1+64*1024) // 64 KB line
	err = os.WriteFile(releaseFilePath, largeLine, 0644)
	if err != nil {
		t.Fatalf("Failed to create release file: %v", err)
	}

	// Set JAVA_HOME to the temporary directory
	os.Setenv("JAVA_HOME", tempDir)
	defer os.Unsetenv("JAVA_HOME")

	// get JDK version
	_, version := GetJDKmajorVersion()
	if version != "" {
		t.Errorf("Expected empty JAVA_VERSION for scanner error, got '%s'", version)
	}
}

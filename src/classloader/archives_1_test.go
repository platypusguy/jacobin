/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
    "archive/zip"
    "io"
	"os"
	"path/filepath"
	"testing"
)

var GOOD_JAR_NAME = "hello.jar"
var NO_MANIFEST_JAR_NAME = "nomanifest.jar"

func getJarFileName(name string) (string, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(pwd, "..", "..", "testdata", name), nil
}

func getJar(name string, t *testing.T) (*Archive, error) {
    fileName, err := getJarFileName(name)

	if err != nil {
		t.Error("Unable to get jar file", err)
		return nil, err
	}

	return OpenArchive(fileName)
}

// test helper: create a temporary JAR file with provided entries
// - manifest: map of key->value to include in MANIFEST.MF (nil for none)
// - files: map[path]content for arbitrary entries
// Returns full path to the jar and a cleanup func.
func makeTempJar1(t *testing.T, manifest map[string]string, files map[string][]byte) (string, func()) {
    t.Helper()

    dir := t.TempDir()
    jarPath := filepath.Join(dir, "temp1.jar")
    f, err := os.Create(jarPath)
    if err != nil {
        t.Fatalf("failed creating temp jar: %v", err)
    }
    zw := zip.NewWriter(f)

    if manifest != nil {
        mw, err := zw.Create("META-INF/MANIFEST.MF")
        if err != nil {
            t.Fatalf("failed creating manifest entry: %v", err)
        }
        for k, v := range manifest {
            if _, err := io.WriteString(mw, k+": "+v+"\n"); err != nil {
                t.Fatalf("failed writing manifest: %v", err)
            }
        }
    }

    for path, content := range files {
        w, err := zw.Create(path)
        if err != nil {
            t.Fatalf("failed creating entry %s: %v", path, err)
        }
        if _, err := w.Write(content); err != nil {
            t.Fatalf("failed writing entry %s: %v", path, err)
        }
    }

    if err := zw.Close(); err != nil {
        t.Fatalf("failed closing zip writer: %v", err)
    }
    if err := f.Close(); err != nil {
        t.Fatalf("failed closing file: %v", err)
    }
    cleanup := func() { _ = os.Remove(jarPath) }
    return jarPath, cleanup
}

func TestGoodJarFile(t *testing.T) {
	jar, err := getJar(GOOD_JAR_NAME, t)

	if err != nil {
		return
	}

	if err := jar.scanArchive(); err != nil {
		t.Error("Error scanning archive", err)
	}
}

func TestManifestParsing(t *testing.T) {
	jar, err := getJar(GOOD_JAR_NAME, t)

	if err != nil {
		return
	}

	if err := jar.scanArchive(); err != nil {
		t.Error("Error scanning archive", err)
		return
	}

	value, ok := jar.Manifest["Main-Class"]
	if !ok {
		t.Error("Main-Class attribute should have been there, but wasn't")
	}
	if value != "jacobin.HelloWorld" {
		t.Error("Expected Main-Class to be 'jacobin.HelloWorld', but was " + value)
	}
}

func TestLoadClassSuccess(t *testing.T) {
	jar, err := getJar(GOOD_JAR_NAME, t)
	if err != nil {
		return
	}

	result, err := jar.loadClass("jacobin.HelloWorld")
	if err != nil {
		t.Error("Error loading class", err)
	}
	if !result.Success {
		t.Error("Loading class was not successful")
	}
}

func TestLoadClassDoesNotExist(t *testing.T) {
    jar, err := getJar(NO_MANIFEST_JAR_NAME, t)
    if err != nil {
        return
    }

	_, err = jar.loadClass("jacobin.HelloWorld")
	if err == nil {
		t.Error("Expected error loading class, but didn't get one, err")
	}
}

// Boundary: OpenArchive with non-existent file should error
func TestOpenArchiveNonExistent(t *testing.T) {
    _, err := OpenArchive(filepath.Join(t.TempDir(), "nope.jar"))
    if err == nil {
        t.Fatalf("expected error for non-existent jar, got nil")
    }
}

// Boundary: Empty jar should open and have no manifest/class entries
func TestEmptyJar(t *testing.T) {
    jarPath, cleanup := makeTempJar1(t, nil, map[string][]byte{})
    defer cleanup()
    jar, err := OpenArchive(jarPath)
    if err != nil {
        t.Fatalf("OpenArchive failed: %v", err)
    }
    if jar.hasResource("META-INF/MANIFEST.MF", TypeManifest) {
        t.Errorf("did not expect manifest in empty jar")
    }
}

// Boundary: Manifest parsing with malformed lines and extra colons
func TestManifestParsing_MalformedAndExtraColons(t *testing.T) {
    // Create a manifest content with various edge cases
    // Note: makeTempJar1 writes k: v per pair, so to create malformed lines
    // we add them as regular files and then overwrite the manifest explicitly.
    dir := t.TempDir()
    jarPath := filepath.Join(dir, "mf.jar")
    f, err := os.Create(jarPath)
    if err != nil { t.Fatalf("create: %v", err) }
    zw := zip.NewWriter(f)
    mw, err := zw.Create("META-INF/MANIFEST.MF")
    if err != nil { t.Fatalf("manifest create: %v", err) }
    content := "NoColonLine\nKeyOnly:\nKey:Value:Extra\n Spaced-Key :  spaced value  \n"
    if _, err := io.WriteString(mw, content); err != nil { t.Fatalf("write: %v", err) }
    if err := zw.Close(); err != nil { t.Fatalf("zw close: %v", err) }
    if err := f.Close(); err != nil { t.Fatalf("file close: %v", err) }

    jar, err := OpenArchive(jarPath)
    if err != nil { t.Fatalf("OpenArchive: %v", err) }

    // NoColonLine ignored
    if _, ok := jar.Manifest["NoColonLine"]; ok {
        t.Errorf("expected NoColonLine to be ignored")
    }
    // KeyOnly should exist with empty value
    if v, ok := jar.Manifest["KeyOnly"]; !ok || v != "" {
        t.Errorf("expected KeyOnly with empty value, got ok=%v v=%q", ok, v)
    }
    // Extra colons: only first value segment kept ("Value")
    if v := jar.Manifest["Key"]; v != "Value" {
        t.Errorf("expected Key to be 'Value', got %q", v)
    }
    // Trimmed key and value
    if v := jar.Manifest["Spaced-Key"]; v != "spaced value" {
        t.Errorf("expected trimmed key/value, got %q", v)
    }
}

// Boundary: getMainClass requires exact-case key
func TestGetMainClass_CaseSensitiveKey(t *testing.T) {
    jarPath, cleanup := makeTempJar1(t, map[string]string{"main-class": "com.example.Main"}, nil)
    defer cleanup()
    jar, err := OpenArchive(jarPath)
    if err != nil { t.Fatalf("OpenArchive: %v", err) }
    if got := jar.getMainClass(); got != "" {
        t.Errorf("expected empty main class for lowercase key, got %q", got)
    }
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// test helper: create a temporary JAR file with provided entries
// - TypeManifest: map of key->value to include in MANIFEST.MF (nil for none)
// - files: map[path]content to include arbitrary entries
// Returns full path to the jar and a cleanup func.
func makeTempJar(t *testing.T, manifest map[string]string, files map[string][]byte) (string, func()) {
	t.Helper()

	dir := t.TempDir()
	jarPath := filepath.Join(dir, "test.jar")
	f, err := os.Create(jarPath)
	if err != nil {
		t.Fatalf("failed creating temp jar: %v", err)
	}
	zw := zip.NewWriter(f)

	// Write TypeManifest if requested
	if manifest != nil {
		mw, err := zw.Create("META-INF/MANIFEST.MF")
		if err != nil {
			t.Fatalf("failed creating TypeManifest entry: %v", err)
		}
		// Build TypeManifest content lines, ensure trailing newline(s)
		for k, v := range manifest {
			line := k + ": " + v + "\n"
			if _, err := io.WriteString(mw, line); err != nil {
				t.Fatalf("failed writing TypeManifest content: %v", err)
			}
		}
	}

	// Write arbitrary files
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

func TestRecordFileAndHasResource_ForClassResourceAndManifest(t *testing.T) {
	manifest := map[string]string{"Main-Class": "com.example.Main"}
	files := map[string][]byte{
		"com/example/Main.class": {0xCA, 0xFE, 0xBA, 0xBE},
		"resources/config.txt":   []byte("hello"),
	}
	jarPath, cleanup := makeTempJar(t, manifest, files)

	jar, err := OpenArchive(jarPath)
	if err != nil {
		t.Fatalf("OpenArchive failed: %v", err)
	}

	// Class entry name should be normalized to dotted form without .class
	if !jar.hasResource("com.example.Main", TypeClassFile) {
		t.Errorf("expected hasResource for class entry to be true")
	}
	// Wrong type should be false
	if jar.hasResource("com.example.Main", TypeManifest) {
		t.Errorf("expected hasResource wrong type to be false")
	}
	// TypeResource should be tracked under original path
	if !jar.hasResource("resources/config.txt", TypeResource) {
		t.Errorf("expected hasResource for resource entry to be true")
	}
	// TypeManifest entry should be present
	if !jar.hasResource("META-INF/MANIFEST.MF", TypeManifest) {
		t.Errorf("expected hasResource for TypeManifest entry to be true")
	}

	cleanup()
}

func TestGetMainClass_PresentAndMissing(t *testing.T) {
	// With Main-Class
	jarWithMC, cleanup1 := makeTempJar(t, map[string]string{"Main-Class": "com.example.Main"}, map[string][]byte{})
	jar1, err := OpenArchive(jarWithMC)
	if err != nil {
		t.Fatalf("OpenArchive failed: %v", err)
	}
	if got := jar1.getMainClass(); got != "com.example.Main" {
		t.Errorf("getMainClass mismatch, got %q", got)
	}
	cleanup1()

	// Without Main-Class
	jarNoMC, cleanup2 := makeTempJar(t, map[string]string{"Class-Path": "lib/a.jar lib/b.jar"}, map[string][]byte{})
	jar2, err := OpenArchive(jarNoMC)
	if err != nil {
		t.Fatalf("OpenArchive failed: %v", err)
	}
	if got := jar2.getMainClass(); got != "" {
		t.Errorf("expected empty main class when missing, got %q", got)
	}
	cleanup2()
}

func TestGetClassPath_WithAndWithoutManifestEntry(t *testing.T) {
    manifest := map[string]string{"Class-Path": "lib/a.jar lib/b.jar"}
    jarPathWithCP, cleanup1 := makeTempJar(t, manifest, map[string][]byte{})
    jar1, err := OpenArchive(jarPathWithCP)
    if err != nil {
        t.Fatalf("OpenArchive failed: %v", err)
    }
    jar1.UpdateArchiveWithClassPath()
    cp := jar1.Classpath
    if len(cp) != 3 {
        t.Fatalf("expected classpath length 3, got %d: %#v", len(cp), cp)
    }
    baseDir := filepath.Dir(jarPathWithCP)
    expected1 := filepath.Join(baseDir, "lib/a.jar")
    expected2 := filepath.Join(baseDir, "lib/b.jar")
    if cp[0] != jarPathWithCP || cp[1] != expected1 || cp[2] != expected2 {
        t.Errorf("unexpected classpath: %#v (want [%q %q %q])", cp, jarPathWithCP, expected1, expected2)
    }
    cleanup1()

	jarPathNoCP, cleanup2 := makeTempJar(t, nil, map[string][]byte{})
	jar2, err := OpenArchive(jarPathNoCP)
	if err != nil {
		t.Fatalf("OpenArchive failed: %v", err)
	}
	jar2.UpdateArchiveWithClassPath()
	cp2 := jar2.Classpath
	if len(cp2) != 1 || cp2[0] != jarPathNoCP {
		t.Errorf("expected only jar filename in classpath, got %#v", cp2)
	}
	cleanup2()

	manifest = map[string]string{"ABC": "DEF"}
	path3, cleanup3 := makeTempJar(t, manifest, map[string][]byte{})
	jar3, err := OpenArchive(path3)
	if err != nil {
		t.Fatalf("OpenArchive failed: %v", err)
	}
	jar3.UpdateArchiveWithClassPath()
	cp3 := jar3.Classpath
	if len(cp3) != 1 {
		t.Fatalf("expected classpath length 1, got %d: %#v", len(cp3), cp3)
	}
	if !strings.HasSuffix(cp3[0], ".jar") {
		t.Errorf("unexpected classpath: %#v", cp3)
	}
	cleanup3()

}

func TestLoadClass_SuccessAndErrors(t *testing.T) {
	classBytes := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00}
	files := map[string][]byte{
		"com/example/Main.class": classBytes,
		"resources/readme.txt":   []byte("data"),
	}
	jarPath, cleanup := makeTempJar(t, map[string]string{"Main-Class": "com.example.Main"}, files)

	jar, err := OpenArchive(jarPath)
	if err != nil {
		t.Fatalf("OpenArchive failed: %v", err)
	}

	// Success for class load
	res, err := jar.loadClass("com.example.Main")
	if err != nil {
		t.Fatalf("loadClass failed: %v", err)
	}
	if !res.Success {
		t.Errorf("expected Success=true from loadClass")
	}
	if res.ResourceEntry.Name != "com.example.Main" || res.ResourceEntry.Type != TypeClassFile {
		t.Errorf("unexpected ResourceEntry: %#v", res.ResourceEntry)
	}
	if res.Data == nil || len(*res.Data) != len(classBytes) {
		t.Fatalf("unexpected data length: got %d want %d", len(*res.Data), len(classBytes))
	}
	for i := range classBytes {
		if (*res.Data)[i] != classBytes[i] {
			t.Fatalf("data byte %d mismatch: got %x want %x", i, (*res.Data)[i], classBytes[i])
		}
	}

	// Error: loading a non-class resource by its name should fail with type error
	if _, err := jar.loadClass("resources/readme.txt"); err == nil {
		t.Errorf("expected error when loading non-class resource as class")
	}

	// Error: missing class
	if _, err := jar.loadClass("com.example.Missing"); err == nil {
		t.Errorf("expected error when loading missing class")
	}

	cleanup()
}

// Boundary: CRLF in manifest and zero-length class handling
func TestManifestCRLFAndZeroLengthClass(t *testing.T) {
    // Build a jar with CRLF manifest and a zero-length class file
    dir := t.TempDir()
    jarPath := filepath.Join(dir, "crlf.jar")
    f, err := os.Create(jarPath)
    if err != nil { t.Fatalf("create: %v", err) }
    zw := zip.NewWriter(f)

    // CRLF manifest
    mw, err := zw.Create("META-INF/MANIFEST.MF")
    if err != nil { t.Fatalf("manifest create: %v", err) }
    // Note: use CRLF line endings
    content := "Main-Class: com.example.Main\r\nClass-Path: lib/a.jar lib/b.jar\r\n"
    if _, err := io.WriteString(mw, content); err != nil { t.Fatalf("write: %v", err) }

    // zero-length class entry
    cw, err := zw.Create("com/example/Empty.class")
    if err != nil { t.Fatalf("class create: %v", err) }
    if _, err := cw.Write([]byte{}); err != nil { t.Fatalf("class write: %v", err) }

    if err := zw.Close(); err != nil { t.Fatalf("close zw: %v", err) }
    if err := f.Close(); err != nil { t.Fatalf("close file: %v", err) }

    jar, err := OpenArchive(jarPath)
    if err != nil { t.Fatalf("OpenArchive: %v", err) }

    // Ensure CRLF parsing trims properly
    if got := jar.getMainClass(); got != "com.example.Main" {
        t.Fatalf("CRLF Main-Class parse mismatch: %q", got)
    }
    jar.UpdateArchiveWithClassPath()
    if len(jar.Classpath) != 3 {
        t.Fatalf("expected 3 classpath entries, got %d: %#v", len(jar.Classpath), jar.Classpath)
    }

    // Zero-length class should still load successfully, with 0 bytes
    res, err := jar.loadClass("com.example.Empty")
    if err != nil { t.Fatalf("loadClass zero-length: %v", err) }
    if !res.Success || res.Data == nil || len(*res.Data) != 0 {
        t.Fatalf("unexpected zero-length class load result: success=%v len=%d", res.Success, len(*res.Data))
    }
}

// Boundary: hasResource is case-sensitive for class names and types must match
func TestHasResource_CaseSensitivityAndType(t *testing.T) {
    files := map[string][]byte{
        "com/example/Main.class": {0xCA, 0xFE, 0xBA, 0xBE},
    }
    jarPath, cleanup := makeTempJar(t, map[string]string{"Main-Class": "com.example.Main"}, files)
    defer cleanup()

    jar, err := OpenArchive(jarPath)
    if err != nil { t.Fatalf("OpenArchive failed: %v", err) }

    // exact name true
    if !jar.hasResource("com.example.Main", TypeClassFile) {
        t.Fatalf("expected exact class to exist")
    }
    // different case false
    if jar.hasResource("com.example.main", TypeClassFile) {
        t.Fatalf("did not expect case-insensitive match")
    }
    // correct name but wrong type false
    if jar.hasResource("com.example.Main", TypeResource) {
        t.Fatalf("did not expect resource type match for class")
    }
}

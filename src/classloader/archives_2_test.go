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
	"testing"
)

// test helper: create a temporary JAR file with provided entries
// - manifest: map of key->value to include in MANIFEST.MF (nil for none)
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

	// Write manifest if requested
	if manifest != nil {
		mw, err := zw.Create("META-INF/MANIFEST.MF")
		if err != nil {
			t.Fatalf("failed creating manifest entry: %v", err)
		}
		// Build manifest content lines, ensure trailing newline(s)
		for k, v := range manifest {
			line := k + ": " + v + "\n"
			if _, err := io.WriteString(mw, line); err != nil {
				t.Fatalf("failed writing manifest content: %v", err)
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
	jarPath, _ := makeTempJar(t, manifest, files)

	jar, err := NewJarFile(jarPath)
	if err != nil {
		t.Fatalf("NewJarFile failed: %v", err)
	}

	// Class entry name should be normalized to dotted form without .class
	if !jar.hasResource("com.example.Main", ClassFile) {
		t.Errorf("expected hasResource for class entry to be true")
	}
	// Wrong type should be false
	if jar.hasResource("com.example.Main", Manifest) {
		t.Errorf("expected hasResource wrong type to be false")
	}
	// Resource should be tracked under original path
	if !jar.hasResource("resources/config.txt", Resource) {
		t.Errorf("expected hasResource for resource entry to be true")
	}
	// Manifest entry should be present
	if !jar.hasResource("META-INF/MANIFEST.MF", Manifest) {
		t.Errorf("expected hasResource for manifest entry to be true")
	}
}

func TestGetMainClass_PresentAndMissing(t *testing.T) {
	// With Main-Class
	jarWithMC, _ := makeTempJar(t, map[string]string{"Main-Class": "com.example.Main"}, map[string][]byte{})
	jar1, err := NewJarFile(jarWithMC)
	if err != nil { t.Fatalf("NewJarFile failed: %v", err) }
	if got := jar1.getMainClass(); got != "com.example.Main" {
		t.Errorf("getMainClass mismatch, got %q", got)
	}

	// Without Main-Class
	jarNoMC, _ := makeTempJar(t, map[string]string{"Class-Path": "lib/a.jar lib/b.jar"}, map[string][]byte{})
	jar2, err := NewJarFile(jarNoMC)
	if err != nil { t.Fatalf("NewJarFile failed: %v", err) }
	if got := jar2.getMainClass(); got != "" {
		t.Errorf("expected empty main class when missing, got %q", got)
	}
}

func TestGetClassPath_WithAndWithoutManifestEntry(t *testing.T) {
	manifest := map[string]string{"Class-Path": "lib/a.jar lib/b.jar"}
	jarPathWithCP, _ := makeTempJar(t, manifest, map[string][]byte{})
	jar1, err := NewJarFile(jarPathWithCP)
	if err != nil { t.Fatalf("NewJarFile failed: %v", err) }
	cp := jar1.getClassPath()
	if len(cp) != 3 {
		t.Fatalf("expected classpath length 3, got %d: %#v", len(cp), cp)
	}
	if cp[0] != jarPathWithCP || cp[1] != "lib/a.jar" || cp[2] != "lib/b.jar" {
		t.Errorf("unexpected classpath: %#v", cp)
	}

	jarPathNoCP, _ := makeTempJar(t, nil, map[string][]byte{})
	jar2, err := NewJarFile(jarPathNoCP)
	if err != nil { t.Fatalf("NewJarFile failed: %v", err) }
	cp2 := jar2.getClassPath()
	if len(cp2) != 1 || cp2[0] != jarPathNoCP {
		t.Errorf("expected only jar filename in classpath, got %#v", cp2)
	}
}

func TestLoadClass_SuccessAndErrors(t *testing.T) {
	classBytes := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00}
	files := map[string][]byte{
		"com/example/Main.class": classBytes,
		"resources/readme.txt":   []byte("data"),
	}
	jarPath, _ := makeTempJar(t, map[string]string{"Main-Class": "com.example.Main"}, files)

	jar, err := NewJarFile(jarPath)
	if err != nil { t.Fatalf("NewJarFile failed: %v", err) }

	// Success for class load
	res, err := jar.loadClass("com.example.Main")
	if err != nil {
		t.Fatalf("loadClass failed: %v", err)
	}
	if !res.Success {
		t.Errorf("expected Success=true from loadClass")
	}
	if res.ResourceEntry.Name != "com.example.Main" || res.ResourceEntry.Type != ClassFile {
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
}

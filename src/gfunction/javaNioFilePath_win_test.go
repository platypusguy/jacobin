package gfunction

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestWindowsPaths(t *testing.T) {
	// Save original globals
	oldOnWindows := globals.OnWindows

	// Simulate Windows
	globals.OnWindows = true
	defer func() {
		globals.OnWindows = oldOnWindows
	}()

	tests := []struct {
		path     string
		isAbs    bool
		root     string
		parent   string
		absMatch string
	}{
		{`C:\foo\bar`, true, `C:\`, `C:\foo`, `C:\foo\bar`},
		{`C:\`, true, `C:\`, "", `C:\`},
		{`\\server\share\file`, true, `\\server\share\`, `\\server\share`, `\\server\share\file`},
		{`\foo\bar`, false, `\`, `\foo`, `C:\foo\bar`}, // Rooted but not absolute on Windows
		{`foo\bar`, false, "", `foo`, `C:\foo\bar`},
	}

	globals.InitGlobals("test")
	globals.SetSystemProperty("user.dir", `C:\`)

	for _, tt := range tests {
		p := newPath(tt.path)

		// Test isAbsolute
		isAbsRaw := filePathIsAbsolute([]interface{}{p}).(int64)
		isAbs := isAbsRaw != 0
		if isAbs != tt.isAbs {
			t.Errorf("isAbsolute(%q) = %v; want %v", tt.path, isAbs, tt.isAbs)
		}

		// Test getRoot
		rootObj := filePathGetRoot([]interface{}{p})
		if tt.root == "" {
			if !object.IsNull(rootObj) {
				t.Errorf("getRoot(%q) = %v; want null", tt.path, rootObj)
			}
		} else {
			rootStr := object.GoStringFromStringObject(rootObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
			if rootStr != tt.root {
				t.Errorf("getRoot(%q) = %q; want %q", tt.path, rootStr, tt.root)
			}
		}

		// Test getParent
		parentObj := filePathGetParent([]interface{}{p})
		if tt.parent == "" {
			if !object.IsNull(parentObj) {
				t.Errorf("getParent(%q) = %v; want null", tt.path, parentObj)
			}
		} else {
			parentStr := object.GoStringFromStringObject(parentObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
			if parentStr != tt.parent {
				t.Errorf("getParent(%q) = %q; want %q", tt.path, parentStr, tt.parent)
			}
		}

		// Test toAbsolutePath
		absObj := filePathToAbsolutePath([]interface{}{p}).(*object.Object)
		absStr := object.GoStringFromStringObject(absObj.FieldTable["value"].Fvalue.(*object.Object))
		if absStr != tt.absMatch {
			t.Errorf("toAbsolutePath(%q) = %q; want %q", tt.path, absStr, tt.absMatch)
		}
	}

	// Test Normalize
	normTests := []struct {
		path     string
		expected string
	}{
		{`C:\foo\.\bar\..`, `C:\foo`},
		{`C:\foo\..\bar`, `C:\bar`},
		{`\foo\bar\..`, `\foo`},
		{`\\server\share\foo\..\bar`, `\\server\share\bar`},
	}

	for _, tt := range normTests {
		p := newPath(tt.path)
		normObj := filePathNormalize([]interface{}{p}).(*object.Object)
		normStr := object.GoStringFromStringObject(normObj.FieldTable["value"].Fvalue.(*object.Object))
		if normStr != tt.expected {
			t.Errorf("normalize(%q) = %q; want %q", tt.path, normStr, tt.expected)
		}
	}

	// Test Resolve
	resolveTests := []struct {
		base     string
		other    string
		expected string
	}{
		{`C:\abc`, `def`, `C:\abc\def`},
		{`C:\abc`, `\def`, `C:\def`},
		{`C:\abc`, `D:\def`, `D:\def`},
	}
	for _, tt := range resolveTests {
		p := newPath(tt.base)
		o := object.StringObjectFromGoString(tt.other)
		resObj := filePathResolve([]interface{}{p, o}).(*object.Object)
		resStr := object.GoStringFromStringObject(resObj.FieldTable["value"].Fvalue.(*object.Object))
		if resStr != tt.expected {
			t.Errorf("resolve(%q, %q) = %q; want %q", tt.base, tt.other, resStr, tt.expected)
		}
	}

	// Test Relativize
	relTests := []struct {
		base     string
		other    string
		expected string
	}{
		{`C:\a\b`, `C:\a\b\c\d`, `c\d`},
		{`C:\a\b\c\d`, `C:\a\b`, `..\..`},
		{`\\server\share\a`, `\\server\share\a\b`, `b`},
	}
	for _, tt := range relTests {
		p1 := newPath(tt.base)
		p2 := newPath(tt.other)
		resObj := filePathRelativize([]interface{}{p1, p2}).(*object.Object)
		resStr := object.GoStringFromStringObject(resObj.FieldTable["value"].Fvalue.(*object.Object))
		if resStr != tt.expected {
			t.Errorf("relativize(%q, %q) = %q; want %q", tt.base, tt.other, resStr, tt.expected)
		}
	}

	// Test Equals and CompareTo
	p1 := newPath(`C:\FOO`)
	p2 := newPath(`C:\foo`)
	if filePathEquals([]interface{}{p1, p2}) != types.JavaBoolTrue {
		t.Errorf("Expected C:\\FOO equals C:\\foo on Windows")
	}
	cmp := filePathCompareTo([]interface{}{p1, p2}).(int64)
	if cmp != 0 {
		t.Errorf("Expected C:\\FOO compareTo C:\\foo == 0 on Windows, got %v", cmp)
	}
}

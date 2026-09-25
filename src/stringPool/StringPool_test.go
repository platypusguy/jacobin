/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-6 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package stringPool

import (
	"io"
	"jacobin/src/globals"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

func TestGetStringIndexNilArg(t *testing.T) {
	globals.InitGlobals("test")

	index := GetStringIndex(nil)
	str := GetStringPointer(index)
	if str == nil || *str != "" {
		t.Errorf("Expected GetStringIndex(nil) to resolve to empty string, got: %v", str)
	}
}

func TestGetStringPointerOutOfRange(t *testing.T) {
	globals.InitGlobals("test")

	ptr := GetStringPointer(uint32(1_000_000))
	if ptr != nil {
		t.Errorf("Expected GetStringPointer with out-of-range index to return nil, got: %v", *ptr)
	}
}

func TestGetStringPointerReturnsCopy(t *testing.T) {
	globals.InitGlobals("test")

	str1 := "some unique test string"
	index := GetStringIndex(&str1)

	ptr1 := GetStringPointer(index)
	ptr2 := GetStringPointer(index)
	if ptr1 == ptr2 {
		t.Errorf("Expected GetStringPointer to return distinct pointers (copies) on each call")
	}
	if *ptr1 != *ptr2 {
		t.Errorf("Expected both pointers to point to equal string values, got '%s' and '%s'", *ptr1, *ptr2)
	}
}

func TestGetStringPoolSize(t *testing.T) {
	globals.InitGlobals("test")

	initialSize := GetStringPoolSize()

	str1 := "TestGetStringPoolSize unique string"
	GetStringIndex(&str1)

	newSize := GetStringPoolSize()
	if newSize != initialSize+1 {
		t.Errorf("Expected pool size to grow by 1, got initial=%d, new=%d", initialSize, newSize)
	}
}

func TestDumpStringPoolWithContext(t *testing.T) {
	globals.InitGlobals("test")

	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DumpStringPool("myContext")

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stdout = normalStdout

	output := string(msg)
	if !strings.Contains(output, "myContext") {
		t.Errorf("Expected DumpStringPool output to contain the context string 'myContext', got: %s", output)
	}
	if !strings.Contains(output, "DumpStringPool BEGIN") || !strings.Contains(output, "DumpStringPool END") {
		t.Errorf("Expected DumpStringPool output to contain BEGIN/END markers, got: %s", output)
	}
	if !strings.Contains(output, "java/lang/String") {
		t.Errorf("Expected DumpStringPool output to list pre-stored entries like 'java/lang/String', got: %s", output)
	}
}

func TestDumpStringPoolWithoutContext(t *testing.T) {
	globals.InitGlobals("test")

	normalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DumpStringPool("")

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stdout = normalStdout

	output := string(msg)
	if !strings.Contains(output, "DumpStringPool BEGIN") || !strings.Contains(output, "DumpStringPool END") {
		t.Errorf("Expected DumpStringPool output to contain BEGIN/END markers, got: %s", output)
	}
}

func TestPreloadArrayClassesToStringPool(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	PreloadArrayClassesToStringPool()

	arrayClasses := []string{
		types.BoolArray,
		types.JavaByteArray,
		types.DoubleArray,
		types.FloatArray,
		types.IntArray,
		types.LongArray,
		types.RefArray,
		types.RuneArray,
	}

	for _, className := range arrayClasses {
		name := className
		index := GetStringIndex(&name)
		str := GetStringPointer(index)
		if str == nil || *str != className {
			t.Errorf("Expected array class '%s' to be present in string pool, got: %v", className, str)
		}
	}

	// re-running PreloadArrayClassesToStringPool should not duplicate entries
	sizeBefore := GetStringPoolSize()
	PreloadArrayClassesToStringPool()
	sizeAfter := GetStringPoolSize()
	if sizeBefore != sizeAfter {
		t.Errorf("Expected PreloadArrayClassesToStringPool to be idempotent, got size before=%d, after=%d",
			sizeBefore, sizeAfter)
	}
}

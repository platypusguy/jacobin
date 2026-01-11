/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func TestMaxMemory(t *testing.T) {
	mem := maxMemory(nil)
	if mem.(int64) <= 1_000_000 {
		t.Errorf("maxMemory() = %d bytes; expected > 1,000,000", mem)
	}
}

func TestTotalMemory(t *testing.T) {
	mem := totalMemory(nil)
	if mem.(int64) <= 1_000_000 {
		t.Errorf("totalMemory() = %d bytes; expected > 1,000,000", mem)
	}
}

func TestFreeMemory(t *testing.T) {
	mem := freeMemory(nil)
	// It's hard to guarantee free memory > 0 without allocations,
	// but it should at least return a non-negative value.
	if mem.(int64) < 0 {
		t.Errorf("freeMemory() = %d bytes; expected >= 0", mem)
	}
}

func TestGC(t *testing.T) {
	res := runtimeGC(nil)
	if res != nil {
		t.Errorf("runtimeGC() returned %v; expected nil", res)
	}
}

func TestRuntimeCPUs(t *testing.T) {
	cpus := runtimeAvailableProcessors(nil)
	if cpus.(int64) < 1 {
		t.Errorf("runtimeCPUs() = %d; expected >= 1", cpus)
	}
}

func TestRuntimeVersion(t *testing.T) {
	globals.InitStringPool()
	res := runtimeVersion(nil)
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("runtimeVersion() did not return *object.Object, got %T", res)
	}
	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	if className != "java/lang/Runtime$Version" {
		t.Errorf("expected class java/lang/Runtime$Version, got %s", className)
	}
}

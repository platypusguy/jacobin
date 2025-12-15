/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package types

import "testing"

func TestClinitStatusConstants(t *testing.T) {
	if NoClInit != 0x00 {
		t.Fatalf("NoClInit expected 0x00, got 0x%02x", NoClInit)
	}
	if ClInitNotRun != 0x01 {
		t.Fatalf("ClInitNotRun expected 0x01, got 0x%02x", ClInitNotRun)
	}
	if ClInitInProgress != 0x02 {
		t.Fatalf("ClInitInProgress expected 0x02, got 0x%02x", ClInitInProgress)
	}
	if ClInitRun != 0x03 {
		t.Fatalf("ClInitRun expected 0x03, got 0x%02x", ClInitRun)
	}
}

func TestStringPoolRelatedConstants(t *testing.T) {
	// String/Object class names
	if ObjectClassName != "java/lang/Object" {
		t.Fatalf("ObjectClassName mismatch: %q", ObjectClassName)
	}
	if PtrToJavaLangObject == nil || *PtrToJavaLangObject != ObjectClassName {
		t.Fatalf("PtrToJavaLangObject does not point to ObjectClassName")
	}
	if StringClassName != "java/lang/String" {
		t.Fatalf("StringClassName mismatch: %q", StringClassName)
	}
	if StringClassRef != "Ljava/lang/String;" {
		t.Fatalf("StringClassRef mismatch: %q", StringClassRef)
	}
	if ModuleClassRef != "Ljava/lang/Module;" {
		t.Fatalf("ModuleClassRef mismatch: %q", ModuleClassRef)
	}

	// Pool indices per globals.InitStringPool contract
	if StringPoolStringIndex != 1 {
		t.Fatalf("StringPoolStringIndex expected 1, got %d", StringPoolStringIndex)
	}
	if StringPoolObjectIndex != 2 {
		t.Fatalf("StringPoolObjectIndex expected 2, got %d", StringPoolObjectIndex)
	}
}

func TestMiscStringConstants(t *testing.T) {
	if EmptyString != "" {
		t.Fatalf("EmptyString expected empty, got %q", EmptyString)
	}
	if NullString != "null" {
		t.Fatalf("NullString expected 'null', got %q", NullString)
	}
}

func TestInvalidStringIndexAndStackInflator(t *testing.T) {
	if InvalidStringIndex != 0xffffffff {
		t.Fatalf("InvalidStringIndex expected 0xffffffff, got 0x%x", InvalidStringIndex)
	}
	if StackInflator != 2 {
		t.Fatalf("StackInflator expected 2, got %d", StackInflator)
	}
}

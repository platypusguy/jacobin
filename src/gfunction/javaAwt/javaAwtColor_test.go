/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaAwt

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"reflect"
	"testing"
)

func TestLoad_Awt_Color_RegistersMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Awt_Color()

	expectedMethods := []string{
		"java/awt/Color.<clinit>()V",
		"java/awt/Color.<init>(FFF)V",
		"java/awt/Color.<init>(FFFF)V",
		"java/awt/Color.<init>(I)V",
		"java/awt/Color.<init>(IZ)V",
		"java/awt/Color.<init>(II)V",
		"java/awt/Color.<init>(III)V",
		"java/awt/Color.<init>(Ljava/awt/color/ColorSpace;[FF)V",
		"java/awt/Color.initIDs()V",
		"java/awt/Color.hashCode()I",
	}

	for _, sig := range expectedMethods {
		meth, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("missing signature %s", sig)
			continue
		}
		if meth.GFunction == nil {
			t.Errorf("expected non-nil GFunction for %s", sig)
		}
	}
}

func TestInitColorRGB(t *testing.T) {
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: "I", Fvalue: int64(0)}

	// Valid RGB
	res := initColorRGB([]any{obj, int64(255), int64(0), int64(128)})
	if res != nil {
		t.Fatalf("expected nil result on valid initColorRGB, got %v", res)
	}

	valField := obj.FieldTable["value"]
	packed := int64((255 << 24) | (255 << 16) | (0 << 8) | 128)
	expectedPacked := int64(int32(packed))
	if valField.Fvalue != expectedPacked {
		t.Errorf("expected packed value %d, got %v", expectedPacked, valField.Fvalue)
	}

	// Invalid param 0 type
	err := initColorRGB([]any{"not-an-object", int64(1), int64(2), int64(3)})
	errBlk, ok := err.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", err)
	}

	// Invalid range
	errRange := initColorRGB([]any{obj, int64(300), int64(0), int64(0)})
	errBlkRange, ok := errRange.(*ghelpers.GErrBlk)
	if !ok || errBlkRange.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for out-of-range component, got %v", errRange)
	}
}

func TestColorHashCode(t *testing.T) {
	obj := object.MakeEmptyObject()
	packed := int64(0xFF000002)
	expectedHash := int64(int32(packed))
	obj.FieldTable["value"] = object.Field{Ftype: "I", Fvalue: expectedHash}

	res := colorHashCode([]any{obj})
	if !reflect.DeepEqual(res, expectedHash) {
		t.Errorf("expected hash %d, got %v", expectedHash, res)
	}

	// Invalid target
	err := colorHashCode([]any{"invalid"})
	if _, ok := err.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk for invalid object, got %v", err)
	}

	// Missing field
	objNoField := object.MakeEmptyObject()
	errMissing := colorHashCode([]any{objNoField})
	if _, ok := errMissing.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk for missing value field, got %v", errMissing)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction_test

import (
	"reflect"
	"testing"

	"jacobin/src/gfunction"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestNewRunnable_PopulatesClassAndFields(t *testing.T) {
	cl := object.JavaByteArrayFromGoString("java/lang/Thread")
	mn := object.JavaByteArrayFromGoString("run")
	sg := object.JavaByteArrayFromGoString("()V")

	obj := gfunction.NewRunnable(cl, mn, sg)
	if obj == nil {
		t.Fatalf("TestNewRunnable returned nil")
	}

	// Verify class name
	gotClass := object.GoStringFromStringPoolIndex(obj.KlassName)
	wantClass := "java/lang/Runnable"
	if gotClass != wantClass {
		t.Fatalf("class name = %q; want %q", gotClass, wantClass)
	}

	// Helper to check one field
	checkField := func(fieldName string, want []types.JavaByte) {
		fld, ok := obj.FieldTable[fieldName]
		if !ok {
			t.Fatalf("missing field %q", fieldName)
		}
		if fld.Ftype != types.ByteArray {
			t.Fatalf("field %q type = %q; want %q", fieldName, fld.Ftype, types.ByteArray)
		}
		gotSlice, ok := fld.Fvalue.([]types.JavaByte)
		if !ok {
			t.Fatalf("field %q value has unexpected type %T; want []types.JavaByte", fieldName, fld.Fvalue)
		}
		if !reflect.DeepEqual(gotSlice, want) {
			t.Fatalf("field %q value = %v; want %v", fieldName, gotSlice, want)
		}
	}

	checkField("clName", cl)
	checkField("methName", mn)
	checkField("signature", sg)
}

func TestNewRunnable_AllowsEmptyAndNil(t *testing.T) {
	empty := []types.JavaByte{}
	var nilSlice []types.JavaByte

	obj1 := gfunction.NewRunnable(empty, empty, empty)
	if obj1 == nil {
		t.Fatalf("TestNewRunnable with empty slices returned nil")
	}

	// empty slices should be preserved
	if got := obj1.FieldTable["clName"].Fvalue.([]types.JavaByte); !reflect.DeepEqual(got, empty) {
		t.Fatalf("empty clName not preserved; got %v", got)
	}
	if got := obj1.FieldTable["methName"].Fvalue.([]types.JavaByte); !reflect.DeepEqual(got, empty) {
		t.Fatalf("empty methName not preserved; got %v", got)
	}
	if got := obj1.FieldTable["signature"].Fvalue.([]types.JavaByte); !reflect.DeepEqual(got, empty) {
		t.Fatalf("empty signature not preserved; got %v", got)
	}

	obj2 := gfunction.NewRunnable(nilSlice, nilSlice, nilSlice)
	if obj2 == nil {
		t.Fatalf("NewRunnable with nil slices returned nil")
	}

	// nil slices are allowed; values may be stored as a nil slice or an empty slice
	// Verify type is []types.JavaByte and length is zero
	for _, name := range []string{"clName", "methName", "signature"} {
		fld := obj2.FieldTable[name]
		if fld.Ftype != types.ByteArray {
			t.Fatalf("field %q type = %q; want %q", name, fld.Ftype, types.ByteArray)
		}
		gotSlice, ok := fld.Fvalue.([]types.JavaByte)
		if !ok {
			t.Fatalf("field %q value type %T; want []types.JavaByte", name, fld.Fvalue)
		}
		if len(gotSlice) != 0 {
			t.Fatalf("field %q length = %d; want 0 for nil/empty input", name, len(gotSlice))
		}
	}
}

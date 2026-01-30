package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/types"
	"testing"
)

func TestObjects_CheckFromIndexSize_CurrentBehavior(t *testing.T) {
	globals.InitStringPool()

	// Valid under current impl only when fromIndex == 0 (length bug uses size as length)
	if v := objectsCheckFromIndexSize([]interface{}{int64(0), int64(3), int64(10)}); v != int64(0) {
		t.Fatalf("checkFromIndexSize expected 0, got %v", v)
	}

	// fromIndex > 0 should trigger IndexOutOfBoundsException due to buggy length usage
	if err := objectsCheckFromIndexSize([]interface{}{int64(1), int64(3), int64(10)}); err == nil {
		t.Fatalf("expected error for fromIndex > 0")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}

	// Negative size -> error
	if err := objectsCheckFromIndexSize([]interface{}{int64(0), int64(-1), int64(10)}); err == nil {
		t.Fatalf("expected error for negative size")
	}
}

func TestObjects_CheckFromToIndex_CurrentBehavior(t *testing.T) {
	globals.InitStringPool()

	// Valid when 0 <= fromIndex <= toIndex; length check uses toIndex as length
	if v := objectsCheckFromToIndex([]interface{}{int64(2), int64(5), int64(10)}); v != int64(2) {
		t.Fatalf("checkFromToIndex expected 2, got %v", v)
	}

	// fromIndex > toIndex -> error
	if err := objectsCheckFromToIndex([]interface{}{int64(6), int64(5), int64(10)}); err == nil {
		t.Fatalf("expected error for fromIndex > toIndex")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}
}

func TestObjects_CheckIndex(t *testing.T) {
	globals.InitStringPool()

	// Valid index within [0, length)
	if v := objectsCheckIndex([]interface{}{int64(3), int64(10)}); v != int64(3) {
		t.Fatalf("checkIndex expected 3, got %v", v)
	}

	// index == length -> error
	if err := objectsCheckIndex([]interface{}{int64(10), int64(10)}); err == nil {
		t.Fatalf("expected error for index == length")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}

	// negative index -> error
	if err := objectsCheckIndex([]interface{}{int64(-1), int64(10)}); err == nil {
		t.Fatalf("expected error for negative index")
	}
}

func TestObjects_IsNull_And_NonNull_CurrentBehavior(t *testing.T) {
	globals.InitStringPool()

	// Non-object (e.g., int64) -> isNull true, nonNull false
	if v := objectsIsNull([]interface{}{int64(5)}); v.(types.JavaBool) != types.JavaBoolTrue {
		t.Fatalf("isNull for non-object expected true (1), got %v", v)
	}
	if v := objectsNonNull([]interface{}{int64(5)}); v.(types.JavaBool) != types.JavaBoolFalse {
		t.Fatalf("nonNull for non-object expected false, got %v", v)
	}

	// Null object (typed-nil) -> current impl treats it as non-null (type assertion ok)
	if v := objectsIsNull([]interface{}{nil}); v.(types.JavaBool) != types.JavaBoolTrue { // nil is not *object.Object, so true
		t.Fatalf("isNull(nil) expected true (1), got %v", v)
	}
	if v := objectsNonNull([]interface{}{nil}); v.(types.JavaBool) != types.JavaBoolFalse {
		t.Fatalf("nonNull(nil) expected false, got %v", v)
	}
}

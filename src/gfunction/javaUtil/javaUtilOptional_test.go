package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helper to make an Optional object with a specific field value
func makeOptionalWithValue(val interface{}) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameOptional)
	obj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: val}
	return obj
}

func TestOptional_Empty_Behavior(t *testing.T) {
	globals.InitStringPool()

	// empty()
	opt := optionalEmpty([]interface{}{}).(*object.Object)

	// isEmpty -> true; isPresent -> false
	if v := optionalIsEmpty([]interface{}{opt}).(int64); v != types.JavaBoolTrue {
		t.Fatalf("isEmpty expected true, got %d", v)
	}
	if v := optionalIsPresent([]interface{}{opt}).(int64); v != types.JavaBoolFalse {
		t.Fatalf("isPresent expected false, got %d", v)
	}

	// get on empty -> NoSuchElementException
	if err := optionalGet([]interface{}{opt}); err == nil {
		t.Fatalf("expected error for get on empty optional")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.NoSuchElementException {
			t.Fatalf("expected NoSuchElementException, got %d", geb.ExceptionType)
		}
	}

	// toString -> "empty"
	sObj := optionalToString([]interface{}{opt}).(*object.Object)
	if s := object.GoStringFromStringObject(sObj); s != "empty" {
		t.Fatalf("toString for empty mismatch: %q", s)
	}

	// orElse with default should return that object's value
	def := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(7))
	res := optionalOrElse([]interface{}{opt, def})
	if res.(int64) != 7 {
		t.Fatalf("orElse on empty expected 7, got %v", res)
	}

	// orElseThrow on empty -> NoSuchElementException
	if err := optionalOrElseThrow([]interface{}{opt}); err == nil {
		t.Fatalf("expected NoSuchElementException from orElseThrow on empty")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.NoSuchElementException {
			t.Fatalf("expected NoSuchElementException, got %d", geb.ExceptionType)
		}
	}
}

func TestOptional_Present_Behavior(t *testing.T) {
	globals.InitStringPool()

	// Build Optional with an int64 value 42 (predictable toString formatting)
	opt := object.MakeEmptyObjectWithClassName(&types.ClassNameOptional)
	opt.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(42)}

	if v := optionalIsEmpty([]interface{}{opt}).(int64); v != types.JavaBoolFalse {
		t.Fatalf("isEmpty expected false, got %d", v)
	}
	if v := optionalIsPresent([]interface{}{opt}).(int64); v != types.JavaBoolTrue {
		t.Fatalf("isPresent expected true, got %d", v)
	}

	// get returns the raw stored value (int64)
	if v := optionalGet([]interface{}{opt}).(int64); v != 42 {
		t.Fatalf("get expected 42, got %d", v)
	}

	// toString reflects type and value
	sObj := optionalToString([]interface{}{opt}).(*object.Object)
	if s := object.GoStringFromStringObject(sObj); s != "Optional[int64 :: 42]" {
		t.Fatalf("toString mismatch: %q", s)
	}

	// orElse with default should return the present value, not the default
	def := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(7))
	if v := optionalOrElse([]interface{}{opt, def}).(int64); v != 42 {
		t.Fatalf("orElse expected 42, got %d", v)
	}

	// orElseThrow should return the value, not throw
	if v := optionalOrElseThrow([]interface{}{opt}).(int64); v != 42 {
		t.Fatalf("orElseThrow expected 42, got %d", v)
	}
}

func TestOptional_Equals_Cases(t *testing.T) {
	globals.InitStringPool()

	// Two empty optionals -> equal (true)
	a := optionalEmpty([]interface{}{}).(*object.Object)
	b := optionalEmpty([]interface{}{}).(*object.Object)
	if v := optionalEquals([]interface{}{a, b}).(int64); v != types.JavaBoolTrue {
		t.Fatalf("equals(empty, empty) expected true, got %d", v)
	}

	// Two present with same primitive value -> true per current impl (value comparison)
	o1 := object.MakeEmptyObjectWithClassName(&types.ClassNameOptional)
	o1.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(5)}
	o2 := object.MakeEmptyObjectWithClassName(&types.ClassNameOptional)
	o2.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(5)}
	if v := optionalEquals([]interface{}{o1, o2}).(int64); v != types.JavaBoolTrue {
		t.Fatalf("equals(5,5) expected true, got %d", v)
	}

	// Present vs empty -> false
	if v := optionalEquals([]interface{}{o1, a}).(int64); v != types.JavaBoolFalse {
		t.Fatalf("equals(present, empty) expected false, got %d", v)
	}

	// Non-object argument -> error
	if err := optionalEquals([]interface{}{o1, int64(3)}); err == nil {
		t.Fatalf("expected error for equals with non-object param")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
		}
	}
}

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func newLevelObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&levelClassName)
}

func TestLoggingLevelInit_And_Getters(t *testing.T) {
	globals.InitStringPool()

	lvl := newLevelObj()
	name := object.StringObjectFromGoString("INFO")
	ret := loggingLevelInit([]interface{}{lvl, name, int64(800)})
	if ret != nil {
		t.Fatalf("loggingLevelInit returned error: %v", ret)
	}

	if got := object.GoStringFromStringObject(loggingLevelGetName([]interface{}{lvl}).(*object.Object)); got != "INFO" {
		t.Fatalf("expected name INFO, got %q", got)
	}
	if got := loggingLevelIntValue([]interface{}{lvl}).(int64); got != 800 {
		t.Fatalf("expected intValue 800, got %d", got)
	}
	if ret := loggingLevelGetResourceBundleName([]interface{}{lvl}); ret != object.Null {
		t.Fatalf("expected null resource bundle name, got %v", ret)
	}
	if got := object.GoStringFromStringObject(loggingLevelToString([]interface{}{lvl}).(*object.Object)); got != "INFO" {
		t.Fatalf("expected toString INFO, got %q", got)
	}
	if got := loggingLevelHashCode([]interface{}{lvl}).(int64); got != 800 {
		t.Fatalf("expected hashCode 800, got %d", got)
	}
}

func TestLoggingLevelInitWithResourceBundle(t *testing.T) {
	globals.InitStringPool()

	lvl := newLevelObj()
	name := object.StringObjectFromGoString("CONFIG")
	rb := object.StringObjectFromGoString("MyBundle")
	ret := loggingLevelInitWithResourceBundle([]interface{}{lvl, name, int64(700), rb})
	if ret != nil {
		t.Fatalf("loggingLevelInitWithResourceBundle returned error: %v", ret)
	}

	if got := object.GoStringFromStringObject(loggingLevelGetResourceBundleName([]interface{}{lvl}).(*object.Object)); got != "MyBundle" {
		t.Fatalf("expected resource bundle name MyBundle, got %q", got)
	}
	if got := object.GoStringFromStringObject(loggingLevelGetLocalizedName([]interface{}{lvl}).(*object.Object)); got != "CONFIG" {
		t.Fatalf("expected localized name CONFIG, got %q", got)
	}
}

func TestLoggingLevelEquals(t *testing.T) {
	globals.InitStringPool()

	a := makeLevelObject("INFO", 800, "")
	b := makeLevelObject("INFO", 800, "")
	c := makeLevelObject("WARNING", 900, "")

	if ret := loggingLevelEquals([]interface{}{a, b}); ret != types.JavaBoolTrue {
		t.Fatalf("expected equal levels to be equal")
	}
	if ret := loggingLevelEquals([]interface{}{a, c}); ret != types.JavaBoolFalse {
		t.Fatalf("expected different levels to be unequal")
	}
	if ret := loggingLevelEquals([]interface{}{a, object.Null}); ret != types.JavaBoolFalse {
		t.Fatalf("expected null comparison to be false")
	}
}

func TestLoggingLevelParse_StandardAndNumeric(t *testing.T) {
	globals.InitStringPool()

	// Standard name
	ret := loggingLevelParse([]interface{}{object.StringObjectFromGoString("SEVERE")})
	lvl, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret)
	}
	if got := loggingLevelIntValue([]interface{}{lvl}).(int64); got != 1000 {
		t.Fatalf("expected intValue 1000 for SEVERE, got %d", got)
	}

	// Numeric string
	ret2 := loggingLevelParse([]interface{}{object.StringObjectFromGoString("123")})
	lvl2, ok := ret2.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret2)
	}
	if got := loggingLevelIntValue([]interface{}{lvl2}).(int64); got != 123 {
		t.Fatalf("expected intValue 123, got %d", got)
	}

	// Invalid string
	ret3 := loggingLevelParse([]interface{}{object.StringObjectFromGoString("NOT_A_LEVEL")})
	geb, ok := ret3.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", ret3)
	}
	if geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
	}
}

func TestLoggingLevel_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingLevelGetName([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingLevelIntValue([]interface{}{nil}); ret == nil {
		t.Fatalf("expected error for nil param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}

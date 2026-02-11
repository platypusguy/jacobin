/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestNamedParameterSpecInit(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/spec/NamedParameterSpec"
	specObj := object.MakeEmptyObjectWithClassName(&className)
	nameStr := "X25519"
	nameObj := object.StringObjectFromGoString(nameStr)

	// Test success
	params := []interface{}{specObj, nameObj}
	result := namedParameterSpecInit(params)
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	fieldEntry, exists := specObj.FieldTable["name"]
	if !exists {
		t.Fatal("Expected 'name' field to exist")
	}
	if fieldEntry.Ftype != types.StringClassName {
		t.Errorf("Expected field type %s, got %s", types.StringClassName, fieldEntry.Ftype)
	}
	if fieldEntry.Fvalue.(*object.Object) != nameObj {
		t.Errorf("Expected field value %v, got %v", nameObj, fieldEntry.Fvalue)
	}

	// Test error: wrong number of arguments
	result = namedParameterSpecInit([]interface{}{specObj})
	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for wrong number of arguments, got %v", result)
	}

	// Test error: invalid self object
	result = namedParameterSpecInit([]interface{}{"not an object", nameObj})
	errBlk, ok = result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid self object, got %v", result)
	}

	// Test error: invalid name object
	result = namedParameterSpecInit([]interface{}{specObj, "not an object"})
	errBlk, ok = result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid name object, got %v", result)
	}

	// Test error: name is not a String object
	notStringObj := object.MakeEmptyObjectWithClassName(&className)
	result = namedParameterSpecInit([]interface{}{specObj, notStringObj})
	errBlk, ok = result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for non-string name object, got %v", result)
	}
}

func TestNamedParameterSpecGetName(t *testing.T) {
	globals.InitGlobals("test")

	className := "java/security/spec/NamedParameterSpec"
	specObj := object.MakeEmptyObjectWithClassName(&className)
	nameStr := "Ed25519"
	nameObj := object.StringObjectFromGoString(nameStr)

	// Test success
	specObj.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: nameObj,
	}

	result := namedParameterSpecGetName([]interface{}{specObj})
	resObj, ok := result.(*object.Object)
	if !ok || resObj != nameObj {
		t.Errorf("Expected name object %v, got %v", nameObj, result)
	}

	// Test error: missing arguments
	result = namedParameterSpecGetName([]interface{}{})
	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for missing arguments, got %v", result)
	}

	// Test error: invalid self object
	result = namedParameterSpecGetName([]interface{}{"not an object"})
	errBlk, ok = result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid self object, got %v", result)
	}

	// Test error: name field not set
	emptyObj := object.MakeEmptyObjectWithClassName(&className)
	result = namedParameterSpecGetName([]interface{}{emptyObj})
	errBlk, ok = result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for unset name field, got %v", result)
	}

	// Test error: name field has invalid type
	invalidFieldObj := object.MakeEmptyObjectWithClassName(&className)
	invalidFieldObj.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: "not an object",
	}
	result = namedParameterSpecGetName([]interface{}{invalidFieldObj})
	errBlk, ok = result.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid name field type, got %v", result)
	}
}

func TestNamedParameterSpecClinit(t *testing.T) {
	globals.InitGlobals("test")

	result := namedParameterSpecClinit([]interface{}{})
	if result != nil {
		t.Errorf("Expected nil result from clinit, got %v", result)
	}

	className := "java/security/spec/NamedParameterSpec"
	fields := []string{"X25519", "X448", "ED25519", "ED448"}

	for _, field := range fields {
		fullFieldName := className + "." + field
		staticVar, ok := statics.QueryStatic(className, field)
		if !ok {
			t.Errorf("Expected static field %s to exist", fullFieldName)
			continue
		}
		if staticVar.Type != types.Ref {
			t.Errorf("Expected static field %s to have type Ref, got %v", fullFieldName, staticVar.Type)
		}
		specObj, ok := staticVar.Value.(*object.Object)
		if !ok || specObj == nil {
			t.Errorf("Expected static field %s value to be an *object.Object", fullFieldName)
			continue
		}

		// Verify the name inside the spec object
		nameField, exists := specObj.FieldTable["name"]
		if !exists {
			t.Errorf("Spec object for %s missing name field", fullFieldName)
			continue
		}
		nameObj := nameField.Fvalue.(*object.Object)
		goName := object.GoStringFromStringObject(nameObj)

		// Note: the code uses "X25519", "X448", "Ed25519", "Ed448" as names
		expectedName := field
		if field == "ED25519" {
			expectedName = "Ed25519"
		} else if field == "ED448" {
			expectedName = "Ed448"
		}

		if goName != expectedName {
			t.Errorf("Expected spec name %s, got %s", expectedName, goName)
		}
	}
}

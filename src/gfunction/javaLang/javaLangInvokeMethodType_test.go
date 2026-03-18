/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/classloader"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"testing"
)

func TestResolveTypeDescriptorPrimitives(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()

	// Initialize the wrapper class for 'int' so its TYPE field is populated.
	integerClinit(nil)

	// Test resolving the primitive descriptor 'I'
	classObj, err := resolveTypeDescriptor("I")
	if err != nil {
		t.Fatalf("Unexpected error resolving primitive type 'I': %v", err)
	}
	if classObj == nil || object.IsNull(classObj) {
		t.Fatalf("resolveTypeDescriptor returned nil for primitive type 'I'")
	}

	// Verify the returned object represents the primitive class "int"
	classObj.ThMutex.RLock()
	nameField, ok := classObj.FieldTable["name"]
	classObj.ThMutex.RUnlock()

	if !ok {
		t.Fatalf("Primitive Class object is missing the 'name' field")
	}

	nameStr, ok := nameField.Fvalue.(string)
	if !ok {
		t.Fatalf("Primitive Class object 'name' field is not a string, got %T", nameField.Fvalue)
	}

	if nameStr != "int" {
		t.Errorf("Expected primitive class name to be 'int', got '%s'", nameStr)
	}
}

func TestResolveTypeDescriptorObjects(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()

	// Test resolving an object descriptor
	descriptor := "Ljava/lang/String;"
	classObj, err := resolveTypeDescriptor(descriptor)
	if err != nil {
		t.Fatalf("Unexpected error resolving object type '%s': %v", descriptor, err)
	}
	if classObj == nil || object.IsNull(classObj) {
		t.Fatalf("resolveTypeDescriptor returned nil for object type '%s'", descriptor)
	}

	// Verify the returned object represents the class "java/lang/String"
	classObj.ThMutex.RLock()
	nameField, ok := classObj.FieldTable["name"]
	classObj.ThMutex.RUnlock()

	if !ok {
		t.Fatalf("Class object is missing the 'name' field")
	}

	nameStr, ok := nameField.Fvalue.(string)
	if !ok {
		t.Fatalf("Class object 'name' field is not a string, got %T", nameField.Fvalue)
	}

	expectedName := "java/lang/String"
	if nameStr != expectedName {
		t.Errorf("Expected class name to be '%s', got '%s'", expectedName, nameStr)
	}
}

func TestResolveTypeDescriptorArrays(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()

	// Test resolving an array descriptor
	descriptor := "[I"
	classObj, err := resolveTypeDescriptor(descriptor)
	if err != nil {
		t.Fatalf("Unexpected error resolving array type '%s': %v", descriptor, err)
	}
	if classObj == nil || object.IsNull(classObj) {
		t.Fatalf("resolveTypeDescriptor returned nil for array type '%s'", descriptor)
	}

	// Verify the returned object represents the array class "[I"
	classObj.ThMutex.RLock()
	nameField, ok := classObj.FieldTable["name"]
	classObj.ThMutex.RUnlock()

	if !ok {
		t.Fatalf("Array Class object is missing the 'name' field")
	}

	nameStr, ok := nameField.Fvalue.(string)
	if !ok {
		t.Fatalf("Array Class object 'name' field is not a string, got %T", nameField.Fvalue)
	}

	expectedName := "[I"
	if nameStr != expectedName {
		t.Errorf("Expected array class name to be '%s', got '%s'", expectedName, nameStr)
	}
}

func TestParseDescriptorToClasses_Invalid(t *testing.T) {
	globals.InitGlobals("test")
	
	// Test invalid descriptors
	invalidDescriptors := []string{
		"",
		"()", // Missing return type
		"(I", // Missing closing paren
		"I)V", // Missing opening paren
		"(Ljava/lang/String)V", // Missing semicolon
	}

	for _, desc := range invalidDescriptors {
		_, _, err := parseDescriptorToClasses(desc)
		if err == nil {
			t.Errorf("Expected error for invalid descriptor: %s", desc)
		}
	}
}

// Get the tokens of the parameter types for methods that don't take arrays
func TestGetNextTypeDescriptorNonArray(t *testing.T) {
	// (ZILjava/lang/String;)V is passed in as: ZILjava/lang/String;
	paramStr := "ZILjava/lang/String;"
	results := make([]string, 0)
	for i := 0; i < len(paramStr); {
		typeStr, width := getNextTypeDescriptor(paramStr[i:])
		results = append(results, typeStr)
		i += width
	}
	if results[0] != "Z" || results[1] != "I" ||
		results[2] != "Ljava/lang/String;" {
		t.Errorf("Expected Z, I, and Ljava/lang/String;, got %s", results)
	}
}

// Get the tokens of the parameter types for methods that take arrays
func TestGetNextTypeDescriptorWithArray(t *testing.T) {
	// (ZI[[Ljava/lang/String;)V is passed in as: ZI[[Ljava/lang/String;
	paramStr := "ZI[[Ljava/lang/String;"
	results := make([]string, 0)
	for i := 0; i < len(paramStr); {
		typeStr, width := getNextTypeDescriptor(paramStr[i:])
		results = append(results, typeStr)
		i += width
	}
	if results[0] != "Z" || results[1] != "I" ||
		results[2] != "[[Ljava/lang/String;" {
		t.Errorf("Expected Z, I, and [[Ljava/lang/String;, got %s", results)
	}
}

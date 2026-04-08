/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"testing"
)

func TestMakeJlcObject(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	className := "java/lang/String"
	jlcObj := MakeJlcObject(className)

	if jlcObj == nil {
		t.Fatal("MakeJlcObject returned nil")
	}

	// Verify the KlassName of the returned object is java/lang/Class
	// (StringPoolJavaLangClassIndex is used in MakeJlcObject)
	classNameStrPtr := stringPool.GetStringPointer(jlcObj.KlassName)
	if classNameStrPtr == nil || *classNameStrPtr != "java/lang/Class" {
		t.Errorf("Expected KlassName to point to 'java/lang/Class', got '%v'",
			*classNameStrPtr)
	}

	// Verify the "name" field contains a String object with the correct class name
	jlcObj.ThMutex.RLock()
	nameField, ok := jlcObj.FieldTable["name"]
	jlcObj.ThMutex.RUnlock()

	if !ok {
		t.Fatal("Jlc object missing 'name' field")
	}

	nameObj, ok := nameField.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("Jlc object 'name' field is not an *object.Object, got %T", nameField.Fvalue)
	}

	actualName := object.GoStringFromStringObject(nameObj)
	if actualName != className {
		t.Errorf("Expected Jlc name to be '%s', got '%s'", className, actualName)
	}

	// Verify $statics field is an empty string slice
	jlcObj.ThMutex.RLock()
	staticsField, ok := jlcObj.FieldTable["$statics"]
	jlcObj.ThMutex.RUnlock()

	if !ok {
		t.Fatal("Jlc object missing '$statics' field")
	}

	staticsSlice, ok := staticsField.Fvalue.([]string)
	if !ok {
		t.Fatalf("Jlc object '$statics' field is not a []string, got %T", staticsField.Fvalue)
	}

	if len(staticsSlice) != 0 {
		t.Errorf("Expected empty '$statics' slice, got length %d", len(staticsSlice))
	}
}

// Test that the String class is loaded and initialized correctly
// and that the java/lang/Class instance in the loaded class is correct
func TestStringClassStaticsLoaded(t *testing.T) {
	// Full initialization required to load base classes and string pool
	globals.InitGlobals("test")
	trace.Init()
	Init() // This calls InitMethodArea, JmodMapInit, GetBaseJmodBytes

	// Load the base classes which includes java.lang.String
	LoadBaseClasses()

	className := "java/lang/String"

	// Fetch the class from the Method Area
	k := MethAreaFetch(className)
	if k == nil {
		t.Fatalf("Class %s not found in MethArea after LoadBaseClasses", className)
	}

	// Retrieve the ClassObject (the JLC instance) from the class data
	classObj := k.Data.ClassObject
	if classObj == nil {
		t.Fatalf("ClassObject for %s is nil", className)
	}

	// Verify the "name" field contains a String object with the correct class name
	classObj.ThMutex.RLock()
	nameField, ok := classObj.FieldTable["name"]
	classObj.ThMutex.RUnlock()

	if !ok {
		t.Fatal("Class object missing 'name' field")
	}

	// Note: in convertToPostableClass, the 'name' field is set as a String object
	nameObj, ok := nameField.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("Class object 'name' field is not an *object.Object, got %T", nameField.Fvalue)
	}

	actualName := object.GoStringFromStringObject(nameObj)
	if actualName != className {
		t.Errorf("Expected Class object name to be '%s', got '%s'", className, actualName)
	}

	// Verify that the $statics slice contains expected static fields from java.lang.String
	// serialVersionUID is universally present in Serializable classes like String.
	expectedStatic := "serialVersionUIDJ" // Name + Descriptor

	classObj.ThMutex.RLock()
	staticsField, ok := classObj.FieldTable["$statics"]
	classObj.ThMutex.RUnlock()

	if !ok {
		t.Fatal("Class object missing '$statics' field")
	}

	staticsSlice, ok := staticsField.Fvalue.([]string)
	if !ok {
		t.Fatalf("Class object '$statics' field is not a []string, got %T", staticsField.Fvalue)
	}

	found := false
	for _, s := range staticsSlice {
		if s == expectedStatic {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected static field '%s' not found in java.lang.String $statics array", expectedStatic)
	}

	// Verify the $klass pointer points back to the metadata
	classObj.ThMutex.RLock()
	klassField, ok := classObj.FieldTable["$klass"]
	classObj.ThMutex.RUnlock()

	if !ok {
		t.Fatal("Class object missing '$klass' field")
	}

	klassDataPtr, ok := klassField.Fvalue.(*ClData)
	if !ok {
		t.Fatalf("Class object '$klass' field is not a *ClData, got %T", klassField.Fvalue)
	}

	if klassDataPtr.Name != className {
		t.Errorf("Expected $klass pointer to point to ClData for '%s', got '%s'", className, klassDataPtr.Name)
	}
}

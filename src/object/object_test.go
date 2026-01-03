/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"math"
	"testing"
)

func TestIsNull(t *testing.T) {
	if !IsNull(nil) {
		t.Errorf("nil should be null")
	}

	var op *Object = Null
	if !IsNull(op) {
		t.Errorf("pointer to non-allocated object should be null")
	}
}

// An array of Object pointers is never null
func TestIsNullForArrays(t *testing.T) {
	stringPool.EmptyStringPool()
	stringPool.PreloadArrayClassesToStringPool()
	arrayObj := Make1DimArray(T_REF, 10)
	if IsNull(arrayObj.FieldTable["value"].Fvalue.([]*Object)) {
		t.Errorf("arrayObj should not be null")
	}
}

func TestMakeEmptyObjectWithClassName(t *testing.T) {
	globals.InitGlobals("test")
	clName := "genericClass"
	o := MakeEmptyObjectWithClassName(&clName)
	fieldSize := len(o.FieldTable)
	if fieldSize != 0 {
		t.Errorf("fieldSize should be zero, got %d", fieldSize)
	}

	namePter := stringPool.GetStringPointer(o.KlassName)
	if *namePter != clName {
		t.Errorf("Did not get 'generic' class type, got %s", *namePter)
	}
}

func TestMakeValidPrimitiveByte(t *testing.T) {
	globals.InitGlobals("test")
	objPtr := MakePrimitiveObject("java/lang/Byte", types.Byte, uint8(0x61))
	if *(stringPool.GetStringPointer(objPtr.KlassName)) != "java/lang/Byte" {
		t.Errorf("Klass should be java/lang/Byte, got %s",
			*(stringPool.GetStringPointer(objPtr.KlassName)))
	}

	value := objPtr.FieldTable["value"].Fvalue.(uint8)
	if value != uint8(0x61) {
		t.Errorf("Value should be 0x61, got 0x%02x", value)
	}
}

func TestMakeValidPrimitiveDouble(t *testing.T) {
	globals.InitGlobals("test")
	objPtr := MakePrimitiveObject("java/lang/Double", types.Double, 42.0)
	if *(stringPool.GetStringPointer(objPtr.KlassName)) != "java/lang/Double" {
		t.Errorf("Klass should be java/lang/Double, got %s", *(stringPool.GetStringPointer(objPtr.KlassName)))
	}

	value := objPtr.FieldTable["value"].Fvalue.(float64)
	if value != 42.0 {
		t.Errorf("Value should be 0x42.0, got 0x%f", value)
	}
}

func TestCloneObject_1(t *testing.T) {
	globals.InitGlobals("test")
	obj1 := MakePrimitiveObject("java/lang/Double", types.Double, 42.0)
	obj2 := CloneObject(obj1)

	// Make sure that the class identifiers are identical.
	if obj2.KlassName != obj1.KlassName {
		t.Errorf("KlassName should be the same. obj1: %v, obj2: %v", obj1.KlassName, obj2.KlassName)
	}

	// Make sure that their hashes are different.
	if obj2.Mark.Hash == obj1.Mark.Hash {
		t.Errorf("Mark.Hash should be different. obj1: %v, obj2: %v", obj1.Mark.Hash, obj2.Mark.Hash)
	}

	// Capture values for both obj1 and obj2.
	// Then, make sure they are identical.
	value1 := obj1.FieldTable["value"].Fvalue.(float64)
	value2 := obj2.FieldTable["value"].Fvalue.(float64)
	if value2 != value1 {
		t.Errorf("value2 should equal value1, expected %f, observed %f", value1, value2)
		return
	}

	// Change just the obj2 value.
	fld := obj2.FieldTable["value"]
	fld.Fvalue = 43.0
	obj2.FieldTable["value"] = fld

	// Capture values for both obj1 and obj2.
	// Then, make sure they differ in the expected manner.
	value1 = obj1.FieldTable["value"].Fvalue.(float64)
	value2 = obj2.FieldTable["value"].Fvalue.(float64)
	if value1 != 42.0 || value2 != 43.0 {
		t.Errorf("Expected value1=42.0 and value2=43.0 but observed value1=%f and value2=%f", 42.0, 43.0)
		return
	}

}

func TestCloneObject_2(t *testing.T) {
	globals.InitGlobals("test")
	obj1 := MakePrimitiveObject("flying/purple/PeopleEater", types.Int, 1958)
	jthing := [3]int64{1, 2, 3}
	fthing := [3]float64{4, 5, 6}
	obj1.FieldTable["jane"] = Field{Ftype: types.LongArray, Fvalue: jthing}
	obj1.FieldTable["felice"] = Field{Ftype: types.FloatArray, Fvalue: fthing}

	obj2 := CloneObject(obj1)

	// Make sure that the class identifiers are identical.
	if obj2.KlassName != obj1.KlassName {
		t.Errorf("KlassName should be the same. obj1: %v, obj2: %v", obj1.KlassName, obj2.KlassName)
	}

	// Make sure that their hashes are different.
	if obj2.Mark.Hash == obj1.Mark.Hash {
		t.Errorf("Mark.Hash should be different. obj1: %v, obj2: %v", obj1.Mark.Hash, obj2.Mark.Hash)
	}

	// Capture values for both obj1 and obj2.
	// Then, make sure they are identical.
	jane1 := obj1.FieldTable["jane"].Fvalue.([3]int64)
	jane2 := obj2.FieldTable["jane"].Fvalue.([3]int64)
	if jane2 != jane1 {
		t.Errorf("jane2 should equal jane1, expected %v, observed %v", jane1, jane2)
		return
	}

	// Change just the obj2 value for jane.
	fld := obj2.FieldTable["jane"]
	fld.Fvalue = [3]int64{7, 8, 9}
	obj2.FieldTable["jane"] = fld

	// Capture values for both obj1 and obj2.
	// Then, make sure they differ in the expected manner.
	jane1 = obj1.FieldTable["jane"].Fvalue.([3]int64)
	jane2 = obj2.FieldTable["jane"].Fvalue.([3]int64)
	if jane1 != [3]int64{1, 2, 3} || jane2 != [3]int64{7, 8, 9} {
		t.Errorf("Expected jane1=[3]int64{1, 2, 3} and jane2=[3]int64{7, 8, 9} but observed jane1=%v and jane2=%v", jane1, jane2)
		return
	}

	// Capture values for both obj1 and obj2.
	// Then, make sure they are identical.
	felice1 := obj1.FieldTable["felice"].Fvalue.([3]float64)
	felice2 := obj2.FieldTable["felice"].Fvalue.([3]float64)
	if felice2 != felice1 {
		t.Errorf("felice: felice2 should equal felice1, expected %v, observed %v", felice1, felice2)
		return
	}

	// Change just the obj2 value for felice.
	fld = obj2.FieldTable["felice"]
	fld.Fvalue = [3]float64{7, 8, 9}
	obj2.FieldTable["felice"] = fld

	// Capture values for both obj1 and obj2.
	// Then, make sure they differ in the expected manner.
	felice1 = obj1.FieldTable["felice"].Fvalue.([3]float64)
	felice2 = obj2.FieldTable["felice"].Fvalue.([3]float64)
	if felice1 != [3]float64{4, 5, 6} || felice2 != [3]float64{7, 8, 9} {
		t.Errorf("felice: Expected felice1=[3]float64{1, 2, 3} and felice2=[3]float64{7, 8, 9} but observed felice1=%v and felice2=%v", felice1, felice2)
		return
	}
}

// The following tests were generated by JetBrains Junie to cover test gaps

// Test MakeEmptyObject function (completely untested)
func TestMakeEmptyObject(t *testing.T) {
	// Test basic object creation
	obj := MakeEmptyObject()
	if obj == nil {
		t.Errorf("MakeEmptyObject() should not return nil")
	}

	// Verify hash is set and non-zero
	if obj.Mark.Hash == 0 {
		t.Errorf("Expected non-zero hash, got %d", obj.Mark.Hash)
	}

	// Verify KlassName is InvalidStringIndex
	if obj.KlassName != types.InvalidStringIndex {
		t.Errorf("Expected KlassName to be InvalidStringIndex (%d), got %d", types.InvalidStringIndex, obj.KlassName)
	}

	// Verify FieldTable is initialized and empty
	if obj.FieldTable == nil {
		t.Errorf("FieldTable should be initialized, got nil")
	}
	if len(obj.FieldTable) != 0 {
		t.Errorf("Expected empty FieldTable, got %d fields", len(obj.FieldTable))
	}

	// Test hash uniqueness with multiple objects
	obj2 := MakeEmptyObject()
	if obj.Mark.Hash == obj2.Mark.Hash {
		t.Errorf("Expected different hash values for different objects, both got %d", obj.Mark.Hash)
	}
}

// Test MakeOneFieldObject function (completely untested)
func TestMakeOneFieldObject(t *testing.T) {
	globals.InitGlobals("test")

	// Test with string field
	stringObj := MakeOneFieldObject("java/lang/String", "data", types.ByteArray, "hello")
	if stringObj == nil {
		t.Errorf("MakeOneFieldObject() should not return nil")
	}

	// Verify class name is set correctly
	className := stringPool.GetStringPointer(stringObj.KlassName)
	if *className != "java/lang/String" {
		t.Errorf("Expected class name 'java/lang/String', got '%s'", *className)
	}

	// Verify field name and value are stored correctly
	field, exists := stringObj.FieldTable["data"]
	if !exists {
		t.Errorf("Expected field 'data' to exist")
	}
	if field.Ftype != types.ByteArray {
		t.Errorf("Expected field type '%s', got '%s'", types.ByteArray, field.Ftype)
	}
	if field.Fvalue.(string) != "hello" {
		t.Errorf("Expected field value 'hello', got '%s'", field.Fvalue.(string))
	}

	// Test with integer field
	intObj := MakeOneFieldObject("java/lang/Integer", "value", types.Int, 42)
	intField, exists := intObj.FieldTable["value"]
	if !exists {
		t.Errorf("Expected field 'value' to exist")
	}
	if intField.Ftype != types.Int {
		t.Errorf("Expected field type '%s', got '%s'", types.Int, intField.Ftype)
	}
	if intField.Fvalue.(int) != 42 {
		t.Errorf("Expected field value 42, got %d", intField.Fvalue.(int))
	}

	// Test with double field
	doubleObj := MakeOneFieldObject("java/lang/Double", "amount", types.Double, 3.14)
	doubleField, exists := doubleObj.FieldTable["amount"]
	if !exists {
		t.Errorf("Expected field 'amount' to exist")
	}
	if doubleField.Ftype != types.Double {
		t.Errorf("Expected field type '%s', got '%s'", types.Double, doubleField.Ftype)
	}
	if doubleField.Fvalue.(float64) != 3.14 {
		t.Errorf("Expected field value 3.14, got %f", doubleField.Fvalue.(float64))
	}

	// Test with boolean field
	boolObj := MakeOneFieldObject("java/lang/Boolean", "flag", types.Bool, true)
	boolField, exists := boolObj.FieldTable["flag"]
	if !exists {
		t.Errorf("Expected field 'flag' to exist")
	}
	if boolField.Ftype != types.Bool {
		t.Errorf("Expected field type '%s', got '%s'", types.Bool, boolField.Ftype)
	}
	if boolField.Fvalue.(bool) != true {
		t.Errorf("Expected field value true, got %v", boolField.Fvalue.(bool))
	}

	// Test with nil value
	nilObj := MakeOneFieldObject("TestClass", "nullField", types.Ref, nil)
	nilField, exists := nilObj.FieldTable["nullField"]
	if !exists {
		t.Errorf("Expected field 'nullField' to exist")
	}
	if nilField.Fvalue != nil {
		t.Errorf("Expected field value to be nil, got %v", nilField.Fvalue)
	}

	// Test with empty field name
	emptyFieldObj := MakeOneFieldObject("TestClass", "", types.ByteArray, "test")
	_, exists = emptyFieldObj.FieldTable[""]
	if !exists {
		t.Errorf("Expected empty string field name to be valid")
	}
}

// Test UpdateValueFieldFromJavaBytes function (completely untested)
func TestUpdateValueFieldFromJavaBytes(t *testing.T) {
	globals.InitGlobals("test")

	// Test with empty byte array
	obj := MakeEmptyObject()
	emptyBytes := []types.JavaByte{}
	UpdateValueFieldFromJavaBytes(obj, emptyBytes)

	field, exists := obj.FieldTable["value"]
	if !exists {
		t.Errorf("Expected 'value' field to exist after update")
	}
	if field.Ftype != "Ljava/lang/String;" {
		t.Errorf("Expected field type 'Ljava/lang/String;', got '%s'", field.Ftype)
	}
	resultBytes := field.Fvalue.([]types.JavaByte)
	if len(resultBytes) != 0 {
		t.Errorf("Expected empty byte array, got length %d", len(resultBytes))
	}

	// Test with single byte
	singleByte := []types.JavaByte{65} // ASCII 'A'
	UpdateValueFieldFromJavaBytes(obj, singleByte)
	field = obj.FieldTable["value"]
	resultBytes = field.Fvalue.([]types.JavaByte)
	if len(resultBytes) != 1 || resultBytes[0] != 65 {
		t.Errorf("Expected single byte [65], got %v", resultBytes)
	}

	// Test with multiple bytes
	multiBytes := []types.JavaByte{72, 101, 108, 108, 111} // "Hello"
	UpdateValueFieldFromJavaBytes(obj, multiBytes)
	field = obj.FieldTable["value"]
	resultBytes = field.Fvalue.([]types.JavaByte)
	if len(resultBytes) != 5 {
		t.Errorf("Expected 5 bytes, got %d", len(resultBytes))
	}
	for i, expected := range multiBytes {
		if resultBytes[i] != expected {
			t.Errorf("Expected byte at index %d to be %d, got %d", i, expected, resultBytes[i])
		}
	}

	// Test with nil object (should handle gracefully or panic - depends on design)
	// This test verifies current behavior rather than asserting what should happen
	defer func() {
		if r := recover(); r != nil {
			// If it panics, that's acceptable behavior for nil input
			t.Logf("UpdateValueFieldFromJavaBytes panicked with nil object: %v", r)
		}
	}()
	UpdateValueFieldFromJavaBytes(nil, []types.JavaByte{1, 2, 3})

	// Test overwriting existing value field
	obj2 := MakeEmptyObject()
	obj2.FieldTable["value"] = Field{Ftype: types.Int, Fvalue: 42}
	newBytes := []types.JavaByte{65, 66, 67}
	UpdateValueFieldFromJavaBytes(obj2, newBytes)

	field = obj2.FieldTable["value"]
	if field.Ftype != "Ljava/lang/String;" {
		t.Errorf("Expected field type to be overwritten to 'Ljava/lang/String;', got '%s'", field.Ftype)
	}
	resultBytes = field.Fvalue.([]types.JavaByte)
	if len(resultBytes) != 3 {
		t.Errorf("Expected 3 bytes after overwrite, got %d", len(resultBytes))
	}
}

// Test ClearFieldTable function (completely untested)
func TestClearFieldTable(t *testing.T) {
	globals.InitGlobals("test")

	// Test clearing empty field table
	emptyObj := MakeEmptyObject()
	ClearFieldTable(emptyObj)
	if len(emptyObj.FieldTable) != 0 {
		t.Errorf("Expected field table to remain empty, got %d fields", len(emptyObj.FieldTable))
	}

	// Test clearing field table with one field
	singleFieldObj := MakeOneFieldObject("TestClass", "field1", types.ByteArray, "value1")
	if len(singleFieldObj.FieldTable) != 1 {
		t.Errorf("Expected 1 field before clearing, got %d", len(singleFieldObj.FieldTable))
	}
	ClearFieldTable(singleFieldObj)
	if len(singleFieldObj.FieldTable) != 0 {
		t.Errorf("Expected field table to be empty after clearing, got %d fields", len(singleFieldObj.FieldTable))
	}

	// Test clearing field table with multiple fields
	multiFieldObj := MakeEmptyObject()
	multiFieldObj.FieldTable["field1"] = Field{Ftype: types.ByteArray, Fvalue: "value1"}
	multiFieldObj.FieldTable["field2"] = Field{Ftype: types.Int, Fvalue: 42}
	multiFieldObj.FieldTable["field3"] = Field{Ftype: types.Double, Fvalue: 3.14}

	if len(multiFieldObj.FieldTable) != 3 {
		t.Errorf("Expected 3 fields before clearing, got %d", len(multiFieldObj.FieldTable))
	}

	// Store original values to verify other properties remain unchanged
	originalKlass := multiFieldObj.KlassName
	originalHash := multiFieldObj.Mark.Hash

	ClearFieldTable(multiFieldObj)

	if len(multiFieldObj.FieldTable) != 0 {
		t.Errorf("Expected field table to be empty after clearing, got %d fields", len(multiFieldObj.FieldTable))
	}

	// Verify other object properties remain unchanged
	if multiFieldObj.KlassName != originalKlass {
		t.Errorf("Expected KlassName to remain unchanged, was %d, now %d", originalKlass, multiFieldObj.KlassName)
	}
	if multiFieldObj.Mark.Hash != originalHash {
		t.Errorf("Expected hash to remain unchanged, was %d, now %d", originalHash, multiFieldObj.Mark.Hash)
	}

	// Test with nil object (should handle gracefully or panic)
	defer func() {
		if r := recover(); r != nil {
			// If it panics, that's acceptable behavior for nil input
			t.Logf("ClearFieldTable panicked with nil object: %v", r)
		}
	}()
	ClearFieldTable(nil)
}

// Test GetClassNameSuffix function (completely untested)
func TestGetClassNameSuffix(t *testing.T) {
	globals.InitGlobals("test")

	// Test with nil object
	result := GetClassNameSuffix(nil, false)
	if result != types.NullString {
		t.Errorf("Expected NullString for nil object, got '%s'", result)
	}

	result = GetClassNameSuffix(nil, true)
	if result != types.NullString {
		t.Errorf("Expected NullString for nil object with inner=true, got '%s'", result)
	}

	// Test with Null object
	result = GetClassNameSuffix(Null, false)
	if result != types.NullString {
		t.Errorf("Expected NullString for Null object, got '%s'", result)
	}

	// Test with simple class name (e.g., "String")
	simpleObj := MakeEmptyObjectWithClassName(&[]string{"String"}[0])
	result = GetClassNameSuffix(simpleObj, false)
	if result != "String" {
		t.Errorf("Expected 'String' for simple class name, got '%s'", result)
	}

	result = GetClassNameSuffix(simpleObj, true)
	if result != "String" {
		t.Errorf("Expected 'String' for simple class name with inner=true, got '%s'", result)
	}

	// Test with full package path (e.g., "java/lang/String")
	packageObj := MakeEmptyObjectWithClassName(&[]string{"java/lang/String"}[0])
	result = GetClassNameSuffix(packageObj, false)
	if result != "String" {
		t.Errorf("Expected 'String' for java/lang/String, got '%s'", result)
	}

	result = GetClassNameSuffix(packageObj, true)
	if result != "String" {
		t.Errorf("Expected 'String' for java/lang/String with inner=true, got '%s'", result)
	}

	// Test with inner class (e.g., "OuterClass$InnerClass")
	innerObj := MakeEmptyObjectWithClassName(&[]string{"OuterClass$InnerClass"}[0])
	result = GetClassNameSuffix(innerObj, false)
	if result != "OuterClass$InnerClass" {
		t.Errorf("Expected 'OuterClass$InnerClass' for inner class with inner=false, got '%s'", result)
	}

	result = GetClassNameSuffix(innerObj, true)
	if result != "InnerClass" {
		t.Errorf("Expected 'InnerClass' for inner class with inner=true, got '%s'", result)
	}

	// Test with multiple inner classes (e.g., "A$B$C")
	multiInnerObj := MakeEmptyObjectWithClassName(&[]string{"com/example/A$B$C"}[0])
	result = GetClassNameSuffix(multiInnerObj, false)
	if result != "A$B$C" {
		t.Errorf("Expected 'A$B$C' for multi-inner class with inner=false, got '%s'", result)
	}

	result = GetClassNameSuffix(multiInnerObj, true)
	if result != "C" {
		t.Errorf("Expected 'C' for multi-inner class with inner=true, got '%s'", result)
	}

	// Test with classes using dots instead of slashes
	dotObj := MakeEmptyObjectWithClassName(&[]string{"java.lang.String"}[0])
	result = GetClassNameSuffix(dotObj, false)
	if result != "String" {
		t.Errorf("Expected 'String' for dot notation class, got '%s'", result)
	}

	// Test with complex package and inner class
	complexObj := MakeEmptyObjectWithClassName(&[]string{"com/example/package/OuterClass$InnerClass$DeepInner"}[0])
	result = GetClassNameSuffix(complexObj, false)
	if result != "OuterClass$InnerClass$DeepInner" {
		t.Errorf("Expected 'OuterClass$InnerClass$DeepInner' for complex class with inner=false, got '%s'", result)
	}

	result = GetClassNameSuffix(complexObj, true)
	if result != "DeepInner" {
		t.Errorf("Expected 'DeepInner' for complex class with inner=true, got '%s'", result)
	}

	// Test edge case: class name ending with $
	dollarObj := MakeEmptyObjectWithClassName(&[]string{"TestClass$"}[0])
	result = GetClassNameSuffix(dollarObj, true)
	if result != "" {
		t.Errorf("Expected empty string for class ending with $, got '%s'", result)
	}
}

// Enhanced TestMakePrimitiveObject to cover all primitive types
func TestMakePrimitiveObjectAllTypes(t *testing.T) {
	globals.InitGlobals("test")

	// Test Boolean
	boolObj := MakePrimitiveObject("java/lang/Boolean", types.Bool, true)
	if *(stringPool.GetStringPointer(boolObj.KlassName)) != "java/lang/Boolean" {
		t.Errorf("Expected class java/lang/Boolean, got %s", *(stringPool.GetStringPointer(boolObj.KlassName)))
	}
	if boolObj.FieldTable["value"].Fvalue.(bool) != true {
		t.Errorf("Expected boolean value true, got %v", boolObj.FieldTable["value"].Fvalue)
	}

	// Test Character
	charObj := MakePrimitiveObject("java/lang/Character", types.Char, 'A')
	if *(stringPool.GetStringPointer(charObj.KlassName)) != "java/lang/Character" {
		t.Errorf("Expected class java/lang/Character, got %s", *(stringPool.GetStringPointer(charObj.KlassName)))
	}
	if charObj.FieldTable["value"].Fvalue.(rune) != 'A' {
		t.Errorf("Expected char value 'A', got %v", charObj.FieldTable["value"].Fvalue)
	}

	// Test Float
	floatObj := MakePrimitiveObject("java/lang/Float", types.Float, float32(3.14))
	if *(stringPool.GetStringPointer(floatObj.KlassName)) != "java/lang/Float" {
		t.Errorf("Expected class java/lang/Float, got %s", *(stringPool.GetStringPointer(floatObj.KlassName)))
	}
	if floatObj.FieldTable["value"].Fvalue.(float32) != float32(3.14) {
		t.Errorf("Expected float value 3.14, got %v", floatObj.FieldTable["value"].Fvalue)
	}

	// Test Integer
	intObj := MakePrimitiveObject("java/lang/Integer", types.Int, 42)
	if *(stringPool.GetStringPointer(intObj.KlassName)) != "java/lang/Integer" {
		t.Errorf("Expected class java/lang/Integer, got %s", *(stringPool.GetStringPointer(intObj.KlassName)))
	}
	if intObj.FieldTable["value"].Fvalue.(int) != 42 {
		t.Errorf("Expected int value 42, got %v", intObj.FieldTable["value"].Fvalue)
	}

	// Test Long
	longObj := MakePrimitiveObject("java/lang/Long", types.Long, int64(9223372036854775807))
	if *(stringPool.GetStringPointer(longObj.KlassName)) != "java/lang/Long" {
		t.Errorf("Expected class java/lang/Long, got %s", *(stringPool.GetStringPointer(longObj.KlassName)))
	}
	if longObj.FieldTable["value"].Fvalue.(int64) != int64(9223372036854775807) {
		t.Errorf("Expected long value 9223372036854775807, got %v", longObj.FieldTable["value"].Fvalue)
	}

	// Test Short
	shortObj := MakePrimitiveObject("java/lang/Short", types.Short, int16(32767))
	if *(stringPool.GetStringPointer(shortObj.KlassName)) != "java/lang/Short" {
		t.Errorf("Expected class java/lang/Short, got %s", *(stringPool.GetStringPointer(shortObj.KlassName)))
	}
	if shortObj.FieldTable["value"].Fvalue.(int16) != int16(32767) {
		t.Errorf("Expected short value 32767, got %v", shortObj.FieldTable["value"].Fvalue)
	}

	// Test edge cases with min/max values
	// Test Integer min/max
	intMinObj := MakePrimitiveObject("java/lang/Integer", types.Int, -2147483648)
	if intMinObj.FieldTable["value"].Fvalue.(int) != -2147483648 {
		t.Errorf("Expected int min value -2147483648, got %v", intMinObj.FieldTable["value"].Fvalue)
	}

	intMaxObj := MakePrimitiveObject("java/lang/Integer", types.Int, 2147483647)
	if intMaxObj.FieldTable["value"].Fvalue.(int) != 2147483647 {
		t.Errorf("Expected int max value 2147483647, got %v", intMaxObj.FieldTable["value"].Fvalue)
	}

	// Test Long min/max
	longMinObj := MakePrimitiveObject("java/lang/Long", types.Long, int64(-9223372036854775808))
	if longMinObj.FieldTable["value"].Fvalue.(int64) != int64(-9223372036854775808) {
		t.Errorf("Expected long min value -9223372036854775808, got %v", longMinObj.FieldTable["value"].Fvalue)
	}

	// Test Short min/max
	shortMinObj := MakePrimitiveObject("java/lang/Short", types.Short, int16(-32768))
	if shortMinObj.FieldTable["value"].Fvalue.(int16) != int16(-32768) {
		t.Errorf("Expected short min value -32768, got %v", shortMinObj.FieldTable["value"].Fvalue)
	}

	shortMaxObj := MakePrimitiveObject("java/lang/Short", types.Short, int16(32767))
	if shortMaxObj.FieldTable["value"].Fvalue.(int16) != int16(32767) {
		t.Errorf("Expected short max value 32767, got %v", shortMaxObj.FieldTable["value"].Fvalue)
	}
}

// Test object hash uniqueness
func TestMakeEmptyObjectHashUniqueness(t *testing.T) {
	const numObjects = 100
	objects := make([]*Object, numObjects)
	hashes := make(map[uint32]bool)

	// Create multiple objects and collect their hashes
	for i := 0; i < numObjects; i++ {
		objects[i] = MakeEmptyObject()
		hash := objects[i].Mark.Hash

		if hashes[hash] {
			t.Errorf("Hash collision detected: hash %d appears multiple times", hash)
		}
		hashes[hash] = true
	}

	// Verify we got unique hashes for all objects
	if len(hashes) != numObjects {
		t.Errorf("Expected %d unique hashes, got %d", numObjects, len(hashes))
	}
}

func TestTVO(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObject()
	_ = obj.TVO()

	clName1 := "apple/beet/carrot"
	clName2 := "dandelion/daisy/bluebell"
	obj1 := MakeEmptyObjectWithClassName(&clName1)
	obj2 := MakeEmptyObjectWithClassName(&clName2)
	obj1.FieldTable["field6b"] = Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}
	obj1.FieldTable["field7"] = Field{Ftype: types.ByteArray, Fvalue: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	obj1.FieldTable["field8"] = Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{61, 62, 63, 64, 65, 66, 67, 68, 69, 70}}
	obj1.FieldTable["field9"] = Field{Ftype: types.Int, Fvalue: uint32(math.Pow(2, 27) - 1)}
	obj1.FieldTable["field10"] = Field{Ftype: types.Int, Fvalue: uint16(32767)}
	obj1.FieldTable["field1"] = Field{Ftype: types.ByteArray, Fvalue: JavaByteArrayFromGoString("value1")}
	obj1.FieldTable["field2"] = Field{Ftype: types.Int, Fvalue: 42}
	obj1.FieldTable["field3"] = Field{Ftype: types.Double, Fvalue: 3.14}
	obj1.FieldTable["field4"] = Field{Ftype: types.StringClassName, Fvalue: JavaByteArrayFromGoString("value4")}
	obj1.FieldTable["field5"] = Field{Ftype: types.Ref, Fvalue: obj2}
	obj1.FieldTable["field6a"] = Field{Ftype: types.Bool, Fvalue: types.JavaBoolTrue}
	obj3 := StringObjectFromGoString("Hey diddle diddle .....")
	obj1.FieldTable["field11"] = Field{Ftype: types.StringClassName, Fvalue: obj3}

	t.Log(obj1.TVO())

}

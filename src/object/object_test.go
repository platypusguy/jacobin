/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/globals"
	"jacobin/stringPool"
	"jacobin/types"
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

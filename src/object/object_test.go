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

	var op *Object
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

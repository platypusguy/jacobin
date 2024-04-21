/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/globals"
	"jacobin/object"
	"jacobin/stringPool"
	"strings"
	"testing"
)

func TestArrayCopyNonOverlapping(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err != nil {
		e := err.(error)
		t.Errorf("Unexpected error in test of arrayCopy(): %s", error.Error(e))
	}

	rawDestArray := dest.FieldTable["value"].Fvalue.([]int64)
	j := int64(0)
	for i := 0; i < 10; i++ {
		j += rawDestArray[i]
	}

	if j != 5 {
		t.Errorf("Expected total to be 5, got %d", j)
	}

	if rawDestArray[0] != 1 || rawDestArray[5] != 0 {
		t.Errorf("Expedting [0] to be 1, [5] to be 0, got %d, %d",
			rawDestArray[0], rawDestArray[5])
	}
}

func TestArrayInvalidParmCount(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 4)
	params[0] = src
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	// params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Expecting error, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Expected 5 parameters") {
		t.Errorf("Expected error re 5 parameters, got %s", errMsg)
	}
}

func TestArrayCopyInvalidPos(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(-1) // this is an invalid position in the array
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Negative position") {
		t.Errorf("Expected error re invalid position, got %s", errMsg)
	}
}

func TestArrayCopyNullArray(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = object.Null // clearly invalid
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "null src or dest") {
		t.Errorf("Expected error re null array, got %s", errMsg)
	}
}

func TestArrayCopyInvalidObject(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	o := object.MakeEmptyObject()
	objType := "invalid object"
	o.KlassName = stringPool.GetStringIndex(&objType)
	params[0] = o
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "invalid src or dest array") {
		t.Errorf("Expected error re invalid array type, got %s", errMsg)
	}
}

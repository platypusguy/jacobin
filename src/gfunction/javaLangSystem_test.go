/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/globals"
	"jacobin/object"
	"testing"
)

func TestStringArrayCopyNonOverlapping(t *testing.T) {
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
}

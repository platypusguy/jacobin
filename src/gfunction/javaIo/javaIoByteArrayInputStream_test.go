/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"container/list"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func setupBAISTest() {
	globals.InitGlobals("test")
	g := globals.GetGlobalRef()
	g.FuncInstantiateClass = func(classname string, frameStack *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		obj.FieldTable = make(map[string]object.Field)
		return obj, nil
	}
}

func TestByteArrayInputStream_Basic(t *testing.T) {
	setupBAISTest()

	data := []types.JavaByte{10, 20, 30, 40, 50}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	params := []interface{}{bais, bufObj}
	ByteArrayInputStreamInit(params)

	// available
	avail := ByteArrayInputStreamAvailable([]interface{}{bais})
	if avail.(int64) != 5 {
		t.Errorf("Expected available 5, got %d", avail)
	}

	// read one by one
	for i := 0; i < 5; i++ {
		res := ByteArrayInputStreamRead([]interface{}{bais})
		if res.(int64) != int64(data[i]) {
			t.Errorf("At index %d, expected %d, got %d", i, data[i], res)
		}
	}

	// EOF
	res := ByteArrayInputStreamRead([]interface{}{bais})
	if res.(int64) != -1 {
		t.Errorf("Expected -1 at EOF, got %d", res)
	}

	avail = ByteArrayInputStreamAvailable([]interface{}{bais})
	if avail.(int64) != 0 {
		t.Errorf("Expected available 0 at EOF, got %d", avail)
	}
}

func TestByteArrayInputStream_Range(t *testing.T) {
	setupBAISTest()

	data := []types.JavaByte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	// Init with offset 2, length 5 (elements 3, 4, 5, 6, 7)
	bais := object.MakeEmptyObject()
	params := []interface{}{bais, bufObj, int64(2), int64(5)}
	ByteArrayInputStreamInit(params)

	avail := ByteArrayInputStreamAvailable([]interface{}{bais})
	if avail.(int64) != 5 {
		t.Errorf("Expected available 5, got %d", avail)
	}

	dest := make([]types.JavaByte, 10)
	destObj := object.MakeEmptyObject()
	destObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: dest}

	// read 3 bytes
	readRes := ByteArrayInputStreamReadRange([]interface{}{bais, destObj, int64(1), int64(3)})
	if readRes.(int64) != 3 {
		t.Errorf("Expected to read 3 bytes, got %d", readRes)
	}

	expectedDest := []types.JavaByte{0, 3, 4, 5, 0, 0, 0, 0, 0, 0}
	for i, v := range expectedDest {
		if dest[i] != v {
			t.Errorf("At index %d, expected %d, got %d", i, v, dest[i])
		}
	}

	// read remaining (should be 2)
	readRes = ByteArrayInputStreamReadRange([]interface{}{bais, destObj, int64(5), int64(2)})
	if readRes.(int64) != 2 {
		t.Errorf("Expected to read 2 remaining bytes, got %d", readRes)
	}
	if dest[5] != 6 || dest[6] != 7 {
		t.Errorf("Unexpected values in destination: %d, %d", dest[5], dest[6])
	}

	// read EOF
	readRes = ByteArrayInputStreamReadRange([]interface{}{bais, destObj, int64(0), int64(1)})
	if readRes.(int64) != -1 {
		t.Errorf("Expected -1 at EOF, got %d", readRes)
	}
}

func TestByteArrayInputStream_MarkReset(t *testing.T) {
	setupBAISTest()

	data := []types.JavaByte{100, 101, 102, 103, 104}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	// Read 2
	ByteArrayInputStreamRead([]interface{}{bais})
	ByteArrayInputStreamRead([]interface{}{bais})

	// Mark at pos 2
	ByteArrayInputStreamMark([]interface{}{bais, int64(100)})

	// Read 2 more
	r1 := ByteArrayInputStreamRead([]interface{}{bais})
	r2 := ByteArrayInputStreamRead([]interface{}{bais})
	if r1.(int64) != 102 || r2.(int64) != 103 {
		t.Errorf("Unexpected read values: %d, %d", r1, r2)
	}

	// Reset
	ByteArrayInputStreamReset([]interface{}{bais})

	// Read again from mark
	r1 = ByteArrayInputStreamRead([]interface{}{bais})
	if r1.(int64) != 102 {
		t.Errorf("Expected 102 after reset, got %d", r1)
	}
}

func TestByteArrayInputStream_Skip(t *testing.T) {
	setupBAISTest()

	data := []types.JavaByte{1, 2, 3, 4, 5}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	// Skip 2
	skipped := ByteArrayInputStreamSkip([]interface{}{bais, int64(2)})
	if skipped.(int64) != 2 {
		t.Errorf("Expected skipped 2, got %d", skipped)
	}

	// Read 1
	res := ByteArrayInputStreamRead([]interface{}{bais})
	if res.(int64) != 3 {
		t.Errorf("Expected 3 after skipping 2, got %d", res)
	}

	// Skip more than available
	skipped = ByteArrayInputStreamSkip([]interface{}{bais, int64(10)})
	if skipped.(int64) != 2 { // only 4, 5 left
		t.Errorf("Expected skipped 2 more, got %d", skipped)
	}

	// EOF
	res = ByteArrayInputStreamRead([]interface{}{bais})
	if res.(int64) != -1 {
		t.Errorf("Expected -1 at EOF, got %d", res)
	}
}

func TestByteArrayInputStream_ReadAllBytes(t *testing.T) {
	setupBAISTest()

	data := []types.JavaByte{10, 20, 30, 40, 50}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	// Read 1
	ByteArrayInputStreamRead([]interface{}{bais})

	// Read All remaining (20, 30, 40, 50)
	res := ByteArrayInputStreamReadAllBytes([]interface{}{bais})
	resObj := res.(*object.Object)
	resData := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)

	if len(resData) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(resData))
	}
	if resData[0] != 20 || resData[3] != 50 {
		t.Errorf("Unexpected data: %v", resData)
	}

	// Read All again (empty)
	res = ByteArrayInputStreamReadAllBytes([]interface{}{bais})
	resObj = res.(*object.Object)
	resData = resObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(resData) != 0 {
		t.Errorf("Expected 0 bytes, got %d", len(resData))
	}
}

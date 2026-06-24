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
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"testing"
)

func setupBufferedInputStreamTest() {
	globals.InitGlobals("test")
	Load_Io_ByteArrayInputStream()
	Load_Io_BufferedInputStream()
	g := globals.GetGlobalRef()
	g.FuncInstantiateClass = func(classname string, frameStack *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		obj.FieldTable = make(map[string]object.Field)
		obj.KlassName = stringPool.GetStringIndex(&classname)
		return obj, nil
	}
	globals.InitStringPool()
}

func TestBufferedInputStream_Basic(t *testing.T) {
	setupBufferedInputStreamTest()

	// 1. Setup underlying ByteArrayInputStream
	data := []types.JavaByte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	nameBAIS := "java/io/ByteArrayInputStream"
	bais.KlassName = stringPool.GetStringIndex(&nameBAIS)
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	// 2. Setup BufferedInputStream wrapping bais with small buffer (size 3)
	bis := object.MakeEmptyObject()
	nameBIS := "java/io/BufferedInputStream"
	bis.KlassName = stringPool.GetStringIndex(&nameBIS)
	if res := BufferedInputStreamInit([]interface{}{bis, bais, int64(3)}); res != nil {
		t.Fatalf("Init failed: %v", res)
	}

	// 3. Read first byte (should trigger fill)
	b1 := BufferedInputStreamRead([]interface{}{bis})
	if b1.(int64) != 1 {
		t.Errorf("Expected 1, got %d", b1)
	}

	// count should be 3 now (initial fill)
	count := bis.FieldTable["count"].Fvalue.(int64)
	if count != 3 {
		t.Errorf("After read(1), expected count 3, got %d", count)
	}

	// 4. Read remaining bytes in first buffer
	b2 := BufferedInputStreamRead([]interface{}{bis})
	b3 := BufferedInputStreamRead([]interface{}{bis})
	if b2.(int64) != 2 || b3.(int64) != 3 {
		t.Errorf("Expected 2, 3, got %d, %d", b2, b3)
	}

	// 5. Read next byte (should trigger second fill)
	b4 := BufferedInputStreamRead([]interface{}{bis})
	if b4.(int64) != 4 {
		t.Errorf("Expected 4, got %d", b4)
	}
	// count should be 3 now (it's reset to 0 and refilled with 3 bytes)
	count = bis.FieldTable["count"].Fvalue.(int64)
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestBufferedInputStream_ReadRange(t *testing.T) {
	setupBufferedInputStreamTest()

	data := make([]types.JavaByte, 20)
	for i := range 20 {
		data[i] = types.JavaByte(i + 1)
	}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	nameBAIS := "java/io/ByteArrayInputStream"
	bais.KlassName = stringPool.GetStringIndex(&nameBAIS)
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	bis := object.MakeEmptyObject()
	nameBIS := "java/io/BufferedInputStream"
	bis.KlassName = stringPool.GetStringIndex(&nameBIS)
	BufferedInputStreamInit([]interface{}{bis, bais, int64(5)})

	dest := make([]types.JavaByte, 10)
	destObj := object.MakeEmptyObject()
	destObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: dest}

	// Read 7 bytes (will require multiple fills or direct bypass if optimized, but our implementation fills)
	n := BufferedInputStreamReadRange([]interface{}{bis, destObj, int64(0), int64(7)})
	if n.(int64) != 7 {
		t.Errorf("Expected to read 7 bytes, got %d", n)
	}

	for i := range 7 {
		if dest[i] != types.JavaByte(i+1) {
			t.Errorf("At index %d, expected %d, got %d", i, i+1, dest[i])
		}
	}
}

func TestBufferedInputStream_MarkReset(t *testing.T) {
	setupBufferedInputStreamTest()

	data := []types.JavaByte{1, 2, 3, 4, 5, 6, 7, 8}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	nameBAIS := "java/io/ByteArrayInputStream"
	bais.KlassName = stringPool.GetStringIndex(&nameBAIS)
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	bis := object.MakeEmptyObject()
	nameBIS := "java/io/BufferedInputStream"
	bis.KlassName = stringPool.GetStringIndex(&nameBIS)
	BufferedInputStreamInit([]interface{}{bis, bais, int64(4)})

	// Read 2 bytes: 1, 2
	BufferedInputStreamRead([]interface{}{bis})
	BufferedInputStreamRead([]interface{}{bis})

	// Mark at pos 2, limit 10
	BufferedInputStreamMark([]interface{}{bis, int64(10)})

	// Read 4 more bytes: 3, 4, 5, 6 (triggers fill)
	BufferedInputStreamRead([]interface{}{bis}) // 3
	BufferedInputStreamRead([]interface{}{bis}) // 4
	BufferedInputStreamRead([]interface{}{bis}) // 5 (fill happens here)
	BufferedInputStreamRead([]interface{}{bis}) // 6

	// Reset
	BufferedInputStreamReset([]interface{}{bis})

	// Read next: should be 3
	b := BufferedInputStreamRead([]interface{}{bis})
	if b.(int64) != 3 {
		t.Errorf("Expected 3 after reset, got %d", b)
	}
}

func TestBufferedInputStream_Available(t *testing.T) {
	setupBufferedInputStreamTest()

	data := []types.JavaByte{1, 2, 3, 4, 5}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	nameBAIS := "java/io/ByteArrayInputStream"
	bais.KlassName = stringPool.GetStringIndex(&nameBAIS)
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	bis := object.MakeEmptyObject()
	nameBIS := "java/io/BufferedInputStream"
	bis.KlassName = stringPool.GetStringIndex(&nameBIS)
	BufferedInputStreamInit([]interface{}{bis, bais, int64(10)})

	avail := BufferedInputStreamAvailable([]interface{}{bis})
	if avail.(int64) != 5 {
		t.Errorf("Expected 5 available, got %d", avail)
	}

	// Read 1 (triggers fill of all 5)
	BufferedInputStreamRead([]interface{}{bis})

	avail = BufferedInputStreamAvailable([]interface{}{bis})
	if avail.(int64) != 4 {
		t.Errorf("Expected 4 available, got %d", avail)
	}
}

func TestBufferedInputStream_Skip(t *testing.T) {
	setupBufferedInputStreamTest()

	data := []types.JavaByte{1, 2, 3, 4, 5, 6, 7, 8}
	bufObj := object.MakeEmptyObject()
	bufObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}

	bais := object.MakeEmptyObject()
	nameBAIS := "java/io/ByteArrayInputStream"
	bais.KlassName = stringPool.GetStringIndex(&nameBAIS)
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	bis := object.MakeEmptyObject()
	nameBIS := "java/io/BufferedInputStream"
	bis.KlassName = stringPool.GetStringIndex(&nameBIS)
	BufferedInputStreamInit([]interface{}{bis, bais, int64(4)})

	// Skip 2
	n := BufferedInputStreamSkip([]interface{}{bis, int64(2)})
	if n.(int64) != 2 {
		t.Errorf("Expected skip 2, got %d", n)
	}

	// Read next: should be 3
	b := BufferedInputStreamRead([]interface{}{bis})
	if b.(int64) != 3 {
		t.Errorf("Expected 3, got %d", b)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"bytes"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

func initTestStreams() {
	Load_Io_InputStream()
	Load_Io_ByteArrayInputStream()
	Load_Io_ByteArrayOutputStream()
	Load_Io_OutputStream()
}

func createTestByteArrayInputStream(data []byte) *object.Object {
	initTestStreams()
	jba := make([]types.JavaByte, len(data))
	for i, b := range data {
		jba[i] = types.JavaByte(b)
	}
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, jba)

	bais := object.MakeEmptyObject()
	name := "java/io/ByteArrayInputStream"
	bais.KlassName = stringPool.GetStringIndex(&name)
	ByteArrayInputStreamInit([]any{bais, bufObj})
	return bais
}

func createTestByteArrayOutputStream() *object.Object {
	initTestStreams()
	baos := object.MakeEmptyObject()
	name := "java/io/ByteArrayOutputStream"
	baos.KlassName = stringPool.GetStringIndex(&name)
	ByteArrayOutputStreamInit([]any{baos})
	return baos
}

func getByteArrayOutputStreamData(baos *object.Object) []byte {
	res := ByteArrayOutputStreamToByteArray([]any{baos})
	resObj := res.(*object.Object)
	jba := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	resBytes := make([]byte, len(jba))
	for i, b := range jba {
		resBytes[i] = byte(b)
	}
	return resBytes
}

func TestLoad_Io_InputStream(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Io_InputStream()
	signatures := []string{
		"java/io/InputStream.<clinit>()V",
		"java/io/InputStream.<init>()V",
		"java/io/InputStream.available()I",
		"java/io/InputStream.close()V",
		"java/io/InputStream.mark(I)V",
		"java/io/InputStream.markSupported()Z",
		"java/io/InputStream.nullInputStream()Ljava/io/InputStream;",
		"java/io/InputStream.read()I",
		"java/io/InputStream.read([B)I",
		"java/io/InputStream.read([BII)I",
		"java/io/InputStream.readAllBytes()[B",
		"java/io/InputStream.readNBytes(I)[B",
		"java/io/InputStream.readNBytes([BII)I",
		"java/io/InputStream.reset()V",
		"java/io/InputStream.skip(J)J",
		"java/io/InputStream.skipNBytes(J)V",
		"java/io/InputStream.transferTo(Ljava/io/OutputStream;)J",
	}

	for _, sig := range signatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s to be registered", sig)
		}
	}
}

func TestInputStreamReadAllBytes(t *testing.T) {
	data := []byte("Hello, Jacobin VM! Standard stream tests.")
	bais := createTestByteArrayInputStream(data)

	res := inputStreamReadAllBytes([]any{bais})
	resObj, ok := res.(*object.Object)
	if !ok || resObj == nil {
		t.Fatalf("Expected *object.Object, got %T (%v)", res, res)
	}

	jba := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	outBytes := make([]byte, len(jba))
	for i, b := range jba {
		outBytes[i] = byte(b)
	}

	if !bytes.Equal(outBytes, data) {
		t.Errorf("readAllBytes returned %q, expected %q", outBytes, data)
	}

	// Reading again at EOF should return empty array
	res2 := inputStreamReadAllBytes([]any{bais})
	resObj2 := res2.(*object.Object)
	jba2 := resObj2.FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(jba2) != 0 {
		t.Errorf("Expected empty byte array at EOF, got len %d", len(jba2))
	}
}

func TestInputStreamReadNBytesLen(t *testing.T) {
	data := []byte("0123456789abcdefghij")
	bais := createTestByteArrayInputStream(data)

	// Read 5 bytes
	res5 := inputStreamReadNBytesLen([]any{bais, int64(5)})
	resObj5 := res5.(*object.Object)
	jba5 := resObj5.FieldTable["value"].Fvalue.([]types.JavaByte)
	if string(object.GoByteArrayFromJavaByteArray(jba5)) != "01234" {
		t.Errorf("Expected '01234', got %s", object.GoByteArrayFromJavaByteArray(jba5))
	}

	// Read 0 bytes
	res0 := inputStreamReadNBytesLen([]any{bais, int64(0)})
	resObj0 := res0.(*object.Object)
	jba0 := resObj0.FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(jba0) != 0 {
		t.Errorf("Expected 0 bytes, got %d", len(jba0))
	}

	// Read 100 bytes (more than remaining 15 bytes)
	resMore := inputStreamReadNBytesLen([]any{bais, int64(100)})
	resObjMore := resMore.(*object.Object)
	jbaMore := resObjMore.FieldTable["value"].Fvalue.([]types.JavaByte)
	if string(object.GoByteArrayFromJavaByteArray(jbaMore)) != "56789abcdefghij" {
		t.Errorf("Expected remaining 15 bytes '56789abcdefghij', got %s", object.GoByteArrayFromJavaByteArray(jbaMore))
	}

	// Negative length should throw IllegalArgumentException
	resNeg := inputStreamReadNBytesLen([]any{bais, int64(-1)})
	errBlk, ok := resNeg.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for negative len, got %v", resNeg)
	}
}

func TestInputStreamReadNBytesBuf(t *testing.T) {
	data := []byte("abcdefghijklmnopqrstuvwxyz")
	bais := createTestByteArrayInputStream(data)

	dest := make([]types.JavaByte, 10)
	destObj := object.MakePrimitiveObject("[B", types.JavaByteArray, dest)

	// Read 4 bytes at offset 2
	res := inputStreamReadNBytesBuf([]any{bais, destObj, int64(2), int64(4)})
	n, ok := res.(int64)
	if !ok || n != 4 {
		t.Fatalf("Expected 4 bytes read, got %v", res)
	}

	actual := destObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if actual[2] != types.JavaByte('a') || actual[5] != types.JavaByte('d') {
		t.Errorf("Expected bytes written at offset 2..5, got %v", actual)
	}

	// Read 0 length
	res0 := inputStreamReadNBytesBuf([]any{bais, destObj, int64(0), int64(0)})
	if res0.(int64) != 0 {
		t.Errorf("Expected 0 for len=0, got %v", res0)
	}

	// Out of bounds check
	resOOB := inputStreamReadNBytesBuf([]any{bais, destObj, int64(8), int64(5)})
	if errBlk, ok := resOOB.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Errorf("Expected IndexOutOfBoundsException, got %v", resOOB)
	}

	resNegOff := inputStreamReadNBytesBuf([]any{bais, destObj, int64(-1), int64(2)})
	if errBlk, ok := resNegOff.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Errorf("Expected IndexOutOfBoundsException for negative offset, got %v", resNegOff)
	}
}

func TestInputStreamSkipNBytes(t *testing.T) {
	data := []byte("0123456789abcdefghij")
	bais := createTestByteArrayInputStream(data)

	// Skip 0 or negative
	if res := inputStreamSkipNBytes([]any{bais, int64(0)}); res != nil {
		t.Errorf("Expected nil for skipNBytes(0), got %v", res)
	}
	if res := inputStreamSkipNBytes([]any{bais, int64(-5)}); res != nil {
		t.Errorf("Expected nil for skipNBytes(-5), got %v", res)
	}

	// Skip 5 bytes ("01234")
	if res := inputStreamSkipNBytes([]any{bais, int64(5)}); res != nil {
		t.Fatalf("skipNBytes(5) failed: %v", res)
	}

	// Next byte should be '5'
	nextRes := inputStreamReadNBytesLen([]any{bais, int64(1)})
	nextObj := nextRes.(*object.Object)
	nextBytes := nextObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if nextBytes[0] != types.JavaByte('5') {
		t.Errorf("Expected '5', got %c", nextBytes[0])
	}

	// Skip 14 bytes (remaining 14 bytes: '6789abcdefghij')
	if res := inputStreamSkipNBytes([]any{bais, int64(14)}); res != nil {
		t.Fatalf("skipNBytes(14) failed: %v", res)
	}

	// Skipping beyond EOF should return IOException / EOFException
	resEOF := inputStreamSkipNBytes([]any{bais, int64(1)})
	if errBlk, ok := resEOF.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException for skip past EOF, got %v", resEOF)
	}
}

func TestInputStreamTransferTo(t *testing.T) {
	data := []byte("Testing InputStream.transferTo() with Jacobin virtual machine streams!")
	bais := createTestByteArrayInputStream(data)
	baos := createTestByteArrayOutputStream()

	res := inputStreamTransferTo([]any{bais, baos})
	transferred, ok := res.(int64)
	if !ok || transferred != int64(len(data)) {
		t.Fatalf("Expected transferred %d, got %v", len(data), res)
	}

	written := getByteArrayOutputStreamData(baos)
	if !bytes.Equal(written, data) {
		t.Errorf("Transferred data mismatch: got %q, expected %q", written, data)
	}

	// Null destination stream should throw NullPointerException
	resNull := inputStreamTransferTo([]any{bais, object.Null})
	if errBlk, ok := resNull.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null out, got %v", resNull)
	}
}

func TestInputStreamMiscellaneousMethods(t *testing.T) {
	initTestStreams()
	inName := "java/io/InputStream"
	inObj := object.MakeEmptyObject()
	inObj.KlassName = stringPool.GetStringIndex(&inName)

	// available()
	if val := inputStreamAvailable([]any{inObj}); val.(int64) != 0 {
		t.Errorf("Expected available 0, got %v", val)
	}

	// nullInputStream()
	nullStream := inputStreamNullInputStream([]any{})
	if object.IsNull(nullStream) {
		t.Errorf("Expected non-null nullInputStream object")
	}

	// reset() throws IOException
	resReset := inputStreamReset([]any{inObj})
	if errBlk, ok := resReset.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException for reset(), got %v", resReset)
	}

	// Null receivers check
	if errBlk, ok := inputStreamReadAllBytes([]any{object.Null}).(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE on null receiver for readAllBytes, got %v", errBlk)
	}
	if errBlk, ok := inputStreamReadNBytesLen([]any{object.Null, int64(1)}).(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE on null receiver for readNBytes(I), got %v", errBlk)
	}
	if errBlk, ok := inputStreamReadNBytesBuf([]any{object.Null, inObj, int64(0), int64(0)}).(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE on null receiver for readNBytes([BII), got %v", errBlk)
	}
	if errBlk, ok := inputStreamSkipNBytes([]any{object.Null, int64(1)}).(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE on null receiver for skipNBytes, got %v", errBlk)
	}
	if errBlk, ok := inputStreamTransferTo([]any{object.Null, inObj}).(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE on null receiver for transferTo, got %v", errBlk)
	}
}

func TestInputStreamReadIntoByteArray_ZeroLength(t *testing.T) {
	inputStreamObj := &object.Object{}
	byteArrayObj := &object.Object{}

	params := []interface{}{inputStreamObj, byteArrayObj}
	result := inputStreamReadIntoByteArray(params)

	resInt, ok := result.(int64)
	if !ok {
		t.Errorf("Expected result to be int64, got %T", result)
	} else if resInt != 0 {
		t.Errorf("Expected 0 bytes read for zero length request, got %d", resInt)
	}
}

func TestInputStreamReadIntoByteArray_NilArray(t *testing.T) {
	inputStreamObj := &object.Object{}
	params := []interface{}{inputStreamObj, nil}

	result := inputStreamReadIntoByteArray(params)
	errBlock, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Errorf("Expected an *ghelpers.GErrBlk representing NullPointerException, got %T", result)
	} else if errBlock == nil {
		t.Error("Expected an error block, got nil")
	}
}

func TestInputStreamReadIntoByteArray_ValidRead(t *testing.T) {
	data := []byte("hello")
	bais := createTestByteArrayInputStream(data)
	dest := make([]types.JavaByte, 5)
	destObj := object.MakePrimitiveObject("[B", types.JavaByteArray, dest)

	params := []interface{}{bais, destObj}
	result := inputStreamReadIntoByteArray(params)

	if n, ok := result.(int64); !ok || n != 5 {
		t.Errorf("Expected 5 bytes read, got %v", result)
	}
}

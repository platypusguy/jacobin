/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"hash/adler32"
	"io"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

func initDeflaterTest() {
	globals.InitGlobals("test")
	Load_Util_Zip_Deflater()
}

func newTestDeflaterObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&deflaterClassName)
}

func makeJavaByteArrayObj(b []byte) *object.Object {
	jBytes := make([]types.JavaByte, len(b))
	for i, v := range b {
		jBytes[i] = types.JavaByte(v)
	}
	return object.MakePrimitiveObject("[B", types.JavaByteArray, jBytes)
}

func getBytesFromJavaByteArrayObj(obj *object.Object) []byte {
	valFld, ok := obj.FieldTable["value"]
	if !ok || valFld.Fvalue == nil {
		return nil
	}
	javaBytes, ok := valFld.Fvalue.([]types.JavaByte)
	if !ok {
		return nil
	}
	return object.GoByteArrayFromJavaByteArray(javaBytes)
}

func TestDeflaterClinitAndConstants(t *testing.T) {
	initDeflaterTest()
	res := deflaterClinit(nil)
	if res != nil {
		t.Fatalf("expected nil return from deflaterClinit, got: %v", res)
	}

	tests := []struct {
		name     string
		expected int64
	}{
		{"DEFLATED", DeflaterDeflated},
		{"NO_COMPRESSION", DeflaterNoCompression},
		{"BEST_SPEED", DeflaterBestSpeed},
		{"BEST_COMPRESSION", DeflaterBestCompression},
		{"DEFAULT_COMPRESSION", DeflaterDefaultCompression},
		{"FILTERED", DeflaterFiltered},
		{"HUFFMAN_ONLY", DeflaterHuffmanOnly},
		{"DEFAULT_STRATEGY", DeflaterDefaultStrategy},
		{"NO_FLUSH", DeflaterNoFlush},
		{"SYNC_FLUSH", DeflaterSyncFlush},
		{"FULL_FLUSH", DeflaterFullFlush},
	}

	for _, tt := range tests {
		staticVal, ok := statics.Statics["java/util/zip/Deflater."+tt.name]
		if !ok {
			t.Errorf("static field %s not found", tt.name)
			continue
		}
		if staticVal.Value.(int64) != tt.expected {
			t.Errorf("static field %s expected %d, got %v", tt.name, tt.expected, staticVal.Value)
		}
	}
}

func TestDeflaterConstructors(t *testing.T) {
	initDeflaterTest()

	// 1. Default constructor
	obj1 := newTestDeflaterObj()
	res1 := deflaterInitDefault([]any{obj1})
	if res1 != nil {
		t.Fatalf("expected nil from deflaterInitDefault, got %v", res1)
	}
	state1 := getDeflaterState(obj1)
	if state1 == nil || state1.level != int(DeflaterDefaultCompression) || state1.nowrap != false {
		t.Fatalf("unexpected state for default deflater: %+v", state1)
	}

	// 2. Level constructor
	obj2 := newTestDeflaterObj()
	res2 := deflaterInitLevel([]any{obj2, int64(9)})
	if res2 != nil {
		t.Fatalf("expected nil from deflaterInitLevel, got %v", res2)
	}
	state2 := getDeflaterState(obj2)
	if state2 == nil || state2.level != 9 || state2.nowrap != false {
		t.Fatalf("unexpected state for level deflater: %+v", state2)
	}

	// 3. Level + Nowrap constructor
	obj3 := newTestDeflaterObj()
	res3 := deflaterInitLevelNowrap([]any{obj3, int64(1), types.JavaBoolTrue})
	if res3 != nil {
		t.Fatalf("expected nil from deflaterInitLevelNowrap, got %v", res3)
	}
	state3 := getDeflaterState(obj3)
	if state3 == nil || state3.level != 1 || state3.nowrap != true {
		t.Fatalf("unexpected state for nowrap deflater: %+v", state3)
	}

	// 4. Invalid level
	obj4 := newTestDeflaterObj()
	res4 := deflaterInitLevelNowrap([]any{obj4, int64(10), types.JavaBoolFalse})
	errBlk, ok := asGErrBlk(res4)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException for level 10, got %v", res4)
	}

	// 5. Null receiver
	res5 := deflaterInitDefault([]any{object.Null})
	errBlk5, ok := asGErrBlk(res5)
	if !ok || errBlk5.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NullPointerException for null receiver, got %v", res5)
	}
}

func TestDeflaterStandardZlibCompression(t *testing.T) {
	initDeflaterTest()
	originalData := []byte("The quick brown fox jumps over the lazy dog. Repetition repetition repetition to test compression compression!")

	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	inputObj := makeJavaByteArrayObj(originalData)
	resInput := deflaterSetInputArray([]any{defObj, inputObj})
	if resInput != nil {
		t.Fatalf("unexpected error from setInput: %v", resInput)
	}

	if deflaterNeedsInput([]any{defObj}) != types.JavaBoolFalse {
		t.Fatalf("expected needsInput() to be false after setInput")
	}

	deflaterFinish([]any{defObj})

	outBufObj := makeJavaByteArrayObj(make([]byte, 1024))
	nVal := deflaterDeflateArray([]any{defObj, outBufObj})
	n, ok := nVal.(int64)
	if !ok || n <= 0 {
		t.Fatalf("expected positive bytes deflated, got: %v", nVal)
	}

	if deflaterFinished([]any{defObj}) != types.JavaBoolTrue {
		t.Fatalf("expected finished() to be true")
	}

	compressedData := getBytesFromJavaByteArrayObj(outBufObj)[:n]

	// Decompress with Go's zlib reader and verify
	zr, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		t.Fatalf("failed to create zlib reader: %v", err)
	}
	decompressed, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed to decompress with zlib: %v", err)
	}
	zr.Close()

	if !bytes.Equal(decompressed, originalData) {
		t.Fatalf("decompressed data mismatch!\nExpected: %s\nGot:      %s", originalData, decompressed)
	}

	// Verify getAdler() matches
	expectedAdler := int64(int32(adler32.Checksum(originalData)))
	actualAdler := deflaterGetAdler([]any{defObj}).(int64)
	if actualAdler != expectedAdler {
		t.Fatalf("expected Adler %d, got %d", expectedAdler, actualAdler)
	}

	// Verify bytes read and written
	if deflaterGetBytesRead([]any{defObj}).(int64) != int64(len(originalData)) {
		t.Fatalf("expected bytesRead == %d", len(originalData))
	}
	if deflaterGetBytesWritten([]any{defObj}).(int64) != n {
		t.Fatalf("expected bytesWritten == %d", n)
	}
}

func TestDeflaterNowrapRawCompression(t *testing.T) {
	initDeflaterTest()
	originalData := []byte("Testing raw deflate without zlib header or trailer wrapper for ZIP entries.")

	defObj := newTestDeflaterObj()
	deflaterInitLevelNowrap([]any{defObj, int64(DeflaterBestCompression), types.JavaBoolTrue})

	inputObj := makeJavaByteArrayObj(originalData)
	deflaterSetInputRange([]any{defObj, inputObj, int64(0), int64(len(originalData))})
	deflaterFinish([]any{defObj})

	outBufObj := makeJavaByteArrayObj(make([]byte, 512))
	nVal := deflaterDeflateRange([]any{defObj, outBufObj, int64(0), int64(512)})
	n := nVal.(int64)
	if n <= 0 {
		t.Fatalf("expected positive bytes deflated, got: %d", n)
	}

	compressedData := getBytesFromJavaByteArrayObj(outBufObj)[:n]

	// Decompress with Go's raw flate reader
	fr := flate.NewReader(bytes.NewReader(compressedData))
	decompressed, err := io.ReadAll(fr)
	if err != nil {
		t.Fatalf("failed to decompress raw flate: %v", err)
	}
	fr.Close()

	if !bytes.Equal(decompressed, originalData) {
		t.Fatalf("decompressed raw flate mismatch!\nExpected: %s\nGot:      %s", originalData, decompressed)
	}
}

func TestDeflaterChunkedStreaming(t *testing.T) {
	initDeflaterTest()
	data1 := []byte("First chunk of stream. ")
	data2 := []byte("Second chunk of stream with additional text. ")
	data3 := []byte("Third and final chunk of stream.")
	fullData := append(append(append([]byte(nil), data1...), data2...), data3...)

	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	var compressedResult []byte
	smallBufObj := makeJavaByteArrayObj(make([]byte, 16))

	// Feed chunk 1
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(data1)})
	for deflaterNeedsInput([]any{defObj}) == types.JavaBoolFalse {
		n := deflaterDeflateArray([]any{defObj, smallBufObj}).(int64)
		if n > 0 {
			compressedResult = append(compressedResult, getBytesFromJavaByteArrayObj(smallBufObj)[:n]...)
		}
	}

	// Feed chunk 2
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(data2)})
	for deflaterNeedsInput([]any{defObj}) == types.JavaBoolFalse {
		n := deflaterDeflateArray([]any{defObj, smallBufObj}).(int64)
		if n > 0 {
			compressedResult = append(compressedResult, getBytesFromJavaByteArrayObj(smallBufObj)[:n]...)
		}
	}

	// Feed chunk 3 and finish
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(data3)})
	deflaterFinish([]any{defObj})

	for deflaterFinished([]any{defObj}) == types.JavaBoolFalse {
		n := deflaterDeflateArray([]any{defObj, smallBufObj}).(int64)
		if n > 0 {
			compressedResult = append(compressedResult, getBytesFromJavaByteArrayObj(smallBufObj)[:n]...)
		} else if n == 0 && deflaterFinished([]any{defObj}) == types.JavaBoolFalse {
			t.Fatalf("deflate returned 0 while finished is false")
		}
	}

	// Decompress and verify
	zr, err := zlib.NewReader(bytes.NewReader(compressedResult))
	if err != nil {
		t.Fatalf("failed to create zlib reader: %v", err)
	}
	decompressed, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed to decompress chunked: %v", err)
	}
	zr.Close()

	if !bytes.Equal(decompressed, fullData) {
		t.Fatalf("chunked stream mismatch!\nExpected: %s\nGot:      %s", fullData, decompressed)
	}
}

func TestDeflaterSyncFlush(t *testing.T) {
	initDeflaterTest()
	data1 := []byte("Message block 1 for sync flush.")
	data2 := []byte("Message block 2 after sync flush.")

	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	outBufObj := makeJavaByteArrayObj(make([]byte, 1024))

	// 1. Write block 1 and sync flush
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(data1)})
	n1 := deflaterDeflateRangeFlush([]any{defObj, outBufObj, int64(0), int64(1024), DeflaterSyncFlush}).(int64)
	if n1 <= 0 {
		t.Fatalf("expected positive bytes from SYNC_FLUSH, got %d", n1)
	}
	part1 := append([]byte(nil), getBytesFromJavaByteArrayObj(outBufObj)[:n1]...)

	// 2. Write block 2 and finish
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(data2)})
	deflaterFinish([]any{defObj})

	n2 := deflaterDeflateRangeFlush([]any{defObj, outBufObj, int64(0), int64(1024), DeflaterNoFlush}).(int64)
	part2 := append([]byte(nil), getBytesFromJavaByteArrayObj(outBufObj)[:n2]...)

	allCompressed := append(part1, part2...)

	zr, err := zlib.NewReader(bytes.NewReader(allCompressed))
	if err != nil {
		t.Fatalf("failed zlib reader on sync flush data: %v", err)
	}
	decompressed, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed decompression: %v", err)
	}
	zr.Close()

	expectedFull := append(append([]byte(nil), data1...), data2...)
	if !bytes.Equal(decompressed, expectedFull) {
		t.Fatalf("mismatch on sync flush stream!\nExpected: %s\nGot:      %s", expectedFull, decompressed)
	}
}

func TestDeflaterResetAndReuse(t *testing.T) {
	initDeflaterTest()
	msg1 := []byte("First message to compress")
	msg2 := []byte("Second completely independent message")

	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	// Pass 1
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(msg1)})
	deflaterFinish([]any{defObj})
	outBuf1 := makeJavaByteArrayObj(make([]byte, 256))
	n1 := deflaterDeflateArray([]any{defObj, outBuf1}).(int64)
	comp1 := getBytesFromJavaByteArrayObj(outBuf1)[:n1]

	// Reset
	deflaterReset([]any{defObj})
	if deflaterNeedsInput([]any{defObj}) != types.JavaBoolTrue {
		t.Fatalf("expected needsInput() == true after reset")
	}
	if deflaterFinished([]any{defObj}) != types.JavaBoolFalse {
		t.Fatalf("expected finished() == false after reset")
	}
	if deflaterGetBytesRead([]any{defObj}).(int64) != 0 {
		t.Fatalf("expected bytesRead == 0 after reset")
	}
	if deflaterGetBytesWritten([]any{defObj}).(int64) != 0 {
		t.Fatalf("expected bytesWritten == 0 after reset")
	}
	if deflaterGetAdler([]any{defObj}).(int64) != 1 {
		t.Fatalf("expected adler == 1 after reset")
	}

	// Pass 2
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(msg2)})
	deflaterFinish([]any{defObj})
	outBuf2 := makeJavaByteArrayObj(make([]byte, 256))
	n2 := deflaterDeflateArray([]any{defObj, outBuf2}).(int64)
	comp2 := getBytesFromJavaByteArrayObj(outBuf2)[:n2]

	// Verify both pass 1 and pass 2 decompress correctly
	zr1, _ := zlib.NewReader(bytes.NewReader(comp1))
	dec1, _ := io.ReadAll(zr1)
	zr1.Close()
	if !bytes.Equal(dec1, msg1) {
		t.Fatalf("pass 1 mismatch: got %s", dec1)
	}

	zr2, _ := zlib.NewReader(bytes.NewReader(comp2))
	dec2, _ := io.ReadAll(zr2)
	zr2.Close()
	if !bytes.Equal(dec2, msg2) {
		t.Fatalf("pass 2 mismatch: got %s", dec2)
	}
}

func TestDeflaterPresetDictionary(t *testing.T) {
	initDeflaterTest()
	dict := []byte("this is a shared compression dictionary with frequent tokens ")
	msg := []byte("frequent tokens are in this is a shared compression dictionary!")

	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	dictObj := makeJavaByteArrayObj(dict)
	deflaterSetDictionaryArray([]any{defObj, dictObj})

	// getAdler() should return adler of the dictionary before compression
	expectedDictAdler := int64(int32(adler32.Checksum(dict)))
	if deflaterGetAdler([]any{defObj}).(int64) != expectedDictAdler {
		t.Fatalf("expected dict adler %d, got %d", expectedDictAdler, deflaterGetAdler([]any{defObj}))
	}

	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(msg)})
	deflaterFinish([]any{defObj})

	outBufObj := makeJavaByteArrayObj(make([]byte, 512))
	n := deflaterDeflateArray([]any{defObj, outBufObj}).(int64)
	comp := getBytesFromJavaByteArrayObj(outBufObj)[:n]

	// A zlib reader validates the trailer Adler-32 against uncompressed data only.
	zr, err := zlib.NewReaderDict(bytes.NewReader(comp), dict)
	if err != nil {
		t.Fatalf("failed to create zlib reader with dictionary: %v", err)
	}
	zlibDecoded, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed to decompress zlib stream with dictionary: %v", err)
	}
	if err := zr.Close(); err != nil {
		t.Fatalf("failed to close zlib reader with dictionary: %v", err)
	}
	if !bytes.Equal(zlibDecoded, msg) {
		t.Fatalf("zlib dictionary decompression mismatch: got %s", zlibDecoded)
	}

	// Decompress with flate with dictionary
	// In zlib stream with dictionary, header is 6 bytes (2 byte header + 4 byte dict id)
	if len(comp) < 6 {
		t.Fatalf("compressed data too short for zlib dict: %d", len(comp))
	}
	rawFlateData := comp[6 : len(comp)-4]
	fr := flate.NewReaderDict(bytes.NewReader(rawFlateData), dict)
	dec, err := io.ReadAll(fr)
	if err != nil {
		t.Fatalf("flate dict decompression failed: %v", err)
	}
	fr.Close()

	if !bytes.Equal(dec, msg) {
		t.Fatalf("decompressed dict mismatch!\nExpected: %s\nGot:      %s", msg, dec)
	}
}

func TestDeflaterNowrapDictionaryAdler(t *testing.T) {
	initDeflaterTest()
	dict := []byte("a preset dictionary for raw deflate")
	msg := []byte("raw deflate data that shares dictionary words")
	defObj := newTestDeflaterObj()
	deflaterInitLevelNowrap([]any{defObj, int64(DeflaterDefaultCompression), types.JavaBoolTrue})
	deflaterSetDictionaryArray([]any{defObj, makeJavaByteArrayObj(dict)})
	deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj(msg)})
	deflaterFinish([]any{defObj})
	outBuf := makeJavaByteArrayObj(make([]byte, 512))
	deflaterDeflateArray([]any{defObj, outBuf})

	expectedAdler := int64(int32(adler32.Checksum(msg)))
	if actualAdler := deflaterGetAdler([]any{defObj}).(int64); actualAdler != expectedAdler {
		t.Fatalf("expected Adler %d for uncompressed input, got %d", expectedAdler, actualAdler)
	}
}

func TestDeflaterSetLevelAndStrategy(t *testing.T) {
	initDeflaterTest()
	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	// Valid setLevel
	res := deflaterSetLevel([]any{defObj, int64(DeflaterBestCompression)})
	if res != nil {
		t.Fatalf("expected nil from setLevel, got: %v", res)
	}

	// Invalid setLevel
	resInvalidLvl := deflaterSetLevel([]any{defObj, int64(15)})
	errBlk, ok := asGErrBlk(resInvalidLvl)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException on level 15, got: %v", resInvalidLvl)
	}

	// Valid setStrategy
	resStrat := deflaterSetStrategy([]any{defObj, int64(DeflaterHuffmanOnly)})
	if resStrat != nil {
		t.Fatalf("expected nil from setStrategy, got: %v", resStrat)
	}

	// Invalid setStrategy
	resInvalidStrat := deflaterSetStrategy([]any{defObj, int64(99)})
	errBlk2, ok := asGErrBlk(resInvalidStrat)
	if !ok || errBlk2.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException on strategy 99, got: %v", resInvalidStrat)
	}
}

func TestDeflaterEndAndExceptions(t *testing.T) {
	initDeflaterTest()
	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	deflaterEnd([]any{defObj})
	if state := getDeflaterState(defObj); state != nil {
		t.Fatalf("expected state to be removed after end(), got %+v", state)
	}

	// After end(), operations should return NullPointerException
	res := deflaterSetInputArray([]any{defObj, makeJavaByteArrayObj([]byte{1, 2, 3})})
	errBlk, ok := asGErrBlk(res)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NullPointerException after end(), got: %v", res)
	}

	resDeflate := deflaterDeflateArray([]any{defObj, makeJavaByteArrayObj(make([]byte, 10))})
	errBlk2, ok := asGErrBlk(resDeflate)
	if !ok || errBlk2.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NullPointerException after end(), got: %v", resDeflate)
	}
}

func TestDeflaterTrappedMethods(t *testing.T) {
	initDeflaterTest()

	trappedMethods := []string{
		"java/util/zip/Deflater.setInput(Ljava/nio/ByteBuffer;)V",
		"java/util/zip/Deflater.setDictionary(Ljava/nio/ByteBuffer;)V",
		"java/util/zip/Deflater.deflate(Ljava/nio/ByteBuffer;)I",
		"java/util/zip/Deflater.deflate(Ljava/nio/ByteBuffer;I)I",
		"java/util/zip/Deflater.getClass()Ljava/lang/Class;",
		"java/util/zip/Deflater.notify()V",
		"java/util/zip/Deflater.notifyAll()V",
		"java/util/zip/Deflater.wait()V",
		"java/util/zip/Deflater.wait(J)V",
		"java/util/zip/Deflater.wait(JI)V",
	}

	for _, meth := range trappedMethods {
		gmeth, ok := ghelpers.MethodSignatures[meth]
		if !ok {
			t.Errorf("method signature %s not registered", meth)
			continue
		}
		res := gmeth.GFunction([]any{})
		errBlk, ok := asGErrBlk(res)
		if !ok || errBlk.ExceptionType != excNames.UnsupportedOperationException {
			t.Errorf("expected trap UnsupportedOperationException for %s, got: %v", meth, res)
		}
	}
}

func TestDeflaterBoundsAndEdgeCases(t *testing.T) {
	initDeflaterTest()
	defObj := newTestDeflaterObj()
	deflaterInitDefault([]any{defObj})

	bufObj := makeJavaByteArrayObj(make([]byte, 20))

	// 1. Negative offset
	res := deflaterSetInputRange([]any{defObj, bufObj, int64(-1), int64(10)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ArrayIndexOutOfBoundsException {
		t.Fatalf("expected ArrayIndexOutOfBoundsException for negative offset, got %v", res)
	}

	// 2. Negative length
	res = deflaterSetInputRange([]any{defObj, bufObj, int64(0), int64(-5)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ArrayIndexOutOfBoundsException {
		t.Fatalf("expected ArrayIndexOutOfBoundsException for negative length, got %v", res)
	}

	// 3. Offset + Length overflow
	res = deflaterSetInputRange([]any{defObj, bufObj, int64(15), int64(10)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ArrayIndexOutOfBoundsException {
		t.Fatalf("expected ArrayIndexOutOfBoundsException for out of bounds, got %v", res)
	}

	// 4. Dictionary bounds
	res = deflaterSetDictionaryRange([]any{defObj, bufObj, int64(10), int64(20)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ArrayIndexOutOfBoundsException {
		t.Fatalf("expected ArrayIndexOutOfBoundsException for dictionary out of bounds, got %v", res)
	}

	// 5. Deflate bounds
	res = deflaterDeflateRangeFlush([]any{defObj, bufObj, int64(15), int64(10), DeflaterNoFlush})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ArrayIndexOutOfBoundsException {
		t.Fatalf("expected ArrayIndexOutOfBoundsException for deflate out of bounds, got %v", res)
	}

	// 6. Invalid flush flag
	res = deflaterDeflateRangeFlush([]any{defObj, bufObj, int64(0), int64(10), int64(99)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException for flush 99, got %v", res)
	}

	// 7. Deflate with zero length
	nZero := deflaterDeflateRange([]any{defObj, bufObj, int64(0), int64(0)}).(int64)
	if nZero != 0 {
		t.Fatalf("expected 0 for zero length deflate, got %d", nZero)
	}

	// 8. Empty input compression
	defEmpty := newTestDeflaterObj()
	deflaterInitDefault([]any{defEmpty})
	deflaterFinish([]any{defEmpty})
	outEmpty := makeJavaByteArrayObj(make([]byte, 64))
	nEmpty := deflaterDeflateArray([]any{defEmpty, outEmpty}).(int64)
	if nEmpty <= 0 {
		t.Fatalf("expected positive bytes for empty zlib deflate, got %d", nEmpty)
	}
	if deflaterFinished([]any{defEmpty}) != types.JavaBoolTrue {
		t.Fatalf("expected finished == true for empty deflate")
	}

	zr, err := zlib.NewReader(bytes.NewReader(getBytesFromJavaByteArrayObj(outEmpty)[:nEmpty]))
	if err != nil {
		t.Fatalf("failed zlib reader on empty stream: %v", err)
	}
	decEmpty, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed to read empty stream: %v", err)
	}
	zr.Close()
	if len(decEmpty) != 0 {
		t.Fatalf("expected 0 decompressed bytes, got %d", len(decEmpty))
	}

	// 9. Finalize
	resFinalize := deflaterFinalize([]any{defEmpty})
	if resFinalize != nil {
		t.Fatalf("expected nil from finalize, got %v", resFinalize)
	}
	if state := getDeflaterState(defEmpty); state != nil {
		t.Fatalf("expected state to be removed after finalize, got %+v", state)
	}
}

func TestDeflaterObjectOverrides(t *testing.T) {
	initDeflaterTest()
	def1 := newTestDeflaterObj()
	deflaterInitDefault([]any{def1})
	def2 := newTestDeflaterObj()
	deflaterInitDefault([]any{def2})

	// equals
	if deflaterEquals([]any{def1, def1}) != types.JavaBoolTrue {
		t.Fatalf("expected def1.equals(def1) == true")
	}
	if deflaterEquals([]any{def1, def2}) != types.JavaBoolFalse {
		t.Fatalf("expected def1.equals(def2) == false")
	}

	// hashCode
	hc1 := deflaterHashCode([]any{def1}).(int64)
	if hc1 < 0 {
		t.Fatalf("expected non-negative hash code, got %d", hc1)
	}

	// toString
	strObj := deflaterToString([]any{def1}).(*object.Object)
	str := object.GoStringFromStringObject(strObj)
	if len(str) == 0 || !bytes.Contains([]byte(str), []byte("java.util.zip.Deflater@")) {
		t.Fatalf("unexpected toString(): %s", str)
	}
}

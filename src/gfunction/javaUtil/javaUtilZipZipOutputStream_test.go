/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"archive/zip"
	"bytes"
	"hash/crc32"
	"strings"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaIo"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

func initZipOutputStreamTests() {
	globals.InitGlobals("test")
	javaIo.Load_Io_InputStream()
	javaIo.Load_Io_ByteArrayInputStream()
	javaIo.Load_Io_ByteArrayOutputStream()
	Load_Util_Zip_ZipEntry()
	Load_Util_Zip_ZipInputStream()
	Load_Util_Zip_ZipOutputStream()
}

func newTestZipOutputStream(outObj *object.Object) *object.Object {
	zosObj := object.MakeEmptyObjectWithClassName(&zipOutputStreamClassName)
	zipOutputStreamInit([]any{zosObj, outObj})
	return zosObj
}

func javaByteArrayFromGoBytes(data []byte) *object.Object {
	jBytes := make([]types.JavaByte, len(data))
	for i, b := range data {
		jBytes[i] = types.JavaByte(b)
	}
	return object.MakePrimitiveObject("[B", types.JavaByteArray, jBytes)
}

func TestZipOutputStream_Clinit(t *testing.T) {
	initZipOutputStreamTests()

	zipOutputStreamClinit(nil)

	sStored, ok := statics.Statics[zipOutputStreamClassName+".STORED"]
	if !ok || sStored.Value.(int64) != ZipEntryStored {
		t.Fatalf("STORED static field incorrect: %v, ok=%v", sStored, ok)
	}

	sDeflated, ok := statics.Statics[zipOutputStreamClassName+".DEFLATED"]
	if !ok || sDeflated.Value.(int64) != ZipEntryDeflated {
		t.Fatalf("DEFLATED static field incorrect: %v, ok=%v", sDeflated, ok)
	}
}

func TestZipOutputStream_Init(t *testing.T) {
	initZipOutputStreamTests()

	// 1. Successful initialization
	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)
	if zos == nil {
		t.Fatal("Failed to create ZipOutputStream")
	}

	state := getZipOutStreamState(zos)
	if state == nil || state.closed || state.finished {
		t.Fatal("State should be active and not closed")
	}
	if state.defaultMethod != ZipEntryDeflated {
		t.Fatalf("Expected default method DEFLATED (8), got %d", state.defaultMethod)
	}

	// 2. Null receiver
	res := zipOutputStreamInit([]any{nil, baos})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("Expected NPE for null receiver, got: %v", res)
	}

	// 3. Null OutputStream
	res = zipOutputStreamInit([]any{zos, nil})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("Expected NPE for null OutputStream, got: %v", res)
	}

	// 4. Trapped 2-arg constructor
	meth, ok := ghelpers.MethodSignatures["java/util/zip/ZipOutputStream.<init>(Ljava/io/OutputStream;Ljava/nio/charset/Charset;)V"]
	if !ok || meth.GFunction == nil {
		t.Fatal("Trapped constructor not registered")
	}
	trapRes := meth.GFunction(nil)
	if trapRes == nil {
		t.Fatal("Expected trap function return")
	}
}

func TestZipOutputStream_SetComment(t *testing.T) {
	initZipOutputStreamTests()

	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	// Valid comment
	commentStr := "Sample ZIP comment"
	res := zipOutputStreamSetComment([]any{zos, object.StringObjectFromGoString(commentStr)})
	if res != nil {
		t.Fatalf("Unexpected error from setComment: %v", res)
	}
	state := getZipOutStreamState(zos)
	if state.comment != commentStr {
		t.Fatalf("Expected comment %s, got %s", commentStr, state.comment)
	}

	// Null comment (clears comment)
	res = zipOutputStreamSetComment([]any{zos, object.Null})
	if res != nil {
		t.Fatalf("Unexpected error: %v", res)
	}
	if state.comment != "" {
		t.Fatalf("Expected empty comment, got %s", state.comment)
	}

	// Too long comment (> 65535 bytes)
	longComment := strings.Repeat("A", 65536)
	res = zipOutputStreamSetComment([]any{zos, object.StringObjectFromGoString(longComment)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("Expected IllegalArgumentException for long comment, got %v", res)
	}

	// Null receiver
	res = zipOutputStreamSetComment([]any{nil, object.StringObjectFromGoString("test")})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("Expected NPE, got %v", res)
	}

	// Closed stream
	zipOutputStreamClose([]any{zos})
	res = zipOutputStreamSetComment([]any{zos, object.StringObjectFromGoString("test")})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IOException {
		t.Fatalf("Expected IOException on closed stream, got %v", res)
	}
}

func TestZipOutputStream_SetMethod_And_Level(t *testing.T) {
	initZipOutputStreamTests()

	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	// Valid methods
	res := zipOutputStreamSetMethod([]any{zos, ZipEntryStored})
	if res != nil {
		t.Fatalf("Unexpected error: %v", res)
	}
	state := getZipOutStreamState(zos)
	if state.defaultMethod != ZipEntryStored {
		t.Fatalf("Expected default method STORED (0), got %d", state.defaultMethod)
	}

	res = zipOutputStreamSetMethod([]any{zos, ZipEntryDeflated})
	if res != nil {
		t.Fatalf("Unexpected error: %v", res)
	}
	if state.defaultMethod != ZipEntryDeflated {
		t.Fatalf("Expected default method DEFLATED (8), got %d", state.defaultMethod)
	}

	// Invalid method
	res = zipOutputStreamSetMethod([]any{zos, int64(3)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("Expected IllegalArgumentException for invalid method, got %v", res)
	}

	// Valid levels
	for lvl := int64(-1); lvl <= 9; lvl++ {
		res = zipOutputStreamSetLevel([]any{zos, lvl})
		if res != nil {
			t.Fatalf("Unexpected error for level %d: %v", lvl, res)
		}
		if state.defaultLevel != int(lvl) {
			t.Fatalf("Expected default level %d, got %d", lvl, state.defaultLevel)
		}
	}

	// Invalid levels
	res = zipOutputStreamSetLevel([]any{zos, int64(-2)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("Expected IllegalArgumentException for level -2, got %v", res)
	}
	res = zipOutputStreamSetLevel([]any{zos, int64(10)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("Expected IllegalArgumentException for level 10, got %v", res)
	}

	// Closed stream checks
	zipOutputStreamClose([]any{zos})
	res = zipOutputStreamSetMethod([]any{zos, ZipEntryDeflated})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IOException {
		t.Fatalf("Expected IOException on closed stream for setMethod, got %v", res)
	}
	res = zipOutputStreamSetLevel([]any{zos, int64(5)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IOException {
		t.Fatalf("Expected IOException on closed stream for setLevel, got %v", res)
	}
}

func TestZipOutputStream_PutNextEntry_Errors(t *testing.T) {
	initZipOutputStreamTests()

	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	// Null receiver
	res := zipOutputStreamPutNextEntry([]any{nil, nil})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("Expected NPE for null receiver, got %v", res)
	}

	// Null entry
	res = zipOutputStreamPutNextEntry([]any{zos, nil})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("Expected NPE for null entry, got %v", res)
	}

	// Duplicate entry name
	entry1, _ := CreateZipEntryObject("test.txt")
	res = zipOutputStreamPutNextEntry([]any{zos, entry1})
	if res != nil {
		t.Fatalf("Unexpected error for entry1: %v", res)
	}
	entry2, _ := CreateZipEntryObject("test.txt")
	res = zipOutputStreamPutNextEntry([]any{zos, entry2})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException for duplicate entry, got %v", res)
	}

	// STORED entry missing size/crc
	zos2 := newTestZipOutputStream(newTestByteArrayOutputStream())
	storedEntry, _ := CreateZipEntryObject("stored.txt")
	zipEntrySetMethod([]any{storedEntry, ZipEntryStored})
	res = zipOutputStreamPutNextEntry([]any{zos2, storedEntry})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException for STORED entry missing size/crc, got %v", res)
	}

	// STORED entry with invalid compressed size
	storedEntry2, _ := CreateZipEntryObject("stored2.txt")
	zipEntrySetMethod([]any{storedEntry2, ZipEntryStored})
	zipEntrySetSize([]any{storedEntry2, int64(10)})
	zipEntrySetCompressedSize([]any{storedEntry2, int64(20)})
	zipEntrySetCrc([]any{storedEntry2, int64(12345)})
	res = zipOutputStreamPutNextEntry([]any{zos2, storedEntry2})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException for STORED entry with mismatching csize, got %v", res)
	}
}

func TestZipOutputStream_Write_Errors(t *testing.T) {
	initZipOutputStreamTests()

	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	buf := javaByteArrayFromGoBytes([]byte("hello world"))

	// Write without active entry
	res := zipOutputStreamWriteRange([]any{zos, buf, int64(0), int64(5)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException when writing with no active entry, got %v", res)
	}

	// Write null buffer
	entry, _ := CreateZipEntryObject("file.txt")
	zipOutputStreamPutNextEntry([]any{zos, entry})
	res = zipOutputStreamWriteRange([]any{zos, nil, int64(0), int64(5)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("Expected NPE for null buffer, got %v", res)
	}

	// Write out of bounds
	res = zipOutputStreamWriteRange([]any{zos, buf, int64(-1), int64(5)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Fatalf("Expected IndexOutOfBoundsException for negative offset, got %v", res)
	}
	res = zipOutputStreamWriteRange([]any{zos, buf, int64(0), int64(100)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Fatalf("Expected IndexOutOfBoundsException for length > size, got %v", res)
	}

	// Write on closed stream
	zipOutputStreamClose([]any{zos})
	if state := getZipOutStreamState(zos); state != nil {
		t.Fatalf("expected stream state to be removed after close, got %+v", state)
	}
	res = zipOutputStreamWriteRange([]any{zos, buf, int64(0), int64(5)})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.IOException {
		t.Fatalf("Expected IOException on closed stream, got %v", res)
	}
}

func TestZipOutputStream_Overrides(t *testing.T) {
	initZipOutputStreamTests()

	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	// toString
	strRes := zipOutputStreamToString([]any{zos})
	strObj, ok := strRes.(*object.Object)
	if !ok {
		t.Fatalf("Expected String object, got %v", strRes)
	}
	strVal := object.GoStringFromStringObject(strObj)
	if !strings.HasPrefix(strVal, "java.util.zip.ZipOutputStream@") {
		t.Fatalf("Unexpected toString result: %s", strVal)
	}

	// equals and hashCode
	eqRes := zipOutputStreamEquals([]any{zos, zos})
	if eqRes != types.JavaBoolTrue {
		t.Fatalf("Expected JavaBoolTrue, got %v", eqRes)
	}
	eqRes2 := zipOutputStreamEquals([]any{zos, baos})
	if eqRes2 != types.JavaBoolFalse {
		t.Fatalf("Expected JavaBoolFalse, got %v", eqRes2)
	}
	hcRes := zipOutputStreamHashCode([]any{zos})
	if _, ok := hcRes.(int64); !ok {
		t.Fatalf("Expected int64 hashcode, got %v", hcRes)
	}

	// write(int b) and write(byte[] b) overrides
	entry, _ := CreateZipEntryObject("overrides.txt")
	zipOutputStreamPutNextEntry([]any{zos, entry})

	// write(I)V
	res := zipOutputStreamWriteByte([]any{zos, int64('H')})
	if res != nil {
		t.Fatalf("write(int) failed: %v", res)
	}

	// write([B)V
	byteBuf := javaByteArrayFromGoBytes([]byte("ello"))
	res = zipOutputStreamWriteByteArray([]any{zos, byteBuf})
	if res != nil {
		t.Fatalf("write(byte[]) failed: %v", res)
	}

	// deflate() and flush()
	res = zipOutputStreamDeflate([]any{zos})
	if res != nil {
		t.Fatalf("deflate() failed: %v", res)
	}
	res = zipOutputStreamFlush([]any{zos})
	if res != nil {
		t.Fatalf("flush() failed: %v", res)
	}

	zipOutputStreamCloseEntry([]any{zos})
	zipOutputStreamClose([]any{zos})

	// Verify content with archive/zip
	zipBytes := getByteArrayOutputStreamBytes(baos)
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("Failed to open zip reader: %v", err)
	}
	if len(zr.File) != 1 || zr.File[0].Name != "overrides.txt" {
		t.Fatalf("Unexpected files in zip: %v", zr.File)
	}
	rc, err := zr.File[0].Open()
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer rc.Close()
	var uncompressed bytes.Buffer
	_, err = uncompressed.ReadFrom(rc)
	if err != nil {
		t.Fatalf("Failed to read uncompressed: %v", err)
	}
	if uncompressed.String() != "Hello" {
		t.Fatalf("Expected 'Hello', got '%s'", uncompressed.String())
	}
}

func TestZipOutputStream_Stored_Validation(t *testing.T) {
	initZipOutputStreamTests()

	// 1. Successful STORED entry
	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	data := []byte("Exact STORED Content 12345")
	crc := crc32.ChecksumIEEE(data)

	entry, _ := CreateZipEntryObject("stored_exact.txt")
	zipEntrySetMethod([]any{entry, ZipEntryStored})
	zipEntrySetSize([]any{entry, int64(len(data))})
	zipEntrySetCompressedSize([]any{entry, int64(len(data))})
	zipEntrySetCrc([]any{entry, int64(crc)})

	res := zipOutputStreamPutNextEntry([]any{zos, entry})
	if res != nil {
		t.Fatalf("putNextEntry failed: %v", res)
	}

	res = zipOutputStreamWriteByteArray([]any{zos, javaByteArrayFromGoBytes(data)})
	if res != nil {
		t.Fatalf("write failed: %v", res)
	}

	res = zipOutputStreamCloseEntry([]any{zos})
	if res != nil {
		t.Fatalf("closeEntry failed: %v", res)
	}

	zipOutputStreamClose([]any{zos})

	// Verify with archive/zip
	zipBytes := getByteArrayOutputStreamBytes(baos)
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}
	if zr.File[0].Method != zip.Store {
		t.Fatalf("Expected STORED method, got %d", zr.File[0].Method)
	}
	rc, _ := zr.File[0].Open()
	var readBuf bytes.Buffer
	readBuf.ReadFrom(rc)
	rc.Close()
	if readBuf.String() != string(data) {
		t.Fatalf("Expected '%s', got '%s'", string(data), readBuf.String())
	}

	// 2. STORED entry writing too many bytes
	baos2 := newTestByteArrayOutputStream()
	zos2 := newTestZipOutputStream(baos2)
	entry2, _ := CreateZipEntryObject("overflow.txt")
	zipEntrySetMethod([]any{entry2, ZipEntryStored})
	zipEntrySetSize([]any{entry2, int64(5)})
	zipEntrySetCrc([]any{entry2, int64(crc)})
	zipOutputStreamPutNextEntry([]any{zos2, entry2})
	res = zipOutputStreamWriteByteArray([]any{zos2, javaByteArrayFromGoBytes([]byte("123456"))})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException for writing more than STORED size, got %v", res)
	}

	// 3. STORED entry closed with too few bytes
	baos3 := newTestByteArrayOutputStream()
	zos3 := newTestZipOutputStream(baos3)
	entry3, _ := CreateZipEntryObject("underflow.txt")
	zipEntrySetMethod([]any{entry3, ZipEntryStored})
	zipEntrySetSize([]any{entry3, int64(10)})
	zipEntrySetCrc([]any{entry3, int64(crc)})
	zipOutputStreamPutNextEntry([]any{zos3, entry3})
	zipOutputStreamWriteByteArray([]any{zos3, javaByteArrayFromGoBytes([]byte("12345"))})
	res = zipOutputStreamCloseEntry([]any{zos3})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException for closing STORED entry with underflow, got %v", res)
	}

	// 4. STORED entry closed with CRC mismatch
	baos4 := newTestByteArrayOutputStream()
	zos4 := newTestZipOutputStream(baos4)
	entry4, _ := CreateZipEntryObject("bad_crc.txt")
	zipEntrySetMethod([]any{entry4, ZipEntryStored})
	zipEntrySetSize([]any{entry4, int64(5)})
	zipEntrySetCrc([]any{entry4, int64(999999)})
	zipOutputStreamPutNextEntry([]any{zos4, entry4})
	zipOutputStreamWriteByteArray([]any{zos4, javaByteArrayFromGoBytes([]byte("12345"))})
	res = zipOutputStreamCloseEntry([]any{zos4})
	if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.ZipException {
		t.Fatalf("Expected ZipException for STORED entry CRC mismatch, got %v", res)
	}
}

func TestZipOutputStream_RoundTrip_MultiEntries(t *testing.T) {
	initZipOutputStreamTests()

	baos := newTestByteArrayOutputStream()
	zos := newTestZipOutputStream(baos)

	zipOutputStreamSetComment([]any{zos, object.StringObjectFromGoString("Archive-wide comment")})

	// Entry 1: Deflated file
	e1, _ := CreateZipEntryObject("dir/file1.txt")
	zipEntrySetComment([]any{e1, object.StringObjectFromGoString("Comment for file1")})
	extraData := []byte{0x01, 0x02, 0x03, 0x04}
	zipEntrySetExtra([]any{e1, javaByteArrayFromGoBytes(extraData)})
	zipOutputStreamPutNextEntry([]any{zos, e1})
	f1Content := []byte("The quick brown fox jumps over the lazy dog. Repeat " + strings.Repeat("hello ", 50))
	zipOutputStreamWriteByteArray([]any{zos, javaByteArrayFromGoBytes(f1Content)})

	// Entry 2: Stored file (auto closes entry 1)
	f2Content := []byte("Fixed uncompressed stored data")
	e2, _ := CreateZipEntryObject("dir/stored.dat")
	zipEntrySetMethod([]any{e2, ZipEntryStored})
	zipEntrySetSize([]any{e2, int64(len(f2Content))})
	zipEntrySetCrc([]any{e2, int64(crc32.ChecksumIEEE(f2Content))})
	zipOutputStreamPutNextEntry([]any{zos, e2})
	zipOutputStreamWriteByteArray([]any{zos, javaByteArrayFromGoBytes(f2Content)})

	// Entry 3: Empty file
	e3, _ := CreateZipEntryObject("dir/empty.txt")
	zipOutputStreamPutNextEntry([]any{zos, e3})

	// Finish and close
	resFinish := zipOutputStreamFinish([]any{zos})
	if resFinish != nil {
		t.Fatalf("Finish failed: %v", resFinish)
	}
	resClose := zipOutputStreamClose([]any{zos})
	if resClose != nil {
		t.Fatalf("Close failed: %v", resClose)
	}

	zipBytes := getByteArrayOutputStreamBytes(baos)
	if len(zipBytes) == 0 {
		t.Fatal("Zip bytes empty")
	}

	// 1. Verify with Go's archive/zip
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("archive/zip NewReader failed: %v", err)
	}
	if zr.Comment != "Archive-wide comment" {
		t.Fatalf("Expected archive comment 'Archive-wide comment', got '%s'", zr.Comment)
	}
	if len(zr.File) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(zr.File))
	}

	// File 1 verification
	if zr.File[0].Name != "dir/file1.txt" || zr.File[0].Comment != "Comment for file1" {
		t.Fatalf("File 1 header mismatch: %+v", zr.File[0].FileHeader)
	}
	if !bytes.Equal(zr.File[0].Extra, extraData) {
		t.Fatalf("File 1 extra mismatch: got %v, expected %v", zr.File[0].Extra, extraData)
	}
	rc1, err := zr.File[0].Open()
	if err != nil {
		t.Fatalf("Failed to open file1: %v", err)
	}
	var f1Read bytes.Buffer
	f1Read.ReadFrom(rc1)
	rc1.Close()
	if !bytes.Equal(f1Read.Bytes(), f1Content) {
		t.Fatalf("File 1 content mismatch")
	}

	// File 2 verification
	if zr.File[1].Name != "dir/stored.dat" || zr.File[1].Method != zip.Store {
		t.Fatalf("File 2 header mismatch: %+v", zr.File[1].FileHeader)
	}
	rc2, err := zr.File[1].Open()
	if err != nil {
		t.Fatalf("Failed to open file2: %v", err)
	}
	var f2Read bytes.Buffer
	f2Read.ReadFrom(rc2)
	rc2.Close()
	if !bytes.Equal(f2Read.Bytes(), f2Content) {
		t.Fatalf("File 2 content mismatch")
	}

	// File 3 verification
	if zr.File[2].Name != "dir/empty.txt" {
		t.Fatalf("File 3 header mismatch: %+v", zr.File[2].FileHeader)
	}
	rc3, err := zr.File[2].Open()
	if err != nil {
		t.Fatalf("Failed to open file3: %v", err)
	}
	var f3Read bytes.Buffer
	f3Read.ReadFrom(rc3)
	rc3.Close()
	if f3Read.Len() != 0 {
		t.Fatalf("File 3 should be empty, got %d bytes", f3Read.Len())
	}

	// 2. Verify with Jacobin's ZipInputStream
	bais := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(bais)

	// Entry 1
	e1ReadRes := zipInputStreamGetNextEntry([]any{zis})
	e1ReadObj, ok := e1ReadRes.(*object.Object)
	if !ok || e1ReadObj == nil {
		t.Fatalf("ZipInputStream failed to get entry 1: %v", e1ReadRes)
	}
	name1Obj := zipEntryGetName([]any{e1ReadObj}).(*object.Object)
	if object.GoStringFromStringObject(name1Obj) != "dir/file1.txt" {
		t.Fatalf("ZipInputStream entry 1 name mismatch: %s", object.GoStringFromStringObject(name1Obj))
	}
	e1BytesRes := zipInputStreamReadAllBytes([]any{zis})
	e1BytesObj, ok := e1BytesRes.(*object.Object)
	if !ok {
		t.Fatalf("ZipInputStream readAllBytes failed for entry 1: %v", e1BytesRes)
	}
	e1ReadData := object.GoByteArrayFromJavaByteArray(e1BytesObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(e1ReadData, f1Content) {
		t.Fatalf("ZipInputStream entry 1 data mismatch")
	}

	// Entry 2
	e2ReadRes := zipInputStreamGetNextEntry([]any{zis})
	e2ReadObj, ok := e2ReadRes.(*object.Object)
	if !ok || e2ReadObj == nil {
		t.Fatalf("ZipInputStream failed to get entry 2: %v", e2ReadRes)
	}
	name2Obj := zipEntryGetName([]any{e2ReadObj}).(*object.Object)
	if object.GoStringFromStringObject(name2Obj) != "dir/stored.dat" {
		t.Fatalf("ZipInputStream entry 2 name mismatch: %s", object.GoStringFromStringObject(name2Obj))
	}
	e2BytesRes := zipInputStreamReadAllBytes([]any{zis})
	e2BytesObj, ok := e2BytesRes.(*object.Object)
	if !ok {
		t.Fatalf("ZipInputStream readAllBytes failed for entry 2: %v", e2BytesRes)
	}
	e2ReadData := object.GoByteArrayFromJavaByteArray(e2BytesObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(e2ReadData, f2Content) {
		t.Fatalf("ZipInputStream entry 2 data mismatch")
	}

	// Entry 3
	e3ReadRes := zipInputStreamGetNextEntry([]any{zis})
	e3ReadObj, ok := e3ReadRes.(*object.Object)
	if !ok || e3ReadObj == nil {
		t.Fatalf("ZipInputStream failed to get entry 3: %v", e3ReadRes)
	}
	name3Obj := zipEntryGetName([]any{e3ReadObj}).(*object.Object)
	if object.GoStringFromStringObject(name3Obj) != "dir/empty.txt" {
		t.Fatalf("ZipInputStream entry 3 name mismatch: %s", object.GoStringFromStringObject(name3Obj))
	}

	// End of archive
	e4ReadRes := zipInputStreamGetNextEntry([]any{zis})
	if !object.IsNull(e4ReadRes) && e4ReadRes != nil {
		t.Fatalf("Expected null for end of archive, got %v", e4ReadRes)
	}

	zipInputStreamClose([]any{zis})
}

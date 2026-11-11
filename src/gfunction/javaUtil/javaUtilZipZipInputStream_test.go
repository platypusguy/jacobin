/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"hash/crc32"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaIo"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func initZipInputStreamTests() {
	globals.InitGlobals("test")
	javaIo.Load_Io_InputStream()
	javaIo.Load_Io_ByteArrayInputStream()
	javaIo.Load_Io_ByteArrayOutputStream()
	Load_Util_Zip_ZipEntry()
	Load_Util_Zip_ZipInputStream()
}

var testByteArrayInputStreamClassName = "java/io/ByteArrayInputStream"
var testByteArrayOutputStreamClassName = "java/io/ByteArrayOutputStream"

func newTestByteArrayInputStream(data []byte) *object.Object {
	jBytes := make([]types.JavaByte, len(data))
	for i, b := range data {
		jBytes[i] = types.JavaByte(b)
	}
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, jBytes)
	streamObj := object.MakeEmptyObjectWithClassName(&testByteArrayInputStreamClassName)
	javaIo.ByteArrayInputStreamInit([]any{streamObj, bufObj})
	return streamObj
}

func newTestByteArrayOutputStream() *object.Object {
	streamObj := object.MakeEmptyObjectWithClassName(&testByteArrayOutputStreamClassName)
	javaIo.ByteArrayOutputStreamInit([]any{streamObj})
	return streamObj
}

func getByteArrayOutputStreamBytes(baos *object.Object) []byte {
	res := javaIo.ByteArrayOutputStreamToByteArray([]any{baos})
	resObj, ok := res.(*object.Object)
	if !ok || resObj == nil {
		return nil
	}
	fld, ok := resObj.FieldTable["value"]
	if !ok || fld.Fvalue == nil {
		return nil
	}
	if jBytes, ok := fld.Fvalue.([]types.JavaByte); ok {
		return goBytesFromJavaBytes(jBytes)
	}
	return nil
}

func newTestZipInputStream(inObj *object.Object) *object.Object {
	zisObj := object.MakeEmptyObjectWithClassName(&zipInputStreamClassName)
	zipInputStreamInit([]any{zisObj, inObj})
	return zisObj
}

func goBytesFromJavaBytes(jBytes []types.JavaByte) []byte {
	b := make([]byte, len(jBytes))
	for i, jb := range jBytes {
		b[i] = byte(jb)
	}
	return b
}

func createSampleZipData(t *testing.T) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Entry 1: Deflated file
	f1, err := w.Create("file1.txt")
	if err != nil {
		t.Fatalf("Create file1.txt failed: %v", err)
	}
	if _, err := f1.Write([]byte("Hello, Deflated World!")); err != nil {
		t.Fatalf("Write file1.txt failed: %v", err)
	}

	// Entry 2: Stored file. Stored entries need sizes and CRC in the header,
	// so use CreateRaw to avoid a data descriptor.
	storedData := []byte("Hello, Stored World!")
	hdr2 := &zip.FileHeader{
		Name:               "stored.txt",
		Method:             zip.Store,
		CRC32:              crc32.ChecksumIEEE(storedData),
		CompressedSize64:   uint64(len(storedData)),
		UncompressedSize64: uint64(len(storedData)),
	}
	f2, err := w.CreateRaw(hdr2)
	if err != nil {
		t.Fatalf("Create stored.txt failed: %v", err)
	}
	if _, err := f2.Write(storedData); err != nil {
		t.Fatalf("Write stored.txt failed: %v", err)
	}

	// Entry 3: Directory
	hdr3 := &zip.FileHeader{
		Name:               "subfolder/",
		Method:             zip.Store,
		CRC32:              crc32.ChecksumIEEE(nil),
		CompressedSize64:   0,
		UncompressedSize64: 0,
	}
	if _, err := w.CreateRaw(hdr3); err != nil {
		t.Fatalf("Create subfolder/ failed: %v", err)
	}

	// Entry 4: Empty file
	hdr4 := &zip.FileHeader{
		Name:   "empty.bin",
		Method: zip.Deflate,
	}
	if _, err := w.CreateHeader(hdr4); err != nil {
		t.Fatalf("Create empty.bin failed: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("zip.Writer.Close failed: %v", err)
	}
	return buf.Bytes()
}

func TestZipInputStream_MethodRegistration(t *testing.T) {
	initZipInputStreamTests()

	expectedSlots := map[string]int{
		"java/util/zip/ZipInputStream.<clinit>()V":                                                0,
		"java/util/zip/ZipInputStream.<init>(Ljava/io/InputStream;)V":                             1,
		"java/util/zip/ZipInputStream.<init>(Ljava/io/InputStream;Ljava/nio/charset/Charset;)V":   2,
		"java/util/zip/ZipInputStream.available()I":                                               0,
		"java/util/zip/ZipInputStream.close()V":                                                   0,
		"java/util/zip/ZipInputStream.closeEntry()V":                                              0,
		"java/util/zip/ZipInputStream.createZipEntry(Ljava/lang/String;)Ljava/util/zip/ZipEntry;": 1,
		"java/util/zip/ZipInputStream.getNextEntry()Ljava/util/zip/ZipEntry;":                     0,
		"java/util/zip/ZipInputStream.mark(I)V":                                                   1,
		"java/util/zip/ZipInputStream.markSupported()Z":                                           0,
		"java/util/zip/ZipInputStream.read()I":                                                    0,
		"java/util/zip/ZipInputStream.read([B)I":                                                  1,
		"java/util/zip/ZipInputStream.read([BII)I":                                                3,
		"java/util/zip/ZipInputStream.readAllBytes()[B":                                           0,
		"java/util/zip/ZipInputStream.readNBytes(I)[B":                                            1,
		"java/util/zip/ZipInputStream.reset()V":                                                   0,
		"java/util/zip/ZipInputStream.skip(J)J":                                                   1,
		"java/util/zip/ZipInputStream.transferTo(Ljava/io/OutputStream;)J":                        1,
	}

	for sig, slots := range expectedSlots {
		meth, exists := ghelpers.MethodSignatures[sig]
		if !exists {
			t.Errorf("Method signature missing: %s", sig)
			continue
		}
		if meth.ParamSlots != slots {
			t.Errorf("Method %s ParamSlots = %d, want %d", sig, meth.ParamSlots, slots)
		}
		if meth.GFunction == nil {
			t.Errorf("Method %s GFunction is nil", sig)
		}
	}
}

func TestZipInputStream_ReadEntriesAndContent(t *testing.T) {
	initZipInputStreamTests()

	zipBytes := createSampleZipData(t)
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)

	// Entry 1: file1.txt
	e1 := zipInputStreamGetNextEntry([]any{zis})
	entry1, ok := e1.(*object.Object)
	if !ok || object.IsNull(entry1) {
		t.Fatalf("getNextEntry 1 failed: %v", e1)
	}
	name1 := object.GoStringFromStringObject(zipEntryGetName([]any{entry1}).(*object.Object))
	if name1 != "file1.txt" {
		t.Errorf("entry 1 name = %s, want 'file1.txt'", name1)
	}
	if avail := zipInputStreamAvailable([]any{zis}).(int64); avail != 1 {
		t.Errorf("available() = %d, want 1", avail)
	}

	// Read content of entry 1 using read([BII)I
	readBuf := make([]types.JavaByte, 64)
	readBufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, readBuf)
	n1 := zipInputStreamReadRange([]any{zis, readBufObj, int64(0), int64(len(readBuf))})
	nBytes, ok := n1.(int64)
	if !ok || nBytes != 22 {
		t.Fatalf("read entry 1 bytes count = %v, want 22", n1)
	}
	content1 := string(goBytesFromJavaBytes(readBuf[:nBytes]))
	if content1 != "Hello, Deflated World!" {
		t.Errorf("entry 1 content = %q, want 'Hello, Deflated World!'", content1)
	}

	// Read past EOF
	eofRes := zipInputStreamReadRange([]any{zis, readBufObj, int64(0), int64(len(readBuf))})
	if eofRes.(int64) != -1 {
		t.Errorf("read at EOF = %v, want -1", eofRes)
	}

	// Entry 2: stored.txt
	e2 := zipInputStreamGetNextEntry([]any{zis})
	entry2, ok := e2.(*object.Object)
	if !ok || object.IsNull(entry2) {
		t.Fatalf("getNextEntry 2 failed: %v", e2)
	}
	name2 := object.GoStringFromStringObject(zipEntryGetName([]any{entry2}).(*object.Object))
	if name2 != "stored.txt" {
		t.Errorf("entry 2 name = %s, want 'stored.txt'", name2)
	}

	// Read entry 2 using read() byte-by-byte
	var content2Bytes []byte
	for {
		bRes := zipInputStreamReadByte([]any{zis})
		bVal, ok := bRes.(int64)
		if !ok {
			t.Fatalf("read() returned error: %v", bRes)
		}
		if bVal == -1 {
			break
		}
		content2Bytes = append(content2Bytes, byte(bVal))
	}
	content2 := string(content2Bytes)
	if content2 != "Hello, Stored World!" {
		t.Errorf("entry 2 content = %q, want 'Hello, Stored World!'", content2)
	}

	// Entry 3: subfolder/ (directory)
	e3 := zipInputStreamGetNextEntry([]any{zis})
	entry3, ok := e3.(*object.Object)
	if !ok || object.IsNull(entry3) {
		t.Fatalf("getNextEntry 3 failed: %v", e3)
	}
	name3 := object.GoStringFromStringObject(zipEntryGetName([]any{entry3}).(*object.Object))
	if name3 != "subfolder/" {
		t.Errorf("entry 3 name = %s, want 'subfolder/'", name3)
	}
	if isDir := zipEntryIsDirectory([]any{entry3}).(int64); isDir != types.JavaBoolTrue {
		t.Errorf("subfolder/ isDirectory = %d, want true", isDir)
	}

	// Entry 4: empty.bin
	e4 := zipInputStreamGetNextEntry([]any{zis})
	entry4, ok := e4.(*object.Object)
	if !ok || object.IsNull(entry4) {
		t.Fatalf("getNextEntry 4 failed: %v", e4)
	}
	name4 := object.GoStringFromStringObject(zipEntryGetName([]any{entry4}).(*object.Object))
	if name4 != "empty.bin" {
		t.Errorf("entry 4 name = %s, want 'empty.bin'", name4)
	}

	// Read empty file
	emptyRead := zipInputStreamReadByte([]any{zis})
	if emptyRead.(int64) != -1 {
		t.Errorf("read empty file = %v, want -1", emptyRead)
	}

	// No more entries
	eEnd := zipInputStreamGetNextEntry([]any{zis})
	if !object.IsNull(eEnd) {
		t.Errorf("getNextEntry after end should be Null, got %v", eEnd)
	}

	// Close stream
	if res := zipInputStreamClose([]any{zis}); res != nil {
		t.Errorf("close returned error: %v", res)
	}
}

func TestZipInputStream_SkipAndCloseEntry(t *testing.T) {
	initZipInputStreamTests()

	zipBytes := createSampleZipData(t)
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)

	// Entry 1: skip part of content
	e1 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e1) {
		t.Fatalf("getNextEntry 1 failed")
	}

	skipped := zipInputStreamSkip([]any{zis, int64(7)})
	if skipped.(int64) != 7 {
		t.Errorf("skip(7) = %v, want 7", skipped)
	}

	// Read remaining 15 bytes ("Deflated World!")
	readBuf := make([]types.JavaByte, 32)
	readBufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, readBuf)
	n := zipInputStreamReadByteArray([]any{zis, readBufObj})
	nBytes := n.(int64)
	if nBytes != 15 {
		t.Fatalf("read after skip = %d, want 15", nBytes)
	}
	remStr := string(goBytesFromJavaBytes(readBuf[:nBytes]))
	if remStr != "Deflated World!" {
		t.Errorf("remaining content = %q, want 'Deflated World!'", remStr)
	}

	// Entry 2: closeEntry() without reading
	e2 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e2) {
		t.Fatalf("getNextEntry 2 failed")
	}
	if res := zipInputStreamCloseEntry([]any{zis}); res != nil {
		t.Errorf("closeEntry() error: %v", res)
	}

	// Entry 3: next entry should still be readable correctly
	e3 := zipInputStreamGetNextEntry([]any{zis})
	entry3 := e3.(*object.Object)
	name3 := object.GoStringFromStringObject(zipEntryGetName([]any{entry3}).(*object.Object))
	if name3 != "subfolder/" {
		t.Errorf("entry 3 name = %s, want 'subfolder/'", name3)
	}
}

func TestZipInputStream_CreateZipEntryAndMarkReset(t *testing.T) {
	initZipInputStreamTests()

	zipBytes := createSampleZipData(t)
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)

	// createZipEntry
	entryObjRes := zipInputStreamCreateZipEntry([]any{zis, object.StringObjectFromGoString("customEntry.txt")})
	entryObj, ok := entryObjRes.(*object.Object)
	if !ok || object.IsNull(entryObj) {
		t.Fatalf("createZipEntry failed: %v", entryObjRes)
	}
	if name := object.GoStringFromStringObject(zipEntryGetName([]any{entryObj}).(*object.Object)); name != "customEntry.txt" {
		t.Errorf("custom entry name = %s, want 'customEntry.txt'", name)
	}

	// markSupported -> false
	if ms := ghelpers.ReturnFalse([]any{zis}).(int64); ms != types.JavaBoolFalse {
		t.Errorf("markSupported = %d, want 0", ms)
	}

	// mark -> just returns
	if m := ghelpers.JustReturn([]any{zis, int64(10)}); m != nil {
		t.Errorf("mark returned %v, want nil", m)
	}

	// reset -> throws IOException("mark/reset not supported")
	res := zipInputStreamReset([]any{zis})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IOException {
		t.Errorf("reset() expected IOException, got: %v", res)
	}
}

func TestZipInputStream_RejectStoredEntryWithDataDescriptor(t *testing.T) {
	initZipInputStreamTests()

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	stored, err := w.CreateHeader(&zip.FileHeader{Name: "stored.txt", Method: zip.Store})
	if err != nil {
		t.Fatalf("Create stored.txt failed: %v", err)
	}
	if _, err := stored.Write([]byte("stored data containing PK\x03\x04 bytes")); err != nil {
		t.Fatalf("Write stored.txt failed: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zip.Writer.Close failed: %v", err)
	}

	zis := newTestZipInputStream(newTestByteArrayInputStream(buf.Bytes()))
	got := zipInputStreamGetNextEntry([]any{zis})
	errBlk, ok := got.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("getNextEntry = %v, want IOException", got)
	}
	if errBlk.ExceptionType != excNames.IOException {
		t.Errorf("getNextEntry exception = %v, want IOException", errBlk.ExceptionType)
	}
}

func TestZipInputStream_ErrorsAndExceptions(t *testing.T) {
	initZipInputStreamTests()

	// Null receiver / null constructor arguments
	rawObj := object.MakeEmptyObjectWithClassName(&zipInputStreamClassName)
	assertZipErr(t, zipInputStreamInit([]any{rawObj, nil}), excNames.NullPointerException, "init null in")
	assertZipErr(t, zipInputStreamInit([]any{rawObj, object.Null}), excNames.NullPointerException, "init Null in")
	assertZipErr(t, zipInputStreamInit([]any{rawObj, rawObj, nil}), excNames.NullPointerException, "init null charset")

	assertZipErr(t, zipInputStreamGetNextEntry([]any{nil}), excNames.NullPointerException, "getNextEntry null this")
	assertZipErr(t, zipInputStreamCloseEntry([]any{nil}), excNames.NullPointerException, "closeEntry null this")
	assertZipErr(t, zipInputStreamClose([]any{nil}), excNames.NullPointerException, "close null this")
	assertZipErr(t, zipInputStreamAvailable([]any{nil}), excNames.NullPointerException, "available null this")
	assertZipErr(t, zipInputStreamReadByte([]any{nil}), excNames.NullPointerException, "read byte null this")
	assertZipErr(t, zipInputStreamSkip([]any{nil, int64(10)}), excNames.NullPointerException, "skip null this")
	assertZipErr(t, zipInputStreamReset([]any{nil}), excNames.NullPointerException, "reset null this")
	assertZipErr(t, zipInputStreamCreateZipEntry([]any{nil, object.StringObjectFromGoString("a")}), excNames.NullPointerException, "createZipEntry null this")

	// Closed stream operations
	zipBytes := createSampleZipData(t)
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)
	zipInputStreamClose([]any{zis})

	assertZipErr(t, zipInputStreamGetNextEntry([]any{zis}), excNames.IOException, "getNextEntry on closed")
	assertZipErr(t, zipInputStreamAvailable([]any{zis}), excNames.IOException, "available on closed")
	assertZipErr(t, zipInputStreamCloseEntry([]any{zis}), excNames.IOException, "closeEntry on closed")
	assertZipErr(t, zipInputStreamReset([]any{zis}), excNames.IOException, "reset on closed")

	readBuf := make([]types.JavaByte, 16)
	readBufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, readBuf)
	assertZipErr(t, zipInputStreamReadRange([]any{zis, readBufObj, int64(0), int64(16)}), excNames.IOException, "read on closed")
	assertZipErr(t, zipInputStreamSkip([]any{zis, int64(5)}), excNames.IOException, "skip on closed")

	// Array bounds errors
	inStream2 := newTestByteArrayInputStream(zipBytes)
	zis2 := newTestZipInputStream(inStream2)
	zipInputStreamGetNextEntry([]any{zis2})

	assertZipErr(t, zipInputStreamReadRange([]any{zis2, readBufObj, int64(-1), int64(5)}), excNames.IndexOutOfBoundsException, "negative offset")
	assertZipErr(t, zipInputStreamReadRange([]any{zis2, readBufObj, int64(0), int64(-1)}), excNames.IndexOutOfBoundsException, "negative length")
	assertZipErr(t, zipInputStreamReadRange([]any{zis2, readBufObj, int64(10), int64(10)}), excNames.IndexOutOfBoundsException, "offset+len > cap")
	assertZipErr(t, zipInputStreamSkip([]any{zis2, int64(-5)}), excNames.IllegalArgumentException, "negative skip length")
}

func TestZipInputStream_InitWithCharset(t *testing.T) {
	initZipInputStreamTests()

	zipBytes := createSampleZipData(t)
	inStream := newTestByteArrayInputStream(zipBytes)

	charsetClassName := "java/nio/charset/Charset"
	charsetObj := object.MakeEmptyObjectWithClassName(&charsetClassName)

	zisObj := object.MakeEmptyObjectWithClassName(&zipInputStreamClassName)
	res := zipInputStreamInit([]any{zisObj, inStream, charsetObj})
	if res != nil {
		t.Fatalf("init with charset failed: %v", res)
	}

	csFld, ok := zisObj.FieldTable["charset"]
	if !ok || csFld.Fvalue != charsetObj {
		t.Errorf("charset field not set correctly: %v", csFld)
	}

	// Verify we can read entries normally
	e := zipInputStreamGetNextEntry([]any{zisObj})
	if object.IsNull(e) {
		t.Fatalf("getNextEntry failed with charset stream")
	}
}

func TestZipInputStream_CorruptedStream(t *testing.T) {
	initZipInputStreamTests()

	// Truncated header
	shortBytes := []byte{'P', 'K', 0x03, 0x04, 0x01, 0x02}
	inStream := newTestByteArrayInputStream(shortBytes)
	zis := newTestZipInputStream(inStream)
	assertZipErr(t, zipInputStreamGetNextEntry([]any{zis}), excNames.IOException, "truncated local header")
}

func TestZipInputStream_WrappedErrors(t *testing.T) {
	initZipInputStreamTests()

	geb := ghelpers.GetGErrBlk(excNames.IOException, "base io exception")
	wrapped := fmt.Errorf("wrapped error: %w", geb)

	unwrapped, ok := asGErrBlk(wrapped)
	if !ok || unwrapped == nil {
		t.Fatalf("asGErrBlk failed to unwrap wrapped GErrBlk")
	}
	if unwrapped.ExceptionType != excNames.IOException {
		t.Errorf("unwrapped exception type = %d, want %d", unwrapped.ExceptionType, excNames.IOException)
	}

	stdErr := fmt.Errorf("standard io error")
	unwrappedStd, okStd := asGErrBlk(stdErr)
	if !okStd || unwrappedStd == nil || unwrappedStd.ExceptionType != excNames.IOException {
		t.Errorf("asGErrBlk for standard error failed: %v", unwrappedStd)
	}
}

func TestZipInputStream_TransferTo(t *testing.T) {
	initZipInputStreamTests()

	zipBytes := createSampleZipData(t)
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)

	// Call transferTo before getNextEntry -> should transfer 0 bytes
	baosPre := newTestByteArrayOutputStream()
	resPre := zipInputStreamTransferTo([]any{zis, baosPre})
	if nPre, ok := resPre.(int64); !ok || nPre != 0 {
		t.Errorf("transferTo before getNextEntry returned %v, expected 0", resPre)
	}

	// Entry 1: Deflated file ("Hello, Deflated World!")
	e1 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e1) {
		t.Fatalf("getNextEntry failed for entry 1")
	}

	baos1 := newTestByteArrayOutputStream()
	res1 := ghelpers.Invoke("java/util/zip/ZipInputStream.transferTo(Ljava/io/OutputStream;)J", []any{zis, baos1})
	n1, ok1 := res1.(int64)
	if !ok1 || n1 != int64(len("Hello, Deflated World!")) {
		t.Fatalf("transferTo for entry 1 returned %v, expected %d", res1, len("Hello, Deflated World!"))
	}
	data1 := getByteArrayOutputStreamBytes(baos1)
	if string(data1) != "Hello, Deflated World!" {
		t.Errorf("entry 1 transferred content = %q, want 'Hello, Deflated World!'", string(data1))
	}

	// Draining entry 1: subsequent read on this entry should return -1 (EOF)
	readBuf := make([]types.JavaByte, 16)
	readBufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, readBuf)
	eofRead := zipInputStreamReadRange([]any{zis, readBufObj, int64(0), int64(16)})
	if eofRead.(int64) != -1 {
		t.Errorf("read after transferTo = %v, want -1 (EOF)", eofRead)
	}

	// Entry 2: Stored file ("Hello, Stored World!")
	e2 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e2) {
		t.Fatalf("getNextEntry failed for entry 2")
	}

	baos2 := newTestByteArrayOutputStream()
	res2 := zipInputStreamTransferTo([]any{zis, baos2})
	n2, ok2 := res2.(int64)
	if !ok2 || n2 != int64(len("Hello, Stored World!")) {
		t.Fatalf("transferTo for entry 2 returned %v, expected %d", res2, len("Hello, Stored World!"))
	}
	data2 := getByteArrayOutputStreamBytes(baos2)
	if string(data2) != "Hello, Stored World!" {
		t.Errorf("entry 2 transferred content = %q, want 'Hello, Stored World!'", string(data2))
	}

	// Entry 3: Directory (0 bytes)
	e3 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e3) {
		t.Fatalf("getNextEntry failed for entry 3")
	}
	baos3 := newTestByteArrayOutputStream()
	res3 := ghelpers.Invoke("java/util/zip/ZipInputStream.transferTo(Ljava/io/OutputStream;)J", []any{zis, baos3})
	if n3, ok := res3.(int64); !ok || n3 != 0 {
		t.Errorf("transferTo for directory entry returned %v, want 0", res3)
	}

	// Entry 4: Empty file (0 bytes)
	e4 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e4) {
		t.Fatalf("getNextEntry failed for entry 4")
	}
	baos4 := newTestByteArrayOutputStream()
	res4 := zipInputStreamTransferTo([]any{zis, baos4})
	if n4, ok := res4.(int64); !ok || n4 != 0 {
		t.Errorf("transferTo for empty entry returned %v, want 0", res4)
	}

	// EOF: no more entries
	eEOF := zipInputStreamGetNextEntry([]any{zis})
	if !object.IsNull(eEOF) {
		t.Errorf("expected null for EOF entry, got %v", eEOF)
	}

	// Error cases:
	// 1. Null destination stream throws NullPointerException
	inStreamNullTest := newTestByteArrayInputStream(zipBytes)
	zisNullTest := newTestZipInputStream(inStreamNullTest)
	zipInputStreamGetNextEntry([]any{zisNullTest})
	resNull := zipInputStreamTransferTo([]any{zisNullTest, object.Null})
	assertZipErr(t, resNull, excNames.NullPointerException, "transferTo null destination")

	// 2. Null receiver throws NullPointerException
	resNullReceiver := zipInputStreamTransferTo([]any{object.Null, baos1})
	assertZipErr(t, resNullReceiver, excNames.NullPointerException, "transferTo null receiver")

	// 3. Closed stream throws IOException
	inStreamClosedTest := newTestByteArrayInputStream(zipBytes)
	zisClosedTest := newTestZipInputStream(inStreamClosedTest)
	zipInputStreamClose([]any{zisClosedTest})
	resClosed := zipInputStreamTransferTo([]any{zisClosedTest, baos1})
	assertZipErr(t, resClosed, excNames.IOException, "transferTo on closed ZipInputStream")

	// 4. Large entry transfer (>8192 bytes)
	bufLarge := new(bytes.Buffer)
	wLarge := zip.NewWriter(bufLarge)
	fLarge, err := wLarge.Create("large.bin")
	if err != nil {
		t.Fatalf("Create large.bin failed: %v", err)
	}
	largeData := make([]byte, 25000)
	for i := range largeData {
		largeData[i] = byte((i * 37) % 256)
	}
	if _, err := fLarge.Write(largeData); err != nil {
		t.Fatalf("Write large.bin failed: %v", err)
	}
	if err := wLarge.Close(); err != nil {
		t.Fatalf("wLarge.Close failed: %v", err)
	}

	zisLarge := newTestZipInputStream(newTestByteArrayInputStream(bufLarge.Bytes()))
	zipInputStreamGetNextEntry([]any{zisLarge})
	baosLarge := newTestByteArrayOutputStream()
	resLarge := ghelpers.Invoke("java/util/zip/ZipInputStream.transferTo(Ljava/io/OutputStream;)J", []any{zisLarge, baosLarge})
	nLarge, okLarge := resLarge.(int64)
	if !okLarge || nLarge != int64(len(largeData)) {
		t.Fatalf("transferTo for large entry returned %v, expected %d", resLarge, len(largeData))
	}
	dataLarge := getByteArrayOutputStreamBytes(baosLarge)
	if !bytes.Equal(dataLarge, largeData) {
		t.Errorf("large entry transferred data mismatch, len got %d, want %d", len(dataLarge), len(largeData))
	}
}

func TestZipInputStream_ReadAllBytes(t *testing.T) {
	initZipInputStreamTests()

	// Create sample zip with small entries and a large entry (>8192 bytes)
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Entry 1: Deflated file
	f1, err := w.Create("file1.txt")
	if err != nil {
		t.Fatalf("Create file1.txt failed: %v", err)
	}
	content1Expected := []byte("Hello, Deflated World! Testing ZipInputStream.readAllBytes().")
	if _, err := f1.Write(content1Expected); err != nil {
		t.Fatalf("Write file1.txt failed: %v", err)
	}

	// Entry 2: Stored file with its metadata in the header, so no data descriptor is needed.
	content2Expected := []byte("Hello, Stored World! Testing stored entry readAllBytes.")
	hdr2 := &zip.FileHeader{
		Name:               "stored.txt",
		Method:             zip.Store,
		CRC32:              crc32.ChecksumIEEE(content2Expected),
		CompressedSize64:   uint64(len(content2Expected)),
		UncompressedSize64: uint64(len(content2Expected)),
	}
	f2, err := w.CreateRaw(hdr2)
	if err != nil {
		t.Fatalf("Create stored.txt failed: %v", err)
	}
	if _, err := f2.Write(content2Expected); err != nil {
		t.Fatalf("Write stored.txt failed: %v", err)
	}

	// Entry 3: Large file (>8192 bytes, e.g. 25000 bytes)
	f3, err := w.Create("large.txt")
	if err != nil {
		t.Fatalf("Create large.txt failed: %v", err)
	}
	largeData := make([]byte, 25000)
	for i := range largeData {
		largeData[i] = byte((i * 37) % 256)
	}
	if _, err := f3.Write(largeData); err != nil {
		t.Fatalf("Write large.txt failed: %v", err)
	}

	// Entry 4: Empty file
	hdr4 := &zip.FileHeader{
		Name:   "empty.bin",
		Method: zip.Deflate,
	}
	if _, err := w.CreateHeader(hdr4); err != nil {
		t.Fatalf("Create empty.bin failed: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("zip.Writer.Close failed: %v", err)
	}

	zipBytes := buf.Bytes()
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)

	// Call readAllBytes before getNextEntry -> should return empty byte array
	resPre := zipInputStreamReadAllBytes([]any{zis})
	resPreObj, okPre := resPre.(*object.Object)
	if !okPre || resPreObj == nil {
		t.Fatalf("readAllBytes before getNextEntry failed: %v", resPre)
	}
	if preBytes := goBytesFromJavaBytes(resPreObj.FieldTable["value"].Fvalue.([]types.JavaByte)); len(preBytes) != 0 {
		t.Errorf("readAllBytes before getNextEntry expected 0 bytes, got %d", len(preBytes))
	}

	// Entry 1: Deflated file
	e1 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e1) {
		t.Fatalf("getNextEntry 1 failed")
	}
	res1 := zipInputStreamReadAllBytes([]any{zis})
	res1Obj, ok1 := res1.(*object.Object)
	if !ok1 || res1Obj == nil {
		t.Fatalf("readAllBytes entry 1 failed: %v", res1)
	}
	data1 := goBytesFromJavaBytes(res1Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(data1, content1Expected) {
		t.Errorf("readAllBytes entry 1 got %q, want %q", data1, content1Expected)
	}

	// Read again on exhausted entry 1 -> should return empty byte array
	res1Again := zipInputStreamReadAllBytes([]any{zis})
	res1AgainObj, ok1Again := res1Again.(*object.Object)
	if !ok1Again || res1AgainObj == nil {
		t.Fatalf("readAllBytes on exhausted entry 1 failed: %v", res1Again)
	}
	if exhaustedBytes := goBytesFromJavaBytes(res1AgainObj.FieldTable["value"].Fvalue.([]types.JavaByte)); len(exhaustedBytes) != 0 {
		t.Errorf("readAllBytes on exhausted entry expected 0 bytes, got %d", len(exhaustedBytes))
	}

	// Entry 2: Stored file
	e2 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e2) {
		t.Fatalf("getNextEntry 2 failed")
	}
	res2 := zipInputStreamReadAllBytes([]any{zis})
	res2Obj, ok2 := res2.(*object.Object)
	if !ok2 || res2Obj == nil {
		t.Fatalf("readAllBytes entry 2 failed: %v", res2)
	}
	data2 := goBytesFromJavaBytes(res2Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(data2, content2Expected) {
		t.Errorf("readAllBytes entry 2 got %q, want %q", data2, content2Expected)
	}

	// Entry 3: Large file
	e3 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e3) {
		t.Fatalf("getNextEntry 3 failed")
	}
	res3 := zipInputStreamReadAllBytes([]any{zis})
	res3Obj, ok3 := res3.(*object.Object)
	if !ok3 || res3Obj == nil {
		t.Fatalf("readAllBytes entry 3 failed: %v", res3)
	}
	data3 := goBytesFromJavaBytes(res3Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(data3, largeData) {
		t.Errorf("readAllBytes entry 3 large data mismatch, len got %d, want %d", len(data3), len(largeData))
	}

	// Entry 4: Empty file
	e4 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e4) {
		t.Fatalf("getNextEntry 4 failed")
	}
	res4 := zipInputStreamReadAllBytes([]any{zis})
	res4Obj, ok4 := res4.(*object.Object)
	if !ok4 || res4Obj == nil {
		t.Fatalf("readAllBytes entry 4 failed: %v", res4)
	}
	data4 := goBytesFromJavaBytes(res4Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if len(data4) != 0 {
		t.Errorf("readAllBytes entry 4 got %d bytes, want 0", len(data4))
	}

	// EOF entry
	eEOF := zipInputStreamGetNextEntry([]any{zis})
	if !object.IsNull(eEOF) {
		t.Errorf("expected null for EOF entry, got %v", eEOF)
	}
	resEOF := zipInputStreamReadAllBytes([]any{zis})
	resEOFObj, okEOF := resEOF.(*object.Object)
	if !okEOF || resEOFObj == nil {
		t.Fatalf("readAllBytes at EOF failed: %v", resEOF)
	}
	if eofBytes := goBytesFromJavaBytes(resEOFObj.FieldTable["value"].Fvalue.([]types.JavaByte)); len(eofBytes) != 0 {
		t.Errorf("readAllBytes at EOF expected 0 bytes, got %d", len(eofBytes))
	}

	// Test invocation through ghelpers.Invoke
	inStreamInvoke := newTestByteArrayInputStream(zipBytes)
	zisInvoke := newTestZipInputStream(inStreamInvoke)
	zipInputStreamGetNextEntry([]any{zisInvoke})
	resInvoke := ghelpers.Invoke("java/util/zip/ZipInputStream.readAllBytes()[B", []any{zisInvoke})
	resInvokeObj, okInvoke := resInvoke.(*object.Object)
	if !okInvoke || resInvokeObj == nil {
		t.Fatalf("ghelpers.Invoke readAllBytes()[B failed: %v", resInvoke)
	}
	dataInvoke := goBytesFromJavaBytes(resInvokeObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(dataInvoke, content1Expected) {
		t.Errorf("ghelpers.Invoke readAllBytes()[B got %q, want %q", dataInvoke, content1Expected)
	}

	// Test NullPointerException for null receiver
	resNull := zipInputStreamReadAllBytes([]any{object.Null})
	assertZipErr(t, resNull, excNames.NullPointerException, "readAllBytes with null receiver")

	// Test IOException for closed stream
	inStreamClosed := newTestByteArrayInputStream(zipBytes)
	zisClosed := newTestZipInputStream(inStreamClosed)
	zipInputStreamClose([]any{zisClosed})
	resClosedZis := zipInputStreamReadAllBytes([]any{zisClosed})
	assertZipErr(t, resClosedZis, excNames.IOException, "readAllBytes on closed stream")
}

func TestZipInputStream_ReadNBytes(t *testing.T) {
	initZipInputStreamTests()

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Entry 1: Deflated file
	f1, err := w.Create("file1.txt")
	if err != nil {
		t.Fatalf("Create file1.txt failed: %v", err)
	}
	content1Expected := []byte("Hello, Deflated World! Testing ZipInputStream.readNBytes(I)[B.")
	if _, err := f1.Write(content1Expected); err != nil {
		t.Fatalf("Write file1.txt failed: %v", err)
	}

	// Entry 2: Large file (>8192 bytes, e.g. 20000 bytes)
	f2, err := w.Create("large.txt")
	if err != nil {
		t.Fatalf("Create large.txt failed: %v", err)
	}
	largeData := make([]byte, 20000)
	for i := range largeData {
		largeData[i] = byte((i * 41) % 256)
	}
	if _, err := f2.Write(largeData); err != nil {
		t.Fatalf("Write large.txt failed: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("zip.Writer.Close failed: %v", err)
	}

	zipBytes := buf.Bytes()
	inStream := newTestByteArrayInputStream(zipBytes)
	zis := newTestZipInputStream(inStream)

	// Call readNBytes before getNextEntry -> should return empty byte array
	resPre := zipInputStreamReadNBytes([]any{zis, int64(10)})
	resPreObj, okPre := resPre.(*object.Object)
	if !okPre || resPreObj == nil {
		t.Fatalf("readNBytes before getNextEntry failed: %v", resPre)
	}
	if preBytes := goBytesFromJavaBytes(resPreObj.FieldTable["value"].Fvalue.([]types.JavaByte)); len(preBytes) != 0 {
		t.Errorf("readNBytes before getNextEntry expected 0 bytes, got %d", len(preBytes))
	}

	// Entry 1: Deflated file
	e1 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e1) {
		t.Fatalf("getNextEntry 1 failed")
	}

	// Read 0 bytes -> returns byte[0]
	res0 := zipInputStreamReadNBytes([]any{zis, int64(0)})
	res0Obj, ok0 := res0.(*object.Object)
	if !ok0 || res0Obj == nil {
		t.Fatalf("readNBytes 0 bytes failed: %v", res0)
	}
	if b0 := goBytesFromJavaBytes(res0Obj.FieldTable["value"].Fvalue.([]types.JavaByte)); len(b0) != 0 {
		t.Errorf("readNBytes 0 bytes expected len 0, got %d", len(b0))
	}

	// Read first 5 bytes ("Hello")
	res5 := zipInputStreamReadNBytes([]any{zis, int64(5)})
	res5Obj, ok5 := res5.(*object.Object)
	if !ok5 || res5Obj == nil {
		t.Fatalf("readNBytes 5 bytes failed: %v", res5)
	}
	data5 := goBytesFromJavaBytes(res5Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(data5, content1Expected[:5]) {
		t.Errorf("readNBytes 5 bytes got %q, want %q", data5, content1Expected[:5])
	}

	// Read next 10 bytes (", Deflated")
	res10 := zipInputStreamReadNBytes([]any{zis, int64(10)})
	res10Obj, ok10 := res10.(*object.Object)
	if !ok10 || res10Obj == nil {
		t.Fatalf("readNBytes 10 bytes failed: %v", res10)
	}
	data10 := goBytesFromJavaBytes(res10Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(data10, content1Expected[5:15]) {
		t.Errorf("readNBytes 10 bytes got %q, want %q", data10, content1Expected[5:15])
	}

	// Read more bytes than remain (e.g. 1000 bytes) -> should read the rest
	resRest := zipInputStreamReadNBytes([]any{zis, int64(1000)})
	resRestObj, okRest := resRest.(*object.Object)
	if !okRest || resRestObj == nil {
		t.Fatalf("readNBytes rest failed: %v", resRest)
	}
	dataRest := goBytesFromJavaBytes(resRestObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(dataRest, content1Expected[15:]) {
		t.Errorf("readNBytes rest got %q, want %q", dataRest, content1Expected[15:])
	}

	// Read again on exhausted entry -> returns empty array
	resExhausted := zipInputStreamReadNBytes([]any{zis, int64(10)})
	resExhaustedObj, okExhausted := resExhausted.(*object.Object)
	if !okExhausted || resExhaustedObj == nil {
		t.Fatalf("readNBytes exhausted failed: %v", resExhausted)
	}
	if bExhausted := goBytesFromJavaBytes(resExhaustedObj.FieldTable["value"].Fvalue.([]types.JavaByte)); len(bExhausted) != 0 {
		t.Errorf("readNBytes on exhausted entry expected 0 bytes, got %d", len(bExhausted))
	}

	// Entry 2: Large file across multiple chunks (read 12000 bytes)
	e2 := zipInputStreamGetNextEntry([]any{zis})
	if object.IsNull(e2) {
		t.Fatalf("getNextEntry 2 failed")
	}

	resChunk1 := zipInputStreamReadNBytes([]any{zis, int64(12000)})
	resChunk1Obj, okChunk1 := resChunk1.(*object.Object)
	if !okChunk1 || resChunk1Obj == nil {
		t.Fatalf("readNBytes 12000 bytes failed: %v", resChunk1)
	}
	dataChunk1 := goBytesFromJavaBytes(resChunk1Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(dataChunk1, largeData[:12000]) {
		t.Errorf("readNBytes chunk1 mismatch, len got %d, want 12000", len(dataChunk1))
	}

	// Read remaining 8000 bytes via ghelpers.Invoke
	resChunk2 := ghelpers.Invoke("java/util/zip/ZipInputStream.readNBytes(I)[B", []any{zis, int64(8000)})
	resChunk2Obj, okChunk2 := resChunk2.(*object.Object)
	if !okChunk2 || resChunk2Obj == nil {
		t.Fatalf("ghelpers.Invoke readNBytes(I)[B failed: %v", resChunk2)
	}
	dataChunk2 := goBytesFromJavaBytes(resChunk2Obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(dataChunk2, largeData[12000:]) {
		t.Errorf("readNBytes chunk2 mismatch, len got %d, want 8000", len(dataChunk2))
	}

	// Test negative length -> IllegalArgumentException
	resNeg := zipInputStreamReadNBytes([]any{zis, int64(-1)})
	assertZipErr(t, resNeg, excNames.IllegalArgumentException, "readNBytes with negative length")

	// Test null receiver -> NullPointerException
	resNull := zipInputStreamReadNBytes([]any{object.Null, int64(5)})
	assertZipErr(t, resNull, excNames.NullPointerException, "readNBytes with null receiver")

	// Test closed stream -> IOException
	inStreamClosed := newTestByteArrayInputStream(zipBytes)
	zisClosed := newTestZipInputStream(inStreamClosed)
	zipInputStreamClose([]any{zisClosed})
	resClosedZis := zipInputStreamReadNBytes([]any{zisClosed, int64(5)})
	assertZipErr(t, resClosedZis, excNames.IOException, "readNBytes on closed stream")
}

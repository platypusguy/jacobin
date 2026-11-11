/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"strings"
	"testing"
	"time"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

func newTestZipEntry(name string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&zipEntryClassName)
	strObj := object.StringObjectFromGoString(name)
	zipEntryInitString([]any{obj, strObj})
	return obj
}

func assertZipErr(t *testing.T, got any, expectedExc int, msg string) {
	t.Helper()
	geb, ok := got.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("%s: expected *ghelpers.GErrBlk, got %T (%v)", msg, got, got)
	}
	if geb.ExceptionType != expectedExc {
		t.Fatalf("%s: expected exception type %d, got %d (msg: %s)", msg, expectedExc, geb.ExceptionType, geb.ErrMsg)
	}
}

func TestZipEntry_MethodRegistration(t *testing.T) {
	Load_Util_Zip_ZipEntry()

	expectedMethods := []struct {
		sig   string
		slots int
	}{
		{"java/util/zip/ZipEntry.<clinit>()V", 0},
		{"java/util/zip/ZipEntry.<init>(Ljava/lang/String;)V", 1},
		{"java/util/zip/ZipEntry.<init>(Ljava/util/zip/ZipEntry;)V", 1},
		{"java/util/zip/ZipEntry.getName()Ljava/lang/String;", 0},
		{"java/util/zip/ZipEntry.setTime(J)V", 1},
		{"java/util/zip/ZipEntry.getTime()J", 0},
		{"java/util/zip/ZipEntry.setTimeLocal(Ljava/time/LocalDateTime;)V", 1},
		{"java/util/zip/ZipEntry.getTimeLocal()Ljava/time/LocalDateTime;", 0},
		{"java/util/zip/ZipEntry.setLastModifiedTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;", 1},
		{"java/util/zip/ZipEntry.getLastModifiedTime()Ljava/nio/file/attribute/FileTime;", 0},
		{"java/util/zip/ZipEntry.setLastAccessTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;", 1},
		{"java/util/zip/ZipEntry.getLastAccessTime()Ljava/nio/file/attribute/FileTime;", 0},
		{"java/util/zip/ZipEntry.setCreationTime(Ljava/nio/file/attribute/FileTime;)Ljava/util/zip/ZipEntry;", 1},
		{"java/util/zip/ZipEntry.getCreationTime()Ljava/nio/file/attribute/FileTime;", 0},
		{"java/util/zip/ZipEntry.setSize(J)V", 1},
		{"java/util/zip/ZipEntry.getSize()J", 0},
		{"java/util/zip/ZipEntry.getCompressedSize()J", 0},
		{"java/util/zip/ZipEntry.setCompressedSize(J)V", 1},
		{"java/util/zip/ZipEntry.setCrc(J)V", 1},
		{"java/util/zip/ZipEntry.getCrc()J", 0},
		{"java/util/zip/ZipEntry.setMethod(I)V", 1},
		{"java/util/zip/ZipEntry.getMethod()I", 0},
		{"java/util/zip/ZipEntry.setExtra([B)V", 1},
		{"java/util/zip/ZipEntry.getExtra()[B", 0},
		{"java/util/zip/ZipEntry.setComment(Ljava/lang/String;)V", 1},
		{"java/util/zip/ZipEntry.getComment()Ljava/lang/String;", 0},
		{"java/util/zip/ZipEntry.isDirectory()Z", 0},
		{"java/util/zip/ZipEntry.toString()Ljava/lang/String;", 0},
		{"java/util/zip/ZipEntry.hashCode()I", 0},
		{"java/util/zip/ZipEntry.clone()Ljava/lang/Object;", 0},
		{"java/util/zip/ZipEntry.equals(Ljava/lang/Object;)Z", 1},
	}

	for _, m := range expectedMethods {
		meth, exists := ghelpers.MethodSignatures[m.sig]
		if !exists {
			t.Errorf("Method signature missing: %s", m.sig)
			continue
		}
		if meth.ParamSlots != m.slots {
			t.Errorf("Method %s ParamSlots = %d, want %d", m.sig, meth.ParamSlots, m.slots)
		}
		if meth.GFunction == nil {
			t.Errorf("Method %s GFunction is nil", m.sig)
		}
	}
}

func TestZipEntry_StaticConstants(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	zipEntryClinit(nil)

	sStored, ok := statics.Statics["java/util/zip/ZipEntry.STORED"]
	if !ok || sStored.Value.(int64) != 0 {
		t.Errorf("ZipEntry.STORED = %v, want 0", sStored.Value)
	}

	sDeflated, ok := statics.Statics["java/util/zip/ZipEntry.DEFLATED"]
	if !ok || sDeflated.Value.(int64) != 8 {
		t.Errorf("ZipEntry.DEFLATED = %v, want 8", sDeflated.Value)
	}
}

func TestZipEntry_ConstructorsAndName(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	// Normal construction
	e := newTestZipEntry("test/entry.txt")
	nameRes := zipEntryGetName([]any{e})
	strObj, ok := nameRes.(*object.Object)
	if !ok || object.GoStringFromStringObject(strObj) != "test/entry.txt" {
		t.Errorf("getName() = %v, want 'test/entry.txt'", nameRes)
	}

	// Initial default values
	if sz := zipEntryGetSize([]any{e}).(int64); sz != -1 {
		t.Errorf("initial getSize() = %d, want -1", sz)
	}
	if csz := zipEntryGetCompressedSize([]any{e}).(int64); csz != -1 {
		t.Errorf("initial getCompressedSize() = %d, want -1", csz)
	}
	if crc := zipEntryGetCrc([]any{e}).(int64); crc != -1 {
		t.Errorf("initial getCrc() = %d, want -1", crc)
	}
	if meth := zipEntryGetMethod([]any{e}).(int64); meth != -1 {
		t.Errorf("initial getMethod() = %d, want -1", meth)
	}
	if tm := zipEntryGetTime([]any{e}).(int64); tm != -1 {
		t.Errorf("initial getTime() = %d, want -1", tm)
	}
	if ext := zipEntryGetExtra([]any{e}); !object.IsNull(ext) {
		t.Errorf("initial getExtra() should be Null, got %v", ext)
	}
	if comm := zipEntryGetComment([]any{e}); !object.IsNull(comm) {
		t.Errorf("initial getComment() should be Null, got %v", comm)
	}

	// CreateZipEntryObject helper
	obj, errBlk := CreateZipEntryObject("myFile.dat")
	if errBlk != nil || obj == nil {
		t.Fatalf("CreateZipEntryObject failed: %v", errBlk)
	}
	if name := object.GoStringFromStringObject(zipEntryGetName([]any{obj}).(*object.Object)); name != "myFile.dat" {
		t.Errorf("CreateZipEntryObject name = %s, want 'myFile.dat'", name)
	}

	// Error handling: null name
	rawObj := object.MakeEmptyObjectWithClassName(&zipEntryClassName)
	assertZipErr(t, zipEntryInitString([]any{rawObj, nil}), excNames.NullPointerException, "init null name")
	assertZipErr(t, zipEntryInitString([]any{rawObj, object.Null}), excNames.NullPointerException, "init Null name")

	// Error handling: oversized name
	longName := strings.Repeat("A", 65536)
	longStrObj := object.StringObjectFromGoString(longName)
	assertZipErr(t, zipEntryInitString([]any{rawObj, longStrObj}), excNames.IllegalArgumentException, "init long name")
}

func TestZipEntry_CopyConstructor(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	src := newTestZipEntry("source.txt")
	zipEntrySetSize([]any{src, int64(1024)})
	zipEntrySetCompressedSize([]any{src, int64(512)})
	zipEntrySetCrc([]any{src, int64(0x12345678)})
	zipEntrySetMethod([]any{src, ZipEntryDeflated})
	zipEntrySetTime([]any{src, int64(1700000000000)})
	zipEntrySetComment([]any{src, object.StringObjectFromGoString("test comment")})

	extraBytes := []types.JavaByte{1, 2, 3, 4}
	extraObj := object.MakePrimitiveObject("[B", types.JavaByteArray, extraBytes)
	zipEntrySetExtra([]any{src, extraObj})

	// Copy constructor
	dst := object.MakeEmptyObjectWithClassName(&zipEntryClassName)
	res := zipEntryInitCopy([]any{dst, src})
	if res != nil {
		t.Fatalf("copy constructor returned error: %v", res)
	}

	// Verify all fields copied
	if name := object.GoStringFromStringObject(zipEntryGetName([]any{dst}).(*object.Object)); name != "source.txt" {
		t.Errorf("copy name = %s, want 'source.txt'", name)
	}
	if sz := zipEntryGetSize([]any{dst}).(int64); sz != 1024 {
		t.Errorf("copy size = %d, want 1024", sz)
	}
	if csz := zipEntryGetCompressedSize([]any{dst}).(int64); csz != 512 {
		t.Errorf("copy csize = %d, want 512", csz)
	}
	if crc := zipEntryGetCrc([]any{dst}).(int64); crc != 0x12345678 {
		t.Errorf("copy crc = %x, want 12345678", crc)
	}
	if m := zipEntryGetMethod([]any{dst}).(int64); m != ZipEntryDeflated {
		t.Errorf("copy method = %d, want %d", m, ZipEntryDeflated)
	}
	if tm := zipEntryGetTime([]any{dst}).(int64); tm != 1700000000000 {
		t.Errorf("copy time = %d, want 1700000000000", tm)
	}
	if comm := object.GoStringFromStringObject(zipEntryGetComment([]any{dst}).(*object.Object)); comm != "test comment" {
		t.Errorf("copy comment = %s, want 'test comment'", comm)
	}

	// Verify extra byte array is a cloned copy
	dstExtraObj := zipEntryGetExtra([]any{dst}).(*object.Object)
	dstBytes, ok := getByteArrayFromObject(dstExtraObj)
	if !ok || len(dstBytes) != 4 || dstBytes[0] != 1 {
		t.Fatalf("copy extra bytes invalid: %v", dstBytes)
	}

	// Null source entry check
	assertZipErr(t, zipEntryInitCopy([]any{dst, nil}), excNames.NullPointerException, "copy null source")
	assertZipErr(t, zipEntryInitCopy([]any{dst, object.Null}), excNames.NullPointerException, "copy Null source")
}

func TestZipEntry_TimeAndFileTime(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	e := newTestZipEntry("file.txt")

	// setTime / getTime
	zipEntrySetTime([]any{e, int64(1600000000123)})
	if tm := zipEntryGetTime([]any{e}).(int64); tm != 1600000000123 {
		t.Errorf("getTime() = %d, want 1600000000123", tm)
	}

	// getLastModifiedTime when set via setTime
	lmtObj := zipEntryGetLastModifiedTime([]any{e}).(*object.Object)
	if ms, ok := fileTimeToMillisHelper(lmtObj); !ok || ms != 1600000000123 {
		t.Errorf("getLastModifiedTime() ms = %d, want 1600000000123", ms)
	}

	// setLastModifiedTime
	ft := makeFileTimeFromMillis(1700000000999)
	setRes := zipEntrySetLastModifiedTime([]any{e, ft})
	if setRes != e {
		t.Errorf("setLastModifiedTime should return this, got %v", setRes)
	}
	if tm := zipEntryGetTime([]any{e}).(int64); tm != 1700000000999 {
		t.Errorf("after setLastModifiedTime, getTime() = %d, want 1700000000999", tm)
	}
	if gotFt := zipEntryGetLastModifiedTime([]any{e}).(*object.Object); gotFt != ft {
		t.Errorf("getLastModifiedTime() did not return set FileTime instance")
	}

	// setLastAccessTime / getLastAccessTime
	ftAccess := makeFileTimeFromMillis(1700000010000)
	zipEntrySetLastAccessTime([]any{e, ftAccess})
	if gotAtime := zipEntryGetLastAccessTime([]any{e}).(*object.Object); gotAtime != ftAccess {
		t.Errorf("getLastAccessTime() mismatch")
	}

	// setCreationTime / getCreationTime
	ftCreate := makeFileTimeFromMillis(1700000020000)
	zipEntrySetCreationTime([]any{e, ftCreate})
	if gotCtime := zipEntryGetCreationTime([]any{e}).(*object.Object); gotCtime != ftCreate {
		t.Errorf("getCreationTime() mismatch")
	}

	// setTimeLocal / getTimeLocal
	testEpochMs := time.Date(2026, 11, 10, 15, 30, 45, 0, time.Local).UnixMilli()
	ldt := localDateTimeFromMillis(testEpochMs)
	zipEntrySetTimeLocal([]any{e, ldt})

	gotLdt := zipEntryGetTimeLocal([]any{e}).(*object.Object)
	if gotMs, ok := millisFromLocalDateTime(gotLdt); !ok || gotMs != testEpochMs {
		t.Errorf("getTimeLocal() ms = %d, want %d", gotMs, testEpochMs)
	}

	// Error handling: null FileTime arguments
	assertZipErr(t, zipEntrySetLastModifiedTime([]any{e, nil}), excNames.NullPointerException, "setLastModifiedTime null")
	assertZipErr(t, zipEntrySetLastAccessTime([]any{e, nil}), excNames.NullPointerException, "setLastAccessTime null")
	assertZipErr(t, zipEntrySetCreationTime([]any{e, nil}), excNames.NullPointerException, "setCreationTime null")
	assertZipErr(t, zipEntrySetTimeLocal([]any{e, nil}), excNames.NullPointerException, "setTimeLocal null")
}

func TestZipEntry_SizeCompressedSizeCrcMethod(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	e := newTestZipEntry("data.bin")

	// Size
	zipEntrySetSize([]any{e, int64(4096)})
	if sz := zipEntryGetSize([]any{e}).(int64); sz != 4096 {
		t.Errorf("getSize() = %d, want 4096", sz)
	}
	assertZipErr(t, zipEntrySetSize([]any{e, int64(-5)}), excNames.IllegalArgumentException, "setSize negative")

	// CompressedSize
	zipEntrySetCompressedSize([]any{e, int64(1024)})
	if csz := zipEntryGetCompressedSize([]any{e}).(int64); csz != 1024 {
		t.Errorf("getCompressedSize() = %d, want 1024", csz)
	}

	// CRC
	zipEntrySetCrc([]any{e, int64(0xDEADBEEF)})
	if crc := zipEntryGetCrc([]any{e}).(int64); crc != 0xDEADBEEF {
		t.Errorf("getCrc() = %x, want DEADBEEF", crc)
	}
	assertZipErr(t, zipEntrySetCrc([]any{e, int64(-1)}), excNames.IllegalArgumentException, "setCrc negative")
	assertZipErr(t, zipEntrySetCrc([]any{e, int64(0x100000000)}), excNames.IllegalArgumentException, "setCrc > 0xFFFFFFFF")

	// Method
	zipEntrySetMethod([]any{e, ZipEntryStored})
	if m := zipEntryGetMethod([]any{e}).(int64); m != ZipEntryStored {
		t.Errorf("getMethod() = %d, want STORED (0)", m)
	}
	zipEntrySetMethod([]any{e, ZipEntryDeflated})
	if m := zipEntryGetMethod([]any{e}).(int64); m != ZipEntryDeflated {
		t.Errorf("getMethod() = %d, want DEFLATED (8)", m)
	}
	assertZipErr(t, zipEntrySetMethod([]any{e, int64(1)}), excNames.IllegalArgumentException, "setMethod invalid (1)")
	assertZipErr(t, zipEntrySetMethod([]any{e, int64(9)}), excNames.IllegalArgumentException, "setMethod invalid (9)")
}

func TestZipEntry_ExtraAndComment(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	e := newTestZipEntry("doc.pdf")

	// Extra bytes
	b := []types.JavaByte{10, 20, 30, 40, 50}
	extraObj := object.MakePrimitiveObject("[B", types.JavaByteArray, b)
	zipEntrySetExtra([]any{e, extraObj})

	gotExtra := zipEntryGetExtra([]any{e}).(*object.Object)
	bytes, ok := getByteArrayFromObject(gotExtra)
	if !ok || len(bytes) != 5 || bytes[2] != 30 {
		t.Errorf("getExtra() mismatch: %v", bytes)
	}

	// Null extra resets field
	zipEntrySetExtra([]any{e, nil})
	if ext := zipEntryGetExtra([]any{e}); !object.IsNull(ext) {
		t.Errorf("after setExtra(null), expected Null, got %v", ext)
	}

	// Comment
	commStr := "A test comment"
	zipEntrySetComment([]any{e, object.StringObjectFromGoString(commStr)})
	gotComm := object.GoStringFromStringObject(zipEntryGetComment([]any{e}).(*object.Object))
	if gotComm != commStr {
		t.Errorf("getComment() = %s, want '%s'", gotComm, commStr)
	}

	// Null comment resets field
	zipEntrySetComment([]any{e, nil})
	if comm := zipEntryGetComment([]any{e}); !object.IsNull(comm) {
		t.Errorf("after setComment(null), expected Null, got %v", comm)
	}

	// Oversized comment
	oversizedComm := object.StringObjectFromGoString(strings.Repeat("C", 65536))
	assertZipErr(t, zipEntrySetComment([]any{e, oversizedComm}), excNames.IllegalArgumentException, "oversized comment")
}

func TestZipEntry_IsDirectory(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	dirEntry := newTestZipEntry("META-INF/")
	if isDir := zipEntryIsDirectory([]any{dirEntry}).(int64); isDir != types.JavaBoolTrue {
		t.Errorf("META-INF/ isDirectory() = %d, want true", isDir)
	}

	fileEntry := newTestZipEntry("META-INF/MANIFEST.MF")
	if isDir := zipEntryIsDirectory([]any{fileEntry}).(int64); isDir != types.JavaBoolFalse {
		t.Errorf("MANIFEST.MF isDirectory() = %d, want false", isDir)
	}
}

func TestZipEntry_ToStringHashCodeCloneEquals(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	e := newTestZipEntry("test/file.txt")

	// toString
	str := object.GoStringFromStringObject(zipEntryToString([]any{e}).(*object.Object))
	if str != "test/file.txt" {
		t.Errorf("toString() = %s, want 'test/file.txt'", str)
	}

	// hashCode
	h := zipEntryHashCode([]any{e}).(int64)
	expectedH := int64(javaStringHashCode("test/file.txt"))
	if h != expectedH {
		t.Errorf("hashCode() = %d, want %d", h, expectedH)
	}

	// clone
	zipEntrySetSize([]any{e, int64(12345)})
	cloneRes := zipEntryClone([]any{e}).(*object.Object)
	if cloneRes == e {
		t.Errorf("clone() returned same pointer")
	}
	if name := object.GoStringFromStringObject(zipEntryGetName([]any{cloneRes}).(*object.Object)); name != "test/file.txt" {
		t.Errorf("cloned entry name = %s, want 'test/file.txt'", name)
	}
	if sz := zipEntryGetSize([]any{cloneRes}).(int64); sz != 12345 {
		t.Errorf("cloned entry size = %d, want 12345", sz)
	}

	// equals
	if eq := zipEntryEquals([]any{e, e}).(int64); eq != types.JavaBoolTrue {
		t.Errorf("e.equals(e) = %d, want true", eq)
	}
	if eq := zipEntryEquals([]any{e, cloneRes}).(int64); eq != types.JavaBoolFalse {
		t.Errorf("e.equals(cloneRes) = %d, want false", eq)
	}
}

func TestZipEntry_NullReceiverAndErrorHandling(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Zip_ZipEntry()

	assertZipErr(t, zipEntryGetName([]any{nil}), excNames.NullPointerException, "getName null this")
	assertZipErr(t, zipEntryGetSize([]any{nil}), excNames.NullPointerException, "getSize null this")
	assertZipErr(t, zipEntryGetTime([]any{nil}), excNames.NullPointerException, "getTime null this")
	assertZipErr(t, zipEntrySetSize([]any{nil, int64(10)}), excNames.NullPointerException, "setSize null this")
	assertZipErr(t, zipEntrySetTime([]any{nil, int64(10)}), excNames.NullPointerException, "setTime null this")
	assertZipErr(t, zipEntryClone([]any{nil}), excNames.NullPointerException, "clone null this")
	assertZipErr(t, zipEntryHashCode([]any{nil}), excNames.NullPointerException, "hashCode null this")
	assertZipErr(t, zipEntryIsDirectory([]any{nil}), excNames.NullPointerException, "isDirectory null this")
}

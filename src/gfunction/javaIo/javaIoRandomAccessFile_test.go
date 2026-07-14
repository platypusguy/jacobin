package javaIo

import (
	"bytes"
	"encoding/binary"
	"io"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"testing"
)

func TestRafOpen0(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_open0_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	rafObj := newRAFObject()
	pathStrObj := object.StringObjectFromGoString(tmpFile.Name())

	// Test read mode (O_RDONLY = 1)
	params := []interface{}{rafObj, pathStrObj, int64(1)}
	ret := rafOpen0(params)
	if ret != nil {
		t.Fatalf("rafOpen0 (read) returned error: %v", ret)
	}
	fld, _ := rafObj.FieldTable[ghelpers.FileHandle]
	fh := fld.Fvalue.(*os.File)
	fh.Close()

	// Test read-write mode (O_RDWR = 2)
	params = []interface{}{rafObj, pathStrObj, int64(2)}
	ret = rafOpen0(params)
	if ret != nil {
		t.Fatalf("rafOpen0 (read-write) returned error: %v", ret)
	}
	fld, _ = rafObj.FieldTable[ghelpers.FileHandle]
	fh = fld.Fvalue.(*os.File)
	fh.Close()
}

// helper to create a new RandomAccessFile object with initialized FieldTable
func newRAFObject() *object.Object {
	return &object.Object{FieldTable: make(map[string]object.Field)}
}

func TestClinitGeneric(t *testing.T) {
	ret := ghelpers.ClinitGeneric(nil)
	if ret != nil {
		t.Errorf("ghelpers.ClinitGeneric should return nil, got %v", ret)
	}
}

func TestJustReturn(t *testing.T) {
	ret := ghelpers.JustReturn(nil)
	if ret != nil {
		t.Errorf("ghelpers.JustReturn should return nil, got %v", ret)
	}
}

func TestRafInitStringAndGetFilePointer(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rafObj := newRAFObject()

	pathStrObj := object.StringObjectFromGoString(tmpFile.Name())
	modeStrObj := object.StringObjectFromGoString("r")

	params := []interface{}{rafObj, pathStrObj, modeStrObj}
	ret := rafInitString(params)
	if ret != nil {
		t.Fatalf("rafInitString returned error: %v", ret)
	}

	fld, ok := rafObj.FieldTable[ghelpers.FilePath]
	if !ok {
		t.Fatal("ghelpers.FilePath field not set")
	}
	gotPath := string(fld.Fvalue.([]byte))
	if gotPath != tmpFile.Name() {
		t.Fatalf("ghelpers.FilePath mismatch, want %s, got %s", tmpFile.Name(), gotPath)
	}

	fld, ok = rafObj.FieldTable[ghelpers.FileHandle]
	if !ok {
		t.Fatal("ghelpers.FileHandle field not set")
	}
	fh, ok := fld.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("ghelpers.FileHandle field has wrong type %T", fld.Fvalue)
	}

	getPointerParams := []interface{}{rafObj}
	pos := rafGetFilePointer(getPointerParams)
	offset, ok := pos.(int64)
	if !ok {
		t.Fatalf("rafGetFilePointer returned wrong type %T", pos)
	}
	if offset != 0 {
		t.Errorf("Initial file pointer expected 0, got %d", offset)
	}

	fh.Close()
}

func TestRafInitFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	fileObj.FieldTable[ghelpers.FilePath] = object.Field{
		Ftype:  types.JavaByteArray,
		Fvalue: object.JavaByteArrayFromGoString(tmpFile.Name()),
	}

	rafObj := newRAFObject()

	modeStrObj := object.StringObjectFromGoString("r")

	params := []interface{}{rafObj, fileObj, modeStrObj}
	ret := rafInitFile(params)
	if ret != nil {
		t.Fatalf("rafInitFile returned error: %v", ret)
	}

	fld, ok := rafObj.FieldTable[ghelpers.FilePath]
	if !ok {
		t.Fatal("ghelpers.FilePath field not set")
	}
	gotPath := string(fld.Fvalue.([]byte))
	if gotPath != tmpFile.Name() {
		t.Fatalf("ghelpers.FilePath mismatch, want %s, got %s", tmpFile.Name(), gotPath)
	}

	fld, ok = rafObj.FieldTable[ghelpers.FileHandle]
	if !ok {
		t.Fatal("ghelpers.FileHandle field not set")
	}
	fh, ok := fld.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("ghelpers.FileHandle field has wrong type %T", fld.Fvalue)
	}

	fh.Close()
}

func TestFisClose(t *testing.T) {
	rafObj := newRAFObject()

	// Set ghelpers.FileHandle with a pipe writer to avoid closing os.Stdout accidentally
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: w}

	ret := fisClose([]interface{}{rafObj})
	if ret != nil {
		t.Errorf("fisClose returned error: %v", ret)
	}

	// Writing after close should fail
	_, err = w.Write([]byte("test"))
	if err == nil {
		t.Errorf("Write succeeded after close, expected failure")
	}
}

func TestFisReadOne(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_read_one")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte{0x42}
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	ret := fisReadOne([]interface{}{rafObj})

	intRet, ok := ret.(int64)
	if !ok {
		t.Fatalf("fisReadOne returned wrong type %T", ret)
	}
	if intRet != int64(content[0]) {
		t.Errorf("fisReadOne expected %d, got %d", content[0], intRet)
	}
}

func TestFisReadByteArray(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_ba")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	byteArray := make([]types.JavaByte, len(content))

	javaByteArrayObj := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray)
	params := []interface{}{rafObj, javaByteArrayObj}
	ret := fisReadByteArray(params)

	numRead, ok := ret.(int64)
	if !ok {
		t.Fatalf("fisReadByteArray returned wrong type %T", ret)
	}
	if numRead != int64(len(content)) {
		t.Errorf("fisReadByteArray expected read %d bytes, got %d", len(content), numRead)
	}
}

func TestFisReadByteArrayOffset(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_ba_offset")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello world")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	byteArray := make([]types.JavaByte, len(content))
	javaByteArrayObj := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray)

	offset := int64(2)
	length := int64(5)

	params := []interface{}{rafObj, javaByteArrayObj, offset, length}
	ret := fisReadByteArrayOffset(params)

	numRead, ok := ret.(int64)
	if !ok {
		t.Fatalf("fisReadByteArrayOffset returned wrong type %T", ret)
	}
	if numRead != length {
		t.Errorf("fisReadByteArrayOffset expected read %d bytes, got %d", length, numRead)
	}
}

func TestRafReadFully(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_fully")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello readFully")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	// Case 1: Read fully successfully
	byteArray := make([]types.JavaByte, 5)
	javaByteArrayObj := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray)
	params := []interface{}{rafObj, javaByteArrayObj}
	ret := rafReadFully(params)

	if ret != nil {
		t.Errorf("rafReadFully failed: %v", ret)
	}

	readVal := javaByteArrayObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	goReadVal := object.GoByteArrayFromJavaByteArray(readVal)
	if !bytes.Equal(goReadVal, content[:5]) {
		t.Errorf("Expected %v, got %v", content[:5], goReadVal)
	}

	// Case 2: Read fully to the end
	byteArray2 := make([]types.JavaByte, len(content)-5)
	javaByteArrayObj2 := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray2)
	params2 := []interface{}{rafObj, javaByteArrayObj2}
	ret2 := rafReadFully(params2)

	if ret2 != nil {
		t.Errorf("rafReadFully failed at step 2: %v", ret2)
	}

	readVal2 := javaByteArrayObj2.FieldTable["value"].Fvalue.([]types.JavaByte)
	goReadVal2 := object.GoByteArrayFromJavaByteArray(readVal2)
	if !bytes.Equal(goReadVal2, content[5:]) {
		t.Errorf("Expected %v, got %v", content[5:], goReadVal2)
	}

	// Case 3: Read fully past EOF - should return error
	byteArray3 := make([]types.JavaByte, 1)
	javaByteArrayObj3 := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray3)
	params3 := []interface{}{rafObj, javaByteArrayObj3}
	ret3 := rafReadFully(params3)

	if ret3 == nil {
		t.Errorf("Expected error for readFully past EOF, got nil")
	}
}

func TestRafSetLength(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_set_length")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello setLength")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	// Case 1: Shorten the file
	newLen := int64(5)
	params := []interface{}{rafObj, newLen}
	ret := rafSetLength(params)

	if ret != nil {
		t.Errorf("rafSetLength failed: %v", ret)
	}

	fi, _ := tmpFile.Stat()
	if fi.Size() != newLen {
		t.Errorf("Expected size %d, got %d", newLen, fi.Size())
	}

	// Case 2: Lengthen the file
	newLen = int64(20)
	params = []interface{}{rafObj, newLen}
	ret = rafSetLength(params)

	if ret != nil {
		t.Errorf("rafSetLength failed: %v", ret)
	}

	fi, _ = tmpFile.Stat()
	if fi.Size() != newLen {
		t.Errorf("Expected size %d, got %d", newLen, fi.Size())
	}
}

func TestRafLengthAndSeek(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_length_seek")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello length and seek")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	// Test length()
	ret := rafLength([]interface{}{rafObj})
	if ret.(int64) != int64(len(content)) {
		t.Errorf("Expected length %d, got %d", len(content), ret)
	}

	// Test seek()
	rafSeek([]interface{}{rafObj, int64(6)})
	pos := rafGetFilePointer([]interface{}{rafObj})
	if pos.(int64) != 6 {
		t.Errorf("Expected position 6, got %d", pos)
	}
}

func TestRafReadFullyOffset(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_fully_offset")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("0123456789")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	byteArray := make([]types.JavaByte, 10)
	javaByteArrayObj := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray)

	// Read 4 bytes starting from index 2 in the array
	params := []interface{}{rafObj, javaByteArrayObj, int64(2), int64(4)}
	ret := rafReadFullyOffset(params)

	if ret != nil {
		t.Errorf("rafReadFullyOffset failed: %v", ret)
	}

	readVal := javaByteArrayObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if readVal[2] != types.JavaByte('0') || readVal[5] != types.JavaByte('3') {
		t.Errorf("Unexpected read value: %v", readVal)
	}
}

func TestRafWriteMethods(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_write_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	// Test write(I)V
	rafWrite([]interface{}{rafObj, int64('A')})

	// Test write([B)V
	byteArray := []types.JavaByte{types.JavaByte('B'), types.JavaByte('C')}
	javaByteArrayObj := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray)
	rafWriteByteArray([]interface{}{rafObj, javaByteArrayObj})

	// Test write([BII)V
	byteArray2 := []types.JavaByte{types.JavaByte('X'), types.JavaByte('D'), types.JavaByte('E'), types.JavaByte('Y')}
	javaByteArrayObj2 := object.MakePrimitiveObject(types.JavaByteArray, types.JavaByteArray, byteArray2)
	rafWriteByteArrayOffset([]interface{}{rafObj, javaByteArrayObj2, int64(1), int64(2)})

	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)
	written, _ := io.ReadAll(tmpFile)
	expected := []byte("ABCDE")
	if !bytes.Equal(written, expected) {
		t.Errorf("Expected %s, got %s", expected, written)
	}
}

func TestRafReadDataMethods(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_data_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Prepare data
	// Boolean(true), Byte(0x12), Char('A'), Double(1.23), Float(4.56), Int(0x12345678), Long(0x1122334455667788), Short(0x1234), UnsignedByte(0xFE), UnsignedShort(0xFEDC)
	var buf bytes.Buffer
	buf.WriteByte(1)                                    // Boolean true
	buf.WriteByte(0x12)                                 // Byte
	binary.Write(&buf, binary.BigEndian, uint16('A'))   // Char
	binary.Write(&buf, binary.BigEndian, 1.23)          // Double
	binary.Write(&buf, binary.BigEndian, float32(4.56)) // Float
	binary.Write(&buf, binary.BigEndian, int32(0x12345678))
	binary.Write(&buf, binary.BigEndian, int64(0x1122334455667788))
	binary.Write(&buf, binary.BigEndian, int16(0x1234))
	buf.WriteByte(0xFE)                                  // UnsignedByte
	binary.Write(&buf, binary.BigEndian, uint16(0xFEDC)) // UnsignedShort
	buf.WriteString("line1\nline2\rline3\r\n")           // readLine data

	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	params := []interface{}{rafObj}

	// readBoolean
	if res := rafReadBoolean(params); res.(int64) != 1 {
		t.Errorf("readBoolean: expected 1, got %v", res)
	}
	// readByte
	if res := rafReadByte(params); res.(int64) != 0x12 {
		t.Errorf("readByte: expected 0x12, got %v", res)
	}
	// readChar
	if res := rafReadChar(params); res.(int64) != int64('A') {
		t.Errorf("readChar: expected %d, got %v", int64('A'), res)
	}
	// readDouble
	if res := rafReadDouble(params); res.(float64) != 1.23 {
		t.Errorf("readDouble: expected 1.23, got %v", res)
	}
	// readFloat
	if res := rafReadFloat(params); float32(res.(float64)) != 4.56 {
		t.Errorf("readFloat: expected 4.56, got %v", res)
	}
	// readInt
	if res := rafReadInt(params); res.(int64) != 0x12345678 {
		t.Errorf("readInt: expected 0x12345678, got %v", res)
	}
	// readLong
	if res := rafReadLong(params); res.(int64) != 0x1122334455667788 {
		t.Errorf("readLong: expected 0x1122334455667788, got %v", res)
	}
	// readShort
	if res := rafReadShort(params); res.(int64) != 0x1234 {
		t.Errorf("readShort: expected 0x1234, got %v", res)
	}
	// readUnsignedByte
	if res := rafReadUnsignedByte(params); res.(int64) != 0xFE {
		t.Errorf("readUnsignedByte: expected 0xFE, got %v", res)
	}
	// readUnsignedShort
	if res := rafReadUnsignedShort(params); res.(int64) != 0xFEDC {
		t.Errorf("readUnsignedShort: expected 0xFEDC, got %v", res)
	}

	// readLine
	if res := rafReadLine(params); object.GoStringFromStringObject(res.(*object.Object)) != "line1" {
		t.Errorf("readLine1: expected 'line1', got '%s'", object.GoStringFromStringObject(res.(*object.Object)))
	}
	if res := rafReadLine(params); object.GoStringFromStringObject(res.(*object.Object)) != "line2" {
		t.Errorf("readLine2: expected 'line2', got '%s'", object.GoStringFromStringObject(res.(*object.Object)))
	}
	if res := rafReadLine(params); object.GoStringFromStringObject(res.(*object.Object)) != "line3" {
		t.Errorf("readLine3: expected 'line3', got '%s'", object.GoStringFromStringObject(res.(*object.Object)))
	}
}

func TestRafReadUTF(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_utf_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Prepare data: 2-byte length, then modified UTF-8
	// Modified UTF-8 for \u0000 is 0xC0 0x80
	// But let's use what Java would produce.
	// "Hello, 世界! " is standard UTF-8.
	// \u0000 is 0xC0 0x80.
	utfBytes := append([]byte("Hello, 世界! "), 0xC0, 0x80)
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint16(len(utfBytes)))
	buf.Write(utfBytes)

	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	res := rafReadUTF([]interface{}{rafObj})
	gotStr := object.GoStringFromStringObject(res.(*object.Object))
	expectedStr := "Hello, 世界! \u0000"
	if gotStr != expectedStr {
		t.Errorf("readUTF: expected %q, got %q", expectedStr, gotStr)
	}
}

func TestRafWriteDataMethods(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_write_data_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rafObj := newRAFObject()
	rafObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: tmpFile}

	// writeBoolean
	rafWriteBoolean([]interface{}{rafObj, int64(1)})
	// writeByte
	rafWriteByte([]interface{}{rafObj, int64(0x12)})
	// writeShort
	rafWriteShort([]interface{}{rafObj, int64(0x1234)})
	// writeChar
	rafWriteChar([]interface{}{rafObj, int64('A')})
	// writeInt
	rafWriteInt([]interface{}{rafObj, int64(0x12345678)})
	// writeLong
	rafWriteLong([]interface{}{rafObj, int64(0x1122334455667788)})
	// writeFloat
	rafWriteFloat([]interface{}{rafObj, 4.56})
	// writeDouble
	rafWriteDouble([]interface{}{rafObj, 1.23})
	// writeBytes
	rafWriteBytes([]interface{}{rafObj, object.StringObjectFromGoString("abc")})
	// writeChars
	rafWriteChars([]interface{}{rafObj, object.StringObjectFromGoString("ABC")})
	// writeUTF
	rafWriteUTF([]interface{}{rafObj, object.StringObjectFromGoString("Hello, 世界! \u0000")})

	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	// Now read back and verify
	params := []interface{}{rafObj}

	// readBoolean
	if res := rafReadBoolean(params); res.(int64) != 1 {
		t.Errorf("readBoolean: expected 1, got %v", res)
	}
	// readByte
	if res := rafReadByte(params); res.(int64) != 0x12 {
		t.Errorf("readByte: expected 0x12, got %v", res)
	}
	// readShort
	if res := rafReadShort(params); res.(int64) != 0x1234 {
		t.Errorf("readShort: expected 0x1234, got %v", res)
	}
	// readChar
	if res := rafReadChar(params); res.(int64) != int64('A') {
		t.Errorf("readChar: expected %d, got %v", int64('A'), res)
	}
	// readInt
	if res := rafReadInt(params); res.(int64) != 0x12345678 {
		t.Errorf("readInt: expected 0x12345678, got %v", res)
	}
	// readLong
	if res := rafReadLong(params); res.(int64) != 0x1122334455667788 {
		t.Errorf("readLong: expected 0x1122334455667788, got %v", res)
	}
	// readFloat
	if res := rafReadFloat(params); float32(res.(float64)) != 4.56 {
		t.Errorf("readFloat: expected 4.56, got %v", res)
	}
	// readDouble
	if res := rafReadDouble(params); res.(float64) != 1.23 {
		t.Errorf("readDouble: expected 1.23, got %v", res)
	}

	// writeBytes verification (read back 3 bytes)
	b3 := make([]byte, 3)
	io.ReadFull(tmpFile, b3)
	if string(b3) != "abc" {
		t.Errorf("writeBytes: expected 'abc', got %q", string(b3))
	}

	// writeChars verification (read back 3 chars = 6 bytes)
	b6 := make([]byte, 6)
	io.ReadFull(tmpFile, b6)
	if string(b6) != "\x00A\x00B\x00C" {
		t.Errorf("writeChars: unexpected content %v", b6)
	}

	// writeUTF verification
	res := rafReadUTF(params)
	gotStr := object.GoStringFromStringObject(res.(*object.Object))
	expectedStr := "Hello, 世界! \u0000"
	if gotStr != expectedStr {
		t.Errorf("writeUTF/readUTF: expected %q, got %q", expectedStr, gotStr)
	}
}

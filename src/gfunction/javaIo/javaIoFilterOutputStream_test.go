package javaIo

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"testing"
)

func TestFilterOutputStream_MethodRegistration(t *testing.T) {
	globals.InitStringPool()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Io_FilterOutputStream()

	cases := []struct {
		key   string
		slots int
	}{
		{"java/io/FilterOutputStream.<clinit>()V", 0},
		{"java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V", 1},
		{"java/io/FilterOutputStream.close()V", 0},
		{"java/io/FilterOutputStream.flush()V", 0},
		{"java/io/FilterOutputStream.write(I)V", 1},
		{"java/io/FilterOutputStream.write([B)V", 1},
		{"java/io/FilterOutputStream.write([BII)V", 3},
	}
	for _, c := range cases {
		gm, ok := ghelpers.MethodSignatures[c.key]
		if !ok {
			t.Fatalf("method not registered: %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("ParamSlots mismatch for %s: want %d got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("GFunction nil for %s", c.key)
		}
	}
}

func TestFilterOutputStream_Init_Write_Flush_Close(t *testing.T) {
	globals.InitStringPool()
	// Prepare temp file and underlying OutputStream object using existing helper
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "fostream_test1.bin")
	defer os.Remove(filePath)
	under := makeOutputStreamObjForFile(t, filePath)

	// target FilterOutputStream object
	target := object.MakeEmptyObject()
	if res := filteroutputstreamInit([]interface{}{target, under}); res != nil {
		t.Fatalf("filteroutputstreamInit error: %v", res)
	}
	// write one byte 'Z'
	if res := filteroutputstreamWrite([]interface{}{target, int64('Z')}); res != nil {
		t.Fatalf("write(int) error: %v", res)
	}
	// write byte array [0x41,0x42]
	buf := []types.JavaByte{types.JavaByte('A'), types.JavaByte('B')}
	arr := &object.Object{FieldTable: map[string]object.Field{"value": {Ftype: types.ByteArray, Fvalue: buf}}}
	if res := filteroutputstreamWriteBytes([]interface{}{target, arr}); res != nil {
		t.Fatalf("write(byte[]) error: %v", res)
	}
	// write range from ["hello"], offset 1, len 3 => "ell"
	str := object.StringObjectFromGoString("hello")
	if res := filteroutputstreamWriteBytesRange([]interface{}{target, str, int64(1), int64(3)}); res != nil {
		t.Fatalf("write(byte[],off,len) error: %v", res)
	}
	// flush
	if res := filteroutputstreamFlush([]interface{}{target}); res != nil {
		t.Fatalf("flush error: %v", res)
	}
	// verify file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != "ZABell" {
		t.Fatalf("content mismatch: got %q want %q", string(data), "ZABell")
	}
	// close
	if res := filteroutputstreamClose([]interface{}{target}); res != nil {
		t.Fatalf("close error: %v", res)
	}
}

func TestFilterOutputStream_WriteRange_ParamError(t *testing.T) {
	// Out-of-bounds should raise IndexOutOfBoundsException
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "fostream_test2.bin")
	defer os.Remove(filePath)
	under := makeOutputStreamObjForFile(t, filePath)
	target := object.MakeEmptyObject()
	_ = filteroutputstreamInit([]interface{}{target, under})

	str := object.StringObjectFromGoString("xyz")
	res := filteroutputstreamWriteBytesRange([]interface{}{target, str, int64(2), int64(5)})
	if res == nil {
		t.Fatalf("expected error for out-of-bounds write range")
	}
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	} else {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	_ = filteroutputstreamClose([]interface{}{target})
}

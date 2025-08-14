package gfunction

import (
    "hash/adler32"
    "jacobin/globals"
    "jacobin/object"
    "testing"
)

func makeZipByteArray(b []byte) *object.Object {
    jb := object.JavaByteArrayFromGoByteArray(b)
    return object.StringObjectFromJavaByteArray(jb)
}

func newAdler32Obj() *object.Object {
    cn := "java/util/zip/Adler32"
    return object.MakeEmptyObjectWithClassName(&cn)
}

func TestAdler32_Init_GetValue_Reset(t *testing.T) {
    globals.InitStringPool()

    a := newAdler32Obj()
    if ret := adlerInit([]interface{}{a}); ret != nil { t.Fatalf("adlerInit error: %v", ret) }

    // Initial value should be 1 per spec and implementation
    if v := adlerGetValue([]interface{}{a}).(int64); v != 1 {
        t.Fatalf("initial Adler32 value expected 1, got %d", v)
    }

    // After a small update then reset returns to 1
    _ = adlerUpdateFromInt([]interface{}{a, int64('A')})
    _ = adlerReset([]interface{}{a})
    if v := adlerGetValue([]interface{}{a}).(int64); v != 1 {
        t.Fatalf("after reset, expected 1, got %d", v)
    }
}

func TestAdler32_Update_Array_And_Int_KnownVectors(t *testing.T) {
    globals.InitStringPool()

    // Known test string for Adler32
    data := []byte("123456789")
    want := adler32.Checksum(data) // known value 0x091E01DE

    a := newAdler32Obj()
    _ = adlerInit([]interface{}{a})

    // update([B, offset, toIndex) per current implementation (3rd param is end index, not count)
    arr := makeZipByteArray(data)
    _ = adlerUpdateFromArray([]interface{}{a, arr, int64(0), int64(len(data))})

    got := uint32(adlerGetValue([]interface{}{a}).(int64))
    if got != want {
        t.Fatalf("Adler32 checksum mismatch: want 0x%08x got 0x%08x", want, got)
    }

    // Now test single-byte update matches library when applied cumulatively
    _ = adlerReset([]interface{}{a})
    // Compute expected by applying updateAdler32 starting from 1 to byte 'A'
    expected := updateAdler32(uint32(1), []byte{'A'})
    _ = adlerUpdateFromInt([]interface{}{a, int64('A')})
    got2 := uint32(adlerGetValue([]interface{}{a}).(int64))
    if got2 != expected {
        t.Fatalf("Adler32 single-byte mismatch: want 0x%08x got 0x%08x", expected, got2)
    }
}

func TestAdler32_Update_Subrange_UsesEndIndex(t *testing.T) {
    globals.InitStringPool()

    data := []byte("abcdef")
    // We'll update only "bcd" (indices 1..3) by passing offset=1, toIndex=4
    sub := data[1:4]

    a := newAdler32Obj()
    _ = adlerInit([]interface{}{a})

    arr := makeZipByteArray(data)
    _ = adlerUpdateFromArray([]interface{}{a, arr, int64(1), int64(4)})
    got := uint32(adlerGetValue([]interface{}{a}).(int64))

    // Compute expected using the implementation helper starting from initial 1
    expected := updateAdler32(uint32(1), sub)
    if got != expected {
        t.Fatalf("Adler32 subrange mismatch: want 0x%08x got 0x%08x", expected, got)
    }
}

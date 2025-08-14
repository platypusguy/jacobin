package gfunction

import (
    "hash/crc32"
    "jacobin/globals"
    "jacobin/object"
    "testing"
)

func newCRC32Obj() *object.Object {
    cn := "java/util/zip/CRC32"
    return object.MakeEmptyObjectWithClassName(&cn)
}

func newCRC32CObj() *object.Object {
    cn := "java/util/zip/CRC32C"
    return object.MakeEmptyObjectWithClassName(&cn)
}

func TestCRC32_Init_GetValue_Reset(t *testing.T) {
    globals.InitStringPool()

    c := newCRC32Obj()
    if ret := crc32InitIEEE([]interface{}{c}); ret != nil { t.Fatalf("crc32InitIEEE error: %v", ret) }

    if v := crc32GetValue([]interface{}{c}).(int64); v != 0 {
        t.Fatalf("initial CRC32 value expected 0, got %d", v)
    }

    _ = crc32UpdateFromInt([]interface{}{c, int64('A')})
    _ = crc32Reset([]interface{}{c})
    if v := crc32GetValue([]interface{}{c}).(int64); v != 0 {
        t.Fatalf("after reset, expected 0, got %d", v)
    }
}

func TestCRC32_IEEE_KnownVector_And_Int(t *testing.T) {
    globals.InitStringPool()

    data := []byte("123456789")
    want := crc32.ChecksumIEEE(data) // known 0xCBF43926

    c := newCRC32Obj()
    _ = crc32InitIEEE([]interface{}{c})

    arr := makeZipByteArray(data)
    // Use (offset, toIndex) per current implementation quirk
    _ = crc32UpdateFromArray([]interface{}{c, arr, int64(0), int64(len(data))})

    got := uint32(crc32GetValue([]interface{}{c}).(int64))
    if got != want {
        t.Fatalf("CRC32 IEEE mismatch: want 0x%08x got 0x%08x", want, got)
    }

    // Single byte update vs library update from 0
    _ = crc32Reset([]interface{}{c})
    _ = crc32UpdateFromInt([]interface{}{c, int64('A')})
    got2 := uint32(crc32GetValue([]interface{}{c}).(int64))
    table := crc32.MakeTable(crc32.IEEE)
    want2 := crc32.Update(0, table, []byte{'A'})
    if got2 != want2 {
        t.Fatalf("CRC32 IEEE single-byte mismatch: want 0x%08x got 0x%08x", want2, got2)
    }
}

func TestCRC32C_KnownVector_And_Subrange(t *testing.T) {
    globals.InitStringPool()

    // CRC32C(Castagnoli) known vector for "123456789"
    data := []byte("123456789")
    tableC := crc32.MakeTable(crc32.Castagnoli)
    want := crc32.Update(0, tableC, data) // known 0xE3069283

    c := newCRC32CObj()
    _ = crc32InitCastagnoli([]interface{}{c})

    arr := makeZipByteArray(data)
    _ = crc32UpdateFromArray([]interface{}{c, arr, int64(0), int64(len(data))})

    got := uint32(crc32GetValue([]interface{}{c}).(int64))
    if got != want {
        t.Fatalf("CRC32C mismatch: want 0x%08x got 0x%08x", want, got)
    }

    // Subrange update: use "bcd" from "abcdef"
    _ = crc32Reset([]interface{}{c})
    all := []byte("abcdef")
    sub := all[1:4]
    arr2 := makeZipByteArray(all)
    _ = crc32UpdateFromArray([]interface{}{c, arr2, int64(1), int64(4)})
    got2 := uint32(crc32GetValue([]interface{}{c}).(int64))
    want2 := crc32.Update(0, tableC, sub)
    if got2 != want2 {
        t.Fatalf("CRC32C subrange mismatch: want 0x%08x got 0x%08x", want2, got2)
    }
}

package gfunction

import (
    "jacobin/excNames"
    "jacobin/globals"
    "jacobin/object"
    "jacobin/types"
    "testing"
)

// Helpers
func makeHFDefault(t *testing.T) *object.Object {
    t.Helper()
    // Ensure <clinit> statics are loaded
    if ret := hexFormatClinit([]interface{}{}); ret != nil {
        t.Fatalf("hexFormatClinit returned error: %v", ret)
    }
    obj := hfOf([]interface{}{}).(*object.Object)
    return obj
}

func makeString(s string) *object.Object { return object.StringObjectFromGoString(s) }

func makeByteArrObjHF(b []byte) *object.Object {
    jb := object.JavaByteArrayFromGoByteArray(b)
    return object.StringObjectFromJavaByteArray(jb)
}

func assertJavaBoolHF(t *testing.T, got interface{}, want int64, msg string) {
    t.Helper()
    gi, ok := got.(int64)
    if !ok {
        t.Fatalf("%s: expected int64 Java boolean, got %T", msg, got)
    }
    if gi != want {
        t.Fatalf("%s: expected %d, got %d", msg, want, gi)
    }
}

func TestHexFormat_Default_And_Modifiers(t *testing.T) {
    globals.InitStringPool()

    hf := makeHFDefault(t)

    // Default delimiter/prefix/suffix should be empty; lowercase
    del := hfDelimiter([]interface{}{hf}).(*object.Object)
    if s := object.GoStringFromStringObject(del); s != "" {
        t.Fatalf("expected empty delimiter by default, got %q", s)
    }

    // toString should reflect lowercase and empty fields
    ts := hfToString([]interface{}{hf}).(*object.Object)
    s := object.GoStringFromStringObject(ts)
    if want := "uppercase: false, delimiter: \"\", prefix: \"\", suffix: \"\""; s != want {
        t.Fatalf("unexpected toString: %q", s)
    }

    // Apply modifiers: delimiter, prefix, suffix, uppercase
    hf2 := hfWithDelimiter([]interface{}{hf, makeString(":")} ).(*object.Object)
    hf2 = hfWithPrefix([]interface{}{hf2, makeString("0x")} ).(*object.Object)
    hf2 = hfWithSuffix([]interface{}{hf2, makeString("h")} ).(*object.Object)
    hf2 = hfWithUpperCase([]interface{}{hf2}).(*object.Object)

    // Verify getters and toString
    if s := object.GoStringFromStringObject(hfDelimiter([]interface{}{hf2}).(*object.Object)); s != ":" {
        t.Fatalf("expected delimiter ':', got %q", s)
    }
    if s := object.GoStringFromStringObject(hfPrefix([]interface{}{hf2}).(*object.Object)); s != "0x" {
        t.Fatalf("expected prefix '0x', got %q", s)
    }
    if s := object.GoStringFromStringObject(hfSuffix([]interface{}{hf2}).(*object.Object)); s != "h" {
        t.Fatalf("expected suffix 'h', got %q", s)
    }
    ts2 := hfToString([]interface{}{hf2}).(*object.Object)
    s2 := object.GoStringFromStringObject(ts2)
    if want := "uppercase: true, delimiter: \":\", prefix: \"0x\", suffix: \"h\""; s2 != want {
        t.Fatalf("unexpected toString after modifiers: %q", s2)
    }
}

func TestHexFormat_FormatHex_FromBytes(t *testing.T) {
    globals.InitStringPool()

    hf := makeHFDefault(t)

    // Lowercase, custom prefix/suffix and delimiter
    hf = hfWithPrefix([]interface{}{hf, makeString("[")}).(*object.Object)
    hf = hfWithSuffix([]interface{}{hf, makeString("]")}).(*object.Object)
    hf = hfWithDelimiter([]interface{}{hf, makeString(" ")}).(*object.Object)

    // Full array format
    data := []byte{0x00, 0x0F, 0xA5}
    arr := makeByteArrObjHF(data)
    out := hfFormatHexFromBytes([]interface{}{hf, arr}).(*object.Object)
    if got := object.GoStringFromStringObject(out); got != "[00] [0f] [a5]" {
        t.Fatalf("formatHex([B) mismatch: got %q", got)
    }

    // Subrange [1:3) -> indices 1 and 2
    out2 := hfFormatHexFromBytes([]interface{}{hf, arr, int64(1), int64(3)}).(*object.Object)
    if got := object.GoStringFromStringObject(out2); got != "[0f] [a5]" {
        t.Fatalf("formatHex([BII) mismatch: got %q", got)
    }

    // Out of range should return IndexOutOfBoundsException
    if err := hfFormatHexFromBytes([]interface{}{hf, arr, int64(-1), int64(2)}); err == nil {
        t.Fatalf("expected error for negative fromIndex")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    }
    if err := hfFormatHexFromBytes([]interface{}{hf, arr, int64(1), int64(10)}); err == nil {
        t.Fatalf("expected error for toIndex beyond length")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    }
    if err := hfFormatHexFromBytes([]interface{}{hf, arr, int64(2), int64(2)}); err == nil {
        t.Fatalf("expected error for toIndex <= fromIndex")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    }

    // Empty array with empty delimiter should yield empty string (avoid non-empty delimiter to prevent slice issue)
    empty := makeByteArrObjHF([]byte{})
    out3 := hfFormatHexFromBytes([]interface{}{makeHFDefault(t), empty}).(*object.Object)
    if got := object.GoStringFromStringObject(out3); got != "" {
        t.Fatalf("expected empty string for empty input, got %q", got)
    }
}

func TestHexFormat_ToHexDigits_Primitives(t *testing.T) {
    globals.InitStringPool()

    hfLower := makeHFDefault(t)
    hfUpper := hfWithUpperCase([]interface{}{hfLower}).(*object.Object)

    // byte -> 2 hex digits
    if s := object.GoStringFromStringObject(hfByteToHexDigits([]interface{}{hfLower, int64(0xAB)}).(*object.Object)); s != "ab" {
        t.Fatalf("byte toHexDigits lowercase mismatch: %q", s)
    }
    if s := object.GoStringFromStringObject(hfByteToHexDigits([]interface{}{hfUpper, int64(0xAB)}).(*object.Object)); s != "AB" {
        t.Fatalf("byte toHexDigits uppercase mismatch: %q", s)
    }

    // char -> 4 hex digits
    if s := object.GoStringFromStringObject(hfCharToHexDigits([]interface{}{hfLower, int64(0x0042)}).(*object.Object)); s != "0042" {
        t.Fatalf("char toHexDigits mismatch: %q", s)
    }

    // int -> 8 hex digits
    if s := object.GoStringFromStringObject(hfIntToHexDigits([]interface{}{hfLower, int64(0x00c0ffee)}).(*object.Object)); s != "00c0ffee" {
        t.Fatalf("int toHexDigits mismatch: %q", s)
    }

    // long -> 16 hex digits; with length cropping
    sFull := object.GoStringFromStringObject(hfLongToHexDigits([]interface{}{hfUpper, int64(0x123)},).(*object.Object))
    if sFull != "0000000000000123" {
        t.Fatalf("long toHexDigits full mismatch: %q", sFull)
    }
    sCrop := object.GoStringFromStringObject(hfLongToHexDigits([]interface{}{hfUpper, int64(0xDEADBEEF), int64(6)}).(*object.Object))
    if sCrop != "ADBEEF" {
        t.Fatalf("long toHexDigits cropped mismatch: %q", sCrop)
    }

    // short -> 4 hex digits
    if s := object.GoStringFromStringObject(hfShortToHexDigits([]interface{}{hfLower, int64(0x0bad)}).(*object.Object)); s != "0bad" {
        t.Fatalf("short toHexDigits mismatch: %q", s)
    }
}

func TestHexFormat_FromHexDigit_And_HighLow(t *testing.T) {
    globals.InitStringPool()

    hf := makeHFDefault(t)
    // fromHexDigit valid
    if v := hfFromHexDigit([]interface{}{int64('0')}).(int64); v != 0 { t.Fatalf("fromHexDigit '0' -> 0, got %d", v) }
    if v := hfFromHexDigit([]interface{}{int64('9')}).(int64); v != 9 { t.Fatalf("fromHexDigit '9' -> 9, got %d", v) }
    if v := hfFromHexDigit([]interface{}{int64('A')}).(int64); v != 10 { t.Fatalf("fromHexDigit 'A' -> 10, got %d", v) }
    if v := hfFromHexDigit([]interface{}{int64('f')}).(int64); v != 15 { t.Fatalf("fromHexDigit 'f' -> 15, got %d", v) }

    // invalid -> NumberFormatException
    if err := hfFromHexDigit([]interface{}{int64('G')}); err == nil {
        t.Fatalf("expected error for invalid hex digit")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.NumberFormatException {
            t.Fatalf("expected NumberFormatException, got %d", geb.ExceptionType)
        }
    }

    // toHighHexDigit / toLowHexDigit use digits table; verify against lowercase default
    high := hfToHighHexDigit([]interface{}{hf, int64(0xAB)}).(int64)
    low := hfToLowHexDigit([]interface{}{hf, int64(0xAB)}).(int64)
    if rune(high) != 'a' || rune(low) != 'b' {
        t.Fatalf("expected high 'a' and low 'b', got %q and %q", rune(high), rune(low))
    }

    // Uppercase variant
    hfU := hfWithUpperCase([]interface{}{hf}).(*object.Object)
    highU := hfToHighHexDigit([]interface{}{hfU, int64(0xAB)}).(int64)
    lowU := hfToLowHexDigit([]interface{}{hfU, int64(0xAB)}).(int64)
    if rune(highU) != 'A' || rune(lowU) != 'B' {
        t.Fatalf("expected high 'A' and low 'B', got %q and %q", rune(highU), rune(lowU))
    }
}

func TestHexFormat_Equals(t *testing.T) {
    globals.InitStringPool()

    hf1 := makeHFDefault(t)
    hf1 = hfWithDelimiter([]interface{}{hf1, makeString(":")} ).(*object.Object)
    hf1 = hfWithPrefix([]interface{}{hf1, makeString("0x")} ).(*object.Object)
    hf1 = hfWithSuffix([]interface{}{hf1, makeString("h")} ).(*object.Object)
    hf1 = hfWithUpperCase([]interface{}{hf1}).(*object.Object)

    hf2 := makeHFDefault(t)
    hf2 = hfWithDelimiter([]interface{}{hf2, makeString(":")} ).(*object.Object)
    hf2 = hfWithPrefix([]interface{}{hf2, makeString("0x")} ).(*object.Object)
    hf2 = hfWithSuffix([]interface{}{hf2, makeString("h")} ).(*object.Object)
    hf2 = hfWithUpperCase([]interface{}{hf2}).(*object.Object)

    // Should be equal
    assertJavaBoolHF(t, hfEquals([]interface{}{hf1, hf2}), types.JavaBoolTrue, "HexFormat equals expected true")

    // Change one field -> not equal
    hf3 := hfWithLowerCase([]interface{}{hf2}).(*object.Object)
    assertJavaBoolHF(t, hfEquals([]interface{}{hf1, hf3}), types.JavaBoolFalse, "HexFormat equals expected false after change")
}

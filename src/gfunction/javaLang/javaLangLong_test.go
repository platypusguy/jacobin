package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestLong_ValueOf_And_DoubleValue(t *testing.T) {
	// valueOf(J) -> Long object, then doubleValue() -> primitive double
	obj := longValueOf([]interface{}{int64(42)})
	lobj, ok := obj.(*object.Object)
	if !ok {
		t.Fatalf("valueOf did not return object, got %T", obj)
	}
	out := longDoubleValue([]interface{}{lobj})
	if got := out.(float64); got != 42.0 {
		t.Fatalf("doubleValue mismatch: got %v want %v", got, 42.0)
	}
}

func TestLong_ParseLong_Valid_And_Invalid(t *testing.T) {
	// NOTE: fixed implementation reads params[0] for the string argument
	s := object.StringObjectFromGoString("12345")
	out := longParseLong([]interface{}{s})
	if got := out.(int64); got != 12345 {
		t.Fatalf("parseLong valid: got %d", got)
	}
	// invalid -> NumberFormatException
	sinv := object.StringObjectFromGoString("abc")
	out = longParseLong([]interface{}{sinv})
	if geb, ok := out.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NumberFormatException {
		if !ok {
			t.Fatalf("parseLong invalid: expected *ghelpers.GErrBlk, got %T", out)
		}
		t.Fatalf("parseLong invalid: expected NumberFormatException, got %v", geb)
	}
}

func TestLong_RotateLeft_Right(t *testing.T) {
	// Use a known pattern and compare by unsigned bit pattern
	val := int64(0x0123456789abcdef)
	// rotate left by 8 -> 0x23456789abcdef01
	outL := longRotateLeft([]interface{}{val, int64(8)})
	gotL := uint64(outL.(int64))
	var expectedL uint64 = 0x23456789abcdef01
	if gotL != expectedL {
		t.Fatalf("rotateLeft mismatch: got 0x%016x want 0x%016x", gotL, expectedL)
	}
	// rotate right by 12 -> expected using uint64 rotation
	outR := longRotateRight([]interface{}{val, int64(12)})
	gotR := uint64(outR.(int64))
	// compute expected separately without constant overflow
	var base uint64 = 0x0123456789abcdef
	expectedR := (base >> 12) | (base << (64 - 12))
	if gotR != expectedR {
		t.Fatalf("rotateRight mismatch: got 0x%016x want 0x%016x", gotR, expectedR)
	}
}

func TestLong_Decode(t *testing.T) {
	// hex
	out := longDecode([]interface{}{object.StringObjectFromGoString("0x123")})
	if got := out.(*object.Object).FieldTable["value"].Fvalue.(int64); got != 0x123 {
		t.Fatalf("decode hex got %d", got)
	}
	// octal
	out = longDecode([]interface{}{object.StringObjectFromGoString("0123")})
	if got := out.(*object.Object).FieldTable["value"].Fvalue.(int64); got != 0123 {
		t.Fatalf("decode octal got %d", got)
	}
	// decimal
	out = longDecode([]interface{}{object.StringObjectFromGoString("123")})
	if got := out.(*object.Object).FieldTable["value"].Fvalue.(int64); got != 123 {
		t.Fatalf("decode decimal got %d", got)
	}
}

func TestLong_CompareTo_Equals_HashCode(t *testing.T) {
	l1 := longValueOf([]interface{}{int64(100)}).(*object.Object)
	l2 := longValueOf([]interface{}{int64(100)}).(*object.Object)
	l3 := longValueOf([]interface{}{int64(200)}).(*object.Object)

	// equals
	if longEquals([]interface{}{l1, l2}) != types.JavaBoolTrue {
		t.Fatal("l1 should equal l2")
	}
	if longEquals([]interface{}{l1, l3}) != types.JavaBoolFalse {
		t.Fatal("l1 should not equal l3")
	}

	// compareTo
	if got := longCompareTo([]interface{}{l1, l2}).(int64); got != 0 {
		t.Fatalf("compareTo same got %d", got)
	}
	if got := longCompareTo([]interface{}{l1, l3}).(int64); got != -1 {
		t.Fatalf("compareTo smaller got %d", got)
	}

	// hashCode
	h1 := longHashCode([]interface{}{l1}).(int64)
	h2 := longHashCodeStatic([]interface{}{int64(100)}).(int64)
	if h1 != h2 {
		t.Fatalf("hashCode mismatch: %d vs %d", h1, h2)
	}
}

func TestLong_Conversions(t *testing.T) {
	val := int64(0x123456789abcdef0)
	l := longValueOf([]interface{}{val}).(*object.Object)
	if got := longByteValue([]interface{}{l}).(int64); int8(got) != int8(val) {
		t.Fatalf("byteValue got %x", got)
	}
	if got := longShortValue([]interface{}{l}).(int64); int16(got) != int16(val) {
		t.Fatalf("shortValue got %x", got)
	}
	if got := longIntValue([]interface{}{l}).(int64); int32(got) != int32(val) {
		t.Fatalf("intValue got %x", got)
	}
	if got := longFloatValue([]interface{}{l}).(float32); got != float32(val) {
		t.Fatalf("floatValue got %v", got)
	}
}

func TestLong_ParseMethods(t *testing.T) {
	// parseLong with radix
	if got := longParseLongRadix([]interface{}{object.StringObjectFromGoString("123"), int64(8)}).(int64); got != 0123 {
		t.Fatalf("parseLongRadix got %d", got)
	}

	// parseUnsignedLong
	if got := longParseUnsignedLong([]interface{}{object.StringObjectFromGoString("18446744073709551615")}).(int64); uint64(got) != 0xffffffffffffffff {
		t.Fatalf("parseUnsignedLong got %x", got)
	}

	// parseLong CharSequence
	cs := object.StringObjectFromGoString("abc123def")
	if got := longParseLongCharSequence([]interface{}{cs, int64(3), int64(6), int64(10)}).(int64); got != 123 {
		t.Fatalf("parseLongCharSequence got %d", got)
	}
}

func TestLong_ToStringRadix(t *testing.T) {
	if got := object.GoStringFromStringObject(longToStringRadix([]interface{}{int64(255), int64(16)}).(*object.Object)); got != "ff" {
		t.Fatalf("toStringRadix got %s", got)
	}
}

func TestLong_AdditionalMethods(t *testing.T) {
	// bitCount
	if got := longBitCount([]interface{}{int64(7)}).(int64); got != 3 {
		t.Fatalf("bitCount(7) got %d", got)
	}

	// compare
	if got := longCompare([]interface{}{int64(10), int64(20)}).(int64); got != -1 {
		t.Fatalf("compare(10,20) got %d", got)
	}
	if got := longCompare([]interface{}{int64(20), int64(10)}).(int64); got != 1 {
		t.Fatalf("compare(20,10) got %d", got)
	}

	// compareUnsigned
	if got := longCompareUnsigned([]interface{}{int64(-1), int64(0)}).(int64); got != 1 {
		t.Fatalf("compareUnsigned(-1,0) got %d", got)
	}

	// divideUnsigned
	if got := longDivideUnsigned([]interface{}{int64(-2), int64(2)}).(int64); uint64(got) != 0x7fffffffffffffff {
		t.Fatalf("divideUnsigned(-2,2) got 0x%x", got)
	}

	// highestOneBit
	if got := longHighestOneBit([]interface{}{int64(0xF0)}).(int64); got != 0x80 {
		t.Fatalf("highestOneBit(0xF0) got %x", got)
	}

	// numberOfLeadingZeros
	if got := longNumberOfLeadingZeros([]interface{}{int64(1)}).(int64); got != 63 {
		t.Fatalf("numberOfLeadingZeros(1) got %d", got)
	}

	// reverse
	if got := longReverse([]interface{}{int64(1)}).(int64); uint64(got) != 0x8000000000000000 {
		t.Fatalf("reverse(1) got %x", got)
	}

	// toUnsignedString
	out := longToUnsignedString([]interface{}{int64(-1)})
	s := object.GoStringFromStringObject(out.(*object.Object))
	if s != "18446744073709551615" {
		t.Fatalf("toUnsignedString(-1) got %q", s)
	}

	// compress and expand
	if got := longCompress([]interface{}{int64(-1), int64(1)}).(int64); got != 1 {
		t.Fatalf("longCompress(-1, 1) got %d", got)
	}
	if got := longExpand([]interface{}{int64(1), int64(1)}).(int64); got != 1 {
		t.Fatalf("longExpand(1, 1) got %d", got)
	}
}

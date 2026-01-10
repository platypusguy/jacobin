package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
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

func TestLong_ToHexString_And_ToString(t *testing.T) {
	// toHexString should not pad per Java behavior
	out := longToHexString([]interface{}{int64(1)})
	sObj := out.(*object.Object)
	if got := object.GoStringFromStringObject(sObj); got != "1" {
		t.Fatalf("toHexString(1) got %q", got)
	}
	out = longToHexString([]interface{}{int64(-1)})
	sObj = out.(*object.Object)
	if got := object.GoStringFromStringObject(sObj); got != "ffffffffffffffff" {
		t.Fatalf("toHexString(-1) got %q", got)
	}

	// toString decimal
	out = longToString([]interface{}{int64(-123)})
	sObj = out.(*object.Object)
	if got := object.GoStringFromStringObject(sObj); got != "-123" {
		t.Fatalf("toString(-123) got %q", got)
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

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helper to extract class name string from an object
func classNameOf(obj *object.Object) string {
	return object.GoStringFromStringPoolIndex(obj.KlassName)
}

func TestShortFloatDoubleValue(t *testing.T) {
	globals.InitStringPool()

	cases := []int64{0, 1, -1, 42, 127, -128, 32767, -32768}
	for _, v := range cases {
		shortObj := object.MakePrimitiveObject("java/lang/Short", types.Short, v)

		// Test floatValue
		resF := shortFloatValue([]interface{}{shortObj})
		f, okF := resF.(float64)
		if !okF {
			t.Fatalf("expected float64 from shortFloatValue, got %T", resF)
		}
		if f != float64(v) {
			t.Fatalf("floatValue mismatch: expected %v, got %v", float64(v), f)
		}

		// Test doubleValue
		resD := shortDoubleValue([]interface{}{shortObj})
		d, okD := resD.(float64)
		if !okD {
			t.Fatalf("expected float64 from shortDoubleValue, got %T", resD)
		}
		if d != float64(v) {
			t.Fatalf("doubleValue mismatch: expected %v, got %v", float64(v), d)
		}
	}
}

func TestShortValueOf(t *testing.T) {
	globals.InitStringPool()

	cases := []int64{0, 1, -1, 12345, -12345, 32767, -32768}
	for _, v := range cases {
		res := shortValueOf([]interface{}{v})
		obj, ok := res.(*object.Object)
		if !ok {
			t.Fatalf("expected *object.Object from shortValueOf, got %T", res)
		}
		if cn := classNameOf(obj); cn != "java/lang/Short" {
			t.Fatalf("expected class java/lang/Short, got %s", cn)
		}
		// Check the boxed value
		val, ok := obj.FieldTable["value"].Fvalue.(int64)
		if !ok {
			t.Fatalf("expected int64 value field, got %T", obj.FieldTable["value"].Fvalue)
		}
		if val != v {
			t.Fatalf("valueOf mismatch: expected %d, got %d", v, val)
		}
	}
}

func TestShortRoundTrip_ValueOfThenDoubleValue(t *testing.T) {
	globals.InitStringPool()

	cases := []int64{7, -7, 30000, -30000}
	for _, v := range cases {
		obj := shortValueOf([]interface{}{v}).(*object.Object)
		res := shortDoubleValue([]interface{}{obj})
		d := res.(float64)
		if d != float64(v) {
			t.Fatalf("round-trip mismatch: expected %v, got %v", float64(v), d)
		}
	}
}

func TestShort_AdditionalMethods(t *testing.T) {
	globals.InitStringPool()

	// compare
	if res := shortCompare([]interface{}{int64(123), int64(-123)}).(int64); res != 246 {
		t.Errorf("compare(123, -123) expected 246, got %d", res)
	}
	if res := shortCompare([]interface{}{int64(-123), int64(123)}).(int64); res != -246 {
		t.Errorf("compare(-123, 123) expected -246, got %d", res)
	}
	if res := shortCompare([]interface{}{int64(10), int64(20)}).(int64); res >= 0 {
		t.Errorf("compare(10, 20) expected < 0, got %d", res)
	}
	if res := shortCompare([]interface{}{int64(20), int64(10)}).(int64); res <= 0 {
		t.Errorf("compare(20, 10) expected > 0, got %d", res)
	}

	// compareUnsigned
	if res := shortCompareUnsigned([]interface{}{int64(-1), int64(1)}).(int64); res != 65534 {
		t.Errorf("compareUnsigned(-1, 1) expected 65534, got %d", res)
	}
	if res := shortCompareUnsigned([]interface{}{int64(123), int64(-256)}).(int64); res != -65157 {
		t.Errorf("compareUnsigned(123, -256) expected -65157, got %d", res)
	}

	// parseShort
	sObj := object.StringObjectFromGoString("1234")
	if res := shortParseShort([]interface{}{sObj}).(int64); res != 1234 {
		t.Errorf("parseShort('1234') expected 1234, got %d", res)
	}

	// parseShortRadix
	sObjHex := object.StringObjectFromGoString("1A")
	if res := shortParseShortRadix([]interface{}{sObjHex, int64(16)}).(int64); res != 26 {
		t.Errorf("parseShortRadix('1A', 16) expected 26, got %d", res)
	}

	// reverseBytes
	if res := shortReverseBytes([]interface{}{int64(0x1234)}).(int64); uint16(res) != 0x3412 {
		t.Errorf("reverseBytes(0x1234) expected 0x3412, got 0x%x", uint16(res))
	}

	// byteValue
	shortObj := object.MakePrimitiveObject("java/lang/Short", types.Short, int64(12345))
	if res := shortByteValue([]interface{}{shortObj}).(int64); res != 57 {
		t.Errorf("byteValue(12345) expected 57, got %d", res)
	}

	// toStringS
	resStr := shortToStringS([]interface{}{int64(123)}).(*object.Object)
	if s := object.GoStringFromStringObject(resStr); s != "123" {
		t.Errorf("toString(123) expected '123', got '%s'", s)
	}

	// toUnsignedInt
	if res := shortToUnsignedInt([]interface{}{int64(-1)}).(int64); res != 65535 {
		t.Errorf("toUnsignedInt(-1) expected 65535, got %d", res)
	}

	// toUnsignedLong
	if res := shortToUnsignedLong([]interface{}{int64(-1)}).(int64); res != 65535 {
		t.Errorf("toUnsignedLong(-1) expected 65535, got %d", res)
	}

	// equals
	s1 := shortValueOf([]interface{}{int64(42)}).(*object.Object)
	s2 := shortValueOf([]interface{}{int64(42)}).(*object.Object)
	s3 := shortValueOf([]interface{}{int64(43)}).(*object.Object)
	if res := shortEquals([]interface{}{s1, s2}); res != types.JavaBoolTrue {
		t.Errorf("equals(42, 42) expected true, got %v", res)
	}
	if res := shortEquals([]interface{}{s1, s3}); res != types.JavaBoolFalse {
		t.Errorf("equals(42, 43) expected false, got %v", res)
	}
}

func TestShortCompareTo(t *testing.T) {
	globals.InitStringPool()

	s1 := shortValueOf([]interface{}{int64(10)}).(*object.Object)
	s2 := shortValueOf([]interface{}{int64(20)}).(*object.Object)
	s3 := shortValueOf([]interface{}{int64(10)}).(*object.Object)

	if res := shortCompareTo([]interface{}{s1, s2}).(int64); res >= 0 {
		t.Errorf("compareTo(10, 20) expected negative, got %d", res)
	}
	if res := shortCompareTo([]interface{}{s2, s1}).(int64); res <= 0 {
		t.Errorf("compareTo(20, 10) expected positive, got %d", res)
	}
	if res := shortCompareTo([]interface{}{s1, s3}).(int64); res != 0 {
		t.Errorf("compareTo(10, 10) expected 0, got %d", res)
	}

	// null receiver
	res := shortCompareTo([]interface{}{object.Null, s1})
	gErr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk for null receiver, got %T", res)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}

	// null argument
	res = shortCompareTo([]interface{}{s1, object.Null})
	gErr, ok = res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk for null argument, got %T", res)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestShortDecode(t *testing.T) {
	globals.InitStringPool()

	// decimal
	res := shortDecode([]interface{}{object.StringObjectFromGoString("123")})
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	if v := obj.FieldTable["value"].Fvalue.(int64); v != 123 {
		t.Errorf("decode('123') expected 123, got %d", v)
	}

	// negative decimal
	res = shortDecode([]interface{}{object.StringObjectFromGoString("-45")})
	obj = res.(*object.Object)
	if v := obj.FieldTable["value"].Fvalue.(int64); v != -45 {
		t.Errorf("decode('-45') expected -45, got %d", v)
	}

	// "#" hex prefix
	res = shortDecode([]interface{}{object.StringObjectFromGoString("#1A")})
	obj, ok = res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object for hex decode, got %T", res)
	}
	if v := obj.FieldTable["value"].Fvalue.(int64); v != 26 {
		t.Errorf("decode('#1A') expected 26, got %d", v)
	}

	// empty string
	res = shortDecode([]interface{}{object.StringObjectFromGoString("")})
	gErr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk for empty string, got %T", res)
	}
	if gErr.ExceptionType != excNames.NumberFormatException {
		t.Errorf("expected NumberFormatException, got %d", gErr.ExceptionType)
	}

	// invalid string
	res = shortDecode([]interface{}{object.StringObjectFromGoString("abc")})
	gErr, ok = res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk for invalid string, got %T", res)
	}
	if gErr.ExceptionType != excNames.NumberFormatException {
		t.Errorf("expected NumberFormatException, got %d", gErr.ExceptionType)
	}
}

func TestShortHashCode(t *testing.T) {
	globals.InitStringPool()

	shortObj := object.MakePrimitiveObject("java/lang/Short", types.Short, int64(12345))
	if res := shortHashCode([]interface{}{shortObj}).(int64); res != 12345 {
		t.Errorf("hashCode() expected 12345, got %d", res)
	}

	if res := shortHashCodeStatic([]interface{}{int64(-999)}).(int64); res != -999 {
		t.Errorf("hashCode(S) expected -999, got %d", res)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"reflect"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestLoad_Lang_Byte_RegistersMethods(t *testing.T) {
	saved := MethodSignatures
	defer func() { MethodSignatures = saved }()
	MethodSignatures = make(map[string]GMeth)

	Load_Lang_Byte()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/lang/Byte.<clinit>()V", 0, clinitGeneric},
		{"java/lang/Byte.byteValue()B", 0, byteIntLongShortByteValue},
		{"java/lang/Byte.compare(BB)I", 2, byteCompare},
		{"java/lang/Byte.compareUnsigned(BB)I", 2, byteCompareUnsigned},
		{"java/lang/Byte.decode(Ljava/lang/String;)Ljava/lang/Byte;", 1, byteDecode},
		{"java/lang/Byte.doubleValue()D", 0, byteFloatDoubleValue},
		{"java/lang/Byte.equals(Ljava/lang/Object;)Z", 1, byteEquals},
		{"java/lang/Byte.floatValue()F", 0, byteFloatDoubleValue},
		{"java/lang/Byte.hashCode()I", 0, byteHashCode},
		{"java/lang/Byte.hashCode(B)I", 1, byteHashCodeStatic},
		{"java/lang/Byte.intValue()I", 0, byteIntLongShortByteValue},
		{"java/lang/Byte.longValue()J", 0, byteIntLongShortByteValue},
		{"java/lang/Byte.parseByte(Ljava/lang/String;)B", 1, byteParseByte},
		{"java/lang/Byte.parseByte(Ljava/lang/String;I)B", 2, byteParseByteRadix},
		{"java/lang/Byte.shortValue()S", 0, byteIntLongShortByteValue},
		{"java/lang/Byte.toString()Ljava/lang/String;", 0, byteToString},
		{"java/lang/Byte.toString(B)Ljava/lang/String;", 1, byteToStringStatic},
		{"java/lang/Byte.toUnsignedInt(B)I", 1, byteToUnsignedInt},
		{"java/lang/Byte.toUnsignedLong(B)J", 1, byteToUnsignedLong},
		{"java/lang/Byte.valueOf(B)Ljava/lang/Byte;", 1, byteValueOf},
		{"java/lang/Byte.valueOf(Ljava/lang/String;)Ljava/lang/Byte;", 1, byteValueOfString},
		{"java/lang/Byte.valueOf(Ljava/lang/String;I)Ljava/lang/Byte;", 2, byteValueOfString},
	}

	for _, c := range checks {
		got, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if got.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
		}
		if got.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}

func TestByteDecode_Various(t *testing.T) {
	globals.InitGlobals("test")

	// valid with leading #
	s := object.StringObjectFromGoString("#0a")
	ret := byteDecode([]interface{}{s})
	obj := ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != 10 {
		t.Fatalf("decode #0a expected 10, got %v", obj.FieldTable["value"].Fvalue)
	}

	// valid with 0x
	s = object.StringObjectFromGoString("0x2f")
	ret = byteDecode([]interface{}{s})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != 47 {
		t.Fatalf("decode 0x2f expected 47, got %v", obj.FieldTable["value"].Fvalue)
	}

	// too large
	s = object.StringObjectFromGoString("1ff") // 511
	if blk, ok := byteDecode([]interface{}{s}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Fatalf("decode too-large expected NFE, got %T", blk)
	}

	// invalid hex
	s = object.StringObjectFromGoString("zz")
	if blk, ok := byteDecode([]interface{}{s}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Fatalf("decode invalid hex expected NFE, got %T", blk)
	}

	// octal
	s = object.StringObjectFromGoString("010")
	ret = byteDecode([]interface{}{s})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != 8 {
		t.Fatalf("decode 010 expected 8, got %v", obj.FieldTable["value"].Fvalue)
	}

	// decimal
	s = object.StringObjectFromGoString("127")
	ret = byteDecode([]interface{}{s})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != 127 {
		t.Fatalf("decode 127 expected 127, got %v", obj.FieldTable["value"].Fvalue)
	}

	// negative
	s = object.StringObjectFromGoString("-128")
	ret = byteDecode([]interface{}{s})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != -128 {
		t.Fatalf("decode -128 expected -128, got %v", obj.FieldTable["value"].Fvalue)
	}
}

func TestByteDoubleValue_ToString_ValueOf(t *testing.T) {
	globals.InitGlobals("test")

	b := object.MakePrimitiveObject("java/lang/Byte", types.Byte, int64(127))

	// doubleValue
	if v := byteFloatDoubleValue([]interface{}{b}); v.(float64) != float64(127) {
		t.Fatalf("doubleValue wrong")
	}

	// toString
	s := byteToString([]interface{}{b}).(*object.Object)
	if str := object.GoStringFromStringObject(s); str != "127" {
		t.Fatalf("toString wrong: %q", str)
	}

	// valueOf
	vobj := byteValueOf([]interface{}{int64(5)}).(*object.Object)
	if v := vobj.FieldTable["value"].Fvalue.(int64); v != 5 {
		t.Fatalf("valueOf 5 wrong: %v", v)
	}
}

func TestByte_AdditionalMethods(t *testing.T) {
	globals.InitGlobals("test")

	// compare
	if res := byteCompare([]interface{}{int64(10), int64(20)}).(int64); res >= 0 {
		t.Errorf("compare(10, 20) expected < 0, got %d", res)
	}
	if res := byteCompare([]interface{}{int64(20), int64(10)}).(int64); res <= 0 {
		t.Errorf("compare(20, 10) expected > 0, got %d", res)
	}

	// compareUnsigned
	if res := byteCompareUnsigned([]interface{}{int64(-1), int64(1)}).(int64); res != 1 {
		t.Errorf("compareUnsigned(-1, 1) expected 1, got %d", res)
	}

	// parseByte
	sObj := object.StringObjectFromGoString("123")
	if res := byteParseByte([]interface{}{sObj}).(int64); res != 123 {
		t.Errorf("parseByte('123') expected 123, got %d", res)
	}

	// parseByteRadix
	sObjHex := object.StringObjectFromGoString("1A")
	if res := byteParseByteRadix([]interface{}{sObjHex, int64(16)}).(int64); res != 26 {
		t.Errorf("parseByteRadix('1A', 16) expected 26, got %d", res)
	}

	// toStringStatic
	resStr := byteToStringStatic([]interface{}{int64(123)}).(*object.Object)
	if s := object.GoStringFromStringObject(resStr); s != "123" {
		t.Errorf("toString(123) expected '123', got '%s'", s)
	}

	// toUnsignedInt
	if res := byteToUnsignedInt([]interface{}{int64(-1)}).(int64); res != 255 {
		t.Errorf("toUnsignedInt(-1) expected 255, got %d", res)
	}

	// toUnsignedLong
	if res := byteToUnsignedLong([]interface{}{int64(-1)}).(int64); res != 255 {
		t.Errorf("toUnsignedLong(-1) expected 255, got %d", res)
	}

	// equals
	b1 := byteValueOf([]interface{}{int64(42)}).(*object.Object)
	b2 := byteValueOf([]interface{}{int64(42)}).(*object.Object)
	b3 := byteValueOf([]interface{}{int64(43)}).(*object.Object)
	if res := byteEquals([]interface{}{b1, b2}); res != types.JavaBoolTrue {
		t.Errorf("equals(42, 42) expected true, got %v", res)
	}
	if res := byteEquals([]interface{}{b1, b3}); res != types.JavaBoolFalse {
		t.Errorf("equals(42, 43) expected false, got %v", res)
	}

	// valueOf(String)
	sObjVal := object.StringObjectFromGoString("42")
	vobj := byteValueOfString([]interface{}{sObjVal}).(*object.Object)
	if v := vobj.FieldTable["value"].Fvalue.(int64); v != 42 {
		t.Errorf("valueOf('42') expected 42, got %v", v)
	}

	// hashCode
	if h := byteHashCode([]interface{}{b1}); h.(int64) != 42 {
		t.Errorf("hashCode(42) expected 42, got %v", h)
	}
	if h := byteHashCodeStatic([]interface{}{int64(42)}); h.(int64) != 42 {
		t.Errorf("hashCodeStatic(42) expected 42, got %v", h)
	}
}

func TestByteIntLongShortByteValue(t *testing.T) {
	globals.InitGlobals("test")
	b := object.MakePrimitiveObject("java/lang/Byte", types.Byte, int64(42))

	// byteValue
	if v := byteIntLongShortByteValue([]interface{}{b}); v.(int64) != 42 {
		t.Errorf("byteValue expected 42, got %v", v)
	}

	// intValue
	if v := byteIntLongShortByteValue([]interface{}{b}); v.(int64) != 42 {
		t.Errorf("intValue expected 42, got %v", v)
	}

	// longValue
	if v := byteIntLongShortByteValue([]interface{}{b}); v.(int64) != 42 {
		t.Errorf("longValue expected 42, got %v", v)
	}

	// shortValue
	if v := byteIntLongShortByteValue([]interface{}{b}); v.(int64) != 42 {
		t.Errorf("shortValue expected 42, got %v", v)
	}
}

func TestByte_CoverageExt(t *testing.T) {
	globals.InitGlobals("test")

	// byteDecode empty string
	sEmpty := object.StringObjectFromGoString("")
	if blk, ok := byteDecode([]interface{}{sEmpty}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("byteDecode('') expected NFE, got %T", blk)
	}
	sPlus := object.StringObjectFromGoString("+42")
	if ret := byteDecode([]interface{}{sPlus}); ret.(*object.Object).FieldTable["value"].Fvalue.(int64) != 42 {
		t.Errorf("byteDecode('+42') expected 42, got %v", ret)
	}
	sOutOfRange := object.StringObjectFromGoString("128")
	if blk, ok := byteDecode([]interface{}{sOutOfRange}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("byteDecode('128') expected NFE, got %T", blk)
	}

	// byteCompareUnsigned branches
	if res := byteCompareUnsigned([]interface{}{int64(1), int64(1)}).(int64); res != 0 {
		t.Errorf("byteCompareUnsigned(1, 1) expected 0, got %d", res)
	}
	if res := byteCompareUnsigned([]interface{}{int64(1), int64(2)}).(int64); res != -1 {
		t.Errorf("byteCompareUnsigned(1, 2) expected -1, got %d", res)
	}

	// byteEquals branches
	b1 := byteValueOf([]interface{}{int64(42)}).(*object.Object)
	if res := byteEquals([]interface{}{b1, object.Null}); res != types.JavaBoolFalse {
		t.Errorf("byteEquals(b1, null) expected false, got %v", res)
	}
	otherClass := "java/lang/Integer"
	otherObj := object.MakeEmptyObjectWithClassName(&otherClass)
	if res := byteEquals([]interface{}{b1, otherObj}); res != types.JavaBoolFalse {
		t.Errorf("byteEquals(b1, Integer) expected false, got %v", res)
	}

	// byteParseByte error
	sInvalid := object.StringObjectFromGoString("abc")
	if blk, ok := byteParseByte([]interface{}{sInvalid}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("byteParseByte('abc') expected NFE, got %T", blk)
	}

	// byteParseByteRadix error & invalid radix
	if blk, ok := byteParseByteRadix([]interface{}{sInvalid, int64(16)}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("byteParseByteRadix('abc', 16) expected NFE, got %T", blk)
	}
	if blk, ok := byteParseByteRadix([]interface{}{sInvalid, int64(1)}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("byteParseByteRadix(..., 1) expected NFE, got %T", blk)
	}

	// byteValueOfString with radix and error
	sVal := object.StringObjectFromGoString("123")
	vobj := byteValueOfString([]interface{}{sVal, int64(10)}).(*object.Object)
	if v := vobj.FieldTable["value"].Fvalue.(int64); v != 123 {
		t.Errorf("byteValueOfString('123', 10) expected 123, got %v", v)
	}
	if blk, ok := byteValueOfString([]interface{}{sInvalid, int64(10)}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("byteValueOfString('abc', 10) expected NFE, got %T", blk)
	}
}

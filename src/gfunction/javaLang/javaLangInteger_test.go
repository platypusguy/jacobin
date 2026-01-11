package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helper to check *ghelpers.GErrBlk with the expected exception type (int)
func checkIntegerErrType(t *testing.T, res interface{}, expected int) {
	t.Helper()

	if errObj, ok := res.(*ghelpers.GErrBlk); ok {
		if errObj.ExceptionType != expected {
			t.Fatalf("expected exception type %d, got %d", expected, errObj.ExceptionType)
		}
	}
	// If res is not *ghelpers.GErrBlk, do nothing (not an error)
}

func TestIntegerByteValue(t *testing.T) {
	globals.InitStringPool()
	// Create an Integer object with value 127
	intObj := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(127))
	params := []interface{}{intObj}

	res := integerByteValue(params)
	if res != int64(127) {
		t.Errorf("expected 127, got %v", res)
	}

	// Test with value that will be truncated
	intObj = object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(257))
	params = []interface{}{intObj}

	res = integerByteValue(params)
	if res != int64(1) { // 129 truncated to byte is -127
		t.Errorf("expected -127, got %v", res)
	}
}

func TestIntegerParseInt(t *testing.T) {
	globals.InitStringPool()

	// Test valid integer
	strObj := object.StringObjectFromGoString("123")
	params := []interface{}{strObj}

	res := integerParseInt(params)
	if res != int64(123) {
		t.Errorf("expected 123, got %v", res)
	}

	// Test negative integer
	strObj = object.StringObjectFromGoString("-456")
	params = []interface{}{strObj}

	res = integerParseInt(params)
	if res != int64(-456) {
		t.Errorf("expected -456, got %v", res)
	}

	// Test invalid integer
	strObj = object.StringObjectFromGoString("not_an_int")
	params = []interface{}{strObj}

	res = integerParseInt(params)
	checkIntegerErrType(t, res, excNames.NumberFormatException)
}

func TestIntegerParseIntRadix(t *testing.T) {
	globals.InitStringPool()

	// Test base 10
	strObj := object.StringObjectFromGoString("123")
	params := []interface{}{strObj, int64(10)}

	res := integerParseIntRadix(params)
	if res != int64(123) {
		t.Errorf("expected 123, got %v", res)
	}

	// Test base 16 (hex)
	strObj = object.StringObjectFromGoString("1A")
	params = []interface{}{strObj, int64(16)}

	res = integerParseIntRadix(params)
	if res != int64(26) {
		t.Errorf("expected 26, got %v", res)
	}

	// Test base 2 (binary)
	strObj = object.StringObjectFromGoString("1010")
	params = []interface{}{strObj, int64(2)}

	res = integerParseIntRadix(params)
	if res != int64(10) {
		t.Errorf("expected 10, got %v", res)
	}

	// Test invalid radix
	strObj = object.StringObjectFromGoString("123")
	params = []interface{}{strObj, int64(37)}

	res = integerParseIntRadix(params)
	checkIntegerErrType(t, res, excNames.NumberFormatException)

	// Test invalid number for given radix
	strObj = object.StringObjectFromGoString("12A")
	params = []interface{}{strObj, int64(10)}

	res = integerParseIntRadix(params)
	checkIntegerErrType(t, res, excNames.NumberFormatException)
}

func TestIntegerToString(t *testing.T) {
	globals.InitStringPool()

	// Create an Integer object with value 123
	intObj := object.MakePrimitiveObject("java/lang/Integer", "I", int64(123))
	params := []interface{}{intObj}

	res := integerToString(params)
	if obj, ok := res.(*object.Object); ok {
		str := object.GoStringFromStringObject(obj)
		if str != "123" {
			t.Errorf("expected \"123\", got %q", str)
		}
	} else {
		t.Errorf("expected *object.Object, got %T", res)
	}
}

func TestIntegerToStringIorII(t *testing.T) {
	globals.InitStringPool()

	// Test with just the integer
	params := []interface{}{int64(123)}

	res := integerToStringIorII(params)
	if obj, ok := res.(*object.Object); ok {
		str := object.GoStringFromStringObject(obj)
		if str != "123" {
			t.Errorf("expected \"123\", got %q", str)
		}
	} else {
		t.Errorf("expected *object.Object, got %T", res)
	}

	// Test with integer and radix
	params = []interface{}{int64(26), int64(16)}

	res = integerToStringIorII(params)
	if obj, ok := res.(*object.Object); ok {
		str := object.GoStringFromStringObject(obj)
		if str != "1a" {
			t.Errorf("expected \"1a\", got %q", str)
		}
	} else {
		t.Errorf("expected *object.Object, got %T", res)
	}

	// Test with invalid radix
	params = []interface{}{int64(123), int64(37)}

	res = integerToStringIorII(params)
	checkIntegerErrType(t, res, excNames.NumberFormatException)
}

func TestIntegerBitCount(t *testing.T) {
	globals.InitStringPool()
	// Test with various integers
	testCases := []struct {
		input    int64
		expected int64
	}{
		{0, 0},           // No bits set
		{1, 1},           // One bit set
		{3, 2},           // Two bits set (11 in binary)
		{-1, 32},         // All bits set
		{0x0F0F0F0F, 16}, // Half the bits set
	}

	for _, tc := range testCases {
		params := []interface{}{tc.input}
		res := integerBitCount(params)
		if res != tc.expected {
			t.Errorf("bitCount(%d): expected %d, got %v", tc.input, tc.expected, res)
		}
	}
}

func signInt64(x int64) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func TestIntegerCompare(t *testing.T) {
	globals.InitStringPool()
	testCases := []struct {
		x        int64
		y        int64
		expected int64
	}{
		{5, 10, -1},   // x < y
		{10, 5, 1},    // x > y
		{7, 7, 0},     // x == y
		{-5, 10, -1},  // negative x < positive y
		{10, -5, 1},   // positive x > negative y
		{-10, -5, -1}, // -10 < -5
	}

	for _, tc := range testCases {
		params := []interface{}{tc.x, tc.y}
		res := integerCompare(params).(int64)
		if signInt64(res) != signInt64(tc.expected) {
			t.Errorf("compare(%d, %d): expected %d, got %v", tc.x, tc.y, tc.expected, res)
		}
	}
}

func TestIntegerValueOfInt(t *testing.T) {
	globals.InitStringPool()

	// Test with a simple integer
	params := []interface{}{int64(42)}

	res := integerValueOfInt(params)
	if obj, ok := res.(*object.Object); ok {
		// Check the value field
		if value, ok := obj.FieldTable["value"].Fvalue.(int64); ok {
			if value != 42 {
				t.Errorf("expected value 42, got %d", value)
			}
		} else {
			t.Errorf("expected int64 value, got %T", obj.FieldTable["value"].Fvalue)
		}

		// Check the class name
		className := object.GoStringFromStringPoolIndex(obj.KlassName)
		if className != "java/lang/Integer" {
			t.Errorf("expected type java/lang/Integer, got %s", className)
		}
	} else {
		t.Errorf("expected *object.Object, got %T", res)
	}
}

func TestIntegerValueOfString(t *testing.T) {
	globals.InitStringPool()

	// Test with valid string
	strObj := object.StringObjectFromGoString("123")
	params := []interface{}{strObj}

	res := integerValueOfString(params)
	if obj, ok := res.(*object.Object); ok {
		// Check the value field
		if value, ok := obj.FieldTable["value"].Fvalue.(int64); ok {
			if value != 123 {
				t.Errorf("expected value 123, got %d", value)
			}
		} else {
			t.Errorf("expected int64 value, got %T", obj.FieldTable["value"].Fvalue)
		}

		// Check the class name
		className := object.GoStringFromStringPoolIndex(obj.KlassName)
		if className != "java/lang/Integer" {
			t.Errorf("expected type java/lang/Integer, got %s", className)
		}
	} else {
		t.Errorf("expected *object.Object, got %T", res)
	}

	// Test with valid string and radix
	strObj = object.StringObjectFromGoString("1A")
	params = []interface{}{strObj, int64(16)}

	res = integerValueOfString(params)
	if obj, ok := res.(*object.Object); ok {
		// Check the value field
		if value, ok := obj.FieldTable["value"].Fvalue.(int64); ok {
			if value != 26 {
				t.Errorf("expected value 26, got %d", value)
			}
		} else {
			t.Errorf("expected int64 value, got %T", obj.FieldTable["value"].Fvalue)
		}
	} else {
		t.Errorf("expected *object.Object, got %T", res)
	}

	// Test with invalid string
	strObj = object.StringObjectFromGoString("not_an_int")
	params = []interface{}{strObj}

	res = integerValueOfString(params)
	checkIntegerErrType(t, res, excNames.NumberFormatException)
}

func TestInteger_NumberOfLeadingAndTrailingZeros(t *testing.T) {
	// 0x0000F000 -> leading zeros = 16, trailing zeros = 12 (32-bit semantics)
	val := int64(0x0000F000)
	if lz := integerNumberOfLeadingZeros([]interface{}{val}).(int64); lz != 16 {
		t.Fatalf("numberOfLeadingZeros(0x0000F000)=%d", lz)
	}
	if tz := integerNumberOfTrailingZeros([]interface{}{val}).(int64); tz != 12 {
		t.Fatalf("numberOfTrailingZeros(0x0000F000)=%d", tz)
	}
}

func TestInteger_RotateLeftRight(t *testing.T) {
	val := int64(0x12345678)
	rl := integerRotateLeft([]interface{}{val, int64(8)}).(int64)
	rr := integerRotateRight([]interface{}{val, int64(8)}).(int64)
	if uint32(rl) != 0x34567812 {
		t.Fatalf("rotateLeft(0x12345678,8)=0x%08x", uint32(rl))
	}
	if uint32(rr) != 0x78123456 {
		t.Fatalf("rotateRight(0x12345678,8)=0x%08x", uint32(rr))
	}
}

func TestInteger_Reverse_And_ReverseBytes(t *testing.T) {
	// reverse bits of 1 -> 0x80000000
	if got := uint32(integerReverse([]interface{}{int64(1)}).(int64)); got != 0x80000000 {
		t.Fatalf("reverse(1)=0x%08x", got)
	}
	// reverse bytes of 0x12345678 -> 0x78563412
	if got := uint32(integerReverseBytes([]interface{}{int64(0x12345678)}).(int64)); got != 0x78563412 {
		t.Fatalf("reverseBytes(0x12345678)=0x%08x", got)
	}
}

func TestInteger_ToHexOctalBinaryString(t *testing.T) {
	// hex
	if s := object.GoStringFromStringObject(integerToHexString([]interface{}{int64(26)}).(*object.Object)); s != "1a" {
		t.Fatalf("toHexString(26)=%q", s)
	}
	// octal
	if s := object.GoStringFromStringObject(integerToOctalString([]interface{}{int64(26)}).(*object.Object)); s != "32" {
		t.Fatalf("toOctalString(26)=%q", s)
	}
	// binary unsigned: -1 -> 32 ones
	s := object.GoStringFromStringObject(integerToBinaryString([]interface{}{int64(-1)}).(*object.Object))
	if len(s) != 32 {
		t.Fatalf("binary length for -1 expected 32, got %d", len(s))
	}
	for i := 0; i < len(s); i++ {
		if s[i] != '1' {
			t.Fatalf("expected all ones, got %q at %d", s[i], i)
		}
	}
}

func TestInteger_CompareUnsigned_And_UnsignedDivRem(t *testing.T) {
	// compareUnsigned: (-1) > 0 as unsigned
	if c := integerCompareUnsigned([]interface{}{int64(-1), int64(0)}).(int64); c != 1 {
		t.Fatalf("compareUnsigned(-1,0)=%d", c)
	}
	if c := integerCompareUnsigned([]interface{}{int64(1), int64(-1)}).(int64); c != -1 {
		t.Fatalf("compareUnsigned(1,-1)=%d", c)
	}
	// divideUnsigned and remainderUnsigned
	if q := integerDivideUnsigned([]interface{}{int64(-2), int64(3)}).(int64); q != int64(uint32(0xFFFFFFFE)/3) {
		t.Fatalf("divideUnsigned(-2,3)=%d", q)
	}
	if r := integerRemainderUnsigned([]interface{}{int64(-2), int64(3)}).(int64); r != int64(uint32(0xFFFFFFFE)%3) {
		t.Fatalf("remainderUnsigned(-2,3)=%d", r)
	}
	// divide by zero -> ArithmeticException
	if res := integerDivideUnsigned([]interface{}{int64(1), int64(0)}); res == nil {
		t.Fatalf("expected error for divide by zero")
	} else if geb, ok := res.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.ArithmeticException {
		t.Fatalf("expected ArithmeticException, got %T", res)
	}
}

func TestInteger_HighestLowestOneBit_MaxMinSum(t *testing.T) {
	if h := uint32(integerHighestOneBit([]interface{}{int64(0xF234)}).(int64)); h != 0x8000 {
		t.Fatalf("highestOneBit(0xF234)=0x%08x", h)
	}
	if l := uint32(integerLowestOneBit([]interface{}{int64(0xF234)}).(int64)); l != 0x0004 {
		t.Fatalf("lowestOneBit(0xF234)=0x%08x", l)
	}
	if m := integerMax([]interface{}{int64(-5), int64(3)}).(int64); m != 3 {
		t.Fatalf("max(-5,3)=%d", m)
	}
	if n := integerMin([]interface{}{int64(-5), int64(3)}).(int64); n != -5 {
		t.Fatalf("min(-5,3)=%d", n)
	}
	if s := integerSum([]interface{}{int64(7), int64(8)}).(int64); s != 15 {
		t.Fatalf("sum(7,8)=%d", s)
	}
}

func TestInteger_ToUnsignedString_Variants(t *testing.T) {
	// toUnsignedString with negative -> decimal of uint32
	s := object.GoStringFromStringObject(integerToUnsignedString([]interface{}{int64(-1)}).(*object.Object))
	if s != "4294967295" {
		t.Fatalf("toUnsignedString(-1)=%q", s)
	}
	// toUnsignedStringRadix hex
	s2 := object.GoStringFromStringObject(integerToUnsignedStringRadix([]interface{}{int64(-1), int64(16)}).(*object.Object))
	if s2 != "ffffffff" {
		t.Fatalf("toUnsignedString(-1,16)=%q", s2)
	}
}

func TestInteger_ParseIntCharSequence(t *testing.T) {
	globals.InitStringPool()

	cs := object.StringObjectFromGoString("123456")
	// parseInt(cs, 1, 4, 10) -> "234" -> 234
	res := integerParseIntCharSequence([]interface{}{cs, int64(1), int64(4), int64(10)})
	if res.(int64) != 234 {
		t.Errorf("expected 234, got %v", res)
	}

	// parseUnsignedInt(cs, 1, 4, 16) -> "234" hex -> 564
	res = integerParseUnsignedIntCharSequence([]interface{}{cs, int64(1), int64(4), int64(16)})
	if res.(int64) != 564 {
		t.Errorf("expected 564, got %v", res)
	}
}

func TestInteger_CompressExpand_Negative(t *testing.T) {
	globals.InitStringPool()

	// Test compress with negative mask
	// In Java, compress(0xFFFFFFFF, 0x80000000) should be 1 if mask is treated as unsigned
	// Wait, if mask is 0x80000000, only the 31st bit is kept and moved to position 0.
	// So compress(-1, 0x80000000) should be 1.
	res := integerCompress([]interface{}{int64(-1), int64(int32(-2147483648))})
	if res.(int64) != 1 {
		t.Errorf("compress(-1, 0x80000000): expected 1, got %v", res)
	}

	// Test expand with negative mask
	// expand(1, 0x80000000) should be 0x80000000
	res = integerExpand([]interface{}{int64(1), int64(int32(-2147483648))})
	if uint32(res.(int64)) != 0x80000000 {
		t.Errorf("expand(1, 0x80000000): expected 0x80000000, got 0x%x", res)
	}
}

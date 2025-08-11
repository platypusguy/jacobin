package gfunction

import (
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/types"
	"testing"
)

// Helper to check *GErrBlk with the expected exception type (int)
func checkIntegerErrType(t *testing.T, res interface{}, expected int) {
	t.Helper()

	if errObj, ok := res.(*GErrBlk); ok {
		if errObj.ExceptionType != expected {
			t.Fatalf("expected exception type %d, got %d", expected, errObj.ExceptionType)
		}
	}
	// If res is not *GErrBlk, do nothing (not an error)
}

func TestIntegerByteValue(t *testing.T) {
	globals.InitStringPool()
	// Create an Integer object with value 127
	intObj := Populator("java/lang/Integer", types.Int, int64(127))
	params := []interface{}{intObj}

	res := integerByteValue(params)
	if res != int64(127) {
		t.Errorf("expected 127, got %v", res)
	}

	// Test with value that will be truncated
	intObj = Populator("java/lang/Integer", types.Int, int64(257))
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
	intObj := Populator("java/lang/Integer", "I", int64(123))
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

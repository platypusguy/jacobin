package gfunction

import (
	"jacobin/object"
	"math/big"
	"strconv"
	"testing"
)

// Helper to create BigDecimal from string for tests
func makeBigDecimalFromString(t *testing.T, s string) *object.Object {
	t.Helper()
	bd := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	ret := bigdecimalInitString([]interface{}{bd, object.StringObjectFromGoString(s)})
	if ret != nil {
		t.Fatalf("makeBigDecimalFromString: failed to init BigDecimal from string %s", s)
	}
	return bd
}

// Helper to compare BigDecimal unscaled value string and scale
func assertBigDecimalUnscaledScale(t *testing.T, bd *object.Object, expectedUnscaled string, expectedScale int64) {
	t.Helper()
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	bigInt := intValObj.FieldTable["value"].Fvalue.(*big.Int)
	if bigInt.String() != expectedUnscaled {
		t.Fatalf("expected unscaled %s, got %s", expectedUnscaled, bigInt.String())
	}
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	if scale != expectedScale {
		t.Fatalf("expected scale %d, got %d", expectedScale, scale)
	}
}

func buildStripCvt(t *testing.T, startString string, scale1, scale2 int64) {
	t.Logf("buildStripCvt: start string: %s", startString)
	startDouble, err := strconv.ParseFloat(startString, 64)
	if err != nil {
		t.Errorf("ERROR buildStripCvt: failed to parse start string %q: %v", startString, err)
	}
	// Make BigDecimal from a Go string.
	original := makeBigDecimalFromString(t, startString)
	// Convert to a float64 and compare to the expected value.
	dbl := bigdecimalDoubleValue([]interface{}{original})
	if dbl != startDouble {
		t.Errorf("ERROR buildStripCvt: original expected value: %f, observed: %f", startDouble, dbl)
	}
	// Get its scale.
	scale := original.FieldTable["scale"].Fvalue.(int64)
	// Compare to the expected scale.
	if scale != scale1 {
		t.Errorf("ERROR buildStripCvt: original expected scale: %d, observed: %d", scale1, scale)
	}
	// Strip trailing zeros.
	stripped := bigdecimalStripTrailingZeros([]interface{}{original}).(*object.Object)
	// Convert to a float64.
	dbl = bigdecimalDoubleValue([]interface{}{stripped})
	if dbl != startDouble {
		t.Errorf("ERROR buildStripCvt: after stripping zeros, expected value: %f, observed: %f", startDouble, dbl)
	}
	// Get its scale.
	scale = stripped.FieldTable["scale"].Fvalue.(int64)
	// Compare to the expected scale.
	if scale != scale2 {
		t.Errorf("ERROR buildStripCvt: after stripping zeros, expected scale: %d, observed: %d", scale2, scale)
	}
}

func TestBigDecimalNtoZFunctions(t *testing.T) {

	t.Run("bigdecimalStripTrailingZeros", func(t *testing.T) {
		bdutInit()
		buildStripCvt(t, "123.45", 2, 2)
		buildStripCvt(t, "3.1416", 4, 4)
		buildStripCvt(t, "-3.141600", 6, 4)
		buildStripCvt(t, "-3.141600E0", 6, 4)
		buildStripCvt(t, "-3.141600e+0", 6, 4)
		buildStripCvt(t, "-3.141600e-0", 6, 4)
		buildStripCvt(t, "-31.41600e-1", 6, 4)
		buildStripCvt(t, "-.3141600E+001", 6, 4)
		buildStripCvt(t, "-.31416E+001", 4, 4)
	})

	t.Run("bigdecimalNegate", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "12345")
		res := bigdecimalNegate([]interface{}{bd}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "-12345", 0)
	})

	t.Run("bigdecimalPlus", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123")
		res := bigdecimalPlus([]interface{}{bd}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "123", 0)
	})

	t.Run("bigdecimalPow", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "2")
		res := bigdecimalPow([]interface{}{bd, int64(10)}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "1024", 0)
	})

	t.Run("bigdecimalPrecision", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "12345")
		got := bigdecimalPrecision([]interface{}{bd}).(int64)
		if got != 5 {
			t.Fatalf("expected precision 5, got %d", got)
		}
	})

	t.Run("bigdecimalRemainder", func(t *testing.T) {
		bdutInit()
		dividend := makeBigDecimalFromString(t, "10")
		divisor := makeBigDecimalFromString(t, "3")
		res := bigdecimalRemainder([]interface{}{dividend, divisor}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "1", 0)

		divZero := makeBigDecimalFromString(t, "0")
		err := bigdecimalRemainder([]interface{}{dividend, divZero})
		if err == nil {
			t.Fatalf("expected error for divide by zero remainder")
		}
	})

	t.Run("bigdecimalScale", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123.45")
		scale := bigdecimalScale([]interface{}{bd}).(int64)
		if scale != 2 {
			t.Fatalf("expected scale 2, got %d", scale)
		}
	})

	t.Run("bigdecimalScaleByPowerOfTen", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "12345")
		res := bigdecimalScaleByPowerOfTen([]interface{}{bd, int64(2)}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "12345", -2)
	})

	t.Run("bigdecimalSetScale", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "12345")
		res := bigdecimalSetScale([]interface{}{bd, int64(2)}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "1234500", 2)
	})

	t.Run("bigdecimalShortValueExact", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "32767")
		res := bigdecimalShortValueExact([]interface{}{bd}).(int64)
		if res != 32767 {
			t.Fatalf("expected short 32767, got %d", res)
		}

		bigBd := makeBigDecimalFromString(t, "40000")
		err := bigdecimalShortValueExact([]interface{}{bigBd})
		if err == nil {
			t.Fatalf("expected error for shortValueExact overflow")
		}
	})

	t.Run("bigdecimalSignum", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "-123")
		res := bigdecimalSignum([]interface{}{bd}).(int64)
		if res != -1 {
			t.Fatalf("expected signum -1, got %d", res)
		}
	})

	t.Run("bigdecimalSubtract", func(t *testing.T) {
		bdutInit()
		a := makeBigDecimalFromString(t, "1000")
		b := makeBigDecimalFromString(t, "100")
		res := bigdecimalSubtract([]interface{}{a, b}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "900", 0)
	})

	t.Run("bigdecimalToBigInteger", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123")
		res := bigdecimalToBigInteger([]interface{}{bd}).(*object.Object)
		intVal := res.FieldTable["value"].Fvalue.(*big.Int)
		if intVal.String() != "123" {
			t.Fatalf("expected BigInteger 123, got %s", intVal.String())
		}
	})

	t.Run("bigdecimalToBigIntegerExact", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123")
		res := bigdecimalToBigIntegerExact([]interface{}{bd}).(*object.Object)
		intVal := res.FieldTable["value"].Fvalue.(*big.Int)
		if intVal.String() != "123" {
			t.Fatalf("expected BigIntegerExact 123, got %s", intVal.String())
		}

		bdFrac := makeBigDecimalFromString(t, "123.45")
		err := bigdecimalToBigIntegerExact([]interface{}{bdFrac})
		if err == nil {
			t.Fatalf("expected error for toBigIntegerExact fractional part")
		}
	})

	t.Run("bigdecimalToString", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123.45")
		strObj := bigdecimalToString([]interface{}{bd}).(*object.Object)
		str := object.GoStringFromStringObject(strObj)
		if str != "123.45" {
			t.Fatalf("expected string '123.45', got %q", str)
		}
	})

	t.Run("bigdecimalUnscaledValue", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123.45")
		res := bigdecimalUnscaledValue([]interface{}{bd}).(*object.Object)
		intVal := res.FieldTable["value"].Fvalue.(*big.Int)
		if intVal.String() != "12345" {
			t.Fatalf("expected unscaled value 12345, got %s", intVal.String())
		}
	})

	t.Run("bigdecimalValueOfDouble", func(t *testing.T) {
		bdutInit()
		val := 12.34
		res := bigdecimalValueOfDouble([]interface{}{val}).(*object.Object)
		scale := res.FieldTable["scale"].Fvalue.(int64)
		if scale < 0 {
			t.Fatalf("unexpected negative scale %d", scale)
		}
	})

	t.Run("bigdecimalValueOfLong", func(t *testing.T) {
		bdutInit()
		val := int64(123456)
		res := bigdecimalValueOfLong([]interface{}{val}).(*object.Object)
		intVal := res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt := intVal.FieldTable["value"].Fvalue.(*big.Int)
		if bigInt.Int64() != val {
			t.Fatalf("expected BigDecimal valueOfLong %d, got %d", val, bigInt.Int64())
		}
	})

	t.Run("bigdecimalValueOfLongInt", func(t *testing.T) {
		bdutInit()
		val := int64(123456)
		scale := int64(2)
		res := bigdecimalValueOfLongInt([]interface{}{val, scale}).(*object.Object)
		intVal := res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt := intVal.FieldTable["value"].Fvalue.(*big.Int)
		if bigInt.Int64() != val {
			t.Fatalf("expected BigDecimal valueOfLongInt %d, got %d", val, bigInt.Int64())
		}
		s := res.FieldTable["scale"].Fvalue.(int64)
		if s != scale {
			t.Fatalf("expected scale %d, got %d", scale, s)
		}
	})
}

func TestBigDecimalNtoZFunctions_EdgeCases(t *testing.T) {
	// Pow edge cases
	t.Run("bigdecimalPow_zeroExponent", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123")
		res := bigdecimalPow([]interface{}{bd, int64(0)}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "1", 0) // x^0 == 1
	})

	t.Run("bigdecimalPow_negativeExponent", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "2")
		// Your implementation might not support negative exponent - expect error or specific behavior
		res := bigdecimalPow([]interface{}{bd, int64(-2)})
		if res == nil {
			t.Fatalf("expected error or specific behavior for negative exponent")
		}
	})

	// bigdecimalSetScale edge cases
	t.Run("bigdecimalSetScale_negativeScale", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "12345")
		res := bigdecimalSetScale([]interface{}{bd, int64(-2)}).(*object.Object)
		assertBigDecimalUnscaledScale(t, res, "123", -2) // scale reduced, value adjusted
	})

	t.Run("bigdecimalSetScale_sameScale", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "12345")
		scale := bd.FieldTable["scale"].Fvalue.(int64)
		res := bigdecimalSetScale([]interface{}{bd, scale})
		if res != bd {
			t.Fatalf("expected same object when scale unchanged")
		}
	})

	// bigdecimalShortValueExact boundary values
	t.Run("bigdecimalShortValueExact_minMax", func(t *testing.T) {
		bdutInit()
		minShort := makeBigDecimalFromString(t, "-32768")
		res := bigdecimalShortValueExact([]interface{}{minShort}).(int64)
		if res != -32768 {
			t.Fatalf("expected short min -32768, got %d", res)
		}

		maxShort := makeBigDecimalFromString(t, "32767")
		res2 := bigdecimalShortValueExact([]interface{}{maxShort}).(int64)
		if res2 != 32767 {
			t.Fatalf("expected short max 32767, got %d", res2)
		}
	})

	// bigdecimalRemainder with negatives
	t.Run("bigdecimalRemainder_negativeValues", func(t *testing.T) {
		bdutInit()

		// dividend negative, divisor positive
		dividend := makeBigDecimalFromString(t, "-10")
		divisor := makeBigDecimalFromString(t, "3")
		res := bigdecimalRemainder([]interface{}{dividend, divisor}).(*object.Object)

		intValObj := res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt := intValObj.FieldTable["value"].Fvalue.(*big.Int)

		expected := "-1"
		if bigInt.String() != expected {
			t.Fatalf("expected remainder %s, got %s", expected, bigInt.String())
		}

		// dividend positive, divisor negative
		dividend = makeBigDecimalFromString(t, "10")
		divisor = makeBigDecimalFromString(t, "-3")
		res = bigdecimalRemainder([]interface{}{dividend, divisor}).(*object.Object)

		intValObj = res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt = intValObj.FieldTable["value"].Fvalue.(*big.Int)

		expected = "1"
		if bigInt.String() != expected {
			t.Fatalf("expected remainder %s, got %s", expected, bigInt.String())
		}

		// dividend negative, divisor negative
		dividend = makeBigDecimalFromString(t, "-10")
		divisor = makeBigDecimalFromString(t, "-3")
		res = bigdecimalRemainder([]interface{}{dividend, divisor}).(*object.Object)

		intValObj = res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt = intValObj.FieldTable["value"].Fvalue.(*big.Int)

		expected = "-1"
		if bigInt.String() != expected {
			t.Fatalf("expected remainder %s, got %s", expected, bigInt.String())
		}
	})

	// bigdecimalStripTrailingZeros edge cases
	t.Run("bigdecimalStripTrailingZeros_zeroValue", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "0.000")
		res := bigdecimalStripTrailingZeros([]interface{}{bd}).(*object.Object)
		scale := res.FieldTable["scale"].Fvalue.(int64)
		if scale != 0 {
			t.Fatalf("expected scale 0 for zero value after stripTrailingZeros, got %d", scale)
		}
	})

	t.Run("bigdecimalStripTrailingZeros_noTrailingZeros", func(t *testing.T) {
		bdutInit()
		bd := makeBigDecimalFromString(t, "123.45")
		res := bigdecimalStripTrailingZeros([]interface{}{bd}).(*object.Object)
		if res == bd {
			t.Fatalf("expected new object when no trailing zeros")
		}
	})

	// bigdecimalValueOfDouble edge cases
	t.Run("bigdecimalValueOfDouble_zeroAndNegative", func(t *testing.T) {
		bdutInit()
		res := bigdecimalValueOfDouble([]interface{}{float64(0)}).(*object.Object)
		scale := res.FieldTable["scale"].Fvalue.(int64)
		if scale < 0 {
			t.Fatalf("unexpected negative scale for zero value")
		}

		resNeg := bigdecimalValueOfDouble([]interface{}{-123.456}).(*object.Object)
		intValObj := resNeg.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt := intValObj.FieldTable["value"].Fvalue.(*big.Int)
		if bigInt.Sign() >= 0 {
			t.Fatalf("expected negative BigDecimal for negative input")
		}
	})

	t.Run("bigdecimalValueOfDouble_largeValue", func(t *testing.T) {
		bdutInit()
		largeVal := 1e18
		res := bigdecimalValueOfDouble([]interface{}{largeVal}).(*object.Object)
		intValObj := res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt := intValObj.FieldTable["value"].Fvalue.(*big.Int)
		if bigInt.Cmp(big.NewInt(0)) <= 0 {
			t.Fatalf("expected positive BigDecimal for large positive input")
		}
	})

	// bigdecimalValueOfLongInt edge cases
	t.Run("bigdecimalValueOfLongInt_zeroScale", func(t *testing.T) {
		bdutInit()
		val := int64(12345)
		res := bigdecimalValueOfLongInt([]interface{}{val, int64(0)}).(*object.Object)
		scale := res.FieldTable["scale"].Fvalue.(int64)
		if scale != 0 {
			t.Fatalf("expected scale 0, got %d", scale)
		}
	})

	t.Run("bigdecimalValueOfLongInt_negativeValue", func(t *testing.T) {
		bdutInit()
		val := int64(-12345)
		scale := int64(3)
		res := bigdecimalValueOfLongInt([]interface{}{val, scale}).(*object.Object)
		intValObj := res.FieldTable["intVal"].Fvalue.(*object.Object)
		bigInt := intValObj.FieldTable["value"].Fvalue.(*big.Int)
		if bigInt.Int64() != val {
			t.Fatalf("expected BigDecimal valueOfLongInt %d, got %d", val, bigInt.Int64())
		}
		if res.FieldTable["scale"].Fvalue.(int64) != scale {
			t.Fatalf("expected scale %d, got %d", scale, res.FieldTable["scale"].Fvalue.(int64))
		}
	})
}

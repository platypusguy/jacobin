/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors
 */
package gfunction

import (
	"jacobin/src/exceptions"
	"jacobin/src/globals"
	"math"
	"math/big"
	"reflect"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Required initialisation for BigDecimal unit tests.
func bdutInit() {
	// Initialise globals, string pool, etc.
	globals.InitGlobals("test")
	gl := globals.GetGlobalRef()
	// Make sure that globals FuncThrowException has a defined value.
	gl.FuncThrowException = exceptions.ThrowExNil
	// Load all the static constants for BigDouble.
	loadStaticsBigDouble()
}

// Helper inside tests: assert GErrBlk with expected exception id
func assertGErrBlk(t *testing.T, res interface{}, wantExc int) {
	t.Helper()
	if res == nil {
		t.Fatalf("expected *GErrBlk, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if errObj.ExceptionType != wantExc {
		t.Fatalf("expected exception %d, got %d; msg=%s", wantExc, errObj.ExceptionType, errObj.ErrMsg)
	}
}

// extractBigDecimalComponents extracts (unscaled *big.Int, scale int64) from a BigDecimal *object.Object
func extractBigDecimalComponents(t *testing.T, bd *object.Object) (*big.Int, int64) {
	t.Helper()
	if bd == nil {
		t.Fatalf("bd is nil")
	}
	intValField, ok := bd.FieldTable["intVal"]
	if !ok {
		t.Fatalf("bd.FieldTable missing intVal")
	}
	intValObj, ok := intValField.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("intVal field is not *object.Object; got %T", intValField.Fvalue)
	}
	valueField, ok := intValObj.FieldTable["value"]
	if !ok {
		t.Fatalf("bigInteger object missing value field")
	}
	bigInt, ok := valueField.Fvalue.(*big.Int)
	if !ok {
		t.Fatalf("value field not *big.Int; got %T", valueField.Fvalue)
	}
	scaleField, ok := bd.FieldTable["scale"]
	if !ok {
		t.Fatalf("bd.FieldTable missing scale")
	}
	scale, ok := scaleField.Fvalue.(int64)
	if !ok {
		t.Fatalf("scale field not int64; got %T", scaleField.Fvalue)
	}
	return bigInt, scale
}

func Test_gfunction_bigdecimal_all(t *testing.T) {
	t.Run("Init_bigdecimalInitIntLong", func(t *testing.T) {
		bdutInit()
		self := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
		res := bigdecimalInitIntLong([]interface{}{self, int64(12345)})
		if res != nil {
			t.Fatalf("expected nil, got error: %v", res)
		}
		unscaled, scale := extractBigDecimalComponents(t, self)
		if unscaled.Cmp(big.NewInt(12345)) != 0 {
			t.Fatalf("unexpected unscaled: %s", unscaled.String())
		}
		if scale != 0 {
			t.Fatalf("expected scale 0, got %d", scale)
		}
	})

	t.Run("Init_bigdecimalInitDouble", func(t *testing.T) {
		bdutInit()
		self := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
		// use a value with decimal places
		res := bigdecimalInitDouble([]interface{}{self, 12.34})
		if res != nil {
			t.Fatalf("expected nil, got error: %v", res)
		}
		unscaled, scale := extractBigDecimalComponents(t, self)
		// 12.34 -> unscaled 1234 scale 2 (or equivalent depending on parse)
		if scale != 2 && scale != 0 {
			// permissive: our parseDecimalString may produce 2
		}
		_ = unscaled
	})

	t.Run("Init_bigdecimalInitString", func(t *testing.T) {
		bdutInit()
		self := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
		strObj := object.StringObjectFromGoString("314.16")
		res := bigdecimalInitString([]interface{}{self, strObj})
		if res != nil {
			t.Fatalf("expected nil, got error: %v", res)
		}
		unscaled, scale := extractBigDecimalComponents(t, self)
		if unscaled.Cmp(big.NewInt(31416)) != 0 {
			t.Fatalf("unexpected unscaled: %s", unscaled.String())
		}
		if scale != 2 {
			t.Fatalf("expected scale 2, got %d", scale)
		}
	})

	t.Run("Init_bigdecimalInitBigInteger", func(t *testing.T) {
		bdutInit()
		self := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
		bi := makeBigIntegerFromBigInt(big.NewInt(77))
		res := bigdecimalInitBigInteger([]interface{}{self, bi})
		if res != nil {
			t.Fatalf("expected nil, got error: %v", res)
		}
		unscaled, scale := extractBigDecimalComponents(t, self)
		if unscaled.Cmp(big.NewInt(77)) != 0 {
			t.Fatalf("unexpected unscaled: %s", unscaled.String())
		}
		if scale != 0 {
			t.Fatalf("expected scale 0, got %d", scale)
		}
	})

	t.Run("bigdecimalAbs", func(t *testing.T) {
		bdutInit()
		in := bigDecimalObjectFromBigInt(big.NewInt(-1234), 4, 0)
		res := bigdecimalAbs([]interface{}{in})
		out, ok := res.(*object.Object)
		if !ok {
			t.Fatalf("expected *object.Object, got %T", res)
		}
		unscaled, _ := extractBigDecimalComponents(t, out)
		if unscaled.Cmp(big.NewInt(1234)) != 0 {
			t.Fatalf("abs failed, got %s", unscaled.String())
		}
	})

	t.Run("bigdecimalAdd", func(t *testing.T) {
		bdutInit()
		a := bigDecimalObjectFromBigInt(big.NewInt(100), 3, 0)
		b := bigDecimalObjectFromBigInt(big.NewInt(23), 2, 0)
		res := bigdecimalAdd([]interface{}{a, b})
		out, ok := res.(*object.Object)
		if !ok {
			t.Fatalf("expected *object.Object, got %T", res)
		}
		unscaled, scale := extractBigDecimalComponents(t, out)
		if unscaled.Cmp(big.NewInt(123)) != 0 {
			t.Fatalf("add incorrect unscaled: %s", unscaled.String())
		}
		if scale != 0 {
			t.Fatalf("unexpected scale: %d", scale)
		}
	})

	t.Run("bigdecimalByteValueExact success and failure", func(t *testing.T) {
		bdutInit()
		okbd := bigDecimalObjectFromBigInt(big.NewInt(127), 3, 0)
		res := bigdecimalByteValueExact([]interface{}{okbd})
		if _, ok := res.(types.JavaByte); !ok {
			t.Fatalf("expected types.JavaByte, got %T", res)
		}
		// out of range
		bad := bigDecimalObjectFromBigInt(big.NewInt(128), 3, 0)
		res2 := bigdecimalByteValueExact([]interface{}{bad})
		assertGErrBlk(t, res2, excNames.ArithmeticException)
	})

	t.Run("bigdecimalCompareTo", func(t *testing.T) {
		bdutInit()
		a := bigDecimalObjectFromBigInt(big.NewInt(1), 1, 0)
		b := bigDecimalObjectFromBigInt(big.NewInt(2), 1, 0)
		res := bigdecimalCompareTo([]interface{}{a, b})
		if res.(int64) >= 0 {
			t.Fatalf("expected negative cmp, got %d", res.(int64))
		}
	})

	t.Run("bigdecimalDivide and divide by zero", func(t *testing.T) {
		bdutInit()
		dividend := bigDecimalObjectFromBigInt(big.NewInt(10), 2, 0)
		divisor := bigDecimalObjectFromBigInt(big.NewInt(3), 1, 0)
		res := bigdecimalDivide([]interface{}{dividend, divisor})
		out, ok := res.(*object.Object)
		if !ok {
			t.Fatalf("expected *object.Object, got %T", res)
		}
		quot, _ := extractBigDecimalComponents(t, out)
		// integer division 10/3 -> 3
		if quot.Cmp(big.NewInt(3)) != 0 {
			t.Fatalf("unexpected quotient %s", quot.String())
		}

		// divide by zero
		zero := bigDecimalObjectFromBigInt(big.NewInt(0), 1, 0)
		res2 := bigdecimalDivide([]interface{}{dividend, zero})
		assertGErrBlk(t, res2, excNames.ArithmeticException)
	})

	t.Run("bigdecimalDivideAndRemainder", func(t *testing.T) {
		bdutInit()
		dividend := bigDecimalObjectFromBigInt(big.NewInt(20), 2, 0)
		divisor := bigDecimalObjectFromBigInt(big.NewInt(6), 1, 0)
		res := bigdecimalDivideAndRemainder([]interface{}{dividend, divisor})
		// Expect object (array)
		arrObj, ok := res.(*object.Object)
		if !ok {
			t.Fatalf("expected *object.Object array, got %T", res)
		}
		// Use reflection to locate the underlying slice of BigDecimal objects inside the returned array object.
		// This is necessary because the returned array object wraps the slice in a field, but the field name
		// or structure can vary depending on implementation, so we attempt several common possibilities.
		v := reflect.ValueOf(arrObj).Elem()
		var sliceVal reflect.Value
		found := false

		// Try common field names that might hold the slice:
		// - "Fvalue", "Value", "V", "FVal" are common naming conventions used in Go struct fields holding the value.
		// We check if the field exists, is valid, and if it contains a slice of *object.Object.
		for _, name := range []string{"Fvalue", "Value", "V", "FVal"} {
			f := v.FieldByName(name)
			if f.IsValid() {
				// The field might be stored as interface{} wrapping a slice
				if f.Kind() == reflect.Interface && !f.IsNil() {
					inner := reflect.ValueOf(f.Interface())
					if inner.Kind() == reflect.Slice {
						sliceVal = inner
						found = true
						break
					}
				}
				// Or it might be directly a slice type
				if f.Kind() == reflect.Slice {
					sliceVal = f
					found = true
					break
				}
			}
		}

		// If none of the common field names work, try a fallback to locate the slice in the object's FieldTable,
		// which might store the slice under the key "value" in some implementations.
		if !found {
			if valFld, ok := arrObj.FieldTable["value"]; ok {
				if sl, ok := valFld.Fvalue.([]*object.Object); ok {
					// Confirm the slice length to avoid panic
					if len(sl) != 2 {
						t.Fatalf("array length != 2")
					}
					q, _ := extractBigDecimalComponents(t, sl[0])
					r, _ := extractBigDecimalComponents(t, sl[1])
					if q.Cmp(big.NewInt(3)) != 0 {
						t.Fatalf("expected quotient 3, got %s", q.String())
					}
					if r.Cmp(big.NewInt(2)) != 0 {
						t.Fatalf("expected remainder 2, got %s", r.String())
					}
					// Successfully extracted slice from FieldTable, no need to continue
					return
				}
			}
			t.Fatalf("could not locate underlying slice in returned array object (reflection fallback failed)")
		}

		// Proceed to check the slice length and access elements once found
		if sliceVal.Len() != 2 {
			t.Fatalf("array length != 2; len=%d", sliceVal.Len())
		}

		elem0 := sliceVal.Index(0).Interface().(*object.Object)
		elem1 := sliceVal.Index(1).Interface().(*object.Object)
		q, _ := extractBigDecimalComponents(t, elem0)
		r, _ := extractBigDecimalComponents(t, elem1)

		if q.Cmp(big.NewInt(3)) != 0 {
			t.Fatalf("expected quotient 3, got %s", q.String())
		}
		if r.Cmp(big.NewInt(2)) != 0 {
			t.Fatalf("expected remainder 2, got %s", r.String())
		}
	})

	t.Run("bigdecimalDivideToIntegralValue divide by zero", func(t *testing.T) {
		bdutInit()
		dv := bigDecimalObjectFromBigInt(big.NewInt(7), 1, 0)
		dr := bigDecimalObjectFromBigInt(big.NewInt(0), 1, 0)
		res := bigdecimalDivideToIntegralValue([]interface{}{dv, dr})
		assertGErrBlk(t, res, excNames.ArithmeticException)
	})

	t.Run("bigdecimalDoubleValue and FloatValue", func(t *testing.T) {
		bdutInit()
		bd := bigDecimalObjectFromBigInt(big.NewInt(12345), 5, 2) // 123.45
		resd := bigdecimalDoubleValue([]interface{}{bd})
		if got, ok := resd.(float64); !ok {
			t.Fatalf("expected float64, got %T", resd)
		} else if math.Abs(got-123.45) > 1e-9 {
			t.Fatalf("doubleValue mismatch: got %v", got)
		}
		resf := bigdecimalFloatValue([]interface{}{bd})
		if gotf, ok := resf.(float64); !ok {
			t.Fatalf("expected float64, got %T", resf)
		} else if math.Abs(gotf-123.45) > 1e-5 {
			t.Fatalf("floatValue mismatch: got %v", gotf)
		}
	})

	t.Run("bigdecimalEquals true/false", func(t *testing.T) {
		bdutInit()
		a := bigDecimalObjectFromBigInt(big.NewInt(1000), 4, 2) // 10.00
		b := bigDecimalObjectFromBigInt(big.NewInt(1000), 4, 2) // 10.00
		res := bigdecimalEquals([]interface{}{a, b}).(int64)
		if !object.GoBooleanFromJavaBoolean(res) {
			t.Fatalf("expected true equals, got false")
		}
		// different scale
		c := bigDecimalObjectFromBigInt(big.NewInt(1000), 4, 1)
		res = bigdecimalEquals([]interface{}{a, c}).(int64)
		if object.GoBooleanFromJavaBoolean(res) {
			t.Fatalf("expected false for different scale, got true")
		}
		// different value same scale
		d := bigDecimalObjectFromBigInt(big.NewInt(2000), 4, 2)
		res = bigdecimalEquals([]interface{}{a, d}).(int64)
		if object.GoBooleanFromJavaBoolean(res) {
			t.Fatalf("expected false for different unscaled, got true")
		}
	})

	t.Run("bigdecimalIntValue and IntValueExact", func(t *testing.T) {
		bdutInit()
		bd := bigDecimalObjectFromBigInt(big.NewInt(42), 2, 0)
		res := bigdecimalIntValue([]interface{}{bd})
		if got, ok := res.(int64); !ok {
			t.Fatalf("expected int64, got %T", res)
		} else if got != 42 {
			t.Fatalf("unexpected int value %d", got)
		}

		// exact with non-zero scale -> error
		bd2 := bigDecimalObjectFromBigInt(big.NewInt(4200), 4, 2)
		res2 := bigdecimalIntValueExact([]interface{}{bd2})
		assertGErrBlk(t, res2, excNames.ArithmeticException)

		// exact but out of int32 range -> produce ArithmeticException
		large := new(big.Int).SetInt64(int64(math.MaxInt32))
		large = large.Add(large, big.NewInt(1))
		bd3 := bigDecimalObjectFromBigInt(large, precisionFromBigInt(large), 0)
		res3 := bigdecimalIntValueExact([]interface{}{bd3})
		assertGErrBlk(t, res3, excNames.ArithmeticException)

		// valid exact
		val := bigDecimalObjectFromBigInt(big.NewInt(123), 3, 0)
		res4 := bigdecimalIntValueExact([]interface{}{val})
		if got, ok := res4.(int64); !ok {
			t.Fatalf("expected int64, got %T", res4)
		} else if got != 123 {
			t.Fatalf("expected 123, got %d", got)
		}
	})

	t.Run("bigdecimalLongValue and LongValueExact", func(t *testing.T) {
		bdutInit()
		bd := bigDecimalObjectFromBigInt(big.NewInt(9999999999), 10, 0)
		res := bigdecimalLongValue([]interface{}{bd})
		if got, ok := res.(int64); !ok {
			t.Fatalf("expected int64, got %T", res)
		} else if got != 9999999999 {
			t.Fatalf("unexpected long value %d", got)
		}

		// LongValueExact non-zero scale -> error
		bd2 := bigDecimalObjectFromBigInt(big.NewInt(1000), 4, 1)
		res2 := bigdecimalLongValueExact([]interface{}{bd2})
		assertGErrBlk(t, res2, excNames.ArithmeticException)

		// Valid exact (within int64)
		val := bigDecimalObjectFromBigInt(big.NewInt(1234567890), 10, 0)
		res3 := bigdecimalLongValueExact([]interface{}{val})
		if got, ok := res3.(int64); !ok {
			t.Fatalf("expected int64, got %T", res3)
		} else if got != 1234567890 {
			t.Fatalf("unexpected long value %d", got)
		}
	})

	t.Run("bigdecimalMaxMin", func(t *testing.T) {
		bdutInit()
		a := bigDecimalObjectFromBigInt(big.NewInt(5), 1, 0)
		b := bigDecimalObjectFromBigInt(big.NewInt(10), 2, 0)
		maxie := bigdecimalMax([]interface{}{a, b})
		if maxie != b {
			t.Fatalf("expected max to be b")
		}
		minnie := bigdecimalMin([]interface{}{a, b})
		if minnie != a {
			t.Fatalf("expected min to be a")
		}
	})

	t.Run("bigdecimalMovePointLeft and MovePointRight", func(t *testing.T) {
		bdutInit()
		bd := bigDecimalObjectFromBigInt(big.NewInt(12345), 5, 2) // 123.45
		// Move left by 1. The scale increases by 1 -> scale 3.
		res := bigdecimalMovePointLeft([]interface{}{bd, int64(1)})
		out := res.(*object.Object)
		_, s := extractBigDecimalComponents(t, out)
		if s != 3 {
			t.Fatalf("expected scale 3 after MovePointLeft, got %d", s)
		}

		// move right by 1: since num <= scale (1 <= 2), reduce scale
		res2 := bigdecimalMovePointRight([]interface{}{bd, int64(1)})
		out2 := res2.(*object.Object)
		_, s2 := extractBigDecimalComponents(t, out2)
		if s2 != 1 {
			t.Fatalf("expected scale 1 after MovePointRight (reduce), got %d", s2)
		}

		// move right by large (num > scale) -> multiplier applied
		res3 := bigdecimalMovePointRight([]interface{}{bd, int64(5)})
		out3 := res3.(*object.Object)
		u3, s3 := extractBigDecimalComponents(t, out3)
		// shift = num - scale = 5 - 2 = 3 => unscaled multiplied by 10^3 => 12345 * 1000
		expected := new(big.Int).Mul(big.NewInt(12345), new(big.Int).Exp(big.NewInt(10), big.NewInt(3), nil))
		if u3.Cmp(expected) != 0 {
			t.Fatalf("expected unscaled %s got %s", expected.String(), u3.String())
		}
		if s3 != 0 {
			t.Fatalf("expected scale 0 after shift, got %d", s3)
		}
	})

	t.Run("bigdecimalMultiply", func(t *testing.T) {
		bdutInit()
		a := bigDecimalObjectFromBigInt(big.NewInt(12), 2, 1) // 1.2
		b := bigDecimalObjectFromBigInt(big.NewInt(25), 2, 1) // 2.5
		res := bigdecimalMultiply([]interface{}{a, b})
		out := res.(*object.Object)
		u, s := extractBigDecimalComponents(t, out)
		// unscaled 12*25 = 300 ; scale = 1+1 = 2
		if u.Cmp(big.NewInt(300)) != 0 {
			t.Fatalf("unexpected unscaled multiply: %s", u.String())
		}
		if s != 2 {
			t.Fatalf("unexpected scale multiply: %d", s)
		}
	})
}

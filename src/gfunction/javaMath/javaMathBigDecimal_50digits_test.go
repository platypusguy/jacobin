/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
	"testing"
)

// Helpers local to this file
func bdFromStr50(t *testing.T, s string) *object.Object {
	t.Helper()
	bd := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
	if res := BigdecimalInitString([]interface{}{bd, object.StringObjectFromGoString(s)}); res != nil {
		t.Fatalf("bdFromStr50 calling BigdecimalInitString failed for %s: %v", s, res)
	}
	return bd
}

func bdToString(t *testing.T, bd *object.Object) string {
	t.Helper()
	params := []interface{}{bd}
	return object.GoStringFromStringObject(BigdecimalToString(params).(*object.Object))
}

func Test_BigDecimal_Add_50Digits(t *testing.T) {
	bdutInit()
	// Two 50-digit numbers with decimals: a has 3 decimal places, b has 5 decimal places
	aStr := "12345678901234567890123456789012345678901234567.123" // 47 digits before dot + 3 after = 50
	bStr := "987654321098765432109876543210987654321098765.43210" // 45 digits before dot + 5 after = 50
	A := bdFromStr50(t, aStr)
	B := bdFromStr50(t, bStr)
	C := bigdecimalAdd([]interface{}{A, B}).(*object.Object)
	// Expected: align scales to max(sA,sB)=5
	aU, aS := extractBigDecimalComponents(t, A)
	bU, bS := extractBigDecimalComponents(t, B)
	s := aS
	if bS > s {
		s = bS
	}
	adjA := new(big.Int).Set(aU)
	if s > aS {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-aS), nil)
		adjA.Mul(adjA, mul)
	}
	adjB := new(big.Int).Set(bU)
	if s > bS {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-bS), nil)
		adjB.Mul(adjB, mul)
	}
	expected := new(big.Int).Add(adjA, adjB)
	u, gotS := extractBigDecimalComponents(t, C)
	if gotS != s {
		t.Fatalf("add scale expected %d, got %d", s, gotS)
	}
	if u.Cmp(expected) != 0 {
		t.Fatalf("add unscaled mismatch: got %s want %s", u.String(), expected.String())
	}
}

func Test_BigDecimal_Subtract_50Digits(t *testing.T) {
	bdutInit()
	// Two 50-digit numbers with decimals; verify subtraction. a has 3 dp, b has 5 dp
	aStr := "98765432109876543210987654321098765432109876543.210" // 47+3
	bStr := "123456789012345678901234567890123456789012345.67890" // 45+5
	A := bdFromStr50(t, aStr)
	B := bdFromStr50(t, bStr)
	res := bigdecimalSubtract([]interface{}{A, B}).(*object.Object)
	// Expected with aligned scales
	aU, aS := extractBigDecimalComponents(t, A)
	bU, bS := extractBigDecimalComponents(t, B)
	s := aS
	if bS > s {
		s = bS
	}
	adjA := new(big.Int).Set(aU)
	if s > aS {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-aS), nil)
		adjA.Mul(adjA, mul)
	}
	adjB := new(big.Int).Set(bU)
	if s > bS {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-bS), nil)
		adjB.Mul(adjB, mul)
	}
	expected := new(big.Int).Sub(adjA, adjB)
	u, gotS := extractBigDecimalComponents(t, res)
	if gotS != s {
		t.Fatalf("subtract scale expected %d, got %d", s, gotS)
	}
	if u.Cmp(expected) != 0 {
		t.Fatalf("subtract unscaled mismatch: got %s want %s", u.String(), expected.String())
	}
}

func Test_BigDecimal_Multiply_50Digits(t *testing.T) {
	bdutInit()
	// Two 50-digit numbers with decimals; product scale should be sum of scales (3+5=8)
	aStr := "12345678901234567890123456789012345678901234567.890" // 47+3
	bStr := "987654321098765432109876543210987654321098765.43210" // 45+5
	A := bdFromStr50(t, aStr)
	B := bdFromStr50(t, bStr)
	res := bigdecimalMultiply([]interface{}{A, B}).(*object.Object)
	// Expected via unscaled multiplication; scale = sA + sB
	aU, aS := extractBigDecimalComponents(t, A)
	bU, bS := extractBigDecimalComponents(t, B)
	expected := new(big.Int).Mul(new(big.Int).Set(aU), new(big.Int).Set(bU))
	u, s := extractBigDecimalComponents(t, res)
	if s != aS+bS {
		t.Fatalf("multiply scale expected %d, got %d", aS+bS, s)
	}
	if u.Cmp(expected) != 0 {
		t.Fatalf("multiply unscaled mismatch: got %s want %s", u.String(), expected.String())
	}
}

func Test_BigDecimal_Divide_50Digits_ExactMatch(t *testing.T) {
	bdutInit()

	// dividend
	aStr := "963963963963963963963963963963963963963963963963963963963.963" // 57+3
	Div := bdFromStr50(t, aStr)

	// divisor
	bStr := "2.87654876548765487654876548765487654876548765487654"
	Dsr := bdFromStr50(t, bStr)

	// Perform division with specified scale and rounding (3 dp, HALF_UP)
	rm := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("HALF_UP")})
	if blk, ok := rm.(*ghelpers.GErrBlk); ok {
		t.Fatalf("failed to get RoundingMode.HALF_UP: %v", blk)
	}
	Quotient := bigdecimalDivideScaleRoundingMode([]interface{}{Div, Dsr, int64(3), rm.(*object.Object)})

	// Expected 3-decimal HALF_UP rounding of approximately a/3
	expected := "335111288753189383117212577810800663414238150377653296938.250"

	// Check quotient or other result.
	switch Quotient.(type) {
	case *object.Object:
		observed := bdToString(t, Quotient.(*object.Object))
		if observed != expected {
			t.Fatalf("BigDecimal divide result mismatch: got %s want %s", observed, expected)
		}
		t.Logf("BigDecimal divide result: %s", observed)
	case *ghelpers.GErrBlk:
		t.Fatalf("unexpected G function error block, got %s", Quotient.(*ghelpers.GErrBlk).ErrMsg)
	default:
		t.Fatalf("unexpected G function return, got %T", Quotient)
	}
}

func Test_BigDecimal_DivideAndRemainder_IssueCase(t *testing.T) {
	bdutInit()
	A := bdFromStr50(t, "963963963963963963963963963963963963963963963963963963963.963")
	B := bdFromStr50(t, "2.87654876548765487654876548765487654876548765487654")
	arr := bigdecimalDivideAndRemainder([]interface{}{A, B})
	if blk, ok := arr.(*ghelpers.GErrBlk); ok {
		t.Fatalf("divideAndRemainder returned error: %v", blk)
	}
	arrObj := arr.(*object.Object)
	vals, _ := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(vals) != 2 {
		t.Fatalf("expected array of length 2, got %d", len(vals))
	}
	qStr := bdToString(t, vals[0])
	rStr := bdToString(t, vals[1])
	expectedQ := "335111288753189383117212577810800663414238150377653296938"
	expectedR := "0.71804475666634227855232831541824534553677564996548"
	if qStr != expectedQ {
		t.Fatalf("divideAndRemainder quotient mismatch: got %s want %s", qStr, expectedQ)
	}
	if rStr != expectedR {
		t.Fatalf("divideAndRemainder remainder mismatch: got %s want %s", rStr, expectedR)
	}
}

func Test_BigDecimal_Pow_5_Scaled_String_IssueCase(t *testing.T) {
	bdutInit()
	B := bdFromStr50(t, "2.87654876548765487654876548765487654876548765487654")
	res := bigdecimalPow([]interface{}{B, int64(5)})
	if blk, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("pow returned error: %v", blk)
	}
	C := res.(*object.Object)
	observed := bdToString(t, C)
	expected := "196.9512332632041425648903952282597858571550912228881576467059003975279194556093059395958737271613605818429385286431679948696259759708157256485445979026081037557688179566936520311094050547107546448817694942476182957384867462249103941371639497016916093024"
	if observed != expected {
		t.Fatalf("BigDecimal pow(5) mismatch: got %s want %s", observed, expected)
	}
}

func Test_BigDecimal_LongValue_ScaledOverflow_IssueCase(t *testing.T) {
	bdutInit()
	// A from the issue description
	A := bdFromStr50(t, "963963963963963963963963963963963963963963963963963963963.963")
	res := bigdecimalLongValue([]interface{}{A})
	got, ok := res.(int64)
	if !ok {
		t.Fatalf("expected int64 from longValue, got %T", res)
	}
	var expected int64 = -5163478405204313541
	if got != expected {
		t.Fatalf("A.longValue() mismatch: got %d want %d", got, expected)
	}
}

func Test_BigDecimal_Divide_50Digits_ArithmeticException(t *testing.T) {
	bdutInit()
	// These specific operands do not divide to a terminating decimal; expect ArithmeticException
	// dividend
	aStr := "987654321098765432109876543210987654321098765.43210" // 45+5
	Div := bdFromStr50(t, aStr)
	// divisor
	bStr := "12345678901234567890123456789012345678901234567.890" // 47+3
	Dsr := bdFromStr50(t, bStr)
	// expected quotient
	cStr := "3.0"
	expected := bdFromStr50(t, cStr)
	t.Logf("expected quotient: %T", expected)

	res := bigdecimalDivide([]interface{}{Div, Dsr})
	blk, isGErrBlk := res.(*ghelpers.GErrBlk)
	if !isGErrBlk {
		t.Fatalf("expected ArithmeticException for non-terminating division, got %T", blk)
	}
}

// Added test for toBigInteger issue case
func Test_BigDecimal_ToBigInteger_IssueCase(t *testing.T) {
	bdutInit()
	A := bdFromStr50(t, "963963963963963963963963963963963963963963963963963963963.963")
	res := bigdecimalToBigInteger([]interface{}{A})
	bi, ok := res.(*object.Object)
	if !ok || bi == nil {
		t.Fatalf("toBigInteger did not return a BigInteger object: %T", res)
	}
	got := bi.FieldTable["value"].Fvalue.(*big.Int).String()
	expected := "963963963963963963963963963963963963963963963963963963963"
	if got != expected {
		t.Fatalf("BigDecimal.toBigInteger mismatch: got %s want %s", got, expected)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"math"
	"testing"
)

func newLongSummaryStatisticsObj() *object.Object {
	className := "java/util/LongSummaryStatistics"
	return object.MakeEmptyObjectWithClassName(&className)
}

func TestLongSummaryStatistics_InitDefaults(t *testing.T) {
	globals.InitStringPool()

	lss := newLongSummaryStatisticsObj()
	ret := longSummaryStatisticsInit([]interface{}{lss})
	if ret != nil {
		t.Fatalf("expected init to return nil, got %v", ret)
	}

	if v := longSummaryStatisticsGetCount([]interface{}{lss}); v != int64(0) {
		t.Errorf("expected count 0, got %v", v)
	}
	if v := longSummaryStatisticsGetSum([]interface{}{lss}); v != int64(0) {
		t.Errorf("expected sum 0, got %v", v)
	}
	if v := longSummaryStatisticsGetMin([]interface{}{lss}); v != int64(math.MaxInt64) {
		t.Errorf("expected min Long.MAX_VALUE, got %v", v)
	}
	if v := longSummaryStatisticsGetMax([]interface{}{lss}); v != int64(math.MinInt64) {
		t.Errorf("expected max Long.MIN_VALUE, got %v", v)
	}
	if v := longSummaryStatisticsGetAverage([]interface{}{lss}); v != float64(0) {
		t.Errorf("expected average 0, got %v", v)
	}
}

func TestLongSummaryStatistics_InitJJJJ_Valid(t *testing.T) {
	lss := newLongSummaryStatisticsObj()
	ret := longSummaryStatisticsInitJJJJ([]interface{}{lss, int64(3), int64(1), int64(5), int64(9)})
	if ret != nil {
		t.Errorf("expected nil return, got %v", ret)
	}
	if v := longSummaryStatisticsGetCount([]interface{}{lss}); v != int64(3) {
		t.Errorf("expected count 3, got %v", v)
	}
	if v := longSummaryStatisticsGetSum([]interface{}{lss}); v != int64(9) {
		t.Errorf("expected sum 9, got %v", v)
	}
	if v := longSummaryStatisticsGetMin([]interface{}{lss}); v != int64(1) {
		t.Errorf("expected min 1, got %v", v)
	}
	if v := longSummaryStatisticsGetMax([]interface{}{lss}); v != int64(5) {
		t.Errorf("expected max 5, got %v", v)
	}
}

func TestLongSummaryStatistics_InitJJJJ_ZeroCount(t *testing.T) {
	lss := newLongSummaryStatisticsObj()
	ret := longSummaryStatisticsInitJJJJ([]interface{}{lss, int64(0), int64(1), int64(5), int64(9)})
	if ret != nil {
		t.Errorf("expected nil return, got %v", ret)
	}
	if v := longSummaryStatisticsGetCount([]interface{}{lss}); v != int64(0) {
		t.Errorf("expected count 0, got %v", v)
	}
	if v := longSummaryStatisticsGetSum([]interface{}{lss}); v != int64(0) {
		t.Errorf("expected sum 0, got %v", v)
	}
	if v := longSummaryStatisticsGetMin([]interface{}{lss}); v != int64(math.MaxInt64) {
		t.Errorf("expected min MaxInt64, got %v", v)
	}
	if v := longSummaryStatisticsGetMax([]interface{}{lss}); v != int64(math.MinInt64) {
		t.Errorf("expected max MinInt64, got %v", v)
	}
}

func TestLongSummaryStatistics_InitJJJJ_NegativeCount(t *testing.T) {
	lss := newLongSummaryStatisticsObj()
	ret := longSummaryStatisticsInitJJJJ([]interface{}{lss, int64(-1), int64(1), int64(5), int64(9)})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestLongSummaryStatistics_InitJJJJ_MinGreaterThanMax(t *testing.T) {
	lss := newLongSummaryStatisticsObj()
	ret := longSummaryStatisticsInitJJJJ([]interface{}{lss, int64(3), int64(10), int64(5), int64(9)})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestLongSummaryStatistics_InitJJJJ_TooFewArguments(t *testing.T) {
	lss := newLongSummaryStatisticsObj()
	ret := longSummaryStatisticsInitJJJJ([]interface{}{lss, int64(3)})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestLongSummaryStatistics_Accept(t *testing.T) {
	globals.InitStringPool()

	lss := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss})

	_ = longSummaryStatisticsAccept([]interface{}{lss, int64(10)})
	_ = longSummaryStatisticsAccept([]interface{}{lss, int64(20)})
	_ = longSummaryStatisticsAccept([]interface{}{lss, int64(5)})

	if v := longSummaryStatisticsGetCount([]interface{}{lss}); v != int64(3) {
		t.Errorf("expected count 3, got %v", v)
	}
	if v := longSummaryStatisticsGetSum([]interface{}{lss}); v != int64(35) {
		t.Errorf("expected sum 35, got %v", v)
	}
	if v := longSummaryStatisticsGetMin([]interface{}{lss}); v != int64(5) {
		t.Errorf("expected min 5, got %v", v)
	}
	if v := longSummaryStatisticsGetMax([]interface{}{lss}); v != int64(20) {
		t.Errorf("expected max 20, got %v", v)
	}

	avg := longSummaryStatisticsGetAverage([]interface{}{lss}).(float64)
	expectedAvg := float64(35) / float64(3)
	if math.Abs(avg-expectedAvg) > 1e-9 {
		t.Errorf("expected average %f, got %f", expectedAvg, avg)
	}
}

func TestLongSummaryStatistics_Accept_TooFewArguments(t *testing.T) {
	result := longSummaryStatisticsAccept([]interface{}{}).(*ghelpers.GErrBlk)
	if result.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for too few arguments, got %v", result.ExceptionType)
	}
}

func TestLongSummaryStatistics_Accept_InvalidReceiver(t *testing.T) {
	result := longSummaryStatisticsAccept([]interface{}{"not an object", int64(1)}).(*ghelpers.GErrBlk)
	if result.ExceptionType != excNames.ClassCastException {
		t.Errorf("Expected ClassCastException for invalid receiver, got %v", result.ExceptionType)
	}
}

func TestLongSummaryStatistics_Accept_InvalidArgument(t *testing.T) {
	globals.InitStringPool()
	lss := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss})

	result := longSummaryStatisticsAccept([]interface{}{lss, "not a long"}).(*ghelpers.GErrBlk)
	if result.ExceptionType != excNames.ClassCastException {
		t.Errorf("Expected ClassCastException for invalid argument, got %v", result.ExceptionType)
	}
}

func TestLongSummaryStatistics_Combine(t *testing.T) {
	globals.InitStringPool()

	lss1 := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss1})
	_ = longSummaryStatisticsAccept([]interface{}{lss1, int64(1)})
	_ = longSummaryStatisticsAccept([]interface{}{lss1, int64(2)})

	lss2 := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss2})
	_ = longSummaryStatisticsAccept([]interface{}{lss2, int64(10)})
	_ = longSummaryStatisticsAccept([]interface{}{lss2, int64(-5)})

	ret := longSummaryStatisticsCombine([]interface{}{lss1, lss2})
	if ret != nil {
		t.Fatalf("expected combine to return nil, got %v", ret)
	}

	if v := longSummaryStatisticsGetCount([]interface{}{lss1}); v != int64(4) {
		t.Errorf("expected count 4, got %v", v)
	}
	if v := longSummaryStatisticsGetSum([]interface{}{lss1}); v != int64(8) {
		t.Errorf("expected sum 8, got %v", v)
	}
	if v := longSummaryStatisticsGetMin([]interface{}{lss1}); v != int64(-5) {
		t.Errorf("expected min -5, got %v", v)
	}
	if v := longSummaryStatisticsGetMax([]interface{}{lss1}); v != int64(10) {
		t.Errorf("expected max 10, got %v", v)
	}
}

func TestLongSummaryStatistics_Combine_TooFewArguments(t *testing.T) {
	result := longSummaryStatisticsCombine([]interface{}{}).(*ghelpers.GErrBlk)
	if result.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for too few arguments, got %v", result.ExceptionType)
	}
}

func TestLongSummaryStatistics_Combine_InvalidReceiver(t *testing.T) {
	result := longSummaryStatisticsCombine([]interface{}{"not an object", "also not"}).(*ghelpers.GErrBlk)
	if result.ExceptionType != excNames.ClassCastException {
		t.Errorf("Expected ClassCastException for invalid receiver, got %v", result.ExceptionType)
	}
}

func TestLongSummaryStatistics_Combine_InvalidArgument(t *testing.T) {
	globals.InitStringPool()
	lss := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss})

	result := longSummaryStatisticsCombine([]interface{}{lss, "not an object"}).(*ghelpers.GErrBlk)
	if result.ExceptionType != excNames.ClassCastException {
		t.Errorf("Expected ClassCastException for invalid argument, got %v", result.ExceptionType)
	}
}

func TestLongSummaryStatistics_ToString(t *testing.T) {
	globals.InitStringPool()

	lss := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss})
	_ = longSummaryStatisticsAccept([]interface{}{lss, int64(12)})
	_ = longSummaryStatisticsAccept([]interface{}{lss, int64(23)})
	_ = longSummaryStatisticsAccept([]interface{}{lss, int64(10)})

	res := longSummaryStatisticsToString([]interface{}{lss})
	strObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}

	str := object.GoStringFromStringObject(strObj)
	expected := "LongSummaryStatistics{count=3, sum=45, min=10, average=15.000000, max=23}"
	if str != expected {
		t.Errorf("expected %q, got %q", expected, str)
	}
}

func TestLongSummaryStatistics_ToString_EmptyStatistics(t *testing.T) {
	globals.InitStringPool()

	lss := newLongSummaryStatisticsObj()
	_ = longSummaryStatisticsInit([]interface{}{lss})

	res := longSummaryStatisticsToString([]interface{}{lss})
	strObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}

	str := object.GoStringFromStringObject(strObj)
	expected := "LongSummaryStatistics{count=0, sum=0, min=9223372036854775807, average=0.000000, max=-9223372036854775808}"
	if str != expected {
		t.Errorf("expected %q, got %q", expected, str)
	}
}

func TestLongSummaryStatistics_MethodSignatures(t *testing.T) {
	Load_Util_LongSummaryStatistics()

	expectedMethods := []string{
		"java/util/LongSummaryStatistics.<clinit>()V",
		"java/util/LongSummaryStatistics.<init>()V",
		"java/util/LongSummaryStatistics.accept(J)V",
		"java/util/LongSummaryStatistics.combine(Ljava/util/LongSummaryStatistics;)V",
		"java/util/LongSummaryStatistics.getCount()J",
		"java/util/LongSummaryStatistics.getSum()J",
		"java/util/LongSummaryStatistics.getMin()J",
		"java/util/LongSummaryStatistics.getMax()J",
		"java/util/LongSummaryStatistics.getAverage()D",
		"java/util/LongSummaryStatistics.toString()Ljava/lang/String;",
	}

	for _, m := range expectedMethods {
		if _, ok := ghelpers.MethodSignatures[m]; !ok {
			t.Errorf("expected method signature %q to be registered", m)
		}
	}
}

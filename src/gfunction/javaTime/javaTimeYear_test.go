/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"strconv"
	"strings"
	"testing"
)

func yearTestInit() {
	globals.InitGlobals("test")
	Load_Time_Year()
}

func yearTestObject(t *testing.T, value any) *object.Object {
	t.Helper()
	result, ok := value.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", value)
	}
	return result
}

func yearTestValue(t *testing.T, value any) int64 {
	t.Helper()
	result, ok := getInt64Val(value)
	if !ok {
		t.Fatalf("expected integer value, got %T (%v)", value, value)
	}
	return result
}

func yearTestString(t *testing.T, value any) string {
	t.Helper()
	result, ok := value.(*object.Object)
	if !ok || !object.IsStringObject(result) {
		t.Fatalf("expected String object, got %T", value)
	}
	return object.GoStringFromStringObject(result)
}

func yearTestError(t *testing.T, value any, exception int) {
	t.Helper()
	result, ok := value.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected error block, got %T (%v)", value, value)
	}
	if result.ExceptionType != exception {
		t.Errorf("expected exception type %d, got %d (%s)", exception, result.ExceptionType, result.ErrMsg)
	}
}

func yearTestMonthDay(month, day int64) *object.Object {
	className := "java/time/MonthDay"
	monthDay := object.MakeEmptyObjectWithClassName(&className)
	monthDay.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: month}
	monthDay.FieldTable["day"] = object.Field{Ftype: types.Int, Fvalue: day}
	return monthDay
}

func TestYearClinitAndRegistration(t *testing.T) {
	yearTestInit()

	for _, signature := range []string{
		"java/time/Year.<clinit>()V",
		"java/time/Year.of(I)Ljava/time/Year;",
		"java/time/Year.parse(Ljava/lang/CharSequence;)Ljava/time/Year;",
		"java/time/Year.plusYears(J)Ljava/time/Year;",
		"java/time/Year.atDay(I)Ljava/time/LocalDate;",
		"java/time/Year.toString()Ljava/lang/String;",
	} {
		if _, ok := ghelpers.MethodSignatures[signature]; !ok {
			t.Errorf("Year method %q was not registered", signature)
		}
	}

	if result := yearClinit(nil); result != nil {
		t.Errorf("Year.<clinit> returned %v, expected nil", result)
	}
	for name, want := range map[string]int64{"MIN_VALUE": yearMinValue, "MAX_VALUE": yearMaxValue} {
		constant, ok := statics.QueryStatic(classNameYearVal, name)
		if !ok {
			t.Fatalf("Year.%s was not initialized", name)
		}
		if got := yearTestValue(t, constant.Value); got != want {
			t.Errorf("Year.%s = %d, want %d", name, got, want)
		}
	}

	for _, value := range []int64{yearMinValue, -1, 0, 1, yearMaxValue} {
		year := yearTestObject(t, yearOf([]any{value}))
		if got := yearTestValue(t, yearGetValue([]any{year})); got != value {
			t.Errorf("Year.of(%d).getValue() = %d", value, got)
		}
	}
	for _, value := range []int64{yearMinValue - 1, yearMaxValue + 1} {
		yearTestError(t, yearOf([]any{value}), excNames.DateTimeException)
	}
	yearTestError(t, yearOf(nil), excNames.IllegalArgumentException)

	instance := object.MakeEmptyObjectWithClassName(&classNameYearVal)
	if result := yearInit([]any{instance, int64(42)}); result != nil {
		t.Fatalf("Year.<init> returned %v", result)
	}
	if got := yearTestValue(t, yearGetValue([]any{instance})); got != 42 {
		t.Errorf("constructed Year value = %d, want 42", got)
	}
	yearTestError(t, yearInit([]any{instance, int64(yearMaxValue + 1)}), excNames.DateTimeException)
}

func TestYearLeapYearAndFormatting(t *testing.T) {
	yearTestInit()

	for _, test := range []struct {
		year   int64
		leap   bool
		month  int64
		length int64
	}{
		{2000, true, 2, 366},
		{1900, false, 1, 365},
		{2024, true, 2, 366},
		{-400, true, 2, 366},
		{-100, false, 1, 365},
	} {
		year := yearTestObject(t, yearOf([]any{test.year}))
		wantLeap := types.ConvertGoBoolToJavaBool(test.leap)
		if got := yearIsLeapStatic([]any{test.year}); got != wantLeap {
			t.Errorf("isLeap(%d) = %v, want %v", test.year, got, wantLeap)
		}
		if got := yearIsLeapInstance([]any{year}); got != wantLeap {
			t.Errorf("Year(%d).isLeap() = %v, want %v", test.year, got, wantLeap)
		}
		if got := yearLength([]any{year}); yearTestValue(t, got) != test.length {
			t.Errorf("Year(%d).length() = %v, want %d", test.year, got, test.length)
		}
		if got := yearIsValidMonthDay([]any{year, yearTestMonthDay(2, 29)}); got != wantLeap {
			t.Errorf("Year(%d).isValidMonthDay(February 29) = %v, want %v", test.year, got, wantLeap)
		}
	}

	for _, test := range []struct {
		year int64
		text string
	}{
		{0, "0000"},
		{1, "0001"},
		{-1, "-0001"},
		{999, "0999"},
		{1000, "1000"},
		{9999, "9999"},
		{10000, "+10000"},
		{-10000, "-10000"},
	} {
		year := yearTestObject(t, yearOf([]any{test.year}))
		if got := yearTestString(t, yearToString([]any{year})); got != test.text {
			t.Errorf("Year(%d).toString() = %q, want %q", test.year, got, test.text)
		}
	}

	for _, text := range []string{"2024", " -1 ", "0"} {
		parsed := yearTestObject(t, yearParse([]any{object.StringObjectFromGoString(text)}))
		want, _ := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
		if got := yearTestValue(t, yearGetValue([]any{parsed})); got != want {
			t.Errorf("Year.parse(%q) = %d, want %d", text, got, want)
		}
	}
	formatter := object.MakeEmptyObjectWithClassName(&classNameYearVal)
	if got := yearTestString(t, yearFormat([]any{yearOf([]any{2024}), formatter})); got != "2024" {
		t.Errorf("Year.format returned %q, want %q", got, "2024")
	}
	yearTestError(t, yearParse([]any{object.StringObjectFromGoString("not-a-year")}), excNames.DateTimeParseException)
	yearTestError(t, yearParse([]any{object.StringObjectFromGoString("1000000000")}), excNames.DateTimeParseException)
	yearTestError(t, yearParse([]any{object.StringObjectFromGoString("")}), excNames.DateTimeParseException)
	parsed := yearTestObject(t, yearParseFormatter([]any{object.StringObjectFromGoString("2024"), formatter}))
	if got := yearTestValue(t, yearGetValue([]any{parsed})); got != 2024 {
		t.Errorf("Year.parse with formatter = %d, want 2024", got)
	}
}

func TestYearTemporalFieldsAndArithmetic(t *testing.T) {
	yearTestInit()

	year := yearTestObject(t, yearOf([]any{2020}))
	bceYear := yearTestObject(t, yearOf([]any{-1}))
	yearField := getChronoField("YEAR")
	yearOfEraField := getChronoField("YEAR_OF_ERA")
	eraField := getChronoField("ERA")
	monthField := getChronoField("MONTH_OF_YEAR")
	for _, field := range []*object.Object{yearField, yearOfEraField, eraField} {
		if got := yearIsSupportedField([]any{year, field}); got != types.JavaBoolTrue {
			t.Errorf("Year should support field %q", getTemporalFieldName(field))
		}
	}
	if got := yearIsSupportedField([]any{year, monthField}); got != types.JavaBoolFalse {
		t.Error("Year should not support MONTH_OF_YEAR")
	}
	if got := yearTestValue(t, yearGetLong([]any{year, yearField})); got != 2020 {
		t.Errorf("YEAR = %d, want 2020", got)
	}
	if got := yearTestValue(t, yearGetLong([]any{bceYear, yearOfEraField})); got != 2 {
		t.Errorf("BCE YEAR_OF_ERA = %d, want 2", got)
	}
	if got := yearTestValue(t, yearGetLong([]any{bceYear, eraField})); got != 0 {
		t.Errorf("BCE ERA = %d, want 0", got)
	}

	for _, test := range []struct {
		field string
		min   int64
		max   int64
	}{
		{"YEAR", yearMinValue, yearMaxValue},
		{"YEAR_OF_ERA", 1, 1000000000},
		{"ERA", 0, 1},
	} {
		rangeObj := yearTestObject(t, yearRange([]any{bceYear, getChronoField(test.field)}))
		if got := yearTestValue(t, rangeObj.FieldTable["minSmallest"].Fvalue); got != test.min {
			t.Errorf("%s minSmallest = %d, want %d", test.field, got, test.min)
		}
		if got := yearTestValue(t, rangeObj.FieldTable["maxLargest"].Fvalue); got != test.max {
			t.Errorf("%s maxLargest = %d, want %d", test.field, got, test.max)
		}
	}
	if got := yearTestObject(t, yearWithField([]any{bceYear, yearOfEraField, int64(3)})); yearTestValue(t, yearGetValue([]any{got})) != -2 {
		t.Errorf("setting BCE YEAR_OF_ERA to 3 produced %d, want -2", yearTestValue(t, yearGetValue([]any{got})))
	}
	if got := yearTestObject(t, yearWithField([]any{year, eraField, int64(0)})); yearTestValue(t, yearGetValue([]any{got})) != -2019 {
		t.Errorf("setting ERA to 0 produced %d, want -2019", yearTestValue(t, yearGetValue([]any{got})))
	}
	yearTestError(t, yearWithField([]any{year, yearField, int64(yearMaxValue + 1)}), excNames.DateTimeException)
	yearTestError(t, yearGetLong([]any{year, monthField}), excNames.DateTimeException)

	decades := ghelpers.CreateTemporalUnit("Decades")
	if got := yearTestObject(t, yearPlusUnit([]any{year, int64(2), decades})); yearTestValue(t, yearGetValue([]any{got})) != 2040 {
		t.Errorf("2020 plus two decades = %d, want 2040", yearTestValue(t, yearGetValue([]any{got})))
	}
	if got := yearTestObject(t, yearMinusYears([]any{year, int64(3)})); yearTestValue(t, yearGetValue([]any{got})) != 2017 {
		t.Errorf("2020 minus three years = %d, want 2017", yearTestValue(t, yearGetValue([]any{got})))
	}
	if got := yearTestValue(t, yearUntil([]any{year, yearOf([]any{2050}), decades})); got != 3 {
		t.Errorf("years until 2050 in decades = %d, want 3", got)
	}
	if got := yearTestValue(t, yearUntil([]any{bceYear, yearOf([]any{1}), ghelpers.CreateTemporalUnit("Eras")})); got != 1 {
		t.Errorf("BCE-to-CE eras until = %d, want 1", got)
	}
	if got := yearTestObject(t, yearPlusYears([]any{year, int64(0)})); got != year {
		t.Error("adding zero years should return the original Year object")
	}
	yearTestError(t, yearPlusUnit([]any{year, int64(1), ghelpers.CreateTemporalUnit("Days")}), excNames.DateTimeException)
	yearTestError(t, yearPlusYears([]any{year, int64(yearMaxValue)}), excNames.DateTimeException)
}

func TestYearDatesAndComparisons(t *testing.T) {
	yearTestInit()

	leapYear := yearTestObject(t, yearOf([]any{2024}))
	commonYear := yearTestObject(t, yearOf([]any{2023}))
	leapDay := yearTestObject(t, yearAtDay([]any{leapYear, int64(60)}))
	if got := yearTestValue(t, leapDay.FieldTable["month"].Fvalue); got != 2 {
		t.Errorf("2024 day 60 month = %d, want February", got)
	}
	if got := yearTestValue(t, leapDay.FieldTable["day"].Fvalue); got != 29 {
		t.Errorf("2024 day 60 = %d, want 29", got)
	}
	lastDay := yearTestObject(t, yearAtDay([]any{commonYear, int64(365)}))
	if got := yearTestValue(t, lastDay.FieldTable["month"].Fvalue); got != 12 {
		t.Errorf("2023 day 365 month = %d, want December", got)
	}
	if got := yearTestValue(t, lastDay.FieldTable["day"].Fvalue); got != 31 {
		t.Errorf("2023 day 365 = %d, want 31", got)
	}
	yearTestError(t, yearAtDay([]any{commonYear, int64(366)}), excNames.DateTimeException)

	month := yearTestObject(t, yearAtMonthInt([]any{leapYear, int64(2)}))
	if got := yearTestValue(t, month.FieldTable["year"].Fvalue); got != 2024 {
		t.Errorf("atMonth year = %d, want 2024", got)
	}
	if got := yearTestValue(t, month.FieldTable["month"].Fvalue); got != 2 {
		t.Errorf("atMonth month = %d, want 2", got)
	}
	monthObject := object.MakeEmptyObjectWithClassName(&classNameMonthVal)
	monthObject.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	if got := yearTestObject(t, yearAtMonthMonth([]any{leapYear, monthObject})); yearTestValue(t, got.FieldTable["month"].Fvalue) != 2 {
		t.Errorf("atMonth(Month.FEBRUARY) returned month %d, want 2", yearTestValue(t, got.FieldTable["month"].Fvalue))
	}
	monthDay := yearTestMonthDay(2, 29)
	if got := yearTestObject(t, yearAtMonthDay([]any{leapYear, monthDay})); yearTestValue(t, got.FieldTable["day"].Fvalue) != 29 {
		t.Errorf("atMonthDay February 29 returned day %d", yearTestValue(t, got.FieldTable["day"].Fvalue))
	}
	yearTestError(t, yearAtMonthDay([]any{commonYear, monthDay}), excNames.DateTimeException)

	if got := yearTestObject(t, yearFrom([]any{leapDay})); yearTestValue(t, yearGetValue([]any{got})) != 2024 {
		t.Errorf("Year.from(date) = %d, want 2024", yearTestValue(t, yearGetValue([]any{got})))
	}
	invalidYear := object.MakeEmptyObjectWithClassName(&classNameYearVal)
	if got := yearTestObject(t, yearFrom([]any{invalidYear})); got != invalidYear {
		t.Error("Year.from should leave an invalid Year instance unchanged")
	}

	if got := yearTestValue(t, yearCompareTo([]any{commonYear, leapYear})); got != -1 {
		t.Errorf("2023.compareTo(2024) = %d, want -1", got)
	}
	if got := yearTestValue(t, yearCompareTo([]any{leapYear, commonYear})); got != 1 {
		t.Errorf("2024.compareTo(2023) = %d, want 1", got)
	}
	if yearEquals([]any{leapYear, yearOf([]any{2024})}) != types.JavaBoolTrue {
		t.Error("Years with the same value should be equal")
	}
	if yearIsAfter([]any{leapYear, commonYear}) != types.JavaBoolTrue || yearIsBefore([]any{commonYear, leapYear}) != types.JavaBoolTrue {
		t.Error("Year before/after comparisons returned false for ordered years")
	}
	if got := yearTestValue(t, yearHashCode([]any{leapYear})); got != 2024 {
		t.Errorf("Year hash code = %d, want 2024", got)
	}
}

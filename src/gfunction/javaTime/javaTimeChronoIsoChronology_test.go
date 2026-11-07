/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestIsoChronology(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Chrono_IsoChronology()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ValueRange()

	// <clinit> and INSTANCE static field
	_ = isoChronologyClinit([]any{})
	staticVal := statics.GetStaticValue("java/time/chrono/IsoChronology", "INSTANCE")
	isoChronologyObject, okChrono := staticVal.(*object.Object)
	if !okChrono || object.IsNull(isoChronologyObject) {
		t.Fatalf("Expected non-null Object for IsoChronology INSTANCE, got %v", staticVal)
	}

	// getId & getCalendarType
	idResult := chronologyGetId([]any{isoChronologyObject})
	idString := object.GoStringFromStringObject(idResult.(*object.Object))
	if idString != "ISO" {
		t.Errorf("Expected getId to return 'ISO', got '%s'", idString)
	}

	calendarTypeResult := chronologyGetCalendarType([]any{isoChronologyObject})
	calendarTypeString := object.GoStringFromStringObject(calendarTypeResult.(*object.Object))
	if calendarTypeString != "iso8601" {
		t.Errorf("Expected getCalendarType to return 'iso8601', got '%s'", calendarTypeString)
	}

	// isLeapYear
	if !isJavaBoolTrue(isoChronologyIsLeapYear([]any{isoChronologyObject, int64(2024)})) {
		t.Errorf("Expected 2024 to be a leap year")
	}
	if !isJavaBoolTrue(isoChronologyIsLeapYear([]any{isoChronologyObject, int64(2000)})) {
		t.Errorf("Expected 2000 to be a leap year")
	}
	if isJavaBoolTrue(isoChronologyIsLeapYear([]any{isoChronologyObject, int64(1900)})) {
		t.Errorf("Expected 1900 not to be a leap year")
	}
	if isJavaBoolTrue(isoChronologyIsLeapYear([]any{isoChronologyObject, int64(2023)})) {
		t.Errorf("Expected 2023 not to be a leap year")
	}

	// eraOf
	ceEraResult := isoChronologyEraOf([]any{isoChronologyObject, int64(1)})
	ceEraObject, okCe := ceEraResult.(*object.Object)
	if !okCe || object.IsNull(ceEraObject) {
		t.Fatalf("Expected non-null Era object for eraOf(1)")
	}

	bceEraResult := isoChronologyEraOf([]any{isoChronologyObject, int64(0)})
	bceEraObject, okBce := bceEraResult.(*object.Object)
	if !okBce || object.IsNull(bceEraObject) {
		t.Fatalf("Expected non-null Era object for eraOf(0)")
	}

	invalidEraResult := isoChronologyEraOf([]any{isoChronologyObject, int64(2)})
	if !isGErrBlk(invalidEraResult) {
		t.Errorf("Expected error block for eraOf(2), got %v", invalidEraResult)
	}

	// eras
	erasResult := chronologyEras([]any{isoChronologyObject})
	erasListObject, okErasList := erasResult.(*object.Object)
	if !okErasList || object.IsNull(erasListObject) {
		t.Errorf("Expected non-null eras list")
	}

	// prolepticYear
	ceProlepticYear := isoChronologyProlepticYear([]any{isoChronologyObject, ceEraObject, int64(2026)})
	if ceProlepticYear.(int64) != 2026 {
		t.Errorf("Expected prolepticYear(CE, 2026) to be 2026, got %v", ceProlepticYear)
	}

	bceProlepticYear := isoChronologyProlepticYear([]any{isoChronologyObject, bceEraObject, int64(1)})
	if bceProlepticYear.(int64) != 0 {
		t.Errorf("Expected prolepticYear(BCE, 1) to be 0, got %v", bceProlepticYear)
	}

	// range
	dayOfMonthField := getChronoFieldInstance("DAY_OF_MONTH")
	rangeResult := isoChronologyRange([]any{isoChronologyObject, dayOfMonthField})
	rangeObject, okRange := rangeResult.(*object.Object)
	if !okRange || object.IsNull(rangeObject) {
		t.Fatalf("Expected non-null ValueRange for DAY_OF_MONTH")
	}
	minSmallest := rangeObject.FieldTable["minSmallest"].Fvalue.(int64)
	maxLargest := rangeObject.FieldTable["maxLargest"].Fvalue.(int64)
	if minSmallest != 1 || maxLargest != 31 {
		t.Errorf("Expected DAY_OF_MONTH range 1..31, got %d..%d", minSmallest, maxLargest)
	}

	// date(prolepticYear, month, dayOfMonth)
	dateYmdResult := chronologyDateYMD([]any{isoChronologyObject, int64(2026), int64(11), int64(5)})
	dateYmdObject, okDateYmd := dateYmdResult.(*object.Object)
	if !okDateYmd || object.IsNull(dateYmdObject) {
		t.Fatalf("Expected non-null LocalDate from date(2026, 11, 5)")
	}

	// date(Era, yearOfEra, month, dayOfMonth)
	dateEraResult := chronologyDateEra([]any{isoChronologyObject, ceEraObject, int64(2026), int64(11), int64(5)})
	dateEraObject, okDateEra := dateEraResult.(*object.Object)
	if !okDateEra || object.IsNull(dateEraObject) {
		t.Fatalf("Expected non-null LocalDate from date(CE, 2026, 11, 5)")
	}

	// date(TemporalAccessor)
	dateTemporalResult := chronologyDateTemporal([]any{isoChronologyObject, dateYmdObject})
	if object.IsNull(dateTemporalResult.(*object.Object)) {
		t.Errorf("Expected non-null LocalDate from date(TemporalAccessor)")
	}

	// dateEpochDay
	dateEpochDayResult := chronologyDateEpochDay([]any{isoChronologyObject, int64(0)})
	if object.IsNull(dateEpochDayResult.(*object.Object)) {
		t.Errorf("Expected non-null LocalDate from dateEpochDay(0)")
	}

	// dateNow
	dateNowResult := chronologyDateNow([]any{isoChronologyObject})
	if object.IsNull(dateNowResult.(*object.Object)) {
		t.Errorf("Expected non-null LocalDate from dateNow()")
	}

	// dateYearDay
	dateYearDayResult := chronologyDateYearDay([]any{isoChronologyObject, int64(2026), int64(309)})
	if object.IsNull(dateYearDayResult.(*object.Object)) {
		t.Errorf("Expected non-null LocalDate from dateYearDay(2026, 309)")
	}

	// dateYearDay(Era, yearOfEra, dayOfYear)
	dateYearDayEraResult := chronologyDateYearDayEra([]any{isoChronologyObject, ceEraObject, int64(2026), int64(309)})
	if object.IsNull(dateYearDayEraResult.(*object.Object)) {
		t.Errorf("Expected non-null LocalDate from dateYearDay(CE, 2026, 309)")
	}

	// epochSecond
	epochSecondResult := chronologyEpochSecond([]any{
		isoChronologyObject,
		int64(2026),
		int64(11),
		int64(5),
		int64(12),
		int64(30),
		int64(0),
		object.Null,
	})
	if _, okSecond := epochSecondResult.(int64); !okSecond {
		t.Errorf("Expected int64 epochSecond result, got %T", epochSecondResult)
	}

	// epochSecond(Era, ...)
	epochSecondEraResult := chronologyEpochSecondEra([]any{
		isoChronologyObject,
		ceEraObject,
		int64(2026),
		int64(11),
		int64(5),
		int64(12),
		int64(30),
		int64(0),
		object.Null,
	})
	if _, okSecondEra := epochSecondEraResult.(int64); !okSecondEra {
		t.Errorf("Expected int64 epochSecondEra result, got %T", epochSecondEraResult)
	}

	// period
	periodResult := chronologyPeriod([]any{isoChronologyObject, int64(1), int64(2), int64(3)})
	if object.IsNull(periodResult.(*object.Object)) {
		t.Errorf("Expected non-null Period from period(1, 2, 3)")
	}

	// localDateTime
	localDateTimeResult := chronologyLocalDateTime([]any{isoChronologyObject, dateYmdObject})
	if object.IsNull(localDateTimeResult.(*object.Object)) {
		t.Errorf("Expected non-null LocalDateTime from localDateTime(TemporalAccessor)")
	}

	// zonedDateTime
	zonedDateTimeResult := chronologyZonedDateTimeTemporal([]any{isoChronologyObject, dateYmdObject})
	if object.IsNull(zonedDateTimeResult.(*object.Object)) {
		t.Errorf("Expected non-null ZonedDateTime from zonedDateTime(TemporalAccessor)")
	}

	// equals, hashCode, compareTo, toString
	if !isJavaBoolTrue(chronologyEquals([]any{isoChronologyObject, isoChronologyObject})) {
		t.Errorf("Expected IsoChronology to equal itself")
	}

	if chronologyCompareTo([]any{isoChronologyObject, isoChronologyObject}).(int64) != 0 {
		t.Errorf("Expected compareTo with self to be 0")
	}

	hashCodeResult := chronologyHashCode([]any{isoChronologyObject})
	if _, okHash := hashCodeResult.(int64); !okHash {
		t.Errorf("Expected int64 hashCode result, got %T", hashCodeResult)
	}

	toStringResult := chronologyToString([]any{isoChronologyObject})
	if object.GoStringFromStringObject(toStringResult.(*object.Object)) != "ISO" {
		t.Errorf("Expected toString 'ISO', got '%s'", object.GoStringFromStringObject(toStringResult.(*object.Object)))
	}
}

func getChronoFieldInstance(fieldName string) *object.Object {
	ensureChronoFieldInited()
	fieldVal := statics.GetStaticValue("java/time/temporal/ChronoField", fieldName)
	if fieldObject, okObject := fieldVal.(*object.Object); okObject && !object.IsNull(fieldObject) {
		return fieldObject
	}
	fallbackObject := object.MakeEmptyObjectWithClassName(&classNameChronoField)
	fallbackObject.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString(fieldName),
	}
	return fallbackObject
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestTemporalAccessor(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_TemporalAccessor()
	Load_Time_Temporal_ChronoField()

	chronoFieldClinit(nil)

	classNameTemporalAccessor := "java/time/temporal/TemporalAccessor"
	accessorObject := object.MakeEmptyObjectWithClassName(&classNameTemporalAccessor)
	accessorObject.FieldTable["DAY_OF_MONTH"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(25),
	}

	dayOfMonthName := object.StringObjectFromGoString("DAY_OF_MONTH")
	dayOfMonthField := chronoFieldValueOf([]any{dayOfMonthName}).(*object.Object)

	// isSupported
	isSupported := temporalAccessorIsSupported([]any{accessorObject, dayOfMonthField})
	if !isJavaBoolTrue(isSupported) {
		t.Errorf("Expected DAY_OF_MONTH to be supported on accessor")
	}

	// get and getLong
	getIntValue := temporalAccessorGet([]any{accessorObject, dayOfMonthField})
	if getIntValue.(int64) != 25 {
		t.Errorf("Expected get to return 25, got %v", getIntValue)
	}

	getLongValue := temporalAccessorGetLong([]any{accessorObject, dayOfMonthField})
	if getLongValue.(int64) != 25 {
		t.Errorf("Expected getLong to return 25, got %v", getLongValue)
	}

	// range
	rangeResult := temporalAccessorRange([]any{accessorObject, dayOfMonthField})
	if rangeResult == nil || object.IsNull(rangeResult.(*object.Object)) {
		t.Errorf("Expected non-null ValueRange for supported field")
	}

	// query
	zoneIdQuery := temporalQueriesZoneId(nil).(*object.Object)
	queryResult := temporalAccessorQuery([]any{accessorObject, zoneIdQuery})
	if queryResult != nil {
		t.Errorf("Expected query to return nil for unset zoneId, got %v", queryResult)
	}

	// Unsupported field error handling
	yearName := object.StringObjectFromGoString("YEAR")
	yearField := chronoFieldValueOf([]any{yearName}).(*object.Object)
	unsupportedResult := temporalAccessorGetLong([]any{accessorObject, yearField})
	if !isGErrBlk(unsupportedResult) {
		t.Errorf("Expected UnsupportedTemporalTypeException for unset field, got %v", unsupportedResult)
	}
}

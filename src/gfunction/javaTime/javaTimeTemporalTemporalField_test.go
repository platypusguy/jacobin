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

func TestTemporalField(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_TemporalField()
	Load_Time_Temporal_ChronoField()

	chronoFieldClinit(nil)

	dayOfMonthName := object.StringObjectFromGoString("DAY_OF_MONTH")
	fieldObject := chronoFieldValueOf([]any{dayOfMonthName}).(*object.Object)

	// isDateBased and isTimeBased
	if !isJavaBoolTrue(temporalFieldIsDateBased([]any{fieldObject})) {
		t.Errorf("Expected isDateBased to be true for DAY_OF_MONTH")
	}
	if isJavaBoolTrue(temporalFieldIsTimeBased([]any{fieldObject})) {
		t.Errorf("Expected isTimeBased to be false for DAY_OF_MONTH")
	}

	// range and getBaseUnit / getRangeUnit
	rangeObj := temporalFieldRange([]any{fieldObject})
	if rangeObj == nil || object.IsNull(rangeObj.(*object.Object)) {
		t.Errorf("Expected non-null ValueRange")
	}

	baseUnit := temporalFieldGetBaseUnit([]any{fieldObject})
	if baseUnit == nil || object.IsNull(baseUnit.(*object.Object)) {
		t.Errorf("Expected non-null baseUnit")
	}

	rangeUnit := temporalFieldGetRangeUnit([]any{fieldObject})
	if rangeUnit == nil || object.IsNull(rangeUnit.(*object.Object)) {
		t.Errorf("Expected non-null rangeUnit")
	}

	// toString
	strResult := temporalFieldToString([]any{fieldObject})
	if object.IsNull(strResult.(*object.Object)) {
		t.Errorf("Expected non-null String from toString")
	}

	// getFrom and adjustInto with mock Temporal
	classNameTemporal := "java/time/temporal/Temporal"
	temporalObject := object.MakeEmptyObjectWithClassName(&classNameTemporal)
	temporalObject.FieldTable["DAY_OF_MONTH"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(10),
	}

	fromVal := temporalFieldGetFrom([]any{fieldObject, temporalObject})
	if fromVal.(int64) != 10 {
		t.Errorf("Expected getFrom to return 10, got %v", fromVal)
	}

	adjusted := temporalFieldAdjustInto([]any{fieldObject, temporalObject, int64(20)})
	adjustedObj := adjusted.(*object.Object)
	if adjustedObj.FieldTable["DAY_OF_MONTH"].Fvalue.(int64) != 20 {
		t.Errorf("Expected DAY_OF_MONTH to be 20 after adjustInto, got %v", adjustedObj.FieldTable["DAY_OF_MONTH"].Fvalue)
	}
}

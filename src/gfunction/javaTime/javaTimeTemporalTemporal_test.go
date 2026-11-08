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

func TestTemporalAndAdjuster(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_Temporal()
	Load_Time_Temporal_TemporalAdjuster()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()

	chronoFieldClinit(nil)
	chronoUnitClinit(nil)

	classNameTemporal := "java/time/temporal/Temporal"
	temporalObject := object.MakeEmptyObjectWithClassName(&classNameTemporal)
	temporalObject.FieldTable["DAYS"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(10),
	}
	temporalObject.FieldTable["seconds"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(86400 * 10),
	}
	temporalObject.FieldTable["DAY_OF_MONTH"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(10),
	}

	dayOfMonthName := object.StringObjectFromGoString("DAY_OF_MONTH")
	dayOfMonthField := chronoFieldValueOf([]any{dayOfMonthName}).(*object.Object)

	daysName := object.StringObjectFromGoString("DAYS")
	daysUnit := chronoUnitValueOf([]any{daysName}).(*object.Object)

	// isSupported (unit)
	if !isJavaBoolTrue(temporalIsSupportedUnit([]any{temporalObject, daysUnit})) {
		t.Errorf("Expected DAYS to be supported on temporal")
	}

	// plus(long, TemporalUnit)
	plusResult := temporalPlusUnit([]any{temporalObject, int64(5), daysUnit})
	plusObject, okPlus := plusResult.(*object.Object)
	if !okPlus || object.IsNull(plusObject) {
		t.Fatalf("Expected non-null object from plusUnit, got %v", plusResult)
	}
	daysValue := plusObject.FieldTable["DAYS"].Fvalue.(int64)
	if daysValue != 15 {
		t.Errorf("Expected DAYS to be 15 after plus 5 days, got %d", daysValue)
	}

	// minus(long, TemporalUnit)
	minusResult := temporalMinusUnit([]any{temporalObject, int64(3), daysUnit})
	minusObject := minusResult.(*object.Object)
	minusDaysValue := minusObject.FieldTable["DAYS"].Fvalue.(int64)
	if minusDaysValue != 7 {
		t.Errorf("Expected DAYS to be 7 after minus 3 days, got %d", minusDaysValue)
	}

	// with(TemporalField, long)
	withResult := temporalWithField([]any{temporalObject, dayOfMonthField, int64(22)})
	withObject := withResult.(*object.Object)
	withDayValue := withObject.FieldTable["DAY_OF_MONTH"].Fvalue.(int64)
	if withDayValue != 22 {
		t.Errorf("Expected DAY_OF_MONTH to be 22 after withField, got %d", withDayValue)
	}

	// until
	endTemporal := object.MakeEmptyObjectWithClassName(&classNameTemporal)
	endTemporal.FieldTable["DAYS"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(25),
	}
	endTemporal.FieldTable["seconds"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(86400 * 25),
	}
	untilResult := temporalUntil([]any{temporalObject, endTemporal, daysUnit})
	if untilResult.(int64) != 15 {
		t.Errorf("Expected until to return 15 days, got %v", untilResult)
	}

	// adjustInto
	adjustedResult := temporalAdjusterAdjustInto([]any{temporalObject, endTemporal})
	if adjustedResult == nil || object.IsNull(adjustedResult.(*object.Object)) {
		t.Errorf("Expected non-null result from adjustInto, got %v", adjustedResult)
	}
}

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

func TestTemporalUnit(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_TemporalUnit()
	Load_Time_Temporal_ChronoUnit()

	chronoUnitClinit(nil)

	daysName := object.StringObjectFromGoString("DAYS")
	unitObject := chronoUnitValueOf([]any{daysName}).(*object.Object)

	// isDateBased and isTimeBased
	if !isJavaBoolTrue(temporalUnitIsDateBased([]any{unitObject})) {
		t.Errorf("Expected isDateBased to be true for DAYS")
	}
	if isJavaBoolTrue(temporalUnitIsTimeBased([]any{unitObject})) {
		t.Errorf("Expected isTimeBased to be false for DAYS")
	}

	// getDuration
	durationObj := temporalUnitGetDuration([]any{unitObject})
	if durationObj == nil || object.IsNull(durationObj.(*object.Object)) {
		t.Errorf("Expected non-null Duration for DAYS")
	}

	// toString
	strResult := temporalUnitToString([]any{unitObject})
	if object.IsNull(strResult.(*object.Object)) {
		t.Errorf("Expected non-null String from toString")
	}

	// isSupportedBy, addTo, and between with mock Temporal
	classNameTemporal := "java/time/temporal/Temporal"
	temporalStart := object.MakeEmptyObjectWithClassName(&classNameTemporal)
	temporalStart.FieldTable["DAYS"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(5),
	}
	temporalStart.FieldTable["seconds"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(86400 * 5),
	}

	temporalEnd := object.MakeEmptyObjectWithClassName(&classNameTemporal)
	temporalEnd.FieldTable["DAYS"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(15),
	}
	temporalEnd.FieldTable["seconds"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(86400 * 15),
	}

	if !isJavaBoolTrue(temporalUnitIsSupportedBy([]any{unitObject, temporalStart})) {
		t.Errorf("Expected DAYS to be supported by temporal")
	}

	added := temporalUnitAddTo([]any{unitObject, temporalStart, int64(10)})
	addedObj := added.(*object.Object)
	if addedObj.FieldTable["seconds"].Fvalue.(int64) != 86400*15 {
		t.Errorf("Expected seconds to be %d after addTo, got %v", 86400*15, addedObj.FieldTable["seconds"].Fvalue)
	}

	betweenVal := temporalUnitBetween([]any{unitObject, temporalStart, temporalEnd})
	if betweenVal.(int64) != 10 {
		t.Errorf("Expected between to return 10, got %v", betweenVal)
	}
}

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

func TestTemporalAmount(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_TemporalAmount()
	Load_Time_Temporal_ChronoUnit()

	chronoUnitClinit(nil)

	classNameTemporalAmount := "java/time/temporal/TemporalAmount"
	amountObject := object.MakeEmptyObjectWithClassName(&classNameTemporalAmount)
	amountObject.FieldTable["seconds"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(100),
	}

	secondsName := object.StringObjectFromGoString("SECONDS")
	secondsUnit := chronoUnitValueOf([]any{secondsName}).(*object.Object)

	// get
	getValue := temporalAmountGet([]any{amountObject, secondsUnit})
	if getValue.(int64) != 100 {
		t.Errorf("Expected get SECONDS to return 100, got %v", getValue)
	}

	// getUnits
	unitsList := temporalAmountGetUnits([]any{amountObject})
	if unitsList == nil || object.IsNull(unitsList.(*object.Object)) {
		t.Errorf("Expected non-null List for getUnits")
	}

	// addTo and subtractFrom with target Temporal
	classNameTemporal := "java/time/temporal/Temporal"
	targetTemporal := object.MakeEmptyObjectWithClassName(&classNameTemporal)
	targetTemporal.FieldTable["seconds"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(500),
	}

	addResult := temporalAmountAddTo([]any{amountObject, targetTemporal})
	addedObject, okAdded := addResult.(*object.Object)
	if !okAdded || object.IsNull(addedObject) {
		t.Fatalf("Expected non-null object from addTo, got %v", addResult)
	}
	if addedObject.FieldTable["seconds"].Fvalue.(int64) != 600 {
		t.Errorf("Expected 600 seconds after addTo, got %v", addedObject.FieldTable["seconds"].Fvalue)
	}

	subResult := temporalAmountSubtractFrom([]any{amountObject, targetTemporal})
	subbedObject := subResult.(*object.Object)
	if subbedObject.FieldTable["seconds"].Fvalue.(int64) != 400 {
		t.Errorf("Expected 400 seconds after subtractFrom, got %v", subbedObject.FieldTable["seconds"].Fvalue)
	}
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func TestTemporalQueriesAndQuery(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_TemporalQuery()
	Load_Time_Temporal_TemporalQueries()

	temporalQueriesClinit(nil)

	// Test query factories
	zoneIdQuery := temporalQueriesZoneId(nil)
	chronologyQuery := temporalQueriesChronology(nil)
	precisionQuery := temporalQueriesPrecision(nil)
	zoneQuery := temporalQueriesZone(nil)
	offsetQuery := temporalQueriesOffset(nil)
	localDateQuery := temporalQueriesLocalDate(nil)
	localTimeQuery := temporalQueriesLocalTime(nil)

	if object.IsNull(zoneIdQuery.(*object.Object)) {
		t.Errorf("Expected non-null zoneIdQuery")
	}
	if object.IsNull(chronologyQuery.(*object.Object)) {
		t.Errorf("Expected non-null chronologyQuery")
	}
	if object.IsNull(precisionQuery.(*object.Object)) {
		t.Errorf("Expected non-null precisionQuery")
	}
	if object.IsNull(zoneQuery.(*object.Object)) {
		t.Errorf("Expected non-null zoneQuery")
	}
	if object.IsNull(offsetQuery.(*object.Object)) {
		t.Errorf("Expected non-null offsetQuery")
	}
	if object.IsNull(localDateQuery.(*object.Object)) {
		t.Errorf("Expected non-null localDateQuery")
	}
	if object.IsNull(localTimeQuery.(*object.Object)) {
		t.Errorf("Expected non-null localTimeQuery")
	}

	// Test queryFrom with mock TemporalAccessor
	classNameTemporalAccessor := "java/time/temporal/TemporalAccessor"
	mockAccessor := object.MakeEmptyObjectWithClassName(&classNameTemporalAccessor)

	classNameZoneId := "java/time/ZoneId"
	mockZoneId := object.MakeEmptyObjectWithClassName(&classNameZoneId)
	mockAccessor.FieldTable["zoneId"] = object.Field{
		Ftype:  "Ljava/time/ZoneId;",
		Fvalue: mockZoneId,
	}

	classNameChronology := "java/time/chrono/Chronology"
	mockChronology := object.MakeEmptyObjectWithClassName(&classNameChronology)
	mockAccessor.FieldTable["chronology"] = object.Field{
		Ftype:  "Ljava/time/chrono/Chronology;",
		Fvalue: mockChronology,
	}

	// Query zoneId
	queriedZoneId := temporalQueryQueryFrom([]any{zoneIdQuery, mockAccessor})
	if queriedZoneId != mockZoneId {
		t.Errorf("Expected queriedZoneId to be mockZoneId, got %v", queriedZoneId)
	}

	// Query chronology
	queriedChronology := temporalQueryQueryFrom([]any{chronologyQuery, mockAccessor})
	if queriedChronology != mockChronology {
		t.Errorf("Expected queriedChronology to be mockChronology, got %v", queriedChronology)
	}

	// Query not present should return nil
	queriedPrecision := temporalQueryQueryFrom([]any{precisionQuery, mockAccessor})
	if queriedPrecision != nil {
		t.Errorf("Expected nil when field not present, got %v", queriedPrecision)
	}

	// Null query or null temporal should return error
	if !isGErrBlk(temporalQueryQueryFrom([]any{object.Null, mockAccessor})) {
		t.Errorf("Expected error for null query")
	}
	if !isGErrBlk(temporalQueryQueryFrom([]any{zoneIdQuery, object.Null})) {
		t.Errorf("Expected error for null temporal")
	}
}

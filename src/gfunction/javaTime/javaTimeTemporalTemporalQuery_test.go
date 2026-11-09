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
	"testing"
)

func TestTemporalQuery(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_TemporalQuery()
	Load_Time_Temporal_TemporalQueries()
	Load_Time_LocalDate()

	// Null checks
	err := temporalQueryQueryFrom(nil)
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("queryFrom(nil) expected NullPointerException, got %v", err)
	}

	err = temporalQueryQueryFrom([]any{nil, nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("queryFrom(nil, nil) expected NullPointerException, got %v", err)
	}

	queryChronology := temporalQueriesChronology(nil).(*object.Object)
	err = temporalQueryQueryFrom([]any{queryChronology, nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("queryFrom(query, nil) expected NullPointerException, got %v", err)
	}

	// Valid query
	date := CreateLocalDate(2026, 11, 8)
	res := temporalQueryQueryFrom([]any{queryChronology, date})
	if resObj, ok := res.(*object.Object); !ok || resObj == nil {
		t.Errorf("queryFrom(chronology, date) failed: %v", res)
	}
}

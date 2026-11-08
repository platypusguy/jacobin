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

func TestChronoZonedDateTime(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Chrono_Chronology()
	Load_Time_Chrono_ChronoLocalDate()
	Load_Time_Chrono_ChronoLocalDateTime()
	Load_Time_Chrono_ChronoZonedDateTime()

	isoName := object.StringObjectFromGoString("ISO")
	chronoObj := chronologyOf([]any{isoName}).(*object.Object)

	dateObj := chronologyDateYMD([]any{chronoObj, int64(2026), int64(11), int64(5)}).(*object.Object)

	classNameLocalTime := "java/time/LocalTime"
	timeObj := object.MakeEmptyObjectWithClassName(&classNameLocalTime)
	timeObj.FieldTable["hour"] = object.Field{Fvalue: int64(10)}
	timeObj.FieldTable["minute"] = object.Field{Fvalue: int64(30)}
	timeObj.FieldTable["second"] = object.Field{Fvalue: int64(0)}
	timeObj.FieldTable["nano"] = object.Field{Fvalue: int64(0)}

	dateTimeObj := chronoLocalDateAtTime([]any{dateObj, timeObj}).(*object.Object)

	classNameZoneId := "java/time/ZoneId"
	zoneIdObj := object.MakeEmptyObjectWithClassName(&classNameZoneId)
	zoneIdObj.FieldTable["id"] = object.Field{Fvalue: object.StringObjectFromGoString("UTC")}

	zonedDateTimeObj := chronoLocalDateTimeAtZone([]any{dateTimeObj, zoneIdObj}).(*object.Object)

	// getZone
	zoneRetrieved := chronoZonedDateTimeGetZone([]any{zonedDateTimeObj})
	if object.IsNull(zoneRetrieved.(*object.Object)) {
		t.Errorf("Expected non-null ZoneId from getZone")
	}

	// toLocalDate, toLocalTime, toLocalDateTime
	if object.IsNull(chronoZonedDateTimeToLocalDate([]any{zonedDateTimeObj}).(*object.Object)) {
		t.Errorf("Expected non-null LocalDate")
	}
	if object.IsNull(chronoZonedDateTimeToLocalTime([]any{zonedDateTimeObj}).(*object.Object)) {
		t.Errorf("Expected non-null LocalTime")
	}
	if object.IsNull(chronoZonedDateTimeToLocalDateTime([]any{zonedDateTimeObj}).(*object.Object)) {
		t.Errorf("Expected non-null LocalDateTime")
	}

	// toEpochSecond
	epochSec := chronoZonedDateTimeToEpochSecond([]any{zonedDateTimeObj})
	if _, okSec := epochSec.(int64); !okSec {
		t.Errorf("Expected int64 from toEpochSecond, got %T", epochSec)
	}

	// toString
	strResult := chronoZonedDateTimeToString([]any{zonedDateTimeObj})
	if strResult == nil || object.IsNull(strResult.(*object.Object)) {
		t.Errorf("Expected non-null String from toString")
	}
}

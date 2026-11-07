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

func TestChronoLocalDateTime(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Chrono_Chronology()
	Load_Time_Chrono_ChronoLocalDate()
	Load_Time_Chrono_ChronoLocalDateTime()
	Load_Time_Temporal_ChronoUnit()

	chronoUnitClinit(nil)

	isoName := object.StringObjectFromGoString("ISO")
	chronoObj := chronologyOf([]any{isoName}).(*object.Object)

	dateObj := chronologyDateYMD([]any{chronoObj, int64(2026), int64(11), int64(5)}).(*object.Object)

	classNameLocalTime := "java/time/LocalTime"
	timeObj1 := object.MakeEmptyObjectWithClassName(&classNameLocalTime)
	timeObj1.FieldTable["hour"] = object.Field{Fvalue: int64(10)}
	timeObj1.FieldTable["minute"] = object.Field{Fvalue: int64(30)}
	timeObj1.FieldTable["second"] = object.Field{Fvalue: int64(0)}
	timeObj1.FieldTable["nano"] = object.Field{Fvalue: int64(0)}

	timeObj2 := object.MakeEmptyObjectWithClassName(&classNameLocalTime)
	timeObj2.FieldTable["hour"] = object.Field{Fvalue: int64(14)}
	timeObj2.FieldTable["minute"] = object.Field{Fvalue: int64(45)}
	timeObj2.FieldTable["second"] = object.Field{Fvalue: int64(0)}
	timeObj2.FieldTable["nano"] = object.Field{Fvalue: int64(0)}

	dateTimeObj1 := chronoLocalDateAtTime([]any{dateObj, timeObj1}).(*object.Object)
	dateTimeObj2 := chronoLocalDateAtTime([]any{dateObj, timeObj2}).(*object.Object)

	// isBefore, isAfter, isEqual
	if !isJavaBoolTrue(chronoLocalDateTimeIsBefore([]any{dateTimeObj1, dateTimeObj2})) {
		t.Errorf("Expected dateTime1 to be before dateTime2")
	}
	if !isJavaBoolTrue(chronoLocalDateTimeIsAfter([]any{dateTimeObj2, dateTimeObj1})) {
		t.Errorf("Expected dateTime2 to be after dateTime1")
	}
	if !isJavaBoolTrue(chronoLocalDateTimeIsEqual([]any{dateTimeObj1, dateTimeObj1})) {
		t.Errorf("Expected dateTime1 to equal itself")
	}

	// toLocalDate and toLocalTime
	retrievedDate := chronoLocalDateTimeToLocalDate([]any{dateTimeObj1})
	if object.IsNull(retrievedDate.(*object.Object)) {
		t.Errorf("Expected non-null LocalDate from toLocalDate")
	}

	retrievedTime := chronoLocalDateTimeToLocalTime([]any{dateTimeObj1})
	if object.IsNull(retrievedTime.(*object.Object)) {
		t.Errorf("Expected non-null LocalTime from toLocalTime")
	}

	// plus and minus
	hoursName := object.StringObjectFromGoString("HOURS")
	hoursUnit := chronoUnitValueOf([]any{hoursName}).(*object.Object)

	plusDT := chronoLocalDateTimePlusUnit([]any{dateTimeObj1, int64(2), hoursUnit})
	if object.IsNull(plusDT.(*object.Object)) {
		t.Errorf("Expected non-null ChronoLocalDateTime from plusUnit")
	}

	// toString
	strResult := chronoLocalDateTimeToString([]any{dateTimeObj1})
	if object.GoStringFromStringObject(strResult.(*object.Object)) != "2026-11-05T10:30:00" {
		t.Errorf("Expected '2026-11-05T10:30:00', got '%s'", object.GoStringFromStringObject(strResult.(*object.Object)))
	}
}

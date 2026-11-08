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
	"jacobin/src/statics"
	"testing"
)

func dateTimeFormatterTestInit() {
	globals.InitGlobals("test")
	Load_Time_Format_DateTimeFormatter()
	dateTimeFormatterClinit(nil)
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_LocalDateTime()
	Load_Time_Instant()
	Load_Time_ZoneOffset()
}

func TestDateTimeFormatterConstantsAndFactories(t *testing.T) {
	dateTimeFormatterTestInit()

	expectedConstants := []string{
		"ISO_LOCAL_DATE", "ISO_OFFSET_DATE", "ISO_DATE",
		"ISO_LOCAL_TIME", "ISO_OFFSET_TIME", "ISO_TIME",
		"ISO_LOCAL_DATE_TIME", "ISO_OFFSET_DATE_TIME", "ISO_ZONED_DATE_TIME", "ISO_DATE_TIME",
		"ISO_ORDINAL_DATE", "ISO_WEEK_DATE", "ISO_INSTANT", "BASIC_ISO_DATE", "RFC_1123_DATE_TIME",
	}

	for _, constName := range expectedConstants {
		st, ok := statics.QueryStatic(classNameDateTimeFormatter, constName)
		if !ok || st.Value == nil {
			t.Errorf("DateTimeFormatter.%s not found", constName)
		}
	}

	// ofPattern(String)
	patStr := object.StringObjectFromGoString("yyyy/MM/dd")
	fmtObj := dateTimeFormatterOfPatternString([]any{patStr}).(*object.Object)
	if fmtObj == nil {
		t.Fatalf("ofPattern('yyyy/MM/dd') returned nil")
	}

	// Null check
	err := dateTimeFormatterOfPatternString([]any{nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("ofPattern(nil) expected NullPointerException, got %v", err)
	}
}

func TestDateTimeFormatterFormatting(t *testing.T) {
	dateTimeFormatterTestInit()

	date := CreateLocalDate(2026, 11, 8)
	timeVal := CreateLocalTime(14, 30, 45, 0)
	dt := CreateLocalDateTime(date, timeVal)

	// ISO_LOCAL_DATE formatting
	isoDateFmt := statics.GetStaticValue(classNameDateTimeFormatter, "ISO_LOCAL_DATE").(*object.Object)
	res := dateTimeFormatterFormat([]any{isoDateFmt, date})
	if strObj, ok := res.(*object.Object); !ok || object.GoStringFromStringObject(strObj) != "2026-11-08" {
		t.Errorf("ISO_LOCAL_DATE.format(2026-11-08) = %v, want '2026-11-08'", res)
	}

	// ISO_LOCAL_TIME formatting
	isoTimeFmt := statics.GetStaticValue(classNameDateTimeFormatter, "ISO_LOCAL_TIME").(*object.Object)
	res = dateTimeFormatterFormat([]any{isoTimeFmt, timeVal})
	if strObj, ok := res.(*object.Object); !ok || object.GoStringFromStringObject(strObj) != "14:30:45" {
		t.Errorf("ISO_LOCAL_TIME.format(14:30:45) = %v, want '14:30:45'", res)
	}

	// ISO_LOCAL_DATE_TIME formatting
	isoLdtFmt := statics.GetStaticValue(classNameDateTimeFormatter, "ISO_LOCAL_DATE_TIME").(*object.Object)
	res = dateTimeFormatterFormat([]any{isoLdtFmt, dt})
	if strObj, ok := res.(*object.Object); !ok || object.GoStringFromStringObject(strObj) != "2026-11-08T14:30:45" {
		t.Errorf("ISO_LOCAL_DATE_TIME.format = %v, want '2026-11-08T14:30:45'", res)
	}

	// BASIC_ISO_DATE formatting
	basicFmt := statics.GetStaticValue(classNameDateTimeFormatter, "BASIC_ISO_DATE").(*object.Object)
	res = dateTimeFormatterFormat([]any{basicFmt, date})
	if strObj, ok := res.(*object.Object); !ok || object.GoStringFromStringObject(strObj) != "20261108" {
		t.Errorf("BASIC_ISO_DATE.format(2026-11-08) = %v, want '20261108'", res)
	}

	// Custom pattern formatting
	patFmt := dateTimeFormatterOfPatternString([]any{object.StringObjectFromGoString("yyyy/MM/dd")}).(*object.Object)
	res = dateTimeFormatterFormat([]any{patFmt, date})
	if strObj, ok := res.(*object.Object); !ok || object.GoStringFromStringObject(strObj) != "2026/11/08" {
		t.Errorf("yyyy/MM/dd format = %v, want '2026/11/08'", res)
	}
}

func TestDateTimeFormatterParsing(t *testing.T) {
	dateTimeFormatterTestInit()

	isoDateFmt := statics.GetStaticValue(classNameDateTimeFormatter, "ISO_LOCAL_DATE").(*object.Object)
	strToParse := object.StringObjectFromGoString("2026-11-08")
	res := dateTimeFormatterParseCharSequence([]any{isoDateFmt, strToParse})
	if resObj, ok := res.(*object.Object); !ok || resObj == nil {
		t.Fatalf("ISO_LOCAL_DATE.parse('2026-11-08') failed: %v", res)
	}

	// Parse invalid string
	badStr := object.StringObjectFromGoString("not-a-date")
	err := dateTimeFormatterParseCharSequence([]any{isoDateFmt, badStr})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeParseException {
		t.Errorf("parse('not-a-date') expected DateTimeParseException, got %v", err)
	}
}

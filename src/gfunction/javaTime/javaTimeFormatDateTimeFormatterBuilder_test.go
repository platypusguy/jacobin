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

func dateTimeFormatterBuilderTestInit() {
	globals.InitGlobals("test")
	Load_Time_Format_DateTimeFormatter()
	Load_Time_Format_DateTimeFormatterBuilder()
	Load_Time_LocalDate()
}

func TestDateTimeFormatterBuilder(t *testing.T) {
	dateTimeFormatterBuilderTestInit()

	builderObj := object.MakeEmptyObjectWithClassName(&classNameDateTimeFormatterBuilder)
	dateTimeFormatterBuilderInit([]any{builderObj})

	// appendPattern
	patStr := object.StringObjectFromGoString("yyyy-MM-dd")
	res := dateTimeFormatterBuilderAppendPattern([]any{builderObj, patStr})
	if res != builderObj {
		t.Errorf("appendPattern did not return self")
	}

	// appendLiteral
	litStr := object.StringObjectFromGoString(" ")
	res = dateTimeFormatterBuilderAppendLiteralString([]any{builderObj, litStr})
	if res != builderObj {
		t.Errorf("appendLiteralString did not return self")
	}

	// parseCaseInsensitive / parseCaseSensitive / parseLenient / parseStrict
	_ = dateTimeFormatterBuilderParseCaseInsensitive([]any{builderObj})
	_ = dateTimeFormatterBuilderParseCaseSensitive([]any{builderObj})
	_ = dateTimeFormatterBuilderParseLenient([]any{builderObj})
	_ = dateTimeFormatterBuilderParseStrict([]any{builderObj})

	// optionalStart & optionalEnd
	_ = dateTimeFormatterBuilderOptionalStart([]any{builderObj})
	_ = dateTimeFormatterBuilderAppendPattern([]any{builderObj, object.StringObjectFromGoString("HH:mm:ss")})
	_ = dateTimeFormatterBuilderOptionalEnd([]any{builderObj})

	// toFormatter
	fmtRes := dateTimeFormatterBuilderToFormatterNoArgs([]any{builderObj})
	fmtObj, ok := fmtRes.(*object.Object)
	if !ok || fmtObj == nil {
		t.Fatalf("toFormatter returned non-object: %v", fmtRes)
	}

	// Format with built formatter
	date := CreateLocalDate(2026, 11, 8)
	formatRes := dateTimeFormatterFormat([]any{fmtObj, date})
	if strObj, okStr := formatRes.(*object.Object); !okStr || object.GoStringFromStringObject(strObj) == "" {
		t.Errorf("formatting with built formatter failed: %v", formatRes)
	}
}

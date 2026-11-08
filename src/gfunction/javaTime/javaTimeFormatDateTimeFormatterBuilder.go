/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
)

var classNameDateTimeFormatterBuilder = "java/time/format/DateTimeFormatterBuilder"
var classNameDateTimeFormatter = "java/time/format/DateTimeFormatter"

var dtFormatterBuilderMutex sync.Mutex

// Load_Time_Format_DateTimeFormatterBuilder registers all java/time/format/DateTimeFormatterBuilder methods.
func Load_Time_Format_DateTimeFormatterBuilder() {
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderInit,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.append(Ljava/time/format/DateTimeFormatter;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendFormatter,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendChronologyId()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderAppendChronologyId,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendChronologyText(Ljava/time/format/TextStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendChronologyText,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendFraction(Ljava/time/temporal/TemporalField;IIZ)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  dateTimeFormatterBuilderAppendFraction,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendGenericZoneText(Ljava/time/format/TextStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendGenericZoneText,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendGenericZoneText(Ljava/time/format/TextStyle;Ljava/util/Set;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendGenericZoneTextWithSet,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendInstant()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderAppendInstant,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendInstant(I)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendInstantFractionalDigits,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendLiteral(C)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendLiteralChar,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendLiteral(Ljava/lang/String;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendLiteralString,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendLocalized(Ljava/time/format/FormatStyle;Ljava/time/format/FormatStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendLocalizedStyles,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendLocalizedOffset(Ljava/time/format/TextStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendLocalizedOffset,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendOffset(Ljava/lang/String;Ljava/lang/String;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendOffsetPattern,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendOffsetId()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderAppendOffsetId,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendOptional(Ljava/time/format/DateTimeFormatter;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendOptional,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendPattern(Ljava/lang/String;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendPattern,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendText(Ljava/time/temporal/TemporalField;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendTextShort,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendText(Ljava/time/temporal/TemporalField;Ljava/time/format/TextStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendTextStyle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendText(Ljava/time/temporal/TemporalField;Ljava/util/Map;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendTextMap,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendValue(Ljava/time/temporal/TemporalField;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendValueShort,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendValue(Ljava/time/temporal/TemporalField;I)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendValueWidth,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendValue(Ljava/time/temporal/TemporalField;IILjava/time/format/SignStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  dateTimeFormatterBuilderAppendValueWidthSign,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendValueReduced(Ljava/time/temporal/TemporalField;IIILjava/time/chrono/ChronoLocalDate;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  dateTimeFormatterBuilderAppendValueReducedChronoDate,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendValueReduced(Ljava/time/temporal/TemporalField;III)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  dateTimeFormatterBuilderAppendValueReducedBaseValue,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendValueReduced(Ljava/time/temporal/TemporalField;IILjava/time/chrono/ChronoLocalDate;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  dateTimeFormatterBuilderAppendValueReducedDate,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendZoneId()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderAppendZoneId,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendZoneOrOffsetId()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderAppendZoneOrOffsetId,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendZoneRegionId()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderAppendZoneRegionId,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendZoneText(Ljava/time/format/TextStyle;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderAppendZoneTextStyle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.appendZoneText(Ljava/time/format/TextStyle;Ljava/util/Set;)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderAppendZoneTextSet,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.getLocalizedDateTimePattern(Ljava/time/format/FormatStyle;Ljava/time/format/FormatStyle;Ljava/time/chrono/Chronology;Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  dateTimeFormatterBuilderGetLocalizedDateTimePattern,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.optionalEnd()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderOptionalEnd,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.optionalStart()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderOptionalStart,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.padNext(I)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderPadNextPadWidth,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.padNext(IC)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderPadNextPadWidthPadChar,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.parseCaseInsensitive()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderParseCaseInsensitive,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.parseCaseSensitive()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderParseCaseSensitive,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.parseDefaulting(Ljava/time/temporal/TemporalField;J)Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderParseDefaulting,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.parseLenient()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderParseLenient,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.parseStrict()Ljava/time/format/DateTimeFormatterBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderParseStrict,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.toFormatter()Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterBuilderToFormatterNoArgs,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.toFormatter(Ljava/util/Locale;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterBuilderToFormatterLocale,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatterBuilder.toFormatter(Ljava/time/format/ResolverStyle;Ljava/time/chrono/Chronology;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterBuilderToFormatterResolverChronology,
		}
}

func initBuilderFields(obj *object.Object) {
	if obj.FieldTable == nil {
		obj.FieldTable = make(map[string]object.Field)
	}
	if _, ok := obj.FieldTable["pattern"]; !ok {
		obj.FieldTable["pattern"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(""),
		}
	}
}

func getBuilderPattern(obj *object.Object) string {
	if obj == nil || obj.FieldTable == nil {
		return ""
	}
	if f, ok := obj.FieldTable["pattern"]; ok {
		if strObj, okStr := f.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
			return object.GoStringFromStringObject(strObj)
		}
	}
	return ""
}

func appendBuilderPattern(obj *object.Object, part string) {
	initBuilderFields(obj)
	curr := getBuilderPattern(obj)
	obj.FieldTable["pattern"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString(curr + part),
	}
}

// <init>()
func dateTimeFormatterBuilderInit(params []any) any {
	if len(params) > 0 && params[0] != nil {
		if self, ok := params[0].(*object.Object); ok && !object.IsNull(self) {
			initBuilderFields(self)
		}
	}
	return nil
}

// append(DateTimeFormatter formatter)
func dateTimeFormatterBuilderAppendFormatter(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.append: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.append: formatter is null")
	}
	self := params[0].(*object.Object)
	formatter := params[1].(*object.Object)
	pat := ""
	if formatter.FieldTable != nil {
		if f, ok := formatter.FieldTable["pattern"]; ok {
			if strObj, okStr := f.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				pat = object.GoStringFromStringObject(strObj)
			}
		}
	}
	appendBuilderPattern(self, pat)
	return self
}

// appendChronologyId()
func dateTimeFormatterBuilderAppendChronologyId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendChronologyId: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "G")
	return self
}

// appendChronologyText(TextStyle textStyle)
func dateTimeFormatterBuilderAppendChronologyText(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendChronologyText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendChronologyText: textStyle is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "G")
	return self
}

// appendFraction(TemporalField field, int minWidth, int maxWidth, boolean decimalPoint)
func dateTimeFormatterBuilderAppendFraction(params []any) any {
	if len(params) < 5 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendFraction: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendFraction: field is null")
	}
	self := params[0].(*object.Object)
	minWidth, okMin := getInt64Val(params[2])
	maxWidth, okMax := getInt64Val(params[3])
	if !okMin || !okMax {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendFraction: invalid widths")
	}
	if minWidth < 0 || minWidth > 9 || maxWidth < 1 || maxWidth > 9 || maxWidth < minWidth {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendFraction: invalid min/max width")
	}
	appendBuilderPattern(self, "S")
	return self
}

// appendGenericZoneText(TextStyle textStyle)
func dateTimeFormatterBuilderAppendGenericZoneText(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendGenericZoneText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendGenericZoneText: textStyle is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "vv")
	return self
}

// appendGenericZoneText(TextStyle textStyle, Set preferredZones)
func dateTimeFormatterBuilderAppendGenericZoneTextWithSet(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendGenericZoneText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendGenericZoneText: textStyle is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "vv")
	return self
}

// appendInstant()
func dateTimeFormatterBuilderAppendInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendInstant: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "yyyy-MM-dd'T'HH:mm:ss.SSSX")
	return self
}

// appendInstant(int fractionalDigits)
func dateTimeFormatterBuilderAppendInstantFractionalDigits(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendInstant: self is null")
	}
	digits, ok := getInt64Val(params[1])
	if !ok || digits < -1 || digits > 9 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendInstant: invalid fractional digits")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "yyyy-MM-dd'T'HH:mm:ssX")
	return self
}

// appendLiteral(char literal)
func dateTimeFormatterBuilderAppendLiteralChar(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendLiteral: self is null")
	}
	chVal, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendLiteral: invalid char")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, string(rune(chVal)))
	return self
}

// appendLiteral(String literal)
func dateTimeFormatterBuilderAppendLiteralString(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendLiteral: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendLiteral: literal is null")
	}
	strObj := params[1].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendLiteral: expected String")
	}
	lit := object.GoStringFromStringObject(strObj)
	self := params[0].(*object.Object)
	appendBuilderPattern(self, lit)
	return self
}

// appendLocalized(FormatStyle dateStyle, FormatStyle timeStyle)
func dateTimeFormatterBuilderAppendLocalizedStyles(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendLocalized: self is null")
	}
	dateStyleNull := params[1] == nil || object.IsNull(params[1].(*object.Object))
	timeStyleNull := params[2] == nil || object.IsNull(params[2].(*object.Object))
	if dateStyleNull && timeStyleNull {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendLocalized: both styles cannot be null")
	}
	self := params[0].(*object.Object)
	if !dateStyleNull && !timeStyleNull {
		appendBuilderPattern(self, "yyyy-MM-dd HH:mm:ss")
	} else if !dateStyleNull {
		appendBuilderPattern(self, "yyyy-MM-dd")
	} else {
		appendBuilderPattern(self, "HH:mm:ss")
	}
	return self
}

// appendLocalizedOffset(TextStyle style)
func dateTimeFormatterBuilderAppendLocalizedOffset(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendLocalizedOffset: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendLocalizedOffset: style is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "ZZZZ")
	return self
}

// appendOffset(String pattern, String noOffsetText)
func dateTimeFormatterBuilderAppendOffsetPattern(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendOffset: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendOffset: pattern is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendOffset: noOffsetText is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "XXX")
	return self
}

// appendOffsetId()
func dateTimeFormatterBuilderAppendOffsetId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendOffsetId: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "XXX")
	return self
}

// appendOptional(DateTimeFormatter formatter)
func dateTimeFormatterBuilderAppendOptional(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendOptional: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendOptional: formatter is null")
	}
	self := params[0].(*object.Object)
	formatter := params[1].(*object.Object)
	pat := ""
	if formatter.FieldTable != nil {
		if f, ok := formatter.FieldTable["pattern"]; ok {
			if strObj, okStr := f.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				pat = object.GoStringFromStringObject(strObj)
			}
		}
	}
	appendBuilderPattern(self, "["+pat+"]")
	return self
}

// appendPattern(String pattern)
func dateTimeFormatterBuilderAppendPattern(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendPattern: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendPattern: pattern is null")
	}
	strObj := params[1].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendPattern: expected String")
	}
	pat := object.GoStringFromStringObject(strObj)
	self := params[0].(*object.Object)
	appendBuilderPattern(self, pat)
	return self
}

// appendText(TemporalField field)
func dateTimeFormatterBuilderAppendTextShort(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: field is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "a")
	return self
}

// appendText(TemporalField field, TextStyle textStyle)
func dateTimeFormatterBuilderAppendTextStyle(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: field is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: textStyle is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "MMMM")
	return self
}

// appendText(TemporalField field, Map textLookup)
func dateTimeFormatterBuilderAppendTextMap(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: field is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendText: textLookup is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "text")
	return self
}

// appendValue(TemporalField field)
func dateTimeFormatterBuilderAppendValueShort(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: field is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "v")
	return self
}

// appendValue(TemporalField field, int width)
func dateTimeFormatterBuilderAppendValueWidth(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: field is null")
	}
	width, ok := getInt64Val(params[2])
	if !ok || width < 1 || width > 19 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendValue: invalid width")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "v")
	return self
}

// appendValue(TemporalField field, int minWidth, int maxWidth, SignStyle signStyle)
func dateTimeFormatterBuilderAppendValueWidthSign(params []any) any {
	if len(params) < 5 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: field is null")
	}
	if params[4] == nil || object.IsNull(params[4].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValue: signStyle is null")
	}
	minWidth, okMin := getInt64Val(params[2])
	maxWidth, okMax := getInt64Val(params[3])
	if !okMin || !okMax || minWidth < 1 || minWidth > 19 || maxWidth < 1 || maxWidth > 19 || maxWidth < minWidth {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendValue: invalid min/max width")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "v")
	return self
}

// appendValueReduced(TemporalField field, int width, int maxWidth, int baseValue)
func dateTimeFormatterBuilderAppendValueReducedBaseValue(params []any) any {
	if len(params) < 5 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: field is null")
	}
	width, okW := getInt64Val(params[2])
	maxWidth, okM := getInt64Val(params[3])
	if !okW || !okM || width < 1 || width > 10 || maxWidth < 1 || maxWidth > 10 || maxWidth < width {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendValueReduced: invalid width/maxWidth")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "yy")
	return self
}

// appendValueReduced(TemporalField field, int width, int maxWidth, ChronoLocalDate baseDate)
func dateTimeFormatterBuilderAppendValueReducedDate(params []any) any {
	if len(params) < 5 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: field is null")
	}
	if params[4] == nil || object.IsNull(params[4].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: baseDate is null")
	}
	width, okW := getInt64Val(params[2])
	maxWidth, okM := getInt64Val(params[3])
	if !okW || !okM || width < 1 || width > 10 || maxWidth < 1 || maxWidth > 10 || maxWidth < width {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendValueReduced: invalid width/maxWidth")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "yy")
	return self
}

// appendValueReduced(TemporalField field, int width, int maxWidth, int baseValue, ChronoLocalDate baseDate)
func dateTimeFormatterBuilderAppendValueReducedChronoDate(params []any) any {
	if len(params) < 6 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendValueReduced: field is null")
	}
	width, okW := getInt64Val(params[2])
	maxWidth, okM := getInt64Val(params[3])
	if !okW || !okM || width < 1 || width > 10 || maxWidth < 1 || maxWidth > 10 || maxWidth < width {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.appendValueReduced: invalid width/maxWidth")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "yy")
	return self
}

// appendZoneId()
func dateTimeFormatterBuilderAppendZoneId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneId: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "VV")
	return self
}

// appendZoneOrOffsetId()
func dateTimeFormatterBuilderAppendZoneOrOffsetId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneOrOffsetId: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "VV")
	return self
}

// appendZoneRegionId()
func dateTimeFormatterBuilderAppendZoneRegionId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneRegionId: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "z")
	return self
}

// appendZoneText(TextStyle textStyle)
func dateTimeFormatterBuilderAppendZoneTextStyle(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneText: textStyle is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "zzzz")
	return self
}

// appendZoneText(TextStyle textStyle, Set preferredZones)
func dateTimeFormatterBuilderAppendZoneTextSet(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneText: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.appendZoneText: textStyle is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "zzzz")
	return self
}

// getLocalizedDateTimePattern(FormatStyle dateStyle, FormatStyle timeStyle, Chronology chrono, Locale locale)
func dateTimeFormatterBuilderGetLocalizedDateTimePattern(params []any) any {
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.getLocalizedDateTimePattern: insufficient params")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.getLocalizedDateTimePattern: chrono is null")
	}
	if params[3] == nil || object.IsNull(params[3].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.getLocalizedDateTimePattern: locale is null")
	}
	dateStyleNull := params[0] == nil || object.IsNull(params[0].(*object.Object))
	timeStyleNull := params[1] == nil || object.IsNull(params[1].(*object.Object))
	if dateStyleNull && timeStyleNull {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.getLocalizedDateTimePattern: both dateStyle and timeStyle cannot be null")
	}
	if !dateStyleNull && !timeStyleNull {
		return object.StringObjectFromGoString("yyyy-MM-dd HH:mm:ss")
	} else if !dateStyleNull {
		return object.StringObjectFromGoString("yyyy-MM-dd")
	}
	return object.StringObjectFromGoString("HH:mm:ss")
}

// optionalEnd()
func dateTimeFormatterBuilderOptionalEnd(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.optionalEnd: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "]")
	return self
}

// optionalStart()
func dateTimeFormatterBuilderOptionalStart(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.optionalStart: self is null")
	}
	self := params[0].(*object.Object)
	appendBuilderPattern(self, "[")
	return self
}

// padNext(int padWidth)
func dateTimeFormatterBuilderPadNextPadWidth(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.padNext: self is null")
	}
	width, ok := getInt64Val(params[1])
	if !ok || width < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.padNext: invalid padWidth")
	}
	return params[0].(*object.Object)
}

// padNext(int padWidth, char padChar)
func dateTimeFormatterBuilderPadNextPadWidthPadChar(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.padNext: self is null")
	}
	width, ok := getInt64Val(params[1])
	if !ok || width < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatterBuilder.padNext: invalid padWidth")
	}
	return params[0].(*object.Object)
}

// parseCaseInsensitive()
func dateTimeFormatterBuilderParseCaseInsensitive(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.parseCaseInsensitive: self is null")
	}
	return params[0].(*object.Object)
}

// parseCaseSensitive()
func dateTimeFormatterBuilderParseCaseSensitive(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.parseCaseSensitive: self is null")
	}
	return params[0].(*object.Object)
}

// parseDefaulting(TemporalField field, long value)
func dateTimeFormatterBuilderParseDefaulting(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.parseDefaulting: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.parseDefaulting: field is null")
	}
	return params[0].(*object.Object)
}

// parseLenient()
func dateTimeFormatterBuilderParseLenient(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.parseLenient: self is null")
	}
	return params[0].(*object.Object)
}

// parseStrict()
func dateTimeFormatterBuilderParseStrict(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.parseStrict: self is null")
	}
	return params[0].(*object.Object)
}

func createDateTimeFormatterObject(pattern string, localeObj *object.Object, resolverObj *object.Object, chronoObj *object.Object) *object.Object {
	formatter := object.MakeEmptyObjectWithClassName(&classNameDateTimeFormatter)
	formatter.FieldTable["pattern"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString(pattern),
	}
	if localeObj != nil && !object.IsNull(localeObj) {
		formatter.FieldTable["locale"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: localeObj,
		}
	}
	if resolverObj != nil && !object.IsNull(resolverObj) {
		formatter.FieldTable["resolverStyle"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: resolverObj,
		}
	}
	if chronoObj != nil && !object.IsNull(chronoObj) {
		formatter.FieldTable["chronology"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: chronoObj,
		}
	}
	return formatter
}

// toFormatter()
func dateTimeFormatterBuilderToFormatterNoArgs(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.toFormatter: self is null")
	}
	self := params[0].(*object.Object)
	pattern := getBuilderPattern(self)
	return createDateTimeFormatterObject(pattern, nil, nil, nil)
}

// toFormatter(Locale locale)
func dateTimeFormatterBuilderToFormatterLocale(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.toFormatter: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.toFormatter: locale is null")
	}
	self := params[0].(*object.Object)
	locale := params[1].(*object.Object)
	pattern := getBuilderPattern(self)
	return createDateTimeFormatterObject(pattern, locale, nil, nil)
}

// toFormatter(ResolverStyle resolverStyle, Chronology chronology)
func dateTimeFormatterBuilderToFormatterResolverChronology(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.toFormatter: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatterBuilder.toFormatter: resolverStyle is null")
	}
	self := params[0].(*object.Object)
	resolver := params[1].(*object.Object)
	var chrono *object.Object
	if params[2] != nil && !object.IsNull(params[2].(*object.Object)) {
		chrono = params[2].(*object.Object)
	}
	pattern := getBuilderPattern(self)
	return createDateTimeFormatterObject(pattern, nil, resolver, chrono)
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"strings"
	"sync"
	"time"
)

var dtFormatterInitOnce sync.Once

// Standard DateTimeFormatter constants
var (
	formatterISO_LOCAL_DATE       *object.Object
	formatterISO_OFFSET_DATE      *object.Object
	formatterISO_DATE             *object.Object
	formatterISO_LOCAL_TIME       *object.Object
	formatterISO_OFFSET_TIME      *object.Object
	formatterISO_TIME             *object.Object
	formatterISO_LOCAL_DATE_TIME  *object.Object
	formatterISO_OFFSET_DATE_TIME *object.Object
	formatterISO_ZONED_DATE_TIME  *object.Object
	formatterISO_DATE_TIME        *object.Object
	formatterISO_ORDINAL_DATE     *object.Object
	formatterISO_WEEK_DATE        *object.Object
	formatterISO_INSTANT          *object.Object
	formatterBASIC_ISO_DATE       *object.Object
	formatterRFC_1123_DATE_TIME   *object.Object
)

func initStandardFormatters() {
	dtFormatterInitOnce.Do(func() {
		formatterISO_LOCAL_DATE = createDateTimeFormatterObject("yyyy-MM-dd", nil, nil, nil)
		formatterISO_OFFSET_DATE = createDateTimeFormatterObject("yyyy-MM-ddXXX", nil, nil, nil)
		formatterISO_DATE = createDateTimeFormatterObject("yyyy-MM-dd[XXX]", nil, nil, nil)
		formatterISO_LOCAL_TIME = createDateTimeFormatterObject("HH:mm:ss", nil, nil, nil)
		formatterISO_OFFSET_TIME = createDateTimeFormatterObject("HH:mm:ssXXX", nil, nil, nil)
		formatterISO_TIME = createDateTimeFormatterObject("HH:mm:ss[XXX]", nil, nil, nil)
		formatterISO_LOCAL_DATE_TIME = createDateTimeFormatterObject("yyyy-MM-dd'T'HH:mm:ss", nil, nil, nil)
		formatterISO_OFFSET_DATE_TIME = createDateTimeFormatterObject("yyyy-MM-dd'T'HH:mm:ssXXX", nil, nil, nil)
		formatterISO_ZONED_DATE_TIME = createDateTimeFormatterObject("yyyy-MM-dd'T'HH:mm:ssXXX'['VV']'", nil, nil, nil)
		formatterISO_DATE_TIME = createDateTimeFormatterObject("yyyy-MM-dd'T'HH:mm:ss[XXX]['['VV']']", nil, nil, nil)
		formatterISO_ORDINAL_DATE = createDateTimeFormatterObject("yyyy-DDD", nil, nil, nil)
		formatterISO_WEEK_DATE = createDateTimeFormatterObject("YYYY-'W'ww-e", nil, nil, nil)
		formatterISO_INSTANT = createDateTimeFormatterObject("yyyy-MM-dd'T'HH:mm:ss.SSSX", nil, nil, nil)
		formatterBASIC_ISO_DATE = createDateTimeFormatterObject("yyyyMMdd", nil, nil, nil)
		formatterRFC_1123_DATE_TIME = createDateTimeFormatterObject("EEE, d MMM yyyy HH:mm:ss z", nil, nil, nil)

		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_LOCAL_DATE", statics.Static{Type: types.Ref, Value: formatterISO_LOCAL_DATE})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_OFFSET_DATE", statics.Static{Type: types.Ref, Value: formatterISO_OFFSET_DATE})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_DATE", statics.Static{Type: types.Ref, Value: formatterISO_DATE})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_LOCAL_TIME", statics.Static{Type: types.Ref, Value: formatterISO_LOCAL_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_OFFSET_TIME", statics.Static{Type: types.Ref, Value: formatterISO_OFFSET_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_TIME", statics.Static{Type: types.Ref, Value: formatterISO_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_LOCAL_DATE_TIME", statics.Static{Type: types.Ref, Value: formatterISO_LOCAL_DATE_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_OFFSET_DATE_TIME", statics.Static{Type: types.Ref, Value: formatterISO_OFFSET_DATE_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_ZONED_DATE_TIME", statics.Static{Type: types.Ref, Value: formatterISO_ZONED_DATE_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_DATE_TIME", statics.Static{Type: types.Ref, Value: formatterISO_DATE_TIME})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_ORDINAL_DATE", statics.Static{Type: types.Ref, Value: formatterISO_ORDINAL_DATE})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_WEEK_DATE", statics.Static{Type: types.Ref, Value: formatterISO_WEEK_DATE})
		_ = statics.AddStatic(classNameDateTimeFormatter+".ISO_INSTANT", statics.Static{Type: types.Ref, Value: formatterISO_INSTANT})
		_ = statics.AddStatic(classNameDateTimeFormatter+".BASIC_ISO_DATE", statics.Static{Type: types.Ref, Value: formatterBASIC_ISO_DATE})
		_ = statics.AddStatic(classNameDateTimeFormatter+".RFC_1123_DATE_TIME", statics.Static{Type: types.Ref, Value: formatterRFC_1123_DATE_TIME})
	})
}

// Load_Time_Format_DateTimeFormatter registers all java/time/format/DateTimeFormatter methods.
func Load_Time_Format_DateTimeFormatter() {
	initStandardFormatters()

	// Factories
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.ofPattern(Ljava/lang/String;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterOfPatternString,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.ofPattern(Ljava/lang/String;Ljava/util/Locale;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterOfPatternStringLocale,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.ofLocalizedDate(Ljava/time/format/FormatStyle;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterOfLocalizedDate,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.ofLocalizedTime(Ljava/time/format/FormatStyle;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterOfLocalizedTime,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.ofLocalizedDateTime(Ljava/time/format/FormatStyle;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterOfLocalizedDateTimeSingle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.ofLocalizedDateTime(Ljava/time/format/FormatStyle;Ljava/time/format/FormatStyle;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterOfLocalizedDateTimeDouble,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parsedExcessDays()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterParsedExcessDays,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parsedLeapSecond()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterParsedLeapSecond,
		}

	// Instance methods - Getters / Accessors
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.getLocale()Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterGetLocale,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.getDecimalStyle()Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterGetDecimalStyle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterGetChronology,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.getZone()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterGetZone,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.getResolverStyle()Ljava/time/format/ResolverStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterGetResolverStyle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.getResolverFields()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterGetResolverFields,
		}

	// Withers
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withLocale(Ljava/util/Locale;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithLocale,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withDecimalStyle(Ljava/time/format/DecimalStyle;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithDecimalStyle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withChronology(Ljava/time/chrono/Chronology;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithChronology,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withZone(Ljava/time/ZoneId;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithZone,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withResolverStyle(Ljava/time/format/ResolverStyle;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithResolverStyle,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withResolverFields([Ljava/time/temporal/TemporalField;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithResolverFieldsArray,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.withResolverFields(Ljava/util/Set;)Ljava/time/format/DateTimeFormatter;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterWithResolverFieldsSet,
		}

	// Formatting
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.format(Ljava/time/temporal/TemporalAccessor;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterFormat,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.formatTo(Ljava/time/temporal/TemporalAccessor;Ljava/lang/Appendable;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterFormatTo,
		}

	// Parsing
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parse(Ljava/lang/CharSequence;)Ljava/time/temporal/TemporalAccessor;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterParseCharSequence,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parse(Ljava/lang/CharSequence;Ljava/text/ParsePosition;)Ljava/time/temporal/TemporalAccessor;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterParseCharSequencePosition,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parse(Ljava/lang/CharSequence;Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterParseCharSequenceQuery,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parseBest(Ljava/lang/CharSequence;[Ljava/time/temporal/TemporalQuery;)Ljava/time/temporal/TemporalAccessor;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterParseBest,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.parseUnresolved(Ljava/lang/CharSequence;Ljava/text/ParsePosition;)Ljava/time/temporal/TemporalAccessor;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dateTimeFormatterParseUnresolved,
		}

	// Conversions
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.toFormat()Ljava/text/Format;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterToFormatNoArgs,
		}

	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.toFormat(Ljava/time/temporal/TemporalQuery;)Ljava/text/Format;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dateTimeFormatterToFormatQuery,
		}

	// Object methods
	ghelpers.MethodSignatures["java/time/format/DateTimeFormatter.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dateTimeFormatterToString,
		}
}

// ofPattern(String pattern)
func dateTimeFormatterOfPatternString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofPattern: pattern is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatter.ofPattern: expected String")
	}
	pat := object.GoStringFromStringObject(strObj)
	return createDateTimeFormatterObject(pat, nil, nil, nil)
}

// ofPattern(String pattern, Locale locale)
func dateTimeFormatterOfPatternStringLocale(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofPattern: pattern is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofPattern: locale is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DateTimeFormatter.ofPattern: expected String")
	}
	pat := object.GoStringFromStringObject(strObj)
	localeObj := params[1].(*object.Object)
	return createDateTimeFormatterObject(pat, localeObj, nil, nil)
}

// ofLocalizedDate(FormatStyle dateStyle)
func dateTimeFormatterOfLocalizedDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofLocalizedDate: dateStyle is null")
	}
	return createDateTimeFormatterObject("yyyy-MM-dd", nil, nil, nil)
}

// ofLocalizedTime(FormatStyle timeStyle)
func dateTimeFormatterOfLocalizedTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofLocalizedTime: timeStyle is null")
	}
	return createDateTimeFormatterObject("HH:mm:ss", nil, nil, nil)
}

// ofLocalizedDateTime(FormatStyle dateTimeStyle)
func dateTimeFormatterOfLocalizedDateTimeSingle(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofLocalizedDateTime: dateTimeStyle is null")
	}
	return createDateTimeFormatterObject("yyyy-MM-dd HH:mm:ss", nil, nil, nil)
}

// ofLocalizedDateTime(FormatStyle dateStyle, FormatStyle timeStyle)
func dateTimeFormatterOfLocalizedDateTimeDouble(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofLocalizedDateTime: dateStyle is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.ofLocalizedDateTime: timeStyle is null")
	}
	return createDateTimeFormatterObject("yyyy-MM-dd HH:mm:ss", nil, nil, nil)
}

// parsedExcessDays()
func dateTimeFormatterParsedExcessDays(params []any) any {
	return getExcessDaysTemporalQuery()
}

// parsedLeapSecond()
func dateTimeFormatterParsedLeapSecond(params []any) any {
	return getLeapSecondTemporalQuery()
}

func getExcessDaysTemporalQuery() *object.Object {
	queryObject := object.MakeEmptyObjectWithClassName(&classNameTemporalQueryVal)
	queryObject.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString("excessDays"),
	}
	return queryObject
}

func getLeapSecondTemporalQuery() *object.Object {
	queryObject := object.MakeEmptyObjectWithClassName(&classNameTemporalQueryVal)
	queryObject.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString("leapSecond"),
	}
	return queryObject
}

// getLocale()
func dateTimeFormatterGetLocale(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.getLocale: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if f, ok := self.FieldTable["locale"]; ok && f.Fvalue != nil {
			if lObj, okObj := f.Fvalue.(*object.Object); okObj && !object.IsNull(lObj) {
				return lObj
			}
		}
	}
	loc := object.MakeEmptyObjectWithClassName(&types.ClassNameLocale)
	loc.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("en_US"),
	}
	return loc
}

// getDecimalStyle()
func dateTimeFormatterGetDecimalStyle(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.getDecimalStyle: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if f, ok := self.FieldTable["decimalStyle"]; ok && f.Fvalue != nil {
			if dsObj, okObj := f.Fvalue.(*object.Object); okObj && !object.IsNull(dsObj) {
				return dsObj
			}
		}
	}
	dsClassName := "java/time/format/DecimalStyle"
	ds := object.MakeEmptyObjectWithClassName(&dsClassName)
	return ds
}

// getChronology()
func dateTimeFormatterGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.getChronology: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if f, ok := self.FieldTable["chronology"]; ok && f.Fvalue != nil {
			if cObj, okObj := f.Fvalue.(*object.Object); okObj && !object.IsNull(cObj) {
				return cObj
			}
		}
	}
	return nil
}

// getZone()
func dateTimeFormatterGetZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.getZone: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if f, ok := self.FieldTable["zone"]; ok && f.Fvalue != nil {
			if zObj, okObj := f.Fvalue.(*object.Object); okObj && !object.IsNull(zObj) {
				return zObj
			}
		}
	}
	return nil
}

// getResolverStyle()
func dateTimeFormatterGetResolverStyle(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.getResolverStyle: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if f, ok := self.FieldTable["resolverStyle"]; ok && f.Fvalue != nil {
			if rObj, okObj := f.Fvalue.(*object.Object); okObj && !object.IsNull(rObj) {
				return rObj
			}
		}
	}
	rsClassName := "java/time/format/ResolverStyle"
	rs := object.MakeEmptyObjectWithClassName(&rsClassName)
	rs.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString("SMART"),
	}
	return rs
}

// getResolverFields()
func dateTimeFormatterGetResolverFields(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.getResolverFields: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if f, ok := self.FieldTable["resolverFields"]; ok && f.Fvalue != nil {
			if rfObj, okObj := f.Fvalue.(*object.Object); okObj && !object.IsNull(rfObj) {
				return rfObj
			}
		}
	}
	return nil
}

// cloneDateTimeFormatter clones the object and returns a new copy
func cloneDateTimeFormatter(self *object.Object) *object.Object {
	clone := object.MakeEmptyObjectWithClassName(&classNameDateTimeFormatter)
	if self != nil && self.FieldTable != nil {
		for k, v := range self.FieldTable {
			clone.FieldTable[k] = v
		}
	}
	return clone
}

// withLocale(Locale locale)
func dateTimeFormatterWithLocale(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withLocale: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withLocale: locale is null")
	}
	self := params[0].(*object.Object)
	locale := params[1].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	clone.FieldTable["locale"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: locale,
	}
	return clone
}

// withDecimalStyle(DecimalStyle decimalStyle)
func dateTimeFormatterWithDecimalStyle(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withDecimalStyle: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withDecimalStyle: decimalStyle is null")
	}
	self := params[0].(*object.Object)
	decimalStyle := params[1].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	clone.FieldTable["decimalStyle"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: decimalStyle,
	}
	return clone
}

// withChronology(Chronology chrono)
func dateTimeFormatterWithChronology(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withChronology: self is null")
	}
	self := params[0].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		delete(clone.FieldTable, "chronology")
	} else {
		chrono := params[1].(*object.Object)
		clone.FieldTable["chronology"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: chrono,
		}
	}
	return clone
}

// withZone(ZoneId zone)
func dateTimeFormatterWithZone(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withZone: self is null")
	}
	self := params[0].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		delete(clone.FieldTable, "zone")
	} else {
		zone := params[1].(*object.Object)
		clone.FieldTable["zone"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: zone,
		}
	}
	return clone
}

// withResolverStyle(ResolverStyle resolverStyle)
func dateTimeFormatterWithResolverStyle(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withResolverStyle: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withResolverStyle: resolverStyle is null")
	}
	self := params[0].(*object.Object)
	resolverStyle := params[1].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	clone.FieldTable["resolverStyle"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: resolverStyle,
	}
	return clone
}

// withResolverFields(TemporalField... fields)
func dateTimeFormatterWithResolverFieldsArray(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withResolverFields: self is null")
	}
	self := params[0].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		delete(clone.FieldTable, "resolverFields")
	} else {
		fieldsObj := params[1].(*object.Object)
		clone.FieldTable["resolverFields"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: fieldsObj,
		}
	}
	return clone
}

// withResolverFields(Set fields)
func dateTimeFormatterWithResolverFieldsSet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.withResolverFields: self is null")
	}
	self := params[0].(*object.Object)
	clone := cloneDateTimeFormatter(self)
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		delete(clone.FieldTable, "resolverFields")
	} else {
		fieldsObj := params[1].(*object.Object)
		clone.FieldTable["resolverFields"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: fieldsObj,
		}
	}
	return clone
}

func getFormatterPattern(formatter *object.Object) string {
	if formatter == nil || formatter.FieldTable == nil {
		return ""
	}
	if f, ok := formatter.FieldTable["pattern"]; ok && f.Fvalue != nil {
		if strObj, okStr := f.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
			return object.GoStringFromStringObject(strObj)
		}
	}
	return ""
}

// format(TemporalAccessor temporal)
func dateTimeFormatterFormat(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.format: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.format: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	formatted, err := formatTemporalWithPattern(self, temporal)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "DateTimeFormatter.format: "+err.Error())
	}
	return object.StringObjectFromGoString(formatted)
}

// formatTo(TemporalAccessor temporal, Appendable appendable)
func dateTimeFormatterFormatTo(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.formatTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.formatTo: temporal is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.formatTo: appendable is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	appendable := params[2].(*object.Object)

	formatted, err := formatTemporalWithPattern(self, temporal)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "DateTimeFormatter.formatTo: "+err.Error())
	}

	appendStringToStringBufferOrBuilder(appendable, formatted)
	return nil
}

func appendStringToStringBufferOrBuilder(sb *object.Object, s string) {
	if sb == nil || sb.FieldTable == nil {
		return
	}
	if f, ok := sb.FieldTable["value"]; ok && f.Fvalue != nil {
		if curStr, okStr := f.Fvalue.(*object.Object); okStr && object.IsStringObject(curStr) {
			newStr := object.GoStringFromStringObject(curStr) + s
			sb.FieldTable["value"] = object.Field{
				Ftype:  types.StringClassRef,
				Fvalue: object.StringObjectFromGoString(newStr),
			}
		}
	}
}

func getObjectClassName(obj *object.Object) string {
	if obj == nil {
		return ""
	}
	return object.GoStringFromStringPoolIndex(obj.KlassName)
}

func formatTemporalWithPattern(formatter *object.Object, temporal *object.Object) (string, error) {
	pat := getFormatterPattern(formatter)
	className := getObjectClassName(temporal)

	// Extract standard components from temporal if possible
	hasDate := false
	hasTime := false
	var year, month, day, hour, minute, second, nano int64

	if strings.Contains(className, "LocalDateTime") {
		ld := getLocalDateTimeDate(temporal)
		lt := getLocalDateTimeTime(temporal)
		if ld != nil {
			if y, m, d, ok := getLocalDateYMD(ld); ok {
				hasDate = true
				year, month, day = y, m, d
			}
		}
		if lt != nil {
			hasTime = true
			hour, minute, second, nano = getLocalTimeHMSN(lt)
		}
	} else if strings.Contains(className, "LocalDate") {
		if y, m, d, ok := getLocalDateYMD(temporal); ok {
			hasDate = true
			year, month, day = y, m, d
		}
	} else if strings.Contains(className, "LocalTime") {
		hasTime = true
		hour, minute, second, nano = getLocalTimeHMSN(temporal)
	}

	if pat == "" {
		if hasDate && hasTime {
			return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, minute, second), nil
		} else if hasDate {
			return fmt.Sprintf("%04d-%02d-%02d", year, month, day), nil
		} else if hasTime {
			return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second), nil
		}
		return "", nil
	}

	res := pat
	if hasDate {
		res = strings.ReplaceAll(res, "yyyy", fmt.Sprintf("%04d", year))
		res = strings.ReplaceAll(res, "MM", fmt.Sprintf("%02d", month))
		res = strings.ReplaceAll(res, "dd", fmt.Sprintf("%02d", day))
		res = strings.ReplaceAll(res, "d", fmt.Sprintf("%d", day))
		res = strings.ReplaceAll(res, "M", fmt.Sprintf("%d", month))
	}
	if hasTime {
		res = strings.ReplaceAll(res, "HH", fmt.Sprintf("%02d", hour))
		res = strings.ReplaceAll(res, "mm", fmt.Sprintf("%02d", minute))
		res = strings.ReplaceAll(res, "ss", fmt.Sprintf("%02d", second))
		if nano > 0 {
			res = strings.ReplaceAll(res, "SSS", fmt.Sprintf("%03d", nano/1000000))
			res = strings.ReplaceAll(res, "SSSSSSSSS", fmt.Sprintf("%09d", nano))
		} else {
			res = strings.ReplaceAll(res, ".SSS", "")
			res = strings.ReplaceAll(res, ".SSSSSSSSS", "")
			res = strings.ReplaceAll(res, "SSS", "000")
		}
	}
	res = strings.ReplaceAll(res, "'T'", "T")
	res = strings.ReplaceAll(res, "[", "")
	res = strings.ReplaceAll(res, "]", "")
	res = strings.ReplaceAll(res, "XXX", "Z")
	res = strings.ReplaceAll(res, "X", "Z")
	res = strings.ReplaceAll(res, "VV", "UTC")
	res = strings.ReplaceAll(res, "z", "UTC")

	return res, nil
}

// parse(CharSequence text)
func dateTimeFormatterParseCharSequence(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: text is null")
	}
	self := params[0].(*object.Object)
	textObj := params[1].(*object.Object)
	text := object.GoStringFromStringObject(textObj)

	accessor, err := parseTextToTemporalAccessor(self, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "DateTimeFormatter.parse: "+err.Error())
	}
	return accessor
}

// parse(CharSequence text, ParsePosition position)
func dateTimeFormatterParseCharSequencePosition(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: text is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: position is null")
	}
	self := params[0].(*object.Object)
	textObj := params[1].(*object.Object)
	text := object.GoStringFromStringObject(textObj)

	accessor, err := parseTextToTemporalAccessor(self, text)
	if err != nil {
		return nil
	}
	return accessor
}

// parse(CharSequence text, TemporalQuery query)
func dateTimeFormatterParseCharSequenceQuery(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: text is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parse: query is null")
	}
	self := params[0].(*object.Object)
	textObj := params[1].(*object.Object)
	query := params[2].(*object.Object)
	text := object.GoStringFromStringObject(textObj)

	accessor, err := parseTextToTemporalAccessor(self, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "DateTimeFormatter.parse: "+err.Error())
	}

	res := queryFromTemporalAccessor(accessor, query)
	return res
}

// parseBest(CharSequence text, TemporalQuery... queries)
func dateTimeFormatterParseBest(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parseBest: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parseBest: text is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parseBest: queries is null")
	}
	self := params[0].(*object.Object)
	textObj := params[1].(*object.Object)
	text := object.GoStringFromStringObject(textObj)

	accessor, err := parseTextToTemporalAccessor(self, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "DateTimeFormatter.parseBest: "+err.Error())
	}
	return accessor
}

// parseUnresolved(CharSequence text, ParsePosition position)
func dateTimeFormatterParseUnresolved(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parseUnresolved: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parseUnresolved: text is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.parseUnresolved: position is null")
	}
	self := params[0].(*object.Object)
	textObj := params[1].(*object.Object)
	text := object.GoStringFromStringObject(textObj)

	accessor, err := parseTextToTemporalAccessor(self, text)
	if err != nil {
		return nil
	}
	return accessor
}

func parseTextToTemporalAccessor(formatter *object.Object, text string) (*object.Object, error) {
	text = strings.TrimSpace(text)
	// Try ISO LocalDateTime first: 2023-11-05T08:14:00 or 2023-11-05 08:14:00
	if strings.Contains(text, "T") || (strings.Contains(text, "-") && strings.Contains(text, ":")) {
		var year, month, day, hour, minute, second int
		n, _ := fmt.Sscanf(text, "%d-%d-%dT%d:%d:%d", &year, &month, &day, &hour, &minute, &second)
		if n < 5 {
			n, _ = fmt.Sscanf(text, "%d-%d-%d %d:%d:%d", &year, &month, &day, &hour, &minute, &second)
		}
		if n >= 5 {
			ld := CreateLocalDate(int64(year), int64(month), int64(day))
			lt := createLocalTime(int64(hour), int64(minute), int64(second), 0)
			return createLocalDateTime(ld, lt), nil
		}
	}

	// Try LocalDate: 2023-11-05
	if strings.Contains(text, "-") {
		var year, month, day int
		n, _ := fmt.Sscanf(text, "%d-%d-%d", &year, &month, &day)
		if n == 3 {
			return CreateLocalDate(int64(year), int64(month), int64(day)), nil
		}
	}

	// Try LocalTime: 08:14:00 or 08:14
	if strings.Contains(text, ":") {
		var hour, minute, second int
		n, _ := fmt.Sscanf(text, "%d:%d:%d", &hour, &minute, &second)
		if n >= 2 {
			return createLocalTime(int64(hour), int64(minute), int64(second), 0), nil
		}
	}

	// Try Basic ISO Date: 20231105
	if len(text) == 8 {
		var year, month, day int
		n, _ := fmt.Sscanf(text, "%04d%02d%02d", &year, &month, &day)
		if n == 3 && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
			return CreateLocalDate(int64(year), int64(month), int64(day)), nil
		}
	}

	// Try RFC 1123 or time.RFC1123 parsing
	if t, err := time.Parse(time.RFC1123, text); err == nil {
		ld := CreateLocalDate(int64(t.Year()), int64(t.Month()), int64(t.Day()))
		lt := createLocalTime(int64(t.Hour()), int64(t.Minute()), int64(t.Second()), int64(t.Nanosecond()))
		return createLocalDateTime(ld, lt), nil
	}

	return nil, fmt.Errorf("Text '%s' could not be parsed", text)
}

func queryFromTemporalAccessor(accessor *object.Object, query *object.Object) *object.Object {
	if accessor == nil || query == nil {
		return nil
	}
	className := getObjectClassName(accessor)

	qName := ""
	if query.FieldTable != nil {
		if f, ok := query.FieldTable["name"]; ok && f.Fvalue != nil {
			if strObj, okStr := f.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				qName = object.GoStringFromStringObject(strObj)
			}
		}
	}

	switch qName {
	case "localDate":
		if strings.Contains(className, "LocalDate") {
			return accessor
		} else if strings.Contains(className, "LocalDateTime") {
			return getLocalDateTimeDate(accessor)
		}
	case "localTime":
		if strings.Contains(className, "LocalTime") {
			return accessor
		} else if strings.Contains(className, "LocalDateTime") {
			return getLocalDateTimeTime(accessor)
		}
	case "chronology":
		return getIsoChronologySingleton()
	case "precision":
		if strings.Contains(className, "LocalDate") {
			unitObj := object.MakeEmptyObjectWithClassName(&classNameChronoUnitVal)
			unitObj.FieldTable["name"] = object.Field{
				Ftype:  types.StringClassRef,
				Fvalue: object.StringObjectFromGoString("DAYS"),
			}
			return unitObj
		} else if strings.Contains(className, "LocalTime") || strings.Contains(className, "LocalDateTime") {
			unitObj := object.MakeEmptyObjectWithClassName(&classNameChronoUnitVal)
			unitObj.FieldTable["name"] = object.Field{
				Ftype:  types.StringClassRef,
				Fvalue: object.StringObjectFromGoString("NANOS"),
			}
			return unitObj
		}
	}

	return accessor
}

// toFormat()
func dateTimeFormatterToFormatNoArgs(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.toFormat: self is null")
	}
	formatClassName := "java/text/Format"
	fmtObj := object.MakeEmptyObjectWithClassName(&formatClassName)
	return fmtObj
}

// toFormat(TemporalQuery query)
func dateTimeFormatterToFormatQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.toFormat: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.toFormat: query is null")
	}
	formatClassName := "java/text/Format"
	fmtObj := object.MakeEmptyObjectWithClassName(&formatClassName)
	return fmtObj
}

// toString()
func dateTimeFormatterToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DateTimeFormatter.toString: self is null")
	}
	self := params[0].(*object.Object)
	pat := getFormatterPattern(self)
	return object.StringObjectFromGoString(pat)
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
	"time"
)

var (
	classNameChronoLocalDateTimeComparator = "java/time/chrono/ChronoLocalDateTime$1"
	classNameInstant                       = "java/time/Instant"
)

// Load_Time_Chrono_ChronoLocalDateTime loads method signatures for java/time/chrono/ChronoLocalDateTime.
func Load_Time_Chrono_ChronoLocalDateTime() {
	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.atZone(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeAtZone,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.compareTo(Ljava/time/chrono/ChronoLocalDateTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeEquals,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeFormat,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeFrom,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimeGetChronology,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimeHashCode,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.isAfter(Ljava/time/chrono/ChronoLocalDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeIsAfter,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.isBefore(Ljava/time/chrono/ChronoLocalDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeIsBefore,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.isEqual(Ljava/time/chrono/ChronoLocalDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeIsEqual,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoLocalDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoLocalDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeQuery,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.timeLineOrder()Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimeTimelineOrder,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.toEpochSecond(Ljava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.toInstant(Ljava/time/ZoneOffset;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeToInstant,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.toLocalDate()Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.toLocalTime()Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimeToLocalTime,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimeToString,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronoLocalDateTimeUntilUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoLocalDateTimeWithField,
		}
}

// Helpers for ChronoLocalDateTime

func createChronoLocalDateTime(dateObj *object.Object, timeObj *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoLocalDateTime)
	obj.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: dateObj}
	obj.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: timeObj}
	return obj
}

func getChronoLocalDateTimeDate(obj *object.Object) *object.Object {
	if obj != nil && obj.FieldTable != nil {
		if dFld, ok := obj.FieldTable["date"]; ok {
			if dObj, ok := dFld.Fvalue.(*object.Object); ok && dObj != nil {
				return dObj
			}
		}
	}
	return createChronoLocalDate(GetIsoChronology(), 1970, 1, 1)
}

func getChronoLocalDateTimeTime(obj *object.Object) *object.Object {
	if obj != nil && obj.FieldTable != nil {
		if tFld, ok := obj.FieldTable["time"]; ok {
			if tObj, ok := tFld.Fvalue.(*object.Object); ok && tObj != nil {
				return tObj
			}
		}
	}
	return createLocalTime(0, 0, 0, 0)
}

func getZoneOffsetTotalSeconds(offset *object.Object) int64 {
	if offset == nil || offset.FieldTable == nil {
		return 0
	}
	if secFld, ok := offset.FieldTable["totalSeconds"]; ok {
		if s, ok := getInt64Val(secFld.Fvalue); ok {
			return s
		}
	}
	if secFld, ok := offset.FieldTable["seconds"]; ok {
		if s, ok := getInt64Val(secFld.Fvalue); ok {
			return s
		}
	}
	return 0
}

func createChronoZonedDateTime(dtObj *object.Object, zoneObj *object.Object, offsetObj *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoZonedDateTime)
	obj.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: dtObj}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zoneObj}
	obj.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return obj
}

// chronoLocalDateTimeAdjustInto adjusts the specified temporal object to have the same date and time as this object.
// default Temporal adjustInto(Temporal temporal)
func chronoLocalDateTimeAdjustInto(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	target := params[1].(*object.Object)

	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(dateObj)
	year, month, day := ghelpers.GetChronoLocalDateYMD(dateObj)
	hour, minutes, sec, nano := getLocalTimeHMSN(timeObj)
	nanoOfDay := getLocalTimeNanoOfDay(timeObj)

	if target.FieldTable != nil {
		target.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: dateObj}
		target.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: timeObj}
		target.FieldTable["epochDay"] = object.Field{Ftype: types.Long, Fvalue: epochDay}
		target.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: year}
		target.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: month}
		target.FieldTable["day"] = object.Field{Ftype: types.Int, Fvalue: day}
		target.FieldTable["hour"] = object.Field{Ftype: types.Int, Fvalue: hour}
		target.FieldTable["minute"] = object.Field{Ftype: types.Int, Fvalue: minutes}
		target.FieldTable["second"] = object.Field{Ftype: types.Int, Fvalue: sec}
		target.FieldTable["nano"] = object.Field{Ftype: types.Int, Fvalue: nano}
		target.FieldTable["nanoOfDay"] = object.Field{Ftype: types.Long, Fvalue: nanoOfDay}
	}
	return target
}

// chronoLocalDateTimeAtZone combines this date-time with a time-zone to create a ChronoZonedDateTime.
// ChronoZonedDateTime<D> atZone(ZoneId zone)
func chronoLocalDateTimeAtZone(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.atZone: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.atZone: zone is null")
	}
	self := params[0].(*object.Object)
	zone := params[1].(*object.Object)

	return createChronoZonedDateTime(self, zone, nil)
}

// chronoLocalDateTimeCompareTo compares this date-time to another date-time, including the chronology.
// default int compareTo(ChronoLocalDateTime<?> other)
func chronoLocalDateTimeCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	date1 := getChronoLocalDateTimeDate(self)
	date2 := getChronoLocalDateTimeDate(other)

	res := chronoLocalDateCompareTo([]any{date1, date2})
	if cmp, ok := res.(int64); ok && cmp != 0 {
		return cmp
	}

	time1 := getChronoLocalDateTimeTime(self)
	time2 := getChronoLocalDateTimeTime(other)
	timeCmp := compareLocalTime(time1, time2)
	if timeCmp != 0 {
		return timeCmp
	}

	chrono1 := GetChronoLocalDateChronology(date1)
	chrono2 := GetChronoLocalDateChronology(date2)
	id1 := getChronologyId(chrono1)
	id2 := getChronologyId(chrono2)
	if id1 < id2 {
		return int64(-1)
	} else if id1 > id2 {
		return int64(1)
	}
	return int64(0)
}

// chronoLocalDateTimeEquals checks if this date-time is equal to another date-time, including the chronology.
// boolean equals(Object obj)
func chronoLocalDateTimeEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.equals: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	res := chronoLocalDateTimeCompareTo([]any{self, other})
	if cmp, ok := res.(int64); ok && cmp == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateTimeFormat formats this date-time using the specified formatter.
// default String format(DateTimeFormatter formatter)
func chronoLocalDateTimeFormat(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.format: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.format: formatter is null")
	}
	self := params[0].(*object.Object)

	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)
	y, m, d := ghelpers.GetChronoLocalDateYMD(dateObj)
	hr, minutes, sec, _ := getLocalTimeHMSN(timeObj)

	str := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", y, m, d, hr, minutes, sec)
	return object.StringObjectFromGoString(str)
}

// chronoLocalDateTimeFrom obtains an instance of ChronoLocalDateTime from a temporal object.
// static ChronoLocalDateTime<?> from(TemporalAccessor temporal)
func chronoLocalDateTimeFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.from: temporal is null")
	}
	temporal := params[0].(*object.Object)

	if temporal.FieldTable != nil {
		if dFld, ok := temporal.FieldTable["date"]; ok {
			if tFld, ok2 := temporal.FieldTable["time"]; ok2 {
				if dObj, ok := dFld.Fvalue.(*object.Object); ok {
					if tObj, ok := tFld.Fvalue.(*object.Object); ok {
						return createChronoLocalDateTime(dObj, tObj)
					}
				}
			}
		}
	}

	chrono := GetChronoLocalDateChronology(temporal)
	dateRes := chronoLocalDateFrom([]any{temporal})
	if dateObj, ok := dateRes.(*object.Object); ok {
		h, m, s, n := getLocalTimeHMSN(temporal)
		timeObj := createLocalTime(h, m, s, n)
		return createChronoLocalDateTime(dateObj, timeObj)
	}

	y, m, d := ghelpers.GetChronoLocalDateYMD(temporal)
	h, minutes, s, n := getLocalTimeHMSN(temporal)
	dateObj := createChronoLocalDate(chrono, y, m, d)
	timeObj := createLocalTime(h, minutes, s, n)
	return createChronoLocalDateTime(dateObj, timeObj)
}

// chronoLocalDateTimeGetChronology gets the chronology of this date-time.
// default Chronology getChronology()
func chronoLocalDateTimeGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.getChronology: self is null")
	}
	self := params[0].(*object.Object)
	dateObj := getChronoLocalDateTimeDate(self)
	return GetChronoLocalDateChronology(dateObj)
}

// chronoLocalDateTimeHashCode returns a hash code for this date-time.
// int hashCode()
func chronoLocalDateTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)

	dateHash := int64(0)
	if hRes := chronoLocalDateHashCode([]any{dateObj}); hRes != nil {
		if h, ok := getInt64Val(hRes); ok {
			dateHash = h
		}
	}
	timeHash := getLocalTimeNanoOfDay(timeObj)
	return int64(int32(dateHash) ^ int32(timeHash^(timeHash>>32)))
}

// chronoLocalDateTimeIsAfter checks if this date-time is after the specified date-time ignoring the chronology.
// default boolean isAfter(ChronoLocalDateTime<?> other)
func chronoLocalDateTimeIsAfter(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isAfter: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isAfter: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	date1 := getChronoLocalDateTimeDate(self)
	date2 := getChronoLocalDateTimeDate(other)
	epoch1 := ghelpers.GetChronoLocalDateEpochDay(date1)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(date2)

	if epoch1 > epoch2 {
		return types.JavaBoolTrue
	} else if epoch1 < epoch2 {
		return types.JavaBoolFalse
	}

	time1 := getChronoLocalDateTimeTime(self)
	time2 := getChronoLocalDateTimeTime(other)
	if getLocalTimeNanoOfDay(time1) > getLocalTimeNanoOfDay(time2) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateTimeIsBefore checks if this date-time is before the specified date-time ignoring the chronology.
// default boolean isBefore(ChronoLocalDateTime<?> other)
func chronoLocalDateTimeIsBefore(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isBefore: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isBefore: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	date1 := getChronoLocalDateTimeDate(self)
	date2 := getChronoLocalDateTimeDate(other)
	epoch1 := ghelpers.GetChronoLocalDateEpochDay(date1)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(date2)

	if epoch1 < epoch2 {
		return types.JavaBoolTrue
	} else if epoch1 > epoch2 {
		return types.JavaBoolFalse
	}

	time1 := getChronoLocalDateTimeTime(self)
	time2 := getChronoLocalDateTimeTime(other)
	if getLocalTimeNanoOfDay(time1) < getLocalTimeNanoOfDay(time2) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateTimeIsEqual checks if this date-time is equal to the specified date-time ignoring the chronology.
// default boolean isEqual(ChronoLocalDateTime<?> other)
func chronoLocalDateTimeIsEqual(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isEqual: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isEqual: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	date1 := getChronoLocalDateTimeDate(self)
	date2 := getChronoLocalDateTimeDate(other)
	epoch1 := ghelpers.GetChronoLocalDateEpochDay(date1)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(date2)

	if epoch1 != epoch2 {
		return types.JavaBoolFalse
	}

	time1 := getChronoLocalDateTimeTime(self)
	time2 := getChronoLocalDateTimeTime(other)
	if getLocalTimeNanoOfDay(time1) == getLocalTimeNanoOfDay(time2) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateTimeIsSupportedField checks if the specified field is supported.
// boolean isSupported(TemporalField field)
func chronoLocalDateTimeIsSupportedField(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return types.JavaBoolTrue
	}

	switch fieldName {
	case "NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY", "SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY",
		"DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY", "ALIGNED_WEEK_OF_MONTH",
		"ALIGNED_WEEK_OF_YEAR", "MONTH_OF_YEAR", "PROLEPTIC_MONTH", "YEAR_OF_ERA",
		"YEAR", "ERA":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// chronoLocalDateTimeIsSupportedUnit checks if the specified unit is supported.
// default boolean isSupported(TemporalUnit unit)
func chronoLocalDateTimeIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	unitObj := params[1].(*object.Object)
	unitName := getTemporalUnitName(unitObj)
	if unitName == "" {
		return types.JavaBoolTrue
	}

	switch unitName {
	case "NANOS", "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS",
		"DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// chronoLocalDateTimeMinusUnit returns an object of the same type with the specified period subtracted.
// default ChronoLocalDateTime<D> minus(long amountToSubtract, TemporalUnit unit)
func chronoLocalDateTimeMinusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.minus: self is null")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoLocalDateTime.minus: invalid amount")
	}
	unit := params[2].(*object.Object)
	return chronoLocalDateTimePlusUnit([]any{params[0], -amount, unit})
}

// chronoLocalDateTimeMinusAmount returns an object of the same type with the specified amount subtracted.
// default ChronoLocalDateTime<D> minus(TemporalAmount amount)
func chronoLocalDateTimeMinusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	days := ghelpers.GetTemporalAmountDays(amount)
	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(dateObj) - days
	chrono := GetChronoLocalDateChronology(dateObj)
	t := time.Unix(epochDay*86400, 0).UTC()
	newDate := createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
	return createChronoLocalDateTime(newDate, timeObj)
}

// chronoLocalDateTimePlusUnit returns an object of the same type with the specified period added.
// ChronoLocalDateTime<D> plus(long amountToAdd, TemporalUnit unit)
func chronoLocalDateTimePlusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.plus: self is null")
	}
	self := params[0].(*object.Object)
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoLocalDateTime.plus: invalid amount")
	}
	unit := params[2].(*object.Object)

	unitName := strings.ToUpper(getTemporalUnitName(unit))
	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)
	nod := getLocalTimeNanoOfDay(timeObj)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(dateObj)
	chrono := GetChronoLocalDateChronology(dateObj)

	const nanosPerDay = int64(86400000000000)

	switch unitName {
	case "NANOS":
		totalNanos := nod + amount
		dayDelta := totalNanos / nanosPerDay
		newNod := totalNanos % nanosPerDay
		if newNod < 0 {
			newNod += nanosPerDay
			dayDelta--
		}
		newEpoch := epochDay + dayDelta
		t := time.Unix(newEpoch*86400, 0).UTC()
		newDate := createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
		h := newNod / 3600000000000
		newNod %= 3600000000000
		m := newNod / 60000000000
		newNod %= 60000000000
		s := newNod / 1000000000
		n := newNod % 1000000000
		return createChronoLocalDateTime(newDate, createLocalTime(h, m, s, n))
	case "MICROS":
		return chronoLocalDateTimePlusUnit([]any{self, amount * 1000, ghelpers.CreateTemporalUnit("NANOS")})
	case "MILLIS":
		return chronoLocalDateTimePlusUnit([]any{self, amount * 1000000, ghelpers.CreateTemporalUnit("NANOS")})
	case "SECONDS":
		return chronoLocalDateTimePlusUnit([]any{self, amount * 1000000000, ghelpers.CreateTemporalUnit("NANOS")})
	case "MINUTES":
		return chronoLocalDateTimePlusUnit([]any{self, amount * 60000000000, ghelpers.CreateTemporalUnit("NANOS")})
	case "HOURS":
		return chronoLocalDateTimePlusUnit([]any{self, amount * 3600000000000, ghelpers.CreateTemporalUnit("NANOS")})
	case "HALF_DAYS":
		return chronoLocalDateTimePlusUnit([]any{self, amount * 12 * 3600000000000, ghelpers.CreateTemporalUnit("NANOS")})
	default:
		// Date-based units delegate to dateObj plusUnit
		newDateRes := chronoLocalDatePlusUnit([]any{dateObj, amount, unit})
		if newDateObj, ok := newDateRes.(*object.Object); ok {
			return createChronoLocalDateTime(newDateObj, timeObj)
		}
		return self
	}
}

// chronoLocalDateTimePlusAmount returns an object of the same type with the specified amount added.
// default ChronoLocalDateTime<D> plus(TemporalAmount amount)
func chronoLocalDateTimePlusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	days := ghelpers.GetTemporalAmountDays(amount)
	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(dateObj) + days
	chrono := GetChronoLocalDateChronology(dateObj)
	t := time.Unix(epochDay*86400, 0).UTC()
	newDate := createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
	return createChronoLocalDateTime(newDate, timeObj)
}

// chronoLocalDateTimeQuery queries this date-time using the specified query.
// default <R> R query(TemporalQuery<R> query)
func chronoLocalDateTimeQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.query: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)

	if query.FieldTable != nil {
		if nameFld, ok := query.FieldTable["name"]; ok {
			if nameObj, ok := nameFld.Fvalue.(*object.Object); ok && object.IsStringObject(nameObj) {
				qName := object.GoStringFromStringObject(nameObj)
				switch qName {
				case "CHRONOLOGY", "chronology":
					return GetChronoLocalDateChronology(dateObj)
				case "PRECISION", "precision":
					return ghelpers.CreateTemporalUnit("Nanos")
				case "LOCAL_DATE", "localDate":
					return dateObj
				case "LOCAL_TIME", "localTime":
					return timeObj
				case "ZONE_ID", "zoneId", "ZONE", "zone", "OFFSET", "offset":
					return nil
				}
			}
		}
	}

	return temporalQueryQueryFrom([]any{query, self})
}

// chronoLocalDateTimeTimelineOrder gets a comparator that compares ChronoLocalDateTime in time-line order ignoring the chronology.
// static Comparator<ChronoLocalDateTime<?>> timeLineOrder()
func chronoLocalDateTimeTimelineOrder(params []any) any {
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoLocalDateTimeComparator)
	return obj
}

// chronoLocalDateTimeToEpochSecond converts this date-time to the number of seconds from the epoch of 1970-01-01T00:00:00Z.
// default long toEpochSecond(ZoneOffset offset)
func chronoLocalDateTimeToEpochSecond(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toEpochSecond: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toEpochSecond: offset is null")
	}
	self := params[0].(*object.Object)
	offset := params[1].(*object.Object)

	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(dateObj)
	secOfDay := getLocalTimeNanoOfDay(timeObj) / 1000000000
	offsetSec := getZoneOffsetTotalSeconds(offset)

	epochSec := epochDay*86400 + secOfDay - offsetSec
	return epochSec
}

// chronoLocalDateTimeToInstant converts this date-time to an Instant.
// default Instant toInstant(ZoneOffset offset)
func chronoLocalDateTimeToInstant(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toInstant: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toInstant: offset is null")
	}
	self := params[0].(*object.Object)
	offset := params[1].(*object.Object)

	epochSecRes := chronoLocalDateTimeToEpochSecond([]any{self, offset})
	epochSec, ok := epochSecRes.(int64)
	if !ok {
		return epochSecRes
	}

	timeObj := getChronoLocalDateTimeTime(self)
	_, _, _, nano := getLocalTimeHMSN(timeObj)

	obj := object.MakeEmptyObjectWithClassName(&classNameInstant)
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: nano}
	return obj
}

// chronoLocalDateTimeToLocalDate gets the ChronoLocalDate part of this date-time.
// ChronoLocalDate toLocalDate()
func chronoLocalDateTimeToLocalDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toLocalDate: self is null")
	}
	self := params[0].(*object.Object)
	return getChronoLocalDateTimeDate(self)
}

// chronoLocalDateTimeToLocalTime gets the LocalTime part of this date-time.
// LocalTime toLocalTime()
func chronoLocalDateTimeToLocalTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toLocalTime: self is null")
	}
	self := params[0].(*object.Object)
	return getChronoLocalDateTimeTime(self)
}

// chronoLocalDateTimeToString outputs this date-time as a String.
// String toString()
func chronoLocalDateTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.toString: self is null")
	}
	self := params[0].(*object.Object)
	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)

	y, m, d := ghelpers.GetChronoLocalDateYMD(dateObj)
	h, minutes, s, n := getLocalTimeHMSN(timeObj)

	var str string
	if n > 0 {
		str = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%09d", y, m, d, h, minutes, s, n)
	} else {
		str = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", y, m, d, h, minutes, s)
	}
	return object.StringObjectFromGoString(str)
}

// chronoLocalDateTimeUntilUnit calculates the amount of time until another date-time in terms of the specified unit.
// long until(Temporal endExclusive, TemporalUnit unit)
func chronoLocalDateTimeUntilUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.until: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.until: endExclusive is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.until: unit is null")
	}
	self := params[0].(*object.Object)
	end := params[1].(*object.Object)
	unit := params[2].(*object.Object)

	unitName := strings.ToUpper(getTemporalUnitName(unit))

	startDate := getChronoLocalDateTimeDate(self)
	startTime := getChronoLocalDateTimeTime(self)
	endDate := getChronoLocalDateTimeDate(end)
	endTime := getChronoLocalDateTimeTime(end)

	startEpoch := ghelpers.GetChronoLocalDateEpochDay(startDate)
	endEpoch := ghelpers.GetChronoLocalDateEpochDay(endDate)
	startNod := getLocalTimeNanoOfDay(startTime)
	endNod := getLocalTimeNanoOfDay(endTime)

	diffNanos := (endEpoch-startEpoch)*86400000000000 + (endNod - startNod)

	switch unitName {
	case "NANOS":
		return diffNanos
	case "MICROS":
		return diffNanos / 1000
	case "MILLIS":
		return diffNanos / 1000000
	case "SECONDS":
		return diffNanos / 1000000000
	case "MINUTES":
		return diffNanos / 60000000000
	case "HOURS":
		return diffNanos / 3600000000000
	case "HALF_DAYS":
		return diffNanos / (12 * 3600000000000)
	default:
		return chronoLocalDateUntilUnit([]any{startDate, endDate, unit})
	}
}

// chronoLocalDateTimeWithAdjuster returns an adjusted copy of this date-time.
// default ChronoLocalDateTime<D> with(TemporalAdjuster adjuster)
func chronoLocalDateTimeWithAdjuster(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	res := temporalAdjusterAdjustInto([]any{adjuster, self})
	if resObj, ok := res.(*object.Object); ok {
		return resObj
	}
	return self
}

// chronoLocalDateTimeWithField returns an object of the same type with the specified field altered.
// ChronoLocalDateTime<D> with(TemporalField field, long newValue)
func chronoLocalDateTimeWithField(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDateTime.with: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	newValue, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoLocalDateTime.with: invalid newValue")
	}

	fieldName := strings.ToUpper(getTemporalFieldName(field))
	dateObj := getChronoLocalDateTimeDate(self)
	timeObj := getChronoLocalDateTimeTime(self)
	h, m, s, n := getLocalTimeHMSN(timeObj)

	switch fieldName {
	case "NANO_OF_SECOND":
		return createChronoLocalDateTime(dateObj, createLocalTime(h, m, s, newValue))
	case "NANO_OF_DAY":
		h = newValue / 3600000000000
		rem := newValue % 3600000000000
		m = rem / 60000000000
		rem %= 60000000000
		s = rem / 1000000000
		n = rem % 1000000000
		return createChronoLocalDateTime(dateObj, createLocalTime(h, m, s, n))
	case "MICRO_OF_SECOND":
		return createChronoLocalDateTime(dateObj, createLocalTime(h, m, s, newValue*1000))
	case "MILLI_OF_SECOND":
		return createChronoLocalDateTime(dateObj, createLocalTime(h, m, s, newValue*1000000))
	case "SECOND_OF_MINUTE":
		return createChronoLocalDateTime(dateObj, createLocalTime(h, m, newValue, n))
	case "SECOND_OF_DAY":
		h = newValue / 3600
		rem := newValue % 3600
		m = rem / 60
		s = rem % 60
		return createChronoLocalDateTime(dateObj, createLocalTime(h, m, s, n))
	case "MINUTE_OF_HOUR":
		return createChronoLocalDateTime(dateObj, createLocalTime(h, newValue, s, n))
	case "HOUR_OF_DAY":
		return createChronoLocalDateTime(dateObj, createLocalTime(newValue, m, s, n))
	default:
		// Delegate to dateObj withField
		newDateRes := chronoLocalDateWithField([]any{dateObj, field, newValue})
		if newDateObj, ok := newDateRes.(*object.Object); ok {
			return createChronoLocalDateTime(newDateObj, timeObj)
		}
		return self
	}
}

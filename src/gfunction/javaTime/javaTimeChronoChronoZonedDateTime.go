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
	"math"
	"strings"
	"time"
)

var (
	classNameChronoZonedDateTimeComparator = "java/time/chrono/ChronoZonedDateTime$1"
	classNameZoneId                        = "java/time/ZoneId"
	classNameZoneOffset                    = "java/time/ZoneOffset"
)

// Load_Time_Chrono_ChronoZonedDateTime loads method signatures for java/time/chrono/ChronoZonedDateTime.
func Load_Time_Chrono_ChronoZonedDateTime() {
	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.compareTo(Ljava/time/chrono/ChronoZonedDateTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeEquals,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeFormat,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeFrom,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeGetChronology,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.getOffset()Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeGetOffset,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.getZone()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeGetZone,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeHashCode,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.isAfter(Ljava/time/chrono/ChronoZonedDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeIsAfter,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.isBefore(Ljava/time/chrono/ChronoZonedDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeIsBefore,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.isEqual(Ljava/time/chrono/ChronoZonedDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeIsEqual,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3, // long (2 slots) + TemporalUnit (1 slot)
			GFunction:  chronoZonedDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3, // long (2 slots) + TemporalUnit (1 slot)
			GFunction:  chronoZonedDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeQuery,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeRange,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.timeLineOrder()Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeTimelineOrder,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toEpochSecond()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toInstant()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeToInstant,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toLocalDate()Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toLocalDateTime()Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeToLocalDateTime,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toLocalTime()Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeToLocalTime,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeToString,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3, // field (1 slot) + long (2 slots)
			GFunction:  chronoZonedDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.withEarlierOffsetAtOverlap()Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeWithEarlierOffsetAtOverlap,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.withLaterOffsetAtOverlap()Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoZonedDateTimeWithLaterOffsetAtOverlap,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.withZoneSameInstant(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeWithZoneSameInstant,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.withZoneSameLocal(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoZonedDateTimeWithZoneSameLocal,
		}
}

// Helpers for ChronoZonedDateTime

func getChronoZonedDateTimeLocalDateTime(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if ldtFld, ok := obj.FieldTable["dateTime"]; ok {
		if ldtObj, ok := ldtFld.Fvalue.(*object.Object); ok {
			return ldtObj
		}
	}
	return nil
}

func getChronoZonedDateTimeOffset(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if offsetFld, ok := obj.FieldTable["offset"]; ok {
		if offsetObj, ok := offsetFld.Fvalue.(*object.Object); ok {
			return offsetObj
		}
	}
	return nil
}

func getChronoZonedDateTimeZone(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if zoneFld, ok := obj.FieldTable["zone"]; ok {
		if zoneObj, ok := zoneFld.Fvalue.(*object.Object); ok {
			return zoneObj
		}
	}
	return nil
}

func getZoneIdString(zoneObj *object.Object) string {
	if zoneObj == nil || zoneObj.FieldTable == nil {
		return "Z"
	}
	if idFld, ok := zoneObj.FieldTable["id"]; ok {
		if strObj, ok := idFld.Fvalue.(*object.Object); ok {
			return object.GoStringFromStringObject(strObj)
		}
		if s, ok := idFld.Fvalue.(string); ok {
			return s
		}
	}
	return "Z"
}

func createZoneOffset(totalSeconds int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameZoneOffset)
	obj.FieldTable["totalSeconds"] = object.Field{Ftype: types.Int, Fvalue: totalSeconds}
	var idStr string
	if totalSeconds == 0 {
		idStr = "Z"
	} else {
		absSec := totalSeconds
		sign := "+"
		if absSec < 0 {
			sign = "-"
			absSec = -absSec
		}
		hours := absSec / 3600
		minutes := (absSec % 3600) / 60
		secs := absSec % 60
		if secs == 0 {
			idStr = fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
		} else {
			idStr = fmt.Sprintf("%s%02d:%02d:%02d", sign, hours, minutes, secs)
		}
	}
	obj.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString(idStr)}
	return obj
}

func createZoneId(id string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameZoneId)
	obj.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString(id)}
	return obj
}

func getChronoZonedDateTimeYMDHMSN(obj *object.Object) (year, month, day, hour, min, sec, nano int64) {
	ldt := getChronoZonedDateTimeLocalDateTime(obj)
	if ldt != nil {
		dObj := getChronoLocalDateTimeDate(ldt)
		tObj := getChronoLocalDateTimeTime(ldt)
		year, month, day = ghelpers.GetChronoLocalDateYMD(dObj)
		hour, min, sec, nano = getLocalTimeHMSN(tObj)
		return
	}
	if obj != nil && obj.FieldTable != nil {
		if yFld, ok := obj.FieldTable["year"]; ok {
			year, _ = getInt64Val(yFld.Fvalue)
		}
		if mFld, ok := obj.FieldTable["month"]; ok {
			month, _ = getInt64Val(mFld.Fvalue)
		}
		if dFld, ok := obj.FieldTable["day"]; ok {
			day, _ = getInt64Val(dFld.Fvalue)
		}
		if hFld, ok := obj.FieldTable["hour"]; ok {
			hour, _ = getInt64Val(hFld.Fvalue)
		}
		if minFld, ok := obj.FieldTable["minute"]; ok {
			min, _ = getInt64Val(minFld.Fvalue)
		}
		if sFld, ok := obj.FieldTable["second"]; ok {
			sec, _ = getInt64Val(sFld.Fvalue)
		}
		if nFld, ok := obj.FieldTable["nano"]; ok {
			nano, _ = getInt64Val(nFld.Fvalue)
		}
	}
	if year == 0 && month == 0 && day == 0 {
		year, month, day = 1970, 1, 1
	}
	return
}

func getChronoZonedDateTimeEpochSecond(obj *object.Object) int64 {
	if obj == nil || obj.FieldTable == nil {
		return 0
	}
	if epFld, ok := obj.FieldTable["epochSecond"]; ok {
		if val, ok := getInt64Val(epFld.Fvalue); ok {
			return val
		}
	}
	year, month, day, hour, min, sec, _ := getChronoZonedDateTimeYMDHMSN(obj)
	offsetObj := getChronoZonedDateTimeOffset(obj)
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)
	t := time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), 0, time.UTC)
	return t.Unix() - offsetSec
}

func getChronoZonedDateTimeChronology(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if chronoFld, ok := obj.FieldTable["chronology"]; ok {
		if chronoObj, ok := chronoFld.Fvalue.(*object.Object); ok {
			return chronoObj
		}
	}
	ldt := getChronoZonedDateTimeLocalDateTime(obj)
	if ldt != nil {
		dObj := getChronoLocalDateTimeDate(ldt)
		return GetChronoLocalDateChronology(dObj)
	}
	return nil
}

func createChronoZonedDateTimeFromLDT(ldt *object.Object, zone *object.Object, offset *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoZonedDateTime)
	if zone == nil || object.IsNull(zone) {
		zone = createZoneId("UTC")
	}
	if offset == nil || object.IsNull(offset) {
		offset = createZoneOffset(0)
	}
	dObj := getChronoLocalDateTimeDate(ldt)
	tObj := getChronoLocalDateTimeTime(ldt)
	year, month, day := ghelpers.GetChronoLocalDateYMD(dObj)
	hour, min, sec, nano := getLocalTimeHMSN(tObj)
	chrono := GetChronoLocalDateChronology(dObj)
	offsetSec := getZoneOffsetTotalSeconds(offset)
	t := time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), 0, time.UTC)
	epochSecond := t.Unix() - offsetSec

	obj.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zone}
	obj.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offset}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: chrono}
	obj.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: year}
	obj.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: month}
	obj.FieldTable["day"] = object.Field{Ftype: types.Int, Fvalue: day}
	obj.FieldTable["hour"] = object.Field{Ftype: types.Int, Fvalue: hour}
	obj.FieldTable["minute"] = object.Field{Ftype: types.Int, Fvalue: min}
	obj.FieldTable["second"] = object.Field{Ftype: types.Int, Fvalue: sec}
	obj.FieldTable["nano"] = object.Field{Ftype: types.Int, Fvalue: nano}
	obj.FieldTable["epochSecond"] = object.Field{Ftype: types.Long, Fvalue: epochSecond}
	return obj
}

func makeChronoZonedDateTime(chrono *object.Object, year, month, day, hour, min, sec, nano int64, zone *object.Object, offset *object.Object) *object.Object {
	dateObj := createChronoLocalDate(chrono, year, month, day)
	timeObj := createLocalTime(hour, min, sec, nano)
	ldt := createChronoLocalDateTime(dateObj, timeObj)
	return createChronoZonedDateTimeFromLDT(ldt, zone, offset)
}

// chronoZonedDateTimeCompareTo compares this date-time to another date-time, including the chronology.
// default int compareTo(ChronoZonedDateTime<?> other)
func chronoZonedDateTimeCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	ep1 := getChronoZonedDateTimeEpochSecond(self)
	ep2 := getChronoZonedDateTimeEpochSecond(other)
	if ep1 < ep2 {
		return int64(-1)
	} else if ep1 > ep2 {
		return int64(1)
	}

	_, _, _, _, _, _, nano1 := getChronoZonedDateTimeYMDHMSN(self)
	_, _, _, _, _, _, nano2 := getChronoZonedDateTimeYMDHMSN(other)
	if nano1 < nano2 {
		return int64(-1)
	} else if nano1 > nano2 {
		return int64(1)
	}

	ldt1 := getChronoZonedDateTimeLocalDateTime(self)
	ldt2 := getChronoZonedDateTimeLocalDateTime(other)
	if ldt1 != nil && ldt2 != nil {
		cmpLdt := chronoLocalDateTimeCompareTo([]any{ldt1, ldt2})
		if cmpVal, ok := getInt64Val(cmpLdt); ok && cmpVal != 0 {
			return cmpVal
		}
	}

	zone1 := getZoneIdString(getChronoZonedDateTimeZone(self))
	zone2 := getZoneIdString(getChronoZonedDateTimeZone(other))
	if zone1 < zone2 {
		return int64(-1)
	} else if zone1 > zone2 {
		return int64(1)
	}

	chrono1 := getChronoZonedDateTimeChronology(self)
	chrono2 := getChronoZonedDateTimeChronology(other)
	id1 := getChronologyId(chrono1)
	id2 := getChronologyId(chrono2)
	if id1 < id2 {
		return int64(-1)
	} else if id1 > id2 {
		return int64(1)
	}

	return int64(0)
}

// chronoZonedDateTimeEquals checks if this date-time is equal to another date-time.
// boolean equals(Object obj)
func chronoZonedDateTimeEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return false
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return false
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	res := chronoZonedDateTimeCompareTo([]any{self, other})
	if cmpVal, ok := getInt64Val(res); ok {
		return cmpVal == 0
	}
	return false
}

// chronoZonedDateTimeFormat formats this date-time using the specified formatter.
// default String format(DateTimeFormatter formatter)
func chronoZonedDateTimeFormat(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.format: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.format: formatter is null")
	}
	return chronoZonedDateTimeToString(params)
}

// chronoZonedDateTimeFrom obtains an instance of ChronoZonedDateTime from a temporal object.
// static ChronoZonedDateTime<?> from(TemporalAccessor temporal)
func chronoZonedDateTimeFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.from: temporal is null")
	}
	temporal := params[0].(*object.Object)
	if object.IsObjectClass(temporal, classNameChronoZonedDateTime) {
		return temporal
	}

	chrono := GetChronoLocalDateChronology(temporal)
	if chrono == nil {
		chrono = object.MakeEmptyObjectWithClassName(&classNameIsoChronology)
		chrono.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("ISO")}
	}

	var year int64 = 1970
	var month int64 = 1
	var day int64 = 1
	var hour int64 = 0
	var min int64 = 0
	var sec int64 = 0
	var nano int64 = 0
	var zoneObj *object.Object
	var offsetObj *object.Object

	if temporal.FieldTable != nil {
		if yFld, ok := temporal.FieldTable["year"]; ok {
			year, _ = getInt64Val(yFld.Fvalue)
		}
		if mFld, ok := temporal.FieldTable["month"]; ok {
			month, _ = getInt64Val(mFld.Fvalue)
		}
		if dFld, ok := temporal.FieldTable["day"]; ok {
			day, _ = getInt64Val(dFld.Fvalue)
		}
		if hFld, ok := temporal.FieldTable["hour"]; ok {
			hour, _ = getInt64Val(hFld.Fvalue)
		}
		if minFld, ok := temporal.FieldTable["minute"]; ok {
			min, _ = getInt64Val(minFld.Fvalue)
		}
		if sFld, ok := temporal.FieldTable["second"]; ok {
			sec, _ = getInt64Val(sFld.Fvalue)
		}
		if nFld, ok := temporal.FieldTable["nano"]; ok {
			nano, _ = getInt64Val(nFld.Fvalue)
		}
		if zFld, ok := temporal.FieldTable["zone"]; ok {
			zoneObj, _ = zFld.Fvalue.(*object.Object)
		}
		if offFld, ok := temporal.FieldTable["offset"]; ok {
			offsetObj, _ = offFld.Fvalue.(*object.Object)
		}
		if ldtFld, ok := temporal.FieldTable["dateTime"]; ok {
			if ldtObj, ok := ldtFld.Fvalue.(*object.Object); ok {
				dObj := getChronoLocalDateTimeDate(ldtObj)
				tObj := getChronoLocalDateTimeTime(ldtObj)
				year, month, day = ghelpers.GetChronoLocalDateYMD(dObj)
				hour, min, sec, nano = getLocalTimeHMSN(tObj)
			}
		}
	}

	if zoneObj == nil {
		zoneObj = createZoneId("UTC")
	}
	if offsetObj == nil {
		offsetObj = createZoneOffset(0)
	}

	return makeChronoZonedDateTime(chrono, year, month, day, hour, min, sec, nano, zoneObj, offsetObj)
}

// chronoZonedDateTimeGetChronology gets the chronology of this date-time.
// default Chronology getChronology()
func chronoZonedDateTimeGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.getChronology: self is null")
	}
	self := params[0].(*object.Object)
	chrono := getChronoZonedDateTimeChronology(self)
	if chrono == nil {
		chrono = object.MakeEmptyObjectWithClassName(&classNameIsoChronology)
		chrono.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("ISO")}
	}
	return chrono
}

// chronoZonedDateTimeGetOffset gets the zone offset, such as '+01:00'.
// ZoneOffset getOffset()
func chronoZonedDateTimeGetOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.getOffset: self is null")
	}
	self := params[0].(*object.Object)
	offset := getChronoZonedDateTimeOffset(self)
	if offset == nil {
		offset = createZoneOffset(0)
	}
	return offset
}

// chronoZonedDateTimeGetZone gets the zone ID, such as 'Europe/Paris'.
// ZoneId getZone()
func chronoZonedDateTimeGetZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.getZone: self is null")
	}
	self := params[0].(*object.Object)
	zone := getChronoZonedDateTimeZone(self)
	if zone == nil {
		zone = createZoneId("UTC")
	}
	return zone
}

// chronoZonedDateTimeHashCode a hash code for this date-time.
// int hashCode()
func chronoZonedDateTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getChronoZonedDateTimeLocalDateTime(self)
	var ldtHash int64 = 0
	if ldt != nil {
		if h, ok := getInt64Val(chronoLocalDateTimeHashCode([]any{ldt})); ok {
			ldtHash = h
		}
	}
	offsetObj := getChronoZonedDateTimeOffset(self)
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)
	zoneObj := getChronoZonedDateTimeZone(self)
	zoneStr := getZoneIdString(zoneObj)
	var zoneHash int64 = 0
	for _, c := range zoneStr {
		zoneHash = 31*zoneHash + int64(c)
	}
	return ldtHash ^ offsetSec ^ zoneHash
}

// chronoZonedDateTimeIsAfter checks if the instant of this date-time is after that of the specified date-time.
// default boolean isAfter(ChronoZonedDateTime<?> other)
func chronoZonedDateTimeIsAfter(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.isAfter: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.isAfter: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	ep1 := getChronoZonedDateTimeEpochSecond(self)
	ep2 := getChronoZonedDateTimeEpochSecond(other)
	if ep1 > ep2 {
		return true
	} else if ep1 < ep2 {
		return false
	}

	_, _, _, _, _, _, nano1 := getChronoZonedDateTimeYMDHMSN(self)
	_, _, _, _, _, _, nano2 := getChronoZonedDateTimeYMDHMSN(other)
	return nano1 > nano2
}

// chronoZonedDateTimeIsBefore checks if the instant of this date-time is before that of the specified date-time.
// default boolean isBefore(ChronoZonedDateTime<?> other)
func chronoZonedDateTimeIsBefore(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.isBefore: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.isBefore: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	ep1 := getChronoZonedDateTimeEpochSecond(self)
	ep2 := getChronoZonedDateTimeEpochSecond(other)
	if ep1 < ep2 {
		return true
	} else if ep1 > ep2 {
		return false
	}

	_, _, _, _, _, _, nano1 := getChronoZonedDateTimeYMDHMSN(self)
	_, _, _, _, _, _, nano2 := getChronoZonedDateTimeYMDHMSN(other)
	return nano1 < nano2
}

// chronoZonedDateTimeIsEqual checks if the instant of this date-time is equal to that of the specified date-time.
// default boolean isEqual(ChronoZonedDateTime<?> other)
func chronoZonedDateTimeIsEqual(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.isEqual: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.isEqual: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	ep1 := getChronoZonedDateTimeEpochSecond(self)
	ep2 := getChronoZonedDateTimeEpochSecond(other)
	if ep1 != ep2 {
		return false
	}

	_, _, _, _, _, _, nano1 := getChronoZonedDateTimeYMDHMSN(self)
	_, _, _, _, _, _, nano2 := getChronoZonedDateTimeYMDHMSN(other)
	return nano1 == nano2
}

// chronoZonedDateTimeIsSupportedField checks if the specified field is supported.
// boolean isSupported(TemporalField field)
func chronoZonedDateTimeIsSupportedField(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return false
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return false
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "" {
		return false
	}

	switch fieldName {
	case "INSTANT_SECONDS", "OFFSET_SECONDS":
		return true
	default:
		self := params[0].(*object.Object)
		ldt := getChronoZonedDateTimeLocalDateTime(self)
		if ldt != nil {
			return chronoLocalDateTimeIsSupportedField([]any{ldt, field})
		}
		return true
	}
}

// chronoZonedDateTimeIsSupportedUnit checks if the specified unit is supported.
// default boolean isSupported(TemporalUnit unit)
func chronoZonedDateTimeIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return false
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return false
	}
	unit := params[1].(*object.Object)
	unitName := getTemporalUnitName(unit)
	if unitName == "" {
		return false
	}

	switch unitName {
	case "NANOS", "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS",
		"DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return true
	default:
		return false
	}
}

// chronoZonedDateTimeMinusUnit returns an object of the same type as this object with the specified period subtracted.
// default ChronoZonedDateTime<D> minus(long amountToSubtract, TemporalUnit unit)
func chronoZonedDateTimeMinusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.minus: self is null")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoZonedDateTime.minus: invalid amount type")
	}
	unit := params[2]
	return chronoZonedDateTimePlusUnit([]any{params[0], -amount, unit})
}

// chronoZonedDateTimeMinusAmount returns an object of the same type as this object with an amount subtracted.
// default ChronoZonedDateTime<D> minus(TemporalAmount amount)
func chronoZonedDateTimeMinusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		newLdt := chronoLocalDateTimeMinusAmount([]any{ldt, amount})
		if newLdtObj, ok := newLdt.(*object.Object); ok {
			zone := getChronoZonedDateTimeZone(self)
			offset := getChronoZonedDateTimeOffset(self)
			return createChronoZonedDateTimeFromLDT(newLdtObj, zone, offset)
		}
		return newLdt
	}
	return self
}

// chronoZonedDateTimePlusUnit returns an object of the same type as this object with the specified period added.
// ChronoZonedDateTime<D> plus(long amountToAdd, TemporalUnit unit)
func chronoZonedDateTimePlusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.plus: self is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.plus: unit is null")
	}
	amountToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoZonedDateTime.plus: invalid amountToAdd type")
	}
	self := params[0].(*object.Object)
	unit := params[2].(*object.Object)

	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		newLdt := chronoLocalDateTimePlusUnit([]any{ldt, amountToAdd, unit})
		if newLdtObj, ok := newLdt.(*object.Object); ok {
			zone := getChronoZonedDateTimeZone(self)
			offset := getChronoZonedDateTimeOffset(self)
			return createChronoZonedDateTimeFromLDT(newLdtObj, zone, offset)
		}
		return newLdt
	}
	return self
}

// chronoZonedDateTimePlusAmount returns an object of the same type as this object with an amount added.
// default ChronoZonedDateTime<D> plus(TemporalAmount amount)
func chronoZonedDateTimePlusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		newLdt := chronoLocalDateTimePlusAmount([]any{ldt, amount})
		if newLdtObj, ok := newLdt.(*object.Object); ok {
			zone := getChronoZonedDateTimeZone(self)
			offset := getChronoZonedDateTimeOffset(self)
			return createChronoZonedDateTimeFromLDT(newLdtObj, zone, offset)
		}
		return newLdt
	}
	return self
}

// chronoZonedDateTimeQuery queries this date-time using the specified query.
// default <R> R query(TemporalQuery<R> query)
func chronoZonedDateTimeQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.query: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	if query.FieldTable != nil {
		if nameFld, ok := query.FieldTable["name"]; ok {
			if strObj, ok := nameFld.Fvalue.(*object.Object); ok {
				qName := object.GoStringFromStringObject(strObj)
				switch qName {
				case "zone", "zoneId":
					return chronoZonedDateTimeGetZone([]any{self})
				case "chronology":
					return chronoZonedDateTimeGetChronology([]any{self})
				case "precision":
					return ghelpers.CreateTemporalUnit("NANOS")
				case "offset":
					return chronoZonedDateTimeGetOffset([]any{self})
				case "localDate":
					return chronoZonedDateTimeToLocalDate([]any{self})
				case "localTime":
					return chronoZonedDateTimeToLocalTime([]any{self})
				}
			}
		}
	}

	return temporalQueryQueryFrom([]any{query, self})
}

// chronoZonedDateTimeRange gets the range of valid values for the specified field.
// default ValueRange range(TemporalField field)
func chronoZonedDateTimeRange(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.range: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.range: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)

	fieldName := getTemporalFieldName(field)
	switch fieldName {
	case "INSTANT_SECONDS":
		return createValueRange(math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64)
	case "OFFSET_SECONDS":
		return createValueRange(-18*3600, -18*3600, 18*3600, 18*3600)
	default:
		ldt := getChronoZonedDateTimeLocalDateTime(self)
		if ldt != nil {
			return temporalAccessorRange([]any{ldt, field})
		}
		return createValueRange(math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64)
	}
}

// chronoZonedDateTimeTimelineOrder gets a comparator that compares ChronoZonedDateTime in time-line order ignoring the chronology.
// static Comparator<ChronoZonedDateTime<?>> timeLineOrder()
func chronoZonedDateTimeTimelineOrder(params []any) any {
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoZonedDateTimeComparator)
	return obj
}

// chronoZonedDateTimeToEpochSecond converts this date-time to the number of seconds from the epoch of 1970-01-01T00:00:00Z.
// default long toEpochSecond()
func chronoZonedDateTimeToEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.toEpochSecond: self is null")
	}
	self := params[0].(*object.Object)
	return getChronoZonedDateTimeEpochSecond(self)
}

// chronoZonedDateTimeToInstant converts this date-time to an Instant.
// default Instant toInstant()
func chronoZonedDateTimeToInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.toInstant: self is null")
	}
	self := params[0].(*object.Object)
	epochSec := getChronoZonedDateTimeEpochSecond(self)
	_, _, _, _, _, _, nano := getChronoZonedDateTimeYMDHMSN(self)

	obj := object.MakeEmptyObjectWithClassName(&classNameInstant)
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: nano}
	return obj
}

// chronoZonedDateTimeToLocalDate gets the local date part of this date-time.
// default D toLocalDate()
func chronoZonedDateTimeToLocalDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.toLocalDate: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		return chronoLocalDateTimeToLocalDate([]any{ldt})
	}
	year, month, day, _, _, _, _ := getChronoZonedDateTimeYMDHMSN(self)
	chrono := getChronoZonedDateTimeChronology(self)
	return createChronoLocalDate(chrono, year, month, day)
}

// chronoZonedDateTimeToLocalDateTime gets the local date-time part of this date-time.
// ChronoLocalDateTime<D> toLocalDateTime()
func chronoZonedDateTimeToLocalDateTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.toLocalDateTime: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		return ldt
	}
	year, month, day, hour, min, sec, nano := getChronoZonedDateTimeYMDHMSN(self)
	chrono := getChronoZonedDateTimeChronology(self)
	dObj := createChronoLocalDate(chrono, year, month, day)
	tObj := createLocalTime(hour, min, sec, nano)
	return createChronoLocalDateTime(dObj, tObj)
}

// chronoZonedDateTimeToLocalTime gets the local time part of this date-time.
// default LocalTime toLocalTime()
func chronoZonedDateTimeToLocalTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.toLocalTime: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		return chronoLocalDateTimeToLocalTime([]any{ldt})
	}
	_, _, _, hour, min, sec, nano := getChronoZonedDateTimeYMDHMSN(self)
	return createLocalTime(hour, min, sec, nano)
}

// chronoZonedDateTimeToString outputs this date-time as a String.
// String toString()
func chronoZonedDateTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.toString: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getChronoZonedDateTimeLocalDateTime(self)
	var ldtStr string
	if ldt != nil {
		ldtStrObj := chronoLocalDateTimeToString([]any{ldt})
		if strObj, ok := ldtStrObj.(*object.Object); ok {
			ldtStr = object.GoStringFromStringObject(strObj)
		}
	}
	if ldtStr == "" {
		year, month, day, hour, min, sec, nano := getChronoZonedDateTimeYMDHMSN(self)
		if nano > 0 {
			ldtStr = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%09d", year, month, day, hour, min, sec, nano)
			ldtStr = strings.TrimRight(ldtStr, "0")
		} else {
			ldtStr = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, min, sec)
		}
	}

	offsetObj := getChronoZonedDateTimeOffset(self)
	offsetStr := getZoneIdString(offsetObj)
	zoneObj := getChronoZonedDateTimeZone(self)
	zoneStr := getZoneIdString(zoneObj)

	res := ldtStr + offsetStr
	if zoneStr != "" && zoneStr != offsetStr {
		res += "[" + zoneStr + "]"
	}
	return object.StringObjectFromGoString(res)
}

// chronoZonedDateTimeWithAdjuster returns an adjusted copy of this date-time.
// default ChronoZonedDateTime<D> with(TemporalAdjuster adjuster)
func chronoZonedDateTimeWithAdjuster(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	ldt := getChronoZonedDateTimeLocalDateTime(self)
	if ldt != nil {
		newLdt := chronoLocalDateTimeWithAdjuster([]any{ldt, adjuster})
		if newLdtObj, ok := newLdt.(*object.Object); ok {
			zone := getChronoZonedDateTimeZone(self)
			offset := getChronoZonedDateTimeOffset(self)
			return createChronoZonedDateTimeFromLDT(newLdtObj, zone, offset)
		}
		return newLdt
	}
	return self
}

// chronoZonedDateTimeWithField returns a copy of this object with the specified field set to a new value.
// default ChronoZonedDateTime<D> with(TemporalField field, long newValue)
func chronoZonedDateTimeWithField(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.with: field is null")
	}
	newValue, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoZonedDateTime.with: invalid newValue type")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)

	fieldName := getTemporalFieldName(field)
	switch fieldName {
	case "INSTANT_SECONDS":
		_, _, _, _, _, _, nano := getChronoZonedDateTimeYMDHMSN(self)
		zone := getChronoZonedDateTimeZone(self)
		offset := getChronoZonedDateTimeOffset(self)
		chrono := getChronoZonedDateTimeChronology(self)
		offsetSec := getZoneOffsetTotalSeconds(offset)
		localEpoch := newValue + offsetSec
		t := time.Unix(localEpoch, 0).UTC()
		return makeChronoZonedDateTime(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()), int64(t.Hour()), int64(t.Minute()), int64(t.Second()), nano, zone, offset)
	case "OFFSET_SECONDS":
		newOffset := createZoneOffset(newValue)
		ldt := getChronoZonedDateTimeLocalDateTime(self)
		zone := getChronoZonedDateTimeZone(self)
		return createChronoZonedDateTimeFromLDT(ldt, zone, newOffset)
	default:
		ldt := getChronoZonedDateTimeLocalDateTime(self)
		if ldt != nil {
			newLdt := chronoLocalDateTimeWithField([]any{ldt, field, newValue})
			if newLdtObj, ok := newLdt.(*object.Object); ok {
				zone := getChronoZonedDateTimeZone(self)
				offset := getChronoZonedDateTimeOffset(self)
				return createChronoZonedDateTimeFromLDT(newLdtObj, zone, offset)
			}
			return newLdt
		}
		return self
	}
}

// chronoZonedDateTimeWithEarlierOffsetAtOverlap returns a copy of this date-time with the earlier ZoneOffset at overlaps.
// ChronoZonedDateTime<D> withEarlierOffsetAtOverlap()
func chronoZonedDateTimeWithEarlierOffsetAtOverlap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.withEarlierOffsetAtOverlap: self is null")
	}
	return params[0]
}

// chronoZonedDateTimeWithLaterOffsetAtOverlap returns a copy of this date-time with the later ZoneOffset at overlaps.
// ChronoZonedDateTime<D> withLaterOffsetAtOverlap()
func chronoZonedDateTimeWithLaterOffsetAtOverlap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.withLaterOffsetAtOverlap: self is null")
	}
	return params[0]
}

// chronoZonedDateTimeWithZoneSameInstant returns a copy of this date-time with a different time-zone, retaining the instant.
// ChronoZonedDateTime<D> withZoneSameInstant(ZoneId zone)
func chronoZonedDateTimeWithZoneSameInstant(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.withZoneSameInstant: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.withZoneSameInstant: zone is null")
	}
	self := params[0].(*object.Object)
	newZone := params[1].(*object.Object)

	epochSec := getChronoZonedDateTimeEpochSecond(self)
	_, _, _, _, _, _, nano := getChronoZonedDateTimeYMDHMSN(self)
	chrono := getChronoZonedDateTimeChronology(self)

	offsetSec := getZoneOffsetTotalSeconds(newZone)
	localEpoch := epochSec + offsetSec
	t := time.Unix(localEpoch, 0).UTC()
	newOffset := createZoneOffset(offsetSec)

	return makeChronoZonedDateTime(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()), int64(t.Hour()), int64(t.Minute()), int64(t.Second()), nano, newZone, newOffset)
}

// chronoZonedDateTimeWithZoneSameLocal returns a copy of this date-time with a different time-zone, retaining the local date-time if possible.
// ChronoZonedDateTime<D> withZoneSameLocal(ZoneId zone)
func chronoZonedDateTimeWithZoneSameLocal(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.withZoneSameLocal: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoZonedDateTime.withZoneSameLocal: zone is null")
	}
	self := params[0].(*object.Object)
	newZone := params[1].(*object.Object)

	ldt := getChronoZonedDateTimeLocalDateTime(self)
	offsetSec := getZoneOffsetTotalSeconds(newZone)
	newOffset := createZoneOffset(offsetSec)

	return createChronoZonedDateTimeFromLDT(ldt, newZone, newOffset)
}

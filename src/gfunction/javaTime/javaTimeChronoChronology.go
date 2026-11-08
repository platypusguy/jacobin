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
	classNameIsoChronology          = "java/time/chrono/IsoChronology"
	classNameHijrahChronology       = "java/time/chrono/HijrahChronology"
	classNameJapaneseChronology     = "java/time/chrono/JapaneseChronology"
	classNameMinguoChronology       = "java/time/chrono/MinguoChronology"
	classNameThaiBuddhistChronology = "java/time/chrono/ThaiBuddhistChronology"
	classNameChronoLocalDate        = "java/time/chrono/ChronoLocalDate"
	classNameChronoLocalDateTime    = "java/time/chrono/ChronoLocalDateTime"
	classNameChronoZonedDateTime    = "java/time/chrono/ChronoZonedDateTime"
	classNameChronoPeriod           = "java/time/chrono/ChronoPeriod"
	classNameIsoEra                 = "java/time/chrono/IsoEra"
	classNameHashSet                = "java/util/HashSet"
)

// Load_Time_Chrono_Chronology loads method signatures for java/time/chrono/Chronology.
func Load_Time_Chrono_Chronology() {
	ghelpers.MethodSignatures["java/time/chrono/Chronology.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.compareTo(Ljava/time/chrono/Chronology;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.date(III)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyDateYMD,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.date(Ljava/time/chrono/Era;III)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  chronologyDateEra,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.date(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateTemporal,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.dateEpochDay(J)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyDateEpochDay,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.dateNow()Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyDateNow,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.dateNow(Ljava/time/Clock;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateNowClock,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.dateNow(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateNowZone,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.dateYearDay(II)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyDateYearDay,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.dateYearDay(Ljava/time/chrono/Era;II)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyDateYearDayEra,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.epochSecond(IIIIIILjava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  chronologyEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.epochSecond(Ljava/time/chrono/Era;IIIIIILjava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 8,
			GFunction:  chronologyEpochSecondEra,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyEquals,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.eraOf(I)Ljava/time/chrono/Era;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyEraOf,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.eras()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyEras,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyFrom,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.getAvailableChronologies()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyGetAvailableChronologies,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.getCalendarType()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyGetCalendarType,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.getDisplayName(Ljava/time/format/TextStyle;Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.getId()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyGetId,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyHashCode,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.isLeapYear(J)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyIsLeapYear,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.localDateTime(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyLocalDateTime,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.of(Ljava/lang/String;)Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyOf,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.ofLocale(Ljava/util/Locale;)Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyOfLocale,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.period(III)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyPeriod,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.prolepticYear(Ljava/time/chrono/Era;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyProlepticYear,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.range(Ljava/time/temporal/ChronoField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyRange,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.resolveDate(Ljava/util/Map;Ljava/time/format/ResolverStyle;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyResolveDate,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyToString,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.zonedDateTime(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyZonedDateTimeInstant,
		}

	ghelpers.MethodSignatures["java/time/chrono/Chronology.zonedDateTime(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyZonedDateTimeTemporal,
		}
}

func GetIsoChronology() *object.Object {
	return ghelpers.CreateChronology(classNameIsoChronology, "ISO", "iso8601")
}

func getHijrahChronology() *object.Object {
	return ghelpers.CreateChronology(classNameHijrahChronology, "Hijrah-umalqura", "islamic-umalqura")
}

func getJapaneseChronology() *object.Object {
	return ghelpers.CreateChronology(classNameJapaneseChronology, "Japanese", "japanese")
}

func getMinguoChronology() *object.Object {
	return ghelpers.CreateChronology(classNameMinguoChronology, "Minguo", "roc")
}

func getThaiBuddhistChronology() *object.Object {
	return ghelpers.CreateChronology(classNameThaiBuddhistChronology, "ThaiBuddhist", "buddhist")
}

func getChronologyId(obj *object.Object) string {
	if obj == nil || obj.FieldTable == nil {
		return "ISO"
	}
	if idFld, ok := obj.FieldTable["id"]; ok {
		if idObj, ok := idFld.Fvalue.(*object.Object); ok && object.IsStringObject(idObj) {
			return object.GoStringFromStringObject(idObj)
		}
	}
	return "ISO"
}

func getChronologyCalendarType(obj *object.Object) string {
	if obj == nil || obj.FieldTable == nil {
		return "iso8601"
	}
	if calFld, ok := obj.FieldTable["calendarType"]; ok {
		if calObj, ok := calFld.Fvalue.(*object.Object); ok && object.IsStringObject(calObj) {
			return object.GoStringFromStringObject(calObj)
		}
	}
	return "iso8601"
}

func createEra(eraValue int64, name string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameIsoEra)
	nameObj := object.StringObjectFromGoString(name)
	obj.FieldTable["eraValue"] = object.Field{Ftype: types.Int, Fvalue: eraValue}
	obj.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: nameObj}
	return obj
}

func createChronoLocalDate(chrono *object.Object, year, month, day int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoLocalDate)
	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	epochDay := t.Unix() / 86400
	obj.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: year}
	obj.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: month}
	obj.FieldTable["day"] = object.Field{Ftype: types.Int, Fvalue: day}
	obj.FieldTable["epochDay"] = object.Field{Ftype: types.Long, Fvalue: epochDay}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: chrono}
	return obj
}

// chronologyCompareTo compares this chronology to another.
// int compareTo(Chronology other)
func chronologyCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	id1 := getChronologyId(self)
	id2 := getChronologyId(other)
	if id1 < id2 {
		return int64(-1)
	} else if id1 > id2 {
		return int64(1)
	}
	return int64(0)
}

// chronologyDateYMD obtains a local date in this chronology from proleptic-year, month-of-year and day-of-month.
// ChronoLocalDate date(int prolepticYear, int month, int dayOfMonth)
func chronologyDateYMD(params []any) any {
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.date: self is null")
	}
	self := params[0].(*object.Object)
	year, ok1 := getInt64Val(params[1])
	month, ok2 := getInt64Val(params[2])
	day, ok3 := getInt64Val(params[3])
	if !ok1 || !ok2 || !ok3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.date: invalid argument types")
	}

	return createChronoLocalDate(self, year, month, day)
}

// chronologyDateEra obtains a local date in this chronology from the era, year-of-era, month-of-year and day-of-month.
// ChronoLocalDate date(Era era, int yearOfEra, int month, int dayOfMonth)
func chronologyDateEra(params []any) any {
	if len(params) < 5 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.date: self is null")
	}
	self := params[0].(*object.Object)
	eraObj, _ := params[1].(*object.Object)
	yearOfEra, ok1 := getInt64Val(params[2])
	month, ok2 := getInt64Val(params[3])
	day, ok3 := getInt64Val(params[4])
	if !ok1 || !ok2 || !ok3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.date: invalid argument types")
	}

	year := yearOfEra
	if eraObj != nil && !object.IsNull(eraObj) && eraObj.FieldTable != nil {
		if eraValFld, ok := eraObj.FieldTable["eraValue"]; ok {
			if eraVal, ok := getInt64Val(eraValFld.Fvalue); ok && eraVal == 0 {
				year = 1 - yearOfEra
			}
		}
	}

	return createChronoLocalDate(self, year, month, day)
}

// chronologyDateTemporal obtains a local date in this chronology from another temporal object.
// ChronoLocalDate date(TemporalAccessor temporal)
func chronologyDateTemporal(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.date: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.date: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	if temporal.FieldTable != nil {
		var year int64 = 1970
		var month int64 = 1
		var day int64 = 1

		if yFld, ok := temporal.FieldTable["year"]; ok {
			if yVal, ok := getInt64Val(yFld.Fvalue); ok {
				year = yVal
			}
		}
		if mFld, ok := temporal.FieldTable["month"]; ok {
			if mVal, ok := getInt64Val(mFld.Fvalue); ok {
				month = mVal
			}
		}
		if dFld, ok := temporal.FieldTable["day"]; ok {
			if dVal, ok := getInt64Val(dFld.Fvalue); ok {
				day = dVal
			}
		}
		if epochFld, ok := temporal.FieldTable["epochDay"]; ok {
			if epochDay, ok := getInt64Val(epochFld.Fvalue); ok {
				t := time.Unix(epochDay*86400, 0).UTC()
				year = int64(t.Year())
				month = int64(t.Month())
				day = int64(t.Day())
			}
		}

		return createChronoLocalDate(self, year, month, day)
	}

	return createChronoLocalDate(self, 1970, 1, 1)
}

// chronologyDateEpochDay obtains a local date in this chronology from the epoch-day.
// ChronoLocalDate dateEpochDay(long epochDay)
func chronologyDateEpochDay(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateEpochDay: self is null")
	}
	self := params[0].(*object.Object)
	epochDay, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.dateEpochDay: invalid epochDay type")
	}

	t := time.Unix(epochDay*86400, 0).UTC()
	return createChronoLocalDate(self, int64(t.Year()), int64(t.Month()), int64(t.Day()))
}

// chronologyDateNow obtains the current local date in this chronology from the system clock in the default time-zone.
// ChronoLocalDate dateNow()
func chronologyDateNow(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateNow: self is null")
	}
	self := params[0].(*object.Object)
	now := time.Now()
	return createChronoLocalDate(self, int64(now.Year()), int64(now.Month()), int64(now.Day()))
}

// chronologyDateNowClock obtains the current local date in this chronology from the specified clock.
// ChronoLocalDate dateNow(Clock clock)
func chronologyDateNowClock(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateNow: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateNow: clock is null")
	}
	self := params[0].(*object.Object)
	now := time.Now()
	return createChronoLocalDate(self, int64(now.Year()), int64(now.Month()), int64(now.Day()))
}

// chronologyDateNowZone obtains the current local date in this chronology from the system clock in the specified time-zone.
// ChronoLocalDate dateNow(ZoneId zone)
func chronologyDateNowZone(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateNow: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateNow: zone is null")
	}
	self := params[0].(*object.Object)
	now := time.Now()
	return createChronoLocalDate(self, int64(now.Year()), int64(now.Month()), int64(now.Day()))
}

// chronologyDateYearDay obtains a local date in this chronology from the proleptic-year and day-of-year.
// ChronoLocalDate dateYearDay(int prolepticYear, int dayOfYear)
func chronologyDateYearDay(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateYearDay: self is null")
	}
	self := params[0].(*object.Object)
	year, ok1 := getInt64Val(params[1])
	dayOfYear, ok2 := getInt64Val(params[2])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.dateYearDay: invalid argument types")
	}

	start := time.Date(int(year), 1, 1, 0, 0, 0, 0, time.UTC)
	target := start.AddDate(0, 0, int(dayOfYear-1))
	return createChronoLocalDate(self, int64(target.Year()), int64(target.Month()), int64(target.Day()))
}

// chronologyDateYearDayEra obtains a local date in this chronology from the era, year-of-era and day-of-year.
// ChronoLocalDate dateYearDay(Era era, int yearOfEra, int dayOfYear)
func chronologyDateYearDayEra(params []any) any {
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.dateYearDay: self is null")
	}
	self := params[0].(*object.Object)
	eraObj, _ := params[1].(*object.Object)
	yearOfEra, ok1 := getInt64Val(params[2])
	dayOfYear, ok2 := getInt64Val(params[3])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.dateYearDay: invalid argument types")
	}

	year := yearOfEra
	if eraObj != nil && !object.IsNull(eraObj) && eraObj.FieldTable != nil {
		if eraValFld, ok := eraObj.FieldTable["eraValue"]; ok {
			if eraVal, ok := getInt64Val(eraValFld.Fvalue); ok && eraVal == 0 {
				year = 1 - yearOfEra
			}
		}
	}

	start := time.Date(int(year), 1, 1, 0, 0, 0, 0, time.UTC)
	target := start.AddDate(0, 0, int(dayOfYear-1))
	return createChronoLocalDate(self, int64(target.Year()), int64(target.Month()), int64(target.Day()))
}

// chronologyEpochSecond calculates the epoch second from date-time fields and offset.
// long epochSecond(int prolepticYear, int month, int dayOfMonth, int hour, int minute, int second, ZoneOffset zoneOffset)
func chronologyEpochSecond(params []any) any {
	if len(params) < 8 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.epochSecond: self is null")
	}
	year, ok1 := getInt64Val(params[1])
	month, ok2 := getInt64Val(params[2])
	day, ok3 := getInt64Val(params[3])
	hour, ok4 := getInt64Val(params[4])
	minutes, ok5 := getInt64Val(params[5])
	sec, ok6 := getInt64Val(params[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.epochSecond: invalid argument types")
	}

	var offsetSec int64 = 0
	if offsetObj, ok := params[7].(*object.Object); ok && offsetObj != nil && !object.IsNull(offsetObj) {
		if offsetObj.FieldTable != nil {
			if totSecFld, ok := offsetObj.FieldTable["totalSeconds"]; ok {
				if tsVal, ok := getInt64Val(totSecFld.Fvalue); ok {
					offsetSec = tsVal
				}
			}
		}
	}

	t := time.Date(int(year), time.Month(month), int(day), int(hour), int(minutes), int(sec), 0, time.UTC)
	return t.Unix() - offsetSec
}

// chronologyEpochSecondEra calculates the epoch second from era, year-of-era, and date-time fields.
// long epochSecond(Era era, int yearOfEra, int month, int dayOfMonth, int hour, int minute, int second, ZoneOffset zoneOffset)
func chronologyEpochSecondEra(params []any) any {
	if len(params) < 9 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.epochSecond: self is null")
	}
	self := params[0].(*object.Object)
	eraObj, _ := params[1].(*object.Object)
	yearOfEra, ok1 := getInt64Val(params[2])
	month, ok2 := getInt64Val(params[3])
	day, ok3 := getInt64Val(params[4])
	hour, ok4 := getInt64Val(params[5])
	minutes, ok5 := getInt64Val(params[6])
	sec, ok6 := getInt64Val(params[7])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.epochSecond: invalid argument types")
	}

	year := yearOfEra
	if eraObj != nil && !object.IsNull(eraObj) && eraObj.FieldTable != nil {
		if eraValFld, ok := eraObj.FieldTable["eraValue"]; ok {
			if eraVal, ok := getInt64Val(eraValFld.Fvalue); ok && eraVal == 0 {
				year = 1 - yearOfEra
			}
		}
	}

	return chronologyEpochSecond([]any{self, year, month, day, hour, minutes, sec, params[8]})
}

// chronologyEquals checks if this chronology is equal to another.
// boolean equals(Object obj)
func chronologyEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.equals: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	id1 := getChronologyId(self)
	id2 := getChronologyId(other)
	if id1 == id2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronologyEraOf creates the Era object for this chronology.
// Era eraOf(int eraValue)
func chronologyEraOf(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.eraOf: self is null")
	}
	eraVal, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.eraOf: invalid eraValue type")
	}

	name := "CE"
	if eraVal == 0 {
		name = "BCE"
	}
	return createEra(eraVal, name)
}

// chronologyEras gets the list of eras for the chronology.
// List<Era> eras()
func chronologyEras(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.eras: self is null")
	}

	listObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	bce := createEra(0, "BCE")
	ce := createEra(1, "CE")
	erasList := []any{bce, ce}
	listObj.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: erasList,
	}
	return listObj
}

// chronologyFrom obtains an instance of Chronology from a temporal object.
// static Chronology from(TemporalAccessor temporal)
func chronologyFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.from: temporal is null")
	}
	temporal := params[0].(*object.Object)

	if temporal.FieldTable != nil {
		if chronoFld, ok := temporal.FieldTable["chronology"]; ok {
			if chronoObj, ok := chronoFld.Fvalue.(*object.Object); ok && !object.IsNull(chronoObj) {
				return chronoObj
			}
		}
	}
	return GetIsoChronology()
}

// chronologyGetAvailableChronologies returns the available chronologies.
// static Set<Chronology> getAvailableChronologies()
func chronologyGetAvailableChronologies(params []any) any {
	setObj := object.MakeEmptyObjectWithClassName(&classNameHashSet)
	iso := GetIsoChronology()
	hijrah := getHijrahChronology()
	japanese := getJapaneseChronology()
	minguo := getMinguoChronology()
	thai := getThaiBuddhistChronology()

	chronoSet := []any{iso, hijrah, japanese, minguo, thai}
	setObj.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: chronoSet,
	}
	return setObj
}

// chronologyGetCalendarType gets the calendar type of the underlying calendar system.
// String getCalendarType()
func chronologyGetCalendarType(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.getCalendarType: self is null")
	}
	self := params[0].(*object.Object)
	calType := getChronologyCalendarType(self)
	return object.StringObjectFromGoString(calType)
}

// chronologyGetDisplayName gets the textual representation of this chronology.
// String getDisplayName(TextStyle style, Locale locale)
func chronologyGetDisplayName(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.getDisplayName: self is null")
	}
	self := params[0].(*object.Object)
	id := getChronologyId(self)
	return object.StringObjectFromGoString(id)
}

// chronologyGetId gets the ID of the chronology.
// String getId()
func chronologyGetId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.getId: self is null")
	}
	self := params[0].(*object.Object)
	id := getChronologyId(self)
	return object.StringObjectFromGoString(id)
}

// chronologyHashCode returns a hash code for this chronology.
// int hashCode()
func chronologyHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	id := getChronologyId(self)
	var hash int64
	for _, c := range id {
		hash = 31*hash + int64(c)
	}
	return hash
}

// chronologyIsLeapYear checks if the specified year is a leap year.
// boolean isLeapYear(long prolepticYear)
func chronologyIsLeapYear(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.isLeapYear: self is null")
	}
	year, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.isLeapYear: invalid year type")
	}

	isLeap := (year%4 == 0) && ((year%100 != 0) || (year%400 == 0))
	if isLeap {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronologyLocalDateTime obtains a local date-time in this chronology from another temporal object.
// ChronoLocalDateTime<? extends ChronoLocalDate> localDateTime(TemporalAccessor temporal)
func chronologyLocalDateTime(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.localDateTime: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.localDateTime: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	dateObj := chronologyDateTemporal([]any{self, temporal})
	obj := object.MakeEmptyObjectWithClassName(&classNameChronoLocalDateTime)
	obj.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: dateObj}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: self}
	return obj
}

// chronologyOf obtains an instance of Chronology from a chronology ID or calendar system type.
// static Chronology of(String id)
func chronologyOf(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.of: id is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.of: id is not a String")
	}
	id := object.GoStringFromStringObject(strObj)

	switch strings.ToLower(id) {
	case "iso", "iso8601":
		return GetIsoChronology()
	case "hijrah-umalqura", "islamic-umalqura", "hijrah":
		return getHijrahChronology()
	case "japanese":
		return getJapaneseChronology()
	case "minguo", "roc":
		return getMinguoChronology()
	case "thaibuddhist", "buddhist":
		return getThaiBuddhistChronology()
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unknown chronology: %s", id))
	}
}

// chronologyOfLocale obtains an instance of Chronology from a locale.
// static Chronology ofLocale(Locale locale)
func chronologyOfLocale(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.ofLocale: locale is null")
	}
	locale := params[0].(*object.Object)

	if locale.FieldTable != nil {
		if langFld, ok := locale.FieldTable["language"]; ok {
			if langObj, ok := langFld.Fvalue.(*object.Object); ok && object.IsStringObject(langObj) {
				lang := strings.ToLower(object.GoStringFromStringObject(langObj))
				switch lang {
				case "th":
					return getThaiBuddhistChronology()
				case "ja":
					return getJapaneseChronology()
				}
			}
		}
		if countryFld, ok := locale.FieldTable["country"]; ok {
			if countryObj, ok := countryFld.Fvalue.(*object.Object); ok && object.IsStringObject(countryObj) {
				country := strings.ToLower(object.GoStringFromStringObject(countryObj))
				switch country {
				case "th":
					return getThaiBuddhistChronology()
				case "jp":
					return getJapaneseChronology()
				case "tw":
					return getMinguoChronology()
				}
			}
		}
	}

	return GetIsoChronology()
}

// chronologyPeriod obtains a period for this chronology from years, months and days.
// ChronoPeriod period(int years, int months, int days)
func chronologyPeriod(params []any) any {
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.period: self is null")
	}
	self := params[0].(*object.Object)
	years, ok1 := getInt64Val(params[1])
	months, ok2 := getInt64Val(params[2])
	days, ok3 := getInt64Val(params[3])
	if !ok1 || !ok2 || !ok3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.period: invalid argument types")
	}

	obj := object.MakeEmptyObjectWithClassName(&classNameChronoPeriod)
	obj.FieldTable["years"] = object.Field{Ftype: types.Int, Fvalue: years}
	obj.FieldTable["months"] = object.Field{Ftype: types.Int, Fvalue: months}
	obj.FieldTable["days"] = object.Field{Ftype: types.Int, Fvalue: days}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: self}
	return obj
}

// chronologyProlepticYear calculates the proleptic-year given the era and year-of-era.
// int prolepticYear(Era era, int yearOfEra)
func chronologyProlepticYear(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.prolepticYear: self is null")
	}
	eraObj, _ := params[1].(*object.Object)
	yearOfEra, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Chronology.prolepticYear: invalid yearOfEra type")
	}

	if eraObj != nil && !object.IsNull(eraObj) && eraObj.FieldTable != nil {
		if eraValFld, ok := eraObj.FieldTable["eraValue"]; ok {
			if eraVal, ok := getInt64Val(eraValFld.Fvalue); ok && eraVal == 0 {
				return 1 - yearOfEra
			}
		}
	}
	return yearOfEra
}

// chronologyRange gets the range of valid values for the specified field.
// ValueRange range(ChronoField field)
func chronologyRange(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.range: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.range: field is null")
	}
	field := params[1].(*object.Object)

	if field.FieldTable != nil {
		if r, ok := field.FieldTable["range"].Fvalue.(*object.Object); ok && !object.IsNull(r) {
			return r
		}
	}
	return createValueRange(math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64)
}

// chronologyResolveDate resolves parsed ChronoField values into a date during parsing.
// ChronoLocalDate resolveDate(Map<TemporalField,Long> fieldValues, ResolverStyle resolverStyle)
func chronologyResolveDate(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.resolveDate: self is null")
	}
	self := params[0].(*object.Object)
	return createChronoLocalDate(self, 1970, 1, 1)
}

// chronologyToString outputs this chronology as a String.
// String toString()
func chronologyToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.toString: self is null")
	}
	self := params[0].(*object.Object)
	id := getChronologyId(self)
	return object.StringObjectFromGoString(id)
}

// chronologyZonedDateTimeInstant obtains a ChronoZonedDateTime in this chronology from an Instant and ZoneId.
// ChronoZonedDateTime<? extends ChronoLocalDate> zonedDateTime(Instant instant, ZoneId zone)
func chronologyZonedDateTimeInstant(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.zonedDateTime: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.zonedDateTime: instant is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.zonedDateTime: zone is null")
	}
	self := params[0].(*object.Object)
	instant := params[1].(*object.Object)
	zone := params[2].(*object.Object)

	obj := object.MakeEmptyObjectWithClassName(&classNameChronoZonedDateTime)
	obj.FieldTable["instant"] = object.Field{Ftype: types.Ref, Fvalue: instant}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zone}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: self}
	return obj
}

// chronologyZonedDateTimeTemporal obtains a ChronoZonedDateTime in this chronology from another temporal object.
// ChronoZonedDateTime<? extends ChronoLocalDate> zonedDateTime(TemporalAccessor temporal)
func chronologyZonedDateTimeTemporal(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.zonedDateTime: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Chronology.zonedDateTime: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	obj := object.MakeEmptyObjectWithClassName(&classNameChronoZonedDateTime)
	obj.FieldTable["temporal"] = object.Field{Ftype: types.Ref, Fvalue: temporal}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: self}
	return obj
}

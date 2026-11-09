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

var classNameLocalDateTime = "java/time/LocalDateTime"

var localDateTimeMutex sync.Mutex
var localDateTimeInitialized bool = false

var (
	localDateTimeMin *object.Object
	localDateTimeMax *object.Object
)

func ensureLocalDateTimeInited() {
	localDateTimeMutex.Lock()
	defer localDateTimeMutex.Unlock()
	if localDateTimeInitialized {
		return
	}

	ensureLocalDateInited()
	ensureLocalTimeInited()

	localDateTimeMin = createLocalDateTimeInternal(localDateMin, localTimeMin)
	localDateTimeMax = createLocalDateTimeInternal(localDateMax, localTimeMax)

	_ = statics.AddStatic(classNameLocalDateTime+".MIN", statics.Static{
		Type:  "Ljava/time/LocalDateTime;",
		Value: localDateTimeMin,
	})
	_ = statics.AddStatic(classNameLocalDateTime+".MAX", statics.Static{
		Type:  "Ljava/time/LocalDateTime;",
		Value: localDateTimeMax,
	})

	localDateTimeInitialized = true
}

func createLocalDateTimeInternal(dateObj, timeObj *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameLocalDateTime)
	obj.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: dateObj}
	obj.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: timeObj}
	return obj
}

// createLocalDateTime creates a LocalDateTime object ensuring statics are initialized.
func createLocalDateTime(dateObj, timeObj *object.Object) *object.Object {
	ensureLocalDateTimeInited()
	return createLocalDateTimeInternal(dateObj, timeObj)
}

// CreateLocalDateTime is an exported helper for other packages.
func CreateLocalDateTime(dateObj, timeObj *object.Object) *object.Object {
	return createLocalDateTime(dateObj, timeObj)
}

func createInstantInternal(epochSec, nano int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameInstant)
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: nano}
	return obj
}

func getLocalDateTimeDate(obj *object.Object) *object.Object {
	if obj != nil && obj.FieldTable != nil {
		if dFld, ok := obj.FieldTable["date"]; ok {
			if dObj, ok2 := dFld.Fvalue.(*object.Object); ok2 && dObj != nil {
				return dObj
			}
		}
		if y, m, d, ok := getLocalDateYMD(obj); ok {
			return CreateLocalDate(y, m, d)
		}
	}
	return CreateLocalDate(1970, 1, 1)
}

func getLocalDateTimeTime(obj *object.Object) *object.Object {
	if obj != nil && obj.FieldTable != nil {
		if tFld, ok := obj.FieldTable["time"]; ok {
			if tObj, ok2 := tFld.Fvalue.(*object.Object); ok2 && tObj != nil {
				return tObj
			}
		}
		h, m, s, n := getLocalTimeHMSN(obj)
		return createLocalTime(h, m, s, n)
	}
	return createLocalTime(0, 0, 0, 0)
}

func compareLocalDateInternal(d1, d2 *object.Object) int64 {
	y1, m1, day1, _ := getLocalDateYMD(d1)
	y2, m2, day2, _ := getLocalDateYMD(d2)
	if y1 != y2 {
		if y1 < y2 {
			return -1
		}
		return 1
	}
	if m1 != m2 {
		if m1 < m2 {
			return -1
		}
		return 1
	}
	if day1 != day2 {
		if day1 < day2 {
			return -1
		}
		return 1
	}
	return 0
}

func compareLocalDateTime(dt1, dt2 *object.Object) int64 {
	d1 := getLocalDateTimeDate(dt1)
	d2 := getLocalDateTimeDate(dt2)
	cmp := compareLocalDateInternal(d1, d2)
	if cmp != 0 {
		return cmp
	}
	t1 := getLocalDateTimeTime(dt1)
	t2 := getLocalDateTimeTime(dt2)
	return compareLocalTime(t1, t2)
}

func floorDiv(x, y int64) int64 {
	r := x / y
	if (x^y) < 0 && (r*y != x) {
		r--
	}
	return r
}

func floorMod(x, y int64) int64 {
	r := x % y
	if (x^y) < 0 && r != 0 {
		r += y
	}
	return r
}

func localDateTimePlusDaysInternal(dateObj, timeObj *object.Object, days int64) *object.Object {
	if days == 0 {
		return createLocalDateTime(dateObj, timeObj)
	}
	y, m, d, _ := getLocalDateYMD(dateObj)
	epochDay := ymdToEpochDay(y, m, d)
	newEpochDay := epochDay + days
	ny, nm, nd := epochDayToYMD(newEpochDay)
	newDate := CreateLocalDate(ny, nm, nd)
	return createLocalDateTime(newDate, timeObj)
}

func localDateTimePlusWithOverflow(dateObj, timeObj *object.Object, hours, minutes, seconds, nanos, sign int64) *object.Object {
	if (hours | minutes | seconds | nanos) == 0 {
		return createLocalDateTime(dateObj, timeObj)
	}
	totDays := (nanos / nanosPerDay) + (seconds / secondsPerDay) + (minutes / minutesPerDay) + (hours / hoursPerDay)
	totDays *= sign

	totNanos := (nanos % nanosPerDay) + (seconds%secondsPerDay)*nanosPerSecond + (minutes%minutesPerDay)*nanosPerMinute + (hours%hoursPerDay)*nanosPerHour
	totNanos *= sign

	curNod := getLocalTimeNanoOfDay(timeObj)
	totNanos += curNod

	totDays += floorDiv(totNanos, nanosPerDay)
	newNod := floorMod(totNanos, nanosPerDay)

	var newTime *object.Object
	if newNod == curNod {
		newTime = timeObj
	} else {
		newTime = ofNanoOfDay(newNod)
	}
	return localDateTimePlusDaysInternal(dateObj, newTime, totDays)
}

// Load_Time_LocalDateTime registers all java/time/LocalDateTime methods.
func Load_Time_LocalDateTime() {
	ghelpers.MethodSignatures["java/time/LocalDateTime.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeClinit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.atOffset(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeAtOffset,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.atZone(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeAtZone,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.atZone(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeAtZone,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.compareTo(Ljava/time/chrono/ChronoLocalDateTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeEquals,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeFormat,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeFrom,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeGet,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getDayOfMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getDayOfWeek()Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetDayOfWeek,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getDayOfYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getHour()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetHour,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeGetLong,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getMinute()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetMinute,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getMonth()Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getMonthValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetMonthValue,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getNano()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetNano,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getSecond()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.getYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeGetYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeHashCode,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.isAfter(Ljava/time/chrono/ChronoLocalDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeIsAfter,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.isBefore(Ljava/time/chrono/ChronoLocalDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeIsBefore,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.isEqual(Ljava/time/chrono/ChronoLocalDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeIsEqual,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusDays(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusDays,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusHours(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusHours,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusMinutes(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusMinutes,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusMonths(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusMonths,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusNanos(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusNanos,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusSeconds(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusSeconds,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusWeeks(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusWeeks,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.minusYears(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeMinusYears,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.now()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeNow,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.now(Ljava/time/Clock;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeNowClock,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.now(Ljava/time/ZoneId;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeNowZone,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(IIIII)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  localDateTimeOfHM,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(IIIIII)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  localDateTimeOfHMS,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(IIIIIII)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  localDateTimeOfHMSN,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(ILjava/time/Month;III)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  localDateTimeOfYMonthHM,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(ILjava/time/Month;IIII)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  localDateTimeOfYMonthHMS,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(ILjava/time/Month;IIIII)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  localDateTimeOfYMonthHMSN,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.of(Ljava/time/LocalDate;Ljava/time/LocalTime;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeOfDateTime,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.ofEpochSecond(JILjava/time/ZoneOffset;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localDateTimeOfEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeOfInstant,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.parse(Ljava/lang/CharSequence;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeParse,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeParseFormatter,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusDays(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusDaysMethod,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusHours(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusHours,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusMinutes(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusMinutes,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusMonths(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusMonths,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusNanos(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusNanos,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusSeconds(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusSeconds,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusWeeks(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusWeeks,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.plusYears(J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimePlusYears,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeQuery,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeRange,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.toEpochSecond(Ljava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.toInstant(Ljava/time/ZoneOffset;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeToInstant,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.toLocalDate()Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.toLocalDate()Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.toLocalTime()Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeToLocalTime,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateTimeToString,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeTruncatedTo,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeUntil,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withDayOfMonth(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withDayOfYear(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withHour(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithHour,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withMinute(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithMinute,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withMonth(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withNano(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithNano,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withSecond(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalDateTime.withYear(I)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateTimeWithYear,
		}
}

func localDateTimeClinit(params []any) any {
	ensureLocalDateTimeInited()
	return nil
}

func localDateTimeNow(params []any) any {
	ensureLocalDateTimeInited()
	now := time.Now()
	y := int64(now.Year())
	m := int64(now.Month())
	d := int64(now.Day())
	h := int64(now.Hour())
	minutes := int64(now.Minute())
	s := int64(now.Second())
	n := int64(now.Nanosecond())

	dObj := CreateLocalDate(y, m, d)
	tObj := createLocalTime(h, minutes, s, n)
	return createLocalDateTime(dObj, tObj)
}

func localDateTimeNowZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.now: zone is null")
	}
	return localDateTimeNow(nil)
}

func localDateTimeNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.now: clock is null")
	}
	clock := params[0].(*object.Object)
	instRes := clockInstant([]any{clock})
	if err, isErr := instRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	zoneRes := clockGetZone([]any{clock})
	if err, isErr := zoneRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	zdt := zonedDateTimeOfInstant(instRes.(*object.Object), zoneRes.(*object.Object))
	if err, isErr := zdt.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return zonedDateTimeToLocalDateTime([]any{zdt.(*object.Object)})
}

func localDateTimeOfHM(params []any) any {
	if len(params) < 5 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: missing arguments")
	}
	y, okY := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	d, okD := getInt64Val(params[2])
	h, okH := getInt64Val(params[3])
	min, okMin := getInt64Val(params[4])
	if !okY || !okM || !okD || !okH || !okMin {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: invalid arguments")
	}
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	if err := validateLocalTime(h, min, 0, 0); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, 0, 0))
}

func localDateTimeOfHMS(params []any) any {
	if len(params) < 6 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: missing arguments")
	}
	y, okY := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	d, okD := getInt64Val(params[2])
	h, okH := getInt64Val(params[3])
	min, okMin := getInt64Val(params[4])
	s, okS := getInt64Val(params[5])
	if !okY || !okM || !okD || !okH || !okMin || !okS {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: invalid arguments")
	}
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	if err := validateLocalTime(h, min, s, 0); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, s, 0))
}

func localDateTimeOfHMSN(params []any) any {
	if len(params) < 7 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: missing arguments")
	}
	y, okY := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	d, okD := getInt64Val(params[2])
	h, okH := getInt64Val(params[3])
	min, okMin := getInt64Val(params[4])
	s, okS := getInt64Val(params[5])
	n, okN := getInt64Val(params[6])
	if !okY || !okM || !okD || !okH || !okMin || !okS || !okN {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: invalid arguments")
	}
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	if err := validateLocalTime(h, min, s, n); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, s, n))
}

func localDateTimeOfYMonthHM(params []any) any {
	if len(params) < 5 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: missing arguments")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.of: month is null")
	}
	mOrd, okM := getMonthOrdinal(params[1].(*object.Object))
	if !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: invalid month")
	}
	return localDateTimeOfHM([]any{params[0], mOrd + 1, params[2], params[3], params[4]})
}

func localDateTimeOfYMonthHMS(params []any) any {
	if len(params) < 6 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: missing arguments")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.of: month is null")
	}
	mOrd, okM := getMonthOrdinal(params[1].(*object.Object))
	if !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: invalid month")
	}
	return localDateTimeOfHMS([]any{params[0], mOrd + 1, params[2], params[3], params[4], params[5]})
}

func localDateTimeOfYMonthHMSN(params []any) any {
	if len(params) < 7 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: missing arguments")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.of: month is null")
	}
	mOrd, okM := getMonthOrdinal(params[1].(*object.Object))
	if !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.of: invalid month")
	}
	return localDateTimeOfHMSN([]any{params[0], mOrd + 1, params[2], params[3], params[4], params[5], params[6]})
}

func localDateTimeOfDateTime(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.of: date is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.of: time is null")
	}
	return createLocalDateTime(params[0].(*object.Object), params[1].(*object.Object))
}

func localDateTimeOfEpochSecond(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.ofEpochSecond: missing arguments")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.ofEpochSecond: offset is null")
	}
	epochSec, okSec := getInt64Val(params[0])
	nanoOfSec, okNano := getInt64Val(params[1])
	if !okSec || !okNano {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.ofEpochSecond: invalid arguments")
	}
	if nanoOfSec < 0 || nanoOfSec > 999999999 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for NanoOfSecond (valid values 0 - 999999999): %d", nanoOfSec))
	}
	offsetSec := getZoneOffsetTotalSeconds(params[2].(*object.Object))
	localEpochSec := epochSec + offsetSec
	epochDay := floorDiv(localEpochSec, secondsPerDay)
	secOfDay := floorMod(localEpochSec, secondsPerDay)

	y, m, d := epochDayToYMD(epochDay)
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	dateObj := CreateLocalDate(y, m, d)
	timeObj := ofNanoOfDay(secOfDay*nanosPerSecond + nanoOfSec)
	return createLocalDateTime(dateObj, timeObj)
}

func localDateTimeOfInstant(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.ofInstant: instant is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.ofInstant: zone is null")
	}
	instant := params[0].(*object.Object)
	var epochSec, nano int64
	if sFld, ok := instant.FieldTable["seconds"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			epochSec = s
		}
	} else if sFld, ok := instant.FieldTable["epochSecond"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			epochSec = s
		}
	}
	if nFld, ok := instant.FieldTable["nanos"]; ok {
		if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
			nano = n
		}
	} else if nFld, ok := instant.FieldTable["nano"]; ok {
		if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
			nano = n
		}
	}

	zoneObj := params[1].(*object.Object)
	offsetSec := getZoneOffsetTotalSeconds(zoneObj)
	localEpochSec := epochSec + offsetSec
	epochDay := floorDiv(localEpochSec, secondsPerDay)
	secOfDay := floorMod(localEpochSec, secondsPerDay)

	y, m, d := epochDayToYMD(epochDay)
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	dateObj := CreateLocalDate(y, m, d)
	timeObj := ofNanoOfDay(secOfDay*nanosPerSecond + nano)
	return createLocalDateTime(dateObj, timeObj)
}

func localDateTimeFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.from: temporal is null")
	}
	temporal := params[0].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(temporal.KlassName)
	if kName == classNameLocalDateTime {
		return temporal
	}
	dRes := localDateFrom([]any{temporal})
	if err, isErr := dRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	tRes := localTimeFrom([]any{temporal})
	if err, isErr := tRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	dateObj, okD := dRes.(*object.Object)
	timeObj, okT := tRes.(*object.Object)
	if !okD || !okT {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDateTime.from: unable to obtain LocalDateTime from TemporalAccessor")
	}
	return createLocalDateTime(dateObj, timeObj)
}

func localDateTimeParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.parse: text is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.parse: invalid string")
	}
	text := object.GoStringFromStringObject(strObj)
	parts := strings.Split(text, "T")
	if len(parts) != 2 {
		parts = strings.Split(text, "t")
	}
	if len(parts) != 2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed: missing 'T' separator", text))
	}
	y, m, d, errDate := parseLocalDateString(parts[0])
	if errDate != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed: %v", text, errDate))
	}
	h, min, s, n, errTime := parseLocalTimeString(parts[1])
	if errTime != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed: %v", text, errTime))
	}
	return createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, s, n))
}

func localDateTimeParseFormatter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.parse: text is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.parse: formatter is null")
	}
	textObj := params[0].(*object.Object)
	formatterObj := params[1].(*object.Object)
	if !object.IsStringObject(textObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.parse: invalid string")
	}
	text := object.GoStringFromStringObject(textObj)
	accessor, err := parseTextToTemporalAccessor(formatterObj, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, err.Error())
	}
	return localDateTimeFrom([]any{accessor})
}

func localDateTimeIsSupportedField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isSupported: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return types.JavaBoolFalse
	}
	switch strings.ToUpper(fieldName) {
	case "NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY", "SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY",
		"DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "ERA":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

func localDateTimeIsSupportedUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isSupported: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	unitObj := params[1].(*object.Object)
	unitName := getTemporalUnitName(unitObj)
	if unitName == "" {
		return types.JavaBoolFalse
	}
	switch strings.ToUpper(unitName) {
	case "NANOS", "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS",
		"DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return types.JavaBoolTrue
	case "FOREVER":
		return types.JavaBoolFalse
	default:
		return types.JavaBoolFalse
	}
}

func localDateTimeRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.range: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY", "SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY":
		tObj := getLocalDateTimeTime(self)
		return localTimeRange([]any{tObj, fieldObj})
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "ERA":
		dObj := getLocalDateTimeDate(self)
		return localDateRange([]any{dObj, fieldObj})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func localDateTimeGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.get: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "NANO_OF_SECOND", "MICRO_OF_SECOND", "MILLI_OF_SECOND",
		"SECOND_OF_MINUTE", "SECOND_OF_DAY", "MINUTE_OF_HOUR", "MINUTE_OF_DAY",
		"HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM", "HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY":
		tObj := getLocalDateTimeTime(self)
		return localTimeGet([]any{tObj, fieldObj})
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "ERA":
		dObj := getLocalDateTimeDate(self)
		return localDateGet([]any{dObj, fieldObj})
	default:
		return localDateTimeGetLong(params)
	}
}

func localDateTimeGetLong(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getLong: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY", "SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY":
		tObj := getLocalDateTimeTime(self)
		return localTimeGetLong([]any{tObj, fieldObj})
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "ERA":
		dObj := getLocalDateTimeDate(self)
		return localDateGetLong([]any{dObj, fieldObj})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func localDateTimeToLocalDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toLocalDate: self is null")
	}
	return getLocalDateTimeDate(params[0].(*object.Object))
}

func localDateTimeToLocalTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toLocalTime: self is null")
	}
	return getLocalDateTimeTime(params[0].(*object.Object))
}

func localDateTimeGetYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getYear: self is null")
	}
	dObj := getLocalDateTimeDate(params[0].(*object.Object))
	return localDateGetYear([]any{dObj})
}

func localDateTimeGetMonthValue(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getMonthValue: self is null")
	}
	dObj := getLocalDateTimeDate(params[0].(*object.Object))
	return localDateGetMonthValue([]any{dObj})
}

func localDateTimeGetMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getMonth: self is null")
	}
	dObj := getLocalDateTimeDate(params[0].(*object.Object))
	return localDateGetMonth([]any{dObj})
}

func localDateTimeGetDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getDayOfMonth: self is null")
	}
	dObj := getLocalDateTimeDate(params[0].(*object.Object))
	return localDateGetDayOfMonth([]any{dObj})
}

func localDateTimeGetDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getDayOfYear: self is null")
	}
	dObj := getLocalDateTimeDate(params[0].(*object.Object))
	return localDateGetDayOfYear([]any{dObj})
}

func localDateTimeGetDayOfWeek(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getDayOfWeek: self is null")
	}
	dObj := getLocalDateTimeDate(params[0].(*object.Object))
	return localDateGetDayOfWeek([]any{dObj})
}

func localDateTimeGetHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getHour: self is null")
	}
	tObj := getLocalDateTimeTime(params[0].(*object.Object))
	return localTimeGetHour([]any{tObj})
}

func localDateTimeGetMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getMinute: self is null")
	}
	tObj := getLocalDateTimeTime(params[0].(*object.Object))
	return localTimeGetMinute([]any{tObj})
}

func localDateTimeGetSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getSecond: self is null")
	}
	tObj := getLocalDateTimeTime(params[0].(*object.Object))
	return localTimeGetSecond([]any{tObj})
}

func localDateTimeGetNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.getNano: self is null")
	}
	tObj := getLocalDateTimeTime(params[0].(*object.Object))
	return localTimeGetNano([]any{tObj})
}

func localDateTimeWithAdjuster(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	kName := object.GoStringFromStringPoolIndex(adjuster.KlassName)
	switch kName {
	case classNameLocalDateVal:
		return createLocalDateTime(adjuster, tObj)
	case classNameLocalTime:
		return createLocalDateTime(dObj, adjuster)
	case classNameLocalDateTime:
		return adjuster
	}

	res := temporalAdjusterAdjustInto([]any{adjuster, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localDateTimeFrom([]any{resObj})
	}
	return res
}

func localDateTimeWithField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.with: field is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.with: missing value")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	newValue, okVal := getInt64Val(params[2])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.with: invalid value")
	}
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	switch strings.ToUpper(fieldName) {
	case "NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY", "SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY":
		newTime := localTimeWithField([]any{tObj, fieldObj, newValue})
		if err, isErr := newTime.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createLocalDateTime(dObj, newTime.(*object.Object))
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "ERA":
		newDate := localDateWithField([]any{dObj, fieldObj, newValue})
		if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createLocalDateTime(newDate.(*object.Object), tObj)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func localDateTimeWithYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withYear: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withYear: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDateWithYear([]any{dObj, params[1]})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimeWithMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withMonth: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withMonth: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDateWithMonth([]any{dObj, params[1]})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimeWithDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withDayOfMonth: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withDayOfMonth: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDateWithDayOfMonth([]any{dObj, params[1]})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimeWithDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withDayOfYear: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withDayOfYear: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDateWithDayOfYear([]any{dObj, params[1]})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimeWithHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withHour: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withHour: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newTime := localTimeWithHour([]any{tObj, params[1]})
	if err, isErr := newTime.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(dObj, newTime.(*object.Object))
}

func localDateTimeWithMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withMinute: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withMinute: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newTime := localTimeWithMinute([]any{tObj, params[1]})
	if err, isErr := newTime.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(dObj, newTime.(*object.Object))
}

func localDateTimeWithSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withSecond: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withSecond: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newTime := localTimeWithSecond([]any{tObj, params[1]})
	if err, isErr := newTime.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(dObj, newTime.(*object.Object))
}

func localDateTimeWithNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.withNano: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.withNano: missing argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newTime := localTimeWithNano([]any{tObj, params[1]})
	if err, isErr := newTime.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(dObj, newTime.(*object.Object))
}

func localDateTimeTruncatedTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.truncatedTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.truncatedTo: unit is null")
	}
	self := params[0].(*object.Object)
	unitObj := params[1].(*object.Object)
	unitName := getTemporalUnitName(unitObj)
	if unitName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported unit: null")
	}
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	switch strings.ToUpper(unitName) {
	case "NANOS":
		return self
	case "DAYS":
		return createLocalDateTime(dObj, createLocalTime(0, 0, 0, 0))
	case "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS":
		newTime := localTimeTruncatedTo([]any{tObj, unitObj})
		if err, isErr := newTime.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createLocalDateTime(dObj, newTime.(*object.Object))
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localDateTimePlusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plus: self is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plus: unit is null")
	}
	amountToAdd, okVal := getInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plus: invalid amount")
	}
	if amountToAdd == 0 {
		return params[0].(*object.Object)
	}
	self := params[0].(*object.Object)
	unitObj := params[2].(*object.Object)
	unitName := getTemporalUnitName(unitObj)
	if unitName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported unit: null")
	}
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	switch strings.ToUpper(unitName) {
	case "NANOS":
		return localDateTimePlusNanos([]any{self, amountToAdd})
	case "MICROS":
		return localDateTimePlusWithOverflow(dObj, tObj, 0, 0, 0, amountToAdd*1000, 1)
	case "MILLIS":
		return localDateTimePlusWithOverflow(dObj, tObj, 0, 0, 0, amountToAdd*1000000, 1)
	case "SECONDS":
		return localDateTimePlusSeconds([]any{self, amountToAdd})
	case "MINUTES":
		return localDateTimePlusMinutes([]any{self, amountToAdd})
	case "HOURS":
		return localDateTimePlusHours([]any{self, amountToAdd})
	case "HALF_DAYS":
		return localDateTimePlusWithOverflow(dObj, tObj, amountToAdd*12, 0, 0, 0, 1)
	case "DAYS":
		return localDateTimePlusDaysMethod([]any{self, amountToAdd})
	case "WEEKS":
		return localDateTimePlusWeeks([]any{self, amountToAdd})
	case "MONTHS":
		return localDateTimePlusMonths([]any{self, amountToAdd})
	case "YEARS":
		return localDateTimePlusYears([]any{self, amountToAdd})
	case "DECADES":
		return localDateTimePlusYears([]any{self, amountToAdd * 10})
	case "CENTURIES":
		return localDateTimePlusYears([]any{self, amountToAdd * 100})
	case "MILLENNIA":
		return localDateTimePlusYears([]any{self, amountToAdd * 1000})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localDateTimePlusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	res := temporalAmountAddTo([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localDateTimeFrom([]any{resObj})
	}
	return res
}

func localDateTimePlusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusYears: self is null")
	}
	years, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusYears: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDatePlusYears([]any{dObj, years})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimePlusMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusMonths: self is null")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusMonths: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDatePlusMonths([]any{dObj, months})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimePlusWeeks(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusWeeks: self is null")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusWeeks: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDatePlusWeeks([]any{dObj, weeks})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimePlusDaysMethod(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusDays: self is null")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusDays: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	newDate := localDatePlusDays([]any{dObj, days})
	if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createLocalDateTime(newDate.(*object.Object), tObj)
}

func localDateTimePlusHours(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusHours: self is null")
	}
	hours, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusHours: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, hours, 0, 0, 0, 1)
}

func localDateTimePlusMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusMinutes: self is null")
	}
	minutes, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusMinutes: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, 0, minutes, 0, 0, 1)
}

func localDateTimePlusSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusSeconds: self is null")
	}
	seconds, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusSeconds: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, 0, 0, seconds, 0, 1)
}

func localDateTimePlusNanos(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.plusNanos: self is null")
	}
	nanos, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.plusNanos: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, 0, 0, 0, nanos, 1)
}

func localDateTimeMinusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minus: self is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minus: unit is null")
	}
	amountToSubtract, okVal := getInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minus: invalid amount")
	}
	return localDateTimePlusUnit([]any{params[0], -amountToSubtract, params[2]})
}

func localDateTimeMinusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	res := temporalAmountSubtractFrom([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localDateTimeFrom([]any{resObj})
	}
	return res
}

func localDateTimeMinusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusYears: self is null")
	}
	years, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusYears: invalid argument")
	}
	return localDateTimePlusYears([]any{params[0], -years})
}

func localDateTimeMinusMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusMonths: self is null")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusMonths: invalid argument")
	}
	return localDateTimePlusMonths([]any{params[0], -months})
}

func localDateTimeMinusWeeks(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusWeeks: self is null")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusWeeks: invalid argument")
	}
	return localDateTimePlusWeeks([]any{params[0], -weeks})
}

func localDateTimeMinusDays(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusDays: self is null")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusDays: invalid argument")
	}
	return localDateTimePlusDaysMethod([]any{params[0], -days})
}

func localDateTimeMinusHours(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusHours: self is null")
	}
	hours, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusHours: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, hours, 0, 0, 0, -1)
}

func localDateTimeMinusMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusMinutes: self is null")
	}
	minutes, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusMinutes: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, 0, minutes, 0, 0, -1)
}

func localDateTimeMinusSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusSeconds: self is null")
	}
	seconds, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusSeconds: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, 0, 0, seconds, 0, -1)
}

func localDateTimeMinusNanos(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.minusNanos: self is null")
	}
	nanos, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDateTime.minusNanos: invalid argument")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)
	return localDateTimePlusWithOverflow(dObj, tObj, 0, 0, 0, nanos, -1)
}

func localDateTimeQuery(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	queryName := ""
	if query.FieldTable != nil {
		if nameFld, okName := query.FieldTable["name"]; okName {
			if strObj, okStr := nameFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				queryName = object.GoStringFromStringObject(strObj)
			}
		}
	}

	switch queryName {
	case "chronology":
		return GetIsoChronology()
	case "precision":
		return ghelpers.CreateTemporalUnit("Nanos")
	case "localDate":
		return getLocalDateTimeDate(self)
	case "localTime":
		return getLocalDateTimeTime(self)
	case "zoneId", "zone", "offset":
		return nil
	default:
		return executeTemporalQuery(query, self)
	}
}

func localDateTimeAdjustInto(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	target := params[1].(*object.Object)

	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(dObj)
	y, m, d, _ := getLocalDateYMD(dObj)
	h, min, s, n := getLocalTimeHMSN(tObj)
	nod := getLocalTimeNanoOfDay(tObj)

	if target.FieldTable != nil {
		target.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: dObj}
		target.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: tObj}
		target.FieldTable["epochDay"] = object.Field{Ftype: types.Long, Fvalue: epochDay}
		target.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: y}
		target.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: m}
		target.FieldTable["day"] = object.Field{Ftype: types.Int, Fvalue: d}
		target.FieldTable["hour"] = object.Field{Ftype: types.Int, Fvalue: h}
		target.FieldTable["minute"] = object.Field{Ftype: types.Int, Fvalue: min}
		target.FieldTable["second"] = object.Field{Ftype: types.Int, Fvalue: s}
		target.FieldTable["nano"] = object.Field{Ftype: types.Int, Fvalue: n}
		target.FieldTable["nanoOfDay"] = object.Field{Ftype: types.Long, Fvalue: nod}
	}
	return target
}

func localDateTimeUntil(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.until: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.until: unit is null")
	}
	self := params[0].(*object.Object)
	endExclusive := params[1].(*object.Object)
	unitObj := params[2].(*object.Object)

	endLdtRes := localDateTimeFrom([]any{endExclusive})
	if err, isErr := endLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	endLdt := endLdtRes.(*object.Object)

	unitName := getTemporalUnitName(unitObj)
	if unitName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported unit: null")
	}

	d1 := getLocalDateTimeDate(self)
	t1 := getLocalDateTimeTime(self)
	d2 := getLocalDateTimeDate(endLdt)
	t2 := getLocalDateTimeTime(endLdt)

	y1, m1, day1, _ := getLocalDateYMD(d1)
	y2, m2, day2, _ := getLocalDateYMD(d2)
	epochDay1 := ymdToEpochDay(y1, m1, day1)
	epochDay2 := ymdToEpochDay(y2, m2, day2)
	daysUntil := epochDay2 - epochDay1
	nod1 := getLocalTimeNanoOfDay(t1)
	nod2 := getLocalTimeNanoOfDay(t2)
	timeNanos := nod2 - nod1

	switch strings.ToUpper(unitName) {
	case "NANOS":
		return daysUntil*nanosPerDay + timeNanos
	case "MICROS":
		return (daysUntil*nanosPerDay + timeNanos) / 1000
	case "MILLIS":
		return (daysUntil*nanosPerDay + timeNanos) / 1000000
	case "SECONDS":
		sec1 := nod1 / nanosPerSecond
		sec2 := nod2 / nanosPerSecond
		return daysUntil*secondsPerDay + (sec2 - sec1)
	case "MINUTES":
		min1 := nod1 / nanosPerMinute
		min2 := nod2 / nanosPerMinute
		return daysUntil*minutesPerDay + (min2 - min1)
	case "HOURS":
		h1 := nod1 / nanosPerHour
		h2 := nod2 / nanosPerHour
		return daysUntil*hoursPerDay + (h2 - h1)
	case "HALF_DAYS":
		h1 := nod1 / nanosPerHour
		h2 := nod2 / nanosPerHour
		return (daysUntil*hoursPerDay + (h2 - h1)) / 12
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		adjD2 := d2
		if epochDay2 > epochDay1 && nod2 < nod1 {
			adjD2 = localDateMinusDays([]any{d2, int64(1)}).(*object.Object)
		} else if epochDay2 < epochDay1 && nod2 > nod1 {
			adjD2 = localDatePlusDays([]any{d2, int64(1)}).(*object.Object)
		}
		return localDateUntilUnit([]any{d1, adjD2, unitObj})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localDateTimeFormat(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.format: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.format: formatter is null")
	}
	return dateTimeFormatterFormat([]any{params[1], params[0]})
}

func localDateTimeAtOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.atOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.atOffset: offset is null")
	}
	self := params[0].(*object.Object)
	offsetObj := params[1].(*object.Object)
	classNameOffsetDateTime := "java/time/OffsetDateTime"
	odt := object.MakeEmptyObjectWithClassName(&classNameOffsetDateTime)
	odt.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: self}
	odt.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return odt
}

func localDateTimeAtZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.atZone: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.atZone: zone is null")
	}
	self := params[0].(*object.Object)
	zoneObj := params[1].(*object.Object)
	return zonedDateTimeOfLocal(self, zoneObj, nil)
}

func localDateTimeCompareTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	return compareLocalDateTime(self, other)
}

func localDateTimeIsAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isAfter: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isAfter: other is null")
	}
	cmp := compareLocalDateTime(params[0].(*object.Object), params[1].(*object.Object))
	return types.ConvertGoBoolToJavaBool(cmp > 0)
}

func localDateTimeIsBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isBefore: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isBefore: other is null")
	}
	cmp := compareLocalDateTime(params[0].(*object.Object), params[1].(*object.Object))
	return types.ConvertGoBoolToJavaBool(cmp < 0)
}

func localDateTimeIsEqual(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isEqual: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.isEqual: other is null")
	}
	cmp := compareLocalDateTime(params[0].(*object.Object), params[1].(*object.Object))
	return types.ConvertGoBoolToJavaBool(cmp == 0)
}

func localDateTimeEquals(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	if self == other {
		return types.JavaBoolTrue
	}
	kName := object.GoStringFromStringPoolIndex(other.KlassName)
	if kName != classNameLocalDateTime {
		return types.JavaBoolFalse
	}
	d1 := getLocalDateTimeDate(self)
	t1 := getLocalDateTimeTime(self)
	d2 := getLocalDateTimeDate(other)
	t2 := getLocalDateTimeTime(other)

	return types.ConvertGoBoolToJavaBool(compareLocalDateInternal(d1, d2) == 0 && compareLocalTime(t1, t2) == 0)
}

func localDateTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	hD := localDateHashCode([]any{dObj})
	hT := localTimeHashCode([]any{tObj})
	var vD, vT int64
	if v, ok := getInt64Val(hD); ok {
		vD = v
	}
	if v, ok := getInt64Val(hT); ok {
		vT = v
	}
	return vD ^ vT
}

func localDateTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toString: self is null")
	}
	self := params[0].(*object.Object)
	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	y, m, d, _ := getLocalDateYMD(dObj)
	h, minutes, s, n := getLocalTimeHMSN(tObj)

	str := formatLocalDateString(y, m, d) + "T" + formatLocalTimeString(h, minutes, s, n)
	return object.StringObjectFromGoString(str)
}

func localDateTimeToEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toEpochSecond: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toEpochSecond: offset is null")
	}
	self := params[0].(*object.Object)
	offsetObj := params[1].(*object.Object)

	dObj := getLocalDateTimeDate(self)
	tObj := getLocalDateTimeTime(self)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(dObj)
	secOfDay := getLocalTimeNanoOfDay(tObj) / nanosPerSecond
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)

	return epochDay*secondsPerDay + secOfDay - offsetSec
}

func localDateTimeToInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toInstant: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDateTime.toInstant: offset is null")
	}
	self := params[0].(*object.Object)
	offsetObj := params[1].(*object.Object)

	epochSecVal := localDateTimeToEpochSecond([]any{self, offsetObj})
	if err, isErr := epochSecVal.(*ghelpers.GErrBlk); isErr {
		return err
	}
	epochSec, _ := getInt64Val(epochSecVal)
	tObj := getLocalDateTimeTime(self)
	_, _, _, nano := getLocalTimeHMSN(tObj)

	return createInstantInternal(epochSec, nano)
}

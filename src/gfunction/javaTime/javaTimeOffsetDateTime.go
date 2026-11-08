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
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	classNameOffsetDateTime           = "java/time/OffsetDateTime"
	classNameOffsetDateTimeComparator = "java/time/OffsetDateTime$1"
	offsetDateTimeMutex               sync.Mutex
	offsetDateTimeInitialized         bool = false
)

func ensureOffsetDateTimeInited() {
	offsetDateTimeMutex.Lock()
	defer offsetDateTimeMutex.Unlock()
	if offsetDateTimeInitialized {
		return
	}
	ensureLocalDateTimeInited()
	ensureZoneOffsetInited()
	ensureZoneIdInited()

	minLdt := localDateTimeMin
	maxLdt := localDateTimeMax
	maxOffset := zoneOffsetMax
	minOffset := zoneOffsetMin

	minObj := createOffsetDateTime(minLdt, maxOffset)
	maxObj := createOffsetDateTime(maxLdt, minOffset)

	_ = statics.AddStatic(classNameOffsetDateTime+".MIN", statics.Static{
		Type:  "Ljava/time/OffsetDateTime;",
		Value: minObj,
	})
	_ = statics.AddStatic(classNameOffsetDateTime+".MAX", statics.Static{
		Type:  "Ljava/time/OffsetDateTime;",
		Value: maxObj,
	})

	offsetDateTimeInitialized = true
}

// Load_Time_OffsetDateTime registers method signatures for java/time/OffsetDateTime.
func Load_Time_OffsetDateTime() {
	ensureOffsetDateTimeInited()

	ghelpers.MethodSignatures["java/time/OffsetDateTime.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeClinit,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeInit,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.<init>(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeInitWithLdtOffset,
		}

	// Static Factories
	ghelpers.MethodSignatures["java/time/OffsetDateTime.now()Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeNow,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.now(Ljava/time/ZoneId;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeNowZone,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.now(Ljava/time/Clock;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeNowClock,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.of(Ljava/time/LocalDate;Ljava/time/LocalTime;Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  offsetDateTimeOfDate,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.of(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeOfLDT,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.of(IIIIIIILjava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 8,
			GFunction:  offsetDateTimeOfInts,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeOfInstantMethod,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeFrom,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.parse(Ljava/lang/CharSequence;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeParse,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeParseFormatter,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.timeLineOrder()Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeTimelineOrder,
		}

	// Instance Queries & Getters
	ghelpers.MethodSignatures["java/time/OffsetDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeRange,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeGet,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeGetLong,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getOffset()Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetOffset,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withOffsetSameLocal(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithOffsetSameLocal,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withOffsetSameInstant(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithOffsetSameInstant,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toLocalDateTime()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToLocalDateTime,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toLocalDate()Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetYear,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getMonthValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetMonthValue,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getMonth()Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetMonth,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getDayOfMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getDayOfYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getDayOfWeek()Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetDayOfWeek,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toLocalTime()Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToLocalTime,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getHour()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetHour,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getMinute()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetMinute,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getSecond()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetSecond,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.getNano()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeGetNano,
		}

	// With / Alterations
	ghelpers.MethodSignatures["java/time/OffsetDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withYear(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithYear,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withMonth(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithMonth,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withDayOfMonth(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withDayOfYear(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withHour(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithHour,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withMinute(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithMinute,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withSecond(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithSecond,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.withNano(I)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWithNano,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeTruncatedTo,
		}

	// Plus
	ghelpers.MethodSignatures["java/time/OffsetDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusYears(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusYears,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusMonths(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusMonths,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusWeeks(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusWeeks,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusDays(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusDays,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusHours(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusHours,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusMinutes(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusMinutes,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusSeconds(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusSeconds,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.plusNanos(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimePlusNanos,
		}

	// Minus
	ghelpers.MethodSignatures["java/time/OffsetDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusYears(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusYears,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusMonths(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusMonths,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusWeeks(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusWeeks,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusDays(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusDays,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusHours(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusHours,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusMinutes(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusMinutes,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusSeconds(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusSeconds,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.minusNanos(J)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeMinusNanos,
		}

	// Temporal Queries / AdjustInto / Until / Format
	ghelpers.MethodSignatures["java/time/OffsetDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeQuery,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeUntil,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeFormat,
		}

	// Conversions
	ghelpers.MethodSignatures["java/time/OffsetDateTime.atZoneSameInstant(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeAtZoneSameInstant,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.atZoneSimilarLocal(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeAtZoneSimilarLocal,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toOffsetTime()Ljava/time/OffsetTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToOffsetTime,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toEpochSecond()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toInstant()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToInstant,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toZonedDateTime()Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToZonedDateTime,
		}

	// Comparisons & Object methods
	ghelpers.MethodSignatures["java/time/OffsetDateTime.compareTo(Ljava/time/OffsetDateTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.isAfter(Ljava/time/OffsetDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeIsAfter,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.isBefore(Ljava/time/OffsetDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeIsBefore,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.isEqual(Ljava/time/OffsetDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeIsEqual,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeEquals,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeHashCode,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  offsetDateTimeToString,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.writeExternal(Ljava/io/DataOutput;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWriteExternal,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.readExternal(Ljava/io/DataInput;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeReadExternal,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime.write(Ljava/io/DataOutput;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  offsetDateTimeWriteExternal,
		}

	// Comparator signatures for OffsetDateTime$1
	ghelpers.MethodSignatures["java/time/OffsetDateTime$1.compare(Ljava/lang/Object;Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeTimelineCompare,
		}

	ghelpers.MethodSignatures["java/time/OffsetDateTime$1.compare(Ljava/time/OffsetDateTime;Ljava/time/OffsetDateTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  offsetDateTimeTimelineCompare,
		}
}

// Helpers

func createOffsetDateTime(dateTime, offset *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameOffsetDateTime)
	obj.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: dateTime}
	obj.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offset}
	return obj
}

func getOffsetDateTimeDateTime(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["dateTime"]; ok {
		if ldt, okLdt := fld.Fvalue.(*object.Object); okLdt {
			return ldt
		}
	}
	return nil
}

func getOffsetDateTimeOffset(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["offset"]; ok {
		if off, okOff := fld.Fvalue.(*object.Object); okOff {
			return off
		}
	}
	return nil
}

func getOffsetDateTimeEpochSecond(obj *object.Object) int64 {
	ldt := getOffsetDateTimeDateTime(obj)
	offset := getOffsetDateTimeOffset(obj)
	if ldt == nil || offset == nil {
		return 0
	}
	secVal := localDateTimeToEpochSecond([]any{ldt, offset})
	if s, ok := secVal.(int64); ok {
		return s
	}
	return 0
}

func getOffsetDateTimeNano(obj *object.Object) int64 {
	ldt := getOffsetDateTimeDateTime(obj)
	if ldt == nil {
		return 0
	}
	res := localDateTimeGetNano([]any{ldt})
	if n, ok := res.(int64); ok {
		return n
	}
	return 0
}

func createOffsetTime(timeObj, offsetObj *object.Object) *object.Object {
	classNameOffsetTime := "java/time/OffsetTime"
	ot := object.MakeEmptyObjectWithClassName(&classNameOffsetTime)
	ot.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: timeObj}
	ot.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return ot
}

// Implementation of functions

func offsetDateTimeClinit([]any) any {
	ensureOffsetDateTimeInited()
	return nil
}

func offsetDateTimeInit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.<init>: self is null")
	}
	self := params[0].(*object.Object)
	epochDate := CreateLocalDate(1970, 1, 1)
	midnight := createLocalTime(0, 0, 0, 0)
	ldt := createLocalDateTime(epochDate, midnight)
	offset := zoneOffsetUTC
	self.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	self.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offset}
	return nil
}

func offsetDateTimeInitWithLdtOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.<init>: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.<init>: invalid arguments")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.<init>: dateTime is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.<init>: offset is null")
	}
	self := params[0].(*object.Object)
	self.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: params[1].(*object.Object)}
	self.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: params[2].(*object.Object)}
	return nil
}

func offsetDateTimeNow([]any) any {
	defZone := zoneIdSystemDefault(nil)
	if zObj, ok := defZone.(*object.Object); ok && zObj != nil && !object.IsNull(zObj) {
		return offsetDateTimeNowZone([]any{zObj})
	}
	zoneObj := makeZoneIdFromID("UTC")
	now := time.Now().UTC()
	return offsetDateTimeFromTimeAndZone(now, zoneObj)
}

func offsetDateTimeNowZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.now: zone is null")
	}
	zoneObj := params[0].(*object.Object)
	now := time.Now()
	return offsetDateTimeOfInstant(createInstant(now.Unix(), int64(now.Nanosecond())), zoneObj)
}

func offsetDateTimeNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.now: clock is null")
	}
	clockObj := params[0].(*object.Object)
	instRes := clockInstant([]any{clockObj})
	if err, isErr := instRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	instObj := instRes.(*object.Object)
	zoneRes := clockGetZone([]any{clockObj})
	if err, isErr := zoneRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	zoneObj := zoneRes.(*object.Object)
	return offsetDateTimeOfInstant(instObj, zoneObj)
}

func offsetDateTimeFromTimeAndZone(t time.Time, zoneObj *object.Object) *object.Object {
	year := int64(t.Year())
	month := int64(t.Month())
	day := int64(t.Day())
	hour := int64(t.Hour())
	min := int64(t.Minute())
	sec := int64(t.Second())
	nano := int64(t.Nanosecond())

	dateObj := CreateLocalDate(year, month, day)
	timeObj := createLocalTime(hour, min, sec, nano)
	ldt := createLocalDateTime(dateObj, timeObj)

	_, offsetSec := t.Zone()
	offsetObj := createZoneOffset(int64(offsetSec))

	return createOffsetDateTime(ldt, offsetObj)
}

func offsetDateTimeOfDate(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.of: invalid argument count")
	}
	if params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.of: date is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.of: time is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.of: offset is null")
	}
	date := params[0].(*object.Object)
	timeObj := params[1].(*object.Object)
	offset := params[2].(*object.Object)
	ldt := createLocalDateTime(date, timeObj)
	return createOffsetDateTime(ldt, offset)
}

func offsetDateTimeOfLDT(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.of: invalid argument count")
	}
	if params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.of: dateTime is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.of: offset is null")
	}
	ldt := params[0].(*object.Object)
	offset := params[1].(*object.Object)
	return createOffsetDateTime(ldt, offset)
}

func offsetDateTimeOfInts(params []any) any {
	if len(params) < 8 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.of: invalid argument count")
	}
	year, ok1 := getInt64Val(params[0])
	month, ok2 := getInt64Val(params[1])
	day, ok3 := getInt64Val(params[2])
	hour, ok4 := getInt64Val(params[3])
	min, ok5 := getInt64Val(params[4])
	sec, ok6 := getInt64Val(params[5])
	nano, ok7 := getInt64Val(params[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.of: invalid integer parameters")
	}
	if params[7] == nil || object.IsNull(params[7].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.of: offset is null")
	}
	offset := params[7].(*object.Object)

	ldtRes := localDateTimeOfHMSN([]any{year, month, day, hour, min, sec, nano})
	if err, isErr := ldtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	ldt := ldtRes.(*object.Object)
	return createOffsetDateTime(ldt, offset)
}

func offsetDateTimeOfInstant(instant, zone *object.Object) any {
	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rules := rulesRes.(*object.Object)
	offsetRes := zoneRulesGetOffsetInstant([]any{rules, instant})
	if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	offset := offsetRes.(*object.Object)

	epochSec := getInstantEpochSecondVal(instant)
	nano := getInstantNanoVal(instant)

	ldtRes := localDateTimeOfEpochSecond([]any{epochSec, nano, offset})
	if err, isErr := ldtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	ldt := ldtRes.(*object.Object)
	return createOffsetDateTime(ldt, offset)
}

func offsetDateTimeOfInstantMethod(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.ofInstant: invalid argument count")
	}
	if params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.ofInstant: instant is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.ofInstant: zone is null")
	}
	instant := params[0].(*object.Object)
	zone := params[1].(*object.Object)
	return offsetDateTimeOfInstant(instant, zone)
}

func offsetDateTimeFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.from: temporal is null")
	}
	accessor := params[0].(*object.Object)
	className := getObjectClassName(accessor)
	if strings.Contains(className, "OffsetDateTime") {
		return accessor
	}
	if strings.Contains(className, "ZonedDateTime") {
		ldt := getZonedDateTimeLocalDateTime(accessor)
		offset := getZonedDateTimeOffset(accessor)
		return createOffsetDateTime(ldt, offset)
	}

	// Try extracting offset and local date-time
	ensureTemporalQueriesInited()
	offsetRes := temporalAccessorQuery([]any{accessor, offsetQueryInstance})
	if offObj, ok := offsetRes.(*object.Object); ok && offObj != nil && !object.IsNull(offObj) {
		ldtRes := localDateTimeFrom([]any{accessor})
		if ldtObj, okLdt := ldtRes.(*object.Object); okLdt && ldtObj != nil && !object.IsNull(ldtObj) {
			return createOffsetDateTime(ldtObj, offObj)
		}
	}

	// Try extracting instant and zone
	instRes := instantFrom([]any{accessor})
	if instObj, okInst := instRes.(*object.Object); okInst && instObj != nil && !object.IsNull(instObj) {
		zoneRes := temporalAccessorQuery([]any{accessor, zoneIdQueryInstance})
		if zoneObj, okZone := zoneRes.(*object.Object); okZone && zoneObj != nil && !object.IsNull(zoneObj) {
			return offsetDateTimeOfInstant(instObj, zoneObj)
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unable to obtain OffsetDateTime from TemporalAccessor: %s", className))
}

func parseOffsetDateTimeString(text string) any {
	text = strings.TrimSpace(text)
	if text == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text cannot be parsed to OffsetDateTime: empty string")
	}

	// Format: YYYY-MM-DDTHH:MM:SS[.fraction][+HH:MM|-HH:MM|Z]
	tIdx := strings.IndexByte(text, 'T')
	if tIdx < 0 {
		tIdx = strings.IndexByte(text, 't')
	}
	if tIdx < 0 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed to OffsetDateTime: missing 'T'", text))
	}

	dateStr := text[:tIdx]
	timeAndOffsetStr := text[tIdx+1:]

	// Find where offset starts
	offsetIdx := -1
	for i := len(timeAndOffsetStr) - 1; i >= 0; i-- {
		ch := timeAndOffsetStr[i]
		if ch == '+' || ch == '-' || ch == 'Z' || ch == 'z' {
			offsetIdx = i
			break
		}
	}
	if offsetIdx < 0 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed to OffsetDateTime: missing offset", text))
	}

	timeStr := timeAndOffsetStr[:offsetIdx]
	offsetStr := timeAndOffsetStr[offsetIdx:]

	// Parse date
	dateParts := strings.Split(dateStr, "-")
	if len(dateParts) != 3 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Invalid date component in '%s'", text))
	}
	y, err1 := strconv.ParseInt(dateParts[0], 10, 64)
	m, err2 := strconv.ParseInt(dateParts[1], 10, 64)
	d, err3 := strconv.ParseInt(dateParts[2], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Invalid date numbers in '%s'", text))
	}

	// Parse time
	var h, min, s, n int64
	timeParts := strings.Split(timeStr, ":")
	if len(timeParts) < 2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Invalid time component in '%s'", text))
	}
	h, err1 = strconv.ParseInt(timeParts[0], 10, 64)
	min, err2 = strconv.ParseInt(timeParts[1], 10, 64)
	if err1 != nil || err2 != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Invalid time numbers in '%s'", text))
	}
	if len(timeParts) >= 3 {
		secPart := timeParts[2]
		if dotIdx := strings.IndexByte(secPart, '.'); dotIdx >= 0 {
			s, _ = strconv.ParseInt(secPart[:dotIdx], 10, 64)
			nanoStr := secPart[dotIdx+1:]
			if len(nanoStr) > 9 {
				nanoStr = nanoStr[:9]
			} else {
				nanoStr = nanoStr + strings.Repeat("0", 9-len(nanoStr))
			}
			n, _ = strconv.ParseInt(nanoStr, 10, 64)
		} else {
			s, _ = strconv.ParseInt(secPart, 10, 64)
		}
	}

	offsetObjRes := zoneOffsetOf([]any{object.StringObjectFromGoString(offsetStr)})
	if err, isErr := offsetObjRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	offsetObj := offsetObjRes.(*object.Object)

	ldt := createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, s, n))
	return createOffsetDateTime(ldt, offsetObj)
}

func offsetDateTimeParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.parse: text is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.parse: invalid string")
	}
	text := object.GoStringFromStringObject(strObj)
	return parseOffsetDateTimeString(text)
}

func offsetDateTimeParseFormatter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.parse: text is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.parse: formatter is null")
	}
	textObj := params[0].(*object.Object)
	formatterObj := params[1].(*object.Object)
	if !object.IsStringObject(textObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.parse: invalid string")
	}
	text := object.GoStringFromStringObject(textObj)
	accessor, err := parseTextToTemporalAccessor(formatterObj, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, err.Error())
	}
	return offsetDateTimeFrom([]any{accessor})
}

func compareInstant(dt1, dt2 *object.Object) int64 {
	off1 := getOffsetDateTimeOffset(dt1)
	off2 := getOffsetDateTimeOffset(dt2)
	eqOff := zoneOffsetEquals([]any{off1, off2})
	if eqOff == types.JavaBoolTrue {
		ldt1 := getOffsetDateTimeDateTime(dt1)
		ldt2 := getOffsetDateTimeDateTime(dt2)
		res := localDateTimeCompareTo([]any{ldt1, ldt2})
		if r, ok := res.(int64); ok {
			return r
		}
		return 0
	}

	sec1 := getOffsetDateTimeEpochSecond(dt1)
	sec2 := getOffsetDateTimeEpochSecond(dt2)
	if sec1 < sec2 {
		return int64(-1)
	}
	if sec1 > sec2 {
		return int64(1)
	}
	nano1 := getOffsetDateTimeNano(dt1)
	nano2 := getOffsetDateTimeNano(dt2)
	if nano1 < nano2 {
		return int64(-1)
	}
	if nano1 > nano2 {
		return int64(1)
	}
	return int64(0)
}

func offsetDateTimeTimelineOrder([]any) any {
	obj := object.MakeEmptyObjectWithClassName(&classNameOffsetDateTimeComparator)
	return obj
}

func offsetDateTimeTimelineCompare(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTimeComparator.compare: invalid argument count")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTimeComparator.compare: null argument")
	}
	dt1 := params[1].(*object.Object)
	dt2 := params[2].(*object.Object)
	return compareInstant(dt1, dt2)
}

// Instance queries & getters

func offsetDateTimeIsSupportedField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isSupported: self is null")
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
	case "INSTANT_SECONDS", "OFFSET_SECONDS",
		"NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
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

func offsetDateTimeIsSupportedUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isSupported: self is null")
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

func offsetDateTimeRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.range: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "INSTANT_SECONDS":
		return createValueRange(instantMinSecond, instantMinSecond, instantMaxSecond, instantMaxSecond)
	case "OFFSET_SECONDS":
		return createValueRange(-18*3600, -18*3600, 18*3600, 18*3600)
	default:
		ldt := getOffsetDateTimeDateTime(self)
		return localDateTimeRange([]any{ldt, fieldObj})
	}
}

func offsetDateTimeGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.get: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "OFFSET_SECONDS":
		offset := getOffsetDateTimeOffset(self)
		sec, _ := getZoneOffsetTotalSecondsVal(offset)
		return int64(sec)
	case "INSTANT_SECONDS":
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Field INSTANT_SECONDS exceeds int range; use getLong()")
	default:
		ldt := getOffsetDateTimeDateTime(self)
		return localDateTimeGet([]any{ldt, fieldObj})
	}
}

func offsetDateTimeGetLong(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getLong: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "INSTANT_SECONDS":
		return getOffsetDateTimeEpochSecond(self)
	case "OFFSET_SECONDS":
		offset := getOffsetDateTimeOffset(self)
		sec, _ := getZoneOffsetTotalSecondsVal(offset)
		return sec
	default:
		ldt := getOffsetDateTimeDateTime(self)
		return localDateTimeGetLong([]any{ldt, fieldObj})
	}
}

func offsetDateTimeGetOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getOffset: self is null")
	}
	self := params[0].(*object.Object)
	return getOffsetDateTimeOffset(self)
}

func offsetDateTimeWithOffsetSameLocal(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.withOffsetSameLocal: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.withOffsetSameLocal: offset is null")
	}
	self := params[0].(*object.Object)
	newOffset := params[1].(*object.Object)
	curOffset := getOffsetDateTimeOffset(self)
	eq := zoneOffsetEquals([]any{curOffset, newOffset})
	if eq == types.JavaBoolTrue {
		return self
	}
	ldt := getOffsetDateTimeDateTime(self)
	return createOffsetDateTime(ldt, newOffset)
}

func offsetDateTimeWithOffsetSameInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.withOffsetSameInstant: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.withOffsetSameInstant: offset is null")
	}
	self := params[0].(*object.Object)
	newOffset := params[1].(*object.Object)
	curOffset := getOffsetDateTimeOffset(self)
	eq := zoneOffsetEquals([]any{curOffset, newOffset})
	if eq == types.JavaBoolTrue {
		return self
	}
	curSec, _ := getZoneOffsetTotalSecondsVal(curOffset)
	newSec, _ := getZoneOffsetTotalSecondsVal(newOffset)
	diffSec := newSec - curSec

	ldt := getOffsetDateTimeDateTime(self)
	adjustedLdtRes := localDateTimePlusSeconds([]any{ldt, diffSec})
	if err, isErr := adjustedLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(adjustedLdtRes.(*object.Object), newOffset)
}

func offsetDateTimeToLocalDateTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toLocalDateTime: self is null")
	}
	self := params[0].(*object.Object)
	return getOffsetDateTimeDateTime(self)
}

func offsetDateTimeToLocalDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toLocalDate: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeToLocalDate([]any{ldt})
}

func offsetDateTimeGetYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getYear: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetYear([]any{ldt})
}

func offsetDateTimeGetMonthValue(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getMonthValue: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetMonthValue([]any{ldt})
}

func offsetDateTimeGetMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getMonth: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetMonth([]any{ldt})
}

func offsetDateTimeGetDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getDayOfMonth: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetDayOfMonth([]any{ldt})
}

func offsetDateTimeGetDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getDayOfYear: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetDayOfYear([]any{ldt})
}

func offsetDateTimeGetDayOfWeek(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getDayOfWeek: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetDayOfWeek([]any{ldt})
}

func offsetDateTimeToLocalTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toLocalTime: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeToLocalTime([]any{ldt})
}

func offsetDateTimeGetHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getHour: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetHour([]any{ldt})
}

func offsetDateTimeGetMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getMinute: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetMinute([]any{ldt})
}

func offsetDateTimeGetSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getSecond: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetSecond([]any{ldt})
}

func offsetDateTimeGetNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.getNano: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	return localDateTimeGetNano([]any{ldt})
}

// With / Alterations

func offsetDateTimeWithAdjuster(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)
	className := getObjectClassName(adjuster)

	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)

	if strings.Contains(className, "LocalDate") || strings.Contains(className, "LocalTime") || strings.Contains(className, "LocalDateTime") {
		newLdtRes := localDateTimeWithAdjuster([]any{ldt, adjuster})
		if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createOffsetDateTime(newLdtRes.(*object.Object), offset)
	}
	if strings.Contains(className, "ZoneOffset") {
		return offsetDateTimeWithOffsetSameLocal([]any{self, adjuster})
	}
	if strings.Contains(className, "OffsetDateTime") {
		return adjuster
	}
	if strings.Contains(className, "Instant") {
		return offsetDateTimeOfInstant(adjuster, offset)
	}

	return temporalAdjusterAdjustInto([]any{adjuster, self})
}

func offsetDateTimeWithField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.with: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.with: invalid argument count")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.with: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	val, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.with: invalid value")
	}

	fieldName := getTemporalFieldName(fieldObj)
	switch strings.ToUpper(fieldName) {
	case "INSTANT_SECONDS":
		nano := getOffsetDateTimeNano(self)
		inst := createInstant(val, nano)
		offset := getOffsetDateTimeOffset(self)
		return offsetDateTimeOfInstant(inst, offset)
	case "OFFSET_SECONDS":
		newOffsetRes := zoneOffsetOfTotalSeconds([]any{val})
		if err, isErr := newOffsetRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return offsetDateTimeWithOffsetSameLocal([]any{self, newOffsetRes.(*object.Object)})
	default:
		ldt := getOffsetDateTimeDateTime(self)
		offset := getOffsetDateTimeOffset(self)
		newLdtRes := localDateTimeWithField([]any{ldt, fieldObj, val})
		if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createOffsetDateTime(newLdtRes.(*object.Object), offset)
	}
}

func offsetDateTimeWithYear(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withYear: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithYear([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithMonth(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withMonth: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithMonth([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithDayOfMonth(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withDayOfMonth: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithDayOfMonth([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithDayOfYear(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withDayOfYear: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithDayOfYear([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithHour(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withHour: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithHour([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithMinute(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withMinute: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithMinute([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithSecond(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withSecond: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithSecond([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeWithNano(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.withNano: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeWithNano([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimeTruncatedTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.truncatedTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.truncatedTo: unit is null")
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimeTruncatedTo([]any{ldt, unit})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

// Plus

func offsetDateTimePlusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.plus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)
	return temporalAmountAddTo([]any{amount, self})
}

func offsetDateTimePlusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.plus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plus: invalid argument count")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.plus: unit is null")
	}
	self := params[0].(*object.Object)
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plus: invalid amount")
	}
	unit := params[2].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusUnit([]any{ldt, amount, unit})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusYears(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusYears: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusYears([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusMonths(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusMonths: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusMonths([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusWeeks(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusWeeks: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusWeeks([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusDays(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusDays: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusDaysMethod([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusHours(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusHours: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusHours([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusMinutes(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusMinutes: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusMinutes([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusSeconds(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusSeconds: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusSeconds([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

func offsetDateTimePlusNanos(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.plusNanos: invalid argument count")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	newLdtRes := localDateTimePlusNanos([]any{ldt, params[1]})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return createOffsetDateTime(newLdtRes.(*object.Object), offset)
}

// Minus

func offsetDateTimeMinusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.minus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)
	return temporalAmountSubtractFrom([]any{amount, self})
}

func offsetDateTimeMinusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.minus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minus: invalid argument count")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minus: invalid amount")
	}
	if amount == math.MinInt64 {
		return offsetDateTimePlusUnit([]any{offsetDateTimePlusUnit([]any{params[0], -(amount + 1), params[2]}), int64(1), params[2]})
	}
	return offsetDateTimePlusUnit([]any{params[0], -amount, params[2]})
}

func offsetDateTimeMinusYears(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusYears: invalid argument count")
	}
	years, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusYears: invalid argument")
	}
	if years == math.MinInt64 {
		return offsetDateTimePlusYears([]any{offsetDateTimePlusYears([]any{params[0], -(years + 1)}), int64(1)})
	}
	return offsetDateTimePlusYears([]any{params[0], -years})
}

func offsetDateTimeMinusMonths(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusMonths: invalid argument count")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusMonths: invalid argument")
	}
	if months == math.MinInt64 {
		return offsetDateTimePlusMonths([]any{offsetDateTimePlusMonths([]any{params[0], -(months + 1)}), int64(1)})
	}
	return offsetDateTimePlusMonths([]any{params[0], -months})
}

func offsetDateTimeMinusWeeks(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusWeeks: invalid argument count")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusWeeks: invalid argument")
	}
	if weeks == math.MinInt64 {
		return offsetDateTimePlusWeeks([]any{offsetDateTimePlusWeeks([]any{params[0], -(weeks + 1)}), int64(1)})
	}
	return offsetDateTimePlusWeeks([]any{params[0], -weeks})
}

func offsetDateTimeMinusDays(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusDays: invalid argument count")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusDays: invalid argument")
	}
	if days == math.MinInt64 {
		return offsetDateTimePlusDays([]any{offsetDateTimePlusDays([]any{params[0], -(days + 1)}), int64(1)})
	}
	return offsetDateTimePlusDays([]any{params[0], -days})
}

func offsetDateTimeMinusHours(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusHours: invalid argument count")
	}
	hours, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusHours: invalid argument")
	}
	if hours == math.MinInt64 {
		return offsetDateTimePlusHours([]any{offsetDateTimePlusHours([]any{params[0], -(hours + 1)}), int64(1)})
	}
	return offsetDateTimePlusHours([]any{params[0], -hours})
}

func offsetDateTimeMinusMinutes(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusMinutes: invalid argument count")
	}
	mins, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusMinutes: invalid argument")
	}
	if mins == math.MinInt64 {
		return offsetDateTimePlusMinutes([]any{offsetDateTimePlusMinutes([]any{params[0], -(mins + 1)}), int64(1)})
	}
	return offsetDateTimePlusMinutes([]any{params[0], -mins})
}

func offsetDateTimeMinusSeconds(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusSeconds: invalid argument count")
	}
	secs, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusSeconds: invalid argument")
	}
	if secs == math.MinInt64 {
		return offsetDateTimePlusSeconds([]any{offsetDateTimePlusSeconds([]any{params[0], -(secs + 1)}), int64(1)})
	}
	return offsetDateTimePlusSeconds([]any{params[0], -secs})
}

func offsetDateTimeMinusNanos(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusNanos: invalid argument count")
	}
	nanos, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "OffsetDateTime.minusNanos: invalid argument")
	}
	if nanos == math.MinInt64 {
		return offsetDateTimePlusNanos([]any{offsetDateTimePlusNanos([]any{params[0], -(nanos + 1)}), int64(1)})
	}
	return offsetDateTimePlusNanos([]any{params[0], -nanos})
}

// Queries, AdjustInto, Until, Format

func offsetDateTimeQuery(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	ensureTemporalQueriesInited()
	if query == offsetQueryInstance || query == zoneQueryInstance {
		return getOffsetDateTimeOffset(self)
	}
	if query == zoneIdQueryInstance {
		return nil
	}
	if query == localDateQueryInstance {
		return offsetDateTimeToLocalDate([]any{self})
	}
	if query == localTimeQueryInstance {
		return offsetDateTimeToLocalTime([]any{self})
	}
	if query == chronologyQueryInstance {
		return GetIsoChronology()
	}
	if query == precisionQueryInstance {
		return ghelpers.CreateTemporalUnit("NANOS")
	}

	var queryName string
	if query.FieldTable != nil {
		if nameFld, ok := query.FieldTable["name"]; ok {
			if strObj, okStr := nameFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				queryName = object.GoStringFromStringObject(strObj)
			}
		}
	}

	switch queryName {
	case "offset", "zone":
		return getOffsetDateTimeOffset(self)
	case "zoneId":
		return nil
	case "localDate":
		return offsetDateTimeToLocalDate([]any{self})
	case "localTime":
		return offsetDateTimeToLocalTime([]any{self})
	case "chronology":
		return GetIsoChronology()
	case "precision":
		return ghelpers.CreateTemporalUnit("NANOS")
	default:
		return temporalQueryQueryFrom([]any{query, self})
	}
}

func offsetDateTimeAdjustInto(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.adjustInto: target is null")
	}
	self := params[0].(*object.Object)
	target := params[1].(*object.Object)

	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)

	adjRes := localDateTimeAdjustInto([]any{ldt, target})
	if err, isErr := adjRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	adjTarget := adjRes.(*object.Object)
	return zoneOffsetAdjustInto([]any{offset, adjTarget})
}

func offsetDateTimeUntil(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.until: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.until: unit is null")
	}
	self := params[0].(*object.Object)
	endExclusive := params[1].(*object.Object)
	unitObj := params[2].(*object.Object)

	endOdtRes := offsetDateTimeFrom([]any{endExclusive})
	if err, isErr := endOdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	endOdt := endOdtRes.(*object.Object)

	offset := getOffsetDateTimeOffset(self)
	endOdtSameInstantRes := offsetDateTimeWithOffsetSameInstant([]any{endOdt, offset})
	if err, isErr := endOdtSameInstantRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	endOdtSameInstant := endOdtSameInstantRes.(*object.Object)

	unitName := getTemporalUnitName(unitObj)
	switch strings.ToUpper(unitName) {
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		ldt1 := getOffsetDateTimeDateTime(self)
		ldt2 := getOffsetDateTimeDateTime(endOdtSameInstant)
		return localDateTimeUntil([]any{ldt1, ldt2, unitObj})
	default:
		sec1 := getOffsetDateTimeEpochSecond(self)
		nano1 := getOffsetDateTimeNano(self)
		sec2 := getOffsetDateTimeEpochSecond(endOdt)
		nano2 := getOffsetDateTimeNano(endOdt)

		totalSecDiff := sec2 - sec1
		totalNanoDiff := nano2 - nano1

		switch strings.ToUpper(unitName) {
		case "NANOS":
			return totalSecDiff*nanosPerSecond + totalNanoDiff
		case "MICROS":
			return totalSecDiff*1000000 + totalNanoDiff/1000
		case "MILLIS":
			return totalSecDiff*1000 + totalNanoDiff/1000000
		case "SECONDS":
			return totalSecDiff
		case "MINUTES":
			return totalSecDiff / 60
		case "HOURS":
			return totalSecDiff / 3600
		case "HALF_DAYS":
			return totalSecDiff / (3600 * 12)
		default:
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
		}
	}
}

func offsetDateTimeFormat(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.format: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.format: formatter is null")
	}
	self := params[0].(*object.Object)
	formatter := params[1].(*object.Object)
	return dateTimeFormatterFormat([]any{formatter, self})
}

// Conversions

func offsetDateTimeAtZoneSameInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.atZoneSameInstant: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.atZoneSameInstant: zone is null")
	}
	self := params[0].(*object.Object)
	zone := params[1].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	return zonedDateTimeOfInstantLdtOffsetZone(ldt, offset, zone)
}

func offsetDateTimeAtZoneSimilarLocal(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.atZoneSimilarLocal: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.atZoneSimilarLocal: zone is null")
	}
	self := params[0].(*object.Object)
	zone := params[1].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	return zonedDateTimeOfLocal(ldt, zone, offset)
}

func offsetDateTimeToOffsetTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toOffsetTime: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	localTimeObj := localDateTimeToLocalTime([]any{ldt}).(*object.Object)
	return createOffsetTime(localTimeObj, offset)
}

func offsetDateTimeToEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toEpochSecond: self is null")
	}
	self := params[0].(*object.Object)
	return getOffsetDateTimeEpochSecond(self)
}

func offsetDateTimeToInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toInstant: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	return localDateTimeToInstant([]any{ldt, offset})
}

func offsetDateTimeToZonedDateTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toZonedDateTime: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	offset := getOffsetDateTimeOffset(self)
	return zonedDateTimeOfLDT([]any{ldt, offset})
}

// Comparisons and Object Methods

func offsetDateTimeCompareTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	cmp := compareInstant(self, other)
	if cmp == 0 {
		ldt1 := getOffsetDateTimeDateTime(self)
		ldt2 := getOffsetDateTimeDateTime(other)
		return localDateTimeCompareTo([]any{ldt1, ldt2})
	}
	return cmp
}

func offsetDateTimeIsAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isAfter: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isAfter: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	sec1 := getOffsetDateTimeEpochSecond(self)
	sec2 := getOffsetDateTimeEpochSecond(other)
	if sec1 > sec2 {
		return types.JavaBoolTrue
	}
	if sec1 == sec2 && getOffsetDateTimeNano(self) > getOffsetDateTimeNano(other) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func offsetDateTimeIsBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isBefore: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isBefore: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	sec1 := getOffsetDateTimeEpochSecond(self)
	sec2 := getOffsetDateTimeEpochSecond(other)
	if sec1 < sec2 {
		return types.JavaBoolTrue
	}
	if sec1 == sec2 && getOffsetDateTimeNano(self) < getOffsetDateTimeNano(other) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func offsetDateTimeIsEqual(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isEqual: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.isEqual: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	sec1 := getOffsetDateTimeEpochSecond(self)
	sec2 := getOffsetDateTimeEpochSecond(other)
	if sec1 == sec2 && getOffsetDateTimeNano(self) == getOffsetDateTimeNano(other) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func offsetDateTimeEquals(params []any) any {
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
	if !strings.Contains(getObjectClassName(other), "OffsetDateTime") {
		return types.JavaBoolFalse
	}
	ldt1 := getOffsetDateTimeDateTime(self)
	ldt2 := getOffsetDateTimeDateTime(other)
	off1 := getOffsetDateTimeOffset(self)
	off2 := getOffsetDateTimeOffset(other)

	eqLdt := localDateTimeEquals([]any{ldt1, ldt2})
	eqOff := zoneOffsetEquals([]any{off1, off2})

	if eqLdt == types.JavaBoolTrue && eqOff == types.JavaBoolTrue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func offsetDateTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return int64(0)
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	off := getOffsetDateTimeOffset(self)

	h1 := localDateTimeHashCode([]any{ldt}).(int64)
	h2 := zoneOffsetHashCode([]any{off}).(int64)
	return h1 ^ h2
}

func offsetDateTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "OffsetDateTime.toString: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getOffsetDateTimeDateTime(self)
	off := getOffsetDateTimeOffset(self)

	ldtStrObj := localDateTimeToString([]any{ldt}).(*object.Object)
	offStrObj := zoneOffsetToString([]any{off}).(*object.Object)

	ldtStr := object.GoStringFromStringObject(ldtStrObj)
	offStr := object.GoStringFromStringObject(offStrObj)

	return object.StringObjectFromGoString(ldtStr + offStr)
}

func offsetDateTimeWriteExternal(params []any) any {
	return nil
}

func offsetDateTimeReadExternal(params []any) any {
	return offsetDateTimeNow(nil)
}

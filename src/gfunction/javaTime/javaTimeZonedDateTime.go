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
	"jacobin/src/types"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	zonedDateTimeMutex       sync.Mutex
	zonedDateTimeInitialized bool = false
)

func ensureZonedDateTimeInited() {
	zonedDateTimeMutex.Lock()
	defer zonedDateTimeMutex.Unlock()
	if zonedDateTimeInitialized {
		return
	}
	ensureLocalDateTimeInited()
	ensureZoneOffsetInited()
	ensureZoneIdInited()
	zonedDateTimeInitialized = true
}

// Load_Time_ZonedDateTime registers method signatures for java/time/ZonedDateTime.
func Load_Time_ZonedDateTime() {
	ensureZonedDateTimeInited()

	ghelpers.MethodSignatures["java/time/ZonedDateTime.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeClinit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeInit,
		}

	// Static Factories
	ghelpers.MethodSignatures["java/time/ZonedDateTime.now()Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeNow,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.now(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeNowZone,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.now(Ljava/time/Clock;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeNowClock,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.of(Ljava/time/LocalDate;Ljava/time/LocalTime;Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zonedDateTimeOfDate,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.of(Ljava/time/LocalDateTime;Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeOfLDT,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.of(IIIIIIILjava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 8,
			GFunction:  zonedDateTimeOfInts,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.ofLocal(Ljava/time/LocalDateTime;Ljava/time/ZoneId;Ljava/time/ZoneOffset;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zonedDateTimeOfLocalMethod,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeOfInstantMethod,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.ofInstant(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zonedDateTimeOfInstantLdtOffsetZoneMethod,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.ofStrict(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zonedDateTimeOfStrictMethod,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeFrom,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.parse(Ljava/lang/CharSequence;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeParse,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeParseFormatter,
		}

	// Instance Queries & Getters
	ghelpers.MethodSignatures["java/time/ZonedDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeRange,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeGet,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeGetLong,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getOffset()Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetOffset,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getZone()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetZone,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetChronology,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getChronology()Ljava/time/chrono/IsoChronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetChronology,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toLocalDateTime()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToLocalDateTime,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toLocalDateTime()Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToLocalDateTime,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toLocalDate()Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toLocalDate()Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToLocalDate,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetYear,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getMonthValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetMonthValue,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getMonth()Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetMonth,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getDayOfMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getDayOfYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getDayOfWeek()Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetDayOfWeek,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toLocalTime()Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToLocalTime,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getHour()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetHour,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getMinute()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetMinute,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getSecond()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetSecond,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.getNano()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeGetNano,
		}

	// Overlaps and Zone Adjustments
	ghelpers.MethodSignatures["java/time/ZonedDateTime.withEarlierOffsetAtOverlap()Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeWithEarlierOffsetAtOverlap,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withEarlierOffsetAtOverlap()Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeWithEarlierOffsetAtOverlap,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withLaterOffsetAtOverlap()Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeWithLaterOffsetAtOverlap,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withLaterOffsetAtOverlap()Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeWithLaterOffsetAtOverlap,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withZoneSameLocal(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithZoneSameLocal,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withZoneSameLocal(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithZoneSameLocal,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withZoneSameInstant(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithZoneSameInstant,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withZoneSameInstant(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithZoneSameInstant,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withFixedOffsetZone()Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeWithFixedOffsetZone,
		}

	// with Modifications
	ghelpers.MethodSignatures["java/time/ZonedDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withYear(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithYear,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withMonth(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithMonth,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withDayOfMonth(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withDayOfYear(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withHour(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithHour,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withMinute(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithMinute,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withSecond(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithSecond,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.withNano(I)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeWithNano,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeTruncatedTo,
		}

	// Plus operations
	ghelpers.MethodSignatures["java/time/ZonedDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusYears(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusYears,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusMonths(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusMonths,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusWeeks(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusWeeks,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusDays(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusDays,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusHours(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusHours,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusMinutes(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusMinutes,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusSeconds(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusSeconds,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.plusNanos(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimePlusNanos,
		}

	// Minus operations
	ghelpers.MethodSignatures["java/time/ZonedDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusYears(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusYears,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusMonths(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusMonths,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusWeeks(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusWeeks,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusDays(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusDays,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusHours(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusHours,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusMinutes(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusMinutes,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusSeconds(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusSeconds,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.minusNanos(J)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeMinusNanos,
		}

	// Queries, Adjustments & Conversions
	ghelpers.MethodSignatures["java/time/ZonedDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeQuery,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zonedDateTimeUntil,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeFormat,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toOffsetDateTime()Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToOffsetDateTime,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toInstant()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToInstant,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toEpochSecond()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToEpochSecond,
		}

	// Comparisons, Equality & String
	ghelpers.MethodSignatures["java/time/ZonedDateTime.compareTo(Ljava/time/chrono/ChronoZonedDateTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.isAfter(Ljava/time/chrono/ChronoZonedDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeIsAfter,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.isBefore(Ljava/time/chrono/ChronoZonedDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeIsBefore,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.isEqual(Ljava/time/chrono/ChronoZonedDateTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeIsEqual,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zonedDateTimeEquals,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeHashCode,
		}

	ghelpers.MethodSignatures["java/time/ZonedDateTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zonedDateTimeToString,
		}
}

func zonedDateTimeClinit([]any) any {
	ensureZonedDateTimeInited()
	return nil
}

func zonedDateTimeInit([]any) any {
	return nil
}

// Internal constructor
func createZonedDateTime(ldt, offset, zone *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameZonedDateTimeVal)
	obj.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	obj.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offset}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zone}
	return obj
}

// CreateZonedDateTime is an exported helper for other packages.
func CreateZonedDateTime(ldt, offset, zone *object.Object) *object.Object {
	return createZonedDateTime(ldt, offset, zone)
}

func getZonedDateTimeLocalDateTime(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["dateTime"]; ok {
		if ldt, okLdt := fld.Fvalue.(*object.Object); okLdt && ldt != nil {
			return ldt
		}
	}
	return nil
}

func getZonedDateTimeOffset(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["offset"]; ok {
		if off, okOff := fld.Fvalue.(*object.Object); okOff && off != nil {
			return off
		}
	}
	return nil
}

func getZonedDateTimeZone(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["zone"]; ok {
		if z, okZ := fld.Fvalue.(*object.Object); okZ && z != nil {
			return z
		}
	}
	return nil
}

func getZonedDateTimeEpochSecond(obj *object.Object) int64 {
	ldt := getZonedDateTimeLocalDateTime(obj)
	offset := getZonedDateTimeOffset(obj)
	dObj := getLocalDateTimeDate(ldt)
	tObj := getLocalDateTimeTime(ldt)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(dObj)
	secOfDay := getLocalTimeNanoOfDay(tObj) / nanosPerSecond
	offsetSec := getZoneOffsetTotalSeconds(offset)
	return epochDay*secondsPerDay + secOfDay - offsetSec
}

func getZonedDateTimeNano(obj *object.Object) int64 {
	ldt := getZonedDateTimeLocalDateTime(obj)
	tObj := getLocalDateTimeTime(ldt)
	_, _, _, nano := getLocalTimeHMSN(tObj)
	return nano
}

func getZonedDateTimeYMDHMSN(obj *object.Object) (year, month, day, hour, min, sec, nano int64) {
	ldt := getZonedDateTimeLocalDateTime(obj)
	if ldt != nil {
		dObj := getLocalDateTimeDate(ldt)
		tObj := getLocalDateTimeTime(ldt)
		year, month, day, _ = getLocalDateYMD(dObj)
		hour, min, sec, nano = getLocalTimeHMSN(tObj)
		return
	}
	return 1970, 1, 1, 0, 0, 0, 0
}

func zonedDateTimeOfLocal(ldt, zone, preferredOffset *object.Object) any {
	if ldt == nil || object.IsNull(ldt) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofLocal: localDateTime is null")
	}
	if zone == nil || object.IsNull(zone) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofLocal: zone is null")
	}
	zoneKName := object.GoStringFromStringPoolIndex(zone.KlassName)
	if zoneKName == classNameZoneOffsetVal {
		return createZonedDateTime(ldt, zone, zone)
	}

	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)

	validOffsetsRes := zoneRulesGetValidOffsets([]any{rulesObj, ldt})
	if err, isErr := validOffsetsRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	listObj := validOffsetsRes.(*object.Object)
	var validOffsets []*object.Object
	if elFld, ok := listObj.FieldTable["elements"]; ok {
		if slice, ok2 := elFld.Fvalue.([]any); ok2 {
			for _, item := range slice {
				if offObj, ok3 := item.(*object.Object); ok3 {
					validOffsets = append(validOffsets, offObj)
				}
			}
		}
	}

	var offset *object.Object
	if len(validOffsets) == 1 {
		offset = validOffsets[0]
	} else if len(validOffsets) == 0 {
		transRes := zoneRulesGetTransition([]any{rulesObj, ldt})
		if transObj, ok := transRes.(*object.Object); ok && transObj != nil && !object.IsNull(transObj) {
			durRes := zoneOffsetTransitionGetDuration([]any{transObj})
			if durObj, okDur := durRes.(*object.Object); okDur && durObj != nil {
				secFld := durObj.FieldTable["seconds"]
				if sec, okSec := getInt64Val(secFld.Fvalue); okSec {
					newLdtRes := localDateTimePlusSeconds([]any{ldt, sec})
					if _, isErr := newLdtRes.(*ghelpers.GErrBlk); !isErr {
						ldt = newLdtRes.(*object.Object)
					}
				}
			}
			offAfterRes := zoneOffsetTransitionGetOffsetAfter([]any{transObj})
			if offAfterObj, okOff := offAfterRes.(*object.Object); okOff && offAfterObj != nil {
				offset = offAfterObj
			}
		}
		if offset == nil {
			offsetRes := zoneRulesGetOffsetLocalDateTime([]any{rulesObj, ldt})
			if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
				return err
			}
			offset = offsetRes.(*object.Object)
		}
	} else {
		if preferredOffset != nil && !object.IsNull(preferredOffset) {
			prefSec := getZoneOffsetTotalSeconds(preferredOffset)
			for _, vOff := range validOffsets {
				if getZoneOffsetTotalSeconds(vOff) == prefSec {
					offset = vOff
					break
				}
			}
		}
		if offset == nil {
			offset = validOffsets[0]
		}
	}
	return createZonedDateTime(ldt, offset, zone)
}

func zonedDateTimeOfInstant(instant, zone *object.Object) any {
	if instant == nil || object.IsNull(instant) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: instant is null")
	}
	if zone == nil || object.IsNull(zone) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: zone is null")
	}
	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)
	offsetRes := zoneRulesGetOffsetInstant([]any{rulesObj, instant})
	if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	offsetObj := offsetRes.(*object.Object)

	epochSec, nano := getInstantSecondsAndNanos(instant)
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)
	localSec := epochSec + offsetSec
	epochDay := floorDiv(localSec, secondsPerDay)
	secOfDay := floorMod(localSec, secondsPerDay)

	y, m, d := epochDayToYMD(epochDay)
	dateObj := CreateLocalDate(y, m, d)
	timeObj := ofNanoOfDay(secOfDay*nanosPerSecond + nano)
	ldt := createLocalDateTime(dateObj, timeObj)
	return createZonedDateTime(ldt, offsetObj, zone)
}

func zonedDateTimeOfInstantLdtOffsetZone(ldt, offset, zone *object.Object) any {
	if ldt == nil || object.IsNull(ldt) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: localDateTime is null")
	}
	if offset == nil || object.IsNull(offset) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: offset is null")
	}
	if zone == nil || object.IsNull(zone) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: zone is null")
	}
	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)
	validRes := zoneRulesIsValidOffset([]any{rulesObj, ldt, offset})
	if validRes == types.JavaBoolTrue {
		return createZonedDateTime(ldt, offset, zone)
	}
	dObj := getLocalDateTimeDate(ldt)
	tObj := getLocalDateTimeTime(ldt)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(dObj)
	secOfDay := getLocalTimeNanoOfDay(tObj) / nanosPerSecond
	offsetSec := getZoneOffsetTotalSeconds(offset)
	sec := epochDay*secondsPerDay + secOfDay - offsetSec
	_, _, _, nano := getLocalTimeHMSN(tObj)
	inst := createInstant(sec, nano)
	return zonedDateTimeOfInstant(inst, zone)
}

func zonedDateTimeOfStrict(ldt, offset, zone *object.Object) any {
	if ldt == nil || object.IsNull(ldt) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofStrict: localDateTime is null")
	}
	if offset == nil || object.IsNull(offset) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofStrict: offset is null")
	}
	if zone == nil || object.IsNull(zone) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofStrict: zone is null")
	}
	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)
	validRes := zoneRulesIsValidOffset([]any{rulesObj, ldt, offset})
	if err, isErr := validRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if validRes == types.JavaBoolFalse {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "ZoneRules: offset is invalid for local date time")
	}
	return createZonedDateTime(ldt, offset, zone)
}

// Static Factory Implementations

func zonedDateTimeNow([]any) any {
	defZone := zoneIdSystemDefault(nil)
	if zObj, ok := defZone.(*object.Object); ok && zObj != nil && !object.IsNull(zObj) {
		return zonedDateTimeNowZone([]any{zObj})
	}
	zoneUTC := makeZoneIdFromID("UTC")
	return zonedDateTimeNowZone([]any{zoneUTC})
}

func zonedDateTimeNowZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.now: zone is null")
	}
	zoneObj := params[0].(*object.Object)
	now := time.Now().UTC()
	inst := createInstant(now.Unix(), int64(now.Nanosecond()))
	return zonedDateTimeOfInstant(inst, zoneObj)
}

func zonedDateTimeNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.now: clock is null")
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
	return zonedDateTimeOfInstant(instRes.(*object.Object), zoneRes.(*object.Object))
}

func zonedDateTimeOfDate(params []any) any {
	if len(params) < 3 || params[0] == nil || params[1] == nil || params[2] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.of: null argument")
	}
	dateObj := params[0].(*object.Object)
	timeObj := params[1].(*object.Object)
	zoneObj := params[2].(*object.Object)
	ldt := createLocalDateTime(dateObj, timeObj)
	return zonedDateTimeOfLocal(ldt, zoneObj, nil)
}

func zonedDateTimeOfLDT(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.of: null argument")
	}
	ldt := params[0].(*object.Object)
	zoneObj := params[1].(*object.Object)
	return zonedDateTimeOfLocal(ldt, zoneObj, nil)
}

func zonedDateTimeOfInts(params []any) any {
	if len(params) < 8 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.of: missing arguments")
	}
	y, okY := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	d, okD := getInt64Val(params[2])
	h, okH := getInt64Val(params[3])
	min, okMin := getInt64Val(params[4])
	s, okS := getInt64Val(params[5])
	nano, okN := getInt64Val(params[6])
	if !okY || !okM || !okD || !okH || !okMin || !okS || !okN {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.of: invalid integer parameters")
	}
	if params[7] == nil || object.IsNull(params[7].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.of: zone is null")
	}
	zoneObj := params[7].(*object.Object)
	ldt := createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, s, nano))
	return zonedDateTimeOfLocal(ldt, zoneObj, nil)
}

func zonedDateTimeOfLocalMethod(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofLocal: null argument")
	}
	ldt := params[0].(*object.Object)
	zone := params[1].(*object.Object)
	var prefOffset *object.Object
	if len(params) >= 3 && params[2] != nil && !object.IsNull(params[2].(*object.Object)) {
		prefOffset = params[2].(*object.Object)
	}
	return zonedDateTimeOfLocal(ldt, zone, prefOffset)
}

func zonedDateTimeOfInstantMethod(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: null argument")
	}
	instant := params[0].(*object.Object)
	zone := params[1].(*object.Object)
	return zonedDateTimeOfInstant(instant, zone)
}

func zonedDateTimeOfInstantLdtOffsetZoneMethod(params []any) any {
	if len(params) < 3 || params[0] == nil || params[1] == nil || params[2] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofInstant: null argument")
	}
	ldt := params[0].(*object.Object)
	offset := params[1].(*object.Object)
	zone := params[2].(*object.Object)
	return zonedDateTimeOfInstantLdtOffsetZone(ldt, offset, zone)
}

func zonedDateTimeOfStrictMethod(params []any) any {
	if len(params) < 3 || params[0] == nil || params[1] == nil || params[2] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.ofStrict: null argument")
	}
	ldt := params[0].(*object.Object)
	offset := params[1].(*object.Object)
	zone := params[2].(*object.Object)
	return zonedDateTimeOfStrict(ldt, offset, zone)
}

func zonedDateTimeFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.from: temporal is null")
	}
	temporal := params[0].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(temporal.KlassName)
	if kName == classNameZonedDateTimeVal {
		return temporal
	}

	var zoneObj *object.Object
	ensureTemporalQueriesInited()
	zoneRes := temporalAccessorQuery([]any{temporal, zoneQueryInstance})
	if z, ok := zoneRes.(*object.Object); ok && z != nil && !object.IsNull(z) {
		zoneObj = z
	} else {
		zoneIdRes := temporalAccessorQuery([]any{temporal, zoneIdQueryInstance})
		if z, ok := zoneIdRes.(*object.Object); ok && z != nil && !object.IsNull(z) {
			zoneObj = z
		}
	}

	if zoneObj == nil && temporal.FieldTable != nil {
		if zFld, ok := temporal.FieldTable["zone"]; ok {
			if z, ok2 := zFld.Fvalue.(*object.Object); ok2 && z != nil && !object.IsNull(z) {
				zoneObj = z
			}
		}
	}

	if zoneObj == nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain ZonedDateTime from TemporalAccessor: no ZoneId found")
	}

	if kName == classNameInstantVal || (temporal.FieldTable != nil && temporal.FieldTable["seconds"].Fvalue != nil) {
		sec, nano := getInstantSecondsAndNanos(temporal)
		inst := createInstant(sec, nano)
		return zonedDateTimeOfInstant(inst, zoneObj)
	}

	if kName == classNameLocalDateTime {
		return zonedDateTimeOfLocal(temporal, zoneObj, nil)
	}

	cfInstantSec := getChronoFieldForInstant("INSTANT_SECONDS")
	cfNanoOfSec := getChronoFieldForInstant("NANO_OF_SECOND")
	if cfInstantSec != nil {
		isSupp := temporalAccessorIsSupported([]any{temporal, cfInstantSec})
		if isSupp == types.JavaBoolTrue {
			secRes := temporalAccessorGetLong([]any{temporal, cfInstantSec})
			if secVal, ok := getInt64Val(secRes); ok {
				var nanoVal int64 = 0
				if cfNanoOfSec != nil {
					isSuppNano := temporalAccessorIsSupported([]any{temporal, cfNanoOfSec})
					if isSuppNano == types.JavaBoolTrue {
						nanoRes := temporalAccessorGetLong([]any{temporal, cfNanoOfSec})
						if n, okN := getInt64Val(nanoRes); okN {
							nanoVal = n
						}
					}
				}
				inst := createInstant(secVal, nanoVal)
				return zonedDateTimeOfInstant(inst, zoneObj)
			}
		}
	}

	ldtRes := localDateTimeFrom([]any{temporal})
	if ldt, ok := ldtRes.(*object.Object); ok && ldt != nil && !object.IsNull(ldt) {
		return zonedDateTimeOfLocal(ldt, zoneObj, nil)
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain ZonedDateTime from TemporalAccessor: "+kName)
}

func parseZonedDateTimeString(text string) any {
	text = strings.TrimSpace(text)
	if text == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text cannot be empty")
	}

	var zoneIDStr string
	mainText := text
	if idx := strings.Index(text, "["); idx >= 0 {
		endIdx := strings.Index(text, "]")
		if endIdx < idx {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Invalid zone id bracket in text: "+text)
		}
		zoneIDStr = text[idx+1 : endIdx]
		mainText = text[:idx]
	}

	// Look for 'T' separator
	tIdx := strings.Index(mainText, "T")
	if tIdx < 0 {
		tIdx = strings.Index(mainText, "t")
	}
	if tIdx < 0 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed: missing 'T' separator", text))
	}

	datePart := mainText[:tIdx]
	timeWithOffset := mainText[tIdx+1:]

	y, m, d, errDate := parseLocalDateString(datePart)
	if errDate != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: "+errDate.Error())
	}

	var timeOnly string
	var offsetStr string

	if strings.HasSuffix(timeWithOffset, "Z") || strings.HasSuffix(timeWithOffset, "z") {
		timeOnly = timeWithOffset[:len(timeWithOffset)-1]
		offsetStr = "Z"
	} else {
		plusIdx := strings.LastIndex(timeWithOffset, "+")
		minusIdx := strings.LastIndex(timeWithOffset, "-")
		offIdx := -1
		if plusIdx > offIdx {
			offIdx = plusIdx
		}
		if minusIdx > offIdx {
			offIdx = minusIdx
		}

		if offIdx >= 0 {
			timeOnly = timeWithOffset[:offIdx]
			offsetStr = timeWithOffset[offIdx:]
		} else {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed: missing zone offset", text))
		}
	}

	h, min, s, n, errTime := parseLocalTimeString(timeOnly)
	if errTime != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: "+errTime.Error())
	}

	offsetObjRes := zoneOffsetOf([]any{object.StringObjectFromGoString(offsetStr)})
	if err, isErr := offsetObjRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	offsetObj := offsetObjRes.(*object.Object)

	ldt := createLocalDateTime(CreateLocalDate(y, m, d), createLocalTime(h, min, s, n))

	if zoneIDStr != "" {
		zoneObj := makeZoneIdFromID(zoneIDStr)
		return zonedDateTimeOfLocal(ldt, zoneObj, offsetObj)
	}

	return createZonedDateTime(ldt, offsetObj, offsetObj)
}

func zonedDateTimeParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.parse: text is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.parse: invalid string")
	}
	text := object.GoStringFromStringObject(strObj)
	return parseZonedDateTimeString(text)
}

func zonedDateTimeParseFormatter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.parse: text is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.parse: formatter is null")
	}
	textObj := params[0].(*object.Object)
	formatterObj := params[1].(*object.Object)
	if !object.IsStringObject(textObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.parse: invalid string")
	}
	text := object.GoStringFromStringObject(textObj)
	accessor, err := parseTextToTemporalAccessor(formatterObj, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, err.Error())
	}
	return zonedDateTimeFrom([]any{accessor})
}

// Instance Methods

func zonedDateTimeIsSupportedField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.isSupported: self is null")
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

func zonedDateTimeIsSupportedUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.isSupported: self is null")
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

func zonedDateTimeRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.range: field is null")
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
		ldt := getZonedDateTimeLocalDateTime(self)
		return localDateTimeRange([]any{ldt, fieldObj})
	}
}

func zonedDateTimeGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.get: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "INSTANT_SECONDS":
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Field 'InstantSeconds' exceeds int range, use getLong()")
	case "OFFSET_SECONDS":
		off := getZonedDateTimeOffset(self)
		return int64(getZoneOffsetTotalSeconds(off))
	default:
		ldt := getZonedDateTimeLocalDateTime(self)
		return localDateTimeGet([]any{ldt, fieldObj})
	}
}

func zonedDateTimeGetLong(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getLong: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: null")
	}
	switch strings.ToUpper(fieldName) {
	case "INSTANT_SECONDS":
		return getZonedDateTimeEpochSecond(self)
	case "OFFSET_SECONDS":
		off := getZonedDateTimeOffset(self)
		return int64(getZoneOffsetTotalSeconds(off))
	default:
		ldt := getZonedDateTimeLocalDateTime(self)
		return localDateTimeGetLong([]any{ldt, fieldObj})
	}
}

func zonedDateTimeGetOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getOffset: self is null")
	}
	return getZonedDateTimeOffset(params[0].(*object.Object))
}

func zonedDateTimeGetZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getZone: self is null")
	}
	return getZonedDateTimeZone(params[0].(*object.Object))
}

func zonedDateTimeGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getChronology: self is null")
	}
	return GetIsoChronology()
}

func zonedDateTimeToLocalDateTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toLocalDateTime: self is null")
	}
	return getZonedDateTimeLocalDateTime(params[0].(*object.Object))
}

func zonedDateTimeToLocalDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toLocalDate: self is null")
	}
	ldt := getZonedDateTimeLocalDateTime(params[0].(*object.Object))
	return getLocalDateTimeDate(ldt)
}

func zonedDateTimeGetYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getYear: self is null")
	}
	y, _, _, _, _, _, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return y
}

func zonedDateTimeGetMonthValue(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getMonthValue: self is null")
	}
	_, m, _, _, _, _, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return m
}

func zonedDateTimeGetMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getMonth: self is null")
	}
	_, m, _, _, _, _, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	ensureMonthInited()
	if m >= 1 && m <= 12 {
		return monthInstances[m-1]
	}
	return nil
}

func zonedDateTimeGetDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getDayOfMonth: self is null")
	}
	_, _, d, _, _, _, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return d
}

func zonedDateTimeGetDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getDayOfYear: self is null")
	}
	ldt := getZonedDateTimeLocalDateTime(params[0].(*object.Object))
	dObj := getLocalDateTimeDate(ldt)
	return localDateGetDayOfYear([]any{dObj})
}

func zonedDateTimeGetDayOfWeek(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getDayOfWeek: self is null")
	}
	ldt := getZonedDateTimeLocalDateTime(params[0].(*object.Object))
	dObj := getLocalDateTimeDate(ldt)
	return localDateGetDayOfWeek([]any{dObj})
}

func zonedDateTimeToLocalTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toLocalTime: self is null")
	}
	ldt := getZonedDateTimeLocalDateTime(params[0].(*object.Object))
	return getLocalDateTimeTime(ldt)
}

func zonedDateTimeGetHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getHour: self is null")
	}
	_, _, _, h, _, _, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return h
}

func zonedDateTimeGetMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getMinute: self is null")
	}
	_, _, _, _, min, _, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return min
}

func zonedDateTimeGetSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getSecond: self is null")
	}
	_, _, _, _, _, s, _ := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return s
}

func zonedDateTimeGetNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.getNano: self is null")
	}
	_, _, _, _, _, _, n := getZonedDateTimeYMDHMSN(params[0].(*object.Object))
	return n
}

func zonedDateTimeWithEarlierOffsetAtOverlap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withEarlierOffsetAtOverlap: self is null")
	}
	self := params[0].(*object.Object)
	zone := getZonedDateTimeZone(self)
	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)
	if isZoneRulesFixed(rulesObj) {
		return self
	}

	ldt := getZonedDateTimeLocalDateTime(self)
	validOffsetsRes := zoneRulesGetValidOffsets([]any{rulesObj, ldt})
	if err, isErr := validOffsetsRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	listObj := validOffsetsRes.(*object.Object)
	var validOffsets []*object.Object
	if elFld, ok := listObj.FieldTable["elements"]; ok {
		if slice, ok2 := elFld.Fvalue.([]any); ok2 {
			for _, item := range slice {
				if offObj, ok3 := item.(*object.Object); ok3 {
					validOffsets = append(validOffsets, offObj)
				}
			}
		}
	}

	if len(validOffsets) > 1 {
		return createZonedDateTime(ldt, validOffsets[0], zone)
	}
	return self
}

func zonedDateTimeWithLaterOffsetAtOverlap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withLaterOffsetAtOverlap: self is null")
	}
	self := params[0].(*object.Object)
	zone := getZonedDateTimeZone(self)
	rulesRes := zoneIdGetRules([]any{zone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)
	if isZoneRulesFixed(rulesObj) {
		return self
	}

	ldt := getZonedDateTimeLocalDateTime(self)
	validOffsetsRes := zoneRulesGetValidOffsets([]any{rulesObj, ldt})
	if err, isErr := validOffsetsRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	listObj := validOffsetsRes.(*object.Object)
	var validOffsets []*object.Object
	if elFld, ok := listObj.FieldTable["elements"]; ok {
		if slice, ok2 := elFld.Fvalue.([]any); ok2 {
			for _, item := range slice {
				if offObj, ok3 := item.(*object.Object); ok3 {
					validOffsets = append(validOffsets, offObj)
				}
			}
		}
	}

	if len(validOffsets) > 1 {
		return createZonedDateTime(ldt, validOffsets[len(validOffsets)-1], zone)
	}
	return self
}

func zonedDateTimeWithZoneSameLocal(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withZoneSameLocal: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withZoneSameLocal: zone is null")
	}
	self := params[0].(*object.Object)
	newZone := params[1].(*object.Object)
	curZone := getZonedDateTimeZone(self)
	if curZone == newZone {
		return self
	}
	ldt := getZonedDateTimeLocalDateTime(self)
	offset := getZonedDateTimeOffset(self)
	return zonedDateTimeOfLocal(ldt, newZone, offset)
}

func zonedDateTimeWithZoneSameInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withZoneSameInstant: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withZoneSameInstant: zone is null")
	}
	self := params[0].(*object.Object)
	newZone := params[1].(*object.Object)
	curZone := getZonedDateTimeZone(self)
	if curZone == newZone {
		return self
	}
	instRes := zonedDateTimeToInstant([]any{self})
	if err, isErr := instRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	inst := instRes.(*object.Object)
	return zonedDateTimeOfInstant(inst, newZone)
}

func zonedDateTimeWithFixedOffsetZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withFixedOffsetZone: self is null")
	}
	self := params[0].(*object.Object)
	zone := getZonedDateTimeZone(self)
	offset := getZonedDateTimeOffset(self)
	if zone == offset {
		return self
	}
	ldt := getZonedDateTimeLocalDateTime(self)
	return createZonedDateTime(ldt, offset, offset)
}

// with operations

func resolveLocal(self, newLdt *object.Object) any {
	zone := getZonedDateTimeZone(self)
	offset := getZonedDateTimeOffset(self)
	return zonedDateTimeOfLocal(newLdt, zone, offset)
}

func resolveInstant(self, newLdt *object.Object) any {
	zone := getZonedDateTimeZone(self)
	offset := getZonedDateTimeOffset(self)
	return zonedDateTimeOfInstantLdtOffsetZone(newLdt, offset, zone)
}

func zonedDateTimeWithAdjuster(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	kName := object.GoStringFromStringPoolIndex(adjuster.KlassName)
	ldt := getZonedDateTimeLocalDateTime(self)
	zone := getZonedDateTimeZone(self)

	switch kName {
	case "java/time/LocalDate":
		timeObj := getLocalDateTimeTime(ldt)
		newLdt := createLocalDateTime(adjuster, timeObj)
		return resolveLocal(self, newLdt)
	case "java/time/LocalTime":
		dateObj := getLocalDateTimeDate(ldt)
		newLdt := createLocalDateTime(dateObj, adjuster)
		return resolveLocal(self, newLdt)
	case "java/time/LocalDateTime":
		return resolveLocal(self, adjuster)
	case "java/time/Instant":
		return zonedDateTimeOfInstant(adjuster, zone)
	case "java/time/ZoneOffset":
		return zonedDateTimeOfStrict(ldt, adjuster, zone)
	case "java/time/ZoneId", "java/time/ZoneRegion":
		return zonedDateTimeWithZoneSameLocal([]any{self, adjuster})
	default:
		// TemporalAdjuster.adjustInto(Temporal)
		return adjuster
	}
}

func zonedDateTimeWithField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.with: field is null")
	}
	val, okVal := getInt64Val(params[2])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.with: invalid value")
	}

	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	zone := getZonedDateTimeZone(self)

	switch strings.ToUpper(fieldName) {
	case "INSTANT_SECONDS":
		nano := getZonedDateTimeNano(self)
		inst := createInstant(val, nano)
		return zonedDateTimeOfInstant(inst, zone)
	case "OFFSET_SECONDS":
		newOffset := createZoneOffset(val)
		ldt := getZonedDateTimeLocalDateTime(self)
		return zonedDateTimeOfStrict(ldt, newOffset, zone)
	default:
		ldt := getZonedDateTimeLocalDateTime(self)
		newLdtRes := localDateTimeWithField([]any{ldt, fieldObj, val})
		if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		newLdt := newLdtRes.(*object.Object)
		return resolveLocal(self, newLdt)
	}
}

func zonedDateTimeWithYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withYear: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withYear: invalid year")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithYear([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withMonth: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withMonth: invalid month")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithMonth([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withDayOfMonth: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withDayOfMonth: invalid day")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithDayOfMonth([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withDayOfYear: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withDayOfYear: invalid dayOfYear")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithDayOfYear([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withHour: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withHour: invalid hour")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithHour([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withMinute: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withMinute: invalid minute")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithMinute([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withSecond: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withSecond: invalid second")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithSecond([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeWithNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.withNano: self is null")
	}
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.withNano: invalid nano")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeWithNano([]any{ldt, val})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimeTruncatedTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.truncatedTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.truncatedTo: unit is null")
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimeTruncatedTo([]any{ldt, unit})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

// Plus implementations

func getDurationSecondsAndNanos(obj *object.Object) (int64, int64) {
	if obj == nil || obj.FieldTable == nil {
		return 0, 0
	}
	var sec, nano int64
	if fld, ok := obj.FieldTable["seconds"]; ok {
		sec, _ = getInt64Val(fld.Fvalue)
	}
	if fld, ok := obj.FieldTable["nanos"]; ok {
		nano, _ = getInt64Val(fld.Fvalue)
	}
	return sec, nano
}

func zonedDateTimePlusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(amount.KlassName)

	ldt := getZonedDateTimeLocalDateTime(self)
	if kName == "java/time/Period" {
		years, months, days, _ := getPeriodYMD(amount)
		newLdt := ldt
		if years != 0 {
			resY := localDateTimePlusYears([]any{newLdt, years})
			if err, isErr := resY.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resY.(*object.Object)
		}
		if months != 0 {
			resM := localDateTimePlusMonths([]any{newLdt, months})
			if err, isErr := resM.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resM.(*object.Object)
		}
		if days != 0 {
			resD := localDateTimePlusDaysMethod([]any{newLdt, days})
			if err, isErr := resD.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resD.(*object.Object)
		}
		return resolveLocal(self, newLdt)
	}

	if kName == "java/time/Duration" {
		sec, nano := getDurationSecondsAndNanos(amount)
		newLdt := ldt
		if sec != 0 {
			resS := localDateTimePlusSeconds([]any{newLdt, sec})
			if err, isErr := resS.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resS.(*object.Object)
		}
		if nano != 0 {
			resN := localDateTimePlusNanos([]any{newLdt, nano})
			if err, isErr := resN.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resN.(*object.Object)
		}
		return resolveInstant(self, newLdt)
	}

	res := temporalAmountAddTo([]any{amount, self})
	return res
}

func zonedDateTimePlusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plus: self is null")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plus: invalid amount")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plus: unit is null")
	}
	self := params[0].(*object.Object)
	unitObj := params[2].(*object.Object)
	unitName := getTemporalUnitName(unitObj)

	ldt := getZonedDateTimeLocalDateTime(self)
	switch strings.ToUpper(unitName) {
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		newLdtRes := localDateTimePlusUnit([]any{ldt, amount, unitObj})
		if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return resolveLocal(self, newLdtRes.(*object.Object))
	default:
		newLdtRes := localDateTimePlusUnit([]any{ldt, amount, unitObj})
		if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return resolveInstant(self, newLdtRes.(*object.Object))
	}
}

func zonedDateTimePlusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusYears: self is null")
	}
	years, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusYears: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusYears([]any{ldt, years})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusMonths: self is null")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusMonths: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusMonths([]any{ldt, months})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusWeeks(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusWeeks: self is null")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusWeeks: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusWeeks([]any{ldt, weeks})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusDays(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusDays: self is null")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusDays: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusDaysMethod([]any{ldt, days})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveLocal(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusHours(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusHours: self is null")
	}
	hours, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusHours: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusHours([]any{ldt, hours})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveInstant(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusMinutes: self is null")
	}
	minutes, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusMinutes: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusMinutes([]any{ldt, minutes})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveInstant(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusSeconds: self is null")
	}
	seconds, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusSeconds: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusSeconds([]any{ldt, seconds})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveInstant(self, newLdtRes.(*object.Object))
}

func zonedDateTimePlusNanos(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.plusNanos: self is null")
	}
	nanos, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.plusNanos: invalid argument")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	newLdtRes := localDateTimePlusNanos([]any{ldt, nanos})
	if err, isErr := newLdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return resolveInstant(self, newLdtRes.(*object.Object))
}

// Minus implementations

func zonedDateTimeMinusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(amount.KlassName)

	ldt := getZonedDateTimeLocalDateTime(self)
	if kName == "java/time/Period" {
		years, months, days, _ := getPeriodYMD(amount)
		newLdt := ldt
		if years != 0 {
			resY := localDateTimePlusYears([]any{newLdt, -years})
			if err, isErr := resY.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resY.(*object.Object)
		}
		if months != 0 {
			resM := localDateTimePlusMonths([]any{newLdt, -months})
			if err, isErr := resM.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resM.(*object.Object)
		}
		if days != 0 {
			resD := localDateTimePlusDaysMethod([]any{newLdt, -days})
			if err, isErr := resD.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resD.(*object.Object)
		}
		return resolveLocal(self, newLdt)
	}

	if kName == "java/time/Duration" {
		sec, nano := getDurationSecondsAndNanos(amount)
		newLdt := ldt
		if sec != 0 {
			resS := localDateTimePlusSeconds([]any{newLdt, -sec})
			if err, isErr := resS.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resS.(*object.Object)
		}
		if nano != 0 {
			resN := localDateTimePlusNanos([]any{newLdt, -nano})
			if err, isErr := resN.(*ghelpers.GErrBlk); isErr {
				return err
			}
			newLdt = resN.(*object.Object)
		}
		return resolveInstant(self, newLdt)
	}

	res := temporalAmountSubtractFrom([]any{amount, self})
	return res
}

func zonedDateTimeMinusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minus: self is null")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minus: invalid amount")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minus: unit is null")
	}
	self := params[0].(*object.Object)
	unitObj := params[2].(*object.Object)

	if amount == math.MinInt64 {
		return zonedDateTimePlusUnit([]any{zonedDateTimePlusUnit([]any{self, -(amount + 1), unitObj}), int64(1), unitObj})
	}
	return zonedDateTimePlusUnit([]any{self, -amount, unitObj})
}

func zonedDateTimeMinusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusYears: self is null")
	}
	years, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusYears: invalid argument")
	}
	if years == math.MinInt64 {
		return zonedDateTimePlusYears([]any{zonedDateTimePlusYears([]any{params[0], -(years + 1)}), int64(1)})
	}
	return zonedDateTimePlusYears([]any{params[0], -years})
}

func zonedDateTimeMinusMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusMonths: self is null")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusMonths: invalid argument")
	}
	if months == math.MinInt64 {
		return zonedDateTimePlusMonths([]any{zonedDateTimePlusMonths([]any{params[0], -(months + 1)}), int64(1)})
	}
	return zonedDateTimePlusMonths([]any{params[0], -months})
}

func zonedDateTimeMinusWeeks(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusWeeks: self is null")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusWeeks: invalid argument")
	}
	if weeks == math.MinInt64 {
		return zonedDateTimePlusWeeks([]any{zonedDateTimePlusWeeks([]any{params[0], -(weeks + 1)}), int64(1)})
	}
	return zonedDateTimePlusWeeks([]any{params[0], -weeks})
}

func zonedDateTimeMinusDays(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusDays: self is null")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusDays: invalid argument")
	}
	if days == math.MinInt64 {
		return zonedDateTimePlusDays([]any{zonedDateTimePlusDays([]any{params[0], -(days + 1)}), int64(1)})
	}
	return zonedDateTimePlusDays([]any{params[0], -days})
}

func zonedDateTimeMinusHours(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusHours: self is null")
	}
	hours, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusHours: invalid argument")
	}
	if hours == math.MinInt64 {
		return zonedDateTimePlusHours([]any{zonedDateTimePlusHours([]any{params[0], -(hours + 1)}), int64(1)})
	}
	return zonedDateTimePlusHours([]any{params[0], -hours})
}

func zonedDateTimeMinusMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusMinutes: self is null")
	}
	minutes, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusMinutes: invalid argument")
	}
	if minutes == math.MinInt64 {
		return zonedDateTimePlusMinutes([]any{zonedDateTimePlusMinutes([]any{params[0], -(minutes + 1)}), int64(1)})
	}
	return zonedDateTimePlusMinutes([]any{params[0], -minutes})
}

func zonedDateTimeMinusSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusSeconds: self is null")
	}
	seconds, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusSeconds: invalid argument")
	}
	if seconds == math.MinInt64 {
		return zonedDateTimePlusSeconds([]any{zonedDateTimePlusSeconds([]any{params[0], -(seconds + 1)}), int64(1)})
	}
	return zonedDateTimePlusSeconds([]any{params[0], -seconds})
}

func zonedDateTimeMinusNanos(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.minusNanos: self is null")
	}
	nanos, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZonedDateTime.minusNanos: invalid argument")
	}
	if nanos == math.MinInt64 {
		return zonedDateTimePlusNanos([]any{zonedDateTimePlusNanos([]any{params[0], -(nanos + 1)}), int64(1)})
	}
	return zonedDateTimePlusNanos([]any{params[0], -nanos})
}

// Queries & Until

func zonedDateTimeQuery(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	ensureTemporalQueriesInited()
	if query == zoneIdQueryInstance || query == zoneQueryInstance {
		return getZonedDateTimeZone(self)
	}
	if query == offsetQueryInstance {
		return getZonedDateTimeOffset(self)
	}
	if query == localDateQueryInstance {
		return zonedDateTimeToLocalDate([]any{self})
	}
	if query == localTimeQueryInstance {
		return zonedDateTimeToLocalTime([]any{self})
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
	case "zoneId", "zone":
		return getZonedDateTimeZone(self)
	case "offset":
		return getZonedDateTimeOffset(self)
	case "localDate":
		return zonedDateTimeToLocalDate([]any{self})
	case "localTime":
		return zonedDateTimeToLocalTime([]any{self})
	case "chronology":
		return GetIsoChronology()
	case "precision":
		return ghelpers.CreateTemporalUnit("NANOS")
	default:
		return temporalQueryQueryFrom([]any{query, self})
	}
}

func zonedDateTimeUntil(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.until: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.until: unit is null")
	}
	self := params[0].(*object.Object)
	endExclusive := params[1].(*object.Object)
	unitObj := params[2].(*object.Object)

	endZdtRes := zonedDateTimeFrom([]any{endExclusive})
	if err, isErr := endZdtRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	endZdt := endZdtRes.(*object.Object)

	zone := getZonedDateTimeZone(self)
	endZdtSameInstantRes := zonedDateTimeWithZoneSameInstant([]any{endZdt, zone})
	if err, isErr := endZdtSameInstantRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	endZdtSameInstant := endZdtSameInstantRes.(*object.Object)

	unitName := getTemporalUnitName(unitObj)
	switch strings.ToUpper(unitName) {
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		ldt1 := getZonedDateTimeLocalDateTime(self)
		ldt2 := getZonedDateTimeLocalDateTime(endZdtSameInstant)
		return localDateTimeUntil([]any{ldt1, ldt2, unitObj})
	default:
		// Time based: difference in seconds and nanos
		sec1 := getZonedDateTimeEpochSecond(self)
		nano1 := getZonedDateTimeNano(self)
		sec2 := getZonedDateTimeEpochSecond(endZdt)
		nano2 := getZonedDateTimeNano(endZdt)

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
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported unit in ZonedDateTime.until: "+unitName)
		}
	}
}

// Format, Conversions & Object Methods

func zonedDateTimeFormat(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.format: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.format: formatter is null")
	}
	return dateTimeFormatterFormat([]any{params[1], params[0]})
}

func zonedDateTimeToOffsetDateTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toOffsetDateTime: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	offset := getZonedDateTimeOffset(self)
	odt := object.MakeEmptyObjectWithClassName(&classNameOffsetDateTimeVal)
	odt.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	odt.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offset}
	return odt
}

func zonedDateTimeToInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toInstant: self is null")
	}
	self := params[0].(*object.Object)
	sec := getZonedDateTimeEpochSecond(self)
	nano := getZonedDateTimeNano(self)
	return createInstant(sec, nano)
}

func zonedDateTimeToEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toEpochSecond: self is null")
	}
	self := params[0].(*object.Object)
	return getZonedDateTimeEpochSecond(self)
}

func zonedDateTimeCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.compareTo: null argument")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	s1 := getZonedDateTimeEpochSecond(self)
	s2 := getChronoZonedDateTimeEpochSecond(other)
	if s1 < s2 {
		return int64(-1)
	}
	if s1 > s2 {
		return int64(1)
	}

	n1 := getZonedDateTimeNano(self)
	_, _, _, _, _, _, n2 := getChronoZonedDateTimeYMDHMSN(other)
	if n1 < n2 {
		return int64(-1)
	}
	if n1 > n2 {
		return int64(1)
	}

	ldt1 := getZonedDateTimeLocalDateTime(self)
	ldt2 := getChronoZonedDateTimeLocalDateTime(other)
	if ldt1 != nil && ldt2 != nil {
		cmpLdt := compareLocalDateTime(ldt1, ldt2)
		if cmpLdt != 0 {
			return cmpLdt
		}
	}

	z1 := getZonedDateTimeZone(self)
	z2 := getChronoZonedDateTimeZone(other)
	z1Str := getZoneIdString(z1)
	z2Str := getZoneIdString(z2)
	if z1Str < z2Str {
		return int64(-1)
	}
	if z1Str > z2Str {
		return int64(1)
	}

	return int64(0)
}

func zonedDateTimeIsAfter(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.isAfter: null argument")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	s1 := getZonedDateTimeEpochSecond(self)
	s2 := getChronoZonedDateTimeEpochSecond(other)
	if s1 > s2 {
		return types.JavaBoolTrue
	}
	if s1 == s2 {
		n1 := getZonedDateTimeNano(self)
		_, _, _, _, _, _, n2 := getChronoZonedDateTimeYMDHMSN(other)
		if n1 > n2 {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func zonedDateTimeIsBefore(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.isBefore: null argument")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	s1 := getZonedDateTimeEpochSecond(self)
	s2 := getChronoZonedDateTimeEpochSecond(other)
	if s1 < s2 {
		return types.JavaBoolTrue
	}
	if s1 == s2 {
		n1 := getZonedDateTimeNano(self)
		_, _, _, _, _, _, n2 := getChronoZonedDateTimeYMDHMSN(other)
		if n1 < n2 {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func zonedDateTimeIsEqual(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.isEqual: null argument")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	s1 := getZonedDateTimeEpochSecond(self)
	s2 := getChronoZonedDateTimeEpochSecond(other)
	if s1 == s2 {
		n1 := getZonedDateTimeNano(self)
		_, _, _, _, _, _, n2 := getChronoZonedDateTimeYMDHMSN(other)
		if n1 == n2 {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func zonedDateTimeEquals(params []any) any {
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
	if kName != classNameZonedDateTimeVal {
		return types.JavaBoolFalse
	}

	ldt1 := getZonedDateTimeLocalDateTime(self)
	ldt2 := getZonedDateTimeLocalDateTime(other)
	eqLdt := localDateTimeEquals([]any{ldt1, ldt2})
	if eqLdt != types.JavaBoolTrue {
		return types.JavaBoolFalse
	}

	off1 := getZonedDateTimeOffset(self)
	off2 := getZonedDateTimeOffset(other)
	if getZoneOffsetTotalSeconds(off1) != getZoneOffsetTotalSeconds(off2) {
		return types.JavaBoolFalse
	}

	z1 := getZonedDateTimeZone(self)
	z2 := getZonedDateTimeZone(other)
	eqZ := zoneIdEquals([]any{z1, z2})
	if eqZ != types.JavaBoolTrue {
		return types.JavaBoolFalse
	}

	return types.JavaBoolTrue
}

func zonedDateTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	off := getZonedDateTimeOffset(self)
	zone := getZonedDateTimeZone(self)

	hLdtVal, _ := getInt64Val(localDateTimeHashCode([]any{ldt}))
	hOffVal, _ := getInt64Val(zoneOffsetHashCode([]any{off}))
	hZoneVal, _ := getInt64Val(zoneIdHashCode([]any{zone}))

	// rotate left 16 bits
	rotZone := (hZoneVal << 16) | (hZoneVal >> 16)
	return int64(int32(hLdtVal ^ hOffVal ^ rotZone))
}

func zonedDateTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZonedDateTime.toString: self is null")
	}
	self := params[0].(*object.Object)
	ldt := getZonedDateTimeLocalDateTime(self)
	off := getZonedDateTimeOffset(self)
	zone := getZonedDateTimeZone(self)

	ldtStrVal := localDateTimeToString([]any{ldt})
	var ldtStr string
	if strObj, ok := ldtStrVal.(*object.Object); ok && object.IsStringObject(strObj) {
		ldtStr = object.GoStringFromStringObject(strObj)
	}

	offStrVal := zoneOffsetToString([]any{off})
	var offStr string
	if strObj, ok := offStrVal.(*object.Object); ok && object.IsStringObject(strObj) {
		offStr = object.GoStringFromStringObject(strObj)
	}

	str := ldtStr + offStr
	if zone != nil && !object.IsNull(zone) && !isZoneOffsetObject(zone) && zone != off {
		zoneStrVal := zoneIdToString([]any{zone})
		if strObj, ok := zoneStrVal.(*object.Object); ok && object.IsStringObject(strObj) {
			zoneStr := object.GoStringFromStringObject(strObj)
			str += "[" + zoneStr + "]"
		}
	}
	return object.StringObjectFromGoString(str)
}

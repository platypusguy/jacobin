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
	"strconv"
	"strings"
	"sync"
	"time"
)

var classNameLocalDateVal = "java/time/LocalDate"

const (
	localDateMinYear = -999999999
	localDateMaxYear = 999999999
)

var localDateMutex sync.Mutex
var localDateInitialized bool = false

var (
	localDateMin   *object.Object
	localDateMax   *object.Object
	localDateEpoch *object.Object
)

func ensureLocalDateInited() {
	localDateMutex.Lock()
	defer localDateMutex.Unlock()
	if localDateInitialized {
		return
	}

	localDateMin = createLocalDateInternal(localDateMinYear, 1, 1)
	localDateMax = createLocalDateInternal(localDateMaxYear, 12, 31)
	localDateEpoch = createLocalDateInternal(1970, 1, 1)

	_ = statics.AddStatic(classNameLocalDateVal+".MIN", statics.Static{
		Type:  "Ljava/time/LocalDate;",
		Value: localDateMin,
	})
	_ = statics.AddStatic(classNameLocalDateVal+".MAX", statics.Static{
		Type:  "Ljava/time/LocalDate;",
		Value: localDateMax,
	})
	_ = statics.AddStatic(classNameLocalDateVal+".EPOCH", statics.Static{
		Type:  "Ljava/time/LocalDate;",
		Value: localDateEpoch,
	})

	localDateInitialized = true
}

func createLocalDateInternal(year, month, day int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameLocalDateVal)
	obj.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: year}
	obj.FieldTable["month"] = object.Field{Ftype: types.Short, Fvalue: month}
	obj.FieldTable["day"] = object.Field{Ftype: types.Short, Fvalue: day}
	obj.FieldTable["YEAR"] = object.Field{Ftype: types.Long, Fvalue: year}
	obj.FieldTable["MONTH_OF_YEAR"] = object.Field{Ftype: types.Long, Fvalue: month}
	obj.FieldTable["DAY_OF_MONTH"] = object.Field{Ftype: types.Long, Fvalue: day}
	epochDay := ymdToEpochDay(year, month, day)
	obj.FieldTable["epochDay"] = object.Field{Ftype: types.Long, Fvalue: epochDay}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: GetIsoChronology()}
	return obj
}

// CreateLocalDate public helper function
func CreateLocalDate(year, month, day int64) *object.Object {
	ensureLocalDateInited()
	return createLocalDateInternal(year, month, day)
}

func getLocalDateYMD(targetObject *object.Object) (int64, int64, int64, bool) {
	if targetObject == nil || targetObject.FieldTable == nil {
		return 0, 0, 0, false
	}
	var year, month, day int64
	var okY, okM, okD bool

	if yFld, exists := targetObject.FieldTable["year"]; exists {
		year, okY = getInt64Val(yFld.Fvalue)
	} else if yFld, exists := targetObject.FieldTable["YEAR"]; exists {
		year, okY = getInt64Val(yFld.Fvalue)
	}

	if mFld, exists := targetObject.FieldTable["month"]; exists {
		month, okM = getInt64Val(mFld.Fvalue)
	} else if mFld, exists := targetObject.FieldTable["MONTH_OF_YEAR"]; exists {
		month, okM = getInt64Val(mFld.Fvalue)
	}

	if dFld, exists := targetObject.FieldTable["day"]; exists {
		day, okD = getInt64Val(dFld.Fvalue)
	} else if dFld, exists := targetObject.FieldTable["DAY_OF_MONTH"]; exists {
		day, okD = getInt64Val(dFld.Fvalue)
	}

	if okY && okM && okD {
		return year, month, day, true
	}

	if epochFld, exists := targetObject.FieldTable["epochDay"]; exists {
		if ep, ok := getInt64Val(epochFld.Fvalue); ok {
			y, m, d := epochDayToYMD(ep)
			return y, m, d, true
		}
	}

	return 0, 0, 0, false
}

func isLeapYearDate(year int64) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func daysInMonthDate(year, month int64) int64 {
	switch month {
	case 2:
		if isLeapYearDate(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

func ymdToEpochDay(year, month, day int64) int64 {
	var y = year
	var m = month
	var totalDays int64 = 0
	totalDays += 365 * y
	if y >= 0 {
		totalDays += (y+3)/4 - (y+99)/100 + (y+399)/400
	} else {
		totalDays -= y/-4 - y/-100 + y/-400
	}
	totalDays += (367*m - 362) / 12
	totalDays += day - 1
	if m > 2 {
		totalDays--
		if !isLeapYearDate(year) {
			totalDays--
		}
	}
	return totalDays - 719528
}

func epochDayToYMD(epochDay int64) (int64, int64, int64) {
	var zeroDay = epochDay + 719528 - 60
	var adjust int64 = 0
	if zeroDay < 0 {
		var adjustCycles = (zeroDay+1)/146097 - 1
		adjust = adjustCycles * 400
		zeroDay -= adjustCycles * 146097
	}
	var yearEst = (400*zeroDay + 591) / 146097
	var doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	if doyEst < 0 {
		yearEst--
		doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	}
	yearEst += adjust

	var marchDay0 = doyEst
	var marchMonth0 = (marchDay0*5 + 2) / 153
	var month = (marchMonth0+2)%12 + 1
	var dom = marchDay0 - (marchMonth0*306+5)/10 + 1
	yearEst += (marchMonth0 + 2) / 12

	return yearEst, month, dom
}

func dayOfYearDate(year, month, day int64) int64 {
	var doy int64 = 0
	for m := int64(1); m < month; m++ {
		doy += daysInMonthDate(year, m)
	}
	return doy + day
}

func dayOfWeekDate(year, month, day int64) int64 {
	epochDay := ymdToEpochDay(year, month, day)
	dow := (epochDay + 3) % 7
	if dow < 0 {
		dow += 7
	}
	return dow + 1
}

func validateLocalDate(year, month, day int64) error {
	if year < localDateMinYear || year > localDateMaxYear {
		return fmt.Errorf("Invalid value for Year (valid values %d - %d): %d", localDateMinYear, localDateMaxYear, year)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("Invalid value for MonthOfYear (valid values 1 - 12): %d", month)
	}
	maxDays := daysInMonthDate(year, month)
	if day < 1 || day > maxDays {
		return fmt.Errorf("Invalid value for DayOfMonth (valid values 1 - %d): %d", maxDays, day)
	}
	return nil
}

func formatLocalDateString(year, month, day int64) string {
	yStr := formatYearString(year)
	return fmt.Sprintf("%s-%02d-%02d", yStr, month, day)
}

// Load_Time_LocalDate registers all java/time/LocalDate methods.
func Load_Time_LocalDate() {
	ghelpers.MethodSignatures["java/time/LocalDate.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateClinit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atStartOfDay()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateAtStartOfDay,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atStartOfDay(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateAtStartOfDayZone,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atTime(II)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateAtTimeHM,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atTime(III)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localDateAtTimeHMS,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atTime(IIII)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  localDateAtTimeHMSN,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atTime(Ljava/time/LocalTime;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateAtTimeLocalTime,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.atTime(Ljava/time/OffsetTime;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateAtTimeOffsetTime,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateCompareTo,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.compareTo(Ljava/time/chrono/ChronoLocalDate;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateCompareTo,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.datesUntil(Ljava/time/LocalDate;)Ljava/util/stream/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateDatesUntil,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.datesUntil(Ljava/time/LocalDate;Ljava/time/Period;)Ljava/util/stream/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateDatesUntilPeriod,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateEquals,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateFormat,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateFrom,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateGet,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getChronology()Ljava/time/chrono/IsoChronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetChronology,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetChronology,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getDayOfMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getDayOfWeek()Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetDayOfWeek,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getDayOfYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getEra()Ljava/time/chrono/Era;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetEra,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getEra()Ljava/time/chrono/IsoEra;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetEra,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateGetLong,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getMonth()Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getMonthValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetMonthValue,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.getYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateGetYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateHashCode,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.isAfter(Ljava/time/chrono/ChronoLocalDate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateIsAfter,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.isBefore(Ljava/time/chrono/ChronoLocalDate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateIsBefore,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.isEqual(Ljava/time/chrono/ChronoLocalDate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateIsEqual,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.isLeapYear()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateIsLeapYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.lengthOfMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateLengthOfMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.lengthOfYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateLengthOfYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minusDays(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusDays,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minusMonths(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusMonths,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minusWeeks(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusWeeks,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.minusYears(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateMinusYears,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.now()Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateNow,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.now(Ljava/time/Clock;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateNowClock,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.now(Ljava/time/ZoneId;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateNowZone,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.of(III)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localDateOfYMD,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.of(ILjava/time/Month;I)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localDateOfYMonthD,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.ofEpochDay(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateOfEpochDay,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateOfInstant,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.ofYearDay(II)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateOfYearDay,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.parse(Ljava/lang/CharSequence;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateParse,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateParseFormatter,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDatePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDatePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDatePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plusDays(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusDays,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plusMonths(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusMonths,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plusWeeks(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusWeeks,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.plusYears(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDatePlusYears,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateQuery,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateRange,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.toEpochDay()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateToEpochDay,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.toEpochSecond(Ljava/time/LocalTime;Ljava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localDateToString,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateUntilUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.until(Ljava/time/chrono/ChronoLocalDate;)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateUntilPeriod,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.with(Ljava/time/temporal/TemporalField;J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localDateWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.withDayOfMonth(I)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.withDayOfYear(I)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.withMonth(I)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithMonth,
		}

	ghelpers.MethodSignatures["java/time/LocalDate.withYear(I)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localDateWithYear,
		}
}

func localDateClinit([]any) any {
	ensureLocalDateInited()
	return nil
}

func localDateNow([]any) any {
	now := time.Now()
	return CreateLocalDate(int64(now.Year()), int64(now.Month()), int64(now.Day()))
}

func localDateNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.now: clock is null")
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
	return zonedDateTimeToLocalDate([]any{zdt.(*object.Object)})
}

func localDateNowZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.now: zone is null")
	}
	now := time.Now()
	return CreateLocalDate(int64(now.Year()), int64(now.Month()), int64(now.Day()))
}

func localDateOfYMD(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.of: missing arguments")
	}
	y, okY := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	d, okD := getInt64Val(params[2])
	if !okY || !okM || !okD {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.of: invalid arguments")
	}
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return CreateLocalDate(y, m, d)
}

func localDateOfYMonthD(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.of: missing arguments")
	}
	y, okY := getInt64Val(params[0])
	if !okY {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.of: invalid year")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.of: month is null")
	}
	mObj := params[1].(*object.Object)
	mVal, okM := getMonthValueFromObject(mObj)
	if !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.of: invalid month object")
	}
	d, okD := getInt64Val(params[2])
	if !okD {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.of: invalid day")
	}
	if err := validateLocalDate(y, mVal, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return CreateLocalDate(y, mVal, d)
}

func getMonthValueFromObject(obj *object.Object) (int64, bool) {
	if obj == nil || obj.FieldTable == nil {
		return 0, false
	}
	if valFld, exists := obj.FieldTable["value"]; exists {
		if val, ok := getInt64Val(valFld.Fvalue); ok {
			return val, true
		}
	}
	if valFld, exists := obj.FieldTable["ordinal"]; exists {
		if val, ok := getInt64Val(valFld.Fvalue); ok {
			return val + 1, true
		}
	}
	if valFld, exists := obj.FieldTable["MONTH_OF_YEAR"]; exists {
		if val, ok := getInt64Val(valFld.Fvalue); ok {
			return val, true
		}
	}
	return 0, false
}

func getDayOfWeekValueFromObject(obj *object.Object) (int64, bool) {
	if obj == nil || obj.FieldTable == nil {
		return 0, false
	}
	if valFld, exists := obj.FieldTable["value"]; exists {
		if val, ok := getInt64Val(valFld.Fvalue); ok {
			return val, true
		}
	}
	if valFld, exists := obj.FieldTable["ordinal"]; exists {
		if val, ok := getInt64Val(valFld.Fvalue); ok {
			return val + 1, true
		}
	}
	if valFld, exists := obj.FieldTable["DAY_OF_WEEK"]; exists {
		if val, ok := getInt64Val(valFld.Fvalue); ok {
			return val, true
		}
	}
	return 0, false
}

func localDateOfEpochDay(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.ofEpochDay: missing argument")
	}
	epochDay, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.ofEpochDay: invalid argument")
	}
	y, m, d := epochDayToYMD(epochDay)
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return CreateLocalDate(y, m, d)
}

func localDateOfInstant(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.ofInstant: null argument")
	}
	instantObj := params[0].(*object.Object)
	var epochSec int64 = 0
	if fld, ok := instantObj.FieldTable["seconds"]; ok {
		if s, ok2 := getInt64Val(fld.Fvalue); ok2 {
			epochSec = s
		}
	}
	epochDay := epochSec / 86400
	if epochSec < 0 && epochSec%86400 != 0 {
		epochDay--
	}
	y, m, d := epochDayToYMD(epochDay)
	if err := validateLocalDate(y, m, d); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return CreateLocalDate(y, m, d)
}

func localDateOfYearDay(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.ofYearDay: missing arguments")
	}
	year, okY := getInt64Val(params[0])
	dayOfYear, okD := getInt64Val(params[1])
	if !okY || !okD {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.ofYearDay: invalid arguments")
	}
	if year < localDateMinYear || year > localDateMaxYear {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", localDateMinYear, localDateMaxYear, year))
	}
	maxDays := int64(365)
	if isLeapYearDate(year) {
		maxDays = 366
	}
	if dayOfYear < 1 || dayOfYear > maxDays {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for DayOfYear (valid values 1 - %d): %d", maxDays, dayOfYear))
	}

	curDoy := dayOfYear
	month := int64(1)
	for m := int64(1); m <= 12; m++ {
		dim := daysInMonthDate(year, m)
		if curDoy <= dim {
			month = m
			break
		}
		curDoy -= dim
	}
	return CreateLocalDate(year, month, curDoy)
}

func localDateFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.from: null argument")
	}
	obj := params[0].(*object.Object)
	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	if className == classNameLocalDateVal {
		return obj
	}
	if strings.Contains(className, "LocalDateTime") {
		d := getLocalDateTimeDate(obj)
		if d != nil {
			return d
		}
	}
	if strings.Contains(className, "ZonedDateTime") {
		ldt := getZonedDateTimeLocalDateTime(obj)
		if ldt != nil {
			d := getLocalDateTimeDate(ldt)
			if d != nil {
				return d
			}
		}
	}
	if strings.Contains(className, "OffsetDateTime") {
		ldt := getOffsetDateTimeDateTime(obj)
		if ldt != nil {
			d := getLocalDateTimeDate(ldt)
			if d != nil {
				return d
			}
		}
	}
	if y, m, d, ok := getLocalDateYMD(obj); ok {
		return CreateLocalDate(y, m, d)
	}
	if epFld, ok := obj.FieldTable["EPOCH_DAY"]; ok {
		if ep, ok2 := getInt64Val(epFld.Fvalue); ok2 {
			y, m, d := epochDayToYMD(ep)
			return CreateLocalDate(y, m, d)
		}
	}
	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain LocalDate from TemporalAccessor")
}

func parseLocalDateString(text string) (int64, int64, int64, error) {
	parts := strings.Split(text, "-")
	if len(parts) == 3 {
		y, errY := strconv.ParseInt(parts[0], 10, 64)
		m, errM := strconv.ParseInt(parts[1], 10, 64)
		d, errD := strconv.ParseInt(parts[2], 10, 64)
		if errY == nil && errM == nil && errD == nil {
			if err := validateLocalDate(y, m, d); err == nil {
				return y, m, d, nil
			}
		}
	} else if len(parts) == 4 && parts[0] == "" {
		// negative year: -YYYY-MM-DD
		y, errY := strconv.ParseInt(parts[1], 10, 64)
		m, errM := strconv.ParseInt(parts[2], 10, 64)
		d, errD := strconv.ParseInt(parts[3], 10, 64)
		if errY == nil && errM == nil && errD == nil {
			if err := validateLocalDate(-y, m, d); err == nil {
				return -y, m, d, nil
			}
		}
	}
	return 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
}

func localDateParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.parse: text is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.parse: invalid string")
	}
	goStr := strings.TrimSpace(object.GoStringFromStringObject(strObj))
	y, m, d, err := parseLocalDateString(goStr)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, err.Error())
	}
	return CreateLocalDate(y, m, d)
}

func localDateParseFormatter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.parse: text is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.parse: formatter is null")
	}
	textObj := params[0].(*object.Object)
	formatterObj := params[1].(*object.Object)
	if !object.IsStringObject(textObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.parse: invalid string")
	}
	text := object.GoStringFromStringObject(textObj)
	accessor, err := parseTextToTemporalAccessor(formatterObj, text)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, err.Error())
	}
	return localDateFrom([]any{accessor})
}

func localDateGetYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getYear: self is null")
	}
	y, _, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getYear: invalid state")
	}
	return int64(y)
}

func localDateGetMonthValue(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getMonthValue: self is null")
	}
	_, m, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getMonthValue: invalid state")
	}
	return int64(m)
}

func localDateGetMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getMonth: self is null")
	}
	_, m, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getMonth: invalid state")
	}
	return monthOf([]any{m})
}

func localDateGetDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getDayOfMonth: self is null")
	}
	_, _, d, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getDayOfMonth: invalid state")
	}
	return int64(d)
}

func localDateGetDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getDayOfYear: self is null")
	}
	y, m, d, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getDayOfYear: invalid state")
	}
	return int64(dayOfYearDate(y, m, d))
}

func localDateGetDayOfWeek(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getDayOfWeek: self is null")
	}
	y, m, d, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getDayOfWeek: invalid state")
	}
	dow := dayOfWeekDate(y, m, d)
	return dayOfWeekOf([]any{dow})
}

func localDateGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getChronology: self is null")
	}
	return GetIsoChronology()
}

func localDateGetEra(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getEra: self is null")
	}
	y, _, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getEra: invalid state")
	}
	if y >= 1 {
		return createEra(1, "CE")
	}
	return createEra(0, "BCE")
}

func localDateIsLeapYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isLeapYear: self is null")
	}
	y, _, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.isLeapYear: invalid state")
	}
	return types.ConvertGoBoolToJavaBool(isLeapYearDate(y))
}

func localDateLengthOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.lengthOfMonth: self is null")
	}
	y, m, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.lengthOfMonth: invalid state")
	}
	return int64(daysInMonthDate(y, m))
}

func localDateLengthOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.lengthOfYear: self is null")
	}
	y, _, _, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.lengthOfYear: invalid state")
	}
	if isLeapYearDate(y) {
		return int64(366)
	}
	return int64(365)
}

func localDateIsSupportedField(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	if fieldName == "" {
		return types.JavaBoolFalse
	}
	switch strings.ToUpper(fieldName) {
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "ERA":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

func localDateIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	unitObj := params[1].(*object.Object)
	unitName := getTemporalUnitName(unitObj)
	if unitName == "" {
		return types.JavaBoolFalse
	}
	switch strings.ToUpper(unitName) {
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

func localDateRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.range: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	y, m, _, ok := getLocalDateYMD(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.range: invalid date")
	}

	switch strings.ToUpper(fieldName) {
	case "DAY_OF_MONTH":
		return createValueRange(1, 1, daysInMonthDate(y, m), daysInMonthDate(y, m))
	case "DAY_OF_YEAR":
		maxDoy := int64(365)
		if isLeapYearDate(y) {
			maxDoy = 366
		}
		return createValueRange(1, 1, maxDoy, maxDoy)
	case "ALIGNED_WEEK_OF_MONTH":
		if m == 2 && !isLeapYearDate(y) {
			return createValueRange(1, 1, 4, 4)
		}
		return createValueRange(1, 1, 4, 5)
	case "YEAR_OF_ERA":
		if y <= 0 {
			return createValueRange(1, 1, 1000000000, 1000000000)
		}
		return createValueRange(1, 1, localDateMaxYear, localDateMaxYear)
	case "ERA":
		return createValueRange(0, 0, 1, 1)
	case "YEAR":
		return createValueRange(localDateMinYear, localDateMinYear, localDateMaxYear, localDateMaxYear)
	case "MONTH_OF_YEAR":
		return createValueRange(1, 1, 12, 12)
	case "DAY_OF_WEEK":
		return createValueRange(1, 1, 7, 7)
	case "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR":
		return createValueRange(1, 1, 7, 7)
	case "ALIGNED_WEEK_OF_YEAR":
		return createValueRange(1, 1, 53, 53)
	case "PROLEPTIC_MONTH":
		return createValueRange(localDateMinYear*12, localDateMinYear*12, localDateMaxYear*12+11, localDateMaxYear*12+11)
	case "EPOCH_DAY":
		return createValueRange(-365243219162, -365243219162, 365241780471, 365241780471)
	default:
		if object.GoStringFromStringPoolIndex(fieldObj.KlassName) == classNameChronoField {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
		}
		return temporalAccessorRange(params)
	}
}

func localDateGetLong(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.getLong: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	y, m, d, ok := getLocalDateYMD(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.getLong: invalid date")
	}

	switch strings.ToUpper(fieldName) {
	case "DAY_OF_WEEK":
		return dayOfWeekDate(y, m, d)
	case "ALIGNED_DAY_OF_WEEK_IN_MONTH":
		return ((d - 1) % 7) + 1
	case "ALIGNED_DAY_OF_WEEK_IN_YEAR":
		doy := dayOfYearDate(y, m, d)
		return ((doy - 1) % 7) + 1
	case "DAY_OF_MONTH":
		return d
	case "DAY_OF_YEAR":
		return dayOfYearDate(y, m, d)
	case "EPOCH_DAY":
		return ymdToEpochDay(y, m, d)
	case "ALIGNED_WEEK_OF_MONTH":
		return ((d - 1) / 7) + 1
	case "ALIGNED_WEEK_OF_YEAR":
		doy := dayOfYearDate(y, m, d)
		return ((doy - 1) / 7) + 1
	case "MONTH_OF_YEAR":
		return m
	case "PROLEPTIC_MONTH":
		return y*12 + m - 1
	case "YEAR_OF_ERA":
		if y >= 1 {
			return y
		}
		return 1 - y
	case "YEAR":
		return y
	case "ERA":
		if y >= 1 {
			return int64(1)
		}
		return int64(0)
	default:
		if object.GoStringFromStringPoolIndex(fieldObj.KlassName) == classNameChronoField {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
		}
		return temporalAccessorGetLong(params)
	}
}

func localDateGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.get: field is null")
	}
	res := localDateGetLong(params)
	if _, isErr := res.(*ghelpers.GErrBlk); isErr {
		return res
	}
	if val, ok := getInt64Val(res); ok {
		return int64(int32(val))
	}
	return res
}

func localDateToEpochDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.toEpochDay: self is null")
	}
	y, m, d, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.toEpochDay: invalid date")
	}
	return ymdToEpochDay(y, m, d)
}

func localDateToEpochSecond(params []any) any {
	if len(params) < 3 || params[0] == nil || params[1] == nil || params[2] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.toEpochSecond: null argument")
	}
	self := params[0].(*object.Object)
	timeObj := params[1].(*object.Object)
	offsetObj := params[2].(*object.Object)

	y, m, d, ok := getLocalDateYMD(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.toEpochSecond: invalid date")
	}
	epochDay := ymdToEpochDay(y, m, d)
	var secsOfDay int64 = 0
	if hFld, ok := timeObj.FieldTable["hour"]; ok {
		if h, ok2 := getInt64Val(hFld.Fvalue); ok2 {
			secsOfDay += h * 3600
		}
	}
	if minFld, ok := timeObj.FieldTable["minute"]; ok {
		if min, ok2 := getInt64Val(minFld.Fvalue); ok2 {
			secsOfDay += min * 60
		}
	}
	if sFld, ok := timeObj.FieldTable["second"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			secsOfDay += s
		}
	}
	var offsetSec int64 = 0
	if offFld, ok := offsetObj.FieldTable["totalSeconds"]; ok {
		if off, ok2 := getInt64Val(offFld.Fvalue); ok2 {
			offsetSec = off
		}
	}
	return epochDay*86400 + secsOfDay - offsetSec
}

func localDateToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.toString: self is null")
	}
	y, m, d, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.toString: invalid date")
	}
	return object.StringObjectFromGoString(formatLocalDateString(y, m, d))
}

func localDateFormat(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.format: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.format: formatter is null")
	}
	return dateTimeFormatterFormat([]any{params[1], params[0]})
}

func localDateCompareTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	y1, m1, d1, ok1 := getLocalDateYMD(self)
	y2, m2, d2, ok2 := getLocalDateYMD(other)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "LocalDate.compareTo: invalid operand")
	}
	if y1 != y2 {
		if y1 < y2 {
			return int64(-1)
		}
		return int64(1)
	}
	if m1 != m2 {
		if m1 < m2 {
			return int64(-1)
		}
		return int64(1)
	}
	if d1 != d2 {
		if d1 < d2 {
			return int64(-1)
		}
		return int64(1)
	}
	return int64(0)
}

func localDateEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil {
		return types.JavaBoolFalse
	}
	self, okS := params[0].(*object.Object)
	other, okO := params[1].(*object.Object)
	if !okS || !okO || object.IsNull(self) || object.IsNull(other) {
		return types.JavaBoolFalse
	}
	if self == other {
		return types.JavaBoolTrue
	}
	y1, m1, d1, ok1 := getLocalDateYMD(self)
	y2, m2, d2, ok2 := getLocalDateYMD(other)
	if !ok1 || !ok2 {
		return types.JavaBoolFalse
	}
	return types.ConvertGoBoolToJavaBool(y1 == y2 && m1 == m2 && d1 == d2)
}

func localDateHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.hashCode: self is null")
	}
	y, m, d, ok := getLocalDateYMD(params[0].(*object.Object))
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.hashCode: invalid date")
	}
	yearVal := int32(y)
	monthVal := int32(m)
	dayVal := int32(d)
	return int64((yearVal & -2048) ^ ((yearVal << 11) + (monthVal << 6) + dayVal))
}

func localDateIsAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isAfter: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isAfter: other is null")
	}
	cmp := localDateCompareTo(params)
	if err, isErr := cmp.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if val, ok := getInt64Val(cmp); ok {
		return types.ConvertGoBoolToJavaBool(val > 0)
	}
	return types.JavaBoolFalse
}

func localDateIsBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isBefore: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isBefore: other is null")
	}
	cmp := localDateCompareTo(params)
	if err, isErr := cmp.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if val, ok := getInt64Val(cmp); ok {
		return types.ConvertGoBoolToJavaBool(val < 0)
	}
	return types.JavaBoolFalse
}

func localDateIsEqual(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isEqual: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.isEqual: other is null")
	}
	cmp := localDateCompareTo(params)
	if err, isErr := cmp.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if val, ok := getInt64Val(cmp); ok {
		return types.ConvertGoBoolToJavaBool(val == 0)
	}
	return types.JavaBoolFalse
}

func localDateWithYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.withYear: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withYear: missing argument")
	}
	newYear, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withYear: invalid argument")
	}
	y, m, d, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.withYear: invalid state")
	}
	if y == newYear {
		return params[0]
	}
	if newYear < localDateMinYear || newYear > localDateMaxYear {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", localDateMinYear, localDateMaxYear, newYear))
	}
	maxDays := daysInMonthDate(newYear, m)
	if d > maxDays {
		d = maxDays
	}
	return CreateLocalDate(newYear, m, d)
}

func localDateWithMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.withMonth: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withMonth: missing argument")
	}
	newMonth, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withMonth: invalid argument")
	}
	y, m, d, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.withMonth: invalid state")
	}
	if m == newMonth {
		return params[0]
	}
	if newMonth < 1 || newMonth > 12 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MonthOfYear (valid values 1 - 12): %d", newMonth))
	}
	maxDays := daysInMonthDate(y, newMonth)
	if d > maxDays {
		d = maxDays
	}
	return CreateLocalDate(y, newMonth, d)
}

func localDateWithDayOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.withDayOfMonth: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withDayOfMonth: missing argument")
	}
	newDay, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withDayOfMonth: invalid argument")
	}
	y, m, d, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.withDayOfMonth: invalid state")
	}
	if d == newDay {
		return params[0]
	}
	maxDays := daysInMonthDate(y, m)
	if newDay < 1 || newDay > maxDays {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for DayOfMonth (valid values 1 - %d): %d", maxDays, newDay))
	}
	return CreateLocalDate(y, m, newDay)
}

func localDateWithDayOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.withDayOfYear: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withDayOfYear: missing argument")
	}
	doy, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.withDayOfYear: invalid argument")
	}
	y, _, _, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.withDayOfYear: invalid state")
	}
	return localDateOfYearDay([]any{y, doy})
}

func localDateWithField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.with: field is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.with: missing newValue")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	newValue, okVal := getInt64Val(params[2])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.with: invalid new value")
	}

	fieldName := getTemporalFieldName(fieldObj)
	y, m, d, ok := getLocalDateYMD(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.with: invalid date")
	}

	switch strings.ToUpper(fieldName) {
	case "DAY_OF_WEEK":
		if newValue < 1 || newValue > 7 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for DayOfWeek (valid values 1 - 7): %d", newValue))
		}
		curDow := dayOfWeekDate(y, m, d)
		return localDatePlusDays([]any{self, newValue - curDow})
	case "ALIGNED_DAY_OF_WEEK_IN_MONTH":
		if newValue < 1 || newValue > 7 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for AlignedDayOfWeekInMonth (valid values 1 - 7): %d", newValue))
		}
		cur := ((d - 1) % 7) + 1
		return localDatePlusDays([]any{self, newValue - cur})
	case "ALIGNED_DAY_OF_WEEK_IN_YEAR":
		if newValue < 1 || newValue > 7 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for AlignedDayOfWeekInYear (valid values 1 - 7): %d", newValue))
		}
		curDoy := dayOfYearDate(y, m, d)
		cur := ((curDoy - 1) % 7) + 1
		return localDatePlusDays([]any{self, newValue - cur})
	case "DAY_OF_MONTH":
		return localDateWithDayOfMonth([]any{self, newValue})
	case "DAY_OF_YEAR":
		return localDateWithDayOfYear([]any{self, newValue})
	case "EPOCH_DAY":
		return localDateOfEpochDay([]any{newValue})
	case "ALIGNED_WEEK_OF_MONTH":
		if newValue < 1 || newValue > 5 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for AlignedWeekOfMonth (valid values 1 - 5): %d", newValue))
		}
		curWeek := ((d - 1) / 7) + 1
		return localDatePlusDays([]any{self, (newValue - curWeek) * 7})
	case "ALIGNED_WEEK_OF_YEAR":
		if newValue < 1 || newValue > 53 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for AlignedWeekOfYear (valid values 1 - 53): %d", newValue))
		}
		curDoy := dayOfYearDate(y, m, d)
		curWeek := ((curDoy - 1) / 7) + 1
		return localDatePlusDays([]any{self, (newValue - curWeek) * 7})
	case "MONTH_OF_YEAR":
		return localDateWithMonth([]any{self, newValue})
	case "PROLEPTIC_MONTH":
		curPM := y*12 + m - 1
		return localDatePlusMonths([]any{self, newValue - curPM})
	case "YEAR_OF_ERA":
		if y >= 1 {
			return localDateWithYear([]any{self, newValue})
		}
		return localDateWithYear([]any{self, 1 - newValue})
	case "YEAR":
		return localDateWithYear([]any{self, newValue})
	case "ERA":
		if newValue < 0 || newValue > 1 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Era (valid values 0 - 1): %d", newValue))
		}
		curEra := int64(0)
		if y >= 1 {
			curEra = 1
		}
		if curEra == newValue {
			return self
		}
		return localDateWithYear([]any{self, 1 - y})
	default:
		if object.GoStringFromStringPoolIndex(fieldObj.KlassName) == classNameChronoField {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
		}
		res := temporalWithField([]any{self, fieldObj, newValue})
		if err, isErr := res.(*ghelpers.GErrBlk); isErr {
			return err
		}
		if resObj, ok := res.(*object.Object); ok {
			return localDateFrom([]any{resObj})
		}
		return res
	}
}

func localDateWithAdjuster(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	if y, m, d, ok := getLocalDateYMD(adjuster); ok {
		return CreateLocalDate(y, m, d)
	}
	if monthVal, ok := getMonthValueFromObject(adjuster); ok {
		return localDateWithMonth([]any{self, monthVal})
	}
	if dowVal, ok := getDayOfWeekValueFromObject(adjuster); ok {
		y, m, d, ok2 := getLocalDateYMD(self)
		if ok2 {
			curDow := dayOfWeekDate(y, m, d)
			return localDatePlusDays([]any{self, dowVal - curDow})
		}
	}
	if yearVal, ok := getYearValue(adjuster); ok {
		return localDateWithYear([]any{self, yearVal})
	}

	res := temporalAdjusterAdjustInto([]any{adjuster, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localDateFrom([]any{resObj})
	}
	return res
}

func localDatePlusDays(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plusDays: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusDays: missing daysToAdd")
	}
	daysToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusDays: invalid argument")
	}
	if daysToAdd == 0 {
		return params[0]
	}
	y, m, d, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.plusDays: invalid state")
	}
	epochDay := ymdToEpochDay(y, m, d) + daysToAdd
	newY, newM, newD := epochDayToYMD(epochDay)
	if err := validateLocalDate(newY, newM, newD); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return CreateLocalDate(newY, newM, newD)
}

func localDatePlusWeeks(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plusWeeks: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusWeeks: missing weeksToAdd")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusWeeks: invalid argument")
	}
	return localDatePlusDays([]any{params[0], weeks * 7})
}

func localDatePlusMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plusMonths: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusMonths: missing monthsToAdd")
	}
	monthsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusMonths: invalid argument")
	}
	if monthsToAdd == 0 {
		return params[0]
	}
	y, m, d, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.plusMonths: invalid state")
	}

	monthCount := y*12 + (m - 1) + monthsToAdd
	newY := monthCount / 12
	newM := monthCount % 12
	if newM < 0 {
		newM += 12
		newY--
	}
	newM++

	if newY < localDateMinYear || newY > localDateMaxYear {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", localDateMinYear, localDateMaxYear, newY))
	}
	maxDays := daysInMonthDate(newY, newM)
	if d > maxDays {
		d = maxDays
	}
	return CreateLocalDate(newY, newM, d)
}

func localDatePlusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plusYears: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusYears: missing yearsToAdd")
	}
	yearsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plusYears: invalid argument")
	}
	if yearsToAdd == 0 {
		return params[0]
	}
	y, m, d, okYMD := getLocalDateYMD(params[0].(*object.Object))
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.plusYears: invalid state")
	}
	newY := y + yearsToAdd
	if newY < localDateMinYear || newY > localDateMaxYear {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", localDateMinYear, localDateMaxYear, newY))
	}
	maxDays := daysInMonthDate(newY, m)
	if d > maxDays {
		d = maxDays
	}
	return CreateLocalDate(newY, m, d)
}

func localDatePlusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plus: missing arguments")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.plus: invalid amount")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plus: unit is null")
	}
	unitObj := params[2].(*object.Object)
	unitName := strings.ToUpper(getTemporalUnitName(unitObj))

	switch unitName {
	case "DAYS":
		return localDatePlusDays([]any{params[0], amount})
	case "WEEKS":
		return localDatePlusWeeks([]any{params[0], amount})
	case "MONTHS":
		return localDatePlusMonths([]any{params[0], amount})
	case "YEARS":
		return localDatePlusYears([]any{params[0], amount})
	case "DECADES":
		return localDatePlusYears([]any{params[0], amount * 10})
	case "CENTURIES":
		return localDatePlusYears([]any{params[0], amount * 100})
	case "MILLENNIA":
		return localDatePlusYears([]any{params[0], amount * 1000})
	case "ERAS":
		cfEra := getChronoField("ERA")
		if cfEra == nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "ERA field not found")
		}
		eraRes := localDateGetLong([]any{params[0], cfEra})
		if err, isErr := eraRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		eraVal, _ := getInt64Val(eraRes)
		return localDateWithField([]any{params[0], cfEra, eraVal + amount})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localDatePlusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	res := temporalAmountAddTo([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localDateFrom([]any{resObj})
	}
	return res
}

func localDateMinusDays(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minusDays: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusDays: missing argument")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusDays: invalid argument")
	}
	return localDatePlusDays([]any{params[0], -days})
}

func localDateMinusWeeks(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minusWeeks: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusWeeks: missing argument")
	}
	weeks, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusWeeks: invalid argument")
	}
	return localDatePlusWeeks([]any{params[0], -weeks})
}

func localDateMinusMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minusMonths: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusMonths: missing argument")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusMonths: invalid argument")
	}
	return localDatePlusMonths([]any{params[0], -months})
}

func localDateMinusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minusYears: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusYears: missing argument")
	}
	years, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minusYears: invalid argument")
	}
	return localDatePlusYears([]any{params[0], -years})
}

func localDateMinusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minus: missing arguments")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.minus: invalid amount")
	}
	return localDatePlusUnit([]any{params[0], -amount, params[2]})
}

func localDateMinusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	res := temporalAmountSubtractFrom([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localDateFrom([]any{resObj})
	}
	return res
}

func localDateAtTimeHM(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.atTime: missing arguments")
	}
	h, okH := getInt64Val(params[1])
	m, okM := getInt64Val(params[2])
	if !okH || !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.atTime: invalid arguments")
	}
	tObj := createLocalTime(h, m, 0, 0)
	return createLocalDateTime(params[0].(*object.Object), tObj)
}

func localDateAtTimeHMS(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: self is null")
	}
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.atTime: missing arguments")
	}
	h, okH := getInt64Val(params[1])
	m, okM := getInt64Val(params[2])
	s, okS := getInt64Val(params[3])
	if !okH || !okM || !okS {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.atTime: invalid arguments")
	}
	tObj := createLocalTime(h, m, s, 0)
	return createLocalDateTime(params[0].(*object.Object), tObj)
}

func localDateAtTimeHMSN(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: self is null")
	}
	if len(params) < 5 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.atTime: missing arguments")
	}
	h, okH := getInt64Val(params[1])
	m, okM := getInt64Val(params[2])
	s, okS := getInt64Val(params[3])
	n, okN := getInt64Val(params[4])
	if !okH || !okM || !okS || !okN {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalDate.atTime: invalid arguments")
	}
	tObj := createLocalTime(h, m, s, n)
	return createLocalDateTime(params[0].(*object.Object), tObj)
}

func localDateAtTimeLocalTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: time is null")
	}
	return createLocalDateTime(params[0].(*object.Object), params[1].(*object.Object))
}

func localDateAtTimeOffsetTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atTime: offsetTime is null")
	}
	self := params[0].(*object.Object)
	ot := params[1].(*object.Object)
	var timeObj *object.Object
	var offsetObj *object.Object
	if tFld, ok := ot.FieldTable["time"]; ok {
		if t, ok2 := tFld.Fvalue.(*object.Object); ok2 {
			timeObj = t
		}
	}
	if oFld, ok := ot.FieldTable["offset"]; ok {
		if o, ok2 := oFld.Fvalue.(*object.Object); ok2 {
			offsetObj = o
		}
	}
	ldt := createLocalDateTime(self, timeObj)
	classNameOffsetDateTime := "java/time/OffsetDateTime"
	odt := object.MakeEmptyObjectWithClassName(&classNameOffsetDateTime)
	odt.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	odt.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return odt
}

func localDateAtStartOfDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atStartOfDay: self is null")
	}
	tObj := createLocalTime(0, 0, 0, 0)
	return createLocalDateTime(params[0].(*object.Object), tObj)
}

func localDateAtStartOfDayZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atStartOfDay: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.atStartOfDay: zone is null")
	}
	self := params[0].(*object.Object)
	zoneObj := params[1].(*object.Object)
	tObj := createLocalTime(0, 0, 0, 0)
	ldt := createLocalDateTime(self, tObj)

	classNameZonedDateTime := "java/time/ZonedDateTime"
	zdt := object.MakeEmptyObjectWithClassName(&classNameZonedDateTime)
	zdt.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	zdt.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zoneObj}
	zdt.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: getZoneRulesOffset(zoneObj)}
	return zdt
}

func getZoneRulesOffset(zoneObj *object.Object) *object.Object {
	if zoneObj != nil && zoneObj.FieldTable != nil {
		if offFld, ok := zoneObj.FieldTable["offset"]; ok {
			if off, ok2 := offFld.Fvalue.(*object.Object); ok2 {
				return off
			}
		}
	}
	classNameZoneOffset := "java/time/ZoneOffset"
	zo := object.MakeEmptyObjectWithClassName(&classNameZoneOffset)
	zo.FieldTable["totalSeconds"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	return zo
}

func localDateAdjustInto(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(self)
	cfEpochDay := getChronoField("EPOCH_DAY")
	return temporalWithField([]any{temporal, cfEpochDay, epochDay})
}

func localDateQuery(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.query: query is null")
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
		return ghelpers.CreateTemporalUnit("Days")
	case "localDate":
		return self
	case "zoneId", "zone", "offset", "localTime":
		return nil
	default:
		return executeTemporalQuery(query, self)
	}
}

func localDateUntilUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.until: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.until: unit is null")
	}
	self := params[0].(*object.Object)
	endObj := params[1].(*object.Object)
	unitObj := params[2].(*object.Object)

	y1, m1, d1, ok1 := getLocalDateYMD(self)
	endLocalDate := localDateFrom([]any{endObj})
	if err, isErr := endLocalDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	y2, m2, d2, ok2 := getLocalDateYMD(endLocalDate.(*object.Object))
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.until: invalid operand")
	}

	unitName := strings.ToUpper(getTemporalUnitName(unitObj))
	ep1 := ymdToEpochDay(y1, m1, d1)
	ep2 := ymdToEpochDay(y2, m2, d2)
	daysUntil := ep2 - ep1

	switch unitName {
	case "DAYS":
		return daysUntil
	case "WEEKS":
		return daysUntil / 7
	case "MONTHS":
		totalMonths1 := y1*12 + m1 - 1
		totalMonths2 := y2*12 + m2 - 1
		diff := totalMonths2 - totalMonths1
		if diff > 0 && d2 < d1 {
			diff--
		} else if diff < 0 && d2 > d1 {
			diff++
		}
		return diff
	case "YEARS":
		diffYears := y2 - y1
		if diffYears > 0 && (m2 < m1 || (m2 == m1 && d2 < d1)) {
			diffYears--
		} else if diffYears < 0 && (m2 > m1 || (m2 == m1 && d2 > d1)) {
			diffYears++
		}
		return diffYears
	case "DECADES":
		years := (y2 - y1)
		if years > 0 && (m2 < m1 || (m2 == m1 && d2 < d1)) {
			years--
		} else if years < 0 && (m2 > m1 || (m2 == m1 && d2 > d1)) {
			years++
		}
		return years / 10
	case "CENTURIES":
		years := (y2 - y1)
		if years > 0 && (m2 < m1 || (m2 == m1 && d2 < d1)) {
			years--
		} else if years < 0 && (m2 > m1 || (m2 == m1 && d2 > d1)) {
			years++
		}
		return years / 100
	case "MILLENNIA":
		years := (y2 - y1)
		if years > 0 && (m2 < m1 || (m2 == m1 && d2 < d1)) {
			years--
		} else if years < 0 && (m2 > m1 || (m2 == m1 && d2 > d1)) {
			years++
		}
		return years / 1000
	case "ERAS":
		era1 := int64(0)
		if y1 >= 1 {
			era1 = 1
		}
		era2 := int64(0)
		if y2 >= 1 {
			era2 = 1
		}
		return era2 - era1
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localDateUntilPeriod(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.until: endDateExclusive is null")
	}
	self := params[0].(*object.Object)
	endObj := params[1].(*object.Object)

	y1, m1, d1, ok1 := getLocalDateYMD(self)
	endLocalDate := localDateFrom([]any{endObj})
	if err, isErr := endLocalDate.(*ghelpers.GErrBlk); isErr {
		return err
	}
	y2, m2, d2, ok2 := getLocalDateYMD(endLocalDate.(*object.Object))
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "LocalDate.until: invalid date")
	}

	totalMonths1 := y1*12 + m1 - 1
	totalMonths2 := y2*12 + m2 - 1
	totalMonths := totalMonths2 - totalMonths1
	days := d2 - d1

	if totalMonths > 0 && days < 0 {
		totalMonths--
		prevMonthDate := CreateLocalDate(y2, m2, 1)
		prevMonthDate = localDateMinusMonths([]any{prevMonthDate, int64(1)}).(*object.Object)
		py, pm, _, _ := getLocalDateYMD(prevMonthDate)
		days += daysInMonthDate(py, pm)
	} else if totalMonths < 0 && days > 0 {
		totalMonths++
		days -= daysInMonthDate(y1, m1)
	}

	years := totalMonths / 12
	months := totalMonths % 12

	classNamePeriod := "java/time/Period"
	pObj := object.MakeEmptyObjectWithClassName(&classNamePeriod)
	pObj.FieldTable["years"] = object.Field{Ftype: types.Int, Fvalue: years}
	pObj.FieldTable["months"] = object.Field{Ftype: types.Int, Fvalue: months}
	pObj.FieldTable["days"] = object.Field{Ftype: types.Int, Fvalue: days}
	return pObj
}

func localDateDatesUntil(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.datesUntil: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.datesUntil: endExclusive is null")
	}
	classNameStream := "java/util/stream/Stream"
	streamObj := object.MakeEmptyObjectWithClassName(&classNameStream)
	return streamObj
}

func localDateDatesUntilPeriod(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.datesUntil: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.datesUntil: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalDate.datesUntil: step is null")
	}
	stepObj := params[2].(*object.Object)
	var y, m, d int64
	if fld, ok := stepObj.FieldTable["years"]; ok {
		y, _ = getInt64Val(fld.Fvalue)
	}
	if fld, ok := stepObj.FieldTable["months"]; ok {
		m, _ = getInt64Val(fld.Fvalue)
	}
	if fld, ok := stepObj.FieldTable["days"]; ok {
		d, _ = getInt64Val(fld.Fvalue)
	}
	if y == 0 && m == 0 && d == 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "step cannot be zero")
	}
	classNameStream := "java/util/stream/Stream"
	streamObj := object.MakeEmptyObjectWithClassName(&classNameStream)
	return streamObj
}

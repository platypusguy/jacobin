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
	"time"
)

// Load_Time_Chrono_ChronoLocalDate loads method signatures for java/time/chrono/ChronoLocalDate.
func Load_Time_Chrono_ChronoLocalDate() {
	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.atTime(Ljava/time/LocalTime;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateAtTime,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.compareTo(Ljava/time/chrono/ChronoLocalDate;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateEquals,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateFormat,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateFrom,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateGetChronology,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.getEra()Ljava/time/chrono/Era;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateGetEra,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateHashCode,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.isAfter(Ljava/time/chrono/ChronoLocalDate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateIsAfter,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.isBefore(Ljava/time/chrono/ChronoLocalDate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateIsBefore,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.isEqual(Ljava/time/chrono/ChronoLocalDate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateIsEqual,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.isLeapYear()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateIsLeapYear,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.lengthOfMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateLengthOfMonth,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.lengthOfYear()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateLengthOfYear,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoLocalDateMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoLocalDatePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDatePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateQuery,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.timeLineOrder()Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateTimelineOrder,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.toEpochDay()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateToEpochDay,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoLocalDateToString,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronoLocalDateUntilUnit,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.until(Ljava/time/chrono/ChronoLocalDate;)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateUntilPeriod,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoLocalDateWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoLocalDate.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoLocalDateWithField,
		}
}

func GetChronoLocalDateChronology(obj *object.Object) *object.Object {
	if obj != nil && obj.FieldTable != nil {
		if chronoFld, ok := obj.FieldTable["chronology"]; ok {
			if chronoObj, ok := chronoFld.Fvalue.(*object.Object); ok && chronoObj != nil {
				return chronoObj
			}
		}
	}
	return GetIsoChronology()
}

// chronoLocalDateAdjustInto adjusts the specified temporal object to have the same date as this object.
// default Temporal adjustInto(Temporal temporal)
func chronoLocalDateAdjustInto(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	target := params[1].(*object.Object)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(self)
	year, month, day := ghelpers.GetChronoLocalDateYMD(self)

	if target.FieldTable != nil {
		target.FieldTable["epochDay"] = object.Field{Ftype: types.Long, Fvalue: epochDay}
		target.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: year}
		target.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: month}
		target.FieldTable["day"] = object.Field{Ftype: types.Int, Fvalue: day}
	}
	return target
}

// chronoLocalDateAtTime combines this date with a time to create a ChronoLocalDateTime.
// default ChronoLocalDateTime<?> atTime(LocalTime localTime)
func chronoLocalDateAtTime(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.atTime: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.atTime: localTime is null")
	}
	self := params[0].(*object.Object)
	timeObj := params[1].(*object.Object)

	return ghelpers.CreateChronoLocalDateTimeFromDateAndTime(self, timeObj)
}

// chronoLocalDateCompareTo compares this date to another date, including the chronology.
// default int compareTo(ChronoLocalDate other)
func chronoLocalDateCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	epoch1 := ghelpers.GetChronoLocalDateEpochDay(self)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(other)
	if epoch1 < epoch2 {
		return int64(-1)
	} else if epoch1 > epoch2 {
		return int64(1)
	}

	chrono1 := GetChronoLocalDateChronology(self)
	chrono2 := GetChronoLocalDateChronology(other)
	id1 := getChronologyId(chrono1)
	id2 := getChronologyId(chrono2)
	if id1 < id2 {
		return int64(-1)
	} else if id1 > id2 {
		return int64(1)
	}
	return int64(0)
}

// chronoLocalDateEquals checks if this date is equal to another date, including the chronology.
// boolean equals(Object obj)
func chronoLocalDateEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.equals: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	res := chronoLocalDateCompareTo([]any{self, other})
	if cmp, ok := res.(int64); ok && cmp == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateFormat formats this date using the specified formatter.
// default String format(DateTimeFormatter formatter)
func chronoLocalDateFormat(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.format: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.format: formatter is null")
	}
	self := params[0].(*object.Object)
	year, month, day := ghelpers.GetChronoLocalDateYMD(self)
	formatted := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return object.StringObjectFromGoString(formatted)
}

// chronoLocalDateFrom obtains an instance of ChronoLocalDate from a temporal object.
// static ChronoLocalDate from(TemporalAccessor temporal)
func chronoLocalDateFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.from: temporal is null")
	}
	temporal := params[0].(*object.Object)

	chrono := GetIsoChronology()
	if temporal.FieldTable != nil {
		if cFld, ok := temporal.FieldTable["chronology"]; ok {
			if cObj, ok := cFld.Fvalue.(*object.Object); ok && cObj != nil {
				chrono = cObj
			}
		}
	}

	year, month, day := ghelpers.GetChronoLocalDateYMD(temporal)
	return createChronoLocalDate(chrono, year, month, day)
}

// chronoLocalDateGetChronology gets the chronology of this date.
// Chronology getChronology()
func chronoLocalDateGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.getChronology: self is null")
	}
	self := params[0].(*object.Object)
	return GetChronoLocalDateChronology(self)
}

// chronoLocalDateGetEra gets the era, as defined by the chronology.
// default Era getEra()
func chronoLocalDateGetEra(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.getEra: self is null")
	}
	self := params[0].(*object.Object)
	year, _, _ := ghelpers.GetChronoLocalDateYMD(self)
	if year >= 1 {
		return createEra(1, "CE")
	}
	return createEra(0, "BCE")
}

// chronoLocalDateHashCode a suitable hash code.
// int hashCode()
func chronoLocalDateHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(self)
	chrono := GetChronoLocalDateChronology(self)
	chronoHash := int64(0)
	if chronoId := getChronologyId(chrono); len(chronoId) > 0 {
		var h int64
		for i := 0; i < len(chronoId); i++ {
			h = 31*h + int64(chronoId[i])
		}
		chronoHash = h
	}
	return int64(int32(epochDay^(epochDay>>32)) ^ int32(chronoHash))
}

// chronoLocalDateIsAfter checks if this date is after the specified date ignoring the chronology.
// default boolean isAfter(ChronoLocalDate other)
func chronoLocalDateIsAfter(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isAfter: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isAfter: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	epoch1 := ghelpers.GetChronoLocalDateEpochDay(self)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(other)
	if epoch1 > epoch2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateIsBefore checks if this date is before the specified date ignoring the chronology.
// default boolean isBefore(ChronoLocalDate other)
func chronoLocalDateIsBefore(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isBefore: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isBefore: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	epoch1 := ghelpers.GetChronoLocalDateEpochDay(self)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(other)
	if epoch1 < epoch2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateIsEqual checks if this date is equal to the specified date ignoring the chronology.
// default boolean isEqual(ChronoLocalDate other)
func chronoLocalDateIsEqual(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isEqual: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isEqual: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	epoch1 := ghelpers.GetChronoLocalDateEpochDay(self)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(other)
	if epoch1 == epoch2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateIsLeapYear checks if the year is a leap year, as defined by the chronology.
// default boolean isLeapYear()
func chronoLocalDateIsLeapYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isLeapYear: self is null")
	}
	self := params[0].(*object.Object)
	year, _, _ := ghelpers.GetChronoLocalDateYMD(self)
	if ghelpers.IsLeapYearHelper(year) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoLocalDateIsSupportedField checks if the specified field is supported.
// default boolean isSupported(TemporalField field)
func chronoLocalDateIsSupportedField(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isSupported: self is null")
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
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY", "ALIGNED_WEEK_OF_MONTH",
		"ALIGNED_WEEK_OF_YEAR", "MONTH_OF_YEAR", "PROLEPTIC_MONTH", "YEAR_OF_ERA",
		"YEAR", "ERA":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// chronoLocalDateIsSupportedUnit checks if the specified unit is supported.
// default boolean isSupported(TemporalUnit unit)
func chronoLocalDateIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.isSupported: self is null")
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
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// chronoLocalDateLengthOfMonth returns the length of the month represented by this date.
// int lengthOfMonth()
func chronoLocalDateLengthOfMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.lengthOfMonth: self is null")
	}
	self := params[0].(*object.Object)
	year, month, _ := ghelpers.GetChronoLocalDateYMD(self)
	return ghelpers.DaysInMonthHelper(year, month)
}

// chronoLocalDateLengthOfYear returns the length of the year represented by this date.
// default int lengthOfYear()
func chronoLocalDateLengthOfYear(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.lengthOfYear: self is null")
	}
	self := params[0].(*object.Object)
	year, _, _ := ghelpers.GetChronoLocalDateYMD(self)
	if ghelpers.IsLeapYearHelper(year) {
		return int64(366)
	}
	return int64(365)
}

// chronoLocalDateMinusUnit returns an object of the same type with the specified period subtracted.
// default ChronoLocalDate minus(long amountToSubtract, TemporalUnit unit)
func chronoLocalDateMinusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.minus: self is null")
	}
	amount, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoLocalDate.minus: invalid amount")
	}
	unitObj, _ := params[2].(*object.Object)
	return chronoLocalDatePlusUnit([]any{params[0], -amount, unitObj})
}

// chronoLocalDateMinusAmount returns an object of the same type with the specified amount subtracted.
// default ChronoLocalDate minus(TemporalAmount amount)
func chronoLocalDateMinusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	days := ghelpers.GetTemporalAmountDays(amount)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(self)
	newEpochDay := epochDay - days
	chrono := GetChronoLocalDateChronology(self)

	t := time.Unix(newEpochDay*86400, 0).UTC()
	return createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
}

// chronoLocalDatePlusUnit returns an object of the same type with the specified period added.
// default ChronoLocalDate plus(long amountToAdd, TemporalUnit unit)
func chronoLocalDatePlusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.plus: self is null")
	}
	self := params[0].(*object.Object)
	amount, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoLocalDate.plus: invalid amount")
	}
	unitObj, _ := params[2].(*object.Object)
	unitName := getTemporalUnitName(unitObj)

	year, month, day := ghelpers.GetChronoLocalDateYMD(self)
	chrono := GetChronoLocalDateChronology(self)

	switch unitName {
	case "DAYS":
		epochDay := ghelpers.GetChronoLocalDateEpochDay(self) + amount
		t := time.Unix(epochDay*86400, 0).UTC()
		return createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
	case "WEEKS":
		epochDay := ghelpers.GetChronoLocalDateEpochDay(self) + amount*7
		t := time.Unix(epochDay*86400, 0).UTC()
		return createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
	case "MONTHS":
		totalMonths := year*12 + (month - 1) + amount
		newYear := totalMonths / 12
		newMonth := (totalMonths % 12) + 1
		maxDays := ghelpers.DaysInMonthHelper(newYear, newMonth)
		newDay := min(day, maxDays)
		return createChronoLocalDate(chrono, newYear, newMonth, newDay)
	case "YEARS":
		newYear := year + amount
		maxDays := ghelpers.DaysInMonthHelper(newYear, month)
		newDay := min(day, maxDays)
		return createChronoLocalDate(chrono, newYear, month, newDay)
	case "DECADES":
		newYear := year + amount*10
		return createChronoLocalDate(chrono, newYear, month, day)
	case "CENTURIES":
		newYear := year + amount*100
		return createChronoLocalDate(chrono, newYear, month, day)
	case "MILLENNIA":
		newYear := year + amount*1000
		return createChronoLocalDate(chrono, newYear, month, day)
	case "ERAS":
		return self
	default:
		// Default to days addition
		epochDay := ghelpers.GetChronoLocalDateEpochDay(self) + amount
		t := time.Unix(epochDay*86400, 0).UTC()
		return createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
	}
}

// chronoLocalDatePlusAmount returns an object of the same type with the specified amount added.
// default ChronoLocalDate plus(TemporalAmount amount)
func chronoLocalDatePlusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	days := ghelpers.GetTemporalAmountDays(amount)
	epochDay := ghelpers.GetChronoLocalDateEpochDay(self)
	newEpochDay := epochDay + days
	chrono := GetChronoLocalDateChronology(self)

	t := time.Unix(newEpochDay*86400, 0).UTC()
	return createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
}

// chronoLocalDateQuery queries this date using the specified query.
// default <R> R query(TemporalQuery<R> query)
func chronoLocalDateQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.query: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	if query.FieldTable != nil {
		if nameFld, ok := query.FieldTable["name"]; ok {
			if nameObj, ok := nameFld.Fvalue.(*object.Object); ok && object.IsStringObject(nameObj) {
				qName := object.GoStringFromStringObject(nameObj)
				switch qName {
				case "CHRONOLOGY", "chronology":
					return GetChronoLocalDateChronology(self)
				case "PRECISION", "precision":
					return ghelpers.CreateTemporalUnit("Days")
				case "LOCAL_DATE", "localDate":
					return self
				}
			}
		}
	}

	return temporalQueryQueryFrom([]any{query, self})
}

// chronoLocalDateTimelineOrder gets a comparator that compares ChronoLocalDate in time-line order ignoring the chronology.
// static Comparator<ChronoLocalDate> timeLineOrder()
func chronoLocalDateTimelineOrder([]any) any {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameChronoLocalDateComparator)
	return obj
}

// chronoLocalDateToEpochDay converts this date to the Epoch Day.
// default long toEpochDay()
func chronoLocalDateToEpochDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.toEpochDay: self is null")
	}
	self := params[0].(*object.Object)
	return ghelpers.GetChronoLocalDateEpochDay(self)
}

// chronoLocalDateToString outputs this date as a String.
// String toString()
func chronoLocalDateToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.toString: self is null")
	}
	self := params[0].(*object.Object)
	year, month, day := ghelpers.GetChronoLocalDateYMD(self)
	chrono := GetChronoLocalDateChronology(self)
	chronoId := getChronologyId(chrono)

	var str string
	if chronoId == "ISO" || chronoId == "" {
		str = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	} else {
		str = fmt.Sprintf("%s %04d-%02d-%02d", chronoId, year, month, day)
	}
	return object.StringObjectFromGoString(str)
}

// chronoLocalDateUntilUnit calculates the amount of time until another date in terms of the specified unit.
// long until(Temporal endExclusive, TemporalUnit unit)
func chronoLocalDateUntilUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.until: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.until: endExclusive is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.until: unit is null")
	}
	self := params[0].(*object.Object)
	end := params[1].(*object.Object)
	unit := params[2].(*object.Object)

	epoch1 := ghelpers.GetChronoLocalDateEpochDay(self)
	epoch2 := ghelpers.GetChronoLocalDateEpochDay(end)
	diffDays := epoch2 - epoch1

	unitName := getTemporalUnitName(unit)
	switch unitName {
	case "DAYS":
		return diffDays
	case "WEEKS":
		return diffDays / 7
	case "MONTHS":
		y1, m1, _ := ghelpers.GetChronoLocalDateYMD(self)
		y2, m2, _ := ghelpers.GetChronoLocalDateYMD(end)
		return (y2-y1)*12 + (m2 - m1)
	case "YEARS":
		y1, _, _ := ghelpers.GetChronoLocalDateYMD(self)
		y2, _, _ := ghelpers.GetChronoLocalDateYMD(end)
		return y2 - y1
	case "DECADES":
		y1, _, _ := ghelpers.GetChronoLocalDateYMD(self)
		y2, _, _ := ghelpers.GetChronoLocalDateYMD(end)
		return (y2 - y1) / 10
	case "CENTURIES":
		y1, _, _ := ghelpers.GetChronoLocalDateYMD(self)
		y2, _, _ := ghelpers.GetChronoLocalDateYMD(end)
		return (y2 - y1) / 100
	case "MILLENNIA":
		y1, _, _ := ghelpers.GetChronoLocalDateYMD(self)
		y2, _, _ := ghelpers.GetChronoLocalDateYMD(end)
		return (y2 - y1) / 1000
	default:
		return diffDays
	}
}

// chronoLocalDateUntilPeriod calculates the period between this date and another date as a ChronoPeriod.
// ChronoPeriod until(ChronoLocalDate endDateExclusive)
func chronoLocalDateUntilPeriod(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.until: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.until: endDateExclusive is null")
	}
	self := params[0].(*object.Object)
	end := params[1].(*object.Object)

	y1, m1, d1 := ghelpers.GetChronoLocalDateYMD(self)
	y2, m2, d2 := ghelpers.GetChronoLocalDateYMD(end)

	years := y2 - y1
	months := m2 - m1
	days := d2 - d1

	if days < 0 {
		months--
		prevMonthDays := ghelpers.DaysInMonthHelper(y2, m2-1)
		days += prevMonthDays
	}
	if months < 0 {
		years--
		months += 12
	}

	chrono := GetChronoLocalDateChronology(self)
	return ghelpers.CreateChronoPeriodYMD(chrono, years, months, days)
}

// chronoLocalDateWithAdjuster returns an adjusted object of the same type as this object with the adjustment made.
// default ChronoLocalDate with(TemporalAdjuster adjuster)
func chronoLocalDateWithAdjuster(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	res := temporalAdjusterAdjustInto([]any{adjuster, self})
	if resObj, ok := res.(*object.Object); ok {
		return resObj
	}
	return self
}

// chronoLocalDateWithField returns an object of the same type as this object with the specified field altered.
// default ChronoLocalDate with(TemporalField field, long newValue)
func chronoLocalDateWithField(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoLocalDate.with: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	newValue, ok := ghelpers.GetInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoLocalDate.with: invalid newValue")
	}

	fieldName := getTemporalFieldName(fieldObj)
	year, month, day := ghelpers.GetChronoLocalDateYMD(self)
	chrono := GetChronoLocalDateChronology(self)

	switch fieldName {
	case "DAY_OF_MONTH":
		return createChronoLocalDate(chrono, year, month, newValue)
	case "MONTH_OF_YEAR":
		return createChronoLocalDate(chrono, year, newValue, day)
	case "YEAR", "YEAR_OF_ERA":
		return createChronoLocalDate(chrono, newValue, month, day)
	case "EPOCH_DAY":
		t := time.Unix(newValue*86400, 0).UTC()
		return createChronoLocalDate(chrono, int64(t.Year()), int64(t.Month()), int64(t.Day()))
	default:
		return temporalFieldAdjustInto([]any{fieldObj, self, newValue})
	}
}

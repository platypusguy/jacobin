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

var classNameYearVal = "java/time/Year"

const (
	yearMinValue = -999999999
	yearMaxValue = 999999999
)

var yearMutex sync.Mutex
var yearInitialized bool = false

func ensureYearInited() {
	yearMutex.Lock()
	defer yearMutex.Unlock()
	if yearInitialized {
		return
	}

	_ = statics.AddStatic(classNameYearVal+".MIN_VALUE", statics.Static{
		Type:  "I",
		Value: int64(yearMinValue),
	})
	_ = statics.AddStatic(classNameYearVal+".MAX_VALUE", statics.Static{
		Type:  "I",
		Value: int64(yearMaxValue),
	})

	yearInitialized = true
}

func createYear(year int64) *object.Object {
	ensureYearInited()
	obj := object.MakeEmptyObjectWithClassName(&classNameYearVal)
	obj.FieldTable["year"] = object.Field{
		Ftype:  types.Int,
		Fvalue: year,
	}
	obj.FieldTable["YEAR"] = object.Field{
		Ftype:  types.Long,
		Fvalue: year,
	}
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.Int,
		Fvalue: year,
	}
	obj.FieldTable["precision"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: ghelpers.CreateTemporalUnit("Years"),
	}
	return obj
}

func getYearValue(targetObject *object.Object) (int64, bool) {
	if targetObject == nil || targetObject.FieldTable == nil {
		return 0, false
	}
	if yearField, exists := targetObject.FieldTable["year"]; exists {
		if val, ok := getInt64Val(yearField.Fvalue); ok {
			return val, true
		}
	}
	if valField, exists := targetObject.FieldTable["value"]; exists {
		if val, ok := getInt64Val(valField.Fvalue); ok {
			return val, true
		}
	}
	if yearField, exists := targetObject.FieldTable["YEAR"]; exists {
		if val, ok := getInt64Val(yearField.Fvalue); ok {
			return val, true
		}
	}
	return 0, false
}

func isLeapYearVal(year int64) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func formatYearString(year int64) string {
	absYear := year
	if absYear < 0 {
		absYear = -absYear
	}
	if absYear < 1000 {
		if year < 0 {
			return fmt.Sprintf("-%04d", absYear)
		}
		return fmt.Sprintf("%04d", year)
	} else if absYear <= 9999 {
		if year < 0 {
			return fmt.Sprintf("-%04d", absYear)
		}
		return fmt.Sprintf("%04d", year)
	} else {
		if year > 0 {
			return fmt.Sprintf("+%d", year)
		}
		return fmt.Sprintf("%d", year)
	}
}

func createYearMonth(year, month int64) *object.Object {
	classNameYearMonth := "java/time/YearMonth"
	obj := object.MakeEmptyObjectWithClassName(&classNameYearMonth)
	obj.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: year}
	obj.FieldTable["month"] = object.Field{Ftype: types.Int, Fvalue: month}
	return obj
}

func getChronoField(name string) *object.Object {
	ensureChronoFieldInited()
	for index, definition := range chronoFieldDefinitions {
		if definition.enumName == name {
			return chronoFieldInstances[index]
		}
	}
	return nil
}

// Load_Time_Year registers all java/time/Year methods.
func Load_Time_Year() {
	ghelpers.MethodSignatures["java/time/Year.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearClinit,
		}

	ghelpers.MethodSignatures["java/time/Year.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearInit,
		}

	ghelpers.MethodSignatures["java/time/Year.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/Year.atDay(I)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearAtDay,
		}

	ghelpers.MethodSignatures["java/time/Year.atMonth(I)Ljava/time/YearMonth;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearAtMonthInt,
		}

	ghelpers.MethodSignatures["java/time/Year.atMonth(Ljava/time/Month;)Ljava/time/YearMonth;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearAtMonthMonth,
		}

	ghelpers.MethodSignatures["java/time/Year.atMonthDay(Ljava/time/MonthDay;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearAtMonthDay,
		}

	ghelpers.MethodSignatures["java/time/Year.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Year.compareTo(Ljava/time/Year;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Year.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearEquals,
		}

	ghelpers.MethodSignatures["java/time/Year.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearFormat,
		}

	ghelpers.MethodSignatures["java/time/Year.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearFrom,
		}

	ghelpers.MethodSignatures["java/time/Year.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearGet,
		}

	ghelpers.MethodSignatures["java/time/Year.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearGetLong,
		}

	ghelpers.MethodSignatures["java/time/Year.getValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearGetValue,
		}

	ghelpers.MethodSignatures["java/time/Year.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearHashCode,
		}

	ghelpers.MethodSignatures["java/time/Year.isAfter(Ljava/time/Year;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearIsAfter,
		}

	ghelpers.MethodSignatures["java/time/Year.isBefore(Ljava/time/Year;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearIsBefore,
		}

	ghelpers.MethodSignatures["java/time/Year.isLeap()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearIsLeapInstance,
		}

	ghelpers.MethodSignatures["java/time/Year.isLeap(J)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  yearIsLeapStatic,
		}

	ghelpers.MethodSignatures["java/time/Year.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/Year.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/Year.isValidMonthDay(Ljava/time/MonthDay;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearIsValidMonthDay,
		}

	ghelpers.MethodSignatures["java/time/Year.length()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearLength,
		}

	ghelpers.MethodSignatures["java/time/Year.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  yearMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/Year.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  yearMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/Year.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/Year.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/Year.minusYears(J)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  yearMinusYears,
		}

	ghelpers.MethodSignatures["java/time/Year.now()Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearNow,
		}

	ghelpers.MethodSignatures["java/time/Year.now(Ljava/time/Clock;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearNowClock,
		}

	ghelpers.MethodSignatures["java/time/Year.now(Ljava/time/ZoneId;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearNowZone,
		}

	ghelpers.MethodSignatures["java/time/Year.of(I)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearOf,
		}

	ghelpers.MethodSignatures["java/time/Year.parse(Ljava/lang/CharSequence;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearParse,
		}

	ghelpers.MethodSignatures["java/time/Year.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  yearParseFormatter,
		}

	ghelpers.MethodSignatures["java/time/Year.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  yearPlusUnit,
		}

	ghelpers.MethodSignatures["java/time/Year.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  yearPlusUnit,
		}

	ghelpers.MethodSignatures["java/time/Year.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearPlusAmount,
		}

	ghelpers.MethodSignatures["java/time/Year.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearPlusAmount,
		}

	ghelpers.MethodSignatures["java/time/Year.plusYears(J)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  yearPlusYears,
		}

	ghelpers.MethodSignatures["java/time/Year.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearQuery,
		}

	ghelpers.MethodSignatures["java/time/Year.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearRange,
		}

	ghelpers.MethodSignatures["java/time/Year.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  yearToString,
		}

	ghelpers.MethodSignatures["java/time/Year.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  yearUntil,
		}

	ghelpers.MethodSignatures["java/time/Year.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/Year.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  yearWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/Year.with(Ljava/time/temporal/TemporalField;J)Ljava/time/Year;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  yearWithField,
		}

	ghelpers.MethodSignatures["java/time/Year.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  yearWithField,
		}
}

func yearClinit([]any) any {
	ensureYearInited()
	return nil
}

func yearInit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.<init>: self is null")
	}
	self := params[0].(*object.Object)
	val, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.<init>: invalid year argument")
	}
	if val < yearMinValue || val > yearMaxValue {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", yearMinValue, yearMaxValue, val))
	}
	self.FieldTable["year"] = object.Field{Ftype: types.Int, Fvalue: val}
	self.FieldTable["YEAR"] = object.Field{Ftype: types.Long, Fvalue: val}
	self.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: val}
	self.FieldTable["precision"] = object.Field{Ftype: types.Ref, Fvalue: ghelpers.CreateTemporalUnit("Years")}
	return nil
}

func yearOf(params []any) any {
	ensureYearInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.of: missing argument")
	}
	yearVal, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.of: invalid year argument")
	}
	if yearVal < yearMinValue || yearVal > yearMaxValue {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", yearMinValue, yearMaxValue, yearVal))
	}
	return createYear(yearVal)
}

func yearNow([]any) any {
	return createYear(int64(time.Now().Year()))
}

func yearNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.now: clock is null")
	}
	return createYear(int64(time.Now().Year()))
}

func yearNowZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.now: zone is null")
	}
	return createYear(int64(time.Now().Year()))
}

func yearFrom(params []any) any {
	ensureYearInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.from: temporal is null")
	}
	temporal := params[0].(*object.Object)

	if object.GoStringFromStringPoolIndex(temporal.KlassName) == classNameYearVal {
		if y, ok := getYearValue(temporal); ok {
			return createYear(y)
		}
		return temporal
	}

	cfYear := getChronoField("YEAR")
	if cfYear != nil {
		res := temporalAccessorGet([]any{temporal, cfYear})
		if _, isErr := res.(*ghelpers.GErrBlk); !isErr {
			if val, ok := getInt64Val(res); ok {
				return yearOf([]any{val})
			}
		}
	}

	if temporal.FieldTable != nil {
		if yearFld, okY := temporal.FieldTable["year"]; okY {
			if yearVal, ok := getInt64Val(yearFld.Fvalue); ok {
				return yearOf([]any{yearVal})
			}
		}
		if epochFld, okE := temporal.FieldTable["epochDay"]; okE {
			if epochDay, ok := getInt64Val(epochFld.Fvalue); ok {
				t := time.Unix(epochDay*86400, 0).UTC()
				return createYear(int64(t.Year()))
			}
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain Year from TemporalAccessor")
}

func yearGetValue(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.getValue: self is null")
	}
	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.getValue: invalid instance")
	}
	return year
}

func yearIsLeapInstance(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isLeap: self is null")
	}
	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.isLeap: invalid instance")
	}
	return isLeapYearVal(year)
}

func yearIsLeapStatic(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.isLeap: missing argument")
	}
	year, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.isLeap: invalid argument")
	}
	return isLeapYearVal(year)
}

func yearIsValidMonthDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isValidMonthDay: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return false
	}
	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return false
	}
	monthDay := params[1].(*object.Object)
	var month, day int64 = 1, 1
	if monthDay.FieldTable != nil {
		if mFld, okM := monthDay.FieldTable["month"]; okM {
			month, _ = getInt64Val(mFld.Fvalue)
		}
		if dFld, okD := monthDay.FieldTable["day"]; okD {
			day, _ = getInt64Val(dFld.Fvalue)
		} else if dFld, okD2 := monthDay.FieldTable["dayOfMonth"]; okD2 {
			day, _ = getInt64Val(dFld.Fvalue)
		}
	}
	if month == 2 && day == 29 {
		return isLeapYearVal(year)
	}
	return true
}

func yearLength(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.length: self is null")
	}
	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.length: invalid instance")
	}
	if isLeapYearVal(year) {
		return int64(366)
	}
	return int64(365)
}

func yearIsSupportedField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isSupported: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return false
	}

	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	switch strings.ToUpper(fieldName) {
	case "YEAR", "YEAR_OF_ERA", "ERA":
		return true
	default:
		return false
	}
}

func yearIsSupportedUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isSupported: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return false
	}

	unitObj := params[1].(*object.Object)
	unitName := getTemporalUnitName(unitObj)
	switch strings.ToUpper(unitName) {
	case "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return true
	default:
		return false
	}
}

func yearRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.range: field is null")
	}

	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.range: invalid instance")
	}

	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	switch strings.ToUpper(fieldName) {
	case "YEAR":
		return createValueRange(yearMinValue, yearMinValue, yearMaxValue, yearMaxValue)
	case "YEAR_OF_ERA":
		if year <= 0 {
			return createValueRange(1, 1, 1000000000, 1000000000)
		}
		return createValueRange(1, 1, yearMaxValue, yearMaxValue)
	case "ERA":
		return createValueRange(0, 0, 1, 1)
	default:
		if object.GoStringFromStringPoolIndex(fieldObj.KlassName) == classNameChronoField {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
		}
		return temporalAccessorRange(params)
	}
}

func yearGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.get: field is null")
	}

	res := yearGetLong(params)
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	val, _ := getInt64Val(res)
	return int64(int32(val))
}

func yearGetLong(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.getLong: field is null")
	}

	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.getLong: invalid instance")
	}

	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	switch strings.ToUpper(fieldName) {
	case "YEAR":
		return year
	case "YEAR_OF_ERA":
		if year < 1 {
			return 1 - year
		}
		return year
	case "ERA":
		if year < 1 {
			return int64(0)
		}
		return int64(1)
	default:
		if object.GoStringFromStringPoolIndex(fieldObj.KlassName) == classNameChronoField {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
		}
		return temporalAccessorGetLong(params)
	}
}

func yearWithAdjuster(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	res := temporalAdjusterAdjustInto([]any{adjuster, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return yearFrom([]any{resObj})
	}
	return res
}

func yearWithField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.with: field is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.with: missing newValue")
	}

	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.with: invalid instance")
	}

	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	newValue, okVal := getInt64Val(params[2])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.with: invalid newValue")
	}

	switch strings.ToUpper(fieldName) {
	case "YEAR":
		return yearOf([]any{newValue})
	case "YEAR_OF_ERA":
		maxEra := int64(yearMaxValue)
		if year <= 0 {
			maxEra = 1000000000
		}
		if newValue < 1 || newValue > maxEra {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for YearOfEra: %d", newValue))
		}
		if year < 1 {
			return yearOf([]any{1 - newValue})
		}
		return yearOf([]any{newValue})
	case "ERA":
		if newValue < 0 || newValue > 1 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for Era: %d", newValue))
		}
		currentEra := int64(1)
		if year < 1 {
			currentEra = 0
		}
		if currentEra == newValue {
			return self
		}
		return yearOf([]any{1 - year})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func yearPlusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.plus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	res := temporalAmountAddTo([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return yearFrom([]any{resObj})
	}
	return res
}

func yearPlusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.plus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.plus: missing arguments")
	}
	amountToAdd, okVal := getInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.plus: invalid amount")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.plus: unit is null")
	}

	self := params[0].(*object.Object)
	unitObj := params[2].(*object.Object)
	unitName := getTemporalUnitName(unitObj)

	switch strings.ToUpper(unitName) {
	case "YEARS":
		return yearPlusYears([]any{self, amountToAdd})
	case "DECADES":
		return yearPlusYears([]any{self, amountToAdd * 10})
	case "CENTURIES":
		return yearPlusYears([]any{self, amountToAdd * 100})
	case "MILLENNIA":
		return yearPlusYears([]any{self, amountToAdd * 1000})
	case "ERAS":
		cfEra := getChronoField("ERA")
		if cfEra == nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "ERA field not found")
		}
		eraRes := yearGetLong([]any{self, cfEra})
		if err, isErr := eraRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		eraVal, _ := getInt64Val(eraRes)
		return yearWithField([]any{self, cfEra, eraVal + amountToAdd})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func yearPlusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.plusYears: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.plusYears: missing yearsToAdd")
	}
	self := params[0].(*object.Object)
	yearsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.plusYears: invalid yearsToAdd")
	}
	if yearsToAdd == 0 {
		return self
	}
	year, okY := getYearValue(self)
	if !okY {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.plusYears: invalid instance")
	}
	newYear := year + yearsToAdd
	return yearOf([]any{newYear})
}

func yearMinusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.minus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	res := temporalAmountSubtractFrom([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return yearFrom([]any{resObj})
	}
	return res
}

func yearMinusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.minus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.minus: missing arguments")
	}
	amountToSubtract, okVal := getInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.minus: invalid amount")
	}
	return yearPlusUnit([]any{params[0], -amountToSubtract, params[2]})
}

func yearMinusYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.minusYears: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.minusYears: missing yearsToSubtract")
	}
	yearsToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.minusYears: invalid yearsToSubtract")
	}
	return yearPlusYears([]any{params[0], -yearsToSubtract})
}

func yearQuery(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.query: query is null")
	}

	queryObj := params[1].(*object.Object)
	queryName := ""
	if queryObj.FieldTable != nil {
		if nameFld, okName := queryObj.FieldTable["name"]; okName {
			if strObj, okStr := nameFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				queryName = object.GoStringFromStringObject(strObj)
			}
		}
	}

	if queryName == "precision" {
		return ghelpers.CreateTemporalUnit("Years")
	}
	if queryName == "chronology" {
		return GetIsoChronology()
	}
	if queryName == "zoneId" || queryName == "zone" || queryName == "offset" || queryName == "localDate" || queryName == "localTime" {
		return nil
	}

	return executeTemporalQuery(queryObj, params[0].(*object.Object))
}

func yearAdjustInto(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.adjustInto: temporal is null")
	}

	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.adjustInto: invalid instance")
	}

	cfYear := getChronoField("YEAR")
	return temporalWithField([]any{temporal, cfYear, year})
}

func yearUntil(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.until: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.until: unit is null")
	}

	self := params[0].(*object.Object)
	endExclusive := params[1].(*object.Object)
	unitObj := params[2].(*object.Object)

	selfYear, ok := getYearValue(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.until: invalid instance")
	}

	endYearObj := yearFrom([]any{endExclusive})
	if err, isErr := endYearObj.(*ghelpers.GErrBlk); isErr {
		return err
	}
	endYear, okEnd := getYearValue(endYearObj.(*object.Object))
	if !okEnd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.until: unable to obtain year from endExclusive")
	}

	yearsUntil := endYear - selfYear
	unitName := getTemporalUnitName(unitObj)

	switch strings.ToUpper(unitName) {
	case "YEARS":
		return yearsUntil
	case "DECADES":
		return yearsUntil / 10
	case "CENTURIES":
		return yearsUntil / 100
	case "MILLENNIA":
		return yearsUntil / 1000
	case "ERAS":
		cfEra := getChronoField("ERA")
		if cfEra == nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "ERA field not found")
		}
		endEraRes := yearGetLong([]any{endYearObj, cfEra})
		if err, isErr := endEraRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		endEra, _ := getInt64Val(endEraRes)

		selfEraRes := yearGetLong([]any{self, cfEra})
		if err, isErr := selfEraRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		selfEra, _ := getInt64Val(selfEraRes)

		return endEra - selfEra
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func yearFormat(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.format: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.format: formatter is null")
	}
	return yearToString(params)
}

func yearAtDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.atDay: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atDay: missing dayOfYear")
	}
	self := params[0].(*object.Object)
	dayOfYear, okD := getInt64Val(params[1])
	if !okD {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atDay: invalid dayOfYear")
	}
	year, okY := getYearValue(self)
	if !okY {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atDay: invalid instance")
	}

	maxDays := int64(365)
	if isLeapYearVal(year) {
		maxDays = 366
	}
	if dayOfYear < 1 || dayOfYear > maxDays {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for DayOfYear (valid values 1 - %d): %d", maxDays, dayOfYear))
	}

	t := time.Date(int(year), time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, int(dayOfYear-1))
	return createChronoLocalDate(GetIsoChronology(), year, int64(t.Month()), int64(t.Day()))
}

func yearAtMonthInt(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.atMonth: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atMonth: missing month")
	}
	self := params[0].(*object.Object)
	month, okM := getInt64Val(params[1])
	if !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atMonth: invalid month")
	}
	if month < 1 || month > 12 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MonthOfYear: %d", month))
	}
	year, okY := getYearValue(self)
	if !okY {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atMonth: invalid instance")
	}
	return createYearMonth(year, month)
}

func yearAtMonthMonth(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.atMonth: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.atMonth: month is null")
	}
	self := params[0].(*object.Object)
	monthObj := params[1].(*object.Object)
	ord, okOrd := getMonthOrdinal(monthObj)
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atMonth: invalid Month object")
	}
	year, okY := getYearValue(self)
	if !okY {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atMonth: invalid instance")
	}
	return createYearMonth(year, ord+1)
}

func yearAtMonthDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.atMonthDay: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.atMonthDay: monthDay is null")
	}
	self := params[0].(*object.Object)
	year, okY := getYearValue(self)
	if !okY {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.atMonthDay: invalid instance")
	}
	monthDay := params[1].(*object.Object)
	var month, day int64 = 1, 1
	if monthDay.FieldTable != nil {
		if mFld, okM := monthDay.FieldTable["month"]; okM {
			month, _ = getInt64Val(mFld.Fvalue)
		}
		if dFld, okD := monthDay.FieldTable["day"]; okD {
			day, _ = getInt64Val(dFld.Fvalue)
		} else if dFld, okD2 := monthDay.FieldTable["dayOfMonth"]; okD2 {
			day, _ = getInt64Val(dFld.Fvalue)
		}
	}
	if month == 2 && day == 29 && !isLeapYearVal(year) {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid date 'February 29' as '%d' is not a leap year", year))
	}
	return createChronoLocalDate(GetIsoChronology(), year, month, day)
}

func yearCompareTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	selfYear, okS := getYearValue(self)
	otherYear, okO := getYearValue(other)
	if !okS || !okO {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "Year.compareTo: cannot compare")
	}

	if selfYear < otherYear {
		return int64(-1)
	} else if selfYear > otherYear {
		return int64(1)
	}
	return int64(0)
}

func yearIsAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isAfter: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isAfter: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	selfYear, okS := getYearValue(self)
	otherYear, okO := getYearValue(other)
	if !okS || !okO {
		return false
	}
	return selfYear > otherYear
}

func yearIsBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isBefore: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.isBefore: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	selfYear, okS := getYearValue(self)
	otherYear, okO := getYearValue(other)
	if !okS || !okO {
		return false
	}
	return selfYear < otherYear
}

func yearEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil {
		return false
	}
	self, okS := params[0].(*object.Object)
	other, okO := params[1].(*object.Object)
	if !okS || !okO || object.IsNull(self) || object.IsNull(other) {
		return false
	}
	if self == other {
		return true
	}
	selfYear, okSelf := getYearValue(self)
	otherYear, okOther := getYearValue(other)
	if !okSelf || !okOther {
		return false
	}
	return selfYear == otherYear
}

func yearHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return int64(0)
	}
	return year
}

func yearToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.toString: self is null")
	}
	self := params[0].(*object.Object)
	year, ok := getYearValue(self)
	if !ok {
		return object.StringObjectFromGoString("0000")
	}
	return object.StringObjectFromGoString(formatYearString(year))
}

func yearParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.parse: text is null")
	}
	textObj := params[0].(*object.Object)
	if !object.IsStringObject(textObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Year.parse: invalid text")
	}
	text := strings.TrimSpace(object.GoStringFromStringObject(textObj))
	if len(text) == 0 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text cannot be parsed to a Year: empty text")
	}

	val, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' could not be parsed to Year", text))
	}
	if val < yearMinValue || val > yearMaxValue {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Invalid value for Year (valid values %d - %d): %d", yearMinValue, yearMaxValue, val))
	}
	return createYear(val)
}

func yearParseFormatter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.parse: text is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Year.parse: formatter is null")
	}
	return yearParse([]any{params[0]})
}

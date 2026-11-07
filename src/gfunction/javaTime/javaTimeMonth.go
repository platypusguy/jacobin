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

var classNameMonthVal = "java/time/Month"

type monthDefinition struct {
	enumName    string
	displayName string
}

var monthDefinitions = []monthDefinition{
	{"JANUARY", "January"},
	{"FEBRUARY", "February"},
	{"MARCH", "March"},
	{"APRIL", "April"},
	{"MAY", "May"},
	{"JUNE", "June"},
	{"JULY", "July"},
	{"AUGUST", "August"},
	{"SEPTEMBER", "September"},
	{"OCTOBER", "October"},
	{"NOVEMBER", "November"},
	{"DECEMBER", "December"},
}

var monthFullNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

var monthShortNames = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var monthNarrowNames = []string{
	"J",
	"F",
	"M",
	"A",
	"M",
	"J",
	"J",
	"A",
	"S",
	"O",
	"N",
	"D",
}

var monthFirstDayOfYearStandard = []int64{
	1, 32, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335,
}

var monthFirstDayOfYearLeap = []int64{
	1, 32, 61, 92, 122, 153, 183, 214, 245, 275, 306, 336,
}

var monthMutex sync.Mutex
var monthInitialized bool = false
var monthInstances []*object.Object

func ensureMonthInited() {
	monthMutex.Lock()
	defer monthMutex.Unlock()
	if monthInitialized {
		return
	}

	monthInstances = make([]*object.Object, len(monthDefinitions))
	for index, definition := range monthDefinitions {
		monthObj := object.MakeEmptyObjectWithClassName(&classNameMonthVal)
		monthObj.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(definition.enumName),
		}
		monthObj.FieldTable["displayName"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(definition.displayName),
		}
		monthObj.FieldTable["ordinal"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index),
		}
		monthObj.FieldTable["value"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index + 1),
		}
		monthObj.FieldTable["MONTH_OF_YEAR"] = object.Field{
			Ftype:  types.Long,
			Fvalue: int64(index + 1),
		}
		monthObj.FieldTable["precision"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: ghelpers.CreateTemporalUnit("Months"),
		}

		monthInstances[index] = monthObj
		_ = statics.AddStatic(classNameMonthVal+"."+definition.enumName, statics.Static{
			Type:  "Ljava/time/Month;",
			Value: monthObj,
		})
	}

	monthInitialized = true
}

func getMonthOrdinal(targetObject *object.Object) (int64, bool) {
	if targetObject == nil || targetObject.FieldTable == nil {
		return 0, false
	}
	if ordinalField, exists := targetObject.FieldTable["ordinal"]; exists {
		if ordinalValue, ok := getInt64Val(ordinalField.Fvalue); ok {
			return ordinalValue, true
		}
	}
	if nameField, exists := targetObject.FieldTable["name"]; exists {
		if nameStringObject, ok := nameField.Fvalue.(*object.Object); ok && object.IsStringObject(nameStringObject) {
			goName := object.GoStringFromStringObject(nameStringObject)
			for index, definition := range monthDefinitions {
				if definition.enumName == goName || definition.displayName == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

func getChronoFieldMonthOfYear() *object.Object {
	ensureChronoFieldInited()
	for index, definition := range chronoFieldDefinitions {
		if definition.enumName == "MONTH_OF_YEAR" {
			return chronoFieldInstances[index]
		}
	}
	return nil
}

func parseBoolParam(val any) (bool, bool) {
	switch v := val.(type) {
	case bool:
		return v, true
	case int64:
		return v != 0, true
	case int32:
		return v != 0, true
	case int:
		return v != 0, true
	case uint64:
		return v != 0, true
	case uint32:
		return v != 0, true
	case uint8:
		return v != 0, true
	case int8:
		return v != 0, true
	default:
		return false, false
	}
}

// Load_Time_Month registers all java/time/Month methods.
func Load_Time_Month() {
	ghelpers.MethodSignatures["java/time/Month.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthClinit,
		}

	ghelpers.MethodSignatures["java/time/Month.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  monthInit,
		}

	ghelpers.MethodSignatures["java/time/Month.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/Month.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Month.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Month.compareTo(Ljava/time/Month;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Month.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthEquals,
		}

	ghelpers.MethodSignatures["java/time/Month.firstDayOfYear(Z)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthFirstDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/Month.firstMonthOfQuarter()Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthFirstMonthOfQuarter,
		}

	ghelpers.MethodSignatures["java/time/Month.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthFrom,
		}

	ghelpers.MethodSignatures["java/time/Month.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthGet,
		}

	ghelpers.MethodSignatures["java/time/Month.getDisplayName(Ljava/time/format/TextStyle;Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  monthGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/Month.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthGetLong,
		}

	ghelpers.MethodSignatures["java/time/Month.getValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthGetValue,
		}

	ghelpers.MethodSignatures["java/time/Month.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthHashCode,
		}

	ghelpers.MethodSignatures["java/time/Month.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthIsSupported,
		}

	ghelpers.MethodSignatures["java/time/Month.length(Z)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthLength,
		}

	ghelpers.MethodSignatures["java/time/Month.maxLength()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthMaxLength,
		}

	ghelpers.MethodSignatures["java/time/Month.minLength()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthMinLength,
		}

	ghelpers.MethodSignatures["java/time/Month.minus(J)Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthMinus,
		}

	ghelpers.MethodSignatures["java/time/Month.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthName,
		}

	ghelpers.MethodSignatures["java/time/Month.of(I)Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthOf,
		}

	ghelpers.MethodSignatures["java/time/Month.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthOrdinal,
		}

	ghelpers.MethodSignatures["java/time/Month.plus(J)Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthPlus,
		}

	ghelpers.MethodSignatures["java/time/Month.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthQuery,
		}

	ghelpers.MethodSignatures["java/time/Month.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthRange,
		}

	ghelpers.MethodSignatures["java/time/Month.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthToString,
		}

	ghelpers.MethodSignatures["java/time/Month.valueOf(Ljava/lang/String;)Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  monthValueOf,
		}

	ghelpers.MethodSignatures["java/time/Month.values()[Ljava/time/Month;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  monthValues,
		}
}

func monthClinit([]any) any {
	ensureMonthInited()
	return nil
}

func monthInit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.<init>: self is null")
	}
	self := params[0].(*object.Object)
	nameObj, okName := params[1].(*object.Object)
	if !okName || object.IsNull(nameObj) || !object.IsStringObject(nameObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.<init>: invalid name argument")
	}
	ordinal, okOrdinal := getInt64Val(params[2])
	if !okOrdinal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.<init>: invalid ordinal argument")
	}

	self.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: nameObj}
	self.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: ordinal}
	self.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: ordinal + 1}
	self.FieldTable["MONTH_OF_YEAR"] = object.Field{Ftype: types.Long, Fvalue: ordinal + 1}
	self.FieldTable["precision"] = object.Field{Ftype: types.Ref, Fvalue: ghelpers.CreateTemporalUnit("Months")}
	return nil
}

func monthOf(params []any) any {
	ensureMonthInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.of: missing argument")
	}
	monthVal, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.of: invalid month argument")
	}
	if monthVal < 1 || monthVal > 12 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MonthOfYear: %d", monthVal))
	}
	return monthInstances[monthVal-1]
}

func monthFrom(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.from: temporal is null")
	}
	temporal := params[0].(*object.Object)

	if object.GoStringFromStringPoolIndex(temporal.KlassName) == classNameMonthVal {
		if ord, ok := getMonthOrdinal(temporal); ok && ord >= 0 && ord < int64(len(monthInstances)) {
			return monthInstances[ord]
		}
		return temporal
	}

	// Try getting ChronoField MONTH_OF_YEAR
	cfMonthOfYear := getChronoFieldMonthOfYear()
	if cfMonthOfYear != nil {
		res := temporalAccessorGet([]any{temporal, cfMonthOfYear})
		if _, isErr := res.(*ghelpers.GErrBlk); !isErr {
			if val, ok := getInt64Val(res); ok && val >= 1 && val <= 12 {
				return monthInstances[val-1]
			}
		}
	}

	// Fallback: check if temporal has date fields (e.g. month, epochDay)
	if temporal.FieldTable != nil {
		if monthFld, okM := temporal.FieldTable["month"]; okM {
			if monthVal, ok := getInt64Val(monthFld.Fvalue); ok && monthVal >= 1 && monthVal <= 12 {
				return monthInstances[monthVal-1]
			}
		}
		if epochFld, okE := temporal.FieldTable["epochDay"]; okE {
			if epochDay, ok := getInt64Val(epochFld.Fvalue); ok {
				t := time.Unix(epochDay*86400, 0).UTC()
				m := int64(t.Month())
				if m >= 1 && m <= 12 {
					return monthInstances[m-1]
				}
			}
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain Month from TemporalAccessor")
}

func monthGetValue(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.getValue: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getMonthOrdinal(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.getValue: invalid instance")
	}
	return ordinal + 1
}

func monthGetDisplayName(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.getDisplayName: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.getDisplayName: style is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.getDisplayName: locale is null")
	}

	self := params[0].(*object.Object)
	ordinal, ok := getMonthOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(monthDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.getDisplayName: invalid instance")
	}

	styleObj := params[1].(*object.Object)
	styleName := ""
	if styleObj.FieldTable != nil {
		if nameFld, okName := styleObj.FieldTable["name"]; okName {
			if strObj, okStr := nameFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				styleName = strings.ToUpper(object.GoStringFromStringObject(strObj))
			}
		}
	}

	if strings.Contains(styleName, "SHORT") {
		return object.StringObjectFromGoString(monthShortNames[ordinal])
	} else if strings.Contains(styleName, "NARROW") {
		return object.StringObjectFromGoString(monthNarrowNames[ordinal])
	}
	return object.StringObjectFromGoString(monthFullNames[ordinal])
}

func monthIsSupported(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.isSupported: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "MONTH_OF_YEAR" {
		return types.JavaBoolTrue
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return types.JavaBoolFalse
	}
	return temporalFieldIsSupportedBy([]any{field, params[0]})
}

func monthRange(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.range: field is null")
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "MONTH_OF_YEAR" {
		return createValueRange(1, 1, 12, 12)
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorRange(params)
}

func monthGet(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.get: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "MONTH_OF_YEAR" {
		ordinal, ok := getMonthOrdinal(self)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.get: invalid instance")
		}
		return ordinal + 1
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorGet(params)
}

func monthGetLong(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.getLong: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "MONTH_OF_YEAR" {
		ordinal, ok := getMonthOrdinal(self)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.getLong: invalid instance")
		}
		return ordinal + 1
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorGetLong(params)
}

func monthPlus(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.plus: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.plus: missing months argument")
	}
	self := params[0].(*object.Object)
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.plus: invalid months argument")
	}
	ordinal, okOrd := getMonthOrdinal(self)
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.plus: invalid Month instance")
	}

	amount := months % 12
	newOrdinal := (ordinal + amount + 12) % 12
	return monthInstances[newOrdinal]
}

func monthMinus(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.minus: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.minus: missing months argument")
	}
	months, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.minus: invalid months argument")
	}
	return monthPlus([]any{params[0], -(months % 12)})
}

func monthLength(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.length: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.length: missing leapYear argument")
	}
	leapYear, ok := parseBoolParam(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.length: invalid leapYear argument")
	}
	self := params[0].(*object.Object)
	ordinal, okOrd := getMonthOrdinal(self)
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.length: invalid Month instance")
	}

	switch ordinal {
	case 1: // FEBRUARY
		if leapYear {
			return int64(29)
		}
		return int64(28)
	case 3, 5, 8, 10: // APRIL, JUNE, SEPTEMBER, NOVEMBER
		return int64(30)
	default:
		return int64(31)
	}
}

func monthMinLength(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.minLength: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, okOrd := getMonthOrdinal(self)
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.minLength: invalid Month instance")
	}

	switch ordinal {
	case 1: // FEBRUARY
		return int64(28)
	case 3, 5, 8, 10: // APRIL, JUNE, SEPTEMBER, NOVEMBER
		return int64(30)
	default:
		return int64(31)
	}
}

func monthMaxLength(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.maxLength: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, okOrd := getMonthOrdinal(self)
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.maxLength: invalid Month instance")
	}

	switch ordinal {
	case 1: // FEBRUARY
		return int64(29)
	case 3, 5, 8, 10: // APRIL, JUNE, SEPTEMBER, NOVEMBER
		return int64(30)
	default:
		return int64(31)
	}
}

func monthFirstDayOfYear(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.firstDayOfYear: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.firstDayOfYear: missing leapYear argument")
	}
	leapYear, ok := parseBoolParam(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.firstDayOfYear: invalid leapYear argument")
	}
	self := params[0].(*object.Object)
	ordinal, okOrd := getMonthOrdinal(self)
	if !okOrd || ordinal < 0 || ordinal >= int64(len(monthDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.firstDayOfYear: invalid Month instance")
	}

	if leapYear {
		return monthFirstDayOfYearLeap[ordinal]
	}
	return monthFirstDayOfYearStandard[ordinal]
}

func monthFirstMonthOfQuarter(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.firstMonthOfQuarter: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, okOrd := getMonthOrdinal(self)
	if !okOrd || ordinal < 0 || ordinal >= int64(len(monthDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.firstMonthOfQuarter: invalid Month instance")
	}

	firstMonthOrd := (ordinal / 3) * 3
	return monthInstances[firstMonthOrd]
}

func monthQuery(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.query: query is null")
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
		return ghelpers.CreateTemporalUnit("Months")
	}
	if queryName == "chronology" {
		return getIsoChronologySingleton()
	}
	if queryName == "zoneId" || queryName == "zone" || queryName == "offset" || queryName == "localDate" || queryName == "localTime" {
		return nil
	}
	return executeTemporalQuery(queryObj, params[0].(*object.Object))
}

func monthAdjustInto(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	ordinal, ok := getMonthOrdinal(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.adjustInto: invalid Month instance")
	}
	val := ordinal + 1
	cfMonthOfYear := getChronoFieldMonthOfYear()
	return temporalWithField([]any{temporal, cfMonthOfYear, val})
}

func monthValues([]any) any {
	ensureMonthInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/Month;", int64(len(monthInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, monthInstances)
	return arrayObject
}

func monthValueOf(params []any) any {
	ensureMonthInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.valueOf: missing argument")
	}
	stringObject, okObj := params[0].(*object.Object)
	if !okObj {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.valueOf: argument is not a String")
	}
	if stringObject == nil || object.IsNull(stringObject) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.valueOf: name is null")
	}
	if !object.IsStringObject(stringObject) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Month.valueOf: argument is not a String")
	}
	targetName := object.GoStringFromStringObject(stringObject)
	for index, definition := range monthDefinitions {
		if definition.enumName == targetName {
			return monthInstances[index]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "No enum constant java.time.Month."+targetName)
}

func monthName(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.name: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if fieldEntry, exists := selfObject.FieldTable["name"]; exists {
			if stringObject, okStr := fieldEntry.Fvalue.(*object.Object); okStr && object.IsStringObject(stringObject) {
				return stringObject
			}
		}
	}
	return object.StringObjectFromGoString("")
}

func monthOrdinal(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.ordinal: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getMonthOrdinal(selfObject)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func monthToString(params []any) any {
	return monthName(params)
}

func monthCompareTo(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.compareTo: other is null")
	}
	selfObject := params[0].(*object.Object)
	otherObject := params[1].(*object.Object)

	selfOrdinal, ok1 := getMonthOrdinal(selfObject)
	otherOrdinal, ok2 := getMonthOrdinal(otherObject)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "Month.compareTo: object is not a Month")
	}

	if selfOrdinal < otherOrdinal {
		return int64(-1)
	} else if selfOrdinal > otherOrdinal {
		return int64(1)
	}
	return int64(0)
}

func monthEquals(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	selfObject := params[0].(*object.Object)
	otherObject := params[1].(*object.Object)

	if selfObject == otherObject {
		return types.JavaBoolTrue
	}

	selfOrdinal, ok1 := getMonthOrdinal(selfObject)
	otherOrdinal, ok2 := getMonthOrdinal(otherObject)
	if ok1 && ok2 && selfOrdinal == otherOrdinal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func monthHashCode(params []any) any {
	ensureMonthInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Month.hashCode: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinal, ok := getMonthOrdinal(selfObject)
	if !ok {
		return int64(0)
	}
	return ordinal
}

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

var classNameDayOfWeekVal = "java/time/DayOfWeek"

type dayOfWeekDefinition struct {
	enumName    string
	displayName string
}

var dayOfWeekDefinitions = []dayOfWeekDefinition{
	{"MONDAY", "Monday"},
	{"TUESDAY", "Tuesday"},
	{"WEDNESDAY", "Wednesday"},
	{"THURSDAY", "Thursday"},
	{"FRIDAY", "Friday"},
	{"SATURDAY", "Saturday"},
	{"SUNDAY", "Sunday"},
}

var dayOfWeekFullNames = []string{
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
}

var dayOfWeekShortNames = []string{
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
	"Sun",
}

var dayOfWeekNarrowNames = []string{
	"M",
	"T",
	"W",
	"T",
	"F",
	"S",
	"S",
}

var dayOfWeekMutex sync.Mutex
var dayOfWeekInitialized bool = false
var dayOfWeekInstances []*object.Object

func ensureDayOfWeekInited() {
	dayOfWeekMutex.Lock()
	defer dayOfWeekMutex.Unlock()
	if dayOfWeekInitialized {
		return
	}

	dayOfWeekInstances = make([]*object.Object, len(dayOfWeekDefinitions))
	for index, definition := range dayOfWeekDefinitions {
		dayObj := object.MakeEmptyObjectWithClassName(&classNameDayOfWeekVal)
		dayObj.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(definition.enumName),
		}
		dayObj.FieldTable["displayName"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(definition.displayName),
		}
		dayObj.FieldTable["ordinal"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index),
		}
		dayObj.FieldTable["value"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index + 1),
		}
		dayObj.FieldTable["DAY_OF_WEEK"] = object.Field{
			Ftype:  types.Long,
			Fvalue: int64(index + 1),
		}
		dayObj.FieldTable["precision"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: ghelpers.CreateTemporalUnit("Days"),
		}

		dayOfWeekInstances[index] = dayObj
		_ = statics.AddStatic(classNameDayOfWeekVal+"."+definition.enumName, statics.Static{
			Type:  "Ljava/time/DayOfWeek;",
			Value: dayObj,
		})
	}

	dayOfWeekInitialized = true
}

func getDayOfWeekOrdinal(targetObject *object.Object) (int64, bool) {
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
			for index, definition := range dayOfWeekDefinitions {
				if definition.enumName == goName || definition.displayName == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

func getChronoFieldDayOfWeek() *object.Object {
	ensureChronoFieldInited()
	for index, definition := range chronoFieldDefinitions {
		if definition.enumName == "DAY_OF_WEEK" {
			return chronoFieldInstances[index]
		}
	}
	return nil
}

// Load_Time_DayOfWeek registers all java/time/DayOfWeek methods.
func Load_Time_DayOfWeek() {
	dayOfWeekMutex.Lock()
	dayOfWeekInitialized = false
	dayOfWeekMutex.Unlock()
	ensureDayOfWeekInited()

	ghelpers.MethodSignatures["java/time/DayOfWeek.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekClinit,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dayOfWeekInit,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekCompareTo,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekCompareTo,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.compareTo(Ljava/time/DayOfWeek;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekCompareTo,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekEquals,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekFrom,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekGet,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.getDisplayName(Ljava/time/format/TextStyle;Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dayOfWeekGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekGetLong,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.getValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekGetValue,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekHashCode,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekIsSupported,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.minus(J)Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekMinus,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekName,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.of(I)Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekOf,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekOrdinal,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.plus(J)Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekPlus,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekQuery,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekRange,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekToString,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.valueOf(Ljava/lang/String;)Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dayOfWeekValueOf,
		}

	ghelpers.MethodSignatures["java/time/DayOfWeek.values()[Ljava/time/DayOfWeek;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dayOfWeekValues,
		}
}

func dayOfWeekClinit([]any) any {
	dayOfWeekMutex.Lock()
	dayOfWeekInitialized = false
	dayOfWeekMutex.Unlock()
	ensureDayOfWeekInited()
	return nil
}

func dayOfWeekInit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.<init>: self is null")
	}
	self := params[0].(*object.Object)
	nameObj, okName := params[1].(*object.Object)
	if !okName || object.IsNull(nameObj) || !object.IsStringObject(nameObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.<init>: invalid name argument")
	}
	ordinal, okOrdinal := getInt64Val(params[2])
	if !okOrdinal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.<init>: invalid ordinal argument")
	}

	self.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: nameObj}
	self.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: ordinal}
	self.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: ordinal + 1}
	self.FieldTable["DAY_OF_WEEK"] = object.Field{Ftype: types.Long, Fvalue: ordinal + 1}
	self.FieldTable["precision"] = object.Field{Ftype: types.Ref, Fvalue: ghelpers.CreateTemporalUnit("Days")}
	return nil
}

func dayOfWeekOf(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.of: missing argument")
	}
	dayVal, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.of: invalid day argument")
	}
	if dayVal < 1 || dayVal > 7 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for DayOfWeek: %d", dayVal))
	}
	return dayOfWeekInstances[dayVal-1]
}

func dayOfWeekFrom(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.from: temporal is null")
	}
	temporal := params[0].(*object.Object)

	if object.GoStringFromStringPoolIndex(temporal.KlassName) == classNameDayOfWeekVal {
		if ord, ok := getDayOfWeekOrdinal(temporal); ok && ord >= 0 && ord < int64(len(dayOfWeekInstances)) {
			return dayOfWeekInstances[ord]
		}
		return temporal
	}

	// Try getting ChronoField DAY_OF_WEEK
	cfDayOfWeek := getChronoFieldDayOfWeek()
	if cfDayOfWeek != nil {
		res := temporalAccessorGet([]any{temporal, cfDayOfWeek})
		if _, isErr := res.(*ghelpers.GErrBlk); !isErr {
			if val, ok := getInt64Val(res); ok && val >= 1 && val <= 7 {
				return dayOfWeekInstances[val-1]
			}
		}
	}

	// Fallback: check if temporal has date fields (e.g. year, month, day)
	if temporal.FieldTable != nil {
		if yearFld, okY := temporal.FieldTable["year"]; okY {
			if monthFld, okM := temporal.FieldTable["month"]; okM {
				if dayFld, okD := temporal.FieldTable["day"]; okD {
					year, ok1 := getInt64Val(yearFld.Fvalue)
					month, ok2 := getInt64Val(monthFld.Fvalue)
					day, ok3 := getInt64Val(dayFld.Fvalue)
					if ok1 && ok2 && ok3 && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
						t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
						weekday := t.Weekday()
						var dowVal int64
						if weekday == time.Sunday {
							dowVal = 7
						} else {
							dowVal = int64(weekday)
						}
						return dayOfWeekInstances[dowVal-1]
					}
				}
			}
		}
		if epochFld, okE := temporal.FieldTable["epochDay"]; okE {
			if epochDay, ok := getInt64Val(epochFld.Fvalue); ok {
				t := time.Unix(epochDay*86400, 0).UTC()
				weekday := t.Weekday()
				var dowVal int64
				if weekday == time.Sunday {
					dowVal = 7
				} else {
					dowVal = int64(weekday)
				}
				return dayOfWeekInstances[dowVal-1]
			}
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain DayOfWeek from TemporalAccessor")
}

func dayOfWeekGetValue(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.getValue: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getDayOfWeekOrdinal(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.getValue: invalid instance")
	}
	return ordinal + 1
}

func dayOfWeekGetDisplayName(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.getDisplayName: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.getDisplayName: style is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.getDisplayName: locale is null")
	}

	self := params[0].(*object.Object)
	ordinal, ok := getDayOfWeekOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(dayOfWeekDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.getDisplayName: invalid instance")
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
		return object.StringObjectFromGoString(dayOfWeekShortNames[ordinal])
	} else if strings.Contains(styleName, "NARROW") {
		return object.StringObjectFromGoString(dayOfWeekNarrowNames[ordinal])
	}
	return object.StringObjectFromGoString(dayOfWeekFullNames[ordinal])
}

func dayOfWeekIsSupported(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.isSupported: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "DAY_OF_WEEK" {
		return types.JavaBoolTrue
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return types.JavaBoolFalse
	}
	return temporalFieldIsSupportedBy([]any{field, params[0]})
}

func dayOfWeekRange(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.range: field is null")
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "DAY_OF_WEEK" {
		return createValueRange(1, 1, 7, 7)
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorRange(params)
}

func dayOfWeekGet(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.get: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "DAY_OF_WEEK" {
		ordinal, ok := getDayOfWeekOrdinal(self)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.get: invalid instance")
		}
		return ordinal + 1
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorGet(params)
}

func dayOfWeekGetLong(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.getLong: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if fieldName == "DAY_OF_WEEK" {
		ordinal, ok := getDayOfWeekOrdinal(self)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.getLong: invalid instance")
		}
		return ordinal + 1
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorGetLong(params)
}

func dayOfWeekPlus(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.plus: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.plus: missing days argument")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.plus: invalid days argument")
	}
	self := params[0].(*object.Object)
	ordinal, okOrd := getDayOfWeekOrdinal(self)
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.plus: invalid instance")
	}
	amount := days % 7
	newOrdinal := (ordinal + (amount + 7)) % 7
	return dayOfWeekInstances[newOrdinal]
}

func dayOfWeekMinus(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.minus: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.minus: missing days argument")
	}
	days, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.minus: invalid days argument")
	}
	return dayOfWeekPlus([]any{params[0], -(days % 7)})
}

func dayOfWeekQuery(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.query: query is null")
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
		return ghelpers.CreateTemporalUnit("Days")
	}
	if queryName == "zoneId" || queryName == "chronology" || queryName == "zone" || queryName == "offset" || queryName == "localDate" || queryName == "localTime" {
		return nil
	}
	return executeTemporalQuery(queryObj, params[0].(*object.Object))
}

func dayOfWeekAdjustInto(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	ordinal, ok := getDayOfWeekOrdinal(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.adjustInto: invalid DayOfWeek instance")
	}
	val := ordinal + 1
	cfDayOfWeek := getChronoFieldDayOfWeek()
	return temporalWithField([]any{temporal, cfDayOfWeek, val})
}

func dayOfWeekValues([]any) any {
	ensureDayOfWeekInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/DayOfWeek;", int64(len(dayOfWeekInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, dayOfWeekInstances)
	return arrayObject
}

func dayOfWeekValueOf(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.valueOf: missing argument")
	}
	stringObject, okObj := params[0].(*object.Object)
	if !okObj {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.valueOf: argument is not a String")
	}
	if stringObject == nil || object.IsNull(stringObject) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.valueOf: name is null")
	}
	if !object.IsStringObject(stringObject) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DayOfWeek.valueOf: argument is not a String")
	}
	targetName := object.GoStringFromStringObject(stringObject)
	for index, definition := range dayOfWeekDefinitions {
		if definition.enumName == targetName {
			return dayOfWeekInstances[index]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "No enum constant java.time.DayOfWeek."+targetName)
}

func dayOfWeekName(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.name: self is null")
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

func dayOfWeekOrdinal(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.ordinal: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getDayOfWeekOrdinal(selfObject)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func dayOfWeekToString(params []any) any {
	ensureDayOfWeekInited()
	return dayOfWeekName(params)
}

func dayOfWeekCompareTo(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.compareTo: other is null")
	}
	selfObject := params[0].(*object.Object)
	otherObject := params[1].(*object.Object)

	selfOrdinal, okSelf := getDayOfWeekOrdinal(selfObject)
	otherOrdinal, okOther := getDayOfWeekOrdinal(otherObject)
	if !okSelf || !okOther {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "DayOfWeek.compareTo: invalid enum type")
	}
	if selfOrdinal < otherOrdinal {
		return int64(-1)
	} else if selfOrdinal > otherOrdinal {
		return int64(1)
	}
	return int64(0)
}

func dayOfWeekEquals(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	selfObject := params[0].(*object.Object)
	otherObject := params[1].(*object.Object)

	if selfObject == otherObject {
		return types.JavaBoolTrue
	}

	selfOrdinal, okSelf := getDayOfWeekOrdinal(selfObject)
	otherOrdinal, okOther := getDayOfWeekOrdinal(otherObject)
	if okSelf && okOther && selfOrdinal == otherOrdinal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func dayOfWeekHashCode(params []any) any {
	ensureDayOfWeekInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DayOfWeek.hashCode: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinal, ok := getDayOfWeekOrdinal(selfObject)
	if ok {
		return ordinal
	}
	return int64(0)
}

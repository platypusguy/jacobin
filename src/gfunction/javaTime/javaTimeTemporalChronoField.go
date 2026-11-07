/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"math"
	"sync"
)

var classNameChronoField = "java/time/temporal/ChronoField"

type chronoFieldDefinition struct {
	enumName    string
	displayName string
	baseUnit    string
	rangeUnit   string
	minSmallest int64
	minLargest  int64
	maxSmallest int64
	maxLargest  int64
}

var chronoFieldDefinitions = []chronoFieldDefinition{
	{"NANO_OF_SECOND", "NanoOfSecond", "Nanos", "Seconds", 0, 0, 999_999_999, 999_999_999},
	{"NANO_OF_DAY", "NanoOfDay", "Nanos", "Days", 0, 0, 86_399_999_999_999, 86_399_999_999_999},
	{"MICRO_OF_SECOND", "MicroOfSecond", "Micros", "Seconds", 0, 0, 999_999, 999_999},
	{"MICRO_OF_DAY", "MicroOfDay", "Micros", "Days", 0, 0, 86_399_999_999, 86_399_999_999},
	{"MILLI_OF_SECOND", "MilliOfSecond", "Millis", "Seconds", 0, 0, 999, 999},
	{"MILLI_OF_DAY", "MilliOfDay", "Millis", "Days", 0, 0, 86_399_999, 86_399_999},
	{"SECOND_OF_MINUTE", "SecondOfMinute", "Seconds", "Minutes", 0, 0, 59, 59},
	{"SECOND_OF_DAY", "SecondOfDay", "Seconds", "Days", 0, 0, 86_399, 86_399},
	{"MINUTE_OF_HOUR", "MinuteOfHour", "Minutes", "Hours", 0, 0, 59, 59},
	{"MINUTE_OF_DAY", "MinuteOfDay", "Minutes", "Days", 0, 0, 1439, 1439},
	{"HOUR_OF_AMPM", "HourOfAmPm", "Hours", "HalfDays", 0, 0, 11, 11},
	{"CLOCK_HOUR_OF_AMPM", "ClockHourOfAmPm", "Hours", "HalfDays", 1, 1, 12, 12},
	{"HOUR_OF_DAY", "HourOfDay", "Hours", "Days", 0, 0, 23, 23},
	{"CLOCK_HOUR_OF_DAY", "ClockHourOfDay", "Hours", "Days", 1, 1, 24, 24},
	{"AMPM_OF_DAY", "AmPmOfDay", "HalfDays", "Days", 0, 0, 1, 1},
	{"DAY_OF_WEEK", "DayOfWeek", "Days", "Weeks", 1, 1, 7, 7},
	{"ALIGNED_DAY_OF_WEEK_IN_MONTH", "AlignedDayOfWeekInMonth", "Days", "Weeks", 1, 1, 7, 7},
	{"ALIGNED_DAY_OF_WEEK_IN_YEAR", "AlignedDayOfWeekInYear", "Days", "Weeks", 1, 1, 7, 7},
	{"DAY_OF_MONTH", "DayOfMonth", "Days", "Months", 1, 1, 28, 31},
	{"DAY_OF_YEAR", "DayOfYear", "Days", "Years", 1, 1, 365, 366},
	{"EPOCH_DAY", "EpochDay", "Days", "Forever", -365243219162, -365243219162, 365241780471, 365241780471},
	{"ALIGNED_WEEK_OF_MONTH", "AlignedWeekOfMonth", "Weeks", "Months", 1, 1, 4, 5},
	{"ALIGNED_WEEK_OF_YEAR", "AlignedWeekOfYear", "Weeks", "Years", 1, 1, 53, 53},
	{"MONTH_OF_YEAR", "MonthOfYear", "Months", "Years", 1, 1, 12, 12},
	{"PROLEPTIC_MONTH", "ProlepticMonth", "Months", "Forever", -11999999988, -11999999988, 11999999999, 11999999999},
	{"YEAR_OF_ERA", "YearOfEra", "Years", "Forever", 1, 1, 999999999, 1000000000},
	{"YEAR", "Year", "Years", "Forever", -999999999, -999999999, 999999999, 999999999},
	{"ERA", "Era", "Eras", "Forever", 0, 0, 1, 1},
	{"INSTANT_SECONDS", "InstantSeconds", "Seconds", "Forever", math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64},
	{"OFFSET_SECONDS", "OffsetSeconds", "Seconds", "Forever", -64800, -64800, 64800, 64800},
}

var chronoFieldMutex sync.Mutex
var chronoFieldInitialized bool = false
var chronoFieldInstances []*object.Object

func ensureChronoFieldInited() {
	chronoFieldMutex.Lock()
	defer chronoFieldMutex.Unlock()
	if chronoFieldInitialized {
		return
	}

	chronoFieldInstances = make([]*object.Object, len(chronoFieldDefinitions))
	for index, definition := range chronoFieldDefinitions {
		fieldObject := object.MakeEmptyObjectWithClassName(&classNameChronoField)
		fieldObject.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(definition.enumName)}
		fieldObject.FieldTable["displayName"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(definition.displayName)}
		fieldObject.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(index)}
		fieldObject.FieldTable["baseUnit"] = object.Field{Ftype: types.Ref, Fvalue: ghelpers.CreateTemporalUnit(definition.baseUnit)}
		fieldObject.FieldTable["rangeUnit"] = object.Field{Ftype: types.Ref, Fvalue: ghelpers.CreateTemporalUnit(definition.rangeUnit)}
		fieldObject.FieldTable["range"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: createValueRange(definition.minSmallest, definition.minLargest, definition.maxSmallest, definition.maxLargest),
		}
		chronoFieldInstances[index] = fieldObject
		_ = statics.AddStatic(classNameChronoField+"."+definition.enumName, statics.Static{
			Type:  "Ljava/time/temporal/ChronoField;",
			Value: fieldObject,
		})
	}
	chronoFieldInitialized = true
}

func Load_Time_Temporal_ChronoField() {
	ghelpers.MethodSignatures["java/time/temporal/ChronoField.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldClinit,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.adjustInto(Ljava/time/temporal/Temporal;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronoFieldAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.checkValidIntValue(J)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldCheckValidIntValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.checkValidValue(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldCheckValidValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldCompareTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldCompareTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldEquals,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.getBaseUnit()Ljava/time/temporal/TemporalUnit;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldGetBaseUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.getDisplayName(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.getFrom(Ljava/time/temporal/TemporalAccessor;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldGetFrom,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.getRangeUnit()Ljava/time/temporal/TemporalUnit;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldGetRangeUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldHashCode,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.isDateBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldIsDateBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.isSupportedBy(Ljava/time/temporal/TemporalAccessor;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldIsSupportedBy,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.isTimeBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldIsTimeBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldName,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldOrdinal,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.range()Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldRange,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.rangeRefinedBy(Ljava/time/temporal/TemporalAccessor;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldRangeRefinedBy,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.resolve(Ljava/util/Map;Ljava/time/temporal/TemporalAccessor;Ljava/time/format/ResolverStyle;)Ljava/time/temporal/TemporalAccessor;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronoFieldResolve,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldToString,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.valueOf(Ljava/lang/String;)Ljava/time/temporal/ChronoField;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoFieldValueOf,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoField.values()[Ljava/time/temporal/ChronoField;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoFieldValues,
		}
}

func chronoFieldClinit([]any) any {
	ensureChronoFieldInited()
	return nil
}

func chronoFieldAdjustInto(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.adjustInto: temporal is null")
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)
	return temporalWithField([]any{temporalObject, selfObject, params[2]})
}

func chronoFieldCheckValidIntValue(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.checkValidIntValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	valueObject := params[1]
	rangeObject := chronoFieldRange([]any{selfObject})
	if errBlock, isErr := rangeObject.(*ghelpers.GErrBlk); isErr {
		return errBlock
	}
	return valueRangeCheckValidIntValue([]any{rangeObject, valueObject, selfObject})
}

func chronoFieldCheckValidValue(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.checkValidValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	valueObject := params[1]
	rangeObject := chronoFieldRange([]any{selfObject})
	if errBlock, isErr := rangeObject.(*ghelpers.GErrBlk); isErr {
		return errBlock
	}
	return valueRangeCheckValidValue([]any{rangeObject, valueObject, selfObject})
}

func chronoFieldCompareTo(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.compareTo: other is null")
	}
	selfObject := params[0].(*object.Object)
	otherObject := params[1].(*object.Object)

	selfOrdinal, okSelf := getChronoFieldOrdinal(selfObject)
	otherOrdinal, okOther := getChronoFieldOrdinal(otherObject)
	if !okSelf || !okOther {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "ChronoField.compareTo: invalid enum type")
	}
	if selfOrdinal < otherOrdinal {
		return int64(-1)
	} else if selfOrdinal > otherOrdinal {
		return int64(1)
	}
	return int64(0)
}

func chronoFieldEquals(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if params[1] == nil {
		return types.JavaBoolFalse
	}
	otherObject, okObj := params[1].(*object.Object)
	if !okObj || object.IsNull(otherObject) {
		return types.JavaBoolFalse
	}
	if params[0].(*object.Object) == otherObject {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func chronoFieldGetBaseUnit(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.getBaseUnit: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if fieldEntry, exists := selfObject.FieldTable["baseUnit"]; exists {
			if unitObject, okUnit := fieldEntry.Fvalue.(*object.Object); okUnit && !object.IsNull(unitObject) {
				return unitObject
			}
		}
	}
	return temporalFieldGetBaseUnit(params)
}

func chronoFieldGetDisplayName(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.getDisplayName: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.getDisplayName: locale is null")
	}
	return chronoFieldToString([]any{params[0]})
}

func chronoFieldGetFrom(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.getFrom: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.getFrom: temporal is null")
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)
	return temporalAccessorGetLong([]any{temporalObject, selfObject})
}

func chronoFieldGetRangeUnit(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.getRangeUnit: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if fieldEntry, exists := selfObject.FieldTable["rangeUnit"]; exists {
			if unitObject, okUnit := fieldEntry.Fvalue.(*object.Object); okUnit && !object.IsNull(unitObject) {
				return unitObject
			}
		}
	}
	return temporalFieldGetRangeUnit(params)
}

func chronoFieldHashCode(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return int64(0)
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoFieldOrdinal(selfObject)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func chronoFieldIsDateBased(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.isDateBased: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoFieldOrdinal(selfObject)
	if okOrdinal {
		// DAY_OF_WEEK (15) to ERA (27)
		if ordinalValue >= 15 && ordinalValue <= 27 {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	return temporalFieldIsDateBased(params)
}

func chronoFieldIsSupportedBy(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.isSupportedBy: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)
	return temporalAccessorIsSupported([]any{temporalObject, selfObject})
}

func chronoFieldIsTimeBased(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.isTimeBased: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoFieldOrdinal(selfObject)
	if okOrdinal {
		// NANO_OF_SECOND (0) to AMPM_OF_DAY (14)
		if ordinalValue >= 0 && ordinalValue <= 14 {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	return temporalFieldIsTimeBased(params)
}

func chronoFieldName(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.name: self is null")
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

func chronoFieldOrdinal(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.ordinal: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoFieldOrdinal(selfObject)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func chronoFieldRange(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.range: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if fieldEntry, exists := selfObject.FieldTable["range"]; exists {
			if rangeObject, okRange := fieldEntry.Fvalue.(*object.Object); okRange && !object.IsNull(rangeObject) {
				return rangeObject
			}
		}
	}
	return temporalFieldRange(params)
}

func chronoFieldRangeRefinedBy(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.rangeRefinedBy: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.rangeRefinedBy: temporal is null")
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)
	return temporalAccessorRange([]any{temporalObject, selfObject})
}

func chronoFieldResolve(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.resolve: self is null")
	}
	return nil
}

func chronoFieldToString(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.toString: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if fieldEntry, exists := selfObject.FieldTable["displayName"]; exists {
			if stringObject, okStr := fieldEntry.Fvalue.(*object.Object); okStr && object.IsStringObject(stringObject) {
				return stringObject
			}
		}
	}
	return chronoFieldName(params)
}

func chronoFieldValueOf(params []any) any {
	ensureChronoFieldInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoField.valueOf: missing argument")
	}
	stringObject, okObj := params[0].(*object.Object)
	if !okObj {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoField.valueOf: argument is not a String")
	}
	if stringObject == nil || object.IsNull(stringObject) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoField.valueOf: name is null")
	}
	if !object.IsStringObject(stringObject) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoField.valueOf: argument is not a String")
	}
	targetName := object.GoStringFromStringObject(stringObject)
	for index, definition := range chronoFieldDefinitions {
		if definition.enumName == targetName {
			return chronoFieldInstances[index]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "No enum constant java.time.temporal.ChronoField."+targetName)
}

func chronoFieldValues([]any) any {
	ensureChronoFieldInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/temporal/ChronoField;", int64(len(chronoFieldInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, chronoFieldInstances)
	return arrayObject
}

func getChronoFieldOrdinal(fieldObject *object.Object) (int64, bool) {
	if fieldObject == nil || fieldObject.FieldTable == nil {
		return 0, false
	}
	if fieldEntry, exists := fieldObject.FieldTable["ordinal"]; exists {
		return getInt64Val(fieldEntry.Fvalue)
	}
	return 0, false
}

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"strings"
)

// Load_Time_Temporal_TemporalField loads method signatures for java/time/temporal/TemporalField.
func Load_Time_Temporal_TemporalField() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalField.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.adjustInto(Ljava/time/temporal/Temporal;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalFieldAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.getBaseUnit()Ljava/time/temporal/TemporalUnit;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalFieldGetBaseUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.getDisplayName(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalFieldGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.getFrom(Ljava/time/temporal/TemporalAccessor;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalFieldGetFrom,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.getRangeUnit()Ljava/time/temporal/TemporalUnit;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalFieldGetRangeUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.isDateBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalFieldIsDateBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.isSupportedBy(Ljava/time/temporal/TemporalAccessor;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalFieldIsSupportedBy,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.isTimeBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalFieldIsTimeBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.range()Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalFieldRange,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.rangeRefinedBy(Ljava/time/temporal/TemporalAccessor;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalFieldRangeRefinedBy,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.resolve(Ljava/util/Map;Ljava/time/temporal/TemporalAccessor;Ljava/time/format/ResolverStyle;)Ljava/time/temporal/TemporalAccessor;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  temporalFieldResolve,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalField.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalFieldToString,
		}
}

// temporalFieldAdjustInto returns a copy of the specified temporal object with the value of this field set.
// <R extends Temporal> R adjustInto(R temporal, long newValue)
func temporalFieldAdjustInto(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	return temporalWithField([]any{temporal, self, params[2]})
}

// temporalFieldGetBaseUnit gets the unit that the field is measured in.
// TemporalUnit getBaseUnit()
func temporalFieldGetBaseUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.getBaseUnit: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["baseUnit"]; ok {
			if unitObj, ok := fld.Fvalue.(*object.Object); ok && !object.IsNull(unitObj) {
				return unitObj
			}
		}
	}

	fieldName := getTemporalFieldName(self)
	upperName := strings.ToUpper(fieldName)
	switch upperName {
	case "NANO_OF_SECOND", "NANO_OF_DAY":
		return ghelpers.CreateTemporalUnit("Nanos")
	case "MICRO_OF_SECOND", "MICRO_OF_DAY":
		return ghelpers.CreateTemporalUnit("Micros")
	case "MILLI_OF_SECOND", "MILLI_OF_DAY":
		return ghelpers.CreateTemporalUnit("Millis")
	case "SECOND_OF_MINUTE", "SECOND_OF_DAY", "INSTANT_SECONDS":
		return ghelpers.CreateTemporalUnit("Seconds")
	case "MINUTE_OF_HOUR", "MINUTE_OF_DAY":
		return ghelpers.CreateTemporalUnit("Minutes")
	case "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM", "HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY":
		return ghelpers.CreateTemporalUnit("Hours")
	case "AMPM_OF_DAY":
		return ghelpers.CreateTemporalUnit("HalfDays")
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR", "DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY":
		return ghelpers.CreateTemporalUnit("Days")
	case "ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR":
		return ghelpers.CreateTemporalUnit("Weeks")
	case "MONTH_OF_YEAR", "PROLEPTIC_MONTH":
		return ghelpers.CreateTemporalUnit("Months")
	case "YEAR_OF_ERA", "YEAR":
		return ghelpers.CreateTemporalUnit("Years")
	case "CENTURY_OF_ERA":
		return ghelpers.CreateTemporalUnit("Centuries")
	case "ERA":
		return ghelpers.CreateTemporalUnit("Eras")
	default:
		return ghelpers.CreateTemporalUnit("Forever")
	}
}

// temporalFieldGetDisplayName gets the display name for the field in the requested locale.
// default String getDisplayName(Locale locale)
func temporalFieldGetDisplayName(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.getDisplayName: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.getDisplayName: locale is null")
	}
	return temporalFieldToString([]any{params[0]})
}

// temporalFieldGetFrom gets the value of this field from the specified temporal object.
// long getFrom(TemporalAccessor temporal)
func temporalFieldGetFrom(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.getFrom: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.getFrom: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	return temporalAccessorGetLong([]any{temporal, self})
}

// temporalFieldGetRangeUnit gets the range that the field is bound by.
// TemporalUnit getRangeUnit()
func temporalFieldGetRangeUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.getRangeUnit: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["rangeUnit"]; ok {
			if unitObj, ok := fld.Fvalue.(*object.Object); ok && !object.IsNull(unitObj) {
				return unitObj
			}
		}
	}

	fieldName := getTemporalFieldName(self)
	upperName := strings.ToUpper(fieldName)
	switch upperName {
	case "NANO_OF_SECOND", "MICRO_OF_SECOND", "MILLI_OF_SECOND":
		return ghelpers.CreateTemporalUnit("Seconds")
	case "NANO_OF_DAY", "MICRO_OF_DAY", "MILLI_OF_DAY", "SECOND_OF_DAY", "MINUTE_OF_DAY", "HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY":
		return ghelpers.CreateTemporalUnit("Days")
	case "SECOND_OF_MINUTE":
		return ghelpers.CreateTemporalUnit("Minutes")
	case "MINUTE_OF_HOUR", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM":
		return ghelpers.CreateTemporalUnit("Hours")
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR":
		return ghelpers.CreateTemporalUnit("Weeks")
	case "DAY_OF_MONTH", "ALIGNED_WEEK_OF_MONTH":
		return ghelpers.CreateTemporalUnit("Months")
	case "DAY_OF_YEAR", "ALIGNED_WEEK_OF_YEAR", "MONTH_OF_YEAR":
		return ghelpers.CreateTemporalUnit("Years")
	case "PROLEPTIC_MONTH", "YEAR_OF_ERA", "YEAR", "CENTURY_OF_ERA", "ERA", "INSTANT_SECONDS", "EPOCH_DAY":
		return ghelpers.CreateTemporalUnit("Forever")
	default:
		return ghelpers.CreateTemporalUnit("Forever")
	}
}

// temporalFieldIsDateBased checks if this field represents a component of a date.
// boolean isDateBased()
func temporalFieldIsDateBased(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.isDateBased: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["isDateBased"]; ok {
			if v, ok := fld.Fvalue.(int8); ok {
				return v
			}
		}
	}

	fieldName := getTemporalFieldName(self)
	upperName := strings.ToUpper(fieldName)
	switch upperName {
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR",
		"DAY_OF_MONTH", "DAY_OF_YEAR", "EPOCH_DAY",
		"ALIGNED_WEEK_OF_MONTH", "ALIGNED_WEEK_OF_YEAR",
		"MONTH_OF_YEAR", "PROLEPTIC_MONTH",
		"YEAR_OF_ERA", "YEAR", "CENTURY_OF_ERA", "ERA":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// temporalFieldIsSupportedBy checks if this field is supported by the specified temporal object.
// boolean isSupportedBy(TemporalAccessor temporal)
func temporalFieldIsSupportedBy(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.isSupportedBy: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	return temporalAccessorIsSupported([]any{temporal, self})
}

// temporalFieldIsTimeBased checks if this field represents a component of a time.
// boolean isTimeBased()
func temporalFieldIsTimeBased(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.isTimeBased: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["isTimeBased"]; ok {
			if v, ok := fld.Fvalue.(int8); ok {
				return v
			}
		}
	}

	fieldName := getTemporalFieldName(self)
	upperName := strings.ToUpper(fieldName)
	switch upperName {
	case "NANO_OF_SECOND", "NANO_OF_DAY",
		"MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY",
		"SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY",
		"HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY",
		"AMPM_OF_DAY":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// temporalFieldRange gets the range of valid values for the field.
// ValueRange range()
func temporalFieldRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.range: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["range"]; ok {
			if rangeObj, ok := fld.Fvalue.(*object.Object); ok && !object.IsNull(rangeObj) {
				return rangeObj
			}
		}
	}

	fieldName := getTemporalFieldName(self)
	upperName := strings.ToUpper(fieldName)
	switch upperName {
	case "NANO_OF_SECOND":
		return createValueRange(0, 0, 999_999_999, 999_999_999)
	case "NANO_OF_DAY":
		return createValueRange(0, 0, 86_399_999_999_999, 86_399_999_999_999)
	case "MICRO_OF_SECOND":
		return createValueRange(0, 0, 999_999, 999_999)
	case "MICRO_OF_DAY":
		return createValueRange(0, 0, 86_399_999_999, 86_399_999_999)
	case "MILLI_OF_SECOND":
		return createValueRange(0, 0, 999, 999)
	case "MILLI_OF_DAY":
		return createValueRange(0, 0, 86_399_999, 86_399_999)
	case "SECOND_OF_MINUTE":
		return createValueRange(0, 0, 59, 59)
	case "SECOND_OF_DAY":
		return createValueRange(0, 0, 86_399, 86_399)
	case "MINUTE_OF_HOUR":
		return createValueRange(0, 0, 59, 59)
	case "MINUTE_OF_DAY":
		return createValueRange(0, 0, 1439, 1439)
	case "HOUR_OF_AMPM":
		return createValueRange(0, 0, 11, 11)
	case "CLOCK_HOUR_OF_AMPM":
		return createValueRange(1, 1, 12, 12)
	case "HOUR_OF_DAY":
		return createValueRange(0, 0, 23, 23)
	case "CLOCK_HOUR_OF_DAY":
		return createValueRange(1, 1, 24, 24)
	case "AMPM_OF_DAY":
		return createValueRange(0, 0, 1, 1)
	case "DAY_OF_WEEK", "ALIGNED_DAY_OF_WEEK_IN_MONTH", "ALIGNED_DAY_OF_WEEK_IN_YEAR":
		return createValueRange(1, 1, 7, 7)
	case "DAY_OF_MONTH":
		return createValueRange(1, 1, 28, 31)
	case "DAY_OF_YEAR":
		return createValueRange(1, 1, 365, 366)
	case "EPOCH_DAY":
		return createValueRange(-365243219162, -365243219162, 365241780471, 365241780471)
	case "ALIGNED_WEEK_OF_MONTH":
		return createValueRange(1, 1, 4, 5)
	case "ALIGNED_WEEK_OF_YEAR":
		return createValueRange(1, 1, 53, 53)
	case "MONTH_OF_YEAR":
		return createValueRange(1, 1, 12, 12)
	case "PROLEPTIC_MONTH":
		return createValueRange(-11999999988, -11999999988, 11999999999, 11999999999)
	case "YEAR_OF_ERA":
		return createValueRange(1, 1, 999999999, 1000000000)
	case "YEAR":
		return createValueRange(-999999999, -999999999, 999999999, 999999999)
	case "ERA":
		return createValueRange(0, 0, 1, 1)
	case "INSTANT_SECONDS":
		return createValueRange(math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64)
	default:
		return createValueRange(math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64)
	}
}

// temporalFieldRangeRefinedBy gets the range of valid values for the field using the temporal object to refine the result.
// ValueRange rangeRefinedBy(TemporalAccessor temporal)
func temporalFieldRangeRefinedBy(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.rangeRefinedBy: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.rangeRefinedBy: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	return temporalAccessorRange([]any{temporal, self})
}

// temporalFieldResolve resolves this field to calculate other fields and/or a date-time.
// default TemporalAccessor resolve(Map<TemporalField,Long> fieldValues, TemporalAccessor partialTemporal, ResolverStyle resolverStyle)
func temporalFieldResolve(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.resolve: self is null")
	}
	return nil
}

// temporalFieldToString returns the name of the field.
// String toString()
func temporalFieldToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalField.toString: self is null")
	}
	self := params[0].(*object.Object)

	fieldName := getTemporalFieldName(self)
	if fieldName != "" {
		return object.StringObjectFromGoString(fieldName)
	}
	return object.StringObjectFromGoString("TemporalField")
}

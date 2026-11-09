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
)

func Load_Time_Temporal_Temporal() {
	ghelpers.MethodSignatures["java/time/temporal/Temporal.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalPlusUnit,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalPlusAmount,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalUntil,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/temporal/Temporal.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalWithField,
		}
}

func getTemporalUnitName(unit *object.Object) string {
	if unit == nil || unit.FieldTable == nil {
		return ""
	}
	if nameFld, ok := unit.FieldTable["name"]; ok {
		if nameObj, ok := nameFld.Fvalue.(*object.Object); ok && object.IsStringObject(nameObj) {
			return object.GoStringFromStringObject(nameObj)
		}
	}
	return ""
}

// cloneTemporalObject clones a temporal object and its field table.
func cloneTemporalObject(src *object.Object) *object.Object {
	if src == nil {
		return nil
	}
	return object.CloneObject(src)
}

// temporalIsSupportedUnit checks if the specified unit is supported.
// boolean isSupported(TemporalUnit unit)
func temporalIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)

	if self.FieldTable != nil {
		unitName := getTemporalUnitName(unit)
		if unitName != "" {
			if _, ok := self.FieldTable[unitName]; ok {
				return types.JavaBoolTrue
			}
		}
	}
	return types.JavaBoolFalse
}

// temporalWithAdjuster returns an adjusted object of the same type with the adjustment made.
// default Temporal with(TemporalAdjuster adjuster)
func temporalWithAdjuster(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	// Return self or cloned self
	return cloneTemporalObject(self)
}

// temporalWithField returns an object of the same type with the specified field altered.
// Temporal with(TemporalField field, long newValue)
func temporalWithField(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.with: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)

	newValue, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Temporal.with: invalid newValue type")
	}

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	switch className {
	case "java/time/LocalTime":
		return localTimeWithField([]any{self, field, newValue})
	case "java/time/LocalDate":
		return localDateWithField([]any{self, field, newValue})
	case "java/time/LocalDateTime":
		return localDateTimeWithField([]any{self, field, newValue})
	case "java/time/ZonedDateTime":
		return zonedDateTimeWithField([]any{self, field, newValue})
	case "java/time/OffsetDateTime":
		return offsetDateTimeWithField([]any{self, field, newValue})
	case "java/time/Instant":
		return instantWithField([]any{self, field, newValue})
	case "java/time/Year":
		return yearWithField([]any{self, field, newValue})
	case "java/time/chrono/ChronoLocalDate":
		return chronoLocalDateWithField([]any{self, field, newValue})
	case "java/time/chrono/ChronoLocalDateTime":
		return chronoLocalDateTimeWithField([]any{self, field, newValue})
	case "java/time/chrono/ChronoZonedDateTime":
		return chronoZonedDateTimeWithField([]any{self, field, newValue})
	}

	res := cloneTemporalObject(self)
	fieldName := getTemporalFieldName(field)
	if fieldName != "" {
		res.FieldTable[fieldName] = object.Field{Ftype: types.Long, Fvalue: newValue}
		return res
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Temporal.with: unsupported field")
}

// temporalPlusAmount returns an object of the same type with the specified amount added.
// default Temporal plus(TemporalAmount amount)
func temporalPlusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	switch className {
	case "java/time/LocalDate":
		return localDatePlusAmount([]any{self, amount})
	case "java/time/LocalTime":
		return localTimePlusAmount([]any{self, amount})
	case "java/time/LocalDateTime":
		return localDateTimePlusAmount([]any{self, amount})
	case "java/time/ZonedDateTime":
		return zonedDateTimePlusAmount([]any{self, amount})
	case "java/time/OffsetDateTime":
		return offsetDateTimePlusAmount([]any{self, amount})
	case "java/time/Instant":
		return instantPlusAmount([]any{self, amount})
	case "java/time/Year":
		return yearPlusAmount([]any{self, amount})
	case "java/time/chrono/ChronoLocalDate":
		return chronoLocalDatePlusAmount([]any{self, amount})
	case "java/time/chrono/ChronoLocalDateTime":
		return chronoLocalDateTimePlusAmount([]any{self, amount})
	case "java/time/chrono/ChronoZonedDateTime":
		return chronoZonedDateTimePlusAmount([]any{self, amount})
	}

	res := cloneTemporalObject(self)
	if amount.FieldTable != nil {
		if secondsFld, ok := amount.FieldTable["seconds"]; ok {
			if secVal, ok := getInt64Val(secondsFld.Fvalue); ok {
				if curSec, ok := res.FieldTable["seconds"]; ok {
					if curVal, ok := getInt64Val(curSec.Fvalue); ok {
						res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: curVal + secVal}
					}
				}
			}
		}
		if nanosFld, ok := amount.FieldTable["nanos"]; ok {
			if nanosVal, ok := getInt64Val(nanosFld.Fvalue); ok {
				if curNanos, ok := res.FieldTable["nanos"]; ok {
					if curVal, ok := getInt64Val(curNanos.Fvalue); ok {
						res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: curVal + nanosVal}
					}
				}
			}
		}
	}
	return res
}

// temporalPlusUnit returns an object of the same type with the specified period added.
// Temporal plus(long amountToAdd, TemporalUnit unit)
func temporalPlusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.plus: self is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.plus: unit is null")
	}
	self := params[0].(*object.Object)
	unit := params[2].(*object.Object)

	amountToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Temporal.plus: invalid amountToAdd type")
	}

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	switch className {
	case "java/time/LocalDate":
		return localDatePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/LocalTime":
		return localTimePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/LocalDateTime":
		return localDateTimePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/ZonedDateTime":
		return zonedDateTimePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/OffsetDateTime":
		return offsetDateTimePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/Instant":
		return instantPlusUnit([]any{self, amountToAdd, unit})
	case "java/time/Year":
		return yearPlusUnit([]any{self, amountToAdd, unit})
	case "java/time/chrono/ChronoLocalDate":
		return chronoLocalDatePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/chrono/ChronoLocalDateTime":
		return chronoLocalDateTimePlusUnit([]any{self, amountToAdd, unit})
	case "java/time/chrono/ChronoZonedDateTime":
		return chronoZonedDateTimePlusUnit([]any{self, amountToAdd, unit})
	}

	res := cloneTemporalObject(self)
	unitName := getTemporalUnitName(unit)
	if unitName != "" {
		if fld, ok := res.FieldTable[unitName]; ok {
			if val, ok := getInt64Val(fld.Fvalue); ok {
				res.FieldTable[unitName] = object.Field{Ftype: types.Long, Fvalue: val + amountToAdd}
				return res
			}
		}
	}
	return res
}

// temporalMinusAmount returns an object of the same type with the specified amount subtracted.
// default Temporal minus(TemporalAmount amount)
func temporalMinusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	switch className {
	case "java/time/LocalDate":
		return localDateMinusAmount([]any{self, amount})
	case "java/time/LocalTime":
		return localTimeMinusAmount([]any{self, amount})
	case "java/time/LocalDateTime":
		return localDateTimeMinusAmount([]any{self, amount})
	case "java/time/ZonedDateTime":
		return zonedDateTimeMinusAmount([]any{self, amount})
	case "java/time/OffsetDateTime":
		return offsetDateTimeMinusAmount([]any{self, amount})
	case "java/time/Instant":
		return instantMinusAmount([]any{self, amount})
	case "java/time/Year":
		return yearMinusAmount([]any{self, amount})
	case "java/time/chrono/ChronoLocalDate":
		return chronoLocalDateMinusAmount([]any{self, amount})
	case "java/time/chrono/ChronoLocalDateTime":
		return chronoLocalDateTimeMinusAmount([]any{self, amount})
	case "java/time/chrono/ChronoZonedDateTime":
		return chronoZonedDateTimeMinusAmount([]any{self, amount})
	}

	res := cloneTemporalObject(self)
	if amount.FieldTable != nil {
		if secondsFld, ok := amount.FieldTable["seconds"]; ok {
			if secVal, ok := getInt64Val(secondsFld.Fvalue); ok {
				if curSec, ok := res.FieldTable["seconds"]; ok {
					if curVal, ok := getInt64Val(curSec.Fvalue); ok {
						res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: curVal - secVal}
					}
				}
			}
		}
		if nanosFld, ok := amount.FieldTable["nanos"]; ok {
			if nanosVal, ok := getInt64Val(nanosFld.Fvalue); ok {
				if curNanos, ok := res.FieldTable["nanos"]; ok {
					if curVal, ok := getInt64Val(curNanos.Fvalue); ok {
						res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: curVal - nanosVal}
					}
				}
			}
		}
	}
	return res
}

// temporalMinusUnit returns an object of the same type with the specified period subtracted.
// default Temporal minus(long amountToSubtract, TemporalUnit unit)
func temporalMinusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.minus: self is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.minus: unit is null")
	}

	self := params[0].(*object.Object)
	unit := params[2].(*object.Object)
	amountToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Temporal.minus: invalid amountToSubtract type")
	}

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	switch className {
	case "java/time/LocalDate":
		return localDateMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/LocalTime":
		return localTimeMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/LocalDateTime":
		return localDateTimeMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/ZonedDateTime":
		return zonedDateTimeMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/OffsetDateTime":
		return offsetDateTimeMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/Instant":
		return instantMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/Year":
		return yearMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/chrono/ChronoLocalDate":
		return chronoLocalDateMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/chrono/ChronoLocalDateTime":
		return chronoLocalDateTimeMinusUnit([]any{self, amountToSubtract, unit})
	case "java/time/chrono/ChronoZonedDateTime":
		return chronoZonedDateTimeMinusUnit([]any{self, amountToSubtract, unit})
	}

	if amountToSubtract == math.MinInt64 {
		// handle overflow case as plus(Long.MAX_VALUE, unit).plus(1, unit)
		p1 := temporalPlusUnit([]any{params[0], int64(math.MaxInt64), params[2]})
		if err, ok := p1.(*ghelpers.GErrBlk); ok {
			return err
		}
		return temporalPlusUnit([]any{p1, int64(1), params[2]})
	}

	return temporalPlusUnit([]any{params[0], -amountToSubtract, params[2]})
}

// temporalUntil calculates the amount of time until another temporal in terms of the specified unit.
// long until(Temporal endExclusive, TemporalUnit unit)
func temporalUntil(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.until: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.until: endExclusive is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Temporal.until: unit is null")
	}
	self := params[0].(*object.Object)
	endExclusive := params[1].(*object.Object)
	unit := params[2].(*object.Object)

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	switch className {
	case "java/time/LocalDate":
		return localDateUntilUnit([]any{self, endExclusive, unit})
	case "java/time/LocalTime":
		return localTimeUntilUnit([]any{self, endExclusive, unit})
	case "java/time/LocalDateTime":
		return localDateTimeUntil([]any{self, endExclusive, unit})
	case "java/time/ZonedDateTime":
		return zonedDateTimeUntil([]any{self, endExclusive, unit})
	case "java/time/OffsetDateTime":
		return offsetDateTimeUntil([]any{self, endExclusive, unit})
	case "java/time/Instant":
		return instantUntil([]any{self, endExclusive, unit})
	case "java/time/Year":
		return yearUntil([]any{self, endExclusive, unit})
	case "java/time/chrono/ChronoLocalDate":
		return chronoLocalDateUntilUnit([]any{self, endExclusive, unit})
	case "java/time/chrono/ChronoLocalDateTime":
		return chronoLocalDateTimeUntilUnit([]any{self, endExclusive, unit})
	}

	unitName := getTemporalUnitName(unit)
	if unitName != "" && self.FieldTable != nil && endExclusive.FieldTable != nil {
		startFld, ok1 := self.FieldTable[unitName]
		endFld, ok2 := endExclusive.FieldTable[unitName]
		if ok1 && ok2 {
			startVal, ok1 := getInt64Val(startFld.Fvalue)
			endVal, ok2 := getInt64Val(endFld.Fvalue)
			if ok1 && ok2 {
				return endVal - startVal
			}
		}
	}

	return int64(0)
}

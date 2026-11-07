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

// Load_Time_Temporal_TemporalUnit loads method signatures for java/time/temporal/TemporalUnit.
func Load_Time_Temporal_TemporalUnit() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.addTo(Ljava/time/temporal/Temporal;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalUnitAddTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.between(Ljava/time/temporal/Temporal;Ljava/time/temporal/Temporal;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalUnitBetween,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.getDuration()Ljava/time/Duration;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalUnitGetDuration,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.isDateBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalUnitIsDateBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.isDurationEstimated()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalUnitIsDurationEstimated,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.isSupportedBy(Ljava/time/temporal/Temporal;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalUnitIsSupportedBy,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.isTimeBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalUnitIsTimeBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalUnit.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalUnitToString,
		}
}

// temporalUnitAddTo returns a copy of the specified temporal object with the specified period added.
// <R extends Temporal> R addTo(R temporal, long amount)
func temporalUnitAddTo(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.addTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.addTo: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	amount, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalUnit.addTo: invalid amount type")
	}

	unitName := getTemporalUnitName(self)
	res := cloneTemporalObject(temporal)
	if res.FieldTable == nil {
		res.FieldTable = make(map[string]object.Field)
	}

	upperUnit := strings.ToUpper(unitName)
	switch upperUnit {
	case "NANOS":
		if cur, ok := res.FieldTable["nanos"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: v + amount}
				return res
			}
		}
		res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: amount}
		return res
	case "MICROS":
		nanosToAdd := amount * 1000
		if cur, ok := res.FieldTable["nanos"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: v + nanosToAdd}
				return res
			}
		}
		res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: nanosToAdd}
		return res
	case "MILLIS":
		nanosToAdd := amount * 1_000_000
		if cur, ok := res.FieldTable["nanos"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: v + nanosToAdd}
				return res
			}
		}
		res.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: nanosToAdd}
		return res
	case "SECONDS":
		if cur, ok := res.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: v + amount}
				return res
			}
		}
		res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: amount}
		return res
	case "MINUTES":
		secToAdd := amount * 60
		if cur, ok := res.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: v + secToAdd}
				return res
			}
		}
		res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: secToAdd}
		return res
	case "HOURS":
		secToAdd := amount * 3600
		if cur, ok := res.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: v + secToAdd}
				return res
			}
		}
		res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: secToAdd}
		return res
	case "HALF_DAYS":
		secToAdd := amount * 43200
		if cur, ok := res.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: v + secToAdd}
				return res
			}
		}
		res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: secToAdd}
		return res
	case "DAYS":
		secToAdd := amount * 86400
		if cur, ok := res.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(cur.Fvalue); ok {
				res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: v + secToAdd}
				return res
			}
		}
		res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: secToAdd}
		return res
	default:
		if unitName != "" {
			if cur, ok := res.FieldTable[unitName]; ok {
				if v, ok := getInt64Val(cur.Fvalue); ok {
					res.FieldTable[unitName] = object.Field{Ftype: types.Long, Fvalue: v + amount}
					return res
				}
			}
			res.FieldTable[unitName] = object.Field{Ftype: types.Long, Fvalue: amount}
			return res
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalUnit.addTo: unsupported unit")
}

// temporalUnitBetween calculates the amount of time between two temporal objects.
// long between(Temporal temporal1Inclusive, Temporal temporal2Exclusive)
func temporalUnitBetween(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.between: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.between: temporal1Inclusive is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.between: temporal2Exclusive is null")
	}
	self := params[0].(*object.Object)
	t1 := params[1].(*object.Object)
	t2 := params[2].(*object.Object)

	var s1, s2, n1, n2 int64
	if t1.FieldTable != nil {
		if secFld, ok := t1.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(secFld.Fvalue); ok {
				s1 = v
			}
		}
		if nanoFld, ok := t1.FieldTable["nanos"]; ok {
			if v, ok := getInt64Val(nanoFld.Fvalue); ok {
				n1 = v
			}
		}
	}
	if t2.FieldTable != nil {
		if secFld, ok := t2.FieldTable["seconds"]; ok {
			if v, ok := getInt64Val(secFld.Fvalue); ok {
				s2 = v
			}
		}
		if nanoFld, ok := t2.FieldTable["nanos"]; ok {
			if v, ok := getInt64Val(nanoFld.Fvalue); ok {
				n2 = v
			}
		}
	}

	diffSeconds := s2 - s1
	diffNanos := n2 - n1
	totalNanos := diffSeconds*1_000_000_000 + diffNanos

	unitName := getTemporalUnitName(self)
	upperUnit := strings.ToUpper(unitName)
	switch upperUnit {
	case "NANOS":
		return totalNanos
	case "MICROS":
		return totalNanos / 1000
	case "MILLIS":
		return totalNanos / 1_000_000
	case "SECONDS":
		return diffSeconds
	case "MINUTES":
		return diffSeconds / 60
	case "HOURS":
		return diffSeconds / 3600
	case "HALF_DAYS":
		return diffSeconds / 43200
	case "DAYS":
		return diffSeconds / 86400
	default:
		if unitName != "" && t1.FieldTable != nil && t2.FieldTable != nil {
			v1Fld, ok1 := t1.FieldTable[unitName]
			v2Fld, ok2 := t2.FieldTable[unitName]
			if ok1 && ok2 {
				var v1, v2 int64
				if v, ok := getInt64Val(v1Fld.Fvalue); ok {
					v1 = v
				}
				if v, ok := getInt64Val(v2Fld.Fvalue); ok {
					v2 = v
				}
				return v2 - v1
			}
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalUnit.between: unsupported unit")
}

// temporalUnitGetDuration gets the duration of this unit, which may be an estimate.
// Duration getDuration()
func temporalUnitGetDuration(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.getDuration: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if durFld, ok := self.FieldTable["duration"]; ok {
			if durObj, ok := durFld.Fvalue.(*object.Object); ok && !object.IsNull(durObj) {
				return durObj
			}
		}
	}

	unitName := getTemporalUnitName(self)
	upperUnit := strings.ToUpper(unitName)
	switch upperUnit {
	case "NANOS":
		return createDuration(0, 1)
	case "MICROS":
		return createDuration(0, 1000)
	case "MILLIS":
		return createDuration(0, 1_000_000)
	case "SECONDS":
		return createDuration(1, 0)
	case "MINUTES":
		return createDuration(60, 0)
	case "HOURS":
		return createDuration(3600, 0)
	case "HALF_DAYS":
		return createDuration(43200, 0)
	case "DAYS":
		return createDuration(86400, 0)
	case "WEEKS":
		return createDuration(7*86400, 0)
	case "MONTHS":
		return createDuration(30*86400+43200, 0) // ~30.4375 days
	case "YEARS":
		return createDuration(365*86400+6*3600, 0) // ~365.2425 days
	case "DECADES":
		return createDuration(10*(365*86400+6*3600), 0)
	case "CENTURIES":
		return createDuration(100*(365*86400+6*3600), 0)
	case "MILLENNIA":
		return createDuration(1000*(365*86400+6*3600), 0)
	case "ERAS":
		return createDuration(1_000_000_000*(365*86400), 0)
	case "FOREVER":
		return createDuration(math.MaxInt64, 999_999_999)
	default:
		return createDuration(0, 0)
	}
}

// temporalUnitIsDateBased checks if this unit represents a component of a date.
// boolean isDateBased()
func temporalUnitIsDateBased(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.isDateBased: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["isDateBased"]; ok {
			if v, ok := fld.Fvalue.(int8); ok {
				return v
			}
		}
	}

	unitName := getTemporalUnitName(self)
	upperUnit := strings.ToUpper(unitName)
	switch upperUnit {
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// temporalUnitIsDurationEstimated checks if the duration of the unit is an estimate.
// boolean isDurationEstimated()
func temporalUnitIsDurationEstimated(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.isDurationEstimated: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["isDurationEstimated"]; ok {
			if v, ok := fld.Fvalue.(int8); ok {
				return v
			}
		}
	}

	unitName := getTemporalUnitName(self)
	upperUnit := strings.ToUpper(unitName)
	switch upperUnit {
	case "DAYS", "WEEKS", "MONTHS", "YEARS", "DECADES", "CENTURIES", "MILLENNIA", "ERAS", "FOREVER":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// temporalUnitIsSupportedBy checks if this unit is supported by the specified temporal object.
// default boolean isSupportedBy(Temporal temporal)
func temporalUnitIsSupportedBy(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.isSupportedBy: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	return temporalIsSupportedUnit([]any{temporal, self})
}

// temporalUnitIsTimeBased checks if this unit represents a component of a time.
// boolean isTimeBased()
func temporalUnitIsTimeBased(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.isTimeBased: self is null")
	}
	self := params[0].(*object.Object)

	if self.FieldTable != nil {
		if fld, ok := self.FieldTable["isTimeBased"]; ok {
			if v, ok := fld.Fvalue.(int8); ok {
				return v
			}
		}
	}

	unitName := getTemporalUnitName(self)
	upperUnit := strings.ToUpper(unitName)
	switch upperUnit {
	case "NANOS", "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

// temporalUnitToString returns the name of this unit.
// String toString()
func temporalUnitToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalUnit.toString: self is null")
	}
	self := params[0].(*object.Object)

	unitName := getTemporalUnitName(self)
	if unitName != "" {
		return object.StringObjectFromGoString(unitName)
	}
	return object.StringObjectFromGoString("TemporalUnit")
}

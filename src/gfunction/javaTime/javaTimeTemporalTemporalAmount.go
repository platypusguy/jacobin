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
	"strings"
)

var classNameArrayList = "java/util/ArrayList"

// Load_Time_Temporal_TemporalAmount loads method signatures for java/time/temporal/TemporalAmount.
func Load_Time_Temporal_TemporalAmount() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalAmount.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAmount.addTo(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAmountAddTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAmount.get(Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAmountGet,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAmount.getUnits()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAmountGetUnits,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAmount.subtractFrom(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAmountSubtractFrom,
		}
}

// temporalAmountAddTo adds to the specified temporal object.
// Temporal addTo(Temporal temporal)
func temporalAmountAddTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.addTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.addTo: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	res := cloneTemporalObject(temporal)
	if self.FieldTable != nil && res.FieldTable != nil {
		if secondsFld, ok := self.FieldTable["seconds"]; ok {
			if secVal, ok := getInt64Val(secondsFld.Fvalue); ok {
				if curSec, ok := res.FieldTable["seconds"]; ok {
					if curVal, ok := getInt64Val(curSec.Fvalue); ok {
						res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: curVal + secVal}
					}
				}
			}
		}
		if nanosFld, ok := self.FieldTable["nanos"]; ok {
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

// temporalAmountGet returns the value of the requested unit.
// long get(TemporalUnit unit)
func temporalAmountGet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.get: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.get: unit is null")
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)

	unitName := getTemporalUnitName(unit)
	if unitName != "" && self.FieldTable != nil {
		if fld, ok := self.FieldTable[unitName]; ok {
			if val, ok := getInt64Val(fld.Fvalue); ok {
				return val
			}
		}
		lowerName := strings.ToLower(unitName)
		if fld, ok := self.FieldTable[lowerName]; ok {
			if val, ok := getInt64Val(fld.Fvalue); ok {
				return val
			}
		}
	}
	return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalAmount.get: unsupported unit")
}

// temporalAmountGetUnits returns the list of units supported by this amount.
// List<TemporalUnit> getUnits()
func temporalAmountGetUnits(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.getUnits: self is null")
	}
	self := params[0].(*object.Object)

	listObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	var unitsList []any
	if self.FieldTable != nil {
		if unitsFld, ok := self.FieldTable["units"]; ok {
			if uList, ok := unitsFld.Fvalue.([]any); ok {
				unitsList = append([]any(nil), uList...)
			}
		}
	}
	if unitsList == nil {
		unitsList = make([]any, 0)
	}
	listObj.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: unitsList,
	}
	return listObj
}

// temporalAmountSubtractFrom subtracts this amount from the specified temporal object.
// Temporal subtractFrom(Temporal temporal)
func temporalAmountSubtractFrom(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.subtractFrom: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAmount.subtractFrom: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	res := cloneTemporalObject(temporal)
	if self.FieldTable != nil && res.FieldTable != nil {
		if secondsFld, ok := self.FieldTable["seconds"]; ok {
			if secVal, ok := getInt64Val(secondsFld.Fvalue); ok {
				if curSec, ok := res.FieldTable["seconds"]; ok {
					if curVal, ok := getInt64Val(curSec.Fvalue); ok {
						res.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: curVal - secVal}
					}
				}
			}
		}
		if nanosFld, ok := self.FieldTable["nanos"]; ok {
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

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

var classNameValueRange = "java/time/temporal/ValueRange"

func Load_Time_Temporal_TemporalAccessor() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalAccessor.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAccessor.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAccessorGet,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAccessor.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAccessorGetLong,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAccessor.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAccessorIsSupported,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAccessor.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAccessorQuery,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAccessor.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAccessorRange,
		}
}

// createValueRange allocates a new ValueRange object
func createValueRange(minSmallest, minLargest, maxSmallest, maxLargest int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameValueRange)
	obj.FieldTable["minSmallest"] = object.Field{Ftype: types.Long, Fvalue: minSmallest}
	obj.FieldTable["minLargest"] = object.Field{Ftype: types.Long, Fvalue: minLargest}
	obj.FieldTable["maxSmallest"] = object.Field{Ftype: types.Long, Fvalue: maxSmallest}
	obj.FieldTable["maxLargest"] = object.Field{Ftype: types.Long, Fvalue: maxLargest}
	return obj
}

func getTemporalFieldName(field *object.Object) string {
	if field == nil || field.FieldTable == nil {
		return ""
	}
	if nameFld, ok := field.FieldTable["name"]; ok {
		if nameObj, ok := nameFld.Fvalue.(*object.Object); ok && object.IsStringObject(nameObj) {
			return object.GoStringFromStringObject(nameObj)
		}
	}
	return ""
}

// temporalAccessorIsSupported checks if the specified field is supported.
func temporalAccessorIsSupported(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)

	if self.FieldTable != nil {
		fieldName := getTemporalFieldName(field)
		if fieldName != "" {
			if _, ok := self.FieldTable[fieldName]; ok {
				return types.JavaBoolTrue
			}
		}
	}
	return types.JavaBoolFalse
}

// temporalAccessorGetLong gets the value of the specified field as a long.
func temporalAccessorGetLong(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.getLong: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.getLong: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)

	if self.FieldTable != nil {
		fieldName := getTemporalFieldName(field)
		if fieldName != "" {
			if fld, ok := self.FieldTable[fieldName]; ok {
				if val, ok := getInt64Val(fld.Fvalue); ok {
					return val
				}
			}
		}
	}
	return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalAccessor.getLong: field is not supported")
}

// temporalAccessorRange gets the range of valid values for the specified field.
func temporalAccessorRange(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.range: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.range: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)

	isSupp := temporalAccessorIsSupported([]any{self, field})
	if isSupp == types.JavaBoolTrue {
		if field.FieldTable != nil {
			if r, ok := field.FieldTable["range"].Fvalue.(*object.Object); ok && !object.IsNull(r) {
				return r
			}
		}
		return createValueRange(math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64)
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalAccessor.range: unsupported field")
}

// temporalAccessorGet gets the value of the specified field as an int.
func temporalAccessorGet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.get: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.get: field is null")
	}

	rangeRes := temporalAccessorRange(params)
	if err, ok := rangeRes.(*ghelpers.GErrBlk); ok {
		return err
	}

	longRes := temporalAccessorGetLong(params)
	if err, ok := longRes.(*ghelpers.GErrBlk); ok {
		return err
	}

	val, ok := getInt64Val(longRes)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalAccessor.get: invalid field value type")
	}
	if val < math.MinInt32 || val > math.MaxInt32 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "TemporalAccessor.get: value out of int range")
	}

	return val
}

// temporalAccessorQuery queries this date-time.
func temporalAccessorQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.query: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAccessor.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)
	return temporalQueryQueryFrom([]any{query, self})
}

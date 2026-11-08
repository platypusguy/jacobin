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
	"sync"
)

var classNameTemporalQueriesVal = "java/time/temporal/TemporalQueries"
var classNameTemporalQueryVal = "java/time/temporal/TemporalQuery"

var temporalQueriesMutex sync.Mutex
var temporalQueriesInitialized bool = false

var zoneIdQueryInstance *object.Object
var chronologyQueryInstance *object.Object
var precisionQueryInstance *object.Object
var zoneQueryInstance *object.Object
var offsetQueryInstance *object.Object
var localDateQueryInstance *object.Object
var localTimeQueryInstance *object.Object

func makeTemporalQueryObject(queryName string) *object.Object {
	queryObject := object.MakeEmptyObjectWithClassName(&classNameTemporalQueryVal)
	queryObject.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString(queryName),
	}
	return queryObject
}

func ensureTemporalQueriesInited() {
	temporalQueriesMutex.Lock()
	defer temporalQueriesMutex.Unlock()
	if temporalQueriesInitialized {
		return
	}

	zoneIdQueryInstance = makeTemporalQueryObject("zoneId")
	chronologyQueryInstance = makeTemporalQueryObject("chronology")
	precisionQueryInstance = makeTemporalQueryObject("precision")
	zoneQueryInstance = makeTemporalQueryObject("zone")
	offsetQueryInstance = makeTemporalQueryObject("offset")
	localDateQueryInstance = makeTemporalQueryObject("localDate")
	localTimeQueryInstance = makeTemporalQueryObject("localTime")

	_ = statics.AddStatic(classNameTemporalQueriesVal+".ZONE_ID", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: zoneIdQueryInstance,
	})
	_ = statics.AddStatic(classNameTemporalQueriesVal+".CHRONO", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: chronologyQueryInstance,
	})
	_ = statics.AddStatic(classNameTemporalQueriesVal+".PRECISION", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: precisionQueryInstance,
	})
	_ = statics.AddStatic(classNameTemporalQueriesVal+".ZONE", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: zoneQueryInstance,
	})
	_ = statics.AddStatic(classNameTemporalQueriesVal+".OFFSET", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: offsetQueryInstance,
	})
	_ = statics.AddStatic(classNameTemporalQueriesVal+".LOCAL_DATE", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: localDateQueryInstance,
	})
	_ = statics.AddStatic(classNameTemporalQueriesVal+".LOCAL_TIME", statics.Static{
		Type:  "Ljava/time/temporal/TemporalQuery;",
		Value: localTimeQueryInstance,
	})

	temporalQueriesInitialized = true
}

func Load_Time_Temporal_TemporalQueries() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesClinit,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.chronology()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesChronology,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.localDate()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesLocalDate,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.localTime()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesLocalTime,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.offset()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesOffset,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.precision()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesPrecision,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.zone()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesZone,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQueries.zoneId()Ljava/time/temporal/TemporalQuery;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalQueriesZoneId,
		}
}

func temporalQueriesClinit([]any) any {
	ensureTemporalQueriesInited()
	return nil
}

func temporalQueriesChronology([]any) any {
	ensureTemporalQueriesInited()
	return chronologyQueryInstance
}

func temporalQueriesLocalDate([]any) any {
	ensureTemporalQueriesInited()
	return localDateQueryInstance
}

func temporalQueriesLocalTime([]any) any {
	ensureTemporalQueriesInited()
	return localTimeQueryInstance
}

func temporalQueriesOffset([]any) any {
	ensureTemporalQueriesInited()
	return offsetQueryInstance
}

func temporalQueriesPrecision([]any) any {
	ensureTemporalQueriesInited()
	return precisionQueryInstance
}

func temporalQueriesZone([]any) any {
	ensureTemporalQueriesInited()
	return zoneQueryInstance
}

func temporalQueriesZoneId([]any) any {
	ensureTemporalQueriesInited()
	return zoneIdQueryInstance
}

func executeTemporalQuery(queryObject *object.Object, temporalObject *object.Object) any {
	if queryObject == nil || object.IsNull(queryObject) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalQuery.queryFrom: query is null")
	}
	if temporalObject == nil || object.IsNull(temporalObject) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalQuery.queryFrom: temporal is null")
	}

	queryName := ""
	if queryObject.FieldTable != nil {
		if queryNameField, exists := queryObject.FieldTable["name"]; exists {
			if queryNameStringObject, okString := queryNameField.Fvalue.(*object.Object); okString && object.IsStringObject(queryNameStringObject) {
				queryName = object.GoStringFromStringObject(queryNameStringObject)
			}
		}
	}

	if temporalObject.FieldTable != nil {
		switch queryName {
		case "zoneId":
			if zoneIdField, exists := temporalObject.FieldTable["zoneId"]; exists {
				if zoneIdObject, okZoneId := zoneIdField.Fvalue.(*object.Object); okZoneId && !object.IsNull(zoneIdObject) {
					return zoneIdObject
				}
			}
			if zoneField, exists := temporalObject.FieldTable["zone"]; exists {
				if zoneObject, okZone := zoneField.Fvalue.(*object.Object); okZone && !object.IsNull(zoneObject) {
					return zoneObject
				}
			}
			return nil

		case "zone":
			if zoneField, exists := temporalObject.FieldTable["zone"]; exists {
				if zoneObject, okZone := zoneField.Fvalue.(*object.Object); okZone && !object.IsNull(zoneObject) {
					return zoneObject
				}
			}
			if zoneIdField, exists := temporalObject.FieldTable["zoneId"]; exists {
				if zoneIdObject, okZoneId := zoneIdField.Fvalue.(*object.Object); okZoneId && !object.IsNull(zoneIdObject) {
					return zoneIdObject
				}
			}
			if offsetField, exists := temporalObject.FieldTable["offset"]; exists {
				if offsetObject, okOffset := offsetField.Fvalue.(*object.Object); okOffset && !object.IsNull(offsetObject) {
					return offsetObject
				}
			}
			return nil

		case "chronology":
			if chronologyField, exists := temporalObject.FieldTable["chronology"]; exists {
				if chronologyObject, okChronology := chronologyField.Fvalue.(*object.Object); okChronology && !object.IsNull(chronologyObject) {
					return chronologyObject
				}
			}
			return nil

		case "precision":
			if precisionField, exists := temporalObject.FieldTable["precision"]; exists {
				if precisionObject, okPrecision := precisionField.Fvalue.(*object.Object); okPrecision && !object.IsNull(precisionObject) {
					return precisionObject
				}
			}
			return nil

		case "offset":
			if offsetField, exists := temporalObject.FieldTable["offset"]; exists {
				if offsetObject, okOffset := offsetField.Fvalue.(*object.Object); okOffset && !object.IsNull(offsetObject) {
					return offsetObject
				}
			}
			return nil

		case "localDate":
			if dateField, exists := temporalObject.FieldTable["localDate"]; exists {
				if dateObject, okDate := dateField.Fvalue.(*object.Object); okDate && !object.IsNull(dateObject) {
					return dateObject
				}
			}
			if dateField, exists := temporalObject.FieldTable["date"]; exists {
				if dateObject, okDate := dateField.Fvalue.(*object.Object); okDate && !object.IsNull(dateObject) {
					return dateObject
				}
			}
			return nil

		case "localTime":
			if timeField, exists := temporalObject.FieldTable["localTime"]; exists {
				if timeObject, okTime := timeField.Fvalue.(*object.Object); okTime && !object.IsNull(timeObject) {
					return timeObject
				}
			}
			if timeField, exists := temporalObject.FieldTable["time"]; exists {
				if timeObject, okTime := timeField.Fvalue.(*object.Object); okTime && !object.IsNull(timeObject) {
					return timeObject
				}
			}
			return nil

		default:
			if queryName != "" {
				if resultField, exists := temporalObject.FieldTable[queryName]; exists {
					if resultObject, okResult := resultField.Fvalue.(*object.Object); okResult {
						return resultObject
					}
				}
			}
		}
	}

	return nil
}

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

var classNameChronoUnitVal = "java/time/temporal/ChronoUnit"

type chronoUnitDefinition struct {
	enumName  string
	name      string
	seconds   int64
	nanos     int64
	estimated bool
	dateBased bool
	timeBased bool
}

var chronoUnitDefinitions = []chronoUnitDefinition{
	{"NANOS", "Nanos", 0, 1, false, false, true},
	{"MICROS", "Micros", 0, 1000, false, false, true},
	{"MILLIS", "Millis", 0, 1000000, false, false, true},
	{"SECONDS", "Seconds", 1, 0, false, false, true},
	{"MINUTES", "Minutes", 60, 0, false, false, true},
	{"HOURS", "Hours", 3600, 0, false, false, true},
	{"HALF_DAYS", "HalfDays", 43200, 0, false, false, true},
	{"DAYS", "Days", 86400, 0, true, true, false},
	{"WEEKS", "Weeks", 7 * 86400, 0, true, true, false},
	{"MONTHS", "Months", 31556952 / 12, 0, true, true, false},
	{"YEARS", "Years", 31556952, 0, true, true, false},
	{"DECADES", "Decades", 315569520, 0, true, true, false},
	{"CENTURIES", "Centuries", 3155695200, 0, true, true, false},
	{"MILLENNIA", "Millennia", 31556952000, 0, true, true, false},
	{"ERAS", "Eras", 31556952000000000, 0, true, true, false},
	{"FOREVER", "Forever", 31556952000000000, 0, true, false, false},
}

var chronoUnitMutex sync.Mutex
var chronoUnitInitialized bool = false
var chronoUnitInstances []*object.Object

func ensureChronoUnitInited() {
	chronoUnitMutex.Lock()
	defer chronoUnitMutex.Unlock()
	if chronoUnitInitialized {
		return
	}

	classNameDurationLocal := "java/time/Duration"
	chronoUnitInstances = make([]*object.Object, len(chronoUnitDefinitions))
	for index, definition := range chronoUnitDefinitions {
		unitObject := object.MakeEmptyObjectWithClassName(&classNameChronoUnitVal)
		unitObject.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(definition.enumName)}
		unitObject.FieldTable["displayName"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(definition.name)}
		unitObject.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(index)}
		unitObject.FieldTable["estimated"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		if definition.estimated {
			unitObject.FieldTable["estimated"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
		}
		unitObject.FieldTable["dateBased"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		if definition.dateBased {
			unitObject.FieldTable["dateBased"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
		}
		unitObject.FieldTable["timeBased"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		if definition.timeBased {
			unitObject.FieldTable["timeBased"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
		}

		durationObject := object.MakeEmptyObjectWithClassName(&classNameDurationLocal)
		durationObject.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: definition.seconds}
		durationObject.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: definition.nanos}

		unitObject.FieldTable["duration"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: durationObject,
		}
		chronoUnitInstances[index] = unitObject
		_ = statics.AddStatic(classNameChronoUnitVal+"."+definition.enumName, statics.Static{
			Type:  "Ljava/time/temporal/ChronoUnit;",
			Value: unitObject,
		})
	}
	chronoUnitInitialized = true
}

func getChronoUnitOrdinal(targetObject *object.Object) (int64, bool) {
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
			for index, definition := range chronoUnitDefinitions {
				if definition.enumName == goName || definition.name == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

func Load_Time_Temporal_ChronoUnit() {
	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitClinit,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.addTo(Ljava/time/temporal/Temporal;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronoUnitAddTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.between(Ljava/time/temporal/Temporal;Ljava/time/temporal/Temporal;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronoUnitBetween,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoUnitCompareTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoUnitCompareTo,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoUnitEquals,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.getDuration()Ljava/time/Duration;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitGetDuration,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitHashCode,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.isDateBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitIsDateBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.isDurationEstimated()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitIsDurationEstimated,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.isSupportedBy(Ljava/time/temporal/Temporal;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoUnitIsSupportedBy,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.isTimeBased()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitIsTimeBased,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitName,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitOrdinal,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitToString,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.valueOf(Ljava/lang/String;)Ljava/time/temporal/ChronoUnit;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoUnitValueOf,
		}

	ghelpers.MethodSignatures["java/time/temporal/ChronoUnit.values()[Ljava/time/temporal/ChronoUnit;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoUnitValues,
		}
}

func chronoUnitClinit([]any) any {
	ensureChronoUnitInited()
	return nil
}

func chronoUnitAddTo(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.addTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.addTo: temporal is null")
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)
	amountValue := params[2]
	return temporalPlusUnit([]any{temporalObject, amountValue, selfObject})
}

func chronoUnitBetween(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.between: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.between: temporal1Inclusive is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.between: temporal2Exclusive is null")
	}
	selfObject := params[0].(*object.Object)
	temporalInclusive := params[1].(*object.Object)
	temporalExclusive := params[2].(*object.Object)
	return temporalUntil([]any{temporalInclusive, temporalExclusive, selfObject})
}

func chronoUnitCompareTo(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.compareTo: other is null")
	}
	selfObject := params[0].(*object.Object)
	otherObject := params[1].(*object.Object)

	selfOrdinal, okSelf := getChronoUnitOrdinal(selfObject)
	otherOrdinal, okOther := getChronoUnitOrdinal(otherObject)
	if !okSelf || !okOther {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "ChronoUnit.compareTo: invalid enum type")
	}
	if selfOrdinal < otherOrdinal {
		return int64(-1)
	} else if selfOrdinal > otherOrdinal {
		return int64(1)
	}
	return int64(0)
}

func chronoUnitEquals(params []any) any {
	ensureChronoUnitInited()
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

func chronoUnitGetDuration(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.getDuration: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if durationField, exists := selfObject.FieldTable["duration"]; exists {
			if durationObject, okDuration := durationField.Fvalue.(*object.Object); okDuration && !object.IsNull(durationObject) {
				return durationObject
			}
		}
	}
	return temporalUnitGetDuration(params)
}

func chronoUnitHashCode(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return int64(0)
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoUnitOrdinal(selfObject)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func chronoUnitIsDateBased(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.isDateBased: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoUnitOrdinal(selfObject)
	if okOrdinal {
		// DAYS (7) to ERAS (14) are date-based
		if ordinalValue >= 7 && ordinalValue <= 14 {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	return temporalUnitIsDateBased(params)
}

func chronoUnitIsDurationEstimated(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.isDurationEstimated: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoUnitOrdinal(selfObject)
	if okOrdinal {
		// DAYS (7) to FOREVER (15) are estimated
		if ordinalValue >= 7 {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	return temporalUnitIsDurationEstimated(params)
}

func chronoUnitIsSupportedBy(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.isSupportedBy: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)
	return temporalAccessorIsSupported([]any{temporalObject, selfObject})
}

func chronoUnitIsTimeBased(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.isTimeBased: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoUnitOrdinal(selfObject)
	if okOrdinal {
		// NANOS (0) to HALF_DAYS (6) are time-based
		if ordinalValue >= 0 && ordinalValue <= 6 {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	return temporalUnitIsTimeBased(params)
}

func chronoUnitName(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.name: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if nameField, exists := selfObject.FieldTable["name"]; exists {
			if stringObject, okStr := nameField.Fvalue.(*object.Object); okStr && object.IsStringObject(stringObject) {
				return stringObject
			}
		}
	}
	return object.StringObjectFromGoString("")
}

func chronoUnitOrdinal(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.ordinal: self is null")
	}
	selfObject := params[0].(*object.Object)
	ordinalValue, okOrdinal := getChronoUnitOrdinal(selfObject)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func chronoUnitToString(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.toString: self is null")
	}
	selfObject := params[0].(*object.Object)
	if selfObject.FieldTable != nil {
		if displayField, exists := selfObject.FieldTable["displayName"]; exists {
			if stringObject, okStr := displayField.Fvalue.(*object.Object); okStr && object.IsStringObject(stringObject) {
				return stringObject
			}
		}
		if nameField, exists := selfObject.FieldTable["name"]; exists {
			if stringObject, okStr := nameField.Fvalue.(*object.Object); okStr && object.IsStringObject(stringObject) {
				return stringObject
			}
		}
	}
	return object.StringObjectFromGoString("ChronoUnit")
}

func chronoUnitValueOf(params []any) any {
	ensureChronoUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoUnit.valueOf: name is null")
	}
	nameStringObject, okStr := params[0].(*object.Object)
	if !okStr || !object.IsStringObject(nameStringObject) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ChronoUnit.valueOf: invalid string")
	}
	goName := object.GoStringFromStringObject(nameStringObject)
	for index, definition := range chronoUnitDefinitions {
		if definition.enumName == goName {
			return chronoUnitInstances[index]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "No enum constant java.time.temporal.ChronoUnit."+goName)
}

func chronoUnitValues([]any) any {
	ensureChronoUnitInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/temporal/ChronoUnit;", int64(len(chronoUnitInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, chronoUnitInstances)
	return arrayObject
}
